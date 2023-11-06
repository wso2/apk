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

import (
	"io/ioutil"
	"os"
	"reflect"
	"sync"

	toml "github.com/pelletier/go-toml"
	pkgconf "github.com/wso2/apk/adapter/pkg/config"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	onceConfigRead   sync.Once
	apkHome          string
	logConfigPath    string
	controllerConfig *Config
	envVariableMap   map[string]string
)

const (
	// RelativeConfigPath is the relative file path where the configuration file is.
	relativeConfigPath = "/conf/config.toml"
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
		controllerConfig = defaultConfig
		_, err := os.Stat(pkgconf.GetApkHome() + relativeConfigPath)
		if err != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1000, logging.BLOCKER, "Configuration file not found, error: %v", err.Error()))
		}
		content, readErr := ioutil.ReadFile(pkgconf.GetApkHome() + relativeConfigPath)
		if readErr != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1001, logging.BLOCKER, "Error reading configurations, error: %v", readErr.Error()))
			return
		}
		parseErr := toml.Unmarshal(content, controllerConfig)
		if parseErr != nil {
			loggerConfig.ErrorC(logging.PrintError(logging.Error1002, logging.BLOCKER, "Error parsing the configurations, error: %v", parseErr.Error()))
			return
		}

		pkgconf.ResolveConfigEnvValues(reflect.ValueOf(&(controllerConfig.CommonController)).Elem(), "CommonController", true)
	})
	return controllerConfig
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
