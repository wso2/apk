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
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// ResourceParams contains httproute related parameters
type ResourceParams struct {
	AuthSchemes               map[string]dpv1alpha1.Authentication
	ResourceAuthSchemes       map[string]dpv1alpha1.Authentication
	APIPolicies               map[string]dpv1alpha1.APIPolicy
	ResourceAPIPolicies       map[string]dpv1alpha1.APIPolicy
	InterceptorServiceMapping map[string]dpv1alpha1.InterceptorService
	BackendMapping            dpv1alpha1.BackendMapping
	ResourceScopes            map[string]dpv1alpha1.Scope
	RateLimitPolicies         map[string]dpv1alpha1.RateLimitPolicy
	ResourceRateLimitPolicies map[string]dpv1alpha1.RateLimitPolicy
}

// SetInfoHTTPRouteCR populates resources and endpoints of adapterInternalAPI. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (swagger *AdapterInternalAPI) SetInfoHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, resourceParams ResourceParams) error {
	var resources []*Resource
	//TODO(amali) add gateway level securities after gateway crd has implemented
	outputAuthScheme := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.AuthSchemes)))
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.APIPolicies)))
	outputRatelimitPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(resourceParams.RateLimitPolicies)))

	disableScopes := true
	disableAuthentications := false
	config := config.ReadConfigs()

	var authScheme *dpv1alpha1.Authentication
	if outputAuthScheme != nil {
		authScheme = *outputAuthScheme
		if authScheme.Spec.Override != nil && authScheme.Spec.Override.ExternalService.Disabled != nil {
			disableAuthentications = *authScheme.Spec.Override.ExternalService.Disabled
		}
	}
	var apiPolicy *dpv1alpha1.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}

	var ratelimitPolicy *dpv1alpha1.RateLimitPolicy
	if outputRatelimitPolicy != nil {
		ratelimitPolicy = concatRateLimitPolicies(*outputRatelimitPolicy, nil)
	}

	for _, rule := range httpRoute.Spec.Rules {
		var endPoints []Endpoint
		var policies = OperationPolicies{}
		var circuitBreaker *dpv1alpha1.CircuitBreaker
		var healthCheck *dpv1alpha1.HealthCheck
		resourceAuthScheme := authScheme
		resourceAPIPolicy := apiPolicy
		var resourceRatelimitPolicy *dpv1alpha1.RateLimitPolicy
		hasPolicies := false
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
			resolvedBackend, ok := resourceParams.BackendMapping[backendName]
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
					timeoutInMillis = resolvedBackend.Timeout.RouteTimeoutSeconds * 1000
					idleTimeoutInSeconds = resolvedBackend.Timeout.RouteIdleTimeoutSeconds
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
				endPoints = append(endPoints, GetEndpoints(backendName, httpRouteParams.BackendMapping)...)
				backendBasePath = GetBackendBasePath(backendName, httpRouteParams.BackendMapping)
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
			hasPolicies = true
			switch filter.Type {
			case gwapiv1b1.HTTPRouteFilterURLRewrite:
				policyParameters := make(map[string]interface{})
				policyParameters[constants.RewritePathType] = filter.URLRewrite.Path.Type
				policyParameters[constants.IncludeQueryParams] = true

				switch filter.URLRewrite.Path.Type {
				case gwapiv1b1.FullPathHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = backendBasePath + *filter.URLRewrite.Path.ReplaceFullPath
				case gwapiv1b1.PrefixMatchHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = backendBasePath + *filter.URLRewrite.Path.ReplacePrefixMatch
				}

				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1b1.HTTPRouteFilterURLRewrite),
					Action:     constants.ActionRewritePath,
					Parameters: policyParameters,
				})
				hasURLRewritePolicy = true
			case gwapiv1b1.HTTPRouteFilterExtensionRef:
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
			case gwapiv1b1.HTTPRouteFilterRequestHeaderModifier:
				for _, header := range filter.RequestHeaderModifier.Add {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Request = append(policies.Request, Policy{
						PolicyName: string(gwapiv1b1.HTTPRouteFilterRequestHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.RequestHeaderModifier.Remove {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header)

					policies.Request = append(policies.Request, Policy{
						PolicyName: string(gwapiv1b1.HTTPRouteFilterRequestHeaderModifier),
						Action:     constants.ActionHeaderRemove,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.RequestHeaderModifier.Set {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Request = append(policies.Request, Policy{
						PolicyName: string(gwapiv1b1.HTTPRouteFilterRequestHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
			case gwapiv1b1.HTTPRouteFilterResponseHeaderModifier:
				for _, header := range filter.ResponseHeaderModifier.Add {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Response = append(policies.Response, Policy{
						PolicyName: string(gwapiv1b1.HTTPRouteFilterResponseHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.ResponseHeaderModifier.Remove {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header)

					policies.Response = append(policies.Response, Policy{
						PolicyName: string(gwapiv1b1.HTTPRouteFilterResponseHeaderModifier),
						Action:     constants.ActionHeaderRemove,
						Parameters: policyParameters,
					})
				}
				for _, header := range filter.ResponseHeaderModifier.Set {
					policyParameters := make(map[string]interface{})
					policyParameters[constants.HeaderName] = string(header.Name)
					policyParameters[constants.HeaderValue] = string(header.Value)

					policies.Response = append(policies.Response, Policy{
						PolicyName: string(gwapiv1b1.HTTPRouteFilterResponseHeaderModifier),
						Action:     constants.ActionHeaderAdd,
						Parameters: policyParameters,
					})
				}
			}
		}
		if resourceAPIPolicy == apiPolicy {
			apiPolicySelected := concatAPIPolicies(apiPolicy, nil)
			addOperationLevelInterceptors(&policies, apiPolicySelected, resourceParams.InterceptorServiceMapping, resourceParams.BackendMapping)
		} else {
			addOperationLevelInterceptors(&policies, resourceAPIPolicy, resourceParams.InterceptorServiceMapping, resourceParams.BackendMapping)
		}

		loggers.LoggerOasparser.Debugf("Calculating auths for API ..., API_UUID = %v", logging.GetValueFromLogContext("API_UUID"))
		apiAuth := getSecurity(resourceAuthScheme)
		if len(rule.BackendRefs) < 1 {
			return fmt.Errorf("no backendref were provided")
		}

		for _, match := range rule.Matches {
			if !hasURLRewritePolicy {
				policyParameters := make(map[string]interface{})
				if *match.Path.Type == gwapiv1b1.PathMatchPathPrefix {
					policyParameters[constants.RewritePathType] = gwapiv1b1.PrefixMatchHTTPPathModifier
				} else {
					policyParameters[constants.RewritePathType] = gwapiv1b1.FullPathHTTPPathModifier
				}
				policyParameters[constants.IncludeQueryParams] = true
				policyParameters[constants.RewritePathResourcePath] = strings.TrimSuffix(backendBasePath, "/") + *match.Path.Value
				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1b1.HTTPRouteFilterURLRewrite),
					Action:     constants.ActionRewritePath,
					Parameters: policyParameters,
				})
				hasPolicies = true
			}
			resourcePath := swagger.xWso2Basepath + *match.Path.Value
			resource := &Resource{path: resourcePath,
				methods: getAllowedOperations(match.Method, policies, apiAuth,
					parseRateLimitPolicyToInternal(resourceRatelimitPolicy), scopes),
				pathMatchType: *match.Path.Type,
				hasPolicies:   hasPolicies,
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

	swagger.RateLimitPolicy = parseRateLimitPolicyToInternal(ratelimitPolicy)
	swagger.resources = resources
	apiPolicySelected := concatAPIPolicies(apiPolicy, nil)
	swagger.xWso2Cors = getCorsConfigFromAPIPolicy(apiPolicySelected)
	swagger.disableAuthentications = disableAuthentications
	swagger.disableScopes = disableScopes

	// Check whether the API has a backend JWT token
	if apiPolicySelected != nil && apiPolicySelected.Spec.Override != nil && apiPolicySelected.Spec.Override.BackendJWTToken != nil {
		swagger.backendJWTTokenInfo = parseBackendJWTTokenToInternal(apiPolicySelected.Spec.Override.BackendJWTToken)
	}
	return nil
}

func parseBackendJWTTokenToInternal(backendJWTToken *dpv1alpha1.BackendJWTToken) *BackendJWTTokenInfo {
	var customClaims []ClaimMapping
	for _, value := range backendJWTToken.CustomClaims {
		valType := value.Type
		claim := value.Claim
		value := value.Value
		claimMapping := ClaimMapping{
			Claim: claim,
			Value: ClaimVal{
				Value: value,
				Type:  valType,
			},
		}
		customClaims = append(customClaims, claimMapping)
	}
	backendJWTTokenInternal := &BackendJWTTokenInfo{
		Enabled:          backendJWTToken.Enabled,
		Encoding:         backendJWTToken.Encoding,
		Header:           backendJWTToken.Header,
		SigningAlgorithm: backendJWTToken.SigningAlgorithm,
		CustomClaims:     customClaims,
		TokenTTL:         backendJWTToken.TokenTTL,
	}
	return backendJWTTokenInternal
}

func getCorsConfigFromAPIPolicy(apiPolicy *dpv1alpha1.APIPolicy) *CorsConfig {
	var corsConfig *CorsConfig
	if apiPolicy != nil && apiPolicy.Spec.Override != nil {
		if apiPolicy.Spec.Override.CORSPolicy != nil {
			corsConfig = &CorsConfig{
				Enabled:                       apiPolicy.Spec.Override.CORSPolicy.Enabled,
				AccessControlAllowCredentials: apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowCredentials,
				AccessControlAllowHeaders:     apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowHeaders,
				AccessControlAllowMethods:     apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowMethods,
				AccessControlAllowOrigins:     apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowOrigins,
				AccessControlExposeHeaders:    apiPolicy.Spec.Override.CORSPolicy.AccessControlExposeHeaders,
			}
			if apiPolicy.Spec.Override.CORSPolicy.AccessControlMaxAge != nil {
				corsConfig.AccessControlMaxAge = apiPolicy.Spec.Override.CORSPolicy.AccessControlMaxAge
			}
		}
	}
	return corsConfig
}

func parseRateLimitPolicyToInternal(ratelimitPolicy *dpv1alpha1.RateLimitPolicy) *RateLimitPolicy {
	var rateLimitPolicyInternal *RateLimitPolicy
	if ratelimitPolicy != nil {
		if ratelimitPolicy.Spec.Override.API.RateLimit.RequestsPerUnit > 0 {
			rateLimitPolicyInternal = &RateLimitPolicy{
				Count:    ratelimitPolicy.Spec.Override.API.RateLimit.RequestsPerUnit,
				SpanUnit: ratelimitPolicy.Spec.Override.API.RateLimit.Unit,
			}
		}
	}
	return rateLimitPolicyInternal
}

// addOperationLevelInterceptors add the operation level interceptor policy to the policies
func addOperationLevelInterceptors(policies *OperationPolicies, apiPolicy *dpv1alpha1.APIPolicy,
	interceptorServicesMapping map[string]dpv1alpha1.InterceptorService, backendMapping dpv1alpha1.BackendMapping) {
	if apiPolicy != nil && apiPolicy.Spec.Override != nil {
		if len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
			requestInterceptor := interceptorServicesMapping[types.NamespacedName{
				Name:      apiPolicy.Spec.Override.RequestInterceptors[0].Name,
				Namespace: apiPolicy.Spec.Override.RequestInterceptors[0].Namespace,
			}.String()].Spec
			policyParameters := make(map[string]interface{})
			backendName := types.NamespacedName{
				Name:      requestInterceptor.BackendRef.Name,
				Namespace: utils.GetNamespace((*gwapiv1b1.Namespace)(&requestInterceptor.BackendRef.Namespace), apiPolicy.Namespace),
			}
			endpoints := GetEndpoints(backendName, backendMapping)
			if len(endpoints) > 0 {
				policyParameters[constants.InterceptorEndpoints] = endpoints
				policyParameters[constants.InterceptorServiceIncludes] = requestInterceptor.Includes
				policies.Request = append(policies.Request, Policy{
					PolicyName: constants.PolicyRequestInterceptor,
					Action:     constants.ActionInterceptorService,
					Parameters: policyParameters,
				})
			}
		}
		if len(apiPolicy.Spec.Override.ResponseInterceptors) > 0 {
			responseInterceptor := interceptorServicesMapping[types.NamespacedName{
				Name:      apiPolicy.Spec.Override.ResponseInterceptors[0].Name,
				Namespace: apiPolicy.Spec.Override.ResponseInterceptors[0].Namespace,
			}.String()].Spec
			policyParameters := make(map[string]interface{})
			backendName := types.NamespacedName{
				Name:      responseInterceptor.BackendRef.Name,
				Namespace: utils.GetNamespace((*gwapiv1b1.Namespace)(&responseInterceptor.BackendRef.Namespace), apiPolicy.Namespace),
			}
			endpoints := GetEndpoints(backendName, backendMapping)
			if len(endpoints) > 0 {
				policyParameters[constants.InterceptorEndpoints] = endpoints
				policyParameters[constants.InterceptorServiceIncludes] = responseInterceptor.Includes
				policies.Response = append(policies.Response, Policy{
					PolicyName: constants.PolicyResponseInterceptor,
					Action:     constants.ActionInterceptorService,
					Parameters: policyParameters,
				})
			}
		}
	}
}

// GetEndpoints creates endpoints using resolved backends in backendMapping
func GetEndpoints(backendName types.NamespacedName, backendMapping dpv1alpha1.BackendMapping) []Endpoint {
	endpoints := []Endpoint{}
	backend, ok := backendMapping[backendName]
	if ok && backend != nil {
		if len(backend.Services) > 0 {
			for _, service := range backend.Services {
				endpoints = append(endpoints, Endpoint{
					Host:        service.Host,
					Port:        service.Port,
					URLType:     string(backend.Protocol),
					Certificate: []byte(backend.TLS.ResolvedCertificate),
					AllowedSANs: backend.TLS.AllowedSANs,
				})
			}
		}
	}
	return endpoints
}

// GetBackendBasePath gets basePath of the the Backend
func GetBackendBasePath(backendName types.NamespacedName, backendMapping dpv1alpha1.BackendMapping) string {
	backend, ok := backendMapping[backendName]
	if ok && backend != nil {
		if len(backend.Services) > 0 {
			return backend.BasePath
		}
	}
	return ""
}

func concatRateLimitPolicies(schemeUp *dpv1alpha1.RateLimitPolicy, schemeDown *dpv1alpha1.RateLimitPolicy) *dpv1alpha1.RateLimitPolicy {
	finalRateLimit := dpv1alpha1.RateLimitPolicy{}
	if schemeUp != nil && schemeDown != nil {
		finalRateLimit.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	} else if schemeUp != nil {
		finalRateLimit.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, nil, nil)
	} else if schemeDown != nil {
		finalRateLimit.Spec.Override = utils.SelectPolicy(nil, nil, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	}
	return &finalRateLimit
}

func concatAPIPolicies(schemeUp *dpv1alpha1.APIPolicy, schemeDown *dpv1alpha1.APIPolicy) *dpv1alpha1.APIPolicy {
	apiPolicy := dpv1alpha1.APIPolicy{}
	if schemeUp != nil && schemeDown != nil {
		apiPolicy.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	} else if schemeUp != nil {
		apiPolicy.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, nil, nil)
	} else if schemeDown != nil {
		apiPolicy.Spec.Override = utils.SelectPolicy(nil, nil, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	}
	return &apiPolicy
}

func concatAuthSchemes(schemeUp *dpv1alpha1.Authentication, schemeDown *dpv1alpha1.Authentication) *dpv1alpha1.Authentication {
	finalAuth := dpv1alpha1.Authentication{
		Spec: dpv1alpha1.AuthenticationSpec{},
	}
	if schemeUp != nil && schemeDown != nil {
		finalAuth.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	} else if schemeUp != nil {
		finalAuth.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, nil, nil)
	} else if schemeDown != nil {
		finalAuth.Spec.Override = utils.SelectPolicy(nil, nil, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	}
	return &finalAuth
}

// getSecurity returns security schemes and it's definitions with flag to indicate if security is disabled
// make sure authscheme only has external service override values. (i.e. empty default values)
// tip: use concatScheme method
func getSecurity(authScheme *dpv1alpha1.Authentication) *Authentication {
	auth := &Authentication{Disabled: false,
		TestConsoleKey: &TestConsoleKey{Header: constants.TestConsoleKeyHeader},
		JWT: &JWT{Header: constants.AuthorizationHeader},
	}
	if (authScheme != nil && authScheme.Spec.Override.ExternalService.AuthTypes != nil && authScheme.Spec.Override.ExternalService.AuthTypes.JWT.Disabled) {
		auth = &Authentication{Disabled: false,
			TestConsoleKey: &TestConsoleKey{Header: constants.TestConsoleKeyHeader},
		}
	}
	if authScheme != nil {
		if authScheme.Spec.Override.ExternalService.Disabled != nil && *authScheme.Spec.Override.ExternalService.Disabled {
			loggers.LoggerOasparser.Debug("Disabled security")
			return &Authentication{Disabled: true}
		}
		if authScheme.Spec.Override.ExternalService.AuthTypes.APIKey != nil {
			var apiKeys []APIKey
			for _, apiKey := range authScheme.Spec.Override.ExternalService.AuthTypes.APIKey {
				apiKeys = append(apiKeys, APIKey{
					Name: apiKey.Name,
					In:   apiKey.In,
				})
			}
			auth.APIKey = apiKeys
		}
	}
	return auth
}

// getAllowedOperations retuns a list of allowed operatons, if httpMethod is not specified then all methods are allowed.
func getAllowedOperations(httpMethod *gwapiv1b1.HTTPMethod, policies OperationPolicies, auth *Authentication,
	ratelimitPolicy *RateLimitPolicy, scopes []string) []*Operation {
	if httpMethod != nil {
		return []*Operation{{iD: uuid.New().String(), method: string(*httpMethod), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes}}
	}
	return []*Operation{{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodGet), policies: policies,
		auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPost), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodDelete), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPatch), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPut), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodHead), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodOptions), policies: policies,
			auth: auth, RateLimitPolicy: ratelimitPolicy, scopes: scopes}}
}

// SetInfoAPICR populates ID, ApiType, Version and XWso2BasePath of adapterInternalAPI.
func (swagger *AdapterInternalAPI) SetInfoAPICR(api dpv1alpha1.API) {
	swagger.UUID = string(api.ObjectMeta.UID)
	swagger.title = api.Spec.APIDisplayName
	swagger.apiType = api.Spec.APIType
	swagger.version = api.Spec.APIVersion
	swagger.xWso2Basepath = api.Spec.Context
	swagger.OrganizationID = api.Spec.Organization
	swagger.IsSystemAPI = api.Spec.SystemAPI
	swagger.APIProperties = api.Spec.APIProperties
}
