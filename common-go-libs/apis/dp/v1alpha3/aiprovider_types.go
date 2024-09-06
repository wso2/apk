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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AIProviderSpec defines the desired state of AIProvider
type AIProviderSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:MinLength=1
	ProviderName       string          `json:"providerName"`
	ProviderAPIVersion string          `json:"providerAPIVersion"`
	Organization       string          `json:"organization"`
	Model              ValueDetails    `json:"model"`
	RateLimitFields    RateLimitFields `json:"rateLimitFields"`
}

// RateLimitFields defines the Rate Limit fields
type RateLimitFields struct {
	PromptTokens    ValueDetails `json:"promptTokens"`
	CompletionToken ValueDetails `json:"completionToken"`
	TotalToken      ValueDetails `json:"totalToken"`
}

// ValueDetails defines the value details
type ValueDetails struct {
	In    string `json:"in"`
	Value string `json:"value"`
}

// AIProviderStatus defines the observed state of AIProvider
type AIProviderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AIProvider is the Schema for the aiproviders API
type AIProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AIProviderSpec   `json:"spec,omitempty"`
	Status AIProviderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AIProviderList contains a list of AIProvider
type AIProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AIProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AIProvider{}, &AIProviderList{})
}
