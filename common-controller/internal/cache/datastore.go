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
	"reflect"
	"sync"

	logger "github.com/sirupsen/logrus"
	dpv1alpha1 "github.com/wso2/apk/common-controller/internal/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RatelimitDataStore is a cache for rate limit policies.
type RatelimitDataStore struct {
	ratelimitStore       map[types.NamespacedName]*dpv1alpha1.ResolveRateLimitAPIPolicy
	apisToRateLimit      map[types.NamespacedName]*dpv1alpha1.RateLimitPolicyList
	httpRouteToRateLimit map[types.NamespacedName]*dpv1alpha1.RateLimitPolicyList
	mu                   sync.Mutex
}

// CreateNewOperatorDataStore creates a new RatelimitDataStore.
func CreateNewOperatorDataStore() *RatelimitDataStore {
	return &RatelimitDataStore{
		ratelimitStore:       map[types.NamespacedName]*dpv1alpha1.ResolveRateLimitAPIPolicy{},
		apisToRateLimit:      map[types.NamespacedName]*dpv1alpha1.RateLimitPolicyList{},
		httpRouteToRateLimit: map[types.NamespacedName]*dpv1alpha1.RateLimitPolicyList{},
	}
}

// Resolve Ratelimit cache
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

// RatelimitPloicyForAPI Cache
// GetRatelimitsToAPI returns the list of rate limits for the specified API.
func (ods *RatelimitDataStore) GetRatelimitsToAPI(api types.NamespacedName) *dpv1alpha1.RateLimitPolicyList {
	return ods.apisToRateLimit[api]
}

// DeleteRatelimitToAPI deletes the list of rate limits for the specified API.
func (ods *RatelimitDataStore) DeleteRatelimitToAPI(api types.NamespacedName) {
	delete(ods.apisToRateLimit, api)
}

// AddRatelimitToAPI adds a rate limit to the list of rate limits for the specified API.
func (ods *RatelimitDataStore) AddRatelimitToAPI(key types.NamespacedName, ratelimit dpv1alpha1.RateLimitPolicy) {
	if ods.apisToRateLimit[key] == nil {
		ods.apisToRateLimit[key] = &dpv1alpha1.RateLimitPolicyList{
			Items: []dpv1alpha1.RateLimitPolicy{ratelimit},
		}
	} else {
		ods.apisToRateLimit[key].Items = append(ods.apisToRateLimit[key].Items, ratelimit)
	}
}

// IsRateLimitPolicyAvailble checks whether the specified rate limit policy is available in the cache.
func (ods *RatelimitDataStore) IsRateLimitPolicyAvailble(key types.NamespacedName, desiredPolicy dpv1alpha1.RateLimitPolicy) bool {
	if list, ok := ods.apisToRateLimit[key]; ok {
		for _, policy := range list.Items {
			if policiesEqual(policy, desiredPolicy) {
				return true
			}
		}
	}
	return false
}

// RemoveRatelimitPolicyByNamespacedName removes the rate limit policy with the specified namespaced name from the cache.
func (ods *RatelimitDataStore) RemoveRatelimitPolicyByNamespacedName(namespacedName types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	// Iterate through the rate limit policies in the cache.
	for api, rateLimitPolicyList := range ods.apisToRateLimit {
		for i, policy := range rateLimitPolicyList.Items {
			if policy.ObjectMeta.Namespace == namespacedName.Namespace && policy.ObjectMeta.Name == namespacedName.Name {
				// Found the rate limit policy with the specified namespaced name.
				// Remove the policy from the list.
				ods.apisToRateLimit[api].Items = append(rateLimitPolicyList.Items[:i], rateLimitPolicyList.Items[i+1:]...)
				return // Found and removed the policy, so return.
			}
		}
	}
}

// Helper function to compare policies.
func policiesEqual(policy1, policy2 dpv1alpha1.RateLimitPolicy) bool {
	// Compare relevant fields here.
	return reflect.DeepEqual(policy1.Spec, policy2.Spec)
}

// HTTPRouteToRateLimit Cache
// AddRateLimitToAPI adds a rate limit policy to the list associated with the given key.
func (ods *RatelimitDataStore) AddRateLimitToHTTPRoute(key types.NamespacedName, policy dpv1alpha1.RateLimitPolicy) {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	if ods.httpRouteToRateLimit[key] == nil {
		ods.httpRouteToRateLimit[key] = &dpv1alpha1.RateLimitPolicyList{
			Items: []dpv1alpha1.RateLimitPolicy{policy},
		}
	} else {
		ods.httpRouteToRateLimit[key].Items = append(ods.httpRouteToRateLimit[key].Items, policy)
	}
}

// GetRateLimitPolicyList returns the list of rate limit policies associated with the given key.
// If the key is not found, it returns nil.
func (ods *RatelimitDataStore) GetRateLimitPolicyForHTTPRoute(key types.NamespacedName) *dpv1alpha1.RateLimitPolicyList {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	return ods.httpRouteToRateLimit[key]
}

// DeleteRateLimitPolicyList deletes the rate limit policy list associated with the given key.
func (ods *RatelimitDataStore) DeleteRateLimitPolicyListForHTTPRoute(key types.NamespacedName) {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	delete(ods.httpRouteToRateLimit, key)
}

// RemoveRateLimitPolicyFromList removes a specific rate limit policy from the list associated with the given key.
func (ods *RatelimitDataStore) RemoveRateLimitPolicyFromListHttpRoute(key types.NamespacedName, policyToDelete dpv1alpha1.RateLimitPolicy) {
	ods.mu.Lock()
	defer ods.mu.Unlock()

	if list, ok := ods.httpRouteToRateLimit[key]; ok {
		for i, policy := range list.Items {
			if policiesEqual(policy, policyToDelete) {
				ods.httpRouteToRateLimit[key].Items = append(list.Items[:i], list.Items[i+1:]...)
				return // Found and removed the policy, so return.
			}
		}
	}
}

// NamespacedName generates namespaced name for Kubernetes objects
func NamespacedName(obj client.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
