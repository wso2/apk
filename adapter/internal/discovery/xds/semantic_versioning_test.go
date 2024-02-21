/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package xds

import (
	"regexp"
	"testing"

	"github.com/wso2/apk/adapter/config"
	semantic_version "github.com/wso2/apk/adapter/pkg/semanticversion"
)

func TestGetVersionMatchRegex(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		expectedResult string
	}{
		{
			name:           "Version with single digit components",
			version:        "1.2.3",
			expectedResult: "1\\.2\\.3",
		},
		{
			name:           "Version with multi-digit components",
			version:        "123.456.789",
			expectedResult: "123\\.456\\.789",
		},
		{
			name:           "Version with alpha components",
			version:        "v1.0-alpha",
			expectedResult: "v1\\.0-alpha",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetVersionMatchRegex(tt.version)

			if result != tt.expectedResult {
				t.Errorf("Expected regex: %s, Got: %s", tt.expectedResult, result)
			}

			// Test if the regex works correctly
			match, err := regexp.MatchString(result, tt.version)
			if err != nil {
				t.Errorf("Error when matching regex: %v", err)
			}
			if !match {
				t.Errorf("Regex failed to match the version: %s %s", tt.version, result)
			}
		})
	}
}

func TestGetMajorMinorVersionRangeRegex(t *testing.T) {
	tests := []struct {
		name           string
		semVersion     semantic_version.SemVersion
		expectedResult string
	}{
		{
			name:           "Major and minor version only",
			semVersion:     semantic_version.SemVersion{Major: 1, Minor: 2},
			expectedResult: "v1(?:\\.2)?",
		},
		{
			name:           "Major, minor, and patch version",
			semVersion:     semantic_version.SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedResult: "v1(?:\\.2(?:\\.3)?)?",
		},
		{
			name:           "Major version only",
			semVersion:     semantic_version.SemVersion{Major: 1},
			expectedResult: "v1(?:\\.0)?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMajorMinorVersionRangeRegex(tt.semVersion)

			if result != tt.expectedResult {
				t.Errorf("Expected regex: %s, Got: %s", tt.expectedResult, result)
			}
		})
	}
}

func TestGetMinorVersionRangeRegex(t *testing.T) {
	tests := []struct {
		name           string
		semVersion     semantic_version.SemVersion
		expectedResult string
	}{
		{
			name:           "Major, minor, and patch version",
			semVersion:     semantic_version.SemVersion{Version: "v1.2.3", Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedResult: "v1\\.2(?:\\.3)?",
		},
		{
			name:           "Major and minor version only",
			semVersion:     semantic_version.SemVersion{Version: "v1.2", Major: 1, Minor: 2},
			expectedResult: "v1\\.2",
		},
		{
			name:           "Major version only",
			semVersion:     semantic_version.SemVersion{Version: "v1", Major: 1},
			expectedResult: "v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMinorVersionRangeRegex(tt.semVersion)

			if result != tt.expectedResult {
				t.Errorf("Expected regex: %s, Got: %s", tt.expectedResult, result)
			}
		})
	}
}

func TestGetMajorVersionRange(t *testing.T) {
	tests := []struct {
		name           string
		semVersion     semantic_version.SemVersion
		expectedResult string
	}{
		{
			name:           "Major and minor version 1.2.3",
			semVersion:     semantic_version.SemVersion{Version: "v1.2.3", Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedResult: "v1",
		},
		{
			name:           "Major version 2",
			semVersion:     semantic_version.SemVersion{Major: 2},
			expectedResult: "v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMajorVersionRange(tt.semVersion)

			if result != tt.expectedResult {
				t.Errorf("Expected result: %s, Got: %s", tt.expectedResult, result)
			}
		})
	}
}

func TestGetMinorVersionRange(t *testing.T) {
	tests := []struct {
		name           string
		semVersion     semantic_version.SemVersion
		expectedResult string
	}{
		{
			name:           "Major and minor version 1.2",
			semVersion:     semantic_version.SemVersion{Major: 1, Minor: 2},
			expectedResult: "v1.2",
		},
		{
			name:           "Major and minor version 1.2.3",
			semVersion:     semantic_version.SemVersion{Version: "v1.2.3", Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedResult: "v1.2",
		},
		{
			name:           "Major only",
			semVersion:     semantic_version.SemVersion{Major: 10},
			expectedResult: "v10.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMinorVersionRange(tt.semVersion)

			if result != tt.expectedResult {
				t.Errorf("Expected result: %s, Got: %s", tt.expectedResult, result)
			}
		})
	}
}

func TestIsSemanticVersioningEnabled(t *testing.T) {

	conf := config.ReadConfigs()

	tests := []struct {
		name                      string
		apiName                   string
		apiVersion                string
		intelligentRoutingEnabled bool
		expectedResult            bool
	}{
		{
			name:                      "Semantic versioning enabled and valid version provided",
			apiName:                   "TestAPI",
			apiVersion:                "v1.2.3",
			intelligentRoutingEnabled: true,
			expectedResult:            true,
		},
		{
			name:                      "Semantic versioning enabled and valid version provided",
			apiName:                   "TestAPI",
			apiVersion:                "v1.2",
			intelligentRoutingEnabled: true,
			expectedResult:            true,
		},
		{
			name:                      "Semantic versioning enabled and version only contains major version",
			apiName:                   "TestAPI",
			apiVersion:                "v1",
			intelligentRoutingEnabled: true,
			expectedResult:            false,
		},
		{
			name:                      "Semantic versioning enabled and invalid version provided",
			apiName:                   "TestAPI",
			apiVersion:                "1.2.3",
			intelligentRoutingEnabled: true,
			expectedResult:            false,
		},
		{
			name:                      "Semantic versioning disabled and valid version provided",
			apiName:                   "TestAPI",
			apiVersion:                "v1.2.3",
			intelligentRoutingEnabled: false,
			expectedResult:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conf.Envoy.EnableIntelligentRouting = tt.intelligentRoutingEnabled
			result := isSemanticVersioningEnabled(tt.apiName, tt.apiVersion)

			if result != tt.expectedResult {
				t.Errorf("Expected result: %v, Got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestIsVHostMatched(t *testing.T) {
	// Mock orgIDAPIvHostsMap for testing
	orgIDAPIvHostsMap = map[string]map[string][]string{
		"org1": {
			"api1": {"example.com", "api.example.com"},
			"api2": {"test.com"},
		},
		"org2": {
			"api3": {"example.org"},
			"api4": {"test.org"},
		},
	}

	tests := []struct {
		name           string
		organizationID string
		vHost          string
		expectedResult bool
	}{
		{
			name:           "Matching vHost in org1",
			organizationID: "org1",
			vHost:          "example.com",
			expectedResult: true,
		},
		{
			name:           "Matching vHost in org2",
			organizationID: "org2",
			vHost:          "example.org",
			expectedResult: true,
		},
		{
			name:           "Non-matching vHost in org1",
			organizationID: "org1",
			vHost:          "nonexistent.com",
			expectedResult: false,
		},
		{
			name:           "Non-matching vHost in org2",
			organizationID: "org2",
			vHost:          "nonexistent.org",
			expectedResult: false,
		},
		{
			name:           "VHost not found for organization",
			organizationID: "org3",
			vHost:          "example.com",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVHostMatched(tt.organizationID, tt.vHost)

			if result != tt.expectedResult {
				t.Errorf("Expected result: %v, Got: %v", tt.expectedResult, result)
			}
		})
	}
}

// PtrInt returns a pointer to an integer value
func PtrInt(i int) *int {
	return &i
}
