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
	"math/big"
	"math/rand"
	"sync"
	"time"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoy_cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	wso2_cache "github.com/wso2/apk/adapter/pkg/discovery/protocol/cache/v3"
	eventhubTypes "github.com/wso2/apk/adapter/pkg/eventhub/types"
	"github.com/wso2/apk/common-controller/internal/loggers"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	apimachiner_types "k8s.io/apimachinery/pkg/types"
)

// EnvoyInternalAPI struct use to hold envoy resources and adapter internal resources
type EnvoyInternalAPI struct {
	// commonControllerInternalAPI model.commonControllerInternalAPI
	envoyLabels       []string
	routes            []*routev3.Route
	clusters          []*clusterv3.Cluster
	endpointAddresses []*corev3.Address
	enforcerAPI       types.Resource
}

// EnvoyGatewayConfig struct use to hold envoy gateway resources
type EnvoyGatewayConfig struct {
	listener    *listenerv3.Listener
	routeConfig *routev3.RouteConfiguration
	clusters    []*clusterv3.Cluster
	endpoints   []*corev3.Address
	// customRateLimitPolicies []*model.CustomRateLimitPolicy
}

// EnforcerInternalAPI struct use to hold enforcer resources
type EnforcerInternalAPI struct {
	configs                []types.Resource
	subscriptions          []types.Resource
	applications           []types.Resource
	applicationKeyMappings []types.Resource
	applicationMappings    []types.Resource
}

var (
	// TODO: (VirajSalaka) Remove Unused mutexes.
	mutexForXdsUpdate         sync.Mutex
	mutexForInternalMapUpdate sync.Mutex

	cache                              envoy_cachev3.SnapshotCache
	enforcerCache                      wso2_cache.SnapshotCache
	enforcerSubscriptionCache          wso2_cache.SnapshotCache
	enforcerApplicationCache           wso2_cache.SnapshotCache
	enforcerApplicationKeyMappingCache wso2_cache.SnapshotCache
	enforcerApplicationMappingCache    wso2_cache.SnapshotCache

	orgAPIMap map[string]map[string]*EnvoyInternalAPI // organizationID -> Vhost:API_UUID -> EnvoyInternalAPI struct map

	orgIDvHostBasepathMap map[string]map[string]string   // organizationID -> Vhost:basepath -> Vhost:API_UUID
	orgIDAPIvHostsMap     map[string]map[string][]string // organizationID -> UUID -> prod/sand -> Envoy Vhost Array map

	// Envoy Label as map key
	gatewayLabelConfigMap map[string]*EnvoyGatewayConfig // GW-Label -> EnvoyGatewayConfig struct map

	// Listener as map key
	listenerToRouteArrayMap map[string][]*routev3.Route // Listener -> Routes map

	// Common Enforcer Label as map key
	// TODO(amali) This doesn't have a usage yet. It will be used to handle multiple enforcer labels in future.
	enforcerLabelMap map[string]*EnforcerInternalAPI // Enforcer Label -> EnforcerInternalAPI struct map

	// KeyManagerList to store data
	KeyManagerList = make([]eventhubTypes.KeyManager, 0)
	isReady        = false
)

const (
	maxRandomInt             int    = 999999999
	grpcMaxConcurrentStreams        = 1000000
	apiKeyFieldSeparator     string = ":"
	commonEnforcerLabel      string = "commonEnforcerLabel"
)

func maxRandomBigInt() *big.Int {
	return big.NewInt(int64(maxRandomInt))
}

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

func init() {
	cache = envoy_cachev3.NewSnapshotCache(false, IDHash{}, nil)
	enforcerCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerSubscriptionCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationKeyMappingCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationMappingCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	gatewayLabelConfigMap = make(map[string]*EnvoyGatewayConfig)
	listenerToRouteArrayMap = make(map[string][]*routev3.Route)
	orgAPIMap = make(map[string]map[string]*EnvoyInternalAPI)
	orgIDAPIvHostsMap = make(map[string]map[string][]string) // organizationID -> UUID-prod/sand -> Envoy Vhost Array map
	orgIDvHostBasepathMap = make(map[string]map[string]string)

	enforcerLabelMap = make(map[string]*EnforcerInternalAPI)
	//TODO(amali) currently subscriptions, configs, applications, applicationPolicies, subscriptionPolicies,
	// applicationKeyMappings, keyManagerConfigList, revokedTokens are supported with the hard coded label for Enforcer
	enforcerLabelMap[commonEnforcerLabel] = &EnforcerInternalAPI{}
	rand.Seed(time.Now().UnixNano())
	// go watchEnforcerResponse()
}

// GetRateLimiterCache returns xds server cache for rate limiter service.
func GetRateLimiterCache() envoy_cachev3.SnapshotCache {
	return rlsPolicyCache.xdsCache
}

// UpdateRateLimitXDSCache updates the xDS cache of the RateLimiter.
func UpdateRateLimitXDSCache(resolveRatelimitPolicyList []dpv1alpha1.ResolveRateLimitAPIPolicy) {

	for _, resolveRatelimitPolicy := range resolveRatelimitPolicyList {
		// Add Rate Limit inline policies in API to the cache
		rlsPolicyCache.AddAPILevelRateLimitPolicies(resolveRatelimitPolicy)
	}
}

