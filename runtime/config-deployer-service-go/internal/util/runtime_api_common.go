/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package util

import (
	"config-deployer-service-go/internal/constants"
	"config-deployer-service-go/internal/dto"
	"config-deployer-service-go/internal/model"
	"fmt"
	"strings"
)

func ValidateOpenAPIDefinition(apiType string, inputByteArray []byte, apiDefinition string, fileName string,
	returnContent bool) (*dto.APIDefinitionValidationResponse, error) {
	var validationResponse *dto.APIDefinitionValidationResponse

	switch strings.ToUpper(apiType) {
	case constants.API_TYPE_REST:
		if len(inputByteArray) > 0 {
			if fileName == "" {
				if strings.HasSuffix(fileName, ".zip") {
					validationResponse = ExtractAndValidateOpenAPIArchive(inputByteArray, returnContent)
				} else {
					// Assume UTF-8 by default
					openAPIContent := string(inputByteArray)
					validationResponse = ValidateAPIDefinition(openAPIContent, returnContent)
				}
			}
		} else if apiDefinition != "" {
			validationResponse = ValidateAPIDefinition(apiDefinition, returnContent)
		}
	case constants.API_TYPE_GRAPHQL:
		if strings.HasSuffix(fileName, ".graphql") || strings.HasSuffix(fileName, ".txt") ||
			strings.HasSuffix(fileName, ".sdl") {
			validationResponse = ValidateGraphQLSchema(string(inputByteArray), returnContent)
		} else {
			return nil, fmt.Errorf("invalid definition file type provided: %s", fileName)
		}
	case constants.API_TYPE_GRPC:
		if len(inputByteArray) > 0 {
			if strings.HasSuffix(fileName, ".zip") || strings.HasSuffix(fileName, ".proto") {
				validationResponse = ValidateGRPCAPIDefinition(inputByteArray)
			} else {
				return nil, fmt.Errorf("invalid definition file type provided: %s", fileName)
			}
		} else {
			return nil, fmt.Errorf("invalid definition file type provided: %s", fileName)
		}
	}
	return validationResponse, nil
}

func GetGRPCAPIFromProtoDefinition(definition []byte, fileName string) (*model.API, error) {
	return nil, nil
}

func GetAPIFromDefinition(definition string, apiType string) (*model.API, error) {
	return nil, nil
}
