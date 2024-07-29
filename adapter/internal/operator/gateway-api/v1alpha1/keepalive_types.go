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

import (
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// +k8s:deepcopy-gen=true
// TCPKeepalive define the TCP Keepalive configuration.
type TCPKeepalive struct {
	// The total number of unacknowledged probes to send before deciding
	// the connection is dead.
	// Defaults to 9.
	//
	// +optional
	Probes *uint32 `json:"probes,omitempty"`
	// The duration a connection needs to be idle before keep-alive
	// probes start being sent.
	// The duration format is
	// Defaults to `7200s`.
	//
	// +optional
	IdleTime *gwapiv1.Duration `json:"idleTime,omitempty"`
	// The duration between keep-alive probes.
	// Defaults to `75s`.
	//
	// +optional
	Interval *gwapiv1.Duration `json:"interval,omitempty"`
}
