/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

// This file contains code derived from Contour,
// https://github.com/projectcontour/contour
// from the source file
// https://github.com/projectcontour/contour/blob/main/internal/xds/v3/snapshotter.go
// and is provided here subject to the following:
// Copyright Project Contour Authors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/types"
)

var Hash = cachev3.IDHash{}

// SnapshotCacheWithCallbacks uses the go-control-plane SimpleCache to store snapshots of
// Envoy resources, sliced by Node ID so that we can do incremental xDS properly.
// It does this by also implementing callbacks to make sure that the cache is kept
// up to date for each new node.
//
// Having the cache also implement the callbacks is a little bit hacky, but it makes sure
// that all the required bookkeeping happens.
// TODO(youngnick): Talk to the go-control-plane maintainers and see if we can upstream
// this in a better way.
type SnapshotCacheWithCallbacks interface {
	cachev3.SnapshotCache
	serverv3.Callbacks
	GenerateNewSnapshot(string, types.XdsResources) error
}

type snapshotMap map[string]*cachev3.Snapshot

type nodeInfoMap map[int64]*corev3.Node

type snapshotCache struct {
	cachev3.SnapshotCache
	streamIDNodeInfo nodeInfoMap
	snapshotVersion  int64
	lastSnapshot     snapshotMap
	mu               sync.Mutex
}

// GenerateNewSnapshot takes a table of resources (the output from the IR->xDS
// translator) and updates the snapshot version.
func (s *snapshotCache) GenerateNewSnapshot(irKey string, resources types.XdsResources) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	version := s.newSnapshotVersion()

	// Create a snapshot with all xDS resources.
	snapshot, err := cachev3.NewSnapshot(
		version,
		resources,
	)
	if err != nil {
		return err
	}

	s.lastSnapshot[irKey] = snapshot

	for _, node := range s.getNodeIDs(irKey) {
		loggers.LoggerAPKOperator.Debugf("Generating a snapshot with Node %s", node)
		err := s.SetSnapshot(context.TODO(), node, snapshot)
		if err != nil {
			return err
		}
	}

	return nil
}

// newSnapshotVersion increments the current snapshotVersion
// and returns as a string.
func (s *snapshotCache) newSnapshotVersion() string {
	// Reset the snapshotVersion if it ever hits max size.
	if s.snapshotVersion == math.MaxInt64 {
		s.snapshotVersion = 0
	}

	// Increment the snapshot version & return as string.
	s.snapshotVersion++
	return strconv.FormatInt(s.snapshotVersion, 10)
}

// NewSnapshotCache gives you a fresh SnapshotCache.
// It needs a logger that supports the go-control-plane
// required interface (Debugf, Infof, Warnf, and Errorf).
func NewSnapshotCache(ads bool) SnapshotCacheWithCallbacks {
	// // Set up the nasty wrapper hack.
	// wrappedLogger := logger.Sugar()
	//todo(amali) check the logger passing
	return &snapshotCache{
		SnapshotCache:    cachev3.NewSnapshotCache(ads, &Hash, loggers.LoggerAPKOperator),
		lastSnapshot:     make(snapshotMap),
		streamIDNodeInfo: make(nodeInfoMap),
	}
}

// getNodeIDs retrieves the node ids from the node info map whose
// cluster field matches the ir key
func (s *snapshotCache) getNodeIDs(irKey string) []string {
	var nodeIDs []string
	for _, node := range s.streamIDNodeInfo {
		if node != nil && node.Cluster == irKey {
			nodeIDs = append(nodeIDs, node.Id)
		}
	}

	return nodeIDs
}

// OnStreamOpen and the other OnStream* functions implement the callbacks for the
// state-of-the-world stream types.
func (s *snapshotCache) OnStreamOpen(_ context.Context, streamID int64, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streamIDNodeInfo[streamID] = nil

	return nil
}

