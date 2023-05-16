/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

// Package xds contains the implementation for the xds server cache updates
package xds

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoy_cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	envoy_resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	oasParser "github.com/wso2/apk/adapter/internal/oasparser"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	operatorconsts "github.com/wso2/apk/adapter/internal/operator/constants"
	disc_api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	wso2_cache "github.com/wso2/apk/adapter/pkg/discovery/protocol/cache/v3"
	wso2_resource "github.com/wso2/apk/adapter/pkg/discovery/protocol/resource/v3"
	eventhubTypes "github.com/wso2/apk/adapter/pkg/eventhub/types"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// EnvoyInternalAPI struct use to hold envoy resources and adapter internal resources
type EnvoyInternalAPI struct {
	adapterInternalAPI model.AdapterInternalAPI
	envoyLabels        []string
	routes             []*routev3.Route
	clusters           []*clusterv3.Cluster
	endpointAddresses  []*corev3.Address
	enforcerAPI        types.Resource
}

// EnvoyGatewayConfig struct use to hold envoy gateway resources
type EnvoyGatewayConfig struct {
	listener                *listenerv3.Listener
	routeConfig             *routev3.RouteConfiguration
	clusters                []*clusterv3.Cluster
	endpoints               []*corev3.Address
	customRateLimitPolicies []*model.CustomRateLimitPolicy
}

// EnforcerInternalAPI struct use to hold enforcer resources
type EnforcerInternalAPI struct {
	configs                []types.Resource
	keyManagers            []types.Resource
	subscriptions          []types.Resource
	applications           []types.Resource
	apiList                []types.Resource
	applicationPolicies    []types.Resource
	subscriptionPolicies   []types.Resource
	applicationKeyMappings []types.Resource
	revokedTokens          []types.Resource
}

