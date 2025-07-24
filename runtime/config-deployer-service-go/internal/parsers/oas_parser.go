/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
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

package parsers

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"os"
	"sort"
	"strings"
)

type OAS3Parser struct{}

// ParseOpenAPI3 parses OpenAPI 3.x specification from file
func (oasParser *OAS3Parser) ParseOpenAPI3(filePath string) (*openapi3.T, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading OpenAPI file: %s, %w", filePath, err)
	}
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	doc, err := loader.LoadFromData(content)
	if err != nil {
		return nil, fmt.Errorf("error parsing OpenAPI file: %s, %w", filePath, err)
	}

	// Resolve references (equivalent to setResolve(true))
	err = loader.ResolveRefsIn(doc, nil)
	if err != nil {
		return nil, fmt.Errorf("error resolving OpenAPI file: %s, %w", filePath, err)
	}

	return doc, nil
}

// ParseSwagger2 parses Swagger 2.0 specification from file
func (oasParser *OAS3Parser) ParseSwagger2(filePath string) (*spec.Swagger, error) {
	doc, err := loads.Spec(filePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing Swagger file: %s, %w", filePath, err)
	}

	// Expand the spec (equivalent to read with resolve=true)
	expandedDoc, err := doc.Expanded()
	if err != nil {
		return nil, fmt.Errorf("error resolving Swagger file: %s, %w", filePath, err)
	}

	return expandedDoc.Spec(), nil
}

// RemoveUnsupportedBlocksFromResources removes unsupported blocks from the API definition
func (oasParser *OAS3Parser) RemoveUnsupportedBlocksFromResources(jsonString string) (string, error) {
	var jsonObject map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonObject); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	definitionUpdated := false
	if pathsInterface, exists := jsonObject[constants.OPENAPI_RESOURCE_KEY]; exists {
		if paths, ok := pathsInterface.(map[string]interface{}); ok {
			// Remove unsupported blocks recursively
			for _, unsupportedBlockKey := range constants.UnsupportedResourceBlocks {
				updated := oasParser.removeBlocksRecursivelyFromJSONObject(unsupportedBlockKey, paths, false)
				definitionUpdated = definitionUpdated || updated
			}
		}
	}

	if definitionUpdated {
		jsonBytes, err := json.MarshalIndent(jsonObject, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return string(jsonBytes), nil
	}

	return "", nil // Return empty string if no changes were made
}

// GetAPIFromDefinition parses OpenAPI 3.x definition and returns API object
func (oasParser *OAS3Parser) GetAPIFromDefinition(content string) (*dto.API, error) {
	openAPI, err := oasParser.getOpenAPI(content)
	if err != nil {
		return nil, err
	}
	servers := openAPI.Servers
	api := &dto.API{}
	info := openAPI.Info
	if info != nil {
		api.Name = info.Title
		api.Version = info.Version
	}
	if servers != nil && len(servers) > 0 {
		api.Endpoint = servers[0].URL
	}
	uriTemplates, err := oasParser.getURITemplates(openAPI)
	if err != nil {
		return nil, err
	}
	api.URITemplates = uriTemplates
	return api, nil
}

// getOpenAPI parses OpenAPI definition string and returns OpenAPI object
func (oasParser *OAS3Parser) getOpenAPI(oasDefinition string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	ctx := util.GetOpenAPIValidationContext(loader.Context)

	doc, err := loader.LoadFromData([]byte(oasDefinition))
	if err != nil {
		return nil, fmt.Errorf("errors found when parsing OAS definition: %w", err)
	}
	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("errors found when validating OAS definition: %w", err)
	}
	return doc, nil
}

// getURITemplates returns URI templates according to the given OpenAPI model
func (oasParser *OAS3Parser) getURITemplates(openAPI *openapi3.T) ([]dto.URITemplate, error) {
	var urlTemplates []dto.URITemplate
	scopes, err := oasParser.getScopes(openAPI)
	if err != nil {
		return nil, err
	}
	for pathKey, pathItem := range openAPI.Paths.Map() {
		operations := map[string]*openapi3.Operation{
			"GET":     pathItem.Get,
			"PUT":     pathItem.Put,
			"POST":    pathItem.Post,
			"DELETE":  pathItem.Delete,
			"PATCH":   pathItem.Patch,
			"HEAD":    pathItem.Head,
			"OPTIONS": pathItem.Options,
		}
		for httpMethod, operation := range operations {
			if operation != nil && (constants.SupportedMethods[strings.ToLower(httpMethod)] ||
				constants.GraphQLSupportedMethods[strings.ToUpper(httpMethod)]) {
				template := &dto.URITemplate{
					URITemplate: pathKey,
					Verb:        strings.ToUpper(httpMethod),
					AuthEnabled: true,
					Scopes:      []string{},
				}
				opScopes := oasParser.getScopeOfOperations(constants.OPENAPI_SECURITY_SCHEMA_KEY, operation)
				if len(opScopes) > 0 {
					template, err = setScopesToTemplate(template, opScopes, scopes)
				} else {
					oauth2Scopes := oasParser.getScopeOfOperations(constants.OAUTH2_SECURITY_SCHEMA_KEY, operation)
					if len(oauth2Scopes) > 0 {
						template, err = setScopesToTemplate(template, oauth2Scopes, scopes)
					}
				}
				if operation.Servers != nil && len(*operation.Servers) > 0 {
					template.Endpoint = (*operation.Servers)[0].URL
				}
				urlTemplates = append(urlTemplates, *template)
			}
		}
	}
	return urlTemplates, nil
}

