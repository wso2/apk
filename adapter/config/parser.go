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
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml"
	logger "github.com/sirupsen/logrus"
	pkgconf "github.com/wso2/apk/adapter/pkg/config"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	onceConfigRead      sync.Once
	onceGetDefaultVhost sync.Once
	adapterConfig       *Config
	defaultVhost        map[string]string
)

const (
	// The environtmental variable which represents the path of the distribution in host machine.
	mgwHomeEnvVariable = "MGW_HOME"
	// RelativeConfigPath is the relative file path where the configuration file is.
	relativeConfigPath = "/conf/config.toml"
)

// Constants related to utility functions
const (
	tenantDomainSeparator = "@"
	superTenantDomain     = "carbon.super"
)

// ReadConfigs implements adapter configuration read operation. The read operation will happen only once, hence
// the consistancy is ensured.
//
// If the "MGW_HOME" variable is set, the configuration file location would be picked relative to the
// variable's value ("/conf/config.toml"). otherwise, the "MGW_HOME" variable would be set to the directory
// from where the executable is called from.
//
// Returns the configuration object that is initialized with default values. Changes to the default
// configuration object is achieved through the configuration file.
func ReadConfigs() *Config {
	onceConfigRead.Do(func() {
		adapterConfig = defaultConfig
		_, err := os.Stat(pkgconf.GetMgwHome() + relativeConfigPath)
		if err != nil {
			loggerConfig.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Configuration file not found : %s", err.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1000,
			})
		}
		content, readErr := ioutil.ReadFile(pkgconf.GetMgwHome() + relativeConfigPath)
		if readErr != nil {
			loggerConfig.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error reading configurations : %s", readErr.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1001,
			})
			return
		}
		parseErr := toml.Unmarshal(content, adapterConfig)
		if parseErr != nil {
			loggerConfig.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error parsing the configurations : %s", parseErr.Error()),
				Severity:  logging.BLOCKER,
				ErrorCode: 1002,
			})
			return
		}

		adapterConfig.resolveDeprecatedProperties()
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Adapter)).Elem(), "Adapter", true)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.ControlPlane)).Elem(), "ControlPlane", true)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Envoy)).Elem(), "Router", true)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Enforcer)).Elem(), "Enforcer", false)
		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(adapterConfig.Analytics)).Elem(), "Analytics", false)
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
// If the "MGW_HOME" variable is set, the log configuration file location would be picked relative to the
// variable's value ("/conf/log_config.toml"). otherwise, the "MGW_HOME" variable would be set to the directory
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

// GetMgwHome reads the MGW_HOME environmental variable and returns the value.
// This represent the directory where the distribution is located.
// If the env variable is not present, the directory from which the executable is triggered will be assigned.
func GetMgwHome() string {
	return pkgconf.GetMgwHome()
}

// GetControlPlaneConnectedTenantDomain returns the tenant domain of the user used to authenticate with event hub.
func GetControlPlaneConnectedTenantDomain() string {
	// Read configurations to get the control plane authenticated user
	conf := ReadConfigs()

	// Populate data from the config
	cpTenantAdminUser := conf.ControlPlane.Username
	tenantDomain := strings.Split(cpTenantAdminUser, tenantDomainSeparator)
	if len(tenantDomain) > 1 {
		return tenantDomain[len(tenantDomain)-1]
	}
	return superTenantDomain
}

func (config *Config) resolveDeprecatedProperties() {
	if config.ControlPlane.ServiceURLDeprecated != UnassignedAsDeprecated {
		printDeprecatedWarningLog("controlPlane.serviceUrl", "controlPlane.serviceURL")
		config.ControlPlane.ServiceURL = config.ControlPlane.ServiceURLDeprecated
	}

	if len(config.Enforcer.Throttling.Publisher.URLGroupDeprecated) > 0 {
		printDeprecatedWarningLog("enforcer.throttling.publisher.urlGroup", "enforcer.throttling.publisher.URLGroup")
		config.Enforcer.Throttling.Publisher.URLGroup = config.Enforcer.Throttling.Publisher.URLGroupDeprecated
	}

	// For boolean values, adapter check if the condition is changed by checking against the default value it is originally
	// assigned.
	if !config.Enforcer.RestServer.Enable {
		printDeprecatedWarningLog("enforcer.restServer.enable", "enforcer.restServer.enabled")
		config.Enforcer.RestServer.Enabled = config.Enforcer.RestServer.Enable
	}
	if config.Adapter.Consul.Enable {
		printDeprecatedWarningLog("adapter.consul.enable", "adapter.consul.enabled")
		config.Adapter.Consul.Enabled = config.Adapter.Consul.Enable
	}
	if config.Enforcer.JwtGenerator.Enable {
		printDeprecatedWarningLog("enforcer.jwtGenerator.enable", "enforcer.jwtGenerator.enabled")
		config.Enforcer.JwtGenerator.Enabled = config.Enforcer.JwtGenerator.Enable
	}

}

func printDeprecatedWarningLog(deprecatedTerm, currentTerm string) {
	logger.Warnf("%s is deprecated. Use %s instead", deprecatedTerm, currentTerm)
}
