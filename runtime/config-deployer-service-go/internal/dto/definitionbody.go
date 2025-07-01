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

// DefinitionBody is the body of the API definition.
type DefinitionBody struct {
	Definition Definition `json:"definition"` // API definition (OAS/Graphql/WebSocket)
	URL        string     `json:"url"`        // URL of the API definition
	APIType    string     `json:"apiType"`    // Type of API
}

// Definition of the API definition body.
type Definition struct {
	FileName    string `json:"definitionType"` // Name of the API specification file.
	FileContent []byte `json:"definitionBody"` // Content of the API specification.
}
