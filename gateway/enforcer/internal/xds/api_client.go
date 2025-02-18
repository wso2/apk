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

package xds

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	v3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	api_ads "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	apiTypedURL = "type.googleapis.com/wso2.discovery.api.Api"
)

// APIXDSClient manages the connection to the API Discovery Service via gRPC.
// It supports connection retries, TLS configuration, and handling of API stream data.
type APIXDSClient struct {
	Host           string
	Port           string
	maxRetries     int
	retryInterval  time.Duration
	tlsConfig      *tls.Config
	grpcConn       *grpc.ClientConn
	ctx            context.Context
	cancel         context.CancelFunc
	client         api_ads.ApiDiscoveryServiceClient
	log            logging.Logger
	cfg            *config.Server
	latestReceived *v3.DiscoveryResponse
	latestACKed    *v3.DiscoveryResponse
	stream         api_ads.ApiDiscoveryService_StreamApisClient
	apiDatastore   *datastore.APIStore
}

// NewAPIXDSClient initializes a new instance of APIXDSClient with the given parameters.
// It sets up the host, port, retry logic, TLS configuration, and logger.
func NewAPIXDSClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server, apiDatastore *datastore.APIStore) *APIXDSClient {
	// Create a new APIClient object
	return &APIXDSClient{
		Host:          host,
		Port:          port,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		tlsConfig:     tlsConfig,
		grpcConn:      nil,
		log:           cfg.Logger,
		cfg:           cfg,
		apiDatastore:  apiDatastore,
	}
}

// InitiateAPIXDSConnection establishes a gRPC connection to the API Discovery Service
// and initiates a streaming API configuration. If the connection fails, it will retry
// based on the configured retry policy. Received configuration updates are logged.
func (c *APIXDSClient) InitiateAPIXDSConnection() {
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(nil, c.Host, c.Port, c.tlsConfig, c.maxRetries, c.retryInterval)
	c.grpcConn = grpcConn
	client := api_ads.NewApiDiscoveryServiceClient(grpcConn)
	c.client = client

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	stream, err := client.StreamApis(ctx)
	if err != nil {
		cancel()
		c.grpcConn.Close()
		c.log.Error(err, "failed to initiate XDS connection with API Discovery Service. Retrying the connection.")
		c.waitAndRetry()
		return
	}
	c.stream = stream
	// Send initial request
	dreq := DiscoveryRequestForNode(CreateNode(c.cfg.EnforcerLabel, c.cfg.InstanceIdentifier), "", "", nil, apiTypedURL)
	if stream == nil {
		c.log.Error(fmt.Errorf("failed to initiate XDS connection with API Discovery Service"), "Retrying the connection")
		c.grpcConn.Close()
		c.waitAndRetry()
		return
	}
	if err := stream.Send(dreq); err != nil {
		cancel()
		c.grpcConn.Close()
		panic(fmt.Errorf("failed to send initial discovery request: %v", err))
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.log.Error(err, "Failed to receive API stream data")
				c.nack(err)
				cancel()
				c.grpcConn.Close()
				c.waitAndRetry()
				return
			}
			c.latestReceived = resp
			handleResponseErr := c.handleResponse(resp)
			if handleResponseErr != nil {
				c.nack(handleResponseErr)
				continue
			}
			c.ack()
		}
	}()
}

func (c *APIXDSClient) waitAndRetry() {
	c.log.Sugar().Debug(fmt.Sprintf("Waiting for %d ms before retrying the connection", c.retryInterval.Milliseconds()))
	// Wait for a while before retrying the connection
	time.Sleep(c.retryInterval)
	go c.InitiateAPIXDSConnection()
}

func (c *APIXDSClient) ack() {
	dreq := DiscoveryRequestForNode(CreateNode(c.cfg.EnforcerLabel, c.cfg.InstanceIdentifier), c.latestReceived.GetVersionInfo(), c.latestReceived.GetNonce(), nil, apiTypedURL)
	c.stream.Send(dreq)
	c.latestACKed = c.latestReceived
}

func (c *APIXDSClient) nack(e error) {
	errDetail := &status.Status{
		Message: e.Error(),
	}
	dreq := DiscoveryRequestForNode(CreateNode(c.cfg.EnforcerLabel, c.cfg.InstanceIdentifier), c.latestACKed.GetVersionInfo(), c.latestReceived.GetNonce(), errDetail, apiTypedURL)
	c.stream.Send(dreq)
	c.latestACKed = c.latestReceived
}

func (c *APIXDSClient) handleResponse(response *v3.DiscoveryResponse) error {

	var apis []*api.Api
	for _, res := range response.GetResources() {
		var apiResource api.Api
		if err := proto.Unmarshal(res.GetValue(), &apiResource); err != nil {
			c.log.Sugar().Debug(fmt.Sprintf("Failed to unmarshal API resource: %v", err))
			return err
		}
		apis = append(apis, &apiResource)
	}
	c.apiDatastore.AddAPIs(apis)
	c.log.Sugar().Debug(fmt.Sprintf("Number of APIs received: %d", len(apis)))
	return nil
}
