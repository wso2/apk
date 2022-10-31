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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"net/url"

	k8sUtils "github.com/BLasan/APKCTL-Demo/CTL/k8s"
	"github.com/BLasan/APKCTL-Demo/CTL/utils"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/spec"
	"gopkg.in/yaml.v2"
)

var dirPath string
var desFilePath string

func CreateAPI(filePath, namespace, serviceUrl, apiName, version string, isDryRun, applyNetworkPolicy bool) {

	var apiContent []byte
	var err error

	// Checking if path to API definition is provided. If not specified, use the default OpenAPI definition
	if filePath == "" {
		apiContent = []byte(utils.DefaultSwaggerFile)
	} else {
		apiContent, err = ioutil.ReadFile(filePath)
	}

	if err != nil {
		utils.HandleErrorAndExit("Error encountered while reading API definition file", err)
	}

	definitionJsn, err := utils.ToJSON(apiContent)

	if err != nil {
		utils.HandleErrorAndExit("Error converting API definition file to json", err)
	}

	definitionVersion := utils.FindAPIDefinitionVersion(definitionJsn)

	if definitionVersion == utils.Swagger2 {

		// API definition is a Swagger file
		var swaggerSpec spec.Swagger
		err = json.Unmarshal(definitionJsn, &swaggerSpec)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling swagger", err)
		}

		// If version is provided as a flag, modify the version in the Swagger file accordingly
		if version != "" {
			swaggerSpec.Info.Version = version
		}

		// Modify the API name in the Swagger file with the name from the create api command
		if apiName != "" {
			swaggerSpec.Info.Title = apiName
		}

		// If service URL is provided as a flag, modify the backend URL in the Swagger file accordingly
		if serviceUrl != "" {
			swaggerSpec.Host = serviceUrl
		}

		createAndDeploySwaggerAPI(swaggerSpec, filePath, namespace, serviceUrl, isDryRun, applyNetworkPolicy)

	} else if definitionVersion == utils.OpenAPI3 {

		// API definition is an OpenAPI Definition file
		var openAPISpec openapi3.T
		err = json.Unmarshal(definitionJsn, &openAPISpec)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling OpenAPI Definition", err)
		}

		// If version is provided as a flag, modify the version in the OpenAPI definition accordingly
		if version != "" {
			openAPISpec.Info.Version = version
		}

		// Modify the API name in the OpenAPI definition with the name from the create api command
		if apiName != "" {
			openAPISpec.Info.Title = apiName
		}

		// If service URL is provided as a flag, modify the backend URL in the OpenAPI definition accordingly
		if serviceUrl != "" {
			openAPISpec.Servers[0].URL = serviceUrl
		}

		createAndDeployOpenAPI(openAPISpec, filePath, namespace, serviceUrl, apiName, isDryRun, applyNetworkPolicy)

	} else {
		utils.HandleErrorAndExit("Error resolving API definition. Provided file kind is not supported or not acceptable.", nil)
	}
}