var (
	// TODO: (VirajSalaka) Remove Unused mutexes.
	mutexForXdsUpdate         sync.Mutex
	mutexForInternalMapUpdate sync.Mutex

	cache                              envoy_cachev3.SnapshotCache
	enforcerCache                      wso2_cache.SnapshotCache
	enforcerSubscriptionCache          wso2_cache.SnapshotCache
	enforcerApplicationCache           wso2_cache.SnapshotCache
	enforcerAPICache                   wso2_cache.SnapshotCache
	enforcerApplicationPolicyCache     wso2_cache.SnapshotCache
	enforcerSubscriptionPolicyCache    wso2_cache.SnapshotCache
	enforcerApplicationKeyMappingCache wso2_cache.SnapshotCache
	enforcerKeyManagerCache            wso2_cache.SnapshotCache
	enforcerRevokedTokensCache         wso2_cache.SnapshotCache
	enforcerThrottleDataCache          wso2_cache.SnapshotCache

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

var void struct{}

const (
	commonEnforcerLabel  string = "commonEnforcerLabel"
	maxRandomInt         int    = 999999999
	prototypedAPI        string = "PROTOTYPED"
	apiKeyFieldSeparator string = ":"
	gatewayController    string = "GatewayController"
	apiController        string = "APIController"
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

func init() {
	cache = envoy_cachev3.NewSnapshotCache(false, IDHash{}, nil)
	enforcerCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerSubscriptionCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerAPICache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationPolicyCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerSubscriptionPolicyCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationKeyMappingCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerKeyManagerCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerRevokedTokensCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerThrottleDataCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)

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

// GetXdsCache returns xds server cache.
func GetXdsCache() envoy_cachev3.SnapshotCache {
	return cache
}

// GetRateLimiterCache returns xds server cache for rate limiter service.
func GetRateLimiterCache() envoy_cachev3.SnapshotCache {
	return rlsPolicyCache.xdsCache
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

// GetEnforcerAPICache returns xds server cache.
func GetEnforcerAPICache() wso2_cache.SnapshotCache {
	return enforcerAPICache
}

// GetEnforcerApplicationPolicyCache returns xds server cache.
func GetEnforcerApplicationPolicyCache() wso2_cache.SnapshotCache {
	return enforcerApplicationPolicyCache
}

// GetEnforcerSubscriptionPolicyCache returns xds server cache.
func GetEnforcerSubscriptionPolicyCache() wso2_cache.SnapshotCache {
	return enforcerSubscriptionPolicyCache
}

// GetEnforcerApplicationKeyMappingCache returns xds server cache.
func GetEnforcerApplicationKeyMappingCache() wso2_cache.SnapshotCache {
	return enforcerApplicationKeyMappingCache
}

// GetEnforcerKeyManagerCache returns xds server cache.
func GetEnforcerKeyManagerCache() wso2_cache.SnapshotCache {
	return enforcerKeyManagerCache
}

// GetEnforcerRevokedTokenCache return token cache
func GetEnforcerRevokedTokenCache() wso2_cache.SnapshotCache {
	return enforcerRevokedTokensCache
}

// GetEnforcerThrottleDataCache return throttle data cache
func GetEnforcerThrottleDataCache() wso2_cache.SnapshotCache {
	return enforcerThrottleDataCache
}

// DeleteAPICREvent deletes API with the given UUID from the given gw environments
func DeleteAPICREvent(labels []string, apiUUID string, organizationID string) error {
	mutexForInternalMapUpdate.Lock()
	defer mutexForInternalMapUpdate.Unlock()

	prodvHostIdentifier := GetvHostsIdentifier(apiUUID, operatorconsts.Production)
	sandvHostIdentifier := GetvHostsIdentifier(apiUUID, operatorconsts.Sandbox)
	vHosts := append(orgIDAPIvHostsMap[organizationID][prodvHostIdentifier],
		orgIDAPIvHostsMap[organizationID][sandvHostIdentifier]...)

	delete(orgIDAPIvHostsMap[organizationID], prodvHostIdentifier)
	delete(orgIDAPIvHostsMap[organizationID], sandvHostIdentifier)
	for _, vhost := range vHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vhost, apiUUID)
		if err := deleteAPI(apiIdentifier, labels, organizationID); err != nil {
			logger.LoggerXds.ErrorC(logging.GetErrorByCode(1410, apiIdentifier, organizationID, labels))
			return err
		}
		// if no error, update internal vhost maps
		// error only happens when API not found in deleteAPI func
		logger.LoggerXds.Infof("Successfully undeployed the API %v under Organization %s and environment %s ",
			apiIdentifier, organizationID, labels)
	}
	return nil
}

// deleteAPI deletes an API, its resources and updates the caches of given environments
func deleteAPI(apiIdentifier string, environments []string, organizationID string) error {
	_, exists := orgAPIMap[organizationID][apiIdentifier]
	if !exists {
		logger.LoggerXds.Infof("Unable to delete API: %v from Organization: %v. API Does not exist.", apiIdentifier, organizationID)
		return errors.New(constants.NotFound)
	}

	existingLabels := orgAPIMap[organizationID][apiIdentifier].envoyLabels
	toBeDelEnvs, toBeKeptEnvs := getEnvironmentsToBeDeleted(existingLabels, environments)

	for _, val := range toBeDelEnvs {
		isAllowedToDelete := stringutils.StringInSlice(val, existingLabels)
		if isAllowedToDelete {
			// do not delete from all environments, hence do not clear routes, clusters, endpoints, enforcerAPIs
			orgAPIMap[organizationID][apiIdentifier].envoyLabels = toBeKeptEnvs
			updateXdsCacheOnAPIChange(toBeDelEnvs, []string{})
			if len(toBeKeptEnvs) != 0 {
				return nil
			}
			logger.LoggerXds.Infof("API identifier: %v does not have any gateways. Hence deleting the API from label : %s.",
				apiIdentifier, val)
			cleanMapResources(apiIdentifier, organizationID, toBeDelEnvs)
			return nil
		}
	}

	//clean maps of routes, clusters, endpoints, enforcerAPIs
	if len(environments) == 0 {
		cleanMapResources(apiIdentifier, organizationID, toBeDelEnvs)
	}
	return nil
}

func cleanMapResources(apiIdentifier string, organizationID string, toBeDelEnvs []string) {
	delete(orgAPIMap[organizationID], apiIdentifier)
	vHost, err := ExtractVhostFromAPIIdentifier(apiIdentifier)
	if err != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1713, apiIdentifier, err))
	} else {
		rlsPolicyCache.DeleteAPILevelRateLimitPolicies(organizationID, vHost, apiIdentifier)
	}

	//updateXdsCacheOnAPIAdd is called after cleaning maps of routes, clusters, endpoints, enforcerAPIs.
	//Therefore resources that belongs to the deleting API do not exist. Caches updated only with
	//resources that belongs to the remaining APIs
	updateXdsCacheOnAPIChange(toBeDelEnvs, []string{})

	deleteBasepathForVHost(organizationID, apiIdentifier)
	//TODO: (SuKSW) clean any remaining in label wise maps, if this is the last API of that label
	logger.LoggerXds.Infof("Deleted API %v of organization %v", apiIdentifier, organizationID)
}

