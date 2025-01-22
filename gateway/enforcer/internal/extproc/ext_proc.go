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

package extproc

import (
	"fmt"
	"io"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/authorization"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/ratelimit"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/requesthandler"
	"github.com/wso2/apk/gateway/enforcer/internal/util"

	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// ExternalProcessingServer represents a server for handling external processing requests.
// It contains a logger for logging purposes.
type ExternalProcessingServer struct {
	log                 logging.Logger
	apiStore            *datastore.APIStore
	ratelimitHelper     *ratelimit.AIRatelimitHelper
	requestConfigHolder *requestconfig.Holder
}

const (
	pathAttribute                                   string = "path"
	vHostAttribute                                  string = "vHost"
	basePathAttribute                               string = "basePath"
	methodAttribute                                 string = "method"
	apiVersionAttribute                             string = "version"
	apiNameAttribute                                string = "name"
	clusterNameAttribute                            string = "clusterName"
	enableBackendBasedAIRatelimitAttribute          string = "enableBackendBasedAIRatelimit"
	backendBasedAIRatelimitDescriptorValueAttribute string = "backendBasedAIRatelimitDescriptorValue"
	
)
var httpHandler requesthandler.HTTP = requesthandler.HTTP{}

// StartExternalProcessingServer initializes and starts the external processing server.
// It creates a gRPC server using the provided configuration and registers the external
// processor server with it.
//
// Parameters:
//   - cfg: A pointer to the Server configuration which includes paths to the enforcer's
//     public and private keys, and a logger instance.
//
// If there is an error during the creation of the gRPC server, the function will panic.
func StartExternalProcessingServer(cfg *config.Server, apiStore *datastore.APIStore) {
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

	ratelimitHelper := ratelimit.NewAIRatelimitHelper(cfg)
	envoy_service_proc_v3.RegisterExternalProcessorServer(server, &ExternalProcessingServer{cfg.Logger, apiStore, ratelimitHelper, nil})
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ExternalProcessingPort))
	if err != nil {
		cfg.Logger.Error(err, fmt.Sprintf("Failed to listen on port: %s", cfg.ExternalProcessingPort))
	}
	cfg.Logger.Info("Starting to serve external processing server")
	if err := server.Serve(listener); err != nil {
		cfg.Logger.Error(err, "Failed to serve grpc server")
	}
}

