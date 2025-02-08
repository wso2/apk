/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package requestconfig

import (
	"strings"

	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// API is a struct that represents an API
type API struct {
	Name                              string                       `json:"name"`                  // Name of the API
	Version                           string                       `json:"version"`               // API version
	Vhost                             string                       `json:"vhost"`                 // Virtual host for the API
	BasePath                          string                       `json:"basePath"`              // Base path for the API
	APIType                           string                       `json:"apiType"`               // Type of the API
	EnvType                           string                       `json:"envType"`               // Environment type (e.g., production, sandbox)
	APILifeCycleState                 string                       `json:"apiLifeCycleState"`     // Lifecycle state of the API
	AuthorizationHeader               string                       `json:"authorizationHeader"`   // Authorization header used by the API
	OrganizationID                    string                       `json:"organizationId"`        // Organization ID for the API
	UUID                              string                       `json:"uuid"`                  // Unique identifier for the API
	Tier                              string                       `json:"tier"`                  // API tier (e.g., Unlimited)
	DisableAuthentication             bool                         `json:"disableAuthentication"` // Whether authentication is disabled
	DisableScopes                     bool                         `json:"disableScopes"`         // Whether scopes are disabled
	Resources                         []*Resource                  `json:"resources"`             // List of resources for the API
	ResourceMap                       map[string]*Resource         `json:"resourceMap"`           // Map of resources for the API
	IsMockedAPI                       bool                         `json:"isMockedApi"`           // Whether the API is mocked
	MutualSSL                         string                       `json:"mutualSSL"`             // Mutual SSL configuration
	TransportSecurity                 bool                         `json:"transportSecurity"`     // Whether transport security is enabled
	ApplicationSecurity               map[string]bool              `json:"applicationSecurity"`   // Application security settings
	BackendJwtConfiguration           *dto.BackendJWTConfiguration `json:"jwtConfigurationDto"`   // JWT configuration DTO
	SystemAPI                         bool                         `json:"systemAPI"`             // Whether the API is a system API
	APIDefinition                     []byte                       `json:"apiDefinition"`
	APIDefinitionPath                 string                       `json:"apiDefinitionPath"`
	Environment                       string                       `json:"environment"`                       // API environment (e.g., development, production)
	SubscriptionValidation            bool                         `json:"subscriptionValidation"`            // Whether subscription validation is enabled
	EndpointSecurity                  []EndpointSecurity           `json:"endpointSecurity"`                  // Endpoint security configurations
	Endpoints                         EndpointCluster              `json:"endpoints"`                         // Endpoint cluster for the API
	AiProvider                        *dto.AIProvider              `json:"aiProvider"`                        // AI provider configuration
	AIModelBasedRoundRobin            *dto.AIModelBasedRoundRobin  `json:"aiModelBasedRoundRobin"`            // AI model-based round robin configuration
	DoSubscriptionAIRLInHeaderReponse bool                         `json:"doSubscriptionAIRLInHeaderReponse"` // Whether to include subscription AIRL in header response
	DoSubscriptionAIRLInBodyReponse   bool                         `json:"doSubscriptionAIRLInBodyReponse"`   // Whether to include subscription AIRL in body response
}

// IsGraphQLAPI checks whether the API is graphql
func (api *API) IsGraphQLAPI() bool {
	return strings.ToLower(api.APIType) == "graphql"
}
