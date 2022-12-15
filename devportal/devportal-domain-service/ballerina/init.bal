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
import ballerina/uuid;

configurable DatasourceConfiguration datasourceConfiguration = ?;
configurable ThrottlingConfiguration throttleConfig = ?;
string kid = uuid:createType1AsString();

postgresql:Client|sql:Error dbClient;
APKConfiguration apkConfig;

function init() {
    log:printInfo("Starting APK Devportal Domain Service...");
    apkConfig = {
        throttlingConfiguration: throttleConfig,
        datasourceConfiguration: datasourceConfiguration,
        tokenIssuerConfiguration: {keyId: kid},
        keyStores: {
        signing: {
            path: "/home/wso2apk/devportal/security/wso2carbon.key"
        },
        tls: {
            path: "/home/wso2apk/devportal/security/wso2carbon.key"
        }
    }
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

public function getConnection() returns postgresql:Client | error {
    return dbClient;  
}