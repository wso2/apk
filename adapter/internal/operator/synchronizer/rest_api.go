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
	"errors"
	"strings"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/dataholder"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/logging"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// undeployAPIInGateway undeploys the related API in CREATE and UPDATE events.
func undeployRestAPIInGateway(apiState APIState) error {
	var err error
	if apiState.ProdHTTPRoute != nil {
		err = deleteAPIFromEnv(apiState.ProdHTTPRoute.HTTPRouteCombined, apiState)
	}
	if err != nil {
		loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error2630, logging.MAJOR, "Error undeploying prod httpRoute of API : %v in Organization %v from environments %v."+
			" Hence not checking on deleting the sand httpRoute of the API", string(apiState.APIDefinition.ObjectMeta.UID), apiState.APIDefinition.Spec.Organization,
			getGatewayNameForAPI(apiState.ProdHTTPRoute.HTTPRouteCombined)))
		return err
	}
	if apiState.SandHTTPRoute != nil {
		err = deleteAPIFromEnv(apiState.SandHTTPRoute.HTTPRouteCombined, apiState)
	}
	return err
}

// UpdateInternalMapsFromHTTPRoute extracts the API details from the HTTPRoute.
func UpdateInternalMapsFromHTTPRoute(apiState APIState, httpRoute *HTTPRouteState, envType string) (*model.AdapterInternalAPI, map[string]struct{}, error) {
	adapterInternalAPI, err := generateAdapterInternalAPI(apiState, httpRoute, envType)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2632, logging.MAJOR, "Error generating AdapterInternalAPI for HTTPRoute: %v. %v", httpRoute.HTTPRouteCombined.Name, err))
		return nil, nil, err
	}

	vHosts := getVhostsForAPI(httpRoute.HTTPRouteCombined)
	labels := getGatewayNameForAPI(httpRoute.HTTPRouteCombined)
	listeners, relativeSectionNames := getListenersForAPI(httpRoute.HTTPRouteCombined)
	// We dont have a use case where a perticular API's two different http routes refer to two different gateway. Hence get the first listener name for the list for processing.
	if len(listeners) == 0 || len(relativeSectionNames) == 0 {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2633, logging.MINOR, "Failed to find a matching listener for http route: %v. ",
			httpRoute.HTTPRouteCombined.Name))
		return nil, nil, errors.New("failed to find matching listener name for the provided http route")
	}

	listenerName := listeners[0]
	sectionName := relativeSectionNames[0]
	if len(listeners) > 0 {
		if err := xds.PopulateInternalMaps(adapterInternalAPI, labels, vHosts, sectionName, listenerName); err != nil {
			return nil, nil, err
		}
	}
	return adapterInternalAPI, labels, nil
}

