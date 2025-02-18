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
	config_from_adapter "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/config/enforcer"
	config_ads "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/config"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	configTypedURL      = "type.googleapis.com/wso2.discovery.config.enforcer.Config"
	commonEnforcerLabel = "commonEnforcerLabel"
)

// ConfigXDSClient is a client for managing gRPC connections to the Config Discovery Service (XDS).
// It handles retry logic, TLS configuration, and logging for configuration data streams.
type ConfigXDSClient struct {
	Host            string
	Port            string
	maxRetries      int
	retryInterval   time.Duration
	tlsConfig       *tls.Config
	grpcConn        *grpc.ClientConn
	ctx             context.Context
	cancel          context.CancelFunc
	client          config_ads.ConfigDiscoveryServiceClient
	log             logging.Logger
	cfg             *config.Server
	latestReceived  *v3.DiscoveryResponse
	latestACKed     *v3.DiscoveryResponse
	stream          config_ads.ConfigDiscoveryService_StreamConfigsClient
	configDatastore *datastore.ConfigStore
}

// NewXDSConfigClient creates a new instance of ConfigXDSClient.
// It initializes the client with the given host, port, retry parameters, TLS configuration, and logger.
func NewXDSConfigClient(host string, port string, maxRetries int, retryInterval time.Duration, tlsConfig *tls.Config, cfg *config.Server, configDatastore *datastore.ConfigStore) *ConfigXDSClient {
	// Create a new APIClient object
	return &ConfigXDSClient{
		Host:            host,
		Port:            port,
		maxRetries:      maxRetries,
		retryInterval:   retryInterval,
		tlsConfig:       tlsConfig,
		grpcConn:        nil,
		log:             cfg.Logger,
		cfg:             cfg,
		configDatastore: configDatastore,
	}
}

// InitiateConfigXDSConnection establishes and maintains a gRPC connection to the Config Discovery Service.
// It also handles reconnection logic on errors and listens for incoming configuration streams.
func (c *ConfigXDSClient) InitiateConfigXDSConnection() {
	grpcConn := util.CreateGRPCConnectionWithRetryAndPanic(nil, c.Host, c.Port, c.tlsConfig, c.maxRetries, c.retryInterval)
	c.grpcConn = grpcConn
	client := config_ads.NewConfigDiscoveryServiceClient(grpcConn)

	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel

	stream, err := client.StreamConfigs(ctx)
	if err != nil {
		cancel()
		c.grpcConn.Close()
		c.log.Error(err, "Failed to initiate XDS connection with Config Discovery Service. Retrying the connection.")
		c.waitAndRetry()
		return
	}

	c.stream = stream
	// Send initial request
	if stream == nil {
		c.log.Error(fmt.Errorf("failed to initiate XDS connection with Config Discovery Service"), "Retrying the connection")
		c.grpcConn.Close()
		c.waitAndRetry()
		return
	}
	dreq := DiscoveryRequestForNode(CreateNode(commonEnforcerLabel, c.cfg.InstanceIdentifier), "", "", nil, configTypedURL)
	if err := stream.Send(dreq); err != nil {
		cancel()
		c.grpcConn.Close()
		panic(fmt.Errorf("failed to send initial discovery request: %v", err))
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				c.log.Error(err, "Failed to receive config data")
				c.nack(err)
				cancel()
				c.grpcConn.Close()
				c.waitAndRetry()
				return
			}
			// c.log.Info(fmt.Sprintf("Received config: %v", resp))
			c.latestReceived = resp
			handleRespErr := c.handleResponse(resp)
			if handleRespErr != nil {
				c.nack(handleRespErr)
				continue
			}
			c.ack()
		}
	}()
}

func (c *ConfigXDSClient) ack() {
	dreq := DiscoveryRequestForNode(CreateNode(commonEnforcerLabel, c.cfg.InstanceIdentifier), c.latestReceived.GetVersionInfo(), c.latestReceived.GetNonce(), nil, configTypedURL)
	c.stream.Send(dreq)
	c.latestACKed = c.latestReceived
}

func (c *ConfigXDSClient) nack(e error) {
	errDetail := &status.Status{
		Message: e.Error(),
	}
	dreq := DiscoveryRequestForNode(CreateNode(commonEnforcerLabel, c.cfg.InstanceIdentifier), c.latestACKed.GetVersionInfo(), c.latestReceived.GetNonce(), errDetail, configTypedURL)
	c.stream.Send(dreq)
	c.latestACKed = c.latestReceived
}

func (c *ConfigXDSClient) handleResponse(response *v3.DiscoveryResponse) error {

	var configs []*config_from_adapter.Config
	for _, res := range response.GetResources() {
		var configResource config_from_adapter.Config
		if err := proto.Unmarshal(res.GetValue(), &configResource); err != nil {
			c.log.Sugar().Debug(fmt.Sprintf("Failed to unmarshal Config resource: %v", err))
			return err
		}
		configs = append(configs, &configResource)
	}
	c.configDatastore.AddConfigs(configs)
	c.log.Sugar().Debug(fmt.Sprintf("Number of Configs received: %d", len(configs)))
	return nil
}

func (c *ConfigXDSClient) waitAndRetry() {
	c.log.Sugar().Debug(fmt.Sprintf("Waiting for %d ms before retrying the connection", c.retryInterval.Milliseconds()))
	// Wait for a while before retrying the connection
	time.Sleep(c.retryInterval)
	go c.InitiateConfigXDSConnection()
}
