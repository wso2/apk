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
	logging "github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/internal/management-server/utils"
	cpv1alpha1 "github.com/wso2/apk/adapter/pkg/apis/cp/v1alpha1"

	stub "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/subscription"
	sub_model "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"

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
	xdsStream stub.ApplicationDiscoveryService_StreamApplicationsClient
	// applicationMap contains the application cache
	applicationMap map[string]cpv1alpha1.Application
	// applicationChannel is used to notifiy the application updates
	applicationChannel chan ApplicationEvent
	// subscriptionMap contains the application cache
	subscriptionMap map[string]cpv1alpha1.Subscription
	// subscriptionChannel is used to notifiy the subscription updates
	subscriptionChannel chan SubscriptionEvent
	// XDS stream for streaming Subscriptions from client
	xdsSubStream stub.SubscriptionDiscoveryService_StreamSubscriptionsClient
	// Last Acknowledged Response from the apkmgt server
	lastAckedResponseSub *discovery.DiscoveryResponse
	// Last Received Response from the apkmgt server
	// Last Received Response is always is equal to the lastAckedResponse according to current implementation as there is no
	// validation performed on successfully received response.
	lastReceivedResponseSub *discovery.DiscoveryResponse
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

const (
	// SubscriptionCreate is subscription create event type
	SubscriptionCreate = 0
	// SubscriptionUpdate is subscription update event type
	SubscriptionUpdate = 1
	// SubscriptionDelete is subscription delete event type
	SubscriptionDelete = 2
)

// ApplicationEvent is the application event data holder
type ApplicationEvent struct {
	Type        EventType
	Application *cpv1alpha1.Application
}

// SubscriptionEvent is the subsctiption event data holder
type SubscriptionEvent struct {
	Type         EventType
	Subscription *cpv1alpha1.Subscription
}

const (
	// The type url for requesting Application Entries from apkmgt server.
	applicationTypeURL string = "type.googleapis.com/wso2.discovery.subscription.Application"
	// The type url for requesting Subscription Entries from apkmgt server.
	subscriptionTypeURL string = "type.googleapis.com/wso2.discovery.subscription.Subscription"
)

func init() {
	lastAckedResponse = &discovery.DiscoveryResponse{}
	lastAckedResponseSub = &discovery.DiscoveryResponse{}
	applicationChannel = make(chan ApplicationEvent, 1000)
	applicationMap = make(map[string]cpv1alpha1.Application)
	subscriptionChannel = make(chan SubscriptionEvent, 1000)
	subscriptionMap = make(map[string]cpv1alpha1.Subscription)
}

func initConnection(xdsURL string) error {
	// TODO: (AmaliMatharaarachchi) Bring in connection level configurations
	transportCredentials, err := utils.GenerateTLSCredentials()
	conn, err := grpc.Dial(xdsURL, grpc.WithTransportCredentials(transportCredentials), grpc.WithBlock())
	if err != nil {
		// TODO: (AmaliMatharaarachchi) retries
		loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1700, logging.BLOCKER, "Error while connecting to the APK Management Server. %v", err.Error()))
		return err
	}

	client := stub.NewApplicationDiscoveryServiceClient(conn)
	clientSub := stub.NewSubscriptionDiscoveryServiceClient(conn)
	streamContext := context.Background()
	xdsStream, err = client.StreamApplications(streamContext)
	xdsSubStream, err = clientSub.StreamSubscriptions(streamContext)

	if err != nil {
		// TODO: (AmaliMatharaarachchi) handle error.
		loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1701, logging.BLOCKER, "Error while starting APK Management application stream. %v", err.Error()))
		return err
	}
	loggers.LoggerXds.Infof("Connection to the APK Management Server: %s is successful.", xdsURL)
	return nil
}

func watchApplications() {
	for {
		discoveryResponse, err := xdsStream.Recv()
		if err == io.EOF {
			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1702, logging.CRITICAL, "EOF is received from the APK Management Server application stream. %v", err.Error()))
			return
		}
		if err != nil {
			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1703, logging.CRITICAL, "Failed to receive the discovery response from the APK Management Server application stream. %v", err.Error()))
			errStatus, _ := grpcStatus.FromError(err)
			if errStatus.Code() == codes.Unavailable {
				loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1704, logging.MINOR, "The APK Management Server application stream connection stopped: %v", err.Error()))
				return
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

func watchSubscriptions() {
	for {
		discoveryResponse, err := xdsSubStream.Recv()
		if err == io.EOF {
			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1717, logging.CRITICAL, "EOF is received from the APK Management Server subscription stream. %v", err.Error()))
			return
		}
		if err != nil {
			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1718, logging.CRITICAL, "Failed to receive the discovery response from the APK Management Server subscription stream. %v", err.Error()))
			errStatus, _ := grpcStatus.FromError(err)
			if errStatus.Code() == codes.Unavailable {
				loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1719, logging.MINOR, "The APK Management Server subscription stream connection stopped: %v", err.Error()))
				return
			}
			nackSub(err.Error())
		} else {
			lastReceivedResponseSub = discoveryResponse
			loggers.LoggerXds.Debugf("Discovery response is received : %s", discoveryResponse.VersionInfo)
			addSubscriptionsToChannel(discoveryResponse)
			ackSub()
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

func ackSub() {
	lastAckedResponseSub = lastReceivedResponseSub
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:          getAdapterNode(),
		VersionInfo:   lastAckedResponseSub.VersionInfo,
		TypeUrl:       subscriptionTypeURL,
		ResponseNonce: lastReceivedResponseSub.Nonce,
	}
	xdsSubStream.Send(discoveryRequest)
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

func nackSub(errorMessage string) {
	if lastAckedResponseSub == nil {
		return
	}
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:          getAdapterNode(),
		VersionInfo:   lastAckedResponseSub.VersionInfo,
		TypeUrl:       subscriptionTypeURL,
		ResponseNonce: lastReceivedResponseSub.Nonce,
		ErrorDetail: &status.Status{
			Message: errorMessage,
		},
	}
	xdsSubStream.Send(discoveryRequest)
}

