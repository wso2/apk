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
import runtime_domain_service.org.wso2.apk.runtime.model as runtimeModels;
import runtime_domain_service.java.util as utilapis;
import ballerina/jwt;
import ballerina/regex;
import runtime_domain_service.org.wso2.apk.runtime as runtimeUtil;
import ballerina/mime;
import ballerina/jballerina.java;
import ballerina/lang.value;
import runtime_domain_service.java.lang;
import runtime_domain_service.org.wso2.apk.runtime.api as runtimeapi;
import ballerina/uuid;
import ballerina/file;
import ballerina/io;
import ballerina/crypto;
import wso2/apk_common_lib as commons;

public class APIClient {

    public isolated function getAPIDefinitionByID(string id, commons:Organization organization) returns http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        model:API|error api = getAPI(id, organization);
        if api is model:API {
            json|error definition = self.getDefinition(api);
            if definition is json {
                http:Response response = new;
                response.setJsonPayload(definition);
                response.statusCode = 200;
                return response;
            } else {
                log:printError("Error while reading definition:", definition);
                InternalServerErrorError internalError = {body: {code: 909000, message: "Internal Error Occured while retrieving definition"}};
                return internalError;
            }
        }
        NotFoundError notfound = {body: {code: 909100, message: id + " not found."}};
        return notfound;
    }

    private isolated function getDefinition(model:API api) returns json|commons:APKError {
        do {
            string? definitionFileRef = api.spec.definitionFileRef;
            if definitionFileRef is string && definitionFileRef.length() > 0 {
                model:ConfigMap? configmapDefinition = check self.getDefinitionConfigmap(definitionFileRef, api.metadata.namespace);
                if configmapDefinition is model:ConfigMap {
                    return check self.getDefinitionFromConfigMap(configmapDefinition);
                }
            }
            // definitionFileRef not specified or empty. definitionfile
            return self.retrieveDefaultDefinition(api);
        } on fail var e {
            string message = "Internal Error occured while retrieving api Definition";
            return error(message, e, message = message, description = message, code = 909000, statusCode = 500);
        }

    }
    private isolated function getDefinitionFromConfigMap(model:ConfigMap configmap) returns json|error? {
        map<string>? data = configmap.data;
        if data is map<string> {
            string[] keys = data.keys();
            if keys.length() == 1 {
                return check value:fromJsonString(data.get(keys[0]).toString());
            }
        }
        return;
    }
    private isolated function getDefinitionConfigmap(string name, string namespace) returns model:ConfigMap|error? {
        http:Response response = check getConfigMapValueFromNameAndNamespace(name, namespace);
        if response.statusCode == 200 {
            json configMapValue = check response.getJsonPayload();
            model:ConfigMap configmapDefinition = check configMapValue.cloneWithType(model:ConfigMap);
            return configmapDefinition;
        }
        return;
    }

    //Get APIs deployed in default namespace by APIId.
    public isolated function getAPIById(string id, commons:Organization organization) returns API|NotFoundError|commons:APKError {
        boolean APIIDAvailable = id.length() > 0 ? true : false;
        if (APIIDAvailable && string:length(id.toString()) > 0)
        {
            lock {
                map<model:API>? apiMap = apilist[organization.uuid];
                if apiMap is map<model:API> {
                    model:API? api = apiMap[id];
                    if api != null {
                        API detailedAPI = check convertK8sAPItoAPI(api, false);
                        return detailedAPI.cloneReadOnly();
                    }
                }
            } on fail var e {
                return error("Error while retrieving API", e, message = "Error while retrieving API", description = "Error while retrieving API", code = 909000, statusCode = 500);
            }
        }
        NotFoundError notfound = {body: {code: 909100, message: id + " not found."}};
        return notfound;
    }

    //Delete APIs deployed in a namespace by APIId.
    public isolated function deleteAPIById(string id, commons:Organization organization) returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|commons:APKError {
        boolean APIIDAvailable = id.length() > 0 ? true : false;
        if (APIIDAvailable && string:length(id.toString()) > 0)
        {
            model:API|error api = getAPI(id, organization);
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
                _ = check self.deleteHttpRoutes(api);
                _ = check self.deleteServiceMappings(api);
                _ = check self.deleteAuthneticationCRs(api);
                _ = check self.deleteScopeCrsForAPI(api);
                _ = check self.deleteBackendPolicies(api);
                _ = check self.deleteBackendServices(api);
                _ = check self.deleteInternalAPI(api.metadata.name, api.metadata.namespace);
            } else {
                NotFoundError apiNotfound = {body: {code: 900910, description: "API with " + id + " not found", message: "API not found"}};
                return apiNotfound;
            }
        }
        return http:OK;
    }
    private isolated function deleteHttpRoutes(model:API api) returns commons:APKError? {
        do {
            model:HttprouteList|http:ClientError httpRouteListResponse = check getHttproutesForAPIS(api.spec.apiDisplayName, api.spec.apiVersion, api.metadata.namespace);
            if httpRouteListResponse is model:HttprouteList {
                foreach model:Httproute item in httpRouteListResponse.items {
                    http:Response|http:ClientError httprouteDeletionResponse = deleteHttpRoute(item.metadata.name, item.metadata.namespace);
                    if httprouteDeletionResponse is http:Response {
                        if httprouteDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check httprouteDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting Httproute", httprouteDeletionResponse);
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting httproutes", e);
            return error("Error occured deleting httproutes", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }
    private isolated function deleteInternalAPI(string k8sAPIName, string apiNameSpace) returns commons:APKError? {
        do {
            model:RuntimeAPI|http:ClientError internalAPI = getInternalAPI(k8sAPIName, apiNameSpace);
            if internalAPI is model:RuntimeAPI {
                http:Response internalAPIDeletionResponse = check deleteInternalAPI(k8sAPIName, apiNameSpace);
                if internalAPIDeletionResponse.statusCode != http:STATUS_OK {
                    json responsePayLoad = check internalAPIDeletionResponse.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            }
        } on fail var e {
            log:printError("Error occured deleting Internal API", e);
            return error("Error occured deleting servicemapping", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }
    private isolated function deleteBackendPolicies(model:API api) returns commons:APKError? {
        do {
            model:BackendPolicyList|http:ClientError backendPolicyListResponse = check getBackendPolicyCRsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, api.metadata.namespace);
            if backendPolicyListResponse is model:BackendPolicyList {
                foreach model:BackendPolicy item in backendPolicyListResponse.items {
                    http:Response|http:ClientError serviceDeletionResponse = deleteBackendPolicyCR(item.metadata.name, item.metadata.namespace);
                    if serviceDeletionResponse is http:Response {
                        if serviceDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check serviceDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting service mapping", serviceDeletionResponse);
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting servicemapping", e);
            return error("Error occured deleting servicemapping", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }
    private isolated function deleteBackendServices(model:API api) returns commons:APKError? {
        do {
            model:ServiceList|http:ClientError serviceListResponse = check getBackendServicesForAPI(api.spec.apiDisplayName, api.spec.apiVersion, api.metadata.namespace);
            if serviceListResponse is model:ServiceList {
                foreach model:Service item in serviceListResponse.items {
                    http:Response|http:ClientError serviceDeletionResponse = deleteService(item.metadata.name, item.metadata.namespace);
                    if serviceDeletionResponse is http:Response {
                        if serviceDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check serviceDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting service mapping");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting servicemapping", e);
            return error("Error occured deleting servicemapping", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }
    private isolated function deleteAuthneticationCRs(model:API api) returns commons:APKError? {
        do {
            model:AuthenticationList|http:ClientError authenticationCrListResponse = check getAuthenticationCrsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, api.metadata.namespace);
            if authenticationCrListResponse is model:AuthenticationList {
                foreach model:Authentication item in authenticationCrListResponse.items {
                    http:Response|http:ClientError k8ServiceMappingDeletionResponse = deleteAuthenticationCR(item.metadata.name, item.metadata.namespace);
                    if k8ServiceMappingDeletionResponse is http:Response {
                        if k8ServiceMappingDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check k8ServiceMappingDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting service mapping");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting servicemapping", e);
            return error("Error occured deleting servicemapping", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }
    private isolated function deleteScopeCrsForAPI(model:API api) returns commons:APKError? {
        do {
            model:ScopeList|http:ClientError scopeCrListResponse = check getScopeCrsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, api.metadata.namespace);
            if scopeCrListResponse is model:ScopeList {
                foreach model:Scope item in scopeCrListResponse.items {
                    http:Response|http:ClientError scopeCrDeletionResponse = deleteScopeCr(item.metadata.name, item.metadata.namespace);
                    if scopeCrDeletionResponse is http:Response {
                        if scopeCrDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check scopeCrDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting scopes");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting scope", e);
            return error("Error occured deleting scopes", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }

    # This returns list of APIS.
    #
    # + query - Parameter Description  
    # + 'limit - Parameter Description  
    # + offset - Parameter Description  
    # + sortBy - Parameter Description  
    # + sortOrder - Parameter Description  
    # + organization - Parameter Description
    # + return - Return list of APIS in namsepace.
    public isolated function getAPIList(string? query, int 'limit, int offset, string sortBy, string sortOrder, commons:Organization organization) returns APIList|BadRequestError|commons:APKError {
        API[] apilist = [];
        foreach model:API api in getAPIs(organization) {
            API convertedModel = check convertK8sAPItoAPI(api, true);
            apilist.push(convertedModel);
        } on fail var e {
            return error("Error occured while getting API list", e, message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
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
    public isolated function createAPI(API api, string? definition, commons:Organization organization) returns commons:APKError|CreatedAPI|BadRequestError {
        do {
            if (self.validateName(api.name, organization)) {
                BadRequestError badRequest = {body: {code: 90911, message: "API Name - " + api.name + " already exist.", description: "API Name - " + api.name + " already exist."}};
                return badRequest;
            }
            lock {
                if (ALLOWED_API_TYPES.indexOf(api.'type) is ()) {
                    BadRequestError badRequest = {body: {code: 900912, message: "unsupported API Type."}};
                    return badRequest.clone();
                }
            }
            if self.validateContextAndVersion(api.context, api.'version, organization) {
                BadRequestError badRequest = {body: {code: 90911, message: "API Context - " + api.context + " already exist.", description: "API Context " + api.context + " already exist."}};
                return badRequest;
            }

            self.setDefaultOperationsIfNotExist(api);
            string uniqueId = getUniqueIdForAPI(api.name, api.'version, organization);
            model:APIArtifact apiArtifact = {uniqueId: uniqueId};
            APIOperations[]? operations = api.operations;
            if operations is APIOperations[] {
                if operations.length() == 0 {
                    BadRequestError badRequestError = {body: {code: 90912, message: "Atleast one operation need to specified"}};
                    return badRequestError;
                }
                // Validating operation policies.
                BadRequestError|() badRequestError = self.validateOperationPolicies(api.apiPolicies, operations, organization);
                if (badRequestError is BadRequestError) {
                    return badRequestError;
                }
            } else {
                BadRequestError badRequestError = {body: {code: 90912, message: "Atleast one operation need to specified"}};
                return badRequestError;
            }
            record {}? endpointConfig = api.endpointConfig;
            map<model:Endpoint|()> createdEndpoints = {};
            if endpointConfig is record {} {
                createdEndpoints = check self.createAndAddBackendServics(apiArtifact, api, endpointConfig, (), (), organization);
            }
            _ = check self.setHttpRoute(apiArtifact, api, createdEndpoints.hasKey(PRODUCTION_TYPE) ? createdEndpoints.get(PRODUCTION_TYPE) : (), uniqueId, PRODUCTION_TYPE, organization);
            _ = check self.setHttpRoute(apiArtifact, api, createdEndpoints.hasKey(SANDBOX_TYPE) ? createdEndpoints.get(SANDBOX_TYPE) : (), uniqueId, SANDBOX_TYPE, organization);
            json generatedSwagger = check self.retrieveGeneratedSwaggerDefinition(api, definition);
            self.retrieveGeneratedConfigmapForDefinition(apiArtifact, api, generatedSwagger, uniqueId);
            self.generateAndSetAPICRArtifact(apiArtifact, api, organization);
            self.generateAndSetRuntimeAPIArtifact(apiArtifact, api, (), organization);
            model:API deployAPIToK8sResult = check self.deployAPIToK8s(apiArtifact);
            CreatedAPI createdAPI = {body: check convertK8sAPItoAPI(deployAPIToK8sResult, true)};
            return createdAPI;
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            log:printError("Internal Error occured", e);
            return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = 500);
        }
    }

    isolated function validateOperationPolicies(APIOperationPolicies? apiPolicies, APIOperations[] operations, commons:Organization organization) returns BadRequestError|() {
        foreach APIOperations operation in operations {
            APIOperationPolicies? operationPolicies = operation.operationPolicies;
            if (!self.isPolicyEmpty(operationPolicies)) {
                if (self.isPolicyEmpty(apiPolicies)) {
                    // Validating resource level operation policy data
                    BadRequestError|() badRequestError = self.validateOperationPolicyData(operationPolicies, organization);
                    if (badRequestError is BadRequestError) {
                        return badRequestError;
                    }
                } else {
                     // Presence of both resource level and API level operation policies.
                    BadRequestError badRequestError = {body: {code: 90917, message: "Presence of both resource level and API level operation policies is not allowed"}};
                    return badRequestError;
                }
            }
        }
        if (!self.isPolicyEmpty(apiPolicies)) {
            // Validating API level operation policy data
            return self.validateOperationPolicyData(apiPolicies, organization);
        }
        return ();
    }

    isolated function isPolicyEmpty(APIOperationPolicies? policies) returns boolean {
        return policies == () || policies.length() == 0;
    }

    isolated function validateOperationPolicyData(APIOperationPolicies? operationPolicies, commons:Organization organization) returns BadRequestError|() {
        if operationPolicies is APIOperationPolicies {
            // Validating request operation policy data.
            BadRequestError|() badRequestError = self.validatePolicyDetails(operationPolicies.request, organization);
            if (badRequestError == ()) {
                // Validating response operation policy data.
                return self.validatePolicyDetails(operationPolicies.response, organization);
            } else {
                return badRequestError;
            }
        }
        return ();
    }

    isolated function validatePolicyDetails(OperationPolicy[]? policyData, commons:Organization organization) returns BadRequestError|() {
        if (policyData is OperationPolicy[]) {
            foreach OperationPolicy policy in policyData {
                string policyName = policy.policyName;
                OperationPolicyParameters[]? policyParameters = policy.parameters;
                if (policyParameters is OperationPolicyParameters[]) {
                    string[] allowedPolicyAttributes = [];
                    string[] receivedPolicyParameters = [];
                    any|error mediationPolicyList = self.getMediationPolicyList(SEARCH_CRITERIA_NAME + ":" + policyName, 1, 0,
                        SORT_BY_POLICY_NAME, SORT_ORDER_ASC, organization);
                    if (mediationPolicyList is MediationPolicyList && mediationPolicyList.count > 0) {
                        MediationPolicy[]? listing = mediationPolicyList.list;
                        if (listing is MediationPolicy[]) {
                            MediationPolicySpecAttribute[]? parameters = listing[0].policyAttributes;
                            if (parameters is MediationPolicySpecAttribute[]) {
                                foreach MediationPolicySpecAttribute params in parameters {
                                    allowedPolicyAttributes.push(<string>params.name);
                                }
                            }
                            foreach OperationPolicyParameters policyParam in policyParameters {
                                string[] keys = policyParam.keys();
                                foreach string key in keys {
                                    receivedPolicyParameters.push(key);
                                }
                            }
                            if (allowedPolicyAttributes != receivedPolicyParameters) {
                                // Allowed policy attributes does not match with the parameters provided
                                BadRequestError badRequestError = {body: {code: 90916, message: "Invalid parameters provided for policy " + policyName}};
                                return badRequestError;
                            }
                        }
                    } else {
                        // Invalid operation policy name.
                        BadRequestError badRequestError = {body: {code: 90915, message: "Invalid operation policy name"}};
                        return badRequestError;
                    }
                }
            }
        }
        return ();
    }


    private isolated function generateAndSetRuntimeAPIArtifact(model:APIArtifact apiArtifact, API api, Service? serviceEntry, commons:Organization organization) {

        apiArtifact.runtimeAPI = self.generateRuntimeAPIArtifact(api, serviceEntry, organization);
    }
    public isolated function generateRuntimeAPIArtifact(API api, Service? serviceEntry, commons:Organization organization) returns model:RuntimeAPI {
        model:RuntimeAPI runtimeAPI = {
            metadata: {
                name: getUniqueIdForAPI(api.name, api.'version, organization),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                labels: self.getLabels(api)
            },
            spec: {
                name: api.name,
                context: api.context,
                'version: api.'version,
                'type: api.'type
            }
        };
        APIOperations[]? operations = api.operations;
        if operations is APIOperations[] {
            model:Operations[] runtimeAPIOperations = [];
            foreach APIOperations operation in operations {
                model:Operations runtimeAPIOperation = {
                    target: <string>operation.target,
                    verb: <string>operation.verb,
                    authTypeEnabled: operation.authTypeEnabled ?: true,
                    scopes: operation.scopes ?: []
                };
                APIOperationPolicies? operationPoliciesToUse = ();
                if (api.apiPolicies is APIOperationPolicies) {
                    operationPoliciesToUse = api.apiPolicies;
                } else {
                    operationPoliciesToUse = operation.operationPolicies;
                }
                model:OperationPolicy[] runtimeAPIRequestPolicies = [];
                model:OperationPolicy[] runtimeAPIResponsePolicies = [];

                if operationPoliciesToUse is APIOperationPolicies {
                    OperationPolicy[]? request = operationPoliciesToUse.request;
                    if request is OperationPolicy[] {
                        foreach OperationPolicy policy in request {
                            model:OperationPolicy|error runtimeAPIRequestPolicy = policy.cloneWithType(model:OperationPolicy);
                            if (runtimeAPIRequestPolicy is model:OperationPolicy) {
                                runtimeAPIRequestPolicies.push(runtimeAPIRequestPolicy);
                            }
                        }
                    }
                    OperationPolicy[]? response = operationPoliciesToUse.response;
                    if response is OperationPolicy[] {
                        foreach OperationPolicy policy in response {
                            model:OperationPolicy|error runtimeAPIResponsePolicy = policy.cloneWithType(model:OperationPolicy);
                            if (runtimeAPIResponsePolicy is model:OperationPolicy) {
                                runtimeAPIResponsePolicies.push(runtimeAPIResponsePolicy);
                            }
                        }
                    }
                }
                runtimeAPIOperation.operationPolicies = {
                    request: runtimeAPIRequestPolicies,
                    response: runtimeAPIResponsePolicies
                };
                record {|anydata...;|}? endpointConfig = api.endpointConfig;
                if endpointConfig is record {} {
                    runtimeAPI.spec.endpointConfig = endpointConfig;
                }
                runtimeAPIOperations.push(runtimeAPIOperation);
            }
            runtimeAPI.spec.operations = runtimeAPIOperations;
        }
        record {|anydata...;|}? endpointConfig = api.endpointConfig;
        if endpointConfig is record {} {
            runtimeAPI.spec.endpointConfig = endpointConfig;
        }
        if serviceEntry is Service {
            runtimeAPI.spec.serviceInfo = {
                name: serviceEntry.name,
                namespace: serviceEntry.namespace
            };
        }
        return runtimeAPI;
    }
    private isolated function createBackendPolicy(API api, APIOperations? apiOperation, string endpointType, commons:Organization organization, map<anydata>? endpointConfig, string url, model:Service backendService) returns model:BackendPolicy {
        string nameSpace = getNameSpace(runtimeConfiguration.apiCreationNamespace);
        model:BackendPolicy backendPolicy = {
            metadata: {
                name: getBackendPolicyUid(api, endpointType, organization),
                namespace: nameSpace,
                labels: self.getLabels(api)
            },
            spec: {
                default: {protocol: url.startsWith(PROTOCOL_HTTPS) ? PROTOCOL_HTTPS : PROTOCOL_HTTP},
                targetRef: {
                    group: "",
                    kind: "Service",
                    name: backendService.metadata.name,
                    namespace: backendService.metadata.namespace
                }
            }
        };
        return backendPolicy;
    }
    private isolated function createAndAddBackendServics(model:APIArtifact apiArtifact, API api, record {} endpointConfig, APIOperations? apiOperation, string? endpointType, commons:Organization organization) returns map<model:Endpoint>|commons:APKError|error {
        map<model:Endpoint> endpointIdMap = {};
        anydata|error sandboxEndpointConfig = trap endpointConfig.get("sandbox_endpoints");
        anydata|error productionEndpointConfig = trap endpointConfig.get("production_endpoints");
        if endpointType == () || (endpointType == SANDBOX_TYPE) {
            if sandboxEndpointConfig is map<anydata> {
                if sandboxEndpointConfig.hasKey("url") {
                    anydata url = sandboxEndpointConfig.get("url");
                    model:Service backendService = self.createBackendService(api, apiOperation, SANDBOX_TYPE, organization, <string>url);
                    model:BackendPolicy backendPolicy = self.createBackendPolicy(api, apiOperation, SANDBOX_TYPE, organization, sandboxEndpointConfig, <string>url, backendService);
                    if apiOperation == () {
                        apiArtifact.sandboxEndpointAvailable = true;
                        apiArtifact.sandboxUrl = <string?>url;
                    }
                    apiArtifact.backendServices[backendService.metadata.name] = (backendService);
                    apiArtifact.backendPolicies[backendPolicy.metadata.name] = (backendPolicy);
                    endpointIdMap[SANDBOX_TYPE] = {
                        port: check self.getPort(<string>url),
                        namespace: backendService.metadata.namespace,
                        name: backendService.metadata.name,
                        serviceEntry: false,
                        url: <string?>url
                    };
                } else {
                    commons:APKError e = error("Sandbox Endpoint Not specified", message = "Endpoint Not specified", description = "Sandbox Endpoint Not specified", code = 90911, statusCode = 400);
                    return e;
                }
            }
        }
        if endpointType == () || (endpointType == PRODUCTION_TYPE) {
            if productionEndpointConfig is map<anydata> {
                if productionEndpointConfig.hasKey("url") {
                    anydata url = productionEndpointConfig.get("url");
                    model:Service backendService = self.createBackendService(api, apiOperation, PRODUCTION_TYPE, organization, <string>url);
                    model:BackendPolicy backendPolicy = self.createBackendPolicy(api, apiOperation, SANDBOX_TYPE, organization, productionEndpointConfig, <string>url, backendService);
                    if apiOperation == () {
                        apiArtifact.productionEndpointAvailable = true;
                        apiArtifact.productionUrl = <string?>url;
                    }
                    apiArtifact.backendServices[backendService.metadata.name] = backendService;
                    apiArtifact.backendPolicies[backendPolicy.metadata.name] = (backendPolicy);
                    endpointIdMap[PRODUCTION_TYPE] = {
                        port: check self.getPort(<string>url),
                        namespace: backendService.metadata.namespace,
                        name: backendService.metadata.name,
                        serviceEntry: false,
                        url: <string?>url
                    };
                } else {
                    commons:APKError e = error("Production Endpoint Not specified", message = "Endpoint Not specified", description = "Production Endpoint Not specified", code = 90911, statusCode = 400);
                    return e;
                }
            }
        }
        return endpointIdMap;
    }
    isolated function getLabels(API api) returns map<string> {
        map<string> labels = {"api-name": api.name, "api-version": api.'version};
        return labels;
    }
    isolated function validateContextAndVersion(string context, string 'version, commons:Organization organization) returns boolean {
        foreach model:API k8sAPI in getAPIs(organization) {
            if k8sAPI.spec.context == self.returnFullContext(context, 'version) &&
            k8sAPI.spec.organization == organization.uuid {
                return true;
            }
        }
        return false;
    }

    isolated function validateContext(string context, commons:Organization organization) returns boolean {

        foreach model:API k8sAPI in getAPIs(organization) {
            if (k8sAPI.spec.context.startsWith(context) &&
            k8sAPI.spec.organization == organization.uuid) {
                return true;
            }
        }
        return false;
    }

    isolated function returnFullContext(string context, string 'version) returns string {
        string fullContext = context;
        if (!string:endsWith(context, 'version)) {
            fullContext = string:'join("/", context, 'version);
        }
        return fullContext;
    }

    isolated function validateName(string name, commons:Organization organization) returns boolean {
        foreach model:API k8sAPI in getAPIs(organization) {
            if k8sAPI.spec.apiDisplayName == name && k8sAPI.spec.organization == organization.uuid {
                return true;
            }
        }
        return false;
    }

    isolated function createAPIFromService(string serviceKey, API api, commons:Organization organization) returns CreatedAPI|BadRequestError|InternalServerErrorError|commons:APKError {
        if (self.validateName(api.name, organization)) {
            BadRequestError badRequest = {body: {code: 90911, message: "API Name - " + api.name + " already exist.", description: "API Name - " + api.name + " already exist."}};
            return badRequest;
        }
        lock {
            if (ALLOWED_API_TYPES.indexOf(api.'type) is ()) {
                BadRequestError badRequest = {body: {code: 900912, message: "unsupported API Type."}};
                return badRequest.clone();
            }
        }
        if self.validateContextAndVersion(api.context, api.'version, organization) {
            BadRequestError badRequest = {body: {code: 90911, message: "API Context - " + api.context + " already exist.", description: "API Context " + api.context + " already exist."}};
            return badRequest;
        }
        self.setDefaultOperationsIfNotExist(api);
        APIOperations[]? operations = api.operations;
        if operations is APIOperations[] {
            // Validating operation policies.
            BadRequestError|() badRequestError = self.validateOperationPolicies(api.apiPolicies, operations, organization);
            if (badRequestError is BadRequestError) {
                return badRequestError;
            }
        }
        api.context = self.returnFullContext(api.context, api.'version);
        Service|error serviceRetrieved = getServiceById(serviceKey);
        string uniqueId = getUniqueIdForAPI(api.name, api.'version, organization);
        if serviceRetrieved is Service {
            model:APIArtifact apiArtifact = {uniqueId: uniqueId};
            model:Endpoint endpoint = {
                port: self.retrievePort(serviceRetrieved),
                namespace: serviceRetrieved.namespace,
                name: serviceRetrieved.name,
                serviceEntry: true
            };
            check self.setHttpRoute(apiArtifact, api, endpoint, uniqueId, PRODUCTION_TYPE, organization);
            json generatedSwaggerDefinition = check self.retrieveGeneratedSwaggerDefinition(api, ());
            self.retrieveGeneratedConfigmapForDefinition(apiArtifact, api, generatedSwaggerDefinition, uniqueId);
            self.generateAndSetAPICRArtifact(apiArtifact, api, organization);
            self.generateAndSetK8sServiceMapping(apiArtifact, api, serviceRetrieved, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            self.generateAndSetRuntimeAPIArtifact(apiArtifact, api, serviceRetrieved, organization);
            model:API deployAPIToK8sResult = check self.deployAPIToK8s(apiArtifact);
            CreatedAPI createdAPI = {body: check convertK8sAPItoAPI(deployAPIToK8sResult, false)};
            return createdAPI;
        } else {
            BadRequestError badRequest = {body: {code: 90913, message: "Service from " + serviceKey + " not found."}};
            return badRequest;
        }
    }
    private isolated function deployAPIToK8s(model:APIArtifact apiArtifact) returns commons:APKError|model:API {
        do {
            model:ConfigMap? definition = apiArtifact.definition;
            if definition is model:ConfigMap {
                _ = check self.deployConfigMap(definition);
            }
            model:API? api = apiArtifact.api;
            if api is model:API {
                check self.deleteHttpRoutes(api);
                check self.deleteServiceMappings(api);
                check self.deleteAuthneticationCRs(api);
                _ = check self.deleteScopeCrsForAPI(api);
                check self.deleteBackendPolicies(api);
                check self.deleteBackendServices(api);
                check self.deleteInternalAPI(api.metadata.name, api.metadata.namespace);
            }
            check self.deployScopeCrs(apiArtifact);
            check self.deployBackendServices(apiArtifact);
            check self.deployBackendPolicies(apiArtifact);
            check self.deployAuthneticationCRs(apiArtifact);
            check self.deployHttpRoutes(apiArtifact.productionRoute);
            check self.deployHttpRoutes(apiArtifact.sandboxRoute);
            check self.deployServiceMappings(apiArtifact);
            check self.deployRuntimeAPI(apiArtifact);
            return check self.deployK8sAPICr(apiArtifact);
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            log:printError("Internal Error occured while deploying API", e);
            commons:APKError internalError = error("Internal Error occured while deploying API", code = 909000, statusCode = 500, description = "Internal Error occured while deploying API", message = "Internal Error occured while deploying API");
            return internalError;
        }
    }
    private isolated function deployScopeCrs(model:APIArtifact apiArtifact) returns error? {
        foreach model:Scope scope in apiArtifact.scopes {
            http:Response deployScopeResult = check deployScopeCR(scope, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployScopeResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Scope Successfully" + scope.toString());
            } else {
                json responsePayLoad = check deployScopeResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }
    private isolated function deployK8sAPICr(model:APIArtifact apiArtifact) returns model:API|commons:APKError|error {
        model:API? k8sAPI = apiArtifact.api;
        if k8sAPI is model:API {
            model:API? k8sAPIByNameAndNamespace = check getK8sAPIByNameAndNamespace(k8sAPI.metadata.name, k8sAPI.metadata.namespace);
            if k8sAPIByNameAndNamespace is model:API {
                k8sAPI.metadata.resourceVersion = k8sAPIByNameAndNamespace.metadata.resourceVersion;
                http:Response deployAPICRResult = check updateAPICR(k8sAPI, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployAPICRResult.statusCode == http:STATUS_OK {
                    json responsePayLoad = check deployAPICRResult.getJsonPayload();
                    log:printDebug("Updated K8sAPI Successfully" + responsePayLoad.toJsonString());
                    return check responsePayLoad.cloneWithType(model:API);
                } else {
                    json responsePayLoad = check deployAPICRResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    model:StatusDetails? details = statusResponse.details;
                    if details is model:StatusDetails {
                        model:StatusCause[] 'causes = details.'causes;
                        foreach model:StatusCause 'cause in 'causes {
                            if 'cause.'field == "spec.context" {
                                commons:APKError badeRequestError = error("Invalid API Context", code = 90911, description = "API Context " + k8sAPI.spec.context + " Invalid", message = "Invalid API context", statusCode = 400);
                                return badeRequestError;
                            } else if 'cause.'field == "spec.apiDisplayName" {
                                commons:APKError badeRequestError = error("Invalid API Name", code = 90911, description = "API Name " + k8sAPI.spec.apiDisplayName + " Invalid", message = "Invalid API Name", statusCode = 400);
                                return badeRequestError;
                            }
                        }
                        commons:APKError badeRequestError = error("Invalid API Request", code = 90911, description = "Invalid API Request", message = "Invalid API Request", statusCode = 400);
                        return badeRequestError;
                    }
                    return self.handleK8sTimeout(statusResponse);
                }
            } else {
                http:Response deployAPICRResult = check deployAPICR(k8sAPI, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployAPICRResult.statusCode == http:STATUS_CREATED {
                    json responsePayLoad = check deployAPICRResult.getJsonPayload();
                    log:printDebug("Deployed K8sAPI Successfully" + responsePayLoad.toJsonString());
                    return check responsePayLoad.cloneWithType(model:API);
                } else {
                    json responsePayLoad = check deployAPICRResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    model:StatusDetails? details = statusResponse.details;
                    if details is model:StatusDetails {
                        model:StatusCause[] 'causes = details.'causes;
                        foreach model:StatusCause 'cause in 'causes {
                            if 'cause.'field == "spec.context" {
                                commons:APKError badRequestError = error("Invalid API Context", code = 90911, description = "API Context " + k8sAPI.spec.context + " Invalid", message = "Invalid API context", statusCode = 400);
                                return badRequestError;
                            } else if 'cause.'field == "spec.apiDisplayName" {
                                commons:APKError badRequestError = error("Invalid API Name", code = 90911, description = "API Name " + k8sAPI.spec.apiDisplayName + " Invalid", message = "Invalid API Name", statusCode = 400);
                                return badRequestError;
                            }
                        }
                        commons:APKError badRequestError = error("Invalid API Request", code = 90911, description = "Invalid API Request", message = "Invalid API Request", statusCode = 400);
                        return badRequestError;

                    }
                    return self.handleK8sTimeout(statusResponse);
                }
            }
        } else {
            commons:APKError badRequestError = error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = 500);
            return badRequestError;
        }
    }
    private isolated function deployRuntimeAPI(model:APIArtifact apiArtifact) returns error? {
        model:RuntimeAPI? runtimeapi = apiArtifact.runtimeAPI;
        if runtimeapi is model:RuntimeAPI {
            http:Response deployRuntimeAPICRResult = check createInternalAPI(runtimeapi, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployRuntimeAPICRResult.statusCode == http:STATUS_CREATED {
                json responsePayLoad = check deployRuntimeAPICRResult.getJsonPayload();
                log:printDebug("Deployed RuntimeAPI Successfully" + responsePayLoad.toJsonString());
            } else {
                json responsePayLoad = check deployRuntimeAPICRResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                model:StatusDetails? details = statusResponse.details;
                if details is model:StatusDetails {
                    return error("Invalid API Request", code = 90911, description = "Invalid API Request", message = "Invalid API Request", statusCode = 400);
                }
                return self.handleK8sTimeout(statusResponse);
            }
        } else {
            return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = 500);
        }
    }
    private isolated function deployHttpRoutes(model:Httproute[] httproutes) returns error? {
        foreach model:Httproute httpRoute in httproutes {
            if httpRoute.spec.rules.length() > 0 {
                http:Response deployHttpRouteResult = check deployHttpRoute(httpRoute, getNameSpace(runtimeConfiguration.apiCreationNamespace));
                if deployHttpRouteResult.statusCode == http:STATUS_CREATED {
                    log:printDebug("Deployed HttpRoute Successfully" + httpRoute.toString());
                } else {
                    json responsePayLoad = check deployHttpRouteResult.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            }
        }
    }
    private isolated function deployServiceMappings(model:APIArtifact apiArtifact) returns error? {
        foreach model:K8sServiceMapping k8sServiceMapping in apiArtifact.serviceMapping {
            http:Response deployServiceMappingCRResult = check deployServiceMappingCR(k8sServiceMapping, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployServiceMappingCRResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed K8sAPI Successfully" + k8sServiceMapping.toString());
            } else {
                json responsePayLoad = check deployServiceMappingCRResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }
    private isolated function deployAuthneticationCRs(model:APIArtifact apiArtifact) returns error? {
        string[] keys = apiArtifact.authenticationMap.keys();
        foreach string authenticationCrName in keys {
            model:Authentication authenticationCr = apiArtifact.authenticationMap.get(authenticationCrName);
            http:Response authenticationCrDeployResponse = check deployAuthenticationCR(authenticationCr, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if authenticationCrDeployResponse.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed HttpRoute Successfully" + authenticationCr.toString());
            } else {
                json responsePayLoad = check authenticationCrDeployResponse.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }
    private isolated function deployBackendPolicies(model:APIArtifact apiArtifact) returns error? {
        foreach model:BackendPolicy backendPolicy in apiArtifact.backendPolicies {
            http:Response deployBackendPolicyResult = check deployBackendPolicyCR(backendPolicy, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployBackendPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed BackendPolicy Successfully" + backendPolicy.toString());
            } else {
                json responsePayLoad = check deployBackendPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deployBackendServices(model:APIArtifact apiArtifact) returns error? {
        foreach model:Service backendService in apiArtifact.backendServices {
            http:Response deployServiceResult = check deployService(backendService, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployServiceResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed HttpRoute Successfully" + backendService.toString());
            } else {
                json responsePayLoad = check deployServiceResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }
    private isolated function deployConfigMap(model:ConfigMap definition) returns commons:APKError|error? {
        http:Response configMapRetrieved = check getConfigMapValueFromNameAndNamespace(definition.metadata.name, definition.metadata.namespace);
        if configMapRetrieved.statusCode == 404 {
            http:Response deployConfigMapResult = check deployConfigMap(definition, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployConfigMapResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Configmap Successfully" + definition.toString());
            } else {
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        } else if configMapRetrieved.statusCode == 200 {
            http:Response deployConfigMapResult = check updateConfigMap(definition, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            if deployConfigMapResult.statusCode == http:STATUS_OK {
                log:printDebug("updated Configmap Successfully" + definition.toString());
            } else {
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        } else {
            json responsePayLoad = check configMapRetrieved.getJsonPayload();
            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
            check self.handleK8sTimeout(statusResponse);
        }
    }
    private isolated function retrieveGeneratedConfigmapForDefinition(model:APIArtifact apiArtifact, API api, json generatedSwaggerDefinition, string uniqueId) {
        map<string> configMapData = {};
        if api.'type == API_TYPE_REST {
            configMapData["openapi.json"] = generatedSwaggerDefinition.toJsonString();
        }
        model:ConfigMap configMap = {
            metadata: {
                name: self.retrieveDefinitionName(uniqueId),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: (),
                labels: self.getLabels(api)

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
            if api.'type == API_TYPE_REST {
                foreach string httpverb in HTTP_DEFAULT_METHODS {
                    APIOperations apiOperation = {target: "/*", verb: httpverb.toUpperAscii()};
                    apiOperations.push(apiOperation);
                }
                api.operations = apiOperations;
            }
        }
    }

    private isolated function generateAndSetAPICRArtifact(model:APIArtifact apiArtifact, API api, commons:Organization organization) {
        model:API k8sAPI = {
            metadata: {
                name: apiArtifact.uniqueId,
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                labels: self.getLabels(api)
            },
            spec: {
                apiDisplayName: api.name,
                apiType: api.'type,
                apiVersion: api.'version,
                context: self.returnFullContext(api.context, api.'version),
                organization: organization.uuid
            }
        };
        model:ConfigMap? definition = apiArtifact?.definition;
        if definition is model:ConfigMap {
            k8sAPI.spec.definitionFileRef = definition.metadata.name;
        }
        string[] productionHttpRoutes = [];
        foreach model:Httproute httpRoute in apiArtifact.productionRoute {
            if httpRoute.spec.rules.length() > 0 {
                productionHttpRoutes.push(httpRoute.metadata.name);
            }
        }
        string[] sandBoxHttpRoutes = [];
        foreach model:Httproute httpRoute in apiArtifact.sandboxRoute {
            if httpRoute.spec.rules.length() > 0 {
                sandBoxHttpRoutes.push(httpRoute.metadata.name);
            }
        }
        if productionHttpRoutes.length() > 0 {
            k8sAPI.spec.prodHTTPRouteRefs = productionHttpRoutes;
        }
        if sandBoxHttpRoutes.length() > 0 {
            k8sAPI.spec.sandHTTPRouteRefs = sandBoxHttpRoutes;
        }
        apiArtifact.api = k8sAPI;
    }

    isolated function retrieveDefinitionName(string uniqueId) returns string {
        return uniqueId + "-definition";
    }

    private isolated function retrieveDisableAuthenticationRefName(API api, string 'type, commons:Organization organization) returns string {
        return getUniqueIdForAPI(api.name, api.'version, organization) + "-" + 'type + "-authentication";
    }

    private isolated function setHttpRoute(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, string uniqueId, string endpointType, commons:Organization organization) returns commons:APKError? {
        APIOperations[] apiOperations = api.operations ?: [];
        APIOperations[][] operationsArray = [];
        int row = 0;
        int column = 0;
        foreach APIOperations item in apiOperations {
            if column > 7 {
                row = row + 1;
                column = 0;
            }
            operationsArray[row][column] = item;
            column = column + 1;
        }
        foreach APIOperations[] item in operationsArray {
            API clonedAPI = api.clone();
            clonedAPI.operations = item.clone();
            _ = check self.putHttpRouteForPartition(apiArtifact, clonedAPI, endpoint, uniqueId, endpointType, organization);
        }
    }
    private isolated function putHttpRouteForPartition(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, string uniqueId, string endpointType, commons:Organization organization) returns commons:APKError? {
        model:Httproute httpRoute = {
            metadata:
                {
                name: retrieveHttpRouteRefName(api, endpointType, organization),
                namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace),
                uid: (),
                creationTimestamp: (),
                labels: self.getLabels(api)
            },
            spec: {
                parentRefs: self.generateAndRetrieveParentRefs(api, uniqueId),
                rules: check self.generateHttpRouteRules(apiArtifact, api, endpoint, endpointType, organization),
                hostnames: self.getHostNames(api, uniqueId, endpointType, organization)
            }
        };
        if endpointType == PRODUCTION_TYPE {
            apiArtifact.productionRoute.push(httpRoute);
        } else {
            apiArtifact.sandboxRoute.push(httpRoute);
        }
        return;
    }

    private isolated function getHostNames(API api, string uniqueId, string endpointType, commons:Organization organization) returns string[] {
        //todo: need to implement vhost feature
        Vhost[] vhosts = runtimeConfiguration.vhost;
        string[] hosts = [];
        foreach Vhost vhost in vhosts {
            if vhost.'type == endpointType {
                foreach string host in vhost.hosts {
                    hosts.push(string:concat(organization.uuid, ".", host));
                }
            }
        }
        return hosts;
    }

    private isolated function generateAndRetrieveParentRefs(API api, string uniqueId) returns model:ParentReference[] {
        model:ParentReference[] parentRefs = [];
        model:ParentReference parentRef = {group: "gateway.networking.k8s.io", kind: "Gateway", name: "Default"};
        parentRefs.push(parentRef);
        return parentRefs;
    }

    private isolated function generateHttpRouteRules(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, string endpointType, commons:Organization organization) returns model:HTTPRouteRule[]|commons:APKError {
        model:HTTPRouteRule[] httpRouteRules = [];
        APIOperations[]? operations = api.operations;
        if operations is APIOperations[] {
            foreach APIOperations operation in operations {
                model:HTTPRouteRule|() httpRouteRule = check self.generateHttpRouteRule(apiArtifact, api, endpoint, operation, endpointType, organization);
                if httpRouteRule is model:HTTPRouteRule {
                    model:HTTPRouteFilter[]? filters = httpRouteRule.filters;
                    if filters is () {
                        filters = [];
                        httpRouteRule.filters = filters;
                    }
                    if !(operation.authTypeEnabled ?: true) {
                        string disableAuthenticationRefName = self.retrieveDisableAuthenticationRefName(api, endpointType, organization);
                        if !apiArtifact.authenticationMap.hasKey(disableAuthenticationRefName) {
                            model:Authentication generateDisableAuthenticationCR = self.generateDisableAuthenticationCR(apiArtifact, api, endpointType, organization);
                            apiArtifact.authenticationMap[disableAuthenticationRefName] = generateDisableAuthenticationCR;
                        }
                        model:HTTPRouteFilter disableAuthenticationFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: "Authentication", name: disableAuthenticationRefName}};
                        (<model:HTTPRouteFilter[]>filters).push(disableAuthenticationFilter);
                    }
                    string[]? scopes = operation.scopes;
                    if scopes is string[] {
                        foreach string scope in scopes {
                            model:Scope scopeCr;
                            if apiArtifact.scopes.hasKey(scope) {
                                scopeCr = apiArtifact.scopes.get(scope);
                            } else {
                                scopeCr = self.generateScopeCR(apiArtifact, api, organization, scope);
                            }
                            model:HTTPRouteFilter scopeFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: scopeCr.kind, name: scopeCr.metadata.name}};
                            (<model:HTTPRouteFilter[]>filters).push(scopeFilter);
                        }
                    }
                    httpRouteRules.push(httpRouteRule);
                }
            }
        }
        return httpRouteRules;
    }
    private isolated function generateScopeCR(model:APIArtifact apiArtifact, API api, commons:Organization organization, string scope) returns model:Scope {
        string scopeName = uuid:createType1AsString();
        model:Scope scopeCr = {
            metadata: {name: scopeName, namespace: getNameSpace(runtimeConfiguration.apiCreationNamespace), labels: self.getLabels(api)},
            spec: {
                names: [scope]
            }
        };
        apiArtifact.scopes[scope] = scopeCr;
        return scopeCr;
    }
    private isolated function generateDisableAuthenticationCR(model:APIArtifact apiArtifact, API api, string endpointType, commons:Organization organization) returns model:Authentication {
        string retrieveDisableAuthenticationRefName = self.retrieveDisableAuthenticationRefName(api, endpointType, organization);
        string nameSpace = getNameSpace(runtimeConfiguration.apiCreationNamespace);
        model:Authentication authentication = {
            metadata: {name: retrieveDisableAuthenticationRefName, namespace: nameSpace, labels: self.getLabels(api)},
            spec: {
                targetRef: {
                    group: "",
                    kind: "Resource",
                    name: retrieveHttpRouteRefName(api, endpointType, organization),
                    namespace: nameSpace
                },
                override: {
                    ext: {disabled: true},
                    'type: "ext"
                }
            }
        };
        return authentication;
    }

    private isolated function generateHttpRouteRule(model:APIArtifact apiArtifact, API api, model:Endpoint? endpoint, APIOperations operation, string endpointType, commons:Organization organization) returns model:HTTPRouteRule|()|commons:APKError {
        do {
            record {}? endpointConfig = operation.endpointConfig;
            model:Endpoint? endpointToUse = ();
            if endpointConfig is record {} {
                // endpointConfig presense at Operation Level.
                map<model:Endpoint> operationalLevelBackend = check self.createAndAddBackendServics(apiArtifact, api, endpointConfig, operation, endpointType, organization);
                if operationalLevelBackend.hasKey(endpointType) {
                    endpointToUse = operationalLevelBackend.get(endpointType);
                }
            } else {
                if endpoint is model:Endpoint {
                    endpointToUse = endpoint;
                }
            }
            if endpointToUse != () {
                model:HTTPRouteRule httpRouteRule = {matches: self.retrieveMatches(api, operation, organization), backendRefs: self.retrieveGeneratedBackend(api, endpointToUse, endpointType), filters: self.generateFilters(apiArtifact, api, endpointToUse, operation, endpointType, organization)};
                return httpRouteRule;
            } else {
                return ();
            }
        } on fail var e {
            log:printError("Internal Error occured", e);
            return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = 500);
        }
    }

    private isolated function generateFilters(model:APIArtifact apiArtifact, API api, model:Endpoint endpoint, APIOperations operation, string endpointType, commons:Organization organization) returns model:HTTPRouteFilter[] {
        model:HTTPRouteFilter[] routeFilters = [];
        model:HTTPRouteFilter replacePathFilter = {'type: "URLRewrite", urlRewrite: {path: {'type: "ReplaceFullPath", replaceFullPath: self.generatePrefixMatch(api, endpoint, operation, endpointType)}}};
        routeFilters.push(replacePathFilter);
        APIOperationPolicies? operationPoliciesToUse = ();
        if (api.apiPolicies is APIOperationPolicies) {
            operationPoliciesToUse = api.apiPolicies;
        } else {
            operationPoliciesToUse = operation.operationPolicies;
        }
        if operationPoliciesToUse is APIOperationPolicies {
            OperationPolicy[]? request = operationPoliciesToUse.request;
            if request is OperationPolicy[] {
                model:HTTPRouteFilter requestHeaderFilter = {'type: "RequestHeaderModifier",
                    requestHeaderModifier: self.extractHttpHeaderFilterData(request, organization)};
                routeFilters.push(requestHeaderFilter);
            }
            OperationPolicy[]? response = operationPoliciesToUse.response;
            if response is OperationPolicy[] {
                model:HTTPRouteFilter responseHeaderFilter = {'type: "ResponseHeaderModifier",
                    responseHeaderModifier: self.extractHttpHeaderFilterData(response, organization)};
                routeFilters.push(responseHeaderFilter);
            }
        }
        return routeFilters;
    }

    isolated function extractHttpHeaderFilterData(OperationPolicy[] operationPolicy, commons:Organization organization) returns model:HTTPHeaderFilter {
        model:HTTPHeader[] setPolicies = [];
        string[] removePolicies = [];
        foreach OperationPolicy policy in operationPolicy {
            string policyName = policy.policyName;
            OperationPolicyParameters[]? policyParameters = policy.parameters;
            if (policyParameters is OperationPolicyParameters[]) {
                if (policyName == "addHeader") {
                    model:HTTPHeader httpHeader = {
                        name: <string>policyParameters[0].headerName,
                        value: <string>policyParameters[0].headerValue
                    };
                    setPolicies.push(httpHeader);
                }
                if (policyName == "removeHeader") {
                    string httpHeader = <string>policyParameters[0].headerName;
                    removePolicies.push(httpHeader);
                }
            }
        }
        model:HTTPHeaderFilter headerModifier = {};
        if (setPolicies != []) {
            headerModifier.set = setPolicies;
        }
        if (removePolicies != []) {
            headerModifier.remove = removePolicies;
        }
        return headerModifier;
    }

    isolated function generatePrefixMatch(API api, model:Endpoint endpoint, APIOperations operation, string endpointType) returns string {
        string target = operation.target ?: "/*";
        string[] splitValues = regex:split(target, "/");
        string generatedPath = "";
        int pathparamCount = 1;
        if (target == "/*") {
            generatedPath = "\\1";
        } else {
            foreach int i in 0 ..< splitValues.length() {
                if splitValues[i].trim().length() > 0 {
                    // path contains path param
                    if regex:matches(splitValues[i], "\\{.*\\}") {
                        generatedPath = generatedPath + "/" + regex:replaceAll(splitValues[i].trim(), "\\{.*\\}", "\\" + pathparamCount.toString());
                        pathparamCount += 1;
                    } else {
                        generatedPath = generatedPath + "/" + splitValues[i];
                    }
                }
            }
        }

        if generatedPath.endsWith("/*") {
            int lastSlashIndex = <int>generatedPath.lastIndexOf("/", generatedPath.length());
            generatedPath = generatedPath.substring(0, lastSlashIndex) + "///" + pathparamCount.toString();
        }
        if endpoint.serviceEntry {
            return generatedPath.trim();
        }
        string path = self.getPath(<string>endpoint.url);
        if path.endsWith("/") {
            if generatedPath.startsWith("/") {
                return path.substring(0, path.length() - 1) + generatedPath;
            }
        }
        return path + generatedPath;
    }

    public isolated function retrievePathPrefix(string context, string 'version, string operation, commons:Organization organization) returns string {
        string fullContext = self.returnFullContext(context, 'version);
        string[] splitValues = regex:split(operation, "/");
        string generatedPath = fullContext;
        if (operation == "/*") {
            return generatedPath + "(.*)";
        }
        foreach string pathPart in splitValues {
            if pathPart.trim().length() > 0 {
                // path contains path param
                if regex:matches(pathPart, "\\{.*\\}") {
                    generatedPath = generatedPath + "/" + regex:replaceAll(pathPart.trim(), "\\{.*\\}", "(.*)");
                } else {
                    generatedPath = generatedPath + "/" + pathPart;
                }
            }
        }

        if generatedPath.endsWith("/*") {
            int lastSlashIndex = <int>generatedPath.lastIndexOf("/", generatedPath.length());
            generatedPath = generatedPath.substring(0, lastSlashIndex) + "(.*)";
        }
        return generatedPath.trim();
    }

    private isolated function retrieveGeneratedBackend(API api, model:Endpoint endpoint, string endpointType) returns model:HTTPBackendRef[] {
        model:HTTPBackendRef httpBackend = {
            namespace: <string>endpoint.namespace,
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

    private isolated function retrieveMatches(API api, APIOperations apiOperation, commons:Organization organization) returns model:HTTPRouteMatch[] {
        model:HTTPRouteMatch[] httpRouteMatch = [];
        model:HTTPRouteMatch httpRoute = self.retrieveHttpRouteMatch(api, apiOperation, organization);

        httpRouteMatch.push(httpRoute);
        return httpRouteMatch;
    }
    private isolated function retrieveHttpRouteMatch(API api, APIOperations apiOperation, commons:Organization organization) returns model:HTTPRouteMatch {

        return {method: <string>apiOperation.verb, path: {'type: "RegularExpression", value: self.retrievePathPrefix(api.context, api.'version, apiOperation.target ?: "/*", organization)}};
    }
    isolated function retrieveGeneratedSwaggerDefinition(API api, string? definition) returns json|commons:APKError {
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
                boolean? authTypeEnabled = apiOperation.authTypeEnabled;
                if authTypeEnabled is boolean {
                    uriTemplate.setAuthEnabled(authTypeEnabled);
                } else {
                    uriTemplate.setAuthEnabled(true);
                }
                string[]? scopes = apiOperation.scopes;
                if scopes is string[] {
                    foreach string item in scopes {
                        runtimeModels:Scope scope = runtimeModels:newScope1();
                        scope.setId(item);
                        scope.setName(item);
                        scope.setKey(item);
                        uriTemplate.setScopes(scope);
                    }
                }
                _ = uritemplatesSet.add(uriTemplate);
            }
        }
        api1.setUriTemplates(uritemplatesSet);
        string?|runtimeapi:APIManagementException retrievedDefinition = "";
        if definition is string && definition.toString().trim().length() > 0 {
            retrievedDefinition = runtimeUtil:RuntimeAPICommonUtil_generateDefinition2(api1, definition);
        } else {
            retrievedDefinition = runtimeUtil:RuntimeAPICommonUtil_generateDefinition(api1);
        }
        if retrievedDefinition is string && retrievedDefinition.toString().trim().length() > 0 {
            json|error jsonString = value:fromJsonString(retrievedDefinition);
            if jsonString is json {
                return jsonString;
            } else {
                log:printError("Error on converting to json", jsonString);
                return error("Error occured while generating openapi definition", code = 900920, message = "Error occured while generating openapi definition", statusCode = 500, description = "Error occured while generating openapi definition");
            }
        } else if retrievedDefinition is () {
            return "";
        } else {
            return error("Error occured while generating openapi definition", code = 900920, message = "Error occured while generating openapi definition", statusCode = 500, description = "Error occured while generating openapi definition");
        }
    }

    public isolated function generateAPIKey(string apiId, commons:Organization organization) returns APIKey|BadRequestError|NotFoundError|InternalServerErrorError {
        model:API|error api = getAPI(apiId, organization);
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

    public function retrieveAllApisAtStartup(map<map<model:API>>? apiMap, string? continueValue) returns error? {
        string? resultValue = continueValue;
        model:APIList|http:ClientError retrieveAllAPISResult;
        if resultValue is string {
            retrieveAllAPISResult = retrieveAllAPIS(resultValue);
        } else {
            retrieveAllAPISResult = retrieveAllAPIS(());
        }

        if retrieveAllAPISResult is model:APIList {
            model:ListMeta metadata = retrieveAllAPISResult.metadata;
            model:API[] items = retrieveAllAPISResult.items;
            if apiMap is map<map<model:API>> {
                putallAPIS(apiMap, items.clone());
            } else {
                lock {
                    putallAPIS(apilist, items.clone());
                }
            }

            string? continueElement = metadata.'continue;
            if continueElement is string {
                if continueElement.length() > 0 {
                    _ = check self.retrieveAllApisAtStartup(apiMap, continueElement);
                }
            }
            string? resourceVersion = metadata.'resourceVersion;
            if resourceVersion is string {
                setResourceVersion(resourceVersion);
            }
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
                    creationTimestamp: (),
                    labels: self.getLabels(api)
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

    isolated function deleteServiceMappings(model:API api) returns commons:APKError? {
        do {
            map<model:K8sServiceMapping> retrieveServiceMappingsForAPIResult = retrieveServiceMappingsForAPI(api).clone();
            model:ServiceMappingList|http:ClientError k8sServiceMapingsDeletionResponse = check getK8sServiceMapingsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, api.metadata.namespace);
            if k8sServiceMapingsDeletionResponse is model:ServiceMappingList {
                foreach model:K8sServiceMapping item in k8sServiceMapingsDeletionResponse.items {
                    retrieveServiceMappingsForAPIResult[<string>item.metadata.uid] = item;
                }
            } else {
                log:printError("Error occured while deleting service mapping");
            }
            string[] keys = retrieveServiceMappingsForAPIResult.keys();
            foreach string key in keys {
                model:K8sServiceMapping serviceMapping = retrieveServiceMappingsForAPIResult.get(key);
                http:Response|http:ClientError k8ServiceMappingDeletionResponse = deleteK8ServiceMapping(serviceMapping.metadata.name, serviceMapping.metadata.namespace);
                if k8ServiceMappingDeletionResponse is http:Response {
                    if k8ServiceMappingDeletionResponse.statusCode != http:STATUS_OK {
                        json responsePayLoad = check k8ServiceMappingDeletionResponse.getJsonPayload();
                        model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                        check self.handleK8sTimeout(statusResponse);
                    }
                } else {
                    log:printError("Error occured while deleting service mapping");
                }
            }
            return;
        } on fail var e {
            log:printError("Error occured deleting servicemapping", e);
            return error("Error occured deleting servicemapping", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }

    }

    public isolated function validateDefinition(http:Request message, boolean returnContent) returns InternalServerErrorError|BadRequestError|APIDefinitionValidationResponse|commons:APKError {
        do {
            DefinitionValidationRequest|BadRequestError definitionValidationRequest = check self.mapApiDefinitionPayload(message);
            if definitionValidationRequest is DefinitionValidationRequest {
                runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|error|BadRequestError validationResponse = self.validateAndRetrieveDefinition(definitionValidationRequest.'type, definitionValidationRequest.url, definitionValidationRequest.inlineAPIDefinition, definitionValidationRequest.content, definitionValidationRequest.fileName);
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
            else {
                return definitionValidationRequest;
            }
        } on fail var e {
            commons:APKError apkError = error("Internal Error", e, code = 900900, message = "InternalServer Error", description = "InternalServer Error", statusCode = 500);
            return apkError;
        }
    }

    private isolated function mapApiDefinitionPayload(http:Request message) returns DefinitionValidationRequest|BadRequestError|error {
        string|() url = ();
        string|() fileName = ();
        byte[]|() fileContent = ();
        string definitionType = "REST";
        string|() inlineAPIDefinition = ();
        mime:Entity[] payLoadParts = check message.getBodyParts();
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
        DefinitionValidationRequest definitionValidationRequest = {
            content: fileContent,
            fileName: fileName,
            inlineAPIDefinition: inlineAPIDefinition,
            url: url,
            'type: definitionType
        };
        return definitionValidationRequest;
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

    isolated function handleK8sTimeout(model:Status errorStatus) returns commons:APKError {
        model:StatusDetails? details = errorStatus.details;
        if details is model:StatusDetails {
            if details.retryAfterSeconds is int && details.retryAfterSeconds >= 0 {
                // K8s api level ratelimit hit.
                log:printError("K8s API Timeout happens when invoking k8s api");
            }
        }
        commons:APKError apkError = error("Internal Server Error", code = 900900, message = "Internal Server Error", statusCode = 500, description = "Internal Server Error");
        return apkError;
    }

    isolated function createBackendService(API api, APIOperations? apiOperation, string endpointType, commons:Organization organization, string url) returns model:Service {
        string nameSpace = getNameSpace(runtimeConfiguration.apiCreationNamespace);
        model:Service backendService = {
            metadata: {
                name: getBackendServiceUid(api, apiOperation, endpointType, organization),
                namespace: nameSpace,
                uid: (),
                creationTimestamp: (),
                labels: self.getLabels(api)
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

    public isolated function validateAPIExistence(string query, commons:Organization organization) returns NotFoundError|BadRequestError|http:Ok {
        int? indexOfColon = query.indexOf(":", 0);
        boolean exist = false;
        if indexOfColon is int && indexOfColon > 0 {
            string keyWord = query.substring(0, indexOfColon);
            string keyWordValue = query.substring(keyWord.length() + 1, query.length());
            if keyWord == "name" {
                exist = self.validateName(keyWordValue, organization);
            } else if keyWord == "context" {
                exist = self.validateContext(keyWordValue, organization);
            } else {
                BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord " + keyWord}};
                return badRequest;
            }
        } else {
            // Consider full string as name;
            exist = self.validateName(query, organization);
        }
        if exist {
            http:Ok ok = {};
            return ok;
        } else {
            NotFoundError notFound = {body: {code: 900914, message: "context/name doesn't exist"}};
            return notFound;
        }
    }

    public isolated function importDefinition(http:Request payload, commons:Organization organization) returns commons:APKError|CreatedAPI|InternalServerErrorError|BadRequestError {
        do {
            ImportDefintionRequest|BadRequestError importDefinitionRequest = check self.mapImportDefinitionRequest(payload);
            if importDefinitionRequest is ImportDefintionRequest {
                lock {
                    if (ALLOWED_API_TYPES.indexOf(importDefinitionRequest.'type) is ()) {
                        BadRequestError badRequest = {body: {code: 900912, message: "unsupported API Type."}};
                        return badRequest.clone();
                    }
                }
                runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|BadRequestError validateAndRetrieveDefinitionResult = check self.validateAndRetrieveDefinition(importDefinitionRequest.'type, importDefinitionRequest.url, importDefinitionRequest.inlineAPIDefinition, importDefinitionRequest.content, importDefinitionRequest.fileName);
                if validateAndRetrieveDefinitionResult is runtimeapi:APIDefinitionValidationResponse {
                    if validateAndRetrieveDefinitionResult.isValid() {
                        runtimeapi:APIDefinition parser = validateAndRetrieveDefinitionResult.getParser();
                        log:printDebug("content available ==", contentAvailable = (validateAndRetrieveDefinitionResult.getContent() is string));
                        utilapis:Set|runtimeapi:APIManagementException uRITemplates = parser.getURITemplates(<string>validateAndRetrieveDefinitionResult.getContent());
                        if uRITemplates is utilapis:Set {
                            API additionalPropertes = importDefinitionRequest.additionalPropertes;
                            APIOperations[]? operations = additionalPropertes.operations;
                            if !(operations is APIOperations[]) {
                                operations = [];
                            }
                            lang:Object[] uriTemplates = check uRITemplates.toArray();
                            foreach lang:Object uritemplate in uriTemplates {
                                runtimeModels:URITemplate template = check java:cast(uritemplate);
                                if operations is APIOperations[] {
                                    operations.push({target: template.getUriTemplate(), authTypeEnabled: template.isAuthEnabled(), verb: template.getHTTPVerb().toString().toUpperAscii()});
                                }
                            }
                            additionalPropertes.operations = operations;
                            return self.createAPI(additionalPropertes, validateAndRetrieveDefinitionResult.getContent(), organization);
                        }
                        log:printError("Error occured retrieving uri templates from definition", uRITemplates);
                        runtimeapi:JAPIManagementException excetion = check uRITemplates.ensureType(runtimeapi:JAPIManagementException);
                        runtimeapi:ErrorHandler errorHandler = excetion.getErrorHandler();
                        BadRequestError badeRequest = {body: {code: errorHandler.getErrorCode(), message: errorHandler.getErrorMessage().toString()}};
                        return badeRequest;
                    }
                    // Error definition.
                    ErrorListItem[] errorItems = [];
                    utilapis:ArrayList errorItemsResult = validateAndRetrieveDefinitionResult.getErrorItems();
                    foreach int i in 0 ... errorItemsResult.size() - 1 {
                        runtimeapi:ErrorItem errorItem = check java:cast(errorItemsResult.get(i));
                        ErrorListItem errorListItem = {code: errorItem.getErrorCode().toString(), message: <string>errorItem.getErrorMessage(), description: errorItem.getErrorDescription()};
                        errorItems.push(errorListItem);
                    }
                    BadRequestError badRequest = {body: {code: 90091, message: "Invalid API Definition", 'error: errorItems}};
                    return badRequest;
                } else if validateAndRetrieveDefinitionResult is BadRequestError {
                    return validateAndRetrieveDefinitionResult;
                } else {
                    log:printError("Error occured creating api from defintion", validateAndRetrieveDefinitionResult);
                    runtimeapi:JAPIManagementException excetion = check validateAndRetrieveDefinitionResult.ensureType(runtimeapi:JAPIManagementException);
                    runtimeapi:ErrorHandler errorHandler = excetion.getErrorHandler();
                    BadRequestError badeRequest = {body: {code: errorHandler.getErrorCode(), message: errorHandler.getErrorMessage().toString()}};
                    return badeRequest;
                }
            } else {
                return <BadRequestError>importDefinitionRequest;
            }
        } on fail var e {
            log:printError("Error occured importing API", e);
            InternalServerErrorError internalError = {body: {code: 900900, message: "Internal Error."}};
            return internalError;
        }
    }
    private isolated function validateAndRetrieveDefinition(string 'type, string? url, string? inlineAPIDefinition, byte[]? content, string? fileName) returns runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|error|BadRequestError {
        runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|error validationResponse;
        boolean inlineApiDefinitionAvailable = inlineAPIDefinition is string;
        boolean fileAvailable = fileName is string && content is byte[];
        boolean urlAvailble = url is string;
        boolean typeAvailable = 'type.length() > 0;

        if !typeAvailable {
            BadRequestError badRequest = {body: {code: 90914, message: "type field unavailable"}};
            return badRequest;
        }
        if url is string {
            if (fileAvailable || inlineApiDefinitionAvailable) {
                BadRequestError badRequest = {body: {code: 90914, message: "multiple fields of  url,file,inlineAPIDefinition given"}};
                return badRequest;
            }
            string|error retrieveDefinitionFromUrlResult = self.retrieveDefinitionFromUrl(url);
            if retrieveDefinitionFromUrlResult is string {
                validationResponse = runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition('type, [], retrieveDefinitionFromUrlResult, fileName ?: "", true);
            } else {
                log:printError("Error occured while retrieving definition from url", retrieveDefinitionFromUrlResult);
                BadRequestError badRequest = {body: {code: 900900, message: "retrieveDefinitionFromUrlResult"}};
                return badRequest;
            }
        } else if fileName is string && content is byte[] {
            if (urlAvailble || inlineApiDefinitionAvailable) {
                BadRequestError badRequest = {body: {code: 90914, message: "multiple fields of  url,file,inlineAPIDefinition given"}};
                return badRequest;
            }
            validationResponse = runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition('type, <byte[]>content, "", <string>fileName, true);
        } else if inlineAPIDefinition is string {
            if (fileAvailable || urlAvailble) {
                BadRequestError badRequest = {body: {code: 90914, message: "multiple fields of  url,file,inlineAPIDefinition given"}};
                return badRequest;
            }
            validationResponse = runtimeUtil:RuntimeAPICommonUtil_validateOpenAPIDefinition('type, <byte[]>[], <string>inlineAPIDefinition, "", true);
        } else {
            BadRequestError badRequest = {body: {code: 90914, message: "atleast one of the field required"}};
            return badRequest;
        }
        return validationResponse;
    }
    private isolated function mapImportDefinitionRequest(http:Request message) returns ImportDefintionRequest|error|BadRequestError {
        string|() url = ();
        string|() fileName = ();
        byte[]|() fileContent = ();
        string|() inlineAPIDefinition = ();
        string|() additionalProperties = ();
        string|() 'type = ();
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
                } else if fieldName == "inlineAPIDefinition" {
                    inlineAPIDefinition = check payLoadPart.getText();
                } else if fieldName == "additionalProperties" {
                    additionalProperties = check payLoadPart.getText();
                } else if fieldName == "type" {
                    'type = check payLoadPart.getText();
                }
            }
        }
        if 'type is () {
            BadRequestError badRequest = {body: {code: 90914, message: "type field unavailable"}};
            return badRequest;
        }
        if url is () && fileName is () && inlineAPIDefinition is () && fileContent is () {
            BadRequestError badRequest = {body: {code: 90914, message: "atleast one of the field required (file,inlineApiDefinition,url)."}};
            return badRequest;
        }
        if additionalProperties is () || additionalProperties.length() == 0 {
            BadRequestError badRequest = {body: {code: 90914, message: "additionalProperties not provided."}};
            return badRequest;
        }
        json apiObject = check value:fromJsonString(additionalProperties);
        API api = check apiObject.cloneWithType(API);
        ImportDefintionRequest importDefintionRequest = {
            fileName: fileName,
            inlineAPIDefinition: inlineAPIDefinition,
            additionalPropertes: api,
            url: url,
            content: fileContent,
            'type: 'type
        };
        return importDefintionRequest;
    }
    public isolated function copyAPI(string newVersion, string? serviceId, string apiId, commons:Organization organization) returns CreatedAPI|NotFoundError|BadRequestError|commons:APKError {
        // validating API existence.
        if newVersion.trim().length() == 0 || apiId.trim().length() == 0 {
            BadRequestError badRequest = {body: {code: 900912, message: "new Version/APIID not exist."}};
            return badRequest;
        }
        do {
            API|NotFoundError api = check self.getAPIById(apiId, organization);
            if api is API {
                model:APIArtifact apiArtifact = check self.getApiArtifact(api, organization);
                // validating version
                if isAPIVersionExist(api.name, newVersion, organization) {
                    BadRequestError badRequest = {body: {code: 900920, message: newVersion + " already exist."}};
                    return badRequest;
                }
                //validating serviceuid if exist.
                if apiArtifact.serviceMapping.length() > 0 {
                    if serviceId is string && serviceId.toString().trim().length() > 0 {
                        Service|error serviceById = getServiceById(serviceId);
                        if serviceById is Service {
                            check self.prepareApiArtifactforNewVersion(apiArtifact, serviceById, api, newVersion, organization);
                            model:API deployAPIToK8sResult = check self.deployAPIToK8s(apiArtifact);
                            CreatedAPI createdAPI = {body: {name: deployAPIToK8sResult.spec.apiDisplayName, context: self.returnFullContext(deployAPIToK8sResult.spec.context, deployAPIToK8sResult.spec.apiVersion), 'version: deployAPIToK8sResult.spec.apiVersion, id: deployAPIToK8sResult.metadata.uid}};
                            return createdAPI;
                        } else {
                            BadRequestError badRequest = {body: {code: 900921, message: serviceId + " service not exist."}};
                            return badRequest;
                        }
                    }
                }
                check self.prepareApiArtifactforNewVersion(apiArtifact, (), api, newVersion, organization);
                model:API deployAPIToK8sResult = check self.deployAPIToK8s(apiArtifact);
                CreatedAPI createdAPI = {body: {name: deployAPIToK8sResult.spec.apiDisplayName, context: self.returnFullContext(deployAPIToK8sResult.spec.context, deployAPIToK8sResult.spec.apiVersion), 'version: deployAPIToK8sResult.spec.apiVersion, id: deployAPIToK8sResult.metadata.uid}};
                return createdAPI;
            } else {
                return <NotFoundError>api;
            }
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            log:printError("Internal Error occured", e);
            return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = 500);
        }
    }
    private isolated function prepareApiArtifactforNewVersion(model:APIArtifact apiArtifact, Service? serviceEntry, API oldAPI, string newVersion, commons:Organization organization) returns error? {
        API newAPI = {
            name: oldAPI.name,
            context: regex:replace(oldAPI.context, oldAPI.'version, newVersion),
            'version: newVersion
        };
        check self.prepareConfigMap(apiArtifact, oldAPI, newAPI);
        self.prepareHttpRoute(apiArtifact, serviceEntry, oldAPI, newAPI, PRODUCTION_TYPE, organization);
        self.prepareHttpRoute(apiArtifact, serviceEntry, oldAPI, newAPI, SANDBOX_TYPE, organization);
        self.prepareK8sServiceMapping(apiArtifact, serviceEntry, oldAPI, newAPI, organization);
        self.prepareAPICr(apiArtifact, oldAPI, newAPI, organization);
        apiArtifact.runtimeAPI = self.generateRuntimeAPIArtifact(newAPI, serviceEntry, organization);

    }
    private isolated function prepareK8sServiceMapping(model:APIArtifact apiArtifact, Service? serviceEntry, API oldAPI, API newAPI, commons:Organization organization) {
        model:K8sServiceMapping[] serviceMappings = apiArtifact.serviceMapping;
        foreach model:K8sServiceMapping serviceMapping in serviceMappings {
            serviceMapping.metadata.name = self.getServiceMappingEntryName(apiArtifact.uniqueId);
            serviceMapping.metadata.labels = self.getLabels(newAPI);
            serviceMapping.spec.apiRef.name = apiArtifact.uniqueId;
        }
    }
    private isolated function prepareHttpRoute(model:APIArtifact apiArtifact, Service? serviceEntry, API oldAPI, API newAPI, string endpointType, commons:Organization organization) {
        model:Httproute[] httproutes;
        if endpointType == PRODUCTION_TYPE {
            httproutes = apiArtifact.productionRoute;
        } else {
            httproutes = apiArtifact.sandboxRoute;
        }
        map<string> serviceMapping = {};
        map<string> extenstionRefMappings = {};
        foreach model:Httproute httproute in httproutes {
            httproute.metadata.name = retrieveHttpRouteRefName(newAPI, endpointType, organization);
            httproute.metadata.labels = self.getLabels(newAPI);
            model:HTTPRouteRule[] routeRules = httproute.spec.rules;
            foreach model:HTTPRouteRule routeRule in routeRules {
                model:HTTPBackendRef[]? backendRefs = routeRule.backendRefs;
                if backendRefs is model:HTTPBackendRef[] {
                    foreach model:HTTPBackendRef backendRef in backendRefs {
                        if serviceMapping.hasKey(backendRef.name) {
                            string newServiceName = serviceMapping.get(backendRef.name);
                            backendRef.name = newServiceName;
                        } else {
                            [string, string]? prepareBackendRefResult = self.prepareBackendRef(backendRef, apiArtifact, serviceEntry, oldAPI, newAPI, endpointType, organization);
                            if prepareBackendRefResult is [string, string] {
                                serviceMapping[prepareBackendRefResult[0]] = prepareBackendRefResult[1];
                            }
                        }
                    }
                }
                model:HTTPRouteMatch[]? matches = routeRule.matches;
                if matches is model:HTTPRouteMatch[] {
                    foreach model:HTTPRouteMatch routeMatch in matches {
                        model:HTTPPathMatch? path = routeMatch.path;
                        if path is model:HTTPPathMatch {
                            path.value = regex:replace(path.value, oldAPI.'version, newAPI.'version);
                        }
                    }
                }
                model:HTTPRouteFilter[]? filters = routeRule.filters;
                if filters is model:HTTPRouteFilter[] {
                    foreach model:HTTPRouteFilter filter in filters {
                        model:LocalObjectReference? extensionRef = filter.extensionRef;
                        if extensionRef is model:LocalObjectReference {
                            if extenstionRefMappings.hasKey(extensionRef.name) {
                                extensionRef.name = extenstionRefMappings.get(extensionRef.name);
                            }
                            if extensionRef.kind == "Authentication" {
                                if apiArtifact.authenticationMap.hasKey(extensionRef.name) {
                                    model:Authentication authentication = apiArtifact.authenticationMap.get(extensionRef.name).clone();
                                    model:Authentication newAuthenticationCR = self.prepareAuthenticationCR(apiArtifact, newAPI, authentication, endpointType, organization);
                                    _ = apiArtifact.authenticationMap.remove(extensionRef.name);
                                    apiArtifact.authenticationMap[newAuthenticationCR.metadata.name] = newAuthenticationCR;
                                    extenstionRefMappings[extensionRef.name] = newAuthenticationCR.metadata.name;
                                    extensionRef.name = newAuthenticationCR.metadata.name;
                                }
                            } else if extensionRef.kind == "Scope" {
                                if apiArtifact.scopes.hasKey(extensionRef.name) {
                                    model:Scope scope = apiArtifact.scopes.get(extensionRef.name).clone();
                                    model:Scope newScopeCR = self.prepareScopeCR(apiArtifact, newAPI, scope, organization);
                                    _ = apiArtifact.scopes.remove(extensionRef.name);
                                    apiArtifact.scopes[newScopeCR.metadata.name] = newScopeCR;
                                    extenstionRefMappings[extensionRef.name] = newScopeCR.metadata.name;
                                    extensionRef.name = newScopeCR.metadata.name;
                                }
                            }
                        }
                    }
                }
            }
            self.prepareBackendPolicyCR(apiArtifact, newAPI, endpointType, serviceMapping, organization);
        }
    }
    private isolated function prepareScopeCR(model:APIArtifact apiArtifact, API api, model:Scope scope, commons:Organization organization) returns model:Scope {
        scope.metadata.name = uuid:createType1AsString();
        scope.metadata.labels = self.getLabels(api);
        return scope;
    }
    private isolated function prepareBackendPolicyCR(model:APIArtifact apiArtifact, API api, string endpointType, map<string> backendServiceMapping, commons:Organization organization) {
        foreach model:BackendPolicy backendPolicy in apiArtifact.backendPolicies {
            backendPolicy.metadata.name = getBackendPolicyUid(api, endpointType, organization);
            string oldServiceName = backendPolicy.spec.targetRef.name;
            if backendServiceMapping.hasKey(oldServiceName) {
                backendPolicy.spec.targetRef.name = backendServiceMapping.get(oldServiceName);
            }
        }
    }
    private isolated function prepareAuthenticationCR(model:APIArtifact apiArtifact, API api, model:Authentication authentication, string endpointType, commons:Organization organization) returns model:Authentication {
        authentication.metadata.name = self.retrieveDisableAuthenticationRefName(api, endpointType, organization);
        authentication.metadata.labels = self.getLabels(api);
        authentication.spec.targetRef.name = retrieveHttpRouteRefName(api, endpointType, organization);
        return authentication;
    }
    private isolated function prepareBackendRef(model:HTTPBackendRef backendRef, model:APIArtifact apiArtifact, Service? serviceEntry, API oldAPI, API newAPI, string endpointType, commons:Organization organization) returns [string, string]? {
        if apiArtifact.serviceMapping.length() >= 1 {
            if serviceEntry is Service {
                backendRef.name = serviceEntry.name;
                backendRef.namespace = serviceEntry.namespace;
                apiArtifact.serviceMapping.removeAll();
                self.generateAndSetK8sServiceMapping(apiArtifact, newAPI, serviceEntry, getNameSpace(runtimeConfiguration.apiCreationNamespace));
            }
        } else {
            map<model:Service> backendServices = apiArtifact.backendServices;
            string oldBackendRefName = backendRef.name;
            if backendServices.hasKey(oldBackendRefName) {
                model:Service 'service = backendServices.get(oldBackendRefName).clone();
                if 'service.metadata.name.includes("-api-") {
                    'service.metadata.name = getBackendServiceUid(newAPI, (), endpointType, organization);
                } else {
                    'service.metadata.name = getBackendServiceUid(newAPI, {}, endpointType, organization);
                }
                'service.metadata.labels = self.getLabels(newAPI);
                _ = backendServices.remove(oldBackendRefName);
                backendServices['service.metadata.name] = 'service;
                backendRef.name = 'service.metadata.name;
                return [oldBackendRefName, 'service.metadata.name];
            }
        }
        return;
    }
    private isolated function prepareAPICr(model:APIArtifact apiArtifact, API oldAPI, API newAPI, commons:Organization organization) {
        string uuid = getUniqueIdForAPI(oldAPI.name, newAPI.'version, organization);
        apiArtifact.uniqueId = uuid;
        model:API? api = apiArtifact.api;
        if api is model:API {
            string oldName = api.metadata.name;
            api.spec.apiVersion = newAPI.'version;
            api.metadata.name = uuid;
            api.metadata.labels = self.getLabels(newAPI);
            api.spec.context = newAPI.context;
            string[] prodHTTPRouteRefs = [];
            foreach model:Httproute httpRoute in apiArtifact.productionRoute {
                prodHTTPRouteRefs.push(httpRoute.metadata.name);
            }
            if prodHTTPRouteRefs.length() > 0 {
                api.spec.prodHTTPRouteRefs = prodHTTPRouteRefs;
            }
            string[] sandHTTPRouteRefs = [];
            foreach model:Httproute httpRoute in apiArtifact.sandboxRoute {
                sandHTTPRouteRefs.push(httpRoute.metadata.name);
            }
            if sandHTTPRouteRefs.length() > 0 {
                api.spec.sandHTTPRouteRefs = sandHTTPRouteRefs;
            }
            string? definitionFileRef = api.spec.definitionFileRef;
            if definitionFileRef is string {
                api.spec.definitionFileRef = regex:replaceAll(definitionFileRef, oldName, uuid);
            }
        }
    }
    private isolated function prepareConfigMap(model:APIArtifact apiArtifact, API oldAPI, API newAPI) returns error? {
        model:ConfigMap? definition = apiArtifact.definition;
        if definition is model:ConfigMap {
            json definitionJson = check self.getDefinitionFromConfigMap(definition);
            json info = check definitionJson.info;
            map<json> infoElement = <map<json>>info;
            infoElement["version"] = newAPI.'version;
            map<json> definitionMap = <map<json>>definitionJson;
            definitionMap["info"] = infoElement;
            self.retrieveGeneratedConfigmapForDefinition(apiArtifact, newAPI, definitionMap, apiArtifact.uniqueId);
        }
    }

    private isolated function getApiArtifact(API api, commons:Organization organization) returns model:APIArtifact|error {
        model:API k8sapi = check getAPI(<string>api.id, organization);
        model:APIArtifact apiArtifact = {uniqueId: k8sapi.metadata.name};
        // retrieveConfigmap
        string? definitionFileRef = k8sapi.spec.definitionFileRef;
        if definitionFileRef is string {
            model:ConfigMap|error? definitionConfigmap = check self.getDefinitionConfigmap(definitionFileRef, k8sapi.metadata.namespace);
            if definitionConfigmap is model:ConfigMap {
                apiArtifact.definition = self.sanitizeConfigMapData(definitionConfigmap);
            }
        }
        json[]? prodHTTPRouteRefs = k8sapi.spec.prodHTTPRouteRefs;
        if prodHTTPRouteRefs is json[] && prodHTTPRouteRefs.length() > 0 {
            foreach json prodHTTPRouteRef in prodHTTPRouteRefs {
                model:Httproute httpRoute = check getHttpRoute(prodHTTPRouteRef.toString(), k8sapi.metadata.namespace);
                apiArtifact.productionRoute.push(self.sanitizeHttpRoute(httpRoute));
            }
        }
        json[]? sandHTTPRouteRefs = k8sapi.spec.sandHTTPRouteRefs;
        if sandHTTPRouteRefs is json[] && sandHTTPRouteRefs.length() > 0 {
            foreach json sandHTTPRouteRef in sandHTTPRouteRefs {
                model:Httproute httpRoute = check getHttpRoute(sandHTTPRouteRef.toString(), k8sapi.metadata.namespace);
                apiArtifact.sandboxRoute.push(self.sanitizeHttpRoute(httpRoute));
            }
        }
        model:ServiceMappingList k8sServiceMapingsForAPI = check getK8sServiceMapingsForAPI(api.name, api.'version, k8sapi.metadata.namespace);
        foreach model:K8sServiceMapping serviceMapping in k8sServiceMapingsForAPI.items {
            apiArtifact.serviceMapping.push(self.sanitizeServiceMapping(serviceMapping));
        }
        model:ServiceList backendServicesForAPI = check getBackendServicesForAPI(api.name, api.'version, k8sapi.metadata.namespace);
        foreach model:Service serviceEntry in backendServicesForAPI.items {
            apiArtifact.backendServices[serviceEntry.metadata.name] = self.sanitizeBackendService(serviceEntry);
        }
        model:AuthenticationList authenticationCrsForAPI = check getAuthenticationCrsForAPI(api.name, api.'version, k8sapi.metadata.namespace);
        foreach model:Authentication authentication in authenticationCrsForAPI.items {
            apiArtifact.authenticationMap[authentication.metadata.name] = self.sanitizeAuthenticationCrs(authentication);
        }
        model:BackendPolicyList backendPolicyList = check getBackendPolicyCRsForAPI(api.name, api.'version, k8sapi.metadata.namespace);
        foreach model:BackendPolicy backendPolicy in backendPolicyList.items {
            apiArtifact.backendPolicies[backendPolicy.metadata.name] = self.sanitizeBackendPolicyCrs(backendPolicy);
        }
        model:RuntimeAPI|http:ClientError internalAPI = getInternalAPI(k8sapi.metadata.name, k8sapi.metadata.namespace);
        if internalAPI is model:RuntimeAPI {
            apiArtifact.runtimeAPI = self.sanitizeRuntimeAPI(internalAPI);
        } else if (internalAPI is http:ApplicationResponseError) {
            if internalAPI.detail().statusCode != 404 {
                return internalAPI;
            }
        } else {
            return internalAPI;
        }
        model:ScopeList scopeList = check getScopeCrsForAPI(k8sapi.spec.apiDisplayName, k8sapi.spec.apiVersion, k8sapi.metadata.namespace);
        foreach model:Scope scope in scopeList.items {
            apiArtifact.scopes[scope.metadata.name] = self.sanitizeScopeCR(scope);
        }
        apiArtifact.api = self.sanitizeAPICR(k8sapi);
        return apiArtifact;
    }
    private isolated function sanitizeScopeCR(model:Scope scope) returns model:Scope {
        model:Scope sanitizedScope = {
            metadata: {name: scope.metadata.name, namespace: scope.metadata.namespace, labels: scope.metadata.labels},
            spec: scope.spec
        };
        return sanitizedScope;
    }
    private isolated function sanitizeRuntimeAPI(model:RuntimeAPI runtimeAPI) returns model:RuntimeAPI {
        model:RuntimeAPI sanitizedAPI = {
            metadata: {name: runtimeAPI.metadata.name, namespace: runtimeAPI.metadata.namespace, labels: runtimeAPI.metadata.labels},
            spec: runtimeAPI.spec
        };
        return sanitizedAPI;
    }
    private isolated function sanitizeAPICR(model:API api) returns model:API {
        model:API modifiedAPI = {
            metadata: {name: api.metadata.name, namespace: api.metadata.namespace},
            spec: {apiDisplayName: api.spec.apiDisplayName, apiType: api.spec.apiType, apiVersion: api.spec.apiVersion, context: api.spec.context, organization: api.spec.organization}
        };
        if api.spec.definitionFileRef is string && api.spec.definitionFileRef.toString().trim().length() > 0 {
            modifiedAPI.spec.definitionFileRef = api.spec.definitionFileRef;
        }
        string[]|() prodHTTPRouteRefs = api.spec.prodHTTPRouteRefs;
        if prodHTTPRouteRefs is string[] && prodHTTPRouteRefs.length() > 0 {
            modifiedAPI.spec.prodHTTPRouteRefs = prodHTTPRouteRefs;
        }
        string[]|() sandHTTPRouteRefs = api.spec.sandHTTPRouteRefs;
        if sandHTTPRouteRefs is string[] && sandHTTPRouteRefs.length() > 0 {
            modifiedAPI.spec.sandHTTPRouteRefs = sandHTTPRouteRefs;
        }
        return modifiedAPI;
    }

    private isolated function sanitizeConfigMapData(model:ConfigMap configmap) returns model:ConfigMap {
        return {
            metadata: {
                name: configmap.metadata.name,
                namespace: configmap.metadata.namespace,
                labels: configmap.metadata.labels
            },
            data: configmap.data,
            binaryData: configmap.binaryData
        };
    }

    private isolated function sanitizeHttpRoute(model:Httproute httproute) returns model:Httproute {
        return {
            metadata: {
                name: httproute.metadata.name,
                namespace: httproute.metadata.namespace,
                labels: httproute.metadata.labels
            },
            spec: httproute.spec
        };
    }

    private isolated function sanitizeServiceMapping(model:K8sServiceMapping serviceMapping) returns model:K8sServiceMapping {
        return {
            metadata: {
                name: serviceMapping.metadata.name,
                namespace: serviceMapping.metadata.namespace,
                labels: serviceMapping.metadata.labels
            },
            spec: serviceMapping.spec
        };
    }

    private isolated function sanitizeBackendService(model:Service serviceEntry) returns model:Service {
        return {
            metadata: {
                name: serviceEntry.metadata.name,
                namespace: serviceEntry.metadata.namespace,
                labels: serviceEntry.metadata.labels
            },
            spec: serviceEntry.spec
        };
    }

    private isolated function sanitizeAuthenticationCrs(model:Authentication authentication) returns model:Authentication {
        return {
            metadata: {
                name: authentication.metadata.name,
                namespace: authentication.metadata.namespace,
                labels: authentication.metadata.labels
            },
            spec: authentication.spec
        };
    }
    private isolated function sanitizeBackendPolicyCrs(model:BackendPolicy backendPolicy) returns model:BackendPolicy {
        return {
            metadata: {
                name: backendPolicy.metadata.name,
                namespace: backendPolicy.metadata.namespace,
                labels: backendPolicy.metadata.labels
            },
            spec: backendPolicy.spec
        };
    }
    public isolated function exportAPI(string? apiId, commons:Organization organization) returns commons:APKError|NotFoundError|http:Response|BadRequestError {
        if apiId is string {
            do {
                API|NotFoundError api = check self.getAPIById(apiId, organization);
                if api is API {
                    model:APIArtifact apiArtifact = check self.getApiArtifact(api, organization);
                    string zipDir = check file:createTempDir(uuid:createType1AsString());
                    model:API? k8sAPI = apiArtifact.api;
                    if k8sAPI is model:API {
                        _ = check self.convertAndStoreYamlFile(k8sAPI.toJsonString(), k8sAPI.metadata.name, zipDir, "api");
                    }
                    model:ConfigMap? definition = apiArtifact.definition;
                    if definition is model:ConfigMap {
                        _ = check self.convertAndStoreYamlFile(definition.toJsonString(), definition.metadata.name, zipDir, "definitions");
                    }
                    foreach model:Authentication authenticationCr in apiArtifact.authenticationMap {
                        _ = check self.convertAndStoreYamlFile(authenticationCr.toJsonString(), authenticationCr.metadata.name, zipDir, "policies/authentications");
                    }
                    foreach model:Service backendService in apiArtifact.backendServices {
                        _ = check self.convertAndStoreYamlFile(backendService.toJsonString(), backendService.metadata.name, zipDir, "backends");
                    }
                    foreach model:Httproute httpRoute in apiArtifact.productionRoute {
                        _ = check self.convertAndStoreYamlFile(httpRoute.toJsonString(), httpRoute.metadata.name, zipDir, "httproutes");
                    }
                    foreach model:Httproute httpRoute in apiArtifact.sandboxRoute {
                        _ = check self.convertAndStoreYamlFile(httpRoute.toJsonString(), httpRoute.metadata.name, zipDir, "httproutes");
                    }

                    foreach model:K8sServiceMapping servicemapping in apiArtifact.serviceMapping {
                        _ = check self.convertAndStoreYamlFile(servicemapping.toJsonString(), servicemapping.metadata.name, zipDir, "servicemappings");
                    }
                    foreach model:BackendPolicy backendPolicy in apiArtifact.backendPolicies {
                        _ = check self.convertAndStoreYamlFile(backendPolicy.toJsonString(), backendPolicy.metadata.name, zipDir, "policies/backendPolicies");

                    }
                    foreach model:Scope scope in apiArtifact.scopes {
                        _ = check self.convertAndStoreYamlFile(scope.toJsonString(), scope.metadata.name, zipDir, "scopes");
                    }
                    model:RuntimeAPI? runtimeAPI = apiArtifact.runtimeAPI;
                    if runtimeAPI is model:RuntimeAPI {
                        _ = check self.convertAndStoreYamlFile(runtimeAPI.toJsonString(), runtimeAPI.metadata.name, zipDir, "runtimeapi");
                    }
                    [string, string] zipName = check self.zipDirectory(apiId, zipDir);
                    http:Response response = new;
                    response.setFileAsPayload(zipName[1]);
                    response.addHeader("Content-Disposition", "attachment; filename=" + zipName[0]);
                    return response;
                } else {
                    return <NotFoundError>api;
                }
            } on fail var e {
                commons:APKError apkError = error("Internal Error occured when exporting api", e, message = "Internal Error.", code = 900900, description = "Internal Error.", statusCode = 500);
                return apkError;
            }
        } else {
            BadRequestError badRequest = {body: {code: 90912, message: "apiId not found in request."}};
            return badRequest;
        }
    }
    private isolated function convertAndStoreYamlFile(string jsonString, string fileName, string directroy, string? subDirectory) returns error? {
        runtimeUtil:YamlUtil yamlUtil = runtimeUtil:newYamlUtil1();
        string|() convertedYaml = check yamlUtil.fromJsonStringToYaml(jsonString);
        string fullPath = directroy;
        if convertedYaml is string {
            if subDirectory is string {
                fullPath = fullPath + file:pathSeparator + subDirectory;
                _ = check file:createDir(directroy + file:pathSeparator + subDirectory, file:RECURSIVE);
            }
            fullPath = fullPath + file:pathSeparator + fileName + ".yaml";
            _ = check io:fileWriteString(fullPath, convertedYaml);
        }
    }
    private isolated function zipDirectory(string apiId, string directoryPath) returns [string, string]|error {
        string zipName = apiId + ZIP_FILE_EXTENSTION;
        string zipPath = directoryPath + ZIP_FILE_EXTENSTION;
        _ = check runtimeUtil:ZIPUtils_zipDir(directoryPath, zipPath);
        return [zipName, zipPath];
    }
    public isolated function updateAPI(string apiId, API payload, string? definition, commons:Organization organization) returns API|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        API|NotFoundError api = check self.getAPIById(apiId, organization);
        if api is API {
            if payload.'type != api.'type {
                BadRequestError badRequest = {body: {code: 900930, message: "API Type change not supported from update."}};
                return badRequest;
            }
            if payload.context != api.context {
                BadRequestError badRequest = {body: {code: 900930, message: "Context change not supported from update."}};
                return badRequest;
            }
            if payload.'version != api.'version {
                BadRequestError badRequest = {body: {code: 900930, message: "Version change not supported from update."}};
                return badRequest;
            }
            API_serviceInfo? serviceInfo = payload.serviceInfo;
            model:Endpoint? endpoint = ();
            if serviceInfo is API_serviceInfo {
                if payload.serviceInfo is API_serviceInfo {
                    ServiceClient serviceClient = new;
                    Service|error serviceEntry = serviceClient.getServiceByNameandNamespace(<string>serviceInfo.name, <string>serviceInfo.namespace);
                    if serviceEntry is Service {
                        endpoint = {
                            port: self.retrievePort(serviceEntry),
                            namespace: serviceEntry.namespace,
                            name: serviceEntry.name,
                            serviceEntry: true
                        };
                    } else {
                        BadRequestError badRequest = {body: {code: 900930, message: "Service not found."}};
                        return badRequest;
                    }
                }
            }
            do {
                self.setDefaultOperationsIfNotExist(payload);
                string uniqueId = getUniqueIdForAPI(payload.name, payload.'version, organization);
                model:APIArtifact apiArtifact = {uniqueId: uniqueId};
                APIOperations[]? operations = payload.operations;
                if operations is APIOperations[] {
                    if operations.length() == 0 {
                        BadRequestError badRequestError = {body: {code: 90912, message: "Atleast one operation need to specified"}};
                        return badRequestError;
                    }
                    // Validating operation policies.
                    BadRequestError|() badRequestError = self.validateOperationPolicies(api.apiPolicies, operations, organization);
                    if (badRequestError is BadRequestError) {
                        return badRequestError;
                    }
                } else {
                    BadRequestError badRequestError = {body: {code: 90912, message: "Atleast one operation need to specified"}};
                    return badRequestError;
                }
                if endpoint is model:Endpoint {
                    _ = check self.setHttpRoute(apiArtifact, payload, endpoint, uniqueId, PRODUCTION_TYPE, organization);
                } else {
                    record {}? endpointConfig = payload.endpointConfig;
                    map<model:Endpoint|()> createdEndpoints = {};
                    if endpointConfig is record {} {
                        createdEndpoints = check self.createAndAddBackendServics(apiArtifact, payload, endpointConfig, (), (), organization);
                    }
                    _ = check self.setHttpRoute(apiArtifact, payload, createdEndpoints.hasKey(PRODUCTION_TYPE) ? createdEndpoints.get(PRODUCTION_TYPE) : (), uniqueId, PRODUCTION_TYPE, organization);
                    _ = check self.setHttpRoute(apiArtifact, payload, createdEndpoints.hasKey(SANDBOX_TYPE) ? createdEndpoints.get(SANDBOX_TYPE) : (), uniqueId, SANDBOX_TYPE, organization);
                }
                if definition is string {
                    self.retrieveGeneratedConfigmapForDefinition(apiArtifact, payload, definition, uniqueId);
                } else {
                    json internalDefinition = check self.getDefinition(check getAPI(apiId, organization));
                    json generatedSwagger = check self.retrieveGeneratedSwaggerDefinition(payload, internalDefinition.toJsonString());
                    self.retrieveGeneratedConfigmapForDefinition(apiArtifact, payload, generatedSwagger, uniqueId);
                }
                self.generateAndSetAPICRArtifact(apiArtifact, payload, organization);
                self.generateAndSetRuntimeAPIArtifact(apiArtifact, payload, (), organization);
                model:API deployAPIToK8sResult = check self.deployAPIToK8s(apiArtifact);
                return check convertK8sAPItoAPI(deployAPIToK8sResult, false);
            } on fail var e {
                if e is commons:APKError {
                    return e;
                }
                log:printError("Internal Error occured", e);
                return error("Internal Error occured", code = 909000, message = "Internal Error occured", description = "Internal Error occured", statusCode = 500);
            }
        } else {
            return api;
        }
    }
    # Description
    #
    # + apiId - Parameter Description  
    # + payload - Parameter Description  
    # + organization - Parameter Description
    # + return - Return Value Description
    public isolated function updateAPIDefinition(string apiId, http:Request payload, commons:Organization organization) returns http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError|BadRequestError|commons:APKError {
        do {
            API|NotFoundError api = check self.getAPIById(apiId, organization);
            if api is API {
                API updateAPI = {...api};
                DefinitionValidationRequest|BadRequestError definitionValidationRequest = check self.mapApiDefinitionPayload(payload);
                if definitionValidationRequest is DefinitionValidationRequest {
                    runtimeapi:APIDefinitionValidationResponse|runtimeapi:APIManagementException|BadRequestError validateAndRetrieveDefinitionResult = check self.validateAndRetrieveDefinition(updateAPI.'type, definitionValidationRequest.url, definitionValidationRequest.inlineAPIDefinition, definitionValidationRequest.content, definitionValidationRequest.fileName);
                    if validateAndRetrieveDefinitionResult is runtimeapi:APIDefinitionValidationResponse {
                        if validateAndRetrieveDefinitionResult.isValid() {
                            runtimeapi:APIDefinition parser = validateAndRetrieveDefinitionResult.getParser();
                            log:printDebug("content available ==", contentAvailable = (validateAndRetrieveDefinitionResult.getContent() is string));
                            utilapis:Set|runtimeapi:APIManagementException uRITemplates = parser.getURITemplates(<string>validateAndRetrieveDefinitionResult.getContent());
                            if uRITemplates is utilapis:Set {
                                APIOperations[]? operations = updateAPI.operations;
                                if !(operations is APIOperations[]) {
                                    operations = [];
                                }
                                lang:Object[] uriTemplates = check uRITemplates.toArray();
                                APIOperations[] sortedOperations = [];
                                foreach lang:Object uritemplate in uriTemplates {
                                    runtimeModels:URITemplate template = check java:cast(uritemplate);
                                    boolean found = false;
                                    foreach APIOperations operation in operations {
                                        if operation.target == template.getUriTemplate() && operation.verb == template.getHTTPVerb() {
                                            sortedOperations.push(operation);
                                            found = true;
                                            break;
                                        }
                                    }
                                    if !found {
                                        string[] scopes = [];
                                        utilapis:List scopeSet = template.getScopes();
                                        lang:Object[] scopeArray = check scopeSet.toArray();
                                        foreach lang:Object scope in scopeArray {
                                            scopes.push(scope.toString());
                                        }
                                        sortedOperations.push({
                                            target: template.getUriTemplate(),
                                            verb: template.getHTTPVerb(),
                                            scopes: scopes
                                        });
                                    }

                                }
                                updateAPI.operations = sortedOperations;
                                _ = check self.updateAPI(apiId, updateAPI, validateAndRetrieveDefinitionResult.getContent(), organization);
                                return self.getAPIDefinitionByID(apiId, organization);
                            } else {
                                log:printError("Error occured retrieving uri templates from definition", uRITemplates);
                                runtimeapi:JAPIManagementException excetion = check uRITemplates.ensureType(runtimeapi:JAPIManagementException);
                                runtimeapi:ErrorHandler errorHandler = excetion.getErrorHandler();
                                BadRequestError badeRequest = {body: {code: errorHandler.getErrorCode(), message: errorHandler.getErrorMessage().toString()}};
                                return badeRequest;
                            }
                        }
                        // Error definition.
                        ErrorListItem[] errorItems = [];
                        utilapis:ArrayList errorItemsResult = validateAndRetrieveDefinitionResult.getErrorItems();
                        foreach int i in 0 ... errorItemsResult.size() - 1 {
                            runtimeapi:ErrorItem errorItem = check java:cast(errorItemsResult.get(i));
                            ErrorListItem errorListItem = {code: errorItem.getErrorCode().toString(), message: <string>errorItem.getErrorMessage(), description: errorItem.getErrorDescription()};
                            errorItems.push(errorListItem);
                        }
                        BadRequestError badRequest = {body: {code: 90091, message: "Invalid API Definition", 'error: errorItems}};
                        return badRequest;
                    } else if validateAndRetrieveDefinitionResult is BadRequestError {
                        return validateAndRetrieveDefinitionResult;
                    } else {
                        log:printError("Error occured creating api from defintion", validateAndRetrieveDefinitionResult);
                        runtimeapi:JAPIManagementException excetion = check validateAndRetrieveDefinitionResult.ensureType(runtimeapi:JAPIManagementException);
                        runtimeapi:ErrorHandler errorHandler = excetion.getErrorHandler();
                        BadRequestError badeRequest = {body: {code: errorHandler.getErrorCode(), message: errorHandler.getErrorMessage().toString()}};
                        return badeRequest;
                    }
                } else {
                    return <BadRequestError>definitionValidationRequest;
                }
            } else {
                return <NotFoundError>api;
            }
        } on fail var e {
            log:printError("Error occured importing API", e);
            InternalServerErrorError internalError = {body: {code: 900900, message: "Internal Error."}};
            return internalError;
        }
    }

    private isolated function filterMediationPoliciesBasedOnQuery(MediationPolicy[] mediationPolicyList, string query, int 'limit, int offset, string sortBy, string sortOrder) returns MediationPolicyList|BadRequestError {
        MediationPolicy[] filteredList = [];
        if query.length() > 0 {
            int? semiCollonIndex = string:indexOf(query, ":", 0);
            if semiCollonIndex is int && semiCollonIndex > 0 {
                string keyWord = query.substring(0, semiCollonIndex);
                string keyWordValue = query.substring(keyWord.length() + 1, query.length());
                keyWordValue = keyWordValue + "|\\w+" + keyWordValue + "\\w+" + "|" + keyWordValue + "\\w+" + "|\\w+" + keyWordValue;
                if keyWord.trim() == SEARCH_CRITERIA_NAME {
                    foreach MediationPolicy mediationPolicy in mediationPolicyList {
                        if (regex:matches(mediationPolicy.name, keyWordValue)) {
                            filteredList.push(mediationPolicy);
                        }
                    }
                } else if keyWord.trim() == SEARCH_CRITERIA_TYPE {
                    foreach MediationPolicy mediationPolicy in mediationPolicyList {
                        if (regex:matches(mediationPolicy.'type, keyWordValue)) {
                            filteredList.push(mediationPolicy);
                        }
                    }
                } else {
                    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord " + keyWord}};
                    return badRequest;
                }
            } else {
                string keyWordValue = query + "|\\w+" + query + "\\w+" + "|" + query + "\\w+" + "|\\w+" + query;

                foreach MediationPolicy mediationPolicy in mediationPolicyList {

                    if (regex:matches(mediationPolicy.name, keyWordValue)) {
                        filteredList.push(mediationPolicy);
                    }
                }
            }
        } else {
            filteredList = mediationPolicyList;
        }
        return self.filterMediationPolicies(filteredList, 'limit, offset, sortBy, sortOrder);
    }

    private isolated function filterMediationPolicies(MediationPolicy[] mediationPolicyList, int 'limit, int offset, string sortBy, string sortOrder) returns MediationPolicyList|BadRequestError {
        MediationPolicy[] clonedMediationPolicyList = mediationPolicyList.clone();
        MediationPolicy[] sortedMediationPolicies = [];
        if sortBy == SORT_BY_POLICY_NAME && sortOrder == SORT_ORDER_ASC {
            sortedMediationPolicies = from var mediationPolicy in clonedMediationPolicyList
                order by mediationPolicy.name ascending
                select mediationPolicy;
        } else if sortBy == SORT_BY_POLICY_NAME && sortOrder == SORT_ORDER_DESC {
            sortedMediationPolicies = from var mediationPolicy in clonedMediationPolicyList
                order by mediationPolicy.name descending
                select mediationPolicy;
        } else if sortBy == SORT_BY_ID && sortOrder == SORT_ORDER_ASC {
            sortedMediationPolicies = from var mediationPolicy in clonedMediationPolicyList
                order by mediationPolicy.id ascending
                select mediationPolicy;
        } else if sortBy == SORT_BY_ID && sortOrder == SORT_ORDER_DESC {
            sortedMediationPolicies = from var mediationPolicy in clonedMediationPolicyList
                order by mediationPolicy.id descending
                select mediationPolicy;
        } else {
            BadRequestError badRequest = {body: {code: 90912, message: "Invalid Sort By/Sort Order Value "}};
            return badRequest;
        }
        MediationPolicy[] limitSet = [];
        if sortedMediationPolicies.length() > offset {
            foreach int i in offset ... (sortedMediationPolicies.length() - 1) {
                if limitSet.length() < 'limit {
                    limitSet.push(sortedMediationPolicies[i]);
                }
            }
        }
        return {list: limitSet, count: limitSet.length(), pagination: {total: mediationPolicyList.length(), 'limit: 'limit, offset: offset}};

    }

    # This return a Mediation policy by id.
    #
    # + id - Policy Id
    # + organization - Organization
    # + return - Return a Mediation Policy.
    public isolated function getMediationPolicyById(string id, commons:Organization organization) returns MediationPolicy|NotFoundError|commons:APKError {
        boolean mediationPolicyIDAvailable = id.length() > 0 ? true : false;
        if (mediationPolicyIDAvailable && string:length(id.toString()) > 0)
        {
            lock {
                foreach model:MediationPolicy mediationPolicy in getAvailableMediaionPolicies(organization) {
                    if mediationPolicy.id == id {
                        MediationPolicy matchedPolicy = convertPolicyModeltoPolicy(mediationPolicy);
                        return matchedPolicy.cloneReadOnly();
                    }
                }
            } on fail var e {
                return error("Error while retrieving Mediation policy", e, message = "Error while retrieving Mediation policy", description = "Error while retrieving Mediation policy", code = 909000, statusCode = 500);
            }
        }
        NotFoundError notfound = {body: {code: 909100, message: id + " not found."}};
        return notfound;
    }

    # This returns list of Mediation Policies.
    #
    # + query - Search Query
    # + 'limit - Limit 
    # + offset - Offset 
    # + sortBy - SortBy Parameter
    # + sortOrder - SortOrder  Parameter
    # + organization - Organization
    # + return - Return list of Mediation Policies.
    public isolated function getMediationPolicyList(string? query, int 'limit, int offset, string sortBy, string sortOrder, commons:Organization organization) returns MediationPolicyList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        MediationPolicy[] mediationPolicyList = [];
        foreach model:MediationPolicy mediationPolicy in getAvailableMediaionPolicies(organization) {
            MediationPolicy policyItem = convertPolicyModeltoPolicy(mediationPolicy);
            mediationPolicyList.push(policyItem);
        } on fail var e {
            return error("Error occured while getting Mediation policy list", e, message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
        if query is string && query.toString().trim().length() > 0 {
            return self.filterMediationPoliciesBasedOnQuery(mediationPolicyList, query, 'limit, offset, sortBy, sortOrder);
        } else {
            return self.filterMediationPolicies(mediationPolicyList, 'limit, offset, sortBy, sortOrder);
        }
    }

}

type ImportDefintionRequest record {
    string? url;
    string? fileName;
    byte[]? content;
    string? inlineAPIDefinition;
    API additionalPropertes;
    string 'type;
};

type DefinitionValidationRequest record {|
    string? url;
    string? fileName;
    byte[]? content;
    string? inlineAPIDefinition;
    string 'type;

|};

public isolated function getBackendServiceUid(API api, APIOperations? apiOperation, string endpointType, commons:Organization organization) returns string {
    string concatanatedString = uuid:createType1AsString();
    if (apiOperation is APIOperations) {
        return "backend-" + concatanatedString + "-resource";
    } else {
        concatanatedString = string:'join("-", organization.uuid, api.name, 'api.'version, endpointType);
        byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
        concatanatedString = hashedValue.toBase16();
        return "backend-" + concatanatedString + "-api";
    }
}

public isolated function getBackendPolicyUid(API api, string endpointType, commons:Organization organization) returns string {
    string concatanatedString = uuid:createType1AsString();
    return "backendpolicy-" + concatanatedString;
}

public isolated function getUniqueIdForAPI(string name, string 'version, commons:Organization organization) returns string {
    string concatanatedString = string:'join("-", organization.uuid, name, 'version);
    byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
    return hashedValue.toBase16();
}

public isolated function retrieveHttpRouteRefName(API api, string 'type, commons:Organization organization) returns string {
    return uuid:createType1AsString();
}
