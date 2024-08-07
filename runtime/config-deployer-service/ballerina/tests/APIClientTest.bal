import config_deployer_service.model;
import config_deployer_service.org.wso2.apk.config.model as runtimeModels;

import ballerina/io;
import ballerina/test;

import wso2/apk_common_lib;
import wso2/apk_common_lib as commons;

commons:Organization organization = {
    displayName: "default",
    name: "default",
    organizationClaimValue: "default",
    uuid: "",
    enabled: true
};

@test:Config {dataProvider: APIToAPKConfDataProvider}
public function testFromAPIModelToAPKConf(runtimeModels:API api, APKConf expected) returns error? {
    APIClient apiClient = new;
    APKConf apkConf = check apiClient.fromAPIModelToAPKConf(api);
    test:assertEquals(apkConf, expected, "APKConf is not equal to expected APKConf");
}

@test:Config {}
public function testCORSPolicyGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "API_CORS.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/API_CORS.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:CORSPolicy? corsPolicySpecExpected = {
        accessControlAllowCredentials: true,
        accessControlAllowOrigins: ["wso2.com"],
        accessControlAllowHeaders: ["Content-Type", "Authorization"],
        accessControlAllowMethods: ["GET"],
        accessControlMaxAge: 3600
    };

    foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
        model:APIPolicyData? policyData = apiPolicy.spec.default;
        if (policyData is model:APIPolicyData) {
            test:assertEquals(policyData.cORSPolicy, corsPolicySpecExpected, "CORS Policy is not equal to expected CORS Policy");
        }
    }
}

@test:Config {}
public function testInvalidCORSPolicyGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "invalid_API_CORS.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/invalid_API_CORS.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    map<string> errors = {
        "expected type: Boolean, found: String": "#/corsConfiguration/corsConfigurationEnabled: expected type: Boolean, found: String",
        "expected type: String, found: Integer": "#/corsConfiguration/accessControlAllowMethods/0: expected type: String, found: Integer"
    };

    apk_common_lib:APKError|model:APIArtifact apiArtifact = apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);
    if apiArtifact is model:APIArtifact {
        test:assertFail("Expected an error but got an APIArtifact");
    } else {
        apk_common_lib:ErrorHandler & readonly details = apiArtifact.detail();
        test:assertEquals(details.code, 909029);
        test:assertEquals(details.message, "Invalid apk-conf provided");
        test:assertEquals(details.description, "Invalid apk-conf provided");
        test:assertEquals(details.statusCode, 400);
        test:assertEquals(details.moreInfo, errors);
    }
}

@test:Config {}
public function testBackendJWTConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "API_CORS.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/API_CORS.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:BackendJWTSpec backendJWTConfigSpec = {
        encoding: "Base64",
        signingAlgorithm: "SHA256withRSA",
        header: "X-JWT-Assertion",
        tokenTTL: 3600,
        customClaims: [{claim: "claim1", value: "value1", 'type: "string"}, {claim: "claim2", value: "value2", 'type: "string"}]
    };
    test:assertTrue(apiArtifact.backendJwt is model:BackendJWT);
    model:BackendJWT? backendJwt = apiArtifact.backendJwt;
    if backendJwt is model:BackendJWT {
        test:assertEquals(backendJwt.spec, backendJWTConfigSpec, "Backend JWT Config is not equal to expected Backend JWT Config");
        model:BackendJwtReference backendJwtReference = {name: backendJwt.metadata.name};
        foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
            model:APIPolicyData? policyData = apiPolicy.spec.default;
            if (policyData is model:APIPolicyData) {
                test:assertEquals(policyData.backendJwtPolicy, backendJwtReference, "Backend JWT Config is not equal to expected Backend JWT Config");
            }
        }
    } else {
        test:assertFail("Backend JWT is not equal to expected Backend JWT Config");
    }
}

