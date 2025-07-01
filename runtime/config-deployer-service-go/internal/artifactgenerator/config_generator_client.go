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

package artifactgenerator

import (
	"config-deployer-service-go/internal/model"
	"config-deployer-service-go/internal/util"
	// "crypto/sha256"
	// "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/lestrrat-go/jwx/v2/jwk"
	// "github.com/wso2/apk/common-go-libs/loggers"
	// "config-deployer-service-go/internal/config"
	// "github.com/wso2/apk/gateway/enforcer/internal/util"
	"config-deployer-service-go/internal/constants"
	"config-deployer-service-go/internal/dto"
)

// GetGeneratedAPKConf creates the APK configuration file from api specification.
func GetGeneratedAPKConf(cxt *gin.Context) {
	definitionBody, err := prepareDefinitionBodyFromRequest(cxt)
	var validateAndRetrieveDefinitionResult *dto.APIDefinitionValidationResponse
	var apiType string

	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Failed to parse request: " + err.Error(),
		})
		return
	}

	if definitionBody.Definition.FileName == "" && definitionBody.URL == "" {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Specify either definition or url",
		})
		return
	}

	if definitionBody.APIType == "" {
		// Setting the default API type as REST.
		apiType = constants.API_TYPE_REST
	} else {
		apiType = definitionBody.APIType
	}

	if !slices.Contains(constants.ALLOWED_API_TYPES, strings.ToUpper(apiType)) {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Invalid API Type",
		})
		return
	}

	if definitionBody.URL != "" {
		validateAndRetrieveDefinitionResult, err = validateAndRetrieveDefinition(apiType, definitionBody.URL,
			nil, "")
	} else if definitionBody.Definition.FileName != "" && len(definitionBody.Definition.FileContent) > 0 {
		definition := definitionBody.Definition
		validateAndRetrieveDefinitionResult, err = validateAndRetrieveDefinition(apiType, "",
			definition.FileContent, definition.FileName)
	}

	if validateAndRetrieveDefinitionResult != nil {
		if validateAndRetrieveDefinitionResult.IsValid {
			var apiFromDefinition *model.API
			if strings.ToUpper(apiType) == constants.API_TYPE_GRPC {
				var fileName = ""
				if definitionBody.Definition.FileName != "" {
					definition := definitionBody.Definition
					fileName = definition.FileName
				}
				apiFromDefinition, err = util.GetGRPCAPIFromProtoDefinition(
					validateAndRetrieveDefinitionResult.ProtoContent, fileName)
				if err != nil {
					cxt.JSON(http.StatusInternalServerError, gin.H{
						"code":    909022,
						"message": "Error occurred while retrieving the API from proto definition: " + err.Error(),
					})
					return
				}
			} else {
				apiFromDefinition, err = util.GetAPIFromDefinition(validateAndRetrieveDefinitionResult.Content, apiType)
				if err != nil {
					cxt.JSON(http.StatusInternalServerError, gin.H{
						"code":    909022,
						"message": "Error occurred while retrieving the API from definition: " + err.Error(),
					})
					return
				}
			}
			apiFromDefinition.Type = apiType
		}
	} else {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909022,
			"message": "Error occurred while validating the definition",
		})
		return
	}
}

func prepareDefinitionBodyFromRequest(cxt *gin.Context) (*dto.DefinitionBody, error) {
	definitionBody := &dto.DefinitionBody{}

	// Parse the multipart form with a max memory of 10MB
	if err := cxt.Request.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	defFile, defHeader, err := cxt.Request.FormFile("definition")
	if err == nil {
		err := defFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close definition file: %w", err)
		}
		fileContent, readErr := io.ReadAll(defFile)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read uploaded file: %w", readErr)
		}
		definitionBody.Definition = dto.Definition{
			FileName:    defHeader.Filename,
			FileContent: fileContent,
		}
	}
	if url := cxt.PostForm("url"); url != "" {
		definitionBody.URL = url
	}
	if apiType := cxt.PostForm("apiType"); apiType != "" {
		definitionBody.APIType = apiType
	}

	return definitionBody, nil
}

func validateAndRetrieveDefinition(apiType, url string, content []byte,
	fileName string) (*dto.APIDefinitionValidationResponse, error) {
	if url != "" {
		definition, err := retrieveDefinitionFromUrl(url)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve definition from URL: %w", err)
		}
		return util.ValidateOpenAPIDefinition(apiType, nil, definition, "", true)
	}
	if fileName != "" && len(content) > 0 {
		return util.ValidateOpenAPIDefinition(apiType, content, "", fileName, true)
	}
	return nil, fmt.Errorf("either URL or file content must be provided")
}

func retrieveDefinitionFromUrl(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error occurred while retrieving the definition from the url: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error occurred while closing the response body: %v\n", err)
		}
	}(response.Body)
	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("error occurred while reading the definition from the url: %w", err)
		}
		return string(body), nil
	} else {
		return "", fmt.Errorf("error occurred while retrieving the definition from the url: %s. Status code:"+
			" %d", url, response.StatusCode)
	}
}
