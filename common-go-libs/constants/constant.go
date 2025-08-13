/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package constants

// Controller related constants.
const (
	RatelimitController          string = "RatelimitController"
	AIRatelimitController        string = "AIRatelimitController"
	ApplicationController        string = "ApplicationController"
	SubscriptionController       string = "SubscriptionController"
	ApplicationMappingController string = "ApplicationMappingController"
	GatewayClassController       string = "GatewayClassController"
	RoutePolicyController        string = "RoutePolicyController"
	RouteMetadataController      string = "RouteMetadataController"
)

// API events related constants
const (
	Create string = "CREATED"
	Update string = "UPDATED"
	Delete string = "DELETED"
)

// Subscriprion events related constants
const (
	ApplicationCreated            string = "APPLICATION_CREATED"
	ApplicationUpdated            string = "APPLICATION_UPDATED"
	ApplicationDeleted            string = "APPLICATION_DELETED"
	SubscriptionCreated           string = "SUBSCRIPTION_CREATED"
	SubscriptionUpdated           string = "SUBSCRIPTION_UPDATED"
	SubscriptionDeleted           string = "SUBSCRIPTION_DELETED"
	ApplicationMappingCreated     string = "APPLICATION_MAPPING_CREATED"
	ApplicationMappingUpdated     string = "APPLICATION_MAPPING_UPDATED"
	ApplicationMappingDeleted     string = "APPLICATION_MAPPING_DELETED"
	ApplicationKeyMappingCreated  string = "APPLICATION_KEY_MAPPING_CREATED"
	ApplicationKeyMappingUpdated  string = "APPLICATION_KEY_MAPPING_UPDATED"
	ApplicationKeyMappingDeleted  string = "APPLICATION_KEY_MAPPING_DELETED"
	RoutePolicyCreatedOrUpdated   string = "ROUTE_POLICY_CREATED_OR_UPDATED"
	RoutePolicyDeleted            string = "ROUTE_POLICY_DELETED"
	RouteMetadataCreatedOrUpdated string = "ROUTE_METADATA_CREATED_OR_UPDATED"
	RouteMetadataDeleted          string = "ROUTE_METADATA_DELETED"
	AllEvents                     string = "ALL_EVENTS"
)

// Environment variable names and default values
const (
	OperatorPodNamespace             string = "OPERATOR_POD_NAMESPACE"
	OperatorPodNamespaceDefaultValue string = "default"
)

// CRD Kinds
const (
	KindAuthentication = "Authentication"
	KindAPI            = "API"
	All                = "All"
	//TODO(amali) remove this after fixing the issue in https://github.com/wso2/apk/issues/383
	KindResource        = "Resource"
	KindRateLimitPolicy = "RateLimitPolicy"
)

// Env types
const (
	Production = "PRODUCTION"
	Sandbox    = "SANDBOX"
)

// Header names in runtime
const (
	OrganizationHeader = "X-WSO2-Organization"
)

// Global interceptor cluster names
const (
	GlobalRequestInterceptorClusterName  = "request_interceptor_global_cluster"
	GlobalResponseInterceptorClusterName = "response_interceptor_global_cluster"
)

// Application authentication types
const (
	OAuth2 = "OAuth2"
)

// XDSRoute Metadata
const (
	ExternalProcessingNamespace = "envoy.filters.http.ext_proc"
	JWTAuthnMetadataNamespace = "envoy.filters.http.jwt_authn"
	ExtensionRefs               = "ExtensionRefs"
)

// Metadata keys for AI Token Rate Limit mediation policy
const (
	PromptTokenCountIDMetadataKey     = "promptTokenCount"
	CompletionTokenCountIDMetadataKey = "completionTokenCount"
	TotalTokenCountIDMetadataKey      = "totalTokenCount"
)

// Metadata from external processing
const (
	MetadataNamespace = "com.wso2.kgw.ext_proc"
	// JWTAuthnPayloadInMetadata is the key used to store JWT authentication payload in metadata
	JWTAuthnPayloadInMetadata = "jwt_authn_payload"
)

// SecurityPolicy header claims to header
const (
	ClientIDHeaderKey = "X-WSO2-Clinet-ID"
	ScopesHeaderKey = "X-WSO2-Scopes"
)

// Subscription ratelimit header names
const (
	SubscriptionUUIDHeaderName = "X-WSO2-Subscription-UUID"
)

// Kind constants for various Kubernetes resources
const (
	KindConfigMap            = "ConfigMap"
	KindClientTrafficPolicy  = "ClientTrafficPolicy"
	KindBackendTrafficPolicy = "BackendTrafficPolicy"
	KindBackendTLSPolicy     = "BackendTLSPolicy"
	KindBackend              = "Backend"
	KindEnvoyPatchPolicy     = "EnvoyPatchPolicy"
	KindEnvoyExtensionPolicy = "EnvoyExtensionPolicy"
	KindSecurityPolicy       = "SecurityPolicy"
	KindEnvoyProxy           = "EnvoyProxy"
	KindGateway              = "Gateway"
	KindGatewayClass         = "GatewayClass"
	KindGRPCRoute            = "GRPCRoute"
	KindHTTPRoute            = "HTTPRoute"
	KindNamespace            = "Namespace"
	KindTLSRoute             = "TLSRoute"
	KindTCPRoute             = "TCPRoute"
	KindUDPRoute             = "UDPRoute"
	KindService              = "Service"
	KindServiceImport        = "ServiceImport"
	KindSecret               = "Secret"
	KindHTTPRouteFilter      = "HTTPRouteFilter"
	KindReferenceGrant       = "ReferenceGrant"
)

const (
	// MediationAITokenRatelimit holds the name of the AI Token Rate Limit mediation policy.
	MediationAITokenRatelimit = "AITokenRatelimit"
	// MediationSubscriptionRatelimit holds the name of the Subscription Rate Limit mediation policy.
	MediationSubscriptionRatelimit = "SubscriptionRatelimit"
	// MediationSubscriptionValidation holds the name of the Subscription Validation mediation policy.
	MediationSubscriptionValidation = "SubscriptionValidation"
	// MediationAIModelBasedRoundRobin holds the name of the AI Model Based Round Robin mediation policy.
	MediationAIModelBasedRoundRobin = "AIModelBasedRoundRobin"
	// MediationAnalytics holds the name of the Analytics mediation policy.
	MediationAnalytics = "Analytics"
	// MediationBackendJWT holds the name of the Backend JWT mediation policy.
	MediationBackendJWT = "BackendJWT"
	// MediationGraphQL holds the name of the GraphQL mediation policy.
	MediationGraphQL = "GraphQL"
)

const (
	// GraphQLPolicyKeySchema is the key for specifying the GraphQL schema.
	GraphQLPolicyKeySchema = "Schema"
)

const (
	LabelAPKName = "apk.wso2.com/name"
	LabelAPKVersion = "apk.wso2.com/version"
	LabelAPKOrganization = "apk.wso2.com/organization"
)