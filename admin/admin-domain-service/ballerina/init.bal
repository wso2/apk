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

import ballerina/io;
import ballerinax/java.jdbc;
import ballerina/sql;
import ballerina/http;

configurable DatasourceConfiguration datasourceConfiguration = ?;
jdbc:Client|sql:Error dbClient;
configurable ThrottlingConfiguration throttleConfig = ?;

configurable int ADMIN_PORT = 9443;

listener http:Listener ep0 = new (ADMIN_PORT);

function init() {
    io:println("Starting APK Admin Domain Service...");
    APKConfiguration apkConfig = {
        throttlingConfiguration: throttleConfig,
        datasourceConfiguration: datasourceConfiguration
    };
    dbClient = 
         new(datasourceConfiguration.url, datasourceConfiguration.username, 
             datasourceConfiguration.password,
             connectionPool = { maxOpenConnections: datasourceConfiguration.maxPoolSize });
}

public function getConnection() returns jdbc:Client | error {
     return dbClient;  
 }
