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
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/mediation"

	v31 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	v32 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	types "k8s.io/apimachinery/pkg/types"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// ExternalProcessingServer represents a server for handling external processing requests.
// It contains a logger for logging purposes.
type ExternalProcessingServer struct {
	log                              logging.Logger
	subscriptionApplicationDatastore *datastore.SubscriptionApplicationDataStore
	routePolicyAndMetadataDatastore  *datastore.RoutePolicyAndMetadataDataStore
	// ratelimitHelper                  *ratelimit.AIRatelimitHelper
	cfg             *config.Server
	revokedJTIStore *datastore.RevokedJTIStore
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
	customOrgMetadataKey                            string = "customorg"
	endpointBasepath                                string = "endpointBasepath"
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

// StartExternalProcessingServer initializes and starts the external processing server.
// It creates a gRPC server using the provided configuration and registers the external
// processor server with it.
//
// Parameters:
//   - cfg: A pointer to the Server configuration which includes paths to the enforcer's
//     public and private keys, and a logger instance.
//
// If there is an error during the creation of the gRPC server, the function will panic.
func StartExternalProcessingServer(cfg *config.Server,
	subAppDatastore *datastore.SubscriptionApplicationDataStore,
	routePolicyAndMetadataDS *datastore.RoutePolicyAndMetadataDataStore,
	revokedJTIStore *datastore.RevokedJTIStore) {
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

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	cfg.Logger.Info("Health check added.....")
	envoy_service_proc_v3.RegisterExternalProcessorServer(server,
		&ExternalProcessingServer{cfg.Logger,
			subAppDatastore,
			routePolicyAndMetadataDS,
			cfg,
			revokedJTIStore})
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ExternalProcessingPort))
	if err != nil {
		cfg.Logger.Error(err, fmt.Sprintf("Failed to listen on port: %s", cfg.ExternalProcessingPort))
	}
	cfg.Logger.Info(fmt.Sprintf("Starting to serve external processing server on port: %s", cfg.ExternalProcessingPort))
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
	requestConfigHolder := &requestconfig.Holder{}
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
		s.log.Sugar().Debug(fmt.Sprintf("Attributes: %+v", req.Attributes))
		if requestConfigHolder.AttributesPopulated == false {
			extRefs, requestAttributes, envType := s.extractExtensionRefsAndRouteAttributes(req.Attributes)
			requestConfigHolder.EnvType = envType
			requestConfigHolder.RequestAttributes = requestAttributes
			if len(extRefs) > 0 {
				requestConfigHolder.AttributesPopulated = true
				for _, extRef := range extRefs {
					s.log.Sugar().Debug(fmt.Sprintf("Extension Reference: %s", extRef))
					parts := strings.Split(extRef, "/")
					if len(parts) == 3 {
						kind := parts[0]
						namespace := parts[1]
						name := parts[2]
						s.log.Sugar().Debug(fmt.Sprintf("Kind: %s, Namespace: %s, Name: %s", kind, namespace, name))
						if kind == "RoutePolicy" {
							// Fetch the RoutePolicy from the datastore
							namespacedName := types.NamespacedName{
								Name:      name,
								Namespace: namespace,
							}.String()
							routePolicy := s.routePolicyAndMetadataDatastore.GetRoutePolicy(namespacedName)
							if routePolicy != nil {
								s.log.Sugar().Debugf("Found RoutePolicy: %+v", routePolicy)
								if requestConfigHolder.RoutePolicy == nil {
									requestConfigHolder.RoutePolicy = &dpv2alpha1.RoutePolicy{
										Spec: dpv2alpha1.RoutePolicySpec{
											RequestMediation:  make([]*dpv2alpha1.Mediation, 0),
											ResponseMediation: make([]*dpv2alpha1.Mediation, 0),
										},
									}
								}
								for _, reqPolicy := range routePolicy.Spec.RequestMediation {
									requestConfigHolder.RoutePolicy.Spec.RequestMediation = append(requestConfigHolder.RoutePolicy.Spec.RequestMediation, reqPolicy)
								}
								for _, resPolicy := range routePolicy.Spec.ResponseMediation {
									requestConfigHolder.RoutePolicy.Spec.ResponseMediation = append(requestConfigHolder.RoutePolicy.Spec.ResponseMediation, resPolicy)
								}
							} else {
								s.log.Sugar().Errorf("RoutePolicy %s/%s not found", namespace, name)
							}
						} else if kind == "RouteMetadata" {
							// Fetch the RouteMetadata from the datastore
							namespacedName := types.NamespacedName{
								Name:      name,
								Namespace: namespace,
							}.String()
							routeMetadata := s.routePolicyAndMetadataDatastore.GetRouteMetadata(namespacedName)
							if routeMetadata != nil {
								s.log.Sugar().Debugf("Found RouteMetadata: %+v", routeMetadata)
								// We dont support multiple RouteMetadata for a request, Hence the last one will be used
								requestConfigHolder.RouteMetadata = routeMetadata
							} else {
								s.log.Sugar().Errorf("RouteMetadata %s/%s not found", namespace, name)
							}
						} else {
							s.log.Sugar().Debugf("Unknown kind: %s", kind)
						}
					}
				}
			}
		}

		metadata := make(map[string]*structpb.Value)
		switch v := req.Request.(type) {
		case *envoy_service_proc_v3.ProcessingRequest_RequestHeaders:
			s.log.Sugar().Debug("Request Headers Flow")
			s.log.Sugar().Debug(fmt.Sprintf("request header %+v, ", v.RequestHeaders))
			requestConfigHolder.ProcessingPhase = requestconfig.ProcessingPhaseRequestHeaders
			requestConfigHolder.RequestHeaders = req.GetRequestHeaders()
			requestConfigHolder.JWTAuthnPayloaClaims = s.extractJWTAuthnNamespaceData(req.GetMetadataContext())

			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
					ClearRouteCache: true,
				},
			}
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestHeaders{
				RequestHeaders: rhq,
			}

			// Token revocation
			if s.revokedJTIStore != nil {
				if jti, ok := requestConfigHolder.JWTAuthnPayloaClaims["jti"]; ok && jti != nil {
					if jtiStr, ok := jti.(string); ok {
						if s.revokedJTIStore.IsJTIRevoked(jtiStr) {
							s.log.Sugar().Debug("Token is revoked")
							resp.Response = &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
								ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
									Status: &v32.HttpStatus{
										Code: v32.StatusCode_Unauthorized,
									},
									Body:    []byte("Unauthorized: Token is revoked"),
									Details: "revoked_token",
								},
							}
							return srv.Send(resp)
						}
					} else {
						s.log.Sugar().Debug("JTI claim is not a string")
					}
				} else {
					s.log.Sugar().Debug("JTI claim not found")
				}
			}

			for key, value := range requestConfigHolder.RequestHeaders.Headers.Headers {
				s.log.Sugar().Debugf("Request Header: %s: %s", key, value)
			}
			// Override the processing mode based on the attached api policies
			requestBodyMode := v31.ProcessingMode_NONE
			responseHeaderMode := v31.ProcessingMode_SKIP
			responseBodyMode := v31.ProcessingMode_NONE
			if requestConfigHolder.RoutePolicy != nil {
				s.log.Sugar().Debugf("RoutePolicy: %+v", requestConfigHolder.RoutePolicy)
				if requestConfigHolder.RoutePolicy.Spec.RequestMediation != nil {
					s.log.Sugar().Debugf("Request Mediation Policies: %+v", requestConfigHolder.RoutePolicy.Spec.RequestMediation)
					for _, policy := range requestConfigHolder.RoutePolicy.Spec.RequestMediation {
						reqBodyProcessing := mediation.MediationAndRequestBodyProcessing[policy.PolicyName]
						if reqBodyProcessing {
							requestBodyMode = v31.ProcessingMode_BUFFERED
						}
						s.log.Sugar().Debugf("Processing Mode for Policy %s: RequestBodyMode: %s, ResponseHeaderMode: %s, ResponseBodyMode: %s",
							policy.PolicyName,
							requestBodyMode.String(),
							responseHeaderMode.String(),
							responseBodyMode.String())
					}
				} else {
					s.log.Sugar().Debugf("No Request Mediation Policies found in RoutePolicy")
				}
				if requestConfigHolder.RoutePolicy.Spec.ResponseMediation != nil {
					s.log.Sugar().Debugf("Response Mediation Policies: %+v", requestConfigHolder.RoutePolicy.Spec.ResponseMediation)
					for _, policy := range requestConfigHolder.RoutePolicy.Spec.ResponseMediation {
						if respHeaderProcessing, ok := mediation.MediationAndResponseHeaderProcessing[policy.PolicyName]; ok && respHeaderProcessing {
							responseHeaderMode = v31.ProcessingMode_SEND
						}
						if respBodyProcessing, ok := mediation.MediationAndResponseBodyProcessing[policy.PolicyName]; ok && respBodyProcessing {
							responseBodyMode = v31.ProcessingMode_BUFFERED
						}

						s.log.Sugar().Debugf("Processing Mode for Policy %s: ResponseHeaderMode: %s, ResponseBodyMode: %s",
							policy.PolicyName,
							responseHeaderMode.String(),
							responseBodyMode.String())
					}
				} else {
					s.log.Sugar().Debugf("No Response Mediation Policies found in RoutePolicy")
				}
			}
			resp.ModeOverride = &v31.ProcessingMode{
				RequestBodyMode:    requestBodyMode,
				ResponseHeaderMode: responseHeaderMode,
				ResponseBodyMode:   responseBodyMode,
			}

			if requestConfigHolder.RoutePolicy != nil {
				s.log.Sugar().Debugf("RoutePolicy: %+v", requestConfigHolder.RoutePolicy)
				if requestConfigHolder.RoutePolicy.Spec.RequestMediation != nil {
					s.log.Sugar().Debugf("Request Mediation Policies: %+v", requestConfigHolder.RoutePolicy.Spec.RequestMediation)
					for _, policy := range requestConfigHolder.RoutePolicy.Spec.RequestMediation {
						if mediation.MediationAndRequestHeaderProcessing[policy.PolicyName] {
							mediation := mediation.CreateMediation(policy)
							if mediation == nil {
								s.log.Sugar().Errorf("Failed to create mediation for policy: %+v", policy)
								continue
							}
							mediationResult := mediation.Process(requestConfigHolder)
							s.log.Sugar().Debugf("Mediation Result: %+v", mediationResult)
							s.updateRequestConfigBasedOnMediationResults(mediationResult, requestConfigHolder, requestconfig.ProcessingPhaseRequestHeaders)
							stopProcessingMediations := s.processMediationResultAndPrepareResponse(
								mediationResult,
								resp,
								requestconfig.ProcessingPhaseRequestHeaders,
								metadata)
							if stopProcessingMediations {
								s.log.Sugar().Debug("Stopping further processing of request headers due to immediate response")
								break
							}
						}
					}
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_RequestBody:
			s.log.Sugar().Debug("Request Body Flow")
			s.log.Sugar().Debug(fmt.Sprintf("request body %+v, ", v.RequestBody))
			requestConfigHolder.RequestBody = req.GetRequestBody()
			requestConfigHolder.ProcessingPhase = requestconfig.ProcessingPhaseRequestBody

			rhq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
					ClearRouteCache: true,
				},
			}
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{
				RequestBody: rhq,
			}

			if requestConfigHolder.RoutePolicy != nil {
				s.log.Sugar().Debugf("RoutePolicy: %+v", requestConfigHolder.RoutePolicy)
				if requestConfigHolder.RoutePolicy.Spec.RequestMediation != nil {
					s.log.Sugar().Debugf("Request Mediation Policies: %+v", requestConfigHolder.RoutePolicy.Spec.RequestMediation)
					for _, policy := range requestConfigHolder.RoutePolicy.Spec.RequestMediation {
						if mediation.MediationAndRequestBodyProcessing[policy.PolicyName] {
							mediation := mediation.CreateMediation(policy)
							if mediation == nil {
								s.log.Sugar().Errorf("Failed to create mediation for policy: %+v", policy)
								continue
							}
							mediationResult := mediation.Process(requestConfigHolder)
							s.log.Sugar().Debugf("Mediation Result: %+v", mediationResult)
							s.updateRequestConfigBasedOnMediationResults(mediationResult, requestConfigHolder, requestconfig.ProcessingPhaseRequestBody)
							stopProcessingMediations := s.processMediationResultAndPrepareResponse(
								mediationResult,
								resp,
								requestconfig.ProcessingPhaseRequestBody,
								metadata)
							if stopProcessingMediations {
								s.log.Sugar().Debug("Stopping further processing of request headers due to immediate response")
								break
							}
						}
					}
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_ResponseHeaders:
			s.log.Sugar().Debug("Response Headers Flow")
			s.log.Sugar().Debug(fmt.Sprintf("response header %+v, ", v.ResponseHeaders))
			requestConfigHolder.ResponseHeaders = req.GetResponseHeaders()
			requestConfigHolder.ProcessingPhase = requestconfig.ProcessingPhaseResponseHeaders

			rhq := &envoy_service_proc_v3.HeadersResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
					ClearRouteCache: true,
				},
			}
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_ResponseHeaders{
				ResponseHeaders: rhq,
			}

			if requestConfigHolder.RoutePolicy != nil {
				s.log.Sugar().Debugf("RoutePolicy: %+v", requestConfigHolder.RoutePolicy)
				if requestConfigHolder.RoutePolicy.Spec.ResponseMediation != nil {
					s.log.Sugar().Debugf("Request Mediation Policies: %+v", requestConfigHolder.RoutePolicy.Spec.RequestMediation)
					for _, policy := range requestConfigHolder.RoutePolicy.Spec.ResponseMediation {
						if mediation.MediationAndResponseHeaderProcessing[policy.PolicyName] {
							mediation := mediation.CreateMediation(policy)
							if mediation == nil {
								s.log.Sugar().Errorf("Failed to create mediation for policy: %+v", policy)
								continue
							}
							mediationResult := mediation.Process(requestConfigHolder)
							s.log.Sugar().Debugf("Mediation Result: %+v", mediationResult)
							s.updateRequestConfigBasedOnMediationResults(mediationResult, requestConfigHolder, requestconfig.ProcessingPhaseResponseHeaders)
							stopProcessingMediations := s.processMediationResultAndPrepareResponse(
								mediationResult,
								resp,
								requestconfig.ProcessingPhaseResponseHeaders,
								metadata)
							if stopProcessingMediations {
								s.log.Sugar().Debug("Stopping further processing of request headers due to immediate response")
								break
							}
						}
					}
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_ResponseBody:
			s.log.Sugar().Debug("Response Body Flow")
			s.log.Sugar().Debug(fmt.Sprintf("response body %+v, ", v.ResponseBody))
			requestConfigHolder.ResponseBody = req.GetResponseBody()
			requestConfigHolder.ProcessingPhase = requestconfig.ProcessingPhaseResponseBody

			rhq := &envoy_service_proc_v3.BodyResponse{
				Response: &envoy_service_proc_v3.CommonResponse{
					ClearRouteCache: true,
				},
			}
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_ResponseBody{
				ResponseBody: rhq,
			}

			if requestConfigHolder.RoutePolicy != nil {
				s.log.Sugar().Debugf("RoutePolicy: %+v", requestConfigHolder.RoutePolicy)
				if requestConfigHolder.RoutePolicy.Spec.ResponseMediation != nil {
					s.log.Sugar().Debugf("Request Mediation Policies: %+v", requestConfigHolder.RoutePolicy.Spec.RequestMediation)
					for _, policy := range requestConfigHolder.RoutePolicy.Spec.ResponseMediation {
						if mediation.MediationAndResponseBodyProcessing[policy.PolicyName] {
							mediation := mediation.CreateMediation(policy)
							if mediation == nil {
								s.log.Sugar().Errorf("Failed to create mediation for policy: %+v", policy)
								continue
							}
							mediationResult := mediation.Process(requestConfigHolder)
							s.log.Sugar().Debugf("Mediation Result: %+v", mediationResult)
							s.updateRequestConfigBasedOnMediationResults(mediationResult, requestConfigHolder, requestconfig.ProcessingPhaseResponseBody)
							stopProcessingMediations := s.processMediationResultAndPrepareResponse(
								mediationResult,
								resp,
								requestconfig.ProcessingPhaseResponseBody,
								metadata)
							if stopProcessingMediations {
								s.log.Sugar().Debug("Stopping further processing of request headers due to immediate response")
								break
							}
						}
					}
				}
			}
		default:
			s.log.Sugar().Debug(fmt.Sprintf("Unknown Request type %v\n", v))
		}
		// Set dynamic metadata
		resp.DynamicMetadata = &structpb.Struct{
			Fields: map[string]*structpb.Value{
				constants.MetadataNamespace: {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{Fields: metadata},
					},
				},
			},
		}
		if err := srv.Send(resp); err != nil {
			s.log.Sugar().Debug(fmt.Sprintf("send error %v", err))
		}
	}
}

