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
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	// "github.com/wso2/apk/gateway/enforcer/internal/authentication/authenticator"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/mediation"

	// "github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
	// "github.com/wso2/apk/gateway/enforcer/internal/ratelimit"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	// "github.com/wso2/apk/gateway/enforcer/internal/requesthandler"
	// "github.com/wso2/apk/gateway/enforcer/internal/transformer"
	v31 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	"github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
	types "k8s.io/apimachinery/pkg/types"

	// "net"
	// "time"

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
			extRefs := s.extractExtensionRefs(req.Attributes)
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
									requestConfigHolder.RoutePolicy = routePolicy
								} else {
									for _, reqPolicy := range routePolicy.Spec.RequestMediation {
										requestConfigHolder.RoutePolicy.Spec.RequestMediation = append(requestConfigHolder.RoutePolicy.Spec.RequestMediation, reqPolicy)
									}
									for _, resPolicy := range routePolicy.Spec.ResponseMediation {
										requestConfigHolder.RoutePolicy.Spec.ResponseMediation = append(requestConfigHolder.RoutePolicy.Spec.ResponseMediation, resPolicy)
									}
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

		dynamicMetadataKeyValuePairs := make(map[string]string)
		switch v := req.Request.(type) {
		case *envoy_service_proc_v3.ProcessingRequest_RequestHeaders:
			s.log.Sugar().Debug("Request Headers Flow")
			s.log.Sugar().Debug(fmt.Sprintf("request header %+v, ", v.RequestHeaders))
			requestConfigHolder.RequestHeaders = req.GetRequestHeaders()
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
						if mediation.MediationAndRequestBodyProcessing[policy.PolicyName] {
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
						if mediation.MediationAndResponseHeaderProcessing[policy.PolicyName] {
							responseHeaderMode = v31.ProcessingMode_SEND
						}
						if mediation.MediationAndResponseBodyProcessing[policy.PolicyName] {
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
							
						}
					}
				}
			}

		case *envoy_service_proc_v3.ProcessingRequest_RequestBody:
			s.log.Sugar().Debug("Request Body Flow")
			s.log.Sugar().Debug(fmt.Sprintf("request body %+v, ", v.RequestBody))
			requestConfigHolder.RequestBody = req.GetRequestBody()

			

		case *envoy_service_proc_v3.ProcessingRequest_ResponseHeaders:
			s.log.Sugar().Debug("Response Headers Flow")
			s.log.Sugar().Debug(fmt.Sprintf("response header %+v, ", v.ResponseHeaders))
			requestConfigHolder.ResponseHeaders = req.GetRequestHeaders()

		case *envoy_service_proc_v3.ProcessingRequest_ResponseBody:
			s.log.Sugar().Debug("Response Body Flow")
			s.log.Sugar().Debug(fmt.Sprintf("response body %+v, ", v.ResponseBody))
			requestConfigHolder.ResponseBody = req.GetResponseBody()

		default:
			s.log.Sugar().Debug(fmt.Sprintf("Unknown Request type %v\n", v))
		}
		// Set dynamic metadata
		dynamicMetadata, err := buildDynamicMetadata(prepareMetadataKeyValuePairAndAddTo(dynamicMetadataKeyValuePairs, requestConfigHolder, s.cfg))
		if err != nil {
			s.log.Error(err, "failed to build dynamic metadata")
		} else {
			resp.DynamicMetadata = dynamicMetadata
		}
		if err := srv.Send(resp); err != nil {
			s.log.Sugar().Debug(fmt.Sprintf("send error %v", err))
		}
	}
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

// extractExternalProcessingMetadata extracts the external processing metadata from the given data.
// func extractExternalProcessingMetadata(data *corev3.Metadata) (*dto.ExternalProcessingEnvoyMetadata, error) {
// 	filterMatadata := data.GetFilterMetadata()
// 	if filterMatadata != nil {
// 		externalProcessingEnvoyMetadata := &dto.ExternalProcessingEnvoyMetadata{}
// 		jwtFilterdata := filterMatadata["envoy.filters.http.jwt_authn"]
// 		if jwtFilterdata != nil {
// 			authenticationData := &dto.AuthenticationData{}

