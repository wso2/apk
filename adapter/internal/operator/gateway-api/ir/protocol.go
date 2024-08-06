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

package ir

type AppProtocol string

const (
	// GRPC declares that the port carries gRPC traffic.
	GRPC AppProtocol = "GRPC"
	// GRPCWeb declares that the port carries gRPC traffic.
	GRPCWeb AppProtocol = "GRPC-Web"
	// HTTP declares that the port carries HTTP/1.1 traffic.
	// Note that HTTP/1.0 or earlier may not be supported by the proxy.
	HTTP AppProtocol = "HTTP"
	// HTTP2 declares that the port carries HTTP/2 traffic.
	HTTP2 AppProtocol = "HTTP2"
	// HTTPS declares that the port carries HTTPS traffic.
	HTTPS AppProtocol = "HTTPS"
	// TCP declares the port uses TCP.
	// This is the default protocol for a service port.
	TCP AppProtocol = "TCP"
	// UDP declares that the port uses UDP.
	UDP AppProtocol = "UDP"
)