func deleteBasepathForVHost(organizationID, apiIdentifier string) {
	// Remove the basepath from map (that is used to avoid duplicate basepaths)
	if oldOrgAPIAPI, ok := orgAPIMap[organizationID][apiIdentifier]; ok {
		s := strings.Split(apiIdentifier, apiKeyFieldSeparator)
		vHost := s[0]
		oldBasepath := oldOrgAPIAPI.adapterInternalAPI.GetXWso2Basepath()
		delete(orgIDvHostBasepathMap[organizationID], vHost+":"+oldBasepath)
	}
}

// when this method is called, openAPIEnvoy map is updated.
// Old labels refers to the previously assigned labels
// New labels refers to the the updated labels
func updateXdsCacheOnAPIChange(oldLabels []string, newLabels []string) bool {
	revisionStatus := false
	// TODO: (VirajSalaka) check possible optimizations, Since the number of labels are low by design it should not be an issue
	for _, newLabel := range newLabels {
		listeners, clusters, routes, endpoints, apis := GenerateEnvoyResoucesForGateway(newLabel)
		UpdateEnforcerApis(newLabel, apis, "")
		UpdateRateLimiterPolicies(newLabel)
		success := UpdateXdsCacheWithLock(newLabel, endpoints, clusters, routes, listeners)
		logger.LoggerXds.Debugf("Xds Cache is updated for the newly added label : %v", newLabel)
		if success {
			// if even one label was updated with latest revision, we take the revision as deployed.
			// (other labels also will get updated successfully)
			revisionStatus = success
			continue
		}
	}
	for _, oldLabel := range oldLabels {
		if !stringutils.StringInSlice(oldLabel, newLabels) {
			listeners, clusters, routes, endpoints, apis := GenerateEnvoyResoucesForGateway(oldLabel)
			UpdateEnforcerApis(oldLabel, apis, "")
			UpdateRateLimiterPolicies(oldLabel)
			UpdateXdsCacheWithLock(oldLabel, endpoints, clusters, routes, listeners)
			logger.LoggerXds.Debugf("Xds Cache is updated for the already existing label : %v", oldLabel)
		}
	}
	return revisionStatus
}

// SetReady Method to set the status after the last api is fected and updated in router.
func SetReady() {
	logger.LoggerXds.Infof("Finished deploying startup APIs. Deploying the readiness endpoint...")
	isReady = true
}

