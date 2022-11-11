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
import backoffice_domain_service.org.wso2.apk.apimgt.api as api;
import backoffice_domain_service.org.wso2.apk.apimgt.init as apkinit;

configurable DatasourceConfiguration datasourceConfiguration = ?;

function init() {
    io:println("Starting APK Backoffice Domain Service...");
    APKConfiguration apkConfig = {
        datasourceConfiguration: datasourceConfiguration
    };
    string configJson = apkConfig.toJson().toJsonString();
    // Pass the configurations to java init component
    api:APIManagementException? err = apkinit:APKComponent_activate(configJson);
    if (err != ()) {
        io:println(err);
    }
}
