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

package cmd

import (
	"github.com/BLasan/APKCTL-Demo/CTL/impl"
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
	"github.com/spf13/cobra"
)

var profile string
var namespace string
var helmVersion int

const InstallPlatformCmdLiteral = "install platform"
const InstallPlatformCmdShortDesc = "Install APIM Control Plane component(s) and Data Plane component(s)"
const InstallPlatformCmdLongDesc = `Install APIM Control Plane component(s) and Data Plane component(s)`
const InstallPlatformCmdExamples = utils.ProjectName + ` ` + InstallPlatformCmdLiteral

// InstallPlatformCmd represents the APKCTL platform installation command
var InstallPlatformCmd = &cobra.Command{
	Use:     InstallPlatformCmdLiteral,
	Short:   InstallPlatformCmdShortDesc,
	Long:    InstallPlatformCmdLongDesc,
	Example: InstallPlatformCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		handleInstallPlatform()
	},
}

func handleInstallPlatform() {
	impl.InstallPlatform(profile, namespace, helmVersion)
}

func init() {
	InstallPlatformCmd.Flags().StringVar(&profile, "profile", "", "Name of profile i.e. CP (Control Plane) or DP (Data Plane)")
	InstallPlatformCmd.Flags().StringVarP(&namespace, "namespace", "n", "", `Namespace for the profile.
		Note that CP and DP can be in two different namespaces or in same namespace`)
	InstallPlatformCmd.Flags().IntVar(&helmVersion, "helm-version", 3, "Helm version to use")
}
