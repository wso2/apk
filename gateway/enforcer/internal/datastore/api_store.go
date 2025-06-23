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
	"fmt"
	"sync"

	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/inbuiltpolicy"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// APIStore is a thread-safe store for APIs.
type APIStore struct {
	apis        map[string]*requestconfig.API
	mu          sync.RWMutex
	configStore *ConfigStore
	cfg         *config.Server
}

// NewAPIStore creates a new instance of APIStore.
func NewAPIStore(configStore *ConfigStore, cfg *config.Server) *APIStore {
	return &APIStore{
		configStore: configStore,
		// apis: make(map[string]*api.Api, 0),
		cfg: cfg,
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
			Name:                    api.Title,
			Version:                 api.Version,
			Vhost:                   api.Vhost,
			BasePath:                api.BasePath,
			APIType:                 api.ApiType,
			EnvType:                 api.EnvType,
			APILifeCycleState:       api.ApiLifeCycleState,
			AuthorizationHeader:     "", // You might want to set this field if applicable
			OrganizationID:          api.OrganizationId,
			UUID:                    api.Id,
			Tier:                    api.Tier,
			DisableAuthentication:   api.DisableAuthentications,
			DisableScopes:           api.DisableScopes,
			Resources:               make([]*requestconfig.Resource, 0),
			ResourceMap:             make(map[string]*requestconfig.Resource, 0),
			IsMockedAPI:             false, // You can add logic to determine if the API is mocked
			MutualSSL:               api.MutualSSL,
			TransportSecurity:       api.TransportSecurity,
			ApplicationSecurity:     api.ApplicationSecurity,
			BackendJwtConfiguration: convertBackendJWTTokenInfoToJWTConfig(api.BackendJWTTokenInfo, s.cfg, fmt.Sprintf("%s-%s", api.Title, api.Version)),
			SystemAPI:               api.SystemAPI,
			APIDefinition:           api.ApiDefinitionFile,
			APIDefinitionPath:       api.ApiDefinitionPath,
			Environment:             api.Environment,
			SubscriptionValidation:  api.SubscriptionValidation,
			// Endpoints:                         api.Endpoints,
			EndpointSecurity:                  convertSecurityInfoToEndpointSecurity(api.EndpointSecurity),
			AiProvider:                        convertAIProviderToDTO(api.Aiprovider),
			AIModelBasedRoundRobin:            convertAIModelBasedRoundRobinToDTO(api.AiModelBasedRoundRobin),
			DoSubscriptionAIRLInHeaderReponse: api.Aiprovider != nil && api.Aiprovider.PromptTokens != nil && api.Aiprovider.PromptTokens.In == dto.InHeader,
			DoSubscriptionAIRLInBodyReponse:   api.Aiprovider != nil && api.Aiprovider.PromptTokens != nil && api.Aiprovider.PromptTokens.In == dto.InBody,
			RequestInBuiltPolicies:            covertRequestInBuiltPoliciesToDTO(api.RequestInBuiltPolicies),
			ResponseInBuiltPolicies:           covertResponseInBuiltPoliciesToDTO(api.ResponseInBuiltPolicies),
		}
		for _, resource := range api.Resources {
			for _, operation := range resource.Methods {
				resource := buildResource(operation, resource.Path, resource.Endpoints, convertAIModelBasedRoundRobinToDTO(resource.AiModelBasedRoundRobin), func() []*requestconfig.EndpointSecurity {
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
				customAPI.Resources = append(customAPI.Resources, &resource)
				customAPI.ResourceMap[resource.GetResourceIdentifier()] = &resource
			}

		}
		s.cfg.Logger.Sugar().Debug(fmt.Sprintf("Adding API: %+v", customAPI.BackendJwtConfiguration))
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

// convertRequestInBuiltPoliciesToDTO converts a slice of InBuiltPolicy to a slice of dto.InBuiltPolicy.
func covertRequestInBuiltPoliciesToDTO(requestPolicies []*api.InBuiltPolicy) []dto.InBuiltPolicy {
	if requestPolicies == nil {
		return nil
	}
	dtoPolicies := make([]dto.InBuiltPolicy, 0, len(requestPolicies))
	for _, policy := range requestPolicies {
		switch policy.PolicyName {
		case inbuiltpolicy.RegexGuardrailName:
			dtoPolicies = append(dtoPolicies, &inbuiltpolicy.RegexGuardrail{
				BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
					PolicyName:    policy.PolicyName,
					PolicyID:      policy.PolicyID,
					PolicyVersion: policy.PolicyVersion,
					Parameters:    policy.Parameters,
				},
			})
		case inbuiltpolicy.WordCountGuardrailName:
			dtoPolicies = append(dtoPolicies, &inbuiltpolicy.WordCountGuardrail{
				BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
					PolicyName:    policy.PolicyName,
					PolicyID:      policy.PolicyID,
					PolicyVersion: policy.PolicyVersion,
					Parameters:    policy.Parameters,
				},
			})
		}
	}
	return dtoPolicies
}

// convertResponseInBuiltPoliciesToDTO converts a slice of InBuiltPolicy to a slice of dto.InBuiltPolicy.
func covertResponseInBuiltPoliciesToDTO(responsePolicies []*api.InBuiltPolicy) []dto.InBuiltPolicy {
	if responsePolicies == nil {
		return nil
	}
	dtoPolicies := make([]dto.InBuiltPolicy, 0, len(responsePolicies))
	for _, policy := range responsePolicies {
		switch policy.PolicyName {
		case "RegexGuardrail":
			dtoPolicies = append(dtoPolicies, &inbuiltpolicy.RegexGuardrail{
				BaseInBuiltPolicy: dto.BaseInBuiltPolicy{
					PolicyName:    policy.PolicyName,
					PolicyID:      policy.PolicyID,
					PolicyVersion: policy.PolicyVersion,
					Parameters:    policy.Parameters,
				},
			})
		}
	}
	return dtoPolicies
}

// convertAIModelBasedRoundRobinToDTO converts AIModelBasedRoundRobin to DTO.
func convertAIModelBasedRoundRobinToDTO(aiModelBasedRoundRobin *api.AIModelBasedRoundRobin) *dto.AIModelBasedRoundRobin {
	if aiModelBasedRoundRobin == nil {
		return nil
	}
	return &dto.AIModelBasedRoundRobin{
		Enabled:                      aiModelBasedRoundRobin.Enabled,
		OnQuotaExceedSuspendDuration: int(aiModelBasedRoundRobin.OnQuotaExceedSuspendDuration),
		ProductionModels:             convertModelWeights(aiModelBasedRoundRobin.ProductionModels),
		SandboxModels:                convertModelWeights(aiModelBasedRoundRobin.SandboxModels),
	}
}

// convertModelWeights converts []*api.ModelWeight to []dto.ModelWeight.
func convertModelWeights(apiModelWeights []*api.ModelWeight) []dto.ModelWeight {
	dtoModelWeights := make([]dto.ModelWeight, len(apiModelWeights))
	for i, modelWeight := range apiModelWeights {
		dtoModelWeights[i] = dto.ModelWeight{
			Model:    modelWeight.Model,
			Endpoint: modelWeight.Endpoint,
			Weight:   int(modelWeight.Weight),
		}
	}
	return dtoModelWeights
}

// convertAIProviderToDTO converts AIProvider to DTO.
func convertAIProviderToDTO(aiProvider *api.AIProvider) *dto.AIProvider {
	if aiProvider == nil {
		return nil
	}
	return &dto.AIProvider{
		ProviderName:       aiProvider.ProviderName,
		ProviderAPIVersion: aiProvider.ProviderAPIVersion,
		Organization:       aiProvider.Organization,
		Enabled:            aiProvider.Enabled,
		SupportedModels:    aiProvider.SupportedModels,
		RequestModel:       convertValueDetailsPtr(aiProvider.RequestModel),
		ResponseModel:      convertValueDetailsPtr(aiProvider.ResponseModel),
		PromptTokens:       convertValueDetailsPtr(aiProvider.PromptTokens),
		CompletionToken:    convertValueDetailsPtr(aiProvider.CompletionToken),
		TotalToken:         convertValueDetailsPtr(aiProvider.TotalToken),
	}
}

// convertValueDetailsPtr converts *api.ValueDetails to *dto.ValueDetails.
func convertValueDetailsPtr(valueDetails *api.ValueDetails) *dto.ValueDetails {
	if valueDetails == nil {
		return nil
	}
	return &dto.ValueDetails{
		In:    valueDetails.In,
		Value: valueDetails.Value,
	}
}

// GetMatchedAPI retrieves the API that matches the given API key.
// GetMatchedAPI retrieves the API that matches the given API key.
func (s *APIStore) GetMatchedAPI(apiKey string) *requestconfig.API {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apis[apiKey]
}

// UpdateMatchedAPI updates the API that matches the given API key.
func (s *APIStore) UpdateMatchedAPI(apiKey string, api *requestconfig.API) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apis[apiKey] = api
}

