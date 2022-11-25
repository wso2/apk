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
import ballerina/http;

const string K8S_API_ENDPOINT = "/api/v1";
final http:Client k8sApiServerEp = check initializeK8sClient();
configurable string k8sHost = "kubernetes.default";
configurable string saTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token";
string token = check io:fileReadString(saTokenPath);
configurable string caCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt";

# This initialize the k8s Client.
# + return - k8s http client
function initializeK8sClient() returns http:Client|error {
    http:Client k8sApiClient = check new ("https://" + k8sHost,
    auth = {
        token: token
    },
        secureSocket = {
            cert: caCertPath

        }
    );
    return k8sApiClient;
}

# This returns ConfigMap value according to name and namespace.
#
# + name - Name of ConfigMap
# + namespace - Namespace of Configmap
function getConfigMapValueFromNameAndNamespace(string name, string namespace) returns json|error {
    string endpoint = "/api/v1/namespaces/" + namespace + "/configmaps/" + name;
    return k8sApiServerEp->get(endpoint, targetType = json);
}
