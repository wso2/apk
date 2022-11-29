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

// Package utils holds the implementation for common utility functions
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	// TODO: (VirajSalaka) remove outdated dependency
	"unicode"

	"github.com/ghodss/yaml"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/oasparser/constants"
	"github.com/wso2/apk/adapter/pkg/logging"
)

// ToJSON converts a single YAML document into a JSON document
// or returns an error. If the document appears to be JSON the
// YAML decoding path is not used.
// If the input file is json, it would be returned as it is.
func ToJSON(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}
	return yaml.YAMLToJSON(data)
}

var jsonPrefix = []byte("{")

func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix)
}

func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}

// FindAPIDefinitionVersion finds the API definition version for the given json content.
func FindAPIDefinitionVersion(jsn []byte) string {
	var result map[string]interface{}

	err := json.Unmarshal(jsn, &result)
	if err != nil {
		logger.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error while JSON unmarshalling to find the API definition version. %s", err.Error()),
			Severity:  logging.MINOR,
			ErrorCode: 2209,
		})
	}

	if _, ok := result[constants.Swagger]; ok {
		return constants.Swagger2
	} else if _, ok := result[constants.OpenAPI]; ok {
		return constants.OpenAPI3
	} else if versionNumber, ok := result[constants.AsyncAPI]; ok {
		if strings.HasPrefix(versionNumber.(string), "2") {
			return constants.AsyncAPI2
		}
		logger.LoggerOasparser.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("AsyncAPI version %s is not supported.", versionNumber.(string)),
			Severity:  logging.MINOR,
			ErrorCode: 2210,
		})
		return constants.NotSupported
	}
	logger.LoggerOasparser.ErrorC(logging.ErrorDetails{
		Message:   "API definition version is not defined.",
		Severity:  logging.MINOR,
		ErrorCode: 2211,
	})
	return constants.NotDefined
}

// FileNameWithoutExtension returns the file name without the extension
// ex: when provided the path "/foo/hello.world" it returns "hello"
func FileNameWithoutExtension(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}