// convertSecurityInfoToEndpointSecurity converts SecurityInfo to EndpointSecurity.
func convertSecurityInfoToEndpointSecurity(securityInfo []*api.SecurityInfo) []requestconfig.EndpointSecurity {
	if securityInfo == nil {
		return nil
	}
	endpointSecurities := []requestconfig.EndpointSecurity{}
	for i := range securityInfo {
		security := (securityInfo)[i]
		endpointSecurity := requestconfig.EndpointSecurity{
			Password:         security.Password,
			Enabled:          security.Enabled,
			Username:         security.Username,
			SecurityType:     security.SecurityType,
			CustomParameters: security.CustomParameters,
		}
		endpointSecurities = append(endpointSecurities, endpointSecurity)
	}
	return endpointSecurities
}

// ConvertBackendJWTTokenInfoToJWTConfig converts BackendJWTTokenInfo to JWTConfiguration.
func convertBackendJWTTokenInfoToJWTConfig(info *api.BackendJWTTokenInfo, cfg *config.Server, apiName string) *dto.BackendJWTConfiguration {
	if info == nil {
		return nil
	}

	// Convert CustomClaims from map[string]*Claim to map[string]ClaimValue
	customClaims := make(map[string]*dto.ClaimValue)
	for key, claim := range info.CustomClaims {
		if claim != nil {
			customClaims[key] = &dto.ClaimValue{
				Value: claim.Value,
				Type:  claim.Type,
			}
		}
	}
	publicCert, err := util.LoadCertificate(cfg.JWTGeneratorPublicKeyPath)
	if err != nil {
		cfg.Logger.Error(err, fmt.Sprintf("Error loading public cert. Marking API %s as backend jwt disabled.", apiName))
		info.Enabled = false
	}
	privateKey, err := util.LoadPrivateKey(cfg.JWTGeneratorPrivateKeyPath)
	if err != nil {
		cfg.Logger.Error(err, fmt.Sprintf("Error loading private key. Marking API %s as backend jwt disabled. Path: %s", apiName, cfg.JWTGeneratorPrivateKeyPath))
		info.Enabled = false
	}
	return &dto.BackendJWTConfiguration{
		Enabled:            info.Enabled,
		JWTHeader:          info.Header,
		ConsumerDialectURI: "", // Add a default value or fetch if needed
		SignatureAlgorithm: info.SigningAlgorithm,
		Encoding:           info.Encoding,
		TokenIssuerDtoMap:  make(map[string]dto.TokenIssuer), // Populate if required
		JwtExcludedClaims:  make(map[string]bool),            // Populate if required
		PublicCert:         publicCert,                       // Add conversion logic if needed
		PrivateKey:         privateKey,                       // Add conversion logic if needed
		TTL:                int64(info.TokenTTL),             // Convert int32 to int64
		CustomClaims:       customClaims,
	}
}
