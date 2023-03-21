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
	apiModel "github.com/wso2/apk/adapter/internal/api/models"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	oasParser "github.com/wso2/apk/adapter/internal/oasparser"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/throttle"
	wso2_cache "github.com/wso2/apk/adapter/pkg/discovery/protocol/cache/v3"
	wso2_resource "github.com/wso2/apk/adapter/pkg/discovery/protocol/resource/v3"
	eventhubTypes "github.com/wso2/apk/adapter/pkg/eventhub/types"
	operatorconsts "github.com/wso2/apk/adapter/pkg/operator/constants"
	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

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

	// Vhosts entry maps, these maps updated with delta changes (when an API added, only added its entry only)
	// These maps are managed separately for API-CTL and APIM, since when deploying an project from API-CTL there is no API uuid
	apiUUIDToGatewayToVhosts map[string]map[string]string // API_UUID -> gateway-env -> vhost (for un-deploying APIs from APIM or Choreo)

	orgIDAPIMgwSwaggerMap       map[string]map[string]model.MgwSwagger     // organizationID -> Vhost:API_UUID -> MgwSwagger struct map
	orgIDAPIvHostsMap           map[string]map[string][]string             // organizationID -> UUID -> prod/sand -> Envoy Vhost Array map
	orgIDOpenAPIEnvoyMap        map[string]map[string][]string             // organizationID -> Vhost:API_UUID -> Envoy Label Array map
	orgIDOpenAPIRoutesMap       map[string]map[string][]*routev3.Route     // organizationID -> Vhost:API_UUID -> Envoy Routes map
	orgIDOpenAPIClustersMap     map[string]map[string][]*clusterv3.Cluster // organizationID -> Vhost:API_UUID -> Envoy Clusters map
	orgIDOpenAPIEndpointsMap    map[string]map[string][]*corev3.Address    // organizationID -> Vhost:API_UUID -> Envoy Endpoints map
	orgIDOpenAPIEnforcerApisMap map[string]map[string]types.Resource       // organizationID -> Vhost:API_UUID -> API Resource map
	orgIDvHostBasepathMap       map[string]map[string]string               // organizationID -> Vhost:basepath -> Vhost:API_UUID

	// Envoy Label as map key
	envoyListenerConfigMap     map[string][]*listenerv3.Listener        // GW-Label -> Listener Configuration map
	envoyRouteConfigMap        map[string][]*routev3.RouteConfiguration // GW-Label -> Routes Configuration map
	envoyClusterConfigMap      map[string][]*clusterv3.Cluster          // GW-Label -> Global Cluster Configuration map
	envoyEndpointConfigMap     map[string][]*corev3.Address             // GW-Label -> Global Endpoint Configuration map
	envoySystemListenerNameMap map[string]string                        // GW-Label -> System Listener Name map

	// Listener as map key
	listenerToRouteArrayMap map[string][]*routev3.Route // Listener -> Routes map

	// Common Enforcer Label as map key
	enforcerConfigMap                map[string][]types.Resource
	enforcerKeyManagerMap            map[string][]types.Resource
	enforcerSubscriptionMap          map[string][]types.Resource
	enforcerApplicationMap           map[string][]types.Resource
	enforcerAPIListMap               map[string][]types.Resource
	enforcerApplicationPolicyMap     map[string][]types.Resource
	enforcerSubscriptionPolicyMap    map[string][]types.Resource
	enforcerApplicationKeyMappingMap map[string][]types.Resource
	enforcerRevokedTokensMap         map[string][]types.Resource
	enforcerThrottleData             *throttle.ThrottleData

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

	apiUUIDToGatewayToVhosts = make(map[string]map[string]string)
	envoyListenerConfigMap = make(map[string][]*listenerv3.Listener)
	envoyRouteConfigMap = make(map[string][]*routev3.RouteConfiguration)
	envoyClusterConfigMap = make(map[string][]*clusterv3.Cluster)
	envoyEndpointConfigMap = make(map[string][]*corev3.Address)
	listenerToRouteArrayMap = make(map[string][]*routev3.Route)
	envoySystemListenerNameMap = make(map[string]string)

	orgIDAPIMgwSwaggerMap = make(map[string]map[string]model.MgwSwagger)       // organizationID -> Vhost:API_UUID -> MgwSwagger struct map
	orgIDAPIvHostsMap = make(map[string]map[string][]string)                   // organizationID -> UUID-prod/sand -> Envoy Vhost Array map
	orgIDOpenAPIEnvoyMap = make(map[string]map[string][]string)                // organizationID -> Vhost:API_UUID -> Envoy Label Array map
	orgIDOpenAPIRoutesMap = make(map[string]map[string][]*routev3.Route)       // organizationID -> Vhost:API_UUID -> Envoy Routes map
	orgIDOpenAPIClustersMap = make(map[string]map[string][]*clusterv3.Cluster) // organizationID -> Vhost:API_UUID -> Envoy Clusters map
	orgIDOpenAPIEndpointsMap = make(map[string]map[string][]*corev3.Address)   // organizationID -> Vhost:API_UUID -> Envoy Endpoints map
	orgIDOpenAPIEnforcerApisMap = make(map[string]map[string]types.Resource)   // organizationID -> Vhost:API_UUID -> API Resource map
	orgIDvHostBasepathMap = make(map[string]map[string]string)

	enforcerConfigMap = make(map[string][]types.Resource)
	enforcerKeyManagerMap = make(map[string][]types.Resource)
	enforcerSubscriptionMap = make(map[string][]types.Resource)
	enforcerApplicationMap = make(map[string][]types.Resource)
	enforcerAPIListMap = make(map[string][]types.Resource)
	enforcerApplicationPolicyMap = make(map[string][]types.Resource)
	enforcerSubscriptionPolicyMap = make(map[string][]types.Resource)
	enforcerApplicationKeyMappingMap = make(map[string][]types.Resource)
	enforcerRevokedTokensMap = make(map[string][]types.Resource)
	enforcerThrottleData = &throttle.ThrottleData{}
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
		for _, environment := range labels {
			// delete environment if exists
			delete(apiUUIDToGatewayToVhosts[apiUUID], environment)
		}
	}
	return nil
}