func createAndDeploySwaggerAPI(swaggerSpec spec.Swagger, filePath, namespace, serviceUrl string, isDryRun, applyNetworkPolicy bool) {
	httpRoute := utils.HTTPRouteConfig{}
	var parentRef utils.ParentRef

	httpRoute.ApiVersion = utils.HttpRouteApiVersion
	httpRoute.Kind = utils.HttpRouteKind
	httpRoute.HttpRouteSpec.HostNames = append(httpRoute.HttpRouteSpec.HostNames, "www.apk.com")
	parentRef.Name = "eg"
	httpRoute.HttpRouteSpec.ParentRefs = append(httpRoute.HttpRouteSpec.ParentRefs, parentRef)
	httpRoute.MetaData.Name = swaggerSpec.Info.Title + "-" + swaggerSpec.Info.Version
	// httpRoute.MetaData.Namespace = namespace

	labels := make(map[string]string)
	labels["version"] = swaggerSpec.Info.Version
	httpRoute.MetaData.Labels = labels

	var apiPath utils.Path
	var match utils.Match
	var rule utils.Rule
	var backendRef utils.BackendRef

	// Checking if service URL is provided. If not specified, deduce the service URL using the swagger definition
	if serviceUrl == "" {
		if swaggerSpec.Host != "" {
			urlScheme := ""
			for _, scheme := range swaggerSpec.Schemes {
				if scheme == "https" {
					urlScheme = utils.HttpsURLScheme
					break
				} else if scheme == "http" {
					urlScheme = utils.HttpURLScheme
				} else {
					utils.HandleErrorAndExit("Detected scheme(s) within the swagger definition are not supported", nil)
				}
			}
			serviceUrl = urlScheme + swaggerSpec.Host + swaggerSpec.BasePath
		} else {
			utils.HandleErrorAndExit("Unable to find a valid service URL.", nil)
		}
	}

	parsedURL, err := url.ParseRequestURI(serviceUrl)
	if err != nil {
		utils.HandleErrorAndExit("Error while parsing the service URL.", err)
	}
	basePath := parsedURL.Path

	// If API definition is not specified, provide the wildcard resource as a PathPrefix
	if filePath == "" {
		apiPath.Type = utils.PathPrefix
		apiPath.Value = "/"
		match.Path = apiPath
		rule.Matches = append(rule.Matches, match)
	} else {
		counter := 1

		// path & path item
		for path := range swaggerSpec.Paths.Paths {
			// maximum 8 paths are allowed
			if counter > 8 {
				break
			}

			index := strings.IndexAny(path, "{")
			if index >= 0 {
				path = path[:index-1]
			}

			if strings.Contains(path, "/*") {
				path = strings.ReplaceAll(path, "/*", "")
			}

			path = basePath + path
			if path == "" {
				path = "/"
			}

			// pathArr := strings.Split(path, "/")
			// sort.Strings(pathArr)
			// path = utils.FindPathParam(pathArr)

			apiPath.Type = utils.PathPrefix
			apiPath.Value = path
			match.Path = apiPath

			rule.Matches = append(rule.Matches, match)

			// if pathItem.Post != nil {
			// 	fmt.Println("Description Items: ", pathItem.Post.Description)
			// }

			counter++

		}
	}

	backendRef.Kind = utils.ServiceKind
	backendRef.Name = strings.Split(parsedURL.Host, ".")[0]
	// backendRef.Namespace = serviceUrlArr[1]
	if parsedURL.Port() != "" {
		u32, err := strconv.ParseUint(parsedURL.Port(), 10, 32)
		if err != nil {
			fmt.Println("Endpoint port is not in the expected format.", err)
		}
		backendRef.Port = int(uint32(u32))
	} else {
		backendRef.Port = int(uint32(80))
	}

	rule.BackendRefs = append(rule.BackendRefs, backendRef)
	httpRoute.HttpRouteSpec.Rules = append(httpRoute.HttpRouteSpec.Rules, rule)
	if err != nil {
		utils.HandleErrorAndExit("Error extracting port number", err)
	}

	file, err := yaml.Marshal(&httpRoute)
	if err != nil {
		utils.HandleErrorAndExit("Error marshalling httproute file", err)
	}

	if !isDryRun {
		handleDeploy(file, filePath, namespace, swaggerSpec.Info.Title, swaggerSpec.Info.Version, swaggerSpec, utils.Swagger2, parsedURL, applyNetworkPolicy)
	} else {
		handleDryRun(file, filePath, namespace, swaggerSpec.Info.Title, swaggerSpec.Info.Version, swaggerSpec, utils.Swagger2, parsedURL, applyNetworkPolicy)
	}
}

