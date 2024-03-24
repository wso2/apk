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

	"github.com/google/uuid"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/interceptor"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
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
	resources                []*Resource
	xWso2Basepath            string
	xWso2HTTP2BackendEnabled bool
	xWso2Cors                *CorsConfig
	xWso2ThrottlingTier      string
	xWso2AuthHeader          string
	disableAuthentications   bool
	disableScopes            bool
	disableMtls              bool
	OrganizationID           string
	IsPrototyped             bool
	EndpointType             string
	LifecycleStatus          string
	xWso2RequestBodyPass     bool
	IsDefaultVersion         bool
	clientCertificates       []Certificate
	mutualSSL                string
	xWso2ApplicationSecurity bool
	EnvType                  string
	backendJWTTokenInfo      *BackendJWTTokenInfo
	apiDefinitionFile        []byte
	apiDefinitionEndpoint    string
	subscriptionValidation   bool
	APIProperties            []dpv1alpha2.Property
	// GraphQLSchema              string
	// GraphQLComplexities        GraphQLComplexityYaml
	IsSystemAPI      bool
	RateLimitPolicy  *RateLimitPolicy
	environment      string
	Endpoints        *EndpointCluster
	EndpointSecurity []*EndpointSecurity
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
	Port   uint32
	RawURL string
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
	Content []byte
}

// GetAPIDefinitionFile returns the API Definition File.
func (adapterInternalAPI *AdapterInternalAPI) GetAPIDefinitionFile() []byte {
	return adapterInternalAPI.apiDefinitionFile
}

// GetAPIDefinitionEndpoint returns the API Definition Endpoint.
func (adapterInternalAPI *AdapterInternalAPI) GetAPIDefinitionEndpoint() string {
	return adapterInternalAPI.apiDefinitionEndpoint
}

// GetSubscriptionValidation returns the subscription validation status.
func (adapterInternalAPI *AdapterInternalAPI) GetSubscriptionValidation() bool {
	return adapterInternalAPI.subscriptionValidation
}

// GetBackendJWTTokenInfo returns the BackendJWTTokenInfo Object.
func (adapterInternalAPI *AdapterInternalAPI) GetBackendJWTTokenInfo() *BackendJWTTokenInfo {
	return adapterInternalAPI.backendJWTTokenInfo
}

// GetCorsConfig returns the CorsConfiguration Object.
func (adapterInternalAPI *AdapterInternalAPI) GetCorsConfig() *CorsConfig {
	return adapterInternalAPI.xWso2Cors
}

// GetAPIType returns the openapi version
func (adapterInternalAPI *AdapterInternalAPI) GetAPIType() string {
	return adapterInternalAPI.apiType
}

// GetVersion returns the API version
func (adapterInternalAPI *AdapterInternalAPI) GetVersion() string {
	return adapterInternalAPI.version
}

// GetTitle returns the API Title
func (adapterInternalAPI *AdapterInternalAPI) GetTitle() string {
	return adapterInternalAPI.title
}

// GetXWso2Basepath returns the basepath set via the vendor extension.
func (adapterInternalAPI *AdapterInternalAPI) GetXWso2Basepath() string {
	return adapterInternalAPI.xWso2Basepath
}

// GetXWso2HTTP2BackendEnabled returns the http2 backend enabled set via the vendor extension.
func (adapterInternalAPI *AdapterInternalAPI) GetXWso2HTTP2BackendEnabled() bool {
	return adapterInternalAPI.xWso2HTTP2BackendEnabled
}

// GetVendorExtensions returns the map of vendor extensions which are defined
// at openAPI's root level.
func (adapterInternalAPI *AdapterInternalAPI) GetVendorExtensions() map[string]interface{} {
	return adapterInternalAPI.vendorExtensions
}

// GetResources returns the array of resources (openAPI path level info)
func (adapterInternalAPI *AdapterInternalAPI) GetResources() []*Resource {
	return adapterInternalAPI.resources
}

// GetDescription returns the description of the openapi
func (adapterInternalAPI *AdapterInternalAPI) GetDescription() string {
	return adapterInternalAPI.description
}

