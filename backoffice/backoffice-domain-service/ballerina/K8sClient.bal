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

import ballerina/io;
import ballerina/http;
import ballerina/log;
import wso2/apk_common_lib as commons;

const string K8S_API_ENDPOINT = "/api/v1";
final string token = check io:fileReadString(k8sConfig.serviceAccountPath + "/token");
final string caCertPath = k8sConfig.serviceAccountPath + "/ca.crt";
string namespaceFile = k8sConfig.serviceAccountPath + "/namespace";
final string currentNameSpace = check io:fileReadString(namespaceFile);
final http:Client k8sApiServerEp = check initializeK8sClient();

# This initialize the k8s Client.
# + return - k8s http client
public function initializeK8sClient() returns http:Client|error {
    http:Client k8sApiClient = check new ("https://" + k8sConfig.host,
    auth = {
        token: token
    },
        secureSocket = {
        cert: caCertPath

    }
    );
    return k8sApiClient;
}

# This returns Pod value according to given name and namespace.
#
# + name - Name of Pod  
# + namespace - Namespace of Pod
# + return - Return Pod value for a given name and namespace
isolated function getPodFromNameAndNamespace(string name, string namespace) returns string[]| commons:APKError {
    string endpoint = "/api/v1/namespaces/" + namespace + "/endpoints/" + name;
    http:Response|error response = k8sApiServerEp->get(endpoint, targetType = http:Response);
    if response is http:Response {
        json|http:ClientError podValue = response.getJsonPayload();
        if podValue is json{
            do {
                log:printDebug(podValue.toString());
                json[] subsets = check <json[]|error>podValue.subsets;
                json[] addresses = check <json[]|error>subsets[0].addresses;
                string[] hosts =[];
                foreach json item in addresses {
                    string ip = check item.ip;
                    hosts.push(ip);
                }
                log:printDebug(hosts.toString());
                return hosts;
            } on fail var e {
                string message ="Error while retrieving host. Error while retrieving pod information for pod: " + name;
                log:printError(message + e.toBalString());
                return error(message,e, message = message, description = message, code = 909000, statusCode = 500);
            }
        } else {
            string message ="Response isn't a json. Error while retrieving pod information for pod: " + name;
            log:printError(message);
            return error(message, message = message, description = message, code = 909000, statusCode = 500);
        }
    } else {
        string message ="Error while retrieving pod information for pod" + name;
        log:printError(message);
        return error(message, message = message, description = message, code = 909000, statusCode = 500);
    }
}
