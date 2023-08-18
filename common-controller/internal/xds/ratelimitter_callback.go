/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	"context"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/common-controller/internal/loggers"
	utils "github.com/wso2/apk/common-controller/internal/utils"
)

var nodeQueueInstance *utils.NodeQueue

func init() {
	nodeQueueInstance = utils.GenerateNodeQueue()
}

// Callbacks is used to debug the xds server related communication.
type Callbacks struct {
}

// Report logs the fetches and requests.
func (cb *Callbacks) Report() {}

// OnStreamOpen prints debug logs
func (cb *Callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	loggers.LoggerAPKOperator.Debugf("stream %d open for %s\n", id, typ)
	return nil
}

// OnStreamClosed prints debug logs
func (cb *Callbacks) OnStreamClosed(id int64, node *core.Node) {
	loggers.LoggerAPKOperator.Debugf("stream %d closed\n", id)
}

// OnStreamRequest prints debug logs
func (cb *Callbacks) OnStreamRequest(id int64, request *discovery.DiscoveryRequest) error {
	nodeIdentifier := utils.GetNodeIdentifier(request) // TODO: (renuka) set metadata instanceIdentifier from rate limiter (have to add in ADS Client impl)
	if nodeQueueInstance.IsNewNode(nodeIdentifier) {
		loggers.LoggerAPKOperator.Infof("stream request on stream id: %d, from node: %s, version: %s",
			id, nodeIdentifier, request.VersionInfo)
	}
	loggers.LoggerAPKOperator.Debugf("stream request on stream id: %d, from node: %s, version: %s, for type: %s",
		id, nodeIdentifier, request.VersionInfo, request.TypeUrl)
	if request.ErrorDetail != nil {
		loggers.LoggerAPKOperator.ErrorC(logging.PrintError(logging.Error2300, logging.MAJOR,
			"Stream request for type %s on stream id: %d, from node: %s, Error: %s", request.GetTypeUrl(),
			id, nodeIdentifier, request.ErrorDetail.Message))
	}
	return nil
}

// OnStreamResponse prints debug logs
func (cb *Callbacks) OnStreamResponse(context context.Context, id int64, request *discovery.DiscoveryRequest,
	response *discovery.DiscoveryResponse) {
	nodeIdentifier := utils.GetNodeIdentifier(request)
	loggers.LoggerAPKOperator.Debugf("stream response on stream id: %d, to node: %s, version: %s, for type: %v", id,
		nodeIdentifier, response.VersionInfo, response.TypeUrl)
}

// OnFetchRequest prints debug logs
func (cb *Callbacks) OnFetchRequest(_ context.Context, req *discovery.DiscoveryRequest) error {
	loggers.LoggerAPKOperator.Debugf("fetch request from node %s, version: %s, for type %s", utils.GetNodeIdentifier(req),
		req.VersionInfo, req.TypeUrl)
	return nil
}

// OnFetchResponse prints debug logs
func (cb *Callbacks) OnFetchResponse(req *discovery.DiscoveryRequest, res *discovery.DiscoveryResponse) {
	loggers.LoggerAPKOperator.Debugf("fetch response to node: %s, version: %s, for type %s", utils.GetNodeIdentifier(req),
		req.VersionInfo, res.TypeUrl)
}

// OnDeltaStreamOpen is unused.
func (cb *Callbacks) OnDeltaStreamOpen(_ context.Context, id int64, typ string) error {
	return nil
}

// OnDeltaStreamClosed is unused.
func (cb *Callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
}

// OnStreamDeltaResponse is unused.
func (cb *Callbacks) OnStreamDeltaResponse(id int64, req *discovery.DeltaDiscoveryRequest, res *discovery.DeltaDiscoveryResponse) {
}

// OnStreamDeltaRequest is unused.
func (cb *Callbacks) OnStreamDeltaRequest(id int64, req *discovery.DeltaDiscoveryRequest) error {
	return nil
}
