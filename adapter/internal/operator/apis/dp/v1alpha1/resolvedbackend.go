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

import "k8s.io/apimachinery/pkg/types"

// BackendMapping maps read reconciled Backend and resolve properties into ResolvedBackend struct
type BackendMapping map[types.NamespacedName]*ResolvedBackend

// ResolvedBackend holds backend properties
type ResolvedBackend struct {
	Services   []Service
	Protocol   BackendProtocolType
	TLS        ResolvedTLSConfig
	Security   []ResolvedSecurityConfig
	RetryCount int32
}

// ResolvedTLSConfig defines enpoint TLS configurations
type ResolvedTLSConfig struct {
	ResolvedCertificate string
	AllowedSANs         []string
}

// ResolvedSecurityConfig defines enpoint resolved security configurations
type ResolvedSecurityConfig struct {
	Type  string
	Basic ResolvedBasicSecurityConfig
}

// ResolvedBasicSecurityConfig defines resolved basic security configuration
type ResolvedBasicSecurityConfig struct {
	Username string
	Password string
}
