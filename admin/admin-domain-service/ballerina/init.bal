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
import ballerina/http;

configurable DatasourceConfiguration datasourceConfiguration = ?;
final postgresql:Client|sql:Error dbClient;
configurable ThrottlingConfiguration throttleConfig = ?;

configurable int ADMIN_PORT = 9443;

listener http:Listener ep0 = new (ADMIN_PORT);

function init() {
    log:printInfo("Starting APK Admin Domain Service...");
    APKConfiguration apkConfig = {
        throttlingConfiguration: throttleConfig,
        datasourceConfiguration: datasourceConfiguration
    };
    dbClient = 
        new (host = datasourceConfiguration.host,
            username = datasourceConfiguration.username, 
            password = datasourceConfiguration.password, 
            database = datasourceConfiguration.databaseName, 
            port = datasourceConfiguration.port,
            connectionPool = {maxOpenConnections: datasourceConfiguration.maxPoolSize}
            );
    if dbClient is error {
        return log:printError("Error while connecting to database");
    }
}

public isolated function getConnection() returns postgresql:Client | error {
    return dbClient;  
 }


