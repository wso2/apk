/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

import (
	"sync"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/kelseyhightower/envconfig"
	"github.com/wso2/apk/envoy-gateway-extension-server/internal/logging"
)

// Server holds the configuration parameters for the application.
type Server struct {
	LogLevel                         string `envconfig:"LOG_LEVEL" default:"DEBUG"`
	Logger                           logging.Logger
	ExtensionServerHost            string `envconfig:"EXTENSION_SERVER_HOST" default:"0.0.0.0"`
	ExtensionServerPort            string `envconfig:"EXTENSION_SERVER_PORT" default:"5005"`
}

type metrics struct {
	Enabled            bool   `envconfig:"METRICS_ENABLED" default:"false"`
	Type               string `envconfig:"METRICS_TYPE" default:"prometheus"`
	Port               int32  `envconfig:"METRICS_PORT" default:"18006"`
	CollectionInterval int32  `envconfig:"METRICS_COLLECTION_INTERVAL" default:"15"`
}

// package-level variable and mutex for thread safety
var (
	processOnce     sync.Once
	settingInstance *Server
)

// GetConfig initializes and returns a singleton instance of the Settings struct.
// It uses sync.Once to ensure that the initialization logic is executed only once,
// making it safe for concurrent use. If there is an error during the initialization,
// the function will panic.
//
// Returns:
//
//	*Settings - A pointer to the singleton instance of the Settings struct. from environment variables.
func GetConfig() *Server {
	var err error
	processOnce.Do(func() {
		settingInstance = &Server{}
		err = envconfig.Process("", settingInstance)
	})
	if err != nil {
		panic(err)
	}
	// Create Logger based on the env var
	settingInstance.Logger = logging.NewLogger(&egv1a1.EnvoyGatewayLogging{
		Level: map[egv1a1.EnvoyGatewayLogComponent]egv1a1.LogLevel{
			egv1a1.LogComponentGatewayDefault: egv1a1.LogLevel(settingInstance.LogLevel),
		},
	})
	return settingInstance
}
