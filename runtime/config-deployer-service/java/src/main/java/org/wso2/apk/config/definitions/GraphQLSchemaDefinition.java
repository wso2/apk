//
// Copyright (c) 2024, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package org.wso2.apk.config.definitions;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.config.APIConstants;
import org.wso2.apk.config.model.URITemplate;
import org.wso2.apk.config.queryanalysis.GraphqlSchemaType;

import graphql.language.FieldDefinition;
import graphql.language.ObjectTypeDefinition;
import graphql.language.OperationTypeDefinition;
import graphql.language.SchemaDefinition;
import graphql.language.TypeDefinition;
import graphql.schema.idl.SchemaParser;
import graphql.schema.idl.TypeDefinitionRegistry;

public class GraphQLSchemaDefinition {
    protected Log log = LogFactory.getLog(getClass());

    /**
     * Extract GraphQL Operations from given schema.
     *
     * @param typeRegistry graphQL Schema Type Registry
     * @param type         operation type string
     * @return the arrayList of APIOperationsDTO
     */
    public static List<URITemplate> extractGraphQLOperationList(TypeDefinitionRegistry typeRegistry, String type) {
        List<URITemplate> operationArray = new ArrayList<>();
        Map<java.lang.String, TypeDefinition> operationList = typeRegistry.types();
        for (Map.Entry<String, TypeDefinition> entry : operationList.entrySet()) {
            Optional<SchemaDefinition> schemaDefinition = typeRegistry.schemaDefinition();
            if (schemaDefinition.isPresent()) {
                List<OperationTypeDefinition> operationTypeList = schemaDefinition.get().getOperationTypeDefinitions();
                for (OperationTypeDefinition operationTypeDefinition : operationTypeList) {
                    boolean canAddOperation = entry.getValue().getName()
                            .equalsIgnoreCase(operationTypeDefinition.getTypeName().getName()) &&
                            (type == null || type.equals(operationTypeDefinition.getName().toUpperCase()));
                    if (canAddOperation) {
                        addOperations(entry, operationTypeDefinition.getName().toUpperCase(), operationArray);
                    }
                }
            } else {
                boolean canAddOperation = (entry.getValue().getName().equalsIgnoreCase(APIConstants.GRAPHQL_QUERY) ||
                        entry.getValue().getName().equalsIgnoreCase(APIConstants.GRAPHQL_MUTATION)
                        || entry.getValue().getName().equalsIgnoreCase(APIConstants.GRAPHQL_SUBSCRIPTION)) &&
                        (type == null || type.equals(entry.getValue().getName().toUpperCase()));
                if (canAddOperation) {
                    addOperations(entry, entry.getKey(), operationArray);
                }
            }
        }
        return operationArray;
    }

    /**
     * @param entry          Entry
     * @param operationArray operationArray
     */
    private static void addOperations(Map.Entry<String, TypeDefinition> entry, String graphQLType,
            List<URITemplate> operationArray) {
        for (FieldDefinition fieldDef : ((ObjectTypeDefinition) entry.getValue()).getFieldDefinitions()) {
            URITemplate operation = new URITemplate();
            operation.setVerb(graphQLType);
            operation.setUriTemplate(fieldDef.getName());
            operationArray.add(operation);
        }
    }
}
