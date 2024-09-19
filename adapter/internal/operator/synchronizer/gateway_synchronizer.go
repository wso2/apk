/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	dataHolder "github.com/wso2/apk/adapter/internal/dataholder"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/interceptor"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"github.com/wso2/apk/adapter/pkg/logging"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
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
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2628, logging.CRITICAL, "Gateway definition is nil in the event : %v", event.EventType))
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
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2629, logging.MAJOR, "API deployment failed for %s event : %v", event.EventType, err))
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
		dataHolder.RemoveGateway(*gatewayState.GatewayDefinition)
	}
	return err
}

// AddOrUpdateGateway adds/update a Gateway to the XDS server.
func AddOrUpdateGateway(gatewayState GatewayState, state string) (string, error) {
	gateway := gatewayState.GatewayDefinition
	dataHolder.UpdateGateway(*gateway)
	xds.SanitizeGateway(gateway.Name, true)
	resolvedListenerCerts := gatewayState.GatewayStateData.GatewayResolvedListenerCerts
	customRateLimitPolicies := getCustomRateLimitPolicies(gatewayState.GatewayStateData.GatewayCustomRateLimitPolicies)
	gatewayAPIPolicies := gatewayState.GatewayStateData.GatewayAPIPolicies
	gatewayBackendMapping := gatewayState.GatewayStateData.GatewayBackendMapping
	gatewayInterceptorServiceMapping := gatewayState.GatewayStateData.GatewayInterceptorServiceMapping

	gwLuaScript, gwReqICluster, gwReqIAddresses, gwResICluster, gwResIAddresses :=
		generateGlobalInterceptorResource(gatewayAPIPolicies, gatewayInterceptorServiceMapping, gatewayBackendMapping)

	if state == constants.Create {
		xds.GenerateGlobalClusters(gateway.Name)
	}
	listeners, clusters, routes, endpoints, apis := xds.GenerateEnvoyResoucesForGateway(gateway.Name)
	if !config.ReadConfigs().Adapter.EnableGatewayClassController {
		xds.GenerateInterceptorClusters(gateway.Name, gwReqICluster, gwReqIAddresses, gwResICluster, gwResIAddresses)
		xds.UpdateGatewayCache(gateway, resolvedListenerCerts, gwLuaScript, customRateLimitPolicies)
		loggers.LoggerAPKOperator.Debugf("listeners: %v", listeners)
		loggers.LoggerAPKOperator.Debugf("clusters: %v", clusters)
		loggers.LoggerAPKOperator.Debugf("routes: %v", routes)
		loggers.LoggerAPKOperator.Debugf("endpoints: %v", endpoints)
		loggers.LoggerAPKOperator.Debugf("apis: %v", apis)
		xds.UpdateXdsCacheWithLock(gateway.Name, endpoints, clusters, routes, listeners)
	}
	xds.UpdateEnforcerApis(gateway.Name, apis, "")
	return "", nil
}

// DeleteGateway deletes a Gateway from the XDS server.
func DeleteGateway(gateway *gwapiv1.Gateway) (string, error) {
	xds.UpdateXdsCacheWithLock(gateway.Name, nil, nil, nil, nil)
	xds.UpdateEnforcerApis(gateway.Name, nil, "")
	return "", nil
}

// getCustomRateLimitPolicies returns the custom rate limit policies.
func getCustomRateLimitPolicies(customRateLimitPoliciesDef map[string]*dpv1alpha3.RateLimitPolicy) []*model.CustomRateLimitPolicy {
	var customRateLimitPolicies []*model.CustomRateLimitPolicy
	for _, customRateLimitPolicy := range customRateLimitPoliciesDef {
		customRLPolicy := model.ParseCustomRateLimitPolicy(*customRateLimitPolicy)
		customRateLimitPolicies = append(customRateLimitPolicies, customRLPolicy)
	}
	return customRateLimitPolicies
}

func generateGlobalInterceptorResource(gatewayAPIPolicies map[string]dpv1alpha3.APIPolicy,
	gatewayInterceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	gatewayBackendMapping map[string]*dpv1alpha2.ResolvedBackend) (string, *clusterv3.Cluster, []*corev3.Address,
	*clusterv3.Cluster, []*corev3.Address) {
	var gwLuaScript string
	var gwReqICluster, gwResICluster *clusterv3.Cluster
	var gwReqIAddresses, gwResIAddresses []*corev3.Address

	if len(gatewayAPIPolicies) > 0 && len(gatewayBackendMapping) > 0 {
		gwReqI, gwResI := createInterceptors(gatewayAPIPolicies, gatewayInterceptorServiceMapping, gatewayBackendMapping)
		if len(gwReqI) > 0 {
			gwReqICluster, gwReqIAddresses, _ = envoyconf.CreateLuaCluster(nil, gwReqI[string(gwapiv1.HTTPMethodPost)])
		}
		if len(gwResI) > 0 {
			gwResICluster, gwResIAddresses, _ = envoyconf.CreateLuaCluster(nil, gwResI[string(gwapiv1.HTTPMethodPost)])
		}
		gwLuaScript = getGlobalInterceptorScript(gatewayAPIPolicies, gatewayInterceptorServiceMapping, gatewayBackendMapping)
	}
	return gwLuaScript, gwReqICluster, gwReqIAddresses, gwResICluster, gwResIAddresses
}