@test:Config {}
public function testInterceptorConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "API_Interceptors.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/API_Interceptors.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:InterceptorServiceSpec reqInterceptorServiceSpecExpected = {
        backendRef: {name: "backend-785a745eae91584e351e0f898194b1d40eb2cfe4-interceptor"},
        includes: [
            "request_headers",
            "request_body",
            "request_trailers",
            "invocation_context"
        ]
    };

    model:InterceptorServiceSpec resInterceptorServiceSpecExpected = {
        backendRef: {name: "backend-280575aea49bbdf7dad7437eecf9b3685cdda938-interceptor"},
        includes: [
            "response_body",
            "response_trailers"
        ]
    };

    test:assertEquals(apiArtifact.interceptorServices.length(), 2, "Required Interceptor services not defined");
    foreach model:InterceptorService interceptorService in apiArtifact.interceptorServices {
        string interceptorName = interceptorService.metadata.name;
        model:InterceptorReference interceptorReference = {name: interceptorName};
        if (interceptorName.startsWith("request-interceptor")) {
            test:assertEquals(interceptorService.spec, reqInterceptorServiceSpecExpected, "Request Interceptor is not equal to expected Request Interceptor");
            foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
                model:APIPolicyData? policyData = apiPolicy.spec.default;
                if (policyData is model:APIPolicyData) {
                    model:InterceptorReference[]? requestInterceptors = policyData.requestInterceptors;
                    if (requestInterceptors is model:InterceptorReference[]) {
                        foreach model:InterceptorReference reqInterceptorReference in requestInterceptors {
                            test:assertEquals(reqInterceptorReference, interceptorReference, "Request Interceptor ref is not equal to expected Request Interceptor ref");
                        }
                    } else {
                        test:assertFail("Request Interceptor references not found");
                    }
                }
            }
        } else if (interceptorName.startsWith("response-interceptor")) {
            test:assertEquals(interceptorService.spec, resInterceptorServiceSpecExpected, "Response Interceptor is not equal to expected Response Interceptor");
            foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
                model:APIPolicyData? policyData = apiPolicy.spec.default;
                if (policyData is model:APIPolicyData) {
                    model:InterceptorReference[]? responseInterceptors = policyData.responseInterceptors;
                    if (responseInterceptors is model:InterceptorReference[]) {
                        foreach model:InterceptorReference resInterceptorReference in responseInterceptors {
                            test:assertEquals(resInterceptorReference, interceptorReference, "Response Interceptor ref is not equal to expected Response Interceptor ref");
                        }
                    } else {
                        test:assertFail("Response Interceptor references not found");
                    }
                }
            }
        }
    }
}

@test:Config {}
public function testBackendConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "backends.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/backends.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:BackendSpec prodBackendSpec = {
        services: [
            {
                "host": "backend-prod-test",
                "port": 443
            }
        ],
        basePath: "/v1/",
        protocol: "https"
    };

    model:BackendSpec sandboxBackendSpec = {
        services: [
            {
                "host": "http-bin-backend.apk-test.svc.cluster.local",
                "port": 7676
            }
        ],
        protocol: "http"
    };

    test:assertEquals(apiArtifact.backendServices.length(), 3, "Required number of endpoints not found");
    test:assertTrue(apiArtifact.productionEndpointAvailable, "Production endpoint not defined");
    test:assertEquals(apiArtifact.productionHttpRoutes.length(), 1, "Production endpoint not defined");
    foreach model:HTTPRoute httpRoute in apiArtifact.productionHttpRoutes {
        test:assertEquals(httpRoute.spec.hostnames, ["default.gw.wso2.com"], "Production endpoint vhost mismatch");
        test:assertEquals(httpRoute.spec.rules.length(), 2, "Required number of HTTP Route rules not found");
        model:BackendRef[]? backendRefs = httpRoute.spec.rules[0].backendRefs;
        if backendRefs is model:BackendRef[] {
            string backendUUID = backendRefs[0].name;
            test:assertEquals(apiArtifact.backendServices.get(backendUUID).spec, prodBackendSpec, "Production Backend is not equal to expected Production Backend Config");
        } else {
            test:assertFail("Production backend references not found");
        }
    }

    test:assertTrue(apiArtifact.sandboxEndpointAvailable, "Sandbox endpoint not defined");
    test:assertEquals(apiArtifact.sandboxHttpRoutes.length(), 1, "Sandbox Backend not defined");
    foreach model:HTTPRoute httpRoute in apiArtifact.sandboxHttpRoutes {
        test:assertEquals(httpRoute.spec.hostnames, ["default.sandbox.gw.wso2.com"], "Sandbox vhost mismatch");
        model:BackendRef[]? backendRefs = httpRoute.spec.rules[0].backendRefs;
        if backendRefs is model:BackendRef[] {
            string backendUUID = backendRefs[0].name;
            test:assertEquals(apiArtifact.backendServices.get(backendUUID).spec, sandboxBackendSpec, "Sandbox Backend is not equal to expected Sandbox Backend Config");
        } else {
            test:assertFail("Sandbox backend references not found");
        }
    }
}

