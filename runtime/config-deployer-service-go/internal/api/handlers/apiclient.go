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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/crbuilder"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"github.com/wso2/apk/config-deployer-service-go/internal/services"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// APIClient represents the API client
type APIClient struct{}

// FromAPIModelToAPKConf converts APKInternalAPI model to APKConf
func (apiClient *APIClient) FromAPIModelToAPKConf(api *dto.API) (*model.APKConf, error) {
	apkConfUtil := util.APKConfUtil{}
	generatedBasePath := api.Name + api.Version
	data := []byte(generatedBasePath)
	encodedString := "/" + base64.StdEncoding.EncodeToString(data)
	if strings.HasSuffix(encodedString, "==") {
		encodedString = encodedString[:len(encodedString)-2]
	} else if strings.HasSuffix(encodedString, "=") {
		encodedString = encodedString[:len(encodedString)-1]
	}
	apiType := api.Type
	if apiType == "" {
		apiType = constants.API_TYPE_REST
	} else {
		apiType = strings.ToUpper(apiType)
	}
	basePath := api.BasePath
	if len(basePath) == 0 {
		basePath = encodedString
	}

	apkConf := &model.APKConf{
		Name:                   apkConfUtil.GetAPIName(api.Name, api.Type),
		BasePath:               basePath,
		Version:                api.Version,
		Type:                   apiType,
		DefaultVersion:         false,
		SubscriptionValidation: false,
	}

	endpoint := api.Endpoint
	if len(endpoint) > 0 {
		apkConf.EndpointConfigurations = &model.EndpointConfigurations{
			Production: []model.EndpointConfiguration{
				{Endpoint: endpoint},
			},
		}
	}

	uriTemplates := api.URITemplates
	var operations []model.APKOperations
	for _, uriTemplate := range uriTemplates {
		operation := model.APKOperations{
			Target:  &uriTemplate.URITemplate,
			Verb:    &uriTemplate.Verb,
			Secured: &uriTemplate.AuthEnabled,
			Scopes:  uriTemplate.Scopes,
		}
		resourceEndpoint := uriTemplate.Endpoint
		if len(resourceEndpoint) > 0 {
			operation.EndpointConfigurations = &model.EndpointConfigurations{
				Production: []model.EndpointConfiguration{
					{Endpoint: resourceEndpoint},
				},
			}
		}
		operations = append(operations, operation)
	}
	apkConf.Operations = operations
	return apkConf, nil
}

// PrepareArtifact creates the API artifact based on the provided configuration.
func (apiClient *APIClient) PrepareArtifact(apkConfiguration dto.FileData, definitionFile dto.FileData) (*model.APKConf,
	string, error) {

	var apkConf *model.APKConf = nil
	apkContent := string(apkConfiguration.FileContent)
	convertedJson, err := util.YamlToJSON(apkContent)
	if err != nil {
		return nil, "", fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}
	if convertedJson != "" {
		apkConf, err = services.ValidateAndRetrieveAPKConfiguration(convertedJson)
		if err != nil {
			return nil, "", fmt.Errorf("failed to validate APK configuration: %w", err)
		}
	}
	if apkConf == nil {
		return nil, "", fmt.Errorf("apkConfiguration is not provided")
	}
	var apiDefinition string
	apiType := apkConf.Type
	if apiType == constants.API_TYPE_REST {
		definitionFileContent := string(definitionFile.FileContent)
		if strings.HasSuffix(definitionFile.FileName, ".yaml") {
			apiDefinition, err = util.YamlToJSON(definitionFileContent)
			if err != nil {
				return nil, "", fmt.Errorf("invalid API definiton provided. Failed to convert YAML definition to JSON: %w", err)
			}
		} else if strings.HasSuffix(definitionFile.FileName, ".json") {
			apiDefinition = definitionFileContent
		} else {
			return nil, "", fmt.Errorf("invalid REST API definition file type provided: %s", definitionFile.FileName)
		}
	} else if apiType == constants.API_TYPE_GRAPHQL {
		apiDefinition = string(definitionFile.FileContent)
	} else if apiType == constants.API_TYPE_GRPC {
		if strings.HasSuffix(definitionFile.FileName, ".zip") {
			apiDefinition = base64.StdEncoding.EncodeToString(definitionFile.FileContent)
		} else if strings.HasSuffix(definitionFile.FileName, ".proto") {
			apiDefinition = string(definitionFile.FileContent)
		} else {
			return nil, "", fmt.Errorf("invalid gRPC API definition file type provided: %s", definitionFile.FileName)
		}
	}
	return apkConf, apiDefinition, nil
}

