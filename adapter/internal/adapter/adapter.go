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
 *
 */

// Package adapter contains the implementation to start the adapter
package adapter

import (
	"crypto/tls"
	"time"

	enforcerCallbacks "github.com/wso2/apk/adapter/internal/discovery/xds/enforcercallbacks"
	apiservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/api"
	configservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/config"
	subscriptionservice "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/subscription"
	wso2_server "github.com/wso2/apk/adapter/pkg/discovery/protocol/server/v3"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/tlsutils"

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
	flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "ads", ads, "Management server type (ads, xds, rest)")
}

const grpcMaxConcurrentStreams = 1000000

func runManagementServer(conf *config.Config, enforcerServer wso2_server.Server, enforcerSdsServer wso2_server.Server,
	enforcerAppDsSrv wso2_server.Server, enforcerAppKeyMappingDsSrv wso2_server.Server, port uint) {
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
		logger.LoggerMgw.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Failed to listen on port: %v, error: %v", port, err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1100,
		})
	}

	// register services
	configservice.RegisterConfigDiscoveryServiceServer(grpcServer, enforcerServer)
	apiservice.RegisterApiDiscoveryServiceServer(grpcServer, enforcerServer)
	subscriptionservice.RegisterSubscriptionDiscoveryServiceServer(grpcServer, enforcerSdsServer)
	subscriptionservice.RegisterApplicationDiscoveryServiceServer(grpcServer, enforcerAppDsSrv)
	subscriptionservice.RegisterApplicationKeyMappingDiscoveryServiceServer(grpcServer, enforcerAppKeyMappingDsSrv)

	logger.LoggerMgw.Info("port: ", port, " management server listening")
	go func() {
		logger.LoggerMgw.Info("Starting XDS GRPC server.")
		if err = grpcServer.Serve(lis); err != nil {
			logger.LoggerMgw.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Failed to start XDS GRPC server : %v", err.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1101,
			})
		}
	}()
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
		logger.LoggerMgw.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error reading the log configs. %v", errC.Error()),
			Severity:  logging.CRITICAL,
			ErrorCode: 1102,
		})
	}

	logger.LoggerMgw.Info("Starting adapter ....")
	enforcerCache := xds.GetEnforcerCache()
	enforcerSubscriptionCache := xds.GetEnforcerSubscriptionCache()
	enforcerApplicationCache := xds.GetEnforcerApplicationCache()
	enforcerApplicationKeyMappingCache := xds.GetEnforcerApplicationKeyMappingCache()

	enforcerXdsSrv := wso2_server.NewServer(ctx, enforcerCache, &enforcerCallbacks.Callbacks{})
	enforcerSdsSrv := wso2_server.NewServer(ctx, enforcerSubscriptionCache, &enforcerCallbacks.Callbacks{})
	enforcerAppDsSrv := wso2_server.NewServer(ctx, enforcerApplicationCache, &enforcerCallbacks.Callbacks{})
	enforcerAppKeyMappingDsSrv := wso2_server.NewServer(ctx, enforcerApplicationKeyMappingCache, &enforcerCallbacks.Callbacks{})

	runManagementServer(conf, enforcerXdsSrv, enforcerSdsSrv, enforcerAppDsSrv, enforcerAppKeyMappingDsSrv, port)

	// Set enforcer startup configs
	xds.UpdateEnforcerConfig(conf)

	// Set enforcer startup configs
	envs := conf.ControlPlane.EnvironmentLabels

	// If no environments are configured, default gateway label value is assigned.
	if len(envs) == 0 {
		envs = append(envs, config.DefaultGatewayName)
	}

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
