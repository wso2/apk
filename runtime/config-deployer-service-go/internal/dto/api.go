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

// API represents the structure of an API definition
type API struct {
	Name              string        `json:"name" yaml:"name"`                           // api name
	BasePath          string        `json:"basePath" yaml:"basePath"`                   // api base path
	Version           string        `json:"version" yaml:"version"`                     // api version
	Type              string        `json:"type" yaml:"type"`                           // api type (e.g., REST, GraphQL)
	Endpoint          string        `json:"endpoint" yaml:"endpoint"`                   // Endpoint URL
	URITemplates      []URITemplate `json:"uriTemplates" yaml:"uriTemplates"`           // Array of URI templates
	APISecurity       string        `json:"apiSecurity" yaml:"apiSecurity"`             // Security definition
	Scopes            []string      `json:"scopes" yaml:"scopes"`                       // Array of scopes
	GraphQLSchema     string        `json:"graphQLSchema" yaml:"graphQLSchema"`         // GraphQL schema string
	ProtoDefinition   string        `json:"protoDefinition" yaml:"protoDefinition"`     // gRPC proto content
	SwaggerDefinition string        `json:"swaggerDefinition" yaml:"swaggerDefinition"` // Swagger/OpenAPI content
	Environment       string        `json:"environment" yaml:"environment"`             // Deployment environment
}

// URITemplate represents a URI template
type URITemplate struct {
	URITemplate string   `json:"uriTemplate" yaml:"uriTemplate"`
	ResourceURI string   `json:"resourceURI" yaml:"resourceURI"`
	Verb        string   `json:"verb" yaml:"verb"`
	AuthEnabled bool     `json:"authEnabled" yaml:"authEnabled"`
	Scopes      []string `json:"scopes" yaml:"scopes"`
	ID          int      `json:"id" yaml:"id"`
	Endpoint    string   `json:"endpoint" yaml:"endpoint"`
}
