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

// Target represents target attributes in an analytics event.
type Target struct {
	TargetResponseCode int    `json:"targetResponseCode"`
	ResponseCacheHit   bool   `json:"responseCacheHit"`
	Destination        string `json:"destination"`
	ResponseCodeDetail string `json:"responseCodeDetail"`
}

// GetTargetResponseCode returns the target response code.
func (t *Target) GetTargetResponseCode() int {
	return t.TargetResponseCode
}

// SetTargetResponseCode sets the target response code.
func (t *Target) SetTargetResponseCode(targetResponseCode int) {
	t.TargetResponseCode = targetResponseCode
}

// IsResponseCacheHit returns whether the response was a cache hit.
func (t *Target) IsResponseCacheHit() bool {
	return t.ResponseCacheHit
}

// SetResponseCacheHit sets whether the response was a cache hit.
func (t *Target) SetResponseCacheHit(responseCacheHit bool) {
	t.ResponseCacheHit = responseCacheHit
}

// GetDestination returns the destination.
func (t *Target) GetDestination() string {
	return t.Destination
}

// SetDestination sets the destination.
func (t *Target) SetDestination(destination string) {
	t.Destination = destination
}
