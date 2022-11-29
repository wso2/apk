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

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	apiProtos "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/service/apkmgt"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/adapter/pkg/tlsutils"
	"github.com/wso2/apk/management-server/internal/config"
	"github.com/wso2/apk/management-server/internal/database"
	"github.com/wso2/apk/management-server/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type apiService struct {
	apiProtos.UnimplementedAPIServiceServer
}

func NewApiService() *apiService {
	return &apiService{}
}

func (s *apiService) CreateAPI(ctx context.Context, api *apiProtos.API) (*apiProtos.Response, error) {
	logger.LoggerMGTServer.Infof("Message received : %q", api)
	database.CreateAPI(api)
	return &apiProtos.Response{Result: true}, nil
}

func (s *apiService) UpdateAPI(ctx context.Context, api *apiProtos.API) (*apiProtos.Response, error) {
	logger.LoggerMGTServer.Infof("Message received : %q", api)
	database.UpdateAPI(api)
	return &apiProtos.Response{Result: true}, nil
}

func (s *apiService) DeleteAPI(ctx context.Context, api *apiProtos.API) (*apiProtos.Response, error) {
	logger.LoggerMGTServer.Infof("Message received : %q", api)
	database.DeleteAPI(api)
	return &apiProtos.Response{Result: true}, nil
}

func StartGRPCServer() {
	var grpcOptions []grpc.ServerOption
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
		logger.LoggerMGTServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Failed to initiate the ssl context, error: %v", err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1200,
		})
	}
	grpcOptions = append(grpcOptions, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			Time:    time.Duration(5 * time.Minute),
			Timeout: time.Duration(20 * time.Second),
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)
	conf := config.ReadConfigs()
	port := conf.ManagementServer.GRPCPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.LoggerMGTServer.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Failed to listen on port: %v, error: %v", port, err.Error()),
			Severity:  logging.BLOCKER,
			ErrorCode: 1201,
		})
	}
	// register services
	apiService := NewApiService()
	apiProtos.RegisterAPIServiceServer(grpcServer, apiService)
	logger.LoggerMGTServer.Infof("Management server is listening for GRPC connections on port: %v.", port)
	grpcServer.Serve(lis)
}
