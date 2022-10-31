/*
 * Copyright (c) 2022, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package impl

import (
	"fmt"

	k8sUtils "github.com/BLasan/APKCTL-Demo/CTL/k8s"
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
)

func DeleteAPI(namespace, apiName, version string) {
	var errMsg string
	resourceName := apiName + "-" + version
	resourceHttpRoute := k8sUtils.K8sHttpRoute + "/" + resourceName
	resourceConfigMap := k8sUtils.K8sConfigMap + "/" + resourceName

	// Execute kubernetes command to delete API
	if deleteApiErr := k8sUtils.ExecuteCommand(k8sUtils.Kubectl, k8sUtils.K8sDelete, resourceHttpRoute, resourceConfigMap, "-n", namespace); deleteApiErr != nil {
		if namespace != "" {
			errMsg = fmt.Sprintf("\nCould not find \"%s\" API with version \"%s\" in the \"%s\" namespace\n",
				apiName, version, namespace)
		} else {
			errMsg = fmt.Sprintf("\nCould not find \"%s\" API with version \"%s\"\n", apiName, version)
		}
		fmt.Println(errMsg)
		utils.HandleErrorAndExit("Error executing K8s command ", nil)
	}

	fmt.Println("\nSuccessfully deleted " + apiName + " API from " + namespace + " namespace")
}
