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

package semanticversion

import (
	"errors"
	"testing"
)

func TestSemVersionCompare(t *testing.T) {
	tests := []struct {
		name           string
		baseVersion    SemVersion
		compareVersion SemVersion
		expected       bool
	}{
		{
			name:           "Same versions",
			baseVersion:    SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			expected:       true,
		},
		{
			name:           "Base version major is greater",
			baseVersion:    SemVersion{Major: 2, Minor: 1, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			expected:       false,
		},
		{
			name:           "Base version minor is greater",
			baseVersion:    SemVersion{Major: 1, Minor: 3, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 1, Minor: 4, Patch: PtrInt(3)},
			expected:       true,
		},
		{
			name:           "Base version patch is greater",
			baseVersion:    SemVersion{Major: 1, Minor: 2, Patch: PtrInt(4)},
			compareVersion: SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			expected:       false,
		},
		{
			name:           "Compare version major is greater",
			baseVersion:    SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 2, Minor: 2, Patch: PtrInt(3)},
			expected:       true,
		},
		{
			name:           "Compare version minor is greater",
			baseVersion:    SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 1, Minor: 3, Patch: PtrInt(3)},
			expected:       true,
		},
		{
			name:           "Compare version patch is greater",
			baseVersion:    SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 1, Minor: 2, Patch: PtrInt(4)},
			expected:       true,
		},
		{
			name:           "Base version patch is nil",
			baseVersion:    SemVersion{Major: 1, Minor: 2},
			compareVersion: SemVersion{Major: 1, Minor: 2, Patch: PtrInt(4)},
			expected:       true,
		},
		{
			name:           "Compare version patch is nil",
			baseVersion:    SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			compareVersion: SemVersion{Major: 1, Minor: 2},
			expected:       false,
		},
		{
			name:           "Both patch versions are nil",
			baseVersion:    SemVersion{Major: 1, Minor: 2},
			compareVersion: SemVersion{Major: 1, Minor: 2},
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.baseVersion.Compare(tt.compareVersion)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestValidateAndGetVersionComponents(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		apiName        string
		expectedResult *SemVersion
		expectedError  error
	}{
		{
			name:           "Valid version format",
			version:        "v1.2.3",
			apiName:        "TestAPI",
			expectedResult: &SemVersion{Version: "v1.2.3", Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedError:  nil,
		},
		{
			name:           "Valid version format without patch",
			version:        "v1.2",
			apiName:        "TestAPI",
			expectedResult: &SemVersion{Version: "v1.2", Major: 1, Minor: 2, Patch: nil},
			expectedError:  nil,
		},
		{
			name:           "Invalid version format - missing 'v' prefix",
			version:        "1.2.3",
			apiName:        "TestAPI",
			expectedResult: nil,
			expectedError:  errors.New("invalid version: 1.2.3. API version should be in the format x.y.z, x.y, vx.y.z or vx.y where x,y,z are non-negative integers and v is version prefix"),
		},
		{
			name:           "Invalid version format - negative major version",
			version:        "v-1.2.3",
			apiName:        "TestAPI",
			expectedResult: nil,
			expectedError:  errors.New("invalid version format. API major version should be a non-negative integer in API Version: v-1.2.3, <nil>"),
		},
		{
			name:           "Invalid version format - negative minor version",
			version:        "v1.-2.3",
			apiName:        "TestAPI",
			expectedResult: nil,
			expectedError:  errors.New("invalid version format. API minor version should be a non-negative integer in API Version: v1.-2.3, <nil>"),
		},
		{
			name:           "Invalid version format - patch version not an integer",
			version:        "v1.2.three",
			apiName:        "TestAPI",
			expectedResult: nil,
			expectedError:  errors.New("invalid version format. API patch version should be an integer in API Version: v1.2.three, strconv.Atoi: parsing \"three\": invalid syntax"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndGetVersionComponents(tt.version)

			// Check for errors
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("Unexpected error. Expected: %v, Got: %v", tt.expectedError, err)
			}

			// Check for nil results
			if result == nil && tt.expectedResult != nil {
				t.Errorf("Unexpected nil result")
			} else if result != nil && tt.expectedResult == nil {
				t.Errorf("Unexpected non-nil result")
			}

			// Check for result equality
			if result != nil && tt.expectedResult != nil {
				if result.Version != tt.expectedResult.Version || result.Major != tt.expectedResult.Major || result.Minor != tt.expectedResult.Minor || (result.Patch != nil && (*result.Patch != *tt.expectedResult.Patch)) {
					t.Errorf("Unexpected result. Expected: %v, Got: %v", tt.expectedResult, result)
				}
			}
		})
	}

}

// PtrInt returns a pointer to an integer value
func PtrInt(i int) *int {
	return &i
}
