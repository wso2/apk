import ballerina/mime;
import config_deployer_service.model;
import ballerina/http;
import wso2/apk_common_lib as commons;
import ballerina/log;
import ballerina/lang.value;

public class DeployerClient {
    public isolated function handleAPIDeployment(http:Request request) returns commons:APKError|http:Response {
        do {

            DeployApiBody deployAPIBody = check self.retrieveDeployApiBody(request);
            if deployAPIBody.apkConfiguration is () || deployAPIBody.definitionFile is () {
                return e909017();
            }
            APIClient apiclient = new;
            model:APIArtifact prepareArtifact = check apiclient.prepareArtifact(deployAPIBody?.apkConfiguration, deployAPIBody?.definitionFile);
            model:API deployAPIToK8sResult = check self.deployAPIToK8s(prepareArtifact);
            APKConf aPKConf = check self.getAPKConf(<record {|byte[] fileContent; string fileName; anydata...;|}>deployAPIBody.apkConfiguration);
            aPKConf.id = deployAPIToK8sResult.metadata.name;
            string|() apkYaml = check commons:newYamlUtil1().fromJsonStringToYaml(aPKConf.toJsonString());
            if apkYaml is string {
                http:Response response = new;
                response.setPayload(apkYaml);
                response.setHeader("Content-Type", "application/yaml");
                return response;
            } else {
                return e909022("Error occured while converting APKConf to YAML", ());
            }
        } on fail var e {
            if e is commons:APKError {
                return e;
            } else {
                return e909022("Error occured while converting APKConf to YAML", e);
            }
        }
    }
    public isolated function handleAPIUndeployment(string apiId) returns AcceptedString|BadRequestError|InternalServerErrorError|commons:APKError {
        model:Partition|() availablePartitionForAPI = check partitionResolver.getAvailablePartitionForAPI(apiId, "");
        if availablePartitionForAPI is model:Partition {
            model:API|() api = check getK8sAPIByNameAndNamespace(apiId, availablePartitionForAPI.namespace);
            if api is model:API {
                http:Response|http:ClientError apiCRDeletionResponse = deleteAPICR(api.metadata.name, availablePartitionForAPI.namespace);
                if apiCRDeletionResponse is http:ClientError {
                    log:printError("Error while undeploying API CR ", apiCRDeletionResponse);
                }
                string? definitionFileRef = api.spec.definitionFileRef;
                if definitionFileRef is string {
                    http:Response|http:ClientError apiDefinitionDeletionResponse = deleteConfigMap(definitionFileRef, availablePartitionForAPI.namespace);
                    if apiDefinitionDeletionResponse is http:ClientError {
                        log:printError("Error while undeploying API definition ", apiDefinitionDeletionResponse);
                    }
                }
                string response = string `API with id ${apiId} undeployed successfully`;
                json jsonResponse = {status: response};
                return {body: jsonResponse.toString()};
            } else {
                return e909001(apiId);
            }
        } else {
            return e909001(apiId);
        }
    }
    private isolated function getAPKConf(record {|byte[] fileContent; string fileName; anydata...;|} apkConfiguration) returns APKConf|commons:APKError {
        do {
            string apkConfContent = check string:fromBytes(apkConfiguration.fileContent);
            string convertedJson = check commons:newYamlUtil1().fromYamlStringToJson(apkConfContent) ?: "";
            APKConf apkConf = check value:fromJsonStringWithType(convertedJson, APKConf);
            return apkConf;
        } on fail var e {
            return e909022("Error occured while converting APKConf to YAML", e);
        }
    }

    private isolated function retrieveDeployApiBody(http:Request request) returns DeployApiBody|error {
        mime:Entity[] bodyParts = check request.getBodyParts();
        DeployApiBody deployApiBody = {};
        foreach mime:Entity item in bodyParts {
            mime:ContentDisposition contentDisposition = item.getContentDisposition();
            if contentDisposition.name == "apkConfiguration" {
                deployApiBody.apkConfiguration = {fileName: contentDisposition.fileName, fileContent: check item.getByteArray()};
            }
            if contentDisposition.name == "definitionFile" {
                deployApiBody.definitionFile = {fileName: contentDisposition.fileName, fileContent: check item.getByteArray()};
            }
        }
        return deployApiBody;
    }