func (s *ExternalProcessingServer) updateRequestConfigBasedOnMediationResults(mediationResult *mediation.Result, requestConfigHolder *requestconfig.Holder, processingPhase requestconfig.ProcessingPhase) {
	if len(mediationResult.RemoveHeaders) > 0 {
		s.log.Sugar().Debugf("Removing headers: %v", mediationResult.RemoveHeaders)
		var headerValues []*corev3.HeaderValue
		if processingPhase == requestconfig.ProcessingPhaseRequestHeaders {
			for _, header := range requestConfigHolder.RequestHeaders.Headers.Headers {
				if !util.Contains(mediationResult.RemoveHeaders, header.Key) {
					headerValues = append(headerValues, &corev3.HeaderValue{
						Key:      header.Key,
						RawValue: []byte(header.RawValue),
					})
				} else {
					s.log.Sugar().Debugf("Removing header: %s", header.Key)
				}
			}
		} else if processingPhase == requestconfig.ProcessingPhaseResponseHeaders {
			for _, header := range requestConfigHolder.ResponseHeaders.Headers.Headers {
				if !util.Contains(mediationResult.RemoveHeaders, header.Key) {
					headerValues = append(headerValues, &corev3.HeaderValue{
						Key:      header.Key,
						RawValue: []byte(header.RawValue),
					})
				} else {
					s.log.Sugar().Debugf("Removing header: %s", header.Key)
				}
			}
		}
		requestConfigHolder.RequestHeaders.Headers.Headers = headerValues
	}
	if len(mediationResult.AddHeaders) > 0 {
		s.log.Sugar().Debugf("Adding headers: %v", mediationResult.AddHeaders)
		if processingPhase == requestconfig.ProcessingPhaseRequestHeaders {
			for key, value := range mediationResult.AddHeaders {
				s.log.Sugar().Debugf("Adding header: %s: %s", key, value)
				requestConfigHolder.RequestHeaders.Headers.Headers = append(requestConfigHolder.RequestHeaders.Headers.Headers, &corev3.HeaderValue{
					Key:      key,
					RawValue: []byte(value),
				})
			}
		} else if processingPhase == requestconfig.ProcessingPhaseResponseHeaders {
			for key, value := range mediationResult.AddHeaders {
				s.log.Sugar().Debugf("Adding header: %s: %s", key, value)
				requestConfigHolder.ResponseHeaders.Headers.Headers = append(requestConfigHolder.ResponseHeaders.Headers.Headers, &corev3.HeaderValue{
					Key:      key,
					RawValue: []byte(value),
				})
			}
		}
	}
	if mediationResult.ModifyBody {
		s.log.Sugar().Debugf("Modifying body: %s", mediationResult.Body)
		if processingPhase == requestconfig.ProcessingPhaseRequestBody {
			requestConfigHolder.RequestBody.Body = []byte(mediationResult.Body)
		} else if processingPhase == requestconfig.ProcessingPhaseResponseBody {
			requestConfigHolder.ResponseBody.Body = []byte(mediationResult.Body)
		}
	}
}

