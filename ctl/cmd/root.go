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

// import (
// 	"fmt"
// 	"os"

// 	"github.com/BLasan/APKCTL-Demo/utils"
// 	"github.com/spf13/cobra"
// )

// var verbose bool
// var cfgFile string
// var insecure bool
// var cmdPassword string
// var CmdUsername string
// var CmdExportEnvironment string
// var CmdResourceTenantDomain string
// var CmdForceStartFromBegin bool

// // RootCmd related info
// const rootCmdShortDesc = "CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator"
// const rootCmdLongDesc = utils.ProjectName + ` is a Command Line Tool for Importing and Exporting APIs and Applications between different environments of WSO2 API Manager
//  (Dev, Production, Staging, QA etc.) and Managing WSO2 Micro Integrator`

// // This represents the base command when called without any subcommands
// var RootCmd = &cobra.Command{
// 	Use: utils.ProjectName,
// 	Args: func(cmd *cobra.Command, args []string) error {
// 		fmt.Print("abc")
// 		return cobra.ArbitraryArgs(cmd, args)
// 		//  if isK8sEnabled() {
// 		// 	 return cobra.ArbitraryArgs(cmd, args)
// 		//  } else {
// 		// 	 return cobra.NoArgs(cmd, args)
// 		//  }
// 	},
// 	Short: rootCmdShortDesc,
// 	Long:  rootCmdLongDesc,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		//  if isK8sEnabled() {
// 		// 	 ExecuteKubernetes(args...)
// 		//  } else {
// 		// 	 cmd.Help()
// 		//  }
// 	},
// }

// // Execute adds all child commands to the root command sets flags appropriately.
// // This is called by main.main(). It only needs to happen once to the rootCmd.
// func Execute() {
// 	if err := RootCmd.Execute(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(-1)
// 	}
// }

// // init using Cobra
// func init() {
// 	// createConfigFiles()

// 	// cobra.OnInitialize(initConfig)

// 	cobra.EnableCommandSorting = false
// 	RootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose mode")
// 	RootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "k", false,
// 		"Allow connections to SSL endpoints without certs")
// 	//RootCmd.PersistentFlags().StringP("author", "a", "", "WSO2")

// 	//viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))

// 	// Here you will define your flags and configuration settings.
// 	// Cobra supports Persistent Flags, which, if defined here,
// 	// will be global for your application.

// 	// Cobra also supports local flags, which will only run
// 	// when this action is called directly.

// 	// Init ConfigVars
// 	// err := utils.SetConfigVars(utils.MainConfigFilePath)
// 	// if err != nil {
// 	// 	utils.HandleErrorAndExit("Error reading "+utils.MainConfigFilePath+".", err)
// 	// }
// }

// // // createConfigFiles() creates the ConfigDir and necessary ConfigFiles inside the user's $HOME directory
// // func createConfigFiles() {
// // 	err := utils.CreateDirIfNotExist(utils.ConfigDirPath)
// // 	if err != nil {
// // 		utils.HandleErrorAndExit("Error creating config directory: "+utils.ConfigDirPath, err)
// // 	}

// // 	err = utils.CreateDirIfNotExist(utils.DefaultExportDirPath)
// // 	if err != nil {
// // 		utils.HandleErrorAndExit("Error creating config directory: "+utils.DefaultExportDirPath, err)
// // 	}

// // 	utils.CreateDirIfNotExist(filepath.Join(utils.DefaultExportDirPath, utils.ExportedApisDirName))
// // 	utils.CreateDirIfNotExist(filepath.Join(utils.DefaultExportDirPath, utils.ExportedApiProductsDirName))
// // 	utils.CreateDirIfNotExist(filepath.Join(utils.DefaultExportDirPath, utils.ExportedAppsDirName))
// // 	utils.CreateDirIfNotExist(filepath.Join(utils.DefaultExportDirPath, utils.ExportedMigrationArtifactsDirName))

// // 	utils.CreateDirIfNotExist(utils.DefaultCertDirPath)

// // 	if !utils.IsFileExist(utils.MainConfigFilePath) {
// // 		var mainConfig = new(utils.MainConfig)
// // 		mainConfig.Config = utils.Config{HttpRequestTimeout: utils.DefaultHttpRequestTimeout,
// // 			ExportDirectory:      utils.DefaultExportDirPath,
// // 			KubernetesMode:       "Default Mode",
// // 			TokenType:            utils.DefaultTokenType,
// // 			VCSDeletionEnabled:   false,
// // 			VCSConfigFilePath:    "",
// // 			TLSRenegotiationMode: utils.TLSRenegotiationNever}

// // 		utils.WriteConfigFile(mainConfig, utils.MainConfigFilePath)
// // 	}

// // 	if !utils.IsFileExist(utils.SampleMainConfigFilePath) {
// // 		sampleConfig, _ := box.Get("/sample/sample_config.yaml")
// // 		err = ioutil.WriteFile(utils.SampleMainConfigFilePath, sampleConfig, os.ModePerm)
// // 		if err != nil {
// // 			utils.HandleErrorAndExit("Error creating default api spec file", err)
// // 		}
// // 	}

// // 	err = utils.CreateDirIfNotExist(utils.LocalCredentialsDirectoryPath)
// // 	if err != nil {
// // 		utils.HandleErrorAndExit("Error creating local directory: "+utils.LocalCredentialsDirectoryName, err)
// // 	}

// // 	if !utils.IsFileExist(utils.EnvKeysAllFilePath) {
// // 		os.Create(utils.EnvKeysAllFilePath)
// // 	}
// // }

// // initConfig reads in config file and ENV variables if set.
// // func initConfig() {
// // 	if verbose {
// // 		utils.EnableVerboseMode()
// // 		t := time.Now()
// // 		utils.Logf("Executed ImportExportCLI (%s) on %v\n", utils.ProjectName, t.Format(time.RFC1123))
// // 	}

// // 	utils.Logln(utils.LogPrefixInfo+"Insecure:", insecure)
// // 	if insecure {
// // 		utils.Insecure = true
// // 	}

// // 	/*
// // 		 if cfgFile != "" { // enable ability to specify config file via flag
// // 			 viper.SetConfigFile(cfgFile)
// // 		 }

// // 		 viper.SetConfigName(".wso2apim-cli") // name of config file (without extension)
// // 		 viper.AddConfigPath("$HOME")         // adding home directory as first search path
// // 		 viper.AutomaticEnv()                 // read in environment variables that match

// // 		 // If a config file is found, read it in.
// // 		 if err := viper.ReadInConfig(); err == nil {
// // 			 fmt.Println("Using config file:", viper.ConfigFileUsed())
// // 		 }
// // 	*/
// // }

// //  //disable flags when the mode set to kubernetes
// //  func isK8sEnabled() bool {
// // 	 //Get config to check mode
// // 	 configVars := utils.GetMainConfigFromFileSilently(utils.MainConfigFilePath)
// // 	 if configVars != nil && configVars.Config.KubernetesMode {
// // 		 return true
// // 	 } else {
// // 		 return false
// // 	 }
// //  }

// //  //execute kubernetes commands
// //  func ExecuteKubernetes(arg ...string) {
// // 	 cmd := exec.Command(
// // 		 k8sUtils.Kubectl,
// // 		 arg...,
// // 	 )
// // 	 var errBuf, outBuf bytes.Buffer
// // 	 cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
// // 	 cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
// // 	 err := cmd.Run()
// // 	 if err != nil {
// // 		 utils.HandleErrorAndExit("Error executing kubernetes commands ", err)
// // 	 }
// //  }