// Process handles the external processing server stream. It continuously receives
// requests from the stream, processes them, and sends back appropriate responses.
// The function supports different types of processing requests including request headers,
// response headers, request body, and response body.
//
// Parameters:
// - srv: The stream server for processing external requests.
//
// Returns:
// - error: Returns an error if the context is done or if there is an issue receiving or sending the stream request.
//
// The function processes the following request types:
// - envoy_service_proc_v3.ProcessingRequest_RequestHeaders: Logs and processes request headers.
// - envoy_service_proc_v3.ProcessingRequest_ResponseHeaders: Logs and processes response headers.
// - envoy_service_proc_v3.ProcessingRequest_RequestBody: Logs and processes request body.
// - envoy_service_proc_v3.ProcessingRequest_ResponseBody: Logs and processes response body.
//
// If an unknown request type is received, it logs the unknown request type.
func (s *ExternalProcessingServer) Process(srv envoy_service_proc_v3.ExternalProcessor_ProcessServer) error {
	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}

		resp := &envoy_service_proc_v3.ProcessingResponse{}
		// log req.Attributes
		s.log.Info(fmt.Sprintf("Attributes: %+v", req.Attributes))
		s.requestConfigHolder = &requestconfig.Holder{}
		switch v := req.Request.(type) {
		case *envoy_service_proc_v3.ProcessingRequest_RequestHeaders:
			attributes, err := extractExternalProcessingAttributes(req.GetAttributes())
			if err != nil {
				s.log.Error(err, "failed to extract context attributes")
				resp = &envoy_service_proc_v3.ProcessingResponse{
					Response: &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
						ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
							Status: &v32.HttpStatus{
								Code: v32.StatusCode_NotFound,
							},
							Body: []byte("The requested resource is not available."),
							Details: "Could not find the required attributes in the request.",
						},
					},
				}
				break
			}
			s.requestConfigHolder.MatchedAPI = s.apiStore.GetMatchedAPI(util.PrepareAPIKey(attributes.VHost, attributes.BasePath, attributes.APIVersion))
			s.requestConfigHolder.ExternalProcessingEnvoyAttributes = attributes
			s.requestConfigHolder.MatchedResource = httpHandler.GetMatchedResource(s.requestConfigHolder.MatchedAPI, *s.requestConfigHolder.ExternalProcessingEnvoyAttributes)
			s.log.Info(fmt.Sprintf("Matched Resource: %v", s.requestConfigHolder.MatchedResource))
			
			if immediateResponse := authorization.Do(s.requestConfigHolder); immediateResponse != nil {
				resp = &envoy_service_proc_v3.ProcessingResponse{
					Response: &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
						ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
							Status: &v32.HttpStatus{
								Code: v32.StatusCode(immediateResponse.StatusCode),
							},
							Body: []byte(immediateResponse.Message),
						},
					},
				}
				break
			}

			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
					HeaderMutation: &envoy_service_proc_v3.HeaderMutation{
						SetHeaders: []*corev3.HeaderValueOption{
							{
								Header: &corev3.HeaderValue{
									Key:      "x-wso2-cluster-header",
									RawValue: []byte(attributes.ClusterName),
								},
							},
						},
					},
					// This is necessary if the remote server modified headers that are used to calculate the route.
					ClearRouteCache: true,
				},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_RequestHeaders{
					RequestHeaders: rhq,
				},
			}
			break
		case *envoy_service_proc_v3.ProcessingRequest_ResponseHeaders:
			// s.log.Info(fmt.Sprintf("response header %+v, attributes %+v, addr: %+v", v.ResponseHeaders, s.externalProcessingEnvoyAttributes, s))
			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_ResponseHeaders{
					ResponseHeaders: rhq,
				},
			}
			// s.log.Info(fmt.Sprintf("Matched api: %s", s.matchedAPI))
			if s.requestConfigHolder != nil &&
				s.requestConfigHolder.MatchedAPI != nil && 
				s.requestConfigHolder.MatchedAPI.AiProvider != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken != nil &&
				s.requestConfigHolder.ExternalProcessingEnvoyAttributes.EnableBackendBasedAIRatelimit == "true" &&
				s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken.In == "Header" {
				s.log.Info("Backend based AI rate limit enabled using headers")
				tokenCount, err := ratelimit.ExtractTokenCountFromExternalProcessingResponseHeaders(req.GetResponseHeaders().GetHeaders().GetHeaders(), s.requestConfigHolder.MatchedAPI.AiProvider.PromptTokens.Value, s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken.Value, s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken.Value, s.requestConfigHolder.MatchedAPI.AiProvider.Model.Value)
				if err != nil {
					s.log.Error(err, "failed to extract token count from response headers")
				} else {
					s.ratelimitHelper.DoAIRatelimit(tokenCount, true, false, s.requestConfigHolder.ExternalProcessingEnvoyAttributes.BackendBasedAIRatelimitDescriptorValue)
				}
			}

			break
		case *envoy_service_proc_v3.ProcessingRequest_ResponseBody:
			// httpBody := req.GetResponseBody()
			s.log.Info(fmt.Sprintf("attribute %+v\n", s.requestConfigHolder.ExternalProcessingEnvoyAttributes))

			rbq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
					ResponseBody: rbq,
				},
			}

			if s.requestConfigHolder != nil &&
				s.requestConfigHolder.MatchedAPI != nil && 
				s.requestConfigHolder.MatchedAPI.AiProvider != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken != nil &&
				s.requestConfigHolder.ExternalProcessingEnvoyAttributes.EnableBackendBasedAIRatelimit == "true" &&
				s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken.In == "Body" {
				s.log.Info("Backend based AI rate limit enabled using body")
				tokenCount, err := ratelimit.ExtractTokenCountFromExternalProcessingResponseBody(req.GetResponseBody().Body, s.requestConfigHolder.MatchedAPI.AiProvider.PromptTokens.Value, s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken.Value, s.requestConfigHolder.MatchedAPI.AiProvider.CompletionToken.Value, s.requestConfigHolder.MatchedAPI.AiProvider.Model.Value)
				if err != nil {
					s.log.Error(err, "failed to extract token count from response body")
				} else {
					s.ratelimitHelper.DoAIRatelimit(tokenCount, true, false, s.requestConfigHolder.ExternalProcessingEnvoyAttributes.BackendBasedAIRatelimitDescriptorValue)
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_RequestBody:
			// httpBody := req.GetRequestBody()
			// s.log.Info(fmt.Sprint("request body"))
			rbq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_RequestBody{
					RequestBody: rbq,
				},
			}
		default:
			s.log.Info(fmt.Sprintf("Unknown Request type %v\n", v))
		}
		if err := srv.Send(resp); err != nil {
			s.log.Info(fmt.Sprintf("send error %v", err))
		}
	}
}