@test:Config {}
public function testAPILevelRateLimitConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "backends.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/backends.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:RateLimitData rateLimitData = {
        api: {
            requestsPerUnit: 5,
            unit: "Minute"
        }
    };

    test:assertEquals(apiArtifact.rateLimitPolicies.length(), 1, "Required number of Rate Limit policies not found");
    foreach model:RateLimitPolicy rateLimitPolicy in apiArtifact.rateLimitPolicies {
        test:assertEquals(rateLimitPolicy.spec.'default, rateLimitData, "Rate limit policy is not equal to expected Rate limit config");
        test:assertEquals(rateLimitPolicy.spec.targetRef.kind, "API", "Rate limit type is not equal to expected Rate limit type");
    }
}

@test:Config {}
public function testOperationLevelRateLimitConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "resource-level-rate-limit.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/resource-level-rate-limit.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:RateLimitData rateLimitData = {
        api: {
            requestsPerUnit: 10,
            unit: "Hour"
        }
    };

    test:assertEquals(apiArtifact.rateLimitPolicies.length(), 1, "Required number of Rate Limit policies not found");
    foreach model:RateLimitPolicy rateLimitPolicy in apiArtifact.rateLimitPolicies {
        test:assertEquals(rateLimitPolicy.spec.'default, rateLimitData, "Rate limit policy is not equal to expected Rate limit config");
        test:assertEquals(rateLimitPolicy.spec.targetRef.kind, "Resource", "Rate limit type is not equal to expected Rate limit type");
    }
}

@test:Config {}
public function testScopeConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "backends.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/backends.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    test:assertEquals(apiArtifact.scopes.length(), 3, "Required number of scopes not found");
    string[] scopeUUIDs = [
        apiArtifact.scopes.get("admin").metadata.name,
        apiArtifact.scopes.get("publisher").metadata.name,
        apiArtifact.scopes.get("reader").metadata.name
    ];
    foreach model:HTTPRoute httpRoute in apiArtifact.productionHttpRoutes {
        model:HTTPRouteFilter[]? httpFilters = httpRoute.spec.rules[0].filters;
        if httpFilters is model:HTTPRouteFilter[] {
            foreach model:HTTPRouteFilter httpFilter in httpFilters {
                if httpFilter.'type is string {
                    string httpFilterType = <string>httpFilter.'type;
                    if (httpFilterType.equalsIgnoreCaseAscii("ExtensionRef")) {
                        model:LocalObjectReference? extensionRef = httpFilter.extensionRef;
                        if extensionRef is model:LocalObjectReference {
                            test:assertEquals(extensionRef.kind, "Scope", "ExtensionRef for scope is not equal to expected Config");
                            test:assertTrue(scopeUUIDs.indexOf(extensionRef.name) != (), "Scope not found in the scope resources");
                        }
                    }
                }
            }
        } else {
            test:assertFail("HTTP Route filters not found");
        }
    }
}

