/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package synchronizer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/types"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
)

// APIEvent holds the data structure used for passing API
// lifecycle events from controller go routine to synchronizer
// go routine.
type APIEvent struct {
	EventType     string
	Events        []APIState
	UpdatedEvents []string
}

// SuccessEvent holds the data structure used for aknowledgement of a successful API deployment
type SuccessEvent struct {
	// APINamespacedName updated api namespaced names
	APINamespacedName []types.NamespacedName
	State             string
	Events            []string
}

// PartitionEvent is the event sent to the partition server.
type PartitionEvent struct {
	EventType    string   `json:"eventType"`
	APIName      string   `json:"apiName"`
	APIVersion   string   `json:"apiVersion"`
	BasePath     string   `json:"basePath"`
	Organization string   `json:"organization"`
	Partition    string   `json:"partition"`
	APIUUID      string   `json:"apiId"`
	Vhosts       []string `json:"vhosts"`
}

var (
	// TODO: Decide on a buffer size and add to config.
	paritionCh chan *APIEvent
	// Runtime client connetion
	partitionClient *http.Client
)

func init() {
	if config.ReadConfigs().PartitionServer.Enabled {
		paritionCh = make(chan *APIEvent, 10)
	}
}

// HandleAPILifeCycleEvents handles the API events generated from OperatorDataStore
func HandleAPILifeCycleEvents(ch *chan *APIEvent, successChannel *chan SuccessEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for API lifecycle events...")
	for event := range *ch {
		var err error
		switch event.EventType {
		case constants.Delete:
			loggers.LoggerAPKOperator.Infof("Delete event received for %v", event.Events[0].APIDefinition.Name)
			if err = undeployAPIInGateway(event); err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2629, logging.CRITICAL, "API deployment failed for %s event : %v", event.EventType, err))
			} else {
				if config.ReadConfigs().PartitionServer.Enabled {
					paritionCh <- event
				}
			}
		case constants.Create:
			deployMultipleAPIsInGateway(event, successChannel)
		case constants.Update:
			deployMultipleAPIsInGateway(event, successChannel)
		}
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2629, logging.CRITICAL, "API deployment failed for %s event : %v", event.EventType, err))
		} else if event.EventType != constants.Create {
			if config.ReadConfigs().PartitionServer.Enabled {
				paritionCh <- event
			}
		}
	}
}

func undeployAPIInGateway(apiEvent *APIEvent) error {
	var err error
	apiState := apiEvent.Events[0]
	if apiState.APIDefinition.Spec.APIType == constants.REST {
		err = undeployRestAPIInGateway(apiState)
	}
	if apiState.APIDefinition.Spec.APIType == constants.GRAPHQL {
		err = undeployGQLAPIInGateway(apiState)
	}

	if apiState.APIDefinition.Spec.APIType == constants.GRPC {
		return undeployGRPCAPIInGateway(apiState)
	}
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2629, logging.CRITICAL,
			"API deployment failed for %s event : %v, %v", apiEvent.EventType, apiState.APIDefinition.Name, err))
	} else if config.ReadConfigs().PartitionServer.Enabled {
		paritionCh <- apiEvent
	}
	return nil
}

