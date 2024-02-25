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
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
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
	"github.com/wso2/apk/adapter/internal/dataholder"
	"github.com/wso2/apk/adapter/internal/discovery/xds/common"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	oasParser "github.com/wso2/apk/adapter/internal/oasparser"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	operatorconsts "github.com/wso2/apk/adapter/internal/operator/constants"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	wso2_cache "github.com/wso2/apk/adapter/pkg/discovery/protocol/cache/v3"
	wso2_resource "github.com/wso2/apk/adapter/pkg/discovery/protocol/resource/v3"
	eventhubTypes "github.com/wso2/apk/adapter/pkg/eventhub/types"
	semantic_version "github.com/wso2/apk/adapter/pkg/semanticversion"
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
	listeners               []*listenerv3.Listener
	routeConfigs            map[string]*routev3.RouteConfiguration
	clusters                []*clusterv3.Cluster
	endpoints               []*corev3.Address
	customRateLimitPolicies []*model.CustomRateLimitPolicy
}

// EnforcerInternalAPI struct use to hold enforcer resources
type EnforcerInternalAPI struct {
	configs    []types.Resource
	apiList    []types.Resource
	jwtIssuers []types.Resource
}

var (
	// TODO: (VirajSalaka) Remove Unused mutexes.
	mutexForXdsUpdate         sync.Mutex
	mutexForInternalMapUpdate sync.Mutex

	cache                           envoy_cachev3.SnapshotCache
	enforcerCache                   wso2_cache.SnapshotCache
	enforcerJwtIssuerCache          wso2_cache.SnapshotCache
	enforcerAPICache                wso2_cache.SnapshotCache
	enforcerApplicationPolicyCache  wso2_cache.SnapshotCache
	enforcerSubscriptionPolicyCache wso2_cache.SnapshotCache
	enforcerKeyManagerCache         wso2_cache.SnapshotCache
	enforcerRevokedTokensCache      wso2_cache.SnapshotCache
	enforcerThrottleDataCache       wso2_cache.SnapshotCache

	orgAPIMap             map[string]map[string]*EnvoyInternalAPI // organizationID -> Vhost:API_UUID -> EnvoyInternalAPI struct map
	orgIDvHostBasepathMap map[string]map[string]string            // organizationID -> Vhost:basepath -> Vhost:API_UUID
	orgIDAPIvHostsMap     map[string]map[string][]string          // organizationID -> UUID -> prod/sand -> Envoy Vhost Array map

	orgIDLatestAPIVersionMap map[string]map[string]map[string]semantic_version.SemVersion // organizationID -> Vhost:APIName -> Version Range -> Latest API Version
	// Envoy Label as map key
	// TODO(amali) use this without generating all again.
	gatewayLabelConfigMap map[string]*EnvoyGatewayConfig // GW-Label -> EnvoyGatewayConfig struct map

	// Common Enforcer Label as map key
	// This doesn't have a usage yet. It will be used to handle multiple enforcer labels in future.
	enforcerLabelMap map[string]*EnforcerInternalAPI // Enforcer Label -> EnforcerInternalAPI struct map

	// KeyManagerList to store data
	KeyManagerList = make([]eventhubTypes.KeyManager, 0)
	isReady        = false
)

const (
	commonEnforcerLabel  string = "commonEnforcerLabel"
	maxRandomInt         int    = 999999999
	prototypedAPI        string = "PROTOTYPED"
	apiKeyFieldSeparator string = ":"
	gatewayController    string = "GatewayController"
	apiController        string = "APIController"
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
	enforcerAPICache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationPolicyCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerSubscriptionPolicyCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerKeyManagerCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerRevokedTokensCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerThrottleDataCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerJwtIssuerCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	gatewayLabelConfigMap = make(map[string]*EnvoyGatewayConfig)
	orgAPIMap = make(map[string]map[string]*EnvoyInternalAPI)
	orgIDAPIvHostsMap = make(map[string]map[string][]string) // organizationID -> UUID-prod/sand -> Envoy Vhost Array map
	orgIDvHostBasepathMap = make(map[string]map[string]string)
	orgIDLatestAPIVersionMap = make(map[string]map[string]map[string]semantic_version.SemVersion)

	enforcerLabelMap = make(map[string]*EnforcerInternalAPI)
	// currently subscriptions, configs, applications, applicationPolicies, subscriptionPolicies,
	// applicationKeyMappings, keyManagerConfigList, revokedTokens are supported with the hard coded label for Enforcer
	enforcerLabelMap[commonEnforcerLabel] = &EnforcerInternalAPI{}
	rand.Seed(time.Now().UnixNano())
	// go watchEnforcerResponse()
}

