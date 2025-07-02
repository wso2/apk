/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"k8s.io/apimachinery/pkg/types"
)

// RoutePolicyAndMetadataDataStore holds RoutePolicy and RouteMetadata objects,
// providing thread-safe access using a read-write mutex.
type RoutePolicyAndMetadataDataStore struct {
	routePolicies map[string]*dpv2alpha1.RoutePolicy
	routeMetadata map[string]*dpv2alpha1.RouteMetadata
	mu            sync.RWMutex
	commonControllerRestBaseURL string
	cfg *config.Server
}

// NewRoutePolicyAndMetadataDataStore initializes and returns a new datastore instance.
func NewRoutePolicyAndMetadataDataStore(cfg *config.Server) *RoutePolicyAndMetadataDataStore {
	return &RoutePolicyAndMetadataDataStore{
		routePolicies:              make(map[string]*dpv2alpha1.RoutePolicy),
		routeMetadata:              make(map[string]*dpv2alpha1.RouteMetadata),
		commonControllerRestBaseURL: "https://" + cfg.CommonControllerHost + ":" + cfg.CommonControllerRestPort,
		cfg:                        cfg,
	}
}

// AddRoutePolicy adds or updates a RoutePolicy in the datastore.
func (ds *RoutePolicyAndMetadataDataStore) AddRoutePolicy(policy *dpv2alpha1.RoutePolicy) {
	ds.cfg.Logger.Sugar().Debugf("Adding/Updating RoutePolicy: %s/%s", policy.Namespace, policy.Name)
	ds.mu.Lock()
	defer ds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Name:      policy.Name,
		Namespace: policy.Namespace,
	}
	ds.routePolicies[namespacedName.String()] = policy
}

// GetRoutePolicy returns a RoutePolicy by UUID. Returns nil if not found.
func (ds *RoutePolicyAndMetadataDataStore) GetRoutePolicy(namespacedName string) *dpv2alpha1.RoutePolicy {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if policy, exists := ds.routePolicies[namespacedName]; exists {
		return policy
	}
	return nil
}

// DeleteRoutePolicy deletes a RoutePolicy by UUID.
func (ds *RoutePolicyAndMetadataDataStore) DeleteRoutePolicy(namespacedName string) error {
	ds.cfg.Logger.Sugar().Debugf("Deleting RoutePolicy: %s", namespacedName)
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.routePolicies[namespacedName]; exists {
		delete(ds.routePolicies, namespacedName)
		return nil
	}
	return errors.New("route policy not found")
}

// AddRouteMetadata adds or updates a RouteMetadata in the datastore.
func (ds *RoutePolicyAndMetadataDataStore) AddRouteMetadata(metadata *dpv2alpha1.RouteMetadata) {
	ds.cfg.Logger.Sugar().Debugf("Adding/Updating RouteMetadata: %s/%s", metadata.Namespace, metadata.Name)
	ds.mu.Lock()
	defer ds.mu.Unlock()
	namespacedName := types.NamespacedName{
		Name:      metadata.Name,
		Namespace: metadata.Namespace,
	}
	ds.routeMetadata[namespacedName.String()] = metadata
}

// GetRouteMetadata returns a RouteMetadata by UUID. Returns nil if not found.
func (ds *RoutePolicyAndMetadataDataStore) GetRouteMetadata(namespacedName string) *dpv2alpha1.RouteMetadata {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if metadata, exists := ds.routeMetadata[namespacedName]; exists {
		return metadata
	}
	return nil
}

// DeleteRouteMetadata deletes a RouteMetadata by UUID.
func (ds *RoutePolicyAndMetadataDataStore) DeleteRouteMetadata(namespacedName string) error {
	ds.cfg.Logger.Sugar().Debugf("Deleting RouteMetadata: %s", namespacedName)
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.routeMetadata[namespacedName]; exists {
		delete(ds.routeMetadata, namespacedName)
		return nil
	}
	return errors.New("route metadata not found")
}

// LoadStartupData loads RoutePolicies and RouteMetadata from the common controller REST APIs.
func (ds *RoutePolicyAndMetadataDataStore) LoadStartupData() error {
	ds.cfg.Logger.Sugar().Debug("Loading RoutePolicies and RouteMetadata from common controller REST APIs")
	// Load Route Policies
	routePoliciesList, err := ds.loadAllRoutePolicies()
	if err != nil {
		return err
	}
	ds.cfg.Logger.Sugar().Debugf("Loaded %d RoutePolicies from common controller", len(routePoliciesList.Items))
	for _, policy := range routePoliciesList.Items {
		ds.cfg.Logger.Sugar().Debugf("Loading RoutePolicy: %s/%s", policy.Namespace, policy.Name)
		ds.AddRoutePolicy(&policy)
	}

	// Load Route Metadata
	routeMetadataList, err := ds.loadAllRouteMetadata()
	if err != nil {
		return err
	}
	ds.cfg.Logger.Sugar().Debugf("Loaded %d RouteMetadata from common controller", len(routeMetadataList.Items))
	for _, metadata := range routeMetadataList.Items {
		ds.cfg.Logger.Sugar().Debugf("Loading RouteMetadata: %s/%s", metadata.Namespace, metadata.Name)
		ds.AddRouteMetadata(&metadata)
	}
	return nil
}

// loadAllRoutePolicies fetches all RoutePolicies from the controller.
func (ds *RoutePolicyAndMetadataDataStore) loadAllRoutePolicies() (*dpv2alpha1.RoutePolicyList, error) {
	url := fmt.Sprintf("%s/routepolicies", ds.commonControllerRestBaseURL)
	tlsConfig, err := GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}

	resp, err := util.MakeGETRequest(url, tlsConfig, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dpv2alpha1.RoutePolicyList
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ds.cfg.Logger.Sugar().Debugf("loadAllRoutePolicies Response body: %s", string(body))
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// loadAllRouteMetadata fetches all RouteMetadata from the controller.
func (ds *RoutePolicyAndMetadataDataStore) loadAllRouteMetadata() (*dpv2alpha1.RouteMetadataList, error) {
	url := fmt.Sprintf("%s/routemetadata", ds.commonControllerRestBaseURL)
	tlsConfig, err := GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS config: %w", err)
	}

	resp, err := util.MakeGETRequest(url, tlsConfig, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dpv2alpha1.RouteMetadataList
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRoutePolicyCount returns the total number of stored route policies.
func (ds *RoutePolicyAndMetadataDataStore) GetRoutePolicyCount() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return len(ds.routePolicies)
}

// GetRouteMetadataCount returns the total number of stored route metadata.
func (ds *RoutePolicyAndMetadataDataStore) GetRouteMetadataCount() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return len(ds.routeMetadata)
}
