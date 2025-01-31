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
	"crypto/ecdsa"
	"crypto/x509"
)

// JWTConfiguration represents the JWT configuration
type JWTConfiguration struct {
	Enabled                 bool                   `json:"enabled"`                 // Whether JWT is enabled
	JWTHeader               string                 `json:"jwtHeader"`               // JWT header name
	ConsumerDialectURI      string                 `json:"consumerDialectUri"`      // URI for the consumer dialect
	SignatureAlgorithm      string                 `json:"signatureAlgorithm"`      // Algorithm for signature
	Encoding                string                 `json:"encoding"`                // Encoding type
	TokenIssuerDtoMap       map[string]TokenIssuer `json:"tokenIssuerDtoMap"`       // Map of token issuers
	JwtExcludedClaims       map[string]bool        `json:"jwtExcludedClaims"`       // Excluded claims in JWT
	PublicCert              *x509.Certificate      `json:"publicCert"`              // Public certificate
	PrivateKey              *ecdsa.PrivateKey      `json:"privateKey"`              // Private key for signing JWT
	TTL                     int64                  `json:"ttl"`                     // Time to live for the JWT
	CustomClaims            map[string]ClaimValue  `json:"customClaims"`            // Custom claims
	UseKid                  bool                   `json:"useKid"`                  // Whether to use kid
}
