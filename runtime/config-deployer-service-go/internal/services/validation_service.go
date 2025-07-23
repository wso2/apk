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

package services

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/services/validators"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"strings"
)

type ValidationService struct{}

// ValidateAndRetrieveDefinition validates and retrieves the API definition from the provided URL or file content.
func (validationService *ValidationService) ValidateAndRetrieveDefinition(apiType, url string, content []byte,
	fileName string) (*dto.APIDefinitionValidationResponse, error) {
	if url != "" {
		definition, err := util.RetrieveDefinitionFromUrl(url)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve definition from URL: %w", err)
		}
		return validateOpenAPIDefinition(apiType, nil, definition, "", true)
	}
	if fileName != "" && len(content) > 0 {
		return validateOpenAPIDefinition(apiType, content, "", fileName, true)
	}
	return nil, fmt.Errorf("either URL or file content must be provided")
}

// validateOpenAPIDefinition validates the OpenAPI definition based on the API type and returns the validation response.
func validateOpenAPIDefinition(apiType string, inputByteArray []byte, apiDefinition string, fileName string,
	returnContent bool) (*dto.APIDefinitionValidationResponse, error) {
	var validationResponse *dto.APIDefinitionValidationResponse
	var err error

	switch strings.ToUpper(apiType) {
	case constants.API_TYPE_REST:
		restAPIValidator := &validators.RESTAPIValidator{}
		if len(inputByteArray) > 0 {
			if fileName != "" {
				if strings.HasSuffix(fileName, ".zip") {
					validationResponse, err = restAPIValidator.ExtractAndValidateOpenAPIArchive(inputByteArray, returnContent)
					if err != nil {
						return nil, err
					}
				} else {
					openAPIContent := string(inputByteArray)
					validationResponse, err = restAPIValidator.ValidateAPIDefinition(openAPIContent, returnContent)
					if err != nil {
						return nil, err
					}
				}
			} else {
				openAPIContent := string(inputByteArray)
				validationResponse, err = restAPIValidator.ValidateAPIDefinition(openAPIContent, returnContent)
				if err != nil {
					return nil, err
				}
			}
		} else if apiDefinition != "" {
			validationResponse, err = restAPIValidator.ValidateAPIDefinition(apiDefinition, returnContent)
			if err != nil {
				return nil, err
			}
		}
	case constants.API_TYPE_GRAPHQL:
		if strings.HasSuffix(fileName, ".graphql") || strings.HasSuffix(fileName, ".txt") ||
			strings.HasSuffix(fileName, ".sdl") {
			graphqlAPIValidator := &validators.GraphQLAPIValidator{}
			validationResponse, err = graphqlAPIValidator.ValidateGraphQLSchema(string(inputByteArray), returnContent)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("invalid definition file type provided: %s", fileName)
		}
	case constants.API_TYPE_GRPC:
		if len(inputByteArray) > 0 {
			if strings.HasSuffix(fileName, ".zip") || strings.HasSuffix(fileName, ".proto") {
				grpcAPIValidator := &validators.GRPCAPIValidator{}
				validationResponse, err = grpcAPIValidator.ValidateGRPCAPIDefinition(inputByteArray)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("invalid definition file type provided: %s", fileName)
			}
		} else {
			return nil, fmt.Errorf("invalid definition file type provided: %s", fileName)
		}
	}
	return validationResponse, nil
}
