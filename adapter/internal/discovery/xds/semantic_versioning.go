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

package xds

import (
	"strconv"
	"strings"

	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_type_matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	semantic_version "github.com/wso2/apk/adapter/pkg/semanticversion"
)

// GetVersionMatchRegex returns the regex to match the full version string
func GetVersionMatchRegex(version string) string {
	// Match "." character in the version by replacing it with "\\."
	return strings.ReplaceAll(version, ".", "\\.")
}

// GetMajorMinorVersionRangeRegex generates major and minor version compatible range regex for the given version
func GetMajorMinorVersionRangeRegex(semVersion semantic_version.SemVersion) string {
	majorVersion := strconv.Itoa(semVersion.Major)
	minorVersion := strconv.Itoa(semVersion.Minor)
	if semVersion.Patch == nil {
		return "v" + majorVersion + "(?:\\." + minorVersion + ")?"
	}
	patchVersion := strconv.Itoa(*semVersion.Patch)
	return "v" + majorVersion + "(?:\\." + minorVersion + "(?:\\." + patchVersion + ")?)?"
}

// GetMinorVersionRangeRegex generates minor version compatible range regex for the given version
func GetMinorVersionRangeRegex(semVersion semantic_version.SemVersion) string {
	if semVersion.Patch == nil {
		return GetVersionMatchRegex(semVersion.Version)
	}
	majorVersion := strconv.Itoa(semVersion.Major)
	minorVersion := strconv.Itoa(semVersion.Minor)
	patchVersion := strconv.Itoa(*semVersion.Patch)
	return "v" + majorVersion + "\\." + minorVersion + "(?:\\." + patchVersion + ")?"
}

// GetMajorVersionRange generates major version range for the given version
func GetMajorVersionRange(semVersion semantic_version.SemVersion) string {
	return "v" + strconv.Itoa(semVersion.Major)
}

// GetMinorVersionRange generates minor version range for the given version
func GetMinorVersionRange(semVersion semantic_version.SemVersion) string {
	return "v" + strconv.Itoa(semVersion.Major) + "." + strconv.Itoa(semVersion.Minor)
}

