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

package config

// Config represents the adapter configuration.
// It is created directly from the configuration toml file.
type Config struct {
	CommonController commoncontroller
}

// Config represents the adapter configuration.
// It is created directly from the configuration toml file.
// Note :
//
//	Don't use toml tag for configuration properties as it may affect environment variable based
//	config resolution.

// Common controller related Configurations

type commoncontroller struct {
	// XDSPort    int32    `toml:"xdsPort"`
	// NodeLabels []string `toml:"nodeLabels"`
	Keystore keystore
	Server   server
	Operator operator
	// Trusted Certificates
	Truststore  truststore
	Environment string
}

type keystore struct {
	KeyPath  string
	CertPath string
}

type truststore struct {
	Location string
}

type server struct {
	Label string
}

type operator struct {
	Namespaces []string
}
