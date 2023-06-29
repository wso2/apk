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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackendProtocolType defines the backend protocol type.
type BackendProtocolType string

const (
	// HTTPProtocol is the http protocol
	HTTPProtocol BackendProtocolType = "http"
	// HTTPSProtocol is the https protocol
	HTTPSProtocol BackendProtocolType = "https"
	// WSProtocol is the ws protocol
	WSProtocol BackendProtocolType = "ws"
	// WSSProtocol is the wss protocol
	WSSProtocol BackendProtocolType = "wss"
)

// BackendConfigs holds different backend configurations
type BackendConfigs struct {
	// +kubebuilder:validation:Enum=http;https;ws;wss
	Protocol BackendProtocolType `json:"protocol"`
	TLS      TLSConfig           `json:"tls,omitempty"`
	Security []SecurityConfig    `json:"security,omitempty"`
}

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// +kubebuilder:validation:MinItems=1
	Services []Service `json:"services,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=http;https;ws;wss
	// +kubebuilder:default=http
	Protocol BackendProtocolType `json:"protocol"`

	// +optional
	TLS *TLSConfig `json:"tls,omitempty"`

	// +optional
	Security []SecurityConfig `json:"security,omitempty"`

	// + optional
	CircuitBreaker *CircuitBreaker `json:"circuitBreaker,omitempty"`
}

// CircuitBreaker defines the circuit breaker configurations
type CircuitBreaker struct {
	// +kubebuilder:default=1024
	MaxConnections uint32 `json:"maxConnections"`
	// +kubebuilder:default=1024
	MaxPendingRequests uint32 `json:"maxPendingRequests"`
	// +kubebuilder:default=1024
	MaxRequests uint32 `json:"maxRequests"`
	// +kubebuilder:default=3
	MaxRetries         uint32 `json:"maxRetries"`
	MaxConnectionPools uint32 `json:"maxConnectionPools,omitempty"`
}

// RetryBudget defines the retry budget configurations
type RetryBudget struct {
	BudgetPercent       uint64 `json:"budgetPercent,omitempty"`
	MinRetryConcurrency uint32 `json:"minRetryConcurrency,omitempty"`
}

// Service holds host and port information for the service
type Service struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}

// TLSConfig defines enpoint TLS configurations
type TLSConfig struct {
	// CertificateInline is the Inline Certificate entry
	CertificateInline *string `json:"certificateInline,omitempty"`
	// SecretRef denotes the reference to the Secret that contains the Certificate
	SecretRef *RefConfig `json:"secretRef,omitempty"`
	// ConfigMapRef denotes the reference to the ConfigMap that contains the Certificate
	ConfigMapRef *RefConfig `json:"configMapRef,omitempty"`
	// AllowedCNs is the list of allowed Subject Alternative Names (SANs)
	AllowedSANs []string `json:"allowedSANs,omitempty"`
}

// RefConfig holds a config for a secret or a configmap
type RefConfig struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// SecurityConfig defines enpoint security configurations
type SecurityConfig struct {
	Type  string              `json:"type,omitempty"`
	Basic BasicSecurityConfig `json:"basic,omitempty"`
}

// BasicSecurityConfig defines basic security configurations
type BasicSecurityConfig struct {
	SecretRef SecretRef `json:"secretRef"`
}

// SecretRef to credentials
type SecretRef struct {
	Name        string `json:"name"`
	UsernameKey string `json:"usernameKey"`
	PasswordKey string `json:"passwordKey"`
}

// BackendStatus defines the observed state of Backend
type BackendStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Backend is the Schema for the backends API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec   `json:"spec,omitempty"`
	Status BackendStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