func (s *ExternalProcessingServer) processMediationResultAndPrepareResponse(
	mediationResult *mediation.Result,
	resp *envoy_service_proc_v3.ProcessingResponse,
	processingPhase requestconfig.ProcessingPhase,
	metadata map[string]*structpb.Value) (stopProcessingMediations bool) {
	if mediationResult == nil {
		s.log.Sugar().Debug("Mediation Result is nil, skipping further processing")
		return false
	}
	if len(mediationResult.Metadata) > 0 {
		s.log.Sugar().Debugf("Mediation Result Metadata: %+v", mediationResult.Metadata)
		for key, value := range mediationResult.Metadata {
			metadata[key] = value
		}
	} else {
		s.log.Sugar().Debug("No Metadata found in Mediation Result")
	}
	headerMutation := &envoy_service_proc_v3.HeaderMutation{}
	if len(mediationResult.RemoveHeaders) > 0 {
		s.log.Sugar().Debugf("Removing headers: %v", mediationResult.RemoveHeaders)
		headerMutation.RemoveHeaders = mediationResult.RemoveHeaders
	}
	if len(mediationResult.AddHeaders) > 0 {
		s.log.Sugar().Debugf("Adding headers: %v", mediationResult.AddHeaders)
		for key, value := range mediationResult.AddHeaders {
			headerMutation.SetHeaders = append(headerMutation.SetHeaders, &corev3.HeaderValueOption{
				Header: &corev3.HeaderValue{
					Key:      key,
					RawValue: []byte(value),
				},
				AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
			})
		}
	}
	bodyMutation := &envoy_service_proc_v3.BodyMutation{}
	if mediationResult.ModifyBody {
		bodyMutation = &envoy_service_proc_v3.BodyMutation{
			Mutation: &envoy_service_proc_v3.BodyMutation_Body{
				Body: []byte(mediationResult.Body),
			},
		}
	}
	if mediationResult.ImmediateResponse {
		resp.Response = &envoy_service_proc_v3.ProcessingResponse_ImmediateResponse{
			ImmediateResponse: &envoy_service_proc_v3.ImmediateResponse{
				Status: &v32.HttpStatus{
					Code: mediationResult.ImmediateResponseCode,
				},
				Body:    []byte(mediationResult.ImmediateResponseBody),
				Details: mediationResult.ImmediateResponseDetail,
				Headers: headerMutation,
			},
		}
		return true
	}
	if processingPhase == requestconfig.ProcessingPhaseRequestHeaders {
		if resp.Response == nil {
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestHeaders{}
		} else {
			requestHeaderResp, ok := resp.Response.(*envoy_service_proc_v3.ProcessingResponse_RequestHeaders)
			if !ok {
				s.log.Sugar().Debugf("Unexpected response type: %T, expected RequestHeaders", resp.Response)
				return false
			}
			if requestHeaderResp.RequestHeaders == nil {
				requestHeaderResp.RequestHeaders = &envoy_service_proc_v3.HeadersResponse{}
			}
			if requestHeaderResp.RequestHeaders.Response == nil {
				requestHeaderResp.RequestHeaders.Response = &envoy_service_proc_v3.CommonResponse{}
			}
			if requestHeaderResp.RequestHeaders.Response.HeaderMutation != nil {
				requestHeaderResp.RequestHeaders.Response.HeaderMutation.SetHeaders = append(requestHeaderResp.RequestHeaders.Response.HeaderMutation.SetHeaders, headerMutation.SetHeaders...)
				requestHeaderResp.RequestHeaders.Response.HeaderMutation.RemoveHeaders = append(requestHeaderResp.RequestHeaders.Response.HeaderMutation.RemoveHeaders, headerMutation.RemoveHeaders...)
			} else {
				requestHeaderResp.RequestHeaders.Response.HeaderMutation = headerMutation
			}
		}
	} else if processingPhase == requestconfig.ProcessingPhaseRequestBody {
		if resp.Response == nil {
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_RequestBody{}
		} else {
			requestBodyResp, ok := resp.Response.(*envoy_service_proc_v3.ProcessingResponse_RequestBody)
			if !ok {
				s.log.Sugar().Debugf("Unexpected response type: %T, expected RequestBody", resp.Response)
				return false
			}
			if requestBodyResp.RequestBody == nil {
				requestBodyResp.RequestBody = &envoy_service_proc_v3.BodyResponse{}
			}
			if requestBodyResp.RequestBody.Response == nil {
				requestBodyResp.RequestBody.Response = &envoy_service_proc_v3.CommonResponse{}
			}
			if mediationResult.ModifyBody {
				requestBodyResp.RequestBody.Response.BodyMutation = bodyMutation
			}
			if len(mediationResult.AddHeaders) > 0 || len(mediationResult.RemoveHeaders) > 0 {
				if requestBodyResp.RequestBody.Response.HeaderMutation != nil {
					requestBodyResp.RequestBody.Response.HeaderMutation.SetHeaders = append(requestBodyResp.RequestBody.Response.HeaderMutation.SetHeaders, headerMutation.SetHeaders...)
					requestBodyResp.RequestBody.Response.HeaderMutation.RemoveHeaders = append(requestBodyResp.RequestBody.Response.HeaderMutation.RemoveHeaders, headerMutation.RemoveHeaders...)
				} else {
					requestBodyResp.RequestBody.Response.HeaderMutation = headerMutation
				}
			}
		}
	} else if processingPhase == requestconfig.ProcessingPhaseResponseHeaders {
		if resp.Response == nil {
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_ResponseHeaders{}
		} else {
			responseHeaderResp, ok := resp.Response.(*envoy_service_proc_v3.ProcessingResponse_ResponseHeaders)
			if !ok {
				s.log.Sugar().Debugf("Unexpected response type: %T, expected ResponseHeaders", resp.Response)
				return false
			}
			if responseHeaderResp.ResponseHeaders == nil {
				responseHeaderResp.ResponseHeaders = &envoy_service_proc_v3.HeadersResponse{}
			}
			if responseHeaderResp.ResponseHeaders.Response == nil {
				responseHeaderResp.ResponseHeaders.Response = &envoy_service_proc_v3.CommonResponse{}
			}
			if responseHeaderResp.ResponseHeaders.Response.HeaderMutation != nil {
				responseHeaderResp.ResponseHeaders.Response.HeaderMutation.SetHeaders = append(responseHeaderResp.ResponseHeaders.Response.HeaderMutation.SetHeaders, headerMutation.SetHeaders...)
				responseHeaderResp.ResponseHeaders.Response.HeaderMutation.RemoveHeaders = append(responseHeaderResp.ResponseHeaders.Response.HeaderMutation.RemoveHeaders, headerMutation.RemoveHeaders...)
			} else {
				responseHeaderResp.ResponseHeaders.Response.HeaderMutation = headerMutation
			}
		}
	} else if processingPhase == requestconfig.ProcessingPhaseResponseBody {
		if resp.Response == nil {
			resp.Response = &envoy_service_proc_v3.ProcessingResponse_ResponseBody{}
		} else {
			responseBodyResp, ok := resp.Response.(*envoy_service_proc_v3.ProcessingResponse_ResponseBody)
			if !ok {
				s.log.Sugar().Debugf("Unexpected response type: %T, expected ResponseBody", resp.Response)
				return false
			}
			if responseBodyResp.ResponseBody == nil {
				responseBodyResp.ResponseBody = &envoy_service_proc_v3.BodyResponse{}
			}
			if responseBodyResp.ResponseBody.Response == nil {
				responseBodyResp.ResponseBody.Response = &envoy_service_proc_v3.CommonResponse{}
			}
			if mediationResult.ModifyBody {
				responseBodyResp.ResponseBody.Response.BodyMutation = bodyMutation
			}
			if len(mediationResult.AddHeaders) > 0 || len(mediationResult.RemoveHeaders) > 0 {
				if responseBodyResp.ResponseBody.Response.HeaderMutation != nil {
					responseBodyResp.ResponseBody.Response.HeaderMutation.SetHeaders = append(responseBodyResp.ResponseBody.Response.HeaderMutation.SetHeaders, headerMutation.SetHeaders...)
					responseBodyResp.ResponseBody.Response.HeaderMutation.RemoveHeaders = append(responseBodyResp.ResponseBody.Response.HeaderMutation.RemoveHeaders, headerMutation.RemoveHeaders...)
				} else {
					responseBodyResp.ResponseBody.Response.HeaderMutation = headerMutation
				}
			}
		}
	} else {
		s.log.Sugar().Debugf("Unknown processing phase: %s", processingPhase)
		// Return false to indicate that processing should continue
	}
	return false
}