func createAndDeployOpenAPI(openAPISpec openapi3.T, filePath, namespace, serviceUrl, apiName string, isDryRun, applyNetworkPolicy bool) {
	httpRoute := utils.HTTPRouteConfig{}
	var parentRef utils.ParentRef

	httpRoute.ApiVersion = utils.HttpRouteApiVersion
	httpRoute.Kind = utils.HttpRouteKind
	httpRoute.HttpRouteSpec.HostNames = append(httpRoute.HttpRouteSpec.HostNames, "www.apk.com")
	parentRef.Name = "eg"
	httpRoute.HttpRouteSpec.ParentRefs = append(httpRoute.HttpRouteSpec.ParentRefs, parentRef)
	httpRoute.MetaData.Name = apiName + "-" + openAPISpec.Info.Version

	labels := make(map[string]string)
	labels["version"] = openAPISpec.Info.Version
	httpRoute.MetaData.Labels = labels

	var apiPath utils.Path
	var match utils.Match
	var rule utils.Rule
	var backendRef utils.BackendRef

	// Checking if service URL is provided. If not specified, use the service URLs provided under the OpenAPI definition
	if serviceUrl == "" {
		var serviceUrls []string
		for _, serverEntry := range openAPISpec.Servers {
			serviceUrls = append(serviceUrls, serverEntry.URL)
		}
		// We will use the first URL provided under the servers object
		serviceUrl = serviceUrls[0]
	}

	parsedURL, err := url.ParseRequestURI(serviceUrl)
	if err != nil {
		utils.HandleErrorAndExit("Error while parsing the service URL.", err)
	}
	basePath := parsedURL.Path

	// If API definition is not specified, provide the wildcard resource as a PathPrefix
	if filePath == "" {
		apiPath.Type = utils.PathPrefix
		apiPath.Value = "/"
		match.Path = apiPath
		rule.Matches = append(rule.Matches, match)
	} else {
		counter := 1

		// path & path item
		for path := range openAPISpec.Paths {
			// maximum 8 paths are allowed
			if counter > 8 {
				break
			}

			index := strings.IndexAny(path, "{")
			if index >= 0 {
				path = path[:index-1]
			}

			// remove *
			if strings.Contains(path, "/*") {
				path = strings.ReplaceAll(path, "/*", "")
			}

			path = basePath + path
			if path == "" {
				path = "/"
			}

			apiPath.Type = utils.PathPrefix
			apiPath.Value = path
			match.Path = apiPath

			rule.Matches = append(rule.Matches, match)

			counter++
		}
	}

	backendRef.Kind = utils.ServiceKind
	backendRef.Name = strings.Split(parsedURL.Host, ".")[0]
	if parsedURL.Port() != "" {
		u32, err := strconv.ParseUint(parsedURL.Port(), 10, 32)
		if err != nil {
			fmt.Println("Endpoint port is not in the expected format.", err)
		}
		backendRef.Port = int(uint32(u32))
	} else {
		backendRef.Port = int(uint32(80))
	}

	rule.BackendRefs = append(rule.BackendRefs, backendRef)
	httpRoute.HttpRouteSpec.Rules = append(httpRoute.HttpRouteSpec.Rules, rule)

	file, err := yaml.Marshal(&httpRoute)
	if err != nil {
		utils.HandleErrorAndExit("Error marshalling httproute file.", err)
	}

	version := openAPISpec.Info.Version

	if !isDryRun {
		handleDeploy(file, filePath, namespace, apiName, version, openAPISpec, utils.OpenAPI3, parsedURL, applyNetworkPolicy)
	} else {
		handleDryRun(file, filePath, namespace, apiName, version, openAPISpec, utils.OpenAPI3, parsedURL, applyNetworkPolicy)
	}
}

func handleApplyNetworkPolicy(apiName string, serviceUrl *url.URL, dirPath string) {
	service := strings.Split(serviceUrl.Host, ".")[0]
	namespace := strings.Split(serviceUrl.Host, ".")[1]
	dnsType := strings.Split(serviceUrl.Host, ".")[2]

	if dnsType == "svc" {
		out, err := k8sUtils.GetCommandOutput(k8sUtils.Kubectl, k8sUtils.K8sGet, k8sUtils.K8sService, service, "-n", namespace, "-o", "wide")
		if err != nil {
			utils.HandleErrorAndExit("Error executing K8s command", err)
		} else {
			fmt.Println(out)
			selectors := strings.Fields(strings.SplitAfter(out, "\n")[1])[6]
			stringSlice := strings.Split(selectors, ",")
			labelMap := make(map[string]string)

			for _, element := range stringSlice {
				labelSlice := strings.Split(element, "=")
				labelMap[labelSlice[0]] = labelSlice[1]
			}

			file := createNetworkPolicy(apiName, labelMap)
			nwpFilePath := filepath.Join(dirPath, "NetworkPolicy.yaml")

			err = ioutil.WriteFile(nwpFilePath, file, 0644)

			if err != nil {
				utils.HandleErrorAndExit("Error creating NetworkPolicy file", err)
			}
		}
	}
}

func createNetworkPolicy(apiName string, selectors map[string]string) []byte {
	networkPolicy := utils.NetworkPolicy{}
	var metadata utils.MetaData
	var networkPolicySpec utils.NetworkPolicySpec
	var podSelector utils.PodSelector
	var ingress utils.Ingress
	var from utils.From
	var namespaceSelector utils.NamespaceSelector

	networkPolicy.Kind = "NetworkPolicy"
	networkPolicy.ApiVersion = "networking.k8s.io/v1"

	metadata.Name = strings.ToLower(apiName + "NetworkPolicy")
	networkPolicy.MetaData = metadata

	podSelector.MatchLabels = selectors

	networkPolicySpec.PodSelector = podSelector

	nsSMap := make(map[string]string)
	nsSMap["gw"] = "envoy"

	namespaceSelector.MatchLabels = nsSMap
	from.NamespaceSelector = namespaceSelector
	ingress.From = append(ingress.From, from)
	networkPolicySpec.Ingress = append(networkPolicySpec.Ingress, ingress)

	networkPolicy.NetworkPolicySpec = networkPolicySpec

	file, err := yaml.Marshal(&networkPolicy)
	if err != nil {
		fmt.Println(err)
	}

	return file
}