// UpdateRateLimitXDSCacheForCustomPolicies updates the xDS cache of the RateLimiter for custom policies.
func UpdateRateLimitXDSCacheForCustomPolicies(customRateLimitPolicies dpv1alpha1.CustomRateLimitPolicyDef) {
	if customRateLimitPolicies.Key != "" {
		rlsPolicyCache.AddCustomRateLimitPolicies(customRateLimitPolicies)
	}
}

// UpdateRateLimitXDSCacheForAIRatelimitPolicies updates the xDS cache of the RateLimiter for AI ratelimit policies.
func UpdateRateLimitXDSCacheForAIRatelimitPolicies(aiRatelimitPolicySpecs map[apimachiner_types.NamespacedName]*dpv1alpha3.AIRateLimitPolicySpec) {
	rlsPolicyCache.ProcessAIRatelimitPolicySpecsAndUpdateCache(aiRatelimitPolicySpecs)
}

// DeleteAPILevelRateLimitPolicies delete the ratelimit xds cache
func DeleteAPILevelRateLimitPolicies(resolveRatelimitPolicyList []dpv1alpha1.ResolveRateLimitAPIPolicy) {

	for _, resolveRatelimit := range resolveRatelimitPolicyList {
		var org = resolveRatelimit.Organization
		var environment = resolveRatelimit.Environment
		var basePath = resolveRatelimit.BasePath
		rlsPolicyCache.DeleteAPILevelRateLimitPolicies(org, environment, basePath)
	}
}

// DeleteResourceLevelRateLimitPolicies delete the ratelimit xds cache
func DeleteResourceLevelRateLimitPolicies(resolveRatelimitPolicyList []dpv1alpha1.ResolveRateLimitAPIPolicy) {

	for _, resolveRatelimit := range resolveRatelimitPolicyList {

		if resolveRatelimit.Resources == nil || len(resolveRatelimit.Resources) == 0 {
			continue
		}
		var org = resolveRatelimit.Organization
		var environment = resolveRatelimit.Environment
		var basePath = resolveRatelimit.BasePath
		var path = resolveRatelimit.Resources[0].Path
		var method = resolveRatelimit.Resources[0].Method
		rlsPolicyCache.DeleteResourceLevelRateLimitPolicies(org, environment, basePath, path, method)
	}
}

// DeleteSubscriptionRateLimitPolicies delete the ratelimit xds cache
func DeleteSubscriptionRateLimitPolicies(resolveSubscriptionRatelimit dpv1alpha3.ResolveSubscriptionRatelimitPolicy) {
	rlsPolicyCache.RemoveSubscriptionRateLimitPolicy(resolveSubscriptionRatelimit)
}

// UpdateRateLimitXDSCacheForSubscriptionPolicies updates the xDS cache of the RateLimiter for subscription policies.
func UpdateRateLimitXDSCacheForSubscriptionPolicies(resolveSubscriptionRatelimit dpv1alpha3.ResolveSubscriptionRatelimitPolicy) {
	rlsPolicyCache.AddSubscriptionLevelRateLimitPolicy(resolveSubscriptionRatelimit)
}

// DeleteCustomRateLimitPolicies delete the ratelimit xds cache
func DeleteCustomRateLimitPolicies(customRateLimitPolicy dpv1alpha1.CustomRateLimitPolicyDef) {
	rlsPolicyCache.DeleteCustomRateLimitPolicies(customRateLimitPolicy)
}

// GenerateIdentifierForAPIWithUUID generates an identifier unique to the API
func GenerateIdentifierForAPIWithUUID(vhost, uuid string) string {
	return fmt.Sprint(vhost, apiKeyFieldSeparator, uuid)
}

// UpdateRateLimiterPolicies update the rate limiter xDS cache with latest rate limit policies
func UpdateRateLimiterPolicies(label string) {
	_ = rlsPolicyCache.updateXdsCache(label)
}

// SetEmptySnapshotupdate update empty snapshot
func SetEmptySnapshotupdate(lable string) bool {
	return rlsPolicyCache.SetEmptySnapshot(lable)
}

// GetXdsCache returns xds server cache.
func GetXdsCache() envoy_cachev3.SnapshotCache {
	return cache
}

// GetEnforcerCache returns xds server cache.
func GetEnforcerCache() wso2_cache.SnapshotCache {
	return enforcerCache
}

// GetEnforcerSubscriptionCache returns xds server cache.
func GetEnforcerSubscriptionCache() wso2_cache.SnapshotCache {
	return enforcerSubscriptionCache
}

// GetEnforcerApplicationCache returns xds server cache.
func GetEnforcerApplicationCache() wso2_cache.SnapshotCache {
	return enforcerApplicationCache
}

// GetEnforcerApplicationKeyMappingCache returns xds server cache.
func GetEnforcerApplicationKeyMappingCache() wso2_cache.SnapshotCache {
	return enforcerApplicationKeyMappingCache
}

// GetEnforcerApplicationMappingCache returns xds server cache.
func GetEnforcerApplicationMappingCache() wso2_cache.SnapshotCache {
	return enforcerApplicationMappingCache
}
