
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

import ballerinax/postgresql;
import ballerina/sql;
import ballerina/log;
import wso2/apk_common_lib as commons;

isolated function db_addCategory(string categoryName) returns int|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_CATEGORY_Prefix = `SELECT RESOURCE_CATEGORY_ID FROM RESOURCE_CATEGORIES WHERE RESOURCE_CATEGORY = `;
        sql:ParameterizedQuery values1 = `${categoryName}`;
        sql:ParameterizedQuery sqlQuery1 = sql:queryConcat(GET_CATEGORY_Prefix, values1);
        int|sql:Error result1 =  db_Client->queryRow(sqlQuery1);
        if result1 is int {
            log:printDebug("Resource category " + categoryName + " added successfully");
            return result1;
        } else if result1 is sql:NoRowsError {
            sql:ParameterizedQuery GET_CATEGORY_INSERT_Prefix = `INSERT INTO RESOURCE_CATEGORIES (RESOURCE_CATEGORY) VALUES (`;
            sql:ParameterizedQuery values2 = `${categoryName})`;
            sql:ParameterizedQuery sqlQuery2 = sql:queryConcat(GET_CATEGORY_INSERT_Prefix, values2);
            sql:ExecutionResult|sql:Error result2 = db_Client->execute(sqlQuery2);
            if result2 is sql:ExecutionResult {
                log:printDebug("Resource category added successfully");
            } else {
                return e909618(result2);
            }
        } else if result1 is sql:Error {
            return e909618(result1);
        }
    }
}
