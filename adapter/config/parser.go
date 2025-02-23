/*
 *  Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

// Package config contains the implementation and data structures related to configurations and
// configuration (log and adapter config) parsing. If a new configuration is introduced to the adapter
// configuration file, the corresponding change needs to be added to the relevant data stucture as well.
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sync"

	toml "github.com/pelletier/go-toml"
	pkgconf "github.com/wso2/apk/adapter/pkg/config"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	onceConfigRead sync.Once
	adapterConfig  *Config
)

const (
	// RelativeConfigPath is the relative file path where the configuration file is.
	relativeConfigPath   = "/conf/config.toml"
	gatewayTypeDefault   = "Onprem"
	gatewayTypeChoreo    = "Choreo"
	gatewayTypeChoreoPDP = "Choreo-PDP"
	gatewayTypeValue     = "gatewayType"
)

// ReadConfigs implements adapter configuration read operation. The read operation will happen only once, hence
// the consistancy is ensured.
//
// If the "APK_HOME" variable is set, the configuration file location would be picked relative to the
// variable's value ("/conf/config.toml"). otherwise, the "APK_HOME" variable would be set to the directory
// from where the executable is called from.
//
// Returns the configuration object that is initialized with default values. Changes to the default
// configuration object is achieved through the configuration file.
func ReadConfigs() *Config {
	onceConfigRead.Do(func() {
		adapterConfig = defaultConfig
		_, err := os.Stat(pkgconf.GetApkHome() + relativeConfigPath)
		if err != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1000, logging.BLOCKER, "Configuration file not found, error: %v", err.Error()))
		}
		content, readErr := ioutil.ReadFile(pkgconf.GetApkHome() + relativeConfigPath)
		if readErr != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1001, logging.BLOCKER, "Error reading configurations, error: %v", readErr.Error()))
			return
		}
		parseErr := toml.Unmarshal(content, adapterConfig)
		if parseErr != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1002, logging.BLOCKER, "Error parsing the configurations, error: %v", parseErr))
			return
		}

		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Adapter)).Elem(), "Adapter", true)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Envoy)).Elem(), "Router", true)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Enforcer)).Elem(), "Enforcer", false)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Analytics)).Elem(), "Analytics", false)

		// validate the analytics configuration values
		validationErr := validateAnalyticsConfigs(adapterConfig)
		if validationErr != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1002, logging.BLOCKER, "Error validating the configurations, error: %v", parseErr))
			return
		}
	})
	return adapterConfig
}

// SetConfig sets the given configuration to the adapter configuration
func SetConfig(conf *Config) {
	adapterConfig = conf
}

// SetDefaultConfig sets the default configuration to the adapter configuration
func SetDefaultConfig() {
	adapterConfig = defaultConfig
}

// ReadLogConfigs implements adapter/proxy log-configuration read operation.The read operation will happen only once, hence
// the consistancy is ensured.
//
// If the "APK_HOME" variable is set, the log configuration file location would be picked relative to the
// variable's value ("/conf/log_config.toml"). otherwise, the "APK_HOME" variable would be set to the directory
// from where the executable is called from.
//
// Returns the log configuration object mapped from the configuration file during the startup.
func ReadLogConfigs() *pkgconf.LogConfig {
	return pkgconf.ReadLogConfigs()
}

// ClearLogConfigInstance removes the existing configuration.
// Then the log configuration can be re-initialized.
func ClearLogConfigInstance() {
	pkgconf.ClearLogConfigInstance()
}

// GetLogConfigPath returns the file location of the log-config path
func GetLogConfigPath() (string, error) {
	return pkgconf.GetLogConfigPath()
}

// GetApkHome reads the APK_HOME environmental variable and returns the value.
// This represent the directory where the distribution is located.
// If the env variable is not present, the directory from which the executable is triggered will be assigned.
func GetApkHome() string {
	return pkgconf.GetApkHome()
}

func validateAnalyticsConfigs(conf *Config) error {

	gatewayType := gatewayTypeDefault
	if _, exists := conf.Analytics.Properties[gatewayTypeValue]; exists {
		gatewayType = conf.Analytics.Properties[gatewayTypeValue]
	} else {
		conf.Analytics.Properties[gatewayTypeValue] = gatewayTypeDefault
	}

	allowedValuesForGatewayType := map[string]bool{
		gatewayTypeDefault:   true,
		gatewayTypeChoreo:    true,
		gatewayTypeChoreoPDP: true,
	}

	if _, exists := allowedValuesForGatewayType[gatewayType]; !exists {
		return fmt.Errorf("invalid configuration value for analytics.gatewayType. Allowed values are %s, %s, or %s",
			gatewayTypeDefault, gatewayTypeChoreo, gatewayTypeChoreoPDP)
	}
	return nil
}
