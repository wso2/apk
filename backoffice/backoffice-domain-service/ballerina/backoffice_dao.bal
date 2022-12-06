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

function getAPIsDAO() returns API[]|error? {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM api`;
        stream<API, sql:Error?> apisStream = db_Client->query(query);
        API[]? apis = check from API api in apisStream select api;
        check apisStream.close();
        return apis;
    }
}
