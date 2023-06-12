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
	"math/rand"

	"github.com/google/uuid"
	"github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	dpv1alpha1 "github.com/wso2/apk/adapter/internal/operator/apis/dp/v1alpha1"
	"github.com/wso2/apk/adapter/internal/operator/utils"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// HTTPRouteParams contains httproute related parameters
type HTTPRouteParams struct {
	AuthSchemes               map[string]dpv1alpha1.Authentication
	ResourceAuthSchemes       map[string]dpv1alpha1.Authentication
	APIPolicies               map[string]dpv1alpha1.APIPolicy
	ResourceAPIPolicies       map[string]dpv1alpha1.APIPolicy
	InterceptorServiceMapping map[string]dpv1alpha1.InterceptorService
	BackendMapping            dpv1alpha1.BackendMapping
	ResourceScopes            map[string]dpv1alpha1.Scope
	RateLimitPolicies         map[string]dpv1alpha1.RateLimitPolicy
	ResourceRateLimitPolicies map[string]dpv1alpha1.RateLimitPolicy
	APIProperties             map[string]dpv1alpha1.APIProperty
}

// SetInfoHTTPRouteCR populates resources and endpoints of adapterInternalAPI. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (swagger *AdapterInternalAPI) SetInfoHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, httpRouteParams HTTPRouteParams) error {
	var resources []*Resource
	var securitySchemes []SecurityScheme
	//TODO(amali) add gateway level securities after gateway crd has implemented
	outputAuthScheme := utils.TieBreaker(utils.GetPtrSlice(maps.Values(httpRouteParams.AuthSchemes)))
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(httpRouteParams.APIPolicies)))
	outputRatelimitPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(httpRouteParams.RateLimitPolicies)))
	outputAPIProperty := utils.TieBreaker(utils.GetPtrSlice(maps.Values(httpRouteParams.APIProperties)))

	var authScheme *dpv1alpha1.Authentication
	if outputAuthScheme != nil {
		authScheme = *outputAuthScheme
	}
	var apiPolicy *dpv1alpha1.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	var apiProperty *dpv1alpha1.APIProperty
	if outputAPIProperty != nil {
		apiProperty = *outputAPIProperty
	}

	var ratelimitPolicy *dpv1alpha1.RateLimitPolicy
	if outputRatelimitPolicy != nil {
		ratelimitPolicy = concatRateLimitPolicies(*outputRatelimitPolicy, nil)
	}

	for _, rule := range httpRoute.Spec.Rules {
		var endPoints []Endpoint
		var policies = OperationPolicies{}
		resourceAuthScheme := authScheme
		resourceAPIPolicy := concatAPIPolicies(apiPolicy, nil)
		var resourceRatelimitPolicy *dpv1alpha1.RateLimitPolicy
		hasPolicies := false
		var scopes []string
		for _, filter := range rule.Filters {
			hasPolicies = true
			switch filter.Type {
			case gwapiv1b1.HTTPRouteFilterURLRewrite:
				policyParameters := make(map[string]interface{})
				policyParameters[constants.RewritePathType] = filter.URLRewrite.Path.Type
				policyParameters[constants.IncludeQueryParams] = true

				switch filter.URLRewrite.Path.Type {
				case gwapiv1b1.FullPathHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = *filter.URLRewrite.Path.ReplaceFullPath
				case gwapiv1b1.PrefixMatchHTTPPathModifier:
					policyParameters[constants.RewritePathResourcePath] = *filter.URLRewrite.Path.ReplacePrefixMatch
				}

				policies.Request = append(policies.Request, Policy{
					PolicyName: string(gwapiv1b1.HTTPRouteFilterURLRewrite),
					Action:     constants.ActionRewritePath,
					Parameters: policyParameters,
				})
			case gwapiv1b1.HTTPRouteFilterExtensionRef:
				if filter.ExtensionRef.Kind == constants.KindAuthentication {
					if ref, found := httpRouteParams.ResourceAuthSchemes[types.NamespacedName{
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
					if ref, found := httpRouteParams.ResourceAPIPolicies[types.NamespacedName{
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
					if ref, found := httpRouteParams.ResourceScopes[types.NamespacedName{
						Name:      string(filter.ExtensionRef.Name),
						Namespace: httpRoute.Namespace,
					}.String()]; found {
						scopes = ref.Spec.Names
					} else {
						return fmt.Errorf("scope: %s has not been resolved in namespace %s", filter.ExtensionRef.Name, httpRoute.Namespace)
					}
				}
				if filter.ExtensionRef.Kind == constants.KindRateLimitPolicy {
					if ref, found := httpRouteParams.ResourceRateLimitPolicies[types.NamespacedName{
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

		addOperationLevelInterceptors(&policies, resourceAPIPolicy, httpRouteParams.InterceptorServiceMapping, httpRouteParams.BackendMapping)

		loggers.LoggerOasparser.Debug("Calculating auths for API ...")
		securities, securityDefinitions, disabledSecurity := getSecurity(concatAuthScheme(resourceAuthScheme), scopes)
		securitySchemes = append(securitySchemes, securityDefinitions...)
		if len(rule.BackendRefs) < 1 {
			return fmt.Errorf("no backendref were provided")
		}
		var securityConfig []EndpointSecurity
		for _, backend := range rule.BackendRefs {
			backendName := types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}
			resolvedBackend, ok := httpRouteParams.BackendMapping[backendName]
			if ok {
				endPoints = append(endPoints, GetEndpoints(backendName, httpRouteParams.BackendMapping)...)
				for _, security := range resolvedBackend.Security {
					switch security.Type {
					case "Basic":
						securityConfig = append(securityConfig, EndpointSecurity{
							Password: string(security.Basic.Password),
							Username: string(security.Basic.Username),
							Type:     string(security.Type),
							Enabled:  true,
						})
					}
				}
			}
		}
		for _, match := range rule.Matches {
			resourcePath := *match.Path.Value
			resource := &Resource{path: resourcePath,
				methods: getAllowedOperations(match.Method, policies, resourceAuthScheme, securities, disabledSecurity,
					parseRateLimitPolicyToInternal(resourceRatelimitPolicy)),
				pathMatchType: *match.Path.Type,
				hasPolicies:   hasPolicies,
				iD:            uuid.New().String(),
			}
			resource.endpoints = &EndpointCluster{
				Endpoints: endPoints,
			}
			resource.endpointSecurity = utils.GetPtrSlice(securityConfig)
			resources = append(resources, resource)
		}
	}
	swagger.xWso2Cors = getCorsConfigFromAPIPolicy(apiPolicy)
	swagger.RateLimitPolicy = parseRateLimitPolicyToInternal(ratelimitPolicy)
	swagger.resources = resources
	apiPolicySelected := concatAPIPolicies(apiPolicy, nil)
	swagger.securityScheme = securitySchemes
    swagger.APIProperty = apiProperty

	loggers.LoggerOasparser.Info("==========================================")
	loggers.LoggerOasparser.Info(swagger.APIProperty )
	loggers.LoggerOasparser.Info("==========================================")

	// Check whether the API has a backend JWT token
	if apiPolicySelected != nil && apiPolicySelected.Spec.Override != nil && apiPolicySelected.Spec.Override.BackendJWTToken != nil {
		loggers.LoggerOasparser.Info("Setting API Level Backend JWT Token Enable/Disable property")
		swagger.backendJWTTokenInfo = parseBackendJWTTokenToInternal(apiPolicySelected.Spec.Override.BackendJWTToken)
		fmt.Println("CLAIMS::::  ", swagger.backendJWTTokenInfo.CustomClaims)
	}
	return nil
}

func parseBackendJWTTokenToInternal(backendJWTToken *dpv1alpha1.BackendJWTToken) *BackendJWTTokenInfo {
	var customClaims []ClaimMapping
	for _, value := range backendJWTToken.CustomClaims {
		claim := value.Claim
		value := value.Value
		claimMapping := ClaimMapping{
			Claim: claim,
			Value: value,
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

// concatAuthSchemes concatinate override and default authentication rules to a one authentication override rule
// folowing the hierarchy described in the https://gateway-api.sigs.k8s.io/references/policy-attachment/#hierarchy
func concatAuthSchemes(schemeUp *dpv1alpha1.Authentication, schemeDown *dpv1alpha1.Authentication) *dpv1alpha1.Authentication {
	if schemeUp == nil && schemeDown == nil {
		return nil
	} else if schemeUp == nil {
		return schemeDown
	} else if schemeDown == nil {
		return schemeUp
	}

	finalAuth := dpv1alpha1.Authentication{}
	var jwtConfigured bool
	var apiKeyConfigured bool

	finalAuth.Spec.Override.ExternalService.Disabled = schemeUp.Spec.Override.ExternalService.Disabled
	for _, auth := range schemeUp.Spec.Override.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
			jwtConfigured = true
		case constants.APIKeyTypeInOAS:
			finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
			apiKeyConfigured = true
		}
	}

	if !finalAuth.Spec.Override.ExternalService.Disabled {
		finalAuth.Spec.Override.ExternalService.Disabled = schemeDown.Spec.Override.ExternalService.Disabled
	}
	for _, auth := range schemeDown.Spec.Override.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			if !jwtConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				jwtConfigured = true
			}
		case constants.APIKeyTypeInOAS:
			if !apiKeyConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				apiKeyConfigured = true
			}
		}
	}

	if !finalAuth.Spec.Override.ExternalService.Disabled {
		finalAuth.Spec.Override.ExternalService.Disabled = schemeDown.Spec.Default.ExternalService.Disabled
	}
	for _, auth := range schemeDown.Spec.Default.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			if !jwtConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				jwtConfigured = true
			}
		case constants.APIKeyTypeInOAS:
			if !apiKeyConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				apiKeyConfigured = true
			}
		}
	}

	if !finalAuth.Spec.Override.ExternalService.Disabled {
		finalAuth.Spec.Override.ExternalService.Disabled = schemeUp.Spec.Default.ExternalService.Disabled
	}
	for _, auth := range schemeUp.Spec.Default.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			if !jwtConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				jwtConfigured = true
			}
		case constants.APIKeyTypeInOAS:
			if !apiKeyConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				apiKeyConfigured = true
			}
		}
	}
	loggers.LoggerOasparser.Debug("Schemes Final auth: %v", &finalAuth)
	return &finalAuth
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

