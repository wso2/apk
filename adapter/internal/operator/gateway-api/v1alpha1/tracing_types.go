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
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package v1alpha1

// +k8s:deepcopy-gen=true
type ProxyTracing struct {
	// SamplingRate controls the rate at which traffic will be
	// selected for tracing if no prior sampling decision has been made.
	// Defaults to 100, valid values [0-100]. 100 indicates 100% sampling.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=100
	// +optional
	SamplingRate *uint32 `json:"samplingRate,omitempty"`
	// CustomTags defines the custom tags to add to each span.
	// If provider is kubernetes, pod name and namespace are added by default.
	CustomTags map[string]CustomTag `json:"customTags,omitempty"`
	// Provider defines the tracing provider.
	// Only OpenTelemetry is supported currently.
	Provider TracingProvider `json:"provider"`
}

// +k8s:deepcopy-gen=true
type TracingProviderType string

const (
	TracingProviderTypeOpenTelemetry TracingProviderType = "OpenTelemetry"
)

// +k8s:deepcopy-gen=true
type TracingProvider struct {
	// Type defines the tracing provider type.
	// EG currently only supports OpenTelemetry.
	// +kubebuilder:validation:Enum=OpenTelemetry
	// +kubebuilder:default=OpenTelemetry
	Type TracingProviderType `json:"type"`
	// Host define the provider service hostname.
	Host string `json:"host"`
	// Port defines the port the provider service is exposed on.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=4317
	Port int32 `json:"port,omitempty"`
}

// +k8s:deepcopy-gen=true
type CustomTagType string

const (
	// CustomTagTypeLiteral adds hard-coded value to each span.
	CustomTagTypeLiteral CustomTagType = "Literal"
	// CustomTagTypeEnvironment adds value from environment variable to each span.
	CustomTagTypeEnvironment CustomTagType = "Environment"
	// CustomTagTypeRequestHeader adds value from request header to each span.
	CustomTagTypeRequestHeader CustomTagType = "RequestHeader"
)

// +k8s:deepcopy-gen=true
type CustomTag struct {
	// Type defines the type of custom tag.
	// +kubebuilder:validation:Enum=Literal;Environment;RequestHeader
	// +unionDiscriminator
	// +kubebuilder:default=Literal
	Type CustomTagType `json:"type"`
	// Literal adds hard-coded value to each span.
	// It's required when the type is "Literal".
	Literal *LiteralCustomTag `json:"literal,omitempty"`
	// Environment adds value from environment variable to each span.
	// It's required when the type is "Environment".
	Environment *EnvironmentCustomTag `json:"environment,omitempty"`
	// RequestHeader adds value from request header to each span.
	// It's required when the type is "RequestHeader".
	RequestHeader *RequestHeaderCustomTag `json:"requestHeader,omitempty"`

	// TODO: add support for Metadata tags in the future.
	// EG currently doesn't support metadata for route or cluster.
}

// +k8s:deepcopy-gen=true
// LiteralCustomTag adds hard-coded value to each span.
type LiteralCustomTag struct {
	// Value defines the hard-coded value to add to each span.
	Value string `json:"value"`
}

// +k8s:deepcopy-gen=true
// EnvironmentCustomTag adds value from environment variable to each span.
type EnvironmentCustomTag struct {
	// Name defines the name of the environment variable which to extract the value from.
	Name string `json:"name"`
	// DefaultValue defines the default value to use if the environment variable is not set.
	// +optional
	DefaultValue *string `json:"defaultValue,omitempty"`
}

// +k8s:deepcopy-gen=true
// RequestHeaderCustomTag adds value from request header to each span.
type RequestHeaderCustomTag struct {
	// Name defines the name of the request header which to extract the value from.
	Name string `json:"name"`
	// DefaultValue defines the default value to use if the request header is not set.
	// +optional
	DefaultValue *string `json:"defaultValue,omitempty"`
}
