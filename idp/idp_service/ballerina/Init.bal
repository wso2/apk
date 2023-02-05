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
import ballerina/log;
import ballerina/sql;
import ballerina/http;
configurable IDPConfiguration idpConfiguration = ?;
final postgresql:Client|sql:Error dbClient;
listener http:Listener ep0 = new (9090);
function init() {
    dbClient =
        new (host = idpConfiguration.dataSource.host,
            username = idpConfiguration.dataSource.username,
            password = idpConfiguration.dataSource.password,
            database = idpConfiguration.dataSource.databaseName,
            port = idpConfiguration.dataSource.port,
            connectionPool = {maxOpenConnections: idpConfiguration.dataSource.maxPoolSize}
            );
    log:printInfo("Initialize Non Production OIDC Server..");
}

public isolated function getConnection() returns postgresql:Client|error {
    return dbClient;
}

