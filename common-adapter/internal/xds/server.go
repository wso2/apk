/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org).
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
	"flag"
	"fmt"
	"net"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	envoy_cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xdsv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"
	"github.com/wso2/apk/common-adapter/internal/loggers"
	ratelimiterCallbacks "github.com/wso2/apk/common-adapter/internal/xds/callbacks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	rlsPort uint
	cache   envoy_cachev3.SnapshotCache
)

const (
	maxRandomInt             int = 999999999
	grpcMaxConcurrentStreams     = 1000000
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
	flag.UintVar(&rlsPort, "rls-port", 18001, "Rate Limiter management server port")
}

// InitCommonAdapterServer initializes the gRPC server for the common adapter.
func InitCommonAdapterServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	publicKeyLocation, privateKeyLocation, truststoreLocation := tlsutils.GetKeyLocations()
	cert, err := tlsutils.GetServerCertificate(publicKeyLocation, privateKeyLocation)

	caCertPool := tlsutils.GetTrustedCertPool(truststoreLocation)

	if err == nil {
		grpcOptions = append(grpcOptions, grpc.Creds(
			credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAndVerifyClientCert,
				ClientCAs:    caCertPool,
			}),
		))
	} else {
		loggers.LoggerAPK.Warn("failed to initiate the ssl context: ", err)
		panic(err)
	}

	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:    time.Duration(5 * time.Minute),
			Timeout: time.Duration(20 * time.Second),
		}),
	)
	rlsGrpcServer := grpc.NewServer(grpcOptions...)

	// It is required a separate gRPC server for the rate limit xDS, since it is the same RPC method
	// ADS used in both envoy xDS and rate limiter xDS.
	// According to https://github.com/envoyproxy/ratelimit/pull/368#discussion_r995831078 a separate RPC service is not
	// defined specifically to the rate limit xDS, instead using the ADS.
	loggers.LoggerXds.Info("port: ", rlsPort, " ratelimiter management server listening")
	rlsLis, err := net.Listen("tcp", fmt.Sprintf(":%d", rlsPort))
	if err != nil {
		loggers.LoggerXds.ErrorC(logging.GetErrorByCode(1106, rlsPort, err.Error()))
	}
	rateLimiterCache := GetRateLimiterCache()
	rlsSrv := xdsv3.NewServer(ctx, rateLimiterCache, &ratelimiterCallbacks.Callbacks{})

	discoveryv3.RegisterAggregatedDiscoveryServiceServer(rlsGrpcServer, rlsSrv)
	go func() {
		loggers.LoggerXds.Info("Starting Rate Limiter xDS gRPC server.")
		if err = rlsGrpcServer.Serve(rlsLis); err != nil {
			loggers.LoggerXds.ErrorC(logging.GetErrorByCode(1105, rlsPort, err.Error()))
		}
	}()

}

// GetRateLimiterCache returns xds server cache for rate limiter service.
func GetRateLimiterCache() envoy_cachev3.SnapshotCache {
	return rlsPolicyCache.xdsCache
}
