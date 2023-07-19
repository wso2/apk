/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RateLimitPolicySpec defines the desired state of RateLimitPolicy
type RateLimitPolicySpec struct {
	Default   *RateLimitAPIPolicy             `json:"default,omitempty"`
	Override  *RateLimitAPIPolicy             `json:"override,omitempty"`
	TargetRef gwapiv1b1.PolicyTargetReference `json:"targetRef,omitempty"`
}

// RateLimitAPIPolicy defines the desired state of Policy
type RateLimitAPIPolicy struct {
	// Type of the policy can be either "api" or "application" or "subscription"
	//
	// +kubebuilder:validation:Enum=Api;Application;Subscription;Custom
	Type string `json:"type,omitempty"`

	// API policy
	//
	// +optional
	API APIRateLimitPolicy `json:"api,omitempty"`

	// Custom policy
	//
	// +optional
	Custom CustomRateLimitPolicy `json:"custom,omitempty"`

	// Organization is the organization of the policy
	//
	// +optional
	Organization string `json:"organization,omitempty"`
}

// RateLimit is the rate limit value for the applied policy
type RateLimit struct {
	// RequestPerUnit is the number of requests allowed per unit time
	//
	RequestsPerUnit int `json:"requestsPerUnit,omitempty"`

	// Unit is the unit of the requestsPerUnit
	//
	// +kubebuilder:validation:Enum=Minute;Hour;Day
	Unit string `json:"unit,omitempty"`
}

// APIRateLimitPolicy defines the desired state of APIPolicy
type APIRateLimitPolicy struct {

	// RateLimit is the rate limit for the API
	//
	RateLimit RateLimit `json:"rateLimit,omitempty"`
}

// CustomRateLimitPolicy defines the desired state of CustomPolicy
type CustomRateLimitPolicy struct {
	// RateLimit is the rate limit for the API
	//
	RateLimit RateLimit `json:"rateLimit,omitempty"`

	// Key is the key of the custom policy
	//
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key,omitempty"`

	// Value is the value of the custom policy
	//
	// +optional
	Value string `json:"value,omitempty"`
}

// RateLimitPolicyStatus defines the observed state of RateLimitPolicy
type RateLimitPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RateLimitPolicy is the Schema for the ratelimitpolicies API
type RateLimitPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RateLimitPolicySpec   `json:"spec,omitempty"`
	Status RateLimitPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RateLimitPolicyList contains a list of RateLimitPolicy
type RateLimitPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateLimitPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RateLimitPolicy{}, &RateLimitPolicyList{})
}