// GetXWso2ThrottlingTier returns the Throttling tier via the vendor extension.
func (adapterInternalAPI *AdapterInternalAPI) GetXWso2ThrottlingTier() string {
	return adapterInternalAPI.xWso2ThrottlingTier
}

// GetDisableAuthentications returns the authType via the vendor extension.
func (adapterInternalAPI *AdapterInternalAPI) GetDisableAuthentications() bool {
	return adapterInternalAPI.disableAuthentications
}

// GetDisableScopes returns the authType via the vendor extension.
func (adapterInternalAPI *AdapterInternalAPI) GetDisableScopes() bool {
	return adapterInternalAPI.disableScopes
}

// GetDisableMtls returns whether mTLS is disabled or not
func (adapterInternalAPI *AdapterInternalAPI) GetDisableMtls() bool {
	return adapterInternalAPI.disableMtls
}

// GetID returns the Id of the API
func (adapterInternalAPI *AdapterInternalAPI) GetID() string {
	return adapterInternalAPI.id
}

// GetXWso2RequestBodyPass returns boolean value to indicate
// whether it is allowed to pass request body to the enforcer or not.
func (adapterInternalAPI *AdapterInternalAPI) GetXWso2RequestBodyPass() bool {
	return adapterInternalAPI.xWso2RequestBodyPass
}

// SetXWso2RequestBodyPass returns boolean value to indicate
// whether it is allowed to pass request body to the enforcer or not.
func (adapterInternalAPI *AdapterInternalAPI) SetXWso2RequestBodyPass(passBody bool) {
	adapterInternalAPI.xWso2RequestBodyPass = passBody
}

// GetClientCerts returns the client certificates of the API
func (adapterInternalAPI *AdapterInternalAPI) GetClientCerts() []Certificate {
	return adapterInternalAPI.clientCertificates
}

// SetClientCerts set the client certificates of the API
func (adapterInternalAPI *AdapterInternalAPI) SetClientCerts(apiName string, certs []string) {
	var clientCerts []Certificate
	for i, cert := range certs {
		clientCert := Certificate{
			Alias:   apiName + "-cert-" + strconv.Itoa(i),
			Content: []byte(cert),
		}
		clientCerts = append(clientCerts, clientCert)
	}
	adapterInternalAPI.clientCertificates = clientCerts
}

// SetID set the Id of the API
func (adapterInternalAPI *AdapterInternalAPI) SetID(id string) {
	adapterInternalAPI.id = id
}

// SetAPIDefinitionFile sets the API Definition File.
func (adapterInternalAPI *AdapterInternalAPI) SetAPIDefinitionFile(file []byte) {
	adapterInternalAPI.apiDefinitionFile = file
}

// SetAPIDefinitionEndpoint sets the API Definition Endpoint.
func (adapterInternalAPI *AdapterInternalAPI) SetAPIDefinitionEndpoint(endpoint string) {
	adapterInternalAPI.apiDefinitionEndpoint = endpoint
}

// SetSubscriptionValidation sets the subscription validation status.
func (adapterInternalAPI *AdapterInternalAPI) SetSubscriptionValidation(subscriptionValidation bool) {
	adapterInternalAPI.subscriptionValidation = subscriptionValidation
}

// SetName sets the name of the API
func (adapterInternalAPI *AdapterInternalAPI) SetName(name string) {
	adapterInternalAPI.title = name
}

// SetVersion sets the version of the API
func (adapterInternalAPI *AdapterInternalAPI) SetVersion(version string) {
	adapterInternalAPI.version = version
}

// SetIsDefaultVersion sets whether this API is the default
func (adapterInternalAPI *AdapterInternalAPI) SetIsDefaultVersion(isDefaultVersion bool) {
	adapterInternalAPI.IsDefaultVersion = isDefaultVersion
}

// SetXWso2AuthHeader sets the authHeader of the API
func (adapterInternalAPI *AdapterInternalAPI) SetXWso2AuthHeader(authHeader string) {
	if adapterInternalAPI.xWso2AuthHeader == "" {
		adapterInternalAPI.xWso2AuthHeader = authHeader
	}
}

// GetXWSO2AuthHeader returns the auth header set via the vendor extension.
func (adapterInternalAPI *AdapterInternalAPI) GetXWSO2AuthHeader() string {
	return adapterInternalAPI.xWso2AuthHeader
}