func (s *snapshotCache) OnStreamClosed(streamID int64, _ *corev3.Node) {
	// TODO: something with the node?
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.streamIDNodeInfo, streamID)
}

func (s *snapshotCache) OnStreamRequest(streamID int64, req *discoveryv3.DiscoveryRequest) error {
	s.mu.Lock()
	// We could do this a little earlier than the defer, since the last half of this func is only logging
	// but that seemed like a premature optimization.
	defer s.mu.Unlock()

	// It's possible that only the first discovery request will have a node ID set.
	// We also need to save the node ID to the node list anyway.
	// So check if we have a nodeID for this stream already, then set it if not.
	if s.streamIDNodeInfo[streamID] == nil {
		if req.Node.Id == "" {
			return fmt.Errorf("couldn't get the node ID from the first discovery request on stream %d", streamID)
		}
		loggers.LoggerAPKOperator.Debugf("First discovery request on stream %d, got nodeID %s", streamID, req.Node.Id)
		s.streamIDNodeInfo[streamID] = req.Node
	}
	nodeID := s.streamIDNodeInfo[streamID].Id
	cluster := s.streamIDNodeInfo[streamID].Cluster

	var nodeVersion string

	var errorCode int32
	var errorMessage string

	// If no snapshot has been generated yet, we can't do anything, so don't mess with this request.
	// go-control-plane will respond with an empty response, then send an update when a snapshot is generated.
	if s.lastSnapshot[cluster] == nil {
		return nil
	}

	_, err := s.GetSnapshot(nodeID)
	if err != nil {
		err = s.SetSnapshot(context.TODO(), nodeID, s.lastSnapshot[cluster])
		if err != nil {
			return err
		}
	}

	if req.Node != nil {
		if bv := req.Node.GetUserAgentBuildVersion(); bv != nil && bv.Version != nil {
			nodeVersion = fmt.Sprintf("v%d.%d.%d", bv.Version.MajorNumber, bv.Version.MinorNumber, bv.Version.Patch)
		}
	}

	loggers.LoggerAPKOperator.Debugf("Got a new request, version_info %s, response_nonce %s, nodeID %s, node_version %s", req.VersionInfo, req.ResponseNonce, nodeID, nodeVersion)

	if status := req.ErrorDetail; status != nil {
		// if Envoy rejected the last update log the details here.
		// TODO(youngnick): Handle NACK properly
		errorCode = status.Code
		errorMessage = status.Message
	}

	loggers.LoggerAPKOperator.Debugf("handling v3 xDS resource request, version_info %s, response_nonce %s, nodeID %s, node_version %s, resource_names %v, type_url %s, errorCode %d, errorMessage %s",
		req.VersionInfo, req.ResponseNonce,
		nodeID, nodeVersion, req.ResourceNames, req.GetTypeUrl(),
		errorCode, errorMessage)

	return nil
}

func (s *snapshotCache) OnStreamResponse(_ context.Context, streamID int64, _ *discoveryv3.DiscoveryRequest, _ *discoveryv3.DiscoveryResponse) {
	// No mutex lock required here because no writing to the cache.
	node := s.streamIDNodeInfo[streamID]
	if node == nil {
		loggers.LoggerAPKOperator.Errorf("Tried to send a response to a node we haven't seen yet on stream %d", streamID)
	} else {
		loggers.LoggerAPKOperator.Debugf("Sending Response on stream %d to node %s", streamID, node.Id)
	}
}

// OnDeltaStreamOpen and the other OnDeltaStream*/OnStreamDelta* functions implement
// the callbacks for the incremental xDS versions.
// Yes, the different ordering in the name is part of the go-control-plane interface.
func (s *snapshotCache) OnDeltaStreamOpen(_ context.Context, streamID int64, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Ensure that we're adding the streamID to the Node ID list.
	s.streamIDNodeInfo[streamID] = nil

	return nil
}

