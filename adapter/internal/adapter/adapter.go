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
	routercb "github.com/wso2/apk/adapter/internal/discovery/xds/routercallbacks"
	"github.com/wso2/apk/adapter/internal/operator"
	xdstranslatorrunner "github.com/wso2/apk/adapter/internal/operator/gateway-api/translator/runner"
	xdsserverrunner "github.com/wso2/apk/adapter/internal/operator/gateway-api/xds/runner"
	infrarunner "github.com/wso2/apk/adapter/internal/operator/infrastructure/runner"
	"github.com/wso2/apk/adapter/internal/operator/message"
	"github.com/wso2/apk/adapter/internal/operator/provider-resources/runner"
	providerrunner "github.com/wso2/apk/adapter/internal/operator/provider/runner"
	apiservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	configservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/config"
	subscriptionservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/subscription"
	wso2_server "github.com/wso2/apk/adapter/pkg/discovery/protocol/server/v3"
	"github.com/wso2/apk/adapter/pkg/health"
	healthservice "github.com/wso2/apk/adapter/pkg/health/api/wso2/health/service"
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
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	debug       bool
	onlyLogging bool

	port    uint
	alsPort uint

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
	flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "ads", ads, "Management server type (ads, xds, rest)")
}

const grpcMaxConcurrentStreams = 1000000

func runManagementServer(conf *config.Config, server xdsv3.Server, enforcerServer wso2_server.Server,
	enforcerAPIDsSrv wso2_server.Server, enforcerAppPolicyDsSrv wso2_server.Server, enforcerSubPolicyDsSrv wso2_server.Server,
	enforcerKeyManagerDsSrv wso2_server.Server, enforcerRevokedTokenDsSrv wso2_server.Server,
	enforcerThrottleDataDsSrv wso2_server.Server, enforcerJwtIssuerDsSrv wso2_server.Server, port uint) {
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
		logger.LoggerAPK.Warn("failed to initiate the ssl context: ", err)
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
		logger.LoggerAPK.ErrorC(logging.PrintError(logging.Error1100, logging.BLOCKER, "Failed to listen on port: %v, error: %v", port, err.Error()))
	}

	// register services
	discoveryv3.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	configservice.RegisterConfigDiscoveryServiceServer(grpcServer, enforcerServer)
	apiservice.RegisterApiDiscoveryServiceServer(grpcServer, enforcerServer)
	subscriptionservice.RegisterJWTIssuerDiscoveryServiceServer(grpcServer, enforcerJwtIssuerDsSrv)
	// register health service
	healthservice.RegisterHealthServer(grpcServer, &health.Server{})

	logger.LoggerAPK.Info("port: ", port, " management server listening")
	go func() {
		logger.LoggerAPK.Info("Starting XDS GRPC server.")
		if err = grpcServer.Serve(lis); err != nil {
			logger.LoggerAPK.ErrorC(logging.PrintError(logging.Error1101, logging.BLOCKER, "Failed to start XDS GRPS server, error: %v", err.Error()))
		}
	}()

}