// GetXdsCache returns xds server cache.
func GetXdsCache() envoy_cachev3.SnapshotCache {
	return cache
}

// GetEnforcerCache returns xds server cache.
func GetEnforcerCache() wso2_cache.SnapshotCache {
	return enforcerCache
}

// GetEnforcerJWTIssuerCache returns xds server cache.
func GetEnforcerJWTIssuerCache() wso2_cache.SnapshotCache {
	return enforcerJwtIssuerCache
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
			logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1410, logging.MAJOR, "Error undeploying API %v with UUID %v of Organization %v from environments %v, error: %v",
				apiIdentifier, apiUUID, organizationID, labels, err.Error()))
			return err
		}
		// if no error, update internal vhost maps
		// error only happens when API not found in deleteAPI func
		logger.LoggerXds.Infof("Successfully undeployed the API %v with UUID %v under Organization %s and environment %s",
			apiIdentifier, apiUUID, organizationID, labels)
	}
	return nil
}

// deleteAPI deletes an API, its resources and updates the caches of given environments
func deleteAPI(apiIdentifier string, environments []string, organizationID string) error {
	apiUUID, _ := ExtractUUIDFromAPIIdentifier(apiIdentifier)
	var api *EnvoyInternalAPI

	if _, orgExists := orgAPIMap[organizationID]; orgExists {
		if oldAPI, apiExists := orgAPIMap[organizationID][apiIdentifier]; apiExists {
			api = oldAPI
		} else {
			logger.LoggerXds.Infof("Unable to delete API: %v from Organization: %v. API Does not exist. API_UUID: %v", apiIdentifier, organizationID, apiUUID)
			return errors.New(constants.NotFound)
		}

	} else {
		logger.LoggerXds.Infof("Unable to delete API: %v from Organization: %v. Organization Does not exist. API_UUID: %v", apiIdentifier, organizationID, apiUUID)
		return errors.New(constants.NotFound)
	}

	existingLabels := orgAPIMap[organizationID][apiIdentifier].envoyLabels
	toBeDelEnvs, toBeKeptEnvs := getEnvironmentsToBeDeleted(existingLabels, environments)

	if isSemanticVersioningEnabled(api.adapterInternalAPI.GetTitle(), api.adapterInternalAPI.GetVersion()) {
		updateRoutingRulesOnAPIDelete(organizationID, apiIdentifier, api.adapterInternalAPI)
	}

	var isAllowedToDelete bool
	updatedLabelsMap := make(map[string]struct{})
	for _, val := range toBeDelEnvs {
		updatedLabelsMap[val] = struct{}{}
		if stringutils.StringInSlice(val, existingLabels) {
			isAllowedToDelete = true
		}
	}
	if isAllowedToDelete {
		// do not delete from all environments, hence do not clear routes, clusters, endpoints, enforcerAPIs
		orgAPIMap[organizationID][apiIdentifier].envoyLabels = toBeKeptEnvs
		if len(toBeKeptEnvs) != 0 {
			UpdateXdsCacheOnAPIChange(updatedLabelsMap)
			return nil
		}
	}

	//clean maps of routes, clusters, endpoints, enforcerAPIs
	if len(environments) == 0 || isAllowedToDelete {
		cleanMapResources(apiIdentifier, organizationID, toBeDelEnvs)
	}
	UpdateXdsCacheOnAPIChange(updatedLabelsMap)
	return nil
}

func cleanMapResources(apiIdentifier string, organizationID string, toBeDelEnvs []string) {
	if _, orgExists := orgAPIMap[organizationID]; orgExists {
		delete(orgAPIMap[organizationID], apiIdentifier)
	}

	deleteBasepathForVHost(organizationID, apiIdentifier)
	//TODO: (SuKSW) clean any remaining in label wise maps, if this is the last API of that label
	logger.LoggerXds.Infof("Deleted API %v of organization %v", apiIdentifier, organizationID)
}

