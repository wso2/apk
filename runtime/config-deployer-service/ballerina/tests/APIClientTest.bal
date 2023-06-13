import ballerina/test;
import config_deployer_service.org.wso2.apk.config.model as runtimeModels;

@test:Config {dataProvider: APIToAPKConfDataProvider}
public isolated function testFromAPIModelToAPKConf(runtimeModels:API api, APKConf expected) returns error? {
    APIClient apiClient = new;
    APKConf apkConf = check apiClient.fromAPIModelToAPKConf(api);
    test:assertEquals(apkConf, expected, "APKConf is not equal to expected APKConf");
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
    uriTemplates .push(uriTemplate1);
    _ = check api3.setUriTemplates(uriTemplates);
    map<[runtimeModels:API, APKConf]> apkConfMap = {
        "1": [
            api,
            {
                name: "testAPI",
                context: "/test",
                version: "1.0.0",
                operations: []
            }
        ],
        "2": [
            api2,
            {
                name: "testAPI",
                context: "/test",
                version: "1.0.0",
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
                endpointConfigurations: {production: {endpoint: "http://localhost:9090"}},
                operations: [
                    {target: "/menu", verb: "GET", authTypeEnabled: true,scopes:[]},
                    {target: "/order", verb: "POST", authTypeEnabled: false, endpointConfigurations: {production: {endpoint: "http://localhost:9091"}}, scopes: ["scope1"]}
                ]
            }
        ]
    };
    return apkConfMap;
}
