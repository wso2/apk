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
	"fmt"

	"sync"

	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// OperatorDataStore holds the APIStore and API, HttpRoute mappings
type OperatorDataStore struct {
	APIStore           map[types.NamespacedName]*APIState
	APIToHTTPRouteRefs map[types.NamespacedName]HTTPRouteRefs
	HTTPRouteToAPIRefs map[types.NamespacedName]types.NamespacedName

	mu sync.Mutex
}

// HTTPRouteRefs holds ProdHttpRouteRef and SandHttpRouteRef
type HTTPRouteRefs struct {
	ProdHTTPRouteRef string
	SandHTTPRouteRef string
}

// CreateNewOperatorDataStore creates a new OperatorDataStore.
func CreateNewOperatorDataStore() *OperatorDataStore {
	return &OperatorDataStore{
		APIStore:           map[types.NamespacedName]*APIState{},
		APIToHTTPRouteRefs: map[types.NamespacedName]HTTPRouteRefs{},
		HTTPRouteToAPIRefs: map[types.NamespacedName]types.NamespacedName{},
	}
}

// AddNewAPI stores a new API in the OperatorDataStore.
func (ods *OperatorDataStore) AddNewAPI(api dpv1alpha1.API, prodHTTPRoute gwapiv1b1.HTTPRoute, sandHTTPRoute gwapiv1b1.HTTPRoute) (APIState, error) {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	ods.APIStore[utils.NamespacedName(&api)] = &APIState{
		APIDefinition: &api,
		ProdHTTPRoute: &prodHTTPRoute,
		SandHTTPRoute: &sandHTTPRoute}
	ods.APIToHTTPRouteRefs[utils.NamespacedName(&api)] = HTTPRouteRefs{ProdHTTPRouteRef: prodHTTPRoute.Name, SandHTTPRouteRef: sandHTTPRoute.Name}
	ods.HTTPRouteToAPIRefs[utils.NamespacedName(&prodHTTPRoute)] = utils.NamespacedName(&api)
	ods.HTTPRouteToAPIRefs[utils.NamespacedName(&sandHTTPRoute)] = utils.NamespacedName(&api)
	return *ods.APIStore[utils.NamespacedName(&api)], nil
}

// UpdateHTTPRoute updates the HttpRoute of a stored API.
func (ods *OperatorDataStore) UpdateHTTPRoute(apiName types.NamespacedName, httpRoute gwapiv1b1.HTTPRoute, production bool) (APIState, error) {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	apiState, found := ods.APIStore[apiName]
	if !found {
		return APIState{}, fmt.Errorf("error: No API found for the key: %v", apiName.String())
	}

	if production {
		apiState.ProdHTTPRoute = &httpRoute
	} else {
		apiState.SandHTTPRoute = &httpRoute
	}
	return *ods.APIStore[apiName], nil
}

// UpdateAPIDef updates the APIDef of a stored API.
func (ods *OperatorDataStore) UpdateAPIDef(apiDef dpv1alpha1.API) (APIState, error) {
	api, found := ods.APIStore[utils.NamespacedName(&apiDef)]
	if !found {
		return APIState{}, fmt.Errorf("API not found in the Operator Data store: %v", apiDef.Spec.APIDisplayName)
	}
	api.APIDefinition = &apiDef
	return *ods.APIStore[utils.NamespacedName(&apiDef)], nil
}

// GetAPI returns the APIState for a given key if exists.
func (ods *OperatorDataStore) GetAPI(apiName types.NamespacedName) (APIState, bool) {
	api, found := ods.APIStore[apiName]
	if !found {
		return APIState{}, found
	}
	return *api, found

}

func (ods *OperatorDataStore) GetAPIForHTTPRoute(httpRoute types.NamespacedName) (types.NamespacedName, bool) {
	api, found := ods.HTTPRouteToAPIRefs[httpRoute]
	if !found {
		return types.NamespacedName{}, found
	}
	return api, found
}
