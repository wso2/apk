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
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
	"github.com/spf13/cobra"
)

const GetCmdLiteral = "get"
const GetCmdShortDesc = "Get APIs"
const GetCmdLongDesc = `List all the deployed APIs. Returned list of APIs can either be from a specific namespace or from all namespaces based on the provided flags.`
const GetCmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPICmdLiteral + ` --namespace=wso2

	NOTE: The following flags are considered as optional
	--output, -o			Output format
	--namespace, -n			Namespace of the Data Plane
	--all-namespaces, -A	List the APIs across all namespaces.`

// GetCmd represents the Get command
var GetCmd = &cobra.Command{
	Use:     GetCmdLiteral,
	Short:   GetCmdShortDesc,
	Long:    GetCmdLongDesc,
	Example: GetCmdExamples,
}

// func init() {
// 	RootCmd.AddCommand(GetCmd)
// }