// SetMutualSSL sets the optional or mandatory mTLS
func (adapterInternalAPI *AdapterInternalAPI) SetMutualSSL(mutualSSL string) {
	adapterInternalAPI.mutualSSL = mutualSSL
}

// GetMutualSSL returns the optional or mandatory mTLS
func (adapterInternalAPI *AdapterInternalAPI) GetMutualSSL() string {
	return adapterInternalAPI.mutualSSL
}

// SetDisableMtls returns whether mTLS is disabled or not
func (adapterInternalAPI *AdapterInternalAPI) SetDisableMtls(disableMtls bool) {
	adapterInternalAPI.disableMtls = disableMtls
}

// SetXWSO2ApplicationSecurity sets the optional or mandatory application security
func (adapterInternalAPI *AdapterInternalAPI) SetXWSO2ApplicationSecurity(applicationSecurity bool) {
	adapterInternalAPI.xWso2ApplicationSecurity = applicationSecurity
}

// GetXWSO2ApplicationSecurity returns true if application security is mandatory, and false if optional
func (adapterInternalAPI *AdapterInternalAPI) GetXWSO2ApplicationSecurity() bool {
	return adapterInternalAPI.xWso2ApplicationSecurity
}

// GetOrganizationID returns OrganizationID
func (adapterInternalAPI *AdapterInternalAPI) GetOrganizationID() string {
	return adapterInternalAPI.OrganizationID
}

// SetEnvironment sets the environment of the API.
func (adapterInternalAPI *AdapterInternalAPI) SetEnvironment(environment string) {
	adapterInternalAPI.environment = environment
}

// GetEnvironment returns the environment of the API
func (adapterInternalAPI *AdapterInternalAPI) GetEnvironment() string {
	return adapterInternalAPI.environment
}

// Validate method confirms that the adapterInternalAPI has all required fields in the required format.
// This needs to be checked prior to generate router/enforcer related resources.
func (adapterInternalAPI *AdapterInternalAPI) Validate() error {
	for _, res := range adapterInternalAPI.resources {
		if res.endpoints == nil || len(res.endpoints.Endpoints) == 0 {
			loggers.LoggerOasparser.Errorf("No Endpoints are provided for the resources in %s:%s, API_UUID: %v",
				adapterInternalAPI.title, adapterInternalAPI.version, adapterInternalAPI.UUID)
			return errors.New("no endpoints are provided for the API")
		}
		err := res.endpoints.validateEndpointCluster()
		if err != nil {
			loggers.LoggerOasparser.Errorf("Error while parsing the endpoints of the API %s:%s - %v, API_UUID: %v",
				adapterInternalAPI.title, adapterInternalAPI.version, err, adapterInternalAPI.UUID)
			return err
		}
	}
	return nil
}

