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

// Package adapter contains the implementation to start the adapter
package adapter

import (
	"crypto/tls"
	"time"

	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	xdsv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	enforcerCallbacks "github.com/wso2/apk/adapter/internal/discovery/xds/enforcercallbacks"
	ratelimiterCallbacks "github.com/wso2/apk/adapter/internal/discovery/xds/ratelimitercallbacks"
	routercb "github.com/wso2/apk/adapter/internal/discovery/xds/routercallbacks"
	apiservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	configservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/config"
	keymanagerservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/keymgt"
	subscriptionservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/subscription"
	throttleservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/throttle"
	wso2_server "github.com/wso2/apk/adapter/pkg/discovery/protocol/server/v3"
	"github.com/wso2/apk/adapter/pkg/health"
	healthservice "github.com/wso2/apk/adapter/pkg/health/api/wso2/health/service"
	"github.com/wso2/apk/adapter/pkg/operator"
	"github.com/wso2/apk/adapter/pkg/utils/tlsutils"

	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/discovery/xds"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	debug       bool
	onlyLogging bool

	localhost = "0.0.0.0"

	port        uint
	gatewayPort uint
	rlsPort     uint
	alsPort     uint

	mode string
)

const (
	ads          = "ads"
	amqpProtocol = "amqp"
)

func init() {
	flag.BoolVar(&debug, "debug", true, "Use debug logging")
	flag.BoolVar(&onlyLogging, "onlyLogging", false, "Only demo AccessLogging Service")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&gatewayPort, "gateway", 18001, "Management server port for HTTP gateway")
	flag.UintVar(&rlsPort, "rls-port", 18001, "Rate Limiter management server port")
	flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "ads", ads, "Management server type (ads, xds, rest)")
}

const grpcMaxConcurrentStreams = 1000000

func runManagementServer(conf *config.Config, server xdsv3.Server, rlsServer xdsv3.Server, enforcerServer wso2_server.Server, enforcerSdsServer wso2_server.Server,
	enforcerAppDsSrv wso2_server.Server, enforcerAPIDsSrv wso2_server.Server, enforcerAppPolicyDsSrv wso2_server.Server,
	enforcerSubPolicyDsSrv wso2_server.Server, enforcerAppKeyMappingDsSrv wso2_server.Server,
	enforcerKeyManagerDsSrv wso2_server.Server, enforcerRevokedTokenDsSrv wso2_server.Server,
	enforcerThrottleDataDsSrv wso2_server.Server, port, rlsPort uint) {
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
		logger.LoggerMgw.Warn("failed to initiate the ssl context: ", err)
		panic(err)
	}

	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:    time.Duration(5 * time.Minute),
			Timeout: time.Duration(20 * time.Second),
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.LoggerMgw.ErrorC(logging.GetErrorByCode(1100, port, err.Error()))
	}

	// register services
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	configservice.RegisterConfigDiscoveryServiceServer(grpcServer, enforcerServer)
	apiservice.RegisterApiDiscoveryServiceServer(grpcServer, enforcerServer)
	subscriptionservice.RegisterSubscriptionDiscoveryServiceServer(grpcServer, enforcerSdsServer)
	subscriptionservice.RegisterApplicationDiscoveryServiceServer(grpcServer, enforcerAppDsSrv)
	subscriptionservice.RegisterApiListDiscoveryServiceServer(grpcServer, enforcerAPIDsSrv)
	subscriptionservice.RegisterApplicationPolicyDiscoveryServiceServer(grpcServer, enforcerAppPolicyDsSrv)
	subscriptionservice.RegisterSubscriptionPolicyDiscoveryServiceServer(grpcServer, enforcerSubPolicyDsSrv)
	subscriptionservice.RegisterApplicationKeyMappingDiscoveryServiceServer(grpcServer, enforcerAppKeyMappingDsSrv)
	keymanagerservice.RegisterKMDiscoveryServiceServer(grpcServer, enforcerKeyManagerDsSrv)
	keymanagerservice.RegisterRevokedTokenDiscoveryServiceServer(grpcServer, enforcerRevokedTokenDsSrv)
	throttleservice.RegisterThrottleDataDiscoveryServiceServer(grpcServer, enforcerThrottleDataDsSrv)

	// register health service
	healthservice.RegisterHealthServer(grpcServer, &health.Server{})

	logger.LoggerMgw.Info("port: ", port, " management server listening")
	go func() {
		// if control plane enabled wait until it starts
		if conf.ControlPlane.Enabled {
			// wait current goroutine forever for until control plane starts
			health.WaitForControlPlane()
		}
		logger.LoggerMgw.Info("Starting XDS GRPC server.")
		if err = grpcServer.Serve(lis); err != nil {
			logger.LoggerMgw.ErrorC(logging.GetErrorByCode(1101, err.Error()))
		}
	}()

	if conf.Envoy.RateLimit.Enabled {
		// It is required a separate gRPC server for the rate limit xDS, since it is the same RPC method
		// ADS used in both envoy xDS and rate limiter xDS.
		// According to https://github.com/envoyproxy/ratelimit/pull/368#discussion_r995831078 a separate RPC service is not
		// defined specifically to the rate limit xDS, instead using the ADS.
		logger.LoggerMgw.Info("port: ", rlsPort, " ratelimiter management server listening")
		rlsGrpcServer := grpc.NewServer(grpcOptions...)
		rlsLis, err := net.Listen("tcp", fmt.Sprintf(":%d", rlsPort))
		if err != nil {
			logger.LoggerMgw.ErrorC(logging.GetErrorByCode(1106, port, err.Error()))
		}

		discoveryv3.RegisterAggregatedDiscoveryServiceServer(rlsGrpcServer, rlsServer)
		go func() {
			logger.LoggerMgw.Info("Starting Rate Limiter xDS gRPC server.")
			if err = rlsGrpcServer.Serve(rlsLis); err != nil {
				logger.LoggerMgw.ErrorC(logging.GetErrorByCode(1105, port, err.Error()))
			}
		}()
	}

}