    private isolated function deployAPIToK8s(model:APIArtifact apiArtifact) returns commons:APKError|model:API {
        do {
            model:Partition apiPartition;
            model:API? existingAPI;
            model:Partition|() availablePartitionForAPI = check partitionResolver.getAvailablePartitionForAPI(apiArtifact.uniqueId, apiArtifact.organization);
            if availablePartitionForAPI is model:Partition {
                apiPartition = availablePartitionForAPI;
                existingAPI = check getK8sAPIByNameAndNamespace(apiArtifact.uniqueId, apiPartition.namespace);
            } else {
                apiPartition = check partitionResolver.getDeployablePartition();
                existingAPI = ();
            }
            apiArtifact.namespace = apiPartition.namespace;
            if existingAPI is model:API {
                check self.deleteHttpRoutes(existingAPI, <string>apiArtifact?.organization);
                check self.deleteAuthneticationCRs(existingAPI, <string>apiArtifact?.organization);
                _ = check self.deleteScopeCrsForAPI(existingAPI, <string>apiArtifact?.organization);
                check self.deleteBackends(existingAPI, <string>apiArtifact?.organization);
                check self.deleteRateLimitPolicyCRs(existingAPI, <string>apiArtifact?.organization);
                check self.deleteAPIPolicyCRs(existingAPI, <string>apiArtifact?.organization);
                check self.deleteInterceptorServiceCRs(existingAPI, <string>apiArtifact?.organization);
            }
            model:ConfigMap? definition = apiArtifact.definition;
            if definition is model:ConfigMap {
                definition.metadata.namespace = apiPartition.namespace;
                _ = check self.deployConfigMap(definition);
            }
            check self.deployScopeCrs(apiArtifact);
            check self.deployBackendServices(apiArtifact);
            check self.deployAuthneticationCRs(apiArtifact);
            check self.deployRateLimitPolicyCRs(apiArtifact);
            check self.deployAPIPolicyCRs(apiArtifact);
            check self.deployInterceptorServiceCRs(apiArtifact);
            check self.deployHttpRoutes(apiArtifact.productionRoute, <string>apiArtifact?.namespace);
            check self.deployHttpRoutes(apiArtifact.sandboxRoute, <string>apiArtifact?.namespace);
            return check self.deployK8sAPICr(apiArtifact);
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            log:printError("Internal Error occured while deploying API", e);
            return e909028();
        }
    }

