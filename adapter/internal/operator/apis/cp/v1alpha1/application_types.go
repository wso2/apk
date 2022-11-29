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

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	UUID          string            `json:"uuid"`
	Name          string            `json:"name"`
	Owner         string            `json:"owner"`
	Policy        string            `json:"policy"`
	Organization  string            `json:"organization"`
	Attributes    map[string]string `json:"attributes,omitempty"`
	ConsumerKeys  []ConsumerKey     `json:"consumerKeys,omitempty"`
	Subscriptions []Subscription    `json:"subscriptions,omitempty"`
}

// ConsumerKey defines the consumer keys of Application
type ConsumerKey struct {
	Key        string `json:"key"`
	KeyManager string `json:"keyManager"`
}

// Subscription defines a subscription of Application
type Subscription struct {
	UUID               string `json:"uuid"`
	Name               string `json:"name"`
	APIRef             string `json:"apiRef"`
	PolicyID           string `json:"policyId"`
	SubscriptionStatus string `json:"subscriptionStatus"`
	Subscriber         string `json:"subscriber"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
