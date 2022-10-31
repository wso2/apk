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

package main

import (
	"os"
	"time"

	cmd "github.com/BLasan/APKCTL-Demo/CTL/cmd"
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
	"github.com/spf13/cobra"
)

var cfgFile string
var verbose bool

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "apkctl",
	Short: "apkctl",
	Long:  "apkctl",
}

func execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	checkVerbose()
}

func main() {
	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	RootCmd.AddCommand(cmd.InstallPlatformCmd)
	RootCmd.AddCommand(cmd.CreateCmd)
	RootCmd.AddCommand(cmd.DeleteCmd)
	RootCmd.AddCommand(cmd.GetCmd)
	RootCmd.AddCommand(cmd.UninstallPlatformCmd)
	RootCmd.AddCommand(cmd.VersionCmd)
	execute()
}

// func init() {

// 	// cobra.OnInitialize(initConfig)
// }

// func initConfig() {
// 	if verbose {
// 		fmt.Println("Verbose Enabled")
// 		utils.EnableVerboseMode()
// 		t := time.Now()
// 		utils.Logf("Executed ImportExportCLI (%s) on %v\n", utils.ProjectName, t.Format(time.RFC1123))
// 	}
// }

func checkVerbose() {
	if verbose {
		utils.EnableVerboseMode()
		t := time.Now()
		utils.Logf("Executed ImportExportCLI (%s) on %v\n", utils.ProjectName, t.Format(time.RFC1123))
	}
}
