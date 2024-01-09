/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package server

// TokenIssuer holds the properties of TokenIssuer
type TokenIssuer struct {
	Name                string                      `json:"name"`
	Organization        string                      `json:"organization"`
	Issuer              string                      `json:"issuer"`
	ConsumerKeyClaim    string                      `json:"consumerKeyClaim"`
	ScopesClaim         string                      `json:"scopesClaim"`
	SignatureValidation ResolvedSignatureValidation `json:"signatureValidation"`
	ClaimMappings       map[string]string           `json:"claimMappings"`
	Environments        []string                    `json:"environments"`
}

// ResolvedSignatureValidation holds the resolved properties of SignatureValidation
type ResolvedSignatureValidation struct {
	JWKS        *ResolvedJWKS      `json:"jwks"`
	Certificate *ResolvedTLSConfig `json:"certificate"`
}

// ResolvedJWKS holds the resolved properties of JWKS
type ResolvedJWKS struct {
	URL string             `json:"url"`
	TLS *ResolvedTLSConfig `json:"tls"`
}

// ResolvedTLSConfig defines enpoint TLS configurations
type ResolvedTLSConfig struct {
	ResolvedCertificate string   `json:"resolvedCertificate"`
	AllowedSANs         []string `json:"allowedSANs"`
}

// TokenIssuserList contains a list of TokenIssuser
type TokenIssuserList struct {
	List []TokenIssuer `json:"list"`
}
