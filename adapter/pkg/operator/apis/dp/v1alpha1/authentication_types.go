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
	Default   AuthSpec                        `json:"default,omitempty"`
	Override  AuthSpec                        `json:"override,omitempty"`
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
	Disabled   bool       `json:"disabled,omitempty"`
	AuthTypes  []Auth     `json:"authTypes,omitempty"`
}

// ServiceRef service using for Authentication
type ServiceRef struct {
	Group string `json:"group,omitempty"`
	Kind  string `json:"kind,omitempty"`
	Name  string `json:"name,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

// Auth Authentication scheme type and details
type Auth struct {
	// AuthType is an enum {"jwt", "apikey", "basic", "mtls"}
	AuthType string  `json:"type,omitempty"`
	JWT      JWTAuth `json:"jwt,omitempty"`
}

// JWTAuth JWT Authentication scheme details
type JWTAuth struct {
	AuthorizationHeader string `json:"authorizationHeader,omitempty"`
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
