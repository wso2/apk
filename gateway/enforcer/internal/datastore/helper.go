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
	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	auth "github.com/wso2/apk/gateway/enforcer/internal/authentication/authconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

func buildResource(operation *api.Operation, path string, endpointSecurity []*requestconfig.EndpointSecurity) requestconfig.Resource {
	authConfig := auth.AuthenticationConfig{
		Disabled: operation.ApiAuthentication.Disabled,
	}
	if operation.ApiAuthentication != nil {
		if operation.ApiAuthentication.GetOauth2() != nil {
			authConfig.Oauth2AuthenticationConfig = &auth.Oauth2AuthenticationConfig{
				Header:              operation.ApiAuthentication.GetOauth2().GetHeader(),
				SendTokenToUpstream: operation.ApiAuthentication.GetOauth2().GetSendTokenToUpstream(),
			}
		}
		if operation.ApiAuthentication.GetJwt() != nil {
			authConfig.JWTAuthenticationConfig = &auth.JWTAuthenticationConfig{
				Header:              operation.ApiAuthentication.GetJwt().GetHeader(),
				SendTokenToUpstream: operation.ApiAuthentication.GetJwt().GetSendTokenToUpstream(),
				Audience:            operation.ApiAuthentication.GetJwt().GetAudience(),
			}
		}
		apiKeyAuthConfigs := make([]*auth.APIKeyAuthenticationConfig, len(operation.ApiAuthentication.Apikey))
		for i, apiKey := range operation.ApiAuthentication.Apikey {
			apiKeyAuthConfigs[i] = &auth.APIKeyAuthenticationConfig{
				In:                  apiKey.GetIn(),
				Name:                apiKey.GetName(),
				SendTokenToUpstream: apiKey.GetSendTokenToUpstream(),
			}
		}
		authConfig.APIKeyAuthenticationConfigs = apiKeyAuthConfigs
	}
	return requestconfig.Resource{
		MatchID:              operation.MatchID,
		Path:                 util.NormalizePath(path),
		Method:               requestconfig.HTTPMethods(operation.Method),
		Tier:                 operation.Tier,
		EndpointSecurity:     endpointSecurity,
		AuthenticationConfig: authConfig,
		Scopes:               operation.Scopes,
	}
}

func buildPolicy(policies *api.OperationPolicies) requestconfig.PolicyConfig {
	return requestconfig.PolicyConfig{
		Request:  policies.Request,
		Response: policies.Response,
		Fault:    policies.Fault,
	}
}
