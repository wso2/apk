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

var outputFormat string
var allNamespaces bool

const GetAPICmdLiteral = "apis"
const GetAPICmdShortDesc = "Get APIs"
const GetCAPImdLongDesc = `List all the deployed APIs. Returned list of APIs can either be from a specific namespace or from all namespaces based on the provided flags.`
const GetAPICmdExamples = utils.ProjectName + ` ` + GetCmdLiteral + ` ` + GetAPICmdLiteral + ` --namespace=wso2`

// GetApiCmd represents the Get API command
var GetApiCmd = &cobra.Command{
	Use:     GetAPICmdLiteral,
	Short:   GetAPICmdShortDesc,
	Long:    GetCAPImdLongDesc,
	Example: GetAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		handleGetApis()
	},
}

func handleGetApis() {
	impl.GetAPIs(dpNamespace, outputFormat, allNamespaces)
}

func init() {
	GetCmd.AddCommand(GetApiCmd)
	GetApiCmd.Flags().StringVarP(&dpNamespace, "namespace", "n", "default", "Namespace of the API")
	GetApiCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output Format of APIs")
	GetApiCmd.Flags().BoolVar(&allNamespaces, "all-namespaces", false, "Get APIs in all namespaces")
}