// Handle API deploy
func handleDeploy(file []byte, swaggerFilePath, namespace, apiName, version string, definition interface{}, swaggerVersion string, serviceUrl *url.URL, applyNetworkPolicy bool) {
	var err error
	apiProjectDirName := apiName + "-" + version
	dirPath, err = os.MkdirTemp("", apiProjectDirName)
	if err != nil {
		utils.HandleErrorAndExit("Error creating the temp directory", err)
	}

	defer os.RemoveAll(dirPath)

	desFilePath = filepath.Join(dirPath, "HTTPRouteConfig.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error creating HTTPRouteConfig file", err)
	}

	// set the file name to get the file extension
	if swaggerFilePath == "" {
		swaggerFilePath = utils.DefaultSwaggerFileName
	}

	if applyNetworkPolicy {
		handleApplyNetworkPolicy(apiName, serviceUrl, dirPath)
	}

	createConfigMap(filepath.Ext(swaggerFilePath), dirPath, namespace, apiName, definition, swaggerVersion, version)
	// utils.CreateConfigMapFromTemplate(configmap, dirPath)

	args := []string{k8sUtils.K8sApply, k8sUtils.FilenameFlag, filepath.Join(dirPath, "")}

	err = k8sUtils.ExecuteCommand(k8sUtils.Kubectl, args...)
	if err != nil {
		utils.HandleErrorAndExit("Error Deploying the API", err)
	}
	os.RemoveAll(dirPath)

	fmt.Println("\nSuccessfully deployed " + apiName + " API into the " + namespace + " namespace")
}

// Handle the `Dry Run` option of create API command
// This will generate an API project based on the provided command and flags
func handleDryRun(file []byte, swaggerFilePath, namespace, apiName, version string, definition interface{}, swaggerVersion string, serviceUrl *url.URL, applyNetworkPolicy bool) {
	var err error
	dirPath, err = utils.GetAPKCTLHomeDir()
	if err != nil {
		utils.HandleErrorAndExit("Error getting apkctl home directory", err)
	}

	apiProjectDirName := apiName + "-" + version
	dirPath = path.Join(dirPath, utils.APIProjectsDir, apiProjectDirName)

	os.MkdirAll(dirPath, os.ModePerm)

	desFilePath = filepath.Join(dirPath, "HTTPRouteConfig.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)

	if err != nil {
		utils.HandleErrorAndExit("Error creating HTTPRouteConfig file", err)
	}

	// set the file name to get the file extension
	if swaggerFilePath == "" {
		swaggerFilePath = utils.DefaultSwaggerFileName
	}

	if applyNetworkPolicy {
		handleApplyNetworkPolicy(apiName, serviceUrl, dirPath)
	}

	createConfigMap(filepath.Ext(swaggerFilePath), dirPath, namespace, apiName, definition, swaggerVersion, version)
	// utils.CreateConfigMapFromTemplate(configmap, dirPath)

	fmt.Println("Successfully created API project with HttpRouteConfig and ConfigMap files!")
	fmt.Println("API project directory: " + utils.APIProjectsDir + apiName + "-" + version)
}

func createConfigMap(ext, dirPath, namespace, apiname string, definition interface{}, swaggerVersion, apiversion string) {
	configmap := utils.ConfigMap{}
	configmap.ApiVersion = "v1"
	configmap.Kind = "ConfigMap"
	configmap.MetaData.Name = apiname + "-" + apiversion

	if namespace != "" {
		configmap.MetaData.Namespace = namespace
	}

	// content := readSwaggerDef(filepath)

	// if content == "" {
	// 	fmt.Println("Empty Swagger")
	// 	// handle error and exit
	// }

	data := make(map[string]string)

	if ext == ".yaml" {
		content, err := yaml.Marshal(definition)
		if err != nil {
			utils.HandleErrorAndExit("Error while Marshalling the YAML ", err)
		}
		if swaggerVersion == utils.Swagger2 {
			data["swagger.yaml"] = string(content)
		} else if swaggerVersion == utils.OpenAPI3 {
			data["openapi.yaml"] = string(content)
		}
	} else if ext == ".json" {
		content, err := json.Marshal(definition)
		if err != nil {
			utils.HandleErrorAndExit("Error while Marshalling the JSON ", err)
		}
		if swaggerVersion == utils.Swagger2 {
			data["swagger.json"] = string(content)
		} else if swaggerVersion == utils.OpenAPI3 {
			data["openapi.json"] = string(content)
		}
	}

	configmap.Data = data

	file, err := yaml.Marshal(&configmap)

	if err != nil {
		utils.HandleErrorAndExit("Error Marshaling", err)
	}

	desFilePath := path.Join(dirPath, "ConfigMap.yaml")

	// directory location can be defined in the apkctl config file
	err = ioutil.WriteFile(desFilePath, file, 0644)

	if err != nil {
		utils.HandleErrorAndExit("Error creating config file", err)
	}
}
