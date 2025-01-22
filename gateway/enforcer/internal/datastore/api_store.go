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

package datastore

import (
	"sync"

	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// APIStore is a thread-safe store for APIs.
type APIStore struct {
	apis        map[string]*requestconfig.API
	mu          sync.RWMutex
	configStore *ConfigStore
}

// NewAPIStore creates a new instance of APIStore.
func NewAPIStore(configStore *ConfigStore) *APIStore {
	return &APIStore{
		configStore: configStore,
		// apis: make(map[string]*api.Api, 0),
	}
}

// AddAPIs adds a list of APIs to the store.
// This method is thread-safe.
func (s *APIStore) AddAPIs(apis []*api.Api) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apis = make(map[string]*requestconfig.API, len(apis))
	for _, api := range apis {
		customAPI := requestconfig.API{
			Name:                  api.Title,
			Version:               api.Version,
			Vhost:                 api.Vhost,
			BasePath:              api.BasePath,
			APIType:               api.ApiType,
			EnvType:               api.EnvType,
			APILifeCycleState:     api.ApiLifeCycleState,
			AuthorizationHeader:   "", // You might want to set this field if applicable
			OrganizationID:        api.OrganizationId,
			UUID:                  api.Id,
			Tier:                  api.Tier,
			DisableAuthentication: api.DisableAuthentications,
			DisableScopes:         api.DisableScopes,
			Resources:             make([]requestconfig.Resource, 0),
			IsMockedAPI:           false, // You can add logic to determine if the API is mocked
			MutualSSL:             api.MutualSSL,
			TransportSecurity:     api.TransportSecurity,
			ApplicationSecurity:   api.ApplicationSecurity,
			// JwtConfigurationDto:    convertBackendJWTTokenInfoToJWTConfig(api.BackendJWTTokenInfo),
			SystemAPI:              api.SystemAPI,
			APIDefinition:          api.ApiDefinitionFile,
			Environment:            api.Environment,
			SubscriptionValidation: api.SubscriptionValidation,
			// Endpoints:              api.Endpoints,
			// EndpointSecurity:       convertSecurityInfoToEndpointSecurity(api.EndpointSecurity),
			// AiProvider:             api.Aiprovider,

		}
		for _, resource := range api.Resources {
			for _, operation := range resource.Methods {
				resource := buildResource(operation, resource.Path, func() []*requestconfig.EndpointSecurity {
					endpointSecurity := make([]*requestconfig.EndpointSecurity, len(resource.EndpointSecurity))
					for i, es := range resource.EndpointSecurity {
						endpointSecurity[i] = &requestconfig.EndpointSecurity{
							Password:         es.Password,
							Enabled:          es.Enabled,
							Username:         es.Username,
							SecurityType:     es.SecurityType,
							CustomParameters: es.CustomParameters,
						}
					}
					return endpointSecurity
				}())
				customAPI.Resources = append(customAPI.Resources, resource)
			}
		}

		s.apis[util.PrepareAPIKey(api.Vhost, api.BasePath, api.Version)] = &customAPI
	}
}

// GetAPIs retrieves the list of APIs from the store.
// This method is thread-safe.
func (s *APIStore) GetAPIs() map[string]*requestconfig.API {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apis
}

// GetMatchedAPI retrieves the API that matches the given API key.
func (s *APIStore) GetMatchedAPI(apiKey string) *requestconfig.API {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apis[apiKey]
}