func getAdapterNode() *core.Node {
	config := config.ReadConfigs()
	return &core.Node{
		Id: config.ManagementServer.NodeLabel,
	}
}

// InitApkMgtXDSClient initializes the connection to the apkmgt server.
func InitApkMgtXDSClient() {
	loggers.LoggerXds.Info("Starting the XDS Client connection to APK Management server.")
	config := config.ReadConfigs()
	err := initConnection(fmt.Sprintf("%s:%d", config.ManagementServer.Host, config.ManagementServer.XDSPort))
	if err == nil {
		go watchApplications()
		discoveryRequest := &discovery.DiscoveryRequest{
			Node:        getAdapterNode(),
			VersionInfo: "",
			TypeUrl:     applicationTypeURL,
		}
		xdsStream.Send(discoveryRequest)
		go watchSubscriptions()
		discoveryRequestSub := &discovery.DiscoveryRequest{
			Node:        getAdapterNode(),
			VersionInfo: "",
			TypeUrl:     subscriptionTypeURL,
		}
		xdsSubStream.Send(discoveryRequestSub)
	} else {
		loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1705, logging.BLOCKER, "Error while starting the APK Management Server: %v", err.Error()))
	}
}

func addApplicationsToChannel(resp *discovery.DiscoveryResponse) {
	var newApplicationUUIDs []string

	for _, res := range resp.Resources {
		application := &sub_model.Application{}
		err := ptypes.UnmarshalAny(res, application)

		if err != nil {
			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1706, logging.MINOR, "Error while unmarshalling APK Management Server Application discovery response: %v", err.Error()))
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
				Name:         application.Name,
				Owner:        application.Owner,
				Attributes:   application.Attributes,
				Policy:       application.Policy,
				Organization: application.Organization,
			},
		}

		var consumerKeys []cpv1alpha1.Key
		for _, consumerKey := range application.Keys {
			consumerKeys = append(consumerKeys, cpv1alpha1.Key{Key: consumerKey.Key, KeyManager: consumerKey.KeyManager})
		}
		applicationResource.Spec.Keys = consumerKeys

		// Todo:(Sampath) Need to handle adding the subscriptions coming from management server seperately
		// var subscriptions []cpv1alpha1.Subscription
		// for _, subscription := range application.Subscriptions {
		// 	subscriptions = append(subscriptions, cpv1alpha1.Subscription{
		// 		UUID:               subscription.Name,
		// 		SubscriptionStatus: subscription.Spec.SubscriptionStatus,
		// 		PolicyID:           subscription.PolicyId,
		// 		APIRef:             subscription.ApiUuid,
		// 	})
		// }
		// applicationResource.Spec.Subscriptions = subscriptions

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
	for item := range applicationMap {
		application := applicationMap[item]
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

func addSubscriptionsToChannel(resp *discovery.DiscoveryResponse) {
	var newSubscriptionUUIDs []string

	for _, res := range resp.Resources {
		subscription := &sub_model.Subscription{}
		err := ptypes.UnmarshalAny(res, subscription)

		if err != nil {
			loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error1720, logging.MINOR, "Error while unmarshalling APK Management Server Subscription discovery response: %v", err.Error()))
			continue
		}

		subscriptionUUID := subscription.Uuid
		newSubscriptionUUIDs = append(newSubscriptionUUIDs, subscriptionUUID)

		subscriptionResource := &cpv1alpha1.Subscription{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: operatorutils.GetOperatorPodNamespace(),
				Name:      subscription.Uuid,
			},
			Spec: cpv1alpha1.SubscriptionSpec{
				APIRef:             subscription.ApiRef,
				ApplicationRef:     subscription.ApplicationRef,
				PolicyID:           subscription.PolicyId,
				SubscriptionStatus: subscription.SubStatus,
				Subscriber:         subscription.Subscriber,
				Organization:       subscription.Organization,
			},
		}

		var event SubscriptionEvent

		if currentSubscription, found := subscriptionMap[subscriptionUUID]; found {
			if reflect.DeepEqual(currentSubscription.Spec, subscriptionResource.Spec) {
				continue
			}
			// Subscription update event
			event = SubscriptionEvent{
				Type:         SubscriptionUpdate,
				Subscription: subscriptionResource,
			}
			subscriptionMap[subscriptionUUID] = *subscriptionResource
		} else {
			// Subscription create event
			event = SubscriptionEvent{
				Type:         SubscriptionCreate,
				Subscription: subscriptionResource,
			}
			subscriptionMap[subscriptionUUID] = *subscriptionResource
		}

		subscriptionChannel <- event

	}
	// Send delete events for removed subscriptions
	for item := range subscriptionMap {
		subscription := subscriptionMap[item]
		if !stringutils.StringInSlice(subscription.Name, newSubscriptionUUIDs) {
			// Subscription delete event
			event := SubscriptionEvent{
				Type:         SubscriptionDelete,
				Subscription: &subscription,
			}
			subscriptionChannel <- event
			delete(subscriptionMap, subscription.Name)
		}
	}
}
