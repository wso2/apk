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

package model

import "github.com/wso2/apk/config-deployer-service-go/internal/constants"

// APKConf represents the APK configuration for a given API
type APKConf struct {
	ID                     string                        `json:"id" yaml:"id"`
	Name                   string                        `json:"name" yaml:"name" validate:"required,min=1,max=60"`
	BasePath               string                        `json:"basePath" yaml:"basePath" validate:"required,min=1,max=256"`
	Version                string                        `json:"version" yaml:"version" validate:"required,min=1,max=30"`
	Type                   string                        `json:"type" yaml:"type"`
	DefinitionPath         *string                       `json:"definitionPath,omitempty" yaml:"definitionPath,omitempty"`
	DefaultVersion         bool                          `json:"defaultVersion" yaml:"defaultVersion"`
	SubscriptionValidation bool                          `json:"subscriptionValidation" yaml:"subscriptionValidation"`
	Environment            *string                       `json:"environment,omitempty" yaml:"environment,omitempty"`
	EndpointConfigurations *EndpointConfigurations       `json:"endpointConfigurations,omitempty" yaml:"endpointConfigurations,omitempty"`
	AIProvider             *AIProvider                   `json:"aiProvider,omitempty" yaml:"aiProvider,omitempty"`
	Operations             []APKOperations               `json:"operations,omitempty" yaml:"operations,omitempty"`
	APIPolicies            *APIOperationPolicies         `json:"apiPolicies,omitempty" yaml:"apiPolicies,omitempty"`
	RateLimit              *RateLimit                    `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	Authentication         []AuthenticationRequest       `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	AdditionalProperties   []APKConfAdditionalProperties `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	CorsConfiguration      *CORSConfiguration            `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
	KeyManagers            []KeyManager                  `json:"keyManagers,omitempty" yaml:"keyManagers,omitempty"`
}

// NewAPKConf creates a new APKConf with default values
func NewAPKConf() *APKConf {
	return &APKConf{
		Type:                   constants.API_TYPE_REST,
		DefaultVersion:         false,
		SubscriptionValidation: false,
	}
}

// AIProvider represents configuration for an AI provider
type AIProvider struct {
	Name       string `json:"name" yaml:"name"`
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
}

// NewAIProvider creates a new AIProvider
func NewAIProvider(name, apiVersion string) *AIProvider {
	return &AIProvider{
		Name:       name,
		APIVersion: apiVersion,
	}
}

// EndpointConfigurations represents configuration for production and sandbox endpoints
type EndpointConfigurations struct {
	Production []EndpointConfiguration `json:"production,omitempty" yaml:"production,omitempty"`
	Sandbox    []EndpointConfiguration `json:"sandbox,omitempty" yaml:"sandbox,omitempty"`
}

// NewEndpointConfigurations creates a new EndpointConfigurations
func NewEndpointConfigurations() *EndpointConfigurations {
	return &EndpointConfigurations{}
}

// EndpointConfiguration represents configuration for production and sandbox endpoints
type EndpointConfiguration struct {
	Endpoint         interface{}       `json:"endpoint" yaml:"endpoint"` // can be string or K8sService
	EndpointSecurity *EndpointSecurity `json:"endpointSecurity,omitempty" yaml:"endpointSecurity,omitempty"`
	Certificate      *Certificate      `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	Resiliency       *Resiliency       `json:"resiliency,omitempty" yaml:"resiliency,omitempty"`
	AIRatelimit      *AIRatelimit      `json:"aiRatelimit,omitempty" yaml:"aiRatelimit,omitempty"`
	Weight           *int              `json:"weight,omitempty" yaml:"weight,omitempty"`
}

// NewEndpointConfiguration creates a new EndpointConfiguration
func NewEndpointConfiguration(endpoint interface{}) *EndpointConfiguration {
	return &EndpointConfiguration{
		Endpoint: endpoint,
	}
}

// K8sService represents configuration for a K8s Service
type K8sService struct {
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace *string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Port      *int    `json:"port,omitempty" yaml:"port,omitempty"`
	Protocol  *string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// NewK8sService creates a new K8sService
func NewK8sService() *K8sService {
	return &K8sService{}
}

// EndpointSecurity represents configuration for Endpoint Security
type EndpointSecurity struct {
	Enabled      *bool       `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	SecurityType interface{} `json:"securityType,omitempty" yaml:"securityType,omitempty"` // BasicEndpointSecurity or APIKeyEndpointSecurity
}

// NewEndpointSecurity creates a new EndpointSecurity
func NewEndpointSecurity() *EndpointSecurity {
	return &EndpointSecurity{}
}

