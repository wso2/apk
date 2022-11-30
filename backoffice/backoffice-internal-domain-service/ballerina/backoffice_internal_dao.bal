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
import ballerinax/java.jdbc;
import ballerina/uuid;

# Add API details to the database 
#
# + api - API Parameter
# + return - API | error
public function db_createAPI(API api) returns API | error {
    jdbc:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery values = `${uuid:createType1AsString()},${api.name}, ${api.'version}, ${api.context},${api.provider},${api.lifeCycleStatus}, '{}')`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_API_Suffix, values);
        
        
        sql:ExecutionResult | sql:Error result = check db_client->execute(sqlQuery);
        check db_client.close();
        if result is sql:ExecutionResult {
            return api;
        } else {
            return error("Error while inserting data into Database");  
        }
    }
}