// generateAdapterInternalAPI this will populate a AdapterInternalAPI representation for an HTTPRoute
func generateAdapterInternalAPI(apiState APIState, httpRouteState *HTTPRouteState, envType string) (*model.AdapterInternalAPI, error) {
	var adapterInternalAPI model.AdapterInternalAPI
	adapterInternalAPI.SetIsDefaultVersion(apiState.APIDefinition.Spec.IsDefaultVersion)
	adapterInternalAPI.SetInfoAPICR(*apiState.APIDefinition)
	adapterInternalAPI.SetAPIDefinitionFile(apiState.APIDefinitionFile)
	adapterInternalAPI.SetAPIDefinitionEndpoint(apiState.APIDefinition.Spec.DefinitionPath)
	adapterInternalAPI.SetSubscriptionValidation(apiState.SubscriptionValidation)

	adapterInternalAPI.EnvType = envType

	environment := apiState.APIDefinition.Spec.Environment
	if environment == "" {
		conf := config.ReadConfigs()
		environment = conf.Adapter.Environment
	}
	adapterInternalAPI.SetEnvironment(environment)

	resourceParams := model.ResourceParams{
		AuthSchemes:               apiState.Authentications,
		ResourceAuthSchemes:       apiState.ResourceAuthentications,
		BackendMapping:            httpRouteState.BackendMapping,
		APIPolicies:               apiState.APIPolicies,
		ResourceAPIPolicies:       apiState.ResourceAPIPolicies,
		ResourceScopes:            httpRouteState.Scopes,
		InterceptorServiceMapping: apiState.InterceptorServiceMapping,
		BackendJWTMapping:         apiState.BackendJWTMapping,
		RateLimitPolicies:         apiState.RateLimitPolicies,
		ResourceRateLimitPolicies: apiState.ResourceRateLimitPolicies,
	}
	sendToEnforcer := false
	if config.ReadConfigs().Analytics.Enabled {
		loggers.LoggerAPKOperator.Infof("Analytics is enabled for API: %v", apiState.APIDefinition.Name)
		sendToEnforcer = true
	}
	if  apiState.BackendJWTMapping != nil && len(apiState.BackendJWTMapping) > 0 {
		sendToEnforcer = true
	}
	if apiState.APIDefinition.Spec.APIType == "GRPC" {
		sendToEnforcer = true
	}
	if apiState.MutualSSL != nil {
		sendToEnforcer = !apiState.MutualSSL.Disabled
	}
	if (apiState.AIProvider != nil && apiState.AIProvider.Name != "") || apiState.SubscriptionValidation {
		sendToEnforcer = true
	}
	if err := adapterInternalAPI.SetInfoHTTPRouteCR(httpRouteState.HTTPRouteCombined, resourceParams, httpRouteState.RuleIdxToAiRatelimitPolicyMapping, apiState.AIProvider.Spec.RateLimitFields.PromptTokens.In, sendToEnforcer); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2631, logging.MAJOR, "Error setting HttpRoute CR info to adapterInternalAPI. %v", err))
		return nil, err
	}

	if apiState.MutualSSL != nil && apiState.MutualSSL.Required != "" && !adapterInternalAPI.GetDisableAuthentications() {
		adapterInternalAPI.SetDisableMtls(apiState.MutualSSL.Disabled)
		adapterInternalAPI.SetMutualSSL(apiState.MutualSSL.Required)
		adapterInternalAPI.SetClientCerts(apiState.APIDefinition.Name, apiState.MutualSSL.ClientCertificates)
	} else {
		adapterInternalAPI.SetDisableMtls(true)
	}

	if err := adapterInternalAPI.Validate(); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2632, logging.MAJOR, "Error validating adapterInternalAPI intermediate representation. %v", err))
		return nil, err
	}

	if apiState.AIProvider != nil && apiState.AIProvider.Name != "" {
		loggers.LoggerAPKOperator.Debugf("AI Provider: %+v", *apiState.AIProvider)
		adapterInternalAPI.SetAIProvider(*apiState.AIProvider)
	}

	if apiState.ResolvedModelBasedRoundRobin != nil && (len(apiState.ResolvedModelBasedRoundRobin.ProductionModels) > 0 || len(apiState.ResolvedModelBasedRoundRobin.SandboxModels) > 0) {
		loggers.LoggerAPKOperator.Infof("API State Model Based Round Robin: %+v", *apiState.ResolvedModelBasedRoundRobin)
		resolvedModelBasedRoundRobin := *apiState.ResolvedModelBasedRoundRobin
		var productionModels []model.InternalModelWeight
		var sandboxModels []model.InternalModelWeight
		for _, aiModel := range resolvedModelBasedRoundRobin.ProductionModels {
			endpoints := model.GetEndpointsByResolvedBackend(aiModel.ResolvedBackend)
			endpointCluster := model.EndpointCluster{
				Endpoints: endpoints,
			}
			vhost := ""
			for _, hostName := range httpRouteState.HTTPRouteCombined.Spec.Hostnames {
				vhost = string(hostName)
			}
			clusternName := getClusterName(endpointCluster.EndpointPrefix, adapterInternalAPI.GetOrganizationID(), vhost, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), endpoints[0].Host+endpoints[0].Basepath)
			productionModels = append(productionModels, model.InternalModelWeight{
				Model:               aiModel.Model,
				EndpointClusterName: clusternName,
				Weight:              aiModel.Weight,
			})
		}
		for _, aiModel := range resolvedModelBasedRoundRobin.SandboxModels {
			endpoints := model.GetEndpointsByResolvedBackend(aiModel.ResolvedBackend)
			endpointCluster := model.EndpointCluster{
				Endpoints: endpoints,
			}
			vhost := ""
			for _, hostName := range httpRouteState.HTTPRouteCombined.Spec.Hostnames {
				vhost = string(hostName)
			}
			clusternName := getClusterName(endpointCluster.EndpointPrefix, adapterInternalAPI.GetOrganizationID(), vhost, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), endpoints[0].Host+endpoints[0].Basepath)
			sandboxModels = append(sandboxModels, model.InternalModelWeight{
				Model:               aiModel.Model,
				EndpointClusterName: clusternName,
				Weight:              aiModel.Weight,
			})
		}
		aiModelBasedRoundRobin := model.InternalModelBasedRoundRobin{
			OnQuotaExceedSuspendDuration: resolvedModelBasedRoundRobin.OnQuotaExceedSuspendDuration,
			ProductionModels:             productionModels,
			SandboxModels:                sandboxModels,
		}
		adapterInternalAPI.SetModelBasedRoundRobin(aiModelBasedRoundRobin)
	}

	loggers.LoggerAPKOperator.Infof("AdapterInternalAPI AI Provider: %+v", adapterInternalAPI.GetAIProvider())
	loggers.LoggerAPKOperator.Infof("AdapterInternalAPI Model Based Round Robin: %+v", adapterInternalAPI.GetModelBasedRoundRobin())

	return &adapterInternalAPI, nil
}

