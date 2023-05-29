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

package model

import (
	"errors"
	"fmt"

	"github.com/wso2/apk/adapter/internal/oasparser/constants"
)

// supportedPoliciesMap maps (policy action name) -> (policy layout)
var supportedPoliciesMap = map[string]policyLayout{
	constants.ActionHeaderAdd: {
		RequiredParams:   []string{constants.HeaderName, constants.HeaderValue},
		IsPassToEnforcer: false,
	},
	constants.ActionHeaderRemove: {
		RequiredParams:   []string{constants.HeaderName},
		IsPassToEnforcer: false,
	},
	"ADD_QUERY": {
		RequiredParams:   []string{"queryParamName", "queryParamValue"},
		IsPassToEnforcer: true,
	},
	constants.ActionInterceptorService: {
		RequiredParams:   []string{constants.InterceptorServiceURL, constants.InterceptorServiceIncludes},
		IsPassToEnforcer: false,
	},
	constants.ActionRewriteMethod: {
		RequiredParams:   []string{constants.UpdatedMethod},
		IsPassToEnforcer: true,
	},
	constants.ActionRewritePath: {
		RequiredParams:   []string{constants.RewritePathResourcePath, constants.IncludeQueryParams},
		IsPassToEnforcer: true,
	},
	constants.ActionOPA: {
		RequiredParams:   []string{constants.ServerURL, constants.Policy},
		IsPassToEnforcer: true,
	},
}

// PolicyLayout holds the layout of policy that support by APK
type policyLayout struct {
	RequiredParams   []string
	IsPassToEnforcer bool
}

// validatePolicyAction validates policy against the policy definition that supported by APK
func validatePolicyAction(policy *Policy) error {
	if layout, ok := supportedPoliciesMap[policy.Action]; ok {
		for _, requiredParam := range layout.RequiredParams {
			if params, isMap := policy.Parameters.(map[string]interface{}); isMap {
				if _, ok := params[requiredParam]; !ok {
					return fmt.Errorf("required parameter %q not found for the policy action %q", requiredParam, policy.Action)
				}
			} else {
				return errors.New("policy params required in map format")
			}
		}
		policy.IsPassToEnforcer = layout.IsPassToEnforcer
	} else {
		return fmt.Errorf("policy action %q not supported by APK gateway", policy.Action)
	}
	return nil
}
