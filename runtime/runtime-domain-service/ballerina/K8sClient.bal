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

import runtime_domain_service.model;
import ballerina/io;
import ballerina/http;

const string K8S_API_ENDPOINT = "/api/v1";
final http:Client k8sApiServerEp = check initializeK8sClient();
configurable string k8sHost = "kubernetes.default";
configurable string saTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token";
string token = check io:fileReadString(saTokenPath);
configurable string caCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt";

# This initialize the k8s Client.
# + return - k8s http client
function initializeK8sClient() returns http:Client|error {
    http:Client k8sApiClient = check new ("https://" + k8sHost,
    auth = {
        token: token
    },
        secureSocket = {
            cert: caCertPath

        }
    );
    return k8sApiClient;
}

# This returns services in a namsepace.
#
# + namespace - namespace value
# + return - list of services in namespace.
function getServicesListInNamespace(string namespace) returns ServiceList|error {
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
function getServicesListFromK8s() returns ServiceList|error {
    return {list: getServicesList(), pagination: {total: getServicesList().length()}};
}

# This retrieve specific service from name space.
#
# + name - name of service.
# + namespace - namespace of service.
# + return - service in namespace.
function getServiceFromK8s(string name, string namespace) returns ServiceList|error {
    Service? serviceResult = getService(name, namespace);
    if serviceResult is null {
        return {list: []};
    } else {
        return {list: [serviceResult]};
    }
}

# This returns list of APIS.
#
# + return - Return list of APIS in namsepace.
function getAPIList() returns APIList|error {
    API[] apilist = [];
    foreach model:K8sAPI api in getAPIs() {
        API convertedModel = convertK8sAPItoAPI(api);
        apilist.push(convertedModel);
    }
    APIList APIList = {
        list: apilist
    };
    return APIList;
}

function convertK8sAPItoAPI(model:K8sAPI api) returns API {
    API convetedModel = {
        id: api.uuid,
        name: api.apiDisplayName,
        context: api.context,
        'version: api.apiVersion,
        'type: api.apiType,
        createdTime: api.creationTimestamp
    };
    return convetedModel;
}

//Get APIs deployed in default namespace by APIId.
function getAPIById(string id) returns API|InternalServerErrorError|BadRequestError|NotFoundError|error {
    boolean APIIDAvailable = id.length() > 0 ? true : false;
    if (APIIDAvailable && string:length(id.toString()) > 0)
    {
        model:K8sAPI? api = apilist[id];
        if api != null {
            API detailedAPI = convertK8sAPItoAPI(api);
            return detailedAPI;
        } else {
            NotFoundError notfound = {body: {code: 909100, message: id + "not found."}};
            return notfound;
        }
    }
    BadRequestError badRequestError = {body: {code: 900910, message: "missing required attributes"}};
    return badRequestError;
}

//Delete APIs deployed in a namespace by APIId.
function deleteAPIById(string id) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
    boolean APIIDAvailable = id.length() > 0 ? true : false;
    if (APIIDAvailable && string:length(id.toString()) > 0)
    {
        model:K8sAPI|error api = getAPI(id);
        if api is model:K8sAPI {
            string endpoint = "/apis/dp.wso2.com/v1alpha1/namespaces/" + api.namespace + "/apis/" + api.k8sName;
            error|json APIResp = k8sApiServerEp->delete(endpoint, targetType = json);
            if APIResp is error {
                NotFoundError internalError = {body: {code: 900910, message: "APIResp.message()"}};
                return internalError;
            } else {
                return http:OK;
            }
        } else {
            NotFoundError apiNotfound = {body: {code: 900910, description: "API with " + id + " not found", message: "API not found"}};
            return apiNotfound;
        }
    }
    PreconditionFailedError badRequestError = {body: {code: 900910, message: "missing required attributes"}};
    return badRequestError;
}

//Get all deployed APIs in namespace with specific search query
function getAPIListInNamespaceWithQuery(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc") returns APIList|InternalServerErrorError|BadRequestError|error {
    APIInfo[] apiNames = map:toArray(apilist);
    return {list: apiNames, count: apiNames.length(), pagination: {total: apilist.length()}};
}
