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

package constants

// Controller related constants
const (
	RatelimitController          string = "RatelimitController"
	ApplicationController        string = "ApplicationController"
	SubscriptionController       string = "SubscriptionController"
	ApplicationMappingController string = "ApplicationMappingController"
	GatewayClassController       string = "GatewayClassController"
)

// API events related constants
const (
	Create string = "CREATED"
	Update string = "UPDATED"
	Delete string = "DELETED"
)

// Subscriprion events related constants
const (
	ApplicationCreated           string = "APPLICATION_CREATED"
	ApplicationUpdated           string = "APPLICATION_UPDATED"
	ApplicationDeleted           string = "APPLICATION_DELETED"
	SubscriptionCreated          string = "SUBSCRIPTION_CREATED"
	SubscriptionUpdated          string = "SUBSCRIPTION_UPDATED"
	SubscriptionDeleted          string = "SUBSCRIPTION_DELETED"
	ApplicationMappingCreated    string = "APPLICATION_MAPPING_CREATED"
	ApplicationMappingUpdated    string = "APPLICATION_MAPPING_UPDATED"
	ApplicationMappingDeleted    string = "APPLICATION_MAPPING_DELETED"
	ApplicationKeyMappingCreated string = "APPLICATION_KEY_MAPPING_CREATED"
	ApplicationKeyMappingUpdated string = "APPLICATION_KEY_MAPPING_UPDATED"
	ApplicationKeyMappingDeleted string = "APPLICATION_KEY_MAPPING_DELETED"
	AllEvents                    string = "ALL_EVENTS"
)

// Environment variable names and default values
const (
	OperatorPodNamespace             string = "OPERATOR_POD_NAMESPACE"
	OperatorPodNamespaceDefaultValue string = "default"
)

// CRD Kinds
const (
	KindAuthentication = "Authentication"
	KindAPI            = "API"
	KindService        = "Service"
	All                = "All"
	//TODO(amali) remove this after fixing the issue in https://github.com/wso2/apk/issues/383
	KindResource        = "Resource"
	KindGateway         = "Gateway"
	KindRateLimitPolicy = "RateLimitPolicy"
)

// Env types
const (
	Production = "PRODUCTION"
	Sandbox    = "SANDBOX"
)

// Header names in runtime
const (
	OrganizationHeader = "X-WSO2-Organization"
)

// Global interceptor cluster names
const (
	GlobalRequestInterceptorClusterName  = "request_interceptor_global_cluster"
	GlobalResponseInterceptorClusterName = "response_interceptor_global_cluster"
)

// Application authentication types
const (
	OAuth2 = "OAuth2"
)