// BasicEndpointSecurity represents configuration for Basic Endpoint Security
type BasicEndpointSecurity struct {
	SecretName  string `json:"secretName" yaml:"secretName"`
	UserNameKey string `json:"userNameKey" yaml:"userNameKey"`
	PasswordKey string `json:"passwordKey" yaml:"passwordKey"`
}

// NewBasicEndpointSecurity creates a new BasicEndpointSecurity
func NewBasicEndpointSecurity(secretName, userNameKey, passwordKey string) *BasicEndpointSecurity {
	return &BasicEndpointSecurity{
		SecretName:  secretName,
		UserNameKey: userNameKey,
		PasswordKey: passwordKey,
	}
}

// APIKeyEndpointSecurity represents configuration for API Key Endpoint Security
type APIKeyEndpointSecurity struct {
	SecretName     string `json:"secretName" yaml:"secretName"`
	In             string `json:"in" yaml:"in"`
	APIKeyNameKey  string `json:"apiKeyNameKey" yaml:"apiKeyNameKey"`
	APIKeyValueKey string `json:"apiKeyValueKey" yaml:"apiKeyValueKey"`
}

// NewAPIKeyEndpointSecurity creates a new APIKeyEndpointSecurity
func NewAPIKeyEndpointSecurity(secretName, in, apiKeyNameKey, apiKeyValueKey string) *APIKeyEndpointSecurity {
	return &APIKeyEndpointSecurity{
		SecretName:     secretName,
		In:             in,
		APIKeyNameKey:  apiKeyNameKey,
		APIKeyValueKey: apiKeyValueKey,
	}
}

// Certificate represents configuration for K8s Secret
type Certificate struct {
	SecretName *string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	SecretKey  *string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
}

// NewCertificate creates a new Certificate
func NewCertificate() *Certificate {
	return &Certificate{}
}

