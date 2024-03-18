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
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GQLRouteSpec defines the desired state of GQLRoute
type GQLRouteSpec struct {
	v1.CommonRouteSpec `json:",inline"`

	// Hostnames defines a set of hostname that should match against the HTTP Host
	// header to select a GQLRoute used to process the request.
	// +optional
	// +kubebuilder:validation:MaxItems=16
	Hostnames []v1.Hostname `json:"hostnames,omitempty"`

	// BackendRefs defines the backend(s) where matching requests should be
	// sent.
	BackendRefs []v1.HTTPBackendRef `json:"backendRefs,omitempty"`

	// Rules are a list of GraphQL resources, filters and actions.
	//
	// +optional
	// +kubebuilder:validation:MaxItems=16
	Rules []GQLRouteRules `json:"rules,omitempty"`
}

// GQLRouteRules defines semantics for matching an GraphQL request based on
// conditions (matches), processing it (filters), and forwarding the request to
// an API object (backendRefs).
type GQLRouteRules struct {

	// Matches define conditions used for matching the rule against incoming
	// graphQL requests. Each match is independent, i.e. this rule will be matched
	// if **any** one of the matches is satisfied.
	Matches []GQLRouteMatch `json:"matches,omitempty"`

	// Filters define the filters that are applied to requests that match
	// this rule.
	//
	// +kubebuilder:validation:MaxItems=16
	Filters []GQLRouteFilter `json:"filters,omitempty"`
}

// GQLRouteFilter defines the filter to be applied to a request.
type GQLRouteFilter struct {
	// ExtensionRef is an optional, implementation-specific extension to the
	// "filter" behavior.  For example, resource "myroutefilter" in group
	// "networking.example.net"). ExtensionRef MUST NOT be used for core and
	// extended filters.
	//
	// Support: Implementation-specific
	//
	// +optional
	ExtensionRef *v1.LocalObjectReference `json:"extensionRef,omitempty"`
}

// GQLRouteMatch defines the predicate used to match requests to a given
// action.
type GQLRouteMatch struct {
	// Type specifies GQL typematcher.
	// When specified, this route will be matched only if the request has the
	// specified method.
	//
	// Support: Extended
	//
	// +optional
	// +kubebuilder:validation:Default=QUERY
	Type *GQLType `json:"type,omitempty"`

	// Path specifies a GQL request resource matcher.
	Path *string `json:"path,omitempty"`
}

// GQLType describes how to select a GQL request by matching the GQL Type.
// The value is expected in upper case.
//
// Note that values may be added to this enum, implementations
// must ensure that unknown values will not cause a crash.
//
// Unknown values here must result in the implementation setting the
// Accepted Condition for the Route to `status: False`, with a
// Reason of `UnsupportedValue`.
//
// +kubebuilder:validation:Enum=QUERY;MUTATION
type GQLType string

// GQLRouteStatus defines the observed state of GQLRoute
type GQLRouteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GQLRoute is the Schema for the gqlroutes API
type GQLRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GQLRouteSpec   `json:"spec,omitempty"`
	Status GQLRouteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GQLRouteList contains a list of GQLRoute
type GQLRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GQLRoute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GQLRoute{}, &GQLRouteList{})
}
