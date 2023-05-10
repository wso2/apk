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

configurable commons:DatasourceConfiguration datasourceConfiguration = ?;
configurable K8sConfiguration k8sConfig = ?;
configurable ManagementServerConfiguration & readonly managementServerConfig = ?;
final postgresql:Client|sql:Error dbClient;

configurable commons:IDPConfiguration idpConfiguration = {
    publicKey: {certFilePath: "/home/wso2apk/backoffice/security/mg.pem"}
};
configurable KeyStores & readonly keyStores = {
    tls: {certFilePath: "/home/wso2apk/admin/security/backoffice.pem", keyFilePath: "/home/wso2apk/admin/security/backoffice.key"}
};

commons:DBBasedOrgResolver organizationResolver = new (datasourceConfiguration);
commons:JWTValidationInterceptor jwtValidationInterceptor = new (idpConfiguration, organizationResolver);
commons:RequestErrorInterceptor requestErrorInterceptor = new;
commons:ResponseErrorInterceptor responseErrorInterceptor = new;

function init() returns error? {
    log:printInfo("Starting APK Backoffice Domain Service...");

    dbClient =
        new (host = datasourceConfiguration.host,
        username = datasourceConfiguration.username,
        password = datasourceConfiguration.password,
        database = datasourceConfiguration.databaseName,
        port = datasourceConfiguration.port,
        connectionPool = {maxOpenConnections: datasourceConfiguration.maxPoolSize}
    );

    check db_initialiseResourceCategories();
}

public isolated function getConnection() returns postgresql:Client|error {
    return dbClient;
}

isolated function db_initialiseResourceCategories() returns error? {
    string[] catrgoryList = [RESOURCE_TYPE_THUMBNAIL];
    foreach string category in catrgoryList {
        int|commons:APKError result = db_addCategory(category);
        if result is commons:APKError {
            log:printDebug("Error while adding category: " + category);
            return result;
        }
    }
}
