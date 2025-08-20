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
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// RoutePolicySpec defines the desired state of RoutePolicy
type RoutePolicySpec struct {
	RequestMediation  []*Mediation `json:"requestMediation,omitempty"`
	ResponseMediation []*Mediation `json:"responseMediation,omitempty"`
}

// Mediation represents a policy mediation configuration
// It can be used for both request and response mediation
type Mediation struct {
	PolicyName    string      `json:"policyName"`
	PolicyID      string      ` json:"policyID"`
	PolicyVersion string      `json:"policyVersion,omitempty"`
	Parameters    []*Parameter `json:"parameters,omitempty"`
}

// Parameter represents a key-value or key-valueFrom pair for policy parameters
type Parameter struct {
	// Key is the name of the parameter
	// It is used to identify the parameter in the policy
	Key      string                        `json:"key"`
	// Value is the value of the parameter
	Value    string                        `json:"value,omitempty"` // JSON-encoded value for the parameter
	// ValueRef is used to reference a value from another resource
	// It can be used to reference a value from a ConfigMap, Secret, or other resources
	ValueRef *gwapiv1.LocalObjectReference `json:"valueRef,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RoutePolicy is the Schema for the routepolicies API
type RoutePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoutePolicySpec   `json:"spec,omitempty"`
	Status gwapiv1a2.PolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RoutePolicyList contains a list of RoutePolicy
type RoutePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoutePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RoutePolicy{}, &RoutePolicyList{})
}
