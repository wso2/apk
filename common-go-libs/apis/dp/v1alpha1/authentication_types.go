/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AuthenticationSpec defines the desired state of Authentication
type AuthenticationSpec struct {
	Default   *AuthSpec                       `json:"default,omitempty"`
	Override  *AuthSpec                       `json:"override,omitempty"`
	TargetRef gwapiv1b1.NamespacedPolicyTargetReference `json:"targetRef,omitempty"`
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

	// TestConsoleKey is to specify the Test Console Key authentication scheme details
	//
	// +optional
	TestConsoleKey TestConsoleKeyAuth `json:"testConsoleKey,omitempty"`
}

// TestConsoleKeyAuth Test Console Key Authentication scheme details
type TestConsoleKeyAuth struct {
	// Header is the header name used to pass the Test Console Key
	//
	// +kubebuilder:default:=internal-key
	// +optional
	// +kubebuilder:validation:MinLength=1
	Header string `json:"header,omitempty"`

	// SendTokenToUpstream is to specify whether the Test Console Key should be sent to the upstream
	//
	// +optional
	SendTokenToUpstream bool `json:"sendTokenToUpstream,omitempty"`
}

// Oauth2Auth OAuth2 Authentication scheme details
type Oauth2Auth struct {

	// Disabled is to disable OAuth2 authentication
	//
	// +kubebuilder:default:=false
	// +optional
	Disabled bool `json:"disabled"`

	// Header is the header name used to pass the OAuth2 token
	//
	// +kubebuilder:default:=authorization
	// +optional
	Header string `json:"header,omitempty"`

	// SendTokenToUpstream is to specify whether the OAuth2 token should be sent to the upstream
	//
	// +optional
	SendTokenToUpstream bool `json:"sendTokenToUpstream,omitempty"`
}

// APIKeyAuth APIKey Authentication scheme details
type APIKeyAuth struct {
	//	In is to specify how the APIKey is passed to the request
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
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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