func deleteBasepathForVHost(organizationID, apiIdentifier string) {
	// Remove the basepath from map (that is used to avoid duplicate basepaths)
	if _, orgExists := orgAPIMap[organizationID]; orgExists {
		if oldOrgAPIAPI, ok := orgAPIMap[organizationID][apiIdentifier]; ok {
			s := strings.Split(apiIdentifier, apiKeyFieldSeparator)
			vHost := s[0]
			oldBasepath := oldOrgAPIAPI.adapterInternalAPI.GetXWso2Basepath()
			delete(orgIDvHostBasepathMap[organizationID], vHost+":"+oldBasepath)
		}
	}
}

// UpdateXdsCacheOnAPIChange when this method is called, openAPIEnvoy map is updated.
// Old labels refers to the previously assigned labels
// New labels refers to the the updated labels
func UpdateXdsCacheOnAPIChange(labels map[string]struct{}) bool {
	revisionStatus := false
	// TODO: (VirajSalaka) check possible optimizations, Since the number of labels are low by design it should not be an issue
	for newLabel := range labels {
		listeners, clusters, routes, endpoints, apis := GenerateEnvoyResoucesForGateway(newLabel)
		UpdateEnforcerApis(newLabel, apis, "")
		success := UpdateXdsCacheWithLock(newLabel, endpoints, clusters, routes, listeners)
		logger.LoggerXds.Debugf("Xds Cache is updated for the label : %v", newLabel)
		if success {
			// if even one label was updated with latest revision, we take the revision as deployed.
			// (other labels also will get updated successfully)
			revisionStatus = success
			continue
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
					logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1411, logging.MAJOR, "Error extracting vhost from API identifier: %v for Organization %v. Ignore deploying the API, error: %v", apiKey, organizationID, err))
					continue
				}
				isDefaultVersion := false
				var orgAPI *EnvoyInternalAPI
				// If the adapterInternalAPI is not found, proceed with other APIs. (Unreachable condition at this point)
				// If that happens, there is no purpose in processing clusters too.
				if org, ok := orgAPIMap[organizationID]; !ok {
					continue
				} else if orgAPI, ok = org[apiKey]; !ok {
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
	systemHost := conf.Envoy.SystemHost

	logger.LoggerXds.Debugf("System Host : %v", systemHost)

	// Add health endpoint
	routeHealth := envoyconf.CreateHealthEndpoint()
	vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], routeHealth)

	// Add the readiness endpoint. isReady flag will be set to true once all the apis are fetched from the control plane
	if isReady {
		readynessEndpoint := envoyconf.CreateReadyEndpoint()
		vhostToRouteArrayMap[systemHost] = append(vhostToRouteArrayMap[systemHost], readynessEndpoint)
	}

	envoyGatewayConfig, gwFound := gatewayLabelConfigMap[gatewayName]
	// gwFound means that the gateway is configured in the envoy config.
	listeners := envoyGatewayConfig.listeners
	if !gwFound || listeners == nil || len(listeners) == 0 {
		return nil, nil, nil, nil, nil
	}

	routeConfigs := make(map[string]*routev3.RouteConfiguration, 0)
	for _, route := range envoyGatewayConfig.routeConfigs {
		route.VirtualHosts = []*routev3.VirtualHost{}
	}
	// TODO(amali) Revisit the following
	// Find the matching listener for each vhost and then only add the routes to the routeConfigs
	for _, listener := range listeners {
		for vhost, routes := range vhostToRouteArrayMap {
			// todo(amali) without going through all this pain just to get the listener section name,
			// let the api decide which gateway section it refers to.
			// because it was already there in httproute cr
			listenerSection, found := common.FindElement(dataholder.GetAllGatewayListenerSections(),
				func(listenerSection gwapiv1b1.Listener) bool {
					if listenerSection.Hostname != nil && common.MatchesHostname(vhost, string(*listenerSection.Hostname)) {
						// if the envoy side vhost matches to a hostname in gateway, then it is a match
						if listener.Name == common.GetEnvoyListenerName(string(listenerSection.Protocol), uint32(listenerSection.Port)) {
							return true
						}
					}
					return false
				})
			if found {
				// Prepare the route config name based on the gateway listener section name.
				routeConfigName := common.GetEnvoyRouteConfigName(listener.Name, string(listenerSection.Name))
				routesConfig := oasParser.GetRouteConfigs(map[string][]*routev3.Route{vhost: routes}, routeConfigName, envoyGatewayConfig.customRateLimitPolicies)

				routeConfigMatched, alreadyExistsInRouteConfigList := routeConfigs[routeConfigName]
				if alreadyExistsInRouteConfigList {
					logger.LoggerAPKOperator.Debugf("Route already exists. %v", routeConfigName)
					routeConfigMatched.VirtualHosts = append(routeConfigMatched.VirtualHosts, routesConfig.VirtualHosts...)
				} else {
					logger.LoggerAPKOperator.Debugf("Route does not exist, Hence adding a new config. %v", routeConfigName)
					routeConfigs[routeConfigName] = routesConfig
				}
			} else {
				logger.LoggerAPKOperator.Errorf("Failed to find a matching gateway listener section in gateway CR for this vhost: %s in %v", vhost, listener.Name)
			}
		}
	}

	// Find gateway listeners that has $systemHost as its hostname and add the system routeConfig referencing those listeners
	gatewayListeners := dataholder.GetAllGatewayListenerSections()
	for _, listener := range gatewayListeners {
		if systemHost == string(*listener.Hostname) {
			var vhostToRouteArrayFilteredMapForSystemEndpoints = make(map[string][]*routev3.Route)
			vhostToRouteArrayFilteredMapForSystemEndpoints[systemHost] = vhostToRouteArrayMap[systemHost]
			routeConfigName := common.GetEnvoyRouteConfigName(common.GetEnvoyListenerName(string(listener.Protocol), uint32(listener.Port)), string(listener.Name))
			systemRoutesConfig := oasParser.GetRouteConfigs(vhostToRouteArrayFilteredMapForSystemEndpoints, routeConfigName, envoyGatewayConfig.customRateLimitPolicies)
			routeConfigs[routeConfigName] = systemRoutesConfig
		}
	}

	envoyGatewayConfig.routeConfigs = routeConfigs
	clusterArray = append(clusterArray, envoyGatewayConfig.clusters...)
	endpointArray = append(endpointArray, envoyGatewayConfig.endpoints...)
	generatedListeners, clusters, generatedRouteConfigs, endpoints := oasParser.GetCacheResources(endpointArray, clusterArray, listeners, routeConfigs)
	return generatedListeners, clusters, generatedRouteConfigs, endpoints, apis
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

	if _, ok := gatewayLabelConfigMap[label]; ok {
		gatewayLabelConfigMap[label].clusters = append(gatewayLabelConfigMap[label].clusters, clusters...)
		gatewayLabelConfigMap[label].endpoints = append(gatewayLabelConfigMap[label].endpoints, endpoints...)
	}
}

