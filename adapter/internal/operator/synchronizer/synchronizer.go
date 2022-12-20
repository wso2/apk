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

	"github.com/wso2/apk/adapter/config"
	client "github.com/wso2/apk/adapter/internal/grpc-client"
	"github.com/wso2/apk/adapter/internal/loggers"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// APIEvent holds the data structure used for passing API
// lifecycle events from controller go routine to synchronizer
// go routine.
type APIEvent struct {
	EventType string
	Event     APIState
}

// HandleAPILifeCycleEvents handles the API events generated from OperatorDataStore
func HandleAPILifeCycleEvents(ch *chan APIEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for API lifecycle events...")
	for event := range *ch {
		if event.Event.APIDefinition == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   "API Event is nil",
				Severity:  logging.BLOCKER,
				ErrorCode: 2600,
			})
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
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("API deployment failed for %s event : %v", event.EventType, err),
				ErrorCode: 2616,
				Severity:  logging.MAJOR,
			})
		} else {
			//TODO (amali) fix creating go routine per event
			go sendAPIToAPKMgtServer(event)
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
		loggers.LoggerXds.ErrorC(logging.ErrorDetails{
			Message: fmt.Sprintf("Error undeploying prod httpRoute of API : %v in Organization %v from environments %v."+
				" Hence not checking on deleting the sand httpRoute of the API",
				string(apiState.APIDefinition.ObjectMeta.UID), apiState.APIDefinition.Spec.Organization,
				getLabelsForAPI(apiState.ProdHTTPRoute.HTTPRoute)),
			Severity:  logging.MAJOR,
			ErrorCode: 1415,
		})
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
	vHosts := getVhostsForAPI(httpRoute)
	uuid := string(apiState.APIDefinition.ObjectMeta.UID)
	return xds.DeleteAPICREvent(vHosts, labels, uuid, org)
}

// deployAPIInGateway deploys the related API in CREATE and UPDATE events.
func deployAPIInGateway(apiState APIState) error {
	var err error
	if apiState.ProdHTTPRoute != nil {
		err = generateMGWSwagger(apiState, apiState.ProdHTTPRoute, true)
	}
	if err != nil {
		return err
	}
	if apiState.SandHTTPRoute != nil {
		err = generateMGWSwagger(apiState, apiState.SandHTTPRoute, false)
	}
	return err
}

// generateMGWSwagger this will populate a mgwswagger representation for an HTTPRoute
func generateMGWSwagger(apiState APIState, httpRoute *HTTPRouteState, isProd bool) error {
	var mgwSwagger model.MgwSwagger
	mgwSwagger.SetInfoAPICR(*apiState.APIDefinition)
	if err := mgwSwagger.SetInfoHTTPRouteCR(httpRoute.HTTPRoute, httpRoute.Authentications, httpRoute.ResourceAuthentications,
		isProd); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error setting HttpRoute CR info to mgwSwagger for isProdEndpoint: %v. %v", isProd, err),
			Severity:  logging.MAJOR,
			ErrorCode: 2613,
		})
		return err
	}
	if err := mgwSwagger.Validate(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message: fmt.Sprintf("Error validating mgwSwagger intermediate representation for isProdEndpoint: %v. %v",
				isProd, err),
			Severity:  logging.MAJOR,
			ErrorCode: 2615,
		})
		return err
	}
	vHosts := getVhostsForAPI(httpRoute.HTTPRoute)
	labels := getLabelsForAPI(httpRoute.HTTPRoute)
	for _, vHost := range vHosts {
		err := xds.UpdateAPICache(vHost, labels, mgwSwagger)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message: fmt.Sprintf("Error updating the API : %s:%s in vhost: %s. %v",
					mgwSwagger.GetTitle(), mgwSwagger.GetVersion(), vHost, err),
				Severity:  logging.MAJOR,
				ErrorCode: 2614,
			})
		}
	}
	return nil
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

// sendAPIToAPKMgtServer sends the API create/update/delete event to the APK management server.
func sendAPIToAPKMgtServer(apiEvent APIEvent) {
	loggers.LoggerAPKOperator.Infof("Sending API to APK management server: %v", apiEvent.Event.APIDefinition.Spec.APIDisplayName)
	conf := config.ReadConfigs()
	address := conf.Adapter.GRPCClient.ManagementServerAddress
	conn, err := client.GetConnection(address)
	api := apiEvent.Event
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error creating connection for %v : %v", address, err),
			ErrorCode: 6000,
			Severity:  logging.BLOCKER,
		})
	}
	_, err = client.ExecuteGRPCCall(conn, func() (interface{}, error) {
		apiClient := apiProtos.NewAPIServiceClient(conn)
		if strings.Compare(apiEvent.EventType, constants.Create) == 0 {
			return apiClient.CreateAPI(context.Background(), &apiProtos.API{
				Uuid:           string(api.APIDefinition.GetUID()),
				Version:        api.APIDefinition.Spec.APIVersion,
				Name:           api.APIDefinition.Spec.APIDisplayName,
				Context:        api.APIDefinition.Spec.Context,
				Type:           api.APIDefinition.Spec.APIType,
				OrganizationId: api.APIDefinition.Spec.Organization,
				Resources:      getResourcesForAPI(api),
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
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error sending API to APK management server : %v", err),
			ErrorCode: 6001,
			Severity:  logging.MAJOR,
		})
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
