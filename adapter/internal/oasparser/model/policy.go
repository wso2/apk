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

// OperationPolicies holds policies of the APIM operations
type OperationPolicies struct {
	Request  PolicyList `json:"request,omitempty"`
	Response PolicyList `json:"response,omitempty"`
	Fault    PolicyList `json:"fault,omitempty"`
}

// PolicyList holds list of Polices in a flow of operation
type PolicyList []Policy

// Policy holds APIM policies
type Policy struct {
	PolicyName       string      `json:"policyName,omitempty"`
	PolicyVersion    string      `json:"policyVersion,omitempty"`
	Action           string      `json:"-"` // This is a meta value used in APK, not included in API YAML
	IsPassToEnforcer bool        `json:"-"` // This is a meta value used in APK, not included in API YAML
	Parameters       interface{} `json:"parameters,omitempty"`
}
