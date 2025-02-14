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

// ExternalProcessingEnvoyAttributes represents the attributes extracted from the external processing request.
type ExternalProcessingEnvoyAttributes struct {
	EnableBackendBasedAIRatelimit          string `json:"enableBackendBasedAIRatelimitAttribute"`
	SuspendAIModel                         string `json:"suspendAIModelAttribute"`
	BackendBasedAIRatelimitDescriptorValue string `json:"backendBasedAIRatelimitDescriptorValueAttribute"`
	Path                                   string `json:"pathAttribute"`
	VHost                                  string `json:"vHostAttribute"`
	BasePath                               string `json:"basePathAttribute"`
	Method                                 string `json:"methodAttribute"`
	APIVersion                             string `json:"apiVersionAttribute"`
	APIName                                string `json:"apiNameAttribute"`
	ClusterName                            string `json:"clusterNameAttribute"`
	RequestMethod                          string `json:"requestMethodAttribute"`
	RequestID                              string `json:"requestIdAttribute"`
	Organization                           string `json:"organizationAttribute"`
	ApplicationID                          string `json:"applicationIdAttribute"`
	CorrelationID                          string `json:"correlationIdAttribute"`
}
