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
 
package authconfig

// AuthenticationConfig represents the configuration for authentication
type AuthenticationConfig struct {
	JWTAuthenticationConfig     *JWTAuthenticationConfig     `json:"jwtAuthenticationConfig"`     // JWT Authentication configuration
	APIKeyAuthenticationConfigs []*APIKeyAuthenticationConfig `json:"apiKeyAuthenticationConfigs"` // List of API key authentication configurations
	Oauth2AuthenticationConfig  *Oauth2AuthenticationConfig  `json:"oauth2AuthenticationConfig"`  // OAuth2 authentication configuration
	Disabled                    bool                         `json:"disabled"`                    // Whether authentication is disabled
}
