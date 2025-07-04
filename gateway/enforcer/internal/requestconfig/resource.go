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
	"fmt"

	auth "github.com/wso2/apk/gateway/enforcer/internal/authentication/authconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// HTTPMethods represents the HTTP methods (GET, POST, etc.
type HTTPMethods string

const (
	// GET represents the GET HTTP method
	GET HTTPMethods = "GET"
	// POST represents the POST HTTP method
	POST HTTPMethods = "POST"
	// PUT represents the PUT HTTP method
	PUT HTTPMethods = "PUT"
	// DELETE represents the DELETE HTTP method
	DELETE HTTPMethods = "DELETE"
	// PATCH represents the PATCH HTTP method
	PATCH HTTPMethods = "PATCH"
	// OPTIONS represents the OPTIONS HTTP method
	OPTIONS HTTPMethods = "OPTIONS"
	// HEAD represents the HEAD HTTP method
	HEAD HTTPMethods = "HEAD"
)

// Resource represents the configuration for a resource
type Resource struct {
	Path                    string                                 `json:"path"`                    // The path of the resource
	MatchID                 string                                 `json:"matchID"`                 // The match ID for the resource
	Method                  HTTPMethods                            `json:"method"`                  // The HTTP method (GET, POST, etc.)
	Tier                    string                                 `json:"tier"`                    // The tier of the resource (default is "Unlimited")
	Endpoints               *EndpointCluster                       `json:"endpoints"`               // Endpoint cluster for the resource
	EndpointSecurity        []*EndpointSecurity                    `json:"endpointSecurity"`        // Endpoint security configurations
	PolicyConfig            PolicyConfig                           `json:"policyConfig"`            // Policy configurations for the resource
	AuthenticationConfig    *auth.AuthenticationConfig             `json:"authenticationConfig"`    // Authentication configuration
	Scopes                  []string                               `json:"scopes"`                  // Scopes for the resource
	AIModelBasedRoundRobin  *dto.AIModelBasedRoundRobin            `json:"aiModelBasedRoundRobin"`  // AI model-based round robin configuration
	RouteMetadataAttributes *dto.ExternalProcessingEnvoyAttributes `json:"routeMetadataAttributes"` // Route metadata attributes
	RequestInBuiltPolicies  []dto.InBuiltPolicy                    `json:"requestPolicies"`         // List of request policies for the resource
	ResponseInBuiltPolicies []dto.InBuiltPolicy                    `json:"responsePolicies"`        // List of response policies for the resource
}

// GetResourceIdentifier returns the identifier for the resource
func (r *Resource) GetResourceIdentifier() string {
	return fmt.Sprintf("%s_%s", r.Method, r.Path)
}