// deleteAPI deletes an API, its resources and updates the caches of given environments
func deleteAPI(apiIdentifier string, environments []string, organizationID string) error {
	_, exists := orgIDAPIMgwSwaggerMap[organizationID][apiIdentifier]
	if !exists {
		logger.LoggerXds.Infof("Unable to delete API: %v from Organization: %v. API Does not exist.", apiIdentifier, organizationID)
		return errors.New(constants.NotFound)
	}

	existingLabels := orgIDOpenAPIEnvoyMap[organizationID][apiIdentifier]
	toBeDelEnvs, toBeKeptEnvs := getEnvironmentsToBeDeleted(existingLabels, environments)

	for _, val := range toBeDelEnvs {
		isAllowedToDelete := stringutils.StringInSlice(val, existingLabels)
		if isAllowedToDelete {
			// do not delete from all environments, hence do not clear routes, clusters, endpoints, enforcerAPIs
			orgIDOpenAPIEnvoyMap[organizationID][apiIdentifier] = toBeKeptEnvs
			updateXdsCacheOnAPIChange(toBeDelEnvs, []string{})
			existingLabels = orgIDOpenAPIEnvoyMap[organizationID][apiIdentifier]
			if len(existingLabels) != 0 {
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
	delete(orgIDOpenAPIRoutesMap[organizationID], apiIdentifier)
	delete(orgIDOpenAPIClustersMap[organizationID], apiIdentifier)
	delete(orgIDOpenAPIEndpointsMap[organizationID], apiIdentifier)
	delete(orgIDOpenAPIEnforcerApisMap[organizationID], apiIdentifier)

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
	delete(orgIDOpenAPIEnvoyMap[organizationID], apiIdentifier)  //delete labels
	delete(orgIDAPIMgwSwaggerMap[organizationID], apiIdentifier) //delete mgwSwagger
	//TODO: (SuKSW) clean any remaining in label wise maps, if this is the last API of that label
	logger.LoggerXds.Infof("Deleted API %v of organization %v", apiIdentifier, organizationID)
}

func deleteBasepathForVHost(organizationID, apiIdentifier string) {
	// Remove the basepath from map (that is used to avoid duplicate basepaths)
	if oldMgwSwagger, ok := orgIDAPIMgwSwaggerMap[organizationID][apiIdentifier]; ok {
		s := strings.Split(apiIdentifier, apiKeyFieldSeparator)
		vHost := s[0]
		oldBasepath := oldMgwSwagger.GetXWso2Basepath()
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
		gateway := new(gwapiv1b1.Gateway)
		gateway.Name = newLabel
		listeners, clusters, routes, endpoints, apis := GenerateEnvoyResoucesForGateway(gateway, true, apiController)
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
			gateway := new(gwapiv1b1.Gateway)
			gateway.Name = oldLabel
			listeners, clusters, routes, endpoints, apis := GenerateEnvoyResoucesForGateway(gateway, true, apiController)
			UpdateEnforcerApis(oldLabel, apis, "")
			UpdateRateLimiterPolicies(oldLabel)
			UpdateXdsCacheWithLock(oldLabel, endpoints, clusters, routes, listeners)
			logger.LoggerXds.Debugf("Xds Cache is updated for the already existing label : %v", oldLabel)
		}
	}
	return revisionStatus
}

// GenerateEnvoyResoucesForGateway generates envoy resources for a given gateway
// This method will list out all APIs mapped to the label. and generate envoy resources for all of these APIs.
func GenerateEnvoyResoucesForGateway(gateway *gwapiv1b1.Gateway, isUpdate bool, flow string) ([]types.Resource,
	[]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	var clusterArray []*clusterv3.Cluster
	var vhostToRouteArrayMap = make(map[string][]*routev3.Route)
	var endpointArray []*corev3.Address
	var apis []types.Resource

	for organizationID, entityMap := range orgIDOpenAPIEnvoyMap {
		for apiKey, labels := range entityMap {
			if stringutils.StringInSlice(gateway.Name, labels) {
				vhost, err := ExtractVhostFromAPIIdentifier(apiKey)
				if err != nil {
					logger.LoggerXds.ErrorC(logging.GetErrorByCode(1411, err.Error(), organizationID))
					continue
				}
				isDefaultVersion := false
				if enforcerAPISwagger, ok := orgIDAPIMgwSwaggerMap[organizationID][apiKey]; ok {
					isDefaultVersion = enforcerAPISwagger.IsDefaultVersion
				} else {
					// If the mgwSwagger is not found, proceed with other APIs. (Unreachable condition at this point)
					// If that happens, there is no purpose in processing clusters too.
					continue
				}
				// If it is a default versioned API, the routes are added to the end of the existing array.
				// Otherwise the routes would be added to the front.
				// /fooContext/2.0.0/* resource path should be matched prior to the /fooContext/* .
				if isDefaultVersion {
					vhostToRouteArrayMap[vhost] = append(vhostToRouteArrayMap[vhost], orgIDOpenAPIRoutesMap[organizationID][apiKey]...)
				} else {
					vhostToRouteArrayMap[vhost] = append(orgIDOpenAPIRoutesMap[organizationID][apiKey], vhostToRouteArrayMap[vhost]...)
				}
				clusterArray = append(clusterArray, orgIDOpenAPIClustersMap[organizationID][apiKey]...)
				endpointArray = append(endpointArray, orgIDOpenAPIEndpointsMap[organizationID][apiKey]...)
				enfocerAPI, ok := orgIDOpenAPIEnforcerApisMap[organizationID][apiKey]
				if ok {
					apis = append(apis, enfocerAPI)
				}
				// listenerArrays = append(listenerArrays, openAPIListenersMap[apiKey])
			}
		}
	}

	// If the token endpoint is enabled, the token endpoint also needs to be added.
	conf := config.ReadConfigs()
	enableJwtIssuer := conf.Enforcer.JwtIssuer.Enabled
	systemHost := conf.Envoy.SystemHost
	logger.LoggerXds.Infof("System Host : %v", systemHost)
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

	var listenerArray []*listenerv3.Listener
	var routesConfig []*routev3.RouteConfiguration

	logger.LoggerXds.Debugf("Flow : %v", flow)

	if flow == gatewayController {
		listenerArray = oasParser.GetProductionListener(gateway)
		envoyListenerConfigMap[gateway.Name] = listenerArray

		for _, listenerObj := range gateway.Spec.Listeners {
			if gwapiv1b1.SectionName(*listenerObj.Hostname) == gwapiv1b1.SectionName(systemHost) {
				envoySystemListenerNameMap[gateway.Name] = string(listenerObj.Name)
			}
		}
		logger.LoggerXds.Debugf("systemListenerName : %v", envoySystemListenerNameMap[gateway.Name])

		for _, listener := range listenerArray {
			logger.LoggerXds.Debugf("Listener : %v", listener)
			routesFromListener := listenerToRouteArrayMap[listener.Name]
			logger.LoggerXds.Debugf("Routes from listener : %v", routesFromListener)
			var vhostToRouteArrayFilteredMap = make(map[string][]*routev3.Route)
			for vhost, routes := range vhostToRouteArrayMap {
				logger.LoggerXds.Debugf("Routes from Vhost Map : %v", routes)
				if (vhost == systemHost && listener.Name == envoySystemListenerNameMap[gateway.Name]) || checkRoutes(routes, routesFromListener) {
					logger.LoggerXds.Debugf("Equal routes : %v", routes)
					vhostToRouteArrayFilteredMap[vhost] = routes
				}
			}
			routesConfig = append(routesConfig, oasParser.GetRouteConfigs(vhostToRouteArrayFilteredMap, listener.Name))
			envoyRouteConfigMap[gateway.Name] = routesConfig
			logger.LoggerXds.Debugf("Listener : %v and routes %v", listener, routesConfig)
		}
	} else if flow == apiController {
		listenerArray = envoyListenerConfigMap[gateway.Name]
		for _, listener := range listenerArray {
			logger.LoggerXds.Debugf("Listener : %v", listener)
			routesFromListener := listenerToRouteArrayMap[listener.Name]
			logger.LoggerXds.Debugf("Routes from listener : %v", routesFromListener)
			var vhostToRouteArrayFilteredMap = make(map[string][]*routev3.Route)
			for vhost, routes := range vhostToRouteArrayMap {
				logger.LoggerXds.Debugf("Routes from Vhost Map : %v", routes)
				if (vhost == systemHost && listener.Name == envoySystemListenerNameMap[gateway.Name]) || checkRoutes(routes, routesFromListener) {
					logger.LoggerXds.Debugf("Equal routes : %v", routes)
					vhostToRouteArrayFilteredMap[vhost] = routes
				}
			}
			routesConfig = append(routesConfig, oasParser.GetRouteConfigs(vhostToRouteArrayFilteredMap, listener.Name))
			envoyRouteConfigMap[gateway.Name] = routesConfig
			logger.LoggerXds.Debugf("Listener : %v and routes %v", listener, routesConfig)
		}
	}

	logger.LoggerXds.Debugf("Routes Config : %v", routesConfig)
	clusterArray = append(clusterArray, envoyClusterConfigMap[gateway.Name]...)
	endpointArray = append(endpointArray, envoyEndpointConfigMap[gateway.Name]...)
	endpoints, clusters, listeners, routeConfigs := oasParser.GetCacheResources(endpointArray, clusterArray, listenerArray, routesConfig)
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
	envoyClusterConfigMap[label] = clusters
	envoyEndpointConfigMap[label] = endpoints
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

	enforcerConfigMap[label] = configs
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
}

// UpdateEnforcerSubscriptions sets new update to the enforcer's Subscriptions
func UpdateEnforcerSubscriptions(subscriptions *subscription.SubscriptionList) {
	//TODO: (Dinusha) check this hardcoded value
	logger.LoggerXds.Debug("Updating Enforcer Subscription Cache")
	label := commonEnforcerLabel
	subscriptionList := enforcerSubscriptionMap[label]
	subscriptionList = append(subscriptionList, subscriptions)

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
	enforcerSubscriptionMap[label] = subscriptionList
	logger.LoggerXds.Infof("New Subscription cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApplications sets new update to the enforcer's Applications
func UpdateEnforcerApplications(applications *subscription.ApplicationList) {
	logger.LoggerXds.Debug("Updating Enforcer Application Cache")
	label := commonEnforcerLabel
	applicationList := enforcerApplicationMap[label]
	applicationList = append(applicationList, applications)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ApplicationListType: applicationList,
	})
	snap.Consistent()

	errSetSnap := enforcerApplicationCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerApplicationMap[label] = applicationList
	logger.LoggerXds.Infof("New Application cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerAPIList sets new update to the enforcer's Apis
func UpdateEnforcerAPIList(label string, apis *subscription.APIList) {
	logger.LoggerXds.Debug("Updating Enforcer API Cache")
	apiList := enforcerAPIListMap[label]
	apiList = append(apiList, apis)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.APIListType: apiList,
	})
	snap.Consistent()

	errSetSnap := enforcerAPICache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerAPIListMap[label] = apiList
	logger.LoggerXds.Infof("New API List cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApplicationPolicies sets new update to the enforcer's Application Policies
func UpdateEnforcerApplicationPolicies(applicationPolicies *subscription.ApplicationPolicyList) {
	logger.LoggerXds.Debug("Updating Enforcer Application Policy Cache")
	label := commonEnforcerLabel
	applicationPolicyList := enforcerApplicationPolicyMap[label]
	applicationPolicyList = append(applicationPolicyList, applicationPolicies)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ApplicationPolicyListType: applicationPolicyList,
	})
	snap.Consistent()

	errSetSnap := enforcerApplicationPolicyCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerApplicationPolicyMap[label] = applicationPolicyList
	logger.LoggerXds.Infof("New Application Policy cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerSubscriptionPolicies sets new update to the enforcer's Subscription Policies
func UpdateEnforcerSubscriptionPolicies(subscriptionPolicies *subscription.SubscriptionPolicyList) {
	logger.LoggerXds.Debug("Updating Enforcer Subscription Policy Cache")
	label := commonEnforcerLabel
	subscriptionPolicyList := enforcerSubscriptionPolicyMap[label]
	subscriptionPolicyList = append(subscriptionPolicyList, subscriptionPolicies)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.SubscriptionPolicyListType: subscriptionPolicyList,
	})
	snap.Consistent()

	errSetSnap := enforcerSubscriptionPolicyCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerSubscriptionPolicyMap[label] = subscriptionPolicyList
	logger.LoggerXds.Infof("New Subscription Policy cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerApplicationKeyMappings sets new update to the enforcer's Application Key Mappings
func UpdateEnforcerApplicationKeyMappings(applicationKeyMappings *subscription.ApplicationKeyMappingList) {
	logger.LoggerXds.Debug("Updating Application Key Mapping Cache")
	label := commonEnforcerLabel
	applicationKeyMappingList := enforcerApplicationKeyMappingMap[label]
	applicationKeyMappingList = append(applicationKeyMappingList, applicationKeyMappings)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ApplicationKeyMappingListType: applicationKeyMappingList,
	})
	snap.Consistent()

	errSetSnap := enforcerApplicationKeyMappingCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerApplicationKeyMappingMap[label] = applicationKeyMappingList
	logger.LoggerXds.Infof("New Application Key Mapping cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateXdsCacheWithLock uses mutex and lock to avoid different go routines updating XDS at the same time
func UpdateXdsCacheWithLock(label string, endpoints []types.Resource, clusters []types.Resource, routes []types.Resource,
	listeners []types.Resource) bool {
	mutexForXdsUpdate.Lock()
	defer mutexForXdsUpdate.Unlock()
	return updateXdsCache(label, endpoints, clusters, routes, listeners)
}

// ListApis returns a list of objects that holds info about each API
func ListApis(apiType string, organizationID string, limit *int64) *apiModel.APIMeta {
	var limitValue int
	if limit == nil {
		limitValue = len(orgIDAPIMgwSwaggerMap[organizationID])
	} else {
		limitValue = int(*limit)
	}
	var apisArray []*apiModel.APIMetaListItem
	i := 0
	for apiIdentifier, mgwSwagger := range orgIDAPIMgwSwaggerMap[organizationID] {
		if i == limitValue {
			break
		}
		if apiType == "" || mgwSwagger.GetAPIType() == apiType {
			var apiMetaListItem apiModel.APIMetaListItem
			apiMetaListItem.APIName = mgwSwagger.GetTitle()
			apiMetaListItem.Version = mgwSwagger.GetVersion()
			apiMetaListItem.APIType = mgwSwagger.GetAPIType()
			apiMetaListItem.Context = mgwSwagger.GetXWso2Basepath()
			apiMetaListItem.GatewayEnvs = orgIDOpenAPIEnvoyMap[organizationID][apiIdentifier]
			vhost := "ERROR"
			if vh, err := ExtractVhostFromAPIIdentifier(apiIdentifier); err == nil {
				vhost = vh
			}
			apiMetaListItem.Vhosts = []string{vhost}
			// orgIDAPIvHostsMap[organizationID][apiIdentifier]
			apisArray = append(apisArray, &apiMetaListItem)
			i++
		}
	}
	var apiMetaObject apiModel.APIMeta
	apiMetaObject.Total = int64(len(orgIDAPIMgwSwaggerMap[organizationID]))
	apiMetaObject.Count = int64(len(apisArray))
	apiMetaObject.List = apisArray
	return &apiMetaObject
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
	enforcerKeyManagerMap[label] = keyManagerConfigList
	logger.LoggerXds.Infof("New key manager cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerRevokedTokens method update the revoked tokens
// in the enforcer
func UpdateEnforcerRevokedTokens(revokedTokens []types.Resource) {
	logger.LoggerXds.Debug("Updating enforcer cache for revoked tokens")
	label := commonEnforcerLabel
	tokens := enforcerRevokedTokensMap[label]
	tokens = append(tokens, revokedTokens...)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.RevokedTokensType: revokedTokens,
	})
	snap.Consistent()

	errSetSnap := enforcerRevokedTokensCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.GetErrorByCode(1414, errSetSnap.Error()))
	}
	enforcerRevokedTokensMap[label] = tokens
	logger.LoggerXds.Infof("New Revoked token cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateEnforcerThrottleData update the key template and blocking conditions
// data in the enforcer
func UpdateEnforcerThrottleData(throttleData *throttle.ThrottleData) {
	logger.LoggerXds.Debug("Updating enforcer cache for throttle data")
	label := commonEnforcerLabel
	var data []types.Resource

	// Set new throttle data content based on the already available content in the cache DTO
	// and the new data being requested to add.
	// ex: keytemplates being pressent in the `throttleData` means this method was called
	// after downloading key templates. That means we should populate keytemplates property
	// in the cache DTO, keeping the other properties as it is. This is done this way to avoid
	// the need of two xds services to push keytemplates and blocking conditions.
	templates := throttleData.KeyTemplates
	conditions := throttleData.BlockingConditions
	ipConditions := throttleData.IpBlockingConditions
	if templates == nil {
		templates = enforcerThrottleData.KeyTemplates
	}
	if conditions == nil {
		conditions = enforcerThrottleData.BlockingConditions
	}
	if ipConditions == nil {
		ipConditions = enforcerThrottleData.IpBlockingConditions
	}

	t := &throttle.ThrottleData{
		KeyTemplates:         templates,
		BlockingConditions:   conditions,
		IpBlockingConditions: ipConditions,
	}
	data = append(data, t)

	version := rand.Intn(maxRandomInt)
	snap, _ := wso2_cache.NewSnapshot(fmt.Sprint(version), map[wso2_resource.Type][]types.Resource{
		wso2_resource.ThrottleDataType: data,
	})
	snap.Consistent()

	err := enforcerThrottleDataCache.SetSnapshot(context.Background(), label, snap)
	if err != nil {
		logger.LoggerXds.Error(err)
	}
	enforcerThrottleData = t
	logger.LoggerXds.Infof("New Throttle Data cache update for the label: " + label + " version: " + fmt.Sprint(version))
}

// UpdateRateLimitXDSCache updates the xDS cache of the RateLimiter.
func UpdateRateLimitXDSCache(vHosts []string, mgwSwagger model.MgwSwagger) {
	// Add Rate Limit inline policies in API to the cache
	rlsPolicyCache.AddAPILevelRateLimitPolicies(vHosts, &mgwSwagger)
}

// UpdateAPICache updates the xDS cache related to the API Lifecycle event.
func UpdateAPICache(vHosts []string, newLabels []string, newlistenersForRoutes []string, mgwSwagger model.MgwSwagger) error {
	mutexForInternalMapUpdate.Lock()
	defer mutexForInternalMapUpdate.Unlock()

	vHostIdentifier := GetvHostsIdentifier(mgwSwagger.UUID, mgwSwagger.EnvType)
	var oldvHosts []string
	if _, ok := orgIDAPIvHostsMap[mgwSwagger.OrganizationID]; ok {
		oldvHosts = orgIDAPIvHostsMap[mgwSwagger.GetOrganizationID()][vHostIdentifier]
		orgIDAPIvHostsMap[mgwSwagger.GetOrganizationID()][vHostIdentifier] = vHosts
	} else {
		vHostsMap := make(map[string][]string)
		vHostsMap[vHostIdentifier] = vHosts
		orgIDAPIvHostsMap[mgwSwagger.GetOrganizationID()] = vHostsMap
	}

	// Remove internal mappigs for old vHosts
	for _, oldvhost := range oldvHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(oldvhost, mgwSwagger.UUID)
		delete(orgIDAPIMgwSwaggerMap[mgwSwagger.GetOrganizationID()], apiIdentifier)
		delete(orgIDOpenAPIRoutesMap[mgwSwagger.GetOrganizationID()], apiIdentifier)
		delete(orgIDOpenAPIClustersMap[mgwSwagger.GetOrganizationID()], apiIdentifier)
		delete(orgIDOpenAPIEndpointsMap[mgwSwagger.GetOrganizationID()], apiIdentifier)
		delete(orgIDOpenAPIEnforcerApisMap[mgwSwagger.GetOrganizationID()], apiIdentifier)
		oldLabels := orgIDOpenAPIEnvoyMap[mgwSwagger.GetOrganizationID()][apiIdentifier]
		updateXdsCacheOnAPIChange(oldLabels, newLabels)
	}

	// Create internal mappigs for new vHosts
	for _, vHost := range vHosts {
		apiIdentifier := GenerateIdentifierForAPIWithUUID(vHost, mgwSwagger.UUID)
		oldLabels := orgIDOpenAPIEnvoyMap[mgwSwagger.GetOrganizationID()][apiIdentifier]

		if _, ok := orgIDAPIMgwSwaggerMap[mgwSwagger.OrganizationID]; ok {
			orgIDAPIMgwSwaggerMap[mgwSwagger.GetOrganizationID()][apiIdentifier] = mgwSwagger
		} else {
			mgwSwaggerMap := make(map[string]model.MgwSwagger)
			mgwSwaggerMap[apiIdentifier] = mgwSwagger
			orgIDAPIMgwSwaggerMap[mgwSwagger.GetOrganizationID()] = mgwSwaggerMap
		}

		if _, ok := orgIDOpenAPIEnvoyMap[mgwSwagger.GetOrganizationID()]; ok {
			orgIDOpenAPIEnvoyMap[mgwSwagger.GetOrganizationID()][apiIdentifier] = newLabels
		} else {
			openAPIEnvoyMap := make(map[string][]string)
			openAPIEnvoyMap[apiIdentifier] = newLabels
			orgIDOpenAPIEnvoyMap[mgwSwagger.GetOrganizationID()] = openAPIEnvoyMap
		}

		routes, clusters, endpoints, err := oasParser.GetRoutesClustersEndpoints(mgwSwagger, nil,
			vHost, mgwSwagger.GetOrganizationID())

		if err != nil {
			return fmt.Errorf("error while deploying API. Name: %s Version: %s, OrgID: %s, Error: %s",
				mgwSwagger.GetTitle(), mgwSwagger.GetVersion(), mgwSwagger.GetOrganizationID(), err.Error())
		}

		if _, ok := orgIDOpenAPIRoutesMap[mgwSwagger.GetOrganizationID()]; ok {
			orgIDOpenAPIRoutesMap[mgwSwagger.GetOrganizationID()][apiIdentifier] = routes
		} else {
			routesMap := make(map[string][]*routev3.Route)
			routesMap[apiIdentifier] = routes
			orgIDOpenAPIRoutesMap[mgwSwagger.GetOrganizationID()] = routesMap
		}

		if _, ok := listenerToRouteArrayMap[newlistenersForRoutes[0]]; ok {
			listenerToRouteArrayMap[newlistenersForRoutes[0]] = append(listenerToRouteArrayMap[newlistenersForRoutes[0]], routes...)
		} else {
			listenerToRouteArrayMap[newlistenersForRoutes[0]] = routes
		}

		if _, ok := orgIDOpenAPIClustersMap[mgwSwagger.GetOrganizationID()]; ok {
			orgIDOpenAPIClustersMap[mgwSwagger.GetOrganizationID()][apiIdentifier] = clusters
		} else {
			clustersMap := make(map[string][]*clusterv3.Cluster)
			clustersMap[apiIdentifier] = clusters
			orgIDOpenAPIClustersMap[mgwSwagger.GetOrganizationID()] = clustersMap
		}

		if _, ok := orgIDOpenAPIEndpointsMap[mgwSwagger.GetOrganizationID()]; ok {
			orgIDOpenAPIEndpointsMap[mgwSwagger.GetOrganizationID()][apiIdentifier] = endpoints
		} else {
			endpointMap := make(map[string][]*corev3.Address)
			endpointMap[apiIdentifier] = endpoints
			orgIDOpenAPIEndpointsMap[mgwSwagger.GetOrganizationID()] = endpointMap
		}

		if _, ok := orgIDOpenAPIEnforcerApisMap[mgwSwagger.GetOrganizationID()]; ok {
			orgIDOpenAPIEnforcerApisMap[mgwSwagger.GetOrganizationID()][apiIdentifier] = oasParser.GetEnforcerAPI(mgwSwagger, vHost)
		} else {
			enforcerAPIMap := make(map[string]types.Resource)
			enforcerAPIMap[apiIdentifier] = oasParser.GetEnforcerAPI(mgwSwagger, vHost)
			orgIDOpenAPIEnforcerApisMap[mgwSwagger.GetOrganizationID()] = enforcerAPIMap
		}

		revisionStatus := updateXdsCacheOnAPIChange(oldLabels, newLabels)
		logger.LoggerXds.Infof("Deployed Revision: %v:%v", apiIdentifier, revisionStatus)
	}
	return nil
}
