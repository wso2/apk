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

	client "github.com/wso2/apk/adapter/internal/grpc-client"
	"github.com/wso2/apk/adapter/internal/loggers"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"
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
		loggers.LoggerAPKOperator.Infof("Event received: %v\n", event)
		if err := deployAPIInGateway(event.Event); err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("api deployment failed:%v", err),
				ErrorCode: 2616,
				Severity:  logging.MAJOR,
			})
		}
		go sendAPIToAPKMgtServer(event)
	}
}

// deployAPIInGateway deploys the related API in CREATE and UPDATE events.
func deployAPIInGateway(apiState APIState) error {
	var mgwSwagger model.MgwSwagger
	if err := mgwSwagger.SetInfoAPICR(*apiState.APIDefinition); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error setting API CR info to mgwSwagger: %v", err),
			Severity:  logging.MAJOR,
			ErrorCode: 2612,
		})
		return err
	}
	if err := mgwSwagger.SetInfoHTTPRouteCR(*apiState.ProdHTTPRoute); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error setting HttpRoute CR info to mgwSwagger: %v", err),
			Severity:  logging.MAJOR,
			ErrorCode: 2613,
		})
		return err
	}
	if err := mgwSwagger.ValidateIR(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error validating mgwSwagger intermediate representation: %v", err),
			Severity:  logging.MAJOR,
			ErrorCode: 2615,
		})
		return err
	}
	vHosts := getVhostForAPI(apiState)
	labels := getLabelsForAPI(apiState)
	for _, vHost := range vHosts {
		err := xds.UpdateAPICache(vHost, labels, mgwSwagger)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("error updating the API cache: %v", err),
				Severity:  logging.MAJOR,
				ErrorCode: 2614,
			})
		}
	}
	return nil
}

// getVhostForAPI returns the vHosts related to an API.
func getVhostForAPI(api APIState) []string {
	var vHosts []string
	for _, hostName := range api.ProdHTTPRoute.Spec.Hostnames {
		vHosts = append(vHosts, string(hostName))
	}
	return vHosts
}

// getLabelsForAPI returns the labels related to an API.
func getLabelsForAPI(api APIState) []string {
	var labels []string
	for _, parentRef := range api.ProdHTTPRoute.Spec.ParentRefs {
		labels = append(labels, string(parentRef.Name))
	}
	return labels
}

// sendAPIToAPKMgtServer sends the API create/update/delete event to the APK management server.
func sendAPIToAPKMgtServer(apiEvent APIEvent) {
	loggers.LoggerAPKOperator.Infof("Sending API to APK management server:%v", apiEvent.Event.APIDefinition.Spec.APIDisplayName)
	conn, err := client.GetConnection()
	api := apiEvent.Event
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error creating connection: %v", err),
			ErrorCode: 6000,
			Severity:  logging.BLOCKER,
		})
	}
	res, err := client.ExecuteGRPCCall(conn, func() (interface{}, error) {
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
		}
		return nil, nil
	})
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error sending API to APK management server:%v", err),
			ErrorCode: 6001,
			Severity:  logging.MAJOR,
		})
	}
	loggers.LoggerAPKOperator.Info(res)
}

// getResourcesForAPI returns []*apiProtos.Resource for HTTPRoute
// resources. Temporary method added until a proper implementation is done.
func getResourcesForAPI(api APIState) []*apiProtos.Resource {
	var resources []*apiProtos.Resource
	var hostNames []string
	for _, hostName := range api.ProdHTTPRoute.Spec.Hostnames {
		hostNames = append(hostNames, string(hostName))
	}
	for _, rule := range api.ProdHTTPRoute.Spec.Rules {
		for _, match := range rule.Matches {
			resources = append(resources, &apiProtos.Resource{Path: *match.Path.Value, Hostname: hostNames})
		}
	}
	return resources
}
