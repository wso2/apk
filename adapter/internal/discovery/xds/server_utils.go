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

package xds

import (
	logger "github.com/wso2/apk/adapter/internal/loggers"
)

// getEnvironmentsToBeDeleted returns an slice of environments APIs to be u-deployed from
// by considering existing environments list and environments that APIs are wished to be un-deployed
func getEnvironmentsToBeDeleted(existingEnvs, deleteEnvs []string) (toBeDel []string, toBeKept []string) {
	toBeDel = make([]string, 0, len(deleteEnvs))
	toBeKept = make([]string, 0, len(deleteEnvs))

	// if deleteEnvs is empty (deleteEnvs wished to be deleted), delete all environments
	if len(deleteEnvs) == 0 {
		return existingEnvs, []string{}
	}
	// otherwise delete env if it wished to
	for _, existingEnv := range existingEnvs {
		if arrayContains(deleteEnvs, existingEnv) {
			toBeDel = append(toBeDel, existingEnv)
		} else {
			toBeKept = append(toBeKept, existingEnv)
		}
	}
	return
}

func updateVhostInternalMaps(uuid, name, version, vHost string, gwEnvs []string) {

	uniqueIdentifier := uuid

	if uniqueIdentifier == "" {
		// If API is imported from apictl, get the hash value of API name and version
		uniqueIdentifier = GenerateHashedAPINameVersionIDWithoutVhost(name, version)
	}
	// update internal map: apiToVhostsMap
	if _, ok := apiToVhostsMap[uniqueIdentifier]; ok {
		apiToVhostsMap[uniqueIdentifier][vHost] = void
	} else {
		apiToVhostsMap[uniqueIdentifier] = map[string]struct{}{vHost: void}
	}

	// update internal map: apiUUIDToGatewayToVhosts
	if uuid == "" {
		// may be deployed with API-CTL
		logger.LoggerXds.Debug("No UUID defined, do not update vhosts internal maps with UUIDs")
		return
	}
	logger.LoggerXds.Debugf("Updating Vhost internal map of API with UUID \"%v\" as %v.", uuid, vHost)
	var envToVhostMap map[string]string
	if existingMap, ok := apiUUIDToGatewayToVhosts[uuid]; ok {
		logger.LoggerXds.Debugf("API with UUID \"%v\" already exist in vhosts internal map.", uuid)
		envToVhostMap = existingMap
	} else {
		logger.LoggerXds.Debugf("API with UUID \"%v\" not exist in vhosts internal map and create new entry.",
			uuid)
		envToVhostMap = make(map[string]string)
	}

	// if a vhost is already exists it is replaced
	// only one vhost is supported for environment
	// this map is only used for un-deploying APIs form APIM
	for _, env := range gwEnvs {
		envToVhostMap[env] = vHost
	}
	apiUUIDToGatewayToVhosts[uuid] = envToVhostMap
}
