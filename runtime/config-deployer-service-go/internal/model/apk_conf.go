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

import (
	"encoding/json"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
)

// APKConf represents the APK configuration for a given API
type APKConf struct {
	ID                     string                        `json:"id,omitempty" yaml:"id,omitempty"`
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

// PolicyName represents enum for all possible policy types.
type PolicyName string

const (
	PolicyNameBackendJWT           PolicyName = "BackendJwt"
	PolicyLuaInterceptor           PolicyName = "LuaInterceptor"
	PolicyWASMInterceptor          PolicyName = "WASMInterceptor"
	PolicyNameAddHeader            PolicyName = "AddHeader"
	PolicyNameSetHeader            PolicyName = "SetHeader"
	PolicyNameRemoveHeader         PolicyName = "RemoveHeader"
	PolicyNameRequestMirror        PolicyName = "RequestMirror"
	PolicyNameRequestRedirect      PolicyName = "RequestRedirect"
	PolicyNameModelBasedRoundRobin PolicyName = "ModelBasedRoundRobin"
)

// BaseOperationPolicy represents common configuration of all policies.
type BaseOperationPolicy struct {
	PolicyName    PolicyName `json:"policyName" yaml:"policyName"`
	PolicyVersion string     `json:"policyVersion" yaml:"policyVersion"`
	PolicyID      *string    `json:"policyId,omitempty" yaml:"policyId,omitempty"`
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

// APKOperationPolicy represents a common interface for all operation policies
type APKOperationPolicy interface {
	GetPolicyName() PolicyName
	GetPolicyVersion() string
	GetPolicyID() *string
}

// APKRequestOperationPolicy represents request operation policies
type APKRequestOperationPolicy struct {
	LuaInterceptorPolicy       *LuaInterceptorPolicy       `json:"luaInterceptorPolicy,omitempty" yaml:"luaInterceptorPolicy,omitempty"`
	WASMInterceptorPolicy      *WASMInterceptorPolicy      `json:"wasmInterceptorPolicy,omitempty" yaml:"wasmInterceptorPolicy,omitempty"`
	BackendJWTPolicy           *BackendJWTPolicy           `json:"backendJWTPolicy,omitempty" yaml:"backendJWTPolicy,omitempty"`
	HeaderModifierPolicy       *HeaderModifierPolicy       `json:"headerModifierPolicy,omitempty" yaml:"headerModifierPolicy,omitempty"`
	RequestMirrorPolicy        *RequestMirrorPolicy        `json:"requestMirrorPolicy,omitempty" yaml:"requestMirrorPolicy,omitempty"`
	RequestRedirectPolicy      *RequestRedirectPolicy      `json:"requestRedirectPolicy,omitempty" yaml:"requestRedirectPolicy,omitempty"`
	ModelBasedRoundRobinPolicy *ModelBasedRoundRobinPolicy `json:"modelBasedRoundRobinPolicy,omitempty" yaml:"modelBasedRoundRobinPolicy,omitempty"`
}

// APKResponseOperationPolicy represents response operation policies
type APKResponseOperationPolicy struct {
	LuaInterceptorPolicy  *LuaInterceptorPolicy  `json:"luaInterceptorPolicy,omitempty" yaml:"luaInterceptorPolicy,omitempty"`
	WASMInterceptorPolicy *WASMInterceptorPolicy `json:"wasmInterceptorPolicy,omitempty" yaml:"wasmInterceptorPolicy,omitempty"`
	HeaderModifierPolicy  *HeaderModifierPolicy  `json:"headerModifierPolicy,omitempty" yaml:"headerModifierPolicy,omitempty"`
}

// GetPolicyName implements APKOperationPolicy interface
func (b BaseOperationPolicy) GetPolicyName() PolicyName {
	return b.PolicyName
}

// GetPolicyVersion implements APKOperationPolicy interface
func (b BaseOperationPolicy) GetPolicyVersion() string {
	return b.PolicyVersion
}

// GetPolicyID implements APKOperationPolicy interface
func (b BaseOperationPolicy) GetPolicyID() *string {
	return b.PolicyID
}

// HeaderModifierPolicy represents header modification configuration for an operation.
type HeaderModifierPolicy struct {
	BaseOperationPolicy
	Parameters HeaderModifierPolicyParameters `json:"parameters" yaml:"parameters"`
}

// HeaderModifierPolicyParameters represents configuration for header modifiers as received from the apk-conf file.
type HeaderModifierPolicyParameters struct {
	HeaderName  string  `json:"headerName" yaml:"headerName"`
	HeaderValue *string `json:"headerValue,omitempty" yaml:"headerValue,omitempty"`
}

// RequestMirrorPolicy represents request mirror configuration for an operation.
type RequestMirrorPolicy struct {
	BaseOperationPolicy
	Parameters RequestMirrorPolicyParameters `json:"parameters" yaml:"parameters"`
}

// RequestMirrorPolicyParameters represents configuration containing the different headers.
type RequestMirrorPolicyParameters struct {
	URLs []string `json:"urls" yaml:"urls"`
}

// RequestRedirectPolicy represents request redirect configuration for an operation.
type RequestRedirectPolicy struct {
	BaseOperationPolicy
	Parameters RequestRedirectPolicyParameters `json:"parameters" yaml:"parameters"`
}

// RequestRedirectPolicyParameters represents configuration containing the different headers.
type RequestRedirectPolicyParameters struct {
	URL        string `json:"url" yaml:"url"`
	StatusCode *int   `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
}

// LuaInterceptorPolicy represents Lua interceptor policy configuration for an operation.
type LuaInterceptorPolicy struct {
	BaseOperationPolicy
	Parameters *LuaInterceptorPolicyParameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// LuaInterceptorPolicyParameters represents configuration for Lua Interceptor Policy parameters.
type LuaInterceptorPolicyParameters struct {
	Name             string  `json:"name" yaml:"name"`
	SourceCode       *string `json:"sourceCode,omitempty" yaml:"sourceCode,omitempty"`
	SourceCodeRef    *string `json:"sourceCodeRef,omitempty" yaml:"sourceCodeRef,omitempty"`
	MountInConfigMap *bool   `json:"mountInConfigMap,omitempty" yaml:"mountInConfigMap,omitempty"`
}

// WASMInterceptorPolicy represents Lua interceptor policy configuration for an operation.
type WASMInterceptorPolicy struct {
	BaseOperationPolicy
	Parameters *WASMInterceptorPolicyParameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// WASMInterceptorPolicyParameters represents configuration for Lua Interceptor Policy parameters.
type WASMInterceptorPolicyParameters struct {
	Name            string   `json:"name" yaml:"name"`
	RootID          string   `json:"rootId" yaml:"rootId"`
	URL             *string  `json:"url,omitempty" yaml:"url,omitempty"`
	Image           *string  `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy *string  `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Config          *string  `json:"config,omitempty" yaml:"config,omitempty"`
	FailOpen        *bool    `json:"failOpen,omitempty" yaml:"failOpen,omitempty"`
	HostKeys        []string `json:"hostKeys,omitempty" yaml:"hostKeys,omitempty"`
}

// BackendJWTPolicy represents configuration for Backend JWT Policy.
type BackendJWTPolicy struct {
	BaseOperationPolicy
	Parameters *BackendJWTPolicyParameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// BackendJWTPolicyParameters represents configuration for Backend JWT Policy parameters.
type BackendJWTPolicyParameters struct {
	Encoding         *string        `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	SigningAlgorithm *string        `json:"signingAlgorithm,omitempty" yaml:"signingAlgorithm,omitempty"`
	Header           *string        `json:"header,omitempty" yaml:"header,omitempty"`
	TokenTTL         *int           `json:"tokenTTL,omitempty" yaml:"tokenTTL,omitempty"`
	CustomClaims     []CustomClaims `json:"customClaims,omitempty" yaml:"customClaims,omitempty"`
}

// ModelBasedRoundRobinPolicy represents model based round robin policy configuration for an operation.
type ModelBasedRoundRobinPolicy struct {
	BaseOperationPolicy
	Parameters ModelBasedRoundRobinPolicyParameters `json:"parameters" yaml:"parameters"`
}

// ModelBasedRoundRobinPolicyParameters represents configuration for model based round robin policy parameters.
type ModelBasedRoundRobinPolicyParameters struct {
	OnQuotaExceedSuspendDuration int            `json:"onQuotaExceedSuspendDuration" yaml:"onQuotaExceedSuspendDuration"`
	ProductionModels             []ModelRouting `json:"productionModels" yaml:"productionModels"`
	SandboxModels                []ModelRouting `json:"sandboxModels" yaml:"sandboxModels"`
}

// ModelRouting represents configuration for model routing.
type ModelRouting struct {
	Model    string `json:"model" yaml:"model"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Weight   int    `json:"weight" yaml:"weight"`
}

// CustomClaims represents configuration for Custom Claims.
type CustomClaims struct {
	Claim string `json:"claim" yaml:"claim"`
	Value string `json:"value" yaml:"value"`
	Type  string `json:"type" yaml:"type"`
}

// Helper methods for handling union types

// UnmarshalJSON provides custom unmarshaling for APKRequestOperationPolicy
func (p *APKRequestOperationPolicy) UnmarshalJSON(data []byte) error {
	var base BaseOperationPolicy
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	switch base.PolicyName {
	case PolicyLuaInterceptor:
		p.LuaInterceptorPolicy = &LuaInterceptorPolicy{}
		return json.Unmarshal(data, p.LuaInterceptorPolicy)
	case PolicyWASMInterceptor:
		p.WASMInterceptorPolicy = &WASMInterceptorPolicy{}
		return json.Unmarshal(data, p.WASMInterceptorPolicy)
	case PolicyNameBackendJWT:
		p.BackendJWTPolicy = &BackendJWTPolicy{}
		return json.Unmarshal(data, p.BackendJWTPolicy)
	case PolicyNameAddHeader, PolicyNameSetHeader, PolicyNameRemoveHeader:
		p.HeaderModifierPolicy = &HeaderModifierPolicy{}
		return json.Unmarshal(data, p.HeaderModifierPolicy)
	case PolicyNameRequestMirror:
		p.RequestMirrorPolicy = &RequestMirrorPolicy{}
		return json.Unmarshal(data, p.RequestMirrorPolicy)
	case PolicyNameRequestRedirect:
		p.RequestRedirectPolicy = &RequestRedirectPolicy{}
		return json.Unmarshal(data, p.RequestRedirectPolicy)
	case PolicyNameModelBasedRoundRobin:
		p.ModelBasedRoundRobinPolicy = &ModelBasedRoundRobinPolicy{}
		return json.Unmarshal(data, p.ModelBasedRoundRobinPolicy)
	}

	return nil
}

// UnmarshalJSON provides custom unmarshaling for APKResponseOperationPolicy
func (p *APKResponseOperationPolicy) UnmarshalJSON(data []byte) error {
	var base BaseOperationPolicy
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	switch base.PolicyName {
	case PolicyLuaInterceptor:
		p.LuaInterceptorPolicy = &LuaInterceptorPolicy{}
		return json.Unmarshal(data, p.LuaInterceptorPolicy)
	case PolicyWASMInterceptor:
		p.WASMInterceptorPolicy = &WASMInterceptorPolicy{}
		return json.Unmarshal(data, p.WASMInterceptorPolicy)
	case PolicyNameAddHeader, PolicyNameSetHeader, PolicyNameRemoveHeader:
		p.HeaderModifierPolicy = &HeaderModifierPolicy{}
		return json.Unmarshal(data, p.HeaderModifierPolicy)
	}

	return nil
}

// GetActivePolicy returns the active policy for request operation policy
func (p *APKRequestOperationPolicy) GetActivePolicy() APKOperationPolicy {
	if p.LuaInterceptorPolicy != nil {
		return p.LuaInterceptorPolicy
	}
	if p.WASMInterceptorPolicy != nil {
		return p.WASMInterceptorPolicy
	}
	if p.BackendJWTPolicy != nil {
		return p.BackendJWTPolicy
	}
	if p.HeaderModifierPolicy != nil {
		return p.HeaderModifierPolicy
	}
	if p.RequestMirrorPolicy != nil {
		return p.RequestMirrorPolicy
	}
	if p.RequestRedirectPolicy != nil {
		return p.RequestRedirectPolicy
	}
	if p.ModelBasedRoundRobinPolicy != nil {
		return p.ModelBasedRoundRobinPolicy
	}
	return nil
}

// GetActivePolicy returns the active policy for response operation policy
func (p *APKResponseOperationPolicy) GetActivePolicy() APKOperationPolicy {
	if p.LuaInterceptorPolicy != nil {
		return p.LuaInterceptorPolicy
	}
	if p.WASMInterceptorPolicy != nil {
		return p.WASMInterceptorPolicy
	}
	if p.HeaderModifierPolicy != nil {
		return p.HeaderModifierPolicy
	}
	return nil
}

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
	Name         string      `json:"name" yaml:"name"`
	Issuer       string      `json:"issuer" yaml:"issuer"`
	JWKSEndpoint string      `json:"JWKSEndpoint" yaml:"JWKSEndpoint"`
	ClaimMapping []Claim     `json:"claimMappings" yaml:"claimMappings"`
	K8sBackend   *K8sBackend `json:"k8sBackend,omitempty" yaml:"k8sBackend,omitempty"`
}

// K8sBackend represents the backend configuration for a Key Manager
type K8sBackend struct {
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace *string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Port      *int    `json:"port,omitempty" yaml:"port,omitempty"`
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
