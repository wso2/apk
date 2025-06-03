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
import ballerina/io;
import wso2/apk_common_lib as commons;
import config_deployer_service.org.wso2.apk.config as runtimeUtil;

configurable KeyStores keyStores = {
    tls: {
        keyFilePath: "/home/wso2kgw/config-deployer/security/config.key"
    }
};
configurable (string & readonly) apkSchemaLocation = "/home/wso2kgw/config-deployer/conf/apk-schema.json";
configurable (K8sConfigurations & readonly) k8sConfiguration = {};
configurable (GatewayConfigurations & readonly) gatewayConfiguration = {};
configurable (PartitionServiceConfiguration & readonly) partitionServiceConfiguration = {};
configurable (Vhost[] & readonly) vhosts = [{name: "Default", hosts: ["gw.wso2.com"], 'type: PRODUCTION_TYPE}, {name: "Default", hosts: ["sandbox.gw.wso2.com"], 'type: SANDBOX_TYPE}];
commons:RequestErrorInterceptor requestErrorInterceptor = new;
commons:ResponseErrorInterceptor responseErrorInterceptor = new;
configurable (commons:IDPConfiguration & readonly) idpConfiguration = {publicKey: {certFilePath: "/home/wso2kgw/config-deployer/security/mg.pem"}};
commons:JWTValidationInterceptor jwtValidationInterceptor = new (idpConfiguration, getOrgResolver(), ["/health", "/api/configurator/apis/generate-configuration", "/api/configurator/apis/generate-k8s-resources", ""]);
final commons:JWTBaseOrgResolver jwtBaseOrgResolver = new;
final PartitionResolver partitionResolver;
final string apkConfSchemaContent = check io:fileReadString(apkSchemaLocation);
listener http:Listener ep0 = new (9443, secureSocket = {
        'key: {
            certFile: <string>keyStores.tls.certFilePath,
            keyFile: <string>keyStores.tls.keyFilePath
        }
    }
);
final runtimeUtil:APKConfValidator apkConfValidator;

isolated function getOrgResolver() returns commons:OrganizationResolver {
    return jwtBaseOrgResolver;
}

# Initializing method for runtime
# + return - Return Error if error occured at initialization.
function init() returns error? {
    apkConfValidator = check runtimeUtil:newAPKConfValidator1(apkConfSchemaContent);
    if partitionServiceConfiguration.enabled {
        partitionResolver = check new PartitionServiceBaseResolver(partitionServiceConfiguration);
    } else {
        partitionResolver = new SinglePartitionResolver();
    }
    log:printInfo("Initializing Configuration Deployer Service..");
}
