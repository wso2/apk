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

import gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

// +k8s:deepcopy-gen=true
// Timeout defines configuration for timeouts related to connections.
type Timeout struct {
	// Timeout settings for TCP.
	//
	// +optional
	TCP *TCPTimeout `json:"tcp,omitempty"`

	// Timeout settings for HTTP.
	//
	// +optional
	HTTP *HTTPTimeout `json:"http,omitempty"`
}

// +k8s:deepcopy-gen=true
type TCPTimeout struct {
	// The timeout for network connection establishment, including TCP and TLS handshakes.
	// Default: 10 seconds.
	//
	// +optional
	ConnectTimeout *gwapiv1.Duration `json:"connectTimeout,omitempty"`
}

// +k8s:deepcopy-gen=true
type HTTPTimeout struct {
	// The idle timeout for an HTTP connection. Idle time is defined as a period in which there are no active requests in the connection.
	// Default: 1 hour.
	//
	// +optional
	ConnectionIdleTimeout *gwapiv1.Duration `json:"connectionIdleTimeout,omitempty"`

	// The maximum duration of an HTTP connection.
	// Default: unlimited.
	//
	// +optional
	MaxConnectionDuration *gwapiv1.Duration `json:"maxConnectionDuration,omitempty"`
}

// +k8s:deepcopy-gen=true
type ClientTimeout struct {
	// Timeout settings for HTTP.
	//
	// +optional
	HTTP *HTTPClientTimeout `json:"http,omitempty"`
}

// +k8s:deepcopy-gen=true
type HTTPClientTimeout struct {
	// The duration envoy waits for the complete request reception. This timer starts upon request
	// initiation and stops when either the last byte of the request is sent upstream or when the response begins.
	//
	// +optional
	RequestReceivedTimeout *gwapiv1.Duration `json:"requestReceivedTimeout,omitempty"`
}
