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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// APISpec defines the desired state of API
type APISpec struct {

	// APIName is the unique name of the API
	//can be used to uniquely identify an API.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=60
	// +kubebuilder:validation:Pattern="^[^~!@#;:%^*()+={}|\\<>\"'',&$\\[\\]\\/]*$"
	APIName string `json:"apiName"`

	// APIVersion is the version number of the API.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:Pattern="^[^~!@#;:%^*()+={}|\\<>\"'',&/$\\[\\]\\s+\\/]+$"
	APIVersion string `json:"apiVersion"`

	// IsDefaultVersion indicates whether this API version should be used as a default API
	//
	// +optional
	IsDefaultVersion bool `json:"isDefaultVersion"`

	// DefinitionFileRef contains the
	// definition of the API in a ConfigMap.
	//
	// +optional
	DefinitionFileRef string `json:"definitionFileRef"`

	// DefinitionPath contains the path to expose the API definition.
	//
	// +kubebuilder:default:=/api-definition
	// +kubebuilder:validation:MinLength=1
	DefinitionPath string `json:"definitionPath"`

	// Production contains a list of references to HttpRoutes
	// of type HttpRoute.
	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
	//
	//
	// +optional
	// +nullable
	// +kubebuilder:validation:MaxItems=1
	Production []EnvConfig `json:"production"`

	// Sandbox contains a list of references to HttpRoutes
	// of type HttpRoute.
	// xref: https://github.com/kubernetes-sigs/gateway-api/blob/main/apis/v1beta1/httproute_types.go
	//
	//
	// +optional
	// +nullable
	// +kubebuilder:validation:MaxItems=1
	Sandbox []EnvConfig `json:"sandbox"`

	// APIType denotes the type of the API.
	// Possible values could be REST, GraphQL, Async
	//
	// +kubebuilder:validation:Enum=REST;GraphQL
	APIType string `json:"apiType"`

	// BasePath denotes the basepath of the API.
	// e.g: /pet-store-api/1.0.6
	//
	// +kubectl:validation:MaxLength=232
	// +kubebuilder:validation:Pattern=^[/][a-zA-Z0-9~/_.-]*$
	BasePath string `json:"basePath"`

	// Organization denotes the organization.
	// related to the API
	//
	// +optional
	Organization string `json:"organization"`

	// SystemAPI denotes if it is an internal system API.
	//
	// +optional
	SystemAPI bool `json:"systemAPI"`

	// APIProperties denotes the custom properties of the API.
	//
	// +optional
	// +nullable
	APIProperties []Property `json:"apiProperties,omitempty"`

	// Environment denotes the environment of the API.
	//
	// +optional
	// +nullable
	Environment string `json:"environment,omitempty"`
}

// Property holds key value pair of APIProperties
type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// EnvConfig contains the environment specific configuration
type EnvConfig struct {
	// RouteRefs denotes the environment of the API.
	RouteRefs []string `json:"routeRefs"`
}

// APIStatus defines the observed state of API
type APIStatus struct {
	// DeploymentStatus denotes the deployment status of the API
	//
	// +optional
	DeploymentStatus DeploymentStatus `json:"deploymentStatus"`
}

// DeploymentStatus contains the status of the API deployment
type DeploymentStatus struct {

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

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="API Name",type="string",JSONPath=".spec.apiName"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.apiVersion"
//+kubebuilder:printcolumn:name="BasePath",type="string",JSONPath=".spec.basePath"
//+kubebuilder:printcolumn:name="Organization",type="string",JSONPath=".spec.organization"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

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