// GenerateEnvoyResoucesForGateway generates envoy resources for a given gateway
// This method will list out all APIs mapped to the label. and generate envoy resources for all of these APIs.
func GenerateEnvoyResoucesForGateway(gatewayName string) ([]types.Resource,
	[]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	var clusterArray []*clusterv3.Cluster
	var vhostToRouteArrayMap = make(map[string][]*routev3.Route)
	var endpointArray []*corev3.Address
	var apis []types.Resource

	for organizationID, entityMap := range orgAPIMap {
		for apiKey, envoyInternalAPI := range entityMap {
			if stringutils.StringInSlice(gatewayName, envoyInternalAPI.envoyLabels) {
				vhost, err := ExtractVhostFromAPIIdentifier(apiKey)
				if err != nil {
					logger.LoggerXds.ErrorC(logging.GetErrorByCode(1411, err.Error(), organizationID))
					continue
				}
				isDefaultVersion := false
				var orgAPI *EnvoyInternalAPI
				var ok bool
				if orgAPI, ok = orgAPIMap[organizationID][apiKey]; !ok {
					// If the adapterInternalAPI is not found, proceed with other APIs. (Unreachable condition at this point)
					// If that happens, there is no purpose in processing clusters too.
					continue
				}
				isDefaultVersion = orgAPI.adapterInternalAPI.IsDefaultVersion
				// If it is a default versioned API, the routes are added to the end of the existing array.
				// Otherwise the routes would be added to the front.
				// /fooContext/2.0.0/* resource path should be matched prior to the /fooContext/* .
				if isDefaultVersion {
					vhostToRouteArrayMap[vhost] = append(vhostToRouteArrayMap[vhost], orgAPI.routes...)
				} else {
					vhostToRouteArrayMap[vhost] = append(orgAPI.routes, vhostToRouteArrayMap[vhost]...)
				}
				clusterArray = append(clusterArray, orgAPI.clusters...)
				endpointArray = append(endpointArray, orgAPI.endpointAddresses...)
				apis = append(apis, orgAPI.enforcerAPI)
			}
		}
	}

	// If the token endpoint is enabled, the token endpoint also needs to be added.
	conf := config.ReadConfigs()
	enableJwtIssuer := conf.Enforcer.JwtIssuer.Enabled
	systemHost := conf.Envoy.SystemHost

	logger.LoggerXds.Debugf("System Host : %v", systemHost)
	if enableJwtIssuer {
		routeToken := envoyconf.CreateTokenRoute()
		vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], routeToken)
	}

	// Add health endpoint
	routeHealth := envoyconf.CreateHealthEndpoint()
	vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], routeHealth)

	// Add the readiness endpoint. isReady flag will be set to true once all the apis are fetched from the control plane
	if isReady {
		readynessEndpoint := envoyconf.CreateReadyEndpoint()
		vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], readynessEndpoint)
	}

	envoyGatewayConfig, gwFound := gatewayLabelConfigMap[gatewayName]
	if !gwFound {
		return nil, nil, nil, nil, nil
	}

	listener := envoyGatewayConfig.listener
	logger.LoggerXds.Debugf("Listener : %v", listener)
	routesFromListener := listenerToRouteArrayMap[listener.Name]
	logger.LoggerXds.Debugf("Routes from listener : %v", routesFromListener)
	var vhostToRouteArrayFilteredMap = make(map[string][]*routev3.Route)
	for vhost, routes := range vhostToRouteArrayMap {
		logger.LoggerXds.Debugf("Routes from Vhost Map : %v", routes)
		if vhost == systemHost || checkRoutes(routes, routesFromListener) {
			logger.LoggerXds.Debugf("Equal routes : %v", routes)
			vhostToRouteArrayFilteredMap[vhost] = routes
		}
	}
	routesConfig := oasParser.GetRouteConfigs(vhostToRouteArrayFilteredMap, listener.Name, envoyGatewayConfig.customRateLimitPolicies)
	envoyGatewayConfig.routeConfig = routesConfig
	logger.LoggerXds.Debugf("Listener : %v and routes %v", listener, routesConfig)

	logger.LoggerXds.Debugf("Routes Config : %v", routesConfig)
	clusterArray = append(clusterArray, gatewayLabelConfigMap[gatewayName].clusters...)
	endpointArray = append(endpointArray, gatewayLabelConfigMap[gatewayName].endpoints...)
	endpoints, clusters, listeners, routeConfigs := oasParser.GetCacheResources(endpointArray, clusterArray, listener, routesConfig)
	logger.LoggerXds.Debugf("Routes Config After Get cache : %v", routeConfigs)
	return endpoints, clusters, listeners, routeConfigs, apis
}

