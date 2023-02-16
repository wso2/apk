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
import wso2/apk_common_lib as commons;

configurable RuntimeConfiguratation & readonly runtimeConfiguration = {
    keyStores: {
        signing: {
            path: "/home/wso2apk/runtime/security/wso2carbon.key"
        },
        tls: {
            path: "/home/wso2apk/runtime/security/wso2carbon.key"
        }
    },
    idpConfiguration: {publicKey: {path: "/home/wso2apk/runtime/security/mg.pem"}}
};
K8sBaseOrgResolver organizationResolver = new;
commons:JWTValidationInterceptor jwtValidationInterceptor = new (runtimeConfiguration.idpConfiguration, organizationResolver);
commons:RequestErrorInterceptor requestErrorInterceptor = new;
listener http:Listener ep0 = new (9443, {interceptors: [jwtValidationInterceptor,requestErrorInterceptor]});
string kid = uuid:createType1AsString();

# Initializing method for runtime
# + return - Return Error if error occured at initialization.
function init() returns error? {
    APIClient apiService = new ();
    error? retrieveAllApisAtStartup = apiService.retrieveAllApisAtStartup((), ());
    if retrieveAllApisAtStartup is error {
        log:printError("Error occured while retrieving API List", retrieveAllApisAtStartup);
    }

    ServiceClient servicesService = new ();
    error? retrieveAllServicesAtStartup = servicesService.retrieveAllServicesAtStartup((), ());
    if retrieveAllServicesAtStartup is error {
        log:printError("Error occured while retrieving Service List", retrieveAllServicesAtStartup);
    }
    OrgClient orgClient = new ();
    error? retrieveAllOrganizationsAtStartup = orgClient.retrieveAllOrganizationsAtStartup((), ());
    if retrieveAllOrganizationsAtStartup is error {
        log:printError("Error occured while retrieving Organization List", retrieveAllOrganizationsAtStartup);
    }
    APIListingTask apiListingTask = new (resourceVersion);
    _ = check apiListingTask.startListening();
    ServiceTask serviceTask = new (servicesResourceVersion);
    _ = check serviceTask.startListening();
    _ = check servicesService.retrieveAllServiceMappingsAtStartup((), ());
    ServiceMappingTask serviceMappingTask = new (serviceMappingResourceVersion);
    _ = check serviceMappingTask.startListening();
    OrganizationListingTask organizationListingTask = new (organizationResourceVersion);
    _ = check organizationListingTask.startListening();
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
