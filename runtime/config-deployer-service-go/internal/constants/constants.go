/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package constants

// API Type related constants.
const (
	API_TYPE_REST    string = "REST"
	API_TYPE_GRAPHQL string = "GRAPHQL"
	API_TYPE_GRPC    string = "GRPC"
	API_TYPE_ASYNC   string = "ASYNC"
	API_TYPE_SOAP    string = "SOAP"
	API_TYPE_SSE     string = "SSE"
	API_TYPE_WS      string = "WS"
	API_TYPE_WEBSUB  string = "WEBSUB"
)

// ALLOWED_API_TYPES is a list of allowed API types.
var ALLOWED_API_TYPES = []string{
	API_TYPE_REST,
	API_TYPE_GRAPHQL,
	API_TYPE_GRPC,
}

const (
	JAVA_IO_TMPDIR               string = "java.io.tmpdir"
	OPENAPI_ARCHIVES_TEMP_FOLDER string = "OPENAPI-archives"
	OPENAPI_ARCHIVE_ZIP_FILE     string = "openapi-archive.zip"
	OPENAPI_EXTRACTED_DIRECTORY  string = "extracted"
	OPENAPI_MASTER_JSON          string = "swagger.json"
	OPENAPI_MASTER_YAML          string = "swagger.yaml"
)
const (
	GRAPHQL_QUERY        = "QUERY"
	GRAPHQL_MUTATION     = "MUTATION"
	GRAPHQL_SUBSCRIPTION = "SUBSCRIPTION"
)

// SupportedMethods Supported HTTP methods
var SupportedMethods = map[string]bool{
	"get":     true,
	"put":     true,
	"post":    true,
	"delete":  true,
	"patch":   true,
	"head":    true,
	"options": true,
}

// GraphQLSupportedMethods GraphQL supported methods
var GraphQLSupportedMethods = map[string]bool{
	"QUERY":        true,
	"MUTATION":     true,
	"SUBSCRIPTION": true,
	"HEAD":         true,
	"OPTIONS":      true,
}

const SWAGGER_X_SCOPE = "x-scope"

type SwaggerVersion int

const (
	SWAGGER SwaggerVersion = iota
	OPEN_API
)

const (
	OPENAPI_RESOURCE_KEY = "paths"
)

const (
	OPENAPI_SECURITY_SCHEMA_KEY = "default"
	OAUTH2_SECURITY_SCHEMA_KEY  = "OAuth2Security"
)

var UnsupportedResourceBlocks = []string{"servers"}

// OpenAPI validation constants
const (
	OPENAPI_ALLOWED_EXTRA_SIBLING_FIELDS = "type"
)

const (
	ZIP_FILE_EXTENSION = ".zip"
)

const (
	PRODUCTION_TYPE  = "production"
	SANDBOX_TYPE     = "sandbox"
	INTERCEPTOR_TYPE = "interceptor"
)

const (
	API_NAME_HASH_LABEL         = "api-name"
	API_VERSION_HASH_LABEL      = "api-version"
	ORGANIZATION_HASH_LABEL     = "organization"
	MANAGED_BY_HASH_LABEL       = "managed-by"
	MANAGED_BY_HASH_LABEL_VALUE = "kgw"
	CP_INITIATED_HASH_LABEL     = "cp-initiated"
)

const (
	ValidatedUserContext = "VALIDATED_USER_CONTEXT"
)

const (
	EnvoyGatewayBackendTrafficPolicy           = "BackendTrafficPolicy"
	EnvoyGatewayBackendTrafficPolicyAPIVersion = "gateway.envoyproxy.io/v1alpha1"
	EnvoyGatewayExtensionPolicy                = "EnvoyExtensionPolicy"
	EnvoyGatewayExtensionPolicyAPIVersion      = "gateway.envoyproxy.io/v1alpha1"
	EnvoyGatewayHTTPRouteFilter                = "HTTPRouteFilter"
	EnvoyGatewayHTTPRouteFilterAPIVersion      = "gateway.envoyproxy.io/v1alpha1"

	WSO2KubernetesGatewayRouteMetadataAPIVersion = "dp.wso2.com/v2alpha1"
	WSO2KubernetesGatewayRouteMetadataKind       = "RouteMetadata"
	WSO2KubernetesGatewayRouteMetadataGroup      = "dp.wso2.com"
	WSO2KubernetesGatewayRoutePolicyAPIVersion   = "dp.wso2.com/v2alpha1"
	WSO2KubernetesGatewayRoutePolicyKind         = "RoutePolicy"
	WSO2KubernetesGatewayRoutePolicyGroup        = "dp.wso2.com"

	K8sKindConfigMap              = "ConfigMap"
	K8sKindHTTPRoute              = "HTTPRoute"
	K8sKindGRPCRoute              = "GRPCRoute"
	K8sKindService                = "Service"
	K8sKindGateway                = "Gateway"
	K8sKindBackend                = "Backend"
	K8sKindSecurityPolicy         = "SecurityPolicy"
	K8sKindBackendTLSPolicy       = "BackendTLSPolicy"
	K8sAPIVersionHTTPRoute        = "gateway.networking.k8s.io/v1"
	K8sGroupEnvoyGateway          = "gateway.envoyproxy.io"
	K8sAPIVersionEnvoyGateway     = "gateway.envoyproxy.io/v1alpha1"
	K8sGroupNetworking            = "gateway.networking.k8s.io"
	K8sGatewayNamespace           = "choreo-egress-gateway"
	K8sAPIVersionBackendTLSPolicy = "gateway.networking.k8s.io/v1alpha3"
	K8sKindGatewayExtensionPolicy = "GatewayExtensionPolicy"
	K8sAPIVersionGatewayExtension = "gateway.choreo.dev/v1alpha1"
	K8sGroupGatewayExtension      = "gateway.choreo.dev"
	K8sKindReferenceGrant         = "ReferenceGrant"
	K8sAPIVersionReferenceGrant   = "gateway.networking.k8s.io/v1beta1"

	K8sMaxAnnotationLength = 10000
)
