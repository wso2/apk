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

package validators

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/parsers"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RESTAPIValidator struct{}

// ExtractAndValidateOpenAPIArchive extracts the OpenAPI archive, validates it, and returns the validation response.
func (restAPIValidator *RESTAPIValidator) ExtractAndValidateOpenAPIArchive(inputByteArray []byte, returnContent bool) (*dto.APIDefinitionValidationResponse, error) {
	// Create temporary directory
	tempDir := filepath.Join(os.Getenv(constants.JAVA_IO_TMPDIR), constants.OPENAPI_ARCHIVES_TEMP_FOLDER, uuid.New().String())
	archivePath := filepath.Join(tempDir, constants.OPENAPI_ARCHIVE_ZIP_FILE)
	extractedLocation, err := util.ExtractUploadedArchive(inputByteArray, constants.OPENAPI_EXTRACTED_DIRECTORY, archivePath, tempDir)
	if err != nil {
		deleteErr := util.DeleteDirectory(tempDir)
		if deleteErr != nil {
			return nil, deleteErr
		}
		return nil, fmt.Errorf("error in accessing uploaded API archive: %w", err)
	}

	// Clean up temporary directory after function completes
	defer os.RemoveAll(tempDir)

	// Find archive directory
	var archiveDirectory string
	files, err := os.ReadDir(extractedLocation)
	if err != nil {
		return nil, fmt.Errorf("error reading extracted directory: %s, %w", extractedLocation, err)
	}
	dirCount := 0
	for _, file := range files {
		if file.IsDir() {
			archiveDirectory = filepath.Join(extractedLocation, file.Name())
			dirCount++
		}
	}
	if dirCount > 1 {
		return nil, fmt.Errorf("swagger definitions should be placed under one root folder")
	}
	if archiveDirectory == "" {
		return nil, fmt.Errorf("could not find an archive in the given ZIP file")
	}

	// Find and read master swagger file
	masterSwaggerPath, err := util.CheckMasterSwagger(archiveDirectory)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(masterSwaggerPath)
	if err != nil {
		return nil, fmt.Errorf("error reading master swagger file: %w", err)
	}

	openAPIContent := string(content)
	// Get swagger version
	version, err := util.GetSwaggerVersion(openAPIContent)
	if err != nil {
		return nil, err
	}

	filePath, err := filepath.Abs(masterSwaggerPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path of master swagger: %s, %w", masterSwaggerPath, err)
	}

	oasParser := parsers.OAS3Parser{}
	switch version {
	case constants.OPEN_API:
		openAPI, err := oasParser.ParseOpenAPI3(filePath)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.MarshalIndent(openAPI, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshalling openapi spec: %w", err)
		}
		openAPIContent = string(jsonBytes)
	case constants.SWAGGER:
		swagger, err := oasParser.ParseSwagger2(filePath)
		if err != nil {
			return nil, err
		}
		yamlBytes, err := yaml.Marshal(swagger)
		if err != nil {
			return nil, fmt.Errorf("error marshalling swagger spec: %w", err)
		}
		openAPIContent = string(yamlBytes)
	default:
		return nil, fmt.Errorf("unsupported Swagger version: %d", version)
	}

	return restAPIValidator.ValidateAPIDefinition(openAPIContent, returnContent)
}