// GenerateK8sArtifacts generates Kubernetes artifacts based on the APK configuration and API definition.
func (apiClient *APIClient) GenerateK8sArtifacts(apkConf *model.APKConf, definition string, organization *dto.Organization,
	cpInitiated bool, namespace string) (*dto.APIArtifact, error) {
	apkConfUtil := util.APKConfUtil{}
	uniqueId := apkConfUtil.GetUniqueIdForAPI(apkConf.Name, apkConf.Version, organization)
	if apkConf.ID != "" {
		uniqueId = apkConf.ID
	} else {
		apkConf.ID = uniqueId
	}
	apkConf = apkConfUtil.AddOptionsResources(apkConf)
	apiResourceBundle := apkConfUtil.CreateAPIResourceBundle(apkConf, organization, cpInitiated, namespace, definition)
	k8sArtifacts, err := crbuilder.CreateResources(&apiResourceBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes resources: %w", err)
	}
	apiArtifact := &dto.APIArtifact{
		Name:         apkConf.Name,
		UniqueID:     uniqueId,
		Version:      apkConf.Version,
		K8sArtifacts: k8sArtifacts,
		Organization: organization.Name,
	}
	return apiArtifact, nil
}

// ZipAPIArtifact creates a zip file containing all API artifact resources
func (apiClient *APIClient) ZipAPIArtifact(apiArtifact *dto.APIArtifact) ([2]string, error) {
	// Create temporary directory
	zipDir, err := util.CreateTempDir()
	if err != nil {
		return [2]string{}, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(zipDir) // Clean up temp directory

	for _, artifact := range apiArtifact.K8sArtifacts {
		jsonData, err := json.Marshal(artifact)
		if err != nil {
			return [2]string{}, fmt.Errorf("failed to convert artifact to JSON: %w", err)
		}
		yamlString, err := util.JsonToYaml(string(jsonData))
		if err != nil {
			return [2]string{}, fmt.Errorf("failed to convert artifact to YAML: %w", err)
		}
		fileName := artifact.GetObjectKind().GroupVersionKind().Kind + "-" + artifact.GetName()
		err = util.StoreFile(yamlString, fileName, zipDir)
		if err != nil {
			return [2]string{}, fmt.Errorf("failed to store artifact file: %w", err)
		}
	}

	zipFileName := fmt.Sprintf("%s-%s", apiArtifact.Name, apiArtifact.Version)
	zipName, err := util.ZipDirectory(zipFileName, zipDir)
	if err != nil {
		return [2]string{}, fmt.Errorf("failed to create zip: %w", err)
	}

	return zipName, nil
}

// DeployAPIToK8s deploys the API artifact to Kubernetes and returns the RouteMetadata.
func (apiClient *APIClient) DeployAPIToK8s(apiArtifact *dto.APIArtifact, namespace string,
	k8sClient client.Client) (*v2alpha1.RouteMetadata, error) {
	routeMetadataList, err := util.GetRouteMetadataList(apiArtifact.UniqueID, namespace, k8sClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get RouteMetadata list: %w", err)
	}
	if routeMetadataList != nil && len(routeMetadataList.Items) > 0 {
		for _, routeMetadata := range routeMetadataList.Items {
			unsuedObjectList, err := util.GetCRsUsedByRouteMetadataNotInAPIArtifact(routeMetadata, apiArtifact,
				namespace, k8sClient)
			if err != nil {
				return nil, fmt.Errorf("failed to get unused objects for RouteMetadata %s: %w", routeMetadata.Name, err)
			}
			for _, unusedObject := range unsuedObjectList.Items {
				if err = util.UndeployCR(k8sClient, unusedObject); err != nil {
					return nil, fmt.Errorf("failed to undeploy CR %s: %w", unusedObject.GetName(), err)
				}
			}
		}
	}

	var createdRouteMetadata *v2alpha1.RouteMetadata
	for _, k8sArtifact := range apiArtifact.K8sArtifacts {
		if err = util.ApplyK8sResource(k8sClient, namespace, k8sArtifact); err != nil {
			return nil, fmt.Errorf("failed to apply k8s resource %s: %w", k8sArtifact.GetName(), err)
		}
		if routeMetadata, ok := k8sArtifact.(*v2alpha1.RouteMetadata); ok {
			createdRouteMetadata = routeMetadata
		}
	}
	return createdRouteMetadata, nil
}

// UndeployAPI removes all RouteMetadata Custom Resource from the Kubernetes cluster based on API ID label.
func (apiClient *APIClient) UndeployAPI(routeMetadataList *v2alpha1.RouteMetadataList, namespace string,
	k8sClient client.Client) error {
	//conf, errReadConfig := config.ReadConfigs()
	//if errReadConfig != nil {
	//	return errReadConfig
	//}
	for _, routeMetadata := range routeMetadataList.Items {
		if err := util.UndeployK8sRouteMetadataCR(k8sClient, routeMetadata); err != nil {
			return fmt.Errorf("unable to delete RouteMetadata CRs: %w", err)
		}
		filteredLabels := util.GetFilteredLabels(routeMetadata.GetLabels())
		objectList, err := util.GetCRsFromLabels(filteredLabels, namespace, k8sClient)
		if err != nil {
			return fmt.Errorf("unable to get objects with labels %v: %w", filteredLabels, err)
		}
		for _, object := range objectList.Items {
			if err = util.UndeployCR(k8sClient, object); err != nil {
				return fmt.Errorf("unable to delete CRs: %w", err)
			}
		}
	}
	return nil
}
