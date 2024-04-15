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
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/dataholder"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/logging"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// extract APIDetails from the GQLRoute
func updateInternalMapsFromGRPCRoute(apiState APIState, grpcRoute *GRPCRouteState, envType string) (*model.AdapterInternalAPI, map[string]struct{}, error) {
	adapterInternalAPI, err := generateGRPCAdapterInternalAPI(apiState, grpcRoute, envType)
	if err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2632, logging.MAJOR, "Error generating AdapterInternalAPI for GRPCRoute: %v. %v", grpcRoute.GRPCRouteCombined.Name, err))
		return nil, nil, err
	}

	vHosts := getVhostsForGRPCAPI(grpcRoute.GRPCRouteCombined)
	labels := getLabelsForGRPCAPI(grpcRoute.GRPCRouteCombined)
	listeners, relativeSectionNames := getListenersForGRPCAPI(grpcRoute.GRPCRouteCombined, adapterInternalAPI.UUID)
	// We dont have a use case where a perticular API's two different grpc routes refer to two different gateway. Hence get the first listener name for the list for processing.
	if len(listeners) == 0 || len(relativeSectionNames) == 0 {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2633, logging.MINOR, "Failed to find a matching listener for grpc route: %v. ",
			grpcRoute.GRPCRouteCombined.Name))
		return nil, nil, errors.New("failed to find matching listener name for the provided grpc route")
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

// generateGRPCAdapterInternalAPI this will populate a AdapterInternalAPI representation for an GRPCRoute
func generateGRPCAdapterInternalAPI(apiState APIState, grpcRoute *GRPCRouteState, envType string) (*model.AdapterInternalAPI, error) {
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
	adapterInternalAPI.SetXWso2RequestBodyPass(true)

	resourceParams := model.ResourceParams{
		AuthSchemes:               apiState.Authentications,
		ResourceAuthSchemes:       apiState.ResourceAuthentications,
		BackendMapping:            grpcRoute.BackendMapping,
		APIPolicies:               apiState.APIPolicies,
		ResourceAPIPolicies:       apiState.ResourceAPIPolicies,
		ResourceScopes:            grpcRoute.Scopes,
		InterceptorServiceMapping: apiState.InterceptorServiceMapping,
		BackendJWTMapping:         apiState.BackendJWTMapping,
		RateLimitPolicies:         apiState.RateLimitPolicies,
		ResourceRateLimitPolicies: apiState.ResourceRateLimitPolicies,
	}
	if err := adapterInternalAPI.SetInfoGRPCRouteCR(grpcRoute.GRPCRouteCombined, resourceParams); err != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2631, logging.MAJOR, "Error setting GRPCRoute CR info to adapterInternalAPI. %v", err))
		return nil, err
	}
	if apiState.MutualSSL != nil && apiState.MutualSSL.Required != "" && !adapterInternalAPI.GetDisableAuthentications() {
		adapterInternalAPI.SetDisableMtls(apiState.MutualSSL.Disabled)
		adapterInternalAPI.SetMutualSSL(apiState.MutualSSL.Required)
		adapterInternalAPI.SetClientCerts(apiState.APIDefinition.Name, apiState.MutualSSL.ClientCertificates)
	} else {
		adapterInternalAPI.SetDisableMtls(true)
	}

	return &adapterInternalAPI, nil

}

// getVhostForAPI returns the vHosts related to an API.
func getVhostsForGRPCAPI(grpcRoute *gwapiv1a2.GRPCRoute) map[string]struct{} {
	vHosts := make(map[string]struct{})
	for _, hostName := range grpcRoute.Spec.Hostnames {
		vHosts[string(hostName)] = struct{}{}
	}
	return vHosts
}

// getLabelsForAPI returns the labels related to an API.
func getLabelsForGRPCAPI(grpcRoute *gwapiv1a2.GRPCRoute) map[string]struct{} {
	labels := make(map[string]struct{})
	for _, parentRef := range grpcRoute.Spec.ParentRefs {
		err := xds.SanitizeGateway(string(parentRef.Name), false)
		if err != nil {
			loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2653, logging.CRITICAL, "Gateway Label is invalid: %s", string(parentRef.Name)))
		} else {
			labels[string(parentRef.Name)] = struct{}{}
		}
	}
	return labels
}

// getListenersForGRPCAPI returns the listeners related to an API.
func getListenersForGRPCAPI(grpcRoute *gwapiv1a2.GRPCRoute, apiUUID string) ([]string, []string) {
	var listeners []string
	var sectionNames []string
	for _, parentRef := range grpcRoute.Spec.ParentRefs {
		namespace := grpcRoute.GetNamespace()
		if parentRef.Namespace != nil && *parentRef.Namespace != "" {
			namespace = string(*parentRef.Namespace)
		}
		gateway, found := dataholder.GetGatewayMap()[types.NamespacedName{
			Namespace: namespace,
			Name:      string(parentRef.Name),
		}.String()]
		if found {
			// find the matching listener
			matchedListener, listenerFound := common.FindElement(gateway.Spec.Listeners, func(listener gwapiv1b1.Listener) bool {
				if string(listener.Name) == string(*parentRef.SectionName) {
					return true
				}
				return false
			})
			if listenerFound {
				sectionNames = append(sectionNames, string(matchedListener.Name))
				listeners = append(listeners, common.GetEnvoyListenerName(string(matchedListener.Protocol), uint32(matchedListener.Port)))
				continue
			}
		}
		loggers.LoggerAPKOperator.Errorf("Failed to find matching listeners for the grpcroute: %+v", grpcRoute.Name)
	}
	return listeners, sectionNames
}

func deleteGRPCAPIFromEnv(grpcRoute *gwapiv1a2.GRPCRoute, apiState APIState) error {
	labels := getLabelsForGRPCAPI(grpcRoute)
	uuid := string(apiState.APIDefinition.ObjectMeta.UID)
	return xds.DeleteAPI(uuid, labels)
}

// undeployGRPCAPIInGateway undeploys the related API in CREATE and UPDATE events.
func undeployGRPCAPIInGateway(apiState APIState) error {
	var err error
	if apiState.ProdGRPCRoute != nil {
		err = deleteGRPCAPIFromEnv(apiState.ProdGRPCRoute.GRPCRouteCombined, apiState)
	}
	if err != nil {
		loggers.LoggerXds.ErrorC(logging.PrintError(logging.Error2630, logging.MAJOR, "Error undeploying prod grpcRoute of API : %v in Organization %v from environments."+
			" Hence not checking on deleting the sand grpcRoute of the API", string(apiState.APIDefinition.ObjectMeta.UID), apiState.APIDefinition.Spec.Organization))
		return err
	}
	if apiState.SandGRPCRoute != nil {
		err = deleteGRPCAPIFromEnv(apiState.SandGRPCRoute.GRPCRouteCombined, apiState)
	}
	return err
}