// use UpdateXdsCacheWithLock to avoid race conditions
func updateXdsCache(label string, endpoints []types.Resource, clusters []types.Resource, routes []types.Resource, listeners []types.Resource) bool {
	version, _ := crand.Int(crand.Reader, maxRandomBigInt())
	// TODO: (VirajSalaka) kept same version for all the resources as we are using simple cache implementation.
	// Will be updated once decide to move to incremental XDS
	snap, errNewSnap := envoy_cachev3.NewSnapshot(fmt.Sprint(version), map[envoy_resource.Type][]types.Resource{
		envoy_resource.EndpointType: endpoints,
		envoy_resource.ClusterType:  clusters,
		envoy_resource.ListenerType: listeners,
		envoy_resource.RouteType:    routes,
	})
	if errNewSnap != nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1413, logging.MAJOR, "Error creating new snapshot : %v", errNewSnap.Error()))
		return false
	}
	snap.Consistent()
	//TODO: (VirajSalaka) check
	errSetSnap := cache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1414, logging.MAJOR, "Error while setting the snapshot : %v", errSetSnap.Error()))
		return false
	}
	return true
}

// UpdateEnforcerConfig Sets new update to the enforcer's configuration
func UpdateEnforcerConfig(configFile *config.Config) {
	// TODO: (Praminda) handle labels
	label := commonEnforcerLabel
	configs := []types.Resource{MarshalConfig(configFile)}
	version, _ := crand.Int(crand.Reader, maxRandomBigInt())
	snap, errNewSnap := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ConfigType: configs,
	})
	if errNewSnap != nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1413, logging.MAJOR, "Error creating new snapshot : %v", errNewSnap.Error()))
	}
	snap.Consistent()

	errSetSnap := enforcerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1414, logging.MAJOR, "Error while setting the snapshot : %v", errSetSnap.Error()))
	}

	enforcerLabelMap[label].configs = configs
	logger.LoggerXds.Infof("New Config cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApis Sets new update to the enforcer's Apis
func UpdateEnforcerApis(label string, apis []types.Resource, version string) {

	if version == "" {
		version = fmt.Sprint(crand.Int(crand.Reader, maxRandomBigInt()))
	}

	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.APIType: apis,
	})
	snap.Consistent()

	errSetSnap := enforcerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1414, logging.MAJOR, "Error while setting the snapshot : %v", errSetSnap.Error()))
	}
	logger.LoggerXds.Infof("New API cache update for the label: " + label + " version: " + fmt.Sprint(version))

}