// getClusterName returns the cluster name for the API.
func getClusterName(epPrefix string, organizationID string, vHost string, swaggerTitle string, swaggerVersion string,
	hostname string) string {
	if hostname != "" {
		return strings.TrimSpace(organizationID+"_"+epPrefix+"_"+vHost+"_"+strings.Replace(swaggerTitle, " ", "", -1)+swaggerVersion) +
			"_" + strings.Replace(hostname, " ", "", -1) + "0"
	}
	return strings.TrimSpace(organizationID + "_" + epPrefix + "_" + vHost + "_" + strings.Replace(swaggerTitle, " ", "", -1) +
		swaggerVersion)
}

// getVhostForAPI returns the vHosts related to an API.
func getVhostsForAPI(httpRoute *gwapiv1.HTTPRoute) map[string]struct{} {
	vHosts := make(map[string]struct{})
	for _, hostName := range httpRoute.Spec.Hostnames {
		vHosts[string(hostName)] = struct{}{}
	}
	return vHosts
}

// getGatewayNameForAPI returns the labels related to an API.
func getGatewayNameForAPI(httpRoute *gwapiv1.HTTPRoute) map[string]struct{} {
	labels := make(map[string]struct{})
	var err error
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		err = xds.SanitizeGateway(string(parentRef.Name), false)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2653, logging.CRITICAL, "Gateway Label is invalid: %s", string(parentRef.Name)))
		} else {
			labels[string(parentRef.Name)] = struct{}{}
		}
	}
	return labels
}

// getListenersForAPI returns the listeners related to an API.
func getListenersForAPI(httpRoute *gwapiv1.HTTPRoute) ([]string, []string) {
	var listeners []string
	var sectionNames []string
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		namespace := httpRoute.GetNamespace()
		if parentRef.Namespace != nil && *parentRef.Namespace != "" {
			namespace = string(*parentRef.Namespace)
		}
		gateway, found := dataholder.GetGatewayMap()[types.NamespacedName{
			Namespace: namespace,
			Name:      string(parentRef.Name),
		}.String()]
		if found {
			// find the matching listener
			matchedListener, listenerFound := common.FindElement(gateway.Spec.Listeners, func(listener gwapiv1.Listener) bool {
				return string(listener.Name) == string(*parentRef.SectionName)
			})
			if listenerFound {
				sectionNames = append(sectionNames, string(matchedListener.Name))
				listeners = append(listeners, common.GetEnvoyListenerName(string(matchedListener.Protocol), uint32(matchedListener.Port)))
				continue
			}
		}
		loggers.LoggerAPKOperator.Errorf("Failed to find matching listeners for the httproute: %+v", httpRoute.Name)
	}
	return listeners, sectionNames
}

func deleteAPIFromEnv(httpRoute *gwapiv1.HTTPRoute, apiState APIState) error {
	labels := getGatewayNameForAPI(httpRoute)
	uuid := string(apiState.APIDefinition.ObjectMeta.UID)
	return xds.DeleteAPI(uuid, labels)
}
