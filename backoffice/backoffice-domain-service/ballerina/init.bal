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
import wso2/apk_common_lib as commons;
import ballerina/lang.runtime;

configurable DatasourceConfiguration datasourceConfiguration = ?;
final postgresql:Client|sql:Error dbClient;

configurable commons:IDPConfiguration idpConfiguration = {
        publicKey:{path: "/home/wso2apk/backoffice/security/mg.pem"}
    };
commons:DBBasedOrgResolver organizationResolver = new(datasourceConfiguration);
commons:JWTValidationInterceptor jwtValidationInterceptor = new (idpConfiguration, organizationResolver);
commons:RequestErrorInterceptor requestErrorInterceptor = new;

function init() {
    log:printInfo("Starting APK Backoffice Domain Service...");

    int retryCount = 0;
    int maxRetries = 5;
    while true {
        dbClient = 
            new (host = datasourceConfiguration.host,
                username = datasourceConfiguration.username, 
                password = datasourceConfiguration.password, 
                database = datasourceConfiguration.databaseName, 
                port = datasourceConfiguration.port,
                connectionPool = {maxOpenConnections: datasourceConfiguration.maxPoolSize}
                );
        if dbClient is error {
            if retryCount < maxRetries {
                retryCount = retryCount + 1;
                log:printError("Error while connecting to database. Retrying...");
                runtime:sleep(5); // wait 5 seconds before retrying
            } else {
                return log:printError("Max retries reached. Failed to connect to database.");
            }
        } else {
            log:printInfo("Connected to the database successfully.");
            return;
        }
    }
}
public isolated function getConnection() returns postgresql:Client | error {
    return dbClient;  
}
