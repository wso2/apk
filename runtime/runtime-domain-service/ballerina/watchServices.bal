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
import ballerina/log;

final websocket:Client servicesClient = check new ("wss://" + k8sHost + "/api/v1/watch/services",
auth = {
    token: token
},
secureSocket = {
    cert: caCertPath
});

map<Service> services = {};

class ServiceTask {
    *task:Job;

    public function execute() {
        do {
            string|error message = check servicesClient->readMessage();
            if message is string {
                json value = check value:fromJsonString(message);
                // log:printInfo(value:toJsonString(value));
                string eventType = <string>check value.'type;
                json eventValue = <json>check value.'object;
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
        on fail var e {
            log:printError("Unable to read service messages", e);
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
        portmapping: check mapPortMapping(event)
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

    return;
}

function grtServiceById(string id) returns Service|error {
    return trap services.get(id);
}
