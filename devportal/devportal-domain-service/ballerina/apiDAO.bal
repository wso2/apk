//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
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

import ballerina/log;
import ballerinax/postgresql;
import ballerina/sql;
import wso2/apk_common_lib as commons;
final string PUBLISHED = "PUBLISHED";
final string PROTOTYPED = "PROTOTYPED";
final string DEPRECATED = "DEPRECATED";

isolated function getAPIByIdDAO(string apiId) returns API|commons:APKError|NotFoundError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `SELECT UUID AS ID,
        API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS, API_TYPE as TYPE, string_to_array(SDK::text,',')::text[] AS SDK , ARTIFACT as ARTIFACT
        FROM API WHERE UUID =${apiId} AND
        STATUS IN (${PUBLISHED},${PROTOTYPED},${DEPRECATED})`;
        API | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90915, message: "API Not Found for provided API ID"}};
            return nfe;
        } else if result is API {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving API";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function getAPIsDAO(string org) returns APIInfo[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID AS ID,
            API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,
            API_TYPE as TYPE, ARTIFACT as ARTIFACT FROM API WHERE ORGANIZATION =${org} AND 
            STATUS IN (${PUBLISHED},${PROTOTYPED},${DEPRECATED})`;
            stream<APIInfo, sql:Error?> apisStream = dbClient->query(query);
            APIInfo[] apis = check from APIInfo api in apisStream select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            string message = "Internal Error occured while retrieving APIs";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getAPIsByQueryDAO(string payload, string org) returns APIInfo[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT DISTINCT UUID AS ID,
            API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,
            API_TYPE as TYPE, ARTIFACT as ARTIFACT FROM API JOIN JSONB_EACH_TEXT(ARTIFACT) e ON true 
            WHERE ORGANIZATION =${org} AND e.value LIKE ${payload} AND 
            STATUS IN (${PUBLISHED},${PROTOTYPED},${DEPRECATED})`;
            stream<APIInfo, sql:Error?> apisStream = dbClient->query(query);
            APIInfo[] apis = check from APIInfo api in apisStream select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            string message = "Internal Error occured while retrieving APIs";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getAPIDefinitionDAO(string apiId) returns APIDefinition|NotFoundError|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `SELECT encode(API_DEFINITION, 'escape')::text AS schemaDefinition, MEDIA_TYPE as type
        FROM API_ARTIFACT WHERE API_UUID =${apiId}`;
        APIDefinition | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90915, message: "API Definition Not Found for provided API ID"}};
            return nfe;
        } else if result is APIDefinition {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printError(result.toString());
            string message = "Internal Error while retrieving API Definition";
            return error(message, result, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getResourceCategoryIdByCategoryTypeDAO(string resourceType) returns int|commons:APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery GET_RESOURCE_CATEGORY_Prefix = `SELECT RESOURCE_CATEGORY_ID FROM RESOURCE_CATEGORIES where RESOURCE_CATEGORY = `; 
        sql:ParameterizedQuery values = `${resourceType}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_RESOURCE_CATEGORY_Prefix, values);
        int|sql:Error result =  db_Client->queryRow(sqlQuery);
        if result is int {
            return result;
        } else {
            log:printError(result.toString());
            string message = "Internal Error while retrieving resource category";
            return error(message, result, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getResourceByResourceCategoryDAO(string apiId, int resourceCategoryId) returns Resource|NotFoundError|commons:APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery sqlQuery = `SELECT UUID AS resourceUUID, API_UUID AS apiUuid, RESOURCE_CATEGORY_ID AS resourceCategoryId, DATA_TYPE AS dataType,
        RESOURCE_CONTENT AS resourceContent,  RESOURCE_BINARY_VALUE AS resourceBinaryValue  
        FROM API_RESOURCES where API_UUID = ${apiId} AND RESOURCE_CATEGORY_ID = ${resourceCategoryId}`;
        Resource|sql:Error result =  db_Client->queryRow(sqlQuery);
        
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90915, message: "Thumbnail Not Found for provided API ID"}};
            return nfe;
        } else if result is Resource {
            return result;
        } else {
            log:printError(result.toString());
            string message = "Internal Error while retrieving resource";
            return error(message, result, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getDocumentByDocumentIdDAO(string documentId, string apiId) returns DocumentMetaData|NotFoundError|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery GET_DOCUMENT_Prefix = `SELECT UUID AS documentId, RESOURCE_UUID AS resourceId, NAME AS name, SUMMARY AS summary,
        TYPE AS documentType, OTHER_TYPE_NAME AS otherTypeName, SOURCE_URL AS sourceUrl,
        SOURCE_TYPE AS sourceType FROM API_DOC_META_DATA where UUID = `;
        sql:ParameterizedQuery values = `${documentId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_DOCUMENT_Prefix, values);
        DocumentMetaData|sql:Error result = db_Client->queryRow(sqlQuery);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body: {code: 90915, message: "Document Not Found for provided Document ID"}};
            return nfe;
        } else if result is DocumentMetaData {
            return result;
        } else {
            log:printError(result.toString());
            string message = "Internal Error while retrieving Document";
            return error(message, result, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getDocumentsDAO(string apiId) returns Document[]|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery GET_DOCUMENTS_Query = `SELECT UUID AS documentId, NAME AS name, SUMMARY AS summary,
        TYPE AS documentType, OTHER_TYPE_NAME AS otherTypeName, SOURCE_URL AS sourceUrl,
        SOURCE_TYPE AS sourceType FROM API_DOC_META_DATA where API_UUID = ${apiId}`;
            stream<Document, sql:Error?> documentStream = db_Client->query(GET_DOCUMENTS_Query);
            Document[]|sql:Error documents = from Document document in documentStream
                select document;
            sql:Error?? close = documentStream.close();
            if documents is sql:Error {
                log:printError(documents.toString());
                string message = "Internal Error while retrieving Document List";
                return error(message, documents, message = message, description = message, code = 909001, statusCode = 500);
            } else if close is sql:Error {
                log:printError(close.toString());
                string message = "Internal Error while retrieving Document List";
                return error(message, close, message = message, description = message, code = 909001, statusCode = 500);
            } else {
                return documents;
            }
        }
    }
}

isolated function getResourceByResourceIdDAO(string resourceId) returns Resource|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery GET_RESOURCE_Prefix = `SELECT UUID AS resourceUUID, API_UUID AS apiUuid, RESOURCE_CATEGORY_ID AS resourceCategoryId, DATA_TYPE AS dataType,
        RESOURCE_CONTENT AS resourceContent,  RESOURCE_BINARY_VALUE AS resourceBinaryValue  
        FROM API_RESOURCES where UUID = `;
        sql:ParameterizedQuery values = `${resourceId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_RESOURCE_Prefix, values);
        Resource|sql:Error result = db_Client->queryRow(sqlQuery);
        if result is Resource {
            return result;
        } else {
            log:printError(result.toString());
            string message = "Internal Error while retrieving resource category";
            return error(message, result, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}
