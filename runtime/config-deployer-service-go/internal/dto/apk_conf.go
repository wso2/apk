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

package dto

// APKConf represents the APK configuration
type APKConf struct {
	Id                     string                  `json:"id,omitempty" yaml:"id,omitempty"`
	Name                   string                  `json:"name" yaml:"name"`
	BasePath               string                  `json:"basePath" yaml:"basePath"`
	Version                string                  `json:"version" yaml:"version"`
	Type                   string                  `json:"type" yaml:"type"`
	DefaultVersion         bool                    `json:"defaultVersion" yaml:"defaultVersion"`
	SubscriptionValidation bool                    `json:"subscriptionValidation" yaml:"subscriptionValidation"`
	Environment            string                  `json:"environment,omitempty" yaml:"environment,omitempty"`
	EndpointConfigurations *EndpointConfigurations `json:"endpointConfigurations,omitempty" yaml:"endpointConfigurations,omitempty"`
	AIProvider             *AIProvider             `json:"aiProvider,omitempty" yaml:"aiProvider,omitempty"`
	Operations             []APKOperations         `json:"operations,omitempty" yaml:"operations,omitempty"`
	//ApiPolicies            *APIOperationPolicies          `json:"apiPolicies,omitempty" yaml:"apiPolicies,omitempty"`
	RateLimit *RateLimit `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	//Authentication       []AuthenticationRequest        `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	AdditionalProperties []APKConf_additionalProperties `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	CorsConfiguration    *CORSConfiguration             `json:"corsConfiguration,omitempty" yaml:"corsConfiguration,omitempty"`
}

// EndpointConfigurations represents endpoint configurations
type EndpointConfigurations struct {
	Production []EndpointConfiguration `json:"production,omitempty" yaml:"production,omitempty"`
	Sandbox    []EndpointConfiguration `json:"sandbox,omitempty" yaml:"sandbox,omitempty"`
}

// EndpointConfiguration represents a single endpoint configuration
type EndpointConfiguration struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// AIProvider represents the AI provider configuration
type AIProvider struct {
	Name       string `json:"name" yaml:"name"`
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
}

// APKOperations represents an api operation
type APKOperations struct {
	Target                 string                  `json:"target,omitempty" yaml:"target,omitempty"`
	Verb                   string                  `json:"verb,omitempty" yaml:"verb,omitempty"`
	Secured                bool                    `json:"secured,omitempty" yaml:"secured,omitempty"`
	EndpointConfigurations *EndpointConfigurations `json:"endpointConfigurations,omitempty" yaml:"endpointConfigurations,omitempty"`
	RateLimit              *RateLimit              `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	Scopes                 []string                `json:"scopes" yaml:"scopes"`
}

// RateLimit represents the rate limit configuration
type RateLimit struct {
	RequestsPerUnit int    `json:"requestsPerUnit" yaml:"requestsPerUnit"`
	Unit            string `json:"unit" yaml:"unit"`
}

type APKConf_additionalProperties struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// CORSConfiguration represents the CORS configuration for an api
type CORSConfiguration struct {
	CORSConfigurationEnabled      bool     `json:"corsConfigurationEnabled" yaml:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty" yaml:"accessControlAllowOrigins,omitempty"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty" yaml:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty" yaml:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty" yaml:"accessControlAllowMethods,omitempty"`
	AccessControlAllowMaxAge      int      `json:"accessControlAllowMaxAge,omitempty" yaml:"accessControlAllowMaxAge,omitempty"`
	AccessControlExposeHeaders    []string `json:"accessControlExposeHeaders,omitempty" yaml:"accessControlExposeHeaders,omitempty"`
}
