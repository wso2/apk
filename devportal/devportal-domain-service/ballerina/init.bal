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
import ballerinax/java.jdbc;
import ballerina/sql;

configurable DatasourceConfiguration datasourceConfiguration = ?;
configurable ThrottlingConfiguration throttleConfig = ?;
jdbc:Client|sql:Error dbClient;
sql:ConnectionPool connPool;

function init() {
    log:printInfo("Starting APK Devportal Domain Service...");
    APKConfiguration apkConfig = {
        throttlingConfiguration: throttleConfig,
        datasourceConfiguration: datasourceConfiguration
    };
    connPool = {maxOpenConnections: datasourceConfiguration.maxPoolSize};
    dbClient = new(datasourceConfiguration.url, datasourceConfiguration.username, 
        datasourceConfiguration.password,
        connectionPool = connPool);
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return log:printError("Error while connecting to database");
    }
}

public function getConnection() returns jdbc:Client | error {
    dbClient = new(datasourceConfiguration.url, datasourceConfiguration.username, 
        datasourceConfiguration.password,
        connectionPool = connPool);
    return dbClient;  
}