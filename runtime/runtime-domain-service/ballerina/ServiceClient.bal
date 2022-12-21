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
    public isolated function getServicesListInNamespace(string namespace) returns ServiceList|error {
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
    #
    # + sortBy - sort by to sort services (name,createdTime)  
    # + sortOrder - Order to sort (asc,desc) 
    # + 'limit - no of services to return  
    # + offset - offset value
    # + return - list of services in namespaces.
    public isolated function getServicesListFromK8s(string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError {
        return self.sortAndLimitServices(getServicesList(), sortBy, sortOrder, 'limit, offset);
    }

    # This retrieve specific service from name space.
    #
    # + name - name of service.
    # + namespace - namespace of service.
    # + return - service in namespace.
    public isolated function getServiceFromK8s(string name, string namespace) returns ServiceList|error {
        Service? serviceResult = getService(name, namespace);
        if serviceResult is null {
            return {list: []};
        } else {
            return {list: [serviceResult]};
        }
    }

    public isolated function getServices(string? name, string? namespace, string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError|InternalServerErrorError {
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
        } else {
            if (serviceNameAvailable && string:length(name.toString()) > 0) {
                return self.getServicesListFromK8sSearchByName(name.toString(), sortBy, sortOrder, 'limit, offset);
            }
        }
        return self.getServicesListFromK8s(sortBy, sortOrder, 'limit, offset);
    }

    public isolated function getServiceById(string serviceId) returns Service|BadRequestError|NotFoundError|InternalServerErrorError {
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

    public isolated function retrieveK8sServiceMapping(string name, string namespace) returns Service|error {
        json serviceByNameAndNamespace = check getServiceByNameAndNamespace(name, namespace);
        return createServiceModel(serviceByNameAndNamespace);
    }

    public isolated function getServiceUsageByServiceId(string serviceId) returns APIList|BadRequestError|NotFoundError|InternalServerErrorError {
        APIInfo[] apiInfos = [];
        Service|BadRequestError|NotFoundError|InternalServerErrorError serviceEntry = self.getServiceById(serviceId);
        if serviceEntry is Service {
            model:API[] k8sAPIS = retrieveAPIMappingsForService(serviceEntry);
            foreach model:API k8sAPI in k8sAPIS {
                apiInfos.push({
                    context: k8sAPI.spec.context,
                    createdTime: k8sAPI.metadata.creationTimestamp,
                    name: k8sAPI.spec.apiDisplayName,
                    id: k8sAPI.metadata.uid,
                    'type: k8sAPI.spec.apiType,
                    'version: k8sAPI.apiVersion
                });
            }
            APIList apiList = {list: apiInfos, count: apiInfos.length(), pagination: {total: apiInfos.length()}};
            return apiList;
        } else {
            return serviceEntry;
        }
    }

    private isolated function getServicesListFromK8sSearchByName(string name, string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError {
        Service[] servicesList = getServicesList().cloneReadOnly();
        Service[] filteredList = [];
        foreach Service 'service in servicesList {
            if 'service.name == name {
                filteredList.push('service);
            }
        }
        return self.sortAndLimitServices(filteredList, sortBy, sortOrder, 'limit, offset);
    }

    private isolated function sortAndLimitServices(Service[] servicesList, string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError {
        Service[] clonedServiceList = servicesList.clone();
        Service[] sortedServices = [];
        if sortBy == SORT_BY_SERVICE_NAME && sortOrder == SORT_ORDER_ASC {
            sortedServices = from var 'service in clonedServiceList
                order by 'service.name ascending
                select 'service;
        } else if sortBy == SORT_BY_SERVICE_NAME && sortOrder == SORT_ORDER_DESC {
            sortedServices = from var 'service in clonedServiceList
                order by 'service.name descending
                select 'service;
        } else if sortBy == SORT_BY_CREATED_TIME && sortOrder == SORT_ORDER_ASC {
            sortedServices = from var 'service in clonedServiceList
                order by 'service.createdTime ascending
                select 'service;
        } else if sortBy == SORT_BY_CREATED_TIME && sortOrder == SORT_ORDER_DESC {
            sortedServices = from var 'service in clonedServiceList
                order by 'service.createdTime descending
                select 'service;
        } else {
            BadRequestError badRequest = {body: {code: 90912, message: "Invalid Sort By/Sort Order Value "}};
            return badRequest;
        }
        Service[] limitedServices = [];
        if sortedServices.length() >= offset {
            foreach int i in offset ... (sortedServices.length() - 1) {
                if limitedServices.length() < 'limit {
                    limitedServices.push(sortedServices[i]);
                }
            }
        }
        ServiceList serviceList = {list: limitedServices, pagination: {offset: offset, 'limit: 'limit, total: sortedServices.length()}};
        return serviceList;
    }
}
