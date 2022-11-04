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

import (
	"io/ioutil"
	"os"
	"reflect"
	"sync"

	toml "github.com/pelletier/go-toml"
	logger "github.com/sirupsen/logrus"
	"github.com/wso2/apk/adapter/pkg/config"
)

var (
	onceConfigRead      sync.Once
	grpcAdapterConfig 	*Config
	mgwHome             string
)

//constants related to utility functions
const (
	// RelativeConfigPath is the relative file path where the configuration file is.
	relativeConfigPath = "/conf/config.toml"
)

// ReadConfigs implements adapter configuration read operation. The read operation will happen only once, hence
// the consistancy is ensured.
//
// If the "MGW_HOME" variable is set, the configuration file location would be picked relative to the
// variable's value ("/conf/config.toml"). otherwise, the "MGW_HOME" variable would be set to the directory
// from where the executable is called from.
//
// Returns the configuration object mapped from the configuration file during the startup.
func ReadConfigs() *Config {
	onceConfigRead.Do(func() {
		grpcAdapterConfig = defaultConfig
		mgwHome = config.GetMgwHome()
		_, err := os.Stat(mgwHome + relativeConfigPath)
		if err != nil {
			logger.Fatal("Configuration file not found.", err)
		}
		content, readErr := ioutil.ReadFile(mgwHome + relativeConfigPath)
		if readErr != nil {
			logger.Fatal("Error reading configurations. ", readErr)
			return
		}
		parseErr := toml.Unmarshal(content, grpcAdapterConfig)
		if parseErr != nil {
			logger.Fatal("Error parsing the configuration ", parseErr)
			return
		}
		// Resolve environment variables.
		config.ResolveConfigEnvValues(reflect.ValueOf(&(grpcAdapterConfig.ManagementServer)).Elem(), "Server", true)
	})
	return grpcAdapterConfig
}
