/*
 *  Copyright (c) 2023 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
	ratelimitStore map[types.NamespacedName]*dpv1alpha1.ResolveRateLimitAPIPolicy
	mu             sync.Mutex
}

// CreateNewOperatorDataStore creates a new RatelimitDataStore.
func CreateNewOperatorDataStore() *RatelimitDataStore {
	return &RatelimitDataStore{
		ratelimitStore: map[types.NamespacedName]*dpv1alpha1.ResolveRateLimitAPIPolicy{},
	}
}

// AddorUpdateRatelimitToStore adds a new ratelimit to the RatelimitDataStore.
func (ods *RatelimitDataStore) AddorUpdateRatelimitToStore(rateLimit types.NamespacedName,
	resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Info("Adding/Updating ratelimit to cache")
	ods.ratelimitStore[rateLimit] = &resolveRatelimit
	logger.Info("resolveRatelimit: ", resolveRatelimit)
}

// GetCachedRatelimitPolicy get cached ratelimit
func (ods *RatelimitDataStore) GetCachedRatelimitPolicy(rateLimit types.NamespacedName) (dpv1alpha1.ResolveRateLimitAPIPolicy, bool) {
	var rateLimitPolicy dpv1alpha1.ResolveRateLimitAPIPolicy
	if cachedRatelimit, found := ods.ratelimitStore[rateLimit]; found {
		logger.Info("Found cached ratelimit")
		logger.Info("cachedRatelimit: ", cachedRatelimit)
		return *cachedRatelimit, true
	}
	return rateLimitPolicy, false
}

// DeleteCachedRatelimitPolicy delete from ratelimit cache
func (ods *RatelimitDataStore) DeleteCachedRatelimitPolicy(rateLimit types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()
	logger.Info("Deleting ratelimit from cache")
	delete(ods.ratelimitStore, rateLimit)
}

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
