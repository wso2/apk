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
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/wso2/apk/adapter/config"
	logger "github.com/wso2/apk/adapter/internal/loggers"
	logging "github.com/wso2/apk/adapter/internal/logging"
	"github.com/wso2/apk/adapter/internal/oasparser/model"
	semantic_version "github.com/wso2/apk/adapter/pkg/semanticversion"
)

type oldSemVersion struct {
	Vhost              string
	APIName            string
	OldMajorSemVersion *semantic_version.SemVersion
	OldMinorSemVersion *semantic_version.SemVersion
}

// GetVersionMatchRegex returns the regex to match the full version string
func GetVersionMatchRegex(version string) string {
	// Match "." character in the version by replacing it with "\\."
	return strings.ReplaceAll(version, ".", "\\.")
}

// GetMajorMinorVersionRangeRegex generates major and minor version compatible range regex for the given version
func GetMajorMinorVersionRangeRegex(semVersion *semantic_version.SemVersion) string {
	majorVersion := strconv.Itoa(semVersion.Major)
	minorVersion := strconv.Itoa(semVersion.Minor)
	if semVersion.Patch == nil {
		return "v" + majorVersion + "(?:\\." + minorVersion + ")?"
	}
	patchVersion := strconv.Itoa(*semVersion.Patch)
	return "v" + majorVersion + "(?:\\." + minorVersion + "(?:\\." + patchVersion + ")?)?"
}

// GetMinorVersionRangeRegex generates minor version compatible range regex for the given version
func GetMinorVersionRangeRegex(semVersion *semantic_version.SemVersion) string {
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

// updateSemanticVersioningInMapForUpdateAPI updates the latest version ranges of the APIs in the organization
func updateSemanticVersioningInMapForUpdateAPI(org string, apiRangeIdentifiers map[string]struct{},
	adapterInternalAPI *model.AdapterInternalAPI) {
	oldSemVersions := make([]oldSemVersion, 0)
	if _, exist := orgIDLatestAPIVersionMap[org]; !exist {
		orgIDLatestAPIVersionMap[org] = make(map[string]map[string]semantic_version.SemVersion)
	}
	semVersion, _ := semantic_version.ValidateAndGetVersionComponents(adapterInternalAPI.GetVersion())

	for apiRangeIdentifier := range apiRangeIdentifiers {
		if currentAPISemVersion, exist := orgIDLatestAPIVersionMap[org][apiRangeIdentifier]; !exist {
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier] = make(map[string]semantic_version.SemVersion)
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier][GetMajorVersionRange(*semVersion)] = *semVersion
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier][GetMinorVersionRange(*semVersion)] = *semVersion
		} else {
			var oldVersion *oldSemVersion
			vhost, _ := ExtractVhostFromAPIIdentifier(apiRangeIdentifier)
			if _, ok := currentAPISemVersion[GetMajorVersionRange(*semVersion)]; !ok {
				currentAPISemVersion[GetMajorVersionRange(*semVersion)] = *semVersion
			} else if currentAPISemVersion[GetMajorVersionRange(*semVersion)].Compare(*semVersion) {
				version := currentAPISemVersion[GetMajorVersionRange(*semVersion)]
				currentAPISemVersion[GetMajorVersionRange(*semVersion)] = *semVersion
				oldVersion = &oldSemVersion{
					Vhost:              vhost,
					APIName:            adapterInternalAPI.GetTitle(),
					OldMajorSemVersion: &version,
				}
			}
			if _, ok := currentAPISemVersion[GetMinorVersionRange(*semVersion)]; !ok {
				currentAPISemVersion[GetMinorVersionRange(*semVersion)] = *semVersion
			} else if currentAPISemVersion[GetMinorVersionRange(*semVersion)].Compare(*semVersion) {
				version := currentAPISemVersion[GetMinorVersionRange(*semVersion)]
				currentAPISemVersion[GetMinorVersionRange(*semVersion)] = *semVersion
				if oldVersion != nil {
					oldVersion.OldMinorSemVersion = &version
				} else {
					oldVersion = &oldSemVersion{
						Vhost:              vhost,
						APIName:            adapterInternalAPI.GetTitle(),
						OldMinorSemVersion: &version,
					}
				}
			}
			if oldVersion != nil {
				oldSemVersions = append(oldSemVersions, *oldVersion)
			}
		}
	}
	updateOldRegex(org, oldSemVersions)
}

