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
	"context"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

// CheckMasterSwagger checks for the existence of the master swagger file in the given directory
func CheckMasterSwagger(archiveDirectory string) (string, error) {
	jsonPath := filepath.Join(archiveDirectory, constants.OPENAPI_MASTER_JSON)
	if FileExists(jsonPath) {
		return jsonPath, nil
	}
	yamlPath := filepath.Join(archiveDirectory, constants.OPENAPI_MASTER_YAML)
	if FileExists(yamlPath) {
		return yamlPath, nil
	}
	return "", fmt.Errorf("could not find a master swagger file with the name of swagger.json/swagger.yaml %s",
		archiveDirectory)
}

// GetSwaggerVersion determines the version of the API specification
func GetSwaggerVersion(content string) (constants.SwaggerVersion, error) {
	// Try to parse as JSON first
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		// Try to parse as YAML
		if err := yaml.Unmarshal([]byte(content), &jsonData); err != nil {
			return constants.SwaggerVersion(0), fmt.Errorf("error unmarshalling Swagger version: %w", err)
		}
	}

	// Check for OpenAPI 3.x
	if openapi, exists := jsonData["openapi"]; exists {
		if openapiStr, ok := openapi.(string); ok {
			if len(openapiStr) > 0 && openapiStr[0] == '3' {
				return constants.OPEN_API, nil
			}
		}
	}

	// Check for Swagger 2.0
	if swagger, exists := jsonData["swagger"]; exists {
		if swaggerStr, ok := swagger.(string); ok {
			if swaggerStr == "2.0" {
				return constants.SWAGGER, nil
			}
		}
	}
	return 0, fmt.Errorf("invalid OAS definition provided")
}

// GetOpenAPIValidationContext returns a configured context for OpenAPI validation
// with common validation options applied
func GetOpenAPIValidationContext(baseCtx context.Context) context.Context {
	return openapi3.WithValidationOptions(baseCtx,
		openapi3.DisableExamplesValidation(),
		openapi3.DisableSchemaDefaultsValidation(),
		openapi3.AllowExtraSiblingFields(constants.OPENAPI_ALLOWED_EXTRA_SIBLING_FIELDS),
	)
}
