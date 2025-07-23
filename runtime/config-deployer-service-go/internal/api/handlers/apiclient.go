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
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"strings"
)

// APIClient represents the API client
type APIClient struct{}

// FromAPIModelToAPKConf converts APKInternalAPI model to APKConf
func (client *APIClient) FromAPIModelToAPKConf(api *dto.API) (*dto.APKConf, error) {
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

	apkConf := &dto.APKConf{
		Name:                   getAPIName(api.Name, api.Type),
		BasePath:               basePath,
		Version:                api.Version,
		Type:                   apiType,
		DefaultVersion:         false,
		SubscriptionValidation: false,
	}

	endpoint := api.Endpoint
	if len(endpoint) > 0 {
		apkConf.EndpointConfigurations = &dto.EndpointConfigurations{
			Production: []dto.EndpointConfiguration{
				{Endpoint: endpoint},
			},
		}
	}

	uriTemplates := api.URITemplates
	var operations []dto.APKOperations
	for _, uriTemplate := range uriTemplates {
		operation := dto.APKOperations{
			Target:  uriTemplate.URITemplate,
			Verb:    uriTemplate.Verb,
			Secured: uriTemplate.AuthEnabled,
			Scopes:  uriTemplate.Scopes,
		}
		resourceEndpoint := uriTemplate.Endpoint
		if len(resourceEndpoint) > 0 {
			operation.EndpointConfigurations = &dto.EndpointConfigurations{
				Production: []dto.EndpointConfiguration{
					{Endpoint: resourceEndpoint},
				},
			}
		}
		operations = append(operations, operation)
	}
	apkConf.Operations = operations
	return apkConf, nil
}

func getAPIName(apiName string, apiType string) string {
	if strings.ToUpper(apiType) == constants.API_TYPE_GRPC {
		return getUniqueNameForGrpcApi(apiName)
	}
	return apiName
}

func getUniqueNameForGrpcApi(concatanatedServices string) string {
	hasher := sha1.New()
	hasher.Write([]byte(concatanatedServices))
	hashedValue := hasher.Sum(nil)
	return hex.EncodeToString(hashedValue)
}
