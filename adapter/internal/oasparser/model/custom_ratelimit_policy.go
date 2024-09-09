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

package model

import (
	dpv1alpha3 "github.com/wso2/apk/common-go-libs/apis/dp/v1alpha3"
)

// RateLimit is the rate limit values for a policy
type RateLimit struct {
	// RequestPerUnit is the number of requests allowed per unit time
	//
	RequestsPerUnit uint32 `json:"requestsPerUnit,omitempty"`

	// Unit is the unit of the requestPerUnit
	//
	Unit string `json:"unit,omitempty"`
}

// CustomRateLimitPolicy defines the desired state of CustomPolicy
type CustomRateLimitPolicy struct {
	Key          string    `json:"key,omitempty"`
	Value        string    `json:"value,omitempty"`
	RateLimit    RateLimit `json:"rateLimit,omitempty"`
	Organization string    `json:"organization,omitempty"`
}

// ParseCustomRateLimitPolicy parses the custom rate limit policy
func ParseCustomRateLimitPolicy(customRateLimitCR dpv1alpha3.RateLimitPolicy) *CustomRateLimitPolicy {
	rlPolicy := concatRateLimitPolicies(&customRateLimitCR, nil)
	return &CustomRateLimitPolicy{
		Key:   rlPolicy.Spec.Override.Custom.Key,
		Value: rlPolicy.Spec.Override.Custom.Value,
		RateLimit: RateLimit{
			RequestsPerUnit: rlPolicy.Spec.Override.Custom.RequestsPerUnit,
			Unit:            rlPolicy.Spec.Override.Custom.Unit,
		},
		Organization: rlPolicy.Spec.Override.Custom.Organization,
	}
}
