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

package model

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/interceptor"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
)

// AdapterInternalAPI represents the object structure holding the information related to the
// adapter internal representation. The values are populated from the operator. The pathItem level information is represented
// by the resources array which contains the Resource entries.
type AdapterInternalAPI struct {
	id                       string
	UUID                     string
	apiType                  string
	description              string
	title                    string
	version                  string
	vendorExtensions         map[string]interface{}
	xWso2Endpoints           map[string]*EndpointCluster
	resources                []*Resource
	xWso2Basepath            string
	xWso2HTTP2BackendEnabled bool
	xWso2Cors                *CorsConfig
	xWso2ThrottlingTier      string
	xWso2AuthHeader          string
	disableAuthentications   bool
	disableScopes            bool
	OrganizationID           string
	IsPrototyped             bool
	EndpointType             string
	LifecycleStatus          string
	xWso2RequestBodyPass     bool
	IsDefaultVersion         bool
	clientCertificates       []Certificate
	xWso2MutualSSL           string
	xWso2ApplicationSecurity bool
	EnvType                  string
	backendJWTTokenInfo      *BackendJWTTokenInfo
	apiDefinitionFile        []byte
	apiDefinitionEndpoint    string
	APIProperties            []dpv1alpha1.Property
	// GraphQLSchema              string
	// GraphQLComplexities        GraphQLComplexityYaml
	IsSystemAPI     bool
	RateLimitPolicy *RateLimitPolicy
	environment     string
}

// BackendJWTTokenInfo represents the object structure holding the information related to the JWT Generator
type BackendJWTTokenInfo struct {
	Enabled          bool
	Encoding         string
	Header           string
	SigningAlgorithm string
	TokenTTL         uint32
	CustomClaims     []ClaimMapping
}

// ClaimMapping represents the object structure holding the information related to the JWT Generator Claims
type ClaimMapping struct {
	Claim string
	Value ClaimVal
}

// ClaimVal represents the object structure holding the information related to the JWT Generator Claim Values
type ClaimVal struct {
	Value string
	Type  string
}

// RateLimitPolicy information related to the rate limiting policy
type RateLimitPolicy struct {
	Count    uint32
	SpanUnit string
}

// EndpointCluster represent an upstream cluster
type EndpointCluster struct {
	EndpointPrefix string
	Endpoints      []Endpoint
	// EndpointType enum {failover, loadbalance}. if any other value provided, consider as the default value; which is loadbalance
	EndpointType string
	Config       *EndpointConfig
	HealthCheck  *HealthCheck
	// Is http2 protocol enabled
	HTTP2BackendEnabled bool
}

// Endpoint represents the structure of an endpoint.
type Endpoint struct {
	// Host name
	Host string
	// BasePath (which would be added as prefix to the path mentioned in openapi definition)
	// In openAPI v2, it is determined from the basePath property
	// In openAPi v3, it is determined from the server object's suffix
	Basepath string
	// https, http, ws, wss
	// In openAPI v2, it is fetched from the schemes entry
	// In openAPI v3, it is extracted from the server property under servers object
	// only https and http are supported at the moment.
	URLType string
	// Port of the endpoint.
	// If the port is not specified, 80 is assigned if URLType is http
	// 443 is assigned if URLType is https
	Port uint32
	//ServiceDiscoveryQuery consul query for service discovery
	ServiceDiscoveryString string
	RawURL                 string
	// Trusted CA Cerificate for the endpoint
	Certificate []byte
	// Subject Alternative Names to verify in the public certificate
	AllowedSANs []string
}

// EndpointSecurity contains parameters of endpoint security at api.json
type EndpointSecurity struct {
	Password         string
	Type             string
	Enabled          bool
	Username         string
	CustomParameters map[string]string
}

// EndpointConfig holds the configs such as timeout, retry, etc. for the EndpointCluster
type EndpointConfig struct {
	RetryConfig          *RetryConfig
	TimeoutInMillis      uint32
	IdleTimeoutInSeconds uint32
	CircuitBreakers      *CircuitBreakers
}

// HealthCheck holds the parameters for health check done by apk to the EndpointCluster
type HealthCheck struct {
	Timeout            uint32
	Interval           uint32
	UnhealthyThreshold uint32
	HealthyThreshold   uint32
}

// RetryConfig holds the parameters for retries done by apk to the EndpointCluster
type RetryConfig struct {
	Count                int32
	StatusCodes          []uint32
	BaseIntervalInMillis int32
}

