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

package requestconfig

import (
	envoy_service_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	subscription_model "github.com/wso2/apk/common-go-libs/pkg/server/model"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// Holder is a struct that holds the request configuration.
type Holder struct {
	AttributesPopulated             bool
	RouteMetadata                   *dpv2alpha1.RouteMetadata
	RoutePolicy                     *dpv2alpha1.RoutePolicy
	APIKeyAuthenticationInfo        *dto.APIKeyAuthenticationInfo
	MatchedSubscription             *subscription_model.Subscription
	MatchedApplication              *subscription_model.Application
	AuthenticatedAuthenticationType string
	RequestHeaders                  *envoy_service_proc_v3.HttpHeaders
	ResponseHeaders                 *envoy_service_proc_v3.HttpHeaders
	ResponseBody                    *envoy_service_proc_v3.HttpBody
	RequestBody                     *envoy_service_proc_v3.HttpBody
	ProcessingPhase                 ProcessingPhase
	AI                              AIConfig
	RequestAttributes               *Attributes
	JWTAuthnPayloaClaims            map[string]interface{}
}

// Attributes holds the attributes related to the request configuration.
type Attributes struct {
	// RouteName is the name of the route.
	RouteName string
	// RequestID is the ID of the request.
	RequestID string
}

// ProcessingPhase represents the phase of processing in the request configuration.
type ProcessingPhase string

const (
	// ProcessingPhaseRequestHeaders represents the request headers processing phase.
	ProcessingPhaseRequestHeaders ProcessingPhase = "request_headers"
	// ProcessingPhaseResponseHeaders represents the response headers processing phase.
	ProcessingPhaseResponseHeaders ProcessingPhase = "response_headers"
	// ProcessingPhaseRequestBody represents the request body processing phase.
	ProcessingPhaseRequestBody ProcessingPhase = "request_body"
	// ProcessingPhaseResponseBody represents the response body processing phase.
	ProcessingPhaseResponseBody ProcessingPhase = "response_body"
)

// AIConfig holds the configuration for AI model suspension.
type AIConfig struct {
	// SuspendModel indicates whether the AI model should be suspended.
	SuspendModel bool
}
