import runtime_domain_service.model;
import ballerina/http;
import ballerina/regex;

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



    public isolated function getServices(string? query, string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError|InternalServerErrorError {
        Service[] serviceList = getServicesList().clone();
        if query is string && query.toString().trim().length() > 0 {
            return self.filterServicesBasedOnQuery(serviceList, query, sortBy, sortOrder, 'limit, offset);
        }

        return self.sortAndLimitServices(serviceList, sortBy, sortOrder, 'limit, offset);
    }

    public isolated function getServiceById(string serviceId) returns Service|BadRequestError|NotFoundError|InternalServerErrorError {
        Service|error retrievedService = getServiceById(serviceId);
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
        model:ServiceMappingList|http:ClientError retrieveAllServiceMappingResult;
        if resultValue is string {
            retrieveAllServiceMappingResult = retrieveAllServiceMappings(resultValue);
        } else {
            retrieveAllServiceMappingResult = retrieveAllServiceMappings(());
        }

        if retrieveAllServiceMappingResult is model:ServiceMappingList {
            model:ListMeta metadata = retrieveAllServiceMappingResult.metadata;
            model:K8sServiceMapping[] items =  retrieveAllServiceMappingResult.items;
            _ = check putAllServiceMappings(items);

            json|error continueElement = metadata.'continue;
            if continueElement is json {
                if (<string>continueElement).length() > 0 {
                    _ = check self.retrieveAllServicesAtStartup(<string?>continueElement);
                }
            }
            string? resourceVersion =  metadata.'resourceVersion;
            if resourceVersion is string{
            setServiceMappingResourceVersion(resourceVersion);                
            }
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
                    'version: k8sAPI.spec.apiVersion
                });
            }
            APIList apiList = {list: apiInfos, count: apiInfos.length(), pagination: {total: apiInfos.length()}};
            return apiList;
        } else {
            return serviceEntry;
        }
    }

    public isolated function filterServicesBasedOnQuery(Service[] servicesList, string query, string sortBy, string sortOrder, int 'limit, int offset) returns ServiceList|BadRequestError|InternalServerErrorError {
        Service[] filteredList = [];
        if query.length() > 0 {
            int? semiCollonIndex = string:indexOf(query, ":", 0);
            if semiCollonIndex is int {
                if semiCollonIndex > 0 {
                    string keyWord = query.substring(0, semiCollonIndex);
                    string keyWordValue = query.substring(keyWord.length() + 1, query.length());
                    if keyWord.trim() == SEARCH_CRITERIA_NAME {
                        keyWordValue = keyWordValue + "|\\w+" + keyWordValue + "\\w+" + "|" + keyWordValue + "\\w+" + "|\\w+" + keyWordValue;
                        foreach Service 'service in servicesList {
                            if (regex:matches('service.name, keyWordValue)) {
                                filteredList.push('service);
                            }
                        }
                    } else if keyWord.trim() == SEARCH_CRITERIA_NAMESPACE {
                        foreach Service 'service in servicesList {
                            if (regex:matches('service.namespace, keyWordValue)) {
                                filteredList.push('service);
                            }
                        }
                    } else {
                        BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord " + keyWord}};
                        return badRequest;
                    }
                }
            } else {
                string keyWordValue = query + "|\\w+" + query + "\\w+" + "|" + query + "\\w+" + "|\\w+" + query;
                foreach Service 'service in servicesList {
                    if (regex:matches('service.name, keyWordValue)) {
                        filteredList.push('service);
                    }
                }
            }
        } else {
            filteredList = servicesList;
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
