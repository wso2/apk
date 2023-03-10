/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
	ManagementServer managementServer
	Database         database
	BackOffice       backOffice
}

type managementServer struct {
	XDSPort          int32      `toml:"xdsPort"`
	NodeLabels       []string   `toml:"nodeLabels"`
	GRPCPort         uint       `toml:"gRPCPort"`
	NotificationPort uint       `toml:"notificationPort"`
	Keystore         keystore   `toml:"keystore"`
	Truststore       truststore `toml:"truststore"`
}

type keystore struct {
	KeyPath  string
	CertPath string
}
type truststore struct {
	Location string
}

type backOffice struct {
	Host            string
	Port            int
	ServiceBasePath string
}

type database struct {
	Name        string
	Username    string
	Password    string
	Host        string
	Port        int
	PoolOptions dbPool
	DbCache     dbCache
}

type dbCache struct {
	CleanupInterval string
	TTL             string
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
