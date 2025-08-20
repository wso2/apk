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

import "time"

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
	Truststore        truststore
	Environment       string
	Redis             redis
	Sts               sts
	WebServer         webServer
	InternalAPIServer internalAPIServer
	ControlPlane      controlplane
	Metrics           Metrics
	Database          database
	DeployResourcesWithClusterRoleBindings bool
}
type controlplane struct {
	Enabled              bool
	Host                 string
	EventPort            int
	RestPort             int
	RetryInterval        time.Duration
	Persistence          persistence
	SkipSSLVerification  bool
	EnableAPIPropagation bool
	APIsRestPath         string
}
type persistence struct {
	Type string
}
type internalAPIServer struct {
	Port int64
}
type keystore struct {
	KeyPath  string
	CertPath string
}

type truststore struct {
	Location string
}

type server struct {
	Label      string
	ServerName string
}

type operator struct {
	Namespaces []string
}

type redis struct {
	Host                string
	Port                string
	Username            string
	Password            string
	UserCertPath        string
	UserKeyPath         string
	CACertPath          string
	TLSEnabled          bool
	RevokedTokenChannel string
}

type sts struct {
	AuthKeyPath   string
	AuthKeyHeader string
}

type webServer struct {
	Port int64
}

// Metrics defines the configuration for metrics collection.
type Metrics struct {
	Enabled bool
	Type    string
	Port    int32
}

type database struct {
	Enabled     bool
	Name        string
	Username    string
	Password    string
	Host        string
	Port        int
	PoolOptions dbPool
}

type dbPool struct {
	// PoolMaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU()
	PoolMaxConns int

	// PoolMinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low
	// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance
	// to create new connections.
	PoolMinConns int

	// PoolMaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	PoolMaxConnLifetime string

	// PoolMaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	PoolMaxConnIdleTime string

	// HealthCheckPeriod is the duration between checks of the health of idle connections.
	PoolHealthCheckPeriod string

	// PoolMaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.
	// This helps prevent all connections from being closed at the exact same time, starving the pool.
	PoolMaxConnLifetimeJitter string
}