// Run starts the XDS server and Rest API server.
func Run(conf *config.Config) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt)
	// TODO: (VirajSalaka) Support the REST API Configuration via flags only if it is a valid requirement
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// log config watcher
	watcherLogConf, _ := fsnotify.NewWatcher()
	logConfigPath, errC := config.GetLogConfigPath()
	if errC == nil {
		errC = watcherLogConf.Add(logConfigPath)
	}

	if errC != nil {
		logger.LoggerMgw.ErrorC(logging.GetErrorByCode(1102, errC.Error()))
	}

	logger.LoggerMgw.Info("Starting adapter ....")
	cache := xds.GetXdsCache()
	rateLimiterCache := xds.GetRateLimiterCache()
	enforcerCache := xds.GetEnforcerCache()
	enforcerSubscriptionCache := xds.GetEnforcerSubscriptionCache()
	enforcerApplicationCache := xds.GetEnforcerApplicationCache()
	enforcerAPICache := xds.GetEnforcerAPICache()
	enforcerApplicationPolicyCache := xds.GetEnforcerApplicationPolicyCache()
	enforcerSubscriptionPolicyCache := xds.GetEnforcerSubscriptionPolicyCache()
	enforcerApplicationKeyMappingCache := xds.GetEnforcerApplicationKeyMappingCache()
	enforcerKeyManagerCache := xds.GetEnforcerKeyManagerCache()
	enforcerRevokedTokenCache := xds.GetEnforcerRevokedTokenCache()
	enforcerThrottleDataCache := xds.GetEnforcerThrottleDataCache()

	srv := xdsv3.NewServer(ctx, cache, &routercb.Callbacks{})
	rlsSrv := xdsv3.NewServer(ctx, rateLimiterCache, &ratelimiterCallbacks.Callbacks{})
	enforcerXdsSrv := wso2_server.NewServer(ctx, enforcerCache, &enforcerCallbacks.Callbacks{})
	enforcerSdsSrv := wso2_server.NewServer(ctx, enforcerSubscriptionCache, &enforcerCallbacks.Callbacks{})
	enforcerAppDsSrv := wso2_server.NewServer(ctx, enforcerApplicationCache, &enforcerCallbacks.Callbacks{})
	enforcerAPIDsSrv := wso2_server.NewServer(ctx, enforcerAPICache, &enforcerCallbacks.Callbacks{})
	enforcerAppPolicyDsSrv := wso2_server.NewServer(ctx, enforcerApplicationPolicyCache, &enforcerCallbacks.Callbacks{})
	enforcerSubPolicyDsSrv := wso2_server.NewServer(ctx, enforcerSubscriptionPolicyCache, &enforcerCallbacks.Callbacks{})
	enforcerAppKeyMappingDsSrv := wso2_server.NewServer(ctx, enforcerApplicationKeyMappingCache, &enforcerCallbacks.Callbacks{})
	enforcerKeyManagerDsSrv := wso2_server.NewServer(ctx, enforcerKeyManagerCache, &enforcerCallbacks.Callbacks{})
	enforcerRevokedTokenDsSrv := wso2_server.NewServer(ctx, enforcerRevokedTokenCache, &enforcerCallbacks.Callbacks{})
	enforcerThrottleDataDsSrv := wso2_server.NewServer(ctx, enforcerThrottleDataCache, &enforcerCallbacks.Callbacks{})

	runManagementServer(conf, srv, rlsSrv, enforcerXdsSrv, enforcerSdsSrv, enforcerAppDsSrv, enforcerAPIDsSrv,
		enforcerAppPolicyDsSrv, enforcerSubPolicyDsSrv, enforcerAppKeyMappingDsSrv, enforcerKeyManagerDsSrv,
		enforcerRevokedTokenDsSrv, enforcerThrottleDataDsSrv, port, rlsPort)

	// Set enforcer startup configs
	xds.UpdateEnforcerConfig(conf)

	go operator.InitOperator()

OUTER:
	for {
		select {
		case l := <-watcherLogConf.Events:
			switch l.Op.String() {
			case "WRITE":
				logger.LoggerMgw.Info("Loading updated log config file...")
				config.ClearLogConfigInstance()
				logger.UpdateLoggers()
			}
		case s := <-sig:
			switch s {
			case os.Interrupt:
				logger.LoggerMgw.Info("Shutting down...")
				break OUTER
			}
		}
	}
	logger.LoggerMgw.Info("Bye!")
}
