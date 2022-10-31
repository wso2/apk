/*
 * Copyright (c) 2022, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package utils

import (
	"bytes"
	"encoding/json"
	"unicode"

	"github.com/ghodss/yaml"
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
		HandleErrorAndExit("Error while JSON unmarshalling to find the API definition version.", err)
	}

	if _, ok := result[Swagger]; ok {
		return Swagger2
	} else if _, ok := result[OpenAPI]; ok {
		return OpenAPI3
	}
	HandleErrorAndContinue("API definition version is not defined.", nil)
	return NotDefined
}
