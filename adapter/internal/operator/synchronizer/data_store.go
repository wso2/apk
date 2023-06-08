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
	"reflect"
	"sync"

	"github.com/wso2/apk/adapter/internal/loggers"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// OperatorDataStore holds the APIStore and API, HttpRoute mappings
type OperatorDataStore struct {
	apiStore     map[types.NamespacedName]*APIState
	gatewayStore map[types.NamespacedName]*GatewayState
	mu           sync.Mutex
}

// CreateNewOperatorDataStore creates a new OperatorDataStore.
func CreateNewOperatorDataStore() *OperatorDataStore {
	return &OperatorDataStore{
		apiStore:     map[types.NamespacedName]*APIState{},
		gatewayStore: map[types.NamespacedName]*GatewayState{},
	}
}

// AddAPIState stores a new API in the OperatorDataStore.
func (ods *OperatorDataStore) AddAPIState(api dpv1alpha1.API, prodHTTPRouteState *HTTPRouteState,
	sandHTTPRouteState *HTTPRouteState, apiDefinition []byte) APIState {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	apiNamespacedName := utils.NamespacedName(&api)
	ods.apiStore[apiNamespacedName] = &APIState{
		APIDefinition:     &api,
		ProdHTTPRoute:     prodHTTPRouteState,
		SandHTTPRoute:     sandHTTPRouteState,
		APIDefinitionFile: apiDefinition,
	}
	return *ods.apiStore[apiNamespacedName]
}

// UpdateAPIState update/create the APIState on ref updates
func (ods *OperatorDataStore) UpdateAPIState(apiDef *dpv1alpha1.API, prodHTTPRoute *HTTPRouteState,
	sandHTTPRoute *HTTPRouteState, apiDefinitionFile []byte) (APIState, []string, bool) {
	_, found := ods.apiStore[utils.NamespacedName(apiDef)]
	if !found {
		loggers.LoggerAPKOperator.Infof("Adding new apistate as API : %s has not found in memory datastore.", apiDef.Spec.APIDisplayName)
		apiState := ods.AddAPIState(*apiDef, prodHTTPRoute, sandHTTPRoute, apiDefinitionFile)
		return apiState, []string{"API"}, true
	}
	return ods.processAPIState(apiDef, prodHTTPRoute, sandHTTPRoute)
}

// processAPIState process and update the APIState on ref updates
func (ods *OperatorDataStore) processAPIState(apiDef *dpv1alpha1.API, prodHTTPRoute *HTTPRouteState,
	sandHTTPRoute *HTTPRouteState) (APIState, []string, bool) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	var updated bool
	events := []string{}
	cachedAPI := ods.apiStore[utils.NamespacedName(apiDef)]

	if apiDef.Generation > cachedAPI.APIDefinition.Generation {
		cachedAPI.APIDefinition = apiDef
		updated = true
		events = append(events, "API Definition")
	}
	if prodHTTPRoute != nil {
		if routeEvents, routesUpdated := updateHTTPRoute(prodHTTPRoute, cachedAPI.ProdHTTPRoute, "Production"); routesUpdated {
			updated = true
			events = append(events, routeEvents...)
		}
	} else {
		cachedAPI.ProdHTTPRoute = nil
	}
	if sandHTTPRoute != nil {
		if routeEvents, routesUpdated := updateHTTPRoute(sandHTTPRoute, cachedAPI.SandHTTPRoute, "Sandbox"); routesUpdated {
			updated = true
			events = append(events, routeEvents...)
		}
	} else {
		cachedAPI.SandHTTPRoute = nil
	}

	return *cachedAPI, events, updated
}