// deployMultipleAPIsInGateway deploys the related API in CREATE and UPDATE events.
func deployMultipleAPIsInGateway(event *APIEvent, successChannel *chan SuccessEvent) {
	updatedLabelsMap := make(map[string]struct{})
	var updatedAPIs []types.NamespacedName
	for i, apiState := range event.Events {
		loggers.LoggerAPKOperator.Infof("%s event received for %s", event.EventType, apiState.APIDefinition.Name)
		// Remove the API from the internal maps before adding it again
		oldGatewayNames := xds.RemoveAPIFromAllInternalMaps(string((*apiState.APIDefinition).ObjectMeta.UID))
		for label := range oldGatewayNames {
			updatedLabelsMap[label] = struct{}{}
		}
		if apiState.APIDefinition.Spec.APIType == "REST" {
			if apiState.ProdHTTPRoute != nil {
				_, updatedLabels, err := UpdateInternalMapsFromHTTPRoute(apiState, apiState.ProdHTTPRoute, constants.Production)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2665, logging.CRITICAL,
						"Error deploying prod httpRoute of API : %v in Organization %v from environments %v. Error: %v",
						string(apiState.APIDefinition.Spec.APIName), apiState.APIDefinition.Spec.Organization,
						getGatewayNameForAPI(apiState.ProdHTTPRoute.HTTPRouteCombined), err))
					// removing failed updates from the events list because this will be sent to partition server
					if len(event.Events) > i {
						event.Events = []APIState{}
					} else {
						event.Events = append(event.Events[:i], event.Events[i+1:]...)
					}
					continue
				}
				for label := range updatedLabels {
					updatedLabelsMap[label] = struct{}{}
				}
			}

			if apiState.SandHTTPRoute != nil {
				_, updatedLabels, err := UpdateInternalMapsFromHTTPRoute(apiState, apiState.SandHTTPRoute, constants.Sandbox)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2666, logging.CRITICAL,
						"Error deploying sand httpRoute of API : %v in Organization %v from environments %v. Error: %v",
						string(apiState.APIDefinition.Spec.APIName), apiState.APIDefinition.Spec.Organization,
						getGatewayNameForAPI(apiState.ProdHTTPRoute.HTTPRouteCombined), err))
					// removing failed updates from the events list because this will be sent to partition server
					if len(event.Events) > i {
						event.Events = []APIState{}
					} else {
						event.Events = append(event.Events[:i], event.Events[i+1:]...)
					}
					continue
				}
				for label := range updatedLabels {
					updatedLabelsMap[label] = struct{}{}
				}
			}
		}

		if apiState.APIDefinition.Spec.APIType == "GraphQL" {
			if apiState.ProdGQLRoute != nil {
				_, updatedLabels, err := updateInternalMapsFromGQLRoute(apiState, apiState.ProdGQLRoute, constants.Production)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2665, logging.CRITICAL,
						"Error deploying prod gqlRoute of API : %v in Organization %v from environments %v. Error: %v",
						string(apiState.APIDefinition.Spec.APIName), apiState.APIDefinition.Spec.Organization,
						getLabelsForGQLAPI(apiState.ProdGQLRoute.GQLRouteCombined), err))
					continue
				}
				for label := range updatedLabels {
					updatedLabelsMap[label] = struct{}{}
				}
			}
			if apiState.SandGQLRoute != nil {
				_, updatedLabels, err := updateInternalMapsFromGQLRoute(apiState, apiState.SandGQLRoute, constants.Sandbox)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2665, logging.CRITICAL,
						"Error deploying sand gqlRoute of API : %v in Organization %v from environments %v. Error: %v",
						string(apiState.APIDefinition.Spec.APIName), apiState.APIDefinition.Spec.Organization,
						getLabelsForGQLAPI(apiState.SandGQLRoute.GQLRouteCombined), err))
					continue
				}
				for label := range updatedLabels {
					updatedLabelsMap[label] = struct{}{}
				}
			}
		}

		if apiState.APIDefinition.Spec.APIType == constants.GRPC {
			if apiState.ProdGRPCRoute != nil {
				_, updatedLabels, err := updateInternalMapsFromGRPCRoute(apiState, apiState.ProdGRPCRoute, constants.Production)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2665, logging.CRITICAL,
						"Error deploying prod grpcRoute of API : %v in Organization %v from environments %v. Error: %v",
						string(apiState.APIDefinition.Spec.APIName), apiState.APIDefinition.Spec.Organization,
						getLabelsForGRPCAPI(apiState.ProdGRPCRoute.GRPCRouteCombined), err))
					continue
				}
				for label := range updatedLabels {
					updatedLabelsMap[label] = struct{}{}
				}
			}

			if apiState.SandGRPCRoute != nil {
				_, updatedLabels, err := updateInternalMapsFromGRPCRoute(apiState, apiState.SandGRPCRoute, constants.Sandbox)
				if err != nil {
					loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2665, logging.CRITICAL,
						"Error deploying sand grpcRoute of API : %v in Organization %v from environments %v. Error: %v",
						string(apiState.APIDefinition.Spec.APIName), apiState.APIDefinition.Spec.Organization,
						getLabelsForGRPCAPI(apiState.SandGRPCRoute.GRPCRouteCombined), err))
					continue
				}
				for label := range updatedLabels {
					updatedLabelsMap[label] = struct{}{}
				}
			}
		}
		updatedAPIs = append(updatedAPIs, utils.NamespacedName(apiState.APIDefinition))
	}

	updated := xds.UpdateXdsCacheOnAPIChange(updatedLabelsMap)
	if updated {
		loggers.LoggerAPKOperator.Infof("XDS cache updated for apis: %+v", updatedAPIs)
		*successChannel <- SuccessEvent{
			APINamespacedName: updatedAPIs,
			State:             event.EventType,
			Events:            event.UpdatedEvents,
		}
		if config.ReadConfigs().PartitionServer.Enabled {
			paritionCh <- event
		}
	} else {
		loggers.LoggerAPKOperator.Infof("XDS cache not updated for APIs : %+v", updatedAPIs)
	}
}

