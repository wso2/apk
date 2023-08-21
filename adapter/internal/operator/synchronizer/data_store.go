/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
func (ods *OperatorDataStore) AddAPIState(apiNamespacedName types.NamespacedName, apiState *APIState) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	ods.apiStore[apiNamespacedName] = apiState
}

// UpdateAPIState update/create the APIState on ref updates
func (ods *OperatorDataStore) UpdateAPIState(apiNamespacedName types.NamespacedName, apiState *APIState) (APIState, []string, bool) {
	_, found := ods.apiStore[apiNamespacedName]
	if !found {
		loggers.LoggerAPKOperator.Infof("Adding new apistate as API : %s has not found in memory datastore.",
			apiState.APIDefinition.Name)
		ods.AddAPIState(apiNamespacedName, apiState)
		return *apiState, []string{"API"}, true
	}
	return ods.processAPIState(apiNamespacedName, apiState)
}

// processAPIState process and update the APIState on ref updates
func (ods *OperatorDataStore) processAPIState(apiNamespacedName types.NamespacedName, apiState *APIState) (APIState, []string, bool) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	var updated bool
	events := []string{}
	cachedAPI := ods.apiStore[apiNamespacedName]

	if apiState.APIDefinition.Generation > cachedAPI.APIDefinition.Generation {
		cachedAPI.APIDefinition = apiState.APIDefinition
		updated = true
		events = append(events, "API Definition")
	}
	if apiState.ProdHTTPRoute != nil {
		if cachedAPI.ProdHTTPRoute == nil {
			cachedAPI.ProdHTTPRoute = apiState.ProdHTTPRoute
			updated = true
			events = append(events, "Production")
		} else if routeEvents, routesUpdated := updateHTTPRoute(apiState.ProdHTTPRoute, cachedAPI.ProdHTTPRoute,
			"Production"); routesUpdated {
			updated = true
			events = append(events, routeEvents...)
		}
	} else {
		cachedAPI.ProdHTTPRoute = nil
	}
	if apiState.SandHTTPRoute != nil {
		if cachedAPI.SandHTTPRoute == nil {
			cachedAPI.SandHTTPRoute = apiState.SandHTTPRoute
			updated = true
			events = append(events, "Sandbox")
		} else if routeEvents, routesUpdated := updateHTTPRoute(apiState.SandHTTPRoute, cachedAPI.SandHTTPRoute, "Sandbox"); routesUpdated {
			updated = true
			events = append(events, routeEvents...)
		}
	} else {
		cachedAPI.SandHTTPRoute = nil
	}
	if len(apiState.Authentications) != len(cachedAPI.Authentications) {
		cachedAPI.Authentications = apiState.Authentications
		updated = true
		events = append(events, "Authentications")
	} else {
		for key, auth := range apiState.Authentications {
			if existingAuth, found := cachedAPI.Authentications[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedAPI.Authentications = apiState.Authentications
					updated = true
					events = append(events, "Authentications")
					break
				}
			} else {
				cachedAPI.Authentications = apiState.Authentications
				updated = true
				events = append(events, "Authentications")
				break
			}
		}
	}
	if len(apiState.ResourceAuthentications) != len(cachedAPI.ResourceAuthentications) {
		cachedAPI.ResourceAuthentications = apiState.ResourceAuthentications
		updated = true
		events = append(events, "Resource Authentications")
	} else {
		for key, auth := range apiState.ResourceAuthentications {
			if existingAuth, found := cachedAPI.ResourceAuthentications[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedAPI.ResourceAuthentications = apiState.ResourceAuthentications
					updated = true
					events = append(events, "Resource Authentications")
					break
				}
			} else {
				cachedAPI.ResourceAuthentications = apiState.ResourceAuthentications
				updated = true
				events = append(events, "Resource Authentications")
				break
			}
		}
	}

	if len(apiState.APIPolicies) != len(cachedAPI.APIPolicies) {
		cachedAPI.APIPolicies = apiState.APIPolicies
		updated = true
		events = append(events, "APIPolicies")
	} else {
		for key, auth := range apiState.APIPolicies {
			if existingAuth, found := cachedAPI.APIPolicies[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedAPI.APIPolicies = apiState.APIPolicies
					updated = true
					events = append(events, "APIPolicies")
					break
				}
			} else {
				cachedAPI.APIPolicies = apiState.APIPolicies
				updated = true
				events = append(events, "APIPolicies")
				break
			}
		}
	}
	if len(apiState.ResourceAPIPolicies) != len(cachedAPI.ResourceAPIPolicies) {
		cachedAPI.ResourceAPIPolicies = apiState.ResourceAPIPolicies
		updated = true
		events = append(events, "Resource APIPolicies")
	} else {
		for key, auth := range apiState.ResourceAPIPolicies {
			if existingAuth, found := cachedAPI.ResourceAPIPolicies[key]; found {
				if auth.UID != existingAuth.UID || auth.Generation > existingAuth.Generation {
					cachedAPI.ResourceAPIPolicies = apiState.ResourceAPIPolicies
					updated = true
					events = append(events, "Resource APIPolicies")
					break
				}
			} else {
				cachedAPI.ResourceAPIPolicies = apiState.ResourceAPIPolicies
				updated = true
				events = append(events, "Resource APIPolicies")
				break
			}
		}
	}

	if len(apiState.RateLimitPolicies) != len(cachedAPI.RateLimitPolicies) {
		cachedAPI.RateLimitPolicies = apiState.RateLimitPolicies
		updated = true
		events = append(events, "RateLimitPolicies")
	} else {
		for key, rateLimitPolicy := range apiState.RateLimitPolicies {
			if existingRateLimitPolicy, found := cachedAPI.RateLimitPolicies[key]; found {
				if rateLimitPolicy.UID != existingRateLimitPolicy.UID || rateLimitPolicy.Generation > existingRateLimitPolicy.Generation {
					cachedAPI.RateLimitPolicies = apiState.RateLimitPolicies
					updated = true
					events = append(events, "RateLimitPolicies")
					break
				}
			} else {
				cachedAPI.RateLimitPolicies = apiState.RateLimitPolicies
				updated = true
				events = append(events, "RateLimitPolicies")
				break
			}
		}
	}
	if len(apiState.ResourceRateLimitPolicies) != len(cachedAPI.ResourceRateLimitPolicies) {
		cachedAPI.ResourceRateLimitPolicies = apiState.ResourceRateLimitPolicies
		updated = true
		events = append(events, "Resource RateLimitPolicies")
	} else {
		for key, rateLimitPolicy := range apiState.ResourceRateLimitPolicies {
			if existingRateLimitPolicy, found := cachedAPI.ResourceRateLimitPolicies[key]; found {
				if rateLimitPolicy.UID != existingRateLimitPolicy.UID || rateLimitPolicy.Generation > existingRateLimitPolicy.Generation {
					cachedAPI.ResourceRateLimitPolicies = apiState.ResourceRateLimitPolicies
					updated = true
					events = append(events, "Resource RateLimitPolicies")
					break
				}
			} else {
				cachedAPI.ResourceRateLimitPolicies = apiState.ResourceRateLimitPolicies
				updated = true
				events = append(events, "Resource RateLimitPolicies")
				break
			}
		}
	}
	if len(apiState.InterceptorServiceMapping) != len(cachedAPI.InterceptorServiceMapping) {
		cachedAPI.InterceptorServiceMapping = apiState.InterceptorServiceMapping
		updated = true
		events = append(events, "Interceptor Service")
	} else {
		for key, interceptService := range apiState.InterceptorServiceMapping {
			if existingInterceptService, found := cachedAPI.InterceptorServiceMapping[key]; found {
				if interceptService.UID != existingInterceptService.UID || interceptService.Generation > existingInterceptService.Generation {
					cachedAPI.InterceptorServiceMapping = apiState.InterceptorServiceMapping
					updated = true
					events = append(events, "Interceptor Service")
					break
				}
			} else {
				cachedAPI.InterceptorServiceMapping = apiState.InterceptorServiceMapping
				updated = true
				events = append(events, "Interceptor Service")
				break
			}
		}
	}

	if len(apiState.BackendJWTMapping) != len(cachedAPI.BackendJWTMapping) {
		cachedAPI.BackendJWTMapping = apiState.BackendJWTMapping
		updated = true
		events = append(events, "Backend JWT")
	} else {
		for key, backendJWT := range apiState.BackendJWTMapping {
			if existingBackendJWT, found := cachedAPI.BackendJWTMapping[key]; found {
				if backendJWT.UID != existingBackendJWT.UID || backendJWT.Generation > existingBackendJWT.Generation {
					cachedAPI.BackendJWTMapping = apiState.BackendJWTMapping
					updated = true
					events = append(events, "Backend JWT")
					break
				}
			} else {
				cachedAPI.BackendJWTMapping = apiState.BackendJWTMapping
				updated = true
				events = append(events, "Backend JWT")
				break
			}
		}
	}

	return *cachedAPI, events, updated
}

