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

package analytics

const (
	// UpstreamSuccessResponseDetail is used to indicate a successful response from the upstream server.
	UpstreamSuccessResponseDetail = "via_upstream"
	// ExtAuthDeniedResponseDetail indicates that the external authorization request was denied.
	ExtAuthDeniedResponseDetail = "ext_authz_denied"
	// ExtAuthErrorResponseDetail indicates an error occurred during external authorization.
	ExtAuthErrorResponseDetail = "ext_authz_error"
	// RouteNotFoundResponseDetail indicates that no route was found for the request.
	RouteNotFoundResponseDetail = "route_not_found"
	// GatewayLabel represents the label used to identify the Envoy Gateway.
	GatewayLabel = "ENVOY"

	// TokenEndpointPath is the path for the token endpoint.
	TokenEndpointPath = "/testkey"
	// HealthEndpointPath is the path for the health check endpoint.
	HealthEndpointPath = "/health"
	// JwksEndpointPath is the path for the JWKS (JSON Web Key Set) endpoint.
	JwksEndpointPath = "/.wellknown/jwks"
	// DefaultForUnassigned is the default value used for unassigned properties.
	DefaultForUnassigned = "UNKNOWN"
	// DataProviderClassProperty specifies the property name for the custom data provider class.
	DataProviderClassProperty = "publisher.custom.data.provider.class"

	// APIThrottleOutErrorCode is the error code for API-level throttling.
	APIThrottleOutErrorCode = 900800
	// HardLimitExceededErrorCode is the error code for exceeding a hard limit.
	HardLimitExceededErrorCode = 900801
	// ResourceThrottleOutErrorCode is the error code for resource-level throttling.
	ResourceThrottleOutErrorCode = 900802
	// ApplicationThrottleOutErrorCode is the error code for application-level throttling.
	ApplicationThrottleOutErrorCode = 900803
	// SubscriptionThrottleOutErrorCode is the error code for subscription-level throttling.
	SubscriptionThrottleOutErrorCode = 900804
	// BlockedErrorCode is the error code for blocked requests.
	BlockedErrorCode = 900805
	// CustomPolicyThrottleOutErrorCode is the error code for custom policy throttling.
	CustomPolicyThrottleOutErrorCode = 900806

	// NhttpReceiverInputOutputErrorSending indicates an error while sending data via the NHTTP receiver.
	NhttpReceiverInputOutputErrorSending = 101000
	// NhttpReceiverInputOutputErrorReceiving indicates an error while receiving data via the NHTTP receiver.
	NhttpReceiverInputOutputErrorReceiving = 101001
	// NhttpSenderInputOutputErrorSending indicates an error while sending data via the NHTTP sender.
	NhttpSenderInputOutputErrorSending = 101500
	// NhttpConnectionFailed indicates that the NHTTP connection failed.
	NhttpConnectionFailed = 101503
	// NhttpConnectionTimeout indicates that the NHTTP connection timed out.
	NhttpConnectionTimeout = 101504
	// NhttpConnectionClosed indicates that the NHTTP connection was closed.
	NhttpConnectionClosed = 101505
	// NhttpProtocolViolation indicates a protocol violation in the NHTTP connection.
	NhttpProtocolViolation = 101506
	// NhttpConnectTimeout indicates a timeout occurred while attempting to connect via NHTTP.
	NhttpConnectTimeout = 101508

	// WebsocketHandshakeResourcePrefix is the prefix used for WebSocket handshake resources.
	WebsocketHandshakeResourcePrefix = "init-request:"
	// GatewayURL represents the original Gateway URL header key.
	GatewayURL = "x-original-gw-url"
	// XForwardProtoHeader represents the header for the forwarded protocol.
	XForwardProtoHeader = "x-forwarded-proto"
	// XForwardPortHeader represents the header for the forwarded port.
	XForwardPortHeader = "x-forwarded-port"
)

