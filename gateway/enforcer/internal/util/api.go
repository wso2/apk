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
	"fmt"
	"strings"
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