// 			for key, structValue := range jwtFilterdata.Fields {
// 				if strings.HasSuffix(key, "-payload") {
// 					sucessData := dto.AuthenticationSuccessData{}
// 					jwtPayload := structValue.GetStructValue()
// 					if jwtPayload != nil {
// 						claims := make(map[string]interface{})
// 						for key, value := range jwtPayload.GetFields() {
// 							if value != nil {
// 								if key == "iss" {
// 									sucessData.Issuer = value.GetStringValue()
// 								}
// 								switch value.Kind.(type) {
// 								case *structpb.Value_StringValue:
// 									claims[key] = value.GetStringValue()
// 								case *structpb.Value_NumberValue:
// 									claims[key] = value.GetNumberValue()
// 								case *structpb.Value_BoolValue:
// 									claims[key] = value.GetBoolValue()
// 								case *structpb.Value_ListValue:
// 									jsonData, err := value.MarshalJSON()
// 									if err != nil {
// 										return nil, err
// 									}
// 									var list []interface{}
// 									err = json.Unmarshal(jsonData, &list)
// 									if err != nil {
// 										return nil, err
// 									}
// 									claims[key] = list
// 								case *structpb.Value_StructValue:
// 									jsonData, err := value.MarshalJSON()
// 									if err != nil {
// 										return nil, err
// 									}
// 									var mapData map[string]interface{}
// 									err = json.Unmarshal(jsonData, &mapData)
// 									if err != nil {
// 										return nil, err
// 									}
// 									claims[key] = mapData
// 								}
// 							}
// 						}
// 						sucessData.Claims = claims
// 					}
// 					if authenticationData.SucessData == nil {
// 						authenticationData.SucessData = make(map[string]*dto.AuthenticationSuccessData)
// 					}
// 					authenticationData.SucessData[key] = &sucessData
// 				}
// 				if strings.HasSuffix(key, "-failed") {
// 					failureStatusStruct := structValue.GetStructValue()
// 					if failureStatusStruct != nil {
// 						code := failureStatusStruct.Fields["code"].GetNumberValue()
// 						message := failureStatusStruct.Fields["message"].GetStringValue()
// 						authenticationFailureData := &dto.AuthenticationFailureData{Code: int(code), Message: message}
// 						if authenticationData.FailedData == nil {
// 							authenticationData.FailedData = make(map[string]*dto.AuthenticationFailureData)
// 						}
// 						authenticationData.FailedData[key] = authenticationFailureData
// 					}
// 				}
// 			}
// 			externalProcessingEnvoyMetadata.AuthenticationData = authenticationData
// 		}
// 		if extProcMetadata, exists := filterMatadata[externalProessingMetadataContextKey]; exists {
// 			if matchedAPIKey, exists := extProcMetadata.Fields[matchedAPIMetadataKey]; exists {
// 				externalProcessingEnvoyMetadata.MatchedAPIIdentifier = matchedAPIKey.GetStringValue()
// 			}
// 			if matchedResourceKey, exists := extProcMetadata.Fields[matchedResourceMetadataKey]; exists {
// 				externalProcessingEnvoyMetadata.MatchedResourceIdentifier = matchedResourceKey.GetStringValue()
// 			}
// 			if matchedApplicationKey, exists := extProcMetadata.Fields[matchedApplicationMetadataKey]; exists {
// 				externalProcessingEnvoyMetadata.MatchedApplicationIdentifier = matchedApplicationKey.GetStringValue()
// 			}
// 			if matchedSubscriptionKey, exists := extProcMetadata.Fields[matchedSubscriptionMetadataKey]; exists {
// 				externalProcessingEnvoyMetadata.MatchedSubscriptionIdentifier = matchedSubscriptionKey.GetStringValue()
// 			}

// 		}
// 		return externalProcessingEnvoyMetadata, nil
// 	}
// 	return nil, nil
// }

func readStructData() {

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

func (s *ExternalProcessingServer) extractExtensionRefs(data map[string]*structpb.Struct) []string {
	var extensionRefs []string

	extProcData, exists := data[constants.ExternalProcessingNamespace]
	if !exists || extProcData == nil {
		s.cfg.Logger.Sugar().Debug("External processing data not found in attributes, Returning empty extensionRefs")
		s.cfg.Logger.Sugar().Debugf("Attributes: %+v", data)
		return extensionRefs
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
			}
		}
	}

	return extensionRefs
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
