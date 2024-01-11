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
	Event         APIState
	UpdatedEvents []string
}

// SuccessEvent holds the data structure used for aknowledgement of a successful API deployment
type SuccessEvent struct {
	APINamespacedName types.NamespacedName
	State             string
	Events            []string
}

var (
	// TODO: Decide on a buffer size and add to config.
	paritionCh chan APIEvent
)

func init() {
	paritionCh = make(chan APIEvent, 10)
}

// HandleAPILifeCycleEvents handles the API events generated from OperatorDataStore
func HandleAPILifeCycleEvents(ch *chan APIEvent, successChannel *chan SuccessEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for API lifecycle events...")
	for event := range *ch {
		if event.Event.APIDefinition == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2628, logging.CRITICAL, "API Event is nil"))
		}
		loggers.LoggerAPKOperator.Infof("%s event received for %v", event.EventType, event.Event.APIDefinition.Name)
		var err error
		switch event.EventType {
		case constants.Delete:
			err = undeployAPIInGateway(event.Event)
		case constants.Create:
			err = deployAPIInGateway(event.Event)
		case constants.Update:
			err = deployAPIInGateway(event.Event)
		}
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2629, logging.MAJOR, "API deployment failed for %s event : %v", event.EventType, err))
		} else {
			if event.EventType != constants.Delete {
				*successChannel <- SuccessEvent{
					APINamespacedName: utils.NamespacedName(event.Event.APIDefinition),
					State:             event.EventType,
					Events:            event.UpdatedEvents,
				}
			}
			if config.ReadConfigs().PartitionServer.Enabled {
				paritionCh <- event
			}
		}
	}
}

// deployAPIInGateway deploys the related API in CREATE and UPDATE events.
func deployAPIInGateway(apiState APIState) error {
	if apiState.APIDefinition.Spec.APIType == "REST" {
		return deployRestAPIInGateway(apiState)
	}
	if apiState.APIDefinition.Spec.APIType == "GraphQL" {
		return deployGQLAPIInGateway(apiState)
	}
	return nil
}

func undeployAPIInGateway(apiState APIState) error {
	if apiState.APIDefinition.Spec.APIType == "REST" {
		return undeployRestAPIInGateway(apiState)
	}
	if apiState.APIDefinition.Spec.APIType == "GraphQL" {
		return undeployGQLAPIInGateway(apiState)
	}
	return nil
}

// Runtime client connetion
var partitionClient *http.Client

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
		if !apiEvent.Event.APIDefinition.Spec.SystemAPI {
			apiDefinition := apiEvent.Event.APIDefinition
			loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v, API_UUID: %v", apiDefinition.Spec.APIName, string(apiDefinition.ObjectMeta.UID))
			api := apiEvent.Event
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
