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

package handlers

import (
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/api"
	"github.com/wso2/apk/config-deployer-service-go/internal/services"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/lestrrat-go/jwx/v2/jwk"
	// "github.com/wso2/apk/common-go-libs/loggers"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
)

// GetGeneratedAPKConf creates the APK configuration file from api specification.
func GetGeneratedAPKConf(cxt *gin.Context) {
	definitionBody, err := prepareDefinitionBodyFromRequest(cxt)
	var apiDefinitionValidationResponse *dto.APIDefinitionValidationResponse
	var apiType string

	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Failed to parse request: " + err.Error(),
		})
		return
	}

	if (definitionBody.URL == "" && definitionBody.Definition.FileName == "") ||
		(definitionBody.URL != "" && definitionBody.Definition.FileName != "") {
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

	validationService := &services.ValidationService{}
	if definitionBody.URL != "" {
		apiDefinitionValidationResponse, err = validationService.RetrieveAndValidateDefinitionFromURL(apiType,
			definitionBody.URL)
	} else if definitionBody.Definition.FileName != "" && definitionBody.Definition.FileContent != nil &&
		len(definitionBody.Definition.FileContent) > 0 {
		definition := definitionBody.Definition
		apiDefinitionValidationResponse, err = validationService.RetrieveAndValidateDefinitionFromFile(apiType,
			definition.FileName, definition.FileContent)
	} else {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Either URL or file content must be provided",
		})
		return
	}

	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    909022,
			"message": "Error occurred while validating the definition: " + err.Error(),
		})
		return
	}

	if apiDefinitionValidationResponse != nil {
		if apiDefinitionValidationResponse.IsValid {
			var apiFromDefinition *dto.API
			if strings.ToUpper(apiType) == constants.API_TYPE_GRPC {
				var fileName = ""
				if definitionBody.Definition.FileName != "" {
					definition := definitionBody.Definition
					fileName = definition.FileName
				}
				grpcUtil := util.GRPCUtil{}
				apiFromDefinition, err = grpcUtil.GetGRPCAPIFromProtoDefinition(
					apiDefinitionValidationResponse.ProtoContent, fileName)
				if err != nil {
					cxt.JSON(http.StatusInternalServerError, gin.H{
						"code":    909022,
						"message": "Error occurred while retrieving the API from proto definition: " + err.Error(),
					})
					return
				}
			} else {
				runtimeAPIUtil := api.RuntimeAPICommonUtil{}
				apiFromDefinition, err = runtimeAPIUtil.GetAPIFromDefinition(apiDefinitionValidationResponse.Content, apiType)
				if err != nil {
					cxt.JSON(http.StatusInternalServerError, gin.H{
						"code":    909022,
						"message": "Error occurred while retrieving the API from definition: " + err.Error(),
					})
					return
				}
			}
			apiFromDefinition.Type = apiType
			apiClient := &APIClient{}
			generatedAPKConf, err := apiClient.FromAPIModelToAPKConf(apiFromDefinition)
			if err != nil {
				cxt.JSON(http.StatusInternalServerError, gin.H{
					"code":    909022,
					"message": "Error occurred while converting API model to APK conf: " + err.Error(),
				})
				return
			}

			yamlBytes, err := util.MarshalToYAMLWithIndent(generatedAPKConf, 2)
			cxt.Data(http.StatusOK, "application/yaml", yamlBytes)
			return
		} else {
			cxt.JSON(http.StatusBadRequest, gin.H{
				"code":    90091,
				"message": "Invalid API Definition",
			})
			return
		}
	} else {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909022,
			"message": "Error occurred while validating the definition",
		})
		return
	}
}

// prepareDefinitionBodyFromRequest prepares the definition body from the request context.
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
		definitionBody.Definition = dto.FileData{
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

// GetGeneratedK8sResources generates Kubernetes resources based on the provided APK configuration and API definition.
func GetGeneratedK8sResources(cxt *gin.Context, organization *dto.Organization, cpInitiatedParam string, namespace string) {
	var cpInitiated bool
	if cpInitiatedParam != "" && !slices.Contains([]string{"true", "false"}, strings.ToLower(cpInitiatedParam)) {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    909018,
			"message": "Invalid cpInitiated value. It should be either 'true' or 'false'.",
		})
		return
	} else {
		cpInitiated = strings.ToLower(cpInitiatedParam) == "true"
	}

	definitionBody, err := prepareGenerateK8sResourcesBodyFromRequest(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Failed to parse request: " + err.Error(),
		})
		return
	}
	if definitionBody.APKConfiguration.FileName == "" || definitionBody.DefinitionFile.FileName == "" {
		cxt.JSON(http.StatusNotAcceptable, gin.H{
			"code":    909018,
			"message": "required apkConfiguration and definitionFile are not provided",
		})
		return
	}

	apiClient := &APIClient{}
	apkConf, apiDefinition, err := apiClient.PrepareArtifact(definitionBody.APKConfiguration, definitionBody.DefinitionFile)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Error while validating the definition: " + err.Error(),
		})
		return
	}

	apiArtifact, err := apiClient.GenerateK8sArtifacts(apkConf, apiDefinition, organization, cpInitiated, namespace)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909052,
			"message": "Error while generating k8s artifact: " + err.Error(),
		})
		return
	}

	zipName, err := apiClient.ZipAPIArtifact(apiArtifact)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909052,
			"message": "Error while generating k8s artifact zip: " + err.Error(),
		})
		return
	}

	cxt.Header("Content-Disposition", "attachment; filename="+zipName[0])
	cxt.Header("Content-Type", "application/zip")
	cxt.File(zipName[1])
	cxt.Status(http.StatusOK)
}

// prepareGenerateK8sResourcesBodyFromRequest prepares the definition body from the request context.
func prepareGenerateK8sResourcesBodyFromRequest(cxt *gin.Context) (*dto.GenerateK8sResourcesBody, error) {
	definitionBody := &dto.GenerateK8sResourcesBody{}

	// Parse the multipart form with a max memory of 10MB
	if err := cxt.Request.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	// Parse apkConfiguration file
	apkConfFile, apkConfHeader, err := cxt.Request.FormFile("apkConfiguration")
	if err == nil {
		err := apkConfFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close APK Configuration file: %w", err)
		}
		fileContent, readErr := io.ReadAll(apkConfFile)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read APK Configuration file: %w", readErr)
		}
		definitionBody.APKConfiguration = dto.FileData{
			FileName:    apkConfHeader.Filename,
			FileContent: fileContent,
		}
	}
	// Parse definitionFile
	defFile, defHeader, err := cxt.Request.FormFile("definitionFile")
	if err == nil {
		err := defFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close API definition file: %w", err)
		}
		fileContent, readErr := io.ReadAll(defFile)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read definitionFile: %w", readErr)
		}
		definitionBody.DefinitionFile = dto.FileData{
			FileName:    defHeader.Filename,
			FileContent: fileContent,
		}
	}
	if apiType := cxt.PostForm("apiType"); apiType != "" {
		definitionBody.APIType = apiType
	}

	return definitionBody, nil
}
