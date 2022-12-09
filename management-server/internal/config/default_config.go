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

// Configuration object which is populated with default values.
var defaultConfig = &Config{
	ManagementServer: managementServer{
		XDSPort:          18000,
		NodeLabels:       []string{"default"},
		GRPCPort:         8765,
		NotificationPort: 8766,
	},
	Database: database{
		Name:     "WSO2AM_DB",
		Username: "wso2carbon",
		Password: "wso2carbon",
		Host:     "wso2apk-db-service",
		Port:     5432,
		PoolOptions: dbPool{
			PoolMaxConns:              4,
			PoolMinConns:              0,
			PoolMaxConnLifetime:       "1h",
			PoolMaxConnIdleTime:       "1h",
			PoolHealthCheckPeriod:     "1m",
			PoolMaxConnLifetimeJitter: "1s",
		},
		DbCache: dbCache{
			CleanupInterval: "1h",
			TTL:             "1h",
		},
	},
	BackOffice: backOffice{
		Host:            "localhost",
		Port:            9443,
		ServiceBasePath: "/api/am/backoffice/internal/apis",
	},
}
