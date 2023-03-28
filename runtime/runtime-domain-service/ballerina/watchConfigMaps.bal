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
import ballerina/websocket;
import ballerina/lang.value;
import runtime_domain_service.model as model;
import ballerina/log;
import ballerina/url;
import ballerina/http;

isolated map<model:ConfigMap> configMapList = {};
string configMapResourceVersion = "";

class ConfigMapListingTask {
    private websocket:Client|error watchconfigMapService;

    function init(string configMapResourceVersion) {
        self.watchconfigMapService = getConfigMapWatchClient(configMapResourceVersion);
    }

    public function startListening() returns error? {

        worker WatchConfigMapThread {
            while true {
                do {
                    websocket:Client configMapClientResult = check self.watchconfigMapService;
                    boolean connectionOpen = configMapClientResult.isOpen();
                    if !connectionOpen {
                        log:printDebug("Websocket Client connection closed conectionId: " + configMapClientResult.getConnectionId() + " state: " + connectionOpen.toString());
                        lock {
                            self.watchconfigMapService = getConfigMapWatchClient(configMapResourceVersion);
                        }
                        websocket:Client retryClient = check self.watchconfigMapService;
                        log:printDebug("Reinitializing client..");
                        connectionOpen = retryClient.isOpen();
                        log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + connectionOpen.toString());
                        _ = check self.readConfigMapEvent(retryClient);
                    } else {
                        _ = check self.readConfigMapEvent(configMapClientResult);
                    }
                } on fail var e {
                    log:printError("Unable to read configmap messages", e);
                    lock {
                        self.watchconfigMapService = getConfigMapWatchClient(configMapResourceVersion);
                    }
                }
            }
        }
    }
    function readConfigMapEvent(websocket:Client configMapWebSocketClient) returns error? {
        boolean connectionOpen = configMapWebSocketClient.isOpen();

        log:printDebug("Using Client Connection conectionId: " + configMapWebSocketClient.getConnectionId() + " state: " + connectionOpen.toString());
        if !connectionOpen {
            error err = error("connection closed");
            return err;
        }
        string message = check configMapWebSocketClient->readMessage();
        log:printDebug(message);
        json value = check value:fromJsonString(message);
        string eventType = <string>check value.'type;
        json eventValue = <json>check value.'object;
        json metadata = <json>check eventValue.metadata;
        if eventType == "ERROR" {
            model:Status|error statusEvent = eventValue.cloneWithType(model:Status);
            if (statusEvent is model:Status) {
                _ = check self.handleConfigMapEventsGone(statusEvent);
            }
        } else {
            string latestResourceVersion = <string>check metadata.resourceVersion;
            setConfigmapResourceVersion(latestResourceVersion);
            model:ConfigMap|error configMap = eventValue.cloneWithType(model:ConfigMap);
            if configMap is model:ConfigMap {
                if eventType == "ADDED" {
                    lock {
                        putConfigMapMetaData(configMap);
                    }
                } else if (eventType == "MODIFIED") {
                    lock {
                        updateConfigMapMetaData(configMap);
                    }
                } else if (eventType == "DELETED") {
                    lock {
                        removeConfigmapMetaData(configMap);
                    }
                }
            } else {
                log:printError("error while converting Configmap event", configMap);
            }
        }
    }

    function handleConfigMapEventsGone(model:Status statusEvent) returns error? {
        if statusEvent.code == 410 {
            log:printDebug("Re-initializing watch service for API due to cache clear.");
            map<model:ConfigMap> configMap = {};
            _ = check retrieveAllConfigMapsAtStartup(configMap, ());
            lock {
                configMapList = configMap.clone();
            }
            self.watchconfigMapService = getConfigMapWatchClient(configMapResourceVersion);
        }
    }

}

public function retrieveAllConfigMapsAtStartup(map<model:ConfigMap>? configMap, string? continueValue) returns error? {
    string? resultValue = continueValue;
    model:ConfigMapList|http:ClientError retrieveAllConfigMapsResult;
    if resultValue is string {
        retrieveAllConfigMapsResult = check retrieveAllconfigMaps(resultValue);
    } else {
        retrieveAllConfigMapsResult = check retrieveAllconfigMaps(());
    }

    if retrieveAllConfigMapsResult is model:ConfigMapList {
        model:ListMeta metadata = retrieveAllConfigMapsResult.metadata;
        model:ConfigMap[] configMaps = retrieveAllConfigMapsResult.items;
        if configMap is map<model:ConfigMap> {
            lock {
                putAllConfigMaps(configMap, configMaps.clone());
            }
        } else {
            lock {
                putAllConfigMaps(configMapList, configMaps.clone());
            }
        }
        string? continueElement = metadata.'continue;
        if continueElement is string {
            if continueElement.length() > 0 {
                _ = check retrieveAllConfigMapsAtStartup(configMap, continueElement);
            }
        }
        string? resourceVersion = metadata.'resourceVersion;
        if resourceVersion is string {
            setConfigmapResourceVersion(resourceVersion);
        }
    }
}

# Description Retrieve Websocket client for watch API event.
#
# + resourceVersion - resource Version to watch after.
# + return - Return websocket Client.
public function getConfigMapWatchClient(string resourceVersion) returns websocket:Client|error {
    log:printDebug("Initializing Watch Service for APIS with resource Version " + resourceVersion);
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/api/v1/watch/configmaps?fieldSelector=" + check getEncodedStringForNamespaceSelector() +
    "&labelSelector=" + check getEncodedStringForLabelSelector();
    if resourceVersion.length() > 0 {
        requestURl = requestURl + "&resourceVersion=" + resourceVersion.toString();
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

isolated function getConfigMap(string uuid) returns model:ConfigMap|error? {
    lock {
        model:ConfigMap|error configMap = trap configMapList.get(uuid);
        if configMap is model:ConfigMap {
            return configMap.clone();
        } else {
            return error("ConfigMap not found");
        }
    }
}

isolated function putAllConfigMaps(map<model:ConfigMap> configMapArray, model:ConfigMap[] configMaps) {
    foreach model:ConfigMap configMap in configMaps {
        lock {
            configMapArray[<string>configMap.metadata.uid] = configMap.clone();
        }
    }
}

function setConfigmapResourceVersion(string resourceVersion) {
    lock {
        configMapResourceVersion = resourceVersion;
    }
}

isolated function putConfigMapMetaData(model:ConfigMap configMap) {
    lock {
        configMapList[<string>configMap.metadata.uid] = configMap.clone();
    }
}

isolated function updateConfigMapMetaData(model:ConfigMap configMap) {
    removeConfigmapMetaData(configMap);
    putConfigMapMetaData(configMap);
}

isolated function removeConfigmapMetaData(model:ConfigMap configMap) {
    lock {
        _ = configMapList.remove(<string>configMap.metadata.uid);
    }
}

isolated function getEncodedStringForNamespaceSelector() returns string|error {
    string fieldSelectorQuery = "metadata.namespace=" + currentNameSpace;
    return check url:encode(fieldSelectorQuery, "UTF-8");
}

isolated function getEncodedStringForLabelSelector() returns string|error {
    string labelSelectorQuery = MANAGED_BY_HASH_LABEL + "=" + MANAGED_BY_HASH_LABEL_VALUE;
    return check url:encode(labelSelectorQuery, "UTF-8");
}