func (s *snapshotCache) OnDeltaStreamClosed(streamID int64, _ *corev3.Node) {
	// TODO: something with the node?
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.streamIDNodeInfo, streamID)
}

func (s *snapshotCache) OnStreamDeltaRequest(streamID int64, req *discoveryv3.DeltaDiscoveryRequest) error {
	s.mu.Lock()
	// We could do this a little earlier than with a defer, since the last half of this func is logging
	// but that seemed like a premature optimization.
	defer s.mu.Unlock()

	var nodeVersion string
	var errorCode int32
	var errorMessage string
	// It's possible that only the first incremental discovery request will have a node ID set.
	// We also need to save the node ID to the node list anyway.
	// So check if we have a nodeID for this stream already, then set it if not.
	node := s.streamIDNodeInfo[streamID]
	if node == nil {
		if req.Node.Id == "" {
			return fmt.Errorf("couldn't get the node ID from the first incremental discovery request on stream %d", streamID)
		}
		loggers.LoggerAPKOperator.Debugf("First incremental discovery request on stream %d, got nodeID %s", streamID, req.Node.Id)
		s.streamIDNodeInfo[streamID] = req.Node
	}
	nodeID := s.streamIDNodeInfo[streamID].Id
	cluster := s.streamIDNodeInfo[streamID].Cluster

	// If no snapshot has been written into the snapshotCache yet, we can't do anything, so don't mess with
	// this request. go-control-plane will respond with an empty response, then send an update when a
	// snapshot is generated.
	if s.lastSnapshot[cluster] == nil {
		return nil
	}

	_, err := s.GetSnapshot(nodeID)
	if err != nil {
		err = s.SetSnapshot(context.TODO(), nodeID, s.lastSnapshot[cluster])
		if err != nil {
			return err
		}
	}

	if req.Node != nil {
		if bv := req.Node.GetUserAgentBuildVersion(); bv != nil && bv.Version != nil {
			nodeVersion = fmt.Sprintf("v%d.%d.%d", bv.Version.MajorNumber, bv.Version.MinorNumber, bv.Version.Patch)
		}
	}

	loggers.LoggerAPKOperator.Debugf("Got a new request, response_nonce %s, nodeID %s, node_version %s",
		req.ResponseNonce, nodeID, nodeVersion)
	if status := req.ErrorDetail; status != nil {
		// if Envoy rejected the last update log the details here.
		// TODO(youngnick): Handle NACK properly
		errorCode = status.Code
		errorMessage = status.Message
	}
	loggers.LoggerAPKOperator.Debugf("handling v3 xDS resource request, response_nonce %s, nodeID %s, node_version %s, resource_names_subscribe %v, resource_names_unsubscribe %v, type_url %s, errorCode %d, errorMessage %s",
		req.ResponseNonce,
		nodeID, nodeVersion,
		req.ResourceNamesSubscribe, req.ResourceNamesUnsubscribe,
		req.GetTypeUrl(),
		errorCode, errorMessage)

	return nil
}

func (s *snapshotCache) OnStreamDeltaResponse(streamID int64, _ *discoveryv3.DeltaDiscoveryRequest, _ *discoveryv3.DeltaDiscoveryResponse) {
	// No mutex lock required here because no writing to the cache.
	node := s.streamIDNodeInfo[streamID]
	if node == nil {
		loggers.LoggerAPKOperator.Errorf("Tried to send a response to a node we haven't seen yet on stream %d", streamID)
	} else {
		loggers.LoggerAPKOperator.Debugf("Sending Incremental Response on stream %d to node %s", streamID, node.Id)
	}
}

func (s *snapshotCache) OnFetchRequest(_ context.Context, _ *discoveryv3.DiscoveryRequest) error {
	return nil
}

func (s *snapshotCache) OnFetchResponse(_ *discoveryv3.DiscoveryRequest, _ *discoveryv3.DiscoveryResponse) {
}
