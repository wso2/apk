// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package extensionserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	// "time"

	pb "github.com/envoyproxy/gateway/proto/extension"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	jwt_authnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	constants "github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/envoy-gateway-extension-server/internal/config"
	"google.golang.org/protobuf/types/known/anypb"

	// durationpb "google.golang.org/protobuf/types/known/durationpb"
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
	filterChains := req.Listener.GetFilterChains()
	defaultFC := req.Listener.DefaultFilterChain
	if defaultFC != nil {
		filterChains = append(filterChains, defaultFC)
	}
	// Go over all of the chains, and add the basic authentication http filter
	for _, currChain := range filterChains {
		httpConManager, hcmIndex, err := findHCM(currChain)
		if err != nil {
			s.cfg.Logger.Sugar().Error("failed to find an HCM in the current chain", slog.Any("error", err))
			continue
		}
		// If a basic authentication filter already exists, update it. Otherwise, create it.
		jwtAuthn, jwtAuthnIndex, err := findJWTAuthnFilter(httpConManager.HttpFilters)
		if err != nil {
			s.cfg.Logger.Sugar().Error("failed to unmarshal the existing basicAuth filter", slog.Any("error", err))
			continue
		}
		if jwtAuthnIndex == -1 {

		} else {
			// Update the basic auth filter
			for providerKey, provider := range jwtAuthn.Providers {
				// Update the provider with the new passwords
				jwksSourceSpecifier := provider.GetJwksSourceSpecifier()
				if remoteJwks, ok := jwksSourceSpecifier.(*jwt_authnv3.JwtProvider_RemoteJwks); ok {
					remoteJwks.RemoteJwks.AsyncFetch = nil
					provider.JwksSourceSpecifier = remoteJwks
					jwtAuthn.Providers[providerKey] = provider
				} else {
					// Not a RemoteJwks, maybe it's a LocalJwks or nil
					s.cfg.Logger.Sugar().Info("jwksSourceSpecifier is not of type *JwtProvider_RemoteJwks")
				}
				provider.PayloadInMetadata = constants.JWTAuthnPayloadInMetadata
			}
		}
		// Add or update the Basic Authentication filter in the HCM
		anyJWTAuthnFilter, _ := anypb.New(jwtAuthn)
		if jwtAuthnIndex > -1 {
			httpConManager.HttpFilters[jwtAuthnIndex].ConfigType = &hcm.HttpFilter_TypedConfig{
				TypedConfig: anyJWTAuthnFilter,
			}
		}

		// Write the updated HCM back to the filter chain
		anyConnectionMgr, _ := anypb.New(httpConManager)
		currChain.Filters[hcmIndex].ConfigType = &listenerv3.Filter_TypedConfig{
			TypedConfig: anyConnectionMgr,
		}
	}
	return &pb.PostHTTPListenerModifyResponse{
		Listener: req.Listener,
	}, nil
}

// Tries to find an HTTP connection manager in the provided filter chain.
func findHCM(filterChain *listenerv3.FilterChain) (*hcm.HttpConnectionManager, int, error) {
	for filterIndex, filter := range filterChain.Filters {
		if filter.Name == wellknown.HTTPConnectionManager {
			hcm := new(hcm.HttpConnectionManager)
			if err := filter.GetTypedConfig().UnmarshalTo(hcm); err != nil {
				return nil, -1, err
			}
			return hcm, filterIndex, nil
		}
	}
	return nil, -1, fmt.Errorf("unable to find HTTPConnectionManager in FilterChain: %s", filterChain.Name)
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
	req.Route.Metadata.FilterMetadata[constants.ExternalProcessingNamespace] = &structpb.Struct{
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

// Tries to find the Basic Authentication HTTP filter in the provided chain
func findJWTAuthnFilter(chain []*hcm.HttpFilter) (*jwt_authnv3.JwtAuthentication, int, error) {
	for i, filter := range chain {
		if filter.Name == "envoy.filters.http.jwt_authn" {
			jf := new(jwt_authnv3.JwtAuthentication)
			if err := filter.GetTypedConfig().UnmarshalTo(jf); err != nil {
				return nil, -1, err
			}
			return jf, i, nil
		}
	}
	return nil, -1, nil
}
