/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	apiProtos "github.com/wso2/apk/management-server/internal/grpc/api"
	logger "github.com/wso2/apk/management-server/internal/loggers"
	"github.com/wso2/apk/adapter/pkg/logging"
	"github.com/wso2/apk/management-server/pkg/tlsutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"github.com/wso2/apk/management-server/config"
)

type aPIService struct {
	apiProtos.UnimplementedAPIServiceServer
}

func NewApiService() *aPIService {
	return &aPIService{}
}

func (s *aPIService) CreateAPI(ctx context.Context, api *apiProtos.API ) (*apiProtos.Response, error) {
	logger.LoggerMGTServer.Printf("Message received : %q", api);
	// TODO database calls to persist data
	return &apiProtos.Response{Result : true}, nil
}

func (s *aPIService) UpdateAPI(ctx context.Context, api *apiProtos.API ) (*apiProtos.Response, error) {
	logger.LoggerMGTServer.Printf("Message received : %q", api);
	// TODO database calls to persist data
	return &apiProtos.Response{Result : true}, nil
}

func (s *aPIService) DeleteAPI(ctx context.Context, api *apiProtos.API ) (*apiProtos.Response, error) {
	logger.LoggerMGTServer.Printf("Message received : %q", api);
	// TODO database calls to persist data
	return &apiProtos.Response{Result : true}, nil
}

func RunManagementServer() {
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
		logger.LoggerMGTServer.Warn("failed to initiate the ssl context: ", err)
		panic(err)
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
			ErrorCode: 1100,
		})
	}

	// register services
	apiService := NewApiService();
	apiProtos.RegisterAPIServiceServer(grpcServer, apiService)



	logger.LoggerMGTServer.Info("port: ", port, " management server listening")
	grpcServer.Serve(lis)
}
