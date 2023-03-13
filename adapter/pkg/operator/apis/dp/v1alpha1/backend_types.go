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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// +kubebuilder:validation:MinItems=1
	Services []Service `json:"services,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=http;https;ws;wss
	// +kubebuilder:default=http
	Protocol BackendProtocolType `json:"protocol"`

	// +optional
	TLS *TLSConfig `json:"tls,omitempty"`

	// +optional
	Security []SecurityConfig `json:"security,omitempty"`
}

// Service holds host and port information for the service
type Service struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

// BackendStatus defines the observed state of Backend
type BackendStatus struct {
	// Status denotes the state of the BackendPolicy in its lifecycle.
	// Possible values could be Accepted, Invalid, Deploy etc.
	//
	//
	// +kubebuilder:validation:MinLength=4
	Status string `json:"status"`

	// Message represents a user friendly message that explains the
	// current state of the BackendPolicy.
	//
	//
	// +kubebuilder:validation:MinLength=4
	// +optional
	Message string `json:"message"`

	// Accepted represents whether the BackendPolicy is accepted or not.
	//
	//
	Accepted bool `json:"accepted"`

	// TransitionTime represents the last known transition timestamp.
	//
	//
	TransitionTime *metav1.Time `json:"transitionTime"`

	// Events contains a list of events related to the BackendPolicy.
	//
	//
	// +optional
	Events []string `json:"events,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Backend is the Schema for the backends API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec   `json:"spec,omitempty"`
	Status BackendStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
