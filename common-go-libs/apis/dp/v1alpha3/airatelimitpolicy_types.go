/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// AIRateLimitPolicySpec defines the desired state of AIRateLimitPolicy
type AIRateLimitPolicySpec struct {
	Override  *AIRateLimit                              `json:"override,omitempty"`
	Default   *AIRateLimit                              `json:"default,omitempty"`
	TargetRef gwapiv1b1.NamespacedPolicyTargetReference `json:"targetRef,omitempty"`
}

// AIRateLimit defines the AI ratelimit configuration
type AIRateLimit struct {
	Organization string        `json:"organization,omitempty"`
	TokenCount   *TokenCount   `json:"tokenCount,omitempty"`
	RequestCount *RequestCount `json:"requestCount,omitempty"`
}

// TokenCount defines the Token based ratelimit configuration
type TokenCount struct {
	// Unit is the unit of the requestsPerUnit
	//
	// +kubebuilder:validation:Enum=Minute;Hour;Day
	Unit string `json:"unit,omitempty"`

	// RequestTokenCount specifies the maximum number of tokens allowed
	// in AI requests within a given unit of time. This value limits the
	// token count sent by the client to the AI service over the defined period.
	//
	// +kubebuilder:validation:Minimum=1
	RequestTokenCount uint32 `json:"requestTokenCount,omitempty"`

	// ResponseTokenCount specifies the maximum number of tokens allowed
	// in AI responses within a given unit of time. This value limits the
	// token count received by the client from the AI service over the defined period.
	//
	// +kubebuilder:validation:Minimum=1
	ResponseTokenCount uint32 `json:"responseTokenCount,omitempty"`

	// TotalTokenCount represents the maximum allowable total token count
	// for both AI requests and responses within a specified unit of time.
	// This value sets the limit for the number of tokens exchanged between
	// the client and AI service during the defined period.
	//
	// +kubebuilder:validation:Minimum=1
	TotalTokenCount uint32 `json:"totalTokenCount,omitempty"`
}

// AIRateLimitPolicyStatus defines the observed state of AIRateLimitPolicy
type AIRateLimitPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AIRateLimitPolicy is the Schema for the airatelimitpolicies API
type AIRateLimitPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AIRateLimitPolicySpec   `json:"spec,omitempty"`
	Status AIRateLimitPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AIRateLimitPolicyList contains a list of AIRateLimitPolicy
type AIRateLimitPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AIRateLimitPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AIRateLimitPolicy{}, &AIRateLimitPolicyList{})
}
