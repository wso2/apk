/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"io"
	"regexp"
	"strings"
)

// ProtoFile represents a proto file structure
type ProtoFile struct {
	ApiName     string    `json:"apiName"`
	PackageName string    `json:"packageName"`
	BasePath    string    `json:"basePath"`
	Version     string    `json:"version"`
	Services    []Service `json:"routes"`
}

// Service represents a gRPC service
type Service struct {
	ServiceName string   `json:"serviceName"`
	Methods     []string `json:"methods"`
}

type GRPCUtil struct{}

// GetGRPCAPIFromProtoDefinition processes the proto file content and extracts API information
func (grpcUtil *GRPCUtil) GetGRPCAPIFromProtoDefinition(definition []byte, fileName string) (*dto.API, error) {
	api := &dto.API{}
	protoFile := &ProtoFile{
		PackageName: "packageName",
		BasePath:    "basePath",
		Version:     "version",
		Services:    []Service{},
	}
	var uriTemplates []dto.URITemplate
	if strings.HasSuffix(fileName, ".zip") {
		protoContents, err := extractProtoFilesFromZip(definition)
		if err != nil {
			return nil, err
		}
		for _, protoContent := range protoContents {
			templates, err := processProtoFile(protoContent, protoFile)
			if err != nil {
				return nil, err
			}
			uriTemplates = append(uriTemplates, templates...)
		}
	} else {
		templates, err := processProtoFile(definition, protoFile)
		if err != nil {
			return nil, err
		}
		uriTemplates = templates
	}
	api.BasePath = protoFile.BasePath
	api.ProtoDefinition = string(definition)
	api.Version = protoFile.Version
	api.Name = protoFile.ApiName
	api.URITemplates = uriTemplates

	return api, nil
}

// extractProtoFilesFromZip extracts proto files from a zip archive
func extractProtoFilesFromZip(zipContent []byte) ([][]byte, error) {
	var protoFiles [][]byte
	reader, err := zip.NewReader(bytes.NewReader(zipContent), int64(len(zipContent)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip content: %w", err)
	}
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, ".proto") {
			protoContent, err := readProtoFileBytesFromZip(file)
			if err != nil {
				return nil, err
			}
			protoFiles = append(protoFiles, protoContent)
		}
	}
	return protoFiles, nil
}

// processProtoFile processes a proto file and extracts URI templates
func processProtoFile(definition []byte, protoFile *ProtoFile) ([]dto.URITemplate, error) {
	content := string(definition)
	packageString := getPackageString(content)
	var uriTemplates []dto.URITemplate
	var apiNameBuilder strings.Builder
	if protoFile.ApiName != "" {
		apiNameBuilder.WriteString(protoFile.ApiName)
	}
	if packageString == "" && protoFile.PackageName != "" {
		return nil, fmt.Errorf("package string has not been defined in proto file: %s", protoFile.PackageName)
	}
	packageName := getPackageName(packageString)
	if packageName == "" {
		packageName = packageString
	}
	var services []*Service
	protoFile.PackageName = packageName
	protoFile.Version = getVersion(packageString)
	protoFile.BasePath = getBasePath(packageString)
	serviceBlocks := extractServiceBlocks(content)

	for _, serviceBlock := range serviceBlocks {
		serviceName := getServiceName(serviceBlock)
		methodNames := extractMethodNames(serviceBlock)
		service := &Service{
			ServiceName: serviceName,
			Methods:     methodNames,
		}
		services = append(services, service)
		for _, method := range service.Methods {
			uriTemplate := dto.URITemplate{
				URITemplate: packageName + "." + service.ServiceName,
				Verb:        method,
				AuthEnabled: true,
				Scopes:      []string{},
			}
			uriTemplates = append(uriTemplates, uriTemplate)
		}
	}

	for _, service := range services {
		apiNameBuilder.WriteString(service.ServiceName)
		apiNameBuilder.WriteString("-")
	}

	protoFile.ApiName = apiNameBuilder.String()
	return uriTemplates, nil
}

// readProtoFileBytesFromZip reads bytes from a zip file entry
func readProtoFileBytesFromZip(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", file.Name, err)
	}
	defer reader.Close()
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", file.Name, err)
	}
	return content, nil
}

// getPackageString extracts package string from proto content
func getPackageString(content string) string {
	// package string has the format "package something"
	packagePattern := regexp.MustCompile(`package\s+([\w.]+);`)
	matches := packagePattern.FindStringSubmatch(content)
	if len(matches) > 1 && len(matches[1]) > 0 {
		return matches[1]
	}
	fmt.Println("Package has not been defined in the proto file")
	return ""
}

// getPackageName extracts package name from package string
func getPackageName(packageString string) string {
	namePattern := regexp.MustCompile(`v\d+(\.\d+)*\.(\w+)$`)
	matches := namePattern.FindStringSubmatch(packageString)
	if len(matches) > 2 {
		return matches[2]
	}
	fmt.Println("Package name not found in proto file.")
	return ""
}

// getVersion extracts version from package string
func getVersion(packageString string) string {
	versionPattern := regexp.MustCompile(`v\d+(\.\d+)*`)
	match := versionPattern.FindString(packageString)
	if match == "" {
		fmt.Println("Version not found in proto file")
	}
	return match
}

// getBasePath extracts base path from package string
func getBasePath(packageString string) string {
	basePathPattern := regexp.MustCompile(`^(.*?)v\d`)
	matches := basePathPattern.FindStringSubmatch(packageString)
	if len(matches) > 1 {
		basePath := matches[1]
		if len(basePath) > 0 && basePath[len(basePath)-1] == '.' {
			basePath = basePath[:len(basePath)-1]
		}
		return "/" + basePath
	}
	fmt.Println("Base path not found in proto file")
	return ""
}

// extractServiceBlocks extracts service blocks from proto content
func extractServiceBlocks(text string) []string {
	// Regular expression pattern to match the service blocks
	patternString := `service\s+\w+\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}`
	pattern := regexp.MustCompile(patternString)
	matches := pattern.FindAllString(text, -1)
	return matches
}

// getServiceName extracts service name from service block
func getServiceName(serviceBlock string) string {
	// Regular expression pattern to match "service ServiceName"
	patternString := `service\s+(\w+)`
	pattern := regexp.MustCompile(patternString)
	matches := pattern.FindStringSubmatch(serviceBlock)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractMethodNames extracts method names from service block
func extractMethodNames(serviceBlock string) []string {
	// Regular expression pattern to match "rpc MethodName"
	patternString := `rpc\s+(\w+)`
	pattern := regexp.MustCompile(patternString)
	matches := pattern.FindAllStringSubmatch(serviceBlock, -1)
	var methodNames []string
	for _, match := range matches {
		if len(match) > 1 {
			methodNames = append(methodNames, match[1])
		}
	}
	return methodNames
}
