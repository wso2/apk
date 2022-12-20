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
import ballerina/http;
import ballerina/lang.runtime;
import ballerina/uuid;

listener http:Listener ep0 = new (9443);
string kid = uuid:createType1AsString();

configurable RuntimeConfiguratation & readonly runtimeConfiguration = {
    keyStores: {
        signing: {
            path: "/home/wso2apk/runtime/security/wso2carbon.key"
        },
        tls: {
            path: "/home/wso2apk/runtime/security/wso2carbon.key"
        }
    }
};

# Initializing method for runtime
# + return - Return Error if error occured at initialization.
function init() returns error? {
    APIClient apiService = new ();
    error? retrieveAllApisAtStartup = apiService.retrieveAllApisAtStartup(());
    if retrieveAllApisAtStartup is error {
        log:printError("Error occured while retrieving API List", retrieveAllApisAtStartup);
    }

    ServiceClient servicesService = new ();
    error? retrieveAllServicesAtStartup = servicesService.retrieveAllServicesAtStartup(());
    if retrieveAllServicesAtStartup is error {
        log:printError("Error occured while retrieving Service List", retrieveAllServicesAtStartup);
    }

    APIListingTask apiListingTask = new (resourceVersion);
    _ = check apiListingTask.startListening();
    ServiceTask serviceTask = new (servicesResourceVersion);
    _ = serviceTask.startListening();
    _ = check servicesService.retrieveAllServiceMappingsAtStartup(());
    ServiceMappingTask serviceMappingTask = new (serviceMappingResourceVersion);
    _ = check serviceMappingTask.startListening();
    _ = check startAndAttachServices();
    log:printInfo("Initializing Runtime Domain Service..");
}

public function deRegisterep() returns error? {
    _ = check ep0.gracefulStop();
}

function startAndAttachServices() returns error? {
    check ep0.attach(healthService, "/");
    check ep0.attach(runtimeService, "/api/am/runtime");
    check ep0.'start();
    runtime:registerListener(ep0);
    runtime:onGracefulStop(deRegisterep);
}
