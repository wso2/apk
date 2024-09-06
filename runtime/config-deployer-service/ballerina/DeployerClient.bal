import config_deployer_service.model;

import ballerina/http;
import ballerina/lang.value;
import ballerina/log;
import ballerina/mime;

import wso2/apk_common_lib as commons;

public class DeployerClient {
    public isolated function handleAPIDeployment(http:Request request, commons:Organization organization) returns commons:APKError|http:Response {
        do {

            DeployApiBody deployAPIBody = check self.retrieveDeployApiBody(request);
            if deployAPIBody.apkConfiguration is () || deployAPIBody.definitionFile is () {
                return e909017();
            }
            APIClient apiclient = new;
            model:APIArtifact prepareArtifact = check apiclient.prepareArtifact(deployAPIBody?.apkConfiguration, deployAPIBody?.definitionFile, organization);
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
    public isolated function handleAPIUndeployment(string apiId, commons:Organization organization) returns AcceptedString|BadRequestError|InternalServerErrorError|commons:APKError {
        model:Partition|() availablePartitionForAPI = check partitionResolver.getAvailablePartitionForAPI(apiId, "");
        if availablePartitionForAPI is model:Partition {
            model:API|() api = check getK8sAPIByNameAndNamespace(apiId, availablePartitionForAPI.namespace);
            if api is model:API {
                http:Response|http:ClientError apiCRDeletionResponse = deleteAPICR(api.metadata.name, availablePartitionForAPI.namespace);
                if apiCRDeletionResponse is http:ClientError {
                    log:printError("Error while undeploying API CR ", apiCRDeletionResponse);
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
                check self.deleteGrpcRoutes(existingAPI, <string>apiArtifact?.organization);
                check self.deleteAuthenticationCRs(existingAPI, <string>apiArtifact?.organization);
                _ = check self.deleteScopeCrsForAPI(existingAPI, <string>apiArtifact?.organization);
                check self.deleteBackends(existingAPI, <string>apiArtifact?.organization);
                check self.deleteRateLimitPolicyCRs(existingAPI, <string>apiArtifact?.organization);
                check self.deleteAPIPolicyCRs(existingAPI, <string>apiArtifact?.organization);
                check self.deleteInterceptorServiceCRs(existingAPI, <string>apiArtifact?.organization);
                check self.deleteBackendJWTConfig(existingAPI, <string>apiArtifact?.organization);
            }
            model:API? api = apiArtifact.api;
            if api is model:API {
                do {
                    model:API deployK8sAPICrResult = check self.deployK8sAPICr(apiArtifact);
                    model:OwnerReference ownerReference = {apiVersion: deployK8sAPICrResult.apiVersion, kind: deployK8sAPICrResult.kind, name: deployK8sAPICrResult.metadata.name, uid: <string>deployK8sAPICrResult.metadata.uid};
                    model:ConfigMap? definition = apiArtifact.definition;
                    if definition is model:ConfigMap {
                        definition.metadata.namespace = apiPartition.namespace;
                        definition.metadata.ownerReferences = [ownerReference];
                        _ = check self.deployConfigMap(definition);
                    }
                    check self.deployScopeCrs(apiArtifact, ownerReference);
                    check self.deployBackendServices(apiArtifact, ownerReference);
                    check self.deployAuthenticationCRs(apiArtifact, ownerReference);
                    check self.deployRateLimitPolicyCRs(apiArtifact, ownerReference);
                    check self.deployInterceptorServiceCRs(apiArtifact, ownerReference);
                    check self.deployBackendJWTConfigs(apiArtifact, ownerReference);
                    check self.deployAPIPolicyCRs(apiArtifact, ownerReference);
                    check self.deployRoutes(apiArtifact.productionHttpRoutes, apiArtifact.productionGqlRoutes, apiArtifact.productionGrpcRoutes, <string>apiArtifact?.namespace, ownerReference);
                    check self.deployRoutes(apiArtifact.sandboxHttpRoutes, apiArtifact.sandboxGqlRoutes, apiArtifact.sandboxGrpcRoutes, <string>apiArtifact?.namespace, ownerReference);

                    return deployK8sAPICrResult;
                } on fail var e {
                    http:Response|http:ClientError apiCRDeletionResponse = deleteAPICR(api.metadata.name, apiArtifact.namespace ?: "");
                    if apiCRDeletionResponse is http:ClientError {
                        log:printError("Error while undeploying API CR ", apiCRDeletionResponse);
                    }
                    if e is commons:APKError {
                        return e;
                    }
                    log:printError("Internal Error occured while deploying API", e);
                    return e909028();
                }
            } else {
                return e909028();
            }
        } on fail var e {
            if e is commons:APKError {
                return e;
            }
            log:printError("Internal Error occured while deploying API", e);
            return e909028();
        }
    }

    private isolated function deployAPIPolicyCRs(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
            apiPolicy.metadata.namespace = apiArtifact.namespace;
            http:Response deployAPIPolicyResult = check deployAPIPolicyCR(apiPolicy, <string>apiArtifact?.namespace);
            if deployAPIPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed APIPolicy Successfully" + apiPolicy.toString());
            } else if deployAPIPolicyResult.statusCode == http:STATUS_CONFLICT {
                log:printDebug("APIPolicy already exists" + apiPolicy.toString());
                model:APIPolicy retrievedApiPolicy = check getAPIPolicyCR(apiPolicy.metadata.name, <string>apiArtifact?.namespace);
                apiPolicy.metadata.resourceVersion = retrievedApiPolicy.metadata.resourceVersion;
                http:Response response = check updateAPIPolicyCR(apiPolicy, <string>apiArtifact?.namespace);
                if response.statusCode != http:STATUS_OK {
                    json responsePayLoad = check response.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            } else {
                json responsePayLoad = check deployAPIPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deleteHttpRoutes(model:API api, string organization) returns commons:APKError? {
        do {
            model:HTTPRouteList|http:ClientError httpRouteListResponse = check getHttproutesForAPIS(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if httpRouteListResponse is model:HTTPRouteList {
                foreach model:HTTPRoute item in httpRouteListResponse.items {
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

    private isolated function deleteGrpcRoutes(model:API api, string organization) returns commons:APKError? {
        do {
            model:GRPCRouteList|http:ClientError grpcRouteListResponse = check getGrpcRoutesForAPIs(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if grpcRouteListResponse is model:GRPCRouteList {
                foreach model:GRPCRoute item in grpcRouteListResponse.items {
                    http:Response|http:ClientError grpcRouteDeletionResponse = deleteGrpcRoute(item.metadata.name, <string>api.metadata?.namespace);
                    if grpcRouteDeletionResponse is http:Response {
                        if grpcRouteDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check grpcRouteDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting GrpcRoute", grpcRouteDeletionResponse);
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting grpcRoutes", e);
            return e909022("Error occured deleting grpcRoutes", e);
        }
    }

    private isolated function deleteBackends(model:API api, string organization) returns commons:APKError? {
        do {
            model:BackendList|http:ClientError backendPolicyListResponse = check getBackendPolicyCRsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
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

    private isolated function deleteAuthenticationCRs(model:API api, string organization) returns commons:APKError? {
        do {
            model:AuthenticationList|http:ClientError authenticationCrListResponse = check getAuthenticationCrsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
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
            model:ScopeList|http:ClientError scopeCrListResponse = check getScopeCrsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
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

    private isolated function deployScopeCrs(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        foreach model:Scope scope in apiArtifact.scopes {
            scope.metadata.ownerReferences = [ownerReference];
            http:Response deployScopeResult = check deployScopeCR(scope, <string>apiArtifact?.namespace);
            if deployScopeResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Scope Successfully" + scope.toString());
            } else if deployScopeResult.statusCode == http:STATUS_CONFLICT {
                log:printDebug("Scope already exists" + scope.toString());
                model:Scope scopeFromK8s = check getScopeCR(scope.metadata.name, <string>apiArtifact?.namespace);
                scope.metadata.resourceVersion = scopeFromK8s.metadata.resourceVersion;
                http:Response scopeCR = check updateScopeCR(scope, <string>apiArtifact?.namespace);
                if scopeCR.statusCode != http:STATUS_OK {
                    json responsePayLoad = check scopeCR.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
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
                            if 'cause.'field == "spec.basePath" {
                                log:printError("Error occurred while updating K8sAPI due to base path ", e909015(k8sAPI.spec.basePath));
                                return e909015(k8sAPI.spec.basePath);
                            } else if 'cause.'field == "spec.apiName" {
                                log:printError("Error occurred while updating K8sAPI due to base path ", e909015(k8sAPI.spec.basePath));
                                return e909016(k8sAPI.spec.apiName);
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
                            if 'cause.'field == "spec.basePath" {
                                return e909015(k8sAPI.spec.basePath);
                            } else if 'cause.'field == "spec.apiName" {
                                return e909016(k8sAPI.spec.apiName);
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
    private isolated function deployRoutes(model:HTTPRoute[]? httproutes, model:GQLRoute[]? gqlroutes, model:GRPCRoute[]? grpcroutes,
            string namespace, model:OwnerReference ownerReference) returns error? {
        if httproutes is model:HTTPRoute[] && httproutes.length() > 0 {
            model:HTTPRoute[] deployReadyHttproutes = httproutes;
            model:HTTPRoute[]|commons:APKError orderedHttproutes = self.createHttpRoutesOrder(httproutes);
            if orderedHttproutes is model:HTTPRoute[] {
                deployReadyHttproutes = orderedHttproutes;
            }
            foreach model:HTTPRoute httpRoute in deployReadyHttproutes {
                httpRoute.metadata.ownerReferences = [ownerReference];
                if httpRoute.spec.rules.length() > 0 {
                    http:Response deployHttpRouteResult = check deployHttpRoute(httpRoute, namespace);
                    if deployHttpRouteResult.statusCode == http:STATUS_CREATED {
                        log:printDebug("Deployed HttpRoute Successfully" + httpRoute.toString());
                    } else if deployHttpRouteResult.statusCode == http:STATUS_CONFLICT {
                        log:printDebug("HttpRoute already exists" + httpRoute.toString());
                        model:HTTPRoute httpRouteFromK8s = check getHttpRoute(httpRoute.metadata.name, namespace);
                        httpRoute.metadata.resourceVersion = httpRouteFromK8s.metadata.resourceVersion;
                        http:Response httpRouteCR = check updateHttpRoute(httpRoute, namespace);
                        if httpRouteCR.statusCode != http:STATUS_OK {
                            json responsePayLoad = check httpRouteCR.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        json responsePayLoad = check deployHttpRouteResult.getJsonPayload();
                        model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                        check self.handleK8sTimeout(statusResponse);
                    }
                }
            }
        } else if gqlroutes is model:GQLRoute[] && gqlroutes.length() > 0 {
            model:GQLRoute[] deployReadyGqlRoutes = gqlroutes;
            model:GQLRoute[]|commons:APKError orderedGqlRoutes = self.createGqlRoutesOrder(gqlroutes);
            if orderedGqlRoutes is model:GQLRoute[] {
                deployReadyGqlRoutes = orderedGqlRoutes;
            }
            foreach model:GQLRoute gqlRoute in deployReadyGqlRoutes {
                gqlRoute.metadata.ownerReferences = [ownerReference];
                if gqlRoute.spec.rules.length() > 0 {
                    http:Response deployGqlRouteResult = check deployGqlRoute(gqlRoute, namespace);
                    if deployGqlRouteResult.statusCode == http:STATUS_CREATED {
                        log:printDebug("Deployed GqlRoute Successfully" + gqlRoute.toString());
                    } else if deployGqlRouteResult.statusCode == http:STATUS_CONFLICT {
                        log:printDebug("GqlRoute already exists" + gqlRoute.toString());
                        model:GQLRoute gqlRouteFromK8s = check getGqlRoute(gqlRoute.metadata.name, namespace);
                        gqlRoute.metadata.resourceVersion = gqlRouteFromK8s.metadata.resourceVersion;
                        http:Response gqlRouteCR = check updateGqlRoute(gqlRoute, namespace);
                        if gqlRouteCR.statusCode != http:STATUS_OK {
                            json responsePayLoad = check gqlRouteCR.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        json responsePayLoad = check deployGqlRouteResult.getJsonPayload();
                        model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                        check self.handleK8sTimeout(statusResponse);
                    }
                }
            }
        } else if grpcroutes is model:GRPCRoute[] && grpcroutes.length() > 0 {
            model:GRPCRoute[] deployReadyGrpcRoutes = grpcroutes;
            model:GRPCRoute[]|commons:APKError orderedGrpcRoutes = self.createGrpcRoutesOrder(grpcroutes);
            if orderedGrpcRoutes is model:GRPCRoute[] {
                deployReadyGrpcRoutes = orderedGrpcRoutes;
            }
            foreach model:GRPCRoute grpcRoute in deployReadyGrpcRoutes {
                grpcRoute.metadata.ownerReferences = [ownerReference];
                if grpcRoute.spec.rules.length() > 0 {
                    http:Response deployGrpcRouteResult = check deployGrpcRoute(grpcRoute, namespace);
                    if deployGrpcRouteResult.statusCode == http:STATUS_CREATED {
                        log:printDebug("Deployed GrpcRoute Successfully" + grpcRoute.toString());
                    } else if deployGrpcRouteResult.statusCode == http:STATUS_CONFLICT {
                        log:printDebug("GrpcRoute already exists" + grpcRoute.toString());
                        model:GRPCRoute grpcRouteFromK8s = check getGrpcRoute(grpcRoute.metadata.name, namespace);
                        grpcRoute.metadata.resourceVersion = grpcRouteFromK8s.metadata.resourceVersion;
                        http:Response grpcRouteCR = check updateGrpcRoute(grpcRoute, namespace);
                        if grpcRouteCR.statusCode != http:STATUS_OK {
                            json responsePayLoad = check grpcRouteCR.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        json responsePayLoad = check deployGrpcRouteResult.getJsonPayload();
                        model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                        check self.handleK8sTimeout(statusResponse);
                    }
                }
            }
        }
    }

    public isolated function createGrpcRoutesOrder(model:GRPCRoute[] grpcRoutes) returns model:GRPCRoute[]|commons:APKError {
        do {
            foreach model:GRPCRoute route in grpcRoutes {
                model:GRPCRouteRule[] routeRules = route.spec.rules;
                model:GRPCRouteRule[] sortedRouteRules = from var routeRule in routeRules
                    order by <string>(routeRule.matches[0].method.'service) descending
                    select routeRule;
                route.spec.rules = sortedRouteRules;
            }
            return grpcRoutes;
        } on fail var e {
            log:printError("Error occured while sorting grpcRoutes", e);
            return e909022("Error occured while sorting grpcRoutes", e);
        }
    }

    public isolated function createHttpRoutesOrder(model:HTTPRoute[] httproutes) returns model:HTTPRoute[]|commons:APKError {
        do {
            foreach model:HTTPRoute route in httproutes {
                model:HTTPRouteRule[] routeRules = route.spec.rules;
                model:HTTPRouteRule[] sortedRouteRules = from var routeRule in routeRules
                    order by (<model:HTTPPathMatch>((<model:HTTPRouteMatch[]>routeRule.matches)[0]).path).value descending
                    select routeRule;
                route.spec.rules = sortedRouteRules;
            }
            return httproutes;
        } on fail var e {
            log:printError("Error occured while sorting httpRoutes", e);
            return e909022("Error occured while sorting httpRoutes", e);
        }
    }

    public isolated function createGqlRoutesOrder(model:GQLRoute[] gqlRoutes) returns model:GQLRoute[]|commons:APKError {
        do {
            foreach model:GQLRoute route in gqlRoutes {
                model:GQLRouteRule[] routeRules = route.spec.rules;
                model:GQLRouteRule[] sortedRouteRules = from var routeRule in routeRules
                    order by <string>((<model:GQLRouteMatch[]>routeRule.matches)[0]).path descending
                    select routeRule;
                route.spec.rules = sortedRouteRules;
            }
            return gqlRoutes;
        } on fail var e {
            log:printError("Error occured while sorting gqlRoutes", e);
            return e909022("Error occured while sorting gqlRoutes", e);
        }
    }

    private isolated function deployAuthenticationCRs(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        string[] keys = apiArtifact.authenticationMap.keys();
        log:printDebug("Inside Deploy Authentication CRs" + keys.toString());
        foreach string authenticationCrName in keys {
            model:Authentication authenticationCr = apiArtifact.authenticationMap.get(authenticationCrName);
            authenticationCr.metadata.ownerReferences = [ownerReference];
            log:printDebug("Authentication CR:" + authenticationCr.toString());
            http:Response authenticationCrDeployResponse = check deployAuthenticationCR(authenticationCr, <string>apiArtifact?.namespace);
            if authenticationCrDeployResponse.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Authentication Successfully" + authenticationCr.toString());
            } else if authenticationCrDeployResponse.statusCode == http:STATUS_CONFLICT {
                log:printDebug("Authentication CR already exists" + authenticationCr.toString());
                model:Authentication authenticationCrFromK8s = check getAuthenticationCR(authenticationCr.metadata.name, <string>apiArtifact?.namespace);
                authenticationCr.metadata.resourceVersion = authenticationCrFromK8s.metadata.resourceVersion;
                http:Response authenticationCR = check updateAuthenticationCR(authenticationCr, <string>apiArtifact?.namespace);
                if authenticationCR.statusCode != http:STATUS_OK {
                    json responsePayLoad = check authenticationCR.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            } else {
                log:printError("Error Deploying Authentication" + authenticationCr.toString());
                json responsePayLoad = check authenticationCrDeployResponse.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }
    private isolated function deployBackendServices(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        foreach model:Backend backendService in apiArtifact.backendServices {
            backendService.metadata.ownerReferences = [ownerReference];
            http:Response deployServiceResult = check deployBackendCR(backendService, <string>apiArtifact?.namespace);
            if deployServiceResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed Backend Successfully" + backendService.toString());
            } else if deployServiceResult.statusCode == http:STATUS_CONFLICT {
                log:printDebug("Backend already exists" + backendService.toString());
                model:Backend backendCRFromK8s = check getBackendCR(backendService.metadata.name, <string>apiArtifact?.namespace);
                backendService.metadata.resourceVersion = backendCRFromK8s.metadata.resourceVersion;
                http:Response backendCR = check updateBackendCR(backendService, <string>apiArtifact?.namespace);
                if backendCR.statusCode != http:STATUS_OK {
                    json responsePayLoad = check backendCR.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            } else {
                json responsePayLoad = check deployServiceResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deployConfigMap(model:ConfigMap definition) returns model:ConfigMap|commons:APKError|error {
        string deployableNamespace = <string>definition.metadata?.namespace;
        http:Response deployConfigMapResult = check deployConfigMap(definition, deployableNamespace);
        if deployConfigMapResult.statusCode == http:STATUS_CREATED {
            log:printDebug("Deployed Configmap Successfully" + definition.toString());
            json responsePayLoad = check deployConfigMapResult.getJsonPayload();
            return check responsePayLoad.cloneWithType(model:ConfigMap);
        } else if deployConfigMapResult.statusCode == http:STATUS_CONFLICT {
            log:printDebug("Configmap Already Exists" + definition.toString());
            model:ConfigMap configMapFromK8s = check getConfigMap(definition.metadata.name, deployableNamespace);
            definition.metadata.resourceVersion = configMapFromK8s.metadata.resourceVersion;
            deployConfigMapResult = check updateConfigMap(definition, deployableNamespace);
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
            json responsePayLoad = check deployConfigMapResult.getJsonPayload();
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

    private isolated function deployRateLimitPolicyCRs(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        foreach model:RateLimitPolicy rateLimitPolicy in apiArtifact.rateLimitPolicies {
            rateLimitPolicy.metadata.ownerReferences = [ownerReference];
            http:Response deployRateLimitPolicyResult = check deployRateLimitPolicyCR(rateLimitPolicy, <string>apiArtifact?.namespace);
            if deployRateLimitPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed RateLimitPolicy Successfully" + rateLimitPolicy.toString());
            } else if deployRateLimitPolicyResult.statusCode == http:STATUS_CONFLICT {
                log:printDebug("RateLimitPolicy already exists" + rateLimitPolicy.toString());
                model:RateLimitPolicy rateLimitPolicyFromK8s = check getRateLimitPolicyCR(rateLimitPolicy.metadata.name, <string>apiArtifact?.namespace);
                rateLimitPolicy.metadata.resourceVersion = rateLimitPolicyFromK8s.metadata.resourceVersion;
                http:Response rateLimitPolicyCR = check updateRateLimitPolicyCR(rateLimitPolicy, <string>apiArtifact?.namespace);
                if rateLimitPolicyCR.statusCode != http:STATUS_OK {
                    json responsePayLoad = check rateLimitPolicyCR.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            } else {
                json responsePayLoad = check deployRateLimitPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deleteRateLimitPolicyCRs(model:API api, string organization) returns commons:APKError? {
        do {
            model:RateLimitPolicyList|http:ClientError rateLimitPolicyCrListResponse = check getRateLimitPolicyCRsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
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
            model:APIPolicyList|http:ClientError apiPolicyCrListResponse = check getAPIPolicyCRsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
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

    private isolated function deployInterceptorServiceCRs(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        foreach model:InterceptorService interceptorService in apiArtifact.interceptorServices {
            interceptorService.metadata.ownerReferences = [ownerReference];
            http:Response deployAPIPolicyResult = check deployInterceptorServiceCR(interceptorService, <string>apiArtifact?.namespace);
            if deployAPIPolicyResult.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed InterceptorService Successfully" + interceptorService.toString());
            } else if deployAPIPolicyResult.statusCode == http:STATUS_CONFLICT {
                log:printDebug("InterceptorService already exists" + interceptorService.toString());
                model:InterceptorService interceptorServiceFromK8s = check getInterceptorServiceCR(interceptorService.metadata.name, <string>apiArtifact?.namespace);
                interceptorService.metadata.resourceVersion = interceptorServiceFromK8s.metadata.resourceVersion;
                http:Response interceptorServiceCR = check updateInterceptorServiceCR(interceptorService, <string>apiArtifact?.namespace);
                if interceptorServiceCR.statusCode != http:STATUS_OK {
                    json responsePayLoad = check interceptorServiceCR.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            } else {
                json responsePayLoad = check deployAPIPolicyResult.getJsonPayload();
                model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                check self.handleK8sTimeout(statusResponse);
            }
        }
    }

    private isolated function deployBackendJWTConfigs(model:APIArtifact apiArtifact, model:OwnerReference ownerReference) returns error? {
        model:BackendJWT? backendJwt = apiArtifact.backendJwt;
        if backendJwt is model:BackendJWT {
            backendJwt.metadata.ownerReferences = [ownerReference];
            http:Response backendJWTCrDeployResponse = check deployBackendJWTCr(backendJwt, <string>apiArtifact?.namespace);
            if backendJWTCrDeployResponse.statusCode == http:STATUS_CREATED {
                log:printDebug("Deployed BackendJWT Config Successfully" + backendJwt.toString());
            } else if backendJWTCrDeployResponse.statusCode == http:STATUS_CONFLICT {
                log:printDebug("BackendJWT Config already exists" + backendJwt.toString());
                model:BackendJWT backendJWTCrFromK8s = check getBackendJWTCr(backendJwt.metadata.name, <string>apiArtifact?.namespace);
                backendJwt.metadata.resourceVersion = backendJWTCrFromK8s.metadata.resourceVersion;
                http:Response backendJWTCr = check updateBackendJWTCr(backendJwt, <string>apiArtifact?.namespace);
                if backendJWTCr.statusCode != http:STATUS_OK {
                    json responsePayLoad = check backendJWTCr.getJsonPayload();
                    model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                    check self.handleK8sTimeout(statusResponse);
                }
            } else {
                json responsePayLoad = check backendJWTCrDeployResponse.getJsonPayload();
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
            model:InterceptorServiceList|http:ClientError interceptorServiceListResponse = check getInterceptorServiceCRsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
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

    private isolated function deleteBackendJWTConfig(model:API api, string organization) returns commons:APKError? {
        do {
            model:BackendJWTList|http:ClientError backendJWTlist = check getBackendJWTCrsForAPI(api.spec.apiName, api.spec.apiVersion, <string>api.metadata?.namespace, organization);
            if backendJWTlist is model:BackendJWTList {
                foreach model:BackendJWT item in backendJWTlist.items {
                    http:Response|http:ClientError backendJWTConfigDeletionResponse = deleteBackendJWTCr(item.metadata.name, <string>item.metadata?.namespace);
                    if backendJWTConfigDeletionResponse is http:Response {
                        if backendJWTConfigDeletionResponse.statusCode != http:STATUS_OK {
                            json responsePayLoad = check backendJWTConfigDeletionResponse.getJsonPayload();
                            model:Status statusResponse = check responsePayLoad.cloneWithType(model:Status);
                            check self.handleK8sTimeout(statusResponse);
                        }
                    } else {
                        log:printError("Error occured while deleting BackendJWT Config.");
                    }
                }
                return;
            }
        } on fail var e {
            log:printError("Error occured deleting BackendJWT Config", e);
            return error("Error occured deleting BackendJWT Config", message = "Internal Server Error", code = 909000, description = "Internal Server Error", statusCode = 500);
        }
    }

}

