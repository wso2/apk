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
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/model"
	"sort"
	"strings"
)

type APKConfUtil struct{}

// GetAPKConf parses the APK configuration from the provided APK content string.
func GetAPKConf(apkContent dto.FileData) (*model.APKConf, error) {
	var apkConf *model.APKConf = nil
	apkConfContent := string(apkContent.FileContent)
	apkConfJson, err := YamlToJSON(apkConfContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}
	if apkConfJson != "" {
		err := json.Unmarshal([]byte(apkConfJson), &apkConf)
		if err != nil {
			return nil, fmt.Errorf("failed to parse APK configuration: %w", err)
		}
	}
	if apkConf == nil {
		return nil, fmt.Errorf("apkConfiguration is not provided")
	}
	return apkConf, nil
}

// GetAPIName generates a unique name for the API based on its type.
func (apkConfUtil *APKConfUtil) GetAPIName(apiName string, apiType string) string {
	if strings.ToUpper(apiType) == constants.API_TYPE_GRPC {
		return apkConfUtil.GetUniqueNameForGrpcApi(apiName)
	}
	return apiName
}

// GetUniqueNameForGrpcApi generates a unique name for gRPC APIs by hashing the concatenated service names.
func (apkConfUtil *APKConfUtil) GetUniqueNameForGrpcApi(concatanatedServices string) string {
	hasher := sha1.New()
	hasher.Write([]byte(concatanatedServices))
	hashedValue := hasher.Sum(nil)
	return hex.EncodeToString(hashedValue)
}

// GetUniqueIdForAPI generates a unique identifier for an API based on its name, version, and organization.
func (apkConfUtil *APKConfUtil) GetUniqueIdForAPI(name string, version string, organization *dto.Organization) string {
	uniqueString := fmt.Sprintf("%s-%s-%s", organization.Name, name, version)
	hasher := sha1.New()
	hasher.Write([]byte(uniqueString))
	hashedValue := hasher.Sum(nil)
	return hex.EncodeToString(hashedValue)
}

// GetResourceLevelEndpointConfig extracts endpoint configurations from a list of APK operations.
func (apkConfUtil *APKConfUtil) GetResourceLevelEndpointConfig(operations []model.APKOperations) []model.EndpointConfigurations {
	endpointConfigurationsList := make([]model.EndpointConfigurations, 0)
	for _, operation := range operations {
		endpointConfigurations := operation.EndpointConfigurations
		if endpointConfigurations != nil {
			endpointConfigurationsList = append(endpointConfigurationsList, *endpointConfigurations)
		}
	}
	return endpointConfigurationsList
}

// CreateAPIResourceBundle creates a resource bundle for the API artifact based on the APK configuration and definition.
func (apkConfUtil *APKConfUtil) CreateAPIResourceBundle(apkConf *model.APKConf, organization *dto.Organization,
	cpInitiated bool, namespace, definition string) dto.APIResourceBundle {
	apiResourceBundle := dto.APIResourceBundle{
		Organization: organization.Name,
		Namespace:    namespace,
		CPInitiated:  cpInitiated,
		APKConf:      apkConf,
		Definition:   definition,
	}
	// bundle apk operations into bins based on secured, rate limit and scopes and create combined resources
	combinedResources := generateCombinedResources(apkConf)
	apiResourceBundle.CombinedResources = combinedResources
	return apiResourceBundle
}

// generateCombinedResources groups APKOperations into CombinedResource buckets based on
// unique combinations of Secured, RateLimit, and Scopes attributes
func generateCombinedResources(apkConf *model.APKConf) []dto.CombinedResource {
	groupMap := make(map[string][]model.APKOperations)
	operations := apkConf.Operations
	for _, operation := range operations {
		populatedOperation := populateEndpointConfigurations(operation, apkConf)
		key := generateGroupingKey(populatedOperation)
		groupMap[key] = append(groupMap[key], populatedOperation)
	}

	var combinedResources []dto.CombinedResource
	for _, groupedOperations := range groupMap {
		combinedResource := dto.CombinedResource{
			APKOperations: groupedOperations,
		}
		combinedResources = append(combinedResources, combinedResource)
	}
	return combinedResources
}

// populateEndpointConfigurations populates missing endpoint configurations from APKConf level
func populateEndpointConfigurations(operation model.APKOperations, apkConf *model.APKConf) model.APKOperations {
	populatedOperation := operation
	// If operation doesn't have endpoint configurations but APKConf does
	if populatedOperation.EndpointConfigurations == nil {
		populatedOperation.EndpointConfigurations = apkConf.EndpointConfigurations
	} else if apkConf.EndpointConfigurations != nil {
		// If operation has endpoint configurations but some are missing, populate from APKConf
		if len(populatedOperation.EndpointConfigurations.Production) == 0 && len(apkConf.EndpointConfigurations.Production) > 0 {
			populatedOperation.EndpointConfigurations.Production = apkConf.EndpointConfigurations.Production
		}
		if len(populatedOperation.EndpointConfigurations.Sandbox) == 0 && len(apkConf.EndpointConfigurations.Sandbox) > 0 {
			populatedOperation.EndpointConfigurations.Sandbox = apkConf.EndpointConfigurations.Sandbox
		}
	}

	return populatedOperation
}

// generateGroupingKey creates a unique key for grouping operations based on
// Secured, RateLimit, and Scopes attributes
func generateGroupingKey(operation model.APKOperations) string {
	var keyParts []string

	// Handle Secured field
	if operation.Secured != nil {
		keyParts = append(keyParts, fmt.Sprintf("secured:%t", *operation.Secured))
	} else {
		keyParts = append(keyParts, "secured:false")
	}

	// Handle RateLimit field
	if operation.RateLimit != nil {
		keyParts = append(keyParts, fmt.Sprintf("rateLimit:%d-%s", operation.RateLimit.RequestsPerUnit, operation.RateLimit.Unit))
	} else {
		keyParts = append(keyParts, "rateLimit:nil")
	}

	// Handle Scopes field
	if len(operation.Scopes) > 0 {
		sortedScopes := make([]string, len(operation.Scopes))
		copy(sortedScopes, operation.Scopes)
		sort.Strings(sortedScopes)
		keyParts = append(keyParts, fmt.Sprintf("scopes:%s", strings.Join(sortedScopes, ",")))
	} else {
		keyParts = append(keyParts, "scopes:empty")
	}

	return strings.Join(keyParts, "|")
}
