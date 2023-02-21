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
	AuthSchemes            map[string]dpv1alpha1.Authentication
	ResourceAuthSchemes    map[string]dpv1alpha1.Authentication
	APIPolicies            map[string]dpv1alpha1.APIPolicy
	ResourceAPIPolicies    map[string]dpv1alpha1.APIPolicy
	BackendPropertyMapping dpv1alpha1.BackendPropertyMapping
	ResourceScopes         map[string]dpv1alpha1.Scope
}

// SetInfoHTTPRouteCR populates resources and endpoints of mgwSwagger. httpRoute.Spec.Rules.Matches
// are used to create resources and httpRoute.Spec.Rules.BackendRefs are used to create EndpointClusters.
func (swagger *MgwSwagger) SetInfoHTTPRouteCR(httpRoute *gwapiv1b1.HTTPRoute, httpRouteParams HTTPRouteParams) error {
	var resources []*Resource
	var securitySchemes []SecurityScheme
	//TODO(amali) add gateway level securities after gateway crd has implemented
	outputAuthScheme := utils.TieBreaker(utils.GetPtrSlice(maps.Values(httpRouteParams.AuthSchemes)))
	outputAPIPolicy := utils.TieBreaker(utils.GetPtrSlice(maps.Values(httpRouteParams.APIPolicies)))

	var authScheme *dpv1alpha1.Authentication
	if outputAuthScheme != nil {
		authScheme = *outputAuthScheme
	}
	var apiPolicy *dpv1alpha1.APIPolicy
	if outputAPIPolicy != nil {
		apiPolicy = *outputAPIPolicy
	}
	for _, rule := range httpRoute.Spec.Rules {
		var endPoints []Endpoint
		var policies = OperationPolicies{}
		resourceAuthScheme := authScheme
		resourceAPIPolicy := apiPolicy
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
						resourceAPIPolicy = concatAPIPolicies(resourceAPIPolicy, &ref)
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
						scopes = append(scopes, ref.Spec.Names...)
					} else {
						return fmt.Errorf("scope: %s has not been resolved in namespace %s", filter.ExtensionRef.Name, httpRoute.Namespace)
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
		loggers.LoggerOasparser.Debug("Calculating auths for API ...")
		securities, securityDefinitions, disabledSecurity := getSecurity(concatAuthScheme(resourceAuthScheme), scopes)
		securitySchemes = append(securitySchemes, securityDefinitions...)
		if len(rule.BackendRefs) < 1 {
			return fmt.Errorf("no backendref were provided")
		}
		for _, backend := range rule.BackendRefs {
			backendProperties := httpRouteParams.BackendPropertyMapping[types.NamespacedName{
				Name:      string(backend.Name),
				Namespace: utils.GetNamespace(backend.Namespace, httpRoute.Namespace),
			}]
			endPoints = append(endPoints,
				Endpoint{Host: backendProperties.ResolvedHostname,
					URLType:     string(backendProperties.Protocol),
					Port:        uint32(*backend.Port),
					Certificate: []byte(backendProperties.TLS.CertificateInline),
				})
		}

		for _, match := range rule.Matches {
			resourcePath := *match.Path.Value
			resource := &Resource{path: resourcePath,
				methods:       getAllowedOperations(match.Method, policies, resourceAuthScheme, securities, disabledSecurity),
				pathMatchType: *match.Path.Type,
				hasPolicies:   hasPolicies,
				iD:            uuid.New().String(),
			}
			resource.endpoints = &EndpointCluster{
				Endpoints: endPoints,
			}
			resources = append(resources, resource)
		}
	}
	swagger.resources = resources
	swagger.securityScheme = securitySchemes
	return nil
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

	finalAuth.Spec.Override.ExternalService.Disabled = schemeUp.Spec.Override.ExternalService.Disabled
	for _, auth := range schemeUp.Spec.Override.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
			jwtConfigured = true
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
		}
	}
	return &finalAuth
}

// concatAuthScheme concat override and default auth policies of an authentication CR
// folowing the hierarchy described in the https://gateway-api.sigs.k8s.io/references/policy-attachment/#hierarchy
func concatAuthScheme(scheme *dpv1alpha1.Authentication) *dpv1alpha1.Authentication {
	if scheme == nil || (!scheme.Spec.Default.ExternalService.Disabled && len(scheme.Spec.Default.ExternalService.AuthTypes) < 1) {
		return scheme
	}
	finalAuth := dpv1alpha1.Authentication{}
	var jwtConfigured bool
	finalAuth.Spec.Override.ExternalService.Disabled = scheme.Spec.Override.ExternalService.Disabled
	for _, auth := range scheme.Spec.Override.ExternalService.AuthTypes {
		switch auth.AuthType {
		case constants.JWTAuth:
			finalAuth.Spec.Override.ExternalService.AuthTypes = append(finalAuth.Spec.Override.ExternalService.AuthTypes, auth)
			jwtConfigured = true
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
		}
	}
	return &finalAuth
}

