/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org).
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

package xds

import (
	"fmt"
	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	dpv1alpha1 "github.com/wso2/apk/common-adapter/internal/operator/api/v1alpha1"
)

const (
	maxRandomInt             int    = 999999999
	grpcMaxConcurrentStreams        = 1000000
	apiKeyFieldSeparator     string = ":"
)

// IDHash uses ID field as the node hash.
type IDHash struct{}

// ID uses the node ID field
func (IDHash) ID(node *corev3.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

var _ envoy_cachev3.NodeHash = IDHash{}

// GetRateLimiterCache returns xds server cache for rate limiter service.
func GetRateLimiterCache() envoy_cachev3.SnapshotCache {
	return rlsPolicyCache.xdsCache
}

// UpdateRateLimitXDSCache updates the xDS cache of the RateLimiter.
func UpdateRateLimitXDSCache(vhosts []string, resolveRatelimit dpv1alpha1.ResolveRateLimitAPIPolicy) {
	// Add Rate Limit inline policies in API to the cache
	rlsPolicyCache.AddAPILevelRateLimitPolicies(vhosts, resolveRatelimit)
}

// DeleteAPILevelRateLimitPolicies delete the ratelimit xds cache
func DeleteAPILevelRateLimitPolicies(org string, vhost string, apiID string) {
	rlsPolicyCache.DeleteAPILevelRateLimitPolicies(org, vhost, apiID)
}

// GenerateIdentifierForAPIWithUUID generates an identifier unique to the API
func GenerateIdentifierForAPIWithUUID(vhost, uuid string) string {
	return fmt.Sprint(vhost, apiKeyFieldSeparator, uuid)
}

// UpdateRateLimiterPolicies update the rate limiter xDS cache with latest rate limit policies
func UpdateRateLimiterPolicies(label string) {
	_ = rlsPolicyCache.updateXdsCache(label)
}

// ExtractVhostFromAPIIdentifier extracts vhost from the API identifier
func ExtractVhostFromAPIIdentifier(id string) (string, error) {
	elem := strings.Split(id, apiKeyFieldSeparator)
	if len(elem) == 2 {
		return elem[0], nil
	}
	err := fmt.Errorf("invalid API identifier: %v", id)
	return "", err
}

// SetEmptySnapshotupdate update empty snapshot
func SetEmptySnapshotupdate(lable string) bool {
	return rlsPolicyCache.SetEmptySnapshot(lable)
}
