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

isolated map<model:Organization> organizationList = {};
string organizationResourceVersion = "";

class OrganizationListingTask {
    websocket:Client|error watchOrganizationService;

    function init(string organizationResourceVersion) {
        self.watchOrganizationService = getOrganizationWatchClient(organizationResourceVersion);
    }

    public function startListening() returns error? {

        worker WatchOrganizationThread {
            while true {
                do {
                    websocket:Client orgClientResult = check self.watchOrganizationService;
                    boolean connectionOpen = orgClientResult.isOpen();
                    if !connectionOpen {
                        log:printDebug("Websocket Client connection closed conectionId: " + orgClientResult.getConnectionId() + " state: " + connectionOpen.toString());
                        self.watchOrganizationService = getOrganizationWatchClient(organizationResourceVersion);
                        websocket:Client retryClient = check self.watchOrganizationService;
                        log:printDebug("Reinitializing client..");
                        connectionOpen = retryClient.isOpen();
                        log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + connectionOpen.toString());
                        _ = check self.readOrganizationEvent(retryClient);
                    } else {
                        _ = check self.readOrganizationEvent(orgClientResult);
                    }
                } on fail var e {
                    log:printError("Unable to read organization messages", e);
                    self.watchOrganizationService = getOrganizationWatchClient(organizationResourceVersion);
                }
            }
        }
    }
    function readOrganizationEvent(websocket:Client organizationWebSocketClient) returns error? {
        boolean connectionOpen = organizationWebSocketClient.isOpen();

        log:printDebug("Using Client Connection conectionId: " + organizationWebSocketClient.getConnectionId() + " state: " + connectionOpen.toString());
        if !connectionOpen {
            error err = error("connection closed");
            return err;
        }
        string message = check organizationWebSocketClient->readMessage();
        log:printDebug(message);
        json value = check value:fromJsonString(message);
        string eventType = <string>check value.'type;
        json eventValue = <json>check value.'object;
        json metadata = <json>check eventValue.metadata;
        if eventType == "ERROR" {
            model:Status|error statusEvent = eventValue.cloneWithType(model:Status);
            if (statusEvent is model:Status) {
                _ = check self.handleOganizationWatchGone(statusEvent);
            }
        } else {
            string latestResourceVersion = <string>check metadata.resourceVersion;
            setOrganizationResourceVersion(latestResourceVersion);
            model:Organization|error organization = eventValue.cloneWithType(model:Organization);
            if organization is model:Organization {
                if (organization.metadata.namespace == getNameSpace(runtimeConfiguration.apiCreationNamespace)) {
                    if eventType == "ADDED" {
                        lock {
                            putOrganization(organization);
                        }
                    } else if (eventType == "MODIFIED") {
                        lock {
                            updateOrganization(organization);
                        }
                    } else if (eventType == "DELETED") {
                        lock {
                            removeOrganization(organization);
                        }
                    }
                }
            } else {
                log:printError("error while converting organization event", organization);
            }
        }
    }

    function handleOganizationWatchGone(model:Status statusEvent) returns error? {
        if statusEvent.code == 410 {
            log:printDebug("Re-initializing watch service for API due to cache clear.");
            map<model:Organization> organizationsMap = {};
            OrgClient orgClient = new ();
            _ = check orgClient.retrieveAllOrganizationsAtStartup(organizationsMap, ());
            lock {
                organizationList = organizationsMap.clone();
            }
            self.watchOrganizationService = getOrganizationWatchClient(organizationResourceVersion);
        }
    }
}

# Description Retrieve Websocket client for watch API event.
#
# + resourceVersion - resource Version to watch after.
# + return - Return websocket Client.
public function getOrganizationWatchClient(string resourceVersion) returns websocket:Client|error {
    log:printDebug("Initializing Watch Service for APIS with resource Version " + resourceVersion);
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/apis/cp.wso2.com/v1alpha1/watch/organizations";
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

isolated function getOrganization(string organization) returns model:Organization|() {
    lock {
        model:Organization|error & readonly readOnlyOrganization = trap organizationList.get(organization).cloneReadOnly();
        if readOnlyOrganization is model:Organization {
            return readOnlyOrganization.clone();
        } else {
            return ();
        }
    }
}

isolated function putAllOrganizations(map<model:Organization> organizationMap, model:Organization[] organizations) {
    foreach model:Organization organization in organizations {
        lock {
            organizationMap[organization.spec.uuid] = organization.clone();
        }
    }
}

function setOrganizationResourceVersion(string resourceVersion) {
    organizationResourceVersion = resourceVersion;
}

isolated function putOrganization(model:Organization organization) {
    lock {
        organizationList[organization.spec.uuid] = organization.clone();
    }
}

isolated function updateOrganization(model:Organization organization) {
    removeOrganization(organization);
    putOrganization(organization);
}

isolated function removeOrganization(model:Organization organization) {
    lock {
        _ = organizationList.remove(organization.spec.uuid);
    }
}