func init() {
	conf := config.ReadConfigs()

	_, _, truststoreLocation := tlsutils.GetKeyLocations()
	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)
	transport := &http.Transport{
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
		TLSClientConfig: &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: conf.PartitionServer.DisableSslVerification},
	}
	partitionClient = &http.Client{Transport: transport}
}

// SendEventToPartitionServer sends the API create/update/delete event to the partition server.
func SendEventToPartitionServer() {
	conf := config.ReadConfigs()
	for apiEvent := range paritionCh {
		for _, event := range apiEvent.Events {
			if !event.APIDefinition.Spec.SystemAPI {
				apiDefinition := event.APIDefinition
				loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v, API_UUID: %v", apiDefinition.Spec.APIName, string(apiDefinition.ObjectMeta.UID))
				api := event
				eventType := apiEvent.EventType
				basePath := api.APIDefinition.Spec.BasePath
				organization := api.APIDefinition.Spec.Organization
				version := api.APIDefinition.Spec.APIVersion
				apiName := api.APIDefinition.Spec.APIName
				apiUUID := string(api.APIDefinition.Name)
				var hostNames []string
				httpRoute := api.ProdHTTPRoute
				if httpRoute == nil {
					httpRoute = api.SandHTTPRoute
				}
				for _, hostName := range httpRoute.HTTPRouteCombined.Spec.Hostnames {
					hostNames = append(hostNames, string(hostName))
				}
				grpcRoute := api.ProdGRPCRoute
				if grpcRoute == nil {
					grpcRoute = api.SandGRPCRoute
				}
				for _, hostName := range grpcRoute.GRPCRouteCombined.Spec.Hostnames {
					hostNames = append(hostNames, string(hostName))
				}
				data := PartitionEvent{
					EventType:    eventType,
					BasePath:     basePath,
					Organization: organization,
					APIVersion:   version,
					APIName:      apiName,
					APIUUID:      apiUUID,
					Vhosts:       hostNames,
					Partition:    conf.PartitionServer.PartitionName,
				}
				payload, err := json.Marshal(data)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error creating Event: %v, API_UUID: %v", err, apiUUID)
				}
				req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", conf.PartitionServer.Host, conf.PartitionServer.ServiceBasePath, "/api-deployment"), bytes.NewBuffer(payload))
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error creating api definition request: %v, API_UUID: %v", err, apiUUID)
				}
				req.Header.Set("Content-Type", "application/json; charset=UTF-8")
				resp, err := partitionClient.Do(req)
				if err != nil {
					loggers.LoggerAPKOperator.Errorf("Error sending API Event: %v, API_UUID: %v", err, apiUUID)
				}
				if resp.StatusCode == http.StatusAccepted {
					loggers.LoggerAPKOperator.Info("API Event Accepted", resp.Status)
				}
			}
		}

	}
}
