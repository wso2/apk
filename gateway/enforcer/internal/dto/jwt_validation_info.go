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

// JWTValidationInfo represents the JWT validation info
type JWTValidationInfo struct {
	Valid             bool                   `json:"valid"`             // Valid
	ExpiryTime        int64                  `json:"expiryTime"`        // Expiry time
	IssuedTime        int64                  `json:"issuedTime"`        // Issued time
	JTI               string                 `json:"jti"`               // JTI
	ValidationCode    int                    `json:"validationCode"`    // Validation code
	ValidationMessage string                 `json:"validationMessage"` // Validation message
	Issuer            string                 `json:"issuer"`            // Issuer
	ClientID          string                 `json:"clientId"`          // Client ID
	Subject           string                 `json:"subject"`           // Subject
	Audiences         []string               `json:"audiences"`         // Audiences
	Scopes            []string               `json:"scopes"`            // Scopes
	Claims            map[string]interface{} `json:"claims"`            // Claims
}
