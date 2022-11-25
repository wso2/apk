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

import ballerina/http;
import ballerina/log;
import runtime_domain_service.model;

function getAPIDefinitionByID(string id) returns string|NotFoundError|NotAcceptableError {
    model:K8sAPI|error api = getAPI(id);
    if api is model:K8sAPI {
        if api.definitionFileRef.length() > 0 {
            string|error definition = getDefinition(api);
            if definition is string {
                return definition;
            } else {
                log:printError("Error while reading definition:", definition);
            }
        }
    }
    NotFoundError notfound = {body: {code: 909100, message: id + "not found."}};
    return notfound;
}

function getDefinition(model:K8sAPI api) returns string|error {
    json|error configMapValue = getConfigMapValueFromNameAndNamespace(api.definitionFileRef, api.namespace);
    if configMapValue is json {
        json|error data = configMapValue.data;
        json|error binaryData = configMapValue.binaryData;
        if data is json {
            map<json> dataMap = <map<json>>data;
            string[] keys = dataMap.keys();
            if keys.length() == 1 {
                return dataMap.get(keys[0]).toJsonString();
            }
        } else if binaryData is json {
            map<json> dataMap = <map<json>>binaryData;
            string[] keys = dataMap.keys();
            if keys.length() == 1 {
                return dataMap.get(keys[0]).toJsonString();
            }
        }
        return "";
    } else {
        return configMapValue;
    }
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

function createAPI(API api) returns string|Error {
    if (validateName(api.name)) {
        return {code: 90911, message: "API Name `${api.name}` already exist.", description: "API Name `${api.name}` already exist."};
    }
    if validateContextAndVersion(api.context, api.'version) {
        return {code: 90912, message: "API Context `${api.context}` already exist.", description: "API Context `${api.context}` already exist."};
    }
    return "created";
}

function validateContextAndVersion(string context, string 'version) returns boolean {

    foreach model:K8sAPI k8sAPI in getAPIs() {
        if k8sAPI.context == returnFullContext(context, 'version) {
            return true;
        }
    }
    return false;
}

function returnFullContext(string context, string 'version) returns string {
    string fullContext = context;
    if (!string:endsWith(context, 'version)) {
        fullContext = string:'join("/", context, 'version);
    }
    return fullContext;
}

function validateName(string name) returns boolean {
    foreach model:K8sAPI k8sAPI in getAPIs() {
        if k8sAPI.apiDisplayName == name {
            return true;
        }
    }
    return false;
}

function createAndDeployAPI(API api) {
    model:API k8sAPI = convertK8sCrAPI(api);
    log:printInfo(<string>k8sAPI.toJson());
}

function convertK8sCrAPI(API api) returns model:API {
    model:API apispec = {
        metadata: {name: api.name.concat(api.'version), namespace: runtimeConfiguration.apiCreationNamespace},
        spec: {
            apiDisplayName: api.name,
            apiType: api.'type,
            apiVersion: api.'version,
            context: returnFullContext(api.context, api.'version),
            definitionFileRef: "",
            prodHTTPRouteRef: "",
            sandHTTPRouteRef: ""
        }
    };
    return apispec;
}
