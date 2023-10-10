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
	Default   *PolicySpec                     `json:"default,omitempty"`
	Override  *PolicySpec                     `json:"override,omitempty"`
	TargetRef gwapiv1b1.PolicyTargetReference `json:"targetRef,omitempty"`
}

// PolicySpec contains API policies
type PolicySpec struct {
	// RequestInterceptors referenced to intercetor services to be applied
	// to the request flow.
	//
	// +optional
	// +nullable
	// +kubebuilder:validation:MaxItems=1
	RequestInterceptors []InterceptorReference `json:"requestInterceptors,omitempty"`

	// ResponseInterceptors referenced to intercetor services to be applied
	// to the response flow.
	//
	// +optional
	// +nullable
	// +kubebuilder:validation:MaxItems=1
	ResponseInterceptors []InterceptorReference `json:"responseInterceptors,omitempty"`

	// BackendJWTPolicy holds reference to backendJWT policy configurations
	BackendJWTPolicy *BackendJWTToken `json:"backendJwtPolicy,omitempty"`

	// CORS policy to be applied to the API.
	CORSPolicy *CORSPolicy `json:"cORSPolicy,omitempty"`
}

// BackendJWTToken holds backend JWT token information
type BackendJWTToken struct {
	// Name holds the name of the BackendJWT resource.
	Name string `json:"name,omitempty"`
}

// CORSPolicy holds CORS policy information
type CORSPolicy struct {
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	//
	// +optional
	AccessControlAllowCredentials bool `json:"accessControlAllowCredentials,omitempty"`

	// AccessControlAllowHeaders indicates which headers can be used
	// during the actual request.
	//
	// +optional
	AccessControlAllowHeaders []string `json:"accessControlAllowHeaders,omitempty"`

	// AccessControlAllowMethods indicates which methods can be used
	// during the actual request.
	//
	// +optional
	AccessControlAllowMethods []string `json:"accessControlAllowMethods,omitempty"`

	// AccessControlAllowOrigins indicates which origins can be used
	// during the actual request.
	//
	// +optional
	AccessControlAllowOrigins []string `json:"accessControlAllowOrigins,omitempty"`

	// AccessControlExposeHeaders indicates which headers can be exposed
	// as part of the response by listing their names.
	//
	// +optional
	AccessControlExposeHeaders []string `json:"accessControlExposeHeaders,omitempty"`

	// AccessControlMaxAge indicates how long the results of a preflight request
	// can be cached in a preflight result cache.
	//
	// +optional
	AccessControlMaxAge *int `json:"accessControlMaxAge,omitempty"`
}

// InterceptorReference holds InterceptorService reference using name and namespace
type InterceptorReference struct {
	// Name is the referced CR's name of InterceptorService resource.
	Name string `json:"name"`
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

func init() {
	SchemeBuilder.Register(&APIPolicy{}, &APIPolicyList{})
}
