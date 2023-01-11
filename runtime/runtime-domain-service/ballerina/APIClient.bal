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
import ballerina/uuid;
import runtime_domain_service.model;
import runtime_domain_service.org.wso2.apk.runtime.model as runtimeModels;
import runtime_domain_service.java.util as utilapis;
import ballerina/jwt;
import ballerina/regex;
import runtime_domain_service.org.wso2.apk.runtime as runtimeUtil;
import ballerina/mime;
import ballerina/jballerina.java;
import ballerina/lang.value;
import runtime_domain_service.org.wso2.apk.runtime.api as runtimeapi;

public class APIClient {

    public isolated function getAPIDefinitionByID(string id) returns json|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        model:API|error api = getAPI(id);
        if api is model:API {
            json|error definition = self.getDefinition(api);
            if definition is json {
                return definition;
            } else {
                log:printError("Error while reading definition:", definition);
                InternalServerErrorError internalError = {body: {code: 909000, message: "Internal Error Occured while retrieving definition"}};
                return internalError;
            }
        }
        NotFoundError notfound = {body: {code: 909100, message: id + " not found."}};
        return notfound;
    }

    private isolated function getDefinition(model:API api) returns json|APKError {
        do {
            string? definitionFileRef = api.spec.definitionFileRef;
            if definitionFileRef is string && definitionFileRef.length() > 0 {
                http:Response response = check getConfigMapValueFromNameAndNamespace(definitionFileRef, api.metadata.namespace);
                if response.statusCode == 200 {
                    json configMapValue = check response.getJsonPayload();
                    json|error data = configMapValue.data;
                    if data is json {
                        map<json> dataMap = <map<json>>data;
                        string[] keys = dataMap.keys();
                        if keys.length() == 1 {
                            return check value:fromJsonString(dataMap.get(keys[0]).toString());
                        }
                    }
                }
            }
            // definitionFileRef not specified or empty. definitionfile
            return self.retrieveDefaultDefinition(api);
        } on fail var e {
            string message = "Internal Error occured while retrieving api Definition";
            return error(message, e, message = message, description = message, code = 909000, statusCode = "500");
        }

    }

    //Get APIs deployed in default namespace by APIId.
    public isolated function getAPIById(string id) returns API|NotFoundError {
        boolean APIIDAvailable = id.length() > 0 ? true : false;
        if (APIIDAvailable && string:length(id.toString()) > 0)
        {
            lock {
                model:API? api = apilist[id];
                if api != null {
                    API detailedAPI = convertK8sAPItoAPI(api);
                    return detailedAPI.cloneReadOnly();
                }
            }
        }
        NotFoundError notfound = {body: {code: 909100, message: id + " not found."}};
        return notfound;
    }

    //Delete APIs deployed in a namespace by APIId.
    public isolated function deleteAPIById(string id) returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|APKError {
        boolean APIIDAvailable = id.length() > 0 ? true : false;
        if (APIIDAvailable && string:length(id.toString()) > 0)
        {
            model:API|error api = getAPI(id);
            if api is model:API {
                http:Response|http:ClientError apiCRDeletionResponse = deleteAPICR(api.metadata.name, api.metadata.namespace);
                if apiCRDeletionResponse is http:ClientError {
                    log:printError("Error while undeploying API CR ", apiCRDeletionResponse);
                }
                string? definitionFileRef = api.spec.definitionFileRef;
                if definitionFileRef is string {
                    http:Response|http:ClientError apiDefinitionDeletionResponse = deleteConfigMap(definitionFileRef, api.metadata.namespace);
                    if apiDefinitionDeletionResponse is http:ClientError {
                        log:printError("Error while undeploying API definition ", apiDefinitionDeletionResponse);
                    }
                }
                string? prodHTTPRouteRef = api.spec.prodHTTPRouteRef;
                if prodHTTPRouteRef is string && prodHTTPRouteRef.toString().length() > 0 {
                    http:Response|http:ClientError prodHttpRouteDeletionResponse = deleteHttpRoute(prodHTTPRouteRef, api.metadata.namespace);
                    if prodHttpRouteDeletionResponse is http:ClientError {
                        log:printError("Error while undeploying prod http route ", prodHttpRouteDeletionResponse);
                    }
                }
                string? sandBoxHttpRouteRef = api.spec.sandHTTPRouteRef;
                if sandBoxHttpRouteRef is string && sandBoxHttpRouteRef.toString().length() > 0 {
                    http:Response|http:ClientError sandHttpRouteDeletionResponse = deleteHttpRoute(sandBoxHttpRouteRef, api.metadata.namespace);
                    if sandHttpRouteDeletionResponse is http:ClientError {
                        log:printError("Error while undeploying prod http route ", sandHttpRouteDeletionResponse);
                    }
                }
                APKError? response = check self.deleteServiceMappings(api);
                if response is APKError {
                    return response;
                }
            } else {
                NotFoundError apiNotfound = {body: {code: 900910, description: "API with " + id + " not found", message: "API not found"}};
                return apiNotfound;
            }
        }
        return http:OK;
    }

    # This returns list of APIS.
    #
    # + query - Parameter Description  
    # + 'limit - Parameter Description  
    # + offset - Parameter Description  
    # + sortBy - Parameter Description  
    # + sortOrder - Parameter Description
    # + return - Return list of APIS in namsepace.
    public isolated function getAPIList(string? query, int 'limit, int offset, string sortBy, string sortOrder) returns APIList|BadRequestError {
        API[] apilist = [];
        foreach model:API api in getAPIs() {
            API convertedModel = convertK8sAPItoAPI(api);
            apilist.push(convertedModel);
        }
        if query is string && query.toString().trim().length() > 0 {
            return self.filterAPISBasedOnQuery(apilist, query, 'limit, offset, sortBy, sortOrder);
        } else {
            return self.filterAPIS(apilist, 'limit, offset, sortBy, sortOrder);
        }
    }
    private isolated function filterAPISBasedOnQuery(API[] apilist, string query, int 'limit, int offset, string sortBy, string sortOrder) returns APIList|BadRequestError {
        API[] filteredList = [];
        if query.length() > 0 {
            int? semiCollonIndex = string:indexOf(query, ":", 0);
            if semiCollonIndex is int && semiCollonIndex > 0 {
                string keyWord = query.substring(0, semiCollonIndex);
                string keyWordValue = query.substring(keyWord.length() + 1, query.length());
                keyWordValue = keyWordValue + "|\\w+" + keyWordValue + "\\w+" + "|" + keyWordValue + "\\w+" + "|\\w+" + keyWordValue;
                if keyWord.trim() == SEARCH_CRITERIA_NAME {
                    foreach API api in apilist {
                        if (regex:matches(api.name, keyWordValue)) {
                            filteredList.push(api);
                        }
                    }
                } else if keyWord.trim() == SEARCH_CRITERIA_TYPE {
                    foreach API api in apilist {
                        if (regex:matches(api.'type, keyWordValue)) {
                            filteredList.push(api);
                        }
                    }
                } else {
                    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord " + keyWord}};
                    return badRequest;
                }
            } else {
                string keyWordValue = query + "|\\w+" + query + "\\w+" + "|" + query + "\\w+" + "|\\w+" + query;

                foreach API api in apilist {

                    if (regex:matches(api.name, keyWordValue)) {
                        filteredList.push(api);
                    }
                }
            }
        } else {
            filteredList = apilist;
        }
        return self.filterAPIS(filteredList, 'limit, offset, sortBy, sortOrder);
    }
    private isolated function filterAPIS(API[] apiList, int 'limit, int offset, string sortBy, string sortOrder) returns APIList|BadRequestError {
        API[] clonedAPIList = apiList.clone();
        API[] sortedAPIS = [];
        if sortBy == SORT_BY_API_NAME && sortOrder == SORT_ORDER_ASC {
            sortedAPIS = from var api in clonedAPIList
                order by api.name ascending
                select api;
        } else if sortBy == SORT_BY_API_NAME && sortOrder == SORT_ORDER_DESC {
            sortedAPIS = from var api in clonedAPIList
                order by api.name descending
                select api;
        } else if sortBy == SORT_BY_CREATED_TIME && sortOrder == SORT_ORDER_ASC {
            sortedAPIS = from var api in clonedAPIList
                order by api.createdTime ascending
                select api;
        } else if sortBy == SORT_BY_CREATED_TIME && sortOrder == SORT_ORDER_DESC {
            sortedAPIS = from var api in clonedAPIList
                order by api.createdTime descending
                select api;
        } else {
            BadRequestError badRequest = {body: {code: 90912, message: "Invalid Sort By/Sort Order Value "}};
            return badRequest;
        }
        API[] limitSet = [];
        if sortedAPIS.length() > offset {
            foreach int i in offset ... (sortedAPIS.length() - 1) {
                if limitSet.length() < 'limit {
                    limitSet.push(sortedAPIS[i]);
                }
            }
        }
        return {list: limitSet, count: limitSet.length(), pagination: {total: apiList.length(), 'limit: 'limit, offset: offset}};

    }
    public isolated function createAPI(API api) returns APKError|CreatedAPI|BadRequestError {
        do {

            if (self.validateName(api.name)) {
                BadRequestError badRequest = {body: {code: 90911, message: "API Name - " + api.name + " already exist.", description: "API Name - " + api.name + " already exist."}};
                return badRequest;
            }
            if self.validateContextAndVersion(api.context, api.'version) {
                BadRequestError badRequest = {body: {code: 90911, message: "API Context - " + api.context + " already exist.", description: "API Context " + api.context + " already exist."}};
                return badRequest;
            }

            self.setDefaultOperationsIfNotExist(api);
            string uniqueId = uuid:createType1AsString();
            model:APIArtifact apiArtifact = {uniqueId: uniqueId};
            APIOperations[]? operations = api.operations;
            if operations is APIOperations[] {
                if operations.length() == 0 {
                    BadRequestError badRequestError = {body: {code: 90912, message: "Atleast one operation need to specified"}};
                    return badRequestError;
                }
            } else {
                BadRequestError badRequestError = {body: {code: 90912, message: "Atleast one operation need to specified"}};
                return badRequestError;
            }
            record {}? endpointConfig = api.endpointConfig;
            if endpointConfig is record {} {
                map<model:Endpoint> createdEndpoints = check self.createAndAddBackendServics(apiArtifact, api, endpointConfig, (), ());
                if createdEndpoints.hasKey(PRODUCTION_TYPE) {
                    _ = check self.setHttpRoute(apiArtifact, api, createdEndpoints.get(PRODUCTION_TYPE), uniqueId, PRODUCTION_TYPE);
                } else {
                    _ = check self.setHttpRoute(apiArtifact, api, (), uniqueId, PRODUCTION_TYPE);
                }
                if createdEndpoints.hasKey(SANDBOX_TYPE) {
                    _ = check self.setHttpRoute(apiArtifact, api, createdEndpoints.get(SANDBOX_TYPE), uniqueId, SANDBOX_TYPE);
                } else {
                    _ = check self.setHttpRoute(apiArtifact, api, (), uniqueId, SANDBOX_TYPE);
                }
            } else {
                _ = check self.setHttpRoute(apiArtifact, api, (), uniqueId, PRODUCTION_TYPE);
                _ = check self.setHttpRoute(apiArtifact, api, (), uniqueId, SANDBOX_TYPE);
            }
            json generatedSwagger = check self.retrieveGeneratedSwaggerDefinition(api);
            self.retrieveGeneratedConfigmapForDefinition(apiArtifact, api, generatedSwagger, uniqueId);
            self.generateAndSetAPICRArtifact(apiArtifact, api);
            model:API? deployAPIToK8sResult = check self.deployAPIToK8s(apiArtifact);
            if deployAPIToK8sResult is model:API {
                CreatedAPI createdAPI = {body: {name: api.name, context: self.returnFullContext(api.context, api.'version), 'version: api.'version, id: deployAPIToK8sResult.metadata.uid}};
                return createdAPI;
            } else {
                return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = "500");
            }
        } on fail var e {
            log:printError("Internal Error occured", e);
            return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = "500");
        }
    }

    private isolated function createAndAddBackendServics(model:APIArtifact apiArtifact, API api, record {} endpointConfig, APIOperations? apiOperation, string? endpointType) returns map<model:Endpoint>|APKError|error {
        map<model:Endpoint> endpointIdMap = {};
        anydata|error sandboxEndpointConfig = trap endpointConfig.get("sandbox_endpoints");
        anydata|error productionEndpointConfig = trap endpointConfig.get("production_endpoints");
        if endpointType == () || (endpointType == SANDBOX_TYPE) {
            if sandboxEndpointConfig is map<anydata> {
                if sandboxEndpointConfig.hasKey("url") {
                    anydata url = sandboxEndpointConfig.get("url");
                    model:Service backendService = self.createBackendService(<string>url, self.getLabels(apiArtifact.uniqueId, api, apiOperation));
                    if apiOperation == () {
                        apiArtifact.sandboxEndpointAvailable = true;
                        apiArtifact.sandboxUrl = <string?>url;
                    }
                    apiArtifact.backendServices.push(backendService);
                    endpointIdMap[SANDBOX_TYPE] = {
                        port: check self.getPort(<string>url),
                        namespace: backendService.metadata.namespace,
                        name: backendService.metadata.name,
                        serviceEntry: false,
                        url: <string?>url
                    };
                } else {
                    APKError e = error("Sandbox Endpoint Not specified", message = "Endpoint Not specified", description = "Sandbox Endpoint Not specified", code = 90911, statusCode = "400");
                    return e;
                }
            }
        }
        if endpointType == () || (endpointType == PRODUCTION_TYPE) {
            if productionEndpointConfig is map<anydata> {
                if productionEndpointConfig.hasKey("url") {
                    anydata url = productionEndpointConfig.get("url");
                    model:Service backendService = self.createBackendService(<string>url, self.getLabels(apiArtifact.uniqueId, api, apiOperation));
                    if apiOperation == () {
                        apiArtifact.productionEndpointAvailable = true;
                        apiArtifact.productionUrl = <string?>url;
                    }
                    apiArtifact.backendServices.push(backendService);
                    endpointIdMap[PRODUCTION_TYPE] = {
                        port: check self.getPort(<string>url),
                        namespace: backendService.metadata.namespace,
                        name: backendService.metadata.name,
                        serviceEntry: false,
                        url: <string?>url
                    };
                } else {
                    APKError e = error("Production Endpoint Not specified", message = "Endpoint Not specified", description = "Production Endpoint Not specified", code = 90911, statusCode = "400");
                    return e;
                }
            }
        }
        return endpointIdMap;
    }
    private isolated function getLabels(string uniqueId, API api, APIOperations? apiOperation) returns map<string> {
        map<string> labels = {"api-name": api.name, "api-version": api.'version, "k8sapi-name": uniqueId};
        return labels;
    }
    isolated function validateContextAndVersion(string context, string 'version) returns boolean {

        foreach model:API k8sAPI in getAPIs() {
            if k8sAPI.spec.context == self.returnFullContext(context, 'version) {
                return true;
            }
        }
        return false;
    }
    isolated function validateContext(string context) returns boolean {

        foreach model:API k8sAPI in getAPIs() {
            if k8sAPI.spec.context == context {
                return true;
            }
        }
        return false;
    }
    private isolated function returnFullContext(string context, string 'version) returns string {
        string fullContext = context;
        if (!string:endsWith(context, 'version)) {
            fullContext = string:'join("/", context, 'version);
        }
        return fullContext;
    }

    isolated function validateName(string name) returns boolean {
        foreach model:API k8sAPI in getAPIs() {
            if k8sAPI.spec.apiDisplayName == name {
                return true;
            }
        }
        return false;
    }

    function convertK8sCrAPI(API api) returns model:API {
        model:API apispec = {
            metadata: {
                name: api.name.concat(api.'version),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: ()
            },
            spec: {
                apiDisplayName: api.name,
                apiType: api.'type,
                apiVersion: api.'version,
                context: self.returnFullContext(api.context, api.'version),
                definitionFileRef: "",
                prodHTTPRouteRef: "",
                sandHTTPRouteRef: "",
                organization: ""
            }
        };
        return apispec;
    }

    isolated function createAPIFromService(string serviceKey, API api) returns CreatedAPI|BadRequestError|InternalServerErrorError|APKError {
        if (self.validateName(api.name)) {
            BadRequestError badRequest = {body: {code: 90911, message: "API Name - " + api.name + " already exist.", description: "API Name - " + api.name + " already exist."}};
            return badRequest;
        }
        if self.validateContextAndVersion(api.context, api.'version) {
            BadRequestError badRequest = {body: {code: 90911, message: "API Context - " + api.context + " already exist.", description: "API Context " + api.context + " already exist."}};
            return badRequest;
        }
        self.setDefaultOperationsIfNotExist(api);
        Service|error serviceRetrieved = getServiceById(serviceKey);
        string uniqueId = uuid:createType1AsString();
        if serviceRetrieved is Service {
            model:APIArtifact apiArtifact = {uniqueId: uniqueId};
            model:Endpoint endpoint = {
                port: self.retrievePort(serviceRetrieved),
                namespace: serviceRetrieved.namespace,
                name: serviceRetrieved.name,
                serviceEntry: true
            };
            check self.setHttpRoute(apiArtifact, api, endpoint, uniqueId, PRODUCTION_TYPE);
            json generatedSwaggerDefinition = check self.retrieveGeneratedSwaggerDefinition(api);
            self.retrieveGeneratedConfigmapForDefinition(apiArtifact, api, generatedSwaggerDefinition, uniqueId);
            self.generateAndSetAPICRArtifact(apiArtifact, api);
            self.generateAndSetK8sServiceMapping(apiArtifact, api, serviceRetrieved, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            APKError|model:API deployAPIToK8sResult = self.deployAPIToK8s(apiArtifact);
            if deployAPIToK8sResult is model:API {
                CreatedAPI createdAPI = {body: {name: api.name, context: self.returnFullContext(api.context, api.'version), 'version: api.'version, id: deployAPIToK8sResult.metadata.uid}};
                return createdAPI;
            } else {
                return deployAPIToK8sResult;
            }

        } else {
            BadRequestError badRequest = {body: {code: 90913, message: "Service from " + serviceKey + " not found."}};
            return badRequest;
        }
    }
    private isolated function deployAPIToK8s(model:APIArtifact apiArtifact) returns APKError|model:API {
        do {
            model:ConfigMap? definition = apiArtifact.definition;
            if definition is model:ConfigMap {
                http:Response deployConfigMapResult = check deployConfigMap(definition, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployConfigMapResult.statusCode == http:STATUS_CREATED {
                    log:printDebug("Deployed Configmap Successfully" + definition.toString());
                } else {
                    json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    _ = check self.handleK8sTimeout(statusResponse);
                }
            }
            foreach model:Service backendService in apiArtifact.backendServices {
                http:Response deployServiceResult = check deployService(backendService, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployServiceResult.statusCode == http:STATUS_CREATED {
                    log:printDebug("Deployed HttpRoute Successfully" + backendService.toString());
                } else {
                    json responsePayLoad = check deployServiceResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    _ = check self.handleK8sTimeout(statusResponse);
                }
            }
            model:Httproute? productionRoute = apiArtifact.productionRoute;
            if productionRoute is model:Httproute {
                http:Response deployHttpRouteResult = check deployHttpRoute(productionRoute, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployHttpRouteResult.statusCode == http:STATUS_CREATED {
                    log:printDebug("Deployed HttpRoute Successfully" + productionRoute.toString());
                } else {
                    json responsePayLoad = check deployHttpRouteResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    _ = check self.handleK8sTimeout(statusResponse);
                }
            }
            model:Httproute? sandboxRoute = apiArtifact.sandboxRoute;
            if sandboxRoute is model:Httproute {
                http:Response deployHttpRouteResult = check deployHttpRoute(sandboxRoute, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployHttpRouteResult.statusCode == http:STATUS_CREATED {
                    log:printDebug("Deployed HttpRoute Successfully" + sandboxRoute.toString());
                } else {
                    json responsePayLoad = check deployHttpRouteResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    _ = check self.handleK8sTimeout(statusResponse);

                }
            }
            foreach model:K8sServiceMapping k8sServiceMapping in apiArtifact.serviceMapping {
                http:Response deployServiceMappingCRResult = check deployServiceMappingCR(k8sServiceMapping, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployServiceMappingCRResult.statusCode == http:STATUS_CREATED {
                    log:printDebug("Deployed K8sAPI Successfully" + k8sServiceMapping.toString());
                } else {
                    json responsePayLoad = check deployServiceMappingCRResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    _ = check self.handleK8sTimeout(statusResponse);
                }
            }
            model:API? k8sAPI = apiArtifact.api;
            if k8sAPI is model:API {
                http:Response deployAPICRResult = check deployAPICR(k8sAPI, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployAPICRResult.statusCode == http:STATUS_CREATED {
                    json responsePayLoad = check deployAPICRResult.getJsonPayload();
                    log:printDebug("Deployed K8sAPI Successfully" + responsePayLoad.toJsonString());
                    return check responsePayLoad.cloneWithType(model:API);
                } else {
                    json responsePayLoad = check deployAPICRResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    _ = check self.handleK8sTimeout(statusResponse);
                    model:StatusDetails? details = statusResponse.details;
                    if details is model:StatusDetails {
                        model:StatusCause[] 'causes = details.'causes;
                        foreach model:StatusCause 'cause in 'causes {
                            if 'cause.'field == "spec.context" {
                                return error("Invalid API Context", code = 90911, description = "API Context " + k8sAPI.spec.context + " Invalid", message = "Invalid API context", statusCode = "400");

                            } else if 'cause.'field == "spec.apiDisplayName" {
                                return error("Invalid API Name", code = 90911, description = "API Name " + k8sAPI.spec.apiDisplayName + " Invalid", message = "Invalid API Name", statusCode = "400");
                            }
                        }
                    }
                    return error("Invalid API Request", code = 90911, description = "Invalid API Request", message = "Invalid API Request", statusCode = "400");
                }
            } else {
                return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = "500");
            }
        } on fail var e {
            log:printError("Internal Error occured while deploying API", e);
            APKError internalError = error("Internal Error occured while deploying API", code = 909000, statusCode = "500", description = "Internal Error occured while deploying API", message = "Internal Error occured while deploying API");
            return internalError;
        }
    }

    private isolated function retrieveGeneratedConfigmapForDefinition(model:APIArtifact apiArtifact, API api, json generatedSwaggerDefinition, string uniqueId) {
        map<string> configMapData = {};
        if api.'type == API_TYPE_HTTP {
            configMapData["openapi.json"] = generatedSwaggerDefinition.toJsonString();
        }
        model:ConfigMap configMap = {
            metadata: {
                name: self.retrieveDefinitionName(api, uniqueId),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: ()

            },
            data: configMapData
        };
        apiArtifact.definition = configMap;
    }

    isolated function setDefaultOperationsIfNotExist(API api) {
        APIOperations[]? operations = api.operations;
        boolean operationsAvailable = false;
        if operations is APIOperations[] {
            operationsAvailable = operations.length() > 0;
        }
        if operationsAvailable == false {
            APIOperations[] apiOperations = [];
            if api.'type == API_TYPE_HTTP {
                foreach string httpverb in HTTP_DEFAULT_METHODS {
                    APIOperations apiOperation = {target: "/*", verb: httpverb.toUpperAscii()};
                    apiOperations.push(apiOperation);
                }
                api.operations = apiOperations;
            }
        }
    }

    private isolated function generateAndSetAPICRArtifact(model:APIArtifact apiArtifact, API api) {
        model:API k8sAPI = {
            metadata: {
                name: apiArtifact.uniqueId,
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: ()
            },
            spec: {
                apiDisplayName: api.name,
                apiType: api.'type,
                apiVersion: api.'version,
                context: self.returnFullContext(api.context, api.'version),
                organization: "carbon.super"
            }
        };
        model:ConfigMap? definition = apiArtifact?.definition;
        if definition is model:ConfigMap {
            k8sAPI.spec.definitionFileRef = definition.metadata.name;
        }
        model:Httproute? productionRoute = apiArtifact.productionRoute;
        if productionRoute is model:Httproute {
            k8sAPI.spec.prodHTTPRouteRef = productionRoute.metadata.name;
        }
        model:Httproute? sandboxRoute = apiArtifact.sandboxRoute;
        if sandboxRoute is model:Httproute {
            k8sAPI.spec.sandHTTPRouteRef = sandboxRoute.metadata.name;
        }
        apiArtifact.api = k8sAPI;
    }

    private isolated function retrieveDefinitionName(API api, string uniqueId) returns string {
        return uniqueId + "-definition";
    }

    private isolated function retrieveHttpRouteRefName(API api, string uniqueId, string 'type) returns string {
        return uniqueId + "-" + 'type;
    }

    private isolated function setHttpRoute(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, string uniqueId, string endpointType) returns APKError? {
        model:Httproute httpRoute = {
            metadata:
                {
                name: self.retrieveHttpRouteRefName(api, uniqueId, endpointType),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: ()
            },
            spec: {
                parentRefs: self.generateAndRetrieveParentRefs(api, uniqueId),
                rules: check self.generateHttpRouteRules(apiArtifact, api, endpoint, endpointType),
                hostnames: self.getHostNames(api, uniqueId, endpointType)
            }
        };
        if endpointType == PRODUCTION_TYPE {
            apiArtifact.productionRoute = httpRoute;
        } else {
            apiArtifact.sandboxRoute = httpRoute;
        }
    }

    private isolated function getHostNames(API api, string unoqueId, string 'type) returns string[] {
        return ["gw.wso2.com"];
    }

    private isolated function generateAndRetrieveParentRefs(API api, string uniqueId) returns model:ParentReference[] {
        model:ParentReference[] parentRefs = [];
        model:ParentReference parentRef = {group: "gateway.networking.k8s.io", kind: "Gateway", name: "Default"};
        parentRefs.push(parentRef);
        return parentRefs;
    }

    private isolated function generateHttpRouteRules(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, string endpointType) returns model:HTTPRouteRule[]|APKError {
        model:HTTPRouteRule[] httpRouteRules = [];
        APIOperations[]? operations = api.operations;
        if operations is APIOperations[] {
            foreach APIOperations operation in operations {
                model:HTTPRouteRule httpRouteRule = check self.generateHttpRouteRule(apiArtifact, api, endpoint, operation, endpointType);
                httpRouteRules.push(httpRouteRule);
            }
        }
        return httpRouteRules;
    }

    private isolated function generateHttpRouteRule(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, APIOperations operation, string endpointType) returns APKError|model:HTTPRouteRule {
        do {
            record {}? endpointConfig = operation.endpointConfig;
            model:Endpoint endpointToUse = {};
            if endpointConfig is record {} {
                // endpointConfig presense at Operation Level.
                map<model:Endpoint> operationalLevelBackend = check self.createAndAddBackendServics(apiArtifact, api, endpointConfig, operation, endpointType);
                endpointToUse = operationalLevelBackend.get(endpointType);
            } else {
                if endpoint is model:Endpoint {
                    endpointToUse = endpoint;
                }
            }
            model:HTTPRouteRule httpRouteRule = {matches: self.retrieveMatches(api, operation), backendRefs: self.retrieveGeneratedBackend(api, endpointToUse, endpointType), filters: self.generateFilters(apiArtifact, api, endpointToUse, operation, endpointType)};
            return httpRouteRule;
        } on fail var e {
            log:printError("Internal Error occured", e);
            return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = "500");
        }
    }

    private isolated function generateFilters(model:APIArtifact apiArtifact, API api, model:Endpoint endpoint, APIOperations operation, string endpointType) returns model:HTTPRouteFilter[] {
        model:HTTPRouteFilter[] routeFilters = [];
        model:HTTPRouteFilter replacePathFilter = {'type: "URLRewrite", urlRewrite: {path: {'type: "ReplacePrefixMatch", replacePrefixMatch: self.generatePrefixMatch(api, endpoint, operation, endpointType)}}};
        routeFilters.push(replacePathFilter);
        if !endpoint.serviceEntry {
            model:HTTPRouteFilter backendHostHeaderModifierFilter = {'type: "RequestHeaderModifier", requestHeaderModifier: {add: [{name: "Host", value: self.gethost(<string>endpoint.url)}]}};
            routeFilters.push(backendHostHeaderModifierFilter);
        }
        return routeFilters;
    }

    isolated function generatePrefixMatch(API api, model:Endpoint endpoint, APIOperations operation, string endpointType) returns string {
        string target = operation.target ?: "/*";
        string generatedPath = "";
        if target == "/*" {
            generatedPath = "/";
        } else {
            string[] splitValues = regex:split(target, "/");
            foreach string pathPart in splitValues {
                if pathPart.indexOf("{", 0) >= 0 || pathPart.indexOf("*", 0) >= 0 {
                    break;
                }
                if pathPart.trim().length() > 0 {
                    generatedPath = generatedPath + "/" + pathPart;
                }
            }
        }
        if endpoint.serviceEntry {
            return generatedPath.trim();
        } else {
            string path = self.getPath(<string>endpoint.url);
            if path.endsWith("/") {
                path = path.substring(0, path.length() - 1);
            }
            if generatedPath.startsWith("/") {
                generatedPath = generatedPath.substring(1, generatedPath.length());
            }
            return path.trim() + "/" + generatedPath.trim();
        }
    }

    public isolated function retrievePathPrefix(string context, string 'version, string operation) returns string {
        string fullContext = self.returnFullContext(context, 'version);
        string[] splitValues = regex:split(operation, "/");
        string generatedPath = fullContext;
        if (operation == "/*") {
            return generatedPath;
        }
        foreach string pathPart in splitValues {
            if pathPart.trim().length() > 0 {
                // path contains path param
                if regex:matches(pathPart, "\\{.*\\}") {
                    generatedPath = generatedPath + "/" + regex:replaceAll(pathPart.trim(), "\\{.*\\}", ".*");
                } else {
                    generatedPath = generatedPath + "/" + pathPart;
                }
            }
        }

        if generatedPath.endsWith("/*") {
            int lastSlashIndex = <int>generatedPath.lastIndexOf("/", generatedPath.length());
            generatedPath = generatedPath.substring(0, lastSlashIndex);
        }
        return generatedPath.trim();
    }

    private isolated function retrieveGeneratedBackend(API api, model:Endpoint endpoint, string endpointType) returns model:HTTPBackendRef[] {
        model:HTTPBackendRef httpBackend = {
            namespace:
            <string>endpoint.namespace,
            kind: "Service",
            weight: 1,
            port: <int>endpoint.port,
            name: <string>endpoint.name,
            group: ""
        };
        return [httpBackend];
    }

    private isolated function retrievePort(Service serviceEntry) returns int {
        PortMapping[]? portmappings = serviceEntry.portmapping;
        if portmappings is PortMapping[] {
            if portmappings.length() > 0 {
                return portmappings[0].targetport;
            }
        }

        return 80;
    }

    private isolated function retrieveMatches(API api, APIOperations apiOperation) returns model:HTTPRouteMatch[] {
        model:HTTPRouteMatch[] httpRouteMatch = [];
        model:HTTPRouteMatch httpRoute = {method: <string>apiOperation.verb, path: {'type: "PathPrefix", value: self.retrievePathPrefix(api.context, api.'version, apiOperation.target ?: "/*")}};
        httpRouteMatch.push(httpRoute);
        return httpRouteMatch;
    }

    isolated function retrieveGeneratedSwaggerDefinition(API api) returns json|APKError {
        runtimeModels:API api1 = runtimeModels:newAPI1();
        api1.setName(api.name);
        api1.setType(api.'type);
        api1.setVersion(api.'version);
        utilapis:Set uritemplatesSet = utilapis:newHashSet1();
        if api.operations is APIOperations[] {
            foreach APIOperations apiOperation in <APIOperations[]>api.operations {
                runtimeModels:URITemplate uriTemplate = runtimeModels:newURITemplate1();
                uriTemplate.setUriTemplate(<string>apiOperation.target);
                string? verb = apiOperation.verb;
                if verb is string {
                    uriTemplate.setHTTPVerb(verb.toUpperAscii());
                }
                _ = uritemplatesSet.add(uriTemplate);
            }
        }
        api1.setUriTemplates(uritemplatesSet);
        string?|runtimeapi:APIManagementException retrievedDefinition = runtimeUtil:RuntimeAPICommonUtil_generateDefinition(api1);
        if retrievedDefinition is string {
            json|error jsonString = value:fromJsonString(retrievedDefinition);
            if jsonString is json {
                return jsonString;
            } else {
                log:printError("Error on converting to json", jsonString);
                return error("Error occured while generating openapi definition", code = 900920, message = "Error occured while generating openapi definition", statusCode = "500", description = "Error occured while generating openapi definition");
            }
        } else if retrievedDefinition is () {
            return "";
        } else {
            return error("Error occured while generating openapi definition", code = 900920, message = "Error occured while generating openapi definition", statusCode = "500", description = "Error occured while generating openapi definition");
        }
    }

    public isolated function generateAPIKey(string apiId) returns APIKey|BadRequestError|NotFoundError|InternalServerErrorError {
        model:API|error api = getAPI(apiId);
        if api is model:API {
            InternalTokenGenerator tokenGenerator = new ();
            string|jwt:Error generatedToken = tokenGenerator.generateToken(api, APK_USER);
            if generatedToken is string {
                APIKey apiKey = {apikey: generatedToken, validityTime: <int>runtimeConfiguration.tokenIssuerConfiguration.expTime};
                return apiKey;
            } else {
                log:printError("Error while Genereting token for API : " + apiId, generatedToken);
                InternalServerErrorError internalError = {body: {code: 90911, message: "Error while Generating Token"}};
                return internalError;
            }
        } else {
            NotFoundError notfound = {body: {code: 909100, message: apiId + "not found."}};
            return notfound;
        }
    }

    public function retrieveAllApisAtStartup(string? continueValue) returns error? {
        string? resultValue = continueValue;
        json|http:ClientError retrieveAllAPISResult;
        if resultValue is string {
            retrieveAllAPISResult = retrieveAllAPIS(resultValue);
        } else {
            retrieveAllAPISResult = retrieveAllAPIS(());
        }

        if retrieveAllAPISResult is json {
            json metadata = check retrieveAllAPISResult.metadata;
            json[] items = <json[]>check retrieveAllAPISResult.items;
            putallAPIS(items);

            json|error continueElement = metadata.'continue;
            if continueElement is json {
                if (<string>continueElement).length() > 0 {
                    _ = check self.retrieveAllApisAtStartup(<string?>continueElement);
                }
            }
            string resourceVersion = <string>check metadata.'resourceVersion;
            setResourceVersion(resourceVersion);
        }
    }

    isolated function generateAndSetK8sServiceMapping(model:APIArtifact apiArtifact, API api, Service serviceEntry, string namespace) {
        model:API? k8sAPI = apiArtifact.api;
        if k8sAPI is model:API {
            model:K8sServiceMapping k8sServiceMapping = {
                metadata: {
                    name: self.getServiceMappingEntryName(apiArtifact.uniqueId),
                    namespace: namespace,
                    uid: (),
                    creationTimestamp: ()
                },
                spec: {
                    serviceRef: {
                        namespace: serviceEntry.namespace,
                        name: serviceEntry.name
                    },
                    apiRef: {
                        namespace: k8sAPI.metadata.namespace,
                        name: k8sAPI.metadata.name
                    }
                }
            };
            apiArtifact.serviceMapping.push(k8sServiceMapping);
        }
    }

    isolated function getServiceMappingEntryName(string uniqueId) returns string {
        return uniqueId + "-servicemapping";
    }

    isolated function deleteServiceMappings(model:API api) returns APKError? {
        do {
            model:K8sServiceMapping[] retrieveServiceMappingsForAPIResult = retrieveServiceMappingsForAPI(api);
            foreach model:K8sServiceMapping serviceMapping in retrieveServiceMappingsForAPIResult {
                http:Response|http:ClientError k8ServiceMappingDeletionResponse = deleteK8ServiceMapping(serviceMapping.metadata.name, serviceMapping.metadata.namespace);
                if k8ServiceMappingDeletionResponse is http:Response {
                    if k8ServiceMappingDeletionResponse.statusCode != http:STATUS_OK {
                        json responsePayLoad = check k8ServiceMappingDeletionResponse.getJsonPayload();
                        model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                        _ = check self.handleK8sTimeout(statusResponse);
                    }
                } else {
                    log:printError("Error occured while deleting service mapping");
                }
            }
            return;
        } on fail var e {
            log:printError("Error occured deleting servicemapping", e);
            return error("Error occured deleting servicemapping", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = "500");
        }

    }

    public isolated function validateDefinition(http:Request message, boolean returnContent) returns InternalServerErrorError|error|BadRequestError|APIDefinitionValidationResponse {
        DefinitionValidationRequest|BadRequestError|error definitionValidationRequest = self.mapApiDefinitionPayload(message);
        if definitionValidationRequest is DefinitionValidationRequest {
            boolean inlineApiDefinitionAvailable = definitionValidationRequest.inlineAPIDefinition is string;
            boolean fileAvailable = definitionValidationRequest.fileName is string && definitionValidationRequest.content is byte[];
            boolean urlAvailble = definitionValidationRequest.url is string;
            boolean typeAvailable = definitionValidationRequest.'type.length() > 0;
            if (!inlineApiDefinitionAvailable && !fileAvailable && !urlAvailble && !typeAvailable) {
                BadRequestError badRequest = {body: {code: 90914, message: "atleast one of the field required"}};
                return badRequest;
            }
            runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException validationResponse;
            if !typeAvailable {
                BadRequestError badRequest = {body: {code: 90914, message: "type field unavailable"}};
                return badRequest;
            }
            string? url = definitionValidationRequest.url;
            if url is string {
                if (fileAvailable || inlineApiDefinitionAvailable) {
                    BadRequestError badRequest = {body: {code: 90914, message: "multiple fields of  url,file,inlineAPIDefinition given"}};
                    return badRequest;
                }
                string|error retrieveDefinitionFromUrlResult = self.retrieveDefinitionFromUrl(url);
                if retrieveDefinitionFromUrlResult is string {
                    validationResponse = check runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition(definitionValidationRequest.'type, [], retrieveDefinitionFromUrlResult, definitionValidationRequest.fileName ?: "", returnContent);
                } else {
                    log:printError("Error occured while retrieving definition from url", retrieveDefinitionFromUrlResult);
                    BadRequestError badRequest = {body: {code: 900900, message: "retrieveDefinitionFromUrlResult"}};
                    return badRequest;
                }
            } else if definitionValidationRequest.fileName is string && definitionValidationRequest.content is byte[] {
                if (urlAvailble || inlineApiDefinitionAvailable) {
                    BadRequestError badRequest = {body: {code: 90914, message: "multiple fields of  url,file,inlineAPIDefinition given"}};
                    return badRequest;
                }
                validationResponse = check runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition(definitionValidationRequest.'type, <byte[]>definitionValidationRequest.content, "", <string>definitionValidationRequest.fileName, returnContent);
            } else if definitionValidationRequest.inlineAPIDefinition is string {
                if (fileAvailable || urlAvailble) {
                    BadRequestError badRequest = {body: {code: 90914, message: "multiple fields of  url,file,inlineAPIDefinition given"}};
                    return badRequest;
                }
                validationResponse = check runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition(definitionValidationRequest.'type, [], <string>definitionValidationRequest.inlineAPIDefinition, "", returnContent);
            } else {
                BadRequestError badRequest = {body: {code: 90914, message: "atleast one of the field required"}};
                return badRequest;
            }
            if validationResponse is runtimeapi:APIDefinitionValidationResponse {

                string[] endpoints = [];
                ErrorListItem[] errorItems = [];
                if validationResponse.isValid() {
                    runtimeapi:Info info = validationResponse.getInfo();
                    utilapis:List endpointList = info.getEndpoints();
                    foreach int i in 0 ... endpointList.size() - 1 {
                        endpoints.push(endpointList.get(i).toString());
                    }
                    APIDefinitionValidationResponse_info validationResponseInfo = {
                        context: info.getContext(),
                        description: info.getDescription(),
                        name: info.getName(),
                        'version: info.getVersion(),
                        openAPIVersion: info.getOpenAPIVersion(),
                        endpoints: endpoints
                    };
                    APIDefinitionValidationResponse response = {content: validationResponse.getContent(), isValid: validationResponse.isValid(), info: validationResponseInfo, errors: errorItems};
                    return response;
                }
                utilapis:ArrayList errorItemsResult = validationResponse.getErrorItems();
                foreach int i in 0 ... errorItemsResult.size() - 1 {
                    runtimeapi:ErrorItem errorItem = check java:cast(errorItemsResult.get(i));
                    ErrorListItem errorListItem = {code: errorItem.getErrorCode().toString(), message: <string>errorItem.getErrorMessage(), description: errorItem.getErrorDescription()};
                    errorItems.push(errorListItem);
                }
                APIDefinitionValidationResponse response = {content: validationResponse.getContent(), isValid: validationResponse.isValid(), info: {}, errors: errorItems};
                return response;

            } else {
                runtimeapi:JAPIManagementException excetion = check validationResponse.ensureType(runtimeapi:JAPIManagementException);
                runtimeapi:ErrorHandler errorHandler = excetion.getErrorHandler();
                BadRequestError badeRequest = {body: {code: errorHandler.getErrorCode(), message: errorHandler.getErrorMessage().toString()}};
                return badeRequest;
            }
        }
        else if definitionValidationRequest is BadRequestError {
            return definitionValidationRequest;
        }

        else {
            InternalServerErrorError internalError = {body: {code: 909000, message: "InternalServerError"}};
            return internalError;
        }
    }

    private isolated function mapApiDefinitionPayload(http:Request message) returns DefinitionValidationRequest|BadRequestError|error {
        string|() url = ();
        string|() fileName = ();
        byte[]|() fileContent = ();
        string|() definitionType = ();
        string|() inlineAPIDefinition = ();
        mime:Entity[]|http:ClientError payLoadParts = message.getBodyParts();
        if payLoadParts is mime:Entity[] {
            foreach mime:Entity payLoadPart in payLoadParts {
                mime:ContentDisposition contentDisposition = payLoadPart.getContentDisposition();
                string fieldName = contentDisposition.name;
                if fieldName == "url" {
                    url = check payLoadPart.getText();
                }
                else if fieldName == "file" {
                    fileName = contentDisposition.fileName;
                    fileContent = check payLoadPart.getByteArray();
                } else if fieldName == "type" {
                    definitionType = check payLoadPart.getText();
                } else if fieldName == "inlineAPIDefinition" {
                    inlineAPIDefinition = check payLoadPart.getText();
                }
            }
        }
        if definitionType is () {
            BadRequestError badeRequest = {body: {code: 90914, message: "type not specified in Request"}};
            return badeRequest;
        }
        return {
            content: fileContent,
            fileName: fileName,
            inlineAPIDefinition: inlineAPIDefinition,
            url: url,
            'type: definitionType
        };
    }

    private isolated function retrieveDefinitionFromUrl(string url) returns string|error {
        string domain = self.getDomain(url);
        string path = self.getPath(url);
        if domain.length() > 0 {
            http:Client httpClient = check new (domain);
            http:Response response = check httpClient->get(path, targetType = http:Response);
            return response.getTextPayload();
        } else {
            return error("invalid url " + url);
        }
    }

    isolated function getDomain(string url) returns string {
        string hostPort = "";
        string protocol = "";
        if url.startsWith("https://") {
            hostPort = url.substring(8, url.length());
            protocol = "https";
        } else if url.startsWith("http://") {
            hostPort = url.substring(7, url.length());
            protocol = "http";
        } else {
            return "";
        }
        int? indexOfSlash = hostPort.indexOf("/", 0);
        if indexOfSlash is int {
            return protocol + "://" + hostPort.substring(0, indexOfSlash);
        } else {
            return protocol + "://" + hostPort;
        }
    }

    isolated function gethost(string url) returns string {
        string host = "";
        if url.startsWith("https://") {
            host = url.substring(8, url.length());
        } else if url.startsWith("http://") {
            host = url.substring(7, url.length());
        } else {
            return "";
        }
        int? indexOfColon = host.indexOf(":", 0);
        if indexOfColon is int {
            return host.substring(0, indexOfColon);
        } else {
            int? indexOfSlash = host.indexOf("/", 0);
            if indexOfSlash is int {
                return host.substring(0, indexOfSlash);
            } else {
                return host;
            }
        }
    }
    isolated function getPort(string url) returns int|error {
        string hostPort = "";
        string protocol = "";
        if url.startsWith("https://") {
            hostPort = url.substring(8, url.length());
            protocol = "https";
        } else if url.startsWith("http://") {
            hostPort = url.substring(7, url.length());
            protocol = "http";
        } else {
            return -1;
        }
        int? indexOfSlash = hostPort.indexOf("/", 0);

        if indexOfSlash is int {
            hostPort = hostPort.substring(0, indexOfSlash);
        }
        int? indexOfColon = hostPort.indexOf(":");
        if indexOfColon is int {
            string port = hostPort.substring(indexOfColon + 1, hostPort.length());
            return check int:fromString(port);
        } else {
            if protocol == "https" {
                return 443;
            } else {
                return 80;
            }
        }
    }

    isolated function getPath(string url) returns string {
        string hostPort = "";
        if url.startsWith("https://") {
            hostPort = url.substring(8, url.length());
        } else if url.startsWith("http://") {
            hostPort = url.substring(7, url.length());
        } else {
            return "";
        }
        int? indexOfSlash = hostPort.indexOf("/", 0);
        if indexOfSlash is int {
            return hostPort.substring(indexOfSlash, hostPort.length());
        } else {
            return "";
        }
    }

    isolated function handleK8sTimeout(model:Status errorStatus) returns APKError? {
        do {
            model:StatusDetails? details = errorStatus.details;
            if details is model:StatusDetails {
                if details.retryAfterSeconds is int && details.retryAfterSeconds >= 0 {
                    // K8s api level ratelimit hit.
                    log:printError("Error while deploying API");
                    APKError apkError = error("Error while deploying API", code = 90915, message = "Internal Error while Deploying API", statusCode = "500", description = "Internal Error while Deploying API");
                    return apkError;
                } else {
                    return;
                }
            }
        } on fail var e {
            log:printError("Internal Error occcured on checking K8sTimeout", e);
            return error("Internal Error occcured on checking K8sTimeout", code = 909000, statusCode = "500", message = "Internal Server Error", description = "Internal Server Error");
        }

    }
    isolated function createBackendService(string url, map<string> labels) returns model:Service {
        string nameSpace = getNameSpace(runtimeConfiguration.apiCreationNamespace);
        model:Service backendService = {
            metadata: {
                name: "backend-" + uuid:createType1AsString(),
                namespace: nameSpace,
                uid: (),
                creationTimestamp: (),
                labels: labels
            },
            spec: {
                externalName: self.gethost(url),
                'type: "ExternalName"
            }
        };
        return backendService;
    }
    public isolated function retrieveDefaultDefinition(model:API api) returns json {
        json defaultOpenApiDefinition = {
            "openapi": "3.0.1",
            "info": {
                "title": api.spec.apiDisplayName,
                "version": api.spec.apiVersion
            },
            "servers": [
                {
                    "url": "/"
                }
            ],
            "security": [
                {
                    "default": []
                }
            ],
            "paths": {
                "/*": {
                    "get": {
                        "responses": {
                            "200": {
                                "description": "OK"
                            }
                        },
                        "security": [
                            {
                                "default": []
                            }
                        ],
                        "x-auth-type": "Application & Application User",
                        "x-throttling-tier": "Unlimited",
                        "x-wso2-application-security": {
                            "security-types": [
                                "oauth2"
                            ],
                            "optional": false
                        }
                    },
                    "put": {
                        "responses": {
                            "200": {
                                "description": "OK"
                            }
                        },
                        "security": [
                            {
                                "default": []
                            }
                        ],
                        "x-auth-type": "Application & Application User",
                        "x-throttling-tier": "Unlimited",
                        "x-wso2-application-security": {
                            "security-types": [
                                "oauth2"
                            ],
                            "optional": false
                        }
                    },
                    "post": {
                        "responses": {
                            "200": {
                                "description": "OK"
                            }
                        },
                        "security": [
                            {
                                "default": []
                            }
                        ],
                        "x-auth-type": "Application & Application User",
                        "x-throttling-tier": "Unlimited",
                        "x-wso2-application-security": {
                            "security-types": [
                                "oauth2"
                            ],
                            "optional": false
                        }
                    },
                    "delete": {
                        "responses": {
                            "200": {
                                "description": "OK"
                            }
                        },
                        "security": [
                            {
                                "default": []
                            }
                        ],
                        "x-auth-type": "Application & Application User",
                        "x-throttling-tier": "Unlimited",
                        "x-wso2-application-security": {
                            "security-types": [
                                "oauth2"
                            ],
                            "optional": false
                        }
                    },
                    "patch": {
                        "responses": {
                            "200": {
                                "description": "OK"
                            }
                        },
                        "security": [
                            {
                                "default": []
                            }
                        ],
                        "x-auth-type": "Application & Application User",
                        "x-throttling-tier": "Unlimited",
                        "x-wso2-application-security": {
                            "security-types": [
                                "oauth2"
                            ],
                            "optional": false
                        }
                    }
                }
            },
            "components": {
                "securitySchemes": {
                    "default": {
                        "type": "oauth2",
                        "flows": {
                            "implicit": {
                                "authorizationUrl": "https://test.com",
                                "scopes": {}
                            }
                        }
                    }
                }
            }
        };
        return defaultOpenApiDefinition;
    }
    public isolated function validateAPIExistence(string query) returns NotFoundError|BadRequestError|http:Ok {
        int? indexOfColon = query.indexOf(":", 0);
        boolean exist = false;
        if indexOfColon is int && indexOfColon > 0 {
            string keyWord = query.substring(0, indexOfColon);
            string keyWordValue = query.substring(keyWord.length() + 1, query.length());
            if keyWord == "name" {
                exist = self.validateName(keyWordValue);
            } else if keyWord == "context" {
                exist = self.validateContext(keyWordValue);
            } else {
                BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord " + keyWord}};
                return badRequest;
            }
        } else {
            // Consider full string as name;
            exist = self.validateName(query);
        }
        if exist {
            http:Ok ok = {};
            return ok;
        } else {
            NotFoundError notFound = {body: {code: 900914, message: "context/name doesn't exist"}};
            return notFound;
        }
    }
}

type DefinitionValidationRequest record {|
    string? url;
    string? fileName;
    byte[]? content;
    string? inlineAPIDefinition;
    string 'type;

|};

