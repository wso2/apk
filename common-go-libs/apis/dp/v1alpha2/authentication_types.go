/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AuthenticationSpec defines the desired state of Authentication
type AuthenticationSpec struct {
	Default   *AuthSpec                       `json:"default,omitempty"`
	Override  *AuthSpec                       `json:"override,omitempty"`
	TargetRef gwapiv1b1.PolicyTargetReference `json:"targetRef,omitempty"`
}

// AuthSpec specification of the authentication service
type AuthSpec struct {
	// Disabled is to disable all authentications
	Disabled *bool `json:"disabled,omitempty"`

	// AuthTypes is to specify the authentication scheme types and details
	AuthTypes *APIAuth `json:"authTypes,omitempty"`
}

// APIAuth Authentication scheme type and details
type APIAuth struct {

	// Oauth2 is to specify the Oauth2 authentication scheme details
	//
	// +optional
	Oauth2 Oauth2Auth `json:"oauth2,omitempty"`

	// APIKey is to specify the APIKey authentication scheme details
	//
	// +optional
	// +nullable
	APIKey []APIKeyAuth `json:"apiKey,omitempty"`

	// JWT is to specify the JWT authentication scheme details
	//
	// +optional
	JWT JWT `json:"jwt,omitempty"`

	// MutualSSL is to specify the features and certificates for mutual SSL
	//
	// +optional
	MutualSSL *MutualSSLConfig `json:"mtls,omitempty"`
}

// MutualSSLConfig scheme type and details
type MutualSSLConfig struct {

	// Disabled is to disable mTLS authentication
	//
	// +kubebuilder:default=false
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// Required indicates whether mutualSSL is mandatory or optional
	// +kubebuilder:validation:Enum=mandatory;optional
	// +kubebuilder:default=optional
	// +optional
	Required string `json:"required"`

	// CertificatesInline is the Inline Certificate entry
	CertificatesInline []*string `json:"certificatesInline,omitempty"`

	// SecretRefs denotes the reference to the Secret that contains the Certificate
	SecretRefs []*RefConfig `json:"secretRefs,omitempty"`

	// ConfigMapRefs denotes the reference to the ConfigMap that contains the Certificate
	ConfigMapRefs []*RefConfig `json:"configMapRefs,omitempty"`
}

// JWT Json Web Token Authentication scheme details
type JWT struct {

	// Disabled is to disable JWT authentication
	//
	// +kubebuilder:default=true
	// +optional
	Disabled *bool `json:"disabled"`

	// Header is the header name used to pass the JWT
	//
	// +kubebuilder:default:=internal-key
	// +optional
	// +kubebuilder:validation:MinLength=1
	Header string `json:"header,omitempty"`

	// SendTokenToUpstream is to specify whether the JWT should be sent to the upstream
	//
	// +optional
	SendTokenToUpstream bool `json:"sendTokenToUpstream,omitempty"`

	// Audience who can invoke a corresponding API
	//
	// +optional
	Audience []string `json:"audience,omitempty"`
}

// Oauth2Auth OAuth2 Authentication scheme details
type Oauth2Auth struct {

	// Required indicates whether OAuth2 is mandatory or optional
	// +kubebuilder:validation:Enum=mandatory;optional
	// +kubebuilder:default=mandatory
	// +optional
	Required string `json:"required,omitempty"`

	// Disabled is to disable OAuth2 authentication
	//
	// +kubebuilder:default=false
	// +optional
	Disabled bool `json:"disabled"`

	// Header is the header name used to pass the OAuth2 token
	//
	// +kubebuilder:default=authorization
	// +optional
	Header string `json:"header,omitempty"`

	// SendTokenToUpstream is to specify whether the OAuth2 token should be sent to the upstream
	//
	// +optional
	SendTokenToUpstream bool `json:"sendTokenToUpstream,omitempty"`
}

// APIKeyAuth APIKey Authentication scheme details
type APIKeyAuth struct {

	//  In is to specify how the APIKey is passed to the request
	//
	// +kubebuilder:validation:Enum=Header;Query
	// +kubebuilder:validation:MinLength=1
	In string `json:"in,omitempty"`

	// Name is the name of the header or query parameter to be used
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// SendTokenToUpstream is to specify whether the APIKey should be sent to the upstream
	//
	// +optional
	SendTokenToUpstream bool `json:"sendTokenToUpstream,omitempty"`
}

// AuthenticationStatus defines the observed state of Authentication
type AuthenticationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// Authentication is the Schema for the authentications API
type Authentication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthenticationSpec   `json:"spec,omitempty"`
	Status AuthenticationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AuthenticationList contains a list of Authentication
type AuthenticationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Authentication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Authentication{}, &AuthenticationList{})
}