@test:Config {}
public isolated function testRetrievePathPrefix() returns error? {

    APIClient apiClient = new;
    string pathPrefix;
    string organizaionStr = "wso2";
    string basePath = "test";
    string version = "1.0.0";
    string errorMessage = "Acual path prefix not equal to expected path prefix";
    commons:Organization organization = {
        displayName: organizaionStr,
        name: organizaionStr,
        organizationClaimValue: organizaionStr,
        uuid: "",
        enabled: true
    };

    pathPrefix = apiClient.retrievePathPrefix(basePath, version, "/", organization);
    test:assertEquals(pathPrefix, "/", errorMessage);

    pathPrefix = apiClient.retrievePathPrefix(basePath, version, "/*", organization);
    test:assertEquals(pathPrefix, "(.*)", errorMessage);

    pathPrefix = apiClient.retrievePathPrefix(basePath, version, "/employees/get", organization);
    test:assertEquals(pathPrefix, "/employees/get", errorMessage);

    pathPrefix = apiClient.retrievePathPrefix(basePath, version, "/employees/get/*", organization);
    test:assertEquals(pathPrefix, "/employees/get(.*)", errorMessage);
}

@test:Config {}
public isolated function testGeneratePrefixMatch() returns error? {

    APIClient apiClient = new;
    string prefixMatch;
    string urlStr = "https://backend-prod-test/v1/";
    string apiName = "sample-api";
    string method = "GET";
    string errorMessage = "Acual prefix match not equal to expected prefix match";

    model:Endpoint endpoint = {
        url: urlStr,
        name: apiName
    };
    APKOperations operations = {
        target: "/",
        verb: method,
        scopes: []
    };
    prefixMatch = apiClient.generatePrefixMatch(endpoint, operations);
    test:assertEquals(prefixMatch, "/", errorMessage);

    endpoint = {
        url: urlStr,
        name: apiName
    };
    operations = {
        target: "/*",
        verb: method,
        scopes: []
    };
    prefixMatch = apiClient.generatePrefixMatch(endpoint, operations);
    test:assertEquals(prefixMatch, "\\1", errorMessage);

    endpoint = {
        url: urlStr,
        name: apiName
    };
    operations = {
        target: "/employees/get",
        verb: method,
        scopes: []
    };
    prefixMatch = apiClient.generatePrefixMatch(endpoint, operations);
    test:assertEquals(prefixMatch, "/employees/get", errorMessage);

    endpoint = {
        url: urlStr,
        name: apiName
    };
    operations = {
        target: "/employees/get/*",
        verb: method,
        scopes: []
    };
    prefixMatch = apiClient.generatePrefixMatch(endpoint, operations);
    test:assertEquals(prefixMatch, "/employees/get///1", errorMessage);
}

@test:Config {}
public function testBackendRetryAndTimeoutGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "backend-retry.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/backend-retry.apk-conf")};
    body.definitionFile = {fileName: "backend-retry.yaml", fileContent: check io:fileReadBytes("./tests/resources/backend-retry.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    model:Retry? retryConfigExpected = {
        count: 3,
        baseIntervalMillis: 1000,
        statusCodes: [504]
    };
    model:Timeout? timeoutConfigExpected = {
        downstreamRequestIdleTimeout: 400,
        upstreamResponseTimeout: 40
    };
    model:CircuitBreaker? circuitBreakerConfigExpected = {
        maxConnectionPools: 200,
        maxConnections: 100,
        maxPendingRequests: 100,
        maxRequests: 100,
        maxRetries: 5
    };

    foreach model:Backend backend in apiArtifact.backendServices {
        model:Timeout? timeout = backend.spec.timeout;
        model:Retry? retryPolicy = backend.spec.'retry;
        model:CircuitBreaker? circuitBreakerPolicy = backend.spec.circuitBreaker;
        if (timeout is model:Timeout) {
            test:assertEquals(timeout, timeoutConfigExpected, "Timeout is not equal to expected Timeout Config");
        }
        if (retryPolicy is model:Retry) {
            test:assertEquals(retryPolicy, retryConfigExpected, "Retry Policy is not equal to expected Retry Policy");
        }
        if (circuitBreakerPolicy is model:CircuitBreaker) {
            test:assertEquals(circuitBreakerPolicy, circuitBreakerConfigExpected, "Circuit Breaker Policy is not equal to " +
            "expected Circuit Breaker Policy");
        }
    }
}

