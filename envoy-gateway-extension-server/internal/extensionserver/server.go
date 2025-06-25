// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package extensionserver

import (
	"context"
	"github.com/wso2/apk/envoy-gateway-extension-server/internal/config"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	structpb "google.golang.org/protobuf/types/known/structpb"
	pb "github.com/envoyproxy/gateway/proto/extension"
)

type Server struct {
	pb.UnimplementedEnvoyGatewayExtensionServer

	cfg *config.Server
}

func New(cfg *config.Server) *Server {
	return &Server{
		cfg: cfg,
	}
}

// PostHTTPListenerModify is called after Envoy Gateway is done generating a
// Listener xDS configuration and before that configuration is passed on to
// Envoy Proxy.
// This example adds Basic Authentication on the Listener level as an example.
// Note: This implementation is not secure, and should not be used to protect
// anything important.
func (s *Server) PostHTTPListenerModify(ctx context.Context, req *pb.PostHTTPListenerModifyRequest) (*pb.PostHTTPListenerModifyResponse, error) {
	s.cfg.Logger.Info("postHTTPListenerModify callback was invoked")
	s.cfg.Logger.Sugar().Infof("Received listener: %+v", req.Listener.Name)
	return &pb.PostHTTPListenerModifyResponse{
		Listener: req.Listener,
	}, nil
}

// PostRouteModify is called after Envoy Gateway is done generating a
// Route xDS configuration and before that configuration is passed on to
// Envoy Proxy.
// This example adds a custom header to the Route as an example.
// Note: This implementation is not secure, and should not be used to protect
// anything important.
func (s *Server) PostRouteModify(ctx context.Context, req *pb.PostRouteModifyRequest) (*pb.PostRouteModifyResponse, error) {
	s.cfg.Logger.Info("postRouteModify callback was invoked")
	s.cfg.Logger.Sugar().Infof("Received route: %+v", req.Route.Match)
	s.cfg.Logger.Sugar().Infof("Received Policies: %+v", req.PostRouteContext.ExtensionResources)
	if req.Route.Metadata == nil {
		req.Route.Metadata = &v31.Metadata{}
	}
	if req.Route.Metadata.FilterMetadata == nil {
		req.Route.Metadata.FilterMetadata = make(map[string]*structpb.Struct)
	}

	req.Route.Metadata.FilterMetadata["test-namespace"] =  &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"hello": &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: "World",
				},
			},
		},
	}
	return &pb.PostRouteModifyResponse{
		Route: req.Route,
	}, nil
}