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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	constantscommon "github.com/wso2/apk/common-go-libs/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/config"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"io"
	"net/http"
	"strings"
)

// HandleAPIDeployment handles the deployment of an API based on the provided APK configuration and API definition.
func HandleAPIDeployment(cxt *gin.Context, organization *dto.Organization, cpInitiatedParam string, namespace string) {
	deployAPIBody, err := prepareDeployAPIBodyFromRequest(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{
			"code":    90091,
			"message": "Failed to parse request: " + err.Error(),
		})
		return
	}
	if deployAPIBody.APKConfiguration.FileName == "" || deployAPIBody.DefinitionFile.FileName == "" {
		cxt.JSON(http.StatusNotAcceptable, gin.H{
			"code":    909017,
			"message": "Invalid API request, required apkConfiguration and definitionFile are not provided",
		})
		return
	}
	apiClient := &APIClient{}
	apiArtifact, err := apiClient.PrepareArtifact(deployAPIBody.APKConfiguration,
		deployAPIBody.DefinitionFile, organization, strings.ToLower(cpInitiatedParam) == "true", namespace)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909052,
			"message": "Error while generating k8s artifact: " + err.Error(),
		})
		return
	}
	k8sClient := config.GetManager().GetClient()
	routeMetadata, err := apiClient.DeployAPIToK8s(apiArtifact, namespace, k8sClient)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909028,
			"message": "Internal error occurred while deploying API: " + err.Error(),
		})
		return
	}
	apkConf, err := util.GetAPKConf(deployAPIBody.APKConfiguration)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909022,
			"message": "Failed to parse APK configuration: " + err.Error(),
		})
		return
	}
	apkConf.ID = routeMetadata.Labels[constantscommon.LabelKGWUUID]
	apkYaml, err := util.MarshalToYAMLWithIndent(apkConf, 2)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909022,
			"message": "Error occurred while converting APKConf to YAML: " + err.Error(),
		})
		return
	}
	cxt.Data(http.StatusOK, "application/yaml", apkYaml)
}

// HandleAPIUndeployment handles the undeployment of an API by removing its associated Kubernetes resources.
func HandleAPIUndeployment(cxt *gin.Context, apiId string, organization *dto.Organization, namespace string) {
	apiClient := &APIClient{}
	k8sClient := config.GetManager().GetClient()
	routeMetadataList, err := util.GetRouteMetadataList(apiId, namespace, k8sClient)
	if err != nil {
		cxt.JSON(http.StatusNotFound, gin.H{
			"code":    909001,
			"message": apiId + " not found: " + err.Error(),
		})
		return
	}
	err = apiClient.UndeployAPI(routeMetadataList, namespace, k8sClient)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{
			"code":    909022,
			"message": "Error while undeploying API: " + err.Error(),
		})
		return
	}
	response := fmt.Sprintf("API with id %s undeployed successfully", apiId)
	jsonResponse := map[string]interface{}{
		"status": response,
	}
	jsonBytes, err := json.Marshal(jsonResponse)
	cxt.Data(http.StatusOK, "application/json", jsonBytes)
	return
}

// prepareDeployAPIBodyFromRequest prepares the definition body from the request context.
func prepareDeployAPIBodyFromRequest(cxt *gin.Context) (*dto.DeployAPIBody, error) {
	deployAPIBody := &dto.DeployAPIBody{}

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
		deployAPIBody.APKConfiguration = dto.FileData{
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
		deployAPIBody.DefinitionFile = dto.FileData{
			FileName:    defHeader.Filename,
			FileContent: fileContent,
		}
	}

	return deployAPIBody, nil
}