const (
	// ExtAuthMetadataContextKey is the context key for external authorization metadata.
	ExtAuthMetadataContextKey = "envoy.filters.http.ext_authz"
	// ExtProcMetadataContextKey is the context key for external processing metadata.
	ExtProcMetadataContextKey = "envoy.filters.http.ext_proc"
	// Wso2MetadataPrefix is the prefix for WSO2 metadata.
	Wso2MetadataPrefix = "x-wso2-"
	// APIIDKey is the key for the API ID.
	APIIDKey = Wso2MetadataPrefix + "api-id"
	// APICreatorKey is the key for the API creator.
	APICreatorKey = Wso2MetadataPrefix + "api-creator"
	// APINameKey is the key for the API name.
	APINameKey = Wso2MetadataPrefix + "api-name"
	// APIVersionKey is the key for the API version.
	APIVersionKey = Wso2MetadataPrefix + "api-version"
	// APITypeKey is the key for the API type.
	APITypeKey = Wso2MetadataPrefix + "api-type"
	// APIUserNameKey is the key for the API user name.
	APIUserNameKey = Wso2MetadataPrefix + "username"
	// APIContextKey is the key for the API context.
	APIContextKey = Wso2MetadataPrefix + "api-context"
	// IsMockAPI is the key indicating if the API is a mock API.
	IsMockAPI = Wso2MetadataPrefix + "is-mock-api"
	// APICreatorTenantDomainKey is the key for the API creator tenant domain.
	APICreatorTenantDomainKey = Wso2MetadataPrefix + "api-creator-tenant-domain"
	// APIOrganizationIDKey is the key for the API organization ID.
	APIOrganizationIDKey = Wso2MetadataPrefix + "api-organization-id"

	// AppIDKey is the key for the application ID.
	AppIDKey = Wso2MetadataPrefix + "application-id"
	// AppUUIDKey is the key for the application UUID.
	AppUUIDKey = Wso2MetadataPrefix + "application-uuid"
	// AppKeyTypeKey is the key for the application key type.
	AppKeyTypeKey = Wso2MetadataPrefix + "application-key-type"
	// AppNameKey is the key for the application name.
	AppNameKey = Wso2MetadataPrefix + "application-name"
	// AppOwnerKey is the key for the application owner.
	AppOwnerKey = Wso2MetadataPrefix + "application-owner"

	// CorrelationIDKey is the key for the correlation ID.
	CorrelationIDKey = Wso2MetadataPrefix + "correlation-id"
	// RegionKey is the key for the region.
	RegionKey = Wso2MetadataPrefix + "region"

	// APIResourceTemplateKey is the key for the API resource template.
	APIResourceTemplateKey = Wso2MetadataPrefix + "api-resource-template"

	// DestinationKey is the key for the destination.
	DestinationKey         = Wso2MetadataPrefix + "destination"
	// DefaultForUnknown is the default value used for unassigned properties.
	DefaultForUnknown = "UNKNOWN"

	// UserAgentKey is the key for the user agent.
	UserAgentKey = Wso2MetadataPrefix + "user-agent"
	// ClientIPKey is the key for the client IP.
	ClientIPKey = Wso2MetadataPrefix + "client-ip"

	// ErrorCodeKey is the key for the error code.
	ErrorCodeKey = "ErrorCode"
	// ApkEnforcerReply is the key for the APK enforcer reply.
	ApkEnforcerReply = "apk-enforcer-reply"
	// RatelimitWso2OrgPrefix is the prefix for WSO2 organization rate limit.
	RatelimitWso2OrgPrefix = "customorg"
	// APIEnvironmentKey is the key for the API environment.
	APIEnvironmentKey = Wso2MetadataPrefix + "api-environment"
	// OrganizationAndAirlPolicy is the key for the organization and rate limit policy.
	OrganizationAndAirlPolicy = "ratelimit:organization-and-rlpolicy"
	// Subscription is the key for the subscription.
	Subscription = "ratelimit:subscription"
	// ExtractTokenFrom is the key for extracting the token from.
	ExtractTokenFrom = "aitoken:extracttokenfrom"
	// PromptTokenID is the key for the prompt token ID.
	PromptTokenID = "aitoken:prompttokenid"
	// CompletionTokenID is the key for the completion token ID.
	CompletionTokenID = "aitoken:completiontokenid"
	// TotalTokenID is the key for the total token ID.
	TotalTokenID = "aitoken:totaltokenid"
	// PromptTokenCount is the key for the prompt token count.
	PromptTokenCount = "aitoken:prompttokencount"
	// CompletionTokenCount is the key for the completion token count.
	CompletionTokenCount = "aitoken:completiontokencount"
	// TotalTokenCount is the key for the total token count.
	TotalTokenCount = "aitoken:totaltokencount"
	// ModelID is the key for the model ID.
	ModelID = "aitoken:modelid"
	// Model is the key for the model.
	Model = "aitoken:model"
	// AiProviderName is the key for the AI provider name.
	AiProviderName = "ai:providername"
	// AiProviderAPIVersion is the key for the AI provider API version.
	AiProviderAPIVersion = "ai:providerversion"
	//anonymousValye is the value for anonymous
	anonymousValye = "anonymous"
	// Unknown is the default value used for unassigned properties.
	Unknown = "UNKNOWN"
)
