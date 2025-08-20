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

package parsers

import (
	"fmt"
	"github.com/graphql-go/graphql/language/ast"
	graphqlParser "github.com/graphql-go/graphql/language/parser"
	"github.com/wso2/apk/config-deployer-service-go/internal/constants"
	"github.com/wso2/apk/config-deployer-service-go/internal/dto"
	"strings"
)

type GraphQLParser struct{}

// GetAPIFromDefinition creates an API from GraphQL schema definition
func (graphQLParser *GraphQLParser) GetAPIFromDefinition(definition string) (*dto.API, error) {
	document, err := graphqlParser.Parse(graphqlParser.ParseParams{
		Source: definition,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse GraphQL schema: %w", err)
	}
	var combinedUriTemplates []dto.URITemplate

	// Extract and add all URI templates for query, mutation, and subscription into a combined list
	queryTemplates := graphQLParser.extractGraphQLOperationList(document, constants.GRAPHQL_QUERY)
	mutationTemplates := graphQLParser.extractGraphQLOperationList(document, constants.GRAPHQL_MUTATION)
	subscriptionTemplates := graphQLParser.extractGraphQLOperationList(document, constants.GRAPHQL_SUBSCRIPTION)
	combinedUriTemplates = append(combinedUriTemplates, queryTemplates...)
	combinedUriTemplates = append(combinedUriTemplates, mutationTemplates...)
	combinedUriTemplates = append(combinedUriTemplates, subscriptionTemplates...)

	api := &dto.API{}
	api.URITemplates = combinedUriTemplates
	api.GraphQLSchema = definition

	return api, nil
}

// extractGraphQLOperationList extracts GraphQL operations from given schema
func (graphQLParser *GraphQLParser) extractGraphQLOperationList(document *ast.Document, operationType string) []dto.URITemplate {
	var operationArray []dto.URITemplate

	// Create a map of type definitions for easy lookup
	typeDefinitions := make(map[string]ast.TypeDefinition)
	var schemaDefinition *ast.SchemaDefinition

	// First pass: collect all type definitions and find schema definition
	for _, definition := range document.Definitions {
		switch def := definition.(type) {
		case *ast.ObjectDefinition:
			typeDefinitions[def.Name.Value] = def
		case *ast.SchemaDefinition:
			schemaDefinition = def
		}
	}

	// Process each type definition
	for typeName, typeDef := range typeDefinitions {
		if objDef, ok := typeDef.(*ast.ObjectDefinition); ok {

			// Check if schema definition exists
			if schemaDefinition != nil {
				// Use explicit schema definition
				for _, operationTypeDef := range schemaDefinition.OperationTypes {
					//canAddOperation := strings.EqualFold(typeName, operationTypeDef.Type.(*ast.Named).Name.Value) &&
					//	(operationType == "" || operationType == strings.ToUpper(operationTypeDef.Operation))
					canAddOperation := strings.EqualFold(typeName, operationTypeDef.Type.Name.Value) &&
						(operationType == "" || operationType == strings.ToUpper(operationTypeDef.Operation))
					if canAddOperation {
						graphQLParser.addOperations(typeName, objDef, strings.ToUpper(operationTypeDef.Operation), &operationArray)
					}
				}
			} else {
				// Use implicit schema definition (Query, Mutation, Subscription)
				canAddOperation := (strings.EqualFold(typeName, constants.GRAPHQL_QUERY) ||
					strings.EqualFold(typeName, constants.GRAPHQL_MUTATION) ||
					strings.EqualFold(typeName, constants.GRAPHQL_SUBSCRIPTION)) &&
					(operationType == "" || operationType == strings.ToUpper(typeName))

				if canAddOperation {
					graphQLParser.addOperations(typeName, objDef, strings.ToUpper(typeName), &operationArray)
				}
			}
		}
	}
	return operationArray
}

// addOperations adds operations from a type definition to the operation array
func (graphQLParser *GraphQLParser) addOperations(typeName string, objDef *ast.ObjectDefinition, graphQLType string,
	operationArray *[]dto.URITemplate) {
	for _, fieldDef := range objDef.Fields {
		operation := &dto.URITemplate{
			URITemplate: fieldDef.Name.Value,
			Verb:        graphQLType,
			AuthEnabled: true,
			Scopes:      []string{},
		}
		operation.Verb = graphQLType
		operation.URITemplate = fieldDef.Name.Value
		*operationArray = append(*operationArray, *operation)
	}
}
