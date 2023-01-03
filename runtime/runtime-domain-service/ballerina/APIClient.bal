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
import runtime_domain_service.org.wso2.apk.runtime.api as runtimeapi;

public class APIClient {

    public isolated function getAPIDefinitionByID(string id) returns string|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        model:API|error api = getAPI(id);
        if api is model:API {
            string? definitionFileRef = api.spec.definitionFileRef;
            if definitionFileRef is string {
                if definitionFileRef.length() > 0 {
                    string|error definition = self.getDefinition(api);
                    if definition is string {
                        return definition;
                    } else {
                        log:printError("Error while reading definition:", definition);
                    }
                }
            }
        }
        NotFoundError notfound = {body: {code: 909100, message: id + "not found."}};
        return notfound;
    }

    private isolated function getDefinition(model:API api) returns string|error {
        string? definitionFileRef = api.spec.definitionFileRef;
        if definitionFileRef is string {
            json|error configMapValue = getConfigMapValueFromNameAndNamespace(definitionFileRef, api.metadata.namespace);
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
        } else {
            return "";
        }
    }

    //Get APIs deployed in default namespace by APIId.
    public isolated function getAPIById(string id) returns API|NotFoundError|InternalServerErrorError {
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
        NotFoundError notfound = {body: {code: 909100, message: id + "not found."}};
        return notfound;
    }

    //Delete APIs deployed in a namespace by APIId.
    public isolated function deleteAPIById(string id) returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|error {
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
                InternalServerErrorError? response = check self.deleteServiceMappings(api);
                if response is InternalServerErrorError {
                    return response;
                }
            } else {
                NotFoundError apiNotfound = {body: {code: 900910, description: "API with " + id + " not found", message: "API not found"}};
                return apiNotfound;
            }
        }
        return http:OK;
    }

    //Get all deployed APIs in namespace with specific search query
    public function getAPIListInNamespaceWithQuery(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc") returns APIList|InternalServerErrorError|BadRequestError|error {
        lock {
            APIInfo[] apiNames = map:toArray(apilist);
            return {list: apiNames.clone(), count: apiNames.length(), pagination: {total: apilist.length()}};
        }
    }

    # This returns list of APIS.
    #
    # + query - Parameter Description  
    # + 'limit - Parameter Description  
    # + offset - Parameter Description  
    # + sortBy - Parameter Description  
    # + sortOrder - Parameter Description
    # + return - Return list of APIS in namsepace.
    public isolated function getAPIList(string? query, int 'limit, int offset, string sortBy, string sortOrder) returns APIList|InternalServerErrorError {
        API[] apilist = [];
        foreach model:API api in getAPIs() {
            API convertedModel = convertK8sAPItoAPI(api);
            apilist.push(convertedModel);
        }
        if query is string {
            return self.filterAPISBasedOnQuery(apilist, query, 'limit, offset, sortBy, sortOrder);
        } else {
            return self.filterAPIS(apilist, 'limit, offset, sortBy, sortOrder);
        }
    }
    private isolated function filterAPISBasedOnQuery(API[] apilist, string query, int 'limit, int offset, string sortBy, string sortOrder) returns APIList|InternalServerErrorError {
        API[] filteredList = [];
        if query.length() > 0 {
            int? semiCollonIndex = string:indexOf(query, ":", 0);
            if semiCollonIndex is int {
                if semiCollonIndex > 0 {
                    string keyWord = query.substring(0, semiCollonIndex);
                    string keyWordValue = query.substring(keyWord.length() + 1, query.length());
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
                        // BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord " + keyWord}};
                        // return badRequest;
                    }
                }
            } else {
                foreach API api in apilist {
                    if (regex:matches(api.name, query)) {
                        filteredList.push(api);
                    }
                }
            }
        } else {
            filteredList = apilist;
        }
        return self.filterAPIS(filteredList, 'limit, offset, sortBy, sortOrder);
    }
    private isolated function filterAPIS(API[] apiList, int 'limit, int offset, string sortBy, string sortOrder) returns APIList|InternalServerErrorError {
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
            // BadRequestError badRequest = {body: {code: 90912, message: "Invalid Sort By/Sort Order Value "}};
            // return badRequest;
        }
        API[] limitSet = [];
        if sortedAPIS.length() >= offset {
            foreach int i in offset ... (sortedAPIS.length() - 1) {
                if limitSet.length() < 'limit {
                    limitSet.push(sortedAPIS[i]);
                }
            }
        }
        return {list: limitSet, count: limitSet.length(), pagination: {total: apiList.length(), 'limit: 'limit, offset: offset}};

    }
    public function createAPI(API api) returns string|Error {
        if (self.validateName(api.name)) {
            return {code: 90911, message: "API Name - " + api.name + " already exist.", description: "API Name - " + api.name + " already exist."};
        }
        if self.validateContextAndVersion(api.context, api.'version) {
            return {code: 90912, message: "API Context - " + api.context + " already exist.", description: "API Context - " + api.context + "already exist."};
        }
        return "created";
    }

    isolated function validateContextAndVersion(string context, string 'version) returns boolean {

        foreach model:API k8sAPI in getAPIs() {
            if k8sAPI.spec.context == self.returnFullContext(context, 'version) {
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

    function createAndDeployAPI(API api) {
        model:API k8sAPI = self.convertK8sCrAPI(api);
        log:printInfo(<string>k8sAPI.toJson());
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

    isolated function createAPIFromService(string serviceKey, API api) returns CreatedAPI|BadRequestError|InternalServerErrorError|error {
        if (self.validateName(api.name)) {
            BadRequestError badRequest = {body: {code: 90911, message: "API Name - " + api.name + " already exist.", description: "API Name - " + api.name + " already exist."}};
            return badRequest;
        }
        if self.validateContextAndVersion(api.context, api.'version) {
            BadRequestError badRequest = {body: {code: 90911, message: "API Context - " + api.context + " already exist.", description: "API Context " + api.context + " already exist."}};
            return badRequest;
        }
        self.setDefaultOperationsIfNotExist(api);
        Service|error serviceRetrieved = grtServiceById(serviceKey);
        string uniqueId = uuid:createType1AsString();
        if serviceRetrieved is Service {
            model:Httproute prodHttpRoute = self.retrieveHttpRoute(api, serviceRetrieved, uniqueId, "production");
            model:API k8sAPI = self.generateAPICRArtifact(api, (), prodHttpRoute, uniqueId);
            model:K8sServiceMapping k8sServiceMapping = self.generateK8sServiceMapping(k8sAPI, serviceRetrieved, getNameSpace(runtimeConfiguration.apiCreationNamespace), uniqueId);
            string|error generatedSwaggerDefinition = self.retrieveGeneratedSwaggerDefinition(api);
            model:ConfigMap definitionConnfigMap = self.retrieveGeneratedConfigmapForDefinition(api, check generatedSwaggerDefinition, uniqueId);
            http:Response deployConfigMapResult = check deployConfigMap(definitionConnfigMap, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployConfigMapResult.statusCode == 201 {
                log:printDebug("Deployed Configmap Successfully" + definitionConnfigMap.toString());
            } else {
                InternalServerErrorError? handleK8sTimeoutResult = check self.handleK8sTimeout(deployConfigMapResult);
                if handleK8sTimeoutResult is InternalServerErrorError {
                    return handleK8sTimeoutResult;
                }
            }
            http:Response deployHttpRouteResult = check deployHttpRoute(prodHttpRoute, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployHttpRouteResult.statusCode == 201 {
                log:printDebug("Deployed HttpRoute Successfully" + prodHttpRoute.toString());
            } else {
                InternalServerErrorError? handleK8sTimeoutResult = check self.handleK8sTimeout(deployHttpRouteResult);
                if handleK8sTimeoutResult is InternalServerErrorError {
                    return handleK8sTimeoutResult;
                }
            }
            http:Response deployServiceMappingCRResult = check deployServiceMappingCR(k8sServiceMapping, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployServiceMappingCRResult.statusCode == 201 {
                log:printDebug("Deployed K8sAPI Successfully" + k8sServiceMapping.toString());
            } else {
                InternalServerErrorError? handleK8sTimeoutResult = check self.handleK8sTimeout(deployServiceMappingCRResult);
                if handleK8sTimeoutResult is InternalServerErrorError {
                    return handleK8sTimeoutResult;
                }
            }
            http:Response deployAPICRResult = check deployAPICR(k8sAPI, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployAPICRResult.statusCode == http:STATUS_CREATED {
                json responsePayLoad = check deployAPICRResult.getJsonPayload();
                log:printDebug("Deployed K8sAPI Successfully" + responsePayLoad.toJsonString());
                model:API createdK8sAPI = check responsePayLoad.cloneWithType(model:API);
                CreatedAPI createdAPI = {body: {name: api.name, context: self.returnFullContext(api.context, api.'version), 'version: api.'version, id: createdK8sAPI.metadata.uid}};
                return createdAPI;
            } else {
                InternalServerErrorError? handleK8sTimeoutResult = check self.handleK8sTimeout(deployAPICRResult);
                if handleK8sTimeoutResult is InternalServerErrorError {
                    return handleK8sTimeoutResult;
                }
                json responsePayLoad = check deployAPICRResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                model:StatusDetails details = statusResponse.details;
                model:StatusCause[] 'causes = details.'causes;

                if details.retryAfterSeconds is int && details.retryAfterSeconds >= 0 {
                    // K8s api level ratelimit hit.
                    log:printError("Error while deploying API");
                    InternalServerErrorError internalEror = {body: {code: 90915, message: "Internal Error while Deploying API"}};
                    return internalEror;
                }
                foreach model:StatusCause 'cause in 'causes {
                    if 'cause.'field == "spec.context" {
                        BadRequestError badRequest = {body: {code: 90911, message: "API Context - " + api.context + " Invalid.", description: "API Context " + api.context + " Invalid."}};
                        return badRequest;

                    } else if 'cause.'field == "spec.apiDisplayName" {
                        BadRequestError badRequest = {body: {code: 90911, message: "API Name - " + api.name + " Invalid.", description: "API Name " + api.name + " Invalid."}};
                        return badRequest;
                    }
                }
                BadRequestError badRequest = {body: {code: 90911, message: "Invalid API Request", description: "Invalid API Request"}};
                return badRequest;
            }
        } else {
            BadRequestError badRequest = {body: {code: 90913, message: "Service from " + serviceKey + " not found."}};
            return badRequest;
        }
    }

    private isolated function retrieveGeneratedConfigmapForDefinition(API api, string generatedSwaggerDefinition, string uniqueId) returns model:ConfigMap {
        map<string> configMapData = {};
        if api.'type == API_TYPE_HTTP {
            configMapData["openapi.json"] = generatedSwaggerDefinition;
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
        return configMap;
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

    private isolated function generateAPICRArtifact(API api, model:Httproute? sandboxHttp, model:Httproute? prodHttp, string uniqueId) returns model:API {
        model:API k8sAPI = {
            metadata: {
                name: uniqueId,
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: ()
            },
            spec: {
                apiDisplayName: api.name,
                apiType: api.'type,
                apiVersion: api.'version,
                context: self.returnFullContext(api.context, api.'version),
                definitionFileRef: self.retrieveDefinitionName(api, uniqueId),
                organization: "carbon.super"
            }
        };
        if (prodHttp is model:Httproute) {
            k8sAPI.spec.prodHTTPRouteRef = self.retrieveHttpRouteRefName(api, uniqueId, "production");
        }
        if (sandboxHttp is model:Httproute) {
            k8sAPI.spec.sandHTTPRouteRef = self.retrieveHttpRouteRefName(api, uniqueId, "sandbox");
        }
        return k8sAPI;
    }

    private isolated function retrieveDefinitionName(API api, string uniqueId) returns string {
        return uniqueId + "-definition";
    }

    private isolated function retrieveHttpRouteRefName(API api, string uniqueId, string 'type) returns string {
        return uniqueId + "-" + 'type;
    }

    private isolated function retrieveHttpRoute(API api, Service? serviceEntry, string uniqueId, string endpointType) returns model:Httproute {
        model:Httproute httpRoute = {
            metadata:
                {
                name: self.retrieveHttpRouteRefName(api, uniqueId, endpointType),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: ()
            },
            spec: {
                parentRefs: self.generateAndRetrieveParentRefs(api, serviceEntry, uniqueId),
                rules: self.generateHttpRouteRules(api, serviceEntry, endpointType),
                hostnames: self.getHostNames(api, uniqueId, endpointType)
            }
        };
        return httpRoute;
    }

    private isolated function getHostNames(API api, string unoqueId, string 'type) returns string[] {
        return ["gw.wso2.com"];
    }

    private isolated function generateAndRetrieveParentRefs(API api, Service? serviceEntry, string uniqueId) returns model:ParentReference[] {
        model:ParentReference[] parentRefs = [];
        model:ParentReference parentRef = {group: "gateway.networking.k8s.io", kind: "Gateway", name: "Default"};
        parentRefs.push(parentRef);
        return parentRefs;
    }

    private isolated function generateHttpRouteRules(API api, Service? serviceEntry, string endpointType) returns model:HTTPRouteRule[] {
        model:HTTPRouteRule[] httpRouteRules = [];
        APIOperations[]? operations = api.operations;
        if operations is APIOperations[] {
            foreach APIOperations operation in operations {
                model:HTTPRouteRule httpRouteRule = self.generateHttpRouteRule(api, serviceEntry, operation, endpointType);
                httpRouteRules.push(httpRouteRule);
            }
        }
        return httpRouteRules;
    }

    private isolated function generateHttpRouteRule(API api, Service? serviceEntry, APIOperations operation, string endpointType) returns model:HTTPRouteRule {
        model:HTTPRouteRule httpRouteRule = {matches: self.retrieveMatches(api, operation), backendRefs: self.retrieveGeneratedBackend(api, serviceEntry, endpointType), filters: self.generateFilters(api, serviceEntry, operation, endpointType)};
        return httpRouteRule;
    }

    private isolated function generateFilters(API api, Service? serviceEntry, APIOperations operation, string endpointType) returns model:HTTPRouteFilter[] {
        model:HTTPRouteFilter[] routeFilters = [];
        model:HTTPRouteFilter replacePathFilter = {'type: "URLRewrite", urlRewrite: {path: {'type: "ReplacePrefixMatch", replacePrefixMatch: self.generatePrefixMatch(api, serviceEntry, operation, endpointType)}}};
        routeFilters.push(replacePathFilter);
        return routeFilters;
    }

    private isolated function generatePrefixMatch(API api, Service? serviceEntry, APIOperations operation, string endpointType) returns string {
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
        if serviceEntry is Service {
            return generatedPath.trim();
        }
        return generatedPath.trim();
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

    private isolated function retrieveGeneratedBackend(API api, Service? serviceEntry, string endpointType) returns model:HTTPBackendRef[] {
        if serviceEntry is Service {
            model:HTTPBackendRef httpBackend = {
                namespace:
            serviceEntry.namespace,
                kind: "Service",
                weight: 1,
                port: self.retrievePort(serviceEntry),
                name: serviceEntry.name,
                group: ""
            };
            return [httpBackend];

        } else {
            //TODO tharindua@wso2.com need to write once resource level endpoint came.
            return [{port: 0, kind: "", name: "", namespace: "", weight: 0, group: ""}];
        }
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

    private isolated function retrieveGeneratedSwaggerDefinition(API api) returns string|error {
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
            return retrievedDefinition;
        } else if retrievedDefinition is () {
            return "";
        } else {
            return error(retrievedDefinition.message());
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

    isolated function generateK8sServiceMapping(model:API api, Service serviceEntry, string namespace, string uniqueId) returns model:K8sServiceMapping {
        model:K8sServiceMapping k8sServiceMapping = {
            metadata: {
                name: self.getServiceMappingEntryName(uniqueId),
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
                    namespace: api.metadata.namespace,
                    name: api.metadata.name
                }
            }
        };
        return k8sServiceMapping;

    }

    isolated function getServiceMappingEntryName(string uniqueId) returns string {
        return uniqueId + "-servicemapping";
    }

    isolated function deleteServiceMappings(model:API api) returns InternalServerErrorError?|error {
        model:K8sServiceMapping[] retrieveServiceMappingsForAPIResult = retrieveServiceMappingsForAPI(api);
        foreach model:K8sServiceMapping serviceMapping in retrieveServiceMappingsForAPIResult {
            http:Response|http:ClientError k8ServiceMappingDeletionResponse = deleteK8ServiceMapping(serviceMapping.metadata.name, serviceMapping.metadata.namespace);
            if k8ServiceMappingDeletionResponse is http:Response {
                if k8ServiceMappingDeletionResponse.statusCode != http:STATUS_OK {
                    InternalServerErrorError?|error handleK8sTimeoutResult = self.handleK8sTimeout(k8ServiceMappingDeletionResponse);
                    if handleK8sTimeoutResult is InternalServerErrorError {
                        return handleK8sTimeoutResult;
                    } else if handleK8sTimeoutResult is error {
                        return handleK8sTimeoutResult;
                    }

                }
            } else {
                log:printError("Error occured while deleting service mapping");
            }
        }
        return;
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
            InternalServerErrorError internalError = {body: {code: 90900, message: "InternalServerError"}};
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

    isolated function handleK8sTimeout(http:Response response) returns InternalServerErrorError?|error {
        json responsePayLoad = check response.getJsonPayload();
        model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
        model:StatusDetails details = statusResponse.details;
        if details.retryAfterSeconds is int && details.retryAfterSeconds >= 0 {
            // K8s api level ratelimit hit.
            log:printError("Error while deploying API");
            InternalServerErrorError internalEror = {body: {code: 90915, message: "Internal Error while Deploying API"}};
            return internalEror;
        } else {
            return;
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