@test:Config {}
public function testJWTAuthenticationOnlyEnable() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "jwtAuth.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/jwtAuth.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);
    model:AuthenticationData expectedAuthenticationData = {
        disabled: false,
        authTypes: {
            oauth2: {
                required: "mandatory",
                disabled: false,
                header: "Authorization",
                sendTokenToUpstream: false
            }
        }
    };
    model:AuthenticationData expectedNoAuthentication = {
        disabled: true
    };

    foreach model:Authentication item in apiArtifact.authenticationMap {
        if string:endsWith(item.metadata.name, "-no-authentication") {
            test:assertEquals(item.spec.default, expectedNoAuthentication);
        } else {
            test:assertEquals(item.spec.default, expectedAuthenticationData);
        }
    }
}

@test:Config {}
public function testAPIKeyOnlyEnable() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "apiKeyOnly.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/apiKeyOnly.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);
    model:AuthenticationData expectedAuthenticationData = {
        disabled: false,
        authTypes: {
            apiKey: {
                required: "optional",
                keys: [
                    {
                        'in: "Header",
                        name: "apiKey",
                        sendTokenToUpstream: false
                    },
                    {
                        'in: "Query",
                        name: "apiKey",
                        sendTokenToUpstream: false
                    }
                ]
            }
        }
    };

    model:AuthenticationData expectedNoAuthentication = {
        disabled: true
    };
    foreach model:Authentication item in apiArtifact.authenticationMap {
        if string:endsWith(item.metadata.name, "-no-authentication") {
            test:assertEquals(item.spec.default, expectedNoAuthentication);
        } else {
            test:assertEquals(item.spec.default, expectedAuthenticationData);
        }
    }
}

@test:Config {}
public function testAPIKeyAndJWTEnable() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "jwtandAPIKey.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/jwtandAPIKey.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);
    model:AuthenticationData expectedAuthenticationData = {
        disabled: false,
        authTypes: {
            apiKey: {
                required: "optional",
                keys:
                [
                    {
                        'in: "Header",
                        name: "apiKey",
                        sendTokenToUpstream: false
                    },
                    {
                        'in: "Query",
                        name: "apiKey",
                        sendTokenToUpstream: false
                    }
                ]
            },
            oauth2: {
                required: "mandatory",
                disabled: false,
                header: "Authorization",
                sendTokenToUpstream: false
            }
        }
    };
    model:AuthenticationData expectedNoAuthentication = {
        disabled: true
    };
    foreach model:Authentication item in apiArtifact.authenticationMap {
        if string:endsWith(item.metadata.name, "-no-authentication") {
            test:assertEquals(item.spec.default, expectedNoAuthentication);
        } else {
            test:assertEquals(item.spec.default, expectedAuthenticationData);
        }
    }
}

@test:Config {}
public function testEnvironmentGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "multiEnv.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/multiEnv.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);
    model:API? api = apiArtifact.api;

    if api is model:API {
        model:APISpec apiSpec = <model:APISpec>api.spec;

        if apiSpec.environment != () {
            test:assertEquals(<string>apiSpec.environment, "dev", "Environment of the API is not equal to expected environment");
        } else {
            test:assertFail("Environment of the API should not be nil");
        }

    } else {
        test:assertFail("API is not equal to expected API Config");
    }

    model:HTTPRoute[] productionRoutes = apiArtifact.productionHttpRoutes;
    foreach var route in productionRoutes {
        test:assertEquals(route.spec.hostnames, ["default-dev.gw.wso2.com"], "Production endpoint vhost mismatch");
    }

}