    private isolated function deployAPIPolicyCRs(model:APIArtifact apiArtifact) returns error? {
        foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
            http:Response deployAPIPolicyResult = check deployAPIPolicyCR(apiPolicy, <string>apiArtifact?.namespace);
            if deployAPIPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed APIPolicy Successfully" + apiPolicy.toString());
            } else {
                json responsePayLoad = check deployAPIPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deleteHttpRoutes(model:API api, string organization) returns commons:APKError? {
        do {
            model:HttprouteList|http:ClientError httpRouteListResponse = check getHttproutesForAPIS(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if httpRouteListResponse is model:HttprouteList {
                foreach model:Httproute item in httpRouteListResponse.items {
                    http:Response|http:ClientError httprouteDeletionResponse = deleteHttpRoute(item.metadata.name, <string>api.metadata?.namespace);
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
            return e909022("Error occured deleting httproutes", e);
        }
    }

    private isolated function deleteBackends(model:API api, string organization) returns commons:APKError? {
        do {
            model:BackendList|http:ClientError backendPolicyListResponse = check getBackendPolicyCRsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if backendPolicyListResponse is model:BackendList {
                foreach model:Backend item in backendPolicyListResponse.items {
                    http:Response|http:ClientError serviceDeletionResponse = deleteBackendPolicyCR(item.metadata.name, <string>api.metadata?.namespace);
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
            return e909022("Error occured deleting servicemapping", e);
        }
    }

    private isolated function deleteAuthneticationCRs(model:API api, string organization) returns commons:APKError? {
        do {
            model:AuthenticationList|http:ClientError authenticationCrListResponse = check getAuthenticationCrsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if authenticationCrListResponse is model:AuthenticationList {
                foreach model:Authentication item in authenticationCrListResponse.items {
                    http:Response|http:ClientError k8ServiceMappingDeletionResponse = deleteAuthenticationCR(item.metadata.name, <string>api.metadata?.namespace);
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
            return e909022("Error occured deleting servicemapping", e);
        }
    }

    private isolated function deleteScopeCrsForAPI(model:API api, string organization) returns commons:APKError? {
        do {
            model:ScopeList|http:ClientError scopeCrListResponse = check getScopeCrsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if scopeCrListResponse is model:ScopeList {
                foreach model:Scope item in scopeCrListResponse.items {
                    http:Response|http:ClientError scopeCrDeletionResponse = deleteScopeCr(item.metadata.name, <string>api.metadata?.namespace);
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
            return e909022("Error occured deleting scope", e);
        }
    }

    private isolated function deployScopeCrs(model:APIArtifact apiArtifact) returns error? {
        foreach model:Scope scope in apiArtifact.scopes {
            http:Response deployScopeResult = check deployScopeCR(scope, <string>apiArtifact?.namespace);
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
            model:API? k8sAPIByNameAndNamespace = check getK8sAPIByNameAndNamespace(k8sAPI.metadata.name, <string>apiArtifact?.namespace);
            if k8sAPIByNameAndNamespace is model:API {
                k8sAPI.metadata.resourceVersion = k8sAPIByNameAndNamespace.metadata.resourceVersion;
                http:Response deployAPICRResult = check updateAPICR(k8sAPI, <string>apiArtifact?.namespace);
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
                                return e909015(k8sAPI.spec.context);
                            } else if 'cause.'field == "spec.apiDisplayName" {
                                return e909016(k8sAPI.spec.apiDisplayName);
                            }
                        }
                        return e909017();
                    }
                    return self.handleK8sTimeout(statusResponse);
                }
            } else {
                http:Response deployAPICRResult = check deployAPICR(k8sAPI, <string>apiArtifact?.namespace);
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
                                return e909015(k8sAPI.spec.context);
                            } else if 'cause.'field == "spec.apiDisplayName" {
                                return e909016(k8sAPI.spec.apiDisplayName);
                            }
                        }
                        return e909017();
                    }
                    return self.handleK8sTimeout(statusResponse);
                }
            }
        } else {
            return e909022("Internal error occured", e = error("Internal error occured"));
        }
    }

    private isolated function deployHttpRoutes(model:Httproute[] httproutes, string namespace) returns error? {
        foreach model:Httproute httpRoute in httproutes {
            if httpRoute.spec.rules.length() > 0 {
                http:Response deployHttpRouteResult = check deployHttpRoute(httpRoute, namespace);
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

    private isolated function deployAuthneticationCRs(model:APIArtifact apiArtifact) returns error? {
        string[] keys = apiArtifact.authenticationMap.keys();
        log:printDebug("Inside Deploy Authentication CRs" + keys.toString());
        foreach string authenticationCrName in keys {
            model:Authentication authenticationCr = apiArtifact.authenticationMap.get(authenticationCrName);
            log:printDebug("Authentication CR:" + authenticationCr.toString());
            http:Response authenticationCrDeployResponse = check deployAuthenticationCR(authenticationCr, <string>apiArtifact?.namespace);
            if authenticationCrDeployResponse.statusCode == http:STATUS_CREATED {
                log:printInfo("Deployed Authentication Successfully" + authenticationCr.toString());
            } else {
                log:printError("Error Deploying Authentication" + authenticationCr.toString());
                json responsePayLoad = check authenticationCrDeployResponse.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deployBackendServices(model:APIArtifact apiArtifact) returns error? {
        foreach model:Backend backendService in apiArtifact.backendServices {
            http:Response deployServiceResult = check deployBackendCR(backendService, <string>apiArtifact?.namespace);
            if deployServiceResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Backend Successfully" + backendService.toString());
            } else {
                json responsePayLoad = check deployServiceResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deployConfigMap(model:ConfigMap definition) returns model:ConfigMap|commons:APKError|error {
        string deployableNamespace = <string>definition.metadata?.namespace;
        http:Response configMapRetrieved = check getConfigMapValueFromNameAndNamespace(definition.metadata.name, deployableNamespace);
        if configMapRetrieved.statusCode == 404 {
            http:Response deployConfigMapResult = check deployConfigMap(definition, deployableNamespace);
            if deployConfigMapResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Configmap Successfully" + definition.toString());
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                return check responsePayLoad.cloneWithType(model:ConfigMap);
            } else {
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                return self.handleK8sTimeout(statusResponse);
            }
        } else if configMapRetrieved.statusCode == 200 {
            http:Response deployConfigMapResult = check updateConfigMap(definition, deployableNamespace);
            if deployConfigMapResult.statusCode == http:STATUS_OK {
                log:printDebug("updated Configmap Successfully" + definition.toString());
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                return check responsePayLoad.cloneWithType(model:ConfigMap);
            } else {
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                return self.handleK8sTimeout(statusResponse);
            }
        } else {
            json responsePayLoad = check configMapRetrieved.getJsonPayload();
            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
            return self.handleK8sTimeout(statusResponse);
        }
    }

    private isolated function updateConfigMap(model:ConfigMap configMap) returns model:ConfigMap|commons:APKError|error {
        http:Response configMapRetrieved = check getConfigMapValueFromNameAndNamespace(configMap.metadata.name, <string>configMap.metadata?.namespace);
        if configMapRetrieved.statusCode == 200 {
            http:Response deployConfigMapResult = check updateConfigMap(configMap, <string>configMap.metadata?.namespace);
            if deployConfigMapResult.statusCode == http:STATUS_OK {
                log:printDebug("updated Configmap Successfully" + configMap.toString());
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                return check responsePayLoad.cloneWithType(model:ConfigMap);
            } else {
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                return self.handleK8sTimeout(statusResponse);
            }
        } else {
            json responsePayLoad = check configMapRetrieved.getJsonPayload();
            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
            return self.handleK8sTimeout(statusResponse);
        }
    }

    private isolated function deleteConfigMap(model:ConfigMap configMap) returns boolean|commons:APKError|error {
        http:Response configMapRetrieved = check getConfigMapValueFromNameAndNamespace(configMap.metadata.name, <string>configMap.metadata?.namespace);
        if configMapRetrieved.statusCode == 200 {
            http:Response deployConfigMapResult = check deleteConfigMap(configMap.metadata.name, <string>configMap.metadata?.namespace);
            if deployConfigMapResult.statusCode == http:STATUS_OK {
                log:printDebug("Configmap deleted Successfully" + configMap.toString());
                return true;
            } else {
                json responsePayLoad = check deployConfigMapResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                return self.handleK8sTimeout(statusResponse);
            }
        } else {
            json responsePayLoad = check configMapRetrieved.getJsonPayload();
            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
            return self.handleK8sTimeout(statusResponse);
        }
    }

    private isolated function deployRateLimitPolicyCRs(model:APIArtifact apiArtifact) returns error? {
        foreach model:RateLimitPolicy rateLimitPolicy in apiArtifact.rateLimitPolicies {
            http:Response deployRateLimitPolicyResult = check deployRateLimitPolicyCR(rateLimitPolicy, <string>apiArtifact?.namespace);
            if deployRateLimitPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed RateLimitPolicy Successfully" + rateLimitPolicy.toString());
            } else {
                json responsePayLoad = check deployRateLimitPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deleteRateLimitPolicyCRs(model:API api, string organization) returns commons:APKError? {
        do {
            model:RateLimitPolicyList|http:ClientError rateLimitPolicyCrListResponse = check getRateLimitPolicyCRsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if rateLimitPolicyCrListResponse is model:RateLimitPolicyList {
                foreach model:RateLimitPolicy item in rateLimitPolicyCrListResponse.items {
                    http:Response|http:ClientError rateLimitPolicyCRDeletionResponse = deleteRateLimitPolicyCR(item.metadata.name, <string>item.metadata?.namespace);
                    if rateLimitPolicyCRDeletionResponse is http:Response {
                        if rateLimitPolicyCRDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check rateLimitPolicyCRDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting rate limit policy");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting rate limit policy", e);
            return e909022("Error occured deleting rate limit policy", e);
        }
    }

    private isolated function deleteAPIPolicyCRs(model:API api, string organization) returns commons:APKError? {
        do {
            model:APIPolicyList|http:ClientError apiPolicyCrListResponse = check getAPIPolicyCRsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if apiPolicyCrListResponse is model:APIPolicyList {
                foreach model:APIPolicy item in apiPolicyCrListResponse.items {
                    http:Response|http:ClientError apiPolicyCRDeletionResponse = deleteAPIPolicyCR(item.metadata.name, <string>item.metadata?.namespace);
                    if apiPolicyCRDeletionResponse is http:Response {
                        if apiPolicyCRDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check apiPolicyCRDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting API policy");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting rate limit policy", e);
            return error("Error occured deleting rate limit policy", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }

    private isolated function deployInterceptorServiceCRs(model:APIArtifact apiArtifact) returns error? {
        foreach model:InterceptorService interceptorService in apiArtifact.interceptorServices {
            http:Response deployAPIPolicyResult = check deployInterceptorServiceCR(interceptorService, <string>apiArtifact?.namespace);
            if deployAPIPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed InterceptorService Successfully" + interceptorService.toString());
            } else {
                json responsePayLoad = check deployAPIPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
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

    private isolated function deleteInterceptorServiceCRs(model:API api, string organization) returns commons:APKError? {
        do {
            model:InterceptorServiceList|http:ClientError interceptorServiceListResponse = check getInterceptorServiceCRsForAPI(api.spec.apiDisplayName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if interceptorServiceListResponse is model:InterceptorServiceList {
                foreach model:InterceptorService item in interceptorServiceListResponse.items {
                    http:Response|http:ClientError interceptorServiceCRDeletionResponse = deleteInterceptorServiceCR(item.metadata.name, <string>item.metadata?.namespace);
                    if interceptorServiceCRDeletionResponse is http:Response {
                        if interceptorServiceCRDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check interceptorServiceCRDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting Interceptor Service");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting Interceptor Service", e);
            return error("Error occured deleting Interceptor Service", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }

}

