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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackendJWTSpec defines the desired state of BackendJWT
type BackendJWTSpec struct {
	// Encoding of the JWT token
	//
	// +optional
	// +kubebuilder:default=Base64
	// +kubebuilder:validation:Enum=Base64;Base64url
	Encoding string `json:"encoding,omitempty"`

	// Header of the JWT token
	//
	// +optional
	// +kubebuilder:default=X-JWT-Assertion
	// +kubebuilder:validation:MinLength=1
	Header string `json:"header,omitempty"`

	// Signing algorithm of the JWT token
	//
	// +optional
	// +kubebuilder:default=SHA256withRSA
	// +kubeBuilder:validation:Enum=SHA256withRSA;SHA384withRSA;SHA512withRSA;SHA256withECDSA;SHA384withECDSA;SHA512withECDSA;SHA256withHMAC;SHA384withHMAC;SHA512withHMAC
	SigningAlgorithm string `json:"signingAlgorithm,omitempty"`

	// TokenTTL time to live for the backend JWT token in seconds
	//
	// +optional
	// +kubebuilder:default=3600
	TokenTTL uint32 `json:"tokenTTL,omitempty"`

	// CustomClaims holds custom claims that needs to be added to the jwt
	//
	// +optional
	// +nullable
	CustomClaims []CustomClaim `json:"customClaims,omitempty"`
}

// CustomClaim holds custom claim information
type CustomClaim struct {
	// Claim name
	//
	// +kubebuilder:validation:MinLength=1
	Claim string `json:"claim,omitempty"`

	// Claim value
	//
	// +kubebuilder:validation:MinLength=1
	Value string `json:"value,omitempty"`

	// Claim type
	//
	// +kubebuilder:default=string
	// +kubebuilder:validation:Enum=string;int;float;bool;long;date
	Type string `json:"type"`
}

// BackendJWTStatus defines the observed state of BackendJWT
type BackendJWTStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BackendJWT is the Schema for the backendjwts API
type BackendJWT struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendJWTSpec   `json:"spec,omitempty"`
	Status BackendJWTStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BackendJWTList contains a list of BackendJWT
type BackendJWTList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackendJWT `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackendJWT{}, &BackendJWTList{})
}
