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
	TargetRef gwapiv1b1.PolicyTargetReference `json:"targetRef,omitempty"`
}

// AuthSpec specification of the authentication service
type AuthSpec struct {
	AuthServerType  string         `json:"type,omitempty"`
	ExternalService ExtAuthService `json:"ext,omitempty"`
}

// ExtAuthService external authentication related information
type ExtAuthService struct {
	ServiceRef ServiceRef `json:"serviceRef,omitempty"`
	// Disabled is to disable all authentications
	//
	// +nullable
	Disabled  *bool    `json:"disabled,omitempty"`
	AuthTypes *APIAuth `json:"authTypes,omitempty"`
}

// ServiceRef service using for Authentication
type ServiceRef struct {
	Group string `json:"group,omitempty"`
	Kind  string `json:"kind,omitempty"`
	Name  string `json:"name,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

// APIAuth Authentication scheme type and details
type APIAuth struct {
	JWT            JWTAuth            `json:"jwt,omitempty"`
	APIKey         []APIKeyAuth       `json:"apiKey,omitempty"`
	TestConsoleKey TestConsoleKeyAuth `json:"testConsoleKey,omitempty"`
}

// TestConsoleKeyAuth Test Console Key Authentication scheme details
type TestConsoleKeyAuth struct {
	// Header is the header name used to pass the Test Console Key
	//
	// +kubebuilder:default:=internal-key
	Header              string `json:"header,omitempty"`
	SendTokenToUpstream bool   `json:"sendTokenToUpstream,omitempty"`
}

// JWTAuth JWT Authentication scheme details
type JWTAuth struct {
	// +kubebuilder:default:=false
	Disabled  						bool   `json:"disabled"`
	// Header is the header name used to pass the JWT token
	//
	// +kubebuilder:default:=authorization
	Header              string `json:"header,omitempty"`
	SendTokenToUpstream bool   `json:"sendTokenToUpstream,omitempty"`
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
	Name                string `json:"name,omitempty"`
	SendTokenToUpstream bool   `json:"sendTokenToUpstream,omitempty"`
}

// AuthenticationStatus defines the observed state of Authentication
type AuthenticationStatus struct {
	// Status denotes the state of the Authentication in its lifecycle.
	// Possible values could be Accepted, Invalid, Deploy etc.
	//
	//
	// +kubebuilder:validation:MinLength=4
	Status string `json:"status"`

	// Message represents a user friendly message that explains the
	// current state of the Authentication.
	//
	//
	// +kubebuilder:validation:MinLength=4
	// +optional
	Message string `json:"message"`

	// Accepted represents whether the Authentication is accepted or not.
	//
	//
	Accepted bool `json:"accepted"`

	// TransitionTime represents the last known transition timestamp.
	//
	//
	TransitionTime *metav1.Time `json:"transitionTime"`

	// Events contains a list of events related to the Authentication.
	//
	//
	// +optional
	Events []string `json:"events,omitempty"`
}

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
