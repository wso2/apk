/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
package controlplane

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/wso2/apk/common-controller/internal/config"
	"github.com/wso2/apk/common-controller/internal/loggers"
	"github.com/wso2/apk/common-controller/internal/server"
	"github.com/wso2/apk/common-controller/internal/utils"
	"github.com/wso2/apk/common-go-libs/constants"
	apkmgt "github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/common-go-libs/pkg/discovery/api/wso2/discovery/subscription"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	grpcStatus "google.golang.org/grpc/status"
)

// Agent is a struct that implements Agent interface
type Agent struct {
	hostname         string
	port             int
	controlPlaneID   string
	artifactDeployer ArtifactDeployer
}

var (
	subsriptionList        *SubscriptionList
	applicationList        *ApplicationList
	appMappingList         *ApplicationMappingList
	connectionFaultChannel chan bool
	eventStreamingClient   apkmgt.EventStreamService_StreamEventsClient
	resources              = []resource{
		{
			endpoint:     "/subscriptions",
			responseType: subsriptionList,
		},
		{
			endpoint:     "/applications",
			responseType: applicationList,
		},
		{endpoint: "/applicationmappings",
			responseType: appMappingList,
		},
	}
)

func init() {
	connectionFaultChannel = make(chan bool, 10)
}

// NewControlPlaneAgent creates a new ControlPlaneAgent
func NewControlPlaneAgent(hostname string, port int, controlPlaneID string, artifactDeployer ArtifactDeployer) *Agent {
	return &Agent{hostname: hostname, port: port, controlPlaneID: controlPlaneID, artifactDeployer: artifactDeployer}
}

func (controlPlaneGrpcClient *Agent) initGrpcConnection() (*grpc.ClientConn, error) {
	config := config.ReadConfigs()
	publicKeyLocation, privateKeyLocation, truststoreLocation := utils.GetKeyLocations()
	cert, err := utils.GetServerCertificate(publicKeyLocation, privateKeyLocation)
	if err != nil {
		return nil, err
	}
	caCertPool := utils.GetTrustedCertPool(truststoreLocation)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   config.CommonController.Server.ServerName,
	})
	hostname := fmt.Sprintf("%s:%d", config.CommonController.ControlPlane.Host, config.CommonController.ControlPlane.EventPort)
	backOff := grpc_retry.BackoffLinearWithJitter(config.CommonController.ControlPlane.RetryInterval*time.Second, 0.5)
	conection, err := grpc.Dial(hostname, grpc.WithTransportCredentials(creds), grpc.WithStreamInterceptor(
		grpc_retry.StreamClientInterceptor(grpc_retry.WithBackoff(backOff))))
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while connecting to the control plane %s", err.Error())
		return nil, err
	}
	md := metadata.Pairs("common-controller-uuid", controlPlaneGrpcClient.controlPlaneID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	client := apkmgt.NewEventStreamServiceClient(conection)
	eventStreamingClient, err = client.StreamEvents(ctx, &apkmgt.Request{Event: controlPlaneGrpcClient.controlPlaneID})
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while initializing streaming %s", err.Error())
		return nil, err
	}
	return conection, nil
}

// StartEventStreaming starts event streaming
func (controlPlaneGrpcClient *Agent) StartEventStreaming() {
	conn := controlPlaneGrpcClient.initializeGrpcStreaming()
	for retryTrueReceived := range connectionFaultChannel {
		// event is always true
		if !retryTrueReceived {
			continue
		}
		time.Sleep(config.ReadConfigs().CommonController.ControlPlane.RetryInterval * time.Second)
		if conn != nil {
			conn.Close()
		}
		loggers.LoggerAPKOperator.Error("Connection lost. Retrying to connect to the control plane")
		conn = controlPlaneGrpcClient.initializeGrpcStreaming()
	}
}

