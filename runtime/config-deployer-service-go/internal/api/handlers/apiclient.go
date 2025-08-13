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
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/crbuilder"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"github.com/wso2/apk/config-deployer-service-go/internal/services"
	"github.com/wso2/apk/config-deployer-service-go/internal/services/validators"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"os"
	"strings"
)

// APIClient represents the API client
type APIClient struct{}

// FromAPIModelToAPKConf converts APKInternalAPI model to APKConf
func (client *APIClient) FromAPIModelToAPKConf(api *dto.API) (*model.APKConf, error) {
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
func (client *APIClient) PrepareArtifact(apkConfiguration dto.FileData, definitionFile dto.FileData,
	organization *dto.Organization, cpInitiated bool, namespace string) (*dto.APIArtifact, error) {

	var apkConf *model.APKConf = nil
	apkContent := string(apkConfiguration.FileContent)
	convertedJson, err := util.YamlToJSON(apkContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}
	if convertedJson != "" {
		apkConf, err = services.ValidateAndRetrieveAPKConfiguration(convertedJson)
		if err != nil {
			return nil, fmt.Errorf("failed to validate APK configuration: %w", err)
		}
	}
	if apkConf == nil {
		return nil, fmt.Errorf("apkConfiguration is not provided")
	}
	var apiDefinition string
	apiType := apkConf.Type
	if apiType == constants.API_TYPE_REST {
		definitionFileContent := string(definitionFile.FileContent)
		if strings.HasSuffix(definitionFile.FileName, ".yaml") {
			apiDefinition, err = util.YamlToJSON(definitionFileContent)
			if err != nil {
				return nil, fmt.Errorf("invalid API definiton provided. Failed to convert YAML definition to JSON: %w", err)
			}
		} else if strings.HasSuffix(definitionFile.FileName, ".json") {
			apiDefinition = definitionFileContent
		} else {
			return nil, fmt.Errorf("invalid REST API definition file type provided: %s", definitionFile.FileName)
		}
	} else if apiType == constants.API_TYPE_GRAPHQL {
		apiDefinition = string(definitionFile.FileContent)
	} else if apiType == constants.API_TYPE_GRPC {
		if strings.HasSuffix(definitionFile.FileName, ".zip") {
			apiDefinition = base64.StdEncoding.EncodeToString(definitionFile.FileContent)
		} else if strings.HasSuffix(definitionFile.FileName, ".proto") {
			apiDefinition = string(definitionFile.FileContent)
		} else {
			return nil, fmt.Errorf("invalid gRPC API definition file type provided: %s", definitionFile.FileName)
		}
	}

	return GenerateK8sArtifacts(apkConf, apiDefinition, organization, cpInitiated, namespace)
}

// GenerateK8sArtifacts generates Kubernetes artifacts based on the APK configuration and API definition.
func GenerateK8sArtifacts(apkConf *model.APKConf, definition string, organization *dto.Organization,
	cpInitiated bool, namespace string) (*dto.APIArtifact, error) {
	apkConfValidator := &validators.APKConfValidator{}
	apkConfUtil := util.APKConfUtil{}
	//uniqueId := apkConfUtil.GetUniqueIdForAPI(apkConf.Name, apkConf.Version, organization)
	//if apkConf.ID != "" {
	//	uniqueId = apkConf.ID
	//}
	var resourceLevelEndpointConfigList []model.EndpointConfigurations
	operations := apkConf.Operations
	if operations != nil {
		if len(operations) == 0 {
			return nil, fmt.Errorf("atleast one operation need to specified")
		}
		err := apkConfValidator.ValidateRateLimit(apkConf.RateLimit, operations)
		if err != nil {
			return nil, fmt.Errorf("failed to validate rate limit: %w", err)
		}
		resourceLevelEndpointConfigList = apkConfUtil.GetResourceLevelEndpointConfig(operations)
	} else {
		return nil, fmt.Errorf("atleast one operation need to specified")
	}
	_ = resourceLevelEndpointConfigList
	var createdEndpoints map[string][]*dto.Endpoint
	var err error
	endpointConfigurations := apkConf.EndpointConfigurations
	//if endpointConfigurations != nil {
	//	createdEndpoints, err = apkConfUtil.CreateAndAddBackendServices(apiArtifact, apkConf, endpointConfigurations,
	//		nil, nil, organization)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to create and add backend services: %w", err)
	//	}
	//}
	_ = createdEndpoints
	_ = endpointConfigurations
	// TODO - aiRateLimit - create AIRateLimitPolicies in RoutePolicy CR and BackendTrafficPolicy CR and attach to httproutes
	// TODO - EndpointSecurity - Create BackendJWT Policy in RoutePolicy CR and attach to httproutes
	// TODO - Handle Resiliency - apply to all httproutes that use this backend service using BackendTrafficPolicy CR
	apiResourceBundle := apkConfUtil.CreateAPIResourceBundle(apkConf, organization, cpInitiated, namespace, definition)
	k8sArtifacts, err := crbuilder.CreateResources(&apiResourceBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes resources: %w", err)
	}
	apiArtifact := &dto.APIArtifact{
		Name:         apkConf.Name,
		Version:      apkConf.Version,
		K8sArtifacts: k8sArtifacts,
	}
	return apiArtifact, nil
}

// ZipAPIArtifact creates a zip file containing all API artifact resources
func (client *APIClient) ZipAPIArtifact(apiArtifact *dto.APIArtifact) ([2]string, error) {
	// Create temporary directory
	zipDir, err := util.CreateTempDir()
	if err != nil {
		return [2]string{}, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(zipDir) // Clean up temp directory

	//definition := apiArtifact.Definition
	//if definition != nil {
	//	yamlString, err := util.MarshalToYAMLWithIndent(definition, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert definition to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), definition.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store definition file: %w", err)
	//	}
	//}
	//for _, authenticationCr := range apiArtifact.AuthenticationMap {
	//	yamlString, err := util.MarshalToYAMLWithIndent(authenticationCr, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert authentication CR to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), authenticationCr.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store authentication CR file: %w", err)
	//	}
	//}
	//for _, httpRoute := range apiArtifact.ProductionHttpRoutes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(httpRoute, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert HTTP route to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), httpRoute.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store HTTP route file: %w", err)
	//	}
	//}
	//for _, httpRoute := range apiArtifact.SandboxHttpRoutes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(httpRoute, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert HTTP route to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), httpRoute.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store HTTP route file: %w", err)
	//	}
	//}
	//for _, gqlRoute := range apiArtifact.ProductionGqlRoutes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(gqlRoute, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert GraphQL route to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), gqlRoute.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store GraphQL route file: %w", err)
	//	}
	//}
	//for _, gqlRoute := range apiArtifact.SandboxGqlRoutes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(gqlRoute, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert GraphQL route to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), gqlRoute.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store GraphQL route file: %w", err)
	//	}
	//}
	//for _, grpcRoute := range apiArtifact.ProductionGrpcRoutes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(grpcRoute, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert gRPC route to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), grpcRoute.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store gRPC route file: %w", err)
	//	}
	//}
	//for _, grpcRoute := range apiArtifact.SandboxGrpcRoutes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(grpcRoute, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert gRPC route to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), grpcRoute.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store gRPC route file: %w", err)
	//	}
	//}
	//for _, backend := range apiArtifact.BackendServices {
	//	yamlString, err := util.MarshalToYAMLWithIndent(backend, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert backend service to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), backend.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store backend service file: %w", err)
	//	}
	//}
	//for _, scope := range apiArtifact.Scopes {
	//	yamlString, err := util.MarshalToYAMLWithIndent(scope, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert scope to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), scope.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store scope file: %w", err)
	//	}
	//}
	//for _, rateLimitPolicy := range apiArtifact.RateLimitPolicies {
	//	yamlString, err := util.MarshalToYAMLWithIndent(rateLimitPolicy, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert rate limit policy to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), rateLimitPolicy.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store rate limit policy file: %w", err)
	//	}
	//}
	//for _, apiPolicy := range apiArtifact.ApiPolicies {
	//	yamlString, err := util.MarshalToYAMLWithIndent(apiPolicy, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert API policy to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), apiPolicy.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store API policy file: %w", err)
	//	}
	//}
	//for _, interceptorService := range apiArtifact.InterceptorServices {
	//	yamlString, err := util.MarshalToYAMLWithIndent(interceptorService, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert interceptor service to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), interceptorService.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store interceptor service file: %w", err)
	//	}
	//}
	//backendJWT := apiArtifact.BackendJwt
	//if backendJWT != nil {
	//	yamlString, err := util.MarshalToYAMLWithIndent(backendJWT, 2)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to convert backend JWT to YAML: %w", err)
	//	}
	//	err = util.StoreFile(string(yamlString), backendJWT.ObjectMeta.Name, zipDir)
	//	if err != nil {
	//		return [2]string{}, fmt.Errorf("failed to store backend JWT file: %w", err)
	//	}
	//}

	for _, artifact := range apiArtifact.K8sArtifacts {
		yamlString, err := util.MarshalToYAMLWithIndent(artifact, 2)
		if err != nil {
			return [2]string{}, fmt.Errorf("failed to convert artifact to YAML: %w", err)
		}
		fileName := artifact.GetObjectKind().GroupVersionKind().Kind + "-" + artifact.GetName()
		err = util.StoreFile(string(yamlString), fileName, zipDir)
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
