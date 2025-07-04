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
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	dpv1alpha5 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha5"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// AdapterInternalAPI represents the object structure holding the information related to the
// adapter internal representation. The values are populated from the operator. The pathItem level information is represented
// by the resources array which contains the Resource entries.
type AdapterInternalAPI struct {
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
	applicationSecurity      map[string]bool
	EnvType                  string
	backendJWTTokenInfo      *BackendJWTTokenInfo
	apiDefinitionFile        []byte
	apiDefinitionEndpoint    string
	subscriptionValidation   bool
	APIProperties            []dpv1alpha3.Property
	// GraphQLSchema              string
	// GraphQLComplexities        GraphQLComplexityYaml
	IsSystemAPI             bool
	RateLimitPolicy         *RateLimitPolicy
	environment             string
	Endpoints               *EndpointCluster
	EndpointSecurity        []*EndpointSecurity
	AIProvider              InternalAIProvider
	AIModelBasedRoundRobin  InternalModelBasedRoundRobin
	HTTPRouteIDs            []string
	RequestInBuiltPolicies  []InternalInBuiltPolicy
	ResponseInBuiltPolicies []InternalInBuiltPolicy
}

// InternalInBuiltPolicy represents the in-built policy configurations
type InternalInBuiltPolicy struct {
	PolicyName    string            `json:"policyName"`
	PolicyID      string            `json:"policyID"`
	PolicyVersion string            `json:"policyVersion"`
	Parameters    map[string]string `json:"parameters,omitempty"`
	PolicyOrder   int               `json:"policyOrder,omitempty"`
}

// InternalModelBasedRoundRobin holds the model based round robin configurations
type InternalModelBasedRoundRobin struct {
	OnQuotaExceedSuspendDuration int                   `json:"onQuotaExceedSuspendDuration,omitempty"`
	ProductionModels             []InternalModelWeight `json:"productionModels"`
	SandboxModels                []InternalModelWeight `json:"sandboxModels"`
}

