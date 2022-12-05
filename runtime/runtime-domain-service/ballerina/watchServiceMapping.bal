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

map<map<model:K8sAPI>> serviceMappings = {};
string serviceMappingResourceVersion = "";
websocket:Client|error|() serviceMappingClient = ();

class ServiceMappingTask {
    function init(string resourceVersion) {
        serviceMappingClient = getServiceMappingClient(resourceVersion);
    }

    public function startListening() returns error? {

        worker WatchServiceMappingThread {
            while true {
                do {
                    websocket:Client|error|() apiClientResult = serviceMappingClient;
                    if apiClientResult is websocket:Client {
                        if !apiClientResult.isOpen() {
                            log:printDebug("Websocket Client connection closed conectionId: " + apiClientResult.getConnectionId() + " state: " + apiClientResult.isOpen().toString());
                            serviceMappingClient = getServiceMappingClient(resourceVersion);
                            websocket:Client|error|() retryClient = serviceMappingClient;
                            if retryClient is websocket:Client {
                                log:printDebug("Reinitializing client..");
                                log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + retryClient.isOpen().toString());
                                _ = check readServiceMappingEvent(retryClient);
                            } else if retryClient is error {
                                log:printError("error while reading message", retryClient);
                            }
                        } else {
                            _ = check readServiceMappingEvent(apiClientResult);
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

function getServiceMappingClient(string resourceVersion) returns websocket:Client|error|() {
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/apis/dp.wso2.com/v1alpha1/watch/servicemappings";
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

function putallServiceMappings(json[] apiData) {
    foreach json api in apiData {
        model:K8sAPI|error k8sAPI = createAPImodel(api);
        if k8sAPI is model:K8sAPI {
            apilist[k8sAPI.uuid] = k8sAPI;
        }
    }
}

function setServiceMappingResourceVersion(string resourceVersionValue) {
    serviceMappingResourceVersion = resourceVersionValue;
}

function readServiceMappingEvent(websocket:Client apiWebsocketClient) returns error? {
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
        setServiceMappingResourceVersion(latestResourceVersion);
        json clonedEvent = eventValue.cloneReadOnly();
        model:K8sServiceMapping|error serviceMapping = <model:K8sServiceMapping>clonedEvent;
        if serviceMapping is model:K8sServiceMapping {
            if serviceMapping.metadata.namespace == getNameSpace(runtimeConfiguration.apiCreationNamespace) {
                if eventType == "ADDED" {
                    addServiceMapping(serviceMappings, serviceMapping);
                } else if (eventType == "MODIFIED") {
                    deleteServiceMapping(serviceMappings, serviceMapping);
                    addServiceMapping(serviceMappings, serviceMapping);
                } else if (eventType == "DELETED") {
                    deleteServiceMapping(serviceMappings, serviceMapping);
                }
            }
        } else {
            log:printError("error while converting");
        }
    } else {
        log:printError("error while reading message", message);
    }

}

function addServiceMapping(map<map<model:K8sAPI>> serviceMappings, model:K8sServiceMapping serviceMapping) {
    model:ServiceReference serviceRef = serviceMapping.spec.serviceRef;
    Service? serviceResult = getService(serviceRef.name, serviceRef.namespace);
    if serviceResult is Service {
        map<model:K8sAPI>? apiList = serviceMappings[serviceResult.id];
        map<model:K8sAPI> apis;
        if apiList is map<model:K8sAPI> {
            apis = apilist;
        } else {
            apis = {};
        }
        serviceMappings[serviceResult.id] = apis;

        model:APIReference apiRef = serviceMapping.spec.apiRef;
        model:K8sAPI? api = getAPIByNameAndNamespace(apiRef.name, apiRef.namespace);
        if api is model:K8sAPI {
            apis[api.uuid] = api;
        }
    }
}

function deleteServiceMapping(map<map<model:K8sAPI>> serviceMappings, model:K8sServiceMapping serviceMapping) {
    model:ServiceReference serviceRef = serviceMapping.spec.serviceRef;
    Service? serviceResult = getService(serviceRef.name, serviceRef.namespace);
    if serviceResult is Service {
        map<model:K8sAPI>? apiList = serviceMappings[serviceResult.id];
        map<model:K8sAPI> apis;
        if apiList is map<model:K8sAPI> {
            apis = apilist;
        } else {
            apis = {};
        }
        serviceMappings[serviceResult.id] = apis;
        model:APIReference apiRef = serviceMapping.spec.apiRef;
        model:K8sAPI? api = getAPIByNameAndNamespace(apiRef.name, apiRef.namespace);
        if api is model:K8sAPI {
            _ = apis.remove(api.uuid);
        }
    }
}

function putAllServiceMappings(json[] events) returns error? {
    foreach json event in events {
        json clonedEvent = event.cloneReadOnly();
        model:K8sServiceMapping|error serviceMapping = trap <model:K8sServiceMapping>clonedEvent;
        if serviceMapping is model:K8sServiceMapping {
            addServiceMapping(serviceMappings, serviceMapping);
        }
    }
}
