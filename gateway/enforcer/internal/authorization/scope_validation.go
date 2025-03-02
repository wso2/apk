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

package authorization

import (
	"encoding/json"

	"github.com/wso2/apk/gateway/enforcer/internal/authentication/authenticator"
	"github.com/wso2/apk/gateway/enforcer/internal/config"
	"github.com/wso2/apk/gateway/enforcer/internal/datastore"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
)

// ValidateScopes validates the scopes of the user against the required scopes.
func ValidateScopes(rch *requestconfig.Holder, subAppDataStore *datastore.SubscriptionApplicationDataStore, cfg *config.Server) *dto.ImmediateResponse {
	requiredScopes := rch.MatchedResource.Scopes
	if rch.AuthenticatedAuthenticationType == authenticator.Oauth2AuthType {
		scopes := rch.JWTValidationInfo.Scopes
		if len(requiredScopes) == 0 {
			return nil
		}
		if len(scopes) == 0 {
			scopeValidationErrorMessage := dto.ErrorResponse{Code: 900910, ErrorMessage: "The access token does not allow you to access the requested resource", ErrorDescription: "User is NOT authorized to access the Resource: " + rch.MatchedResource.Path + ". Scope validation failed."}
			forbiddenJSONMessage, _ := json.MarshalIndent(scopeValidationErrorMessage, "", "  ")

			return &dto.ImmediateResponse{
				StatusCode: 403,
				Message:    string(forbiddenJSONMessage),
			}
		}
		found := false
		for _, requiredScope := range requiredScopes {
			for _, scope := range scopes {
				if requiredScope == scope {
					found = true
					break
				}
			}
		}
		if !found {
			scopeValidationErrorMessage := dto.ErrorResponse{Code: 900910, ErrorMessage: "The access token does not allow you to access the requested resource", ErrorDescription: "User is NOT authorized to access the Resource: " + rch.MatchedResource.Path + ". Scope validation failed."}
			forbiddenJSONMessage, _ := json.MarshalIndent(scopeValidationErrorMessage, "", "  ")

			return &dto.ImmediateResponse{
				StatusCode: 403,
				Message:    string(forbiddenJSONMessage),
			}
		}
	}
	return nil
}