func SetupRunners(conf *config.Config) {
	ctx := ctrl.SetupSignalHandler()

	// Step 1: Start the Kubernetes Provider Service
	// It fetches the resources from the kubernetes
	// and publishes it
	// It also subscribes to status resources and once it receives
	// a status resource back, it writes it out.
	// Final processed crs will be stored in following pResources.
	pResources := new(message.ProviderResources)
	providerRunner := providerrunner.New(&providerrunner.Config{
		ProviderResources: pResources,
	})
	if err := providerRunner.Start(ctx); err != nil {
		logger.LoggerAPKOperator.Error("Error while starting provider service ", err)
	}

	// Step 2: Start the GatewayAPI Translator Runner
	// It subscribes to the provider resources, translates it to xDS IR
	// and infra IR resources and publishes them.
	// Final processed structs will be in pResources, xdsIR, and infraIR
	xdsIR := new(message.XdsIR)
	infraIR := new(message.InfraIR)
	gwRunner := runner.New(&runner.Config{
		ProviderResources: pResources,
		XdsIR:             xdsIR,
		InfraIR:           infraIR,
	})
	if err := gwRunner.Start(ctx); err != nil {
		logger.LoggerAPKOperator.Error("Error while starting translation service ", err)
	}

	// Step 3: Start the Xds Translator Service
	// It subscribes to the xdsIR, translates it into xds Resources and publishes it.
	// Final xds configs are in xds.
	xds := new(message.Xds)
	xdsTranslatorRunner := xdstranslatorrunner.New(&xdstranslatorrunner.Config{
		XdsIR:             xdsIR,
		Xds:               xds,
		ProviderResources: pResources,
	})
	if err := xdsTranslatorRunner.Start(ctx); err != nil {
		logger.LoggerAPKOperator.Error("Error while starting xds translator service ", err)
	}

	// Step 4: Start the Infra Manager Runner
	// It subscribes to the infraIR, translates it into Envoy Proxy infrastructure
	// resources such as K8s deployment and services.
	infraRunner := infrarunner.New(&infrarunner.Config{
		InfraIR: infraIR,
	})
	if err := infraRunner.Start(ctx); err != nil {
		logger.LoggerAPKOperator.Error("Error while starting infrastructure service ", err)
	}

	// Step 5: Start the xDS Server
	// It subscribes to the xds Resources and configures the remote Envoy Proxy
	// via the xDS Protocol.
	xdsServerRunner := xdsserverrunner.New(&xdsserverrunner.Config{
		Xds: xds,
	})
	if err := xdsServerRunner.Start(ctx); err != nil {
		logger.LoggerAPKOperator.Error("Error while starting xds service ", err)
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
		logger.LoggerAPK.ErrorC(logging.PrintError(logging.Error1102, logging.CRITICAL, "Error reading the log configs, error: %v", errC.Error()))
	}

	logger.LoggerAPK.Info("Starting adapter ....")

	cache := xds.GetXdsCache()
	enforcerCache := xds.GetEnforcerCache()
	enforcerAPICache := xds.GetEnforcerAPICache()
	enforcerApplicationPolicyCache := xds.GetEnforcerApplicationPolicyCache()
	enforcerSubscriptionPolicyCache := xds.GetEnforcerSubscriptionPolicyCache()
	enforcerKeyManagerCache := xds.GetEnforcerKeyManagerCache()
	enforcerRevokedTokenCache := xds.GetEnforcerRevokedTokenCache()
	enforcerThrottleDataCache := xds.GetEnforcerThrottleDataCache()
	enforcerJWtIssuerCache := xds.GetEnforcerJWTIssuerCache()
	srv := xdsv3.NewServer(ctx, cache, &routercb.Callbacks{})
	enforcerXdsSrv := wso2_server.NewServer(ctx, enforcerCache, &enforcerCallbacks.Callbacks{})
	enforcerJwtIssuerDsSrv := wso2_server.NewServer(ctx, enforcerJWtIssuerCache, &enforcerCallbacks.Callbacks{})
	enforcerAPIDsSrv := wso2_server.NewServer(ctx, enforcerAPICache, &enforcerCallbacks.Callbacks{})
	enforcerAppPolicyDsSrv := wso2_server.NewServer(ctx, enforcerApplicationPolicyCache, &enforcerCallbacks.Callbacks{})
	enforcerSubPolicyDsSrv := wso2_server.NewServer(ctx, enforcerSubscriptionPolicyCache, &enforcerCallbacks.Callbacks{})
	enforcerKeyManagerDsSrv := wso2_server.NewServer(ctx, enforcerKeyManagerCache, &enforcerCallbacks.Callbacks{})
	enforcerRevokedTokenDsSrv := wso2_server.NewServer(ctx, enforcerRevokedTokenCache, &enforcerCallbacks.Callbacks{})
	enforcerThrottleDataDsSrv := wso2_server.NewServer(ctx, enforcerThrottleDataCache, &enforcerCallbacks.Callbacks{})

	runManagementServer(conf, srv, enforcerXdsSrv, enforcerAPIDsSrv, enforcerAppPolicyDsSrv, enforcerSubPolicyDsSrv,
		enforcerKeyManagerDsSrv, enforcerRevokedTokenDsSrv, enforcerThrottleDataDsSrv, enforcerJwtIssuerDsSrv, port)

	// Set enforcer startup configs
	xds.UpdateEnforcerConfig(conf)
	if !conf.Adapter.EnableGatewayClassController {
		go operator.InitOperator(conf.Adapter.Metrics)
	} else {
		go SetupRunners(conf)
	}

OUTER:
	for {
		select {
		case l := <-watcherLogConf.Events:
			switch l.Op.String() {
			case "WRITE":
				logger.LoggerAPK.Info("Loading updated log config file...")
				config.ClearLogConfigInstance()
				logger.UpdateLoggers()
			}
		case s := <-sig:
			switch s {
			case os.Interrupt:
				logger.LoggerAPK.Info("Shutting down...")
				break OUTER
			}
		}
	}
	logger.LoggerAPK.Info("Bye!")
}
