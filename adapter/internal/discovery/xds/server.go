/*
 *  Copyright (c) 2020, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
 */

// Package xds contains the implementation for the xds server cache updates
package xds

import (
	"context"
	"fmt"
	"math/rand"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoy_cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/subscription"
	wso2_cache "github.com/wso2/apk/adapter/pkg/discovery/protocol/cache/v3"
	wso2_resource "github.com/wso2/apk/adapter/pkg/discovery/protocol/resource/v3"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	enforcerCache                      wso2_cache.SnapshotCache
	enforcerSubscriptionCache          wso2_cache.SnapshotCache
	enforcerApplicationCache           wso2_cache.SnapshotCache
	enforcerApplicationKeyMappingCache wso2_cache.SnapshotCache

	// Common Enforcer Label as map key
	enforcerConfigMap                map[string][]types.Resource
	enforcerSubscriptionMap          map[string][]types.Resource
	enforcerApplicationMap           map[string][]types.Resource
	enforcerApplicationKeyMappingMap map[string][]types.Resource

	// KeyManagerList to store data
	isReady = false
)

var void struct{}

const (
	commonEnforcerLabel  string = "commonEnforcerLabel"
	maxRandomInt         int    = 999999999
	prototypedAPI        string = "PROTOTYPED"
	apiKeyFieldSeparator string = ":"
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
	enforcerCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerSubscriptionCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)
	enforcerApplicationKeyMappingCache = wso2_cache.NewSnapshotCache(false, IDHash{}, nil)

	enforcerConfigMap = make(map[string][]types.Resource)
	enforcerSubscriptionMap = make(map[string][]types.Resource)
	enforcerApplicationMap = make(map[string][]types.Resource)
	enforcerApplicationKeyMappingMap = make(map[string][]types.Resource)
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
		logger.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error creating new snapshot : %v", errNewSnap.Error()),
			Severity:  logging.MAJOR,
			ErrorCode: 1413,
		})
	}
	snap.Consistent()

	errSetSnap := enforcerCache.SetSnapshot(context.Background(), label, snap)
	if errSetSnap != nil {
		logger.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while setting the snapshot : %v", errSetSnap.Error()),
			Severity:  logging.MAJOR,
			ErrorCode: 1414,
		})
	}

	enforcerConfigMap[label] = configs
	logger.LoggerXds.Infof("New Config cache update for the label: " + label + " version: " + fmt.Sprint(version))
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
		logger.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while setting the snapshot : %v", errSetSnap.Error()),
			Severity:  logging.MAJOR,
			ErrorCode: 1414,
		})
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
		logger.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while setting the snapshot : %v", errSetSnap.Error()),
			Severity:  logging.MAJOR,
			ErrorCode: 1414,
		})
	}
	enforcerApplicationMap[label] = applicationList
	logger.LoggerXds.Infof("New Application cache update for the label: " + label + " version: " + fmt.Sprint(version))
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
		logger.LoggerXds.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while setting the snapshot : %v", errSetSnap.Error()),
			Severity:  logging.MAJOR,
			ErrorCode: 1414,
		})
	}
	enforcerApplicationKeyMappingMap[label] = applicationKeyMappingList
	logger.LoggerXds.Infof("New Application Key Mapping cache update for the label: " + label + " version: " + fmt.Sprint(version))
}