// UpdateAPIState update the APIState on ref updates
func updateHTTPRoute(httpRoute *HTTPRouteState, cachedHTTPRoute *HTTPRouteState, endpointType string) ([]string, bool) {
	var updated bool
	events := []string{}
	if cachedHTTPRoute.HTTPRouteCombined == nil || !isEqualHTTPRoutes(cachedHTTPRoute.HTTPRoutePartitions, httpRoute.HTTPRoutePartitions) {
		cachedHTTPRoute.HTTPRouteCombined = httpRoute.HTTPRouteCombined
		cachedHTTPRoute.HTTPRoutePartitions = httpRoute.HTTPRoutePartitions
		updated = true
		events = append(events, endpointType+" Endpoint")
	}

	if len(httpRoute.Scopes) != len(cachedHTTPRoute.Scopes) {
		cachedHTTPRoute.Scopes = httpRoute.Scopes
		updated = true
		events = append(events, "Resource Scopes")
	} else {
		for key, scope := range httpRoute.Scopes {
			if existingScope, found := cachedHTTPRoute.Scopes[key]; found {
				if scope.UID != existingScope.UID || scope.Generation > existingScope.Generation {
					cachedHTTPRoute.Scopes = httpRoute.Scopes
					updated = true
					events = append(events, "Resource Scopes")
					break
				}
			} else {
				cachedHTTPRoute.Scopes = httpRoute.Scopes
				updated = true
				events = append(events, "Resource Scopes")
				break
			}
		}
	}

	if len(httpRoute.BackendMapping) != len(cachedHTTPRoute.BackendMapping) {
		cachedHTTPRoute.BackendMapping = httpRoute.BackendMapping
		updated = true
		events = append(events, endpointType+" Backend Properties")
	} else {
		for key, backend := range httpRoute.BackendMapping {
			if existingBackend, found := cachedHTTPRoute.BackendMapping[key]; found {
				if backend.Backend.UID != existingBackend.Backend.UID || backend.Backend.Generation > existingBackend.Backend.Generation {
					cachedHTTPRoute.BackendMapping = httpRoute.BackendMapping
					updated = true
					events = append(events, endpointType+" Backend Properties")
					break
				}
			} else {
				cachedHTTPRoute.BackendMapping = httpRoute.BackendMapping
				updated = true
				events = append(events, endpointType+" Backend Properties")
				break
			}
		}
	}
	return events, updated
}

