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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TokenIssuerSpec defines the desired state of TokenIssuer
type TokenIssuerSpec struct {
	// Name is the unique name of the Token Issuer in
	// the Organization defined . "Organization/Name" can
	// be used to uniquely identify an Issuer.
	//
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Organization denotes the organization of the Token Issuer.
	//
	// +kubebuilder:validation:MinLength=1
	Organization string `json:"organization"`

	// Issuer denotes the issuer of the Token Issuer.
	//
	// +kubebuilder:validation:MinLength=1
	Issuer string `json:"issuer"`

	// ConsumerKeyClaim denotes the claim key of the consumer key.
	//
	// +kubebuilder:validation:MinLength=1
	ConsumerKeyClaim string `json:"consumerKeyClaim"`

	// ScopesClaim denotes the claim key of the scopes.
	//
	// +kubebuilder:validation:MinLength=1
	ScopesClaim string `json:"scopesClaim"`

	// SignatureValidation denotes the signature validation method of jwt
	SignatureValidation *SignatureValidation `json:"signatureValidation"`

	// ClaimMappings denotes the claim mappings of the jwt
	ClaimMappings *[]ClaimMapping `json:"claimMappings,omitempty"`

	// TargetRef denotes the reference to the which gateway it applies to
	TargetRef *gwapiv1b1.NamespacedPolicyTargetReference `json:"targetRef,omitempty"`
}

// ClaimMapping defines the reference configuration
type ClaimMapping struct {
	// RemoteClaim denotes the remote claim
	RemoteClaim string `json:"remoteClaim"`
	// LocalClaim denotes the local claim
	LocalClaim string `json:"localClaim"`
}

// SignatureValidation defines the signature validation method
type SignatureValidation struct {
	// JWKS denotes the JWKS endpoint information
	JWKS *JWKS `json:"jwks,omitempty"`
	// Certificate denotes the certificate information
	Certificate *CERTConfig `json:"certificate,omitempty"`
}

// JWKS defines the JWKS endpoint
type JWKS struct {
	// URL is the URL of the JWKS endpoint
	URL string `json:"url"`
	// TLS denotes the TLS configuration of the JWKS endpoint
	TLS *CERTConfig `json:"tls,omitempty"`
}

// CERTConfig defines the certificate configuration
type CERTConfig struct {
	// CertificateInline is the Inline Certificate entry
	CertificateInline *string `json:"certificateInline,omitempty"`
	// SecretRef denotes the reference to the Secret that contains the Certificate
	SecretRef *RefConfig `json:"secretRef,omitempty"`
	// ConfigMapRef denotes the reference to the ConfigMap that contains the Certificate
	ConfigMapRef *RefConfig `json:"configMapRef,omitempty"`
}

// TokenIssuerStatus defines the observed state of TokenIssuer
type TokenIssuerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TokenIssuer is the Schema for the tokenIssuer API
type TokenIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TokenIssuerSpec   `json:"spec,omitempty"`
	Status TokenIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TokenIssuerList contains a list of TokenIssuer
type TokenIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TokenIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TokenIssuer{}, &TokenIssuerList{})
}
