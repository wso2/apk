/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
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

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

// YamlToJSON converts YAML string to JSON string
func YamlToJSON(yamlContent string) (string, error) {
	var data interface{}

	// Parse YAML
	err := yaml.Unmarshal([]byte(yamlContent), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Convert interface{} keys to string keys for JSON compatibility
	convertedData := convertInterfaceKeysToString(data)

	// Convert to JSON
	jsonBytes, err := json.Marshal(convertedData)
	if err != nil {
		return "", fmt.Errorf("failed to convert to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

// JsonToYaml converts a JSON string to YAML format
func JsonToYaml(jsonString string) (string, error) {
	// Parse JSON string into a map
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	yamlData, err := MarshalToYAMLWithIndent(jsonData, 2)
	if err != nil {
		return "", fmt.Errorf("failed to convert to YAML: %w", err)
	}

	return string(yamlData), nil
}

// convertInterfaceKeysToString recursively converts map[interface{}]interface{} to map[string]interface{}
func convertInterfaceKeysToString(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// Handle map[string]interface{} that might contain nested map[interface{}]interface{}
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = convertInterfaceKeysToString(value)
		}
		return result
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			strKey := fmt.Sprintf("%v", key)
			result[strKey] = convertInterfaceKeysToString(value)
		}
		return result
	case []interface{}:
		for i, value := range v {
			v[i] = convertInterfaceKeysToString(value)
		}
		return v
	default:
		return v
	}
}

// MarshalToYAMLWithIndent marshals a struct to YAML with custom indentation
func MarshalToYAMLWithIndent(data interface{}, indent int) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(indent)
	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error occurred while encoding to YAML: %w", err)
	}
	encoder.Close()
	return buf.Bytes(), nil
}
