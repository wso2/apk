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

// AIProvider represents the AI provider details.
type AIProvider struct {
	ProviderName       string        `json:"providerName"`       // Name of the AI provider
	ProviderAPIVersion string        `json:"providerAPIVersion"` // API version of the AI provider
	Organization       string        `json:"organization"`       // Organization associated with the provider
	Enabled            bool          `json:"enabled"`            // Whether the provider is enabled
	SupportedModels    []string      `json:"supportedModels"`    // Supported models
	Model              *ValueDetails `json:"model"`              // Model details
	PromptTokens       *ValueDetails `json:"promptTokens"`       // Prompt token details
	CompletionToken    *ValueDetails `json:"completionToken"`    // Completion token details
	TotalToken         *ValueDetails `json:"totalToken"`         // Total token details
}
