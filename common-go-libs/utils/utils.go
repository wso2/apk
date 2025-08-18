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

// Package common includes the common functions shared between enforcer and router callbacks.
package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/wso2/apk/common-go-libs/constants"
)

// GetOperatorPodNamespace returns the namesapce of the operator pod
func GetOperatorPodNamespace() string {
	return GetEnv(constants.OperatorPodNamespace,
		constants.OperatorPodNamespaceDefaultValue)
}

// GetEnv lookup environment variable with key,
// if not defined returns default value
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// HashLast50SHA1 returns the last 50 characters of the SHA-1 hash of the input string.
// Since SHA-1 is only 40 hex chars, this will just return the whole hash.
func HashLast50SHA1(input string) string {
	hash := sha1.Sum([]byte(input))
	hexStr := hex.EncodeToString(hash[:]) // SHA-1 produces 40 hex chars
	if len(hexStr) <= 50 {
		return hexStr
	}
	return hexStr[len(hexStr)-50:]
}

// CreateAIProviderName creates a unique name for the AI provider based on the provider name and API version.
func CreateAIProviderName(providerName string, ProviderAPIVersion string) string {
	// Create a unique name for the AI provider by hashing the provider name
	// This ensures that the name is unique and does not exceed length limits
	return "ai-provider-" + HashLast50SHA1(providerName+ProviderAPIVersion)
}