// InternalModelWeight holds the model configurations
type InternalModelWeight struct {
	Model               string `json:"model"`
	EndpointClusterName string `json:"endpointClusterName"`
	Weight              int    `json:"weight,omitempty"`
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

// InternalAIProvider represents the object structure holding the information related to the AI Provider
type InternalAIProvider struct {
	Enabled            bool
	ProviderName       string
	ProviderAPIVersion string
	Organization       string
	SupportedModels    []string
	RequestModel       ValueDetails
	ResponseModel      ValueDetails
	PromptTokens       ValueDetails
	CompletionToken    ValueDetails
	TotalToken         ValueDetails
}

// ValueDetails defines the value details
type ValueDetails struct {
	In    string `json:"in"`
	Value string `json:"value"`
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
	// Weight assigned for the endpoint (optional)
	Weight int32
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

// SetApplicationSecurity sets the optional or mandatory application security for each security type
// true means mandatory
func (adapterInternalAPI *AdapterInternalAPI) SetApplicationSecurity(key string, value bool) {
	if adapterInternalAPI.applicationSecurity == nil {
		adapterInternalAPI.applicationSecurity = make(map[string]bool)
	}
	adapterInternalAPI.applicationSecurity[key] = value
}

// GetApplicationSecurity returns true if application security is mandatory, and false if optional
func (adapterInternalAPI *AdapterInternalAPI) GetApplicationSecurity() map[string]bool {
	return adapterInternalAPI.applicationSecurity
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

// SetAIProvider sets the AIProvider of the API.
func (adapterInternalAPI *AdapterInternalAPI) SetAIProvider(aiProvider dpv1alpha4.AIProvider) {
	adapterInternalAPI.AIProvider = InternalAIProvider{
		Enabled:            true,
		ProviderName:       aiProvider.Spec.ProviderName,
		ProviderAPIVersion: aiProvider.Spec.ProviderAPIVersion,
		Organization:       aiProvider.Spec.Organization,
		SupportedModels:    aiProvider.Spec.SupportedModels,
		RequestModel: ValueDetails{
			In:    aiProvider.Spec.RequestModel.In,
			Value: aiProvider.Spec.RequestModel.Value,
		},
		ResponseModel: ValueDetails{
			In:    aiProvider.Spec.ResponseModel.In,
			Value: aiProvider.Spec.ResponseModel.Value,
		},
		PromptTokens: ValueDetails{
			In:    aiProvider.Spec.RateLimitFields.PromptTokens.In,
			Value: aiProvider.Spec.RateLimitFields.PromptTokens.Value,
		},
		CompletionToken: ValueDetails{
			In:    aiProvider.Spec.RateLimitFields.CompletionToken.In,
			Value: aiProvider.Spec.RateLimitFields.CompletionToken.Value,
		},
		TotalToken: ValueDetails{
			In:    aiProvider.Spec.RateLimitFields.TotalToken.In,
			Value: aiProvider.Spec.RateLimitFields.TotalToken.Value,
		},
	}
}

// GetAIProvider returns the AIProvider of the API
func (adapterInternalAPI *AdapterInternalAPI) GetAIProvider() InternalAIProvider {
	return adapterInternalAPI.AIProvider
}

// SetModelBasedRoundRobin sets the ModelBasedRoundRobin of the API.
func (adapterInternalAPI *AdapterInternalAPI) SetModelBasedRoundRobin(modelBasedRoundRobin InternalModelBasedRoundRobin) {
	adapterInternalAPI.AIModelBasedRoundRobin = modelBasedRoundRobin
}

// GetModelBasedRoundRobin returns the ModelBasedRoundRobin of the API
func (adapterInternalAPI *AdapterInternalAPI) GetModelBasedRoundRobin() InternalModelBasedRoundRobin {
	return adapterInternalAPI.AIModelBasedRoundRobin
}

// GetRequestInBuiltPolicies returns the in-built policies that are applied to the request of the API.
func (adapterInternalAPI *AdapterInternalAPI) GetRequestInBuiltPolicies() []InternalInBuiltPolicy {
	return adapterInternalAPI.RequestInBuiltPolicies
}

// SetRequestInBuiltPolicies sets the in-built policies that are applied to the request of the API.
func (adapterInternalAPI *AdapterInternalAPI) SetRequestInBuiltPolicies(policies []InternalInBuiltPolicy) {
	if policies == nil {
		adapterInternalAPI.RequestInBuiltPolicies = []InternalInBuiltPolicy{}
	} else {
		adapterInternalAPI.RequestInBuiltPolicies = policies
	}
}

// GetResponseInBuiltPolicies returns the in-built policies that are applied to the response of the API.
func (adapterInternalAPI *AdapterInternalAPI) GetResponseInBuiltPolicies() []InternalInBuiltPolicy {
	return adapterInternalAPI.ResponseInBuiltPolicies
}

// SetResponseInBuiltPolicies sets the in-built policies that are applied to the response of the API.
func (adapterInternalAPI *AdapterInternalAPI) SetResponseInBuiltPolicies(policies []InternalInBuiltPolicy) {
	if policies == nil {
		adapterInternalAPI.ResponseInBuiltPolicies = []InternalInBuiltPolicy{}
	} else {
		adapterInternalAPI.ResponseInBuiltPolicies = policies
	}
}

// Validate method confirms that the adapterInternalAPI has all required fields in the required format.
// This needs to be checked prior to generate router/enforcer related resources.
func (adapterInternalAPI *AdapterInternalAPI) Validate() error {
	for _, res := range adapterInternalAPI.resources {
		if res.endpoints == nil || (len(res.endpoints.Endpoints) == 0 && !res.hasRequestRedirectFilter) {
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
func (adapterInternalAPI *AdapterInternalAPI) SetInfoHTTPRouteCR(httpRoute *gwapiv1.HTTPRoute, resourceParams ResourceParams, ruleIdxToAiRatelimitPolicyMapping map[int]*dpv1alpha3.AIRateLimitPolicy, extractTokenFrom string) error {
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
	var apiPolicy *dpv1alpha5.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	var ratelimitPolicy *dpv1alpha3.RateLimitPolicy
	if outputRatelimitPolicy != nil {
		ratelimitPolicy = *outputRatelimitPolicy
	}

	for ruleID, rule := range httpRoute.Spec.Rules {
		var endPoints []Endpoint
		var policies = OperationPolicies{}
		var circuitBreaker *dpv1alpha2.CircuitBreaker
		var healthCheck *dpv1alpha2.HealthCheck
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
		hasRequestRedirectPolicy := false
		var securityConfig []EndpointSecurity
		var mirrorEndpointClusters []*EndpointCluster

		enableBackendBasedAIRatelimit := false
		descriptorValue := ""
		if aiRatelimitPolicy, exists := ruleIdxToAiRatelimitPolicyMapping[ruleID]; exists {
			loggers.LoggerAPI.Debugf("Found AI ratelimit mapping for ruleId: %d, related api: %s", ruleID, adapterInternalAPI.UUID)
			enableBackendBasedAIRatelimit = true
			descriptorValue = prepareAIRatelimitIdentifier(adapterInternalAPI.OrganizationID, utils.NamespacedName(aiRatelimitPolicy), &aiRatelimitPolicy.Spec)
		} else {
			loggers.LoggerAPI.Debugf("Could not find AIratelimit for ruleId: %d, len of map: %d, related api: %s", ruleID, len(ruleIdxToAiRatelimitPolicyMapping), adapterInternalAPI.UUID)
		}

		backendBasePath := ""
		for _, backend := range rule.BackendRefs {
			backendName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}
			resolvedBackend, ok := resourceParams.BackendMapping[backendName.String()]
			if ok {
				if resolvedBackend.CircuitBreaker != nil {
					circuitBreaker = &dpv1alpha2.CircuitBreaker{
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
					healthCheck = &dpv1alpha2.HealthCheck{
						Interval:           resolvedBackend.HealthCheck.Interval,
						Timeout:            resolvedBackend.HealthCheck.Timeout,
						UnhealthyThreshold: resolvedBackend.HealthCheck.UnhealthyThreshold,
						HealthyThreshold:   resolvedBackend.HealthCheck.HealthyThreshold,
					}
				}
				if backend.Weight != nil {
					// Extracting weights from HTTPRoute if weights are defined
					resolvedBackend.Weight = *backend.Weight
					loggers.LoggerAPI.Debugf("Weighted Routing Capability is enabled for the Resolved Backend %s with weight %d", backendName.String(), resolvedBackend.Weight)
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
				case "APIKey":
					securityConfig = append(securityConfig, EndpointSecurity{
						Type:    string(resolvedBackend.Security.Type),
						Enabled: true,
						CustomParameters: map[string]string{
							"in":    string(resolvedBackend.Security.APIKey.In),
							"key":   string(resolvedBackend.Security.APIKey.Name),
							"value": string(resolvedBackend.Security.APIKey.Value),
						},
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
			case gwapiv1.HTTPRouteFilterRequestHeaderModifier:
				for _, header := range filter.RequestHeaderModifier.Add {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Request = append(policies.Request, Policy{
						PolicyName: string(gwapiv1.HTTPRouteFilterRequestHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.RequestHeaderModifier.Remove {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header)

					policies.Request = append(policies.Request, Policy{
						PolicyName: string(gwapiv1.HTTPRouteFilterRequestHeaderModifier),
						Action:     constants.ActionHeaderRemove,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.RequestHeaderModifier.Set {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Request = append(policies.Request, Policy{
						PolicyName: string(gwapiv1.HTTPRouteFilterRequestHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
			case gwapiv1.HTTPRouteFilterResponseHeaderModifier:
				for _, header := range filter.ResponseHeaderModifier.Add {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Response = append(policies.Response, Policy{
						PolicyName: string(gwapiv1.HTTPRouteFilterResponseHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.ResponseHeaderModifier.Remove {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header)

					policies.Response = append(policies.Response, Policy{
						PolicyName: string(gwapiv1.HTTPRouteFilterResponseHeaderModifier),
						Action:     constants.ActionHeaderRemove,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.ResponseHeaderModifier.Set {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Response = append(policies.Response, Policy{
						PolicyName: string(gwapiv1.HTTPRouteFilterResponseHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
			case gwapiv1.HTTPRouteFilterRequestRedirect:
				var requestRedirectEndpoint Endpoint
				hasRequestRedirectPolicy = true
				policyParameters := make(map[string]interface{})
				scheme := *filter.RequestRedirect.Scheme
				host := string(*filter.RequestRedirect.Hostname)
				port := filter.RequestRedirect.Port
				code := filter.RequestRedirect.StatusCode

				policyParameters[constants.RedirectScheme] = scheme
				requestRedirectEndpoint.URLType = scheme
				policyParameters[constants.RedirectHostname] = host
				requestRedirectEndpoint.Host = host

				if port != nil {
					policyParameters[constants.RedirectPort] = strconv.Itoa(int(*port))
					requestRedirectEndpoint.Port = uint32(*port)
				} else {
					if requestRedirectEndpoint.URLType == "http" {
						requestRedirectEndpoint.Port = 80
					} else if requestRedirectEndpoint.URLType == "https" {
						requestRedirectEndpoint.Port = 443
					}
				}

				if code != nil {
					policyParameters[constants.RedirectStatusCode] = *code
				}

				switch filter.RequestRedirect.Path.Type {
				case gwapiv1.FullPathHTTPPathModifier:
					policyParameters[constants.RedirectPath] = backendBasePath + *filter.RequestRedirect.Path.ReplaceFullPath
				case gwapiv1.PrefixMatchHTTPPathModifier:
					policyParameters[constants.RedirectPath] = backendBasePath + *filter.RequestRedirect.Path.ReplacePrefixMatch
				}
				requestRedirectEndpoint.Basepath = policyParameters[constants.RedirectPath].(string)

				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1.HTTPRouteFilterRequestRedirect),
					Action:     constants.ActionRedirectRequest,
					Parameters: policyParameters,
				})
				endPoints = append(endPoints, requestRedirectEndpoint)

			case gwapiv1.HTTPRouteFilterRequestMirror:
				var mirrorTimeoutInMillis uint32
				var mirrorIdleTimeoutInSeconds uint32
				var mirrorCircuitBreaker *dpv1alpha1.CircuitBreaker
				var mirrorHealthCheck *dpv1alpha1.HealthCheck
				isMirrorRetryConfig := false
				isMirrorRouteTimeout := false
				var mirrorBackendRetryCount uint32
				var mirrorStatusCodes []uint32
				mirrorStatusCodes = append(mirrorStatusCodes, config.Envoy.Upstream.Retry.StatusCodes...)
				var mirrorBaseIntervalInMillis uint32
				policyParameters := make(map[string]interface{})
				mirrorBackend := &filter.RequestMirror.BackendRef
				mirrorBackendName := types.NamespacedName{
					Name:      string(mirrorBackend.Name),
					Namespace: utils.GetNamespace(mirrorBackend.Namespace, httpRoute.Namespace),
				}
				resolvedMirrorBackend, ok := resourceParams.BackendMapping[mirrorBackendName.String()]

				if ok {
					if resolvedMirrorBackend.CircuitBreaker != nil {
						mirrorCircuitBreaker = &dpv1alpha1.CircuitBreaker{
							MaxConnections:     resolvedMirrorBackend.CircuitBreaker.MaxConnections,
							MaxPendingRequests: resolvedMirrorBackend.CircuitBreaker.MaxPendingRequests,
							MaxRequests:        resolvedMirrorBackend.CircuitBreaker.MaxRequests,
							MaxRetries:         resolvedMirrorBackend.CircuitBreaker.MaxRetries,
							MaxConnectionPools: resolvedMirrorBackend.CircuitBreaker.MaxConnectionPools,
						}
					}

					if resolvedMirrorBackend.Timeout != nil {
						isMirrorRouteTimeout = true
						mirrorTimeoutInMillis = resolvedMirrorBackend.Timeout.UpstreamResponseTimeout * 1000
						mirrorIdleTimeoutInSeconds = resolvedMirrorBackend.Timeout.DownstreamRequestIdleTimeout
					}

					if resolvedMirrorBackend.Retry != nil {
						isMirrorRetryConfig = true
						mirrorBackendRetryCount = resolvedMirrorBackend.Retry.Count
						mirrorBaseIntervalInMillis = resolvedMirrorBackend.Retry.BaseIntervalMillis
						if len(resolvedMirrorBackend.Retry.StatusCodes) > 0 {
							mirrorStatusCodes = resolvedMirrorBackend.Retry.StatusCodes
						}
					}

					if resolvedMirrorBackend.HealthCheck != nil {
						mirrorHealthCheck = &dpv1alpha1.HealthCheck{
							Interval:           resolvedMirrorBackend.HealthCheck.Interval,
							Timeout:            resolvedMirrorBackend.HealthCheck.Timeout,
							UnhealthyThreshold: resolvedMirrorBackend.HealthCheck.UnhealthyThreshold,
							HealthyThreshold:   resolvedMirrorBackend.HealthCheck.HealthyThreshold,
						}
					}
				} else {
					return fmt.Errorf("backend: %s has not been resolved", mirrorBackendName)
				}

				mirrorEndpoints := GetEndpoints(mirrorBackendName, resourceParams.BackendMapping)
				if len(mirrorEndpoints) > 0 {
					mirrorEndpointCluster := &EndpointCluster{
						Endpoints: mirrorEndpoints,
					}
					mirrorEndpointConfig := &EndpointConfig{}
					if isMirrorRouteTimeout {
						mirrorEndpointConfig.TimeoutInMillis = mirrorTimeoutInMillis
						mirrorEndpointConfig.IdleTimeoutInSeconds = mirrorIdleTimeoutInSeconds
					}
					if mirrorCircuitBreaker != nil {
						mirrorEndpointConfig.CircuitBreakers = &CircuitBreakers{
							MaxConnections:     int32(mirrorCircuitBreaker.MaxConnections),
							MaxRequests:        int32(mirrorCircuitBreaker.MaxRequests),
							MaxPendingRequests: int32(mirrorCircuitBreaker.MaxPendingRequests),
							MaxRetries:         int32(mirrorCircuitBreaker.MaxRetries),
							MaxConnectionPools: int32(mirrorCircuitBreaker.MaxConnectionPools),
						}
					}
					if isMirrorRetryConfig {
						mirrorEndpointConfig.RetryConfig = &RetryConfig{
							Count:                int32(mirrorBackendRetryCount),
							StatusCodes:          mirrorStatusCodes,
							BaseIntervalInMillis: int32(mirrorBaseIntervalInMillis),
						}
					}
					if mirrorHealthCheck != nil {
						mirrorEndpointCluster.HealthCheck = &HealthCheck{
							Interval:           mirrorHealthCheck.Interval,
							Timeout:            mirrorHealthCheck.Timeout,
							UnhealthyThreshold: mirrorHealthCheck.UnhealthyThreshold,
							HealthyThreshold:   mirrorHealthCheck.HealthyThreshold,
						}
					}
					if isMirrorRouteTimeout || mirrorCircuitBreaker != nil || mirrorHealthCheck != nil || isMirrorRetryConfig {
						mirrorEndpointCluster.Config = mirrorEndpointConfig
					}
					mirrorEndpointClusters = append(mirrorEndpointClusters, mirrorEndpointCluster)
				}
				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1.HTTPRouteFilterRequestMirror),
					Action:     constants.ActionMirrorRequest,
					Parameters: policyParameters,
				})
			}
		}
		resourceAPIPolicy = concatAPIPolicies(resourceAPIPolicy, nil)
		resourceAuthScheme = concatAuthSchemes(resourceAuthScheme, nil)
		resourceRatelimitPolicy = concatRateLimitPolicies(resourceRatelimitPolicy, nil)
		addOperationLevelInterceptors(&policies, resourceAPIPolicy, resourceParams.InterceptorServiceMapping, resourceParams.BackendMapping, httpRoute.Namespace)

		loggers.LoggerOasparser.Debugf("Calculating auths for API ..., API_UUID = %v", adapterInternalAPI.UUID)
		apiAuth := getSecurity(resourceAuthScheme)

		if !hasRequestRedirectPolicy && len(rule.BackendRefs) < 1 {
			return fmt.Errorf("no backendref were provided")
		}

		for matchID, match := range rule.Matches {
			if hasURLRewritePolicy && hasRequestRedirectPolicy {
				return fmt.Errorf("cannot have URL Rewrite and Request Redirect under the same rule")
			}
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
			matchID := getMatchID(httpRoute.Namespace, httpRoute.Name, ruleID, matchID)
			operations := getAllowedOperations(matchID, match.Method, policies, apiAuth,
				parseRateLimitPolicyToInternal(resourceRatelimitPolicy), scopes, mirrorEndpointClusters)

			vhost := ""
			for _, hostName := range httpRoute.Spec.Hostnames {
				vhost = string(hostName)
			}
			var modelBasedRoundRobin *InternalModelBasedRoundRobin
			if extracted := extractModelBasedRoundRobinFromPolicy(resourceAPIPolicy, resourceParams.BackendMapping, adapterInternalAPI, resourcePath, vhost, httpRoute.Namespace); extracted != nil {
				loggers.LoggerAPI.Debugf("ModelBasedRoundRobin extracted %v", extracted)
				modelBasedRoundRobin = extracted
			}

			var requestInBuiltPolicies []*InternalInBuiltPolicy
			var responseInBuiltPolicies []*InternalInBuiltPolicy
			if extracted := extractRequestInBuiltPolicies(resourceAPIPolicy); extracted != nil {
				loggers.LoggerAPI.Debugf("Request In-Built Policies extracted %v", extracted)
				requestInBuiltPolicies = extracted
			}

			if extracted := extractResponseInBuiltPolicies(resourceAPIPolicy); extracted != nil {
				loggers.LoggerAPI.Debugf("Response In-Built Policies extracted %v", extracted)
				responseInBuiltPolicies = extracted
			}

			resource := &Resource{
				path:                                   resourcePath,
				methods:                                operations,
				pathMatchType:                          *match.Path.Type,
				hasPolicies:                            true,
				iD:                                     uuid.New().String(),
				hasRequestRedirectFilter:               hasRequestRedirectPolicy,
				enableBackendBasedAIRatelimit:          enableBackendBasedAIRatelimit,
				backendBasedAIRatelimitDescriptorValue: descriptorValue,
				extractTokenFrom:                       extractTokenFrom,
				AIModelBasedRoundRobin:                 modelBasedRoundRobin,
				RequestInBuiltPolicies:                 requestInBuiltPolicies,
				ResponseInBuiltPolicies:                responseInBuiltPolicies,
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
	if authSpec != nil && authSpec.AuthTypes != nil {
		var required bool
		var oauth2Enabled bool
		var apiKeyEnabled bool
		var jwtEnabled bool
		if authSpec.AuthTypes.OAuth2.Required != "" && authSpec.AuthTypes.OAuth2.Disabled != true {
			oauth2Enabled = true
			required = required || authSpec.AuthTypes.OAuth2.Required == "mandatory"
		}

		if authSpec.AuthTypes.APIKey != nil {
			apiKeyEnabled = true
			required = required || authSpec.AuthTypes.APIKey.Required == "mandatory"
		}
		if !*authSpec.AuthTypes.JWT.Disabled {
			jwtEnabled = true
			required = required || false
		}
		if jwtEnabled {
			adapterInternalAPI.SetApplicationSecurity(constants.JWT, false)
		}
		if oauth2Enabled {
			adapterInternalAPI.SetApplicationSecurity(constants.OAuth2, false)
		}
		if apiKeyEnabled {
			adapterInternalAPI.SetApplicationSecurity(constants.APIKey, false)
		}
	} else {
		adapterInternalAPI.SetApplicationSecurity(constants.OAuth2, true)
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

// ExtractModelBasedRoundRobinFromPolicy extracts the ModelBasedRoundRobin from the API Policy
func extractModelBasedRoundRobinFromPolicy(apiPolicy *dpv1alpha5.APIPolicy, backendMapping map[string]*dpv1alpha4.ResolvedBackend, adapterInternalAPI *AdapterInternalAPI, resourcePath string, vHost string, namespace string) *InternalModelBasedRoundRobin {
	if apiPolicy == nil {
		return nil
	}
	resolvedModelBasedRoundRobin := &InternalModelBasedRoundRobin{}
	loggers.LoggerAPI.Debugf("Extracting ModelBasedRoundRobin from API Policy %v", apiPolicy)
	loggers.LoggerAPI.Debugf("Backend Mapping %v", backendMapping)
	loggers.LoggerAPI.Debugf("ResourcePath %v", resourcePath)

	// Safely access Override section
	if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.ModelBasedRoundRobin != nil {
		loggers.LoggerAPI.Debugf("ModelBasedRoundRobin Override section  %v", apiPolicy.Spec.Override.ModelBasedRoundRobin)
		modelBasedRoundRobin := apiPolicy.Spec.Override.ModelBasedRoundRobin
		resolvedModelBasedRoundRobin = &InternalModelBasedRoundRobin{
			OnQuotaExceedSuspendDuration: modelBasedRoundRobin.OnQuotaExceedSuspendDuration,
		}
		if modelBasedRoundRobin.ProductionModels != nil {
			productionModels := apiPolicy.Spec.Override.ModelBasedRoundRobin.ProductionModels
			for _, model := range productionModels {
				if model.BackendRef.Name != "" {
					if namespace == "" {
						namespace = "default"
					} else if apiPolicy.Namespace != "" {
						namespace = apiPolicy.Namespace
					}
					backendNamespacedName := types.NamespacedName{
						Name:      string(model.BackendRef.Name),
						Namespace: utils.GetNamespace(model.BackendRef.Namespace, namespace),
					}
					loggers.LoggerAPI.Debugf("Backend NamespacedKey %v", backendNamespacedName.String())
					if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
						loggers.LoggerAPI.Debugf("Backend not found %v", backendNamespacedName)
						continue
					}
					endpoints := GetEndpoints(backendNamespacedName, backendMapping)

					clusternName := getClusterName("", adapterInternalAPI.GetOrganizationID(), vHost, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), endpoints[0].Host+endpoints[0].Basepath)

					resolvedModelWeight := InternalModelWeight{
						Model:               model.Model,
						Weight:              model.Weight,
						EndpointClusterName: clusternName,
					}
					resolvedModelBasedRoundRobin.ProductionModels = append(resolvedModelBasedRoundRobin.ProductionModels, resolvedModelWeight)
				}
			}
		}
		if modelBasedRoundRobin.SandboxModels != nil {
			sandboxModels := apiPolicy.Spec.Override.ModelBasedRoundRobin.SandboxModels
			for _, model := range sandboxModels {
				if model.BackendRef.Name != "" {
					if namespace == "" {
						namespace = "default"
					} else if apiPolicy.Namespace != "" {
						namespace = apiPolicy.Namespace
					}
					backendNamespacedName := types.NamespacedName{
						Name:      string(model.BackendRef.Name),
						Namespace: utils.GetNamespace(model.BackendRef.Namespace, namespace),
					}
					loggers.LoggerAPI.Debugf("Backend NamespacedKey %v", backendNamespacedName.String())
					if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
						loggers.LoggerAPI.Debugf("Backend not found %v", backendNamespacedName)
						continue
					}
					endpoints := GetEndpoints(backendNamespacedName, backendMapping)

					clusternName := getClusterName("", adapterInternalAPI.GetOrganizationID(), vHost, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), endpoints[0].Host+endpoints[0].Basepath)

					resolvedModelWeight := InternalModelWeight{
						Model:               model.Model,
						Weight:              model.Weight,
						EndpointClusterName: clusternName,
					}
					resolvedModelBasedRoundRobin.SandboxModels = append(resolvedModelBasedRoundRobin.SandboxModels, resolvedModelWeight)
				}
			}
		}
		return resolvedModelBasedRoundRobin
	}
	// Safely access Default section
	if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.ModelBasedRoundRobin != nil {
		loggers.LoggerAPI.Debugf("ModelBasedRoundRobin Default section  %v", apiPolicy.Spec.Default.ModelBasedRoundRobin)
		modelBasedRoundRobin := apiPolicy.Spec.Default.ModelBasedRoundRobin
		resolvedModelBasedRoundRobin = &InternalModelBasedRoundRobin{
			OnQuotaExceedSuspendDuration: modelBasedRoundRobin.OnQuotaExceedSuspendDuration,
		}
		if modelBasedRoundRobin.ProductionModels != nil {
			loggers.LoggerAPI.Debugf("ModelBasedRoundRobin Default section ProductionModels %v", modelBasedRoundRobin.ProductionModels)
			productionModels := apiPolicy.Spec.Default.ModelBasedRoundRobin.ProductionModels
			for _, model := range productionModels {
				if model.BackendRef.Name != "" {
					if namespace == "" {
						namespace = "default"
					} else if apiPolicy.Namespace != "" {
						namespace = apiPolicy.Namespace
					}
					backendNamespacedName := types.NamespacedName{
						Name:      string(model.BackendRef.Name),
						Namespace: utils.GetNamespace(model.BackendRef.Namespace, namespace),
					}
					if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
						loggers.LoggerAPI.Debugf("Backend not found %v", backendNamespacedName)
						continue
					}
					endpoints := GetEndpoints(backendNamespacedName, backendMapping)

					clusternName := getClusterName("", adapterInternalAPI.GetOrganizationID(), vHost, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), endpoints[0].Host+endpoints[0].Basepath)

					resolvedModelWeight := InternalModelWeight{
						Model:               model.Model,
						Weight:              model.Weight,
						EndpointClusterName: clusternName,
					}
					resolvedModelBasedRoundRobin.ProductionModels = append(resolvedModelBasedRoundRobin.ProductionModels, resolvedModelWeight)
				}
			}
		}
		if modelBasedRoundRobin.SandboxModels != nil {
			loggers.LoggerAPI.Debugf("ModelBasedRoundRobin Default section SandboxModels %v", modelBasedRoundRobin.SandboxModels)
			sandboxModels := apiPolicy.Spec.Default.ModelBasedRoundRobin.SandboxModels
			for _, model := range sandboxModels {
				if model.BackendRef.Name != "" {
					if namespace == "" {
						namespace = "default"
					} else if apiPolicy.Namespace != "" {
						namespace = apiPolicy.Namespace
					}
					backendNamespacedName := types.NamespacedName{
						Name:      string(model.BackendRef.Name),
						Namespace: utils.GetNamespace(model.BackendRef.Namespace, namespace),
					}
					if _, exists := backendMapping[backendNamespacedName.String()]; !exists {
						loggers.LoggerAPI.Debugf("Backend not found %v", backendNamespacedName)
						continue
					}
					endpoints := GetEndpoints(backendNamespacedName, backendMapping)

					clusternName := getClusterName("", adapterInternalAPI.GetOrganizationID(), vHost, adapterInternalAPI.GetTitle(), adapterInternalAPI.GetVersion(), endpoints[0].Host+endpoints[0].Basepath)

					resolvedModelWeight := InternalModelWeight{
						Model:               model.Model,
						Weight:              model.Weight,
						EndpointClusterName: clusternName,
					}
					resolvedModelBasedRoundRobin.SandboxModels = append(resolvedModelBasedRoundRobin.SandboxModels, resolvedModelWeight)
				}
			}
		}
		return resolvedModelBasedRoundRobin
	}

	loggers.LoggerAPI.Debugf("ModelBasedRoundRobin not found in API Policy %v", apiPolicy)
	// Return nil if nothing matches
	return nil
}

// extractRequestInBuiltPolicies extracts the request in-built policies from the API Policy
func extractRequestInBuiltPolicies(apiPolicy *dpv1alpha5.APIPolicy) []*InternalInBuiltPolicy {
	if apiPolicy == nil {
		return nil
	}
	resolvedRequestInBuiltPolicies := []*InternalInBuiltPolicy{}
	// Safely access Override section
	if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.RequestPolicies != nil && len(apiPolicy.Spec.Override.RequestPolicies) > 0 {
		index := 0
		loggers.LoggerAPI.Debugf("RequestPolicies Override section  %v", apiPolicy.Spec.Override.RequestPolicies)
		for _, policy := range apiPolicy.Spec.Override.RequestPolicies {
			resolvedParameters, err := getResolvedPolicyParameters(policy, apiPolicy.Namespace)
			if err != nil {
				loggers.LoggerAPI.Errorf("Error resolving parameters for policy %s: %v", policy.PolicyName, err)
				continue
			}
			resolvedRequestInBuiltPolicies = append(resolvedRequestInBuiltPolicies, &InternalInBuiltPolicy{
				PolicyName:    policy.PolicyName,
				PolicyID:      policy.PolicyID,
				PolicyVersion: policy.PolicyVersion,
				Parameters:    resolvedParameters,
				PolicyOrder:   index,
			})
			index++
		}
	}

	// Safely access Default section
	if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.RequestPolicies != nil && len(apiPolicy.Spec.Default.RequestPolicies) > 0 {
		loggers.LoggerAPI.Debugf("RequestPolicies Default section  %v", apiPolicy.Spec.Default.RequestPolicies)
		index := 0
		for _, policy := range apiPolicy.Spec.Default.RequestPolicies {
			resolvedParameters, err := getResolvedPolicyParameters(policy, apiPolicy.Namespace)
			if err != nil {
				loggers.LoggerAPI.Errorf("Error resolving parameters for policy %s: %v", policy.PolicyName, err)
				continue
			}
			resolvedRequestInBuiltPolicies = append(resolvedRequestInBuiltPolicies, &InternalInBuiltPolicy{
				PolicyName:    policy.PolicyName,
				PolicyID:      policy.PolicyID,
				PolicyVersion: policy.PolicyVersion,
				Parameters:    resolvedParameters,
				PolicyOrder:   index,
			})
			index++
		}
	}
	if len(resolvedRequestInBuiltPolicies) > 0 {
		loggers.LoggerAPI.Debugf("RequestPolicies found in API Policy %v", apiPolicy.Name)
		return resolvedRequestInBuiltPolicies
	}
	loggers.LoggerAPI.Debugf("RequestPolicies not found in API Policy %v", apiPolicy)
	// Return nil if nothing matches
	return nil
}

// extractResponseInBuiltPolicies extracts the response in-built policies from the API Policy
func extractResponseInBuiltPolicies(apiPolicy *dpv1alpha5.APIPolicy) []*InternalInBuiltPolicy {
	if apiPolicy == nil {
		return nil
	}

	resolvedResponseInBuiltPolicies := []*InternalInBuiltPolicy{}

	if apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.ResponsePolicies != nil && len(apiPolicy.Spec.Override.ResponsePolicies) > 0 {
		loggers.LoggerAPI.Debugf("ResponsePolicies Override section  %v", apiPolicy.Spec.Override.ResponsePolicies)
		index := 0
		for _, policy := range apiPolicy.Spec.Override.ResponsePolicies {
			resolvedParameters, err := getResolvedPolicyParameters(policy, apiPolicy.Namespace)
			if err != nil {
				loggers.LoggerAPI.Errorf("Error resolving parameters for policy %s: %v", policy.PolicyName, err)
				continue
			}
			resolvedResponseInBuiltPolicies = append(resolvedResponseInBuiltPolicies, &InternalInBuiltPolicy{
				PolicyName:    policy.PolicyName,
				PolicyID:      policy.PolicyID,
				PolicyVersion: policy.PolicyVersion,
				Parameters:    resolvedParameters,
				PolicyOrder:   index,
			})
			index++
		}
	}

	// Safely access Default section
	if apiPolicy.Spec.Default != nil && apiPolicy.Spec.Default.ResponsePolicies != nil && len(apiPolicy.Spec.Default.ResponsePolicies) > 0 {
		loggers.LoggerAPI.Debugf("ResponsePolicies Default section  %v", apiPolicy.Spec.Default.ResponsePolicies)
		index := 0
		for _, policy := range apiPolicy.Spec.Default.ResponsePolicies {
			resolvedParameters, err := getResolvedPolicyParameters(policy, apiPolicy.Namespace)
			if err != nil {
				loggers.LoggerAPI.Errorf("Error resolving parameters for policy %s: %v", policy.PolicyName, err)
				continue
			}
			resolvedResponseInBuiltPolicies = append(resolvedResponseInBuiltPolicies, &InternalInBuiltPolicy{
				PolicyName:    policy.PolicyName,
				PolicyID:      policy.PolicyID,
				PolicyVersion: policy.PolicyVersion,
				Parameters:    resolvedParameters,
				PolicyOrder:   index,
			})
			index++
		}
	}

	if len(resolvedResponseInBuiltPolicies) > 0 {
		loggers.LoggerAPI.Debugf("ResponsePolicies found in API Policy %v", apiPolicy.Name)
		return resolvedResponseInBuiltPolicies
	}
	loggers.LoggerAPI.Debugf("ResponsePolicies not found in API Policy %v", apiPolicy)
	return nil
}

// getClusterName returns the cluster name for the API.
func getClusterName(epPrefix string, organizationID string, vHost string, swaggerTitle string, swaggerVersion string,
	hostname string) string {
	if hostname != "" {
		return strings.TrimSpace(organizationID+"_"+epPrefix+"_"+vHost+"_"+strings.Replace(swaggerTitle, " ", "", -1)+swaggerVersion) +
			"_" + strings.Replace(hostname, " ", "", -1) + "0"
	}
	return strings.TrimSpace(organizationID + "_" + epPrefix + "_" + vHost + "_" + strings.Replace(swaggerTitle, " ", "", -1) +
		swaggerVersion)
}

// getResolvedPolicyParameters resolves the policy parameters of policies
func getResolvedPolicyParameters(policy dpv1alpha5.Policy, namespace string) (map[string]string, error) {
	resolvedParams := make(map[string]string)

	for _, paramValue := range policy.Parameters {
		value, err := getResolvedParameterValue(namespace, paramValue)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parameter %s: %w", paramValue.Key, err)
		}
		resolvedParams[paramValue.Key] = value
	}

	return resolvedParams, nil
}

// ResolveParameterValue resolves a ParameterValue to its actual string value
func getResolvedParameterValue(namespace string, paramValue dpv1alpha5.Parameter) (string, error) {
	// If it's a direct value, return it
	if paramValue.Value != nil {
		return *paramValue.Value, nil
	}

	// If it's a reference, resolve it
	if paramValue.ValueRef != nil {
		return "", fmt.Errorf("ValueRef is not supported yet")
	}

	return "", nil
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
	var apiPolicy *dpv1alpha5.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	var ratelimitPolicy *dpv1alpha3.RateLimitPolicy
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
	authSpec := utils.SelectPolicy(&authScheme.Spec.Override, &authScheme.Spec.Default, nil, nil)
	if authSpec != nil && authSpec.AuthTypes != nil {
		var required bool
		var oauth2Enabled bool
		var apiKeyEnabled bool
		var jwtEnabled bool
		if authSpec.AuthTypes.OAuth2.Required != "" {
			oauth2Enabled = true
			required = required || authSpec.AuthTypes.OAuth2.Required == "mandatory"
		}
		if authSpec.AuthTypes.APIKey != nil {
			apiKeyEnabled = true
			required = required || authSpec.AuthTypes.APIKey.Required == "mandatory"
		}
		if !*authSpec.AuthTypes.JWT.Disabled {
			jwtEnabled = true
			required = required || false
		}

		if jwtEnabled {
			adapterInternalAPI.SetApplicationSecurity(constants.JWT, false)
		}
		if oauth2Enabled {
			adapterInternalAPI.SetApplicationSecurity(constants.OAuth2, false)
		}
		if apiKeyEnabled {
			adapterInternalAPI.SetApplicationSecurity(constants.APIKey, false)
		}
	} else {
		adapterInternalAPI.SetApplicationSecurity(constants.OAuth2, true)
	}
	adapterInternalAPI.disableScopes = disableScopes
	return nil
}

// SetInfoGRPCRouteCR populates resources and endpoints of adapterInternalAPI. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (adapterInternalAPI *AdapterInternalAPI) SetInfoGRPCRouteCR(grpcRoute *gwapiv1.GRPCRoute, resourceParams ResourceParams) error {
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
	var apiPolicy *dpv1alpha5.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	var ratelimitPolicy *dpv1alpha3.RateLimitPolicy
	if outputRatelimitPolicy != nil {
		ratelimitPolicy = *outputRatelimitPolicy
	}

	//We are only supporting one backend for now
	backend := grpcRoute.Spec.Rules[0].BackendRefs[0]
	backendName := types.NamespacedName{
		Name:      string(backend.Name),
		Namespace: utils.GetNamespace(backend.Namespace, grpcRoute.Namespace),
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

	for _, rule := range grpcRoute.Spec.Rules {
		var policies = OperationPolicies{}
		var endPoints []Endpoint
		resourceAuthScheme := authScheme
		resourceAPIPolicy := apiPolicy
		resourceRatelimitPolicy := ratelimitPolicy
		var scopes []string
		for _, filter := range rule.Filters {
			switch filter.Type {
			case gwapiv1.GRPCRouteFilterExtensionRef:
				if filter.ExtensionRef.Kind == constants.KindAuthentication {
					if ref, found := resourceParams.ResourceAuthSchemes[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: grpcRoute.Namespace,
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
						Namespace: grpcRoute.Namespace,
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
						Namespace: grpcRoute.Namespace,
					}.String()]; found {
						scopes = ref.Spec.Names
						disableScopes = false
					} else {
						return fmt.Errorf("scope: %s has not been resolved in namespace %s", filter.ExtensionRef.Name, grpcRoute.Namespace)
					}
				}
				if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
					if ref, found := resourceParams.ResourceRateLimitPolicies[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: grpcRoute.Namespace,
					}.String()]; found {
						resourceRatelimitPolicy = concatRateLimitPolicies(ratelimitPolicy, &ref)
					} else {
						return fmt.Errorf(`ratelimitpolicy: %s has not been resolved, spec.targetRef.kind should be 
					 'Resource' in resource level RateLimitPolicies`, filter.ExtensionRef.Name)
					}
				}
			}
		}

		resourceAPIPolicy = concatAPIPolicies(resourceAPIPolicy, nil)
		resourceAuthScheme = concatAuthSchemes(resourceAuthScheme, nil)
		resourceRatelimitPolicy = concatRateLimitPolicies(resourceRatelimitPolicy, nil)
		addOperationLevelInterceptors(&policies, resourceAPIPolicy, resourceParams.InterceptorServiceMapping, resourceParams.BackendMapping, grpcRoute.Namespace)

		loggers.LoggerOasparser.Debugf("Calculating auths for API ..., API_UUID = %v", adapterInternalAPI.UUID)
		apiAuth := getSecurity(resourceAuthScheme)

		for _, match := range rule.Matches {
			resourcePath := adapterInternalAPI.GetXWso2Basepath() + "." + *match.Method.Service + "/" + *match.Method.Method
			endPoints = append(endPoints, GetEndpoints(backendName, resourceParams.BackendMapping)...)
			resource := &Resource{path: resourcePath, pathMatchType: "Exact",
				methods: []*Operation{{iD: uuid.New().String(), method: "POST", policies: policies,
					auth: apiAuth, rateLimitPolicy: parseRateLimitPolicyToInternal(resourceRatelimitPolicy), scopes: scopes}},
				iD: uuid.New().String(),
			}
			endpoints := GetEndpoints(backendName, resourceParams.BackendMapping)
			resource.endpoints = &EndpointCluster{
				Endpoints: endpoints,
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
	authSpec := utils.SelectPolicy(&authScheme.Spec.Override, &authScheme.Spec.Default, nil, nil)
	if authSpec != nil && authSpec.AuthTypes != nil {
		var required bool
		var oauth2Enabled bool
		var apiKeyEnabled bool
		var jwtEnabled bool
		if authSpec.AuthTypes.OAuth2.Required != "" {
			oauth2Enabled = true
			required = required || authSpec.AuthTypes.OAuth2.Required == "mandatory"
		}

		if authSpec.AuthTypes.APIKey != nil {
			apiKeyEnabled = true
			required = required || authSpec.AuthTypes.APIKey.Required == "mandatory"
		}
		if !*authSpec.AuthTypes.JWT.Disabled {
			jwtEnabled = true
			required = required || false
		}
		if jwtEnabled {
			adapterInternalAPI.SetApplicationSecurity(constants.JWT, false)
		}
		if oauth2Enabled {
			adapterInternalAPI.SetApplicationSecurity(constants.OAuth2, false)
		}
		if apiKeyEnabled {
			adapterInternalAPI.SetApplicationSecurity(constants.APIKey, false)
		}
	} else {
		adapterInternalAPI.SetApplicationSecurity(constants.OAuth2, true)
	}
	adapterInternalAPI.disableScopes = disableScopes
	// Check whether the API has a backend JWT token
	if apiPolicy != nil && apiPolicy.Spec.Override != nil && apiPolicy.Spec.Override.BackendJWTPolicy != nil {
		backendJWTPolicy := resourceParams.BackendJWTMapping[types.NamespacedName{
			Name:      apiPolicy.Spec.Override.BackendJWTPolicy.Name,
			Namespace: grpcRoute.Namespace,
		}.String()].Spec
		adapterInternalAPI.backendJWTTokenInfo = parseBackendJWTTokenToInternal(backendJWTPolicy)
	}
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

func prepareAIRatelimitIdentifier(org string, namespacedName types.NamespacedName, spec *dpv1alpha3.AIRateLimitPolicySpec) string {
	targetNamespace := string(namespacedName.Namespace)
	if spec.TargetRef.Namespace != nil && string(*spec.TargetRef.Namespace) != "" {
		targetNamespace = string(*spec.TargetRef.Namespace)
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", org, string(namespacedName.Namespace), string(namespacedName.Name), targetNamespace, string(spec.TargetRef.Name))
}
