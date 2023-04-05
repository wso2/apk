/*
 *  Copyright (c) 2023, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/logging"
	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// GatewayEvent holds the data structure used for passing Gateway
// events from controller go routine to synchronizer
// go routine.
type GatewayEvent struct {
	EventType string
	Event     GatewayState
}

// HandleGatewayLifeCycleEvents handles the Gateway events generated from OperatorDataStore
func HandleGatewayLifeCycleEvents(ch *chan GatewayEvent) {
	loggers.LoggerAPKOperator.Info("Operator synchronizer listening for Gateway lifecycle events...")
	for event := range *ch {
		if event.Event.GatewayDefinition == nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2628))
		}
		loggers.LoggerAPKOperator.Infof("%s event received for %v", event.EventType, event.Event.GatewayDefinition.Name)
		var err error
		switch event.EventType {
		case constants.Delete:
			err = undeployGateway(event.Event)
		case constants.Create:
			err = deployGateway(event.Event, constants.Create)
		case constants.Update:
			err = deployGateway(event.Event, constants.Update)
		}
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.GetErrorByCode(2629, event.EventType, err))
		}
	}
}

// deployGateway deploys the related Gateway in CREATE and UPDATE events.
func deployGateway(gatewayState GatewayState, state string) error {
	var err error
	if gatewayState.GatewayDefinition != nil {
		_, err = AddOrUpdateGateway(gatewayState, state)
	}
	return err
}

// undeployGateway undeploys the related Gateway in DELETE events.
func undeployGateway(gatewayState GatewayState) error {
	var err error
	if gatewayState.GatewayDefinition != nil {
		_, err = DeleteGateway(gatewayState.GatewayDefinition)
	}
	return err
}

// AddOrUpdateGateway adds/update a Gateway to the XDS server.
func AddOrUpdateGateway(gatewayState GatewayState,state string) (string, error) {
	gateway := gatewayState.GatewayDefinition
	customRateLimitPolicies := getCustomRateLimitPolicies(gatewayState.CustomRateLimitPolicies)
	if state == constants.Create {
		xds.GenerateGlobalClusters(gateway.Name)
	}
	xds.UpdateGatewayCache(gateway, customRateLimitPolicies)
	listeners, clusters, routes, endpoints, apis := xds.GenerateEnvoyResoucesForGateway(gateway.Name)
	loggers.LoggerAPKOperator.Debugf("listeners: %v", listeners)
	loggers.LoggerAPKOperator.Debugf("clusters: %v", clusters)
	loggers.LoggerAPKOperator.Debugf("routes: %v", routes)
	loggers.LoggerAPKOperator.Debugf("endpoints: %v", endpoints)
	loggers.LoggerAPKOperator.Debugf("apis: %v", apis)
	xds.UpdateXdsCacheWithLock(gateway.Name, endpoints, clusters, routes, listeners)
	xds.UpdateEnforcerApis(gateway.Name, apis, "")
	conf := config.ReadConfigs()
	if conf.Envoy.RateLimit.Enabled {
		xds.UpdateRateLimitXDSCacheForCustomPolicies(gateway.Name,customRateLimitPolicies)
	}
	return "", nil
}

// DeleteGateway deletes a Gateway from the XDS server.
func DeleteGateway(gateway *gwapiv1b1.Gateway) (string, error) {
	xds.UpdateXdsCacheWithLock(gateway.Name, nil, nil, nil, nil)
	xds.UpdateEnforcerApis(gateway.Name, nil, "")
	return "", nil
}

// getCustomRateLimitPolicies returns the custom rate limit policies.
func getCustomRateLimitPolicies(customRateLimitPoliciesDef []*dpv1alpha1.RateLimitPolicy) []*model.CustomRateLimitPolicy {
	var customRateLimitPolicies []*model.CustomRateLimitPolicy
	for _, customRateLimitPolicy := range customRateLimitPoliciesDef {
		customRLPolicy := model.ParseCustomRateLimitPolicy(*customRateLimitPolicy)
		customRateLimitPolicies = append(customRateLimitPolicies, customRLPolicy)
	}
	return customRateLimitPolicies
}