// extractExternalProcessingAttributes extracts the external processing attributes from the given data.
func extractExternalProcessingAttributes(data map[string]*structpb.Struct) (*dto.ExternalProcessingEnvoyAttributes, error) {

	// Get the fields from the map
	extProcData, exists := data["envoy.filters.http.ext_proc"]
	if !exists {
		return nil, fmt.Errorf("key envoy.filters.http.ext_proc not found")
	}

	// Extract the "fields" and iterate over them
	attributes := &dto.ExternalProcessingEnvoyAttributes{}
	fields := extProcData.Fields

	if field, ok := fields["request.method"]; ok {
		method := field.GetStringValue()
		attributes.RequestMehod = method
		fmt.Printf("*******   %s\n", method)
	}

	// We need to navigate through the nested fields to get the actual values
	if field, ok := fields["xds.route_metadata"]; ok {

		filterMetadata := field.GetStringValue()
		var structData corev3.Metadata
		err := prototext.Unmarshal([]byte(filterMetadata), &structData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Protobuf text: %v", err)
		}

		// Extract values for predefined keys
		extractedValues := make(map[string]string)

		keysToExtract := []string{
			pathAttribute,
			vHostAttribute,
			basePathAttribute,
			methodAttribute,
			apiVersionAttribute,
			apiNameAttribute,
			clusterNameAttribute,
			enableBackendBasedAIRatelimitAttribute,
			backendBasedAIRatelimitDescriptorValueAttribute,
		}

		for _, key := range keysToExtract {
			if field, exists := structData.FilterMetadata["envoy.filters.http.ext_proc"]; exists {
				extractedValues[key] = field.Fields[key].GetStringValue()
				// case condition to populate ExternalProcessingEnvoyAttributes
				switch key {
				case pathAttribute:
					attributes.Path = extractedValues[key]
				case vHostAttribute:
					attributes.VHost = extractedValues[key]
				case basePathAttribute:
					attributes.BasePath = extractedValues[key]
				case methodAttribute:
					attributes.Method = extractedValues[key]
				case apiVersionAttribute:
					attributes.APIVersion = extractedValues[key]
				case apiNameAttribute:
					attributes.APIName = extractedValues[key]
				case clusterNameAttribute:
					attributes.ClusterName = extractedValues[key]
				case enableBackendBasedAIRatelimitAttribute:
					attributes.EnableBackendBasedAIRatelimit = extractedValues[key]
				case backendBasedAIRatelimitDescriptorValueAttribute:
					attributes.BackendBasedAIRatelimitDescriptorValue = extractedValues[key]
				}
			}
		}

		// Print extracted values
		for key, value := range extractedValues {
			fmt.Printf("%s: %s\n", key, value)
		}
		// Return the populated struct
		return attributes, nil
	}

	// Key not found
	return nil, fmt.Errorf("key xds.route_metadata not found")
}
