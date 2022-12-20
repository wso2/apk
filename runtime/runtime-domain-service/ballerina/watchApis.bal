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

isolated map<model:K8sAPI> apilist = {};
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
                        boolean connectionOpen = apiClientResult.isOpen();
                        if !connectionOpen {
                            log:printDebug("Websocket Client connection closed conectionId: " + apiClientResult.getConnectionId() + " state: " + connectionOpen.toString());
                            apiClient = getClient(resourceVersion);
                            websocket:Client|error|() retryClient = apiClient;
                            if retryClient is websocket:Client {
                                log:printDebug("Reinitializing client..");
                                connectionOpen = retryClient.isOpen();
                                log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + connectionOpen.toString());
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

# Description Retrieve Websocket client for watch API event.
#
# + resourceVersion - resource Version to watch after.
# + return - Return websocket Client.
public function getClient(string resourceVersion) returns websocket:Client|error {
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

public isolated function createAPImodel(json event) returns model:K8sAPI|error {
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

isolated function getValue(json|error value) returns string {
    if value is json {
        return value.toString();
    } else {
        return "";

    }
}

isolated function getAPIs() returns model:K8sAPI[] {
    lock {
        model:K8sAPI[] & readonly readOnlyAPIList = apilist.toArray().cloneReadOnly();
        return readOnlyAPIList;
    }
}

isolated function getAPI(string id) returns model:K8sAPI|error {
    lock {
        map<model:K8sAPI> & readonly readOnlyAPIMap = apilist.cloneReadOnly();
        return check trap readOnlyAPIMap.get(id);
    }
}

function putallAPIS(json[] apiData) {
    foreach json api in apiData {
        model:K8sAPI|error k8sAPI = createAPImodel(api);
        if k8sAPI is model:K8sAPI {
            lock {
                apilist[k8sAPI.uuid] = k8sAPI.clone();
            }
        }
    }
}

function setResourceVersion(string resourceVersionValue) {
    resourceVersion = resourceVersionValue;
}

function readAPIEvent(websocket:Client apiWebsocketClient) returns error? {
    boolean connectionOpen = apiWebsocketClient.isOpen();

    log:printDebug("Using Client Connection conectionId: " + apiWebsocketClient.getConnectionId() + " state: " + connectionOpen.toString());
    if !connectionOpen {
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
                    lock {
                        apilist[apiModel.uuid] = apiModel.clone();
                    }
                } else if (eventType == "MODIFIED") {
                    lock {
                        _ = apilist.remove(apiModel.uuid);
                        apilist[apiModel.uuid] = apiModel.clone();
                    }
                } else if (eventType == "DELETED") {
                    lock {
                        _ = apilist.remove(apiModel.uuid);
                    }
                }
            }
        } else {
            log:printError("error while converting");
        }
    } else {
        log:printError("error while reading message", message);
    }

}

isolated function getAPIByNameAndNamespace(string name, string namespace) returns model:K8sAPI|() {
    foreach model:K8sAPI api in getAPIs() {
        if (api.k8sName == name && api.namespace == namespace) {
            return api;
        }
    }
    json|error k8sAPIByNameAndNamespace = getK8sAPIByNameAndNamespace(name, namespace);
    if k8sAPIByNameAndNamespace is json {
        model:K8sAPI|error k8sAPI = createAPImodel(k8sAPIByNameAndNamespace);
        if k8sAPI is model:K8sAPI {
            return k8sAPI;
        } else {
            log:printError("Error occued while converting json", k8sAPI);
        }
    }
    return ();
}
