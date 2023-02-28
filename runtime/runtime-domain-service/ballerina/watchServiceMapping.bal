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
import wso2/apk_common_lib as commons;

string serviceMappingResourceVersion = "";
isolated map<model:K8sServiceMapping> k8sServiceMappings = {};

class ServiceMappingTask {
    websocket:Client|error watchK8sServiceClient;

    function init(string resourceVersion) {
        self.watchK8sServiceClient = getServiceMappingClient(resourceVersion);
    }

    public function startListening() returns error? {

        worker WatchServiceMappingThread {
            while true {
                do {
                    websocket:Client apiClientResult = check self.watchK8sServiceClient;
                    boolean connectionOpen = apiClientResult.isOpen();
                    if !connectionOpen {
                        log:printDebug("Websocket Client connection closed conectionId: " + apiClientResult.getConnectionId() + " state: " + connectionOpen.toString());
                        self.watchK8sServiceClient = getServiceMappingClient(serviceMappingResourceVersion);
                        websocket:Client retryClient = check self.watchK8sServiceClient;
                        log:printDebug("Reinitializing client..");
                        connectionOpen = retryClient.isOpen();
                        log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + connectionOpen.toString());
                        _ = check self.readServiceMappingEvent(retryClient);
                    } else {
                        _ = check self.readServiceMappingEvent(apiClientResult);
                    }
                } on fail var e {
                    log:printError("Unable to read ServiceMapping messages", e);
                    self.watchK8sServiceClient = getServiceMappingClient(serviceMappingResourceVersion);
                }
            }
        }
    }
    function readServiceMappingEvent(websocket:Client apiWebsocketClient) returns error? {
        boolean connectionOpen = apiWebsocketClient.isOpen();
        log:printDebug("Using Client Connection conectionId: " + apiWebsocketClient.getConnectionId() + " state: " + connectionOpen.toString());
        if !connectionOpen {
            return error("connection closed");
        }
        string message = check apiWebsocketClient->readMessage();
        log:printDebug(message);
        json value = check value:fromJsonString(message);
        string eventType = <string>check value.'type;
        json eventValue = <json>check value.'object;
        json metadata = <json>check eventValue.metadata;
        if eventType == "ERROR" {
            model:Status|error statusEvent = eventValue.cloneWithType(model:Status);
            if (statusEvent is model:Status) {
                _ = check self.handleWatchServiceMappingsGone(statusEvent);
            }
        } else {
            string latestResourceVersion = <string>check metadata.resourceVersion;
            setServiceMappingResourceVersion(latestResourceVersion);
            json clonedEvent = eventValue.cloneReadOnly();
            model:K8sServiceMapping serviceMapping = <model:K8sServiceMapping>clonedEvent;
            if serviceMapping.metadata.namespace == getNameSpace(runtimeConfiguration.apiCreationNamespace) {
                if eventType == "ADDED" {
                    addServiceMapping(serviceMapping);
                } else if (eventType == "MODIFIED") {
                    deleteServiceMapping(serviceMapping);
                    addServiceMapping(serviceMapping);
                } else if (eventType == "DELETED") {
                    deleteServiceMapping(serviceMapping);
                }
            }
        }
    }
    function handleWatchServiceMappingsGone(model:Status statusEvent) returns error? {
        if statusEvent.code == 410 {
            log:printDebug("Re-initializing watch service for ServiceMapping due to cache clear.");
            map<model:K8sServiceMapping> serviceMappingMap = {};
            ServiceClient serviceClient = new ();
            _ = check serviceClient.retrieveAllServiceMappingsAtStartup(serviceMappingMap, ());
            lock {
                k8sServiceMappings = serviceMappingMap.clone();
            }
            self.watchK8sServiceClient = getServiceMappingClient(serviceMappingResourceVersion);
        }
    }
}

public function getServiceMappingClient(string resourceVersion) returns websocket:Client|error {
    log:printDebug("Initializing Watch Service for ServiceMappings with resource Version " + resourceVersion);
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

function setServiceMappingResourceVersion(string resourceVersionValue) {
    serviceMappingResourceVersion = resourceVersionValue;
}

function addServiceMapping(model:K8sServiceMapping serviceMapping) {
    lock {
        k8sServiceMappings[serviceMapping.metadata.uid ?: ""] = serviceMapping.clone();
    }
}

function deleteServiceMapping(model:K8sServiceMapping serviceMapping) {
    lock {
        _ = k8sServiceMappings.remove(serviceMapping.metadata.uid ?: "");
    }
}

isolated function putAllServiceMappings(map<model:K8sServiceMapping> serviceMappings, model:K8sServiceMapping[] events) {
    foreach model:K8sServiceMapping serviceMapping in events {
        serviceMappings[serviceMapping.metadata.uid ?: ""] = serviceMapping.clone();
    }
}

isolated function retrieveAPIMappingsForService(Service serviceEntry, commons:Organization organization) returns model:API[]|commons:APKError {
    lock {
        string[] keys = k8sServiceMappings.keys();
        model:API[] apis = [];
        foreach string key in keys {
            model:K8sServiceMapping serviceMapping = k8sServiceMappings.get(key);
            model:ServiceReference serviceRef = serviceMapping.spec.serviceRef;
            if (serviceRef.name == serviceEntry.name && serviceRef.namespace == serviceEntry.namespace) {
                model:APIReference apiRef = serviceMapping.spec.apiRef;
                model:API? k8sAPI = check getAPIByNameAndNamespace(apiRef.name, apiRef.namespace, organization.clone());
                if k8sAPI is model:API {
                    apis.push(k8sAPI);
                }
            }
        }
        return apis.clone();
    }
}

isolated function retrieveServiceMappingsForAPI(model:API api) returns map<model:K8sServiceMapping> {
    lock {
        map<model:K8sServiceMapping> sortedk8sServiceMappings = {};
        string[] keys = k8sServiceMappings.keys();
        foreach string key in keys {
            model:K8sServiceMapping serviceMapping = k8sServiceMappings.get(key);
            model:APIReference apiRef = serviceMapping.spec.apiRef;
            if (apiRef.name == api.metadata.name && apiRef.namespace == api.metadata.namespace) {
                sortedk8sServiceMappings[<string>serviceMapping.metadata.uid] = serviceMapping;
            }
        }
        return sortedk8sServiceMappings.clone();
    }
}