// func getFileNameAndContentTypeForDef(matchedAPI *requestconfig.API) (string, string) {
// 	fileName := "attachment; filename=\"api_definition.json\""
// 	contentType := "application/octet-stream"
// 	if matchedAPI.IsGraphQLAPI() {
// 		fileName = "attachment; filename=\"api_definition.graphql\""
// 	}
// 	if matchedAPI.IsgRPCAPI() {
// 		fileType, _ := DetectFileType([]byte(matchedAPI.APIDefinition))

// 		if fileType == "proto" {
// 			return "attachment; filename=\"api_definition.proto\"", contentType
// 		}
// 		if fileType == "zip" {
// 			return "attachment; filename=\"api_definition.zip\"", "application/zip"
// 		}
// 	}

// 	return fileName, contentType
// }

// DetectFileType detects if the file is a .proto or .zip
func DetectFileType(data []byte) (string, error) {

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	if bytes.Contains(decompressedData, []byte("syntax = ")) {
		return "proto", nil
	}

	return "zip", nil
}

// ReadGzip decompresses a GZIP-compressed byte slice and returns the string output
func ReadGzip(gzipData []byte) (string, error) {
	// Create a bytes.Reader from the gzip data
	byteReader := bytes.NewReader(gzipData)

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(byteReader)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// Read the uncompressed data
	var result bytes.Buffer
	_, err = io.Copy(&result, gzipReader)
	if err != nil {
		return "", err
	}

	// Convert bytes buffer to string
	return result.String(), nil
}

