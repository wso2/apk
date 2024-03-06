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
import ballerina/jwt;
import ballerinax/prometheus as _;

configurable IDPConfiguration & readonly idpConfiguration = ?;
final postgresql:Client|sql:Error dbClient;
listener http:Listener ep0 = getListener();
final jwt:ValidatorConfig & readonly validatorConfig;
configurable DatasourceConfiguration datasourceConfiguration = ?;

function getListener() returns http:Listener|error {
    return check new (9443, secureSocket = {
        key: {
            certFile: idpConfiguration.keyStores.tls.certFile,
            keyFile: idpConfiguration.keyStores.tls.keyFile
        }
    });
}

function init() {
    dbClient =
        new (host = datasourceConfiguration.host,
    username = datasourceConfiguration.username,
    password = datasourceConfiguration.password,
    database = datasourceConfiguration.databaseName,
    port = datasourceConfiguration.port,
        connectionPool = {maxOpenConnections: datasourceConfiguration.maxPoolSize}
            );
    validatorConfig = {
        issuer: idpConfiguration.tokenIssuerConfiguration.issuer,
        signatureConfig: {
            certFile: idpConfiguration.keyStores.signing.certFile
        }
    };
    log:printInfo("Initialize Non Production OIDC Server..");
}

public isolated function getConnection() returns postgresql:Client|error {
    return dbClient;
}

public isolated function getValidationConfig() returns jwt:ValidatorConfig {
    return validatorConfig;
}

