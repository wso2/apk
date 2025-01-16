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
	"github.com/google/uuid"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	dpv1alpha2 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha2"
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
	dpv1alpha4 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha4"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// ResourceParams contains httproute related parameters
type ResourceParams struct {
	AuthSchemes               map[string]dpv1alpha2.Authentication
	ResourceAuthSchemes       map[string]dpv1alpha2.Authentication
	APIPolicies               map[string]dpv1alpha4.APIPolicy
	ResourceAPIPolicies       map[string]dpv1alpha4.APIPolicy
	InterceptorServiceMapping map[string]dpv1alpha1.InterceptorService
	BackendJWTMapping         map[string]dpv1alpha1.BackendJWT
	BackendMapping            map[string]*dpv1alpha4.ResolvedBackend
	ResourceScopes            map[string]dpv1alpha1.Scope
	RateLimitPolicies         map[string]dpv1alpha3.RateLimitPolicy
	ResourceRateLimitPolicies map[string]dpv1alpha3.RateLimitPolicy
}

func parseBackendJWTTokenToInternal(backendJWTToken dpv1alpha1.BackendJWTSpec) *BackendJWTTokenInfo {
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
		Enabled:          true,
		Encoding:         backendJWTToken.Encoding,
		Header:           backendJWTToken.Header,
		SigningAlgorithm: backendJWTToken.SigningAlgorithm,
		CustomClaims:     customClaims,
		TokenTTL:         backendJWTToken.TokenTTL,
	}
	return backendJWTTokenInternal
}

func getCorsConfigFromAPIPolicy(apiPolicy *dpv1alpha4.APIPolicy) *CorsConfig {
	globalCorsConfig := config.ReadConfigs().Enforcer.Cors

	var corsConfig = CorsConfig{
		Enabled:                       globalCorsConfig.Enabled,
		AccessControlAllowCredentials: globalCorsConfig.AccessControlAllowCredentials,
		AccessControlAllowHeaders:     globalCorsConfig.AccessControlAllowHeaders,
		AccessControlAllowMethods:     globalCorsConfig.AccessControlAllowMethods,
		AccessControlAllowOrigins:     globalCorsConfig.AccessControlAllowOrigins,
		AccessControlExposeHeaders:    globalCorsConfig.AccessControlExposeHeaders,
		AccessControlMaxAge:           globalCorsConfig.AccessControlMaxAge,
	}
	if apiPolicy != nil && apiPolicy.Spec.Override != nil {
		if apiPolicy.Spec.Override.CORSPolicy != nil {
			corsConfig.Enabled = apiPolicy.Spec.Override.CORSPolicy.Enabled
			corsConfig.AccessControlAllowCredentials = apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowCredentials
			corsConfig.AccessControlAllowHeaders = apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowHeaders
			corsConfig.AccessControlAllowMethods = apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowMethods
			corsConfig.AccessControlAllowOrigins = apiPolicy.Spec.Override.CORSPolicy.AccessControlAllowOrigins
			corsConfig.AccessControlExposeHeaders = apiPolicy.Spec.Override.CORSPolicy.AccessControlExposeHeaders
			if apiPolicy.Spec.Override.CORSPolicy.AccessControlMaxAge != nil {
				corsConfig.AccessControlMaxAge = apiPolicy.Spec.Override.CORSPolicy.AccessControlMaxAge
			}
		}
	}
	return &corsConfig
}

func parseRateLimitPolicyToInternal(ratelimitPolicy *dpv1alpha3.RateLimitPolicy) *RateLimitPolicy {
	var rateLimitPolicyInternal *RateLimitPolicy
	if ratelimitPolicy != nil && ratelimitPolicy.Spec.Override != nil {
		if ratelimitPolicy.Spec.Override.API.RequestsPerUnit > 0 {
			rateLimitPolicyInternal = &RateLimitPolicy{
				Count:    ratelimitPolicy.Spec.Override.API.RequestsPerUnit,
				SpanUnit: ratelimitPolicy.Spec.Override.API.Unit,
			}
		}
	}
	return rateLimitPolicyInternal
}

