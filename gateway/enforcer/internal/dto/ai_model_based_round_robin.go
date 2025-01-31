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

// AIModelBasedRoundRobin represents the AI model-based round robin configuration.
type AIModelBasedRoundRobin struct {
	Enabled                      bool          `json:"enabled"` // Whether AI model-based round robin is enabled
	OnQuotaExceedSuspendDuration int           `json:"onQuotaExceedSuspendDuration,omitempty"`
	Models                       []ModelWeight `json:"models"`
}

// ModelWeight holds the model configurations
type ModelWeight struct {
	Model  string `json:"model"`
	Weight int    `json:"weight,omitempty"`
}
