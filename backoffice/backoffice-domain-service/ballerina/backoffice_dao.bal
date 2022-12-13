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

import ballerina/sql;
import ballerinax/postgresql;

function db_getAPIsDAO() returns API[]|error? {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        stream<API, sql:Error?> apisStream = db_Client->query(GET_API);
        API[]? apis = check from API api in apisStream select api;
        check apisStream.close();
        return apis;
    }
}

function db_changeLCState(string action, string apiId, string organization) returns string|error? {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        string newState = actionToLCState(action);
        sql:ParameterizedQuery values = `${newState}
        WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_LifeCycle_Prefix, values);

        sql:ExecutionResult | sql:Error result = db_Client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return action;
        } else {
            return error("Error while updating LC state into Database");  
        }
    }
}

function db_getCurrentLCStatus(string apiId, string organization) returns string|error? {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_API_LifeCycle_Prefix, values);

        string | sql:Error result =  db_Client->queryRow(sqlQuery);
        
        if result is string {
            return result;
        } else {
            return error("Error while geting LC state from Database");  
        }
    }
}