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

func GetAPIs(namespace, output string, allNamespaces bool) {
	if output == "" {
		output = utils.K8sOutputWithCustomColumns
	}

	if allNamespaces {
		out, err := k8sUtils.GetCommandOutput(k8sUtils.Kubectl, k8sUtils.K8sGet, k8sUtils.K8sHttpRoute, "-o", output)
		if err != nil {
			utils.HandleErrorAndExit("Error executing K8s command", err)
		} else {
			fmt.Println(out)
		}
	} else {
		out, err := k8sUtils.GetCommandOutput(k8sUtils.Kubectl, k8sUtils.K8sGet, k8sUtils.K8sHttpRoute, "-n", namespace, "-o", output)
		if err != nil {
			utils.HandleErrorAndExit("Error executing K8s command", err)
		} else {
			fmt.Println(out)
		}
	}
}
