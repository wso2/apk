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

package dto

type APIDefinitionValidationResponse struct {
	IsValid        bool     `json:"isValid"`                // true if valid
	Content        string   `json:"content,omitempty"`      // Original content
	JSONContent    string   `json:"jsonContent,omitempty"`  // JSON representation
	ProtoContent   []byte   `json:"protoContent,omitempty"` // Proto file content
	Protocol       string   `json:"protocol,omitempty"`     // Protocol type (e.g., HTTP/GRPC)
	OpenAPIVersion string   `json:"openAPIVersion"`         // OpenAPI version (e.g., 3.0.1)
	Name           string   `json:"name"`                   // api name
	Version        string   `json:"version"`                // api version
	Context        string   `json:"context"`                // api context path
	Description    string   `json:"description"`            // api description
	Endpoints      []string `json:"endpoints"`              // List of endpoint URLs
	IsInit         bool     `json:"isInit"`                 // Init status
}
