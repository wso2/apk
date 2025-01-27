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
	log                              logging.Logger
	apiStore                         *datastore.APIStore
	subscriptionApplicationDatastore *datastore.SubscriptionApplicationDataStore
	ratelimitHelper                  *ratelimit.AIRatelimitHelper
	requestConfigHolder              *requestconfig.Holder
	cfg                              *config.Server
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
	promptTokenCountMetadataKey                     string = "aitoken:prompttokencount"
	completionTokenCountMetadataKey                 string = "aitoken:completiontokencount"
	totalTokenCountMetadataKey                      string = "aitoken:totaltokencount"
	modelIDMetadataKey                              string = "aitoken:modelid"
	modelMetadataKey                                string = "aitoken:model"
	aiProviderNameMetadataKey                       string = "ai:providername"
	aiProviderAPIVersionMetadataKey                 string = "ai:providerversion"
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
func StartExternalProcessingServer(cfg *config.Server, apiStore *datastore.APIStore, subAppDatastore *datastore.SubscriptionApplicationDataStore, modelBasedRoundRobinTracker *datastore.ModelBasedRoundRobinTracker) {
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
	envoy_service_proc_v3.RegisterExternalProcessorServer(server, &ExternalProcessingServer{cfg.Logger, apiStore, subAppDatastore, ratelimitHelper, nil, cfg, modelBasedRoundRobinTracker})
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
			s.requestConfigHolder = &requestconfig.Holder{}
			attributes, err := extractExternalProcessingAttributes(req.GetAttributes())
			if err != nil {
				s.log.Error(err, "failed to extract context attributes")
				resp = &envoy_service_proc_v3.ProcessingResponse{
					Response: &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
						ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
							Status: &v32.HttpStatus{
								Code: v32.StatusCode_NotFound,
							},
							Body:    []byte("The requested resource is not available."),
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
			s.log.Info(fmt.Sprintf("req holder: %+v\n s: %+v", &s.requestConfigHolder, &s))

			if immediateResponse := authorization.Validate(s.requestConfigHolder, s.subscriptionApplicationDatastore, s.cfg); immediateResponse != nil {
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
			resp = &envoy_service_proc_v3.ProcessingResponse{}
			if s.requestConfigHolder.MatchedSubscription != nil && s.requestConfigHolder.MatchedSubscription.RatelimitTier != "Unlimited" && s.requestConfigHolder.MatchedSubscription.RatelimitTier != "" {
				dynamicMetadataKeyValuePairs[subscriptionMetadataKey] = s.requestConfigHolder.MatchedSubscription.UUID
				dynamicMetadataKeyValuePairs[usagePolicyMetadataKey] = s.requestConfigHolder.MatchedSubscription.RatelimitTier
				dynamicMetadataKeyValuePairs[organizationMetadataKey] = s.requestConfigHolder.MatchedAPI.OrganizationID
				dynamicMetadataKeyValuePairs[orgAndRLPolicyMetadataKey] = fmt.Sprintf("%s-%s", s.requestConfigHolder.MatchedAPI.OrganizationID, s.requestConfigHolder.MatchedSubscription.RatelimitTier)
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

		case *envoy_service_proc_v3.ProcessingRequest_RequestBody:
			// httpBody := req.GetRequestBody()
			s.log.Info("Request Body Flow")

			s.log.Info(fmt.Sprintf("Matched api: %v", s.requestConfigHolder.MatchedAPI))
			if s.requestConfigHolder != nil &&
				s.requestConfigHolder.MatchedAPI != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider.SupportedModels != nil &&
				s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin != nil &&
				s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin.Enabled {
				s.log.Info("Model Based Round Robin enabled")
				supportedModels := s.requestConfigHolder.MatchedAPI.AiProvider.SupportedModels
				onQuotaExceedSuspendDuration := s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin.OnQuotaExceedSuspendDuration
				modelWeight := s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin.Models
				// convert to datastore.ModelWeight
				var modelWeights []datastore.ModelWeight
				for _, model := range modelWeight {
					modelWeights = append(modelWeights, datastore.ModelWeight{
						Name:   model.Model,
						Weight: model.Weight,
					})
				}
				s.log.Sugar().Debugf(fmt.Sprintf("Supported Models: %v", supportedModels))
				s.log.Sugar().Debugf(fmt.Sprintf("Model Weights: %v", modelWeight))
				s.log.Sugar().Debugf(fmt.Sprintf("On Quota Exceed Suspend Duration: %v", onQuotaExceedSuspendDuration))
				selectedModel := s.modelBasedRoundRobinTracker.GetNextModel(s.requestConfigHolder.MatchedAPI.UUID, s.requestConfigHolder.MatchedResource.Path, modelWeights)
				s.log.Info(fmt.Sprintf("Selected Model: %v", selectedModel))
				if selectedModel == "" {
					s.log.Info("Unable to select a model since all models are suspended. Continue with the user provided model")
				} else {
					// change request body to model to selected model
					httpBody := req.GetRequestBody().Body
					s.log.Info(fmt.Sprintf("request body before %+v\n", httpBody))
					// Define a map to hold the JSON data
					var jsonData map[string]interface{}
					// Unmarshal the JSON data into the map
					err := json.Unmarshal(httpBody, &jsonData)
					if err != nil {
						s.log.Error(err, "Error unmarshaling JSON Reuqest Body")
					}
					s.log.Info(fmt.Sprintf("jsonData %+v\n", jsonData))
					// Change the model to the selected model
					jsonData["model"] = selectedModel
					// Convert the JSON object to a []byte
					newHTTPBody, err := json.Marshal(jsonData)
					if err != nil {
						s.log.Error(err, "Error marshaling JSON")
					}

					// Calculate the new body length
					newBodyLength := len(newHTTPBody)
					s.log.Info(fmt.Sprintf("new body length: %d\n", newBodyLength))

					// Update the Content-Length header
					headers := &envoy_service_proc_v3.HeaderMutation{
						SetHeaders: []*corev3.HeaderValueOption{
							{
								Header: &corev3.HeaderValue{
									Key:      "Content-Length",
									RawValue: []byte(fmt.Sprintf("%d", newBodyLength)), // Set the new Content-Length
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
					s.log.Info(fmt.Sprintf("rbq %+v\n", rbq))
					resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
						RequestBody: rbq,
					}
					s.log.Info(fmt.Sprintf("resp %+v\n", resp))
					//req.GetRequestBody().Body = newHTTPBody
					s.log.Info(fmt.Sprintf("request body after %+v\n", newHTTPBody))
				}
			}
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
			s.log.Info("Response Header Flow")
			matchedAPI := s.requestConfigHolder.MatchedAPI
			if s.requestConfigHolder != nil &&
				matchedAPI != nil &&
				matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.CompletionToken != nil &&
				s.requestConfigHolder.ExternalProcessingEnvoyAttributes.EnableBackendBasedAIRatelimit == "true" &&
				matchedAPI.AiProvider.CompletionToken.In == dto.InHeader {
				s.log.Info("Backend based AI rate limit enabled using headers")
				tokenCount, err := ratelimit.ExtractTokenCountFromExternalProcessingResponseHeaders(req.GetResponseHeaders().GetHeaders().GetHeaders(),
					matchedAPI.AiProvider.PromptTokens.Value,
					matchedAPI.AiProvider.CompletionToken.Value,
					matchedAPI.AiProvider.CompletionToken.Value,
					matchedAPI.AiProvider.Model.Value)
				if err != nil {
					s.log.Error(err, "failed to extract token count from response headers")
				} else {
					s.ratelimitHelper.DoAIRatelimit(tokenCount, true,
						matchedAPI.DoSubscriptionAIRLInHeaderReponse,
						s.requestConfigHolder.ExternalProcessingEnvoyAttributes.BackendBasedAIRatelimitDescriptorValue,
						s.requestConfigHolder.MatchedSubscription, s.requestConfigHolder.MatchedApplication)
					aiProvider := matchedAPI.AiProvider
					dynamicMetadataKeyValuePairs[aiProviderAPIVersionMetadataKey] = aiProvider.ProviderAPIVersion
					dynamicMetadataKeyValuePairs[aiProviderNameMetadataKey] = aiProvider.ProviderName
					dynamicMetadataKeyValuePairs[modelIDMetadataKey] = tokenCount.Model
					dynamicMetadataKeyValuePairs[completionTokenCountMetadataKey] = strconv.Itoa(tokenCount.Completion)
					dynamicMetadataKeyValuePairs[totalTokenCountMetadataKey] = strconv.Itoa(tokenCount.Total)
					dynamicMetadataKeyValuePairs[promptTokenCountMetadataKey] = strconv.Itoa(tokenCount.Prompt)
				}
			}
			if s.requestConfigHolder != nil &&
				s.requestConfigHolder.MatchedAPI != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider.SupportedModels != nil &&
				s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin != nil &&
				s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin.Enabled {
				s.log.Info("Model Based Round Robin enabled")
				headerValues := req.GetResponseHeaders().GetHeaders().GetHeaders()
				s.log.Info(fmt.Sprintf("Header Values: %v", headerValues))
				remainingTokenCount := 100
				remainingRequestCount := 100
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
				}
				if remainingTokenCount <= 50 || remainingRequestCount <= 50 { // Suspend model if token/request count reaches 0
					s.log.Info("Token/request are exhausted. Suspending the model")
					s.requestConfigHolder.ExternalProcessingEnvoyAttributes.SuspendAIModel = "true"
				}
			}
		case *envoy_service_proc_v3.ProcessingRequest_ResponseBody:
			// httpBody := req.GetResponseBody()
			s.log.Info(fmt.Sprintf("req holder: %+v\n s: %+v", &s.requestConfigHolder, &s))
			s.log.Info("Response Body Flow")

			rbq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{},
			}
			resp = &envoy_service_proc_v3.ProcessingResponse{
				Response: &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
					ResponseBody: rbq,
				},
			}
			matchedAPI := s.requestConfigHolder.MatchedAPI
			if s.requestConfigHolder != nil &&
				matchedAPI != nil &&
				matchedAPI.AiProvider != nil &&
				matchedAPI.AiProvider.CompletionToken != nil &&
				s.requestConfigHolder.ExternalProcessingEnvoyAttributes.EnableBackendBasedAIRatelimit == "true" &&
				matchedAPI.AiProvider.CompletionToken.In == dto.InBody {
				s.log.Info("Backend based AI rate limit enabled using body")
				tokenCount, err := ratelimit.ExtractTokenCountFromExternalProcessingResponseBody(req.GetResponseBody().Body,
					matchedAPI.AiProvider.PromptTokens.Value,
					matchedAPI.AiProvider.CompletionToken.Value,
					matchedAPI.AiProvider.CompletionToken.Value,
					matchedAPI.AiProvider.Model.Value)
				if err != nil {
					s.log.Error(err, "failed to extract token count from response body")
				} else {
					s.ratelimitHelper.DoAIRatelimit(tokenCount, true,
						matchedAPI.DoSubscriptionAIRLInBodyReponse,
						s.requestConfigHolder.ExternalProcessingEnvoyAttributes.BackendBasedAIRatelimitDescriptorValue,
						s.requestConfigHolder.MatchedSubscription, s.requestConfigHolder.MatchedApplication)
					aiProvider := matchedAPI.AiProvider
					dynamicMetadataKeyValuePairs[aiProviderAPIVersionMetadataKey] = aiProvider.ProviderAPIVersion
					dynamicMetadataKeyValuePairs[aiProviderNameMetadataKey] = aiProvider.ProviderName
					dynamicMetadataKeyValuePairs[modelIDMetadataKey] = tokenCount.Model
					dynamicMetadataKeyValuePairs[completionTokenCountMetadataKey] = strconv.Itoa(tokenCount.Completion)
					dynamicMetadataKeyValuePairs[totalTokenCountMetadataKey] = strconv.Itoa(tokenCount.Total)
					dynamicMetadataKeyValuePairs[promptTokenCountMetadataKey] = strconv.Itoa(tokenCount.Prompt)

				}
			}

			if s.requestConfigHolder != nil &&
				s.requestConfigHolder.MatchedAPI != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider != nil &&
				s.requestConfigHolder.MatchedAPI.AiProvider.SupportedModels != nil &&
				s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin != nil &&
				s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin.Enabled &&
				s.requestConfigHolder.ExternalProcessingEnvoyAttributes.SuspendAIModel == "true" {
				s.log.Info("Model Based Round Robin enabled")
				httpBody := req.GetResponseBody().Body
				// Define a map to hold the JSON data
				var jsonData map[string]interface{}
				// Unmarshal the JSON data into the map
				err := json.Unmarshal(httpBody, &jsonData)
				if err != nil {
					s.log.Error(err, "Error unmarshaling JSON Response Body")
				}
				s.log.Info(fmt.Sprintf("jsonData %+v\n", jsonData))
				// Retrieve Model from the JSON data
				model := ""
				if modelValue, ok := jsonData["model"].(string); ok {
					model = modelValue
				} else {
					s.log.Error(fmt.Errorf("model is not a string"), "failed to extract model from JSON data")
				}
				s.log.Info("Suspending model: " + model)
				duration := s.requestConfigHolder.MatchedAPI.AIModelBasedRoundRobin.OnQuotaExceedSuspendDuration
				s.modelBasedRoundRobinTracker.SuspendModel(s.requestConfigHolder.MatchedAPI.UUID, s.requestConfigHolder.MatchedResource.Path, model, time.Duration(time.Duration(duration*1000*1000*1000)))
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
	if s.requestConfigHolder.MatchedAPI != nil {
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
