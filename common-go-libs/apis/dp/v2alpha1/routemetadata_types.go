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

package v2alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// RouteMetadataSpec defines the desired state of RouteMetadata
type RouteMetadataSpec struct {
	API API `json:"api,omitempty"`
}

// API represents the API metadata for the RoutePolicy
type API struct {
	// Name is the name of the API
	Name string `json:"name,omitempty"`
	// Version is the version of the API
	Version string `json:"version,omitempty"`
	// Organization is the organization that owns the API
	Organization string `json:"organization,omitempty"`
	// Type is the type of the API
	Type string `json:"type,omitempty"`
	// Environment denotes the environment of the API.
	//
	// +optional
	// +nullable
	Environment string `json:"environment,omitempty"`
	// EnvType denotes the environment type of the API.
	// +optional
	// +kubebuilder:default:=production
	// +kubebuilder:validation:Enum=production;staging;development
	EnvType string `json:"envType,omitempty"`
	// Context is the context of the API
	Context string `json:"context,omitempty"`
	// APIProperties denotes the custom properties of the API.
	//
	// +optional
	// +nullable
	APIProperties []Property `json:"apiProperties,omitempty"`
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

	// Definition is the API definition.
	// +optional
	Definition string `json:"definition,omitempty"`

	// UUID is the unique identifier for the API.
	// +optional
	// +kubebuilder:validation:MinLength=1
	UUID string `json:"uuid,omitempty"`

	// APICreator is the creator of the API.
	// +optional
	// +kubebuilder:validation:MinLength=1
	APICreator string `json:"apiCreator,omitempty"`

	// APICreatorTenantDomain is the tenant domain of the API creator.
	// +optional
	// +kubebuilder:validation:MinLength=1
	APICreatorTenantDomain string `json:"apiCreatorTenantDomain,omitempty"`
}

// Property holds key value pair of APIProperties
type Property struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RouteMetadata is the Schema for the routemetadata API
type RouteMetadata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteMetadataSpec      `json:"spec,omitempty"`
	Status gwapiv1a2.PolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RouteMetadataList contains a list of RouteMetadata
type RouteMetadataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RouteMetadata `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RouteMetadata{}, &RouteMetadataList{})
}