// concatAuthScheme concat override and default auth policies of an authentication CR
// folowing the hierarchy described in the https://gateway-api.sigs.k8s.io/references/policy-attachment/#hierarchy
func concatAuthScheme(scheme *dpv1alpha1.Authentication) *dpv1alpha1.Authentication {
	if scheme == nil || (!scheme.Spec.Default.ExternalService.Disabled && len(scheme.Spec.Default.ExternalService.AuthTypes) < 1) {
		return scheme
	}
	finalAuth := dpv1alpha1.Authentication{}
	var jwtConfigured bool
	var apiKeyConfigured bool
	finalAuth.Spec.Override.ExternalService.Disabled = scheme.Spec.Override.ExternalService.Disabled
	for _, auth := range scheme.Spec.Override.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
			jwtConfigured = true
		case constants.APIKeyTypeInOAS:
			finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
			apiKeyConfigured = true
		}
	}
	if !finalAuth.Spec.Override.ExternalService.Disabled {
		finalAuth.Spec.Override.ExternalService.Disabled = scheme.Spec.Default.ExternalService.Disabled
	}
	for _, auth := range scheme.Spec.Default.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			if !jwtConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				jwtConfigured = true
			}
		case constants.APIKeyTypeInOAS:
			if !apiKeyConfigured {
				finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
				apiKeyConfigured = true
			}
		}
	}
	loggers.LoggerOasparser.Debug("Final auth: %v", &finalAuth)
	return &finalAuth
}

