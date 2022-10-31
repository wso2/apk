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
	"os/exec"
	"strings"

	k8sUtils "github.com/BLasan/APKCTL-Demo/CTL/k8s"
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
)

const gatewayAPICRDsYaml = "https://github.com/envoyproxy/gateway/releases/download/v0.2.0-rc1/gatewayapi-crds.yaml"
const envoyGatewayInstallYaml = "https://github.com/envoyproxy/gateway/releases/download/v0.2.0-rc1/install.yaml"
const gatewayClassYaml = "https://raw.githubusercontent.com/envoyproxy/gateway/v0.2.0-rc1/examples/kubernetes/gatewayclass.yaml"
const gatewayYaml = "https://raw.githubusercontent.com//envoyproxy/gateway/v0.2.0-rc1/examples/kubernetes/gateway.yaml"

func InstallPlatform(profile, namespace string, helmVersion int) {

	if helmVersion == 2 {
		utils.HandleErrorAndExit("Please use Helm version 3 as version 2 is not supported yet", nil)
	}

	// Get default namespace if namespace is not speicified
	if namespace == "" {
		namespace = utils.GetNamespace()
	}

	// If profile is not specified using the --profile flag, install all K8s components
	// (i.e. Components of Control Plane and Data Plane)
	if profile == "" {
		fmt.Printf(
			"Installing APK Platform...\nConnected Cluster Name: %s\nContext: %s\nNamespace: %s\n\n",
			utils.GetClusterName(),
			utils.GetContext(),
			namespace,
		)

		// Install Envoy Gateway within the Data Plane profile
		// NOTE: Need to remove this function when Envoy Gateway is embedded into the helm charts
		installEnvoyGateway()

		// Execute helm commands to add and install Helm Chart for Data Plane and Control Plane profiles
		helmArgs := []string{
			utils.HelmSetFlag,
			"wso2.apk.cp.ipk.enabled=false", // to disable IPK temporarily
		}
		executeHelmCommand(namespace, helmArgs)

	} else if profile == "dp" {
		// TODO: Re-do when the namespace used by the envoy gateway can be overriden
		// Install components in K8s default cluster with default namespace
		fmt.Printf(
			"Installing Data Plane Components...\nConnected Cluster Name: %s\nContext: %s\nNamespace: %s\n\n",
			utils.GetClusterName(),
			utils.GetContext(),
			utils.GetNamespace(),
		)

		// Install Envoy Gateway within the Data Plane profile
		// NOTE: Need to remove this function when Envoy Gateway is embedded into the helm charts
		installEnvoyGateway()

		// Execute helm commands to add and install Helm Chart for Data Plane and Control Plane profiles
		helmArgs := []string{
			utils.HelmSetFlag,
			"wso2.apk.cp.enabled=false",
			utils.HelmSetFlag,
			"wso2.apk.cp.postgresql.enabled=false",
			utils.HelmSetFlag,
			"wso2.apk.cp.ipk.enabled=false", // to disable IPK temporarily
		}
		executeHelmCommand(namespace, helmArgs)

	} else if profile == "cp" {
		if namespace == "" {
			namespace = utils.GetNamespace()
		}

		fmt.Printf(
			"Installing Control Plane Components...\nConnected Cluster Name: %s\nContext: %s\nNamespace: %s\n\n",
			utils.GetClusterName(),
			utils.GetContext(),
			namespace,
		)

		// Execute helm commands to add and install Helm Chart
		helmArgs := []string{
			utils.HelmSetFlag,
			"wso2.apk.dp.enabled=false",
		}
		executeHelmCommand(namespace, helmArgs)
	}
	fmt.Println("\nAll Done! We have configured APK to help you build and manage APIs with ease.")
}

// Function to install and setup Envoy Gateway
func installEnvoyGateway() {
	// Install the Gateway API CRDs
	if err := k8sUtils.ExecuteCommand(
		k8sUtils.Kubectl,
		k8sUtils.K8sApply,
		k8sUtils.FilenameFlag,
		gatewayAPICRDsYaml,
	); err != nil {
		utils.HandleErrorAndExit("Error installing Gateway API CRDs", err)
	}
	// Run Envoy Gateway
	if err := k8sUtils.ExecuteCommand(
		k8sUtils.Kubectl,
		k8sUtils.K8sApply,
		k8sUtils.FilenameFlag,
		envoyGatewayInstallYaml,
	); err != nil {
		utils.HandleErrorAndExit("Error installing Envoy Gateway", err)
	}

	// Check pod status of `gateway-api-admission-server` to determine if it is in Running state
	for {
		podStatus := getPodStatus()
		if strings.Trim(podStatus, "\n") == "Running" {
			break
		}
	}

	// Create the GatewayClass
	if err := k8sUtils.ExecuteCommand(
		k8sUtils.Kubectl,
		k8sUtils.K8sApply,
		k8sUtils.FilenameFlag,
		gatewayClassYaml,
	); err != nil {
		utils.HandleErrorAndExit("Error creating the Gateway Class", err)
	}
	// Create the Gateway
	if err := k8sUtils.ExecuteCommand(
		k8sUtils.Kubectl,
		k8sUtils.K8sApply,
		k8sUtils.FilenameFlag,
		gatewayYaml,
	); err != nil {
		utils.HandleErrorAndExit("Error creating the Gateway", err)
	}
}

// Function to install Helm Chart
// Function to install and setup Data Plane and/or Control Plane profile with a single helm install command
func executeHelmCommand(namespace string, helmArgs []string) {
	// Add bitnami
	if err := k8sUtils.ExecuteCommand(
		utils.Helm,
		utils.HelmRepo,
		utils.HelmAdd,
		"bitnami",
		"https://charts.bitnami.com/bitnami",
	); err != nil {
		utils.HandleErrorAndExit("Error encountered while adding bitnami Helm chart", err)
	}

	// Add chartmuseum
	// if err := k8sUtils.ExecuteCommand(
	// 	utils.Helm,
	// 	utils.HelmRepo,
	// 	utils.HelmAdd,
	// 	"chartmuseum",
	// 	"http://localhost:8080",
	// ); err != nil {
	// 	utils.HandleErrorAndExit("Error encountered while adding chartmuseum Helm chart", err)
	// }

	// Change directory to APK Helm home
	utils.ChangeDirectory("../helm")

	// Download the dependent charts
	if err := k8sUtils.ExecuteCommand(
		utils.Helm,
		utils.HelmDependency,
		utils.HelmBuild,
	); err != nil {
		utils.HandleErrorAndExit("Error encountered while executing the Helm dependency build command", err)
	}

	helmCmd := []string{
		utils.HelmInstall,
		utils.APKHelmChartReleaseName,
		".",
		utils.HelmNamespaceFlag,
		namespace,
		utils.HelmCreateNamespaceFlag,
	}
	helmCmd = append(helmCmd, helmArgs...)

	// Install the APK components
	if err := k8sUtils.ExecuteCommand(
		utils.Helm,
		helmCmd...
	); err != nil {
		utils.HandleErrorAndExit("Error encountered while installing Helm chart", err)
	}

	// Change directory to APKCTL home
	utils.ChangeDirectory("../CTL")
}

func getPodStatus() string {
	podStatus, err := exec.Command(
		"bash", "-c",
		"kubectl get pods -n gateway-system --no-headers | awk '{if ($1 ~ \"gateway-api-admission-server-\") print $3}'",
	).Output()
	if err != nil {
		utils.HandleErrorAndExit("Error while checking the pod status of a pod that is required for the Envoy Gateway", err)
	}
	return string(podStatus)
}
