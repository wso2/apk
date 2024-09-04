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
 */

package v1alpha2

// ResolvedBackend holds backend properties
type ResolvedBackend struct {
	Backend        Backend
	Services       []Service
	Protocol       BackendProtocolType
	TLS            ResolvedTLSConfig
	Security       ResolvedSecurityConfig
	CircuitBreaker *CircuitBreaker
	Timeout        *Timeout
	Retry          *RetryConfig
	BasePath       string `json:"basePath"`
	HealthCheck    *HealthCheck
}

// ResolvedTLSConfig defines enpoint TLS configurations
type ResolvedTLSConfig struct {
	ResolvedCertificate string
	AllowedSANs         []string
}

// ResolvedSecurityConfig defines enpoint resolved security configurations
type ResolvedSecurityConfig struct {
	Type   string
	Basic  ResolvedBasicSecurityConfig
	APIKey ResolvedAPIKeySecurityConfig
}

// ResolvedBasicSecurityConfig defines resolved basic security configuration
type ResolvedBasicSecurityConfig struct {
	Username string
	Password string
}

// ResolvedAPIKeySecurityConfig defines resolved API key security configuration
type ResolvedAPIKeySecurityConfig struct {
	In    string
	Name  string
	Value string
}