// function to check routes []*routev3.Route equlas routes []*routev3.Route
func checkRoutes(routes []*routev3.Route, routesFromListener []*routev3.Route) bool {
	for i := range routes {
		flag := false
		for j := range routesFromListener {
			if routes[i].Name == routesFromListener[j].Name {
				flag = true
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

// GenerateGlobalClusters generates the globally available clusters and endpoints.
func GenerateGlobalClusters(label string) {
	clusters, endpoints := oasParser.GetGlobalClusters()
	gatewayLabelConfigMap[label] = &EnvoyGatewayConfig{
		clusters:  clusters,
		endpoints: endpoints,
	}
}

// GenerateInterceptorClusters generates the globally available clusters and endpoints with interceptors.
func GenerateInterceptorClusters(label string,
	gwReqICluster *clusterv3.Cluster, gwReqIAddresses []*corev3.Address,
	gwResICluster *clusterv3.Cluster, gwResIAddresses []*corev3.Address) {
	var clusters []*clusterv3.Cluster
	var endpoints []*corev3.Address

	if gwReqICluster != nil && len(gwReqIAddresses) > 0 {
		clusters = append(clusters, gwReqICluster)
		endpoints = append(endpoints, gwReqIAddresses...)
	}

	if gwResICluster != nil && len(gwResIAddresses) > 0 {
		clusters = append(clusters, gwResICluster)
		endpoints = append(endpoints, gwResIAddresses...)
	}

	gatewayLabelConfigMap[label].clusters = append(gatewayLabelConfigMap[label].clusters, clusters...)
	gatewayLabelConfigMap[label].endpoints = append(gatewayLabelConfigMap[label].endpoints, endpoints...)
}

// use UpdateXdsCacheWithLock to avoid race conditions
func updateXdsCache(label string, endpoints []types.Resource, clusters []types.Resource, routes []types.Resource, listeners []types.Resource) bool {
	version := rand.Intn(maxRandomInt)
	// TODO: (VirajSalaka) kept same version for all the resources as we are using simple cache implementation.
	// Will be updated once decide to move to incremental XDS
	snap, errNewSnap := envoy_cachev3.NewSnapshot(fmt.Sprint(version), map[envoy_resource.Type][]types.Resource{
		envoy_resource.EndpointType: endpoints,
		envoy_resource.ClusterType:  clusters,
		envoy_resource.ListenerType: listeners,
		envoy_resource.RouteType:    routes,
	})
	if errNewSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1413, errNewSnap.Error()))
		return false
	}
	snap.Consistent()
	//TODO: (VirajSalaka) check
	errSetSnap := cache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
		return false
	}
	logger.LoggerXds.Infof("New Router cache updated for the label: " + label + " version: " + fmt.Sprint(version))
	return true
}

// UpdateRateLimiterPolicies update the rate limiter xDS cache with latest rate limit policies
func UpdateRateLimiterPolicies(label string) {
	_ = rlsPolicyCache.updateXdsCache(label)
}

// UpdateEnforcerConfig Sets new update to the enforcer's configuration
func UpdateEnforcerConfig(configFile *config.Config) {
	// TODO: (Praminda) handle labels
	label := commonEnforcerLabel
	configs := []types.Resource{MarshalConfig(configFile)}
	version := rand.Intn(maxRandomInt)
	snap, errNewSnap := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ConfigType: configs,
	})
	if errNewSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1413, errNewSnap.Error()))
	}
	snap.Consistent()

	errSetSnap := enforcerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}

	enforcerLabelMap[label].configs = configs
	logger.LoggerXds.Infof("New Config cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApis Sets new update to the enforcer's Apis
func UpdateEnforcerApis(label string, apis []types.Resource, version string) {

	if version == "" {
		version = fmt.Sprint(rand.Intn(maxRandomInt))
	}

	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.APIType: apis,
	})
	snap.Consistent()

	errSetSnap := enforcerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	logger.LoggerXds.Infof("New API cache update for the label: " + label + " version: " + fmt.Sprint(version))

	subAPIs := []*subscription.APIs{}
	for _, api := range apis {
		subAPI := subscription.APIs{}
		subAPI.ApiId = api.(*disc_api.Api).GetId()
		subAPI.Name = api.(*disc_api.Api).GetTitle()
		subAPI.Version = api.(*disc_api.Api).GetVersion()
		subAPI.Context = api.(*disc_api.Api).GetBasePath()
		subAPI.Policy = api.(*disc_api.Api).GetTier()
		subAPI.ApiType = api.(*disc_api.Api).GetApiType()
		subAPI.Uuid = api.(*disc_api.Api).GetId()
		subAPIs = append(subAPIs, &subAPI)
	}
	subAPIList := &subscription.APIList{
		List: subAPIs,
	}
	UpdateEnforcerAPIList(label, subAPIList)
}

