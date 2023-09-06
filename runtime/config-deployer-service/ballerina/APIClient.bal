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

import ballerina/log;
import config_deployer_service.model;
import config_deployer_service.org.wso2.apk.config.model as runtimeModels;
import ballerina/regex;
import config_deployer_service.org.wso2.apk.config as runtimeUtil;
import ballerina/lang.value;
import config_deployer_service.org.wso2.apk.config.api as runtimeapi;
import ballerina/uuid;
import ballerina/crypto;
import config_deployer_service.java.io as javaio;
import wso2/apk_common_lib as commons;

#
public class APIClient {

    # This function used to convert APKInternalAPI model to APKConf.
    #
    # + api - APKInternalAPI model
    # + return - APKConf model.
    public isolated function fromAPIModelToAPKConf(runtimeModels:API api) returns APKConf|error {
        string generatedBasePath = api.getName() +  api.getVersion();
        byte[] data = generatedBasePath.toBytes();
        string encodedString = "/" + data.toBase64();
        if (encodedString.endsWith("==")) {
            encodedString = encodedString.substring(0,encodedString.length()-2);
        }
        APKConf apkConf = {
            name: api.getName(),
            basePath: api.getBasePath().length() > 0 ? api.getBasePath() : encodedString,
            version: api.getVersion()
        };
        string endpoint = api.getEndpoint();
        if endpoint.length() > 0 {
            apkConf.endpointConfigurations = {production: {endpoint: endpoint}};
        }

        runtimeModels:URITemplate[]|error uriTemplates = api.getUriTemplates();
        APKOperations[] operations = [];
        if uriTemplates is runtimeModels:URITemplate[] {
            foreach runtimeModels:URITemplate uriTemplate in uriTemplates {
                APKOperations operation = {
                    verb: uriTemplate.getHTTPVerb(),
                    target: uriTemplate.getUriTemplate(),
                    authTypeEnabled: uriTemplate.isAuthEnabled(),
                    scopes: check uriTemplate.getScopes()
                };
                string resourceEndpoint = uriTemplate.getEndpoint();
                if resourceEndpoint.length() > 0 {
                    operation.endpointConfigurations = {production: {endpoint: resourceEndpoint}};
                }
                operations.push(operation);
            }
        }
        apkConf.operations = operations;
        return apkConf;
    }

    public isolated function generateK8sArtifacts(APKConf apkConf, string? definition, commons:Organization organization) returns model:APIArtifact|commons:APKError {
        do {

            string uniqueId = self.getUniqueIdForAPI(apkConf.name, apkConf.version, organization);
            if apkConf.id is string {
                uniqueId = <string>apkConf.id;
            }
            model:APIArtifact apiArtifact = {uniqueId: uniqueId, name: apkConf.name, version: apkConf.version, organization: organization.name};
            EndpointConfigurations[] resourceLevelEndpointConfigList;
            APKOperations[]? operations = apkConf.operations;
            if operations is APKOperations[] {
                if operations.length() == 0 {
                    return e909021();
                }

                // Validating rate limit.
                _ = check self.validateRateLimit(apkConf.apiRateLimit, operations);
                resourceLevelEndpointConfigList = self.getResourceLevelEndpointConfig(operations);
            } else {
                return e909021();
            }
            map<model:Endpoint|()> createdEndpoints = {};
            EndpointConfigurations? endpointConfigurations = apkConf.endpointConfigurations;
            if endpointConfigurations is EndpointConfigurations {
                createdEndpoints = check self.createAndAddBackendServics(apiArtifact, apkConf, endpointConfigurations, (), (), organization);
            }
            AuthenticationRequest[]? authentication = apkConf.authentication;
            if authentication is AuthenticationRequest[] {
                if createdEndpoints != {} {
                _ = check self.populateAuthenticationMap(apiArtifact, apkConf, authentication, createdEndpoints, organization);
                } else {
                    // check if there are resource level endpoints
                    if resourceLevelEndpointConfigList.length() > 0 {
                        foreach EndpointConfigurations resourceEndpointConfigurations in resourceLevelEndpointConfigList {
                            map<model:Endpoint> resourceEndpointIdMap = {};
                            EndpointConfiguration? productionEndpointConfig = resourceEndpointConfigurations.production;
                            EndpointConfiguration? sandboxEndpointConfig = resourceEndpointConfigurations.sandbox;
                            if sandboxEndpointConfig is EndpointConfiguration {
                                resourceEndpointIdMap[SANDBOX_TYPE] = {
                                    name: "",
                                    serviceEntry: false,
                                    url: self.construcURlFromService(sandboxEndpointConfig.endpoint)
                                };
                            }
                            if productionEndpointConfig is EndpointConfiguration {
                                resourceEndpointIdMap[PRODUCTION_TYPE] = {
                                    name: "",
                                    serviceEntry: false,
                                    url: self.construcURlFromService(productionEndpointConfig.endpoint)
                                };
                            }
                            _ = check self.populateAuthenticationMap(apiArtifact, apkConf, authentication, resourceEndpointIdMap, organization);
                        }
                    } else {
                        _ = check self.populateAuthenticationMap(apiArtifact, apkConf, authentication, createdEndpoints, organization);
                    }
                }
            }
            _ = check self.setHttpRoute(apiArtifact, apkConf, createdEndpoints.hasKey(PRODUCTION_TYPE) ? createdEndpoints.get(PRODUCTION_TYPE) : (), uniqueId, PRODUCTION_TYPE, organization);
            _ = check self.setHttpRoute(apiArtifact, apkConf, createdEndpoints.hasKey(SANDBOX_TYPE) ? createdEndpoints.get(SANDBOX_TYPE) : (), uniqueId, SANDBOX_TYPE, organization);

            json generatedSwagger = check self.retrieveGeneratedSwaggerDefinition(apkConf, definition);
            check self.retrieveGeneratedConfigmapForDefinition(apiArtifact, apkConf, generatedSwagger, uniqueId, organization);
            self.generateAndSetAPICRArtifact(apiArtifact, apkConf, organization);
            _ = check self.generateAndSetPolicyCRArtifact(apiArtifact, apkConf, organization);
            apiArtifact.organization = organization.name;
            return apiArtifact;
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            return e909022("Internal Error occured while generating k8s-artifact", e);
        }
    }

    isolated function getResourceLevelEndpointConfig(APKOperations[] operations) returns EndpointConfigurations[] {
        EndpointConfigurations[] endpointConfigurationsList = [];
        foreach APKOperations operation in operations {
            EndpointConfigurations? endpointConfigurations = operation.endpointConfigurations;
            if (endpointConfigurations != ()) {
                // Presence of resource level Endpoint Configuration.
                endpointConfigurationsList.push(endpointConfigurations);
            }
        }
        return endpointConfigurationsList;
    }

