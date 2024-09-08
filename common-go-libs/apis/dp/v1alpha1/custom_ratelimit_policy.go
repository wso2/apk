/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package v1alpha1

import (
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
)

// CustomRateLimitPolicyDef defines the desired state of CustomPolicy
type CustomRateLimitPolicyDef struct {
	Key             string `json:"key,omitempty"`
	Value           string `json:"value,omitempty"`
	RequestsPerUnit uint32 `json:"requestsPerUnit,omitempty"`

	Unit string `json:"unit,omitempty"`
	// RateLimit    RateLimit `json:"rateLimit,omitempty"`
	Organization string `json:"organization,omitempty"`
}

// ParseCustomRateLimitPolicy parses the custom rate limit policy
func ParseCustomRateLimitPolicy(customRateLimitCR dpv1alpha3.RateLimitPolicy) *CustomRateLimitPolicyDef {
	return &CustomRateLimitPolicyDef{
		Key:             customRateLimitCR.Spec.Override.Custom.Key,
		Value:           customRateLimitCR.Spec.Override.Custom.Value,
		RequestsPerUnit: customRateLimitCR.Spec.Override.Custom.RequestsPerUnit,
		Unit:            customRateLimitCR.Spec.Override.Custom.Unit,
		Organization:    customRateLimitCR.Spec.Override.Custom.Organization,
	}
}
