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

isolated function getAPIByIdDAO(string apiId, string org) returns string?|API|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT API_UUID AS ID, API_ID as APIID,
        API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS, API_TYPE as TYPE, ARTIFACT as ARTIFACT
        FROM API WHERE API_UUID =${apiId} AND ORGANIZATION =${org}`;
        API | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is API {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving API");
        }
    }
}

isolated function getAPIsDAO(string org) returns API[]|error? {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT API_UUID AS ID, API_ID as APIID,
        API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS, API_TYPE as TYPE, ARTIFACT as ARTIFACT FROM API WHERE ORGANIZATION =${org}`;
        stream<API, sql:Error?> apisStream = dbClient->query(query);
        API[]? apis = check from API api in apisStream select api;
        check apisStream.close();
        return apis;
    }
}

isolated function getAPIDefinitionDAO(string apiId, string org) returns APIDefinition|NotFoundError|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT encode(API_DEFINITION, 'escape')::text AS schemaDefinition, MEDIA_TYPE as type
        FROM API_ARTIFACT WHERE API_UUID =${apiId} AND ORGANIZATION =${org}`;
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
            return error("Internal Error while retrieving API Definition");
        }
    }
}