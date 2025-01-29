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

import (
	"crypto/x509"
)

// TokenIssuer represents the token issuer
type TokenIssuer struct {
	Issuer                     string                   `json:"issuer"`                     // Issuer of the token
	DisableDefaultClaimMapping bool                     `json:"disableDefaultClaimMapping"` // Whether to disable default claim mapping
	ClaimConfigurations        map[string]*ClaimMapping `json:"claimConfigurations"`        // Claim mappings
	JwksConfigurationDTO       JWKSConfiguration        `json:"jwksConfigurationDTO"`       // JWKS configuration
	Certificate                *x509.Certificate        `json:"certificate"`                // Optional certificate
	ConsumerKeyClaim           string                   `json:"consumerKeyClaim"`           // Claim for the consumer key
	ScopesClaim                string                   `json:"scopesClaim"`                // Claim for scopes
	Audience                   string                   `json:"audience"`                   // Audience of the token
}
