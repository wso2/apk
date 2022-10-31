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

const DeleteCmdLiteral = "delete"
const DeleteCmdShortDesc = "Delete API"
const DeleteCmdLongDesc = `Delete API from Kubernetes Cluster`
const DeleteCmdExamples = utils.ProjectName + ` ` + DeleteCmdLiteral + ` ` + DeleteAPICmdLiteral + ` petstore --version 1.0.0 -n wso2
	
NOTE: The flag --version (-v) is mandatory.
You can optionally provide the --namespace (-n) flag to specify the namespace of the deployed API that you wish to delete.

The API to be deleted is identified using the API name and version.
Optionally, you can specify the namespace that the API resides in.
If the API does not exist, an error is thrown.`

// DeleteCmd represents the Delete command
var DeleteCmd = &cobra.Command{
	Use:     DeleteCmdLiteral,
	Short:   DeleteCmdShortDesc,
	Long:    DeleteCmdLongDesc,
	Example: DeleteCmdExamples,
}
