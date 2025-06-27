// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package extensionserver

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/envoyproxy/gateway/proto/extension"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/wso2/apk/envoy-gateway-extension-server/internal/config"
	structpb "google.golang.org/protobuf/types/known/structpb"
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
	s.cfg.Logger.Sugar().Debugf("Received route: %+v", req.Route.Match)
	s.cfg.Logger.Sugar().Debugf("Received Policies: %+v", req.PostRouteContext.ExtensionResources)
	if req.Route.Metadata == nil {
		req.Route.Metadata = &v31.Metadata{}
	}
	if req.Route.Metadata.FilterMetadata == nil {
		req.Route.Metadata.FilterMetadata = make(map[string]*structpb.Struct)
	}

	// Traverse through all the extension resources and prepare a extension resource identifier list
	extenstionResourceIdentifiers := make([]string, 0, len(req.PostRouteContext.ExtensionResources))
	for _, extenstionResource := range req.PostRouteContext.ExtensionResources {
		// Convert UnstructuredBytes to JSON string
		var jsonObj map[string]interface{}
		if err := json.Unmarshal(extenstionResource.UnstructuredBytes, &jsonObj); err != nil {
			s.cfg.Logger.Sugar().Errorf("Failed to unmarshal extension resource: %v", err)
			continue
		}
		jsonStr, err := json.MarshalIndent(jsonObj, "", "  ")
		if err != nil {
			s.cfg.Logger.Sugar().Errorf("Failed to pretty print extension resource: %v", err)
		} else {
			s.cfg.Logger.Sugar().Infof("Extension resource (pretty):\n%s", string(jsonStr))
		}
		kindValue, ok := jsonObj["kind"]
		var kindStr string
		if ok {
			if str, ok := kindValue.(string); ok {
				kindStr = str
			}
		}
		nameValue, ok := jsonObj["metadata"]
		var nameStr string
		var namespace string
		if ok {
			if metadata, ok := nameValue.(map[string]interface{}); ok {
				nameValue, ok := metadata["name"]
				if ok {
					if str, ok := nameValue.(string); ok {
						nameStr = str
					}
				}
				namespaceValue, ok := metadata["namespace"]
				if ok {
					if str, ok := namespaceValue.(string); ok {
						namespace = str
					}
				}
			}
		}
		s.cfg.Logger.Sugar().Debugf("Extension resource processed > kind: %s, name: %s, namespace: %s", kindStr, nameStr, namespace)
		extenstionResourceIdentifiers = append(extenstionResourceIdentifiers, fmt.Sprintf("%s/%s/%s", kindStr, namespace, nameStr))
	}

	values := make([]*structpb.Value, 0, len(extenstionResourceIdentifiers))
	for _, id := range extenstionResourceIdentifiers {
		values = append(values, structpb.NewStringValue(id))
	}
	req.Route.Metadata.FilterMetadata["ext_proc"] =  &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"ExtensionRefs": structpb.NewListValue(&structpb.ListValue{
				Values: values,
			}),
		},
	}
	return &pb.PostRouteModifyResponse{
		Route: req.Route,
	}, nil
}