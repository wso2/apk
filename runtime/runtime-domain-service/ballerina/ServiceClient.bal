import runtime_domain_service.model;
import ballerina/http;

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

public class ServiceClient {

    # This returns services in a namsepace.
    #
    # + namespace - namespace value
    # + return - list of services in namespace.
    public function getServicesListInNamespace(string namespace) returns ServiceList|error {
        Service[] servicesList = getServicesList();
        Service[] filteredList = [];
        foreach Service item in servicesList {
            if item.namespace == namespace {
                filteredList.push(item);
            }
        }
        return {list: filteredList, pagination: {total: filteredList.length()}};
    }

    # This returns list of services in all namespaces.
    # + return - list of services in namespaces.
    public function getServicesListFromK8s() returns ServiceList|error {
        return {list: getServicesList(), pagination: {total: getServicesList().length()}};
    }

    # This retrieve specific service from name space.
    #
    # + name - name of service.
    # + namespace - namespace of service.
    # + return - service in namespace.
    public function getServiceFromK8s(string name, string namespace) returns ServiceList|error {
        Service? serviceResult = getService(name, namespace);
        if serviceResult is null {
            return {list: []};
        } else {
            return {list: [serviceResult]};
        }
    }

    public function getServices(string? name, string? namespace, string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError|UnauthorizedError|InternalServerErrorError {
        boolean serviceNameAvailable = name == () ? false : true;
        boolean nameSpaceAvailable = namespace == () ? false : true;
        if (nameSpaceAvailable && string:length(namespace.toString()) > 0) {
            if (serviceNameAvailable && string:length(name.toString()) > 0) {
                ServiceList|error serviceList = self.getServiceFromK8s(name.toString(), namespace.toString());
                if serviceList is error {
                    InternalServerErrorError internalError = {body: {code: 900910, message: serviceList.message()}};
                    return internalError;
                } else {
                    return serviceList;
                }
            } else {
                ServiceList|error serviceList = self.getServicesListInNamespace(namespace.toString());
                if serviceList is error {
                    InternalServerErrorError internalError = {body: {code: 900910, message: serviceList.message()}};
                    return internalError;
                } else {
                    return serviceList;
                }
            }
        }
        ServiceList|error serviceList = self.getServicesListFromK8s();
        if serviceList is error {
            InternalServerErrorError internalError = {body: {code: 900910, message: serviceList.message()}};
            return internalError;
        } else {
            return serviceList;
        }
    }

    public function getServiceById(string serviceId) returns Service|BadRequestError|NotFoundError|InternalServerErrorError {
        Service|error retrievedService = grtServiceById(serviceId);
        if retrievedService is Service {
            return retrievedService;
        } else {
            NotFoundError notfound = {body: {code: 90914, message: "Service " + serviceId + " not found"}};
            return notfound;
        }
    }
    public function retrieveAllServicesAtStartup(string? continueValue) returns error? {
        string? resultValue = continueValue;
        json|http:ClientError retrieveAllServicesResult;
        if resultValue is string {
            retrieveAllServicesResult = retrieveAllServices(resultValue);
        } else {
            retrieveAllServicesResult = retrieveAllServices(());
        }

        if retrieveAllServicesResult is json {
            json metadata = check retrieveAllServicesResult.metadata;
            json[] items = <json[]>check retrieveAllServicesResult.items;
            putAllServices(items);

            json|error continueElement = metadata.'continue;
            if continueElement is json {
                if (<string>continueElement).length() > 0 {
                    _ = check self.retrieveAllServicesAtStartup(<string?>continueElement);
                }
            }
            string resourceVersion = <string>check metadata.'resourceVersion;
            setServicesResourceVersion(resourceVersion);
        }
    }
    public function retrieveAllServiceMappingsAtStartup(string? continueValue) returns error? {
        string? resultValue = continueValue;
        json|http:ClientError retrieveAllServiceMappingResult;
        if resultValue is string {
            retrieveAllServiceMappingResult = retrieveAllServiceMappings(resultValue);
        } else {
            retrieveAllServiceMappingResult = retrieveAllServiceMappings(());
        }

        if retrieveAllServiceMappingResult is json {
            json metadata = check retrieveAllServiceMappingResult.metadata;
            json[] items = <json[]>check retrieveAllServiceMappingResult.items;
            _ = check putAllServiceMappings(items);

            json|error continueElement = metadata.'continue;
            if continueElement is json {
                if (<string>continueElement).length() > 0 {
                    _ = check self.retrieveAllServicesAtStartup(<string?>continueElement);
                }
            }
            string resourceVersion = <string>check metadata.'resourceVersion;
            setServiceMappingResourceVersion(resourceVersion);
        }
    }
    public function retrieveK8sServiceMapping(string name, string namespace) returns Service|error {
        json serviceByNameAndNamespace = check getServiceByNameAndNamespace(name, namespace);
        return createServiceModel(serviceByNameAndNamespace);
    }
    public function getServiceUsageByServiceId(string serviceId) returns APIList|BadRequestError|NotFoundError|InternalServerErrorError {
        APIInfo[] apiInfos = [];
        map<model:K8sAPI>|error retrievedUsage = trap serviceMappings.get(serviceId);
        if retrievedUsage is map<model:K8sAPI> {
            string[] keys = retrievedUsage.keys();
            foreach string key in keys {
                model:K8sAPI k8sAPI = retrievedUsage.get(key);
                apiInfos.push({
                    context: k8sAPI.context,
                    createdTime: k8sAPI.creationTimestamp,
                    name: k8sAPI.apiDisplayName,
                    id: k8sAPI.uuid,
                    'type: k8sAPI.apiType,
                    'version: k8sAPI.apiVersion
                });
            }
        }
        APIList apiList = {list: apiInfos, count: apiInfos.length(), pagination: {total: apiInfos.length()}};
        return apiList;
    }
}
