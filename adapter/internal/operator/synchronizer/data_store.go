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
	"sync"

	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"k8s.io/apimachinery/pkg/types"
)

// OperatorDataStore holds the APIStore and API, HttpRoute mappings
type OperatorDataStore struct {
	apiStore map[types.NamespacedName]*APIState
	mu       sync.Mutex
}

// CreateNewOperatorDataStore creates a new OperatorDataStore.
func CreateNewOperatorDataStore() *OperatorDataStore {
	return &OperatorDataStore{
		apiStore: map[types.NamespacedName]*APIState{},
	}
}

// AddNewAPItoODS stores a new API in the OperatorDataStore.
func (ods *OperatorDataStore) AddNewAPItoODS(api dpv1alpha1.API, prodHTTPRouteState *HTTPRouteState,
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

// UpdateAPIState update the APIState on ref updates
func (ods *OperatorDataStore) UpdateAPIState(apiDef *dpv1alpha1.API, prodHTTPRoute *HTTPRouteState,
	sandHTTPRoute *HTTPRouteState) ([]string, APIState, bool) {
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
		if prodHTTPRoute.HTTPRoute.UID != cachedAPI.ProdHTTPRoute.HTTPRoute.UID ||
			prodHTTPRoute.HTTPRoute.Generation > cachedAPI.ProdHTTPRoute.HTTPRoute.Generation {
			cachedAPI.ProdHTTPRoute = prodHTTPRoute
			updated = true
			events = append(events, "Production Endpoint")
		}
	}
	if sandHTTPRoute != nil {
		if sandHTTPRoute.HTTPRoute.UID != cachedAPI.SandHTTPRoute.HTTPRoute.UID ||
			sandHTTPRoute.HTTPRoute.Generation > cachedAPI.SandHTTPRoute.HTTPRoute.Generation {
			cachedAPI.SandHTTPRoute = sandHTTPRoute
			updated = true
			events = append(events, "Sandbox Endpoint")
		}
	}

	return events, *cachedAPI, updated
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