// UpdateAPIState update the APIState on ref updates
func updateHTTPRoute(httpRoute *HTTPRouteState, cachedHTTPRoute *HTTPRouteState, endpointType string) ([]string, bool) {
	var updated bool
	events := []string{}
	if httpRoute.HTTPRoute.UID != cachedHTTPRoute.HTTPRoute.UID ||
		httpRoute.HTTPRoute.Generation > cachedHTTPRoute.HTTPRoute.Generation {
		cachedHTTPRoute.HTTPRoute = httpRoute.HTTPRoute
		updated = true
		events = append(events, endpointType+" Endpoint")
	}
	if len(httpRoute.Authentications) != len(cachedHTTPRoute.Authentications) {
		cachedHTTPRoute.Authentications = httpRoute.Authentications
		updated = true
		events = append(events, endpointType+" Endpoint Authentications")
	} else {
		for key, auth := range httpRoute.Authentications {
			if existingAuth, found := cachedHTTPRoute.Authentications[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedHTTPRoute.Authentications = httpRoute.Authentications
					updated = true
					events = append(events, endpointType+" Endpoint Authentications")
					break
				}
			} else {
				cachedHTTPRoute.Authentications = httpRoute.Authentications
				updated = true
				events = append(events, endpointType+" Endpoint Authentications")
				break
			}
		}
	}
	if len(httpRoute.ResourceAuthentications) != len(cachedHTTPRoute.ResourceAuthentications) {
		cachedHTTPRoute.ResourceAuthentications = httpRoute.ResourceAuthentications
		updated = true
		events = append(events, endpointType+" Endpoint Resource Authentications")
	} else {
		for key, auth := range httpRoute.ResourceAuthentications {
			if existingAuth, found := cachedHTTPRoute.ResourceAuthentications[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedHTTPRoute.ResourceAuthentications = httpRoute.ResourceAuthentications
					updated = true
					events = append(events, endpointType+" Endpoint Resource Authentications")
					break
				}
			} else {
				cachedHTTPRoute.ResourceAuthentications = httpRoute.ResourceAuthentications
				updated = true
				events = append(events, endpointType+" Endpoint Resource Authentications")
				break
			}
		}
	}

	if len(httpRoute.APIPolicies) != len(cachedHTTPRoute.APIPolicies) {
		cachedHTTPRoute.APIPolicies = httpRoute.APIPolicies
		updated = true
		events = append(events, endpointType+" Endpoint APIPolicies")
	} else {
		for key, auth := range httpRoute.APIPolicies {
			if existingAuth, found := cachedHTTPRoute.APIPolicies[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedHTTPRoute.APIPolicies = httpRoute.APIPolicies
					updated = true
					events = append(events, endpointType+" Endpoint APIPolicies")
					break
				}
			} else {
				cachedHTTPRoute.APIPolicies = httpRoute.APIPolicies
				updated = true
				events = append(events, endpointType+" Endpoint APIPolicies")
				break
			}
		}
	}
	if len(httpRoute.ResourceAPIPolicies) != len(cachedHTTPRoute.ResourceAPIPolicies) {
		cachedHTTPRoute.ResourceAPIPolicies = httpRoute.ResourceAPIPolicies
		updated = true
		events = append(events, endpointType+" Endpoint Resource APIPolicies")
	} else {
		for key, auth := range httpRoute.ResourceAPIPolicies {
			if existingAuth, found := cachedHTTPRoute.ResourceAPIPolicies[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedHTTPRoute.ResourceAPIPolicies = httpRoute.ResourceAPIPolicies
					updated = true
					events = append(events, endpointType+" Endpoint Resource APIPolicies")
					break
				}
			} else {
				cachedHTTPRoute.ResourceAPIPolicies = httpRoute.ResourceAPIPolicies
				updated = true
				events = append(events, endpointType+" Endpoint Resource APIPolicies")
				break
			}
		}
	}

	if len(httpRoute.RateLimitPolicies) != len(cachedHTTPRoute.RateLimitPolicies) {
		cachedHTTPRoute.RateLimitPolicies = httpRoute.RateLimitPolicies
		updated = true
		events = append(events, endpointType+" Endpoint RateLimitPolicies")
	} else {
		for key, rateLimitPolicy := range httpRoute.RateLimitPolicies {
			if existingRateLimitPolicy, found := cachedHTTPRoute.RateLimitPolicies[key]; found {
				if rateLimitPolicy.UID != existingRateLimitPolicy.UID || rateLimitPolicy.Generation > existingRateLimitPolicy.Generation {
					cachedHTTPRoute.RateLimitPolicies = httpRoute.RateLimitPolicies
					updated = true
					events = append(events, endpointType+" Endpoint RateLimitPolicies")
					break
				}
			} else {
				cachedHTTPRoute.RateLimitPolicies = httpRoute.RateLimitPolicies
				updated = true
				events = append(events, endpointType+" Endpoint RateLimitPolicies")
				break
			}
		}
	}
	if len(httpRoute.ResourceRateLimitPolicies) != len(cachedHTTPRoute.ResourceRateLimitPolicies) {
		cachedHTTPRoute.ResourceRateLimitPolicies = httpRoute.ResourceRateLimitPolicies
		updated = true
		events = append(events, endpointType+" Endpoint Resource RateLimitPolicies")
	} else {
		for key, rateLimitPolicy := range httpRoute.ResourceRateLimitPolicies {
			if existingRateLimitPolicy, found := cachedHTTPRoute.ResourceRateLimitPolicies[key]; found {
				if rateLimitPolicy.UID != existingRateLimitPolicy.UID || rateLimitPolicy.Generation > existingRateLimitPolicy.Generation {
					cachedHTTPRoute.ResourceRateLimitPolicies = httpRoute.ResourceRateLimitPolicies
					updated = true
					events = append(events, endpointType+" Endpoint Resource RateLimitPolicies")
					break
				}
			} else {
				cachedHTTPRoute.ResourceRateLimitPolicies = httpRoute.ResourceRateLimitPolicies
				updated = true
				events = append(events, endpointType+" Endpoint Resource RateLimitPolicies")
				break
			}
		}
	}

	if len(httpRoute.Scopes) != len(cachedHTTPRoute.Scopes) {
		cachedHTTPRoute.Scopes = httpRoute.Scopes
		updated = true
		events = append(events, endpointType+" Endpoint Resource Scopes")
	} else {
		for key, scope := range httpRoute.Scopes {
			if existingScope, found := cachedHTTPRoute.Scopes[key]; found {
				if scope.UID != existingScope.UID || scope.Generation > existingScope.Generation {
					cachedHTTPRoute.Scopes = httpRoute.Scopes
					updated = true
					events = append(events, endpointType+" Endpoint Resource Scopes")
					break
				}
			} else {
				cachedHTTPRoute.Scopes = httpRoute.Scopes
				updated = true
				events = append(events, endpointType+" Endpoint Resource Scopes")
				break
			}
		}
	}

	if len(httpRoute.InterceptorServiceMapping) != len(cachedHTTPRoute.InterceptorServiceMapping) {
		cachedHTTPRoute.InterceptorServiceMapping = httpRoute.InterceptorServiceMapping
		updated = true
		events = append(events, endpointType+" Interceptor Service")
	} else {
		for key, interceptService := range httpRoute.InterceptorServiceMapping {
			if existingInterceptService, found := cachedHTTPRoute.InterceptorServiceMapping[key]; found {
				if interceptService.UID != existingInterceptService.UID || interceptService.Generation > existingInterceptService.Generation {
					cachedHTTPRoute.InterceptorServiceMapping = httpRoute.InterceptorServiceMapping
					updated = true
					events = append(events, endpointType+" Interceptor Service")
					break
				}
			} else {
				cachedHTTPRoute.InterceptorServiceMapping = httpRoute.InterceptorServiceMapping
				updated = true
				events = append(events, endpointType+" Interceptor Service")
				break
			}
		}
	}

	if !reflect.DeepEqual(cachedHTTPRoute.BackendMapping, httpRoute.BackendMapping) {
		cachedHTTPRoute.BackendMapping = httpRoute.BackendMapping
		updated = true
		events = append(events, endpointType+" Backend Properties")
	}
	return events, updated
}