// Resiliency represents configuration of Resiliency settings
type Resiliency struct {
	CircuitBreaker *CircuitBreaker `json:"circuitBreaker,omitempty" yaml:"circuitBreaker,omitempty"`
	Timeout        *Timeout        `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	RetryPolicy    *RetryPolicy    `json:"retryPolicy,omitempty" yaml:"retryPolicy,omitempty"`
}

// NewResiliency creates a new Resiliency
func NewResiliency() *Resiliency {
	return &Resiliency{}
}

// CircuitBreaker represents configuration of CircuitBreaker settings
type CircuitBreaker struct {
	MaxConnectionPools *int `json:"maxConnectionPools,omitempty" yaml:"maxConnectionPools,omitempty"`
	MaxConnections     *int `json:"maxConnections,omitempty" yaml:"maxConnections,omitempty"`
	MaxPendingRequests *int `json:"maxPendingRequests,omitempty" yaml:"maxPendingRequests,omitempty"`
	MaxRequests        *int `json:"maxRequests,omitempty" yaml:"maxRequests,omitempty"`
	MaxRetries         *int `json:"maxRetries,omitempty" yaml:"maxRetries,omitempty"`
}

// NewCircuitBreaker creates a new CircuitBreaker
func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{}
}

// Timeout represents configuration of timeout
type Timeout struct {
	DownstreamRequestIdleTimeout *int `json:"downstreamRequestIdleTimeout,omitempty" yaml:"downstreamRequestIdleTimeout,omitempty"`
	UpstreamResponseTimeout      *int `json:"upstreamResponseTimeout,omitempty" yaml:"upstreamResponseTimeout,omitempty"`
}

// NewTimeout creates a new Timeout
func NewTimeout() *Timeout {
	return &Timeout{}
}

// RetryPolicy represents configuration for Retry Policy
type RetryPolicy struct {
	Count              *int  `json:"count,omitempty" yaml:"count,omitempty"`
	BaseIntervalMillis *int  `json:"baseIntervalMillis,omitempty" yaml:"baseIntervalMillis,omitempty"`
	StatusCodes        []int `json:"statusCodes,omitempty" yaml:"statusCodes,omitempty"`
}

// NewRetryPolicy creates a new RetryPolicy
func NewRetryPolicy() *RetryPolicy {
	return &RetryPolicy{}
}

// AIRatelimit represents configuration of AIRatelimit settings
type AIRatelimit struct {
	Enabled bool        `json:"enabled" yaml:"enabled"`
	Token   TokenAIRL   `json:"token" yaml:"token"`
	Request RequestAIRL `json:"request" yaml:"request"`
}

// NewAIRatelimit creates a new AIRatelimit with default values
func NewAIRatelimit() *AIRatelimit {
	return &AIRatelimit{
		Enabled: false,
	}
}

// TokenAIRL represents configuration for Token AI rate limit settings
type TokenAIRL struct {
	PromptLimit     int    `json:"promptLimit" yaml:"promptLimit"`
	CompletionLimit int    `json:"completionLimit" yaml:"completionLimit"`
	TotalLimit      int    `json:"totalLimit" yaml:"totalLimit"`
	Unit            string `json:"unit" yaml:"unit"`
}

// NewTokenAIRL creates a new TokenAIRL
func NewTokenAIRL(promptLimit, completionLimit, totalLimit int, unit string) *TokenAIRL {
	return &TokenAIRL{
		PromptLimit:     promptLimit,
		CompletionLimit: completionLimit,
		TotalLimit:      totalLimit,
		Unit:            unit,
	}
}

// RequestAIRL represents configuration for Request AI rate limit settings
type RequestAIRL struct {
	RequestLimit int    `json:"requestLimit" yaml:"requestLimit"`
	Unit         string `json:"unit" yaml:"unit"`
}

// NewRequestAIRL creates a new RequestAIRL
func NewRequestAIRL(requestLimit int, unit string) *RequestAIRL {
	return &RequestAIRL{
		RequestLimit: requestLimit,
		Unit:         unit,
	}
}

// APKOperations represents configuration for APK Operations
type APKOperations struct {
	Target                 *string                 `json:"target,omitempty" yaml:"target,omitempty"`
	Verb                   *string                 `json:"verb,omitempty" yaml:"verb,omitempty"`
	Secured                *bool                   `json:"secured,omitempty" yaml:"secured,omitempty"`
	EndpointConfigurations *EndpointConfigurations `json:"endpointConfigurations,omitempty" yaml:"endpointConfigurations,omitempty"`
	OperationPolicies      *APIOperationPolicies   `json:"operationPolicies,omitempty" yaml:"operationPolicies,omitempty"`
	RateLimit              *RateLimit              `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	Scopes                 []string                `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

// NewAPKOperations creates a new APKOperations
func NewAPKOperations() *APKOperations {
	return &APKOperations{}
}

// APIOperationPolicies represents configuration of APK Operation Policies
type APIOperationPolicies struct {
	Request  []APKRequestOperationPolicy  `json:"request,omitempty" yaml:"request,omitempty"`
	Response []APKResponseOperationPolicy `json:"response,omitempty" yaml:"response,omitempty"`
}

// NewAPIOperationPolicies creates a new APIOperationPolicies
func NewAPIOperationPolicies() *APIOperationPolicies {
	return &APIOperationPolicies{}
}

// APKRequestOperationPolicy represents common type for request operation policies
type APKRequestOperationPolicy interface{}

// APKResponseOperationPolicy represents common type for response operation policies
type APKResponseOperationPolicy interface{}

// RateLimit represents configuration for Rate Limiting
type RateLimit struct {
	RequestsPerUnit int    `json:"requestsPerUnit" yaml:"requestsPerUnit"`
	Unit            string `json:"unit" yaml:"unit"`
}

// NewRateLimit creates a new RateLimit
func NewRateLimit(requestsPerUnit int, unit string) *RateLimit {
	return &RateLimit{
		RequestsPerUnit: requestsPerUnit,
		Unit:            unit,
	}
}

// AuthenticationRequest represents common type for all authentication types
type AuthenticationRequest interface{}

// Authentication represents configuration for authentication types
type Authentication struct {
	AuthType string `json:"authType,omitempty" yaml:"authType,omitempty"`
	Enabled  bool   `json:"enabled" yaml:"enabled"`
}

// NewAuthentication creates a new Authentication with default values
func NewAuthentication() *Authentication {
	return &Authentication{
		Enabled: true,
	}
}

// OAuth2Authentication represents configuration of OAuth2 Authentication type
type OAuth2Authentication struct {
	Authentication
	Required            string `json:"required" yaml:"required"`
	SendTokenToUpstream bool   `json:"sendTokenToUpstream" yaml:"sendTokenToUpstream"`
	HeaderName          string `json:"headerName" yaml:"headerName"`
	HeaderEnable        bool   `json:"headerEnable" yaml:"headerEnable"`
}

// NewOAuth2Authentication creates a new OAuth2Authentication with default values
func NewOAuth2Authentication() *OAuth2Authentication {
	return &OAuth2Authentication{
		Authentication:      *NewAuthentication(),
		Required:            "mandatory",
		SendTokenToUpstream: false,
		HeaderName:          "Authorization",
		HeaderEnable:        true,
	}
}

// JWTAuthentication represents configuration of JWT Authentication type
type JWTAuthentication struct {
	Authentication
	Required            string   `json:"required" yaml:"required"`
	SendTokenToUpstream bool     `json:"sendTokenToUpstream" yaml:"sendTokenToUpstream"`
	HeaderName          string   `json:"headerName" yaml:"headerName"`
	HeaderEnable        bool     `json:"headerEnable" yaml:"headerEnable"`
	Audience            []string `json:"audience" yaml:"audience"`
}

// NewJWTAuthentication creates a new JWTAuthentication with default values
func NewJWTAuthentication() *JWTAuthentication {
	return &JWTAuthentication{
		Authentication:      *NewAuthentication(),
		Required:            "mandatory",
		SendTokenToUpstream: false,
		HeaderName:          "Authorization",
		HeaderEnable:        true,
		Audience:            []string{},
	}
}

// APIKeyAuthentication represents configuration for API Key Auth Type
type APIKeyAuthentication struct {
	Authentication
	Required            string `json:"required" yaml:"required"`
	SendTokenToUpstream bool   `json:"sendTokenToUpstream" yaml:"sendTokenToUpstream"`
	HeaderName          string `json:"headerName" yaml:"headerName"`
	QueryParamName      string `json:"queryParamName" yaml:"queryParamName"`
	HeaderEnable        bool   `json:"headerEnable" yaml:"headerEnable"`
	QueryParamEnable    bool   `json:"queryParamEnable" yaml:"queryParamEnable"`
}

// NewAPIKeyAuthentication creates a new APIKeyAuthentication with default values
func NewAPIKeyAuthentication() *APIKeyAuthentication {
	return &APIKeyAuthentication{
		Authentication:      *NewAuthentication(),
		Required:            "optional",
		SendTokenToUpstream: false,
		HeaderName:          "apiKey",
		QueryParamName:      "apiKey",
		HeaderEnable:        true,
		QueryParamEnable:    false,
	}
}

// MTLSAuthentication represents Mutual SSL configuration of this API
type MTLSAuthentication struct {
	Authentication
	Required     string         `json:"required" yaml:"required"`
	Certificates []ConfigMapRef `json:"certificates" yaml:"certificates"`
}

// NewMTLSAuthentication creates a new MTLSAuthentication with default values
func NewMTLSAuthentication(certificates []ConfigMapRef) *MTLSAuthentication {
	return &MTLSAuthentication{
		Authentication: *NewAuthentication(),
		Required:       "optional",
		Certificates:   certificates,
	}
}

// ConfigMapRef represents configuration for K8s ConfigMap
type ConfigMapRef struct {
	Name string `json:"name" yaml:"name"`
	Key  string `json:"key" yaml:"key"`
}

// NewConfigMapRef creates a new ConfigMapRef
func NewConfigMapRef(name, key string) *ConfigMapRef {
	return &ConfigMapRef{
		Name: name,
		Key:  key,
	}
}

// APKConfAdditionalProperties represents additional properties for APK configuration
type APKConfAdditionalProperties struct {
	Name  *string `json:"name,omitempty" yaml:"name,omitempty"`
	Value *string `json:"value,omitempty" yaml:"value,omitempty"`
}

// NewAPKConfAdditionalProperties creates a new APKConfAdditionalProperties
func NewAPKConfAdditionalProperties() *APKConfAdditionalProperties {
	return &APKConfAdditionalProperties{}
}

// CORSConfiguration represents CORS Configuration of API
type CORSConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled" yaml:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty" yaml:"accessControlAllowOrigins,omitempty"`
	AccessControlAllowCredentials *bool    `json:"accessControlAllowCredentials,omitempty" yaml:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty" yaml:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty" yaml:"accessControlAllowMethods,omitempty"`
	AccessControlAllowMaxAge      *int     `json:"accessControlAllowMaxAge,omitempty" yaml:"accessControlAllowMaxAge,omitempty"`
	AccessControlExposeHeaders    []string `json:"accessControlExposeHeaders,omitempty" yaml:"accessControlExposeHeaders,omitempty"`
}

// NewCORSConfiguration creates a new CORSConfiguration with default values
func NewCORSConfiguration() *CORSConfiguration {
	return &CORSConfiguration{
		CorsConfigurationEnabled: false,
	}
}

// KeyManager represents configuration for a Key Manager
type KeyManager struct {
	Name         string  `json:"name" yaml:"name"`
	Issuer       string  `json:"issuer" yaml:"issuer"`
	JWKSEndpoint string  `json:"JWKSEndpoint" yaml:"JWKSEndpoint"`
	ClaimMapping []Claim `json:"claimMappings" yaml:"claimMappings"`
}

type Claim struct {
	LocalClaim  string `json:"localClaim" yaml:"localClaim"`
	RemoteClaim string `json:"remoteClaim" yaml:"remoteClaim"`
}

// NewKeyManager creates a new KeyManager with default values
func NewKeyManager(name, issuer, jwksEndpoint string) *KeyManager {
	return &KeyManager{
		Name:         name,
		Issuer:       issuer,
		JWKSEndpoint: jwksEndpoint,
	}
}
