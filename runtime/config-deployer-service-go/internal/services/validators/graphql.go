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

package validators

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"github.com/wso2/apk/config-deployer-service-go/internal/util"
	"strings"
)

type GraphQLAPIValidator struct{}

// ValidateGraphQLSchema validates the GraphQL schema and returns the validation response
func (graphQLAPIValidator *GraphQLAPIValidator) ValidateGraphQLSchema(apiDefinition string, returnGraphQLSchemaContent bool) (*dto.APIDefinitionValidationResponse, error) {
	validationResponse := &dto.APIDefinitionValidationResponse{}

	if strings.TrimSpace(apiDefinition) == "" {
		validationResponse.IsValid = false
		return validationResponse, fmt.Errorf("GraphQL Schema cannot be empty or null")
	}
	// Preprocess the schema to remove invalid null default values
	cleanedDefinition := util.PreprocessGraphQLSchema(apiDefinition)
	validationErrors := validateGraphQLSchemaDefinition(cleanedDefinition)
	if len(validationErrors) > 0 {
		validationResponse.IsValid = false
		return validationResponse, fmt.Errorf("this API is not a GraphQL API")
	} else {
		validationResponse.IsValid = true
		validationResponse.Content = cleanedDefinition
	}

	return validationResponse, nil
}

// validateGraphQLSchemaDefinition validates the GraphQL schema and returns validation errors
func validateGraphQLSchemaDefinition(apiDefinition string) []string {
	var validationErrors []string

	document, err := parser.Parse(parser.ParseParams{
		Source: &source.Source{
			Body: []byte(apiDefinition),
			Name: "GraphQL Schema",
		},
	})
	if err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Parse error: %s", err.Error()))
		return validationErrors
	}

	schemaBuilder := util.NewGraphQLSchemaBuilder()
	schemaConfig := graphql.SchemaConfig{}
	typeMap := make(map[string]*graphql.Object)
	inputTypeMap := make(map[string]*graphql.InputObject)
	enumTypeMap := make(map[string]*graphql.Enum)

	var queryType *graphql.Object
	var mutationType *graphql.Object
	var subscriptionType *graphql.Object

	for _, definition := range document.Definitions {
		switch def := definition.(type) {
		case *ast.ObjectDefinition:
			objType := schemaBuilder.BuildObjectType(def, typeMap)
			typeMap[def.Name.Value] = objType
			switch strings.ToLower(def.Name.Value) {
			case "query":
				queryType = objType
			case "mutation":
				mutationType = objType
			case "subscription":
				subscriptionType = objType
			}
		case *ast.InputObjectDefinition:
			inputType := schemaBuilder.BuildInputObjectType(def, inputTypeMap)
			inputTypeMap[def.Name.Value] = inputType
		case *ast.EnumDefinition:
			enumType := schemaBuilder.BuildEnumType(def)
			enumTypeMap[def.Name.Value] = enumType
		case *ast.SchemaDefinition:
			for _, operationTypeDef := range def.OperationTypes {
				typeName := operationTypeDef.Type.Name.Value
				switch operationTypeDef.Operation {
				case "query":
					if objType, exists := typeMap[typeName]; exists {
						queryType = objType
					}
				case "mutation":
					if objType, exists := typeMap[typeName]; exists {
						mutationType = objType
					}
				case "subscription":
					if objType, exists := typeMap[typeName]; exists {
						subscriptionType = objType
					}
				}
			}
		}
	}

	if queryType != nil {
		schemaConfig.Query = queryType
	}
	if mutationType != nil {
		schemaConfig.Mutation = mutationType
	}
	if subscriptionType != nil {
		schemaConfig.Subscription = subscriptionType
	}
	_, err = graphql.NewSchema(schemaConfig)
	if err != nil {
		validationErrors = append(validationErrors, fmt.Sprintf("Schema validation error: %s", err.Error()))
	}

	return validationErrors
}
