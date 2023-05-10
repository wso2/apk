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

final string PUBLISHED = "PUBLISHED";
final string PROTOTYPED = "PROTOTYPED";
final string DEPRECATED = "DEPRECATED";

isolated function getAPIByIdDAO(string apiId) returns API|APKError|NotFoundError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
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
            return error(message, result, message = message, description = message, code = 909000, statusCode = "500");
        }
    }
}

isolated function getAPIsDAO(string org) returns API[]|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID AS ID,
            API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,
            API_TYPE as TYPE, ARTIFACT as ARTIFACT FROM API WHERE ORGANIZATION =${org} AND 
            STATUS IN (${PUBLISHED},${PROTOTYPED},${DEPRECATED})`;
            stream<API, sql:Error?> apisStream = dbClient->query(query);
            API[] apis = check from API api in apisStream select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            string message = "Internal Error occured while retrieving APIs";
            return error(message, e, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function getAPIsByQueryDAO(string payload, string org) returns API[]|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT DISTINCT UUID AS ID,
            API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,
            API_TYPE as TYPE, ARTIFACT as ARTIFACT FROM API JOIN JSONB_EACH_TEXT(ARTIFACT) e ON true 
            WHERE ORGANIZATION =${org} AND e.value LIKE ${payload} AND 
            STATUS IN (${PUBLISHED},${PROTOTYPED},${DEPRECATED})`;
            stream<API, sql:Error?> apisStream = dbClient->query(query);
            API[] apis = check from API api in apisStream select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            string message = "Internal Error occured while retrieving APIs";
            return error(message, e, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function getAPIDefinitionDAO(string apiId) returns APIDefinition|NotFoundError|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
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
            return error(message, result, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function db_getResourceCategoryIdByCategoryType(string resourceType) returns int|APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = "500");
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
            return error(message, result, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function db_getResourceByResourceCategory(string apiId, int resourceCategoryId) returns Resource|boolean|APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        string message = "Error while retrieving connection";
        return error(message, db_Client, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery sqlQuery = `SELECT UUID AS resourceUUID, API_UUID AS apiUuid, RESOURCE_CATEGORY_ID AS resourceCategoryId, DATA_TYPE AS dataType,
        RESOURCE_CONTENT AS resourceContent,  RESOURCE_BINARY_VALUE AS resourceBinaryValue  
        FROM API_RESOURCES where API_UUID = ${apiId} AND RESOURCE_CATEGORY_ID = ${resourceCategoryId}`;
        Resource|sql:Error result =  db_Client->queryRow(sqlQuery);
        
        if result is sql:NoRowsError {
            return false;
        } else if result is Resource {
            return result;
        } else {
            log:printError(result.toString());
            string message = "Internal Error while retrieving resource";
            return error(message, result, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}
