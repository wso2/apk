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
package dto

// ExternalProcessingEnvoyMetadata represents the metadata extracted from the external processing request.
type ExternalProcessingEnvoyMetadata struct {
	AuthenticationData             *AuthenticationData `json:"authenticationData"`
	MatchedAPIIdentifier           string              `json:"matchedAPIIdentifier"`
	MatchedResourceIdentifier      string              `json:"matchedResourceIdentifier"`
	MatchedSubscriptionIdentifier  string              `json:"matchedSubscriptionIdentifier"`
	MatchedApplicationIdentifier   string              `json:"matchedApplicationIdentifier"`
	SemanticEmbedding              string              `json:"semanticEmbedding,omitempty"`
	ResponseStatus                 string              `json:"responseStatus,omitempty"`
	AWSBedrockGuardrailPIIEntities map[string]string   `json:"awsBedrockGuardrailPIIEntities,omitempty"`
	PIIMaskingRegexPIIEntities     map[string]string   `json:"piiMaskingRegexPIIEntities,omitempty"`
}

// AuthenticationData represents the authentication data.
type AuthenticationData struct {
	SucessData map[string]*AuthenticationSuccessData `json:"sucessData"`
	FailedData map[string]*AuthenticationFailureData `json:"failedData"`
}

// AuthenticationSuccessData represents the success data of the JWT authentication filter.
type AuthenticationSuccessData struct {
	Issuer string                 `json:"issuer"`
	Claims map[string]interface{} `json:"claims"`
}

// AuthenticationFailureData represents the status of the JWT authentication filyrt.
type AuthenticationFailureData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
