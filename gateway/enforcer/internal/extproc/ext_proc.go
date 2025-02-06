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
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/gateway/enforcer/internal/analytics"
	"github.com/wso2/apk/gateway/enforcer/internal/authorization"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/graphql"
	"github.com/wso2/apk/gateway/enforcer/internal/jwtbackend"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/ratelimit"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/requesthandler"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
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
	log                              logging.Logger
	apiStore                         *datastore.APIStore
	subscriptionApplicationDatastore *datastore.SubscriptionApplicationDataStore
	ratelimitHelper                  *ratelimit.AIRatelimitHelper
	requestConfigHolder              *requestconfig.Holder
	cfg                              *config.Server
	jwtTransformer                   *transformer.JWTTransformer
	modelBasedRoundRobinTracker      *datastore.ModelBasedRoundRobinTracker
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
	suspendAIModelValueAttribute                    string = "ai:suspendmodel"
	externalProessingMetadataContextKey             string = "envoy.filters.http.ext_proc"
	subscriptionMetadataKey                         string = "ratelimit:subscription"
	usagePolicyMetadataKey                          string = "ratelimit:usage-policy"
	organizationMetadataKey                         string = "ratelimit:organization"
	orgAndRLPolicyMetadataKey                       string = "ratelimit:organization-and-rlpolicy"
	extractTokenFromMetadataKey                     string = "aitoken:extracttokenfrom"
	promptTokenIDMetadataKey                        string = "aitoken:prompttokenid"
	completionTokenIDMetadataKey                    string = "aitoken:completiontokenid"
	totalTokenIDMetadataKey                         string = "aitoken:totaltokenid"
	matchedAPIMetadataKey                           string = "request:matchedapi"
	matchedResourceMetadataKey                      string = "request:matchedresource"
	matchedSubscriptionMetadataKey                  string = "request:matchedsubscription"
	matchedApplicationMetadataKey                   string = "request:matchedapplication"

	modelMetadataKey string = "aitoken:model"
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
func StartExternalProcessingServer(cfg *config.Server, apiStore *datastore.APIStore, subAppDatastore *datastore.SubscriptionApplicationDataStore, jwtTransformer *transformer.JWTTransformer, modelBasedRoundRobinTracker *datastore.ModelBasedRoundRobinTracker) {
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
	envoy_service_proc_v3.RegisterExternalProcessorServer(server, &ExternalProcessingServer{cfg.Logger, apiStore, subAppDatastore, ratelimitHelper, nil, cfg, jwtTransformer, modelBasedRoundRobinTracker})
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
		dynamicMetadataKeyValuePairs := make(map[string]string)
		switch v := req.Request.(type) {
		case *envoy_service_proc_v3.ProcessingRequest_RequestHeaders:
			requestConfigHolder := &requestconfig.Holder{}
			attributes, err := extractExternalProcessingXDSRouteMetadataAttributes(req.GetAttributes())
			if err != nil {
				s.log.Error(err, "failed to extract context attributes")
				resp = &envoy_service_proc_v3.ProcessingResponse{
					Response: &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
						ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
							Status: &v32.HttpStatus{
								Code: v32.StatusCode_NotFound,
							},
							Body:    []byte("The requested resource is not available."),
							Details: "Resource not found",
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
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestHeaders{
				RequestHeaders: rhq,
			}
			apiKey := util.PrepareAPIKey(attributes.VHost, attributes.BasePath, attributes.APIVersion)
			requestConfigHolder.MatchedAPI = s.apiStore.GetMatchedAPI(util.PrepareAPIKey(attributes.VHost, attributes.BasePath, attributes.APIVersion))
			dynamicMetadataKeyValuePairs[matchedAPIMetadataKey] = apiKey
			requestConfigHolder.ExternalProcessingEnvoyAttributes = attributes
			metadata, err := extractExternalProcessingMetadata(req.GetMetadataContext())
			if err != nil {
				s.log.Error(err, "failed to extract context metadata")
				// return status.Errorf(codes.Unknown, "cannot extract metadata: %v", err)
				break
			}
			requestConfigHolder.ExternalProcessingEnvoyMetadata = metadata
			requestConfigHolder.MatchedResource = httpHandler.GetMatchedResource(requestConfigHolder.MatchedAPI, *requestConfigHolder.ExternalProcessingEnvoyAttributes)
			if requestConfigHolder.MatchedResource != nil {
				requestConfigHolder.MatchedResource.RouteMetadataAttributes = attributes
				dynamicMetadataKeyValuePairs[matchedResourceMetadataKey] = requestConfigHolder.MatchedResource.GetResourceIdentifier()
			}
			// s.log.Info(fmt.Sprintf("Matched api bjc: %v", requestConfigHolder.MatchedAPI.BackendJwtConfiguration))
			// s.log.Info(fmt.Sprintf("Matched Resource: %v", requestConfigHolder.MatchedResource))
			// s.log.Info(fmt.Sprintf("req holderrr: %+v\n s: %+v", &requestConfigHolder, &s))
			s.log.Info(fmt.Sprintf("req holderrr: %+v\n s: %+v", requestConfigHolder, s))
			if requestConfigHolder.MatchedResource != nil && requestConfigHolder.MatchedResource.AuthenticationConfig != nil && !requestConfigHolder.MatchedResource.AuthenticationConfig.Disabled && !requestConfigHolder.MatchedAPI.DisableAuthentication {
				jwtValidationInfo := s.jwtTransformer.TransformJWTClaims(requestConfigHolder.MatchedAPI.OrganizationID, requestConfigHolder.ExternalProcessingEnvoyMetadata)
				requestConfigHolder.JWTValidationInfo = &jwtValidationInfo
				s.log.Sugar().Infof("jwtValidation==%v", jwtValidationInfo)
				if immediateResponse := authorization.Validate(requestConfigHolder, s.subscriptionApplicationDatastore, s.cfg); immediateResponse != nil {
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
				if requestConfigHolder.MatchedSubscription != nil && requestConfigHolder.MatchedSubscription.RatelimitTier != "Unlimited" && requestConfigHolder.MatchedSubscription.RatelimitTier != "" {
					dynamicMetadataKeyValuePairs[subscriptionMetadataKey] = requestConfigHolder.MatchedSubscription.UUID
					dynamicMetadataKeyValuePairs[usagePolicyMetadataKey] = requestConfigHolder.MatchedSubscription.RatelimitTier
					dynamicMetadataKeyValuePairs[organizationMetadataKey] = requestConfigHolder.MatchedAPI.OrganizationID
					dynamicMetadataKeyValuePairs[orgAndRLPolicyMetadataKey] = fmt.Sprintf("%s-%s", requestConfigHolder.MatchedAPI.OrganizationID, requestConfigHolder.MatchedSubscription.RatelimitTier)
				}
			}
			backendJWT := ""
			if requestConfigHolder.MatchedAPI != nil && requestConfigHolder.MatchedAPI.BackendJwtConfiguration != nil && requestConfigHolder.MatchedAPI.BackendJwtConfiguration.Enabled {
				backendJWT = jwtbackend.CreateBackendJWT(requestConfigHolder, s.cfg)
				s.log.Sugar().Infof("generated backendJWT==%v", backendJWT)
			}

			if backendJWT != "" {
				rhq.Response.HeaderMutation.SetHeaders = append(rhq.Response.HeaderMutation.SetHeaders, &corev3.HeaderValueOption{
					Header: &corev3.HeaderValue{
						Key:      requestConfigHolder.MatchedAPI.BackendJwtConfiguration.JWTHeader,
						RawValue: []byte(backendJWT),
					},
				})
				s.cfg.Logger.Info(fmt.Sprintf("Added backend JWT to the header: %s, header name: %s", backendJWT, requestConfigHolder.MatchedAPI.BackendJwtConfiguration.JWTHeader))
			}
			if requestConfigHolder.MatchedApplication != nil {
				dynamicMetadataKeyValuePairs[matchedApplicationMetadataKey] = requestConfigHolder.MatchedApplication.UUID
			}
			if requestConfigHolder.MatchedSubscription != nil {
				dynamicMetadataKeyValuePairs[matchedSubscriptionMetadataKey] = requestConfigHolder.MatchedSubscription.UUID
			}

			if requestConfigHolder.MatchedAPI != nil && requestConfigHolder.MatchedAPI.EndpointSecurity != nil {
				s.cfg.Logger.Info(fmt.Sprintf("Inside API Level Endpoint Security: %+v", requestConfigHolder.MatchedAPI.EndpointSecurity))
				for _, es := range requestConfigHolder.MatchedAPI.EndpointSecurity {
					if es.Enabled {
						s.cfg.Logger.Info(fmt.Sprintf("Enabled API Level Endpoint Security: %+v", es))
						s.cfg.Logger.Info(fmt.Sprintf("Enabled API Level Security Type: %s", es.SecurityType))
						if es.SecurityType == "Basic" {
							basicValue := fmt.Sprintf("Basic %s", util.Base64Encode([]byte(fmt.Sprintf("%s:%s", es.Username, es.Password))))
							rhq.Response.HeaderMutation.SetHeaders = append(rhq.Response.HeaderMutation.SetHeaders, &corev3.HeaderValueOption{
								Header: &corev3.HeaderValue{
									Key:      "Authorization",
									RawValue: []byte(basicValue),
								},
							})
						} else if es.SecurityType == "APIKey" {
							rhq.Response.HeaderMutation.SetHeaders = append(rhq.Response.HeaderMutation.SetHeaders, &corev3.HeaderValueOption{
								Header: &corev3.HeaderValue{
									Key:      es.CustomParameters["key"],
									RawValue: []byte(es.CustomParameters["value"]),
								},
							})
						}
					}
				}
			}

			if requestConfigHolder.MatchedResource != nil && requestConfigHolder.MatchedResource.EndpointSecurity != nil {
				s.cfg.Logger.Info(fmt.Sprintf("Resource Level Endpoint Security: %+v", requestConfigHolder.MatchedResource.EndpointSecurity))
				for _, es := range requestConfigHolder.MatchedResource.EndpointSecurity {
					if es.Enabled {
						s.cfg.Logger.Info(fmt.Sprintf("Resource Level Endpoint Security: %+v", es))
						s.cfg.Logger.Info(fmt.Sprintf("Resource Level Security Type: %s", es.SecurityType))
						if es.SecurityType == "Basic" {
							basicValue := fmt.Sprintf("Basic %s", util.Base64Encode([]byte(fmt.Sprintf("%s:%s", es.Username, es.Password))))
							rhq.Response.HeaderMutation.SetHeaders = append(rhq.Response.HeaderMutation.SetHeaders, &corev3.HeaderValueOption{
								Header: &corev3.HeaderValue{
									Key:      "Authorization",
									RawValue: []byte(basicValue),
								},
							})
						} else if es.SecurityType == "APIKey" {
							rhq.Response.HeaderMutation.SetHeaders = append(rhq.Response.HeaderMutation.SetHeaders, &corev3.HeaderValueOption{
								Header: &corev3.HeaderValue{
									Key:      es.CustomParameters["key"],
									RawValue: []byte(es.CustomParameters["value"]),
								},
							})
						}
					}
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_RequestBody:
			// httpBody := req.GetRequestBody()
			s.log.Info("Request Body Flow")
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
				RequestBody: &envoy_service_proc_v3.BodyResponse{
					Response: &envoy_service_proc_v3.CommonResponse{},
				},
			}
			s.cfg.Logger.Info("Request Body Flow")
			metadata, err := extractExternalProcessingMetadata(req.GetMetadataContext())
			if err != nil {
				s.log.Error(err, "failed to extract context metadata")
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("metadata: %v", metadata))
			matchedAPI := s.apiStore.GetMatchedAPI(metadata.MatchedAPIIdentifier)
			if matchedAPI == nil {
				s.cfg.Logger.Info(fmt.Sprintf("Matched API not found: %s", metadata.MatchedAPIIdentifier))
				break
			}

			if matchedAPI.IsGraphQLAPI() {
				if immediateResponse := graphql.ValidateGraphQLOperation(matchedAPI, s.jwtTransformer, metadata, s.subscriptionApplicationDatastore, s.cfg, string(req.GetRequestBody().Body)); immediateResponse != nil {
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
			}
			matchedResource := matchedAPI.ResourceMap[metadata.MatchedResourceIdentifier]
			if matchedResource == nil {
				s.cfg.Logger.Info(fmt.Sprintf("Matched Resource not found: %s", metadata.MatchedResourceIdentifier))
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("Matched Resource: %v", matchedResource.RouteMetadataAttributes))

			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.SupportedModels != nil &&
				matchedAPI.AIModelBasedRoundRobin != nil &&
				matchedAPI.AIModelBasedRoundRobin.Enabled {
				s.cfg.Logger.Info("API Level Model Based Round Robin enabled")
				supportedModels := matchedAPI.AiProvider.SupportedModels
				onQuotaExceedSuspendDuration := matchedAPI.AIModelBasedRoundRobin.OnQuotaExceedSuspendDuration
				s.cfg.Logger.Info(fmt.Sprintf("EnvType :%+v", matchedAPI.EnvType))
				var modelWeight []dto.ModelWeight
				if matchedAPI.EnvType != "" && matchedAPI.EnvType == "PRODUCTION" {
					modelWeight = matchedAPI.AIModelBasedRoundRobin.ProductionModels
				} else if matchedAPI.EnvType != "" && matchedAPI.EnvType == "SANDBOX" {
					modelWeight = matchedAPI.AIModelBasedRoundRobin.SandboxModels
				}
				// convert to datastore.ModelWeight
				var modelWeights []datastore.ModelWeight
				for _, model := range modelWeight {
					modelWeights = append(modelWeights, datastore.ModelWeight{
						Name:     model.Model,
						Endpoint: model.Endpoint,
						Weight:   model.Weight,
					})
				}
				s.log.Sugar().Debugf(fmt.Sprintf("Supported Models: %v", supportedModels))
				s.log.Sugar().Debugf(fmt.Sprintf("Model Weights: %v", modelWeight))
				s.log.Sugar().Debugf(fmt.Sprintf("On Quota Exceed Suspend Duration: %v", onQuotaExceedSuspendDuration))
				selectedModel, selectedEndpoint := s.modelBasedRoundRobinTracker.GetNextModel(matchedAPI.UUID, matchedResource.Path, modelWeights)
				s.cfg.Logger.Info(fmt.Sprintf("Selected Model: %v", selectedModel))
				s.cfg.Logger.Info(fmt.Sprintf("Selected Endpoint: %v", selectedEndpoint))
				if selectedModel == "" || selectedEndpoint == "" {
					s.cfg.Logger.Info("Unable to select a model since all models are suspended. Continue with the user provided model")
				} else {
					// change request body to model to selected model
					httpBody := req.GetRequestBody().Body
					s.cfg.Logger.Info(fmt.Sprintf("request body before %+v\n", httpBody))
					// Define a map to hold the JSON data
					var jsonData map[string]interface{}
					// Unmarshal the JSON data into the map
					err := json.Unmarshal(httpBody, &jsonData)
					if err != nil {
						s.log.Error(err, "Error unmarshaling JSON Reuqest Body")
					}
					s.cfg.Logger.Info(fmt.Sprintf("jsonData %+v\n", jsonData))
					// Change the model to the selected model
					jsonData["model"] = selectedModel
					// Convert the JSON object to a []byte
					newHTTPBody, err := json.Marshal(jsonData)
					if err != nil {
						s.log.Error(err, "Error marshaling JSON")
					}

					// Calculate the new body length
					newBodyLength := len(newHTTPBody)
					s.cfg.Logger.Info(fmt.Sprintf("new body length: %d\n", newBodyLength))

					// Update the Content-Length header
					headers := &envoy_service_proc_v3.HeaderMutation{
						SetHeaders: []*corev3.HeaderValueOption{
							{
								Header: &corev3.HeaderValue{
									Key:      "Content-Length",
									RawValue: []byte(fmt.Sprintf("%d", newBodyLength)), // Set the new Content-Length
								},
							},
							{
								Header: &corev3.HeaderValue{
									Key:      "x-wso2-cluster-header",
									RawValue: []byte(selectedEndpoint),
								},
							},
						},
					}

					rbq := &envoy_service_proc_v3.BodyResponse{
						Response: &envoy_service_proc_v3.CommonResponse{
							Status:         envoy_service_proc_v3.CommonResponse_CONTINUE_AND_REPLACE,
							HeaderMutation: headers, // Add header mutation here
							BodyMutation: &envoy_service_proc_v3.BodyMutation{
								Mutation: &envoy_service_proc_v3.BodyMutation_Body{
									Body: newHTTPBody,
								},
							},
						},
					}
					s.cfg.Logger.Info(fmt.Sprintf("rbq %+v\n", rbq))
					resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
						RequestBody: rbq,
					}
					s.cfg.Logger.Info(fmt.Sprintf("resp %+v\n", resp))
					//req.GetRequestBody().Body = newHTTPBody
					s.cfg.Logger.Info(fmt.Sprintf("request body after %+v\n", newHTTPBody))
				}
			}
			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.SupportedModels != nil &&
				matchedAPI.AIModelBasedRoundRobin == nil &&
				matchedResource.AIModelBasedRoundRobin != nil &&
				matchedResource.AIModelBasedRoundRobin.Enabled {
				s.cfg.Logger.Info("Resource Level Model Based Round Robin enabled")
				supportedModels := matchedAPI.AiProvider.SupportedModels
				onQuotaExceedSuspendDuration := matchedResource.AIModelBasedRoundRobin.OnQuotaExceedSuspendDuration
				s.cfg.Logger.Info(fmt.Sprintf("EnvType :%+v", matchedAPI.EnvType))
				var modelWeight []dto.ModelWeight
				if matchedAPI.EnvType != "" && matchedAPI.EnvType == "PRODUCTION" {
					modelWeight = matchedResource.AIModelBasedRoundRobin.ProductionModels
				} else if matchedAPI.EnvType != "" && matchedAPI.EnvType == "SANDBOX" {
					modelWeight = matchedResource.AIModelBasedRoundRobin.SandboxModels
				}
				// convert to datastore.ModelWeight
				var modelWeights []datastore.ModelWeight
				for _, model := range modelWeight {
					modelWeights = append(modelWeights, datastore.ModelWeight{
						Name:     model.Model,
						Endpoint: model.Endpoint,
						Weight:   model.Weight,
					})
				}
				s.log.Sugar().Debugf(fmt.Sprintf("Supported Models: %v", supportedModels))
				s.log.Sugar().Debugf(fmt.Sprintf("Model Weights: %v", modelWeight))
				s.log.Sugar().Debugf(fmt.Sprintf("On Quota Exceed Suspend Duration: %v", onQuotaExceedSuspendDuration))
				selectedModel, selectedEndpoint := s.modelBasedRoundRobinTracker.GetNextModel(matchedAPI.UUID, matchedResource.Path, modelWeights)
				s.cfg.Logger.Info(fmt.Sprintf("Selected Model: %v", selectedModel))
				s.cfg.Logger.Info(fmt.Sprintf("Selected Endpoint: %v", selectedEndpoint))
				if selectedModel == "" || selectedEndpoint == "" {
					s.cfg.Logger.Info("Unable to select a model since all models are suspended. Continue with the user provided model")
				} else {
					// change request body to model to selected model
					httpBody := req.GetRequestBody().Body
					s.cfg.Logger.Info(fmt.Sprintf("request body before %+v\n", httpBody))
					// Define a map to hold the JSON data
					var jsonData map[string]interface{}
					// Unmarshal the JSON data into the map
					err := json.Unmarshal(httpBody, &jsonData)
					if err != nil {
						s.log.Error(err, "Error unmarshaling JSON Reuqest Body")
					}
					s.cfg.Logger.Info(fmt.Sprintf("jsonData %+v\n", jsonData))
					// Change the model to the selected model
					jsonData["model"] = selectedModel
					// Convert the JSON object to a []byte
					newHTTPBody, err := json.Marshal(jsonData)
					if err != nil {
						s.log.Error(err, "Error marshaling JSON")
					}

					// Calculate the new body length
					newBodyLength := len(newHTTPBody)
					s.cfg.Logger.Info(fmt.Sprintf("new body length: %d\n", newBodyLength))

					// Update the Content-Length header
					headers := &envoy_service_proc_v3.HeaderMutation{
						SetHeaders: []*corev3.HeaderValueOption{
							{
								Header: &corev3.HeaderValue{
									Key:      "Content-Length",
									RawValue: []byte(fmt.Sprintf("%d", newBodyLength)), // Set the new Content-Length
								},
							},
							{
								Header: &corev3.HeaderValue{
									Key:      "x-wso2-cluster-header",
									RawValue: []byte(selectedEndpoint),
								},
							},
						},
					}

					rbq := &envoy_service_proc_v3.BodyResponse{
						Response: &envoy_service_proc_v3.CommonResponse{
							Status:         envoy_service_proc_v3.CommonResponse_CONTINUE_AND_REPLACE,
							HeaderMutation: headers, // Add header mutation here
							BodyMutation: &envoy_service_proc_v3.BodyMutation{
								Mutation: &envoy_service_proc_v3.BodyMutation_Body{
									Body: newHTTPBody,
								},
							},
						},
					}
					s.cfg.Logger.Info(fmt.Sprintf("rbq %+v\n", rbq))
					resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
						RequestBody: rbq,
					}
					s.cfg.Logger.Info(fmt.Sprintf("resp %+v\n", resp))
					//req.GetRequestBody().Body = newHTTPBody
					s.cfg.Logger.Info(fmt.Sprintf("request body after %+v\n", newHTTPBody))
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_ResponseHeaders:
			s.log.Info(fmt.Sprintf("response header %+v, ", v.ResponseHeaders))
			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_ResponseHeaders{
					ResponseHeaders: rhq,
				},
			}
			metadata, err := extractExternalProcessingMetadata(req.GetMetadataContext())
			if err != nil {
				s.log.Error(err, "failed to extract context metadata")
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("metadata: %+v", metadata))
			matchedAPI := s.apiStore.GetMatchedAPI(metadata.MatchedAPIIdentifier)
			if matchedAPI == nil {
				s.cfg.Logger.Info(fmt.Sprintf("Matched API not found: %s", metadata.MatchedAPIIdentifier))
				break
			}
			matchedResource := matchedAPI.ResourceMap[metadata.MatchedResourceIdentifier]
			if matchedResource == nil {
				s.cfg.Logger.Info(fmt.Sprintf("Matched Resource not found: %s", metadata.MatchedResourceIdentifier))
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("Matched Resource: %v", matchedResource.RouteMetadataAttributes))
			matchedSubscription := s.subscriptionApplicationDatastore.GetSubscription(matchedAPI.OrganizationID, metadata.MatchedSubscriptionIdentifier)
			matchedApplication := s.subscriptionApplicationDatastore.GetApplication(matchedAPI.OrganizationID, metadata.MatchedApplicationIdentifier)
			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.PromptTokens != nil &&
				matchedAPI.AiProvider.CompletionToken != nil &&
				matchedAPI.AiProvider.TotalToken != nil &&
				matchedResource.RouteMetadataAttributes != nil &&
				matchedResource.RouteMetadataAttributes.EnableBackendBasedAIRatelimit == "true" &&
				matchedAPI.AiProvider.CompletionToken.In == dto.InHeader {
				s.log.Info("Backend based AI rate limit enabled using headers")
				tokenCount, err := ratelimit.ExtractTokenCountFromExternalProcessingResponseHeaders(req.GetResponseHeaders().GetHeaders().GetHeaders(),
					matchedAPI.AiProvider.PromptTokens.Value,
					matchedAPI.AiProvider.CompletionToken.Value,
					matchedAPI.AiProvider.TotalToken.Value,
					matchedAPI.AiProvider.ResponseModel.Value)
				if err != nil {
					s.log.Error(err, "failed to extract token count from response headers")
				} else {
					go s.ratelimitHelper.DoAIRatelimit(*tokenCount, true,
						matchedAPI.DoSubscriptionAIRLInHeaderReponse,
						matchedResource.RouteMetadataAttributes.BackendBasedAIRatelimitDescriptorValue,
						matchedSubscription, matchedApplication)
					aiProvider := matchedAPI.AiProvider
					dynamicMetadataKeyValuePairs[analytics.AIProviderAPIVersionMetadataKey] = aiProvider.ProviderAPIVersion
					dynamicMetadataKeyValuePairs[analytics.AIProviderNameMetadataKey] = aiProvider.ProviderName
					dynamicMetadataKeyValuePairs[analytics.ModelIDMetadataKey] = tokenCount.Model
					dynamicMetadataKeyValuePairs[analytics.CompletionTokenCountMetadataKey] = strconv.Itoa(tokenCount.Completion)
					dynamicMetadataKeyValuePairs[analytics.TotalTokenCountMetadataKey] = strconv.Itoa(tokenCount.Total)
					dynamicMetadataKeyValuePairs[analytics.PromptTokenCountMetadataKey] = strconv.Itoa(tokenCount.Prompt)
				}
			}
			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.SupportedModels != nil &&
				matchedAPI.AIModelBasedRoundRobin != nil &&
				matchedAPI.AIModelBasedRoundRobin.Enabled {
				s.log.Info("API Level Model Based Round Robin enabled")
				headerValues := req.GetResponseHeaders().GetHeaders().GetHeaders()
				s.log.Info(fmt.Sprintf("Header Values: %v", headerValues))
				remainingTokenCount := 100
				remainingRequestCount := 100
				status := 200
				for _, headerValue := range headerValues {
					if headerValue.Key == "x-ratelimit-remaining-tokens" {
						value, err := util.ConvertStringToInt(string(headerValue.RawValue))
						if err != nil {
							s.log.Error(err, "Unable to retrieve remaining token count by header")
						}
						remainingTokenCount = value
					}
					if headerValue.Key == "x-ratelimit-remaining-requests" {
						value, err := util.ConvertStringToInt(string(headerValue.RawValue))
						if err != nil {
							s.log.Error(err, "Unable to retrieve remaining request count by header")
						}
						remainingRequestCount = value
					}
					if headerValue.Key == "status" {
						status, err = util.ConvertStringToInt(string(headerValue.RawValue))
						if err != nil {
							s.log.Error(err, "Unable to retrieve status code by header")
						}
					}
				}
				if remainingTokenCount <= 0 || remainingRequestCount <= 0 || status == 429 { // Suspend model if token/request count reaches 0 or status code is 429
					s.log.Info("Token/request are exhausted. Suspending the model")
					matchedResource.RouteMetadataAttributes.SuspendAIModel = "true"
					matchedAPI.ResourceMap[metadata.MatchedResourceIdentifier] = matchedResource
					s.apiStore.UpdateMatchedAPI(metadata.MatchedAPIIdentifier, matchedAPI)
				}
			}
			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.SupportedModels != nil &&
				matchedAPI.AIModelBasedRoundRobin == nil &&
				matchedResource.AIModelBasedRoundRobin != nil &&
				matchedResource.AIModelBasedRoundRobin.Enabled {
				s.log.Info("Resource Level Model Based Round Robin enabled")
				headerValues := req.GetResponseHeaders().GetHeaders().GetHeaders()
				s.log.Info(fmt.Sprintf("Header Values: %v", headerValues))
				remainingTokenCount := 100
				remainingRequestCount := 100
				status := 200
				for _, headerValue := range headerValues {
					if headerValue.Key == "x-ratelimit-remaining-tokens" {
						value, err := util.ConvertStringToInt(string(headerValue.RawValue))
						if err != nil {
							s.log.Error(err, "Unable to retrieve remaining token count by header")
						}
						remainingTokenCount = value
					}
					if headerValue.Key == "x-ratelimit-remaining-requests" {
						value, err := util.ConvertStringToInt(string(headerValue.RawValue))
						if err != nil {
							s.log.Error(err, "Unable to retrieve remaining request count by header")
						}
						remainingRequestCount = value
					}
					if headerValue.Key == "status" {
						status, err = util.ConvertStringToInt(string(headerValue.RawValue))
						if err != nil {
							s.log.Error(err, "Unable to retrieve status code by header")
						}
					}
				}
				if remainingTokenCount <= 0 || remainingRequestCount <= 0 || status == 429 { // Suspend model if token/request count reaches 0 or status code is 429
					s.log.Info("Token/request are exhausted. Suspending the model")
					matchedResource.RouteMetadataAttributes.SuspendAIModel = "true"
					matchedAPI.ResourceMap[metadata.MatchedResourceIdentifier] = matchedResource
					s.apiStore.UpdateMatchedAPI(metadata.MatchedAPIIdentifier, matchedAPI)
				}
			}
		case *envoy_service_proc_v3.ProcessingRequest_ResponseBody:
			// httpBody := req.GetResponseBody()
			// s.log.Info(fmt.Sprintf("req holder: %+v\n s: %+v", &s.requestConfigHolder, &s))
			s.log.Info("Response Body Flow")

			rbq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
					ResponseBody: rbq,
				},
			}
			metadata, err := extractExternalProcessingMetadata(req.GetMetadataContext())
			if err != nil {
				s.log.Error(err, "failed to extract context metadata")
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("metadata: %v", metadata))
			matchedAPI := s.apiStore.GetMatchedAPI(metadata.MatchedAPIIdentifier)
			if matchedAPI == nil {
				s.cfg.Logger.Info(fmt.Sprintf("Matched API not found: %s", metadata.MatchedAPIIdentifier))
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("Matched API: %+v", matchedAPI))
			matchedResource := matchedAPI.ResourceMap[metadata.MatchedResourceIdentifier]
			if matchedResource == nil {
				s.cfg.Logger.Info(fmt.Sprintf("Matched Resource not found: %s", metadata.MatchedResourceIdentifier))
				break
			}
			s.cfg.Logger.Info(fmt.Sprintf("Matched resource: %+v", matchedResource))
			matchedSubscription := s.subscriptionApplicationDatastore.GetSubscription(matchedAPI.OrganizationID, metadata.MatchedSubscriptionIdentifier)
			matchedApplication := s.subscriptionApplicationDatastore.GetApplication(matchedAPI.OrganizationID, metadata.MatchedApplicationIdentifier)
			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.CompletionToken != nil &&
				matchedAPI.AiProvider.PromptTokens != nil &&
				matchedAPI.AiProvider.TotalToken != nil &&
				matchedResource.RouteMetadataAttributes != nil &&
				matchedResource.RouteMetadataAttributes.EnableBackendBasedAIRatelimit == "true" &&
				matchedAPI.AiProvider.CompletionToken.In == dto.InBody {
				s.log.Info("Backend based AI rate limit enabled using body")
				tokenCount, err := ratelimit.ExtractTokenCountFromExternalProcessingResponseBody(req.GetResponseBody().Body,
					matchedAPI.AiProvider.PromptTokens.Value,
					matchedAPI.AiProvider.CompletionToken.Value,
					matchedAPI.AiProvider.TotalToken.Value,
					matchedAPI.AiProvider.ResponseModel.Value)
				if err != nil {
					s.log.Error(err, "failed to extract token count from response body")
				} else {
					go s.ratelimitHelper.DoAIRatelimit(*tokenCount, true,
						matchedAPI.DoSubscriptionAIRLInBodyReponse,
						matchedResource.RouteMetadataAttributes.BackendBasedAIRatelimitDescriptorValue,
						matchedSubscription, matchedApplication)
					aiProvider := matchedAPI.AiProvider
					dynamicMetadataKeyValuePairs[analytics.AIProviderAPIVersionMetadataKey] = aiProvider.ProviderAPIVersion
					dynamicMetadataKeyValuePairs[analytics.AIProviderNameMetadataKey] = aiProvider.ProviderName
					dynamicMetadataKeyValuePairs[analytics.ModelIDMetadataKey] = tokenCount.Model
					dynamicMetadataKeyValuePairs[analytics.CompletionTokenCountMetadataKey] = strconv.Itoa(tokenCount.Completion)
					dynamicMetadataKeyValuePairs[analytics.TotalTokenCountMetadataKey] = strconv.Itoa(tokenCount.Total)
					dynamicMetadataKeyValuePairs[analytics.PromptTokenCountMetadataKey] = strconv.Itoa(tokenCount.Prompt)

				}
			}

			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.SupportedModels != nil &&
				matchedAPI.AIModelBasedRoundRobin != nil &&
				matchedAPI.AIModelBasedRoundRobin.Enabled &&
				matchedResource.RouteMetadataAttributes.SuspendAIModel == "true" {
				s.cfg.Logger.Info("API Level Model Based Round Robin enabled")
				httpBody := req.GetResponseBody().Body
				// Define a map to hold the JSON data
				var jsonData map[string]interface{}
				// Unmarshal the JSON data into the map
				err := json.Unmarshal(httpBody, &jsonData)
				if err != nil {
					s.cfg.Logger.Error(err, "Error unmarshaling JSON Response Body")
				}
				s.cfg.Logger.Info(fmt.Sprintf("jsonData %+v\n", jsonData))
				// Retrieve Model from the JSON data
				model := ""
				if modelValue, ok := jsonData["model"].(string); ok {
					model = modelValue
				} else {
					s.cfg.Logger.Error(fmt.Errorf("model is not a string"), "failed to extract model from JSON data")
				}
				s.cfg.Logger.Info("Suspending model: " + model)
				duration := matchedAPI.AIModelBasedRoundRobin.OnQuotaExceedSuspendDuration
				s.modelBasedRoundRobinTracker.SuspendModel(matchedAPI.UUID, matchedResource.Path, model, time.Duration(time.Duration(duration*1000*1000*1000)))
			}
			if matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.SupportedModels != nil &&
				matchedAPI.AIModelBasedRoundRobin == nil &&
				matchedResource.AIModelBasedRoundRobin != nil &&
				matchedResource.AIModelBasedRoundRobin.Enabled &&
				matchedResource.RouteMetadataAttributes.SuspendAIModel == "true" {
				s.cfg.Logger.Info("Resource Level Model Based Round Robin enabled")
				httpBody := req.GetResponseBody().Body
				// Define a map to hold the JSON data
				var jsonData map[string]interface{}
				// Unmarshal the JSON data into the map
				err := json.Unmarshal(httpBody, &jsonData)
				if err != nil {
					s.cfg.Logger.Error(err, "Error unmarshaling JSON Response Body")
				}
				s.cfg.Logger.Info(fmt.Sprintf("jsonData %+v\n", jsonData))
				// Retrieve Model from the JSON data
				model := ""
				if modelValue, ok := jsonData["model"].(string); ok {
					model = modelValue
				} else {
					s.cfg.Logger.Error(fmt.Errorf("model is not a string"), "failed to extract model from JSON data")
				}
				s.cfg.Logger.Info("Suspending model: " + model)
				duration := matchedResource.AIModelBasedRoundRobin.OnQuotaExceedSuspendDuration
				s.modelBasedRoundRobinTracker.SuspendModel(matchedAPI.UUID, matchedResource.Path, model, time.Duration(time.Duration(duration*1000*1000*1000)))
			}
		default:
			s.log.Info(fmt.Sprintf("Unknown Request type %v\n", v))
		}
		// Set dynamic metadata
		dynamicMetadata, err := buildDynamicMetadata(s.prepareMetadataKeyValuePairAndAddTo(dynamicMetadataKeyValuePairs))
		if err != nil {
			s.log.Error(err, "failed to build dynamic metadata")
		} else {
			resp.DynamicMetadata = dynamicMetadata
		}
		if err := srv.Send(resp); err != nil {
			s.log.Info(fmt.Sprintf("send error %v", err))
		}
	}
}

// extractExternalProcessingMetadata extracts the external processing metadata from the given data.
func extractExternalProcessingMetadata(data *corev3.Metadata) (*dto.ExternalProcessingEnvoyMetadata, error) {
	filterMatadata := data.GetFilterMetadata()
	if filterMatadata != nil {
		externalProcessingEnvoyMetadata := &dto.ExternalProcessingEnvoyMetadata{}
		jwtFilterdata := filterMatadata["envoy.filters.http.jwt_authn"]
		if jwtFilterdata != nil {
			jwtFields, exists := jwtFilterdata.Fields["payload_in_metadata"]
			authenticationData := &dto.JwtAuthenticationData{}
			if exists {
				jwtPayload := jwtFields.GetStructValue()
				if jwtPayload != nil {
					claims := make(map[string]interface{})
					for key, value := range jwtPayload.GetFields() {
						if value != nil {
							if key == "iss" {
								authenticationData.Issuer = value.GetStringValue()
							}
							switch value.Kind.(type) {
							case *structpb.Value_StringValue:
								claims[key] = value.GetStringValue()
							case *structpb.Value_NumberValue:
								claims[key] = value.GetNumberValue()
							case *structpb.Value_BoolValue:
								claims[key] = value.GetBoolValue()
							case *structpb.Value_ListValue:
								claims[key] = value.GetListValue()
							}
						}
					}
					authenticationData.Claims = claims
					fmt.Printf("claims: %v\n", claims)
				}
			}
			failureStatusFields, exists := jwtFilterdata.Fields["failed_status"]
			if exists {
				failureStatusStruct := failureStatusFields.GetStructValue()
				if failureStatusStruct != nil {
					code := failureStatusStruct.Fields["code"].GetNumberValue()
					message := failureStatusStruct.Fields["message"].GetStringValue()
					authenticationData.Status = &dto.Status{Code: int(code), Message: message}
				}
			}
			externalProcessingEnvoyMetadata.JwtAuthenticationData = authenticationData
		}
		if extProcMetadata, exists := filterMatadata[externalProessingMetadataContextKey]; exists {
			if matchedAPIKey, exists := extProcMetadata.Fields[matchedAPIMetadataKey]; exists {
				externalProcessingEnvoyMetadata.MatchedAPIIdentifier = matchedAPIKey.GetStringValue()
			}
			if matchedResourceKey, exists := extProcMetadata.Fields[matchedResourceMetadataKey]; exists {
				externalProcessingEnvoyMetadata.MatchedResourceIdentifier = matchedResourceKey.GetStringValue()
			}
			if matchedApplicationKey, exists := extProcMetadata.Fields[matchedApplicationMetadataKey]; exists {
				externalProcessingEnvoyMetadata.MatchedApplicationIdentifier = matchedApplicationKey.GetStringValue()
			}
			if matchedSubscriptionKey, exists := extProcMetadata.Fields[matchedSubscriptionMetadataKey]; exists {
				externalProcessingEnvoyMetadata.MatchedSubscriptionIdentifier = matchedSubscriptionKey.GetStringValue()
			}

		}
		return externalProcessingEnvoyMetadata, nil
	}
	return nil, fmt.Errorf("could not find the filter metadata")
}

// extractExternalProcessingXDSRouteMetadataAttributes extracts the external processing attributes from the given data.
func extractExternalProcessingXDSRouteMetadataAttributes(data map[string]*structpb.Struct) (*dto.ExternalProcessingEnvoyAttributes, error) {

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
		attributes.RequestMethod = method
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
			suspendAIModelValueAttribute,
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
				case suspendAIModelValueAttribute:
					attributes.SuspendAIModel = extractedValues[key]
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

func buildDynamicMetadata(keyValuePairs *map[string]string) (*structpb.Struct, error) {
	// Create the structBuilder
	structBuilder := &structpb.Struct{
		Fields: map[string]*structpb.Value{},
	}

	// Helper function to add metadata
	addMetadata := func(builder *structpb.Struct, key string, value interface{}) error {
		val, err := structpb.NewValue(value)
		if err != nil {
			return err
		}
		builder.Fields[key] = val
		return nil
	}

	for key, value := range *keyValuePairs {
		// Add metadata fields
		if err := addMetadata(structBuilder, key, value); err != nil {
			return nil, err
		}
	}

	// Create the root struct and add the nested struct
	rootStruct := &structpb.Struct{
		Fields: map[string]*structpb.Value{},
	}
	nestedValue := structpb.NewStructValue(structBuilder)
	rootStruct.Fields[externalProessingMetadataContextKey] = nestedValue

	return rootStruct, nil
}

func (s *ExternalProcessingServer) prepareMetadataKeyValuePairAndAddTo(metadataKeyValuePair map[string]string) *map[string]string {
	if s.requestConfigHolder != nil && s.requestConfigHolder.MatchedAPI != nil {
		metadataKeyValuePair[analytics.APIIDKey] = s.requestConfigHolder.MatchedAPI.UUID
		metadataKeyValuePair[analytics.APIContextKey] = s.requestConfigHolder.MatchedAPI.BasePath
		metadataKeyValuePair[organizationMetadataKey] = s.requestConfigHolder.MatchedAPI.OrganizationID
		metadataKeyValuePair[analytics.APINameKey] = s.requestConfigHolder.MatchedAPI.Name
		metadataKeyValuePair[analytics.APIVersionKey] = s.requestConfigHolder.MatchedAPI.Version
		metadataKeyValuePair[analytics.APITypeKey] = s.requestConfigHolder.MatchedAPI.APIType
		// metadataKeyValuePair[analytics.ApiCreatorKey] = s.requestConfigHolder.MatchedAPI.Creator
		// metadataKeyValuePair[analytics.ApiCreatorTenantDomainKey] = s.requestConfigHolder.MatchedAPI.CreatorTenant
		metadataKeyValuePair[analytics.APIOrganizationIDKey] = s.requestConfigHolder.MatchedAPI.OrganizationID

		metadataKeyValuePair[analytics.CorrelationIDKey] = s.requestConfigHolder.ExternalProcessingEnvoyAttributes.CorrelationID
		metadataKeyValuePair[analytics.RegionKey] = s.cfg.EnforcerRegionID
		// metadataKeyValuePair[analytics.UserAgentKey] = s.requestConfigHolder.Metadata.UserAgent
		// metadataKeyValuePair[analytics.ClientIpKey] = s.requestConfigHolder.Metadata.ClientIP
		// metadataKeyValuePair[analytics.ApiResourceTemplateKey] = s.requestConfigHolder.ApiResourceTemplate
		// metadataKeyValuePair[analytics.Destination] = s.requestConfigHolder.Metadata.Destination
		metadataKeyValuePair[analytics.APIEnvironmentKey] = s.requestConfigHolder.MatchedAPI.Environment

		if s.requestConfigHolder.MatchedApplication != nil {
			metadataKeyValuePair[analytics.AppIDKey] = s.requestConfigHolder.MatchedApplication.UUID
			metadataKeyValuePair[analytics.AppUUIDKey] = s.requestConfigHolder.MatchedApplication.UUID
			metadataKeyValuePair[analytics.AppKeyTypeKey] = s.requestConfigHolder.MatchedAPI.EnvType
			metadataKeyValuePair[analytics.AppNameKey] = s.requestConfigHolder.MatchedApplication.Name
			metadataKeyValuePair[analytics.AppOwnerKey] = s.requestConfigHolder.MatchedApplication.Owner
		}
	}
	return &metadataKeyValuePair
}
