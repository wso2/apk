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

package utils

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	k8sUtils "github.com/BLasan/APKCTL-Demo/CTL/k8s"
)

func GetAPKCTLHomeDir() (string, error) {
	return os.Getwd()
}

func GetAPKHelmHomeDir() (string, error) {
	pwd, err := os.Getwd()
	return pwd + "/helm", err
}

// Retrieve name of the connected cluster
func GetClusterName() string {
	clusterName, _ := k8sUtils.GetCommandOutput(
		k8sUtils.Kubectl,
		k8sUtils.K8sConfig,
		k8sUtils.K8sView,
		k8sUtils.MinifyFlag,
		k8sUtils.OutputFormatFlag,
		"jsonpath='{.clusters[].name}'",
	)
	clusterName = strings.ReplaceAll(clusterName, "'", "")
	return clusterName
}

// Retrieve name of the current context
func GetContext() string {
	context, _ := k8sUtils.GetCommandOutput(
		k8sUtils.Kubectl,
		k8sUtils.K8sConfig,
		k8sUtils.K8sView,
		k8sUtils.MinifyFlag,
		k8sUtils.OutputFormatFlag,
		"jsonpath='{.contexts[].name}'",
	)
	context = strings.ReplaceAll(context, "'", "")
	return context
}

// Retrieve the namespace
func GetNamespace() string {
	namespace, _ := k8sUtils.GetCommandOutput(
		k8sUtils.Kubectl,
		k8sUtils.K8sConfig,
		k8sUtils.K8sView,
		k8sUtils.MinifyFlag,
		k8sUtils.OutputFormatFlag,
		"jsonpath='{.contexts[].context.namespace}'",
	)
	namespace = strings.ReplaceAll(namespace, "'", "")
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return namespace
}

func FindPathParam(array []string) string {
	pathPrefix := []string{}
	// low := 0
	// high := len(array) - 1

	// for low <= high {
	// 	mid := (low + high) / 2
	// 	if array[mid] < param {
	// 		low = mid + 1
	// 	} else {
	// 		high = mid - 1
	// 	}
	// }

	// if low == len(array) || strings.ContainsAny(array[low], param) {
	// 	return -1
	// }

	// return low

	for index, item := range array {
		if index == 0 {
			item = "/" + item
		}
		pathPrefix = append(pathPrefix, item)
		if strings.ContainsAny(item, "{}") {
			return strings.Join(pathPrefix, "/")
		}
	}

	return ""
}

func CreateConfigMapFromTemplate(configmap ConfigMap, filepath string) {
	t, err := template.New("APIConfigMap").Parse(configMapTemplate)

	if err != nil {
		HandleErrorAndExit("Error Parsing the template", err)
	}

	f, err := os.Create(path.Join(filepath, "APIConfigMap.yaml"))

	if err != nil {
		HandleErrorAndExit("Error creating configmap", err)
	}

	defer f.Close()

	// var out bytes.Buffer

	templ := template.Must(t, err)
	// err = templ.Execute(&out, configmap)
	err = templ.Execute(f, configmap)

	if err != nil {
		HandleErrorAndExit("Error executing the template", err)
	}

}

// changeDirectory will change the directory to the repoPath specified
func ChangeDirectory(repoPath string) {
	err := os.Chdir(repoPath)
	if err != nil {
		HandleErrorAndExit("Error while changing the current directory to "+repoPath, err)
	}
	pwd, _ := os.Getwd()
	fmt.Println("Changed the current directory to " + pwd)
}