// addOperationLevelInterceptors add the operation level interceptor policy to the policies
func addOperationLevelInterceptors(policies *OperationPolicies, apiPolicy *dpv1alpha4.APIPolicy,
	interceptorServicesMapping map[string]dpv1alpha1.InterceptorService,
	backendMapping map[string]*dpv1alpha4.ResolvedBackend, namespace string) {
	if apiPolicy != nil && apiPolicy.Spec.Override != nil {
		if len(apiPolicy.Spec.Override.RequestInterceptors) > 0 {
			requestInterceptor := interceptorServicesMapping[types.NamespacedName{
				Name:      apiPolicy.Spec.Override.RequestInterceptors[0].Name,
				Namespace: namespace,
			}.String()].Spec
			policyParameters := make(map[string]interface{})
			backendName := types.NamespacedName{
				Name:      requestInterceptor.BackendRef.Name,
				Namespace: namespace,
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
				Namespace: namespace,
			}.String()].Spec
			policyParameters := make(map[string]interface{})
			backendName := types.NamespacedName{
				Name:      responseInterceptor.BackendRef.Name,
				Namespace: namespace,
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
func GetEndpoints(backendName types.NamespacedName, backendMapping map[string]*dpv1alpha4.ResolvedBackend) []Endpoint {
	endpoints := []Endpoint{}
	backend, ok := backendMapping[backendName.String()]
	if ok && backend != nil {
		if len(backend.Services) > 0 {
			for _, service := range backend.Services {
				endpoints = append(endpoints, Endpoint{
					Host:        service.Host,
					Port:        service.Port,
					Basepath:    backend.BasePath,
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
func GetBackendBasePath(backendName types.NamespacedName, backendMapping map[string]*dpv1alpha4.ResolvedBackend) string {
	backend, ok := backendMapping[backendName.String()]
	if ok && backend != nil {
		if len(backend.Services) > 0 {
			return backend.BasePath
		}
	}
	return ""
}

func concatRateLimitPolicies(schemeUp *dpv1alpha3.RateLimitPolicy, schemeDown *dpv1alpha3.RateLimitPolicy) *dpv1alpha3.RateLimitPolicy {
	finalRateLimit := dpv1alpha3.RateLimitPolicy{}
	if schemeUp != nil && schemeDown != nil {
		finalRateLimit.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	} else if schemeUp != nil {
		finalRateLimit.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, nil, nil)
	} else if schemeDown != nil {
		finalRateLimit.Spec.Override = utils.SelectPolicy(nil, nil, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	}
	return &finalRateLimit
}

func concatAPIPolicies(schemeUp *dpv1alpha4.APIPolicy, schemeDown *dpv1alpha4.APIPolicy) *dpv1alpha4.APIPolicy {
	apiPolicy := dpv1alpha4.APIPolicy{}
	if schemeUp != nil && schemeDown != nil {
		apiPolicy.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	} else if schemeUp != nil {
		apiPolicy.Spec.Override = utils.SelectPolicy(&schemeUp.Spec.Override, &schemeUp.Spec.Default, nil, nil)
	} else if schemeDown != nil {
		apiPolicy.Spec.Override = utils.SelectPolicy(nil, nil, &schemeDown.Spec.Override, &schemeDown.Spec.Default)
	}
	return &apiPolicy
}

func concatAuthSchemes(schemeUp *dpv1alpha2.Authentication, schemeDown *dpv1alpha2.Authentication) *dpv1alpha2.Authentication {
	finalAuth := dpv1alpha2.Authentication{
		Spec: dpv1alpha2.AuthenticationSpec{},
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
func getSecurity(authScheme *dpv1alpha2.Authentication) *Authentication {
	authHeader := constants.AuthorizationHeader
	if authScheme != nil && authScheme.Spec.Override != nil && authScheme.Spec.Override.AuthTypes != nil && len(authScheme.Spec.Override.AuthTypes.OAuth2.Header) > 0 {
		authHeader = authScheme.Spec.Override.AuthTypes.OAuth2.Header
	}
	sendTokenToUpstream := false
	if authScheme != nil && authScheme.Spec.Override != nil && authScheme.Spec.Override.AuthTypes != nil {
		sendTokenToUpstream = authScheme.Spec.Override.AuthTypes.OAuth2.SendTokenToUpstream
	}
	auth := &Authentication{Disabled: false,
		Oauth2: &Oauth2{Header: authHeader, SendTokenToUpstream: sendTokenToUpstream},
	}
	if authScheme != nil && authScheme.Spec.Override != nil {
		if authScheme.Spec.Override.Disabled != nil && *authScheme.Spec.Override.Disabled {
			return &Authentication{Disabled: true}
		}
		authFound := false
		if authScheme.Spec.Override.AuthTypes != nil && !authScheme.Spec.Override.AuthTypes.OAuth2.Disabled {
			authFound = true
		} else {
			auth = &Authentication{Disabled: false}
		}
		if authScheme.Spec.Override.AuthTypes != nil && authScheme.Spec.Override.AuthTypes.JWT.Disabled != nil && !*authScheme.Spec.Override.AuthTypes.JWT.Disabled {
			audience := make([]string, 0)
			if len(authScheme.Spec.Override.AuthTypes.JWT.Audience) > 0 {
				audience = authScheme.Spec.Override.AuthTypes.JWT.Audience
			}
			jwtHeader := constants.TestConsoleKeyHeader
			if len(authScheme.Spec.Override.AuthTypes.JWT.Header) > 0 {
				jwtHeader = authScheme.Spec.Override.AuthTypes.JWT.Header
			}
			auth.JWT = &JWT{Header: jwtHeader, SendTokenToUpstream: sendTokenToUpstream, Audience: audience}
			authFound = true
		}
		if authScheme.Spec.Override.AuthTypes != nil && authScheme.Spec.Override.AuthTypes.APIKey != nil {
			authFound = authFound || len(authScheme.Spec.Override.AuthTypes.APIKey.Keys) > 0
			var apiKeys []APIKey
			for _, apiKey := range authScheme.Spec.Override.AuthTypes.APIKey.Keys {
				apiKeys = append(apiKeys, APIKey{
					Name:                apiKey.Name,
					In:                  apiKey.In,
					SendTokenToUpstream: apiKey.SendTokenToUpstream,
				})
			}
			auth.APIKey = apiKeys
		}
		if !authFound {
			return &Authentication{Disabled: true}
		}
	}
	return auth
}

// getAllowedOperations retuns a list of allowed operatons, if httpMethod is not specified then all methods are allowed.
func getAllowedOperations(matchID string, httpMethod *gwapiv1.HTTPMethod, policies OperationPolicies, auth *Authentication,
	ratelimitPolicy *RateLimitPolicy, scopes []string, mirrorEndpointClusters []*EndpointCluster) []*Operation {
	if httpMethod != nil {
		return []*Operation{{iD: uuid.New().String(), method: string(*httpMethod), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID}}
	}
	return []*Operation{{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodGet), policies: policies,
		auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, matchID: matchID},
		{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodPost), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID},
		{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodDelete), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID},
		{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodPatch), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID},
		{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodPut), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID},
		{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodHead), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID},
		{iD: uuid.New().String(), method: string(gwapiv1.HTTPMethodOptions), policies: policies,
			auth: auth, rateLimitPolicy: ratelimitPolicy, scopes: scopes, mirrorEndpointClusters: mirrorEndpointClusters, matchID: matchID}}
}

// SetInfoAPICR populates ID, ApiType, Version and XWso2BasePath of adapterInternalAPI.
func (swagger *AdapterInternalAPI) SetInfoAPICR(api dpv1alpha3.API) {
	swagger.UUID = string(api.ObjectMeta.UID)
	swagger.title = api.Spec.APIName
	swagger.apiType = api.Spec.APIType
	swagger.version = api.Spec.APIVersion
	swagger.xWso2Basepath = api.Spec.BasePath
	swagger.OrganizationID = api.Spec.Organization
	swagger.IsSystemAPI = api.Spec.SystemAPI
	swagger.APIProperties = api.Spec.APIProperties
	httpRouteIDs := []string{}
	for _, route := range api.Spec.Production {
		for _, routeRef := range route.RouteRefs {
			httpRouteIDs = append(httpRouteIDs, getRouteID(api.Namespace, routeRef))
		}
	}
	for _, route := range api.Spec.Sandbox {
		for _, routeRef := range route.RouteRefs {
			httpRouteIDs = append(httpRouteIDs, getRouteID(api.Namespace, routeRef))
		}
	}
	swagger.HTTPRouteIDs = httpRouteIDs
}
