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
import ballerina/websocket;
import ballerina/lang.value;
import ballerina/task;
import runtime_domain_service.model as model;
import ballerina/log;

final websocket:Client apiClient = check new ("wss://" + k8sHost + "/apis/dp.wso2.com/v1alpha1/watch/apis",
auth = {
    token: token
},
secureSocket = {
    cert: caCertPath
});

map<model:K8sAPI> apilist = {};

class APIListingTask {
    *task:Job;

    public function execute() {
        do {
            string|error message = check apiClient->readMessage();
            if message is string {
                json value = check value:fromJsonString(message);
                string eventType = <string>check value.'type;
                json eventValue = <json>check value.'object;
                APIInfo|error apiModel = createAPImodel(eventValue);
                if apiModel is model:K8sAPI {
                    if eventType == "ADDED" {
                        apilist[apiModel.uuid] = apiModel;
                    } else if (eventType == "MODIFIED") {
                        _ = apilist.remove(apiModel.uuid);
                        apilist[apiModel.uuid] = apiModel;
                    } else if (eventType == "DELETED") {
                        _ = apilist.remove(apiModel.uuid);
                    }
                } else {
                    log:printError("error while converting");
                }
            }
        } on fail var e {
            log:printError("Unable to read service messages", e);
        }
    }
}

public function createAPImodel(json event) returns model:K8sAPI|error {
    model:K8sAPI apiInfo = {
        uuid: <string>check event.metadata.uid,
        apiDisplayName: <string>check event.spec.apiDisplayName,
        apiType: <string>check event.spec.apiType,
        apiVersion: <string>check event.spec.apiVersion,
        context: <string>check event.spec.context,
        creationTimestamp: <string>check event.metadata.creationTimestamp,
        definitionFileRef: getValue(event.spec.definitionFileRef),
        sandHTTPRouteRef: getValue(event.spec.sandHTTPRouteRef),
        prodHTTPRouteRef: getValue(event.spec.prodHTTPRouteRef),
        namespace: <string>check event.metadata.namespace,
        k8sName: <string>check event.metadata.name
    };
    return apiInfo;
}

function getValue(json|error value) returns string {
    if value is json {
        return value.toString();
    } else {
        return "";

    }
}

function getAPIs() returns model:K8sAPI[] {
    return apilist.toArray();
}

function getAPI(string id) returns model:K8sAPI|error {
    return check trap apilist.get(id);
}