func getGlobalInterceptorScript(gatewayAPIPolicies map[string]dpv1alpha3.APIPolicy,
	gatewayInterceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	gatewayBackendMapping map[string]*dpv1alpha2.ResolvedBackend) string {
	iInvCtx := &interceptor.InvocationContext{
		OrganizationID:   "",
		BasePath:         "",
		SupportedMethods: "",
		Environment:      "",
		APIName:          "",
		APIVersion:       "",
		PathTemplate:     "",
		Vhost:            "",
		ClusterName:      "",
		APIProperties:    "",
	}
	reqI, resI := createInterceptors(gatewayAPIPolicies, gatewayInterceptorServiceMapping, gatewayBackendMapping)
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

func createInterceptors(gatewayAPIPolicies map[string]dpv1alpha3.APIPolicy,
	gatewayInterceptorServiceMapping map[string]dpv1alpha1.InterceptorService,
	gatewayBackendMapping map[string]*dpv1alpha2.ResolvedBackend) (requestInterceptor map[string]model.InterceptEndpoint, responseInterceptor map[string]model.InterceptEndpoint) {
	requestInterceptorMap := make(map[string]model.InterceptEndpoint)
	responseInterceptorMap := make(map[string]model.InterceptEndpoint)

	var apiPolicy *dpv1alpha3.APIPolicy
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(gatewayAPIPolicies)))
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
		resolvedPolicySpec := utils.SelectPolicy(&apiPolicy.Spec.Override, &apiPolicy.Spec.Default, nil, nil)
		if resolvedPolicySpec != nil {
			if len(resolvedPolicySpec.RequestInterceptors) > 0 {
				reqIEp := getInterceptorEndpoint(apiPolicy.Namespace, &resolvedPolicySpec.RequestInterceptors[0], gatewayInterceptorServiceMapping,
					gatewayBackendMapping, true)
				if reqIEp != nil {
					requestInterceptorMap[string(gwapiv1.HTTPMethodPost)] = *reqIEp
					requestInterceptorMap[string(gwapiv1.HTTPMethodGet)] = *reqIEp
					requestInterceptorMap[string(gwapiv1.HTTPMethodDelete)] = *reqIEp
					requestInterceptorMap[string(gwapiv1.HTTPMethodPatch)] = *reqIEp
					requestInterceptorMap[string(gwapiv1.HTTPMethodPut)] = *reqIEp
					requestInterceptorMap[string(gwapiv1.HTTPMethodHead)] = *reqIEp
					requestInterceptorMap[string(gwapiv1.HTTPMethodOptions)] = *reqIEp
				}
			}
			if len(resolvedPolicySpec.ResponseInterceptors) > 0 {
				resIEp := getInterceptorEndpoint(apiPolicy.Namespace, &resolvedPolicySpec.ResponseInterceptors[0], gatewayInterceptorServiceMapping,
					gatewayBackendMapping, false)
				if resIEp != nil {
					responseInterceptorMap[string(gwapiv1.HTTPMethodPost)] = *resIEp
					responseInterceptorMap[string(gwapiv1.HTTPMethodGet)] = *resIEp
					responseInterceptorMap[string(gwapiv1.HTTPMethodDelete)] = *resIEp
					responseInterceptorMap[string(gwapiv1.HTTPMethodPatch)] = *resIEp
					responseInterceptorMap[string(gwapiv1.HTTPMethodPut)] = *resIEp
					responseInterceptorMap[string(gwapiv1.HTTPMethodHead)] = *resIEp
					responseInterceptorMap[string(gwapiv1.HTTPMethodOptions)] = *resIEp
				}
			}
		}
	}
	return requestInterceptorMap, responseInterceptorMap
}

func getInterceptorEndpoint(namespace string, interceptorRef *dpv1alpha3.InterceptorReference,
	gatewayInterceptorServiceMapping map[string]dpv1alpha1.InterceptorService, gatewayBackendMapping map[string]*dpv1alpha2.ResolvedBackend, isReq bool) *model.InterceptEndpoint {
	interceptor := gatewayInterceptorServiceMapping[types.NamespacedName{
		Namespace: namespace,
		Name:      interceptorRef.Name}.String()].Spec
	endpoints := model.GetEndpoints(types.NamespacedName{Namespace: namespace, Name: interceptor.BackendRef.Name},
		gatewayBackendMapping)
	var clusterName string
	if isReq {
		clusterName = constants.GlobalRequestInterceptorClusterName
	} else {
		clusterName = constants.GlobalResponseInterceptorClusterName
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