func (s *ExternalProcessingServer) extractJWTAuthnNamespaceData(data *corev3.Metadata) map[string]interface{} {
	claims := make(map[string]interface{})
	filterMatadata := data.GetFilterMetadata()
	if filterMatadata != nil {
		jwtAuthnData, exists := filterMatadata[constants.JWTAuthnMetadataNamespace]
		if !exists || jwtAuthnData == nil {
			s.cfg.Logger.Sugar().Debug("JWT Authn data not found")
		} else {
			for _, structValue := range jwtAuthnData.Fields {
				s.cfg.Logger.Sugar().Debugf("JWT Authn Payload: %s", structValue.GetStringValue())
				jwtPayload := structValue.GetStructValue()
				if jwtPayload != nil {
					for key, value := range jwtPayload.GetFields() {
						if value != nil {
							switch value.Kind.(type) {
							case *structpb.Value_StringValue:
								claims[key] = value.GetStringValue()
							case *structpb.Value_NumberValue:
								claims[key] = value.GetNumberValue()
							case *structpb.Value_BoolValue:
								claims[key] = value.GetBoolValue()
							case *structpb.Value_ListValue:
								jsonData, err := value.MarshalJSON()
								if err == nil {
									var list []interface{}
									err = json.Unmarshal(jsonData, &list)
									if err == nil {
										claims[key] = list
									} else {
										s.cfg.Logger.Sugar().Errorf("Failed to unmarshal list value for key %s: %v", key, err)
									}
								} else {
									s.cfg.Logger.Sugar().Errorf("Failed to marshal JSON for list value for key %s: %v", key, err)
								}
							case *structpb.Value_StructValue:
								jsonData, err := value.MarshalJSON()
								if err == nil {
									var mapData map[string]interface{}
									err = json.Unmarshal(jsonData, &mapData)
									if err == nil {
										claims[key] = mapData
									} else {
										s.cfg.Logger.Sugar().Errorf("Failed to unmarshal struct value for key %s: %v", key, err)
									}
								} else {
									s.cfg.Logger.Sugar().Errorf("Failed to marshal JSON for struct value for key %s: %v", key, err)
								}
							}
						}
					}
				}
			}
		}
	}
	s.cfg.Logger.Sugar().Debugf("Extracted JWT Authn Claims: %+v", claims)
	return claims
}

