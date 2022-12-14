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
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
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
func (ods *OperatorDataStore) AddNewAPItoODS(api dpv1alpha1.API, prodHTTPRoute *gwapiv1b1.HTTPRoute,
	sandHTTPRoute *gwapiv1b1.HTTPRoute, authentications map[types.NamespacedName]*dpv1alpha1.Authentication) APIState {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	apiNamespacedName := utils.NamespacedName(&api)
	ods.apiStore[apiNamespacedName] = &APIState{
		APIDefinition:   &api,
		ProdHTTPRoute:   prodHTTPRoute,
		SandHTTPRoute:   sandHTTPRoute,
		Authentications: authentications,
	}
	return *ods.apiStore[apiNamespacedName]
}

// UpdateAPIState update the APIState on ref updates
func (ods *OperatorDataStore) UpdateAPIState(apiDef *dpv1alpha1.API, prodHTTPRoute *gwapiv1b1.HTTPRoute,
	sandHTTPRoute *gwapiv1b1.HTTPRoute, authentications map[types.NamespacedName]*dpv1alpha1.Authentication) ([]string, bool) {
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
	//TODO(amali) remove extensions map related to old routes
	if prodHTTPRoute != nil && (prodHTTPRoute.UID != cachedAPI.ProdHTTPRoute.UID ||
		prodHTTPRoute.Generation > cachedAPI.ProdHTTPRoute.Generation) {
		cachedAPI.ProdHTTPRoute = prodHTTPRoute
		updated = true
		events = append(events, "Production Endpoint")
	}
	if sandHTTPRoute != nil && (sandHTTPRoute.UID != cachedAPI.SandHTTPRoute.UID ||
		sandHTTPRoute.Generation > cachedAPI.SandHTTPRoute.Generation) {
		cachedAPI.SandHTTPRoute = sandHTTPRoute
		updated = true
		events = append(events, "Sandbox Endpoint")
	}
	for name, authentication := range authentications {
		// if existing map has more recent values for auth cr, then keep them
		if existingAuth, found := cachedAPI.Authentications[name]; found &&
			(existingAuth.UID == authentication.UID || existingAuth.Generation >= authentication.Generation) {
			authentications[name] = existingAuth
		}
		updated = true
		events = append(events, "API Authentication Schemes")
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
	//TODO(amali) remove entry from HTTPRouteToAPIRefs and AuthenticationToAPIRefs
}