// CircuitBreakers holds the parameters for retries done by apk to the EndpointCluster
type CircuitBreakers struct {
	MaxConnections     int32
	MaxRequests        int32
	MaxPendingRequests int32
	MaxRetries         int32
	MaxConnectionPools int32
}

// SecurityScheme represents the structure of an security scheme.
type SecurityScheme struct {
	DefinitionName string // Arbitrary name used to define the security scheme. ex: default, myApikey
	Type           string // Type of the security scheme. Valid: apiKey, api_key, oauth2
	Name           string // Used for API key. Name of header or query. ex: x-api-key, apikey
	In             string // Where the api key found in. Valid: query, header
}

// CorsConfig represents the API level Cors Configuration
type CorsConfig struct {
	Enabled                       bool
	AccessControlAllowCredentials bool
	AccessControlAllowHeaders     []string
	AccessControlAllowMethods     []string
	AccessControlAllowOrigins     []string
	AccessControlExposeHeaders    []string
	AccessControlMaxAge           *int
}

// InterceptEndpoint contains the parameters of endpoint security
type InterceptEndpoint struct {
	Enable          bool
	EndpointCluster EndpointCluster
	ClusterName     string
	ClusterTimeout  time.Duration
	RequestTimeout  time.Duration
	// Level this is an enum allowing only values {api, resource, operation}
	// to indicate from which level interceptor is added
	Level string
	// Includes this is an enum allowing only values in
	// {"request_headers", "request_body", "request_trailer", "response_headers", "response_body", "response_trailer",
	//"invocation_context" }
	Includes *interceptor.RequestInclusions
}

// Certificate contains information of a client certificate
type Certificate struct {
	Alias   string
	Tier    string
	Content []byte
}

// GetAPIDefinitionFile returns the API Definition File.
func (swagger *AdapterInternalAPI) GetAPIDefinitionFile() []byte {
	return swagger.apiDefinitionFile
}

// GetAPIDefinitionEndpoint returns the API Definition Endpoint.
func (swagger *AdapterInternalAPI) GetAPIDefinitionEndpoint() string {
	return swagger.apiDefinitionEndpoint
}

// GetBackendJWTTokenInfo returns the BackendJWTTokenInfo Object.
func (swagger *AdapterInternalAPI) GetBackendJWTTokenInfo() *BackendJWTTokenInfo {
	return swagger.backendJWTTokenInfo
}

// GetCorsConfig returns the CorsConfiguration Object.
func (swagger *AdapterInternalAPI) GetCorsConfig() *CorsConfig {
	return swagger.xWso2Cors
}

// GetAPIType returns the openapi version
func (swagger *AdapterInternalAPI) GetAPIType() string {
	return swagger.apiType
}

// GetVersion returns the API version
func (swagger *AdapterInternalAPI) GetVersion() string {
	return swagger.version
}

// GetTitle returns the API Title
func (swagger *AdapterInternalAPI) GetTitle() string {
	return swagger.title
}

// GetXWso2Basepath returns the basepath set via the vendor extension.
func (swagger *AdapterInternalAPI) GetXWso2Basepath() string {
	return swagger.xWso2Basepath
}

// GetXWso2HTTP2BackendEnabled returns the http2 backend enabled set via the vendor extension.
func (swagger *AdapterInternalAPI) GetXWso2HTTP2BackendEnabled() bool {
	return swagger.xWso2HTTP2BackendEnabled
}

// GetVendorExtensions returns the map of vendor extensions which are defined
// at openAPI's root level.
func (swagger *AdapterInternalAPI) GetVendorExtensions() map[string]interface{} {
	return swagger.vendorExtensions
}

// GetXWso2Endpoints returns the array of x wso2 endpoints.
func (swagger *AdapterInternalAPI) GetXWso2Endpoints() map[string]*EndpointCluster {
	return swagger.xWso2Endpoints
}

// GetResources returns the array of resources (openAPI path level info)
func (swagger *AdapterInternalAPI) GetResources() []*Resource {
	return swagger.resources
}

// GetDescription returns the description of the openapi
func (swagger *AdapterInternalAPI) GetDescription() string {
	return swagger.description
}

// GetXWso2ThrottlingTier returns the Throttling tier via the vendor extension.
func (swagger *AdapterInternalAPI) GetXWso2ThrottlingTier() string {
	return swagger.xWso2ThrottlingTier
}

