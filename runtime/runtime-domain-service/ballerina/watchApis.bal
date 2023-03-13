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
import wso2/apk_common_lib as commons;
import ballerina/log;
import ballerina/url;

isolated map<map<model:API>> apilist = {};
string apiResourceVersion = "";

class APIListingTask {
    websocket:Client|error watchAPIService;

    function init(string apiResourceVersion) {
        self.watchAPIService = getAPIClient(apiResourceVersion);
    }

    public function startListening() returns error? {

        worker WatchAPIThread {
            while true {
                do {
                    websocket:Client apiClientResult = check self.watchAPIService;
                    boolean connectionOpen = apiClientResult.isOpen();
                    if !connectionOpen {
                        log:printDebug("Websocket Client connection closed conectionId: " + apiClientResult.getConnectionId() + " state: " + connectionOpen.toString());
                        self.watchAPIService = getAPIClient(apiResourceVersion);
                        websocket:Client retryClient = check self.watchAPIService;
                        log:printDebug("Reinitializing client..");
                        connectionOpen = retryClient.isOpen();
                        log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + connectionOpen.toString());
                        _ = check self.readAPIEvent(retryClient);
                    } else {
                        _ = check self.readAPIEvent(apiClientResult);
                    }
                } on fail var e {
                    log:printError("Unable to read api messages", e);
                    self.watchAPIService = getAPIClient(apiResourceVersion);
                }
            }
        }
    }
    function handleWatchAPIGone(model:Status statusEvent) returns error? {
        if statusEvent.code == 410 {
            log:printDebug("Re-initializing watch service for API due to cache clear.");
            map<map<model:API>> orgApiMap = {};
            APIClient apiClient = new ();
            _ = check apiClient.retrieveAllApisAtStartup(orgApiMap, ());
            lock {
                apilist = orgApiMap.clone();
            }
            self.watchAPIService = getAPIClient(apiResourceVersion);
        }
    }
    function readAPIEvent(websocket:Client apiWebsocketClient) returns error? {
        boolean connectionOpen = apiWebsocketClient.isOpen();

        log:printDebug("Using Client Connection conectionId: " + apiWebsocketClient.getConnectionId() + " state: " + connectionOpen.toString());
        if !connectionOpen {
            error err = error("connection closed");
            return err;
        }
        string message = check apiWebsocketClient->readMessage();
        log:printDebug(message);
        json value = check value:fromJsonString(message);
        string eventType = <string>check value.'type;
        map<json> eventValue = <map<json>>check value.'object;
        json metadata = <json>check eventValue.metadata;
        if eventType == "ERROR" {
            model:Status|error statusEvent = eventValue.cloneWithType(model:Status);
            if (statusEvent is model:Status) {
                _ = check self.handleWatchAPIGone(statusEvent);
            }
        } else {
            string latestResourceVersion = <string>check metadata.resourceVersion;
            setResourceVersion(latestResourceVersion);
            model:API|error apiModel = eventValue.cloneWithType(model:API);
            if apiModel is model:API {
                if (apiModel.metadata.namespace == getNameSpace(runtimeConfiguration.apiCreationNamespace)) {
                    if eventType == "ADDED" {
                        lock {
                            putAPI(apiModel.clone());
                        }
                    } else if (eventType == "MODIFIED") {
                        lock {
                            updateAPI(apiModel.clone());
                        }
                    } else if (eventType == "DELETED") {
                        lock {
                            removeAPI(apiModel);
                        }
                    }
                }
            } else {
                log:printError("error while converting");
            }
        }

    }
}

