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
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"regexp"
)

// GraphQLSchemaBuilder handles GraphQL schema building operations
type GraphQLSchemaBuilder struct {
	typeMap      map[string]*graphql.Object
	inputTypeMap map[string]*graphql.InputObject
	enumTypeMap  map[string]*graphql.Enum
}

// PreprocessGraphQLSchema removes invalid null default values from GraphQL schema
func PreprocessGraphQLSchema(schema string) string {
	// Regular expression to match field definitions with "= null" default values
	// Matches patterns like: "fieldName: [Type!] = null" or "fieldName: Type = null"
	nullDefaultRegex := regexp.MustCompile(`(\w+:\s*(?:\[[A-Za-z_][A-Za-z0-9_]*!?]!?|[A-Za-z_][A-Za-z0-9_]*!?))\s*=\s*null`)

	// Replace "= null" with empty string (removing the default value)
	cleanedSchema := nullDefaultRegex.ReplaceAllString(schema, "$1")

	return cleanedSchema
}

// NewGraphQLSchemaBuilder creates a new GraphQL schema builder
func NewGraphQLSchemaBuilder() *GraphQLSchemaBuilder {
	return &GraphQLSchemaBuilder{
		typeMap:      make(map[string]*graphql.Object),
		inputTypeMap: make(map[string]*graphql.InputObject),
		enumTypeMap:  make(map[string]*graphql.Enum),
	}
}

// BuildObjectType builds a GraphQL object type from AST definition
func (b *GraphQLSchemaBuilder) BuildObjectType(def *ast.ObjectDefinition) *graphql.Object {
	if existingType, exists := b.typeMap[def.Name.Value]; exists {
		return existingType
	}
	fields := make(graphql.Fields)
	for _, field := range def.Fields {
		fieldType := b.resolveFieldType(field.Type)
		fields[field.Name.Value] = &graphql.Field{
			Name: field.Name.Value,
			Type: fieldType,
		}
	}
	objectType := graphql.NewObject(graphql.ObjectConfig{
		Name:   def.Name.Value,
		Fields: fields,
	})
	// Store in typeMap for future reference
	b.typeMap[def.Name.Value] = objectType
	return objectType
}

// BuildInputObjectType builds a GraphQL input object type from AST definition
func (b *GraphQLSchemaBuilder) BuildInputObjectType(def *ast.InputObjectDefinition) *graphql.InputObject {
	if existingType, exists := b.inputTypeMap[def.Name.Value]; exists {
		return existingType
	}
	fields := make(graphql.InputObjectConfigFieldMap)
	for _, field := range def.Fields {
		fieldType := b.resolveInputFieldType(field.Type)
		fields[field.Name.Value] = &graphql.InputObjectFieldConfig{
			Type: fieldType,
		}
	}
	inputObjectType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   def.Name.Value,
		Fields: fields,
	})
	// Store in inputTypeMap for future reference
	b.inputTypeMap[def.Name.Value] = inputObjectType
	return inputObjectType
}

// BuildEnumType builds a GraphQL enum type from AST definition
func (b *GraphQLSchemaBuilder) BuildEnumType(def *ast.EnumDefinition) *graphql.Enum {
	if existingType, exists := b.enumTypeMap[def.Name.Value]; exists {
		return existingType
	}
	values := make(graphql.EnumValueConfigMap)
	for _, value := range def.Values {
		values[value.Name.Value] = &graphql.EnumValueConfig{
			Value: value.Name.Value,
		}
	}
	enumType := graphql.NewEnum(graphql.EnumConfig{
		Name:   def.Name.Value,
		Values: values,
	})
	// Store in enumTypeMap for future reference
	b.enumTypeMap[def.Name.Value] = enumType
	return enumType
}

// resolveFieldType resolves GraphQL field types from AST
func (b *GraphQLSchemaBuilder) resolveFieldType(typeAST ast.Type) graphql.Output {
	switch t := typeAST.(type) {
	case *ast.Named:
		// Check if it's a custom type in our typeMap
		if customType, exists := b.typeMap[t.Name.Value]; exists {
			return customType
		}
		// Check if it's an enum type
		if enumType, exists := b.enumTypeMap[t.Name.Value]; exists {
			return enumType
		}
		// Default to built-in scalar types
		return b.getScalarType(t.Name.Value)
	case *ast.List:
		innerType := b.resolveFieldType(t.Type)
		return graphql.NewList(innerType)
	case *ast.NonNull:
		innerType := b.resolveFieldType(t.Type)
		return graphql.NewNonNull(innerType)
	default:
		return graphql.String
	}
}

// resolveInputFieldType resolves GraphQL input field types from AST
func (b *GraphQLSchemaBuilder) resolveInputFieldType(typeAST ast.Type) graphql.Input {
	switch t := typeAST.(type) {
	case *ast.Named:
		// Check if it's a custom input type in our inputTypeMap
		if customType, exists := b.inputTypeMap[t.Name.Value]; exists {
			return customType
		}
		// Check if it's an enum type
		if enumType, exists := b.enumTypeMap[t.Name.Value]; exists {
			return enumType
		}
		// Default to built-in scalar types
		return b.getScalarType(t.Name.Value)
	case *ast.List:
		innerType := b.resolveInputFieldType(t.Type)
		return graphql.NewList(innerType)
	case *ast.NonNull:
		innerType := b.resolveInputFieldType(t.Type)
		return graphql.NewNonNull(innerType)
	default:
		return graphql.String
	}
}

// getScalarType returns the appropriate GraphQL scalar type
func (b *GraphQLSchemaBuilder) getScalarType(typeName string) graphql.Type {
	switch typeName {
	case "String":
		return graphql.String
	case "Int":
		return graphql.Int
	case "Float":
		return graphql.Float
	case "Boolean":
		return graphql.Boolean
	case "ID":
		return graphql.ID
	default:
		return graphql.String
	}
}
