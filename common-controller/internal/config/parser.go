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
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml"
	pkgconf "github.com/wso2/apk/adapter/pkg/config"
	"github.com/wso2/apk/adapter/pkg/logging"
)

var (
	onceConfigRead   sync.Once
	controllerConfig *Config
	envVariableMap   map[string]string
)

const (
	// RelativeConfigPath is the relative file path where the configuration file is.
	relativeConfigPath = "/conf/config.toml"
	// EnvConfigPrefix is used when configs should be read from environment variables.
	EnvConfigPrefix   = "$env"
	envVariablePrefix = "APK_"
	// envVariableEntrySeparator is used as the separator used to denote nested structured properties.
	envVariableEntrySeparator = "_"
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
		content, readErr := os.ReadFile(pkgconf.GetApkHome() + relativeConfigPath)
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

// ResolveConfigEnvValues looks for the string type config values which should be read from environment variables
// and replace the respective config values from environment variable.
// v - relect.Value of the root level struct
// previousTag - the starting Tag corresponding to the root level struct
// resolveEnvTag - true if $env{} annotation needs to be resolved at adapter level
func ResolveConfigEnvValues(v reflect.Value, previousTag string, resolveEnvTag bool) {
	s := v
	for fieldNum := 0; fieldNum < s.NumField(); fieldNum++ {
		field := s.Field(fieldNum)
		currentTag := previousTag + envVariableEntrySeparator + s.Type().Field(fieldNum).Name
		resolveEnvForReflectValue(field, currentTag, resolveEnvTag)
	}
}

func resolveEnvForReflectValue(field reflect.Value, currentTag string, resolveEnvTag bool) {
	fieldKind := getKind(field)

	switch fieldKind {
	case reflect.String:
		if strings.Contains(fmt.Sprint(field.Interface()), EnvConfigPrefix) && resolveEnvTag {
			field.SetString(ResolveEnvValue(fmt.Sprint(field.Interface())))
		}
		resolveEnvStringValue(currentTag, field)
	case reflect.Slice:
		resolveEnvValueOfArray(field, currentTag, resolveEnvTag)
	case reflect.Array:
		// this condition is never reached.
		resolveEnvValueOfArray(field, currentTag, resolveEnvTag)
	case reflect.Struct:
		if field.Kind() == reflect.Struct {
			ResolveConfigEnvValues(field.Addr().Elem(), currentTag, resolveEnvTag)
		}
	case reflect.Bool:
		resolveEnvBooleanValue(currentTag, field)
	case reflect.Int:
		resolveEnvIntValue(currentTag, field)
	case reflect.Float32:
		resolveEnvFloat32Value(currentTag, field)
	case reflect.Uint:
		resolveEnvUIntValue(currentTag, field)
	case reflect.Map:
		resolveEnvValueOfMap(field, currentTag, resolveEnvTag)
	// this condition is never reached.:
	case reflect.Ptr:
	}
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

func resolveEnvValueOfArray(field reflect.Value, currentTag string, resolveEnvTag bool) {
	var arrayElementType reflect.Kind
	for index := 0; index < field.Len(); index++ {
		arrayElementType = field.Index(index).Kind()
		if field.Index(index).Kind() == reflect.Struct {
			ResolveConfigEnvValues(field.Index(index).Addr().Elem(), currentTag+envVariableEntrySeparator+strconv.Itoa(index), resolveEnvTag)
		} else if field.Index(index).Kind() == reflect.String && strings.Contains(field.Index(index).String(),
			EnvConfigPrefix) && resolveEnvTag {
			field.Index(index).SetString(ResolveEnvValue(field.Index(index).String()))
		}
	}

	if arrayElementType == reflect.Invalid {
		indirectStr := reflect.Indirect(field)
		valueSlice := reflect.MakeSlice(indirectStr.Type(), 1, 1)
		arrayElementType = valueSlice.Index(0).Kind()
	}

	variableValue, exists := envVariableMap[strings.ToUpper(envVariablePrefix+currentTag)]
	if exists {
		elementArrayAsString := splitStringAndTrim(variableValue)
		switch arrayElementType {
		case reflect.String:
			field.Set(reflect.ValueOf(elementArrayAsString))
		case reflect.Int:
			elementArrayAsInt := make([]int, len(elementArrayAsString))
			var parseErr error
			for index, stringElem := range elementArrayAsString {
				elementArrayAsInt[index], parseErr = strconv.Atoi(stringElem)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an integer array", variableValue)
					return
				}
			}
			field.Set(reflect.ValueOf(elementArrayAsInt))
		case reflect.Int32:
			elementArrayAsInt := make([]int32, len(elementArrayAsString))
			for index, stringElem := range elementArrayAsString {
				int64Val, parseErr := strconv.ParseInt(stringElem, 10, 32)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an integer32 array", variableValue)
					return
				}
				elementArrayAsInt[index] = int32(int64Val)
			}
			field.Set(reflect.ValueOf(elementArrayAsInt))
		case reflect.Int64:
			elementArrayAsInt := make([]int64, len(elementArrayAsString))
			for index, stringElem := range elementArrayAsString {
				int64Val, parseErr := strconv.ParseInt(stringElem, 10, 64)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an integer64 array", variableValue)
					return
				}
				elementArrayAsInt[index] = int64Val
			}
			field.Set(reflect.ValueOf(elementArrayAsInt))
		case reflect.Uint:
			elementArrayAsUInt := make([]uint, len(elementArrayAsString))
			for index, stringElem := range elementArrayAsString {
				uint64val, parseErr := strconv.ParseUint(stringElem, 10, 32)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an unsigned integer array", variableValue)
					return
				}
				elementArrayAsUInt[index] = uint(uint64val)
			}
			field.Set(reflect.ValueOf(elementArrayAsUInt))
		case reflect.Uint32:
			elementArrayAsUInt := make([]uint32, len(elementArrayAsString))
			for index, stringElem := range elementArrayAsString {
				uint64Val, parseErr := strconv.ParseUint(stringElem, 10, 32)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an unsigned integer32 array", variableValue)
					return
				}
				elementArrayAsUInt[index] = uint32(uint64Val)
			}
			field.Set(reflect.ValueOf(elementArrayAsUInt))
		case reflect.Uint64:
			elementArrayAsInt := make([]uint64, len(elementArrayAsString))
			for index, stringElem := range elementArrayAsString {
				int64Val, parseErr := strconv.ParseUint(stringElem, 10, 64)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an unsigned integer64 array", variableValue)
					return
				}
				elementArrayAsInt[index] = int64Val
			}
			field.Set(reflect.ValueOf(elementArrayAsInt))
		case reflect.Float32:
			elementArrayAsFloat := make([]float32, len(elementArrayAsString))
			for index, stringElem := range elementArrayAsString {
				float64Val, parseErr := strconv.ParseFloat(stringElem, 32)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an float32 array", variableValue)
					return
				}
				elementArrayAsFloat[index] = float32(float64Val)
			}
			field.Set(reflect.ValueOf(elementArrayAsFloat))
		case reflect.Float64:
			elementArrayAsFloat := make([]float64, len(elementArrayAsString))
			var parseErr error
			for index, stringElem := range elementArrayAsString {
				elementArrayAsFloat[index], parseErr = strconv.ParseFloat(stringElem, 64)
				if parseErr != nil {
					loggerConfig.Logger.Errorf("Error while parsing %s as an float64 array", variableValue)
					return
				}
			}
			field.Set(reflect.ValueOf(elementArrayAsFloat))
		}
	}
}

