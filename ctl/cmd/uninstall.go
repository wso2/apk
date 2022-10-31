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

const UninstallPlatformCmdLiteral = "uninstall platform"
const UninstallPlatformCmdShortDesc = "Uninstall APIM Control Plane component(s) and Data Plane component(s)"
const UninstallPlatformCmdLongDesc = `Uninstall APIM Control Plane component(s) and Data Plane component(s)`
const UninstallPlatformCmdExamples = utils.ProjectName + ` ` + UninstallPlatformCmdLiteral

// UninstallPlatformCmd represents the APKCTL platform uninstall command
var UninstallPlatformCmd = &cobra.Command{
	Use:     UninstallPlatformCmdLiteral,
	Short:   UninstallPlatformCmdShortDesc,
	Long:    UninstallPlatformCmdLongDesc,
	Example: UninstallPlatformCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		handleUninstallPlatform()
	},
}

func handleUninstallPlatform() {
	impl.UninstallPlatform()
}

func init() {

}