// GetDisableAuthentications returns the authType via the vendor extension.
func (swagger *AdapterInternalAPI) GetDisableAuthentications() bool {
	return swagger.disableAuthentications
}

// GetDisableScopes returns the authType via the vendor extension.
func (swagger *AdapterInternalAPI) GetDisableScopes() bool {
	return swagger.disableScopes
}

// GetID returns the Id of the API
func (swagger *AdapterInternalAPI) GetID() string {
	return swagger.id
}

// GetXWso2RequestBodyPass returns boolean value to indicate
// whether it is allowed to pass request body to the enforcer or not.
func (swagger *AdapterInternalAPI) GetXWso2RequestBodyPass() bool {
	return swagger.xWso2RequestBodyPass
}

// GetClientCerts returns the client certificates of the API
func (swagger *AdapterInternalAPI) GetClientCerts() []Certificate {
	return swagger.clientCertificates
}

// SetClientCerts set the client certificates of the API
func (swagger *AdapterInternalAPI) SetClientCerts(certs []Certificate) {
	swagger.clientCertificates = certs
}

// SetID set the Id of the API
func (swagger *AdapterInternalAPI) SetID(id string) {
	swagger.id = id
}

// SetAPIDefinitionFile sets the API Definition File.
func (swagger *AdapterInternalAPI) SetAPIDefinitionFile(file []byte) {
	swagger.apiDefinitionFile = file
}

// SetAPIDefinitionEndpoint sets the API Definition Endpoint.
func (swagger *AdapterInternalAPI) SetAPIDefinitionEndpoint(endpoint string) {
	swagger.apiDefinitionEndpoint = endpoint
}

// SetName sets the name of the API
func (swagger *AdapterInternalAPI) SetName(name string) {
	swagger.title = name
}

// SetVersion sets the version of the API
func (swagger *AdapterInternalAPI) SetVersion(version string) {
	swagger.version = version
}

// SetIsDefaultVersion sets whether this API is the default
func (swagger *AdapterInternalAPI) SetIsDefaultVersion(isDefaultVersion bool) {
	swagger.IsDefaultVersion = isDefaultVersion
}

// SetXWso2AuthHeader sets the authHeader of the API
func (swagger *AdapterInternalAPI) SetXWso2AuthHeader(authHeader string) {
	if swagger.xWso2AuthHeader == "" {
		swagger.xWso2AuthHeader = authHeader
	}
}

// GetXWSO2AuthHeader returns the auth header set via the vendor extension.
func (swagger *AdapterInternalAPI) GetXWSO2AuthHeader() string {
	return swagger.xWso2AuthHeader
}

// SetXWSO2MutualSSL sets the optional or mandatory mTLS
func (swagger *AdapterInternalAPI) SetXWSO2MutualSSL(mutualSSl string) {
	swagger.xWso2MutualSSL = mutualSSl
}

// GetXWSO2MutualSSL returns the optional or mandatory mTLS
func (swagger *AdapterInternalAPI) GetXWSO2MutualSSL() string {
	return swagger.xWso2MutualSSL
}

// SetXWSO2ApplicationSecurity sets the optional or mandatory application security
func (swagger *AdapterInternalAPI) SetXWSO2ApplicationSecurity(applicationSecurity bool) {
	swagger.xWso2ApplicationSecurity = applicationSecurity
}

// GetXWSO2ApplicationSecurity returns the optional or mandatory application security
func (swagger *AdapterInternalAPI) GetXWSO2ApplicationSecurity() bool {
	return swagger.xWso2ApplicationSecurity
}

// GetOrganizationID returns OrganizationID
func (swagger *AdapterInternalAPI) GetOrganizationID() string {
	return swagger.OrganizationID
}

// SetEnvironment sets the environment of the API.
func (swagger *AdapterInternalAPI) SetEnvironment(environment string) {
	swagger.environment = environment
}

// GetEnvironment returns the environment of the API
func (swagger *AdapterInternalAPI) GetEnvironment() string {
	return swagger.environment
}

// Validate method confirms that the adapterInternalAPI has all required fields in the required format.
// This needs to be checked prior to generate router/enforcer related resources.
func (swagger *AdapterInternalAPI) Validate() error {
	for _, res := range swagger.resources {
		if res.endpoints == nil || len(res.endpoints.Endpoints) == 0 {
			logger.LoggerOasparser.Errorf("No Endpoints are provided for the resources in %s:%s, API_UUID: %v",
				swagger.title, swagger.version, swagger.UUID)
			return errors.New("no endpoints are provided for the API")
		}
		err := res.endpoints.validateEndpointCluster()
		if err != nil {
			logger.LoggerOasparser.Errorf("Error while parsing the endpoints of the API %s:%s - %v, API_UUID: %v",
				swagger.title, swagger.version, err, swagger.UUID)
			return err
		}
	}
	return nil
}

