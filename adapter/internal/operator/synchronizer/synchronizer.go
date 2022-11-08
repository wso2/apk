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

	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	model "github.com/wso2/apk/adapter/internal/oasparser/model"
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
		deployAPIInGateway(event.Event)
	}
}

// deployAPIInGateway deploys the related API in CREATE and UPDATE events.
func deployAPIInGateway(apiState APIState) {
	var mgwSwagger model.MgwSwagger
	if err := mgwSwagger.SetInfoAPICR(*apiState.APIDefinition); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error setting API CR info to mgwSwagger: %v", err)
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error setting API CR info to mgwSwagger: %v", err),
			Severity:  logging.MAJOR,
			ErrorCode: 2612,
		})
	}
	if err := mgwSwagger.SetInfoHTTPRouteCR(*apiState.ProdHTTPRoute); err != nil {
		loggers.LoggerAPKOperator.Errorf("Error setting HttpRoute CR info to mgwSwagger: %v", err)
		loggers.LoggerAPKOperator.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("error setting HttpRoute CR info to mgwSwagger: %v", err),
			Severity:  logging.MAJOR,
			ErrorCode: 2613,
		})
	}
	// mgwSwagger.UUID = string(apiState.APIDefinition.UID)
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
