import ballerina/test;
import config_deployer_service.model;
import config_deployer_service.org.wso2.apk.config.model as runtimeModels;
import ballerina/io;

@test:Config {dataProvider: APIToAPKConfDataProvider}
public isolated function testFromAPIModelToAPKConf(runtimeModels:API api, APKConf expected) returns error? {
    APIClient apiClient = new;
    APKConf apkConf = check apiClient.fromAPIModelToAPKConf(api);
    test:assertEquals(apkConf, expected, "APKConf is not equal to expected APKConf");
}

@test:Config {}
public isolated function testCORSPolicyGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "API_CORS.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/API_CORS.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile);

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
public isolated function testBackendJWTConfigGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "API_CORS.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/API_CORS.apk-conf")};
    body.definitionFile = {fileName: "api_cors.yaml", fileContent: check io:fileReadBytes("./tests/resources/api_cors.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile);

    model:BackendJWTSpec backendJWTConfigSpec = {
        encoding: "base64",
        signingAlgorithm: "SHA256withRSA",
        header: "X-JWT-Assertion",
        tokenTTL: 3600,
        customClaims: [{claim: "claim1", value: "value1",'type:"string"}, {claim: "claim2", value: "value2",'type:"string"}]
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
public isolated function testBackendRetryAndTimeoutGenerationFromAPKConf() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "backend-retry.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/backend-retry.apk-conf")};
    body.definitionFile = {fileName: "backend-retry.yaml", fileContent: check io:fileReadBytes("./tests/resources/backend-retry.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile);

    model:Retry? retryConfigExpected = {
        count: 3,
        baseIntervalMillis: 1000,
        statusCodes: [504]
    };
    model:Timeout? timeoutConfigExpected = {
        maxRouteTimeoutSeconds: 60,
        routeIdleTimeoutSeconds: 400,
        routeTimeoutSeconds: 40
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
public isolated function testJWTAuthenticationOnlyEnable() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "jwtAuth.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/jwtAuth.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile);
    model:AuthenticationData expectedAuthenticationData = {
        disabled: false,
        authTypes: {
            jwt: {
                disabled: false,
                header: "Authorization",
                sendTokenToUpstream: false
            }
        }
    };
    foreach model:Authentication item in apiArtifact.authenticationMap {
        test:assertEquals(item.spec.override, expectedAuthenticationData);
    }
}

@test:Config {}
public isolated function testAPIKeyOnlyEnable() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "apiKeyOnly.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/apiKeyOnly.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile);
    model:AuthenticationData expectedAuthenticationData = {
        disabled: false,
        authTypes: {
            apiKey: [
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
    };
    foreach model:Authentication item in apiArtifact.authenticationMap {
        test:assertEquals(item.spec.override, expectedAuthenticationData);
    }
}

@test:Config {}
public isolated function testAPIKeyAndJWTEnable() returns error? {

    GenerateK8sResourcesBody body = {};
    body.apkConfiguration = {fileName: "jwtandAPIKey.apk-conf", fileContent: check io:fileReadBytes("./tests/resources/jwtandAPIKey.apk-conf")};
    body.definitionFile = {fileName: "api.yaml", fileContent: check io:fileReadBytes("./tests/resources/api.yaml")};
    body.apiType = "REST";
    APIClient apiClient = new;

    model:APIArtifact apiArtifact = check apiClient.prepareArtifact(body.apkConfiguration, body.definitionFile);
    model:AuthenticationData expectedAuthenticationData = {
        disabled: false,
        authTypes: {
            apiKey: [
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
            ],
            jwt: {
                disabled: false,
                header: "Authorization",
                sendTokenToUpstream: false
            }
        }
    };
    foreach model:Authentication item in apiArtifact.authenticationMap {
        test:assertEquals(item.spec.override, expectedAuthenticationData);
    }
}

public function APIToAPKConfDataProvider() returns map<[runtimeModels:API, APKConf]>|error {
    runtimeModels:API api = runtimeModels:newAPI1();
    api.setName("testAPI");
    api.setVersion("1.0.0");
    api.setContext("/test");
    runtimeModels:API api2 = runtimeModels:newAPI1();
    api2.setName("testAPI");
    api2.setVersion("1.0.0");
    api2.setContext("/test");
    api2.setEndpoint("http://localhost:9090");
    runtimeModels:API api3 = runtimeModels:newAPI1();
    api3.setName("testAPI");
    api3.setVersion("1.0.0");
    api3.setContext("/test");
    api3.setEndpoint("http://localhost:9090");
    runtimeModels:URITemplate[] uriTemplates = [];
    runtimeModels:URITemplate uriTemplate = runtimeModels:newURITemplate1();
    uriTemplate.setUriTemplate("/menu");
    uriTemplate.setHTTPVerb("GET");
    uriTemplates.push(uriTemplate);
    runtimeModels:URITemplate uriTemplate1 = runtimeModels:newURITemplate1();
    uriTemplate1.setUriTemplate("/order");
    uriTemplate1.setHTTPVerb("POST");
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
                context: "/test",
                version: "1.0.0",
                organization: "",
                operations: []
            }
        ],
        "2": [
            api2,
            {
                name: "testAPI",
                context: "/test",
                version: "1.0.0",
                organization: "",
                endpointConfigurations: {production: {endpoint: "http://localhost:9090"}},
                operations: []
            }

        ],
        "3": [
            api3,
            {
                name: "testAPI",
                context: "/test",
                version: "1.0.0",
                organization: "",
                endpointConfigurations: {production: {endpoint: "http://localhost:9090"}},
                operations: [
                    {target: "/menu", verb: "GET", authTypeEnabled: true, scopes: []},
                    {target: "/order", verb: "POST", authTypeEnabled: false, endpointConfigurations: {production: {endpoint: "http://localhost:9091"}}, scopes: ["scope1"]}
                ]
            }
        ]
    };
    return apkConfMap;
}
