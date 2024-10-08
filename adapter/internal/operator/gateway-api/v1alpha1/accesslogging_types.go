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
type ProxyAccessLog struct {
	// Disable disables access logging for managed proxies if set to true.
	Disable bool `json:"disable,omitempty"`
	// Settings defines accesslog settings for managed proxies.
	// If unspecified, will send default format to stdout.
	// +optional
	Settings []ProxyAccessLogSetting `json:"settings,omitempty"`
}

// +k8s:deepcopy-gen=true
type ProxyAccessLogSetting struct {
	// Format defines the format of accesslog.
	Format ProxyAccessLogFormat `json:"format"`
	// Sinks defines the sinks of accesslog.
	// +kubebuilder:validation:MinItems=1
	Sinks []ProxyAccessLogSink `json:"sinks"`
}

// +k8s:deepcopy-gen=true
type ProxyAccessLogFormatType string

const (
	// ProxyAccessLogFormatTypeText defines the text accesslog format.
	ProxyAccessLogFormatTypeText ProxyAccessLogFormatType = "Text"
	// ProxyAccessLogFormatTypeJSON defines the JSON accesslog format.
	ProxyAccessLogFormatTypeJSON ProxyAccessLogFormatType = "JSON"
	// TODO: support format type "mix" in the future.
)

// +k8s:deepcopy-gen=true
// ProxyAccessLogFormat defines the format of accesslog.
// By default accesslogs are written to standard output.
// +union
//
// +kubebuilder:validation:XValidation:rule="self.type == 'Text' ? has(self.text) : !has(self.text)",message="If AccessLogFormat type is Text, text field needs to be set."
// +kubebuilder:validation:XValidation:rule="self.type == 'JSON' ? has(self.json) : !has(self.json)",message="If AccessLogFormat type is JSON, json field needs to be set."
type ProxyAccessLogFormat struct {
	// Type defines the type of accesslog format.
	// +kubebuilder:validation:Enum=Text;JSON
	// +unionDiscriminator
	Type ProxyAccessLogFormatType `json:"type,omitempty"`
	// Text defines the text accesslog format, following Envoy accesslog formatting,
	// It's required when the format type is "Text".
	// Envoy [command operators](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#command-operators) may be used in the format.
	// The [format string documentation](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#config-access-log-format-strings) provides more information.
	// +optional
	Text *string `json:"text,omitempty"`
	// JSON is additional attributes that describe the specific event occurrence.
	// Structured format for the envoy access logs. Envoy [command operators](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#command-operators)
	// can be used as values for fields within the Struct.
	// It's required when the format type is "JSON".
	// +optional
	JSON map[string]string `json:"json,omitempty"`
}

// +k8s:deepcopy-gen=true
type ProxyAccessLogSinkType string

const (
	// ProxyAccessLogSinkTypeFile defines the file accesslog sink.
	ProxyAccessLogSinkTypeFile ProxyAccessLogSinkType = "File"
	// ProxyAccessLogSinkTypeOpenTelemetry defines the OpenTelemetry accesslog sink.
	// When the provider is Kubernetes, EnvoyGateway always sends `k8s.namespace.name`
	// and `k8s.pod.name` as additional attributes.
	ProxyAccessLogSinkTypeOpenTelemetry ProxyAccessLogSinkType = "OpenTelemetry"
)

// +k8s:deepcopy-gen=true
// ProxyAccessLogSink defines the sink of accesslog.
// +union
//
// +kubebuilder:validation:XValidation:rule="self.type == 'File' ? has(self.file) : !has(self.file)",message="If AccessLogSink type is File, file field needs to be set."
// +kubebuilder:validation:XValidation:rule="self.type == 'OpenTelemetry' ? has(self.openTelemetry) : !has(self.openTelemetry)",message="If AccessLogSink type is OpenTelemetry, openTelemetry field needs to be set."
type ProxyAccessLogSink struct {
	// Type defines the type of accesslog sink.
	// +kubebuilder:validation:Enum=File;OpenTelemetry
	// +unionDiscriminator
	Type ProxyAccessLogSinkType `json:"type,omitempty"`
	// File defines the file accesslog sink.
	// +optional
	File *FileEnvoyProxyAccessLog `json:"file,omitempty"`
	// OpenTelemetry defines the OpenTelemetry accesslog sink.
	// +optional
	OpenTelemetry *OpenTelemetryEnvoyProxyAccessLog `json:"openTelemetry,omitempty"`
}

// +k8s:deepcopy-gen=true
type FileEnvoyProxyAccessLog struct {
	// Path defines the file path used to expose envoy access log(e.g. /dev/stdout).
	// +kubebuilder:validation:MinLength=1
	Path string `json:"path,omitempty"`
}

// +k8s:deepcopy-gen=true
// TODO: consider reuse ExtensionService?
type OpenTelemetryEnvoyProxyAccessLog struct {
	// Host define the extension service hostname.
	Host string `json:"host"`
	// Port defines the port the extension service is exposed on.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=4317
	Port int32 `json:"port,omitempty"`
	// Resources is a set of labels that describe the source of a log entry, including envoy node info.
	// It's recommended to follow [semantic conventions](https://opentelemetry.io/docs/reference/specification/resource/semantic_conventions/).
	// +optional
	Resources map[string]string `json:"resources,omitempty"`

	// TODO: support more OpenTelemetry accesslog options(e.g. TLS, auth etc.) in the future.
}
