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
	"fmt"
	"strings"

	"context"

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	client "github.com/wso2/apk/adapter/internal/management-server/grpc-client"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/operator/services/runtime"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// APIEvent holds the data structure used for passing API
// lifecycle events from controller go routine to synchronizer
// go routine.
type APIEvent struct {
	EventType string
	Event     APIState
}

var (
	// TODO: Decide on a buffer size and add to config.
	mgtServerCh chan APIEvent
)

func init() {
	mgtServerCh = make(chan APIEvent, 10)
}

// HandleAPILifeCycleEvents handles the API events generated from OperatorDataStore
func HandleAPILifeCycleEvents(ch *chan APIEvent) {
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
			if config.ReadConfigs().ManagementServer.Enabled {
				mgtServerCh <- event
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
		_, err = GenerateMGWSwagger(apiState, apiState.ProdHTTPRoute, constants.Production)
	}
	if err != nil {
		return err
	}
	if apiState.SandHTTPRoute != nil {
		_, err = GenerateMGWSwagger(apiState, apiState.SandHTTPRoute, constants.Sandbox)
	}
	return err
}

// GenerateMGWSwagger this will populate a mgwswagger representation for an HTTPRoute
func GenerateMGWSwagger(apiState APIState, httpRoute *HTTPRouteState, envType string) (*model.MgwSwagger, error) {
	var mgwSwagger model.MgwSwagger
	mgwSwagger.SetInfoAPICR(*apiState.APIDefinition)
	mgwSwagger.EnvType = envType
	httpRouteParams := model.HTTPRouteParams{
		AuthSchemes:               httpRoute.Authentications,
		ResourceAuthSchemes:       httpRoute.ResourceAuthentications,
		BackendMapping:            httpRoute.BackendMapping,
		APIPolicies:               httpRoute.APIPolicies,
		ResourceAPIPolicies:       httpRoute.ResourceAPIPolicies,
		ResourceScopes:            httpRoute.Scopes,
		RateLimitPolicies:         httpRoute.RateLimitPolicies,
		ResourceRateLimitPolicies: httpRoute.ResourceRateLimitPolicies,
	}
	if err := mgwSwagger.SetInfoHTTPRouteCR(httpRoute.HTTPRoute, httpRouteParams); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2631, err))
		return nil, err
	}
	if err := mgwSwagger.Validate(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2632, err))
		return nil, err
	}
	vHosts := getVhostsForAPI(httpRoute.HTTPRoute)
	labels := getLabelsForAPI(httpRoute.HTTPRoute)
	listeners := getListenersForAPI(httpRoute.HTTPRoute)

	conf := config.ReadConfigs()
	if conf.Envoy.RateLimit.Enabled {
		xds.UpdateRateLimitXDSCache(vHosts, mgwSwagger)
	}
	err := xds.UpdateAPICache(vHosts, labels, listeners, mgwSwagger)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2633, mgwSwagger.GetTitle(), mgwSwagger.GetVersion(), vHosts, err))
	}
	return &mgwSwagger, nil
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
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		labels = append(labels, string(parentRef.Name))
	}
	return labels
}

// getListenersForAPI returns the listeners related to an API.
func getListenersForAPI(httpRoute *gwapiv1b1.HTTPRoute) []string {
	var listeners []string
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		loggers.LoggerAPKOperator.Info("Recieved Parent Refs:%v", parentRef)
		loggers.LoggerAPKOperator.Info("Recieved Parent Refs Section Name:%v", string(*parentRef.SectionName))
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
				if strings.Compare(apiEvent.EventType, constants.Create) == 0 {
					return apiClient.CreateAPI(context.Background(), &apiProtos.API{
						Uuid:           string(api.APIDefinition.GetUID()),
						Version:        api.APIDefinition.Spec.APIVersion,
						Name:           api.APIDefinition.Spec.APIDisplayName,
						Context:        api.APIDefinition.Spec.Context,
						Type:           api.APIDefinition.Spec.APIType,
						OrganizationId: api.APIDefinition.Spec.Organization,
						Resources:      getResourcesForAPI(api),
						Definition:     runtime.GetAPIDefinition(string(api.APIDefinition.GetUID()), api.APIDefinition.Spec.Organization),
					})
				} else if strings.Compare(apiEvent.EventType, constants.Update) == 0 {
					return apiClient.UpdateAPI(context.Background(), &apiProtos.API{
						Uuid:           string(api.APIDefinition.GetUID()),
						Version:        api.APIDefinition.Spec.APIVersion,
						Name:           api.APIDefinition.Spec.APIDisplayName,
						Context:        api.APIDefinition.Spec.Context,
						Type:           api.APIDefinition.Spec.APIType,
						OrganizationId: api.APIDefinition.Spec.Organization,
						Resources:      getResourcesForAPI(api),
						Definition:     runtime.GetAPIDefinition(string(api.APIDefinition.GetUID()), api.APIDefinition.Spec.Organization),
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
