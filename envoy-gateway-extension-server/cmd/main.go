// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/envoyproxy/gateway/proto/extension"
	"github.com/wso2/apk/envoy-gateway-extension-server/internal/config"
	"github.com/wso2/apk/envoy-gateway-extension-server/internal/extensionserver"
)

func main() {
	startExtensionServer()
	// wait forever
	select {}
}

var grpcServer *grpc.Server
func startExtensionServer() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)
	go func() {
		for range c {
			if grpcServer != nil {
				grpcServer.Stop()
				os.Exit(0)
			}
		}
	}()
	cfg := config.GetConfig()
	address := net.JoinHostPort(cfg.ExtensionServerHost, cfg.ExtensionServerPort)
	cfg.Logger.Sugar().Infof("Starting the extension server", fmt.Sprintf("host, %s", address))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	grpcServer = grpc.NewServer(opts...)
	pb.RegisterEnvoyGatewayExtensionServer(grpcServer, extensionserver.New(cfg))
	return grpcServer.Serve(lis)
}
