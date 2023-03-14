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
	dpv1alpha1 "github.com/wso2/apk/adapter/pkg/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/pkg/operator/utils"
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
	sandHTTPRouteState *HTTPRouteState) APIState {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	apiNamespacedName := utils.NamespacedName(&api)
	ods.apiStore[apiNamespacedName] = &APIState{
		APIDefinition: &api,
		ProdHTTPRoute: prodHTTPRouteState,
		SandHTTPRoute: sandHTTPRouteState,
	}
	return *ods.apiStore[apiNamespacedName]
}

// UpdateAPIState update/create the APIState on ref updates
func (ods *OperatorDataStore) UpdateAPIState(apiDef *dpv1alpha1.API, prodHTTPRoute *HTTPRouteState,
	sandHTTPRoute *HTTPRouteState) (APIState, []string, bool) {
	_, found := ods.apiStore[utils.NamespacedName(apiDef)]
	if !found {
		loggers.LoggerAPKOperator.Infof("Adding new apistate as API : %s has not found in memory datastore.", apiDef.Spec.APIDisplayName)
		apiState := ods.AddAPIState(*apiDef, prodHTTPRoute, sandHTTPRoute)
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
	}
	if sandHTTPRoute != nil {
		if routeEvents, routesUpdated := updateHTTPRoute(sandHTTPRoute, cachedAPI.SandHTTPRoute, "Sandbox"); routesUpdated {
			updated = true
			events = append(events, routeEvents...)
		}
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
func (ods *OperatorDataStore) AddGatewayState(gateway gwapiv1b1.Gateway) GatewayState {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	gatewayNamespacedName := utils.NamespacedName(&gateway)
	ods.gatewayStore[gatewayNamespacedName] = &GatewayState{
		GatewayDefinition: &gateway,
	}
	return *ods.gatewayStore[gatewayNamespacedName]
}

// UpdateGatewayState update/create the GatewayState on ref updates
func (ods *OperatorDataStore) UpdateGatewayState(gatewayDef *gwapiv1b1.Gateway) (GatewayState, []string, bool) {
	_, found := ods.gatewayStore[utils.NamespacedName(gatewayDef)]
	if !found {
		loggers.LoggerAPKOperator.Infof("Adding new gatewaystate as Gateway : %s has not found in memory datastore.", gatewayDef.Name)
		gatewayState := ods.AddGatewayState(*gatewayDef)
		return gatewayState, []string{"GATEWAY"}, true
	}
	return ods.processGatewayState(gatewayDef)
}

// processGatewayState process and update the GatewayState on ref updates
func (ods *OperatorDataStore) processGatewayState(gatewayDef *gwapiv1b1.Gateway) (GatewayState, []string, bool) {
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

	return *cachedGateway, events, updated
}

// GetCachedGateway get cached gatewaystate
func (ods *OperatorDataStore) GetCachedGateway(gatewayName types.NamespacedName) (GatewayState, bool) {
	if cachedGateway, found := ods.gatewayStore[gatewayName]; found {
		return *cachedGateway, true
	}
	return GatewayState{}, false
}

// DeleteCachedGateway delete from gatewaystate cache
func (ods *OperatorDataStore) DeleteCachedGateway(gatewayName types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	delete(ods.gatewayStore, gatewayName)
}