func updateRoutingRulesOnAPIUpdate(organizationID, apiIdentifier, apiName, apiVersion, vHost string) {
	apiSemVersion, err := semantic_version.ValidateAndGetVersionComponents(apiVersion)
	// If the version validation is not success, we just proceed without intelligent version
	// Valid version pattern: vx.y.z or vx.y where x, y and z are non-negative integers and v is a prefix
	if err != nil && apiSemVersion == nil {
		return
	}

	apiRangeIdentifier := generateIdentifierForAPIWithoutVersion(vHost, apiName)
	// Check the major and minor version ranges of the current API
	existingMajorRangeLatestSemVersion, isMajorRangeRegexAvailable :=
		orgIDLatestAPIVersionMap[organizationID][apiRangeIdentifier][GetMajorVersionRange(*apiSemVersion)]
	existingMinorRangeLatestSemVersion, isMinorRangeRegexAvailable :=
		orgIDLatestAPIVersionMap[organizationID][apiRangeIdentifier][GetMinorVersionRange(*apiSemVersion)]

	// Check whether the current API is the latest version in the major and minor version ranges
	isLatestMajorVersion := !isMajorRangeRegexAvailable || existingMajorRangeLatestSemVersion.Compare(*apiSemVersion)
	isLatestMinorVersion := !isMinorRangeRegexAvailable || existingMinorRangeLatestSemVersion.Compare(*apiSemVersion)

	// Remove the existing regexes from the path specifier when latest major and/or minor version is available
	if (isMajorRangeRegexAvailable || isMinorRangeRegexAvailable) && (isLatestMajorVersion || isLatestMinorVersion) {
		// Organization's all apis
		for vuuid, envoyInternalAPI := range orgAPIMap[organizationID] {
			// API's all versions in the same vHost
			if envoyInternalAPI.adapterInternalAPI.GetTitle() == apiName && strings.HasPrefix(vuuid+":", vHost) {

				if (isMajorRangeRegexAvailable && envoyInternalAPI.adapterInternalAPI.GetVersion() == existingMajorRangeLatestSemVersion.Version) ||
					(isMinorRangeRegexAvailable && envoyInternalAPI.adapterInternalAPI.GetVersion() == existingMinorRangeLatestSemVersion.Version) {

					for _, route := range envoyInternalAPI.routes {
						regex := route.GetMatch().GetSafeRegex().GetRegex()
						regexRewritePattern := route.GetRoute().GetRegexRewrite().GetPattern().GetRegex()
						existingMinorRangeLatestVersionRegex := GetVersionMatchRegex(existingMinorRangeLatestSemVersion.Version)
						existingMajorRangeLatestVersionRegex := GetVersionMatchRegex(existingMajorRangeLatestSemVersion.Version)
						if isMinorRangeRegexAvailable && envoyInternalAPI.adapterInternalAPI.GetVersion() == existingMinorRangeLatestSemVersion.Version && isLatestMinorVersion {
							regex = strings.Replace(regex, GetMinorVersionRangeRegex(existingMinorRangeLatestSemVersion), existingMinorRangeLatestVersionRegex, 1)
							regex = strings.Replace(regex, GetMajorMinorVersionRangeRegex(existingMajorRangeLatestSemVersion), existingMajorRangeLatestVersionRegex, 1)
							regexRewritePattern = strings.Replace(regexRewritePattern, GetMinorVersionRangeRegex(existingMinorRangeLatestSemVersion), existingMinorRangeLatestVersionRegex, 1)
							regexRewritePattern = strings.Replace(regexRewritePattern, GetMajorMinorVersionRangeRegex(existingMajorRangeLatestSemVersion), existingMajorRangeLatestVersionRegex, 1)
						}
						if isMajorRangeRegexAvailable && envoyInternalAPI.adapterInternalAPI.GetVersion() == existingMajorRangeLatestSemVersion.Version && isLatestMajorVersion {
							regex = strings.Replace(regex, GetMajorMinorVersionRangeRegex(existingMajorRangeLatestSemVersion), GetMinorVersionRangeRegex(existingMajorRangeLatestSemVersion), 1)
							regexRewritePattern = strings.Replace(regexRewritePattern, GetMajorMinorVersionRangeRegex(existingMajorRangeLatestSemVersion), GetMinorVersionRangeRegex(existingMajorRangeLatestSemVersion), 1)
						}
						pathSpecifier := &routev3.RouteMatch_SafeRegex{
							SafeRegex: &envoy_type_matcherv3.RegexMatcher{
								Regex: regex,
							},
						}
						route.Match.PathSpecifier = pathSpecifier
						action := route.Action.(*routev3.Route_Route)
						action.Route.RegexRewrite.Pattern.Regex = regexRewritePattern
						route.Action = action
					}
				}
			}
		}
	}

	if isLatestMajorVersion || isLatestMinorVersion {
		// Update local memory map with the latest version ranges
		majorVersionRange := GetMajorVersionRange(*apiSemVersion)
		minorVersionRange := GetMinorVersionRange(*apiSemVersion)
		if _, orgExists := orgIDLatestAPIVersionMap[organizationID]; !orgExists {
			orgIDLatestAPIVersionMap[organizationID] = make(map[string]map[string]semantic_version.SemVersion)
		}
		if _, apiRangeExists := orgIDLatestAPIVersionMap[organizationID][apiRangeIdentifier]; !apiRangeExists {
			orgIDLatestAPIVersionMap[organizationID][apiRangeIdentifier] = make(map[string]semantic_version.SemVersion)
		}

		latestVersions := orgIDLatestAPIVersionMap[organizationID][apiRangeIdentifier]
		latestVersions[minorVersionRange] = *apiSemVersion
		if isLatestMajorVersion {
			latestVersions[majorVersionRange] = *apiSemVersion
		}

		// Add the major and/or minor version range matching regexes to the path specifier when
		// latest major and/or minor version is available
		apiRoutes := getRoutesForAPIIdentifier(organizationID, apiIdentifier)

		for _, route := range apiRoutes {
			regex := route.GetMatch().GetSafeRegex().GetRegex()
			regexRewritePattern := route.GetRoute().GetRegexRewrite().GetPattern().GetRegex()
			apiVersionRegex := GetVersionMatchRegex(apiVersion)

			if isLatestMajorVersion {
				regex = strings.Replace(regex, apiVersionRegex, GetMajorMinorVersionRangeRegex(*apiSemVersion), 1)
				regexRewritePattern = strings.Replace(regexRewritePattern, apiVersionRegex, GetMajorMinorVersionRangeRegex(*apiSemVersion), 1)
			} else if isLatestMinorVersion {
				regex = strings.Replace(regex, apiVersionRegex, GetMinorVersionRangeRegex(*apiSemVersion), 1)
				regexRewritePattern = strings.Replace(regexRewritePattern, apiVersionRegex, GetMinorVersionRangeRegex(*apiSemVersion), 1)
			}
			pathSpecifier := &routev3.RouteMatch_SafeRegex{
				SafeRegex: &envoy_type_matcherv3.RegexMatcher{
					Regex: regex,
				},
			}

			route.Match.PathSpecifier = pathSpecifier
			action := route.Action.(*routev3.Route_Route)
			action.Route.RegexRewrite.Pattern.Regex = regexRewritePattern
			route.Action = action
		}

	}
}

