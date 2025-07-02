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

package util

import (
	"encoding/json"
	"fmt"
	"strings"

	dpv2alpha1 "github.com/wso2/apk/common-go-libs/apis/dp/v2alpha1"
)

const delimPeriod = ":"

// PrepareAPIKey prepares the API key using the given vhost, basePath, and version.
func PrepareAPIKey(vhost, basePath, version string) string {
	return fmt.Sprintf("%s:%s:%s", vhost, basePath, version)
}

// NormalizePath normalizes the given path by removing backslashes.
func NormalizePath(input string) string {
	return strings.ReplaceAll(input, "\\", "")
}

// PrepareApplicationKeyMappingCacheKey generates a cache key for application key mapping.
func PrepareApplicationKeyMappingCacheKey(applicationIdentifier, keyType, securityScheme, envID string) string {
	return strings.Join([]string{securityScheme, envID, keyType, applicationIdentifier}, delimPeriod)
}

// ConvertToRoutePolicy converts a JSON string to a RoutePolicy object.
func ConvertToRoutePolicy(jsonStr string) (*dpv2alpha1.RoutePolicy, error) {
	var policy dpv2alpha1.RoutePolicy
	err := json.Unmarshal([]byte(jsonStr), &policy)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RoutePolicy: %w", err)
	}
	return &policy, nil
}

// ConvertToRouteMetadata converts a JSON string to a RouteMetadata object.
func ConvertToRouteMetadata(jsonStr string) (*dpv2alpha1.RouteMetadata, error) {
	var metadata dpv2alpha1.RouteMetadata
	err := json.Unmarshal([]byte(jsonStr), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal RouteMetadata: %w", err)
	}
	return &metadata, nil
}
