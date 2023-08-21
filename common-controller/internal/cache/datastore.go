/*
 *  Copyright (c) 2023 WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package cache

import (
	"sync"

	logger "github.com/sirupsen/logrus"
	dpv1alpha1 "github.com/wso2/apk/common-controller/internal/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RatelimitDataStore is a cache for rate limit policies.
type RatelimitDataStore struct {
	resolveRatelimitStore map[types.NamespacedName]*dpv1alpha1.ResolveRateLimitAPIPolicy
	customRatelimitStore  map[types.NamespacedName]*dpv1alpha1.CustomRateLimitPolicyDef
	mu                    sync.Mutex
}

// CreateNewOperatorDataStore creates a new RatelimitDataStore.
func CreateNewOperatorDataStore() *RatelimitDataStore {
	return &RatelimitDataStore{
		resolveRatelimitStore: map[types.NamespacedName]*dpv1alpha1.ResolveRateLimitAPIPolicy{},
		customRatelimitStore:  map[types.NamespacedName]*dpv1alpha1.CustomRateLimitPolicyDef{},
	}
}

// AddorUpdateResolveRatelimitToStore adds a new ratelimit to the RatelimitDataStore.
func (ods *RatelimitDataStore) AddorUpdateResolveRatelimitToStore(rateLimit types.NamespacedName,
	resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Adding/Updating ratelimit to cache")
	ods.resolveRatelimitStore[rateLimit] = &resolveRatelimit
}

// AddorUpdateCustomRatelimitToStore adds a new ratelimit to the RatelimitDataStore.
func (ods *RatelimitDataStore) AddorUpdateCustomRatelimitToStore(rateLimit types.NamespacedName,
	customRateLimitPolicy dpv1alpha1.CustomRateLimitPolicyDef) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Adding/Updating custom ratelimit to cache")
	ods.customRatelimitStore[rateLimit] = &customRateLimitPolicy
}

// GetResolveRatelimitPolicy get cached ratelimit
func (ods *RatelimitDataStore) GetResolveRatelimitPolicy(rateLimit types.NamespacedName) (dpv1alpha1.ResolveRateLimitAPIPolicy, bool) {
	var rateLimitPolicy dpv1alpha1.ResolveRateLimitAPIPolicy
	if cachedRatelimit, found := ods.resolveRatelimitStore[rateLimit]; found {
		logger.Debug("Found cached ratelimit")
		return *cachedRatelimit, true
	}
	return rateLimitPolicy, false
}

// GetCachedCustomRatelimitPolicy get cached ratelimit
func (ods *RatelimitDataStore) GetCachedCustomRatelimitPolicy(rateLimit types.NamespacedName) (dpv1alpha1.CustomRateLimitPolicyDef, bool) {
	var rateLimitPolicy dpv1alpha1.CustomRateLimitPolicyDef
	if cachedRatelimit, found := ods.customRatelimitStore[rateLimit]; found {
		logger.Debug("Found cached custom ratelimit")
		return *cachedRatelimit, true
	}
	return rateLimitPolicy, false
}

// DeleteResolveRatelimitPolicy delete from ratelimit cache
func (ods *RatelimitDataStore) DeleteResolveRatelimitPolicy(rateLimit types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Deleting ratelimit from cache")
	delete(ods.resolveRatelimitStore, rateLimit)
}

// DeleteCachedCustomRatelimitPolicy delete from ratelimit cache
func (ods *RatelimitDataStore) DeleteCachedCustomRatelimitPolicy(rateLimit types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Debug("Deleting custom ratelimit from cache")
	delete(ods.customRatelimitStore, rateLimit)
}

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