func isEqualHTTPRoutes(cachedHTTPRoutes, newHTTPRoutes map[string]*gwapiv1b1.HTTPRoute) bool {
	for key, cachedHTTPRoute := range cachedHTTPRoutes {
		if newHTTPRoutes[key] == nil {
			return false
		}
		if newHTTPRoutes[key].UID == cachedHTTPRoute.UID &&
			newHTTPRoutes[key].Generation > cachedHTTPRoute.Generation {
			return false
		}
	}
	return true
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
	customRateLimitPolicies map[string]*dpv1alpha1.RateLimitPolicy) (GatewayState, []string, bool) {
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

	if len(customRateLimitPolicies) != len(cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies) {
		cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies = customRateLimitPolicies
		updated = true
		events = append(events, "Gateway Custom RateLimit Policies")
	} else {
		for key, rateLimitPolicy := range customRateLimitPolicies {
			if existingRateLimitPolicy, found := cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies[key]; found {
				if rateLimitPolicy.UID != existingRateLimitPolicy.UID || rateLimitPolicy.Generation > existingRateLimitPolicy.Generation {
					cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies = customRateLimitPolicies
					updated = true
					events = append(events, "Gateway Custom RateLimit Policies")
					break
				}
			} else {
				cachedGateway.GatewayStateData.GatewayCustomRateLimitPolicies = customRateLimitPolicies
				updated = true
				events = append(events, "Gateway Custom RateLimit Policies")
				break
			}
		}
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