// getScopes extracts OAuth scopes from OpenAPI definition
func (oasParser *OAS3Parser) getScopes(openAPI *openapi3.T) ([]string, error) {
	var scopeSet = make(map[string]bool)

	if openAPI.Components != nil && openAPI.Components.SecuritySchemes != nil {
		// Check default security scheme
		if securityScheme, exists := openAPI.Components.SecuritySchemes[constants.OPENAPI_SECURITY_SCHEMA_KEY]; exists {
			if securityScheme.Value != nil && securityScheme.Value.Flows != nil {
				if securityScheme.Value.Flows.Implicit != nil && securityScheme.Value.Flows.Implicit.Scopes != nil {
					for scope := range securityScheme.Value.Flows.Implicit.Scopes {
						scopeSet[scope] = true
					}
				}
			}
		}
		// Check OAuth2Security scheme
		if securityScheme, exists := openAPI.Components.SecuritySchemes[constants.OAUTH2_SECURITY_SCHEMA_KEY]; exists {
			if securityScheme.Value != nil && securityScheme.Value.Flows != nil {
				if securityScheme.Value.Flows.Password != nil && securityScheme.Value.Flows.Password.Scopes != nil {
					for scope := range securityScheme.Value.Flows.Password.Scopes {
						scopeSet[scope] = true
					}
				}
			}
		}
	}
	return sortScopes(scopeSet), nil
}

// getScopeOfOperations gets scopes for a specific operation using security requirements
func (oasParser *OAS3Parser) getScopeOfOperations(oauth2SchemeKey string, operation *openapi3.Operation) []string {
	if operation.Security != nil {
		for _, securityRequirement := range *operation.Security {
			if scopes, exists := securityRequirement[oauth2SchemeKey]; exists {
				return scopes
			}
		}
	}
	return oasParser.getScopeOfOperationsFromExtensions(operation)
}

// getScopeOfOperationsFromExtensions gets scopes from operation extensions
func (oasParser *OAS3Parser) getScopeOfOperationsFromExtensions(operation *openapi3.Operation) []string {
	if operation.Extensions != nil {
		if scopeValue, exists := operation.Extensions[constants.SWAGGER_X_SCOPE]; exists {
			if scopeStr, ok := scopeValue.(string); ok {
				return strings.Split(scopeStr, ",")
			}
		}
	}
	return []string{}
}

// sortScopes sorts the scopes in a map and returns them as a slice of strings
func sortScopes(scopeSet map[string]bool) []string {
	scopes := make([]string, 0, len(scopeSet))
	for scope := range scopeSet {
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)
	return scopes
}

// setScopesToTemplate sets scopes to the URI template
func setScopesToTemplate(template *dto.URITemplate, resourceScopes []string, apiScopes []string) (*dto.URITemplate, error) {
	// Validate that operation scopes exist in the global scopes
	var validScopes []string
	scopeMap := make(map[string]bool)
	for _, scope := range apiScopes {
		scopeMap[scope] = true
	}
	for _, scope := range resourceScopes {
		if scopeMap[scope] {
			validScopes = append(validScopes, scope)
		} else {
			return nil, fmt.Errorf("resource scope '%s' not found", scope)
		}
	}
	template.Scopes = validScopes
	return template, nil
}

// removeBlocksRecursivelyFromJsonObject removes provided key from the json object recursively
func (oasParser *OAS3Parser) removeBlocksRecursivelyFromJSONObject(keyToBeRemoved string, jsonObject map[string]interface{}, definitionUpdated bool) bool {
	if jsonObject == nil {
		return definitionUpdated
	}
	if _, exists := jsonObject[keyToBeRemoved]; exists {
		delete(jsonObject, keyToBeRemoved)
		definitionUpdated = true
	}
	// Recursively check sub-objects
	for _, value := range jsonObject {
		if subObj, ok := value.(map[string]interface{}); ok {
			result := oasParser.removeBlocksRecursivelyFromJSONObject(keyToBeRemoved, subObj, definitionUpdated)
			definitionUpdated = definitionUpdated || result
		}
	}

	return definitionUpdated
}