# Description Retrieve Websocket client for watch API event.
#
# + resourceVersion - resource Version to watch after.
# + return - Return websocket Client.
public function getAPIClient(string resourceVersion) returns websocket:Client|error {
    log:printDebug("Initializing Watch Service for APIS with resource Version " + resourceVersion);
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/apis/dp.wso2.com/v1alpha1/watch/apis";
    boolean questionMark = false;
    if resourceVersion.length() > 0 {
        requestURl = requestURl + "?resourceVersion=" + resourceVersion.toString();
        questionMark = true;
    }
    string fieldSelectorQuery = "metadata.namespace=" + currentNameSpace;

    if questionMark {
        requestURl = requestURl + "&fieldSelector=" + check url:encode(fieldSelectorQuery, "UTF-8");
    } else {
        requestURl = requestURl + "?fieldSelector=" + check url:encode(fieldSelectorQuery, "UTF-8");
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

isolated function getAPIs(commons:Organization organization) returns model:API[] {
    lock {
        map<model:API>|error & readonly readOnlyAPImap = trap apilist.get(organization.uuid).cloneReadOnly();
        if readOnlyAPImap is map<model:API> & readonly {
            return readOnlyAPImap.toArray();
        } else {
            return [];
        }
    }
}

isolated function getAPI(string id, commons:Organization organization) returns model:API|error {
    lock {
        map<model:API> & readonly apiMap = check trap apilist.get(organization.uuid).cloneReadOnly();
        return check trap apiMap.get(id);
    }
}

isolated function putallAPIS(map<map<model:API>> orgApiMap, model:API[] apiData) {
    foreach model:API api in apiData {
        boolean systemAPI = api.spec.systemAPI ?: false;
        if !systemAPI {
            lock {
                map<model:API>|error orgmap = trap orgApiMap.get(api.spec.organization);
                if orgmap is map<model:API> {
                    orgmap[<string>api.metadata.uid] = api.clone();
                } else {
                    map<model:API> apiMap = {};
                    apiMap[<string>api.metadata.uid] = api.clone();
                    orgApiMap[api.spec.organization] = apiMap;
                }
            }
        }

    }
}

function setResourceVersion(string resourceVersionValue) {
    apiResourceVersion = resourceVersionValue;
}

isolated function putAPI(model:API api) {
    boolean systemAPI = api.spec.systemAPI ?: false;
    if !systemAPI {
        lock {
            map<model:API>|error orgapiMap = trap apilist.get(api.spec.organization);
            if orgapiMap is map<model:API> {
                orgapiMap[<string>api.metadata.uid] = api.clone();
            } else {
                map<model:API> apiMap = {};
                apiMap[<string>api.metadata.uid] = api.clone();
                apilist[api.spec.organization] = apiMap;
            }
        }
    }
}

isolated function updateAPI(model:API api) {
    removeAPI(api);
    putAPI(api);
}

isolated function removeAPI(model:API api) {
    lock {
        map<model:API>|error orgapiMap = trap apilist.get(api.spec.organization);
        if orgapiMap is map<model:API> {
            _ = orgapiMap.remove(<string>api.metadata.uid);
        }
    }
}

isolated function getAPIByNameAndNamespace(string name, string namespace, commons:Organization organization) returns model:API|()|commons:APKError {
    foreach model:API api in getAPIs(organization) {
        if (api.metadata.name == name && api.metadata.namespace == namespace) {
            return api;
        }
    }
    model:API? k8sAPIByNameAndNamespace = check getK8sAPIByNameAndNamespace(name, namespace);
    if k8sAPIByNameAndNamespace is model:API {
        return k8sAPIByNameAndNamespace;
    }
    return ();
}

isolated function isAPIVersionExist(string name, string 'newVersion, commons:Organization organization) returns boolean {
    lock {
        map<model:API>|error apiMap = trap apilist.get(organization.uuid);
        if apiMap is map<model:API> {
            model:API[] & readonly readOnlyAPIList = apiMap.toArray().cloneReadOnly();
            foreach model:API & readonly api in readOnlyAPIList {
                if api.spec.apiDisplayName == name && api.spec.apiVersion == 'newVersion {
                    return true;
                }
            }
        }
    }
    return false;
}
