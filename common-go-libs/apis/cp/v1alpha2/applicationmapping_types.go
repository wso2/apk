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

// ApplicationMappingSpec defines the desired state of ApplicationMapping
type ApplicationMappingSpec struct {
	ApplicationRef  string `json:"applicationRef"`
	SubscriptionRef string `json:"subscriptionRef"`
}

// ApplicationMappingStatus defines the observed state of ApplicationMapping
type ApplicationMappingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApplicationMapping is the Schema for the applicationmappings API
type ApplicationMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationMappingSpec   `json:"spec,omitempty"`
	Status ApplicationMappingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationMappingList contains a list of ApplicationMapping
type ApplicationMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationMapping{}, &ApplicationMappingList{})
}
