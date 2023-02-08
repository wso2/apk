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
	"fmt"

	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
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
		if stringutils.StringInSlice(existingEnv, deleteEnvs) {
			toBeDel = append(toBeDel, existingEnv)
		} else {
			toBeKept = append(toBeKept, existingEnv)
		}
	}
	return
}

// GetvHostsIdentifier creates a identifier for vHosts for a API considering prod
// and sand env
func GetvHostsIdentifier(UUID string, envType string) string {
	return fmt.Sprintf("%s-%s", UUID, envType)
}
