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
import wso2/apk_common_lib as commons;

configurable ConfiguratorConfiguration & readonly configuratorConfiguration = {
    keyStores: {
        tls: {
            keyFilePath: "/home/wso2apk/runtime/security/wso2carbon.key"
        }
    }
};




commons:RequestErrorInterceptor requestErrorInterceptor = new;
commons:ResponseErrorInterceptor responseErrorInterceptor = new;
listener http:Listener ep0 = new (9443, secureSocket = {
    'key: {
        certFile: <string>configuratorConfiguration.keyStores.tls.certFilePath,
        keyFile: <string>configuratorConfiguration.keyStores.tls.keyFilePath
    }
},
    interceptors = [requestErrorInterceptor, responseErrorInterceptor]
    );

# Initializing method for runtime
# + return - Return Error if error occured at initialization.
function init() returns error? {
    log:printInfo("Initializing Runtime Domain Service..");
}