func (endpoint *Endpoint) validateEndpoint() error {
	if len(endpoint.ServiceDiscoveryString) > 0 {
		return nil
	}
	if endpoint.Port == 0 || endpoint.Port > 65535 {
		return errors.New("endpoint port value should be between 0 and 65535")
	}
	if len(endpoint.Host) == 0 {
		return errors.New("empty Hostname is provided")
	}
	if strings.HasPrefix(endpoint.Host, "/") {
		return errors.New("relative paths are not supported as endpoint URLs")
	}
	urlString := endpoint.URLType + "://" + endpoint.Host
	_, err := url.ParseRequestURI(urlString)
	return err
}

// GetAuthorityHeader creates the authority header using Host and Port in the form of Host [ ":" Port ]
func (endpoint *Endpoint) GetAuthorityHeader() string {
	return strings.Join([]string{endpoint.Host, strconv.FormatUint(uint64(endpoint.Port), 10)}, ":")
}

func (retryConfig *RetryConfig) validateRetryConfig() {
	conf := config.ReadConfigs()
	var validStatusCodes []uint32
	for _, statusCode := range retryConfig.StatusCodes {
		if statusCode > 598 || statusCode < 401 {
			logger.LoggerOasparser.Errorf("Given status code for the API retry config is invalid." +
				"Must be in the range 401 - 598. Dropping the status code.")
		} else {
			validStatusCodes = append(validStatusCodes, statusCode)
		}
	}
	if len(validStatusCodes) < 1 {
		validStatusCodes = append(validStatusCodes, conf.Envoy.Upstream.Retry.StatusCodes...)
	}
	retryConfig.StatusCodes = validStatusCodes
}

func (endpointCluster *EndpointCluster) validateEndpointCluster() error {
	if endpointCluster != nil && len(endpointCluster.Endpoints) > 0 {
		var err error
		for _, endpoint := range endpointCluster.Endpoints {
			err = endpoint.validateEndpoint()
			if err != nil {
				logger.LoggerOasparser.Errorf("Error while parsing the endpoint. %v", err)
				return err
			}
		}

		if endpointCluster.Config != nil {
			// Validate retry
			if endpointCluster.Config.RetryConfig != nil {
				endpointCluster.Config.RetryConfig.validateRetryConfig()
			}
		}
	}
	return nil
}

func generateEndpointCluster(endpoints []Endpoint, endpointType string) *EndpointCluster {
	if len(endpoints) > 0 {
		endpointCluster := EndpointCluster{
			Endpoints:    endpoints,
			EndpointType: endpointType,
		}
		return &endpointCluster
	}
	return nil
}

// GetOperationInterceptors returns operation interceptors
func (swagger *AdapterInternalAPI) GetOperationInterceptors(apiInterceptor InterceptEndpoint, resourceInterceptor InterceptEndpoint, operations []*Operation, isIn bool) map[string]InterceptEndpoint {
	interceptorOperationMap := make(map[string]InterceptEndpoint)

	for _, op := range operations {
		extensionName := constants.XWso2RequestInterceptor
		// first get operational policies
		operationInterceptor := op.GetCallInterceptorService(isIn)
		// if operational policy interceptor not given check operational level swagger extension
		if !operationInterceptor.Enable {
			if !isIn {
				extensionName = constants.XWso2ResponseInterceptor
			}
			operationInterceptor = swagger.GetInterceptor(op.GetVendorExtensions(), extensionName, constants.OperationLevelInterceptor)
		}
		operationInterceptor.ClusterName = op.iD
		// if operation interceptor not given
		if !operationInterceptor.Enable {
			// assign resource level interceptor
			if resourceInterceptor.Enable {
				operationInterceptor = resourceInterceptor
			} else if apiInterceptor.Enable {
				// if resource interceptor not given add api level interceptor
				operationInterceptor = apiInterceptor
			}
		}
		// add operation to the list only if an interceptor is enabled for the operation
		if operationInterceptor.Enable {
			interceptorOperationMap[strings.ToUpper(op.method)] = operationInterceptor
		}
	}
	return interceptorOperationMap

}