// ValidateAPIDefinition removes unsupported blocks, validates the API definition and returns the validation response.
func (restAPIValidator *RESTAPIValidator) ValidateAPIDefinition(apiDefinition string, returnJsonContent bool) (*dto.APIDefinitionValidationResponse, error) {
	apiDefinitionProcessed := apiDefinition
	oasParser := parsers.OAS3Parser{}
	if !strings.HasPrefix(strings.TrimSpace(apiDefinition), "{") {
		// Convert YAML to JSON
		jsonData, err := util.YamlToJSON(apiDefinition)
		if err != nil {
			return nil, fmt.Errorf("error while reading API definition yaml: %w", err)
		}
		apiDefinitionProcessed = jsonData
	}
	apiDefinitionProcessed, err := oasParser.RemoveUnsupportedBlocksFromResources(apiDefinitionProcessed)
	if err != nil {
		return nil, fmt.Errorf("error while removing unsupported blocks: %w", err)
	}
	if apiDefinitionProcessed != "" {
		apiDefinition = apiDefinitionProcessed
	}
	validationResponse, err := validateAPIDefinition(apiDefinition, returnJsonContent)
	if err != nil {
		return nil, fmt.Errorf("error while validating API definition: %w", err)
	}
	if !validationResponse.IsValid {
		// TODO - IF invalid OAS3 found try OAS2 validation
	}
	return validationResponse, nil
}

// validateAPIDefinition validates the API definition with an optional host and returns the validation response.
func validateAPIDefinition(apiDefinition string, returnJsonContent bool) (*dto.APIDefinitionValidationResponse, error) {
	return validateAPIDefinitionWithHost(apiDefinition, "", returnJsonContent)
}

// validateAPIDefinitionWithHost validates the API definition with a specified host and returns the validation response.
func validateAPIDefinitionWithHost(apiDefinition, host string, returnJsonContent bool) (*dto.APIDefinitionValidationResponse, error) {
	validationResponse := &dto.APIDefinitionValidationResponse{}
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	ctx := util.GetOpenAPIValidationContext(loader.Context)

	doc, err := loader.LoadFromData([]byte(apiDefinition))
	if err != nil {
		validationResponse.IsValid = false
		return validationResponse, fmt.Errorf("error while loading OpenAPI document: %w", err)
	}
	if err := doc.Validate(ctx); err != nil {
		validationResponse.IsValid = false
		return validationResponse, fmt.Errorf("invalid OpenAPI V3 definition found: %w", err)
	}

	// Extract information from the valid OpenAPI document
	info := doc.Info
	var endpoints []string

	// Process servers/endpoints
	if doc.Servers != nil && len(doc.Servers) > 0 {
		for _, server := range doc.Servers {
			endpoint := server.URL
			if strings.HasPrefix(endpoint, "/") {
				var endpointWithHost string
				if host == "" {
					endpointWithHost = "http://api.yourdomain.com" + endpoint
				} else {
					endpointWithHost = host + endpoint
				}
				endpoints = append(endpoints, endpointWithHost)
			} else {
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	// Extract title and context
	var title, context string
	if info.Title != "" {
		title = info.Title
		// Remove whitespace and convert to lowercase for context
		re := regexp.MustCompile(`\s+`)
		context = strings.ToLower(re.ReplaceAllString(info.Title, ""))
	}

	// Extract description
	description := ""
	if info.Description != "" {
		description = info.Description
	}

	// Update validation response as success
	updateValidationResponseAsSuccess(validationResponse, apiDefinition, doc.OpenAPI, title, info.Version,
		context, description, endpoints)

	// Handle JSON content return
	if returnJsonContent {
		if !strings.HasPrefix(strings.TrimSpace(apiDefinition), "{") {
			jsonContent, err := util.YamlToJSON(apiDefinition)
			if err != nil {
				return nil, fmt.Errorf("error while reading API definition yaml: %w", err)
			}
			validationResponse.JSONContent = jsonContent
		} else {
			validationResponse.JSONContent = apiDefinition
		}
	}
	return validationResponse, nil
}

// updateValidationResponseAsSuccess updates the validation response with success information
func updateValidationResponseAsSuccess(validationResponse *dto.APIDefinitionValidationResponse, apiDefinition,
	openAPIVersion, title, version, context, description string, endpoints []string) {
	validationResponse.IsValid = true
	validationResponse.Content = apiDefinition
	validationResponse.OpenAPIVersion = openAPIVersion
	validationResponse.Name = title
	validationResponse.Version = version
	validationResponse.Context = context
	validationResponse.Description = description
	validationResponse.Endpoints = endpoints
}