// UpdateEnforcerSubscriptions sets new update to the enforcer's Subscriptions
func UpdateEnforcerSubscriptions(subscriptions *subscription.SubscriptionList) {
	//TODO: (Dinusha) check this hardcoded value
	logger.LoggerXds.Debug("Updating Enforcer Subscription Cache")
	label := commonEnforcerLabel
	subscriptionList := append(enforcerLabelMap[label].subscriptions, subscriptions)

	// TODO: (VirajSalaka) Decide if a map is required to keep version (just to avoid having the same version)
	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.SubscriptionListType: subscriptionList,
	})
	snap.Consistent()

	errSetSnap := enforcerSubscriptionCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].subscriptions = subscriptionList
	logger.LoggerXds.Infof("New Subscription cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApplications sets new update to the enforcer's Applications
func UpdateEnforcerApplications(applications *subscription.ApplicationList) {
	logger.LoggerXds.Debug("Updating Enforcer Application Cache")
	label := commonEnforcerLabel
	applicationList := append(enforcerLabelMap[label].applications, applications)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ApplicationListType: applicationList,
	})
	snap.Consistent()

	errSetSnap := enforcerApplicationCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].applications = applicationList
	logger.LoggerXds.Infof("New Application cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerAPIList sets new update to the enforcer's Apis
func UpdateEnforcerAPIList(label string, apis *subscription.APIList) {
	logger.LoggerXds.Debug("Updating Enforcer API Cache")
	apiList := append(enforcerLabelMap[label].apiList, apis)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.APIListType: apiList,
	})
	snap.Consistent()

	errSetSnap := enforcerAPICache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].apiList = apiList
	logger.LoggerXds.Infof("New API List cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApplicationPolicies sets new update to the enforcer's Application Policies
func UpdateEnforcerApplicationPolicies(applicationPolicies *subscription.ApplicationPolicyList) {
	logger.LoggerXds.Debug("Updating Enforcer Application Policy Cache")
	label := commonEnforcerLabel
	applicationPolicyList := append(enforcerLabelMap[label].applicationPolicies, applicationPolicies)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ApplicationPolicyListType: applicationPolicyList,
	})
	snap.Consistent()

	errSetSnap := enforcerApplicationPolicyCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].subscriptionPolicies = applicationPolicyList
	logger.LoggerXds.Infof("New Application Policy cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerSubscriptionPolicies sets new update to the enforcer's Subscription Policies
func UpdateEnforcerSubscriptionPolicies(subscriptionPolicies *subscription.SubscriptionPolicyList) {
	logger.LoggerXds.Debug("Updating Enforcer Subscription Policy Cache")
	label := commonEnforcerLabel
	subscriptionPolicyList := append(enforcerLabelMap[label].subscriptionPolicies, subscriptionPolicies)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.SubscriptionPolicyListType: subscriptionPolicyList,
	})
	snap.Consistent()

	errSetSnap := enforcerSubscriptionPolicyCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].subscriptionPolicies = subscriptionPolicyList
	logger.LoggerXds.Infof("New Subscription Policy cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApplicationKeyMappings sets new update to the enforcer's Application Key Mappings
func UpdateEnforcerApplicationKeyMappings(applicationKeyMappings *subscription.ApplicationKeyMappingList) {
	logger.LoggerXds.Debug("Updating Application Key Mapping Cache")
	label := commonEnforcerLabel
	applicationKeyMappingList := append(enforcerLabelMap[label].applicationKeyMappings, applicationKeyMappings)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ApplicationKeyMappingListType: applicationKeyMappingList,
	})
	snap.Consistent()

	errSetSnap := enforcerApplicationKeyMappingCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].applicationKeyMappings = applicationKeyMappingList
	logger.LoggerXds.Infof("New Application Key Mapping cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateXdsCacheWithLock uses mutex and lock to avoid different go routines updating XDS at the same time
func UpdateXdsCacheWithLock(label string, endpoints []types.Resource, clusters []types.Resource, routes []types.Resource,
	listeners []types.Resource) bool {
	mutexForXdsUpdate.Lock()
	defer mutexForXdsUpdate.Unlock()
	return updateXdsCache(label, endpoints, clusters, routes, listeners)
}

// GenerateIdentifierForAPI generates an identifier unique to the API
func GenerateIdentifierForAPI(vhost, name, version string) string {
	return fmt.Sprint(vhost, apiKeyFieldSeparator, name, apiKeyFieldSeparator, version)
}

// GenerateIdentifierForAPIWithUUID generates an identifier unique to the API
func GenerateIdentifierForAPIWithUUID(vhost, uuid string) string {
	return fmt.Sprint(vhost, apiKeyFieldSeparator, uuid)
}

// GenerateIdentifierForAPIWithoutVhost generates an identifier unique to the API name and version
func GenerateIdentifierForAPIWithoutVhost(name, version string) string {
	return fmt.Sprint(name, apiKeyFieldSeparator, version)
}

// GenerateHashedAPINameVersionIDWithoutVhost generates a hashed identifier unique to the API Name and Version
func GenerateHashedAPINameVersionIDWithoutVhost(name, version string) string {
	return generateHashValue(name, version)
}

func generateHashValue(apiName string, apiVersion string) string {
	apiNameVersionHash := sha1.New()
	apiNameVersionHash.Write([]byte(apiName + ":" + apiVersion))
	return hex.EncodeToString(apiNameVersionHash.Sum(nil)[:])
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

// GenerateAndUpdateKeyManagerList converts the data into KeyManager proto type
func GenerateAndUpdateKeyManagerList() {
	var keyManagerConfigList = make([]types.Resource, 0)
	for _, keyManager := range KeyManagerList {
		kmConfig := MarshalKeyManager(&keyManager)
		if kmConfig != nil {
			keyManagerConfigList = append(keyManagerConfigList, kmConfig)
		}
	}
	UpdateEnforcerKeyManagers(keyManagerConfigList)
}

// UpdateEnforcerKeyManagers Sets new update to the enforcer's configuration
func UpdateEnforcerKeyManagers(keyManagerConfigList []types.Resource) {
	logger.LoggerXds.Debug("Updating Key Manager Cache")
	label := commonEnforcerLabel

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.KeyManagerType: keyManagerConfigList,
	})
	snap.Consistent()

	errSetSnap := enforcerKeyManagerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].keyManagers = keyManagerConfigList
	logger.LoggerXds.Infof("New key manager cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerRevokedTokens method update the revoked tokens
// in the enforcer
func UpdateEnforcerRevokedTokens(revokedTokens []types.Resource) {
	logger.LoggerXds.Debug("Updating enforcer cache for revoked tokens")
	label := commonEnforcerLabel
	tokens := append(enforcerLabelMap[label].revokedTokens, revokedTokens...)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.RevokedTokensType: revokedTokens,
	})
	snap.Consistent()

	errSetSnap := enforcerRevokedTokensCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerLabelMap[label].revokedTokens = tokens
	logger.LoggerXds.Infof("New Revoked token cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateRateLimitXDSCache updates the xDS cache of the RateLimiter.
func UpdateRateLimitXDSCache(vHosts []string, adapterInternalAPI model.AdapterInternalAPI) {
	// Add Rate Limit inline policies in API to the cache
	rlsPolicyCache.AddAPILevelRateLimitPolicies(vHosts, &adapterInternalAPI)
}

// UpdateRateLimitXDSCacheForCustomPolicies updates the xDS cache of the RateLimiter for custom policies.
func UpdateRateLimitXDSCacheForCustomPolicies(gwLabel string, customRateLimitPolicies []*model.CustomRateLimitPolicy) {
	rlsPolicyCache.AddCustomRateLimitPolicies(customRateLimitPolicies)
	UpdateRateLimiterPolicies(gwLabel)
}

// UpdateAPICache updates the xDS cache related to the API Lifecycle event.
func UpdateAPICache(vHosts []string, newLabels []string, newlistenersForRoutes []string, adapterInternalAPI model.AdapterInternalAPI) error {
	mutexForInternalMapUpdate.Lock()
	defer mutexForInternalMapUpdate.Unlock()

	vHostIdentifier := GetvHostsIdentifier(adapterInternalAPI.UUID, adapterInternalAPI.EnvType)
	var oldvHosts []string
	if _, ok := orgIDAPIvHostsMap[adapterInternalAPI.OrganizationID]; ok {
		oldvHosts = orgIDAPIvHostsMap[adapterInternalAPI.GetOrganizationID()][vHostIdentifier]
		orgIDAPIvHostsMap[adapterInternalAPI.GetOrganizationID()][vHostIdentifier] = vHosts
	} else {
		vHostsMap := make(map[string][]string)
		vHostsMap[vHostIdentifier] = vHosts
		orgIDAPIvHostsMap[adapterInternalAPI.GetOrganizationID()] = vHostsMap
	}

	// Remove internal mappigs for old vHosts
	for _, oldvhost := range oldvHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(oldvhost, adapterInternalAPI.UUID)
		oldLabels := orgAPIMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier].envoyLabels
		delete(orgAPIMap[adapterInternalAPI.GetOrganizationID()], apiIdentifier)
		updateXdsCacheOnAPIChange(oldLabels, newLabels)
	}

	// Create internal mappigs for new vHosts
	for _, vHost := range vHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, adapterInternalAPI.UUID)
		var oldLabels []string
		if orgAPIMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier] != nil {
			oldLabels = orgAPIMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier].envoyLabels
		}

		routes, clusters, endpoints, err := oasParser.GetRoutesClustersEndpoints(adapterInternalAPI, nil,
			vHost, adapterInternalAPI.GetOrganizationID())

		if err != nil {
			return fmt.Errorf("error while deploying API. Name: %s Version: %s, OrgID: %s, Error: %s",
				adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), adapterInternalAPI.GetOrganizationID(), err.Error())
		}
		if _, ok := orgAPIMap[adapterInternalAPI.OrganizationID]; !ok {
			orgAPIMap[adapterInternalAPI.GetOrganizationID()] = make(map[string]*EnvoyInternalAPI)
		}
		orgAPIMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier] = &EnvoyInternalAPI{
			adapterInternalAPI: adapterInternalAPI,
			envoyLabels:        newLabels,
			routes:             routes,
			clusters:           clusters,
			endpointAddresses:  endpoints,
			enforcerAPI:        oasParser.GetEnforcerAPI(adapterInternalAPI, vHost),
		}

		if _, ok := listenerToRouteArrayMap[newlistenersForRoutes[0]]; ok {
			listenerToRouteArrayMap[newlistenersForRoutes[0]] = append(listenerToRouteArrayMap[newlistenersForRoutes[0]], routes...)
		} else {
			listenerToRouteArrayMap[newlistenersForRoutes[0]] = routes
		}

		revisionStatus := updateXdsCacheOnAPIChange(oldLabels, newLabels)
		logger.LoggerXds.Infof("Deployed Revision: %v:%v", apiIdentifier, revisionStatus)
	}
	return nil
}