func updateOldRegex(org string, oldSemVersions []oldSemVersion) {
	if len(oldSemVersions) < 1 {
		return
	}
	for vuuid, api := range orgAPIMap[org] {
		// get vhost from the api identifier
		vhost, _ := ExtractVhostFromAPIIdentifier(vuuid)
		var oldSelectedSemVersion *oldSemVersion
		for _, oldSemVersion := range oldSemVersions {
			if oldSemVersion.Vhost == vhost && oldSemVersion.APIName == api.adapterInternalAPI.GetTitle() {
				if oldSemVersion.OldMajorSemVersion != nil && oldSemVersion.OldMajorSemVersion.Version ==
					api.adapterInternalAPI.GetVersion() {
					oldSelectedSemVersion = &oldSemVersion
					break
				}
				if oldSemVersion.OldMinorSemVersion != nil &&
					oldSemVersion.OldMinorSemVersion.Version == api.adapterInternalAPI.GetVersion() {
					oldSelectedSemVersion = &oldSemVersion
					break
				}
			}
		}
		logger.LoggerAPI.Error(oldSelectedSemVersion)
		if oldSelectedSemVersion == nil {
			continue
		}

		updateMajor := oldSelectedSemVersion.OldMajorSemVersion != nil && api.adapterInternalAPI.GetVersion() ==
			oldSelectedSemVersion.OldMajorSemVersion.Version
		updateMinor := oldSelectedSemVersion.OldMinorSemVersion != nil && api.adapterInternalAPI.GetVersion() ==
			oldSelectedSemVersion.OldMinorSemVersion.Version

		// apiSemVersion, _ := semantic_version.ValidateAndGetVersionComponents(api.adapterInternalAPI.GetVersion())
		for _, route := range api.routes {
			regex := route.GetMatch().GetSafeRegex().GetRegex()
			regexRewritePattern := route.GetRoute().GetRegexRewrite().GetPattern().GetRegex()
			if updateMajor {
				logger.LoggerAPI.Error(166, regex)
				regex = strings.Replace(regex, GetMajorMinorVersionRangeRegex(oldSelectedSemVersion.OldMajorSemVersion),
					GetMinorVersionRangeRegex(oldSelectedSemVersion.OldMajorSemVersion), 1)
				regexRewritePattern = strings.Replace(regexRewritePattern,
					GetMajorMinorVersionRangeRegex(oldSelectedSemVersion.OldMajorSemVersion),
					GetMinorVersionRangeRegex(oldSelectedSemVersion.OldMajorSemVersion), 1)
				logger.LoggerAPI.Error(166, regex)
			}
			if updateMinor {
				logger.LoggerAPI.Error(175, regex)
				regex = strings.Replace(regex, GetMinorVersionRangeRegex(oldSelectedSemVersion.OldMinorSemVersion),
					GetVersionMatchRegex(oldSelectedSemVersion.OldMinorSemVersion.Version), 1)
				regexRewritePattern = strings.Replace(regexRewritePattern,
					GetMinorVersionRangeRegex(oldSelectedSemVersion.OldMinorSemVersion),
					GetVersionMatchRegex(oldSelectedSemVersion.OldMinorSemVersion.Version), 1)
				logger.LoggerAPI.Error(180, regex)
			}
			pathSpecifier := &routev3.RouteMatch_SafeRegex{
				SafeRegex: &matcherv3.RegexMatcher{
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

// updateSemanticVersioningInMap updates the latest version ranges of the APIs in the organization
func updateSemanticVersioningInMap(org string, apiRangeIdentifiers map[string]struct{}) {
	oldSemVersions := make([]oldSemVersion, 0)
	// Iterate all the APIs in the API range
	for vuuid, api := range orgAPIMap[org] {
		// get vhost from the api identifier
		vhost, _ := ExtractVhostFromAPIIdentifier(vuuid)
		apiName := api.adapterInternalAPI.GetTitle()
		apiRangeIdentifier := GenerateIdentifierForAPIWithoutVersion(vhost, apiName)
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
		if _, exist := orgIDLatestAPIVersionMap[org]; !exist {
			orgIDLatestAPIVersionMap[org] = make(map[string]map[string]semantic_version.SemVersion)
		}
		if currentAPISemVersion, exist := orgIDLatestAPIVersionMap[org][apiRangeIdentifier]; !exist {
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier] = make(map[string]semantic_version.SemVersion)
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier][GetMajorVersionRange(*semVersion)] = *semVersion
			orgIDLatestAPIVersionMap[org][apiRangeIdentifier][GetMinorVersionRange(*semVersion)] = *semVersion
		} else {
			var oldVersion *oldSemVersion
			if _, ok := currentAPISemVersion[GetMajorVersionRange(*semVersion)]; !ok {
				currentAPISemVersion[GetMajorVersionRange(*semVersion)] = *semVersion
			} else if currentAPISemVersion[GetMajorVersionRange(*semVersion)].Compare(*semVersion) {
				version := currentAPISemVersion[GetMajorVersionRange(*semVersion)]
				oldVersion = &oldSemVersion{
					Vhost:              vhost,
					APIName:            apiName,
					OldMajorSemVersion: &version,
				}
				currentAPISemVersion[GetMajorVersionRange(*semVersion)] = *semVersion
			}
			if _, ok := currentAPISemVersion[GetMinorVersionRange(*semVersion)]; !ok {
				currentAPISemVersion[GetMinorVersionRange(*semVersion)] = *semVersion
			} else if currentAPISemVersion[GetMinorVersionRange(*semVersion)].Compare(*semVersion) {
				version := currentAPISemVersion[GetMinorVersionRange(*semVersion)]
				if oldVersion != nil {
					oldVersion.OldMinorSemVersion = &version
				} else {
					oldVersion = &oldSemVersion{
						Vhost:              vhost,
						APIName:            apiName,
						OldMinorSemVersion: &version,
					}
				}
				currentAPISemVersion[GetMinorVersionRange(*semVersion)] = *semVersion
			}
			if oldVersion != nil {
				oldSemVersions = append(oldSemVersions, *oldVersion)
			}
		}
	}
	updateOldRegex(org, oldSemVersions)
}

func updateSemRegexForNewAPI(adapterInternalAPI model.AdapterInternalAPI, routes []*routev3.Route, vhost string) {
	apiIdentifier := GenerateIdentifierForAPIWithoutVersion(vhost, adapterInternalAPI.GetTitle())

	if orgIDLatestAPIVersionMap[adapterInternalAPI.GetOrganizationID()] != nil &&
		orgIDLatestAPIVersionMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier] != nil {
		// get version list for the API
		apiVersionMap := orgIDLatestAPIVersionMap[adapterInternalAPI.GetOrganizationID()][apiIdentifier]
		// get the latest version of the API
		isMajorVersion := false
		update := false
		for versionRange, latestVersion := range apiVersionMap {
			if latestVersion.Version == adapterInternalAPI.GetVersion() {
				update = true
				versionArray := strings.Split(versionRange, ".")
				if len(versionArray) == 1 {
					isMajorVersion = true
					break
				}
			}
		}
		if update {
			// update regex
			for _, route := range routes {
				regex := route.GetMatch().GetSafeRegex().GetRegex()
				regexRewritePattern := route.GetRoute().GetRegexRewrite().GetPattern().GetRegex()
				apiVersionRegex := GetVersionMatchRegex(adapterInternalAPI.GetVersion())
				apiSemVersion, _ := semantic_version.ValidateAndGetVersionComponents(adapterInternalAPI.GetVersion())
				if isMajorVersion {
					regex = strings.Replace(regex, apiVersionRegex, GetMajorMinorVersionRangeRegex(apiSemVersion), 1)
					regexRewritePattern = strings.Replace(regexRewritePattern, apiVersionRegex, GetMajorMinorVersionRangeRegex(apiSemVersion), 1)
				} else {
					regex = strings.Replace(regex, apiVersionRegex, GetMinorVersionRangeRegex(apiSemVersion), 1)
					regexRewritePattern = strings.Replace(regexRewritePattern, apiVersionRegex, GetMinorVersionRangeRegex(apiSemVersion), 1)
				}
				pathSpecifier := &routev3.RouteMatch_SafeRegex{
					SafeRegex: &matcherv3.RegexMatcher{
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

// IsSemanticVersioningEnabled checks whether semantic versioning is enabled for the given API
func IsSemanticVersioningEnabled(apiName, apiVersion string) bool {
	conf := config.ReadConfigs()
	if !conf.Envoy.EnableIntelligentRouting {
		return false
	}

	apiSemVersion, err := semantic_version.ValidateAndGetVersionComponents(apiVersion)
	if err != nil || apiSemVersion == nil {
		logger.LoggerXds.ErrorC(logging.PrintError(logging.Error1411, logging.MAJOR,
			"Error validating the version of the API: %v:%s. Intelligent routing is disabled for the API, %v", apiName, apiVersion, err))
		return false
	}

	return true
}
