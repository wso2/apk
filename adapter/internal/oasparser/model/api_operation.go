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

// Package model contains the implementation of DTOs to convert OpenAPI/Swagger files
// and create a common model which can represent both types.
package model

import (
	"strings"

	"github.com/google/uuid"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/interceptor"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	dpv1alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
)

// Operation type object holds data about each http method in the REST API.
type Operation struct {
	iD     string
	method string
	//security map of security scheme names -> list of scopes
	scopes           []string
	auth             *Authentication
	tier             string
	disableSecurity  bool
	vendorExtensions map[string]interface{}
	policies         OperationPolicies
	mockedAPIConfig  *api.MockedApiConfig
	rateLimitPolicy  *RateLimitPolicy
	mirrorEndpoints  *EndpointCluster
}

// Authentication holds authentication related configurations
type Authentication struct {
	Disabled bool
	JWT      *JWT
	APIKey   []APIKey
	Oauth2   *Oauth2
}

// JWT holds JWT related configurations
type JWT struct {
	Header              string
	SendTokenToUpstream bool
	Audience            []string
}

// Oauth2 holds Oauth2 related configurations
type Oauth2 struct {
	Header              string
	SendTokenToUpstream bool
}

// APIKey holds API Key related configurations
type APIKey struct {
	In                  string
	Name                string
	SendTokenToUpstream bool
}

// SetAuthentication set authentication configurations
func (operation *Operation) SetAuthentication(authentication *Authentication) {
	operation.auth = authentication
}

// GetAuthentication get authentication configurations
func (operation *Operation) GetAuthentication() *Authentication {
	return operation.auth
}

// GetMethod returns the http method name of the give API operation
func (operation *Operation) GetMethod() string {
	return operation.method
}

// GetPolicies returns if the resouce is secured.
func (operation *Operation) GetPolicies() *OperationPolicies {
	return &operation.policies
}

// GetRateLimitPolicy returns the operation level throttling policy
func (operation *Operation) GetRateLimitPolicy() *RateLimitPolicy {
	return operation.rateLimitPolicy
}

// GetScopes returns the security schemas defined for the http opeartion
func (operation *Operation) GetScopes() []string {
	return operation.scopes
}

// GetTier returns the operation level throttling tier
func (operation *Operation) GetTier() string {
	return operation.tier
}

// GetMockedAPIConfig returns the operation level mocked API implementation configs
func (operation *Operation) GetMockedAPIConfig() *api.MockedApiConfig {
	return operation.mockedAPIConfig
}

// GetVendorExtensions returns vendor extensions which are explicitly defined under
// a given resource.
func (operation *Operation) GetVendorExtensions() map[string]interface{} {
	return operation.vendorExtensions
}

// GetID returns the id of a given resource.
// This is a randomly generated UUID
func (operation *Operation) GetID() string {
	return operation.iD
}

// GetMirrorEndpoints returns the endpoints if a mirror filter has been applied.
func (operation *Operation) GetMirrorEndpoints() *EndpointCluster {
	return operation.mirrorEndpoints
}

// GetCallInterceptorService returns the interceptor configs for a given operation.
func (operation *Operation) GetCallInterceptorService(isIn bool) InterceptEndpoint {
	var policies []Policy
	if isIn {
		policies = operation.policies.Request
	} else {
		policies = operation.policies.Response
	}
	if len(policies) > 0 {
		for _, policy := range policies {
			if strings.EqualFold(constants.ActionInterceptorService, policy.Action) {
				if paramMap, isMap := policy.Parameters.(map[string]interface{}); isMap {
					endpoints, endpointsFound := paramMap[constants.InterceptorEndpoints]
					includesValue, includesFound := paramMap[constants.InterceptorServiceIncludes]
					if endpointsFound {
						endpoints, isEndpoints := endpoints.([]Endpoint)
						if isEndpoints {
							conf := config.ReadConfigs()
							clusterTimeoutV := conf.Envoy.ClusterTimeoutInSeconds
							requestTimeoutV := conf.Envoy.ClusterTimeoutInSeconds
							includesV := &interceptor.RequestInclusions{}
							if includesFound {
								includes, ok := includesValue.([]dpv1alpha1.InterceptorInclusion)
								if ok {
									includesV = GenerateInterceptorIncludes(includes)
								}
							}
							return InterceptEndpoint{
								Enable:          true,
								EndpointCluster: EndpointCluster{Endpoints: endpoints},
								ClusterTimeout:  clusterTimeoutV,
								RequestTimeout:  requestTimeoutV,
								Includes:        includesV,
								Level:           constants.OperationLevelInterceptor,
							}
						}
					}
				}
			}
		}
	}
	return InterceptEndpoint{}
}

// NewOperation Creates and returns operation type object
func NewOperation(method string, security []string, extensions map[string]interface{}) *Operation {
	tier := ResolveThrottlingTier(extensions)
	disableSecurity := ResolveDisableSecurity(extensions)
	id := uuid.New().String()
	return &Operation{id, method, security, nil, tier, disableSecurity, extensions, OperationPolicies{}, &api.MockedApiConfig{}, nil, nil}
}

// NewOperationWithPolicies Creates and returns operation with given method and policies
func NewOperationWithPolicies(method string, policies OperationPolicies) *Operation {
	return &Operation{iD: uuid.New().String(), method: method, policies: policies}
}