// initializeGrpcStreaming starts event streaming
func (controlPlaneGrpcClient *Agent) initializeGrpcStreaming() *grpc.ClientConn {
	conn, err := controlPlaneGrpcClient.initGrpcConnection()
	if err != nil {
		loggers.LoggerAPKOperator.Errorf("Error while initializing the connection. error: %s", err.Error())
		connectionFaultChannel <- true
		return conn
	}
	go func() {
		for {
			resp, err := eventStreamingClient.Recv()
			if err == io.EOF {
				connectionFaultChannel <- true
				return
			}
			if err != nil {
				errStatus, _ := grpcStatus.FromError(err)
				if errStatus.Code() == codes.Unavailable ||
					errStatus.Code() == codes.DeadlineExceeded ||
					errStatus.Code() == codes.Canceled ||
					errStatus.Code() == codes.ResourceExhausted ||
					errStatus.Code() == codes.Aborted ||
					errStatus.Code() == codes.Internal {
					loggers.LoggerAPKOperator.Errorf("Connection unavailable. errorCode: %s errorMessage: %s",
						errStatus.Code().String(), errStatus.Message())
					connectionFaultChannel <- true
				}
				return
			}
			controlPlaneGrpcClient.handleEvents(resp)
		}
	}()
	return conn
}
func (controlPlaneGrpcClient *Agent) handleEvents(event *subscription.Event) {
	loggers.LoggerAPKOperator.Infof("Received event %s", event.Type)
	if event.Type == constants.AllEvents {
		go controlPlaneGrpcClient.retrieveAllData()
	} else if event.Type == constants.ApplicationCreated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_CREATED event.")
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
				TimeStamp:      event.TimeStamp,
			}
			loggers.LoggerAPKOperator.Infof("Received Application %s", application.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeployApplication(application)
		}
	} else if event.Type == constants.ApplicationUpdated {
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
			}
			loggers.LoggerAPKOperator.Infof("Received Application %s", application.UUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateApplication(application)
		}
	} else if event.Type == constants.ApplicationDeleted {
		if event.Application != nil {
			application := server.Application{UUID: event.Application.Uuid,
				Name:           event.Application.Name,
				Owner:          event.Application.Owner,
				OrganizationID: event.Application.Organization,
				Attributes:     event.Application.Attributes,
			}
			loggers.LoggerAPKOperator.Infof("Received Application %s", application.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteApplication(application.UUID)
		}
	} else if event.Type == constants.SubscriptionCreated {
		loggers.LoggerAPKOperator.Infof("Received SUBSCRIPTION_CREATED event.")
		if event.Subscription != nil {
			subscription := server.Subscription{UUID: event.Subscription.Uuid,
				Organization:  event.Subscription.Organization,
				SubStatus:     event.Subscription.SubStatus,
				SubscribedAPI: &server.SubscribedAPI{Name: event.Subscription.SubscribedApi.Name, Version: event.Subscription.SubscribedApi.Version},
				RatelimitTier: event.Subscription.RatelimitTier,
			}
			loggers.LoggerAPKOperator.Infof("Received Subscription %s", subscription.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeploySubscription(subscription)
		}
	} else if event.Type == constants.SubscriptionUpdated {
		loggers.LoggerAPKOperator.Infof("Received SUBSCRIPTION_UPDATED event.")
		if event.Subscription != nil {
			subscription := server.Subscription{UUID: event.Subscription.Uuid,
				Organization:  event.Subscription.Organization,
				SubStatus:     event.Subscription.SubStatus,
				SubscribedAPI: &server.SubscribedAPI{Name: event.Subscription.SubscribedApi.Name, Version: event.Subscription.SubscribedApi.Version},
				RatelimitTier: event.Subscription.RatelimitTier,
			}
			loggers.LoggerAPKOperator.Infof("Received Subscription %s", subscription.UUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateSubscription(subscription)
		}
	} else if event.Type == constants.SubscriptionDeleted {
		loggers.LoggerAPKOperator.Infof("Received SUBSCRIPTION_DELETED event.")
		if event.Subscription != nil {
			subscription := server.Subscription{UUID: event.Subscription.Uuid,
				Organization:  event.Subscription.Organization,
				SubStatus:     event.Subscription.SubStatus,
				SubscribedAPI: &server.SubscribedAPI{Name: event.Subscription.SubscribedApi.Name, Version: event.Subscription.SubscribedApi.Version},
			}
			loggers.LoggerAPKOperator.Infof("Received Subscription %s", subscription.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteSubscription(subscription.UUID)
		}
	} else if event.Type == constants.ApplicationKeyMappingCreated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_KEY_MAPPING_CREATED event.")
		if event.ApplicationKeyMapping != nil {
			applicationKeyMapping := server.ApplicationKeyMapping{ApplicationUUID: event.ApplicationKeyMapping.ApplicationUUID,
				SecurityScheme:        event.ApplicationKeyMapping.SecurityScheme,
				ApplicationIdentifier: event.ApplicationKeyMapping.ApplicationIdentifier,
				KeyType:               event.ApplicationKeyMapping.KeyType,
				EnvID:                 event.ApplicationKeyMapping.EnvID,
				OrganizationID:        event.ApplicationKeyMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationKeyMapping %s", applicationKeyMapping.ApplicationUUID)
			controlPlaneGrpcClient.artifactDeployer.DeployKeyMappings(applicationKeyMapping)
		}
	} else if event.Type == constants.ApplicationKeyMappingDeleted {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_KEY_MAPPING_DELETED event.")
		if event.ApplicationKeyMapping != nil {
			applicationKeyMapping := server.ApplicationKeyMapping{ApplicationUUID: event.ApplicationKeyMapping.ApplicationUUID,
				SecurityScheme:        event.ApplicationKeyMapping.SecurityScheme,
				ApplicationIdentifier: event.ApplicationKeyMapping.ApplicationIdentifier,
				KeyType:               event.ApplicationKeyMapping.KeyType,
				EnvID:                 event.ApplicationKeyMapping.EnvID,
				OrganizationID:        event.ApplicationKeyMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationKeyMapping %s", applicationKeyMapping.ApplicationUUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteKeyMappings(applicationKeyMapping)
		}
	} else if event.Type == constants.ApplicationMappingCreated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_MAPPING_CREATED event.")
		if event.ApplicationMapping != nil {
			applicationMapping := server.ApplicationMapping{UUID: event.ApplicationMapping.Uuid,
				ApplicationRef:  event.ApplicationMapping.ApplicationRef,
				SubscriptionRef: event.ApplicationMapping.SubscriptionRef,
				OrganizationID:  event.ApplicationMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationMapping %s", applicationMapping.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeployApplicationMappings(applicationMapping)
		}
	} else if event.Type == constants.ApplicationMappingDeleted {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_MAPPING_DELETED event.")
		if event.ApplicationMapping != nil {
			applicationMapping := server.ApplicationMapping{UUID: event.ApplicationMapping.Uuid,
				ApplicationRef:  event.ApplicationMapping.ApplicationRef,
				SubscriptionRef: event.ApplicationMapping.SubscriptionRef,
				OrganizationID:  event.ApplicationMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationMapping %s", applicationMapping.UUID)
			controlPlaneGrpcClient.artifactDeployer.DeleteApplicationMappings(applicationMapping.UUID)
		}
	} else if event.Type == constants.ApplicationMappingUpdated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_MAPPING_UPDATED event.")
		if event.ApplicationMapping != nil {
			applicationMapping := server.ApplicationMapping{UUID: event.ApplicationMapping.Uuid,
				ApplicationRef:  event.ApplicationMapping.ApplicationRef,
				SubscriptionRef: event.ApplicationMapping.SubscriptionRef,
				OrganizationID:  event.ApplicationMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationMapping %s", applicationMapping.UUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateApplicationMappings(applicationMapping)
		}
	} else if event.Type == constants.ApplicationKeyMappingUpdated {
		loggers.LoggerAPKOperator.Infof("Received APPLICATION_KEY_MAPPING_UPDATED event.")
		if event.ApplicationKeyMapping != nil {
			applicationKeyMapping := server.ApplicationKeyMapping{ApplicationUUID: event.ApplicationKeyMapping.ApplicationUUID,
				SecurityScheme:        event.ApplicationKeyMapping.SecurityScheme,
				ApplicationIdentifier: event.ApplicationKeyMapping.ApplicationIdentifier,
				KeyType:               event.ApplicationKeyMapping.KeyType,
				EnvID:                 event.ApplicationKeyMapping.EnvID,
				OrganizationID:        event.ApplicationKeyMapping.Organization,
			}
			loggers.LoggerAPKOperator.Infof("Received ApplicationKeyMapping %s", applicationKeyMapping.ApplicationUUID)
			controlPlaneGrpcClient.artifactDeployer.UpdateKeyMappings(applicationKeyMapping)
		}
	}
}
func (controlPlaneGrpcClient *Agent) retrieveAllData() {
	var responseChannel = make(chan response)
	config := config.ReadConfigs()
	for _, url := range resources {
		// Create a local copy of the loop variable
		localURL := url

		go InvokeService(localURL.endpoint, localURL.responseType, nil, responseChannel, 0)

		for {
			data := <-responseChannel
			loggers.LoggerAPKOperator.Info("Receiving subscription data for an environment")
			if data.Payload != nil {
				loggers.LoggerAPKOperator.Info("Payload data information received" + string(data.Payload))
				controlPlaneGrpcClient.retrieveDataFromResponseChannel(data)
				break
			} else if data.ErrorCode >= 400 && data.ErrorCode < 500 {
				//Error handle
				loggers.LoggerAPKOperator.Info("Error data information received")
				//health.SetControlPlaneRestAPIStatus(false)
			} else {
				// Keep the iteration going on until a response is received.
				// Error handle
				go func(d response, endpoint string, responseType interface{}) {
					// Retry fetching from control plane after a configured time interval
					if config.CommonController.ControlPlane.RetryInterval == 0 {
						// Assign default retry interval
						config.CommonController.ControlPlane.RetryInterval = 5
					}
					loggers.LoggerAPKOperator.Debugf("Time Duration for retrying: %v", config.CommonController.ControlPlane.RetryInterval*time.Second)
					time.Sleep(config.CommonController.ControlPlane.RetryInterval * time.Second)
					loggers.LoggerAPKOperator.Infof("Retrying to fetch APIs from control plane. Time Duration for the next retry: %v", config.CommonController.ControlPlane.RetryInterval*time.Second)
					go InvokeService(endpoint, responseType, nil, responseChannel, 0)
				}(data, localURL.endpoint, localURL.responseType)
			}
		}
	}
}

type resource struct {
	endpoint     string
	responseType interface{}
}

type response struct {
	Error     error
	Payload   []byte
	ErrorCode int
	Endpoint  string
	Type      interface{}
}

// InvokeService invokes the internal data resource
func InvokeService(endpoint string, responseType interface{}, queryParamMap map[string]string, c chan response,
	retryAttempt int) {
	config := config.ReadConfigs()
	serviceURL := "https://" + config.CommonController.ControlPlane.Host + ":" + strconv.Itoa(config.CommonController.ControlPlane.RestPort) + endpoint
	// Create the request
	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		c <- response{err, nil, 0, endpoint, responseType}
		loggers.LoggerAPKOperator.Errorf("Error occurred while creating an HTTP request for serviceURL: "+serviceURL, err)
		return
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	// Check if TLS is enabled
	skipSSL := config.CommonController.ControlPlane.SkipSSLVerification
	resp, err := InvokeControlPlane(req, skipSSL)

	if err != nil {
		if resp != nil {
			c <- response{err, nil, resp.StatusCode, endpoint, responseType}
		} else {
			c <- response{err, nil, 0, endpoint, responseType}
		}
		loggers.LoggerAPKOperator.Infof("Error occurred while calling the REST API: "+serviceURL, err)
		return
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		if err != nil {
			c <- response{err, nil, resp.StatusCode, endpoint, responseType}
			loggers.LoggerAPKOperator.Infof("Error occurred while reading the response received for: "+serviceURL, err)
			return
		}
		c <- response{nil, responseBytes, resp.StatusCode, endpoint, responseType}
	} else {
		c <- response{errors.New(string(responseBytes)), nil, resp.StatusCode, endpoint, responseType}
		loggers.LoggerAPKOperator.Infof("Failed to fetch data! "+serviceURL+" responded with "+strconv.Itoa(resp.StatusCode),
			err)
	}
}

// InvokeControlPlane sends request to the control plane and returns the response
func InvokeControlPlane(req *http.Request, skipSSL bool) (*http.Response, error) {
	tr := &http.Transport{}
	if !skipSSL {
		_, _, truststoreLocation := utils.GetKeyLocations()
		caCertPool := utils.GetTrustedCertPool(truststoreLocation)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool},
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// Configuring the http client
	client := &http.Client{
		Transport: tr,
	}
	return client.Do(req)
}
func (controlPlaneGrpcClient *Agent) retrieveDataFromResponseChannel(response response) {
	responseType := reflect.TypeOf(response.Type).Elem()
	newResponse := reflect.New(responseType).Interface()
	err := json.Unmarshal(response.Payload, &newResponse)

	if err != nil {
		loggers.LoggerAPI.Infof("Error occurred while unmarshalling the response received for: "+response.Endpoint, err)
	} else {
		switch t := newResponse.(type) {
		case *SubscriptionList:
			loggers.LoggerAPI.Infof("Received Subscription information.")
			subList := newResponse.(*SubscriptionList)
			resolvedSubscriptionList := marshalMultipleSubscriptions(subList)
			if len(resolvedSubscriptionList.List) > 0 {
				controlPlaneGrpcClient.artifactDeployer.DeployAllSubscriptions(resolvedSubscriptionList)
			}

		case *ApplicationList:
			loggers.LoggerAPI.Infof("Received Application information.")
			appList := newResponse.(*ApplicationList)
			resolvedApplicationList := marshalMultipleApplications(appList)
			resolvedApplicationKeyMappingList := marshalMultipleApplicationKeyMappings(appList)
			if len(resolvedApplicationList.List) > 0 {
				controlPlaneGrpcClient.artifactDeployer.DeployAllApplications(resolvedApplicationList)
			}
			if len(resolvedApplicationKeyMappingList.List) > 0 {
				controlPlaneGrpcClient.artifactDeployer.DeployAllKeyMappings(resolvedApplicationKeyMappingList)
			}
		case *ApplicationMappingList:
			loggers.LoggerAPI.Infof("Received Application Mapping information.")
			appMappingList := newResponse.(*ApplicationMappingList)
			resolvedApplicationMappingList := marshalMultipleApplicationMappings(appMappingList)
			if len(resolvedApplicationMappingList.List) > 0 {
				controlPlaneGrpcClient.artifactDeployer.DeployAllApplicationMappings(resolvedApplicationMappingList)
			}
		default:
			loggers.LoggerAPI.Debugf("Unknown type %T", t)
		}
	}
}
func marshalMultipleSubscriptions(subList *SubscriptionList) server.SubscriptionList {
	subscriptionList := server.SubscriptionList{List: []server.Subscription{}}
	for _, subscription := range subList.List {
		loggers.LoggerAPI.Debugf("Subscription: %v", subscription)
		subscriptionList.List = append(subscriptionList.List, server.Subscription{UUID: subscription.UUID, Organization: subscription.Organization, SubStatus: subscription.SubStatus, SubscribedAPI: &server.SubscribedAPI{Name: subscription.SubscribedAPI.Name, Version: subscription.SubscribedAPI.Version}})
	}
	return subscriptionList
}
func marshalMultipleApplications(appList *ApplicationList) server.ApplicationList {
	applicationList := server.ApplicationList{List: []server.Application{}}
	for _, application := range appList.List {
		loggers.LoggerAPI.Debugf("Application: %v", application)
		applicationList.List = append(applicationList.List, server.Application{UUID: application.UUID, Name: application.Name, Owner: application.Owner, OrganizationID: application.Organization, Attributes: application.Attributes})
	}
	return applicationList
}
func marshalMultipleApplicationKeyMappings(appList *ApplicationList) server.ApplicationKeyMappingList {
	applicationKeyMappingList := server.ApplicationKeyMappingList{List: []server.ApplicationKeyMapping{}}
	for _, application := range appList.List {
		loggers.LoggerAPI.Debugf("Application: %v", application)
		for _, securityScheme := range application.SecuritySchemes {
			applicationKeyMappingList.List = append(applicationKeyMappingList.List, server.ApplicationKeyMapping{ApplicationUUID: application.UUID, SecurityScheme: securityScheme.SecurityScheme, ApplicationIdentifier: securityScheme.ApplicationIdentifier, KeyType: securityScheme.KeyType, EnvID: securityScheme.EnvID, OrganizationID: application.Organization})
		}
	}
	return applicationKeyMappingList
}
func marshalMultipleApplicationMappings(appMappingList *ApplicationMappingList) server.ApplicationMappingList {
	applicationMappingList := server.ApplicationMappingList{List: []server.ApplicationMapping{}}
	for _, applicationMapping := range appMappingList.List {
		loggers.LoggerAPI.Debugf("ApplicationMapping: %v", applicationMapping)
		applicationMappingList.List = append(applicationMappingList.List, server.ApplicationMapping{UUID: applicationMapping.UUID, ApplicationRef: applicationMapping.ApplicationRef, SubscriptionRef: applicationMapping.SubscriptionRef, OrganizationID: applicationMapping.Organization})
	}
	return applicationMappingList
}
