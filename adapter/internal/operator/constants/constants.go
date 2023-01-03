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
	APIController         string = "APIController"
	ApplicationController string = "ApplicationController"
)

// API events related constants
const (
	Create string = "CREATE"
	Update string = "UPDATE"
	Delete string = "DELETE"
)

// Environment variable names and default values
const (
	OperatorPodNamespace             string = "OPERATOR_POD_NAMESPACE"
	OperatorPodNamespaceDefaultValue string = "default"
)

// CR Statuses
const (
	DeployedState = "Deployed"
	UpdatedState  = "Updated"
)

// CRD Kinds
const (
	KindAuthentication = "Authentication"
	KindHTTPRoute      = "HTTPRoute"
	//TODO(amali) remove this after fixing the issue in https://github.com/wso2/apk/issues/383
	KindResource = "Resource"
)
