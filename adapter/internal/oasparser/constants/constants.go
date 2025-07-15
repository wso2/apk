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

// endpoint related constants
const (
	Urls                  string = "urls"
	Type                  string = "type"
	HTTP                  string = "http"
	HTTPS                 string = "https"
	LoadBalance           string = "load_balance"
	FailOver              string = "failover"
	AdvanceEndpointConfig string = "advanceEndpointConfig"
	SecurityConfig        string = "securityConfig"
)

// Constants for OpenAPI vendor extension keys and values
const (
	XWso2BasePath                     string = "x-wso2-basePath"
	XWso2HTTP2BackendEnabled          string = "x-wso2-http2-backend-enabled"
	XThrottlingTier                   string = "x-throttling-tier"
	XWso2ThrottlingTier               string = "x-wso2-throttling-tier"
	XAuthHeader                       string = "x-wso2-auth-header"
	XAuthType                         string = "x-auth-type"
	XWso2DisableSecurity              string = "x-wso2-disable-security"
	None                              string = "None"
	DefaultSecurity                   string = "default"
	XMediationScript                  string = "x-mediation-script"
	XScopes                           string = "x-scopes"
	XWso2PassRequestPayloadToEnforcer string = "x-wso2-pass-request-payload-to-enforcer"
	XUriMapping                       string = "x-uri-mapping"
)

// sub-property values and keys relevant for x-wso2-application security extension
const (
	AuthorizationHeader  string = "authorization"
	TestConsoleKeyHeader string = "internal-key"
)

// sub-property keys mentioned under x-wso2-request-interceptor and x-wso2-response-interceptor
const (
	XWso2RequestInterceptor   string = "x-wso2-request-interceptor"
	XWso2ResponseInterceptor  string = "x-wso2-response-interceptor"
	ServiceURL                string = "serviceURL"
	ClusterTimeout            string = "clusterTimeout"
	RequestTimeout            string = "requestTimeout"
	Includes                  string = "includes"
	OperationLevelInterceptor string = "operation"
)

// Constants to represent errors
const (
	AlreadyExists string = "ALREADY_EXISTS"
	NotFound      string = "NOT_FOUND"
)

// operational policy field names
const (
	ActionHeaderAdd          string = "SET_HEADER"
	ActionHeaderRemove       string = "REMOVE_HEADER"
	ActionRewriteMethod      string = "REWRITE_RESOURCE_METHOD"
	ActionInterceptorService string = "CALL_INTERCEPTOR_SERVICE"
	ActionRewritePath        string = "REWRITE_RESOURCE_PATH"
	ActionRedirectRequest    string = "REDIRECT_REQUEST"
	ActionMirrorRequest      string = "MIRROR_REQUEST"

	PolicyRequestInterceptor  string = "PolicyRequestInterceptor"
	PolicyResponseInterceptor string = "PolicyResponseInterceptor"

	RewritePathResourcePath    string = "resourcePath"
	RewritePathType            string = "rewritePathType"
	InterceptorServiceURL      string = "interceptorServiceURL"
	InterceptorEndpoints       string = "interceptorEndpoints"
	InterceptorServiceIncludes string = "includes"
	IncludeQueryParams         string = "includeQueryParams"
	HeaderName                 string = "headerName"
	HeaderValue                string = "headerValue"
	CurrentMethod              string = "currentMethod"
	UpdatedMethod              string = "updatedMethod"
	RedirectScheme             string = "scheme"
	RedirectHostname           string = "hostname"
	RedirectPath               string = "path"
	RedirectPort               string = "port"
	RedirectStatusCode         string = "statusCode"
)

// API Type Constants
const (
	REST                  string = "REST"
	SOAP                  string = "SOAP"
	WS                    string = "WS"
	GRAPHQL               string = "GraphQL"
	GRPC                  string = "GRPC"
	WEBHOOK               string = "WEBHOOK"
	SSE                   string = "SSE"
	Prototyped            string = "prototyped"
	MockedOASEndpointType string = "MOCKED_OAS"
	TemplateEndpointType  string = "TEMPLATE"
	InlineEndpointType    string = "INLINE"
)

// Constants used for version identification of API definitions
const (
	Swagger      string = "swagger"
	OpenAPI      string = "openapi"
	AsyncAPI     string = "asyncapi"
	Swagger2     string = "swagger_2"
	OpenAPI3     string = "openapi_3"
	AsyncAPI2    string = "asyncapi_2"
	NotDefined   string = "not_defined"
	NotSupported string = "not_supported"
)

// Constants used for optionality
const (
	Mandatory string = "mandatory"
	Optional  string = "optional"
)

// CRD Kinds
const (
	KindAuthentication  = "Authentication"
	KindAPIPolicy       = "APIPolicy"
	KindScope           = "Scope"
	KindRateLimitPolicy = "RateLimitPolicy"
	KindService         = "Service"
	KindBackend         = "Backend"
)

// API environment types
const (
	Production = "Production"
	Sandbox    = "Sandbox"
)

// Security scheme constants
const (
	OAuth2 string = "OAuth2"
	APIKey string = "APIKey"
	JWT    string = "JWT"
)

// Guardrail related constants
const (
	// Azure Content Safety Content Moderation related constants
	AzureContentSafetyContentModeration = "AzureContentSafetyContentModeration"
	AzureContentSafetyKey               = "azureContentSafetyKey"

	// Semantic Caching related constants
	SemanticCaching               = "SemanticCache"
	SemanticCacheEmbeddingAPIKey  = "embeddingModelAPIKey"
	SemanticCacheVectorDBPassword = "vectorDBPassword"

	// AWS Bedrock Guardrail related constants
	AWSBedrockGuardrail = "AWSBedrockGuardrail"
	AWSAccessKeyID      = "awsAccessKeyID"
	AWSSecretAccessKey  = "awsSecretAccessKey"
	AWSSessionToken     = "awsSessionToken"
	AWSRoleExternalID   = "awsRoleExternalID"
)