// UpdateEnforcerJWTIssuers sets new update to the enforcer's Applications
func UpdateEnforcerJWTIssuers(jwtIssuers *subscription.JWTIssuerList) {
	logger.LoggerXds.Debug("Updating Enforcer JWT Issuer Cache")
	label := commonEnforcerLabel
	jwtIssuerList := append(enforcerLabelMap[label].jwtIssuers, jwtIssuers)

	version, _ := crand.Int(crand.Reader, maxRandomBigInt())
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.JWTIssuerListType: jwtIssuerList,
	})
	snap.Consistent()

	errSetSnap := enforcerJwtIssuerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1414, logging.MAJOR, "Error while setting the snapshot : %v", errSetSnap.Error()))
	}
	enforcerLabelMap[label].jwtIssuers = jwtIssuerList
	logger.LoggerXds.Infof("New JWTIssuer cache update for the label: " + label + " version: " + fmt.Sprint(version))
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

// GenerateIdentifierForAPIWithoutVersion generates an identifier unique to the API despite of the version
func generateIdentifierForAPIWithoutVersion(vhost, name string) string {
	return fmt.Sprint(vhost, apiKeyFieldSeparator, name)
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

// ExtractUUIDFromAPIIdentifier extracts UUID from the API identifier
func ExtractUUIDFromAPIIdentifier(id string) (string, error) {
	elem := strings.Split(id, apiKeyFieldSeparator)
	if len(elem) == 2 {
		return elem[1], nil
	}
	err := fmt.Errorf("invalid API identifier: %v", id)
	return "", err
}

// RemoveAPICacheForEnv will remove all the internal mappings for a specific environment
func RemoveAPICacheForEnv(adapterInternalAPI model.AdapterInternalAPI, envType string) {
	vHostIdentifier := GetvHostsIdentifier(adapterInternalAPI.UUID, envType)
	var oldvHosts []string
	if _, ok := orgIDAPIvHostsMap[adapterInternalAPI.OrganizationID]; ok {
		oldvHosts = orgIDAPIvHostsMap[adapterInternalAPI.GetOrganizationID()][vHostIdentifier]
		for _, oldvhost := range oldvHosts {
			apiIdentifier := GenerateIdentifierForAPIWithUUID(oldvhost, adapterInternalAPI.UUID)
			if orgMap, orgExists := orgAPIMap[adapterInternalAPI.GetOrganizationID()]; orgExists {
				if _, apiExists := orgMap[apiIdentifier]; apiExists {
					delete(orgAPIMap[adapterInternalAPI.GetOrganizationID()], apiIdentifier)
				}
			}
		}
	}
}

// RemoveAPIFromOrgAPIMap removes api from orgAPI map
func RemoveAPIFromOrgAPIMap(uuid string, orgID string) {
	if orgMap, ok := orgAPIMap[orgID]; ok {
		for apiName := range orgMap {
			if strings.Contains(apiName, uuid) {
				delete(orgMap, apiName)
			}
		}
		if len(orgMap) == 0 {
			delete(orgAPIMap, orgID)
		}
	}
}

// UpdateAPICache updates the xDS cache related to the API Lifecycle event.
func UpdateAPICache(vHosts []string, newLabels []string, listener string, sectionName string,
	adapterInternalAPI model.AdapterInternalAPI) (map[string]struct{}, error) {
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

	updatedLabelsMap := make(map[string]struct{}, 0)

	// Remove internal mappings for old vHosts
	for _, oldvhost := range oldvHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(oldvhost, adapterInternalAPI.UUID)
		if orgMap, orgExists := orgAPIMap[adapterInternalAPI.GetOrganizationID()]; orgExists {
			if _, apiExists := orgMap[apiIdentifier]; apiExists {
				for _, oldLabel := range orgMap[apiIdentifier].envoyLabels {
					updatedLabelsMap[oldLabel] = struct{}{}
				}
				delete(orgAPIMap[adapterInternalAPI.GetOrganizationID()], apiIdentifier)
			}
		}
	}

	// Create internal mappings for new vHosts
	for _, vHost := range vHosts {
		logger.LoggerAPKOperator.Debugf("Creating internal mapping for vhost: %s", vHost)
		apiUUID := adapterInternalAPI.UUID
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, apiUUID)
		var orgExists bool

		// get changing label set
		if _, orgExists = orgAPIMap[adapterInternalAPI.GetOrganizationID()]; orgExists {
			if _, apiExists := orgAPIMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier]; apiExists {
				for _, oldLabel := range orgAPIMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier].envoyLabels {
					updatedLabelsMap[oldLabel] = struct{}{}
				}
			}
		}
		for _, newLabel := range newLabels {
			updatedLabelsMap[newLabel] = struct{}{}
		}

		routes, clusters, endpoints, err := oasParser.GetRoutesClustersEndpoints(&adapterInternalAPI, nil,
			vHost, adapterInternalAPI.GetOrganizationID())

		if err != nil {
			return nil, fmt.Errorf("error while deploying API. Name: %s Version: %s, OrgID: %s, API_UUID: %v, Error: %s",
				adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), adapterInternalAPI.GetOrganizationID(),
				apiUUID, err.Error())
		}
		if !orgExists {
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

		apiVersion := adapterInternalAPI.GetVersion()
		apiName := adapterInternalAPI.GetTitle()
		if isSemanticVersioningEnabled(apiName, apiVersion) {
			updateRoutingRulesOnAPIUpdate(adapterInternalAPI.OrganizationID, apiIdentifier, apiName, apiVersion, vHost)
		}
	}

	return updatedLabelsMap, nil
}