// UpdateGatewayCache updates the xDS cache related to the Gateway Lifecycle event.
func UpdateGatewayCache(gateway *gwapiv1b1.Gateway, resolvedListenerCerts map[string]map[string][]byte,
	gwLuaScript string, customRateLimitPolicies []*model.CustomRateLimitPolicy) error {
	listener := oasParser.GetProductionListener(gateway, resolvedListenerCerts, gwLuaScript)
	gatewayLabelConfigMap[gateway.Name].listener = listener
	conf := config.ReadConfigs()
	if conf.Envoy.RateLimit.Enabled {
		gatewayLabelConfigMap[gateway.Name].customRateLimitPolicies = customRateLimitPolicies
	}
	return nil
}

// SanitizeGateway method sanitizes the gateway name
func SanitizeGateway(gatewayName string, create bool) error {
	if _, exists := enforcerLabelMap[gatewayName]; !exists && create {
		enforcerLabelMap[gatewayName] = &EnforcerInternalAPI{}
	} else if !exists {
		return fmt.Errorf("gateway %v does not exist in enforcerLabelMap", gatewayName)
	}
	if _, exists := gatewayLabelConfigMap[gatewayName]; !exists && create {
		gatewayLabelConfigMap[gatewayName] = &EnvoyGatewayConfig{}
	} else if !exists {
		return fmt.Errorf("gateway %v does not exist in gatewayLabelConfigMap", gatewayName)
	}
	return nil
}