// SetInfoHTTPRouteCR populates resources and endpoints of adapterInternalAPI. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (adapterInternalAPI *AdapterInternalAPI) SetInfoHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, resourceParams ResourceParams) error {
	var resources []*Resource
	outputAuthScheme := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.AuthSchemes)))
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.APIPolicies)))
	outputRatelimitPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.RateLimitPolicies)))

	disableScopes := true
	config := config.ReadConfigs()

	var authScheme *dpv1alpha2.Authentication
	if outputAuthScheme != nil {
		authScheme = *outputAuthScheme
	}
	var apiPolicy *dpv1alpha2.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	var ratelimitPolicy *dpv1alpha1.RateLimitPolicy
	if outputRatelimitPolicy != nil {
		ratelimitPolicy = *outputRatelimitPolicy
	}

	for _, rule := range httpRoute.Spec.Rules {
		var endPoints []Endpoint
		var policies = OperationPolicies{}
		var circuitBreaker *dpv1alpha1.CircuitBreaker
		var healthCheck *dpv1alpha1.HealthCheck
		resourceAuthScheme := authScheme
		resourceAPIPolicy := apiPolicy
		resourceRatelimitPolicy := ratelimitPolicy
		var scopes []string
		var timeoutInMillis uint32
		var idleTimeoutInSeconds uint32
		isRetryConfig := false
		isRouteTimeout := false
		var backendRetryCount uint32
		var statusCodes []uint32
		statusCodes = append(statusCodes, config.Envoy.Upstream.Retry.StatusCodes...)
		var baseIntervalInMillis uint32
		hasURLRewritePolicy := false
		var securityConfig []EndpointSecurity
		backendBasePath := ""
		for _, backend := range rule.BackendRefs {
			backendName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}
			resolvedBackend, ok := resourceParams.BackendMapping[backendName.String()]
			if ok {
				if resolvedBackend.CircuitBreaker != nil {
					circuitBreaker = &dpv1alpha1.CircuitBreaker{
						MaxConnections:     resolvedBackend.CircuitBreaker.MaxConnections,
						MaxPendingRequests: resolvedBackend.CircuitBreaker.MaxPendingRequests,
						MaxRequests:        resolvedBackend.CircuitBreaker.MaxRequests,
						MaxRetries:         resolvedBackend.CircuitBreaker.MaxRetries,
						MaxConnectionPools: resolvedBackend.CircuitBreaker.MaxConnectionPools,
					}
				}
				if resolvedBackend.Timeout != nil {
					isRouteTimeout = true
					timeoutInMillis = resolvedBackend.Timeout.UpstreamResponseTimeout * 1000
					idleTimeoutInSeconds = resolvedBackend.Timeout.DownstreamRequestIdleTimeout
				}

				if resolvedBackend.Retry != nil {
					isRetryConfig = true
					backendRetryCount = resolvedBackend.Retry.Count
					baseIntervalInMillis = resolvedBackend.Retry.BaseIntervalMillis
					if len(resolvedBackend.Retry.StatusCodes) > 0 {
						statusCodes = resolvedBackend.Retry.StatusCodes
					}
				}
				if resolvedBackend.HealthCheck != nil {
					healthCheck = &dpv1alpha1.HealthCheck{
						Interval:           resolvedBackend.HealthCheck.Interval,
						Timeout:            resolvedBackend.HealthCheck.Timeout,
						UnhealthyThreshold: resolvedBackend.HealthCheck.UnhealthyThreshold,
						HealthyThreshold:   resolvedBackend.HealthCheck.HealthyThreshold,
					}
				}
				endPoints = append(endPoints, GetEndpoints(backendName, resourceParams.BackendMapping)...)
				backendBasePath = GetBackendBasePath(backendName, resourceParams.BackendMapping)
				switch resolvedBackend.Security.Type {
				case "Basic":
					securityConfig = append(securityConfig, EndpointSecurity{
						Password: string(resolvedBackend.Security.Basic.Password),
						Username: string(resolvedBackend.Security.Basic.Username),
						Type:     string(resolvedBackend.Security.Type),
						Enabled:  true,
					})
				}
			} else {
				return fmt.Errorf("backend: %s has not been resolved", backendName)
			}
		}
		for _, filter := range rule.Filters {
			switch filter.Type {
			case gwapiv1.HTTPRouteFilterURLRewrite:
				policyParameters := make(map[string]interface{})
				policyParameters[constants.RewritePathType] = filter.URLRewrite.Path.Type
				policyParameters[constants.IncludeQueryParams] = true

				switch filter.URLRewrite.Path.Type {
				case gwapiv1.FullPathHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = backendBasePath + *filter.URLRewrite.Path.ReplaceFullPath
				case gwapiv1.PrefixMatchHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = backendBasePath + *filter.URLRewrite.Path.ReplacePrefixMatch
				}

				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1.HTTPRouteFilterURLRewrite),
					Action:     constants.ActionRewritePath,
					Parameters: policyParameters,
				})
				hasURLRewritePolicy = true
			case gwapiv1.HTTPRouteFilterExtensionRef:
				if filter.ExtensionRef.Kind == constants.KindAuthentication {
					if ref, found := resourceParams.ResourceAuthSchemes[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: httpRoute.Namespace,
					}.String()]; found {
						resourceAuthScheme = concatAuthSchemes(authScheme, &ref)
					} else {
						return fmt.Errorf(`auth scheme: %s has not been resolved, spec.targetRef.kind should be 
						'Resource' in resource level Authentications`, filter.ExtensionRef.Name)
					}
				}
				if filter.ExtensionRef.Kind == constants.KindAPIPolicy {
					if ref, found := resourceParams.ResourceAPIPolicies[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: httpRoute.Namespace,
					}.String()]; found {
						resourceAPIPolicy = concatAPIPolicies(apiPolicy, &ref)
					} else {
						return fmt.Errorf(`apipolicy: %s has not been resolved, spec.targetRef.kind should be 
						'Resource' in resource level APIPolicies`, filter.ExtensionRef.Name)
					}
				}
				if filter.ExtensionRef.Kind == constants.KindScope {
					if ref, found := resourceParams.ResourceScopes[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: httpRoute.Namespace,
					}.String()]; found {
						scopes = ref.Spec.Names
						disableScopes = false
					} else {
						return fmt.Errorf("scope: %s has not been resolved in namespace %s", filter.ExtensionRef.Name, httpRoute.Namespace)
					}
				}
				if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
					if ref, found := resourceParams.ResourceRateLimitPolicies[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: httpRoute.Namespace,
					}.String()]; found {
						resourceRatelimitPolicy = concatRateLimitPolicies(ratelimitPolicy, &ref)
					} else {
						return fmt.Errorf(`ratelimitpolicy: %s has not been resolved, spec.targetRef.kind should be 
						'Resource' in resource level RateLimitPolicies`, filter.ExtensionRef.Name)
					}
				}
		}
		resourceAPIPolicy = concatAPIPolicies(resourceAPIPolicy, nil)
		resourceAuthScheme = concatAuthSchemes(resourceAuthScheme, nil)
		resourceRatelimitPolicy = concatRateLimitPolicies(resourceRatelimitPolicy, nil)
		addOperationLevelInterceptors(&policies, resourceAPIPolicy, resourceParams.InterceptorServiceMapping, resourceParams.BackendMapping, httpRoute.Namespace)

		loggers.LoggerOasparser.Debugf("Calculating auths for API ..., API_UUID = %v", adapterInternalAPI.UUID)
		apiAuth := getSecurity(resourceAuthScheme)
		if len(rule.BackendRefs) < 1 {
			return fmt.Errorf("no backendref were provided")
		}

		for _, match := range rule.Matches {
			if !hasURLRewritePolicy {
				policyParameters := make(map[string]interface{})
				if *match.Path.Type == gwapiv1.PathMatchPathPrefix {
					policyParameters[constants.RewritePathType] = gwapiv1.PrefixMatchHTTPPathModifier
				} else {
					policyParameters[constants.RewritePathType] = gwapiv1.FullPathHTTPPathModifier
				}
				policyParameters[constants.IncludeQueryParams] = true
				policyParameters[constants.RewritePathResourcePath] = strings.TrimSuffix(backendBasePath, "/") + *match.Path.Value
				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1.HTTPRouteFilterURLRewrite),
					Action:     constants.ActionRewritePath,
					Parameters: policyParameters,
				})
			}
			resourcePath := adapterInternalAPI.xWso2Basepath + *match.Path.Value
			resource := &Resource{path: resourcePath,
				methods: getAllowedOperations(match.Method, policies, apiAuth,
					parseRateLimitPolicyToInternal(resourceRatelimitPolicy), scopes),
				pathMatchType: *match.Path.Type,
				hasPolicies:   true,
				iD:            uuid.New().String(),
			}

			resource.endpoints = &EndpointCluster{
				Endpoints: endPoints,
			}

			endpointConfig := &EndpointConfig{}

			if isRouteTimeout {
				endpointConfig.TimeoutInMillis = timeoutInMillis
				endpointConfig.IdleTimeoutInSeconds = idleTimeoutInSeconds
			}
			if circuitBreaker != nil {
				endpointConfig.CircuitBreakers = &CircuitBreakers{
					MaxConnections:     int32(circuitBreaker.MaxConnections),
					MaxRequests:        int32(circuitBreaker.MaxRequests),
					MaxPendingRequests: int32(circuitBreaker.MaxPendingRequests),
					MaxRetries:         int32(circuitBreaker.MaxRetries),
					MaxConnectionPools: int32(circuitBreaker.MaxConnectionPools),
				}
			}
			if isRetryConfig {
				endpointConfig.RetryConfig = &RetryConfig{
					Count:                int32(backendRetryCount),
					StatusCodes:          statusCodes,
					BaseIntervalInMillis: int32(baseIntervalInMillis),
				}
			}
			if healthCheck != nil {
				resource.endpoints.HealthCheck = &HealthCheck{
					Interval:           healthCheck.Interval,
					Timeout:            healthCheck.Timeout,
					UnhealthyThreshold: healthCheck.UnhealthyThreshold,
					HealthyThreshold:   healthCheck.HealthyThreshold,
				}
			}
			if isRouteTimeout || circuitBreaker != nil || healthCheck != nil || isRetryConfig {
				resource.endpoints.Config = endpointConfig
			}
			resource.endpointSecurity = utils.GetPtrSlice(securityConfig)
			resources = append(resources, resource)
		}
	}

	ratelimitPolicy = concatRateLimitPolicies(ratelimitPolicy, nil)
	apiPolicy = concatAPIPolicies(apiPolicy, nil)
	authScheme = concatAuthSchemes(authScheme, nil)

	adapterInternalAPI.RateLimitPolicy = parseRateLimitPolicyToInternal(ratelimitPolicy)
	adapterInternalAPI.resources = resources
	adapterInternalAPI.xWso2Cors = getCorsConfigFromAPIPolicy(apiPolicy)
	if authScheme.Spec.Override != nil && authScheme.Spec.Override.Disabled != nil {
		adapterInternalAPI.disableAuthentications = *authScheme.Spec.Override.Disabled
	}

	authSpec := utils.SelectPolicy(&authScheme.Spec.Override, &authScheme.Spec.Default, nil, nil)
	if authSpec != nil && authSpec.AuthTypes != nil && authSpec.AuthTypes.Oauth2.Required != "" {
		adapterInternalAPI.SetXWSO2ApplicationSecurity(authSpec.AuthTypes.Oauth2.Required == "mandatory")
	} else {
		adapterInternalAPI.SetXWSO2ApplicationSecurity(true)
	}

	adapterInternalAPI.disableScopes = disableScopes

	// Check whether the API has a backend JWT token
	if apiPolicy != nil && apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.BackendJWTPolicy != nil {
		backendJWTPolicy := resourceParams.BackendJWTMapping[types.NamespacedName{
			Name:      apiPolicy.Spec.Override.BackendJWTPolicy.Name,
			Namespace: httpRoute.Namespace,
		}.String()].Spec
		adapterInternalAPI.backendJWTTokenInfo = parseBackendJWTTokenToInternal(backendJWTPolicy)
	}
	return nil
}

