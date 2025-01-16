/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha4

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

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// Services holds hosts and ports
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=1
	Services []Service `json:"services,omitempty"`

	// Protocol defines the backend protocol
	//
	// +optional
	// +kubebuilder:validation:Enum=http;https;ws;wss
	// +kubebuilder:default=http
	Protocol BackendProtocolType `json:"protocol"`

	// BasePath defines the base path of the backend
	// +optional
	BasePath string `json:"basePath"`

	// TLS defines the TLS configurations of the backend
	TLS *TLSConfig `json:"tls,omitempty"`

	// Security defines the security configurations of the backend
	Security *SecurityConfig `json:"security,omitempty"`

	// CircuitBreaker defines the circuit breaker configurations
	CircuitBreaker *CircuitBreaker `json:"circuitBreaker,omitempty"`

	// Timeout configuration for the backend
	Timeout *Timeout `json:"timeout,omitempty"`

	// Retry configuration for the backend
	Retry *RetryConfig `json:"retry,omitempty"`

	// HealthCheck configuration for the backend tcp health check
	HealthCheck *HealthCheck `json:"healthCheck,omitempty"`

	// SupportedModels is the list of supported models
	SupportedModels []string `json:"supportedModels,omitempty"`
}

// HealthCheck defines the health check configurations
type HealthCheck struct {

	// Timeout is the time to wait for a health check response.
	// If the timeout is reached the health check attempt will be considered a failure.
	//
	// +kubebuilder:default=1
	// +optional
	Timeout uint32 `json:"timeout,omitempty"`

	// Interval is the time between health check attempts in seconds.
	//
	// +kubebuilder:default=30
	// +optional
	Interval uint32 `json:"interval,omitempty"`

	// UnhealthyThreshold is the number of consecutive health check failures required
	// before a backend is marked unhealthy.
	//
	// +kubebuilder:default=2
	// +optional
	UnhealthyThreshold uint32 `json:"unhealthyThreshold,omitempty"`

	// HealthyThreshold is the number of healthy health checks required before a host is marked healthy.
	// Note that during startup, only a single successful health check is required to mark a host healthy.
	//
	// +kubebuilder:default=2
	// +optional
	HealthyThreshold uint32 `json:"healthyThreshold,omitempty"`
}

// Timeout defines the timeout configurations
type Timeout struct {
	// UpstreamResponseTimeout spans between the point at which the entire downstream request (i.e. end-of-stream) has been processed and
	// when the upstream response has been completely processed.
	// A value of 0 will disable the route’s timeout.
	//
	// +kubebuilder:default=15
	UpstreamResponseTimeout uint32 `json:"upstreamResponseTimeout"`

	// DownstreamRequestIdleTimeout bounds the amount of time the request's stream may be idle.
	// A value of 0 will completely disable the route's idle timeout.
	//
	// +kubebuilder:default=300
	// +optional
	DownstreamRequestIdleTimeout uint32 `json:"downstreamRequestIdleTimeout"`
}

// CircuitBreaker defines the circuit breaker configurations
type CircuitBreaker struct {

	// MaxConnections is the maximum number of connections that will make to the upstream cluster.
	//
	// +kubebuilder:default=1024
	// +optional
	MaxConnections uint32 `json:"maxConnections"`

	// MaxPendingRequests is the maximum number of pending requests that will allow to the upstream cluster.
	//
	// +kubebuilder:default=1024
	// +optional
	MaxPendingRequests uint32 `json:"maxPendingRequests"`

	// MaxRequests is the maximum number of parallel requests that will make to the upstream cluster.
	//
	// +kubebuilder:default=1024
	// +optional
	MaxRequests uint32 `json:"maxRequests"`

	// MaxRetries is the maximum number of parallel retries that will allow to the upstream cluster.
	//
	// +kubebuilder:default=3
	// +optional
	MaxRetries uint32 `json:"maxRetries"`

	// MaxConnectionPools is the maximum number of parallel connection pools that will allow to the upstream cluster.
	// If not specified, the default is unlimited.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxConnectionPools uint32 `json:"maxConnectionPools"`
}

// RetryConfig defines retry configurations
type RetryConfig struct {

	// Count defines the number of retries.
	// If exceeded, TooEarly(425 response code) response will be sent to the client.
	//
	// +kubebuilder:default=1
	Count uint32 `json:"count"`

	// BaseIntervalMillis is exponential retry back off and it defines the base interval between retries in milliseconds.
	// maximum interval is 10 times of the BaseIntervalMillis
	//
	// +kubebuilder:default=25
	// +kubebuilder:validation:Minimum=1
	// +optional
	BaseIntervalMillis uint32 `json:"baseIntervalMillis"`

	// StatusCodes defines the list of status codes to retry
	//
	// +optional
	StatusCodes []uint32 `json:"statusCodes,omitempty"`
}

// Service holds host and port information for the service
type Service struct {
	// Host is the hostname of the service
	//
	// +kubebuilder:validation:MinLength=1
	Host string `json:"host"`

	// Port of the service
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
	//
	// +optional
	AllowedSANs []string `json:"allowedSANs,omitempty"`
}

// RefConfig holds a config for a secret or a configmap
type RefConfig struct {
	// Name of the secret or configmap
	//
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Key of the secret or configmap
	//
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`
}

// SecurityConfig defines enpoint security configurations
type SecurityConfig struct {
	// Basic security configuration
	Basic *BasicSecurityConfig `json:"basic,omitempty"`
	// APIKey security configuration
	APIKey *APIKeySecurityConfig `json:"apiKey,omitempty"`
}

// APIKeySecurityConfig defines APIKey security configurations
type APIKeySecurityConfig struct {
	//  In is to specify how the APIKey is passed to the request
	//
	// +kubebuilder:validation:Enum=Header;Query
	// +kubebuilder:validation:MinLength=1
	In string `json:"in,omitempty"`

	// Name is the name of the header or query parameter to be used
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// ValueRef to value
	ValueFrom ValueRef `json:"valueFrom"`
}

// ValueRef to value
type ValueRef struct {
	// Name of the secret
	//
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Field Key of the APIKey
	//
	// +kubebuilder:validation:MinLength=1
	ValueKey string `json:"valueKey"`
}

// BasicSecurityConfig defines basic security configurations
type BasicSecurityConfig struct {
	// SecretRef to credentials
	SecretRef SecretRef `json:"secretRef"`
}

// SecretRef to credentials
type SecretRef struct {
	// Name of the secret
	//
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Username Key value
	//
	// +kubebuilder:validation:MinLength=1
	UsernameKey string `json:"usernameKey"`

	// Password Key of the secret
	//
	// +kubebuilder:validation:MinLength=1
	PasswordKey string `json:"passwordKey"`
}

// BackendStatus defines the observed state of Backend
type BackendStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

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
