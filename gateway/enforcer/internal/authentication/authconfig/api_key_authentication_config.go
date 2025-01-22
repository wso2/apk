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

// APIKeyAuthenticationConfig represents the configuration for API key authentication
type APIKeyAuthenticationConfig struct {
    In                 string `json:"in"`                 // The location of the API key (e.g., header, query)
    Name               string `json:"name"`              // The name of the API key field
    SendTokenToUpstream bool   `json:"sendTokenToUpstream"` // Whether to send the token to upstream
}
