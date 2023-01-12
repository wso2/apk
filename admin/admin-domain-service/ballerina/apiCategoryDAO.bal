//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
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

isolated function addAPICategoryDAO(APICategory payload, string org) returns APICategory|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `INSERT INTO API_CATEGORIES (UUID, NAME, 
        DESCRIPTION, ORGANIZATION) VALUES (${payload.id},${payload.name},
        ${payload.description},${org})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return payload;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}

public isolated function checkAPICategoryExistsByNameDAO(string categoryName, string org) returns boolean|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT UUID as ID, NAME, DESCRIPTION 
        FROM API_CATEGORIES WHERE NAME =${categoryName} AND ORGANIZATION =${org}`;
        APICategory | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return false;
        } else if result is APICategory {
            return true;
        } else {
            log:printDebug(result.toString());
            return error("Error while checking API Category existence");
        }
    }
}

isolated function getAPICategoriesDAO(string org) returns APICategory[]|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT UUID as ID, NAME, DESCRIPTION 
        FROM API_CATEGORIES WHERE ORGANIZATION =${org} ORDER BY NAME`;
        stream<APICategory, sql:Error?> apiCategoryStream = dbClient->query(query);
        APICategory[] apiCategoryList = check from APICategory apiCategory in apiCategoryStream select apiCategory;
        check apiCategoryStream.close();
        return apiCategoryList;
    }
}

isolated function getAPICategoryByIdDAO(string id, string org) returns APICategory|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT UUID as ID, NAME, DESCRIPTION 
        FROM API_CATEGORIES WHERE UUID =${id} AND ORGANIZATION =${org}`;
        APICategory | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is APICategory {
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while checking API Category existence");
        }
    }
}

isolated function updateAPICategoryDAO(APICategory body, string org) returns APICategory|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `UPDATE API_CATEGORIES SET NAME = ${body.name},
         DESCRIPTION = ${body.description} WHERE UUID = ${body.id} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return body;
        } else {
            log:printDebug(result.toString());
            return error("Error while updating data record in the Database");  
        }
    }
}

isolated function deleteAPICategoryDAO(string id, string org) returns error?|string {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `DELETE FROM API_CATEGORIES WHERE UUID = ${id} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return ();
        } else {
            log:printDebug(result.toString());
            return error("Error while deleting data record in the Database");  
        }
    }
}