func updateSemanticVersioning(org string, apiRangeIdentifiers map[string]struct{}) {
	// Iterate all the APIs in the API range
	for vuuid, api := range orgAPIMap[org] {
		// get vhost from the api identifier
		vhost, _ := ExtractVhostFromAPIIdentifier(vuuid)
		apiName := api.adapterInternalAPI.GetTitle()
		apiRangeIdentifier := generateIdentifierForAPIWithoutVersion(vhost, apiName)
		if _, ok := apiRangeIdentifiers[apiRangeIdentifier]; !ok {
			continue
		}
		// get sem version from the api in orgmap
		semVersion, err := semantic_version.ValidateAndGetVersionComponents(api.adapterInternalAPI.GetVersion())
		if err != nil {
			logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1410, logging.MAJOR,
				"Error validating the version of the API: %v for Organization: %v. Ignore deploying the API, error: %v",
				vuuid, org, err))
			continue
		}
		if currentAPISemVersion, exist := orgIDLatestAPIVersionMap[org][apiRangeIdentifier]; !exist {
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier] = make(map[string]semantic_version.SemVersion)
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier][GetMajorVersionRange(*semVersion)] = *semVersion
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier][GetMinorVersionRange(*semVersion)] = *semVersion

		} else {
			if _, ok := currentAPISemVersion[GetMajorVersionRange(*semVersion)]; !ok {
				currentAPISemVersion[GetMajorVersionRange(*semVersion)] = *semVersion
			} else {
				if currentAPISemVersion[GetMajorVersionRange(*semVersion)].Compare(*semVersion) {
					currentAPISemVersion[GetMajorVersionRange(*semVersion)] = *semVersion
				}
			}
			if _, ok := currentAPISemVersion[GetMinorVersionRange(*semVersion)]; !ok {
				currentAPISemVersion[GetMinorVersionRange(*semVersion)] = *semVersion
			} else {
				if currentAPISemVersion[GetMinorVersionRange(*semVersion)].Compare(*semVersion) {
					currentAPISemVersion[GetMinorVersionRange(*semVersion)] = *semVersion
				}
			}
		}
	}
}

func getRoutesForAPIIdentifier(organizationID, apiIdentifier string) []*routev3.Route {

	var routes []*routev3.Route
	if _, ok := orgAPIMap[organizationID]; ok {
		if _, ok := orgAPIMap[organizationID][apiIdentifier]; ok {
			routes = orgAPIMap[organizationID][apiIdentifier].routes
		}
	}

	return routes
}

func isSemanticVersioningEnabled(apiName, apiVersion string) bool {

	conf := config.ReadConfigs()
	if !conf.Envoy.EnableIntelligentRouting {
		return false
	}

	apiSemVersion, err := semantic_version.ValidateAndGetVersionComponents(apiVersion)
	if err != nil && apiSemVersion == nil {
		return false
	}

	return true
}