// UpdateGatewayCache updates the xDS cache related to the Gateway Lifecycle event.
func UpdateGatewayCache(gateway *gwapiv1b1.Gateway, resolvedListenerCerts map[string]map[string][]byte,
	gwLuaScript string, customRateLimitPolicies []*model.CustomRateLimitPolicy) error {
	listeners := oasParser.GetProductionListener(gateway, resolvedListenerCerts, gwLuaScript)
	gatewayLabelConfigMap[gateway.Name].listeners = listeners
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

// GetEnvoyGatewayConfigClusters method gets the number of clusters in envoy gateway config
func GetEnvoyGatewayConfigClusters() int {
	totalClusters := 0
	for _, config := range gatewayLabelConfigMap {
		// Add the number of clusters in this EnvoyGatewayConfig instance to the total
		totalClusters += len(config.clusters)
	}
	return totalClusters
}

// GetEnvoyInternalAPIRoutes method gets the number of routes in envoy internal API
func GetEnvoyInternalAPIRoutes() int {
	totalRoutes := 0
	for _, orgMap := range orgAPIMap {
		for _, api := range orgMap {
			// Add the number of routes in this EnvoyInternalAPI instance to the total
			totalRoutes += len(api.routes)
		}
	}
	return totalRoutes
}

// GetEnvoyInternalAPIClusters method gets the number of clusters in envoy internal API
func GetEnvoyInternalAPIClusters() int {
	totalClusters := 0
	for _, orgMap := range orgAPIMap {
		for _, api := range orgMap {
			// Add the number of clusters in this EnvoyInternalAPI instance to the total
			totalClusters += len(api.clusters)
		}
	}
	return totalClusters
}