    private isolated function getHostNames(APKConf apkConf, string uniqueId, string endpointType, commons:Organization organization) returns string[] {
        //todo: need to implement vhost feature
        Vhost[] globalVhosts = vhosts;
        string[] hosts = [];
        foreach Vhost vhost in globalVhosts {
            if vhost.'type == endpointType {
                foreach string host in vhost.hosts {
                    hosts.push(string:concat(organization.name, ".", host));
                }
            }
        }
        return hosts;
    }
    isolated function isPolicyEmpty(APIOperationPolicies? policies) returns boolean {
        if policies is APIOperationPolicies {
            APKOperationPolicy[]? request = policies.request;
            if request is APKOperationPolicy[] {
                if (request.length() > 0) {
                    return false;
                }
            }
            APKOperationPolicy[]? response = policies.response;
            if response is APKOperationPolicy[] {
                if (response.length() > 0) {
                    return false;
                }
            }
        }
        return true;
    }

    isolated function validateRateLimit(RateLimit? apiRateLimit, APKOperations[] operations) returns commons:APKError|() {
        if (apiRateLimit == ()) {
            return ();
        } else {
            foreach APKOperations operation in operations {
                RateLimit? operationRateLimit = operation.operationRateLimit;
                if (operationRateLimit != ()) {
                    // Presence of both resource level and API level rate limits.
                    return e909026();
                }
            }
        }
        return ();
    }

    private isolated function createAndAddBackendServics(model:APIArtifact apiArtifact, APKConf apkConf, EndpointConfigurations endpointConfigurations, APKOperations? apiOperation, string? endpointType, commons:Organization organization) returns map<model:Endpoint>|commons:APKError|error {
        map<model:Endpoint> endpointIdMap = {};
        EndpointConfiguration? productionEndpointConfig = endpointConfigurations.production;
        EndpointConfiguration? sandboxEndpointConfig = endpointConfigurations.sandbox;
        if endpointType == () || (endpointType == SANDBOX_TYPE) {
            if sandboxEndpointConfig is EndpointConfiguration {
                model:Backend backendService = check self.createBackendService(apiArtifact, apkConf, apiOperation, SANDBOX_TYPE, organization, sandboxEndpointConfig);
                if apiOperation == () {
                    apiArtifact.sandboxEndpointAvailable = true;
                }
                apiArtifact.backendServices[backendService.metadata.name] = (backendService);
                endpointIdMap[SANDBOX_TYPE] = {
                    name: backendService.metadata.name,
                    serviceEntry: false,
                    url: self.construcURlFromService(sandboxEndpointConfig.endpoint)
                };
            }
        }
        if (endpointType == () || endpointType == PRODUCTION_TYPE) {
            if productionEndpointConfig is EndpointConfiguration {
                model:Backend backendService = check self.createBackendService(apiArtifact, apkConf, apiOperation, PRODUCTION_TYPE, organization, productionEndpointConfig);
                if apiOperation == () {
                    apiArtifact.productionEndpointAvailable = true;
                }
                apiArtifact.backendServices[backendService.metadata.name] = (backendService);
                endpointIdMap[PRODUCTION_TYPE] = {
                    name: backendService.metadata.name,
                    serviceEntry: false,
                    url: self.construcURlFromService(productionEndpointConfig.endpoint)
                };
            }
        }
        return endpointIdMap;
    }

    isolated function construcURlFromService(string|K8sService endpoint) returns string {
        if endpoint is string {
            return endpoint;
        } else {
            return self.construcURlFromK8sService(endpoint);
        }
    }

    isolated function getLabels(APKConf api, commons:Organization organization) returns map<string> {
        string apiNameHash = crypto:hashSha1(api.name.toBytes()).toBase16();
        string apiVersionHash = crypto:hashSha1(api.'version.toBytes()).toBase16();
        string organizationHash = crypto:hashSha1(organization.name.toBytes()).toBase16();
        map<string> labels = {
            [API_NAME_HASH_LABEL] : apiNameHash,
            [API_VERSION_HASH_LABEL] : apiVersionHash,
            [ORGANIZATION_HASH_LABEL] : organizationHash,
            [MANAGED_BY_HASH_LABEL] : MANAGED_BY_HASH_LABEL_VALUE
        };
        return labels;
    }

    isolated function returnFullBasePath(string basePath, string 'version) returns string {
        string fullBasePath = basePath;
        if (!string:endsWith(basePath, 'version)) {
            fullBasePath = string:'join("/", basePath, 'version);
        }
        return fullBasePath;
    }

    private isolated function construcURlFromK8sService(K8sService 'k8sService) returns string {
        return <string>k8sService.protocol + "://" + string:'join(".", <string>k8sService.name, <string>k8sService.namespace, "svc.cluster.local") + ":" + k8sService.port.toString();
    }

    private isolated function constructURLFromBackendSpec(model:BackendSpec backendSpec) returns string {
        return backendSpec.protocol + "://" + backendSpec.services[0].host + backendSpec.services[0].port.toString();
    }

    private isolated function retrieveGeneratedConfigmapForDefinition(model:APIArtifact apiArtifact, APKConf apkConf, json generatedSwaggerDefinition, string uniqueId, commons:Organization organization) returns error? {
        byte[]|javaio:IOException compressedContent = check commons:GzipUtil_compressGzipFile(generatedSwaggerDefinition.toJsonString().toBytes());
        if compressedContent is byte[] {

            byte[] base64EncodedContent = check commons:EncoderUtil_encodeBase64(compressedContent);
            model:ConfigMap configMap = {
                metadata: {
                    name: self.retrieveDefinitionName(uniqueId),
                    labels: self.getLabels(apkConf, organization)
                }
            };
            configMap.binaryData = {[CONFIGMAP_DEFINITION_KEY] : check string:fromBytes(base64EncodedContent)};
            apiArtifact.definition = configMap;
        } else {
            return compressedContent.cause();
        }
    }

    isolated function setDefaultOperationsIfNotExist(APKConf api) {
        APKOperations[]? operations = api.operations;
        boolean operationsAvailable = false;
        if operations is APKOperations[] {
            operationsAvailable = operations.length() > 0;
        }
        if operationsAvailable == false {
            APKOperations[] apiOperations = [];
            if api.'type == API_TYPE_REST {
                foreach string httpverb in HTTP_DEFAULT_METHODS {
                    APKOperations apiOperation = {target: "/*", verb: httpverb.toUpperAscii()};
                    apiOperations.push(apiOperation);
                }
                api.operations = apiOperations;
            }
        }
    }

    private isolated function generateAndSetPolicyCRArtifact(model:APIArtifact apiArtifact, APKConf apkConf, commons:Organization organization) returns error? {
        if apkConf.apiRateLimit != () {
            model:RateLimitPolicy? rateLimitPolicyCR = self.generateRateLimitPolicyCR(apkConf, apkConf.apiRateLimit, apiArtifact.uniqueId, (), organization);
            if rateLimitPolicyCR != () {
                apiArtifact.rateLimitPolicies[rateLimitPolicyCR.metadata.name] = rateLimitPolicyCR;
            }
        }
        if apkConf.apiPolicies != () || apkConf.corsConfiguration != () {
            model:APIPolicy? apiPolicyCR = check self.generateAPIPolicyAndBackendCR(apiArtifact, apkConf, (), apkConf.apiPolicies, organization, apiArtifact.uniqueId);
            if apiPolicyCR != () {
                apiArtifact.apiPolicies[apiPolicyCR.metadata.name] = apiPolicyCR;
            }
        }
    }

    private isolated function populateAuthenticationMap(model:APIArtifact apiArtifact, APKConf apkConf, AuthenticationRequest[] authentications,
            map<model:Endpoint|()> createdEndpointMap, commons:Organization organization) returns error? {
        map<model:Authentication> authenticationMap = {};
        model:AuthenticationExtenstionType authTypes = {};
        foreach AuthenticationRequest authentication in authentications {
            if authentication.authType == "JWT" {
                JWTAuthentication jwtAuthentication = check authentication.cloneWithType(JWTAuthentication);
                authTypes.jwt = {header: <string>jwtAuthentication.headerName, sendTokenToUpstream: <boolean>jwtAuthentication.sendTokenToUpstream, disabled: !jwtAuthentication.enabled};
            } else if authentication.authType == "APIKey" && authentication is APIKeyAuthentication {
                APIKeyAuthentication apiKeyAuthentication = check authentication.cloneWithType(APIKeyAuthentication);
                authTypes.apiKey = [];
                authTypes.apiKey.push({'in: "Header", name: apiKeyAuthentication.headerName, sendTokenToUpstream: apiKeyAuthentication.sendTokenToUpstream});
                authTypes.apiKey.push({'in: "Query", name: apiKeyAuthentication.queryParamName, sendTokenToUpstream: apiKeyAuthentication.sendTokenToUpstream});
            }
        }
        log:printDebug("Auth Types:" + authTypes.toString());
        string[] keys = createdEndpointMap.keys();
        log:printDebug("createdEndpointMap.keys:" + createdEndpointMap.keys().toString());
        foreach string endpointType in keys {
            string disableAuthenticationRefName = self.retrieveDisableAuthenticationRefName(apkConf, endpointType, organization);
            log:printDebug("disableAuthenticationRefName:" + disableAuthenticationRefName);
            model:Authentication authentication = {
                metadata: {
                    name: disableAuthenticationRefName,
                    labels: self.getLabels(apkConf, organization)
                },
                spec: {
                    override: {
                        disabled: false,
                        authTypes: authTypes
                    },
                    targetRef: {
                        group: "gateway.networking.k8s.io",
                        kind: "Resource",
                        name: apiArtifact.uniqueId
                    }
                }
            };
            log:printDebug("Authentication CR:" + authentication.toString());
            authenticationMap[disableAuthenticationRefName] = authentication;
        }
        log:printDebug("Authentication Map:" + authenticationMap.toString());
        apiArtifact.authenticationMap = authenticationMap;
    }

    private isolated function generateAndSetAPICRArtifact(model:APIArtifact apiArtifact, APKConf apkConf, commons:Organization organization) {
        model:API k8sAPI = {
            metadata: {
                name: apiArtifact.uniqueId,
                labels: self.getLabels(apkConf, organization)
            },
            spec: {
                apiDisplayName: apkConf.name,
                apiType: apkConf.'type,
                apiVersion: apkConf.'version,
                basePath: self.returnFullBasePath(apkConf.basePath, apkConf.'version),
                isDefaultVersion: apkConf.defaultVersion,
                organization: organization.name,
                definitionPath: apkConf.definitionPath
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
            k8sAPI.spec.production = [{httpRouteRefs: productionHttpRoutes}];
        }
        if sandBoxHttpRoutes.length() > 0 {
            k8sAPI.spec.sandbox = [{httpRouteRefs: sandBoxHttpRoutes}];
        }
        if apkConf.id != () {
            k8sAPI.metadata["annotations"] = {[API_UUID_ANNOTATION] : <string>apkConf.id};
        }
        if apkConf.additionalProperties is APKConf_additionalProperties[] {
            model:APIProperties[] properties = [];
            foreach APKConf_additionalProperties additionalProperty in <APKConf_additionalProperties[]>apkConf.additionalProperties {
                properties.push({name: <string>additionalProperty.name, value: <string>additionalProperty.value});
            }
            k8sAPI.spec.apiProperties = properties;
        }
        apiArtifact.api = k8sAPI;
    }

    isolated function retrieveDefinitionName(string uniqueId) returns string {
        return uniqueId + "-definition";
    }

    private isolated function retrieveDisableAuthenticationRefName(APKConf apkConf, string 'type, commons:Organization organization) returns string {
        return self.getUniqueIdForAPI(apkConf.name, apkConf.'version, organization) + "-" + 'type + "-authentication";
    }

    private isolated function setHttpRoute(model:APIArtifact apiArtifact, APKConf apkConf, model:Endpoint? endpoint, string uniqueId, string endpointType, commons:Organization organization) returns commons:APKError|error? {
        APKOperations[] apiOperations = apkConf.operations ?: [];
        APKOperations[][] operationsArray = [];
        int row = 0;
        int column = 0;
        foreach APKOperations item in apiOperations {
            if column > 7 {
                row = row + 1;
                column = 0;
            }
            operationsArray[row][column] = item;
            column = column + 1;
        }
        foreach APKOperations[] item in operationsArray {
            APKConf clonedAPKConf = apkConf.clone();
            clonedAPKConf.operations = item.clone();
            _ = check self.putHttpRouteForPartition(apiArtifact, clonedAPKConf, endpoint, uniqueId, endpointType, organization);
        }
    }

    private isolated function putHttpRouteForPartition(model:APIArtifact apiArtifact, APKConf apkConf, model:Endpoint? endpoint, string uniqueId, string endpointType, commons:Organization organization) returns commons:APKError|error? {
        string httpRouteRefName = self.retrieveHttpRouteRefName(apkConf, endpointType, organization);
        model:Httproute httpRoute = {
            metadata:
                {
                name: httpRouteRefName,
                labels: self.getLabels(apkConf, organization)
            },
            spec: {
                parentRefs: self.generateAndRetrieveParentRefs(apkConf, uniqueId),
                rules: check self.generateHttpRouteRules(apiArtifact, apkConf, endpoint, endpointType, organization, httpRouteRefName),
                hostnames: self.getHostNames(apkConf, uniqueId, endpointType, organization)
            }
        };
        if httpRoute.spec.rules.length() > 0 {
            if endpointType == PRODUCTION_TYPE {
                apiArtifact.productionRoute.push(httpRoute);
            } else {
                apiArtifact.sandboxRoute.push(httpRoute);
            }
        }

        return;
    }

    private isolated function generateAndRetrieveParentRefs(APKConf apkConf, string uniqueId) returns model:ParentReference[] {
        string gatewayName = gatewayConfiguration.name;
        string listenerName = gatewayConfiguration.listenerName;
        model:ParentReference[] parentRefs = [];
        model:ParentReference parentRef = {group: "gateway.networking.k8s.io", kind: "Gateway", name: gatewayName, sectionName: listenerName};
        parentRefs.push(parentRef);
        return parentRefs;
    }

    private isolated function generateHttpRouteRules(model:APIArtifact apiArtifact, APKConf apkConf, model:Endpoint? endpoint, string endpointType, commons:Organization organization, string httpRouteRefName) returns model:HTTPRouteRule[]|commons:APKError|error {
        model:HTTPRouteRule[] httpRouteRules = [];
        APKOperations[]? operations = apkConf.operations;
        if operations is APKOperations[] {
            foreach APKOperations operation in operations {
                model:HTTPRouteRule|() httpRouteRule = check self.generateHttpRouteRule(apiArtifact, apkConf, endpoint, operation, endpointType, organization);
                if httpRouteRule is model:HTTPRouteRule {
                    model:HTTPRouteFilter[]? filters = httpRouteRule.filters;
                    if filters is () {
                        filters = [];
                        httpRouteRule.filters = filters;
                    }
                    string disableAuthenticationRefName = self.retrieveDisableAuthenticationRefName(apkConf, endpointType, organization);
                    if !(operation.authTypeEnabled ?: true) {
                        if !apiArtifact.authenticationMap.hasKey(disableAuthenticationRefName) {
                            model:Authentication generateDisableAuthenticationCR = self.generateDisableAuthenticationCR(apiArtifact, apkConf, endpointType, organization);
                            apiArtifact.authenticationMap[disableAuthenticationRefName] = generateDisableAuthenticationCR;
                        }
                        model:HTTPRouteFilter disableAuthenticationFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: "Authentication", name: disableAuthenticationRefName}};
                        (<model:HTTPRouteFilter[]>filters).push(disableAuthenticationFilter);
                    } else if apiArtifact.authenticationMap.hasKey(disableAuthenticationRefName) {
                        model:HTTPRouteFilter authenticationFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: "Authentication", name: disableAuthenticationRefName}};
                        (<model:HTTPRouteFilter[]>filters).push(authenticationFilter);
                    }
                    string[]? scopes = operation.scopes;
                    if scopes is string[] {
                        foreach string scope in scopes {
                            model:Scope scopeCr;
                            if apiArtifact.scopes.hasKey(scope) {
                                scopeCr = apiArtifact.scopes.get(scope);
                            } else {
                                scopeCr = self.generateScopeCR(apiArtifact, apkConf, organization, scope);
                            }
                            model:HTTPRouteFilter scopeFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: scopeCr.kind, name: scopeCr.metadata.name}};
                            (<model:HTTPRouteFilter[]>filters).push(scopeFilter);
                        }
                    }
                    if operation.operationRateLimit != () {
                        model:RateLimitPolicy? rateLimitPolicyCR = self.generateRateLimitPolicyCR(apkConf, operation.operationRateLimit, apiArtifact.uniqueId, operation, organization);
                        if rateLimitPolicyCR != () {
                            apiArtifact.rateLimitPolicies[rateLimitPolicyCR.metadata.name] = rateLimitPolicyCR;
                            model:HTTPRouteFilter rateLimitPolicyFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: "RateLimitPolicy", name: rateLimitPolicyCR.metadata.name}};
                            (<model:HTTPRouteFilter[]>filters).push(rateLimitPolicyFilter);
                        }
                    }
                    if operation.operationPolicies != () {
                        model:APIPolicy? apiPolicyCR = check self.generateAPIPolicyAndBackendCR(apiArtifact, apkConf, operation, operation.operationPolicies, organization, apiArtifact.uniqueId);
                        if apiPolicyCR != () {
                            apiArtifact.apiPolicies[apiPolicyCR.metadata.name] = apiPolicyCR;
                            model:HTTPRouteFilter apiPolicyFilter = {'type: "ExtensionRef", extensionRef: {group: "dp.wso2.com", kind: "APIPolicy", name: apiPolicyCR.metadata.name}};
                            (<model:HTTPRouteFilter[]>filters).push(apiPolicyFilter);
                        }
                    }
                    httpRouteRules.push(httpRouteRule);
                }
            }
        }
        return httpRouteRules;
    }

    private isolated function generateAPIPolicyAndBackendCR(model:APIArtifact apiArtifact, APKConf apkConf, APKOperations? operations, APIOperationPolicies? policies, commons:Organization organization, string targetRefName) returns model:APIPolicy?|error {
        model:APIPolicyData defaultSpecData = {};
        APKOperationPolicy[]? request = policies?.request;
        any[] requestPolicy = check self.retrieveAPIPolicyDetails(apiArtifact, apkConf, operations, organization, request, "request");
        foreach any item in requestPolicy {
            if item is model:InterceptorReference {
                defaultSpecData.requestInterceptors = [item];
            } else if item is model:BackendJwtReference {
                defaultSpecData.backendJwtPolicy = item;
            }
        }
        APKOperationPolicy[]? response = policies?.response;
        any[] responseInterceptor = check self.retrieveAPIPolicyDetails(apiArtifact, apkConf, operations, organization, response, "response");
        foreach any item in responseInterceptor {
            if item is model:InterceptorReference {
                defaultSpecData.responseInterceptors = [item];
            }
        }
        CORSConfiguration? corsConfiguration = apkConf.corsConfiguration;
        if corsConfiguration is CORSConfiguration {
            model:CORSPolicy? cORSPolicy = self.retrieveCORSPolicyDetails(apiArtifact, apkConf, corsConfiguration, organization);
            if cORSPolicy is model:CORSPolicy {
                defaultSpecData.cORSPolicy = cORSPolicy;
            }
        }
        if defaultSpecData != {} {
            model:APIPolicy? apiPolicyCR = self.generateAPIPolicyCR(apkConf, targetRefName, operations, organization, defaultSpecData);
            if apiPolicyCR != () {
                return apiPolicyCR;
            }
        }
        return ();
    }

    private isolated function generateScopeCR(model:APIArtifact apiArtifact, APKConf apkConf, commons:Organization organization, string scope) returns model:Scope {
        string scopeName = uuid:createType1AsString();
        model:Scope scopeCr = {
            metadata: {
                name: scopeName,
                labels: self.getLabels(apkConf, organization)
            },
            spec: {
                names: [scope]
            }
        };
        apiArtifact.scopes[scope] = scopeCr;
        return scopeCr;
    }

    private isolated function generateDisableAuthenticationCR(model:APIArtifact apiArtifact, APKConf apkConf, string endpointType, commons:Organization organization) returns model:Authentication {
        string retrieveDisableAuthenticationRefName = self.retrieveDisableAuthenticationRefName(apkConf, endpointType, organization);
        model:Authentication authentication = {
            metadata: {name: retrieveDisableAuthenticationRefName, labels: self.getLabels(apkConf, organization)},
            spec: {
                targetRef: {
                    group: "",
                    kind: "Resource",
                    name: self.getUniqueIdForAPI(apkConf.name, apkConf.'version, organization)
                },
                override: {
                    disabled: true
                }
            }
        };
        return authentication;
    }

    private isolated function generateHttpRouteRule(model:APIArtifact apiArtifact, APKConf apkConf, model:Endpoint? endpoint, APKOperations operation, string endpointType, commons:Organization organization) returns model:HTTPRouteRule|()|commons:APKError {
        do {
            EndpointConfigurations? endpointConfig = operation.endpointConfigurations;
            model:Endpoint? endpointToUse = ();
            if endpointConfig is EndpointConfigurations {
                // endpointConfig presense at Operation Level.
                map<model:Endpoint> operationalLevelBackend = check self.createAndAddBackendServics(apiArtifact, apkConf, endpointConfig, operation, endpointType, organization);
                if operationalLevelBackend.hasKey(endpointType) {
                    endpointToUse = operationalLevelBackend.get(endpointType);
                }
            } else {
                if endpoint is model:Endpoint {
                    endpointToUse = endpoint;
                }
            }
            if endpointToUse != () {
                model:HTTPRouteRule httpRouteRule = {matches: self.retrieveMatches(apkConf, operation, organization), backendRefs: self.retrieveGeneratedBackend(apkConf, endpointToUse, endpointType), filters: self.generateFilters(apiArtifact, apkConf, endpointToUse, operation, endpointType, organization)};
                return httpRouteRule;
            } else {
                return ();
            }
        } on fail var e {
            log:printError("Internal Error occured", e);
            return e909022("Internal Error occured", e);
        }
    }

    private isolated function generateFilters(model:APIArtifact apiArtifact, APKConf apkConf, model:Endpoint endpoint, APKOperations operation, string endpointType, commons:Organization organization) returns model:HTTPRouteFilter[] {
        model:HTTPRouteFilter[] routeFilters = [];
        string generatedPath = self.generatePrefixMatch(endpoint, operation);
        if (generatedPath != operation.target) {
            model:HTTPRouteFilter replacePathFilter = {'type: "URLRewrite", urlRewrite: {path: {'type: "ReplaceFullPath", replaceFullPath: generatedPath}}};
            routeFilters.push(replacePathFilter);
        }
        APIOperationPolicies? operationPoliciesToUse = ();
        if (apkConf.apiPolicies is APIOperationPolicies) {
            operationPoliciesToUse = apkConf.apiPolicies;
        } else {
            operationPoliciesToUse = operation.operationPolicies;
        }
        if operationPoliciesToUse is APIOperationPolicies {
            APKOperationPolicy[]? request = operationPoliciesToUse.request;
            if request is APKOperations[] {
            }
        }
        return routeFilters;
    }

    isolated function extractHttpHeaderFilterData(APKOperationPolicy[] operationPolicy, commons:Organization organization) returns model:HTTPHeaderFilter {
        model:HTTPHeader[] setPolicies = [];
        string[] removePolicies = [];
        foreach APKOperationPolicy policy in operationPolicy {
            string policyName = policy.policyName;

            record {}? policyParameters = policy.parameters;
            if (policyParameters is record {}) {
                if (policyName == "addHeader") {
                    model:HTTPHeader httpHeader = {
                        name: <string>policyParameters.get("headerName"),
                        value: <string>policyParameters.get("headerValue")
                    };
                    setPolicies.push(httpHeader);
                }
                if (policyName == "removeHeader") {
                    string httpHeader = <string>policyParameters.get("headerName");
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

    isolated function generatePrefixMatch(model:Endpoint endpoint, APKOperations operation) returns string {
        string target = operation.target ?: "/*";
        string[] splitValues = regex:split(target, "/");
        string generatedPath = "";
        int pathparamCount = 1;
        if (target == "/*") {
            generatedPath = "\\1";
        } else if (target == "/") {
            generatedPath = "/";
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
        return generatedPath;
    }

    public isolated function retrievePathPrefix(string basePath, string 'version, string operation, commons:Organization organization) returns string {
        string[] splitValues = regex:split(operation, "/");
        string generatedPath = "";
        if (operation == "/*") {
            return "(.*)";
        } else if operation == "/" {
            return "/";
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

    private isolated function retrieveGeneratedBackend(APKConf apkConf, model:Endpoint endpoint, string endpointType) returns model:HTTPBackendRef[] {
        model:HTTPBackendRef httpBackend = {
            kind: "Backend",
            name: <string>endpoint.name,
            group: "dp.wso2.com"
        };
        return [httpBackend];
    }

    private isolated function retrieveMatches(APKConf apkConf, APKOperations apiOperation, commons:Organization organization) returns model:HTTPRouteMatch[] {
        model:HTTPRouteMatch[] httpRouteMatch = [];
        model:HTTPRouteMatch httpRoute = self.retrieveHttpRouteMatch(apkConf, apiOperation, organization);

        httpRouteMatch.push(httpRoute);
        return httpRouteMatch;
    }

    private isolated function retrieveHttpRouteMatch(APKConf apkConf, APKOperations apiOperation, commons:Organization organization) returns model:HTTPRouteMatch {

        return {method: <string>apiOperation.verb, path: {'type: "RegularExpression", value: self.retrievePathPrefix(apkConf.basePath, apkConf.'version, apiOperation.target ?: "/*", organization)}};
    }

    isolated function retrieveGeneratedSwaggerDefinition(APKConf apkConf, string? definition) returns json|commons:APKError|error {
        runtimeModels:API api1 = runtimeModels:newAPI1();
        api1.setName(apkConf.name);
        api1.setType(apkConf.'type);
        api1.setVersion(apkConf.'version);
        runtimeModels:URITemplate[] uritemplatesSet = [];
        if apkConf.operations is APKOperations[] {
            foreach APKOperations apiOperation in <APKOperations[]>apkConf.operations {
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
                        uriTemplate.setScopes(item);
                    }
                }
                uritemplatesSet.push(uriTemplate);
            }
        }
        check api1.setUriTemplates(uritemplatesSet);
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
                return e909043();
            }
        } else if retrievedDefinition is () {
            return "";
        } else {
            return e909043();
        }
    }

    isolated function gethost(string|K8sService endpoint) returns string {
        string url;
        if endpoint is string {
            url = endpoint;
        } else {
            url = self.construcURlFromK8sService(endpoint);
        }
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

    isolated function getProtocol(string|K8sService endpoint) returns string {
        if endpoint is string {
            return endpoint.startsWith("https://") ? "https" : "http";
        } else {
            return endpoint.protocol ?: "http";
        }
    }

    isolated function getPort(string|K8sService endpoint) returns int|error {
        string url;
        if endpoint is string {
            url = endpoint;
        } else {
            url = self.construcURlFromK8sService(endpoint);
        }
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

    isolated function handleK8sTimeout(model:Status errorStatus) returns commons:APKError {
        model:StatusDetails? details = errorStatus.details;
        if details is model:StatusDetails {
            if details.retryAfterSeconds is int && details.retryAfterSeconds >= 0 {
                // K8s api level ratelimit hit.
                log:printError("K8s API Timeout happens when invoking k8s api");
            }
        }
        return e909022("Internal server error", e = error("Internal server error"));
    }

    isolated function createBackendService(model:APIArtifact apiArtifact, APKConf apkConf, APKOperations? apiOperation, string endpointType, commons:Organization organization, EndpointConfiguration endpointConfig) returns commons:APKError|model:Backend|error {
        model:SecurityConfig? securityConfig = ();
        EndpointSecurity? endpointSecurity = endpointConfig?.endpointSecurity;
        model:Backend backendService = {
            metadata: {
                name: self.getBackendServiceUid(apkConf, apiOperation, endpointType, organization),
                labels: self.getLabels(apkConf, organization)
            },
            spec: {
                services: [
                    {
                        host: self.gethost(endpointConfig.endpoint),
                        port: check self.getPort(endpointConfig.endpoint)
                    }
                ],
                basePath: endpointConfig.endpoint is string ? getPath(<string>endpointConfig.endpoint) : (),
                protocol: self.getProtocol(endpointConfig.endpoint)
            }
        };
        if endpointType == INTERCEPTOR_TYPE {
            backendService.metadata.name = self.getInterceptorBackendUid(apkConf, endpointType, organization, endpointConfig.endpoint);
        }
        if endpointSecurity is EndpointSecurity {
            if endpointSecurity?.enabled ?: false {
                // When user adds basic auth endpoint security username and password
                BasicEndpointSecurity? securityType = endpointSecurity.securityType;

                if securityType is BasicEndpointSecurity {
                    securityConfig = {
                        basic: {
                            secretRef: {
                                name: <string>securityType.secretName,
                                usernameKey: <string>securityType.userNameKey,
                                passwordKey: <string>securityType.passwordKey
                            }
                        }
                    };
                }
            }
            backendService.spec.security = securityConfig;
        }
        Certificate? endpointCertificate = endpointConfig.certificate;
        if endpointCertificate is Certificate && backendService.spec.protocol == "https" {
            backendService.spec.tls = {
                configMapRef: {
                    key: <string>endpointCertificate.secretKey,
                    name: <string>endpointCertificate.secretName
                }
            };
        }
        Resiliency? resiliency = endpointConfig.resiliency;
        if resiliency is Resiliency {
            backendService.spec.timeout = resiliency.timeout;
            backendService.spec.'retry = resiliency.retryPolicy;
            backendService.spec.circuitBreaker = resiliency.circuitBreaker;
        }
        return backendService;
    }

    public isolated function generateRateLimitPolicyCR(APKConf apkConf, RateLimit? rateLimit, string targetRefName, APKOperations? operation, commons:Organization organization) returns model:RateLimitPolicy? {
        model:RateLimitPolicy? rateLimitPolicyCR = ();
        if rateLimit != () {
            rateLimitPolicyCR = {
                metadata: {
                    name: self.retrieveRateLimitPolicyRefName(operation),
                    labels: self.getLabels(apkConf, organization)
                },
                spec: {
                    default: self.retrieveRateLimitData(rateLimit, organization),
                    targetRef: {
                        group: operation != () ? "dp.wso2.com" : "gateway.networking.k8s.io",
                        kind: operation != () ? "Resource" : "API",
                        name: targetRefName
                    }
                }
            };
        }
        return rateLimitPolicyCR;
    }

    isolated function retrieveRateLimitData(RateLimit rateLimit, commons:Organization organization) returns model:RateLimitData {
        model:RateLimitData rateLimitData = {
            api: {
                requestsPerUnit: rateLimit.requestsPerUnit,
                unit: rateLimit.unit
            }
        };
        return rateLimitData;
    }

    public isolated function generateAPIPolicyCR(APKConf apkConf, string targetRefName, APKOperations? operation, commons:Organization organization, model:APIPolicyData policyData) returns model:APIPolicy? {
        model:APIPolicy? apiPolicyCR = ();
        apiPolicyCR = {
            metadata: {
                name: self.retrieveAPIPolicyRefName(),
                labels: self.getLabels(apkConf, organization)
            },
            spec: {
                default: policyData,
                targetRef: {
                    group: "dp.wso2.com",
                    kind: operation != () ? "Resource" : "API",
                    name: targetRefName
                }
            }
        };
        return apiPolicyCR;
    }

    isolated function retrieveAPIPolicyDetails(model:APIArtifact apiArtifact, APKConf apkConf, APKOperations? operations, commons:Organization organization, APKOperationPolicy[]? policies, string flow) returns any[]|error {
        any[] policyReferences = [];
        if policies is APKOperationPolicy[] {
            foreach APKOperationPolicy policy in policies {
                string policyName = policy.policyName;
                if policy.parameters is record {} {
                    if (policyName == "Interceptor") {
                        InterceptorPolicy interceptorPolicy = check policy.cloneWithType(InterceptorPolicy);
                        InterceptorPolicy_parameters parameters = <InterceptorPolicy_parameters>interceptorPolicy?.parameters;
                        EndpointConfiguration endpointConfig = {endpoint: parameters.backendUrl ?: "", certificate: {secretName: parameters.tlsSecretName, secretKey: parameters.tlsSecretKey}};
                        model:Backend|error backendService = self.createBackendService(apiArtifact, apkConf, operations, INTERCEPTOR_TYPE, organization, endpointConfig);
                        string backendServiceName = "";
                        if backendService is model:Backend {
                            apiArtifact.backendServices[backendService.metadata.name] = (backendService);
                            backendServiceName = backendService.metadata.name;
                        }
                        model:InterceptorService? interceptorService = self.generateInterceptorServiceCR(parameters, backendServiceName, flow, apkConf, organization);
                        model:InterceptorReference? interceptorReference = ();
                        if interceptorService is model:InterceptorService {
                            apiArtifact.interceptorServices[interceptorService.metadata.name] = (interceptorService);
                            interceptorReference = {
                                name: interceptorService.metadata.name
                            };
                        }
                        policyReferences.push(interceptorReference);
                    } else if (policyName == "BackendJwt") {
                        BackendJWTPolicy backendJWTPolicy = check policy.cloneWithType(BackendJWTPolicy);
                        model:BackendJWT backendJwt = self.retrieveBackendJWTPolicy(apkConf, apiArtifact, backendJWTPolicy, organization);
                        apiArtifact.backendJwt = backendJwt;
                        policyReferences.push(<model:BackendJwtReference>{name: backendJwt.metadata.name});
                    } else {
                        return e909052(error("Incorrect API Policy name provided."));
                    }
                }
            }
        }
        return policyReferences;
    }
    private isolated function retrieveBackendJWTPolicy(APKConf apkConf, model:APIArtifact apiArtifact, BackendJWTPolicy backendJWTPolicy, commons:Organization organization) returns model:BackendJWT {
        BackendJWTPolicy_parameters parameters = backendJWTPolicy.parameters ?: {};
        model:BackendJWT backendJwt = {
            metadata: {
                name: self.getBackendJWTPolicyUid(apkConf, organization),
                labels: self.getLabels(apkConf, organization)
            },
            spec: {}
        };
        if parameters.encoding is string {
            backendJwt.spec.encoding = <string>parameters.encoding;
        }
        if parameters.signingAlgorithm is string {
            backendJwt.spec.signingAlgorithm = <string>parameters.signingAlgorithm;
        }
        if parameters.header is string {
            backendJwt.spec.header = <string>parameters.header;
        }
        if parameters.tokenTTL is int {
            backendJwt.spec.tokenTTL = <int>parameters.tokenTTL;
        }
        if parameters.customClaims is CustomClaims[] {
            model:CustomClaims[] backendJWTClaims = [];
            foreach CustomClaims customClaim in <CustomClaims[]>parameters?.customClaims {
                backendJWTClaims.push({
                    claim: customClaim.claim,
                    value: customClaim.value,
                    'type: customClaim.'type
                });
            }
            backendJwt.spec.customClaims = backendJWTClaims;
        }
        return backendJwt;
    }
    private isolated function retrieveCORSPolicyDetails(model:APIArtifact apiArtifact, APKConf apkConf, CORSConfiguration corsConfiguration, commons:Organization organization) returns model:CORSPolicy? {
        model:CORSPolicy corsPolicy = {};
        if corsConfiguration.accessControlAllowCredentials is boolean {
            corsPolicy.accessControlAllowCredentials = <boolean>corsConfiguration.accessControlAllowCredentials;
        }
        if corsConfiguration.accessControlAllowOrigins is string[] {
            corsPolicy.accessControlAllowOrigins = <string[]>corsConfiguration.accessControlAllowOrigins;
        }
        if corsConfiguration.accessControlAllowHeaders is string[] {
            corsPolicy.accessControlAllowHeaders = <string[]>corsConfiguration.accessControlAllowHeaders;
        }
        if corsConfiguration.accessControlAllowMethods is string[] {
            corsPolicy.accessControlAllowMethods = <string[]>corsConfiguration.accessControlAllowMethods;
        }
        if corsConfiguration.accessControlExposeHeaders is string[] {
            corsPolicy.accessControlExposeHeaders = <string[]>corsConfiguration.accessControlExposeHeaders;
        }
        if corsConfiguration.accessControlAllowMaxAge is int {
            corsPolicy.accessControlMaxAge = <int>corsConfiguration.accessControlAllowMaxAge;
        }
        return corsPolicy;
    }

    isolated function generateInterceptorServiceCR(InterceptorPolicy_parameters parameters, string interceptorBackend, string flow, APKConf apkConf, commons:Organization organization) returns model:InterceptorService? {
        model:InterceptorService? interceptorServiceCR = ();
        interceptorServiceCR = {
            metadata: {
                name: self.getInterceptorServiceUid(apkConf, organization, flow, 0),
                labels: self.getLabels(apkConf, organization)
            },
            spec: {
                backendRef: {
                    name: interceptorBackend
                },
                includes: self.getInterceptorIncludes(parameters, flow)
            }
        };
        return interceptorServiceCR;
    }

    isolated function getInterceptorIncludes(InterceptorPolicy_parameters parameters, string flow) returns string[] {
        string[] includes = [];
        if flow == "request" {
            if parameters.headersEnabled ?: false {
                includes.push("request_headers");
            }
            if parameters.bodyEnabled ?: false {
                includes.push("request_body");
            }
            if parameters.trailersEnabled ?: false {
                includes.push("request_trailers");
            }
            if parameters.contextEnabled ?: false {
                includes.push("invocation_context");
            }
        }
        if flow == "response" {
            if parameters.headersEnabled ?: false {
                includes.push("response_headers");
            }
            if parameters.bodyEnabled ?: false {
                includes.push("response_body");
            }
            if parameters.trailersEnabled ?: false {
                includes.push("response_trailers");
            }
            if parameters.contextEnabled ?: false {
                includes.push("invocation_context");
            }
        }
        return includes;
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

    public isolated function getInterceptorBackendUid(APKConf apkConf, string endpointType, commons:Organization organization, string|K8sService backend) returns string {
        string concatanatedString = string:'join("-", organization.name, apkConf.name, 'apkConf.'version, endpointType, self.construcURlFromService(backend));
        byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
        concatanatedString = hashedValue.toBase16();
        return "backend-" + concatanatedString + "-interceptor";
    }

    public isolated function getBackendJWTPolicyUid(APKConf apkConf, commons:Organization organization) returns string {
        string concatanatedString = string:'join("-", organization.name, apkConf.name, 'apkConf.'version);
        byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
        concatanatedString = hashedValue.toBase16();
        return string:'join("-", concatanatedString, "backend-jwt-policy");
    }
    public isolated function getBackendServiceUid(APKConf apkConf, APKOperations? apiOperation, string endpointType, commons:Organization organization) returns string {
        string concatanatedString = uuid:createType1AsString();
        if (apiOperation is APKOperations) {
            return "backend-" + concatanatedString + "-resource";
        } else {
            concatanatedString = string:'join("-", organization.name, apkConf.name, 'apkConf.'version, endpointType);
            byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
            concatanatedString = hashedValue.toBase16();
            return "backend-" + concatanatedString + "-api";
        }
    }

    public isolated function getInterceptorServiceUid(APKConf apkConf, commons:Organization organization, string flow, int interceptorIndex) returns string {
        string concatanatedString = string:'join("-", organization.name, apkConf.name, 'apkConf.'version);
        byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
        concatanatedString = hashedValue.toBase16();
        return flow + "-interceptor-service-" + interceptorIndex.toString() + "-" + concatanatedString + "-resource";
    }

    public isolated function getBackendPolicyUid(APKConf api, string endpointType, commons:Organization organization) returns string {
        string concatanatedString = uuid:createType1AsString();
        return "backendpolicy-" + concatanatedString;
    }

    public isolated function getBackendSecurityUid(string endpointType, commons:Organization organization) returns string {
        string concatanatedString = uuid:createType1AsString();
        return endpointType + "-" + concatanatedString + "-" + organization.name;
    }

    public isolated function getUniqueIdForAPI(string name, string 'version, commons:Organization organization) returns string {
        string concatanatedString = string:'join("-", organization.name, name, 'version);
        byte[] hashedValue = crypto:hashSha1(concatanatedString.toBytes());
        return hashedValue.toBase16();
    }

    public isolated function retrieveHttpRouteRefName(APKConf apkConf, string 'type, commons:Organization organization) returns string {
        return uuid:createType1AsString();
    }

    public isolated function retrieveAPIPolicyRefName() returns string {
        return uuid:createType1AsString();
    }

    public isolated function retrieveInterceptorBackendRefName() returns string {
        return uuid:createType1AsString();
    }

    public isolated function retrieveRateLimitPolicyRefName(APKOperations? operation) returns string {
        if operation is APKOperations {
            return uuid:createType1AsString();
        } else {
            return "api-" + uuid:createType1AsString();
        }
    }
    private isolated function validateAndRetrieveAPKConfiguration(json apkconfJson) returns APKConf|commons:APKError? {
        do {
            runtimeapi:APKConfValidationResponse validationResponse = check apkConfValidator.validate(apkconfJson.toJsonString());
            if validationResponse.isValidated() {
                APKConf apkConf = check apkconfJson.cloneWithType(APKConf);
                map<string> errors = {};
                self.validateEndpointConfigurations(apkConf, errors);
                if (errors.length() > 0) {
                    return e909029(errors);
                }
                return apkConf;
            } else {
                map<string> errorMap = {};
                foreach runtimeapi:ErrorHandler errorItem in check validationResponse.getErrorItems() {
                    errorMap[errorItem.getErrorMessage()] = errorItem.getErrorDescription();
                }
                return e909029(errorMap);
            }
        } on fail var e {
            return e909022("APK configuration is not valid", e);
        }
    }
    private isolated function validateEndpointConfigurations(APKConf apkConf, map<string> errors) {
        EndpointConfigurations? endpointConfigurations = apkConf.endpointConfigurations;
        boolean productionEndpointAvailable = false;
        boolean sandboxEndpointAvailable = false;
        if endpointConfigurations is EndpointConfigurations {
            sandboxEndpointAvailable = endpointConfigurations.sandbox is EndpointConfiguration;
            productionEndpointAvailable = endpointConfigurations.production is EndpointConfiguration;
        }
        APKOperations[]? operations = apkConf.operations;
        if operations is APKOperations[] {
            foreach APKOperations operation in operations {
                boolean operationLevelProductionEndpointAvailable = false;
                boolean operationLevelSandboxEndpointAvailable = false;
                EndpointConfigurations? endpointConfigs = operation.endpointConfigurations;
                if endpointConfigs is EndpointConfigurations {
                    operationLevelProductionEndpointAvailable = endpointConfigs.production is EndpointConfiguration;
                    operationLevelSandboxEndpointAvailable = endpointConfigs.sandbox is EndpointConfiguration;
                }
                    if (!operationLevelProductionEndpointAvailable && !productionEndpointAvailable) && (!operationLevelSandboxEndpointAvailable && !sandboxEndpointAvailable) {
                        errors["endpoint"] = "production/sandbox endpoint not available for " + <string>operation.target;
                    }

            }
        }
    }

    public isolated function prepareArtifact(record {|byte[] fileContent; string fileName; anydata...;|}? apkConfiguration, record {|byte[] fileContent; string fileName; anydata...;|}? definitionFile, commons:Organization organization) returns commons:APKError|model:APIArtifact {
        if apkConfiguration is () || definitionFile is () {
            return e909018("Required apkConfiguration ,definitionFile and apiType are not provided");
        }
        do {
            APKConf? apkConf = ();
            string apkConfContent = check string:fromBytes(apkConfiguration.fileContent);
            string|() convertedJson = check commons:newYamlUtil1().fromYamlStringToJson(apkConfContent);
            if convertedJson is string {
                json apkConfJson = check value:fromJsonString(convertedJson);
                apkConf = check self.validateAndRetrieveAPKConfiguration(apkConfJson);
            }
            string? apiDefinition = ();
            string definitionFileContent = check string:fromBytes(definitionFile.fileContent);
            string apiType = <string>apkConf?.'type;
            if apiType == API_TYPE_REST {
                if definitionFile.fileName.endsWith(".yaml") {
                    apiDefinition = check commons:newYamlUtil1().fromYamlStringToJson(definitionFileContent);
                } else if definitionFile.fileName.endsWith(".json") {
                    apiDefinition = definitionFileContent;
                }
            }
            if apkConf is () {
                return e909022("apkConfiguration is not provided", ());
            }
            APIClient apiclent = new ();
            return check apiclent.generateK8sArtifacts(apkConf, apiDefinition, organization);
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            log:printError("Error occured while prepare artifact", e);
            return e909022("Error occured while prepare artifact", e);
        }

    }
}
