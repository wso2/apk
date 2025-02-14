/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package oasparser

import (
	"strconv"
	"strings"
	"time"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	jwt "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	hcmv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	"github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	envoy "github.com/wso2/apk/adapter/internal/oasparser/envoyconf"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	"github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	"github.com/wso2/apk/common-go-libs/apis/dp/v1alpha1"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// GetGlobalClusters generates initial internal clusters for given environment.
func GetGlobalClusters() ([]*clusterv3.Cluster, []*corev3.Address) {
	var (
		clusters  []*clusterv3.Cluster
		endpoints []*corev3.Address
	)
	conf := config.ReadConfigs()

	if conf.Envoy.RateLimit.Enabled {
		rlCluster, rlEP, errRL := envoy.CreateRateLimitCluster()
		if errRL == nil {
			clusters = append(clusters, rlCluster)
			endpoints = append(endpoints, rlEP...)
		} else {
			logger.LoggerOasparser.ErrorC(logging.PrintError(logging.Error2248, logging.BLOCKER, "Failed to initialize ratelimit cluster. Hence terminating the adapter. Error: %s", errRL))
		}
	}

	if conf.Tracing.Enabled && conf.Tracing.Type != envoy.TracerTypeAzure {
		logger.LoggerOasparser.Debugln("Creating global cluster - Tracing")
		if c, e, err := envoy.CreateTracingCluster(conf); err == nil {
			clusters = append(clusters, c...)
			endpoints = append(endpoints, e...)
		} else {
			logger.LoggerOasparser.ErrorC(logging.PrintError(logging.Error2249, logging.CRITICAL, "Failed to initialize tracer's cluster. Router tracing will be disabled. %v", err.Error()))
			conf.Tracing.Enabled = false
		}
	}

	return clusters, endpoints
}

// GetProductionListener generates the listener configurations.
//
// The VirtualHost is named as "default".
// The provided set of envoy routes will be assigned under the virtual host
//
// The RouteConfiguration is named as "default"
func GetProductionListener(gateway *gwapiv1.Gateway, resolvedListenerCerts map[string]map[string][]byte, gwLuaScript string) []*listenerv3.Listener {
	listeners := envoy.CreateListenerByGateway(gateway, resolvedListenerCerts, gwLuaScript)
	return listeners
}

// GetRouteConfigs generates routesconfiguration configurations.
//
// The VirtualHost is named as "default".
// The provided set of envoy routes will be assigned under the virtual host
//
// The RouteConfiguration is named as "default"
func GetRouteConfigs(vhostToRouteArrayMap map[string][]*routev3.Route, routeConfigName string,
	customRateLimitPolicies []*model.CustomRateLimitPolicy, vhostToSubscriptionAIRL map[string]bool, vhostToSubscriptionRL map[string]bool) *routev3.RouteConfiguration {
	vHosts := envoy.CreateVirtualHosts(vhostToRouteArrayMap, customRateLimitPolicies, vhostToSubscriptionAIRL, vhostToSubscriptionRL)
	routeConfig := envoy.CreateRoutesConfigForRds(vHosts, routeConfigName)
	return routeConfig
}

// GetCacheResources converts the envoy endpoints, clusters, routes, and listener to
// the resource type which is the format required for the Xds cache.
//
// The returned resources are listeners, clusters, routeConfigurations, endpoints
func GetCacheResources(endpoints []*corev3.Address, clusters []*clusterv3.Cluster,
	listeners []*listenerv3.Listener, routeConfigs map[string]*routev3.RouteConfiguration) (
	listenerRes []types.Resource, clusterRes []types.Resource, routeConfigRes []types.Resource,
	endpointRes []types.Resource) {

	clusterRes = []types.Resource{}
	endpointRes = []types.Resource{}
	for _, cluster := range clusters {
		clusterRes = append(clusterRes, cluster)
	}
	for _, endpoint := range endpoints {
		endpointRes = append(endpointRes, endpoint)
	}
	listenerRes = []types.Resource{}
	for _, listener := range listeners {
		listenerRes = append(listenerRes, listener)
	}
	routeConfigRes = []types.Resource{}
	for _, routeConfig := range routeConfigs {
		routeConfigRes = append(routeConfigRes, routeConfig)
	}
	return listenerRes, clusterRes, routeConfigRes, endpointRes
}

// GetEnforcerAPI retrieves the ApiDS object model for a given swagger definition
// along with the vhost to deploy the API.
func GetEnforcerAPI(adapterInternalAPI *model.AdapterInternalAPI, vhost string) *api.Api {
	resources := []*api.Resource{}
	isMockedAPI := false
	clientCertificates := []*api.Certificate{}

	for _, res := range adapterInternalAPI.GetResources() {
		var operations = make([]*api.Operation, len(res.GetMethod()))
		for i, op := range res.GetMethod() {
			operations[i] = GetEnforcerAPIOperation(*op, isMockedAPI)
		}
		resource := &api.Resource{
			Id:      res.GetID(),
			Methods: operations,
			Path:    res.GetPath(),
		}
		if res.GetEndpoints() != nil {
			resource.Endpoints = generateRPCEndpointCluster(res.GetEndpoints())
		}
		if res.GetEndpointSecurity() != nil {
			resource.EndpointSecurity = generateRPCEndpointSecurity(res.GetEndpointSecurity())
		}
		modelBasedRoundRobin := res.GetAIModelBasedRoundRobin()
		logger.LoggerOasparser.Debugf("Get AI Model Based Round Robin: %+v", modelBasedRoundRobin)
		if modelBasedRoundRobin != nil && (modelBasedRoundRobin.ProductionModels != nil || modelBasedRoundRobin.SandboxModels != nil) {
			resource.AiModelBasedRoundRobin = &api.AIModelBasedRoundRobin{
				Enabled:                      true,
				OnQuotaExceedSuspendDuration: int32(modelBasedRoundRobin.OnQuotaExceedSuspendDuration),
				ProductionModels:             convertModelWeights(modelBasedRoundRobin.ProductionModels),
				SandboxModels:                convertModelWeights(modelBasedRoundRobin.SandboxModels),
			}
		}
		logger.LoggerOasparser.Infof("Resource AI Model Based Round Robin: %+v", resource.AiModelBasedRoundRobin)
		resources = append(resources, resource)
	}

	for _, cert := range adapterInternalAPI.GetClientCerts() {
		certificate := &api.Certificate{
			Alias:   cert.Alias,
			Content: cert.Content,
		}
		clientCertificates = append(clientCertificates, certificate)
	}

	var backendJWTTokenInfo *api.BackendJWTTokenInfo

	backendJWTTokenInfoFromInternalAPI := adapterInternalAPI.GetBackendJWTTokenInfo()

	if backendJWTTokenInfoFromInternalAPI != nil {
		customClaims := make(map[string]*api.Claim)
		for _, claim := range backendJWTTokenInfoFromInternalAPI.CustomClaims {
			customClaims[claim.Claim] = &api.Claim{
				Value: claim.Value.Value,
				Type:  claim.Value.Type,
			}
		}
		backendJWTTokenInfo = &api.BackendJWTTokenInfo{
			Enabled:          backendJWTTokenInfoFromInternalAPI.Enabled,
			Encoding:         backendJWTTokenInfoFromInternalAPI.Encoding,
			Header:           backendJWTTokenInfoFromInternalAPI.Header,
			SigningAlgorithm: backendJWTTokenInfoFromInternalAPI.SigningAlgorithm,
			TokenTTL:         int32(backendJWTTokenInfoFromInternalAPI.TokenTTL),
			CustomClaims:     customClaims,
		}
	}

	var aiProvider *api.AIProvider

	aiProviderFromInternalAPI := adapterInternalAPI.GetAIProvider()
	logger.LoggerOasparser.Debugf("Before Internal AI Provider: %+v", aiProviderFromInternalAPI)

	if aiProviderFromInternalAPI.Enabled {
		aiProvider = &api.AIProvider{
			Enabled:            aiProviderFromInternalAPI.Enabled,
			ProviderName:       aiProviderFromInternalAPI.ProviderName,
			ProviderAPIVersion: aiProviderFromInternalAPI.ProviderAPIVersion,
			Organization:       aiProviderFromInternalAPI.Organization,
			SupportedModels:    aiProviderFromInternalAPI.SupportedModels,
			RequestModel: &api.ValueDetails{
				In:    aiProviderFromInternalAPI.RequestModel.In,
				Value: aiProviderFromInternalAPI.RequestModel.Value,
			},
			ResponseModel: &api.ValueDetails{
				In:    aiProviderFromInternalAPI.ResponseModel.In,
				Value: aiProviderFromInternalAPI.ResponseModel.Value,
			},
			PromptTokens: &api.ValueDetails{
				In:    aiProviderFromInternalAPI.PromptTokens.In,
				Value: aiProviderFromInternalAPI.PromptTokens.Value,
			},
			CompletionToken: &api.ValueDetails{
				In:    aiProviderFromInternalAPI.CompletionToken.In,
				Value: aiProviderFromInternalAPI.CompletionToken.Value,
			},
			TotalToken: &api.ValueDetails{
				In:    aiProviderFromInternalAPI.TotalToken.In,
				Value: aiProviderFromInternalAPI.TotalToken.Value,
			},
		}
	}

	logger.LoggerOasparser.Debugf("After Conversion AI Provider: %+v", aiProvider)

	var modelBasedRoundRobin *api.AIModelBasedRoundRobin

	modelBasedRoundRobinFromInternalAPI := adapterInternalAPI.GetModelBasedRoundRobin()
	logger.LoggerOasparser.Debugf("Before Internal Model Based Round Robin: %+v", modelBasedRoundRobinFromInternalAPI)

	if modelBasedRoundRobinFromInternalAPI.ProductionModels != nil || modelBasedRoundRobinFromInternalAPI.SandboxModels != nil {
		modelBasedRoundRobin = &api.AIModelBasedRoundRobin{
			Enabled:                      true,
			OnQuotaExceedSuspendDuration: int32(modelBasedRoundRobinFromInternalAPI.OnQuotaExceedSuspendDuration),
			ProductionModels:             convertModelWeights(modelBasedRoundRobinFromInternalAPI.ProductionModels),
			SandboxModels:                convertModelWeights(modelBasedRoundRobinFromInternalAPI.SandboxModels),
		}
	}

	logger.LoggerOasparser.Debugf("After Conversion Model Based Round Robin: %+v", modelBasedRoundRobin)

	return &api.Api{
		Id:                     adapterInternalAPI.UUID,
		Title:                  adapterInternalAPI.GetTitle(),
		BasePath:               adapterInternalAPI.GetXWso2Basepath(),
		Version:                adapterInternalAPI.GetVersion(),
		ApiType:                adapterInternalAPI.GetAPIType(),
		Resources:              resources,
		ApiLifeCycleState:      adapterInternalAPI.LifecycleStatus,
		Tier:                   adapterInternalAPI.GetXWso2ThrottlingTier(),
		DisableAuthentications: adapterInternalAPI.GetDisableAuthentications(),
		DisableScopes:          adapterInternalAPI.GetDisableScopes(),
		OrganizationId:         adapterInternalAPI.OrganizationID,
		Vhost:                  vhost,
		EnvType:                adapterInternalAPI.EnvType,
		BackendJWTTokenInfo:    backendJWTTokenInfo,
		Endpoints:              generateRPCEndpointCluster(adapterInternalAPI.Endpoints),
		EndpointSecurity:       generateRPCEndpointSecurity(adapterInternalAPI.EndpointSecurity),
		// IsMockedApi:         isMockedAPI,
		ClientCertificates:  clientCertificates,
		MutualSSL:           adapterInternalAPI.GetMutualSSL(),
		ApplicationSecurity: adapterInternalAPI.GetApplicationSecurity(),
		TransportSecurity:   !adapterInternalAPI.GetDisableMtls(),
		// GraphQLSchema:         adapterInternalAPI.GraphQLSchema,
		// GraphqlComplexityInfo: adapterInternalAPI.GraphQLComplexities.Data.List,
		SystemAPI:              adapterInternalAPI.IsSystemAPI,
		ApiDefinitionFile:      adapterInternalAPI.GetAPIDefinitionFile(),
		ApiDefinitionPath:      adapterInternalAPI.GetAPIDefinitionEndpoint(),
		Environment:            adapterInternalAPI.GetEnvironment(),
		SubscriptionValidation: adapterInternalAPI.GetSubscriptionValidation(),
		Aiprovider: func() *api.AIProvider {
			if aiProvider != nil && aiProvider.Enabled {
				return aiProvider
			}
			return nil
		}(),
		HttpRouteIDs: adapterInternalAPI.HTTPRouteIDs,
		AiModelBasedRoundRobin: func() *api.AIModelBasedRoundRobin {
			if modelBasedRoundRobin != nil && modelBasedRoundRobin.Enabled {
				return modelBasedRoundRobin
			}
			return nil
		}(),
	}
}

// GetEnforcerAPIOperation builds the operation object expected by the proto definition
func GetEnforcerAPIOperation(operation model.Operation, isMockedAPI bool) *api.Operation {

	// var mockedAPIConfig *api.MockedApiConfig
	// if isMockedAPI {
	// 	mockedAPIConfig = operation.GetMockedAPIConfig()
	// }

	policies := &api.OperationPolicies{
		Request:  castPoliciesToEnforcerPolicies(operation.GetPolicies().Request),
		Response: castPoliciesToEnforcerPolicies(operation.GetPolicies().Response),
		Fault:    castPoliciesToEnforcerPolicies(operation.GetPolicies().Fault),
	}
	apiOperation := api.Operation{
		Method:            operation.GetMethod(),
		Scopes:            operation.GetScopes(),
		ApiAuthentication: castAPIAuthenticationsToEnforcerAPIAuthentications(operation.GetAuthentication()),
		Tier:              operation.GetTier(),
		Policies:          policies,
		MatchID:           operation.GetMatchID(),
		// MockedApiConfig: mockedAPIConfig,
	}
	return &apiOperation
}

// castAPIAuthenticationsToEnforcerAPIAuthentications casts the APIAuthentications to EnforcerAPIAuthentications
func castAPIAuthenticationsToEnforcerAPIAuthentications(authentication *model.Authentication) *api.APIAuthentication {
	enforcerAuthentication := &api.APIAuthentication{}
	enforcerAuthentication.Disabled = authentication.Disabled
	if authentication.JWT != nil {
		enforcerAuthentication.Jwt = &api.JWT{
			Header:              strings.ToLower(authentication.JWT.Header),
			SendTokenToUpstream: authentication.JWT.SendTokenToUpstream,
			Audience:            authentication.JWT.Audience,
		}
	}
	var apiKeys []*api.APIKey
	for _, apiKey := range authentication.APIKey {
		name := apiKey.Name
		if apiKey.In == "Header" {
			name = strings.ToLower(name)
		}
		apiKeys = append(apiKeys, &api.APIKey{
			In:                  apiKey.In,
			Name:                name,
			SendTokenToUpstream: apiKey.SendTokenToUpstream,
		})
	}
	if len(apiKeys) > 0 {
		enforcerAuthentication.Apikey = apiKeys
	}
	if authentication.Oauth2 != nil {
		enforcerAuthentication.Oauth2 = &api.Oauth2{
			Header:              strings.ToLower(authentication.Oauth2.Header),
			SendTokenToUpstream: authentication.Oauth2.SendTokenToUpstream,
		}
	}
	return enforcerAuthentication
}

func castPoliciesToEnforcerPolicies(policies []model.Policy) []*api.Policy {
	enforcerPolicies := make([]*api.Policy, 0, len(policies))
	for _, policy := range policies {
		if !policy.IsPassToEnforcer {
			// The API Policy do not want support from enforcer to handle the request
			continue
		}
		parameterMap := make(map[string]string)
		if policy.Parameters != nil {
			if params, ok := policy.Parameters.(map[string]interface{}); ok {
				for paramK := range params {
					if paramV, parsed := params[paramK].(string); parsed {
						parameterMap[paramK] = paramV
					} else if paramV, parsed := params[paramK].(bool); parsed {
						parameterMap[paramK] = strconv.FormatBool(paramV)
					}
				}

			}
		}
		enforcerPolicies = append(enforcerPolicies, &api.Policy{
			Action:     policy.Action,
			Parameters: parameterMap,
		})
	}
	return enforcerPolicies
}

func generateRPCEndpointCluster(inputEndpointCluster *model.EndpointCluster) *api.EndpointCluster {
	if inputEndpointCluster == nil || len(inputEndpointCluster.Endpoints) == 0 {
		return nil
	}
	urls := []*api.Endpoint{}
	for _, ep := range inputEndpointCluster.Endpoints {
		endpoint := &api.Endpoint{
			Basepath: ep.Basepath,
			Host:     ep.Host,
			Port:     ep.Port,
			URLType:  ep.URLType,
		}
		urls = append(urls, endpoint)
	}

	endpoints := &api.EndpointCluster{
		Urls: urls,
	}
	if inputEndpointCluster.Config != nil {
		// retry config
		var retryConfig *api.RetryConfig
		if inputEndpointCluster.Config.RetryConfig != nil {
			inputRetryConfig := inputEndpointCluster.Config.RetryConfig
			retryConfig = &api.RetryConfig{
				Count:       uint32(inputRetryConfig.Count),
				StatusCodes: inputRetryConfig.StatusCodes,
			}
		}
		// timeout config
		var timeoutConfig *api.TimeoutConfig
		if inputEndpointCluster.Config.TimeoutInMillis != 0 { // if zero, means not set. Then, global timeout is applied via route configs.
			timeoutConfig = &api.TimeoutConfig{
				RouteTimeoutInMillis: uint32(inputEndpointCluster.Config.TimeoutInMillis),
			}
		}
		// Set all endpoint configs
		endpoints.Config = &api.EndpointClusterConfig{
			RetryConfig:   retryConfig,
			TimeoutConfig: timeoutConfig,
		}
	}
	return endpoints
}

func convertModelWeights(inputModels []model.InternalModelWeight) []*api.ModelWeight {
	var outputModels []*api.ModelWeight
	for _, model := range inputModels {
		outputModels = append(outputModels, &api.ModelWeight{
			Model:    model.Model,
			Endpoint: model.EndpointClusterName,
			Weight:   int32(model.Weight),
		})
	}
	return outputModels
}

func generateRPCEndpointSecurity(inputEndpointSecurity []*model.EndpointSecurity) []*api.SecurityInfo {
	if inputEndpointSecurity == nil {
		return nil
	}
	var securityConfig []*api.SecurityInfo
	for _, security := range inputEndpointSecurity {
		securityConfig = append(securityConfig, &api.SecurityInfo{
			SecurityType:     security.Type,
			Username:         security.Username,
			Password:         security.Password,
			Enabled:          security.Enabled,
			CustomParameters: security.CustomParameters,
		})
	}
	return securityConfig
}

// GetJWTRequirements returns the jwt requirements for the resource

func GetJWTRequirements(adapterAPI *model.AdapterInternalAPI, jwtIssuers map[string]*v1alpha1.ResolvedJWTIssuer) *jwt.JwtRequirement {
	var selectedIssuers []string
	for issuserName, jwtIssuer := range jwtIssuers {
		if contains(jwtIssuer.Environments, "*") {
			selectedIssuers = append(selectedIssuers, issuserName)
		} else if contains(jwtIssuer.Environments, adapterAPI.GetEnvironment()) {
			selectedIssuers = append(selectedIssuers, issuserName)
		}
	}
	if len(selectedIssuers) >= 1 {
		return &jwt.JwtRequirement{
			RequiresType: &jwt.JwtRequirement_RequiresAny{
				RequiresAny: &jwt.JwtRequirementOrList{
					Requirements: func() []*jwt.JwtRequirement {
						var requirements []*jwt.JwtRequirement
						for _, issuer := range selectedIssuers {
							requirements = append(requirements, &jwt.JwtRequirement{
								RequiresType: &jwt.JwtRequirement_ProviderName{
									ProviderName: issuer,
								},
							})
						}
						requirements = append(requirements, &jwt.JwtRequirement{
							RequiresType: &jwt.JwtRequirement_AllowMissingOrFailed{},
						})
						return requirements
					}(),
				},
			}}
	}
	return nil
}

// GenerateJWTPRoviderv3 generates the jwt provider for the resource
func GenerateJWTPRoviders(jwtProviderMap map[string]map[string]*v1alpha1.ResolvedJWTIssuer) (map[string]*jwt.JwtProvider, []*clusterv3.Cluster, []*corev3.Address, error) {
	jwtProviders := map[string]*jwt.JwtProvider{}
	var clusters []*clusterv3.Cluster
	var addresses []*corev3.Address

	for _, orgwizeJWTProviders := range jwtProviderMap {
		for issuerMappingName, jwtIssuer := range orgwizeJWTProviders {
			provider, cluster, address, err := getjwtAuthFilters(jwtIssuer, issuerMappingName)
			if err != nil {
				return nil, nil, nil, err
			}
			jwtProviders[issuerMappingName] = provider
			clusters = append(clusters, cluster...)
			addresses = append(addresses, address...)
		}
	}
	return jwtProviders, clusters, addresses, nil
}

// Function to check if a string array contains a specific string
func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
func getjwtAuthFilters(tokenIssuer *v1alpha1.ResolvedJWTIssuer, issuerName string) (*jwt.JwtProvider, []*clusterv3.Cluster, []*corev3.Address, error) {
	jwksClusters := make([]*clusterv3.Cluster, 0)
	jwksAddresses := make([]*corev3.Address, 0)
	jwtProvider := &jwt.JwtProvider{
		Issuer:                 tokenIssuer.Issuer,
		Forward:                true,
		FailedStatusInMetadata: "failed_status",
		PayloadInMetadata:      "payload_in_metadata",
	}
	if tokenIssuer.SignatureValidation.JWKS != nil {
		logger.LoggerOasparser.Infof("JWKS URL: %s", tokenIssuer.SignatureValidation.JWKS.URL)
		jwksCluster, jwksAddress, err := getRemoteJWKSCluster(*tokenIssuer.SignatureValidation.JWKS, issuerName)
		if err != nil {
			logger.LoggerOasparser.Error(err)
			return nil, nil, nil, err
		}
		jwksClusters = append(jwksClusters, jwksCluster...)
		jwksAddresses = append(jwksAddresses, jwksAddress...)
		jwtProvider.JwksSourceSpecifier = &jwt.JwtProvider_RemoteJwks{
			RemoteJwks: &jwt.RemoteJwks{
				HttpUri: &corev3.HttpUri{
					Uri: tokenIssuer.SignatureValidation.JWKS.URL,
					HttpUpstreamType: &corev3.HttpUri_Cluster{
						Cluster: issuerName,
					},
					Timeout: durationpb.New(2 * time.Second),
				},
				CacheDuration: durationpb.New(2 * time.Hour),
			},
		}
	} else if tokenIssuer.SignatureValidation.Certificate != nil {
		logger.LoggerOasparser.Infof("ResolvedCertificate: %s", tokenIssuer.SignatureValidation.Certificate.ResolvedCertificate)
		jwtProvider.JwksSourceSpecifier = &jwt.JwtProvider_LocalJwks{
			LocalJwks: &corev3.DataSource{
				Specifier: &corev3.DataSource_InlineString{InlineString: tokenIssuer.SignatureValidation.Certificate.ResolvedCertificate},
			},
		}
	}

	return jwtProvider, jwksClusters, jwksAddresses, nil
}
func getRemoteJWKSCluster(jwksInfo v1alpha1.ResolvedJWKS, clusterName string) ([]*clusterv3.Cluster, []*corev3.Address, error) {
	endpoint, err := model.GETHTTPEndpoint(jwksInfo.URL)
	if err != nil {
		return nil, nil, err
	}
	if jwksInfo.TLS != nil {
		endpoint.AllowedSANs = jwksInfo.TLS.AllowedSANs
		endpoint.Certificate = []byte(jwksInfo.TLS.ResolvedCertificate)
	}
	endpoints := make([]model.Endpoint, 0)
	endpoints = append(endpoints, *endpoint)
	endpointCluster := &model.EndpointCluster{
		Endpoints: endpoints,
	}
	return envoy.ProcessEndpoints(clusterName, endpointCluster, 20, endpoint.Basepath)
}

// GetJWTFilter
func GetJWTFilter(jwtRequirement map[string]*jwt.JwtRequirement, jwtProviders map[string]*jwt.JwtProvider) (*hcmv3.HttpFilter, error) {
	if len(jwtProviders) == 0 {
		return nil, nil
	}
	jwtAuthentication := &jwt.JwtAuthentication{
		Providers:      jwtProviders,
		RequirementMap: jwtRequirement,
	}
	// Assuming jwtAuthentication is already defined and initialized
	typedConfig, err := anypb.New(jwtAuthentication)
	if err != nil {
		logger.LoggerOasparser.ErrorC(logging.PrintError(logging.Error2250, logging.CRITICAL, "Failed to parse JWTAuthentication %v", err.Error()))
	}
	return &hcmv3.HttpFilter{
		Name: envoyconf.EnvoyJWT,
		ConfigType: &hcmv3.HttpFilter_TypedConfig{
			TypedConfig: typedConfig,
		},
	}, nil
}
