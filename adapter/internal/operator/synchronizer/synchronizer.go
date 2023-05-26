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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"context"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	client "github.com/wso2/apk/adapter/internal/management-server/grpc-client"
	"k8s.io/apimachinery/pkg/types"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/services/runtime"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"
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
	mgtServerCh chan APIEvent
	paritionCh  chan APIEvent
)

func init() {
	mgtServerCh = make(chan APIEvent, 10)
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
			*successChannel <- SuccessEvent{
				APINamespacedName: utils.NamespacedName(event.Event.APIDefinition),
				State:             event.EventType,
				Events:            event.UpdatedEvents,
			}
			if config.ReadConfigs().ManagementServer.Enabled {
				mgtServerCh <- event
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
	adapterInternalAPI.SetInfoAPICR(*apiState.APIDefinition)
	adapterInternalAPI.SetAPIDefinitionFile(apiState.APIDefinitionFile)
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
		loggers.LoggerAPKOperator.Debugf("Recieved Parent Refs:%v", parentRef)
		loggers.LoggerAPKOperator.Debugf("Recieved Parent Refs Section Name:%v", string(*parentRef.SectionName))
		listeners = append(listeners, string(*parentRef.SectionName))
	}
	return listeners
}

// SendAPIToAPKMgtServer sends the API create/update/delete event to the APK management server.
func SendAPIToAPKMgtServer() {
	loggers.LoggerAPKOperator.Info("Start listening for API to APK management server events")
	conf := config.ReadConfigs()
	address := fmt.Sprintf("%s:%d", conf.ManagementServer.Host, conf.ManagementServer.GRPCClient.Port)
	for apiEvent := range mgtServerCh {
		if !apiEvent.Event.APIDefinition.Spec.SystemAPI {
			loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v", apiEvent.Event.APIDefinition.Spec.APIDisplayName)
			api := apiEvent.Event
			conn, err := client.GetConnection(address)
			apiClient := apiProtos.NewAPIServiceClient(conn)
			if err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2634, address, err))
			}
			_, err = client.ExecuteGRPCCall(func() (interface{}, error) {
				var definition string
				var errDef error
				definition, errDef = runtime.GetAPIDefinition(string(api.APIDefinition.GetUID()), api.APIDefinition.Spec.Organization)
				if strings.Compare(apiEvent.EventType, constants.Create) == 0 {
					if errDef != nil {
						return nil, err
					}
					return apiClient.CreateAPI(context.Background(), &apiProtos.API{
						Uuid:           string(api.APIDefinition.GetUID()),
						Version:        api.APIDefinition.Spec.APIVersion,
						Name:           api.APIDefinition.Spec.APIDisplayName,
						Context:        api.APIDefinition.Spec.Context,
						Type:           api.APIDefinition.Spec.APIType,
						OrganizationId: api.APIDefinition.Spec.Organization,
						Resources:      getResourcesForAPI(api),
						Definition:     definition,
					})
				} else if strings.Compare(apiEvent.EventType, constants.Update) == 0 {
					if errDef != nil {
						return nil, err
					}
					return apiClient.UpdateAPI(context.Background(), &apiProtos.API{
						Uuid:           string(api.APIDefinition.GetUID()),
						Version:        api.APIDefinition.Spec.APIVersion,
						Name:           api.APIDefinition.Spec.APIDisplayName,
						Context:        api.APIDefinition.Spec.Context,
						Type:           api.APIDefinition.Spec.APIType,
						OrganizationId: api.APIDefinition.Spec.Organization,
						Resources:      getResourcesForAPI(api),
						Definition:     definition,
					})
				} else if strings.Compare(apiEvent.EventType, constants.Delete) == 0 {
					return apiClient.DeleteAPI(context.Background(), &apiProtos.API{
						Uuid:           string(api.APIDefinition.GetUID()),
						Version:        api.APIDefinition.Spec.APIVersion,
						Name:           api.APIDefinition.Spec.APIDisplayName,
						Context:        api.APIDefinition.Spec.Context,
						Type:           api.APIDefinition.Spec.APIType,
						OrganizationId: api.APIDefinition.Spec.Organization,
						Resources:      getResourcesForAPI(api),
					})
				}
				return nil, nil
			})
			if err != nil {
				loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2635, err))
			}
		}
	}
}

// SendEventToPartitionServer sends the API create/update/delete event to the partition server.
func SendEventToPartitionServer() {
	conf := config.ReadConfigs()
	for apiEvent := range paritionCh {
		if !apiEvent.Event.APIDefinition.Spec.SystemAPI {
			loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v", apiEvent.Event.APIDefinition.Spec.APIDisplayName)
			api := apiEvent.Event
			eventType := apiEvent.EventType
			context := api.APIDefinition.Spec.Context
			organization := api.APIDefinition.Spec.Organization
			version := api.APIDefinition.Spec.APIVersion
			apiName := api.APIDefinition.Spec.APIDisplayName
			apiUUID := string(api.APIDefinition.GetUID())
			var hostNames []string
			httpRoute := api.ProdHTTPRoute
			if httpRoute == nil {
				httpRoute = api.SandHTTPRoute
			}
			for _, hostName := range httpRoute.HTTPRoute.Spec.Hostnames {
				hostNames = append(hostNames, string(hostName))
			}
			data := PartitionEvent{
				EventType:      eventType,
				APIContext:     context,
				OrganizationID: organization,
				APIVersion:     version,
				APIName:        apiName,
				APIUUID:        apiUUID,
				Vhosts:         hostNames,
				PartitionID:    conf.PartitionServer.PartitionName,
			}
			payload, err := json.Marshal(data)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error creating Event: %v", err)
			}
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", conf.PartitionServer.Host, conf.PartitionServer.ServiceBasePath), bytes.NewBuffer(payload))
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error creating api definition request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				loggers.LoggerAPKOperator.Errorf("Error sending API Event: %v", err)
			}
			if resp.StatusCode != http.StatusAccepted {
				loggers.LoggerAPKOperator.Info("API Event Accepted", resp.Status)
			}
		}
	}
}

// PartitionEvent is the event sent to the partition server.
type PartitionEvent struct {
	EventType      string   `json:"eventType"`
	APIName        string   `json:"apiName"`
	APIVersion     string   `json:"apiVersion"`
	APIContext     string   `json:"apiContext"`
	OrganizationID string   `json:"organizationId"`
	PartitionID    string   `json:"partitionId"`
	APIUUID        string   `json:"apiId"`
	Vhosts         []string `json:"vhosts"`
}

// getResourcesForAPI returns []*apiProtos.Resource for HTTPRoute
// resources. Temporary method added until a proper implementation is done.
func getResourcesForAPI(api APIState) []*apiProtos.Resource {
	var resources []*apiProtos.Resource
	var hostNames []string
	httpRoute := api.ProdHTTPRoute
	if httpRoute == nil {
		httpRoute = api.SandHTTPRoute
	}
	for _, hostName := range httpRoute.HTTPRoute.Spec.Hostnames {
		hostNames = append(hostNames, string(hostName))
	}
	for _, rule := range httpRoute.HTTPRoute.Spec.Rules {
		for _, match := range rule.Matches {
			resource := &apiProtos.Resource{
				Path:     *match.Path.Value,
				Hostname: hostNames,
			}
			if match.Method != nil {
				resource.Verb = string(*match.Method)
			}
			resources = append(resources, resource)
		}
	}
	return resources
}
