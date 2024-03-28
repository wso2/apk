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
 */

package semanticversion

import (
	"fmt"
	"strconv"
	"strings"
)

// SemVersion is the struct to store the version components of an API
type SemVersion struct {
	Version string
	Major   int
	Minor   int
	Patch   *int
}

// ValidateAndGetVersionComponents validates version string and extracts version components
func ValidateAndGetVersionComponents(version string) (*SemVersion, error) {
	versionComponents := strings.Split(version, ".")

	// If the versionComponents length is less than 2, return error
	if len(versionComponents) < 2 || !strings.HasPrefix(versionComponents[0], "v") {
		return nil, fmt.Errorf("invalid version: %v. API version should be in the format x.y.z, x.y, vx.y.z or vx.y where x,y,z are non-negative integers"+
			" and v is version prefix", version)
	}

	majorVersionStr := strings.TrimPrefix(versionComponents[0], "v")

	majorVersion, majorVersionConvErr := strconv.Atoi(majorVersionStr)
	minorVersion, minorVersionConvErr := strconv.Atoi(versionComponents[1])
	if majorVersionConvErr != nil || majorVersion < 0 {
		return nil, fmt.Errorf("invalid version format. API major version should be a non-negative integer in API Version: %v, %v", version, majorVersionConvErr)
	}

	if minorVersionConvErr != nil || minorVersion < 0 {
		return nil, fmt.Errorf("invalid version format. API minor version should be a non-negative integer in API Version: %v, %v", version, minorVersionConvErr)
	}

	if len(versionComponents) == 2 {
		return &SemVersion{
			Version: version,
			Major:   majorVersion,
			Minor:   minorVersion,
			Patch:   nil,
		}, nil
	}

	patchVersion, patchVersionConvErr := strconv.Atoi(versionComponents[2])
	if patchVersionConvErr != nil {
		return nil, fmt.Errorf("invalid version format. API patch version should be an integer in API Version: %v, %v", version, patchVersionConvErr)
	}
	return &SemVersion{
		Version: version,
		Major:   majorVersion,
		Minor:   minorVersion,
		Patch:   &patchVersion,
	}, nil
}

// Compare - compares two semantic versions and returns true
// if `version` is greater or equal than `baseVersion`
func (baseVersion SemVersion) Compare(version SemVersion) bool {
	// Compare major version
	if baseVersion.Major != version.Major {
		return baseVersion.Major < version.Major
	}

	// Compare minor version
	if baseVersion.Minor != version.Minor {
		return baseVersion.Minor < version.Minor
	}

	// Compare patch version
	if baseVersion.Patch != nil && version.Patch != nil {
		return *baseVersion.Patch <= *version.Patch
	} else if baseVersion.Patch != nil {
		return false
	} else if version.Patch != nil {
		return true
	}

	// Versions are equal
	return true
}
