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

package runner

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/keepalive"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimev3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretv3 "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/operator/gateway-api/xds/cache"
	"github.com/wso2/apk/adapter/internal/operator/message"
	"github.com/wso2/apk/adapter/internal/types"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// XdsServerAddress is the listening address of the xds-server.
	XdsServerAddress = "0.0.0.0"
	// xdsTLSCertFilename is the fully qualified path of the file containing the
	// xDS server TLS certificate.
	// xdsTLSCertFilename = "/home/wso2/security/keystore/adapter.crt"
	// // xdsTLSKeyFilename is the fully qualified path of the file containing the
	// // xDS server TLS key.
	// xdsTLSKeyFilename = "/home/wso2/security/keystore/adapter.key"
	// xdsTLSCaFilename is the fully qualified path of the file containing the
	// xDS server trusted CA certificate.
	// xdsTLSCaFilename = "/home/wso2/security/truststore/adapter.crt"
)

type Config struct {
	Xds   *message.Xds
	grpc  *grpc.Server
	cache cache.SnapshotCacheWithCallbacks
}

type Runner struct {
	Config
}

func New(cfg *Config) *Runner {
	return &Runner{Config: *cfg}
}

func (r *Runner) Name() string {
	return string("xds-server")
}

// Start starts the xds-server runner
func (r *Runner) Start(ctx context.Context) (err error) {

	// Set up the gRPC server and register the xDS handler.
	// Create SnapshotCache before start subscribeAndTranslate,
	// prevent panics in case cache is nil.
	cfg := r.tlsConfig(tlsutils.GetKeyLocations())
	r.grpc = grpc.NewServer(grpc.Creds(credentials.NewTLS(cfg)), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             15 * time.Second,
		PermitWithoutStream: true,
	}))

	r.cache = cache.NewSnapshotCache(true)
	registerServer(serverv3.NewServer(ctx, r.cache, r.cache), r.grpc)

	// Start and listen xDS gRPC Server.
	go r.serveXdsServer(ctx)

	// Start message Subscription.
	go r.subscribeAndTranslate(ctx)
	loggers.LoggerAPKOperator.Info("started xds runner")
	return
}

func (r *Runner) serveXdsServer(ctx context.Context) {
	conf := config.ReadConfigs()
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Deployment.Gateway.AdapterXDSPort))
	if err != nil {
		loggers.LoggerAPKOperator.Error(err, "failed to listen on port", "address", conf.Deployment.Gateway.AdapterXDSPort)
		return
	}

	go func() {
		<-ctx.Done()
		loggers.LoggerAPKOperator.Info("grpc server shutting down")
		// We don't use GracefulStop here because envoy
		// has long-lived hanging xDS requests. There's no
		// mechanism to make those pending requests fail,
		// so we forcibly terminate the TCP sessions.
		r.grpc.Stop()
	}()

	if err = r.grpc.Serve(l); err != nil {
		loggers.LoggerAPKOperator.Error(err, "failed to start grpc based xds server")
	}
}

// registerServer registers the given xDS protocol Server with the gRPC
// runtime.
func registerServer(srv serverv3.Server, g *grpc.Server) {
	// register services
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(g, srv)
	secretv3.RegisterSecretDiscoveryServiceServer(g, srv)
	clusterv3.RegisterClusterDiscoveryServiceServer(g, srv)
	endpointv3.RegisterEndpointDiscoveryServiceServer(g, srv)
	listenerv3.RegisterListenerDiscoveryServiceServer(g, srv)
	routev3.RegisterRouteDiscoveryServiceServer(g, srv)
	runtimev3.RegisterRuntimeDiscoveryServiceServer(g, srv)
}

func (r *Runner) subscribeAndTranslate(ctx context.Context) {
	// Subscribe to resources
	message.HandleSubscription(message.Metadata{Runner: string("xds-server"), Message: "xds"}, r.Xds.Subscribe(ctx),
		func(update message.Update[string, *types.ResourceVersionTable], errChan chan error) {
			key := update.Key
			val := update.Value

			loggers.LoggerAPKOperator.Info("Received an update in xds server")
			var err error
			if update.Delete {
				err = r.cache.GenerateNewSnapshot(key, nil)
			} else if val != nil && val.XdsResources != nil {
				if r.cache == nil {
					loggers.LoggerAPKOperator.Error("Failed to init snapshot cache ", err)
					errChan <- err
				} else {
					// Update snapshot cache
					err = r.cache.GenerateNewSnapshot(key, val.XdsResources)
				}
			}
			if err != nil {
				loggers.LoggerAPKOperator.Error("Failed to generate a snapshot ", err)
				errChan <- err
			}
		},
	)

	loggers.LoggerAPKOperator.Info("subscriber shutting down")
}

func (r *Runner) tlsConfig(cert, key, ca string) *tls.Config {
	extCert, err := tlsutils.GetServerCertificate(cert, key)
	if err != nil {
		loggers.LoggerAPKOperator.Error("failed to parse CA certificate")
		return nil
	}
	caCertPool := tlsutils.GetTrustedCertPool(ca)
	return &tls.Config{
		Certificates: []tls.Certificate{extCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS13,
	}
}
