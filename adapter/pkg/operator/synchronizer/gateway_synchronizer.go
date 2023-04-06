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
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/interceptor"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/operator/utils"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
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
	if gatewayState.GatewayDefinition != nil || gatewayState.GatewayStateData.GatewayResolvedListenerCerts != nil {
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
func AddOrUpdateGateway(gatewayState GatewayState, state string) (string, error) {
	gateway := gatewayState.GatewayDefinition
	resolvedListenerCerts := gatewayState.GatewayStateData.GatewayResolvedListenerCerts
	customRateLimitPolicies := getCustomRateLimitPolicies(gatewayState.GatewayStateData.GatewayCustomRateLimitPolicies)
	gatewayAPIPolicies := gatewayState.GatewayStateData.GatewayAPIPolicies
	gatewayBackendMapping := gatewayState.GatewayStateData.GatewayBackendMapping

	gwLuaScript, gwReqICluster, gwReqIAddresses, gwResICluster, gwResIAddresses :=
		generateGlobalInterceptorResource(gatewayAPIPolicies, gatewayBackendMapping)

	if state == constants.Create {
		xds.GenerateGlobalClusters(gatewayState.GatewayDefinition.Name)
	}
	if state == constants.Create || state == constants.Update {
		xds.GenerateGlobalClustersWithInterceptors(gateway.Name,
			gwReqICluster, gwReqIAddresses,
			gwResICluster, gwResIAddresses)
	}

	xds.UpdateGatewayCache(gateway, resolvedListenerCerts, gwLuaScript, customRateLimitPolicies)
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
		xds.UpdateRateLimitXDSCacheForCustomPolicies(gateway.Name, customRateLimitPolicies)
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

func generateGlobalInterceptorResource(gatewayAPIPolicies map[string]v1alpha1.APIPolicy,
	gatewayBackendMapping v1alpha1.BackendMapping) (string, *clusterv3.Cluster, []*corev3.Address,
	*clusterv3.Cluster, []*corev3.Address) {
	var gwLuaScript string
	var gwReqICluster, gwResICluster *clusterv3.Cluster
	var gwReqIAddresses, gwResIAddresses []*corev3.Address

	if len(gatewayAPIPolicies) > 0 && len(gatewayBackendMapping) > 0 {
		gwReqI, gwResI := createInterceptors(gatewayAPIPolicies, gatewayBackendMapping)
		if len(gwReqI) > 0 {
			gwReqICluster, gwReqIAddresses, _ = envoyconf.CreateLuaCluster(nil, gwReqI["POST"])
		}
		if len(gwResI) > 0 {
			gwResICluster, gwResIAddresses, _ = envoyconf.CreateLuaCluster(nil, gwResI["POST"])
		}
		gwLuaScript = getGlobalInterceptorScript(gatewayAPIPolicies, gatewayBackendMapping)
	}
	return gwLuaScript, gwReqICluster, gwReqIAddresses, gwResICluster, gwResIAddresses
}

func getGlobalInterceptorScript(gatewayAPIPolicies map[string]v1alpha1.APIPolicy,
	gatewayBackendMapping v1alpha1.BackendMapping) string {
	iInvCtx := &interceptor.InvocationContext{
		OrganizationID:   "",
		BasePath:         "",
		SupportedMethods: "",
		APIName:          "",
		APIVersion:       "",
		PathTemplate:     "",
		Vhost:            "",
		ClusterName:      "",
	}
	reqI, resI := createInterceptors(gatewayAPIPolicies, gatewayBackendMapping)
	if len(reqI) > 0 || len(resI) > 0 {
		return envoyconf.GetInlineLuaScript(reqI, resI, iInvCtx)
	}
	return `
function envoy_on_request(request_handle)
end
function envoy_on_response(response_handle)
end
`
}

func createInterceptors(gatewayAPIPolicies map[string]v1alpha1.APIPolicy,
	gatewayBackendMapping v1alpha1.BackendMapping) (requestInterceptor map[string]model.InterceptEndpoint, responseInterceptor map[string]model.InterceptEndpoint) {
	requestInterceptorMap := make(map[string]model.InterceptEndpoint)
	responseInterceptorMap := make(map[string]model.InterceptEndpoint)

	var apiPolicy *v1alpha1.APIPolicy
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(gatewayAPIPolicies)))
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
		resolvedPolicySpec := utils.SelectPolicy(&apiPolicy.Spec.Override, &apiPolicy.Spec.Default, nil, nil)
		if resolvedPolicySpec != nil {
			if resolvedPolicySpec.RequestInterceptor != nil {
				reqIEp := getInterceptorEndpoint(resolvedPolicySpec.RequestInterceptor, gatewayBackendMapping, true)
				if reqIEp != nil {
					requestInterceptorMap["POST"] = *reqIEp
				}
			}
			if resolvedPolicySpec.ResponseInterceptor != nil {
				resIEp := getInterceptorEndpoint(resolvedPolicySpec.ResponseInterceptor, gatewayBackendMapping, false)
				if resIEp != nil {
					responseInterceptorMap["POST"] = *resIEp
				}
			}
		}
	}
	return requestInterceptorMap, responseInterceptorMap
}

func getInterceptorEndpoint(interceptor *v1alpha1.InterceptorConfig, gatewayBackendMapping v1alpha1.BackendMapping, isReq bool) *model.InterceptEndpoint {
	endpoints := model.GetEndpoints(types.NamespacedName{Namespace: interceptor.BackendRef.Namespace, Name: interceptor.BackendRef.Name},
		gatewayBackendMapping)
	var clusterName string
	if isReq {
		clusterName = "request_interceptor_global_cluster"
	} else {
		clusterName = "response_interceptor_global_cluster"
	}
	if len(endpoints) > 0 {
		conf := config.ReadConfigs()
		clusterTimeoutV := conf.Envoy.ClusterTimeoutInSeconds
		requestTimeoutV := conf.Envoy.ClusterTimeoutInSeconds
		return &model.InterceptEndpoint{
			Enable:          true,
			ClusterName:     clusterName,
			EndpointCluster: model.EndpointCluster{Endpoints: endpoints},
			ClusterTimeout:  clusterTimeoutV,
			RequestTimeout:  requestTimeoutV,
			Includes:        model.GenerateInterceptorIncludes(interceptor.Includes),
		}
	}
	return nil
}
