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
import ballerina/log;

map<Service> services = {};
string servicesResourceVersion = "";
websocket:Client|error|() servicesClient = ();

class ServiceTask {

    function init(string resourceVersion) {
        servicesClient = getServiceClient(servicesResourceVersion);
    }
    public function startListening() {

        worker WatchServices {
            while true {
                do {
                    websocket:Client|error|() serviceClientResult = servicesClient;
                    if serviceClientResult is websocket:Client {
                        if !serviceClientResult.isOpen() {
                            log:printDebug("ServiceWebsocket Client connection closed conectionId: " + serviceClientResult.getConnectionId());
                            servicesClient = getServiceClient(servicesResourceVersion);
                            websocket:Client|error|() retryClient = servicesClient;
                            if retryClient is websocket:Client {
                                log:printDebug("Reinitializing client..");
                                log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + retryClient.isOpen().toString());
                                _ = check readServiceEvents(retryClient);
                            } else if retryClient is error {
                                log:printError("error while reading message", retryClient);
                            }
                        } else {
                            log:printDebug("Intializd new Client Connection conectionId: " + serviceClientResult.getConnectionId() + " state: " + serviceClientResult.isOpen().toString());
                            _ = check readServiceEvents(serviceClientResult);
                        }

                    } else if serviceClientResult is error {
                        log:printError("error while reading message", serviceClientResult);
                    }
                } on fail var e {
                    log:printError("Unable to read services messages", e);
                }
            }
        }

    }

}

function containsNamespace(string namespace) returns boolean {
    foreach string name in runtimeConfiguration.serviceListingNamespaces {
        if (name == ALL_NAMESPACES || name == namespace) {
            return true;
        }
    }
    return false;
}

public function createServiceModel(json event) returns Service|error {
    Service serviceData = {
        id: <string>check event.metadata.uid,
        name: <string>check event.metadata.name,
        namespace: <string>check event.metadata.namespace,
        'type: <string>check event.spec.'type,
        portmapping: check mapPortMapping(event),
        createdTime: <string>check event.metadata.creationTimestamp
    };
    return serviceData;
}

function mapPortMapping(json event) returns PortMapping[]|error {
    json[] ports = <json[]>check event.spec.ports;
    PortMapping[] portmappings = [];

    foreach json port in ports {
        PortMapping portmapping =
            {
            name: check port.name,
            protocol: check port.protocol,
            port: check port.port,
            targetport: check port.targetPort
        };
        portmappings.push(portmapping);
    }
    return portmappings;
}

function getServicesList() returns Service[] {
    return services.toArray();
}

# This retrieve specific service from name space.
#
# + name - name of service.
# + namespace - namespace of service.
# + return - service in namespace.
function getService(string name, string namespace) returns Service? {
    foreach Service s in getServicesList() {
        if (s.name == name && s.namespace == namespace) {
            return s;
        }
    }
    Service|error retrieveK8sServiceMapping = new ServiceClient().retrieveK8sServiceMapping(name, namespace);
    if retrieveK8sServiceMapping is Service {
        return retrieveK8sServiceMapping;
    }
    return;
}

function grtServiceById(string id) returns Service|error {
    return trap services.get(id);
}

function putAllServices(json[] servicesEntries) {
    foreach json serviceData in servicesEntries {
        Service|error serviceEntry = createServiceModel(serviceData);
        if serviceEntry is Service {
            services[serviceEntry.id] = serviceEntry;
        }
    }
}

function setServicesResourceVersion(string resourceVersionValue) {
    servicesResourceVersion = resourceVersionValue;
}

public isolated function getServiceClient(string resourceVersion) returns websocket:Client|error|() {
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/api/v1/watch/services";
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

function readServiceEvents(websocket:Client serviceWebSocketClient) returns error? {
    log:printDebug("Using Client Connection conectionId: " + serviceWebSocketClient.getConnectionId() + " state: " + serviceWebSocketClient.isOpen().toString());
    if !serviceWebSocketClient.isOpen() {
        error err = error("connection closed");
        return err;
    }
    string|error message = check serviceWebSocketClient->readMessage();
    if message is string {
        json value = check value:fromJsonString(message);
        log:printInfo(value:toJsonString(value));
        string eventType = <string>check value.'type;
        json eventValue = <json>check value.'object;
        json metadata = <json>check eventValue.metadata;
        string latestResourceVersion = <string>check metadata.resourceVersion;
        setServicesResourceVersion(latestResourceVersion);
        Service|error serviceModel = createServiceModel(eventValue);
        if serviceModel is Service {
            if containsNamespace(serviceModel.namespace) {
                if eventType == "ADDED" {
                    services[serviceModel.id] = serviceModel;
                } else if (eventType == "MODIFIED") {
                    _ = services.remove(serviceModel.id);
                    services[serviceModel.id] = serviceModel;
                } else if (eventType == "DELETED") {
                    _ = services.remove(serviceModel.id);
                }
            }
        } else {
            log:printError("Unable to read service messages" + serviceModel.message());
        }
    }
}
