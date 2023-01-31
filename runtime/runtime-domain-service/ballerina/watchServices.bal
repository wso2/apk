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
import ballerina/url;
import runtime_domain_service.model;

isolated map<Service> services = {};
string servicesResourceVersion = "";
websocket:Client|error|() watchServices = ();

class ServiceTask {
    function init(string resourceVersion) {
        watchServices = getServiceClient(servicesResourceVersion);
    }
    public function startListening() returns error? {
        worker WatchServiceThread {
            while true {
                do {
                    websocket:Client|error|() serviceClientResult = watchServices;
                    if serviceClientResult is websocket:Client {
                        boolean connectionOpen = serviceClientResult.isOpen();

                        if !connectionOpen {
                            log:printDebug("ServiceWebsocket Client connection closed conectionId: " + serviceClientResult.getConnectionId());
                            watchServices = getServiceClient(servicesResourceVersion);
                            websocket:Client|error|() retryClient = watchServices;
                            if retryClient is websocket:Client {
                                log:printDebug("Reinitializing client..");
                                connectionOpen = retryClient.isOpen();
                                log:printDebug("Intializd new Client Connection conectionId: " + retryClient.getConnectionId() + " state: " + connectionOpen.toString());
                                _ = check readServiceEvents(retryClient);
                            } else if retryClient is error {
                                log:printError("error while reading message", retryClient);
                            }
                        } else {
                            log:printDebug("Intializd new Client Connection conectionId: " + serviceClientResult.getConnectionId() + " state: " + connectionOpen.toString());
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

public isolated function createServiceModel(model:Service 'service) returns Service|error {
    Service serviceData = {
        id: <string>'service.metadata.uid,
        name: <string>'service.metadata.name,
        namespace: <string>'service.metadata.namespace,
        'type: <string>'service.spec.'type,
        portmapping: check mapPortMapping('service),
        createdTime: <string>'service.metadata.creationTimestamp
    };
    return serviceData;
}

isolated function mapPortMapping(model:Service 'service) returns PortMapping[]|error {
    model:Port[]? ports = 'service.spec.ports;
    PortMapping[] portmappings = [];
    if ports is model:Port[] {
        foreach model:Port port in ports {
            PortMapping portmapping =
            {
                name: port.name ?: "",
                protocol: port.protocol,
                port: port.port,
                targetport: <int>port.targetPort
            };
            portmappings.push(portmapping);
        }
    }

    return portmappings;
}

isolated function getServicesList() returns Service[] {
    lock {
        return services.clone().toArray();
    }
}

# This retrieve specific service from name space.
#
# + name - name of service.
# + namespace - namespace of service.
# + return - service in namespace.
isolated function getService(string name, string namespace) returns Service? {
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

isolated function getServiceById(string id) returns Service|error {
    lock {
        return trap services.cloneReadOnly().get(id);
    }
}

isolated function putAllServices(map<Service> services, model:Service[] servicesEntries) {
    foreach model:Service serviceData in servicesEntries {
        lock {
            if serviceData.spec.'type != "ExternalName" {
                Service|error serviceEntry = createServiceModel(serviceData.clone());
                if serviceEntry is Service {
                    services[serviceEntry.id] = serviceEntry;
                }
            }
        }
    }
}

function setServicesResourceVersion(string resourceVersionValue) {
    servicesResourceVersion = resourceVersionValue;
}

public function getServiceClient(string resourceVersion) returns websocket:Client|error|() {
    log:printDebug("Initializing Watch Service for Services with resource Version " + resourceVersion);
    string requestURl = "wss://" + runtimeConfiguration.k8sConfiguration.host + "/api/v1/watch/services?fieldSelector=" + check getEncodedStringForNamespaces();
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

function getEncodedStringForNamespaces() returns string|error {
    string[] & readonly serviceListingNamespaces = runtimeConfiguration.serviceListingNamespaces;
    string fieldSelectorQuery = "metadata.namespace!=kube-system,metadata.namespace!=kubernetes-dashboard,metadata.namespace!=gateway-system,metadata.namespace!=ingress-nginx,metadata.namespace!=" + currentNameSpace;
    foreach string namespace in serviceListingNamespaces {
        if namespace != ALL_NAMESPACES {
            fieldSelectorQuery += ",metadata.namespace=" + namespace;
        }
    }
    return check url:encode(fieldSelectorQuery, "UTF-8");
}

function readServiceEvents(websocket:Client serviceWebSocketClient) returns error? {
    boolean connectionOpen = serviceWebSocketClient.isOpen();

    log:printDebug("Using Client Connection conectionId: " + serviceWebSocketClient.getConnectionId() + " state: " + connectionOpen.toString());
    if !connectionOpen {
        error err = error("connection closed");
        return err;
    }
    string|error message = check serviceWebSocketClient->readMessage();
    if message is string {
        log:printDebug(message);
        json value = check value:fromJsonString(message);
        string eventType = <string>check value.'type;
        json eventValue = <json>check value.'object;
        json metadata = <json>check eventValue.metadata;
        if eventType == "ERROR" {
            model:Status|error statusEvent = eventValue.cloneWithType(model:Status);
            if (statusEvent is model:Status) {
                _ = check handleWatchServicesGone(statusEvent);
            }
        } else {
            string latestResourceVersion = <string>check metadata.resourceVersion;
            setServicesResourceVersion(latestResourceVersion);
            model:Service|error mappedService = eventValue.cloneWithType(model:Service);
            if mappedService is model:Service {
                if mappedService.spec.'type != "ExternalName" {
                    Service|error serviceModel = createServiceModel(mappedService);
                    if serviceModel is Service {
                        if eventType == "ADDED" {
                            lock {
                                services[serviceModel.id] = serviceModel.clone();
                            }
                        } else if (eventType == "MODIFIED") {
                            lock {
                                _ = services.remove(serviceModel.id);
                                services[serviceModel.id] = serviceModel.clone();
                            }
                        } else if (eventType == "DELETED") {
                            lock {
                                _ = services.remove(serviceModel.id);
                            }
                        }
                    } else {
                        log:printError("Unable to read service messages" + serviceModel.message());
                    }
                }
            }
        }
    }
}

function handleWatchServicesGone(model:Status statusEvent) returns error? {
    if statusEvent.code == 410 {
        log:printDebug("Re-initializing watch service for Services due to cache clear.");
        map<Service> servicesMap = {};
        ServiceClient serviceClient = new ();
        _ = check serviceClient.retrieveAllServicesAtStartup(servicesMap, ());
        lock {
            services = servicesMap.clone();
        }
        watchAPIService = getServiceClient(resourceVersion);
    }
}
