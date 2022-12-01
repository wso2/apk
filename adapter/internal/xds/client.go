/*
 *  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package xds

import (
	"context"
	"fmt"
	"io"
	"reflect"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	cpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/cp/v1alpha1"

	apkmgt_model "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/apkmgt"
	stub "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"

	operatorutils "github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// Last Acknowledged Response from the apkmgt server
	lastAckedResponse *discovery.DiscoveryResponse
	// Last Received Response from the apkmgt server
	// Last Received Response is always is equal to the lastAckedResponse according to current implementation as there is no
	// validation performed on successfully received response.
	lastReceivedResponse *discovery.DiscoveryResponse
	// XDS stream for streaming Aplications from APK Mgt client
	xdsStream stub.APKMgtDiscoveryService_StreamAPKMgtApplicationsClient
	// applicationMap contains the application cache
	applicationMap map[string]cpv1alpha1.Application
	// applicationChannel is used to notifiy the application updates
	applicationChannel chan ApplicationEvent
)

// EventType is the type of the event
type EventType int

const (
	// ApplicationCreate is application create event type
	ApplicationCreate = 0
	// ApplicationUpdate is application update event type
	ApplicationUpdate = 1
	// ApplicationDelete is application delete event type
	ApplicationDelete = 2
)

// ApplicationEvent is the application event data holder
type ApplicationEvent struct {
	Type        EventType
	Application *cpv1alpha1.Application
}

const (
	// The type url for requesting Application Entries from apkmgt server.
	applicationTypeURL string = "type.googleapis.com/wso2.discovery.apkmgt.Application"
)

func init() {
	lastAckedResponse = &discovery.DiscoveryResponse{}
	applicationChannel = make(chan ApplicationEvent, 1000)
	applicationMap = make(map[string]cpv1alpha1.Application)
}

func initConnection(xdsURL string) error {
	// TODO: (AmaliMatharaarachchi) Bring in connection level configurations
	conn, err := grpc.Dial(xdsURL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		// TODO: (AmaliMatharaarachchi) retries
		loggers.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprint("Error while connecting to the APK Management Server. ", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1700,
		})
		return err
	}

	client := stub.NewAPKMgtDiscoveryServiceClient(conn)
	streamContext := context.Background()
	xdsStream, err = client.StreamAPKMgtApplications(streamContext)

	if err != nil {
		// TODO: (AmaliMatharaarachchi) handle error.
		loggers.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprint("Error while starting APK Management application stream. ", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1701,
		})
		return err
	}
	loggers.LoggerXds.Infof("Connection to the APK Management Server: %s is successful.", xdsURL)
	return nil
}

func watchApplications() {
	for {
		discoveryResponse, err := xdsStream.Recv()
		if err == io.EOF {
			loggers.LoggerXds.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprint("EOF is received from the APK Management Server application stream. ", err.Error()),
				Severity:  logging.CRITICAL,
				ErrorCode: 1702,
			})
			return
		}
		if err != nil {
			loggers.LoggerXds.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprint("Failed to receive the discovery response from the APK Management Server application stream. ", err.Error()),
				Severity:  logging.CRITICAL,
				ErrorCode: 1703,
			})
			errStatus, _ := grpcStatus.FromError(err)
			if errStatus.Code() == codes.Unavailable {
				loggers.LoggerXds.ErrorC(logging.ErrorDetails{
					Message:   fmt.Sprint("The APK Management Server application stream connection stopped"),
					Severity:  logging.MINOR,
					ErrorCode: 1704,
				})
			}
			nack(err.Error())
		} else {
			lastReceivedResponse = discoveryResponse
			loggers.LoggerXds.Debugf("Discovery response is received : %s", discoveryResponse.VersionInfo)
			addApplicationsToChannel(discoveryResponse)
			ack()
		}
	}
}

func ack() {
	lastAckedResponse = lastReceivedResponse
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:          getAdapterNode(),
		VersionInfo:   lastAckedResponse.VersionInfo,
		TypeUrl:       applicationTypeURL,
		ResponseNonce: lastReceivedResponse.Nonce,
	}
	xdsStream.Send(discoveryRequest)
}

func nack(errorMessage string) {
	if lastAckedResponse == nil {
		return
	}
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:          getAdapterNode(),
		VersionInfo:   lastAckedResponse.VersionInfo,
		TypeUrl:       applicationTypeURL,
		ResponseNonce: lastReceivedResponse.Nonce,
		ErrorDetail: &status.Status{
			Message: errorMessage,
		},
	}
	xdsStream.Send(discoveryRequest)
}

func getAdapterNode() *core.Node {
	config := config.ReadConfigs()
	return &core.Node{
		Id: config.ManagementServer.NodeLabel,
	}
}

// InitApkMgtClient initializes the connection to the apkmgt server.
func InitApkMgtClient() {
	loggers.LoggerXds.Info("Starting the XDS Client connection to APK Management server.")
	config := config.ReadConfigs()
	err := initConnection(config.ManagementServer.ServiceURL)
	if err == nil {
		go watchApplications()
		discoveryRequest := &discovery.DiscoveryRequest{
			Node:        getAdapterNode(),
			VersionInfo: "",
			TypeUrl:     applicationTypeURL,
		}
		xdsStream.Send(discoveryRequest)
	} else {
		loggers.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprint("Error while starting the APK Management Server. ", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1705,
		})
	}
}

func addApplicationsToChannel(resp *discovery.DiscoveryResponse) {
	var newApplicationUUIDs []string

	for _, res := range resp.Resources {
		application := &apkmgt_model.Application{}
		err := ptypes.UnmarshalAny(res, application)

		if err != nil {
			loggers.LoggerXds.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprint("Error while unmarshalling APK Management Server Application discovery response. ", err.Error()),
				Severity:  logging.MINOR,
				ErrorCode: 1706,
			})
			continue
		}

		applicationUUID := application.Uuid
		newApplicationUUIDs = append(newApplicationUUIDs, applicationUUID)

		applicationResource := &cpv1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: operatorutils.GetOperatorPodNamespace(),
				Name:      application.Uuid,
			},
			Spec: cpv1alpha1.ApplicationSpec{
				UUID:       application.Uuid,
				Name:       application.Name,
				Owner:      application.Owner,
				Attributes: application.Attributes,
			},
		}

		var consumerKeys []cpv1alpha1.ConsumerKey
		for _, consumerKey := range application.ConsumerKeys {
			consumerKeys = append(consumerKeys, cpv1alpha1.ConsumerKey{Key: consumerKey.Key, KeyManager: consumerKey.KeyManager})
		}
		applicationResource.Spec.ConsumerKeys = consumerKeys

		var subscriptions []cpv1alpha1.Subscription
		for _, subscription := range application.Subscriptions {
			subscriptions = append(subscriptions, cpv1alpha1.Subscription{
				UUID:               subscription.Uuid,
				SubscriptionStatus: subscription.SubscriptionStatus,
				PolicyID:           subscription.PolicyId,
				APIRef:             subscription.ApiUuid,
			})
		}
		applicationResource.Spec.Subscriptions = subscriptions

		var event ApplicationEvent

		if currentApplication, found := applicationMap[applicationUUID]; found {
			if reflect.DeepEqual(currentApplication.Spec, applicationResource.Spec) {
				continue
			}
			// Application update event
			event = ApplicationEvent{
				Type:        ApplicationUpdate,
				Application: applicationResource,
			}
			applicationMap[applicationUUID] = *applicationResource
		} else {
			// Application create event
			event = ApplicationEvent{
				Type:        ApplicationCreate,
				Application: applicationResource,
			}
			applicationMap[applicationUUID] = *applicationResource
		}

		applicationChannel <- event

	}
	// Send delete events for removed applications
	for _, application := range applicationMap {
		if !stringutils.StringInSlice(application.Name, newApplicationUUIDs) {
			// Application delete event
			event := ApplicationEvent{
				Type:        ApplicationDelete,
				Application: &application,
			}
			applicationChannel <- event
			delete(applicationMap, application.Name)
		}
	}
}