// concatAPIPolicies concatinate override and default authentication rules to a one authentication override rule
// folowing the hierarchy described in the https://gateway-api.sigs.k8s.io/references/policy-attachment/#hierarchy
// Following code would follow below logic.
// | API override policies | Resource override policies | Resource default policies | API default policies |
// |            1          |              1             |              0            |           0          | API override policies
// |            1          |              0             |              1            |           0          | API override policies
// |            0          |              1             |              0            |           1          | Resource override policies
// |            0          |              0             |              1            |           1          | Resource default policies
func concatAPIPolicies(schemeUp *dpv1alpha1.APIPolicy, schemeDown *dpv1alpha1.APIPolicy) *dpv1alpha1.APIPolicy {
	if schemeUp == nil && schemeDown == nil {
		return nil
	} else if schemeUp == nil {
		return schemeDown
	} else if schemeDown == nil {
		return schemeUp
	}

	finalAPIPolicy := dpv1alpha1.APIPolicy{}

	// api level override policies - must apply
	finalAPIPolicy.Spec.Override.RequestQueryModifier = schemeUp.Spec.Override.RequestQueryModifier

	// resource level override policies - must apply
	// api level override RequestQueryModifier.Add/remove + resource level override RequestQueryModifier.Add/remove
	if len(finalAPIPolicy.Spec.Override.RequestQueryModifier.Add) < 1 {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.Add = schemeDown.Spec.Override.RequestQueryModifier.Add
	}
	if len(finalAPIPolicy.Spec.Override.RequestQueryModifier.Remove) < 1 {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.Remove = schemeDown.Spec.Override.RequestQueryModifier.Remove
	}
	if finalAPIPolicy.Spec.Override.RequestQueryModifier.RemoveAll == "" {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.RemoveAll = schemeDown.Spec.Override.RequestQueryModifier.RemoveAll
	}

	// resource level default policies if above are empty
	if len(finalAPIPolicy.Spec.Override.RequestQueryModifier.Add) < 1 {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.Add = schemeDown.Spec.Default.RequestQueryModifier.Add
	}
	if len(finalAPIPolicy.Spec.Override.RequestQueryModifier.Remove) < 1 {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.Remove = schemeDown.Spec.Default.RequestQueryModifier.Remove
	}
	if finalAPIPolicy.Spec.Override.RequestQueryModifier.RemoveAll == "" {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.RemoveAll = schemeDown.Spec.Default.RequestQueryModifier.RemoveAll
	}

	// API level default policies if above are empty
	if len(finalAPIPolicy.Spec.Override.RequestQueryModifier.Add) < 1 {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.Add = schemeUp.Spec.Default.RequestQueryModifier.Add
	}
	if len(finalAPIPolicy.Spec.Override.RequestQueryModifier.Remove) < 1 {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.Remove = schemeUp.Spec.Default.RequestQueryModifier.Remove
	}
	if finalAPIPolicy.Spec.Override.RequestQueryModifier.RemoveAll == "" {
		finalAPIPolicy.Spec.Override.RequestQueryModifier.RemoveAll = schemeUp.Spec.Default.RequestQueryModifier.RemoveAll
	}

	return &finalAPIPolicy
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
				securityName := fmt.Sprintf("%s_%v", constants.JWTAuth, rand.Intn(999999999))
				authSecurities = append(authSecurities, map[string][]string{securityName: scopes})
				securitySchemes = append(securitySchemes, SecurityScheme{DefinitionName: securityName, Type: constants.Oauth2TypeInOAS})
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
	authScheme *dpv1alpha1.Authentication, securities []map[string][]string, disableSecurity bool) []*Operation {
	if httpMethod != nil {
		return []*Operation{{iD: uuid.New().String(), method: string(*httpMethod), policies: policies,
			disableSecurity: disableSecurity, security: securities}}
	}
	return []*Operation{{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodGet), policies: policies,
		disableSecurity: disableSecurity, security: securities},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPost), policies: policies,
			disableSecurity: disableSecurity, security: securities},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodDelete), policies: policies,
			disableSecurity: disableSecurity, security: securities},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPatch), policies: policies,
			disableSecurity: disableSecurity, security: securities},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodPut), policies: policies,
			disableSecurity: disableSecurity, security: securities},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodHead), policies: policies,
			disableSecurity: disableSecurity, security: securities},
		{iD: uuid.New().String(), method: string(gwapiv1b1.HTTPMethodOptions), policies: policies,
			disableSecurity: disableSecurity, security: securities}}
}

// SetInfoAPICR populates ID, ApiType, Version and XWso2BasePath of mgwSwagger.
func (swagger *MgwSwagger) SetInfoAPICR(api dpv1alpha1.API) {
	swagger.UUID = string(api.ObjectMeta.UID)
	swagger.title = api.Spec.APIDisplayName
	swagger.apiType = api.Spec.APIType
	swagger.version = api.Spec.APIVersion
	swagger.xWso2Basepath = api.Spec.Context
	swagger.OrganizationID = api.Spec.Organization
	swagger.IsSystemAPI = api.Spec.SystemAPI
}
