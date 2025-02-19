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

package grpc

import (
	"fmt"
	"io"
	"net"
	"time"

	v3 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/analytics"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// AccessLogServiceServer is the gRPC server for the Access Log Service.
type AccessLogServiceServer struct {
	cfg         *config.Server
	analytics   *analytics.Analytics
	configStore *datastore.ConfigStore
}

// newAccessLogServiceServer creates a new instance of the Access Log Service Server.
func newAccessLogServiceServer(cfg *config.Server, configStore *datastore.ConfigStore) *AccessLogServiceServer {
	analytics := analytics.NewAnalytics(cfg, configStore)
	return &AccessLogServiceServer{
		cfg:         cfg,
		analytics:   analytics,
		configStore: configStore,
	}
}

// StreamAccessLogs streams access logs to the server.
func (s *AccessLogServiceServer) StreamAccessLogs(stream v3.AccessLogService_StreamAccessLogsServer) error {
	for {
		//s.cfg.Logger.Sugar().Debug("Received a stream of access logs")
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		for _, logEntry := range in.GetHttpLogs().LogEntry {
			s.analytics.Process(logEntry)
		}
	}
}

// StartAccessLogServiceServer starts the Access Log Service Server.
func StartAccessLogServiceServer(cfg *config.Server, configStore *datastore.ConfigStore) {
	// Create a new instance of the Access Log Service Server
	accessLogServiceServer := newAccessLogServiceServer(cfg, configStore)

	kaParams := keepalive.ServerParameters{
		Time:    time.Duration(cfg.ExternalProcessingKeepAliveTime) * time.Hour, // Ping the client if it is idle for 2 hours
		Timeout: 20 * time.Second,
	}
	server, err := util.CreateGRPCServer(cfg.EnforcerPublicKeyPath,
		cfg.EnforcerPrivateKeyPath,
		grpc.MaxRecvMsgSize(cfg.ExternalProcessingMaxMessageSize),
		grpc.MaxHeaderListSize(uint32(cfg.ExternalProcessingMaxHeaderLimit)),
		grpc.KeepaliveParams(kaParams))
	if err != nil {
		panic(err)
	}

	v3.RegisterAccessLogServiceServer(server, accessLogServiceServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.AccessLogServiceServerPort))
	if err != nil {
		cfg.Logger.Error(err, fmt.Sprintf("Failed to listen on port: %s", cfg.AccessLogServiceServerPort))
	}
	cfg.Logger.Sugar().Debug("Starting to serve access log service server")
	if err := server.Serve(listener); err != nil {
		cfg.Logger.Error(err, "Failed to serve access log service server")
	}
}
