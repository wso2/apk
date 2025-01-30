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

package envoyconf

const (
	extAuthzClusterName     string = "ext-authz"
	accessLoggerClusterName string = "access-logger"
	grpcAccessLogLogName    string = "apk_access_logs"
	tracingClusterName      string = "wso2_apk_trace"
	rateLimitClusterName    string = "ratelimit"
)

const (
	httpConManagerStartPrefix  string = "ingress_http"
	extAuthzPerRouteName       string = "type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthzPerRoute"
	extProcPerRouteName        string = "type.googleapis.com/envoy.extensions.filters.http.ext_proc.v3.ExtProcPerRoute"
	ratelimitPerRouteName      string = "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimitPerRoute"
	luaPerRouteName            string = "type.googleapis.com/envoy.extensions.filters.http.lua.v3.LuaPerRoute"
	corsFilterName             string = "type.googleapis.com/envoy.extensions.filters.http.cors.v3.Cors"
	localRateLimitPerRouteName string = "type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit"
	httpProtocolOptionsName    string = "type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions"
	apkWebSocketWASMFilterName string = "envoy.filters.http.mgw_WASM_websocket"
	apkWASMVmID                string = "mgw_WASM_vm"
	apkWASMVmRuntime           string = "envoy.wasm.runtime.v8"
	apkWebSocketWASMFilterRoot string = "mgw_WASM_websocket_root"
	apkWebSocketWASM           string = "/home/wso2/wasm/websocket/mgw-websocket.wasm"
	compressorFilterName       string = "envoy.filters.http.compressor"
)

// cluster prefixes
const (
	requestInterceptClustersNamePrefix  string = "reqInterceptor"
	responseInterceptClustersNamePrefix string = "resInterceptor"
)

// Context Extensions which are set in ExtAuthzPerRoute Config
// These values are shared between the adapter and enforcer, hence if it is required to change
// these values, modifications should be done in the both adapter and enforcer.
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
	retryPolicyRetriableStatusCodes                 string = "retriable-status-codes"
)

const (
	// clusterHeaderName denotes the constant used for header based routing decisions.
	clusterHeaderName string = "x-wso2-cluster-header"
	// xWso2requestInterceptor used to provide request interceptor details for api and resource level
	xWso2requestInterceptor string = "x-wso2-request-interceptor"
	// xWso2responseInterceptor used to provide response interceptor details for api and resource level
	xWso2responseInterceptor string = "x-wso2-response-interceptor"
)

// interceptor levels
const (
	APILevelInterceptor       string = "api"
	ResourceLevelInterceptor  string = "resource"
	OperationLevelInterceptor string = "operation"
)
const (
	httpsURLType     string = "https"
	wssURLType       string = "wss"
	httpMethodHeader string = ":method"
)

// Paths exposed from the router by default
const (
	healthPath              string = "/health"
	readyPath               string = "/ready"
	apiDefinitionQueryParam string = "OAS"
	apiDefinitionPath       string = "/api-definition"
	jwksPath                string = "/.wellknown/jwks"
)

const (
	// healthEndpointResponse - response from the health endpoint
	healthEndpointResponse = "{\"status\": \"healthy\"}"
)

const (
	defaultListenerHostAddress = "0.0.0.0"
)

// tracing configuration constants
const (
	tracerHost              = "host"
	tracerPort              = "port"
	tracerMaxPathLength     = "maxPathLength"
	tracerEndpoint          = "endpoint"
	tracerNameZipkin        = "envoy.tracers.zipkin"
	tracerNameOpenTelemetry = "envoy.tracers.opentelemetry"
	tracerConnectionTimeout = "connectionTimeout"
	tracerServiceNameRouter = "apk_router-default"
	// Azure tracer's name
	TracerTypeAzure = "azure"
	TracerTypeOtlp  = "otlp"
)

// Constants used for SOAP APIs
const (
	contentTypeHeaderName = "content-type"
	contentTypeHeaderXML  = "text/xml"
	contentTypeHeaderSoap = "application/soap+xml"
	soap11ProtocolVersion = "SOAP 1.1 Protocol"
	soap12ProtocolVersion = "SOAP 1.2 Protocol"
	soapActionHeaderName  = "SOAPAction"
)

// metadata keys
const (
	methodRewrite = "method-rewrite"
)

// Enforcer
const (
	apkEnforcerReply = "apk-enforcer-reply"
	uaexCode         = "UAEX"
)

// Constants relevant to the ratelimit service
const (
	RateLimiterDomain                    = "Default"
	RateLimitPolicyOperationLevel string = "OPERATION"
	RateLimitPolicyAPILevel       string = "API"
)

// LuaGlobal is the lua filter name for global lua filter
const LuaGlobal = "envoy.filters.http.lua.global"

// LuaLocal is the lua filter name for local lua filter
const LuaLocal = "envoy.filters.http.lua.local"

// EnvoyJWT is the jwt filter name
const EnvoyJWT = "envoy.filters.http.jwt_authn"
