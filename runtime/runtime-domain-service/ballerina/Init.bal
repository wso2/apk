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
import ballerina/task;
import ballerina/io;

listener http:Listener ep0 = new (9444);
configurable string namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace";
string namespace = check io:fileReadString(namespaceFile);

configurable RuntimeConfiguratation runtimeConfiguration = {serviceListingNamespaces: [ALL_NAMESPACES], apiCreationNamespace: namespace};

# Initializing method for runtime
isolated function init() {
    do {
        _ = check task:scheduleJobRecurByFrequency(new ServiceTask(), 1);
        _ = check task:scheduleJobRecurByFrequency(new APIListingTask(), 1);
    } on fail var e {
        log:printError("Error initializing Task", e);
    }
    log:printInfo("Initializing Runtime Domain Service..");
}

