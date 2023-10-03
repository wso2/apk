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

package v1alpha1

import "k8s.io/apimachinery/pkg/types"

// JWTIssuerMapping maps read reconciled Backend and resolve properties into ResolvedJWTIssuer struct
type JWTIssuerMapping map[types.NamespacedName]*ResolvedJWTIssuer

// ResolvedJWTIssuer holds the resolved properties of JWTIssuer
type ResolvedJWTIssuer struct {
	Name                string
	Organization        string
	Issuer              string
	ConsumerKeyClaim    string
	ScopesClaim         string
	SignatureValidation ResolvedSignatureValidation
	ClaimMappings       map[string]string
	Environments        []string
}

// ResolvedSignatureValidation holds the resolved properties of SignatureValidation
type ResolvedSignatureValidation struct {
	JWKS        *ResolvedJWKS
	Certificate *ResolvedTLSConfig
}

// ResolvedJWKS holds the resolved properties of JWKS
type ResolvedJWKS struct {
	URL string
	TLS *ResolvedTLSConfig
}
