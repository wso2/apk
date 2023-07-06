/*
 *  Copyright (c) 2022, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"k8s.io/apimachinery/pkg/types"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	internalLogging "github.com/wso2/apk/adapter/internal/logging"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
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
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2628))
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
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2629, event.EventType, err))
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

// deleteAPIInGateway undeploys the related API in CREATE and UPDATE events.
func undeployAPIInGateway(apiState APIState) error {
	var err error
	if apiState.ProdHTTPRoute != nil {
		err = deleteAPIFromEnv(apiState.ProdHTTPRoute.HTTPRoute, apiState)
	}
	if err != nil {
		loggers.LoggerXds.ErrorC(logging.GetErrorByCode(2630, string(apiState.APIDefinition.ObjectMeta.UID), apiState.APIDefinition.Spec.Organization,
			getLabelsForAPI(apiState.ProdHTTPRoute.HTTPRoute)))
		return err
	}
	if apiState.SandHTTPRoute != nil {
		err = deleteAPIFromEnv(apiState.SandHTTPRoute.HTTPRoute, apiState)
	}
	return err
}

func deleteAPIFromEnv(httpRoute *gwapiv1b1.HTTPRoute, apiState APIState) error {
	labels := getLabelsForAPI(httpRoute)
	org := apiState.APIDefinition.Spec.Organization
	uuid := string(apiState.APIDefinition.ObjectMeta.UID)
	return xds.DeleteAPICREvent(labels, uuid, org)
}

// deployAPIInGateway deploys the related API in CREATE and UPDATE events.
func deployAPIInGateway(apiState APIState) error {
	var err error
	if apiState.ProdHTTPRoute != nil {
		_, err = GenerateAdapterInternalAPI(apiState, apiState.ProdHTTPRoute, constants.Production)
	}
	if err != nil {
		return err
	}
	if apiState.SandHTTPRoute != nil {
		_, err = GenerateAdapterInternalAPI(apiState, apiState.SandHTTPRoute, constants.Sandbox)
	}
	return err
}

// GenerateAdapterInternalAPI this will populate a AdapterInternalAPI representation for an HTTPRoute
func GenerateAdapterInternalAPI(apiState APIState, httpRoute *HTTPRouteState, envType string) (*model.AdapterInternalAPI, error) {
	var adapterInternalAPI model.AdapterInternalAPI
	adapterInternalAPI.SetIsDefaultVersion(apiState.APIDefinition.Spec.IsDefaultVersion)
	adapterInternalAPI.SetInfoAPICR(*apiState.APIDefinition)
	adapterInternalAPI.SetAPIDefinitionFile(apiState.APIDefinitionFile)
	internalLogging.SetValueToLogContext("API_UUID", adapterInternalAPI.UUID)
	adapterInternalAPI.EnvType = envType
	httpRouteParams := model.HTTPRouteParams{
		AuthSchemes:               httpRoute.Authentications,
		ResourceAuthSchemes:       httpRoute.ResourceAuthentications,
		BackendMapping:            httpRoute.BackendMapping,
		APIPolicies:               httpRoute.APIPolicies,
		ResourceAPIPolicies:       httpRoute.ResourceAPIPolicies,
		ResourceScopes:            httpRoute.Scopes,
		InterceptorServiceMapping: httpRoute.InterceptorServiceMapping,
		RateLimitPolicies:         httpRoute.RateLimitPolicies,
		ResourceRateLimitPolicies: httpRoute.ResourceRateLimitPolicies,
	}
	if err := adapterInternalAPI.SetInfoHTTPRouteCR(httpRoute.HTTPRoute, httpRouteParams); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2631, err))
		return nil, err
	}
	if err := adapterInternalAPI.Validate(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2632, err))
		return nil, err
	}
	vHosts := getVhostsForAPI(httpRoute.HTTPRoute)
	labels := getLabelsForAPI(httpRoute.HTTPRoute)
	listeners := getListenersForAPI(httpRoute.HTTPRoute)

	conf := config.ReadConfigs()
	if conf.Envoy.RateLimit.Enabled {
		xds.UpdateRateLimitXDSCache(vHosts, adapterInternalAPI)
	}
	err := xds.UpdateAPICache(vHosts, labels, listeners, adapterInternalAPI)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2633, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), vHosts, err))
	}
	return &adapterInternalAPI, nil
}

// getVhostForAPI returns the vHosts related to an API.
func getVhostsForAPI(httpRoute *gwapiv1b1.HTTPRoute) []string {
	var vHosts []string
	for _, hostName := range httpRoute.Spec.Hostnames {
		vHosts = append(vHosts, string(hostName))
	}
	fmt.Println("vhosts size: ", len(vHosts))
	return vHosts
}

// getLabelsForAPI returns the labels related to an API.
func getLabelsForAPI(httpRoute *gwapiv1b1.HTTPRoute) []string {
	var labels []string
	var err error
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		err = xds.SanitizeGateway(string(parentRef.Name), false)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2653, string(parentRef.Name)))
		} else {
			labels = append(labels, string(parentRef.Name))
		}
	}
	return labels
}

// getListenersForAPI returns the listeners related to an API.
func getListenersForAPI(httpRoute *gwapiv1b1.HTTPRoute) []string {
	var listeners []string
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		loggers.LoggerAPKOperator.Debugf("Recieved Parent Refs:%v, API_UUID: %v", parentRef, internalLogging.GetValueFromLogContext("API_UUID"))
		loggers.LoggerAPKOperator.Debugf("Recieved Parent Refs Section Name:%v, API_UUID: %v", string(*parentRef.SectionName), internalLogging.GetValueFromLogContext("API_UUID"))
		listeners = append(listeners, string(*parentRef.SectionName))
	}
	return listeners
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
			loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v, API_UUID: %v", apiEvent.Event.APIDefinition.Spec.APIDisplayName, internalLogging.GetValueFromLogContext("API_UUID"))
			api := apiEvent.Event
			eventType := apiEvent.EventType
			context := api.APIDefinition.Spec.Context
			organization := api.APIDefinition.Spec.Organization
			version := api.APIDefinition.Spec.APIVersion
			apiName := api.APIDefinition.Spec.APIDisplayName
			apiUUID := string(api.APIDefinition.Name)
			var hostNames []string
			httpRoute := api.ProdHTTPRoute
			if httpRoute == nil {
				httpRoute = api.SandHTTPRoute
			}
			for _, hostName := range httpRoute.HTTPRoute.Spec.Hostnames {
				hostNames = append(hostNames, string(hostName))
			}
			data := PartitionEvent{
				EventType:    eventType,
				APIContext:   context,
				Organization: organization,
				APIVersion:   version,
				APIName:      apiName,
				APIUUID:      apiUUID,
				Vhosts:       hostNames,
				Partition:    conf.PartitionServer.PartitionName,
			}
			payload, err := json.Marshal(data)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error creating Event: %v, API_UUID: %v", err, internalLogging.GetValueFromLogContext("API_UUID"))
			}
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", conf.PartitionServer.Host, conf.PartitionServer.ServiceBasePath, "/api-deployment"), bytes.NewBuffer(payload))
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error creating api definition request: %v, API_UUID: %v", err, internalLogging.GetValueFromLogContext("API_UUID"))
			}
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			resp, err := partitionClient.Do(req)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error sending API Event: %v, API_UUID: %v", err, internalLogging.GetValueFromLogContext("API_UUID"))
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
	APIContext   string   `json:"apiContext"`
	Organization string   `json:"organization"`
	Partition    string   `json:"partition"`
	APIUUID      string   `json:"apiId"`
	Vhosts       []string `json:"vhosts"`
}