@test:Config {}
public function testBasicAPIFromAPKConf() returns error? {

    string apiType = "REST";
    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "basicAPI.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/basicAPI.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = apiType;
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);
    model:API? api = apiArtifact.api;

    if api is model:API {
        model:APISpec apiSpec = <model:APISpec>api.spec;

        string? definitionFileRef = apiSpec.definitionFileRef;
        if definitionFileRef is string && definitionFileRef == "" {
            test:assertFail("Definition file ref is not equal to expected definition file ref");
        }

        test:assertEquals(<string>apiSpec.apiType, apiType, "API type is not equal to expected API type");

        if apiSpec.isDefaultVersion == () {
            test:assertFail("The field isDefaultVersion of the API should not be nil");
        }
        test:assertFalse(<boolean>apiSpec.isDefaultVersion, "The field isDefaultVersion of the API should be false");

        if apiSpec.systemAPI != () {
            test:assertFail("The field systemAPI of the API should be nil");
        }

        if apiSpec.environment != () {
            test:assertFail("Environment of the API should be nil");
        }

    } else {
        test:assertFail("API is not equal to expected API Config");
    }

    model:HTTPRoute[] productionRoutes = apiArtifact.productionHttpRoutes;
    foreach var route in productionRoutes {
        test:assertEquals(route.spec.hostnames, ["default.gw.wso2.com"], "Production endpoint vhost mismatch");
    }

}

@test:Config {}
public function testSubscriptionAPIPolicyGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "sub-validation.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/sub-validation.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile, organization);

    boolean expectedAPIPolicySpec = true;
    foreach model:APIPolicy apiPolicy in apiArtifact.apiPolicies {
        model:APIPolicyData? policyData = apiPolicy.spec.default;
        if (policyData is model:APIPolicyData) {
            test:assertEquals(policyData.subscriptionValidation, expectedAPIPolicySpec, "Subscription Policy is not equal to expected Subscription Policy");
        }
    }
}

public function APIToAPKConfDataProvider() returns map<[runtimeModels:API, APKConf]>|error {
    runtimeModels:API api = runtimeModels:newAPI1();
    api.setName("testAPI");
    api.setVersion("1.0.0");
    api.setBasePath("/test");
    runtimeModels:API api2 = runtimeModels:newAPI1();
    api2.setName("testAPI");
    api2.setVersion("1.0.0");
    api2.setBasePath("/test");
    api2.setEndpoint("http://localhost:9090");
    runtimeModels:API api3 = runtimeModels:newAPI1();
    api3.setName("testAPI");
    api3.setVersion("1.0.0");
    api3.setBasePath("/test");
    api3.setEndpoint("http://localhost:9090");
    runtimeModels:URITemplate[] uriTemplates = [];
    runtimeModels:URITemplate uriTemplate = runtimeModels:newURITemplate1();
    uriTemplate.setUriTemplate("/menu");
    uriTemplate.setVerb("GET");
    uriTemplates.push(uriTemplate);
    runtimeModels:URITemplate uriTemplate1 = runtimeModels:newURITemplate1();
    uriTemplate1.setUriTemplate("/order");
    uriTemplate1.setVerb("POST");
    uriTemplate1.setAuthEnabled(false);
    uriTemplate1.setEndpoint("http://localhost:9091");
    uriTemplate1.setScopes("scope1");
    uriTemplates.push(uriTemplate1);
    _ = check api3.setUriTemplates(uriTemplates);
    map<[runtimeModels:API, APKConf]> apkConfMap = {
        "1": [
            api,
            {
                name: "testAPI",
                basePath: "/test",
                version: "1.0.0",
                operations: []
            }
        ],
        "2": [
            api2,
            {
                name: "testAPI",
                basePath: "/test",
                version: "1.0.0",
                endpointConfigurations: {production: {endpoint: "http://localhost:9090"}},
                operations: []
            }

        ],
        "3": [
            api3,
            {
                name: "testAPI",
                basePath: "/test",
                version: "1.0.0",
                endpointConfigurations: {production: {endpoint: "http://localhost:9090"}},
                operations: [
                    {target: "/menu", verb: "GET", secured: true, scopes: []},
                    {target: "/order", verb: "POST", secured: false, endpointConfigurations: {production: {endpoint: "http://localhost:9091"}}, scopes: ["scope1"]}
                ]
            }
        ]
    };
    return apkConfMap;
}
