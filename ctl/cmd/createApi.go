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
	"errors"

	"github.com/BLasan/APKCTL-Demo/CTL/impl"
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
	"github.com/spf13/cobra"
)

var dpNamespace string
var serviceUrl string
var file string
var isDryRun bool
var applyNetworkPolicy bool

const CreateAPICmdLiteral = "api"
const CreateAPICmdShortDesc = "Create API and Deploy"
const CreateCAPImdLongDesc = `Create an API and Deploy onto the Kubernetes Cluster`
const createAPICmdExamples = utils.ProjectName + ` ` + CreateCmdLiteral + ` ` + CreateAPICmdLiteral + ` petstore -f swagger.yaml --version 1.0.0 --namespace wso2
` + utils.ProjectName + ` ` + CreateCmdLiteral + ` ` + CreateAPICmdLiteral + ` petstore --service-url http://localhost:9443 --namespace wso2`

// CreateApiCmd represents the create API command
var CreateApiCmd = &cobra.Command{
	Use:     CreateAPICmdLiteral,
	Short:   CreateAPICmdShortDesc,
	Long:    CreateCAPImdLongDesc,
	Example: createAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		apiName := args[0]
		if serviceUrl == "" && file == "" {
			utils.HandleErrorAndExit("Either Swagger Definition or Backend Service URL should be provided", errors.New("backend service URL is mandatory"))
		}
		handleCreateApi(apiName)
	},
}

func handleCreateApi(apiName string) {
	impl.CreateAPI(file, dpNamespace, serviceUrl, apiName, version, isDryRun, applyNetworkPolicy)
}

func init() {
	CreateCmd.AddCommand(CreateApiCmd)
	CreateApiCmd.Flags().StringVarP(&dpNamespace, "namespace", "n", "", "Namespace of the API")
	CreateApiCmd.Flags().StringVar(&serviceUrl, "service-url", "", "Backend Service URL")
	CreateApiCmd.Flags().StringVarP(&file, "file", "f", "", "Path to swagger/OAS definition/GraphQL SDL/WSDL")
	CreateApiCmd.Flags().BoolVar(&isDryRun, "dry-run", false, "Generate API Project inclusive of an HTTPRouteConfig and a ConfigMap")
	CreateApiCmd.Flags().StringVarP(&version, "version", "", "", "Version of the API")
	CreateApiCmd.Flags().BoolVar(&applyNetworkPolicy, "restrict-service-access", false, "Create network policies to restrict access to backend")
}
