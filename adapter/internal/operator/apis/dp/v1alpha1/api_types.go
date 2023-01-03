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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Important: Run "make" to regenerate code after modifying this file

// APISpec defines the desired state of API
type APISpec struct {

	// APIDisplayName is the unique name of the API in
	// the namespace defined. "Namespace/APIDisplayName" can
	// be used to uniquely identify an API.
	//
	//
	APIDisplayName string `json:"apiDisplayName"`

	// APIVersion is the version number of the API.
	//
	//
	APIVersion string `json:"apiVersion"`

	// DefinitionFileRef contains the OpenAPI 3 or Swagger
	// definition of the API in a ConfigMap.
	//
	//
	// +optional
	DefinitionFileRef string `json:"definitionFileRef"`

	// ProdHTTPRouteRefs contains a list of references to HttpRoutes
	// of type HttpRoute.
	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
	//
	//
	// +optional
	ProdHTTPRouteRef string `json:"prodHTTPRouteRef"`

	// SandHTTPRouteRef contains a list of references to HttpRoutes
	// of type HttpRoute.
	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
	//
	//
	// +optional
	SandHTTPRouteRef string `json:"sandHTTPRouteRef"`

	// APIType denotes the type of the API.
	// Possible values could be REST, GraphQL, Async
	//
	APIType string `json:"apiType"`

	// Context denotes the context of the API.
	// e.g: /pet-store-api/1.0.6
	//
	Context string `json:"context"`

	// Organization denotes the organization
	// related to the API
	//
	// +optional
	Organization string `json:"organization"`
}

// APIStatus defines the observed state of API
type APIStatus struct {

	// Status denotes the state of the API in its lifecycle.
	// Possible values could be Accepted, Invalid, Deploy etc.
	//
	//
	Status string `json:"status"`

	// Message represents a user friendly message that explains the
	// current state of the API.
	//
	//
	// +optional
	Message string `json:"message"`

	// Accepted represents whether the API is accepted or not.
	//
	//
	Accepted bool `json:"accepted"`

	// TransitionTime represents the last known transition timestamp.
	//
	//
	TransitionTime *metav1.Time `json:"transitionTime"`

	// Events contains a list of events related to the API.
	//
	//
	// +optional
	Events []string `json:"events,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// API is the Schema for the apis API
type API struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APISpec   `json:"spec,omitempty"`
	Status APIStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// APIList contains a list of API
type APIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []API `json:"items"`
}

func init() {
	SchemeBuilder.Register(&API{}, &APIList{})
}
