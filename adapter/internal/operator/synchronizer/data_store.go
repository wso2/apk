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
	APIStore           map[types.NamespacedName]*APIState
	HTTPRouteToAPIRefs map[types.NamespacedName]types.NamespacedName
	mu                 sync.Mutex
}

// CreateNewOperatorDataStore creates a new OperatorDataStore.
func CreateNewOperatorDataStore() *OperatorDataStore {
	return &OperatorDataStore{
		APIStore:           map[types.NamespacedName]*APIState{},
		HTTPRouteToAPIRefs: map[types.NamespacedName]types.NamespacedName{},
	}
}

// AddNewAPI stores a new API in the OperatorDataStore.
func (ods *OperatorDataStore) AddNewAPI(api dpv1alpha1.API, prodHTTPRoute *gwapiv1b1.HTTPRoute, sandHTTPRoute *gwapiv1b1.HTTPRoute) APIState {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	ods.APIStore[utils.NamespacedName(&api)] = &APIState{
		APIDefinition: &api,
		ProdHTTPRoute: prodHTTPRoute,
		SandHTTPRoute: sandHTTPRoute}

	if prodHTTPRoute != nil {
		ods.HTTPRouteToAPIRefs[utils.NamespacedName(prodHTTPRoute)] = utils.NamespacedName(&api)
	}
	if sandHTTPRoute != nil {
		ods.HTTPRouteToAPIRefs[utils.NamespacedName(sandHTTPRoute)] = utils.NamespacedName(&api)
	}
	return *ods.APIStore[utils.NamespacedName(&api)]
}
