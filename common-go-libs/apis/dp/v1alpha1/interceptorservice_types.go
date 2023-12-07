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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackendReference refers to a Backend resource as the interceptor service.
type BackendReference struct {
	// Name is the name of the Backend resource.
	//
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// InterceptorInclusion defines the type of data which can be included in the interceptor request/response path
type InterceptorInclusion string

const (
	// InterceptorInclusionRequestHeaders is the type to include request headers
	InterceptorInclusionRequestHeaders InterceptorInclusion = "request_headers"
	// InterceptorInclusionRequestBody is the type to include request body
	InterceptorInclusionRequestBody InterceptorInclusion = "request_body"
	// InterceptorInclusionRequestTrailers is the type to include request trailers
	InterceptorInclusionRequestTrailers InterceptorInclusion = "request_trailers"
	// InterceptorInclusionResponseHeaders is the type to include response headers
	InterceptorInclusionResponseHeaders InterceptorInclusion = "response_headers"
	// InterceptorInclusionResponseBody is the type to include response body
	InterceptorInclusionResponseBody InterceptorInclusion = "response_body"
	// InterceptorInclusionResponseTrailers is the type to include response trailers
	InterceptorInclusionResponseTrailers InterceptorInclusion = "response_trailers"
	// InterceptorInclusionInvocationContext is the type to include invocation context
	InterceptorInclusionInvocationContext InterceptorInclusion = "invocation_context"
)

// InterceptorServiceSpec defines the desired state of InterceptorService
type InterceptorServiceSpec struct {
	BackendRef BackendReference `json:"backendRef"`

	// Includes defines the types of data which should be included when calling the interceptor service
	//
	// +optional
	// +kubebuilder:validation:MaxItems=4
	// +nullable
	Includes []InterceptorInclusion `json:"includes,omitempty"`
}

// InterceptorServiceStatus defines the observed state of InterceptorService
type InterceptorServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// InterceptorService is the Schema for the interceptorservices API
type InterceptorService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InterceptorServiceSpec   `json:"spec,omitempty"`
	Status InterceptorServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InterceptorServiceList contains a list of InterceptorService
type InterceptorServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InterceptorService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InterceptorService{}, &InterceptorServiceList{})
}