// SetInfoGQLRouteCR populates resources and endpoints of adapterInternalAPI. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (adapterInternalAPI *AdapterInternalAPI) SetInfoGQLRouteCR(gqlRoute *dpv1alpha2.GQLRoute, resourceParams ResourceParams) error {
	var resources []*Resource
	outputAuthScheme := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.AuthSchemes)))
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.APIPolicies)))
	outputRatelimitPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.RateLimitPolicies)))

	disableScopes := true
	config := config.ReadConfigs()

	var authScheme *dpv1alpha2.Authentication
	if outputAuthScheme != nil {
		authScheme = *outputAuthScheme
	}
	var apiPolicy *dpv1alpha2.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	var ratelimitPolicy *dpv1alpha1.RateLimitPolicy
	if outputRatelimitPolicy != nil {
		ratelimitPolicy = *outputRatelimitPolicy
	}

	//We are only supporting one backend for now
	backend := gqlRoute.Spec.BackendRefs[0]
	backendName := types.NamespacedName{
		Name:      string(backend.Name),
		Namespace: utils.GetNamespace(backend.Namespace, gqlRoute.Namespace),
	}
	resolvedBackend, ok := resourceParams.BackendMapping[backendName.String()]
	if ok {
		endpointConfig := &EndpointConfig{}
		if resolvedBackend.CircuitBreaker != nil {
			endpointConfig.CircuitBreakers = &CircuitBreakers{
				MaxConnections:     int32(resolvedBackend.CircuitBreaker.MaxConnections),
				MaxRequests:        int32(resolvedBackend.CircuitBreaker.MaxRequests),
				MaxPendingRequests: int32(resolvedBackend.CircuitBreaker.MaxPendingRequests),
				MaxRetries:         int32(resolvedBackend.CircuitBreaker.MaxRetries),
				MaxConnectionPools: int32(resolvedBackend.CircuitBreaker.MaxConnectionPools),
			}
		}
		if resolvedBackend.Timeout != nil {
			endpointConfig.TimeoutInMillis = resolvedBackend.Timeout.UpstreamResponseTimeout * 1000
			endpointConfig.IdleTimeoutInSeconds = resolvedBackend.Timeout.DownstreamRequestIdleTimeout
		}
		if resolvedBackend.Retry != nil {
			statusCodes := config.Envoy.Upstream.Retry.StatusCodes
			if len(resolvedBackend.Retry.StatusCodes) > 0 {
				statusCodes = resolvedBackend.Retry.StatusCodes
			}
			endpointConfig.RetryConfig = &RetryConfig{
				Count:                int32(resolvedBackend.Retry.Count),
				StatusCodes:          statusCodes,
				BaseIntervalInMillis: int32(resolvedBackend.Retry.BaseIntervalMillis),
			}
		}
		adapterInternalAPI.Endpoints = &EndpointCluster{
			Endpoints: GetEndpoints(backendName, resourceParams.BackendMapping),
			Config:    endpointConfig,
		}
		if resolvedBackend.HealthCheck != nil {
			adapterInternalAPI.Endpoints.HealthCheck = &HealthCheck{
				Interval:           resolvedBackend.HealthCheck.Interval,
				Timeout:            resolvedBackend.HealthCheck.Timeout,
				UnhealthyThreshold: resolvedBackend.HealthCheck.UnhealthyThreshold,
				HealthyThreshold:   resolvedBackend.HealthCheck.HealthyThreshold,
			}
		}

		var securityConfig []EndpointSecurity
		switch resolvedBackend.Security.Type {
		case "Basic":
			securityConfig = append(securityConfig, EndpointSecurity{
				Password: string(resolvedBackend.Security.Basic.Password),
				Username: string(resolvedBackend.Security.Basic.Username),
				Type:     string(resolvedBackend.Security.Type),
				Enabled:  true,
			})
		}
		adapterInternalAPI.EndpointSecurity = utils.GetPtrSlice(securityConfig)
	} else {
		return fmt.Errorf("backend: %s has not been resolved", backendName)
	}

	for _, rule := range gqlRoute.Spec.Rules {
		var policies = OperationPolicies{}
		resourceAuthScheme := authScheme
		resourceRatelimitPolicy := ratelimitPolicy
		var scopes []string

		for _, filter := range rule.Filters {
			if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindAuthentication {
				if ref, found := resourceParams.ResourceAuthSchemes[types.NamespacedName{
					Name:      string(filter.ExtensionRef.Name),
					Namespace: gqlRoute.Namespace,
				}.String()]; found {
					resourceAuthScheme = concatAuthSchemes(authScheme, &ref)
				} else {
					return fmt.Errorf(`auth scheme: %s has not been resolved, spec.targetRef.kind should be 
						'Resource' in resource level Authentications`, filter.ExtensionRef.Name)
				}
			}
			if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindScope {
				if ref, found := resourceParams.ResourceScopes[types.NamespacedName{
					Name:      string(filter.ExtensionRef.Name),
					Namespace: gqlRoute.Namespace,
				}.String()]; found {
					scopes = ref.Spec.Names
					disableScopes = false
				} else {
					return fmt.Errorf("scope: %s has not been resolved in namespace %s", filter.ExtensionRef.Name, gqlRoute.Namespace)
				}
			}
			if filter.ExtensionRef != nil && filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
				if ref, found := resourceParams.ResourceRateLimitPolicies[types.NamespacedName{
					Name:      string(filter.ExtensionRef.Name),
					Namespace: gqlRoute.Namespace,
				}.String()]; found {
					resourceRatelimitPolicy = concatRateLimitPolicies(ratelimitPolicy, &ref)
				} else {
					return fmt.Errorf(`ratelimitpolicy: %s has not been resolved, spec.targetRef.kind should be 
						'Resource' in resource level RateLimitPolicies`, filter.ExtensionRef.Name)
				}
			}
		}
		resourceAuthScheme = concatAuthSchemes(resourceAuthScheme, nil)
		resourceRatelimitPolicy = concatRateLimitPolicies(resourceRatelimitPolicy, nil)

		apiAuth := getSecurity(resourceAuthScheme)

		for _, match := range rule.Matches {
			resourcePath := *match.Path
			resource := &Resource{path: resourcePath,
				methods: []*Operation{{iD: uuid.New().String(), method: string(*match.Type), policies: policies,
					auth: apiAuth, rateLimitPolicy: parseRateLimitPolicyToInternal(resourceRatelimitPolicy), scopes: scopes}},
				iD: uuid.New().String(),
			}
			resources = append(resources, resource)
		}
	}

	ratelimitPolicy = concatRateLimitPolicies(ratelimitPolicy, nil)
	apiPolicy = concatAPIPolicies(apiPolicy, nil)
	authScheme = concatAuthSchemes(authScheme, nil)

	adapterInternalAPI.RateLimitPolicy = parseRateLimitPolicyToInternal(ratelimitPolicy)
	adapterInternalAPI.resources = resources
	adapterInternalAPI.xWso2Cors = getCorsConfigFromAPIPolicy(apiPolicy)
	if authScheme.Spec.Override != nil && authScheme.Spec.Override.Disabled != nil {
		adapterInternalAPI.disableAuthentications = *authScheme.Spec.Override.Disabled
	}
	adapterInternalAPI.disableScopes = disableScopes
	return nil
}

