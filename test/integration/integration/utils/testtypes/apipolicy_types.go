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
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// APIPolicySpec defines the desired state of APIPolicy
type APIPolicySpec struct {
	// RequestQueryModifier support request query modifications
	//
	//
	// +optional
	Default   *PolicySpec                     `json:"default,omitempty"`
	Override  *PolicySpec                     `json:"override,omitempty"`
	TargetRef gwapiv1b1.PolicyTargetReference `json:"targetRef,omitempty"`
}

// PolicySpec contains API policies
type PolicySpec struct {
	RequestQueryModifier RequestQueryModifier `json:"requestQueryModifier,omitempty"`
	RequestInterceptor   *InterceptorConfig   `json:"requestInterceptor,omitempty"`
	ResponseInterceptor  *InterceptorConfig   `json:"responseInterceptor,omitempty"`
	BackendJWTToken      *BackendJWTToken     `json:"backendJwtToken,omitempty"`
}

// BackendJWTToken holds backend JWT token information
type BackendJWTToken struct {
	Enabled          bool   `json:"enabled,omitempty"`
	Encoding         string `json:"encoding,omitempty"`
	Header           string `json:"header,omitempty"`
	SigningAlgorithm string `json:"signingAlgorithm,omitempty"`
}

// RequestQueryModifier allows to modify request query params
type RequestQueryModifier struct {
	Add       []HTTPQuery `json:"add,omitempty"`
	Remove    []string    `json:"remove,omitempty"`
	RemoveAll string      `json:"removeAll,omitempty"`
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

// InterceptorConfig holds interceptor service information
type InterceptorConfig struct {
	BackendRef BackendReference `json:"backendRef"`

	// +optional
	// +kubebuilder:validation:MaxItems=4
	Includes []InterceptorInclusion `json:"includes,omitempty"`
}

// BackendReference refers to a Backend resource as the interceptor service.
type BackendReference struct {
	// Name is the name of the Backend resource.
	Name string `json:"name"`

	// Namespace is the namespace of the Backend resource.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// HTTPQuery represents an HTTP Header name and value as defined by RFC 7230.
type HTTPQuery struct {
	// Name is the name of the HTTP Header to be matched. Name matching MUST be
	// case insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
	//
	// If multiple entries specify equivalent header names, the first entry with
	// an equivalent name MUST be considered for a match. Subsequent entries
	// with an equivalent header name MUST be ignored. Due to the
	// case-insensitivity of header names, "foo" and "Foo" are considered
	// equivalent.
	Name string `json:"name"`

	// Value is the value of HTTP Header to be matched.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=4096
	Value string `json:"value"`
}

// APIPolicyStatus defines the observed state of APIPolicy
type APIPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// APIPolicy is the Schema for the apipolicies API
type APIPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIPolicySpec   `json:"spec,omitempty"`
	Status APIPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// APIPolicyList contains a list of APIPolicy
type APIPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIPolicy `json:"items"`
}