// GetCachedAPI get cached apistate
func (ods *OperatorDataStore) GetCachedAPI(apiName types.NamespacedName) (APIState, bool) {
	if cachedAPI, found := ods.apiStore[apiName]; found {
		return *cachedAPI, true
	}
	return APIState{}, false
}

// DeleteCachedAPI delete from apistate cache
func (ods *OperatorDataStore) DeleteCachedAPI(apiName types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	delete(ods.apiStore, apiName)
}

// AddGatewayState stores a new Gateway in the OperatorDataStore.
func (ods *OperatorDataStore) AddGatewayState(gateway gwapiv1b1.Gateway,
	gatewayStateData *GatewayStateData) GatewayState {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	gatewayNamespacedName := utils.NamespacedName(&gateway)
	ods.gatewayStore[gatewayNamespacedName] = &GatewayState{
		GatewayDefinition: &gateway,
		GatewayStateData:  gatewayStateData,
	}
	return *ods.gatewayStore[gatewayNamespacedName]
}

// UpdateGatewayState update/create the GatewayState on ref updates
func (ods *OperatorDataStore) UpdateGatewayState(gatewayDef *gwapiv1b1.Gateway,
	gatewayStateData *GatewayStateData) (GatewayState, []string, bool) {
	_, found := ods.gatewayStore[utils.NamespacedName(gatewayDef)]
	if !found {
		loggers.LoggerAPKOperator.Infof("Adding new gatewaystate as Gateway : %s has not found in memory datastore.", gatewayDef.Name)
		gatewayState := ods.AddGatewayState(*gatewayDef, gatewayStateData)
		return gatewayState, []string{"GATEWAY"}, true
	}
	return ods.processGatewayState(gatewayDef, gatewayStateData.GatewayCustomRateLimitPolicies)
}

// processGatewayState process and update the GatewayState on ref updates
func (ods *OperatorDataStore) processGatewayState(gatewayDef *gwapiv1b1.Gateway,
	customRateLimitPolicies []*dpv1alpha1.RateLimitPolicy) (GatewayState, []string, bool) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	var updated bool
	events := []string{}
	cachedGateway := ods.gatewayStore[utils.NamespacedName(gatewayDef)]

	if gatewayDef.Generation > cachedGateway.GatewayDefinition.Generation {
		cachedGateway.GatewayDefinition = gatewayDef
		updated = true
		events = append(events, "Gateway Definition")
	}

	if !reflect.DeepEqual(cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies, customRateLimitPolicies) {
		cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies = customRateLimitPolicies
		updated = true
		events = append(events, "Gateway Custom RateLimit Policies")
	}

	return *cachedGateway, events, updated
}

// GetCachedGateway get cached gatewaystate
func (ods *OperatorDataStore) GetCachedGateway(gatewayName types.NamespacedName) (GatewayState, bool) {
	if cachedGateway, found := ods.gatewayStore[gatewayName]; found {
		return *cachedGateway, true
	}
	return GatewayState{}, false
}

// IsGatewayAvailable get cached gatewaystate
func (ods *OperatorDataStore) IsGatewayAvailable(gatewayName types.NamespacedName) bool {
	_, found := ods.gatewayStore[gatewayName]
	return found
}

// DeleteCachedGateway delete from gatewaystate cache
func (ods *OperatorDataStore) DeleteCachedGateway(gatewayName types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	delete(ods.gatewayStore, gatewayName)
}