// getSecurity returns security schemes and it's definitions with flag to indicate if security is disabled
// make sure authscheme only has external service override values. (i.e. empty default values)
// tip: use concatScheme method
func getSecurity(authScheme *dpv1alpha1.Authentication, scopes []string) ([]map[string][]string, []SecurityScheme, bool) {
	authSecurities := []map[string][]string{}
	securitySchemes := []SecurityScheme{}
	if authScheme != nil {
		if authScheme.Spec.Override.ExternalService.Disabled {
			loggers.LoggerOasparser.Debug("Disabled security")
			return authSecurities, securitySchemes, true
		}
		for _, auth := range authScheme.Spec.Override.ExternalService.AuthTypes {
			switch auth.AuthType {
			case constants.JWTAuth:
				loggers.LoggerOasparser.Debug("Inside JWT auth")
				securityName := fmt.Sprintf("%s_%v", constants.JWTAuth, rand.Intn(999999999))
				authSecurities = append(authSecurities, map[string][]string{securityName: scopes})
				securitySchemes = append(securitySchemes, SecurityScheme{DefinitionName: securityName, Type: constants.Oauth2TypeInOAS})
			case constants.APIKeyTypeInOAS:
				loggers.LoggerOasparser.Debug("Inside API Key auth")
				securityName := fmt.Sprintf("%s_%v", constants.APIKeyTypeInOAS, rand.Intn(999999999))
				authSecurities = append(authSecurities, map[string][]string{securityName: scopes})
				securitySchemes = append(securitySchemes, SecurityScheme{DefinitionName: securityName,
					Type: constants.APIKeyTypeInOAS, In: constants.APIKeyInHeaderOAS, Name: constants.APIKeyNameWithApim})
			}
		}
	} else {
		loggers.LoggerOasparser.Debug("No auths were provided")
		//todo(amali) remove this default security after scope remodelling is done.
		// apply default security
		securityName := fmt.Sprintf("%s_%v", constants.JWTAuth, rand.Intn(999999999))
		authSecurities = append(authSecurities, map[string][]string{securityName: scopes})
		securitySchemes = append(securitySchemes, SecurityScheme{DefinitionName: securityName, Type: constants.Oauth2TypeInOAS})
	}
	return authSecurities, securitySchemes, false
}

// getAllowedOperations retuns a list of allowed operatons, if httpMethod is not specified then all methods are allowed.
func getAllowedOperations(httpMethod *gwapiv1b1.HTTPMethod, policies OperationPolicies,
	authScheme *dpv1alpha1.Authentication, securities []map[string][]string, disableSecurity bool,
	ratelimitPolicy *RateLimitPolicy) []*Operation {
	if httpMethod != nil {
		return []*Operation{{iD: uuid.New().String(), method: string(*httpMethod), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy}}
	}
	return []*Operation{{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodGet), policies: policies,
		disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPost), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodDelete), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPatch), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPut), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodHead), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodOptions), policies: policies,
			disableSecurity: disableSecurity, security: securities, RateLimitPolicy: ratelimitPolicy}}
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
}
