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

// RetrieveAndValidateDefinitionFromURL validates and retrieves the API definition from the provided URL.
func (validationService *ValidationService) RetrieveAndValidateDefinitionFromURL(apiType,
	url string) (*dto.APIDefinitionValidationResponse, error) {
	definition, err := util.RetrieveDefinitionFromUrl(url)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve definition from URL: %w", err)
	} else if definition == "" {
		return nil, fmt.Errorf("API definition is empty or null for the provided URL: %s", url)
	}
	return validateAPIDefinitionFromURL(apiType, definition, true)
}

// RetrieveAndValidateDefinitionFromFile validates and retrieves the API definition from the provided file content.
func (validationService *ValidationService) RetrieveAndValidateDefinitionFromFile(apiType string, fileName string,
	content []byte) (*dto.APIDefinitionValidationResponse, error) {
	return validateAPIDefinitionFromFile(apiType, fileName, content, true)
}

// validateAPIDefinitionFromURL validates the definition extracted from a URL based on the API type and returns the validation response.
func validateAPIDefinitionFromURL(apiType string, apiDefinition string,
	returnContent bool) (*dto.APIDefinitionValidationResponse, error) {
	var validationResponse *dto.APIDefinitionValidationResponse
	var err error

	switch strings.ToUpper(apiType) {
	case constants.API_TYPE_REST:
		restAPIValidator := &validators.RESTAPIValidator{}
		validationResponse, err = restAPIValidator.ValidateOpenAPIDefinition(apiDefinition, returnContent)
		if err != nil {
			return nil, err
		}
	case constants.API_TYPE_GRAPHQL:
		// TODO - Handle GraphQL definition from url
		return nil, fmt.Errorf("handling GraphQL definition from URL is not implemented yet")
	case constants.API_TYPE_GRPC:
		// TODO - Handle gRPC definition from url
		return nil, fmt.Errorf("handling gRPC definition from URL is not implemented yet")
	}

	return validationResponse, nil
}

// validateAPIDefinitionFromFile validates the definition extracted from a file based on the API type and returns the validation response.
func validateAPIDefinitionFromFile(apiType string, fileName string, inputByteArray []byte,
	returnContent bool) (*dto.APIDefinitionValidationResponse, error) {
	var validationResponse *dto.APIDefinitionValidationResponse
	var err error

	switch strings.ToUpper(apiType) {
	case constants.API_TYPE_REST:
		restAPIValidator := &validators.RESTAPIValidator{}
		if strings.HasSuffix(fileName, ".zip") {
			validationResponse, err = restAPIValidator.ExtractAndValidateOpenAPIArchive(inputByteArray, returnContent)
			if err != nil {
				return nil, err
			}
		} else {
			openAPIContent := string(inputByteArray)
			validationResponse, err = restAPIValidator.ValidateOpenAPIDefinition(openAPIContent, returnContent)
			if err != nil {
				return nil, err
			}
		}
	case constants.API_TYPE_GRAPHQL:
		if strings.HasSuffix(fileName, ".graphql") || strings.HasSuffix(fileName, ".txt") ||
			strings.HasSuffix(fileName, ".sdl") {
			graphqlAPIValidator := &validators.GraphQLAPIValidator{}
			validationResponse, err = graphqlAPIValidator.ValidateGraphQLSchema(string(inputByteArray))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("invalid definition file type provided: %s, for graphql API", fileName)
		}
	case constants.API_TYPE_GRPC:
		if strings.HasSuffix(fileName, ".zip") || strings.HasSuffix(fileName, ".proto") {
			grpcAPIValidator := &validators.GRPCAPIValidator{}
			validationResponse, err = grpcAPIValidator.ValidateGRPCAPIDefinition(inputByteArray)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("invalid definition file type provided: %s, for gRPC API", fileName)
		}
	}

	return validationResponse, nil
}
