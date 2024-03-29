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

	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_type_matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/stretchr/testify/assert"
	"github.com/wso2/apk/adapter/config"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
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
		semVersion     *semantic_version.SemVersion
		expectedResult string
	}{
		{
			name:           "Major and minor version only",
			semVersion:     &semantic_version.SemVersion{Major: 1, Minor: 2},
			expectedResult: "v1(?:\\.2)?",
		},
		{
			name:           "Major, minor, and patch version",
			semVersion:     &semantic_version.SemVersion{Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedResult: "v1(?:\\.2(?:\\.3)?)?",
		},
		{
			name:           "Major version only",
			semVersion:     &semantic_version.SemVersion{Major: 1},
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
		semVersion     *semantic_version.SemVersion
		expectedResult string
	}{
		{
			name:           "Major, minor, and patch version",
			semVersion:     &semantic_version.SemVersion{Version: "v1.2.3", Major: 1, Minor: 2, Patch: PtrInt(3)},
			expectedResult: "v1\\.2(?:\\.3)?",
		},
		{
			name:           "Major and minor version only",
			semVersion:     &semantic_version.SemVersion{Version: "v1.2", Major: 1, Minor: 2},
			expectedResult: "v1\\.2",
		},
		{
			name:           "Major version only",
			semVersion:     &semantic_version.SemVersion{Version: "v1", Major: 1},
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
			result := IsSemanticVersioningEnabled(tt.apiName, tt.apiVersion)

			if result != tt.expectedResult {
				t.Errorf("Expected result: %v, Got: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestUpdateRoutingRulesOnAPIUpdate(t *testing.T) {

	var apiID1 model.AdapterInternalAPI
	apiID1.SetName("Test API")
	apiID1.UUID = "apiID1"
	apiID1.OrganizationID = "org1"
	apiID1.SetVersion("v1.0")
	apiID1ResourcePath := "^/test-api/v1\\.0/orders([/]{0,1})"

	var apiID2 model.AdapterInternalAPI
	apiID2.SetName("Mock API")
	apiID2.UUID = "apiID2"
	apiID2.OrganizationID = "org1"
	apiID2.SetVersion("v1.1")
	apiID2ResourcePath := "^/mock-api/v1\\.1/orders([/]{0,1})"

	var apiID3 model.AdapterInternalAPI
	apiID3.SetName("Test API")
	apiID3.SetVersion("v1.1")
	apiID3.OrganizationID = "org1"
	apiID3.UUID = "apiID3"
	apiID3ResourcePath := "^/test-api/v1\\.1/orders([/]{0,1})"

	orgAPIMap = map[string]map[string]*EnvoyInternalAPI{
		"org1": {
			"gw.com:apiID1": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID1,
				routes:             generateRoutes(apiID1ResourcePath),
			},
			"gw.com:apiID2": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID2,
				routes:             generateRoutes(apiID2ResourcePath),
			},
			"gw.com:apiID3": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID3,
				routes:             generateRoutes(apiID3ResourcePath),
			},
		},
	}

	tests := []struct {
		name               string
		api                model.AdapterInternalAPI
		organizationID     string
		apiRangeIdentifier string
		apiIdentifier      string
		vhost              string
		expectedRegex      string
		expectedRewrite    string
		finalRegex         string
		finalRewrite       string
	}{
		{
			name:               "Create an API with major version",
			organizationID:     "org1",
			apiRangeIdentifier: "gw.com:Test API",
			apiIdentifier:      "gw.com:apiID1",
			vhost:              "gw.com",
			api:                apiID1,
			expectedRegex:      "^/test-api/v1(?:\\.0)?/orders([/]{0,1})",
			expectedRewrite:    "^/test-api/v1(?:\\.0)?/orders([/]{0,1})",
			finalRegex:         apiID1ResourcePath,
			finalRewrite:       apiID1ResourcePath,
		},
		{
			name:               "Create an API with major and minor version",
			organizationID:     "org1",
			apiRangeIdentifier: "gw.com:Mock API",
			apiIdentifier:      "gw.com:apiID2",
			vhost:              "gw.com",
			api:                apiID2,
			expectedRegex:      "^/mock-api/v1(?:\\.1)?/orders([/]{0,1})",
			expectedRewrite:    "^/mock-api/v1(?:\\.1)?/orders([/]{0,1})",
			finalRegex:         "^/mock-api/v1(?:\\.1)?/orders([/]{0,1})",
			finalRewrite:       "^/mock-api/v1(?:\\.1)?/orders([/]{0,1})",
		},
		{
			name:               "Create an API with major and minor version",
			organizationID:     "org1",
			apiRangeIdentifier: "gw.com:Test API",
			apiIdentifier:      "gw.com:apiID3",
			vhost:              "gw.com",
			api:                apiID3,
			expectedRegex:      "^/test-api/v1(?:\\.1)?/orders([/]{0,1})",
			expectedRewrite:    "^/test-api/v1(?:\\.1)?/orders([/]{0,1})",
			finalRegex:         "^/test-api/v1(?:\\.1)?/orders([/]{0,1})",
			finalRewrite:       "^/test-api/v1(?:\\.1)?/orders([/]{0,1})",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateSemanticVersioningInMapForUpdateAPI(tt.organizationID,
				map[string]struct{}{tt.apiRangeIdentifier: {}}, &tt.api)
			updateSemRegexForNewAPI(tt.api, orgAPIMap[tt.organizationID][tt.apiIdentifier].routes, tt.vhost)
			api1 := orgAPIMap[tt.organizationID][tt.apiIdentifier]
			routes := api1.routes

			if routes[0].GetMatch().GetSafeRegex().GetRegex() != tt.expectedRegex {
				t.Errorf("Expected regex: %s, Got: %s", tt.expectedRegex, routes[0].GetMatch().GetSafeRegex().GetRegex())
			}
			if routes[0].GetRoute().GetRegexRewrite().GetPattern().GetRegex() != tt.expectedRewrite {
				t.Errorf("Expected rewrite pattern: %s, Got: %s", tt.expectedRewrite, routes[0].GetRoute().GetRegexRewrite().GetPattern().GetRegex())
			}
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api1 := orgAPIMap[tt.organizationID][tt.apiIdentifier]
			routes := api1.routes

			if routes[0].GetMatch().GetSafeRegex().GetRegex() != tt.finalRegex {
				t.Errorf("Expected final regex: %s, Got: %s", tt.finalRegex, routes[0].GetMatch().GetSafeRegex().GetRegex())
			}
			if routes[0].GetRoute().GetRegexRewrite().GetPattern().GetRegex() != tt.finalRewrite {
				t.Errorf("Expected final rewrite pattern: %s, Got: %s", tt.finalRewrite, routes[0].GetRoute().GetRegexRewrite().GetPattern().GetRegex())
			}
		})
	}
}

func generateRoutes(resourcePath string) []*routev3.Route {

	var routes []*routev3.Route
	match := &routev3.RouteMatch{
		PathSpecifier: &routev3.RouteMatch_SafeRegex{
			SafeRegex: &envoy_type_matcherv3.RegexMatcher{
				Regex: resourcePath,
			},
		},
	}

	action := &routev3.Route_Route{
		Route: &routev3.RouteAction{
			RegexRewrite: &envoy_type_matcherv3.RegexMatchAndSubstitute{
				Pattern: &envoy_type_matcherv3.RegexMatcher{
					Regex: resourcePath,
				},
				Substitution: "/bar",
			},
		},
	}

	route := routev3.Route{
		Name:      "example-route",
		Match:     match,
		Action:    action,
		Metadata:  nil,
		Decorator: nil,
	}

	return append(routes, &route)
}

func TestUpdateRoutingRulesOnAPIDelete(t *testing.T) {

	orgIDLatestAPIVersionMap = map[string]map[string]map[string]semantic_version.SemVersion{
		"org3": {
			"gw.com:Test API": {
				"v1": {
					Version: "v1.0",
					Major:   1,
					Minor:   0,
					Patch:   nil,
				},
			},
		},
		"org4": {
			"gw.com:Mock API": {
				"v1.0": {
					Version: "v1.0",
					Major:   1,
					Minor:   0,
					Patch:   nil,
				},
				"v1.5": {
					Version: "v1.5",
					Major:   1,
					Minor:   5,
					Patch:   nil,
				},
				"v1": {
					Version: "v1.5",
					Major:   1,
					Minor:   5,
					Patch:   nil,
				},
			},
		},
	}

	var apiID1 model.AdapterInternalAPI
	apiID1.SetName("Test API")
	apiID1.UUID = "apiID1"
	apiID1.SetVersion("v1.0")
	apiID1ResourcePath := "^/test-api/v1\\.0/orders([/]{0,1})"

	var apiID2 model.AdapterInternalAPI
	apiID2.SetName("Mock API")
	apiID2.UUID = "apiID2"
	apiID2.SetVersion("v1.0")
	apiID2ResourcePath := "^/mock-api/v1\\.0/orders([/]{0,1})"

	var apiID21 model.AdapterInternalAPI
	apiID21.SetName("Mock API")
	apiID21.UUID = "apiID21"
	apiID21.SetVersion("v1.1")
	apiID21ResourcePath := "^/mock-api/v1\\.1/orders([/]{0,1})"

	var apiID25 model.AdapterInternalAPI
	apiID25.SetName("Mock API")
	apiID25.UUID = "apiID25"
	apiID25.SetVersion("v1.5")
	apiID25ResourcePath := "^/mock-api/v1(?:\\.5)?/orders([/]{0,1})"

	orgAPIMap = map[string]map[string]*EnvoyInternalAPI{
		"org3": {
			"gw.com:apiID1": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID1,
				routes:             generateRoutes(apiID1ResourcePath),
			},
		},
		"org4": {
			"gw.com:apiID2": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID2,
				routes:             generateRoutes(apiID2ResourcePath),
			},
			"gw.com:apiID21": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID21,
				routes:             generateRoutes(apiID21ResourcePath),
			},
			"gw.com:apiID25": &EnvoyInternalAPI{
				adapterInternalAPI: &apiID25,
				routes:             generateRoutes(apiID25ResourcePath),
			},
		},
	}

	tests := []struct {
		name           string
		organizationID string
		apiIdentifier  string
		apiCheck       string
		expectedRegex  string
		api            *model.AdapterInternalAPI
		deleteVersion  string
	}{
		{
			name:           "Delete latest major version",
			organizationID: "org3",
			apiIdentifier:  "gw.com:Test API",
			api:            &apiID1,
			deleteVersion:  "v1.0",
			apiCheck:       "gw.com:apiID25",
			expectedRegex:  "^/mock-api/v1(?:\\.5)?/orders([/]{0,1})",
		},
		{
			name:           "Delete latest minor version v1.5",
			organizationID: "org4",
			apiIdentifier:  "gw.com:Mock API",
			api:            &apiID25,
			deleteVersion:  "v1.5",
			apiCheck:       "gw.com:apiID21",
			expectedRegex:  "^/mock-api/v1(?:\\.1)?/orders([/]{0,1})",
		},
		{
			name:           "Delete latest minor version v1.1",
			organizationID: "org4",
			apiIdentifier:  "gw.com:Mock API",
			api:            &apiID21,
			deleteVersion:  "v1.1",
			apiCheck:       "gw.com:apiID2",
			expectedRegex:  "^/mock-api/v1(?:\\.0)?/orders([/]{0,1})",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemoveAPIFromAllInternalMaps(tt.api.UUID)
			if _, ok := orgIDLatestAPIVersionMap[tt.organizationID]; ok {
				if _, ok := orgIDLatestAPIVersionMap[tt.organizationID][tt.apiIdentifier]; ok {
					if _, ok := orgIDLatestAPIVersionMap[tt.organizationID][tt.apiIdentifier][tt.deleteVersion]; ok {
						t.Errorf("API deletion is not successful: %s", tt.deleteVersion)
					}
				}
			}
			if tt.apiCheck != "" {
				routes := orgAPIMap["org4"][tt.apiCheck].routes
				assert.Equal(t, tt.expectedRegex, routes[0].GetMatch().GetSafeRegex().GetRegex(),
					"Expected regex: %s, Got: %s",
				)
			}
		})
	}
}

// PtrInt returns a pointer to an integer value
func PtrInt(i int) *int {
	return &i
}