func splitStringAndTrim(input string) []string {
	ElementArrayAsString := strings.Split(input, ",")
	for index := 0; index < len(ElementArrayAsString); index++ {
		ElementArrayAsString[index] = strings.TrimSpace(ElementArrayAsString[index])
	}
	return ElementArrayAsString
}

func resolveEnvValueOfMap(field reflect.Value, currentTag string, resolveEnvTag bool) {
	// Only map [string]string is supported here.
	for _, key := range field.MapKeys() {
		if field.MapIndex(key).Kind() == reflect.String && key.Kind() == reflect.String {
			variableName := currentTag + envVariableEntrySeparator + key.String()
			variableName = strings.ReplaceAll(variableName, ".", "_")
			variableValue, exists := envVariableMap[envVariablePrefix+strings.ToUpper(variableName)]
			if exists {
				field.SetMapIndex(reflect.ValueOf(key.String()), reflect.ValueOf(variableValue))
			}
		}
	}
}

// ResolveEnvValue replace the respective config values from environment variable.
func ResolveEnvValue(value string) string {
	re := regexp.MustCompile(`(?s)\{(.*)}`) // regex to get everything in between curly brackets
	m := re.FindStringSubmatch(value)
	if len(m) > 1 {
		envValue, exists := os.LookupEnv(m[1])
		if exists {
			return strings.ReplaceAll(re.ReplaceAllString(value, envValue), EnvConfigPrefix, "")
		}
	}
	return value
}
func resolveEnvStringValue(key string, value reflect.Value) {
	variableValue, exists := envVariableMap[envVariablePrefix+strings.ToUpper(key)]
	if exists {
		value.SetString(variableValue)
	}
}

func resolveEnvBooleanValue(key string, value reflect.Value) {
	var resolvedValue bool
	var parseErr error
	variableValue, exists := envVariableMap[envVariablePrefix+strings.ToUpper(key)]
	if exists {
		resolvedValue, parseErr = strconv.ParseBool(variableValue)
		if parseErr != nil {
			loggerConfig.Logger.Errorf("Error while parsing %s as a boolean value.", key)
			return
		}
		value.SetBool(resolvedValue)
	}
}

func resolveEnvIntValue(key string, value reflect.Value) {
	var resolvedValue int
	var parseErr error
	variableValue, exists := envVariableMap[envVariablePrefix+strings.ToUpper(key)]
	if exists {
		resolvedValue, parseErr = strconv.Atoi(variableValue)
		if parseErr != nil {
			loggerConfig.Logger.Errorf("Error while parsing %s as a int value. :", key)
			return
		}
		value.SetInt(int64(resolvedValue))
	}
}

func resolveEnvUIntValue(key string, value reflect.Value) {
	var resolvedValue uint64
	var parseErr error
	variableValue, exists := envVariableMap[envVariablePrefix+strings.ToUpper(key)]
	if exists {
		resolvedValue, parseErr = strconv.ParseUint(variableValue, 10, 32)
		if parseErr != nil {
			loggerConfig.Logger.Errorf("Error while parsing %s as a uint value.", key)
			return
		}
		value.SetUint(resolvedValue)
	}
}

func resolveEnvFloat32Value(key string, value reflect.Value) {
	var resolvedValue float64
	var parseErr error
	variableValue, exists := envVariableMap[envVariablePrefix+strings.ToUpper(key)]
	if exists {
		resolvedValue, parseErr = strconv.ParseFloat(variableValue, 32)
		if parseErr != nil {
			loggerConfig.Logger.Errorf("Error while parsing %s as a float value.", key)
			return
		}
		value.SetFloat(resolvedValue)
	}
}