// GetInterceptor returns interceptors
func (swagger *AdapterInternalAPI) GetInterceptor(vendorExtensions map[string]interface{}, extensionName string, level string) InterceptEndpoint {
	var endpointCluster EndpointCluster
	conf := config.ReadConfigs()
	clusterTimeoutV := conf.Envoy.ClusterTimeoutInSeconds
	requestTimeoutV := conf.Envoy.ClusterTimeoutInSeconds
	includesV := &interceptor.RequestInclusions{}

	if x, found := vendorExtensions[extensionName]; found {
		if val, ok := x.(map[string]interface{}); ok {
			//serviceURL mandatory
			if v, found := val[constants.ServiceURL]; found {
				serviceURLV := v.(string)
				endpoint, err := getHTTPEndpoint(serviceURLV)
				if err != nil {
					logger.LoggerOasparser.Error("Error reading interceptors service url value", err)
					return InterceptEndpoint{}
				}
				if endpoint.Basepath != "" {
					logger.LoggerOasparser.Warnf("Interceptor serviceURL basepath is given as %v but it will be ignored",
						endpoint.Basepath)
				}
				endpointCluster.Endpoints = []Endpoint{*endpoint}

			} else {
				logger.LoggerOasparser.Error("Error reading interceptors service url value")
				return InterceptEndpoint{}
			}
			//clusterTimeout optional
			if v, found := val[constants.ClusterTimeout]; found {
				p, err := strconv.ParseInt(fmt.Sprint(v), 0, 0)
				if err == nil {
					clusterTimeoutV = time.Duration(p)
				} else {
					logger.LoggerOasparser.Errorf("Error reading interceptors %v value : %v", constants.ClusterTimeout, err.Error())
				}
			}
			//requestTimeout optional
			if v, found := val[constants.RequestTimeout]; found {
				p, err := strconv.ParseInt(fmt.Sprint(v), 0, 0)
				if err == nil {
					requestTimeoutV = time.Duration(p)
				} else {
					logger.LoggerOasparser.Errorf("Error reading interceptors %v value : %v", constants.RequestTimeout, err.Error())
				}
			}
			//includes optional
			if v, found := val[constants.Includes]; found {
				includes := v.([]interface{})
				if len(includes) > 0 {
					// convert type of includes from "[]interface{}" to "[]dpv1alpha1.InterceptorInclusion"
					includes := make([]dpv1alpha1.InterceptorInclusion, len(includes))
					includesV = GenerateInterceptorIncludes(includes)
				}
			}

			return InterceptEndpoint{
				Enable:          true,
				EndpointCluster: endpointCluster,
				ClusterTimeout:  clusterTimeoutV,
				RequestTimeout:  requestTimeoutV,
				Includes:        includesV,
				Level:           level,
			}
		}
		logger.LoggerOasparser.Error("Error parsing response interceptors values to adapterInternalAPI")
	}
	return InterceptEndpoint{}
}

// GenerateInterceptorIncludes generate includes
func GenerateInterceptorIncludes(includes []dpv1alpha1.InterceptorInclusion) *interceptor.RequestInclusions {
	includesV := &interceptor.RequestInclusions{}
	for _, include := range includes {
		switch include {
		case dpv1alpha1.InterceptorInclusionRequestHeaders:
			includesV.RequestHeaders = true
		case dpv1alpha1.InterceptorInclusionRequestBody:
			includesV.RequestBody = true
		case dpv1alpha1.InterceptorInclusionRequestTrailers:
			includesV.RequestTrailer = true
		case dpv1alpha1.InterceptorInclusionResponseHeaders:
			includesV.ResponseHeaders = true
		case dpv1alpha1.InterceptorInclusionResponseBody:
			includesV.ResponseBody = true
		case dpv1alpha1.InterceptorInclusionResponseTrailers:
			includesV.ResponseTrailers = true
		case dpv1alpha1.InterceptorInclusionInvocationContext:
			includesV.InvocationContext = true
		}
	}
	return includesV
}

// CreateDummyAdapterInternalAPIForTests creates a dummy AdapterInternalAPI struct to be used for unit tests
func CreateDummyAdapterInternalAPIForTests(title, version, basePath string, resources []*Resource) *AdapterInternalAPI {
	return &AdapterInternalAPI{
		title:         title,
		version:       version,
		xWso2Basepath: basePath,
		resources:     resources,
	}
}
