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

// JWTIssuerSpec defines the desired state of JWTIssuer
type JWTIssuerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the unique name of the JWT Issuer in
	// the Organization defined . "Organization/Name" can
	// be used to uniquely identify an Issuer.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Organization denotes the organization of the JWT Issuer.
	// +kubebuilder:validation:MinLength=1
	Organization string `json:"organization"`
	// Issuer denotes the issuer of the JWT Issuer.
	// +kubebuilder:validation:MinLength=1
	Issuer string `json:"issuer"`
	// ConsumerKeyClaim denotes the claim key of the consumer key.
	// +kubebuilder:validation:MinLength=1
	ConsumerKeyClaim string `json:"consumerKeyClaim"`
	// ScopesClaim denotes the claim key of the scopes.
	// +kubebuilder:validation:MinLength=1
	ScopesClaim string `json:"scopesClaim"`
	// SignatureValidation denotes the signature validation method of jwt
	SignatureValidation *SignatureValidation             `json:"signatureValidation"`
	TargetRef           *gwapiv1b1.PolicyTargetReference `json:"targetRef,omitempty"`
}

// SignatureValidation defines the signature validation method
type SignatureValidation struct {
	JWKS        *JWKS       `json:"jwks,omitempty"`
	Certificate *CERTConfig `json:"certificate,omitempty"`
}

// JWKS defines the JWKS endpoint
type JWKS struct {
	// URL is the URL of the JWKS endpoint
	URL string      `json:"url"`
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

// JWTIssuerStatus defines the observed state of JWTIssuer
type JWTIssuerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// JWTIssuer is the Schema for the jwtissuers API
type JWTIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JWTIssuerSpec   `json:"spec,omitempty"`
	Status JWTIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JWTIssuerList contains a list of JWTIssuer
type JWTIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JWTIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JWTIssuer{}, &JWTIssuerList{})
}