func (s *ExternalProcessingServer) extractExtensionRefsAndRouteAttributes(data map[string]*structpb.Struct) ([]string, *requestconfig.Attributes, string) {
	var extensionRefs []string
	envType := ""
	extProcData, exists := data[constants.ExternalProcessingNamespace]
	if !exists || extProcData == nil {
		s.cfg.Logger.Sugar().Debug("External processing data not found in attributes, Returning empty extensionRefs")
		s.cfg.Logger.Sugar().Debugf("Attributes: %+v", data)
		return extensionRefs, &requestconfig.Attributes{}, envType
	}

	// Check if `xds.route_metadata` is a stringified proto and extract the nested `filter_metadata`
	if rawTextStruct, ok := extProcData.Fields["xds.route_metadata"]; ok {
		rawText := rawTextStruct.GetStringValue()
		if rawText != "" {
			s.cfg.Logger.Sugar().Debugf("Raw xds.route_metadata text: %s", rawText)
			var structFromText corev3.Metadata
			if err := prototext.Unmarshal([]byte(rawText), &structFromText); err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed to unmarshal text proto from xds.route_metadata: %v", err)
			} else {
				filterMetadata := structFromText.GetFilterMetadata()
				// Try to extract ExtensionRefs from parsed struct if available
				if extProcFilter, ok := filterMetadata[constants.ExternalProcessingNamespace]; ok {
					if extProcFilter != nil {
						if field, ok := extProcFilter.Fields[constants.ExtensionRefs]; ok {
							listVal := field.GetListValue()
							for _, val := range listVal.GetValues() {
								if str := val.GetStringValue(); str != "" {
									extensionRefs = append(extensionRefs, str)
								}
							}
						}
					}
				}

				// ----- Extract Annotations from envoy-gateway -----
				if envoyGatewayFilter, ok := filterMetadata["envoy-gateway"]; ok && envoyGatewayFilter != nil {
					if resourcesField, ok := envoyGatewayFilter.Fields["resources"]; ok {
						for _, resVal := range resourcesField.GetListValue().GetValues() {
							if resStruct := resVal.GetStructValue(); resStruct != nil {
								if annotationsField, ok := resStruct.Fields["annotations"]; ok {
									if annotationsStruct := annotationsField.GetStructValue(); annotationsStruct != nil {
										for k, v := range annotationsStruct.Fields {
											if strVal := v.GetStringValue(); strVal != "" {
												if k == "kgw-envtype" {
													envType = strVal
													break
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	routeName := ""
	if stringVal, ok := extProcData.Fields["xds.route_name"]; ok {
		routeName = stringVal.GetStringValue()
	}
	requestID := ""
	if stringVal, ok := extProcData.Fields["request.id"]; ok {
		requestID = stringVal.GetStringValue()
	}

	return extensionRefs, &requestconfig.Attributes{
		RouteName: routeName,
		RequestID: requestID,
	}, envType
}

func getNestedStruct(base *structpb.Value, key string) *structpb.Struct {
	if base == nil || base.GetStructValue() == nil {
		return nil
	}
	val, ok := base.GetStructValue().Fields[key]
	if !ok || val == nil {
		return nil
	}
	return val.GetStructValue()
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

func prepareMetadataKeyValuePairAndAddTo(metadataKeyValuePair map[string]string, requestConfigHolder *requestconfig.Holder, cfg *config.Server) *map[string]string {

	return &metadataKeyValuePair
}
