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

package dto

// FaultSubCategory defines a type for fault subcategories.
type FaultSubCategory string

// Subcategories of the "Other" category.
const (
	OtherMediationError   FaultSubCategory = "MEDIATION_ERROR"
	OtherResourceNotFound FaultSubCategory = "RESOURCE_NOT_FOUND"
	OtherMethodNotAllowed FaultSubCategory = "METHOD_NOT_ALLOWED"
	OtherUnclassified     FaultSubCategory = "UNCLASSIFIED"
)

// Subcategories of the "Throttling" category.
const (
	ThrottlingAPILimitExceeded          FaultSubCategory = "API_LEVEL_LIMIT_EXCEEDED"
	ThrottlingHardLimitExceeded         FaultSubCategory = "HARD_LIMIT_EXCEEDED"
	ThrottlingResourceLimitExceeded     FaultSubCategory = "RESOURCE_LEVEL_LIMIT_EXCEEDED"
	ThrottlingApplicationLimitExceeded  FaultSubCategory = "APPLICATION_LEVEL_LIMIT_EXCEEDED"
	ThrottlingSubscriptionLimitExceeded FaultSubCategory = "SUBSCRIPTION_LIMIT_EXCEEDED"
	ThrottlingBlocked                   FaultSubCategory = "BLOCKED"
	ThrottlingCustomPolicyLimitExceeded FaultSubCategory = "CUSTOM_POLICY_LIMIT_EXCEEDED"
	ThrottlingBurstControlLimitExceeded FaultSubCategory = "BURST_CONTROL_LIMIT_EXCEEDED"
	ThrottlingQueryTooDeep              FaultSubCategory = "QUERY_TOO_DEEP"
	ThrottlingQueryTooComplex           FaultSubCategory = "QUERY_TOO_COMPLEX"
	ThrottlingOther                     FaultSubCategory = "OTHER"
)

// Subcategories of the "TargetConnectivity" category.
const (
	TargetConnectivityConnectionTimeout   FaultSubCategory = "CONNECTION_TIMEOUT"
	TargetConnectivityConnectionSuspended FaultSubCategory = "CONNECTION_SUSPENDED"
	TargetConnectivityOther               FaultSubCategory = "OTHER"
)

// Subcategories of the "Authentication" category.
const (
	AuthenticationFailure                       FaultSubCategory = "AUTHENTICATION_FAILURE"
	AuthenticationAuthorizationFailure          FaultSubCategory = "AUTHORIZATION_FAILURE"
	AuthenticationSubscriptionValidationFailure FaultSubCategory = "SUBSCRIPTION_VALIDATION_FAILURE"
	AuthenticationOther                         FaultSubCategory = "OTHER"
)