func (endpoint *Endpoint) validateEndpoint() error {
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
			loggers.LoggerOasparser.Errorf("Given status code for the API retry config is invalid." +
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
				loggers.LoggerOasparser.Errorf("Error while parsing the endpoint. %v", err)
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
func (adapterInternalAPI *AdapterInternalAPI) GetOperationInterceptors(apiInterceptor InterceptEndpoint, resourceInterceptor InterceptEndpoint, operations []*Operation, isIn bool) map[string]InterceptEndpoint {
	interceptorOperationMap := make(map[string]InterceptEndpoint)

	for _, op := range operations {
		extensionName := constants.XWso2RequestInterceptor
		// first get operational policies
		operationInterceptor := op.GetCallInterceptorService(isIn)
		// if operational policy interceptor not given check operational level adapterInternalAPI extension
		if !operationInterceptor.Enable {
			if !isIn {
				extensionName = constants.XWso2ResponseInterceptor
			}
			operationInterceptor = adapterInternalAPI.GetInterceptor(op.GetVendorExtensions(), extensionName, constants.OperationLevelInterceptor)
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
func (adapterInternalAPI *AdapterInternalAPI) GetInterceptor(vendorExtensions map[string]interface{}, extensionName string, level string) InterceptEndpoint {
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
					loggers.LoggerOasparser.Error("Error reading interceptors service url value", err)
					return InterceptEndpoint{}
				}
				if endpoint.Basepath != "" {
					loggers.LoggerOasparser.Warnf("Interceptor serviceURL basepath is given as %v but it will be ignored",
						endpoint.Basepath)
				}
				endpointCluster.Endpoints = []Endpoint{*endpoint}

			} else {
				loggers.LoggerOasparser.Error("Error reading interceptors service url value")
				return InterceptEndpoint{}
			}
			//clusterTimeout optional
			if v, found := val[constants.ClusterTimeout]; found {
				p, err := strconv.ParseInt(fmt.Sprint(v), 0, 0)
				if err == nil {
					clusterTimeoutV = time.Duration(p)
				} else {
					loggers.LoggerOasparser.Errorf("Error reading interceptors %v value : %v", constants.ClusterTimeout, err.Error())
				}
			}
			//requestTimeout optional
			if v, found := val[constants.RequestTimeout]; found {
				p, err := strconv.ParseInt(fmt.Sprint(v), 0, 0)
				if err == nil {
					requestTimeoutV = time.Duration(p)
				} else {
					loggers.LoggerOasparser.Errorf("Error reading interceptors %v value : %v", constants.RequestTimeout, err.Error())
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
		loggers.LoggerOasparser.Error("Error parsing response interceptors values to adapterInternalAPI")
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
