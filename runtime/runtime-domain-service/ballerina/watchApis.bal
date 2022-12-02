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
import runtime_domain_service.model as model;
import ballerina/log;

map<model:K8sAPI> apilist = {};
string resourceVersion = "";
websocket:Client|error|() apiClient = ();

class APIListingTask {
    function init(string resourceVersion) {
        apiClient = getClient(resourceVersion);
    }

    public function startListening() returns error? {

        worker WatchAPIThread {
            while true {
                do {
                    websocket:Client|error|() apiClientResult = apiClient;
                    if apiClientResult is websocket:Client {
                        if !apiClientResult.isOpen() {
                            log:printDebug("Websocket Client connection closed conectionId: " + apiClientResult.getConnectionId() + " state: " + apiClientResult.isOpen().toString());
                            apiClient = getClient(resourceVersion);
                            websocket:Client|error|() retryClient = apiClient;
                            if retryClient is websocket:Client {
                                log:printDebug("Reinitializing client..");
                                log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + apiClientResult.isOpen().toString());
                                _ = check readAPIEvent(retryClient);
                            } else if retryClient is error {
                                log:printError("error while reading message", retryClient);
                            }
                        } else {
                            _ = check readAPIEvent(apiClientResult);
                        }

                    } else if apiClientResult is error {
                        log:printError("error while reading message", apiClientResult);
                    }
                } on fail var e {
                    log:printError("Unable to read api messages", e);
                }
            }
        }
    }
}

isolated function getClient(string resourceVersion) returns websocket:Client|error|() {
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/apis/dp.wso2.com/v1alpha1/watch/apis";
    if resourceVersion.length() > 0 {
        requestURl = requestURl + "?resourceVersion=" + resourceVersion.toString();
    }
    return new (requestURl,
    auth = {
        token: token
    },
        secureSocket = {
            cert: caCertPath
        }
    );
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

function putallAPIS(json[] apiData) {
    foreach json api in apiData {
        model:K8sAPI|error k8sAPI = createAPImodel(api);
        if k8sAPI is model:K8sAPI {
            apilist[k8sAPI.uuid] = k8sAPI;
        }
    }
}

function setResourceVersion(string resourceVersionValue) {
    resourceVersion = resourceVersionValue;
}

function readAPIEvent(websocket:Client apiWebsocketClient) returns error? {
    log:printDebug("Using Client Connection conectionId: " + apiWebsocketClient.getConnectionId() + " state: " + apiWebsocketClient.isOpen().toString());
    if !apiWebsocketClient.isOpen() {
        error err = error("connection closed");
        return err;
    }
    string|error message = check apiWebsocketClient->readMessage();
    if message is string {
        log:printDebug(message);
        json value = check value:fromJsonString(message);
        string eventType = <string>check value.'type;
        json eventValue = <json>check value.'object;
        json metadata = <json>check eventValue.metadata;
        string latestResourceVersion = <string>check metadata.resourceVersion;
        setResourceVersion(latestResourceVersion);
        APIInfo|error apiModel = createAPImodel(eventValue);
        if apiModel is model:K8sAPI {
            if apiModel.namespace == getNameSpace(runtimeConfiguration.apiCreationNamespace) {
                if eventType == "ADDED" {
                    apilist[apiModel.uuid] = apiModel;
                } else if (eventType == "MODIFIED") {
                    _ = apilist.remove(apiModel.uuid);
                    apilist[apiModel.uuid] = apiModel;
                } else if (eventType == "DELETED") {
                    _ = apilist.remove(apiModel.uuid);
                }
            }
        } else {
            log:printError("error while converting");
        }
    } else {
        log:printError("error while reading message", message);
    }

}
