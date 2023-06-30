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
import ballerina/test;
import ballerina/websocket;
import ballerina/uuid;
import ballerina/http;
import runtime_domain_service.model as model;
import wso2/apk_common_lib as commons;
import runtime_domain_service.java.io;

commons:Organization organiztion1 = {
    name: "org1",
    displayName: "org1",
    uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
    organizationClaimValue: "org1",
    enabled: true,
    serviceListingNamespaces: ["*"],
    properties: []
};
commons:Organization organiztion2 = {
    name: "wso2.com",
    displayName: "wso2.com",
    uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114c",
    organizationClaimValue: "wso2.com",
    enabled: true,
    serviceListingNamespaces: ["*"],
    properties: []
};

@test:Mock {functionName: "getBackendPolicyUid"}
function testgetBackendPolicyUid(API api, string? endpointType, commons:Organization organization) returns string {
    return "backendpolicy-uuid";
}

@test:Mock {functionName: "retrieveHttpRouteRefName"}
function testRetrieveHttpRouteRefName(API api, string 'type, commons:Organization organization) returns string {
    return "http-route-ref-name";
}

@test:Mock {functionName: "retrieveRateLimitPolicyRefName"}
function testRetrieveRateLimitPolicyRefName(APIOperations? operaion) returns string {
    return "rate-limit-policy-ref-name";
}

@test:Mock {functionName: "retrieveAPIPolicyRefName"}
function testRetrieveAPIPolicyRefName() returns string {
    return "api-policy-ref-name";
}

int serviceMappingIndex = 0;

@test:Mock {functionName: "getServiceMappingClient"}
function getMockServiceMappingClient(string resourceVersion) returns websocket:Client|error {
    string initialConectionId = uuid:createType1AsString();
    if resourceVersion == "39433" {
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(getServiceMappingEvent());
        return mock;
    } else if resourceVersion == "5834" {
        string connectionId = uuid:createType1AsString();
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(connectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getNextServiceMappingEvent(), ());
        return mock;
    } else if resourceVersion == "23555" {
        if serviceMappingIndex == 0 {
            websocket:Error websocketError = error("Error", message = "Error");
            serviceMappingIndex += 1;
            return websocketError;
        } else {
            initialConectionId = uuid:createType1AsString();
            websocket:Client mock = test:mock(websocket:Client);
            test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
            test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
            test:prepare(mock).when("readMessage").thenReturn(());
            return mock;
        }
    } else {
        initialConectionId = uuid:createType1AsString();
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(());
        return mock;
    }

}


int configMapWatchIndex = 0;

@test:Mock {functionName: "getConfigMapWatchClient"}
function getTestConfigMapWatchClient(string resourceVersion) returns websocket:Client|error {
    string initialConectionId = uuid:createType1AsString();
    if resourceVersion == "28702" {
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(getConfigMapEvent());
        return mock;
    } else if resourceVersion == "28705" {
        string connectionId = uuid:createType1AsString();
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(connectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getConfigMapUpdateEvent(), ());
        return mock;
    } else if resourceVersion == "28714" {
        if configMapWatchIndex == 0 {
            websocket:Error websocketError = error("Error", message = "Error");
            configMapWatchIndex += 1;
            return websocketError;
        } else {
            initialConectionId = uuid:createType1AsString();
            websocket:Client mock = test:mock(websocket:Client);
            test:prepare(mock).when("isOpen").thenReturn(true);
            test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
            test:prepare(mock).when("readMessage").thenReturnSequence(getConfigMapDeleteEvent(), ());
            return mock;
        }
    } else {
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(());
        return mock;
    }
}

int serviceWatchIndex = 0;

@test:Mock {functionName: "getServiceClient"}
function getMockServiceClient(string resourceVersion) returns websocket:Client|error {
    websocket:Client mock;
    if resourceVersion == "39691" {
        string initialConectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getServiceEvent());
    } else if resourceVersion == "1514" {
        string initialConectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getNextMockServiceEvent(), ());
    } else if resourceVersion == "1517" {
        if serviceWatchIndex == 0 {
            websocket:Error websocketError = error("Error", message = "Error");
            serviceWatchIndex += 1;
            return websocketError;
        } else {
            string initialConectionId = uuid:createType1AsString();
            mock = test:mock(websocket:Client);
            test:prepare(mock).when("isOpen").thenReturn(true);
            test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
            test:prepare(mock).when("readMessage").thenReturn(());
        }
    } else {
        string initialConectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(());
    }
    return mock;
}

int apiIndex = 0;

@test:Mock {functionName: "getAPIClient"}
function getMockClient(string resourceVersion) returns websocket:Client|error {
    websocket:Client mock;
    if resourceVersion == "40316" {
        string initialConectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getMockWatchAPIEvent());
    } else if resourceVersion == "28702" {
        string initialConectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getNextMockWatchAPIEvent(), ());
    } else if resourceVersion == "28712" {
        if apiIndex == 0 {
            websocket:Error websocketError = error("Error", message = "Error");
            apiIndex += 1;
            return websocketError;
        } else {
            string initialConectionId = uuid:createType1AsString();
            mock = test:mock(websocket:Client);
            test:prepare(mock).when("isOpen").thenReturn(true);
            test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
            test:prepare(mock).when("readMessage").thenReturn(());
        }
    } else {
        string initialConectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(());
    }
    return mock;
}

@test:Mock {
    functionName: "initializeK8sClient"
}
function getMockK8sClient() returns http:Client|error {
    http:Client mockK8sClient = test:mock(http:Client);
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis")
        .thenReturn(getMockAPIList());
    string fieldSlector = "metadata.namespace%21%3Dkube-system%2Cmetadata.namespace%21%3Dkubernetes-dashboard%2Cmetadata.namespace%21%3Dgateway-system%2Cmetadata.namespace%21%3Dingress-nginx%2Cmetadata.namespace%21%3Dapk-platform";
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/services?fieldSelector=" + fieldSlector)
        .thenReturn(getMockServiceList());
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps?labelSelector=" + check getEncodedStringForLabelSelector())
        .thenReturn(getMockLabelList());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/servicemappings")
        .thenReturn(getMockServiceMappings());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b08-f2b1-1166-82d5-649ae706d29e").thenReturn(mock404Response());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk/apis/pizzashackAPI1").thenReturn(mock404Response());
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114a-definition").thenReturn(check mockConfigMaps());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/01ed7aca-eb6b-1178-a200-f604a4ce114a").thenReturn(getMockInternalAPI());
    http:ClientError clientError = error("Backend Failure");
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7b08-f2b1-1166-82d5-649ae706d29d-definition").thenReturn(mock404ConfigMap());
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114b-definition").thenReturn(clientError);
    test:prepare(mockK8sClient).when("get").withArguments("/apis/cp.wso2.com/v1alpha1/namespaces/apk-platform/organizations").thenReturn(getMockOrganizationList());
    return mockK8sClient;
}

function getMockOrganizationList() returns model:OrganizationList {
    model:OrganizationList organizationList = {
        apiVersion: "cp.wso2.com/v1alpha1",
        kind: "OrganizationList",
        metadata: {
            resourceVersion: "28702",
            selfLink: "/apis/cp.wso2.com/v1alpha1/organizations"
        },
        items: [
            {
                apiVersion: "cp.wso2.com/v1alpha1",
                kind: "Organization",
                metadata: {
                    name: "org1",
                    namespace: "apk-platform",
                    resourceVersion: "28702",
                    selfLink: "/apis/cp.wso2.com/v1alpha1/namespaces/apk-platform/organizations/org1",
                    uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b"
                },
                spec: {
                    name: "og1",
                    displayName: "org1",
                    uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                    organizationClaimValue: "org1",
                    enabled: true
                }
            }
        ]
    };
    return organizationList;
}

@test:Config {
    dataProvider: pathProvider
}
public function testretrievePathPrefix(string context, string 'version, string path, string expected) {
    APIClient apiclient = new ();
    string retrievePathPrefix = apiclient.retrievePathPrefix(context, 'version, path, organiztion1);
    test:assertEquals(retrievePathPrefix, expected);
}

function pathProvider() returns map<[string, string, string, string]>|error {
    map<[string, string, string, string]> dataSet = {
        "1": ["/abc/1.0.0", "1.0.0", "/abc", "/abc/1.0.0/abc"],
        "2": ["/abc", "1.0.0", "/abc", "/abc/1.0.0/abc"],
        "3": ["/abc", "1.0.0", "/*", "/abc/1.0.0(.*)"],
        "4": ["/abc/1.0.0", "1.0.0", "/*", "/abc/1.0.0(.*)"],
        "5": ["/abc/1.0.0", "1.0.0", "/{path}/abcd", "/abc/1.0.0/(.*)/abcd"],
        "6": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}", "/abc/1.0.0/path1/(.*)"],
        "7": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}/path2", "/abc/1.0.0/path1/(.*)/path2"],
        "8": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}/{pathparam2}", "/abc/1.0.0/path1/(.*)/(.*)"],
        "9": ["/abc", "1.0.0", "/path1/*", "/abc/1.0.0/path1(.*)"]
    };
    return dataSet;
}

@test:Config {dataProvider: contextVersionDataProvider}
public function testValidateContextAndVersion(string context, string 'version, boolean expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.validateContextAndVersion(context, 'version, organiztion1), expected);
}

function contextVersionDataProvider() returns map<[string, string, boolean]>|error {
    map<[string, string, boolean]> dataSet = {
        "1": ["/pizzashack/1.0.0", "1.0.0", true],
        "2": ["/pizzashack", "1.0.0", true],
        "3": ["/pizzashack", "2.0.0", false],
        "4": ["/pizzashack/2.0.0", "2.0.0", false],
        "5": ["/pizzashack3/1.0.0", "1.0.0", false],
        "6": ["/pizzashack3/", "1.0.0", false]
    };
    return dataSet;
}

@test:Config {dataProvider: nameDataProvider}
public function testValidateName(string name, commons:Organization organization, boolean expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.validateName(name, organization), expected);
}

function nameDataProvider() returns map<[string, commons:Organization, boolean]>|error {
    map<[string, commons:Organization, boolean]> dataSet = {
        "1": ["pizzashackAPI1", organiztion1, true],
        "2": ["pizzashackAPInew", organiztion1, false],
        "3": ["pizzashackAPI1", organiztion2, false]

    };
    return dataSet;
}

@test:Config {dataProvider: hostnameDataProvider}
public function testGetDomainPath(string url, string expectedDomain, string expectedPath, int expectedPort, string expectedHost) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getDomain(url), expectedDomain);
    test:assertEquals(apiclient.getPath(url), expectedPath);
    test:assertEquals(apiclient.getPort(url), expectedPort);
    test:assertEquals(apiclient.gethost(url), expectedHost);
}

function hostnameDataProvider() returns map<[string, string, string, int, string]>|error {
    map<[string, string, string, int, string]> dataSet = {
        "1": ["https://localhost/api.json", "https://localhost", "/api.json", 443, "localhost"],
        "2": ["http://localhost/api.json", "http://localhost", "/api.json", 80, "localhost"],
        "3": ["https://localhost:443/api.json", "https://localhost:443", "/api.json", 443, "localhost"],
        "4": ["http://localhost:80/api.json", "http://localhost:80", "/api.json", 80, "localhost"],
        "5": ["https://localhost", "https://localhost", "", 443, "localhost"],
        "6": ["http://localhost", "http://localhost", "", 80, "localhost"],
        "7": ["https://localhost:443", "https://localhost:443", "", 443, "localhost"],
        "8": ["http://localhost:80", "http://localhost:80", "", 80, "localhost"],
        "9": ["tcp://localhost:443", "", "", -1, ""]
    };
    return dataSet;
}

@test:Config {dataProvider: apiNameDataProvider}
public function testGetAPIByNameAndNamespace(string name, string namespace, model:API & readonly|() expected) {
    test:assertEquals(getAPIByNameAndNamespace(name, namespace, organiztion1), expected);
}

function apiNameDataProvider() returns map<[string, string, model:API & readonly|()]>|error {
    map<[string, string, model:API & readonly|()]> dataSet = {
        "1": ["01ed7aca-eb6b-1178-a200-f604a4ce114a", "apk-platform", getMockPizzaShakK8sAPI()],
        "2": ["01ed7b08-f2b1-1166-82d5-649ae706d29e", "apk-platform", ()],
        "3": ["pizzashackAPI1", "apk", ()]
    };
    return dataSet;
}

@test:Config {dataProvider: apiIDDataprovider}
public function testGetAPIById(string id, commons:Organization organization, anydata expected) returns error? {
    model:API? aPI = getAPI(id, organization);
    if aPI is model:API {
        test:assertEquals(aPI, expected);
    } else {
        test:assertEquals(aPI.toBalString(), expected);
    }
}

function apiIDDataprovider() returns map<[string, commons:Organization, anydata]>|error {

    map<[string, commons:Organization, anydata]> dataSet = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", organiztion1, getMockPizzaShakK8sAPI()]
    };
    return dataSet;
}

@test:Config {dataProvider: prefixMatchDataProvider}
public function testGeneratePrefixMatch(API api, model:Endpoint endpoint, APIOperations apiOperation, string expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.generatePrefixMatch(api, endpoint, apiOperation, PRODUCTION_TYPE), expected);

}

function prefixMatchDataProvider() returns map<[API, model:Endpoint, APIOperations, string]> {
    map<[API, model:Endpoint, APIOperations, string]> dataSet = {
        "1": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/order/{orderId}", verb: "POST"}, "/v3/f77cc767/order/\\1"],
        "2": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/menu", verb: "GET"}, "/v3/f77cc767/menu"],
        "3": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/menu", verb: "GET"}, "/v3/f77cc767/menu"],
        "4": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", serviceEntry: true, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/*", verb: "GET"}, "\\1"],
        "5": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/order/{orderId}/details/{item}", verb: "GET"}, "/v3/f77cc767/order/\\1/details/\\2"]
    };
    return dataSet;
}

@test:Config {dataProvider: apiDefinitionDataProvider}
public function testGetAPIDefinitionByID(string apiid, string? accept, anydata expectedResponse) returns error? {
    APIClient apiclient = new ();
    http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError aPIDefinitionByID = apiclient.getAPIDefinitionByID(apiid, organiztion1, accept);
    if aPIDefinitionByID is http:Response {
        string payload;
        if accept == APPLICATION_JSON_MEDIA_TYPE {
            test:assertEquals(aPIDefinitionByID.getContentType(), APPLICATION_JSON_MEDIA_TYPE);
            json jsonPayload = check aPIDefinitionByID.getJsonPayload();
            payload = jsonPayload.toBalString();
        } else if accept == APPLICATION_YAML_MEDIA_TYPE {
            test:assertEquals(aPIDefinitionByID.getContentType(), APPLICATION_YAML_MEDIA_TYPE);
            payload = check aPIDefinitionByID.getTextPayload();
        } else {
            test:assertEquals(aPIDefinitionByID.getContentType(), APPLICATION_JSON_MEDIA_TYPE);
            json jsonPayload = check aPIDefinitionByID.getJsonPayload();
            payload = jsonPayload.toBalString();

        }
        test:assertEquals(payload, expectedResponse);
    } else {
        if aPIDefinitionByID is any {
            test:assertEquals(aPIDefinitionByID.toBalString(), expectedResponse);
        }
        else {
            test:assertEquals(aPIDefinitionByID.toBalString(), expectedResponse);
        }
    }
}

public function apiDefinitionDataProvider() returns map<[string, string?, anydata]> {
    commons:APKError notfound = error commons:APKError("c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found",
        code = 909001,
        message = "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found",
        statusCode = 404,
        description = "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found"
    );
    commons:APKError internalError = error commons:APKError("Internal error occured while retrieving definition",
        error("Internal error occured while retrieving definition"),
        code = 909023,
        message = "Internal error occured while retrieving definition",
        statusCode = 500,
        description = "Internal error occured while retrieving definition"
    );
    commons:APKError error909041 = error commons:APKError("Accept header should be application/json or application/yaml",
        code = 909041,
        message = "Accept header should be application/json or application/yaml",
        statusCode = 406,
        description = "Accept header should be application/json or application/yaml"
    );
    do {

        map<[string, string?, anydata]> dataSet = {
            "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", APPLICATION_JSON_MEDIA_TYPE, mockOpenAPIJson().toBalString()],
            "2": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", (), mockOpenAPIJson().toBalString()],
            "3": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", APPLICATION_YAML_MEDIA_TYPE, check convertJsonToYaml(mockOpenAPIJson().toString())],
            "4": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e9", APPLICATION_JSON_MEDIA_TYPE, notfound.toBalString()],
            "5": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f9", APPLICATION_JSON_MEDIA_TYPE, mockpizzashackAPI11Definition().toBalString()],
            "6": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f1", APPLICATION_JSON_MEDIA_TYPE, mockPizzashackAPI12Definition().toBalString()],
            "7": ["7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1", APPLICATION_JSON_MEDIA_TYPE, mockPizzaShackAPI1Definition(organiztion1.uuid).toBalString()],
            "8": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f3", APPLICATION_JSON_MEDIA_TYPE, internalError.toBalString()],
            "9": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", "text/plain", error909041.toBalString()]

        };
        return dataSet;
    } on fail var e {
        test:assertFail(msg = e.message());
    }
}

@test:Config {dataProvider: apiByIdDataProvider}
public function testgetApiById(string apiid, commons:Organization organization, anydata expectedData) {
    APIClient apiclient = new ();
    API|NotFoundError|commons:APKError aPIById = apiclient.getAPIById(apiid, organization);
    if aPIById is any {
        test:assertEquals(aPIById.toBalString(), expectedData);
    } else {
        test:assertEquals(aPIById.toBalString(), expectedData);
    }
}

public function apiByIdDataProvider() returns map<[string, commons:Organization, anydata]> {
    API & readonly api1 = {
        id: "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
        name: "pizzashackAPI",
        context: "/pizzashack/1.0.0",
        'version: "1.0.0",
        'type: "REST",
        endpointConfig: {"endpoint_type": "http", "sandbox_endpoints": {"url": "https://pizzashack-service:8080/sample/pizzashack/v3/api/"}, "production_endpoints": {"url": "https://pizzashack-service:8080/sample/pizzashack/v3/api/"}},
        operations: [
            {target: "/*", verb: "GET", authTypeEnabled: true, "scopes": []},
            {target: "/*", verb: "PUT", authTypeEnabled: true, "scopes": []},
            {target: "/*", verb: "POST", authTypeEnabled: true, "scopes": []},
            {target: "/*", verb: "DELETE", authTypeEnabled: true, "scopes": []}
        ],
        createdTime: "2022-12-13T09:45:47Z"
    };
    commons:APKError notfound = error commons:APKError("c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found",
        code = 909001,
        message = "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found",
        statusCode = 404,
        description = "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found"
    );
    map<[string, commons:Organization, anydata]> dataset = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", organiztion1, api1.toBalString()],
        "2": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e9", organiztion1, notfound.toBalString()]
    };
    return dataset;
}

@test:Config {dataProvider: getApilistDataProvider}
public function testGetAPIList(string? query, int 'limit, int offset, string sortBy, string sortOrder, anydata expected) {
    APIClient apiclient = new ();
    any|error aPIList = apiclient.getAPIList(query, 'limit, offset, sortBy, sortOrder, organiztion1);
    if aPIList is any {
        test:assertEquals(aPIList.toBalString(), expected);
    } else {
        test:assertEquals(aPIList.toBalString(), expected);
    }
}

function getApilistDataProvider() returns map<[string?, int, int, string, string, anydata]> {
    commons:APKError badRequestError = error commons:APKError("Invalid Sort By/Sort Order value",
        code = 909020,
        message = "Invalid Sort By/Sort Order value",
        statusCode = 406,
        description = "Invalid Sort By/Sort Order value"
    );
    commons:APKError badRequest = error commons:APKError("Invalid keyword type1",
        code = 909019,
        message = "Invalid keyword type1",
        statusCode = 406,
        description = "Invalid keyword type1"
    );
    map<[string?, int, int, string, string, anydata]> dataSet = {
        "1": [
            (),
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 6,
                "list": [
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 6,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "2": [
            (),
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_DESC,
            {
                "count": 6,
                "list": [
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 6,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "3": [
            (),
            10,
            0,
            SORT_BY_CREATED_TIME,
            SORT_ORDER_ASC,
            {
                "count": 6,
                "list": [
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 6,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "4": [
            (),
            10,
            0,
            SORT_BY_CREATED_TIME,
            SORT_ORDER_DESC,
            {
                "count": 6,
                "list": [
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 6,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "5": [(), 10, 0, "description", SORT_ORDER_DESC, badRequestError.toBalString()],
        "6": [
            (),
            3,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 3,
                "list": [
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 3,
                    "total": 6,
                    "next": "/apis?limit=3&offset=3&sortBy=apiName&sortOrder=asc&query=",
                    "previous": ""
                }
            }.toBalString()
        ],
        "7": [
            (),
            3,
            3,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 3,
                "list": [
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "offset": 3,
                    "limit": 3,
                    "total": 6,
                    "next": "",
                    "previous": "/apis?limit=3&offset=0&sortBy=apiName&sortOrder=asc&query="
                }
            }.toBalString()
        ],
        "8": [
            (),
            3,
            6,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 0,
                "list": [],
                "pagination": {
                    "offset": 6,
                    "limit": 3,
                    "total": 6,
                    "next": "",
                    "previous": "/apis?limit=3&offset=3&sortBy=apiName&sortOrder=asc&query="
                }
            }.toBalString()
        ],
        "9": [
            "name:pizzashackAPI",
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 5,
                "list": [
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "10": [
            "pizzashackAPI",
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 5,
                "list": [
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "11": [
            "type:REST",
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 6,
                "list": [
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 6,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "12": [
            "type:WS",
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            {
                "count": 0,
                "list": [],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 0,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "13": [
            "type1:WS",
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            badRequest.toBalString()
        ]
    };
    return dataSet;
}

@test:Config {dataProvider: testDataGeneratedSwaggerDefinition}
public function testRetrieveGeneratedSwaggerDefinition(API api, string? definition, anydata expectedOutput) {
    APIClient apiclient = new;
    test:assertEquals(apiclient.retrieveGeneratedSwaggerDefinition(api, definition), expectedOutput);
}

function testDataGeneratedSwaggerDefinition() returns map<[API, string?, json|commons:APKError]> {
    map<[API, string?, json|commons:APKError]> data = {
        "1": [
            {
                "name": "demoAPI",
                "context": "/demoAPI/1.0.0",
                "version": "1.0.0"
            },
            (),
            {
                "openapi": "3.0.1",
                "info": {"title": "demoAPI", "version": "1.0.0"},
                "security": [{"default": []}],
                "paths": {},
                "components": {
                    "securitySchemes": {
                        "default": {
                            "type": "oauth2",
                            "flows": {"implicit": {"authorizationUrl": "https://test.com", "scopes": {}}}
                        }
                    }
                }
            }
        ],
        "2": [
            {
                "name": "demoAPI",
                "context": "/demoAPI/1.0.0",
                "version": "1.0.0",
                "operations": [{target: "/*", verb: "GET"}, {target: "/*", verb: "POST"}, {target: "/*", verb: "DELETE"}]
            },
            (),
            {
                "openapi": "3.0.1",
                "info": {
                    "title": "demoAPI",
                    "version": "1.0.0"
                },
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
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
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
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
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
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
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
            }
        ]
        ,
        "3": [
            {
                "name": "demoAPI",
                "context": "/demoAPI/1.0.0",
                "version": "1.0.0",
                "type": "REST",
                "operations": [{target: "/menu", verb: "GET"}, {target: "/order", verb: "POST"}, {target: "/order/{orderId}", verb: "GET"}]
            },
            {
                "openapi": "3.0.0",
                "info": {
                    "title": "PizzaShackAPI",
                    "description": "This is a RESTFul API for Pizza Shack online pizza delivery store.\n",
                    "contact": {
                        "name": "John Doe",
                        "url": "http://www.pizzashack.com",
                        "email": "architecture@pizzashack.com"
                    },
                    "license": {
                        "name": "Apache 2.0",
                        "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
                    },
                    "version": "1.0.0"
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
                    "/order": {
                        "post": {
                            "description": "Create a new Order",
                            "requestBody": {
                                "$ref": "#/components/requestBodies/Order"
                            },
                            "responses": {
                                "201": {
                                    "description": "Created. Successful response with the newly created object as entity inthe body.Location header contains URL of newly created entity.",
                                    "headers": {
                                        "Location": {
                                            "description": "The URL of the newly created resource.",
                                            "style": "simple",
                                            "explode": false,
                                            "schema": {
                                                "type": "string"
                                            }
                                        },
                                        "Content-Type": {
                                            "description": "The content type of the body.",
                                            "style": "simple",
                                            "explode": false,
                                            "schema": {
                                                "type": "string"
                                            }
                                        }
                                    },
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Order"
                                            }
                                        }
                                    }
                                },
                                "400": {
                                    "description": "Bad Request. Invalid request or validation error.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                },
                                "415": {
                                    "description": "Unsupported Media Type. The entity of the request was in a not supported format.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                }
                            },
                            "security": [
                                {
                                    "default": []
                                }
                            ]
                        }
                    },
                    "/menu": {
                        "get": {
                            "description": "Return a list of available menu items",
                            "responses": {
                                "200": {
                                    "description": "OK. List of APIs is returned.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "type": "array",
                                                "items": {
                                                    "$ref": "#/components/schemas/MenuItem"
                                                }
                                            }
                                        }
                                    }
                                },
                                "406": {
                                    "description": "Not Acceptable. The requested media type is not supported",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                }
                            },
                            "security": [
                                {
                                    "default": []
                                }
                            ],
                            "x-throttling-tier": "Unlimited"
                        }
                    },
                    "/order/{orderId}": {
                        "get": {
                            "description": "Get details of an Order",
                            "parameters": [
                                {
                                    "name": "orderId",
                                    "in": "path",
                                    "description": "Order Id",
                                    "required": true,
                                    "style": "simple",
                                    "explode": false,
                                    "schema": {
                                        "type": "string",
                                        "format": "string"
                                    }
                                }
                            ],
                            "responses": {
                                "200": {
                                    "description": "OK Requested Order will be returned",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Order"
                                            }
                                        }
                                    }
                                },
                                "404": {
                                    "description": "Not Found. Requested API does not exist.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                },
                                "406": {
                                    "description": "Not Acceptable. The requested media type is not supported",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                }
                            },
                            "security": [
                                {
                                    "default": []
                                }
                            ],
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
                        }
                    }
                },
                "components": {
                    "schemas": {
                        "ErrorListItem": {
                            "title": "Description of individual errors that may have occurred during a request.",
                            "required": [
                                "code",
                                "message"
                            ],
                            "properties": {
                                "message": {
                                    "type": "string",
                                    "description": "Description about individual errors occurred"
                                },
                                "code": {
                                    "type": "integer",
                                    "format": "int64"
                                }
                            }
                        },
                        "MenuItem": {
                            "title": "Pizza menu Item",
                            "required": [
                                "name"
                            ],
                            "properties": {
                                "price": {
                                    "type": "string"
                                },
                                "description": {
                                    "type": "string"
                                },
                                "name": {
                                    "type": "string"
                                },
                                "image": {
                                    "type": "string"
                                }
                            }
                        },
                        "Order": {
                            "title": "Pizza Order",
                            "required": [
                                "orderId"
                            ],
                            "properties": {
                                "customerName": {
                                    "type": "string"
                                },
                                "delivered": {
                                    "type": "boolean"
                                },
                                "address": {
                                    "type": "string"
                                },
                                "pizzaType": {
                                    "type": "string"
                                },
                                "creditCardNumber": {
                                    "type": "string"
                                },
                                "quantity": {
                                    "type": "number"
                                },
                                "orderId": {
                                    "type": "string"
                                }
                            }
                        },
                        "Error": {
                            "title": "Error object returned with 4XX HTTP status",
                            "required": [
                                "code",
                                "message"
                            ],
                            "properties": {
                                "message": {
                                    "type": "string",
                                    "description": "Error message."
                                },
                                "error": {
                                    "type": "array",
                                    "description": "If there are more than one error list them out. Ex. list out validation errors by each field.",
                                    "items": {
                                        "$ref": "#/components/schemas/ErrorListItem"
                                    }
                                },
                                "description": {
                                    "type": "string",
                                    "description": "A detail description about the error message."
                                },
                                "code": {
                                    "type": "integer",
                                    "format": "int64"
                                },
                                "moreInfo": {
                                    "type": "string",
                                    "description": "Preferably an url with more details about the error."
                                }
                            }
                        }
                    },
                    "requestBodies": {
                        "Order": {
                            "description": "Order object that needs to be added",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "$ref": "#/components/schemas/Order"
                                    }
                                }
                            },
                            "required": true
                        }
                    },
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
                },
                "x-wso2-auth-header": "Authorization",
                "x-wso2-cors": {
                    "corsConfigurationEnabled": false,
                    "accessControlAllowOrigins": [
                        "*"
                    ],
                    "accessControlAllowCredentials": false,
                    "accessControlAllowHeaders": [
                        "authorization",
                        "Access-Control-Allow-Origin",
                        "Content-Type",
                        "SOAPAction",
                        "apikey",
                        "Internal-Key"
                    ],
                    "accessControlAllowMethods": [
                        "GET",
                        "PUT",
                        "POST",
                        "DELETE",
                        "PATCH",
                        "OPTIONS"
                    ]
                },
                "x-wso2-production-endpoints": {
                    "urls": [
                        "https://localhost:9443/sample/pizzashack/v1/api/"
                    ],
                    "type": "http"
                },
                "x-wso2-sandbox-endpoints": {
                    "urls": [
                        "https://localhost:9443/sample/pizzashack/v1/api/"
                    ],
                    "type": "http"
                },
                "x-wso2-basePath": "/pizzashack/1.0.0",
                "x-wso2-transports": [
                    "http",
                    "https"
                ],
                "x-wso2-response-cache": {
                    "enabled": false,
                    "cacheTimeoutInSeconds": 300
                }
            }.toJsonString(),
            {
                "openapi": "3.0.0",
                "info": {
                    "title": "demoAPI",
                    "description": "This is a RESTFul API for Pizza Shack online pizza delivery store.\n",
                    "contact": {
                        "name": "John Doe",
                        "url": "http://www.pizzashack.com",
                        "email": "architecture@pizzashack.com"
                    },
                    "license": {
                        "name": "Apache 2.0",
                        "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
                    },
                    "version": "1.0.0"
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
                    "/order": {
                        "post": {
                            "description": "Create a new Order",
                            "requestBody": {
                                "$ref": "#/components/requestBodies/Order"
                            },
                            "responses": {
                                "201": {
                                    "description": "Created. Successful response with the newly created object as entity inthe body.Location header contains URL of newly created entity.",
                                    "headers": {
                                        "Location": {
                                            "description": "The URL of the newly created resource.",
                                            "style": "simple",
                                            "explode": false,
                                            "schema": {
                                                "type": "string"
                                            }
                                        },
                                        "Content-Type": {
                                            "description": "The content type of the body.",
                                            "style": "simple",
                                            "explode": false,
                                            "schema": {
                                                "type": "string"
                                            }
                                        }
                                    },
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Order"
                                            }
                                        }
                                    }
                                },
                                "400": {
                                    "description": "Bad Request. Invalid request or validation error.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                },
                                "415": {
                                    "description": "Unsupported Media Type. The entity of the request was in a not supported format.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                }
                            },
                            "security": [
                                {
                                    "default": []
                                }
                            ],
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
                        }
                    },
                    "/menu": {
                        "get": {
                            "description": "Return a list of available menu items",
                            "responses": {
                                "200": {
                                    "description": "OK. List of APIs is returned.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "type": "array",
                                                "items": {
                                                    "$ref": "#/components/schemas/MenuItem"
                                                }
                                            }
                                        }
                                    }
                                },
                                "406": {
                                    "description": "Not Acceptable. The requested media type is not supported",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                }
                            },
                            "security": [
                                {
                                    "default": []
                                }
                            ],
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
                        }
                    },
                    "/order/{orderId}": {
                        "get": {
                            "description": "Get details of an Order",
                            "parameters": [
                                {
                                    "name": "orderId",
                                    "in": "path",
                                    "description": "Order Id",
                                    "required": true,
                                    "style": "simple",
                                    "explode": false,
                                    "schema": {
                                        "type": "string",
                                        "format": "string"
                                    }
                                }
                            ],
                            "responses": {
                                "200": {
                                    "description": "OK Requested Order will be returned",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Order"
                                            }
                                        }
                                    }
                                },
                                "404": {
                                    "description": "Not Found. Requested API does not exist.",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                },
                                "406": {
                                    "description": "Not Acceptable. The requested media type is not supported",
                                    "content": {
                                        "application/json": {
                                            "schema": {
                                                "$ref": "#/components/schemas/Error"
                                            }
                                        }
                                    }
                                }
                            },
                            "security": [
                                {
                                    "default": []
                                }
                            ],
                            "x-auth-type": true,
                            "x-throttling-tier": "Unlimited"
                        }
                    }
                },
                "components": {
                    "schemas": {
                        "ErrorListItem": {
                            "title": "Description of individual errors that may have occurred during a request.",
                            "required": [
                                "code",
                                "message"
                            ],
                            "properties": {
                                "message": {
                                    "type": "string",
                                    "description": "Description about individual errors occurred"
                                },
                                "code": {
                                    "type": "integer",
                                    "format": "int64"
                                }
                            }
                        },
                        "MenuItem": {
                            "title": "Pizza menu Item",
                            "required": [
                                "name"
                            ],
                            "properties": {
                                "price": {
                                    "type": "string"
                                },
                                "description": {
                                    "type": "string"
                                },
                                "name": {
                                    "type": "string"
                                },
                                "image": {
                                    "type": "string"
                                }
                            }
                        },
                        "Order": {
                            "title": "Pizza Order",
                            "required": [
                                "orderId"
                            ],
                            "properties": {
                                "customerName": {
                                    "type": "string"
                                },
                                "delivered": {
                                    "type": "boolean"
                                },
                                "address": {
                                    "type": "string"
                                },
                                "pizzaType": {
                                    "type": "string"
                                },
                                "creditCardNumber": {
                                    "type": "string"
                                },
                                "quantity": {
                                    "type": "number"
                                },
                                "orderId": {
                                    "type": "string"
                                }
                            }
                        },
                        "Error": {
                            "title": "Error object returned with 4XX HTTP status",
                            "required": [
                                "code",
                                "message"
                            ],
                            "properties": {
                                "message": {
                                    "type": "string",
                                    "description": "Error message."
                                },
                                "error": {
                                    "type": "array",
                                    "description": "If there are more than one error list them out. Ex. list out validation errors by each field.",
                                    "items": {
                                        "$ref": "#/components/schemas/ErrorListItem"
                                    }
                                },
                                "description": {
                                    "type": "string",
                                    "description": "A detail description about the error message."
                                },
                                "code": {
                                    "type": "integer",
                                    "format": "int64"
                                },
                                "moreInfo": {
                                    "type": "string",
                                    "description": "Preferably an url with more details about the error."
                                }
                            }
                        }
                    },
                    "requestBodies": {
                        "Order": {
                            "description": "Order object that needs to be added",
                            "content": {
                                "application/json": {
                                    "schema": {
                                        "$ref": "#/components/schemas/Order"
                                    }
                                }
                            },
                            "required": true
                        }
                    },
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
                },
                "x-wso2-auth-header": "Authorization",
                "x-wso2-cors": {
                    "corsConfigurationEnabled": false,
                    "accessControlAllowOrigins": [
                        "*"
                    ],
                    "accessControlAllowCredentials": false,
                    "accessControlAllowHeaders": [
                        "authorization",
                        "Access-Control-Allow-Origin",
                        "Content-Type",
                        "SOAPAction",
                        "apikey",
                        "Internal-Key"
                    ],
                    "accessControlAllowMethods": [
                        "GET",
                        "PUT",
                        "POST",
                        "DELETE",
                        "PATCH",
                        "OPTIONS"
                    ]
                },
                "x-wso2-production-endpoints": {
                    "urls": [
                        "https://localhost:9443/sample/pizzashack/v1/api/"
                    ],
                    "type": "http"
                },
                "x-wso2-sandbox-endpoints": {
                    "urls": [
                        "https://localhost:9443/sample/pizzashack/v1/api/"
                    ],
                    "type": "http"
                },
                "x-wso2-basePath": "/pizzashack/1.0.0",
                "x-wso2-transports": [
                    "http",
                    "https"
                ],
                "x-wso2-response-cache": {
                    "enabled": false,
                    "cacheTimeoutInSeconds": 300
                }
            }
        ]
    }
;
    return data;
}

@test:Config {dataProvider: validateExistenceDataProvider}
function testValidateAPIExistence(string query, anydata expected) {
    APIClient apiClient = new;
    any|error validateAPIExistence = apiClient.validateAPIExistence(query, organiztion1);
    if validateAPIExistence is any {
        test:assertEquals(validateAPIExistence.toBalString(), expected);
    } else {
        test:assertEquals(validateAPIExistence.toBalString(), expected);
    }
}

function validateExistenceDataProvider() returns map<[string, anydata]> {
    http:Ok ok = {};
    commons:APKError badRequest = error commons:APKError("Invalid keyword type",
        code = 909019,
        message = "Invalid keyword type",
        statusCode = 406,
        description = "Invalid keyword type"
    );
    commons:APKError notFound = error commons:APKError("Context/Name doesn't exist",
        code = 909002,
        message = "Context/Name doesn't exist",
        statusCode = 404,
        description = "Context/Name doesn't exist"
    );
    map<[string, anydata]> data = {
        "1": ["name:pizzashackAPI", ok.toBalString()],
        "2": ["name:mockapi", notFound.toBalString()],
        "3": ["context:/api/v1", notFound.toBalString()],
        "4": ["context:/pizzashack/1.0.0", ok.toBalString()],
        "5": ["pizzashackAPI", ok.toBalString()],
        "6": ["type:pizzashackAPI", badRequest.toBalString()]
    };
    return data;
}

@test:Config {dataProvider: createApiFromServiceDataProvider}
function testCreateAPIFromService(string serviceUUId, string apiUUID, [model:ConfigMap, any] configmapResponse, [model:Httproute, any] httproute, [model:K8sServiceMapping, any] servicemapping, [model:API, any] k8sAPI, [model:RuntimeAPI, any] runtimeAPI, API api, string k8sapiUUID, [model:Backend, any][] backendServices, [model:RateLimitPolicy?, any] rateLimitPolicy, [model:APIPolicy?, any] apiPolicy, [model:InterceptorService, any][] interceptorServices, anydata expected) returns error? {
    APIClient apiClient = new;
    string username = "apkUser";
    http:Response configmapResponse404 = new;
    configmapResponse404.statusCode = 404;
    http:ApplicationResponseError internalApiResponse = error("", statusCode = 404, body = "Not Found", headers = {});
    model:HttprouteList httpRouteList = {metadata: {}, items: []};
    model:ServiceMappingList serviceMappingList = {metadata: {}, items: []};
    model:AuthenticationList authenticationList = {metadata: {}, items: []};
    model:BackendList backendList = {metadata: {}, items: []};
    model:ScopeList scopeList = {metadata: {}, items: []};
    model:RateLimitPolicyList rateLimitPolicyList = {metadata: {}, items: []};
    model:APIPolicyList apiPolicyList = {metadata: {}, items: []};
    model:InterceptorServiceList interceptorServiceList = {metadata: {}, items: []};
    http:Response internalAPIDeletionResponse = new;
    internalAPIDeletionResponse.statusCode = 200;

    foreach [model:Backend, any] backend in backendServices {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backends", backend[0]).thenReturn(backend[1]);
    }
    if rateLimitPolicy[0] is model:RateLimitPolicy {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/ratelimitpolicies", rateLimitPolicy[0]).thenReturn(rateLimitPolicy[1]);
    }
    if apiPolicy[0] is model:APIPolicy {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apipolicies", apiPolicy[0]).thenReturn(apiPolicy[1]);
    }
    foreach [model:InterceptorService, any] [interceptorService, interceptorServiceResponse] in interceptorServices {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/interceptorservices", interceptorService).thenReturn(interceptorServiceResponse);
    }

    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/ratelimitpolicies?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(rateLimitPolicyList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apipolicies?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(apiPolicyList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/interceptorservices?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(interceptorServiceList);
    test:prepare(k8sApiServerEp).when("post").withArguments("/api/v1/namespaces/apk-platform/configmaps", configmapResponse[0]).thenReturn(configmapResponse[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", httproute[0]).thenReturn(httproute[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings", servicemapping[0]).thenReturn(servicemapping[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis", k8sAPI[0]).thenReturn(k8sAPI[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis", runtimeAPI[0]).thenReturn(runtimeAPI[1]);
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/" + apiClient.retrieveDefinitionName(apiUUID)).thenReturn(configmapResponse404);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes/?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(httpRouteList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(serviceMappingList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/authentications?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(authenticationList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backends?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(backendList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/scopes?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(scopeList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sAPI[0].metadata.name).thenReturn(internalApiResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sAPI[0].metadata.name).thenReturn(runtimeAPI[0]);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/" + k8sAPI[0].metadata.name).thenReturn(configmapResponse404);
    test:prepare(k8sApiServerEp).when("delete").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sAPI[0].metadata.name).thenReturn(internalAPIDeletionResponse);
    any|error aPIFromService = apiClient.createAPIFromService(serviceUUId, api, organiztion1, username);
    if aPIFromService is any {
        test:assertEquals(aPIFromService.toBalString(), expected);
    } else {
        test:assertEquals(aPIFromService.toBalString(), expected);
    }
}

function createApiFromServiceDataProvider() returns map<[string, string, [model:ConfigMap, any], [model:Httproute, any], [model:K8sServiceMapping, any], [model:API, any], [model:RuntimeAPI, any], API, string, [model:Backend, any][], [model:RateLimitPolicy?, any], [model:APIPolicy?, any], [model:InterceptorService, any][], anydata]> {
    do {

        string k8sAPIUUID1 = uuid:createType1AsString();
        API api = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0"
        };
        API alreadyNameExist = {
            name: "pizzashackAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0"
        };
        json apiWithOperationPolicies = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "addHeader",
                                "parameters":
                                {
                                    "headerName": "customadd",
                                    "headerValue": "customvalue"
                                }
                            }
                        ],
                        "response": [
                            {
                                "policyName": "removeHeader",
                                "parameters":
                                {
                                    "headerName": "content-length"
                                }

                            }
                        ]
                    }
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ]
        };
        API apiWithInvalidPolicyName = {
            "name": "PizzaAPIOps",
            "context": "/pizzaAPIOps/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/menu",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "addHeader1",
                                "parameters":
                                {
                                    "headerName": "customadd",
                                    "headerValue": "customvalue"
                                }

                            }
                        ]
                    }
                }
            ]
        };
        commons:APKError invalidPolicyNameError = error commons:APKError("Invalid operation policy name",
        code = 909010,
        message = "Invalid operation policy name",
        statusCode = 406,
        description = "Invalid operation policy name"
    );
        API apiWithOperationRateLimits = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationRateLimit": {
                        "requestsPerUnit": 10,
                        "unit": "Minute"
                    }
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ]
        };
        API apiWithAPIRateLimits = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ],
            "apiRateLimit": {
                "requestsPerUnit": 10,
                "unit": "Minute"
            }
        };
        API apiWithBothRateLimits = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/menu",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationRateLimit": {
                        "requestsPerUnit": 10,
                        "unit": "Minute"
                    }
                }
            ],
            "apiRateLimit": {
                "requestsPerUnit": 10,
                "unit": "Minute"
            }
        };
        commons:APKError bothRateLimitsPresentError = error commons:APKError("Presence of both resource level and API level rate limits is not allowed",
            code = 909026,
            message = "Presence of both resource level and API level rate limits is not allowed",
            statusCode = 406,
            description = "Presence of both resource level and API level rate limits is not allowed"
        );
        json apiWithOperationLevelInterceptorPolicy = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "Interceptor",
                                "parameters":
                                {
                                    "headersEnabled": true,
                                    "bodyEnabled": false,
                                    "trailersEnabled": false,
                                    "contextEnabled": true,
                                    "backendUrl": "http://interceptor-backend1.interceptor:9082"
                                }
                            }
                        ],
                        "response": [
                            {
                                "policyName": "Interceptor",
                                "parameters":
                                {
                                    "headersEnabled": false,
                                    "bodyEnabled": true,
                                    "trailersEnabled": false,
                                    "contextEnabled": true,
                                    "backendUrl": "http://interceptor-backend2.interceptor:9083"
                                }
                            }
                        ]
                    }
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ]
        };
        json apiWithAPILevelInterceptorPolicy = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "operations": [
                {
                    target: "/*",
                    verb: "GET",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "PUT",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "POST",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "DELETE",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "PATCH",
                    authTypeEnabled: true
                }
            ],
            "apiPolicies": {
                "request": [
                    {
                        "policyName": "Interceptor",
                        "parameters":
                        {
                            "headersEnabled": true,
                            "bodyEnabled": false,
                            "trailersEnabled": false,
                            "contextEnabled": true,
                            "backendUrl": "http://interceptor-backend1.interceptor:9082"
                        }
                    }
                ],
                "response": [
                    {
                        "policyName": "Interceptor",
                        "parameters":
                        {
                            "headersEnabled": false,
                            "bodyEnabled": true,
                            "trailersEnabled": false,
                            "contextEnabled": true,
                            "backendUrl": "http://interceptor-backend2.interceptor:9083"
                        }
                    }
                ]
            }
        };
        string apiUUID = getUniqueIdForAPI(api.name, api.'version, organiztion1);
        model:ConfigMap configmap = check getMockConfigMap1(apiUUID, api);
        http:Response mockConfigMapResponse = getMockConfigMapResponse(configmap.clone());
        model:Httproute httpRoute = getMockHttpRoute(api, apiUUID, organiztion1);
        http:Response httpRouteResponse = getMockHttpRouteResponse(httpRoute.clone());
        model:Httproute httpRouteWithPolicies = getMockHttpRouteWithOperationPolicies1(api, apiUUID, organiztion1);
        http:Response httpRouteWithPoliciesResponse = getMockHttpRouteResponse(httpRouteWithPolicies.clone());
        model:Httproute httpRouteWithOperationRateLimits = getMockHttpRouteWithOperationRateLimits1(api, apiUUID, organiztion1);
        model:Httproute httpRouteWithInterceptorPolicy = getMockHttpRouteWithOperationInterceptorPolicy1(api, apiUUID, organiztion1);
        http:Response httpRouteWithInterceptorPolicyResponse = getMockHttpRouteResponse(httpRouteWithInterceptorPolicy.clone());
        http:Response httpRouteWithOperationRateLimitsResponse = getMockHttpRouteResponse(httpRouteWithOperationRateLimits.clone());
        model:K8sServiceMapping mockServiceMappingRequest = getMockServiceMappingRequest(api, apiUUID);
        model:API mockAPI = getMockAPI(api, apiUUID, organiztion1.uuid);
        http:Response mockAPIResponse = getMockAPIResponse(mockAPI.clone(), k8sAPIUUID1);
        Service serviceRecord = {
            name: "backend",
            namespace: "apk",
            id: "275b00d1-722c-4df2-b65a-9b14677abe4b",
            'type: "ClusterIP",
            portmapping: [
                {
                    name: "service",
                    protocol: "http",
                    port: 8080,
                    targetport: 8080
                }
            ]
        };
        string backenduuid = getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion1);
        string interceptorBackenduuid1 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend1.interceptor:9082");
        string interceptorBackenduuid2 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend2.interceptor:9083");
        model:Backend backendService = {
            metadata: {name: backenduuid, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: string:'join(".", 'serviceRecord.name, 'serviceRecord.namespace, "svc.cluster.local"), port: 80}], protocol: "http"}
        };
        model:Backend interceptorBackendService1 = {
            metadata: {name: interceptorBackenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend1.interceptor", port: 9082}], protocol: "http"}
        };
        model:Backend interceptorBackendService2 = {
            metadata: {name: interceptorBackenduuid2, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend2.interceptor", port: 9083}], protocol: "http"}
        };
        http:Response backendServiceResponse = getOKBackendServiceResponse(backendService);
        http:Response interceptorBackendServiceResponse1 = getOKBackendServiceResponse(interceptorBackendService1);
        http:Response interceptorBackendServiceResponse2 = getOKBackendServiceResponse(interceptorBackendService2);
        [model:Backend, any][] services = [];
        services.push([backendService, backendServiceResponse]);
        services.push([interceptorBackendService1, interceptorBackendServiceResponse1]);
        services.push([interceptorBackendService2, interceptorBackendServiceResponse2]);

        [model:InterceptorService, any][] interceptorServices = [];
        string interceptorBackendUrl1 =  "http://interceptor-backend1.interceptor:9082";
        string interceptorBackendUrl2 =  "http://interceptor-backend2.interceptor:9083";
        string[] requestIncludes = ["request_headers", "invocation_context"];
        string[] responseIncludes = ["response_body", "invocation_context"];
        model:InterceptorService requestInterceptorService = getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "request", requestIncludes, interceptorBackendUrl1);
        http:Response requestInterceptorServiceResponse = getMockInterceptorServiceResponse(getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "request", requestIncludes, interceptorBackendUrl1).clone());
        interceptorServices.push([requestInterceptorService, requestInterceptorServiceResponse]);
        model:InterceptorService responseInterceptorService = getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "response", responseIncludes, interceptorBackendUrl2);
        http:Response responseInterceptorServiceResponse = getMockInterceptorServiceResponse(getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "response", responseIncludes, interceptorBackendUrl2).clone());
        interceptorServices.push([responseInterceptorService, responseInterceptorServiceResponse]);

        model:RuntimeAPI mockRuntimeAPI = getMockRuntimeAPI(api, apiUUID, organiztion1, serviceRecord);
        http:Response mockRuntimeResponse = getMockRuntimeAPIResponse(mockRuntimeAPI.clone());
        model:RuntimeAPI mockRuntimeAPIWithPolicies = getMockRuntimeAPI(check apiWithOperationPolicies.cloneWithType(API), apiUUID, organiztion1, serviceRecord);
        http:Response mockRuntimeResponseWithPolicies = getMockRuntimeAPIResponse(mockRuntimeAPIWithPolicies.clone());
        model:RuntimeAPI mockRuntimeAOperationRateLimits = getMockRuntimeAPI(apiWithOperationRateLimits, apiUUID, organiztion1, serviceRecord);
        http:Response mockRuntimeResponseWithOperationRateLimits = getMockRuntimeAPIResponse(mockRuntimeAOperationRateLimits.clone());
        model:RuntimeAPI mockRuntimeAPIWithAPIRateLimits = getMockRuntimeAPI(apiWithAPIRateLimits, apiUUID, organiztion1, serviceRecord);
        http:Response mockRuntimeResponseWithAPIRateLimits = getMockRuntimeAPIResponse(mockRuntimeAPIWithAPIRateLimits.clone());
        model:RuntimeAPI mockRuntimeAPIWithOperationInterceptorPolicy = getMockRuntimeAPI(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), apiUUID, organiztion1, serviceRecord);
        http:Response mockRuntimeResponseWithOperationInterceptorPolicy = getMockRuntimeAPIResponse(mockRuntimeAPIWithOperationInterceptorPolicy.clone());
        model:RuntimeAPI mockRuntimeAPIWithAPIInterceptorPolicy = getMockRuntimeAPI(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), apiUUID, organiztion1, serviceRecord);
        http:Response mockRuntimeResponseWithAPIInterceptorPolicy = getMockRuntimeAPIResponse(mockRuntimeAPIWithAPIInterceptorPolicy.clone());
        http:Response serviceMappingResponse = getMockServiceMappingResponse(mockServiceMappingRequest.clone());
        commons:APKError nameAlreadyExistError = error commons:APKError(
            "API Name - " + alreadyNameExist.name + " already exist",
            code = 909011,
            message = "API Name - " + alreadyNameExist.name + " already exist",
            statusCode = 409,
            description = "API Name - " + alreadyNameExist.name + " already exist"
        );
        API contextAlreadyExist = {
            name: "PizzaAPI",
            context: "/pizzashack/1.0.0",
            'version: "1.0.0"
        };
        commons:APKError contextAlreadyExistError = error commons:APKError(
            "API Context - " + contextAlreadyExist.context + " already exist",
            code = 909012,
            message = "API Context - " + contextAlreadyExist.context + " already exist",
            statusCode = 409,
            description = "API Context - " + contextAlreadyExist.context + " already exist"
        );
        commons:APKError serviceNotExist = error commons:APKError("275b00d1-722c-4df2-b65a-9b14677abe4a service does not exist",
            code = 909047,
            message = "275b00d1-722c-4df2-b65a-9b14677abe4a service does not exist",
            statusCode = 404,
            description = "275b00d1-722c-4df2-b65a-9b14677abe4a service does not exist"
        );
        string locationUrl = runtimeConfiguration.baseURl + "/apis/" + k8sAPIUUID1;
        CreatedAPI createdAPI = {
            body: {
                id: k8sAPIUUID1,
                name: "PizzaAPI",
                context: "/pizzaAPI/1.0.0",
                'version: "1.0.0",
                'type: "REST",
                operations: [
                    {target: "/*", verb: "GET", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PUT", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "POST", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "DELETE", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PATCH", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}}
                ],
                serviceInfo: {name: "backend", namespace: "apk"},
                createdTime: "2023-01-17T11:23:49Z"
            },
            headers: {
                location: locationUrl
            }
        };
        json requestPolicy = {
            "policyName": "addHeader",
            "policyVersion": "v1",
            "parameters":
            {
                "headerName": "customadd",
                "headerValue": "customvalue"
            }
        };
        json responsePolicy = {
            "policyName": "removeHeader",
            "policyVersion": "v1",
            "parameters":
            {
                "headerName": "content-length"
            }
        };
        json requestInterceptorPolicy = {
            "policyName": "Interceptor",
            "policyVersion": "v1",
            "parameters":
            {
                "headersEnabled": true,
                "bodyEnabled": false,
                "trailersEnabled": false,
                "contextEnabled": true,
                "backendUrl": "http://interceptor-backend1.interceptor:9082"
            }
        };
        json responseInterceptorPolicy = {
            "policyName": "Interceptor",
            "policyVersion": "v1",
            "parameters":
            {
                "headersEnabled": false,
                "bodyEnabled": true,
                "trailersEnabled": false,
                "contextEnabled": true,
                "backendUrl": "http://interceptor-backend2.interceptor:9083"
            }
        };
        APIRateLimit rateLimit = {
            requestsPerUnit: 10,
            unit: "Minute"
        };
        CreatedAPI createdAPIWithPolicies = {
            body: {
                id: k8sAPIUUID1,
                name: "PizzaAPI",
                context: "/pizzaAPI/1.0.0",
                'version: "1.0.0",
                'type: "REST",
                operations: [
                    {target: "/*", verb: "GET", authTypeEnabled: true, scopes: [], operationPolicies: {request: [check requestPolicy.cloneWithType(OperationPolicy)], response: [check responsePolicy.cloneWithType(OperationPolicy)]}},
                    {target: "/*", verb: "PUT", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "POST", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "DELETE", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PATCH", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}}
                ],
                serviceInfo: {name: "backend", namespace: "apk"},
                createdTime: "2023-01-17T11:23:49Z"
            },
            headers: {
                location: locationUrl
            }
        };
        CreatedAPI createdAPIWithOperationRateLimits = {
            body: {
                id: k8sAPIUUID1,
                name: "PizzaAPI",
                context: "/pizzaAPI/1.0.0",
                'version: "1.0.0",
                'type: "REST",
                operations: [
                    {target: "/*", verb: "GET", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}, operationRateLimit: rateLimit},
                    {target: "/*", verb: "PUT", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "POST", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "DELETE", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PATCH", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}}
                ],
                serviceInfo: {name: "backend", namespace: "apk"},
                createdTime: "2023-01-17T11:23:49Z"
            },
            headers: {
                location: locationUrl
            }
        };
        CreatedAPI createdAPIWithAPIRateLimits = {
            body: {
                id: k8sAPIUUID1,
                name: "PizzaAPI",
                context: "/pizzaAPI/1.0.0",
                'version: "1.0.0",
                'type: "REST",
                operations: [
                    {target: "/*", verb: "GET", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PUT", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "POST", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "DELETE", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PATCH", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}}
                ],
                apiRateLimit: rateLimit,
                serviceInfo: {name: "backend", namespace: "apk"},
                createdTime: "2023-01-17T11:23:49Z"
            },
            headers: {
                location: locationUrl
            }
        };
        CreatedAPI createdAPIWithOperationLevelInterceptorPolicy = {
            body: {
                id: k8sAPIUUID1,
                name: "PizzaAPI",
                context: "/pizzaAPI/1.0.0",
                'version: "1.0.0",
                'type: "REST",
                operations: [
                    {target: "/*", verb: "GET", authTypeEnabled: true, scopes: [], operationPolicies: {request: [check requestInterceptorPolicy.cloneWithType(OperationPolicy)], response: [check responseInterceptorPolicy.cloneWithType(OperationPolicy)]}},
                    {target: "/*", verb: "PUT", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "POST", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "DELETE", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PATCH", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}}
                ],
                serviceInfo: {name: "backend", namespace: "apk"},
                createdTime: "2023-01-17T11:23:49Z"
            },
            headers: {
                location: locationUrl
            }
        };
        CreatedAPI createdAPIWithAPILevelInterceptorPolicy = {
            body: {
                id: k8sAPIUUID1,
                name: "PizzaAPI",
                context: "/pizzaAPI/1.0.0",
                'version: "1.0.0",
                'type: "REST",
                operations: [
                    {target: "/*", verb: "GET", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PUT", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "POST", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "DELETE", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}},
                    {target: "/*", verb: "PATCH", authTypeEnabled: true, scopes: [], operationPolicies: {request: [], response: []}}
                ],
                apiPolicies: {request: [check requestInterceptorPolicy.cloneWithType(OperationPolicy)], response: [check responseInterceptorPolicy.cloneWithType(OperationPolicy)]},
                serviceInfo: {name: "backend", namespace: "apk"},
                createdTime: "2023-01-17T11:23:49Z"
            },
            headers: {
                location: locationUrl
            }
        };

        map<[string, string, [model:ConfigMap, any], [model:Httproute, any], [model:K8sServiceMapping, any], [model:API, any], [model:RuntimeAPI, any], API, string, [model:Backend, any][], [model:RateLimitPolicy|(), any], [model:APIPolicy|(), any], [model:InterceptorService, any][], anydata]> data = {
            "1": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], api, k8sAPIUUID1, services, [(), ()], [(), ()], [], createdAPI.toBalString()],
            "2": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], alreadyNameExist, k8sAPIUUID1, services, [(), ()], [(), ()], [], nameAlreadyExistError.toBalString()],
            "3": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], contextAlreadyExist, k8sAPIUUID1, services, [(), ()], [(), ()], [], contextAlreadyExistError.toBalString()],
            "4": ["275b00d1-722c-4df2-b65a-9b14677abe4a", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], api, k8sAPIUUID1, services, [(), ()], [(), ()], [], serviceNotExist.toBalString()],
            "5": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRouteWithPolicies, httpRouteWithPoliciesResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAPIWithPolicies, mockRuntimeResponseWithPolicies],
                check apiWithOperationPolicies.cloneWithType(API),
                k8sAPIUUID1,
                services,
                [(), ()],
                [(), ()],
                [],
                createdAPIWithPolicies.toBalString()
            ],
            "6": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRouteWithPolicies, httpRouteWithPoliciesResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAPIWithPolicies, mockRuntimeResponseWithPolicies],
                apiWithInvalidPolicyName,
                k8sAPIUUID1,
                services,
                [(), ()],
                [(), ()],
                [],
                invalidPolicyNameError.toBalString()
            ],
            "7": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRouteWithOperationRateLimits, httpRouteWithOperationRateLimitsResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAOperationRateLimits, mockRuntimeResponseWithOperationRateLimits],
                apiWithOperationRateLimits,
                k8sAPIUUID1,
                services,
                [getMockResourceRateLimitPolicy(apiWithOperationRateLimits, organiztion1, apiUUID), getMockRateLimitResponse(getMockResourceRateLimitPolicy(apiWithOperationRateLimits, organiztion1, apiUUID).clone())],
                [(), ()],
                [],
                createdAPIWithOperationRateLimits.toBalString()
            ],
            "8": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRoute, httpRouteResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAPIWithAPIRateLimits, mockRuntimeResponseWithAPIRateLimits],
                apiWithAPIRateLimits,
                k8sAPIUUID1,
                services,
                [getMockAPIRateLimitPolicy(apiWithAPIRateLimits, organiztion1, apiUUID), getMockRateLimitResponse(getMockAPIRateLimitPolicy(apiWithAPIRateLimits, organiztion1, apiUUID).clone())],
                [(), ()],
                [],
                createdAPIWithAPIRateLimits.toBalString()
            ],
            "9": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRouteWithPolicies, httpRouteWithPoliciesResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAPIWithPolicies, mockRuntimeResponseWithPolicies],
                apiWithBothRateLimits,
                k8sAPIUUID1,
                services,
                [getMockAPIRateLimitPolicy(apiWithAPIRateLimits, organiztion1, apiUUID), getMockRateLimitResponse(getMockAPIRateLimitPolicy(apiWithAPIRateLimits, organiztion1, apiUUID).clone())],
                [(), ()],
                [],
                bothRateLimitsPresentError.toBalString()
            ],
            "10": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRouteWithInterceptorPolicy, httpRouteWithInterceptorPolicyResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAPIWithOperationInterceptorPolicy, mockRuntimeResponseWithOperationInterceptorPolicy],
                check apiWithOperationLevelInterceptorPolicy.cloneWithType(API),
                k8sAPIUUID1,
                services,
                [(), ()],
                [getMockResourceLevelPolicy(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID), getMockAPIPolicyResponse(getMockResourceLevelPolicy(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID).clone())],
                interceptorServices,
                createdAPIWithOperationLevelInterceptorPolicy.toBalString()
            ],
            "11": [
                "275b00d1-722c-4df2-b65a-9b14677abe4b",
                apiUUID,
                [configmap, mockConfigMapResponse],
                [httpRoute, httpRouteResponse],
                [mockServiceMappingRequest, serviceMappingResponse],
                [mockAPI, mockAPIResponse],
                [mockRuntimeAPIWithAPIInterceptorPolicy, mockRuntimeResponseWithAPIInterceptorPolicy],
                check apiWithAPILevelInterceptorPolicy.cloneWithType(API),
                k8sAPIUUID1,
                services,
                [(), ()],
                [getMockAPILevelPolicy(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID), getMockAPIPolicyResponse(getMockAPILevelPolicy(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID).clone())],
                interceptorServices,
                createdAPIWithAPILevelInterceptorPolicy.toBalString()
            ]
        };
        return data;
    } on fail var e {
        test:assertFail(msg = e.message());
    }
}

function getMockAPIResponse(model:API api, string k8sAPIUUID) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    api.metadata.uid = k8sAPIUUID;
    api.metadata.creationTimestamp = "2023-01-17T11:23:49Z";
    response.setJsonPayload(api.toJson());
    return response;
}

function getMockAPIResponse1(API api, string apiUUID, string k8sAPIUUID) returns http:Response {
    http:Response response = new;
    response.statusCode = 400;
    model:Status status = {code: 400, message: ""};
    response.setJsonPayload(status.toJson());
    return response;
}

function getMockAPI(API api, string apiUUID, string organization) returns model:API {
    APIClient apiClient = new;
    model:EnvConfig[]? envConfig = [{httpRouteRefs: ["http-route-ref-name"]}];
    model:API k8sapi = {
        "kind": "API",
        "apiVersion": "dp.wso2.com/v1alpha1",
        "metadata": {"name": apiUUID, "namespace": "apk-platform", "labels": getLabels(api, organiztion1)},
        "spec": {
            "apiDisplayName": api.name,
            "apiType": "REST",
            "apiVersion": api.'version,
            "apiProvider": "apkUser",
            "context": apiClient.returnFullContext(api.context, api.'version),
            "organization": organization,
            "definitionFileRef": apiUUID + "-definition",
            "production": envConfig
        },
        "status"
                : null
    };
    return k8sapi;
}

function getMockRuntimeAPI(API api, string apiUUID, commons:Organization organization, Service? serviceEntry) returns model:RuntimeAPI {
    APIClient apiClient = new;
    string userName = "apkUser";
    model:RuntimeAPI runtimeAPI = apiClient.generateRuntimeAPIArtifact(api, serviceEntry, organization, userName);
    if api.operations is () {
        runtimeAPI.spec.operations = [
            {
                target: "/*",
                verb: "GET",
                authTypeEnabled: true,
                operationPolicies: {request: [], response: []}
            },
            {
                target: "/*",
                verb: "PUT",
                authTypeEnabled: true,
                operationPolicies: {request: [], response: []}
            },
            {
                target: "/*",
                verb: "POST",
                authTypeEnabled: true,
                operationPolicies: {request: [], response: []}
            },
            {
                target: "/*",
                verb: "DELETE",
                authTypeEnabled: true,
                operationPolicies: {request: [], response: []}
            },
            {
                target: "/*",
                verb: "PATCH",
                authTypeEnabled: true,
                operationPolicies: {request: [], response: []}
            }
        ];
    }
    return runtimeAPI;
}

function getMockRuntimeAPIResponse(model:RuntimeAPI runtimeApi) returns http:Response {
    http:Response response = new ();
    response.statusCode = 201;
    runtimeApi.metadata.uid = "275b00d1-722c-4df2-b65a-9b14677abe4b";
    response.setJsonPayload(runtimeApi.toJson());
    return response;
}

function getMockAPI1(API api, string apiUUID, string organization) returns model:API {
    APIClient apiClient = new;
    model:EnvConfig[]? envConfig = [{httpRouteRefs: ["http-route-ref-name"]}];
    model:API k8sapi = {
        "kind": "API",
        "apiVersion": "dp.wso2.com/v1alpha1",
        "metadata": {"name": apiUUID, "namespace": "apk-platform", "labels": getLabels(api, organiztion1)},
        "spec": {
            "apiDisplayName": api.name,
            "apiType": "REST",
            "apiVersion": api.'version,
            "apiProvider": "apkUser",
            "context": apiClient.returnFullContext(api.context, api.'version),
            "organization": organization,
            "definitionFileRef": apiUUID + "-definition",
            "sandbox": envConfig
        },
        "status": null
    };
    return k8sapi;
}

function getMockServiceMappingRequest(API api, string apiUUID) returns model:K8sServiceMapping {
    model:K8sServiceMapping serviceMapping = {"kind": "ServiceMapping", "apiVersion": "dp.wso2.com/v1alpha1", "metadata": {"name": apiUUID + "-servicemapping", "namespace": "apk-platform", "labels": getLabels(api, organiztion1)}, "spec": {"serviceRef": {"name": "backend", "namespace": "apk"}, "apiRef": {"name": apiUUID, "namespace": "apk-platform"}}};
    return serviceMapping;
}

function getMockServiceMappingResponse(model:K8sServiceMapping serviceMapping) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    serviceMapping.metadata.uid = uuid:createType1AsString();
    response.setJsonPayload(serviceMapping.toJson());
    return response;
}

function getMockHttpRoute(API api, string apiUUID, commons:Organization organiztion) returns model:Httproute {
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion1)},
        "spec": {
            "hostnames": [string:concat(organiztion.uuid, ".", "gw.wso2.com")],
            "rules": [
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "GET"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "PUT"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "POST"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "DELETE"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "PATCH"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithOperationPolicies1(API api, string apiUUID, commons:Organization organiztion) returns model:Httproute {
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "hostnames": [
                string:concat(organiztion.uuid, ".", "gw.wso2.com")
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithOperationRateLimits1(API api, string apiUUID, commons:Organization organiztion) returns model:Httproute {
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "hostnames": [
                string:concat(organiztion.uuid, ".", "gw.wso2.com")
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "ExtensionRef",
                            "extensionRef": {
                                "group": "dp.wso2.com",
                                "kind": "RateLimitPolicy",
                                "name": "rate-limit-policy-ref-name"
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithOperationInterceptorPolicy1(API api, string apiUUID, commons:Organization organiztion) returns model:Httproute {
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "hostnames": [
                string:concat(organiztion.uuid, ".", "gw.wso2.com")
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "ExtensionRef",
                            "extensionRef": {
                                "group": "dp.wso2.com",
                                "kind": "APIPolicy",
                                "name": "api-policy-ref-name"
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion),
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockConfigMap1(string apiUniqueId, API api) returns model:ConfigMap|error {
    json content = {"openapi": "3.0.1", "info": {"title": "" + api.name + "", "version": "" + api.'version + ""}, "security": [{"default": []}], "paths": {"/*": {"get": {"responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-auth-type": true, "x-throttling-tier": "Unlimited"}, "put": {"responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-auth-type": true, "x-throttling-tier": "Unlimited"}, "post": {"responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-auth-type": true, "x-throttling-tier": "Unlimited"}, "delete": {"responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-auth-type": true, "x-throttling-tier": "Unlimited"}, "patch": {"responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-auth-type": true, "x-throttling-tier": "Unlimited"}}}, "components": {"securitySchemes": {"default": {"type": "oauth2", "flows": {"implicit": {"authorizationUrl": "https://test.com", "scopes": {}}}}}}}
;
    string base64EncodedGzipContent = check getBase64EncodedGzipContent(content.toJsonString().toBytes());
    model:ConfigMap configmap = {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "labels": getLabels(api, organiztion1),
            "name": apiUniqueId + "-definition",
            "namespace": "apk-platform"
        }
    };
    configmap.binaryData = {
        [CONFIGMAP_DEFINITION_KEY] : base64EncodedGzipContent
    };
    return configmap;
}

public function getBase64EncodedGzipContent(byte[] content) returns string|error {
    byte[]|io:IOException gzipUtilCompressGzipFile = check commons:GzipUtil_compressGzipFile(content);
    if gzipUtilCompressGzipFile is byte[] {
        byte[] encoderUtilEncodeBase64 = check commons:EncoderUtil_encodeBase64(gzipUtilCompressGzipFile);
        return string:fromBytes(encoderUtilEncodeBase64);
    } else {
        return error("Error while encoding the content");
    }

}

function getMockConfigMapResponse(model:ConfigMap configmap) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    configmap.metadata.uid = uuid:createType1AsString();
    response.setJsonPayload(configmap.toJson());
    return response;
}

function getMockConfigMapErrorResponse() returns http:Response {
    http:Response response = new;
    response.statusCode = 400;
    model:Status status = {code: 400, message: "configmap already exist"};
    response.setJsonPayload(status.toJson());
    return response;
}

@test:Config {dataProvider: createAPIDataProvider}
function testCreateAPI(string apiUUID, string backenduuid, API api, model:ConfigMap configmap,
        any configmapDeployingResponse, model:Httproute? prodhttpRoute,
        any prodhttpResponse, model:Httproute? sandHttpRoute, any sandhttpResponse,
        [model:Backend, any][] backendServices,
        model:API k8sApi, any k8sapiResponse, model:RuntimeAPI runtimeAPI, any runtimeAPIResponse,
        model:RateLimitPolicy? rateLimitPolicy, any rateLimitPolicyResponse,
        model:APIPolicy? apiPolicy, any apiPolicyResponse,
        [model:InterceptorService, any][] interceptorServices
, string k8sapiUUID, anydata expected) returns error? {
    APIClient apiClient = new;
    string userName = "apkUser";
    test:prepare(k8sApiServerEp).when("post").withArguments("/api/v1/namespaces/apk-platform/configmaps", configmap).thenReturn(configmapDeployingResponse);
    if prodhttpRoute is model:Httproute {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", prodhttpRoute).thenReturn(prodhttpResponse);
    }
    if sandHttpRoute is model:Httproute {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", sandHttpRoute).thenReturn(sandhttpResponse);
    }
    if rateLimitPolicy is model:RateLimitPolicy {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/ratelimitpolicies", rateLimitPolicy).thenReturn(rateLimitPolicyResponse);
    }
    if apiPolicy is model:APIPolicy {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apipolicies", apiPolicy).thenReturn(apiPolicyResponse);
    }
    foreach [model:InterceptorService, any] [interceptorService, interceptorServiceResponse] in interceptorServices {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/interceptorservices", interceptorService).thenReturn(interceptorServiceResponse);
    }
    foreach [model:Backend, any] backend in backendServices {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backends", backend[0]).thenReturn(backend[1]);
    }
    http:Response configmapResponse = new;
    configmapResponse.statusCode = 404;
    http:ApplicationResponseError internalApiResponse = error("internal api not found", statusCode = 404, body = {}, headers = {});
    http:Response internalAPIDeletionResponse = new;
    internalAPIDeletionResponse.statusCode = 200;
    model:HttprouteList httpRouteList = {metadata: {}, items: []};
    model:ServiceMappingList serviceMappingList = {metadata: {}, items: []};
    model:AuthenticationList authenticationList = {metadata: {}, items: []};
    model:BackendList serviceList = {metadata: {}, items: []};
    model:ScopeList scopeList = {metadata: {}, items: []};
    model:RateLimitPolicyList rateLimitPolicyList = {metadata: {}, items: []};
    model:APIPolicyList apiPolicyList = {metadata: {}, items: []};
    model:InterceptorServiceList interceptorServiceList = {metadata: {}, items: []};
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/ratelimitpolicies?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(rateLimitPolicyList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apipolicies?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(apiPolicyList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/interceptorservices?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(interceptorServiceList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/scopes?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(scopeList);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis", k8sApi).thenReturn(k8sapiResponse);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis", runtimeAPI).thenReturn(runtimeAPIResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/" + apiClient.retrieveDefinitionName(apiUUID)).thenReturn(configmapResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes/?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(httpRouteList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(serviceMappingList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/authentications?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(authenticationList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backends?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version, organiztion1)).thenReturn(serviceList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sApi.metadata.name).thenReturn(internalApiResponse);
    test:prepare(k8sApiServerEp).when("delete").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sApi.metadata.name).thenReturn(internalAPIDeletionResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sApi.metadata.name).thenReturn(runtimeAPI);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/" + k8sApi.metadata.name).thenReturn(configmapResponse);
    any|error aPI = apiClient.createAPI(api, (), organiztion1, userName);
    if aPI is any {
        test:assertEquals(aPI.toBalString(), expected);
    } else {
        test:assertEquals(aPI.toBalString(), expected);
    }
}

@test:Config {dataProvider: createAPIRateLimitPolicyProvider}
function testCreateAPIWithRatelimitPolicy(string apiUUID, string backenduuid, API api, model:ConfigMap configmap,
        any configmapDeployingResponse, model:Httproute? prodhttpRoute,
        any prodhttpResponse, model:Httproute? sandHttpRoute, any sandhttpResponse,
        [model:Backend, any][] backendServices,
        model:API k8sApi, any k8sapiResponse, model:RuntimeAPI runtimeAPI, any runtimeAPIResponse,
        model:RateLimitPolicy? rateLimitPolicy, any rateLimitPolicyResponse,
        model:APIPolicy? apiPolicy, any apiPolicyResponse,
        [model:InterceptorService, any][] interceptorServices
, string k8sapiUUID, anydata expected) returns error? {
    return testCreateAPI(apiUUID, backenduuid, api, configmap, configmapDeployingResponse, prodhttpRoute, prodhttpResponse,
            sandHttpRoute, sandhttpResponse, backendServices, k8sApi, k8sapiResponse, runtimeAPI, runtimeAPIResponse,
            rateLimitPolicy, rateLimitPolicyResponse, apiPolicy, apiPolicyResponse, interceptorServices,
            k8sapiUUID, expected);
}

function createAPIRateLimitPolicyProvider() returns map<[string, string, API, model:ConfigMap, any, model:Httproute?, any, model:Httproute?, any, [model:Backend, any][], model:API, any, model:RuntimeAPI, any, model:RateLimitPolicy?, any, model:APIPolicy?, any, [model:InterceptorService, any][], string, anydata]> {
    do {
        API api = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0",
            endpointConfig: {"production_endpoints": {"url": "https://localhost"}}
        };

        API apiWithAPIRateLimits = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ],
            "apiRateLimit": {
                "requestsPerUnit": 10,
                "unit": "Minute"
            }
        };
        API apiWithBothRateLimits = {
            "name": "PizzaAPIPolicies",
            "context": "/pizzaAPIPolcies/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/menu",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationRateLimit": {
                        "requestsPerUnit": 10,
                        "unit": "Minute"
                    }
                }
            ],
            "apiRateLimit": {
                "requestsPerUnit": 10,
                "unit": "Minute"
            }
        };
        API apiWithOperationRateLimits = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationRateLimit": {
                        "requestsPerUnit": 10,
                        "unit": "Minute"
                    }
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ]
        };
        commons:APKError bothRateLimitsPresentError = error commons:APKError("Presence of both resource level and API level rate limits is not allowed",
            code = 909026,
            message = "Presence of both resource level and API level rate limits is not allowed",
            statusCode = 406,
            description = "Presence of both resource level and API level rate limits is not allowed"
        );
        string apiUUID = getUniqueIdForAPI(api.name, api.'version, organiztion1);
        string backenduuid = getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion1);
        string backenduuid1 = getBackendServiceUid(api, (), SANDBOX_TYPE, organiztion1);
        string interceptorBackenduuid1 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend1.interceptor:9082");
        string interceptorBackenduuid2 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend2.interceptor:9083");
        string k8sapiUUID = uuid:createType1AsString();
        model:Backend backendService = {
            metadata: {name: backenduuid, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "localhost", port: 443}], protocol: "https"}
        };
        model:Backend backendService1 = {
            metadata: {name: backenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "localhost", port: 443}], protocol: "https"}
        };
        model:Backend interceptorBackendService1 = {
            metadata: {name: interceptorBackenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend1.interceptor", port: 9082}], protocol: "http"}
        };
        model:Backend interceptorBackendService2 = {
            metadata: {name: interceptorBackenduuid2, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend2.interceptor", port: 9083}], protocol: "http"}
        };
        http:Response backendServiceResponse = getOKBackendServiceResponse(backendService);
        http:Response backendServiceResponse1 = getOKBackendServiceResponse(backendService);
        http:Response interceptorBackendServiceResponse1 = getOKBackendServiceResponse(interceptorBackendService1);
        http:Response interceptorBackendServiceResponse2 = getOKBackendServiceResponse(interceptorBackendService2);
        http:Response backendServiceErrorResponse = new;
        backendServiceErrorResponse.statusCode = 403;
        [model:Backend, any][] services = [];
        services.push([backendService, backendServiceResponse]);
        services.push([backendService1, backendServiceResponse1]);
        services.push([interceptorBackendService1, interceptorBackendServiceResponse1]);
        services.push([interceptorBackendService2, interceptorBackendServiceResponse2]);
        [model:Backend, any][] servicesError = [];
        servicesError.push([backendService, backendServiceErrorResponse]);
        model:ConfigMap configmap = check getMockConfigMap1(apiUUID, api);
        model:Httproute prodhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        model:Httproute prodhttpRouteWithOperationRateLimits = getMockHttpRouteWithOperationRateLimits(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        string locationUrl = runtimeConfiguration.baseURl + "/apis/" + k8sapiUUID;

        CreatedAPI CreatedAPIWithOperationRateLimits = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}, "operationRateLimit": {"requestsPerUnit": 10, "unit": "Minute"}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}]}, headers: {location: locationUrl}};
        CreatedAPI CreatedAPIWithAPIRateLimits = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}], apiRateLimit: {"requestsPerUnit": 10, "unit": "Minute"}}, headers: {location: locationUrl}};

        map<[string, string, API, model:ConfigMap,
    any, model:Httproute|(), any, model:Httproute|(),
    any, [model:Backend, any][], model:API, any, model:RuntimeAPI, any,
    model:RateLimitPolicy|(), any, model:APIPolicy|(), any, [model:InterceptorService, any][], string,
    anydata]> data = {
            "1": [
                apiUUID,
                backenduuid,
                apiWithBothRateLimits,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIErrorNameExist(),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                bothRateLimitsPresentError.toBalString()
            ]
        ,
            "2": [
                apiUUID,
                backenduuid,
                apiWithOperationRateLimits,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRouteWithOperationRateLimits,
                getMockHttpRouteResponse(prodhttpRouteWithOperationRateLimits.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(apiWithOperationRateLimits, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(apiWithOperationRateLimits, apiUUID, organiztion1, ())),
                getMockResourceRateLimitPolicy(apiWithOperationRateLimits, organiztion1, apiUUID),
                getMockRateLimitResponse(getMockResourceRateLimitPolicy(apiWithOperationRateLimits, organiztion1, apiUUID).clone()),
                (),
                (),
                [],
                k8sapiUUID,
                CreatedAPIWithOperationRateLimits.toBalString()
            ]
        ,
            "3": [
                apiUUID,
                backenduuid,
                apiWithAPIRateLimits,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(apiWithAPIRateLimits, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(apiWithAPIRateLimits, apiUUID, organiztion1, ())),
                getMockAPIRateLimitPolicy(apiWithAPIRateLimits, organiztion1, apiUUID),
                getMockRateLimitResponse(getMockAPIRateLimitPolicy(apiWithAPIRateLimits, organiztion1, apiUUID).clone()),
                (),
                (),
                [],
                k8sapiUUID,
                CreatedAPIWithAPIRateLimits.toBalString()
            ]
        };
        return data;
    } on fail var e {
        test:assertFail("tests failed===" + e.message());
    }
}

@test:Config {dataProvider: createAPIWithOperationPolicyProvider}
function testCreateAPIWithOperationPolicy(string apiUUID, string backenduuid, API api, model:ConfigMap configmap,
        any configmapDeployingResponse, model:Httproute? prodhttpRoute,
        any prodhttpResponse, model:Httproute? sandHttpRoute, any sandhttpResponse,
        [model:Backend, any][] backendServices,
        model:API k8sApi, any k8sapiResponse, model:RuntimeAPI runtimeAPI, any runtimeAPIResponse,
        model:RateLimitPolicy? rateLimitPolicy, any rateLimitPolicyResponse,
        model:APIPolicy? apiPolicy, any apiPolicyResponse,
        [model:InterceptorService, any][] interceptorServices
, string k8sapiUUID, anydata expected) returns error? {
    return testCreateAPI(apiUUID, backenduuid, api, configmap, configmapDeployingResponse, prodhttpRoute, prodhttpResponse,
            sandHttpRoute, sandhttpResponse, backendServices, k8sApi, k8sapiResponse, runtimeAPI, runtimeAPIResponse,
            rateLimitPolicy, rateLimitPolicyResponse, apiPolicy, apiPolicyResponse, interceptorServices,
            k8sapiUUID, expected);
}

function createAPIWithOperationPolicyProvider() returns map<[string, string, API, model:ConfigMap, any, model:Httproute?, any, model:Httproute?, any, [model:Backend, any][], model:API, any, model:RuntimeAPI, any, model:RateLimitPolicy?, any, model:APIPolicy?, any, [model:InterceptorService, any][], string, anydata]> {
    do {
        API api = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0",
            endpointConfig: {"production_endpoints": {"url": "https://localhost"}}
        };
        json apiWithOperationPolicies = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "addHeader",
                                "parameters":
                                {
                                    "headerName": "customadd",
                                    "headerValue": "customvalue"
                                }

                            }
                        ],
                        "response": [
                            {
                                "policyName": "removeHeader",
                                "parameters":
                                {
                                    "headerName": "content-length"
                                }

                            }
                        ]
                    }
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ]
        };
        json apiWithAPIPolicies = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    target: "/*",
                    verb: "GET",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "PUT",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "POST",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "DELETE",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "PATCH",
                    authTypeEnabled: true
                }
            ],
            "apiPolicies": {
                request: [
                    {
                        policyName: "addHeader",
                        "parameters":
                        {
                            "headerName": "customadd",
                            "headerValue": "customvalue"
                        }

                    }
                ],
                "response": [
                    {
                        policyName: "removeHeader",
                        "parameters":
                        {
                            "headerName": "content-length"
                        }

                    }
                ]
            }
        };
        json apiWithBothPolicies = {
            "name": "PizzaAPIPolicies",
            "context": "/pizzaAPIPolcies/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/menu",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "addHeader",
                                "parameters":
                                {
                                    "headerName": "customadd",
                                    "headerValue": "customvalue"
                                }
                            }
                        ]
                    }
                }
            ],
            "apiPolicies": {
                "request": [
                    {
                        "policyName": "addHeader",
                        "parameters":
                        {
                            "headerName": "customadd",
                            "headerValue": "customvalue"
                        }

                    }
                ]
            }
        };
        commons:APKError bothPoliciesPresentError = error commons:APKError("Presence of both resource level and API level operation policies is not allowed",
            code = 909025,
            message = "Presence of both resource level and API level operation policies is not allowed",
            statusCode = 406,
            description = "Presence of both resource level and API level operation policies is not allowed"
        );
        json apiWithInvalidPolicyName = {
            "name": "PizzaAPIOps",
            "context": "/pizzaAPIOps/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/menu",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "addHeader1",
                                "parameters":
                                {
                                    "headerName": "customadd",
                                    "headerValue": "customvalue"
                                }

                            }
                        ]
                    }
                }
            ]
        };
        commons:APKError invalidPolicyNameError = error commons:APKError("Invalid operation policy name",
            code = 909010,
            message = "Invalid operation policy name",
            statusCode = 406,
            description = "Invalid operation policy name"
        );
        json apiWithInvalidPolicyParameters = {
            "name": "PizzaAPIOps",
            "context": "/pizzaAPIOps/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/menu",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "addHeader",
                                "parameters":
                                {
                                    "headerName1": "customadd",
                                    "headerValue": "customvalue"
                                }

                            }
                        ]
                    }
                }
            ]
        };
        commons:APKError invalidPolicyParametersError = error commons:APKError("Invalid parameters provided for policy addHeader",
            code = 909024,
            message = "Invalid parameters provided for policy addHeader",
            statusCode = 406,
            description = "Invalid parameters provided for policy addHeader"
        );
        json apiWithOperationLevelInterceptorPolicy = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    "target": "/*",
                    "verb": "GET",
                    "authTypeEnabled": true,
                    "operationPolicies": {
                        "request": [
                            {
                                "policyName": "Interceptor",
                                "parameters":
                                {
                                    "headersEnabled": true,
                                    "bodyEnabled": false,
                                    "trailersEnabled": false,
                                    "contextEnabled": true,
                                    "backendUrl": "http://interceptor-backend1.interceptor:9082"
                                }
                            }
                        ],
                        "response": [
                            {
                                "policyName": "Interceptor",
                                "parameters":
                                {
                                    "headersEnabled": false,
                                    "bodyEnabled": true,
                                    "trailersEnabled": false,
                                    "contextEnabled": true,
                                    "backendUrl": "http://interceptor-backend2.interceptor:9083"
                                }
                            }
                        ]
                    }
                },
                {
                    "target": "/*",
                    "verb": "PUT",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "POST",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "DELETE",
                    "authTypeEnabled": true
                },
                {
                    "target": "/*",
                    "verb": "PATCH",
                    "authTypeEnabled": true
                }
            ]
        };
        json apiWithAPILevelInterceptorPolicy = {
            "name": "PizzaAPI",
            "context": "/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "endpointConfig": {"production_endpoints": {"url": "https://localhost"}},
            "operations": [
                {
                    target: "/*",
                    verb: "GET",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "PUT",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "POST",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "DELETE",
                    authTypeEnabled: true
                },
                {
                    target: "/*",
                    verb: "PATCH",
                    authTypeEnabled: true
                }
            ],
            "apiPolicies": {
                "request": [
                    {
                        "policyName": "Interceptor",
                        "parameters":
                        {
                            "headersEnabled": true,
                            "bodyEnabled": false,
                            "trailersEnabled": false,
                            "contextEnabled": true,
                            "backendUrl": "http://interceptor-backend1.interceptor:9082"
                        }
                    }
                ],
                "response": [
                    {
                        "policyName": "Interceptor",
                        "parameters":
                        {
                            "headersEnabled": false,
                            "bodyEnabled": true,
                            "trailersEnabled": false,
                            "contextEnabled": true,
                            "backendUrl": "http://interceptor-backend2.interceptor:9083"
                        }
                    }
                ]
            }
        };
        string apiUUID = getUniqueIdForAPI(api.name, api.'version, organiztion1);
        string backenduuid = getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion1);
        string backenduuid1 = getBackendServiceUid(api, (), SANDBOX_TYPE, organiztion1);
        string interceptorBackenduuid1 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend1.interceptor:9082");
        string interceptorBackenduuid2 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend2.interceptor:9083");
        string k8sapiUUID = uuid:createType1AsString();
        model:Backend backendService = {
            metadata: {name: backenduuid, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "localhost", port: 443}], protocol: "https"}
        };
        model:Backend backendService1 = {
            metadata: {name: backenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "localhost", port: 443}], protocol: "https"}
        };
        model:Backend interceptorBackendService1 = {
            metadata: {name: interceptorBackenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend1.interceptor", port: 9082}], protocol: "http"}
        };
        model:Backend interceptorBackendService2 = {
            metadata: {name: interceptorBackenduuid2, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend2.interceptor", port: 9083}], protocol: "http"}
        };
        http:Response backendServiceResponse = getOKBackendServiceResponse(backendService);
        http:Response backendServiceResponse1 = getOKBackendServiceResponse(backendService);
        http:Response interceptorBackendServiceResponse1 = getOKBackendServiceResponse(interceptorBackendService1);
        http:Response interceptorBackendServiceResponse2 = getOKBackendServiceResponse(interceptorBackendService2);
        http:Response backendServiceErrorResponse = new;
        backendServiceErrorResponse.statusCode = 403;
        [model:Backend, any][] services = [];
        services.push([backendService, backendServiceResponse]);
        services.push([backendService1, backendServiceResponse1]);
        services.push([interceptorBackendService1, interceptorBackendServiceResponse1]);
        services.push([interceptorBackendService2, interceptorBackendServiceResponse2]);
        [model:Backend, any][] servicesError = [];
        servicesError.push([backendService, backendServiceErrorResponse]);
        
        [model:InterceptorService, any][] interceptorServices = [];
        string interceptorBackendUrl1 =  "http://interceptor-backend1.interceptor:9082";
        string interceptorBackendUrl2 =  "http://interceptor-backend2.interceptor:9083";
        string[] requestIncludes = ["request_headers", "invocation_context"];
        string[] responseIncludes = ["response_body", "invocation_context"];
        model:InterceptorService requestInterceptorService = getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "request", requestIncludes, interceptorBackendUrl1);
        http:Response requestInterceptorServiceResponse = getMockInterceptorServiceResponse(getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "request", requestIncludes, interceptorBackendUrl1).clone());
        interceptorServices.push([requestInterceptorService, requestInterceptorServiceResponse]);
        model:InterceptorService responseInterceptorService = getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "response", responseIncludes, interceptorBackendUrl2);
        http:Response responseInterceptorServiceResponse = getMockInterceptorServiceResponse(getMockInterceptorService(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID, "response", responseIncludes, interceptorBackendUrl2).clone());
        interceptorServices.push([responseInterceptorService, responseInterceptorServiceResponse]);

        model:ConfigMap configmap = check getMockConfigMap1(apiUUID, api);
        model:Httproute prodhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        model:Httproute prodhttpRouteWithOperationPolicies = getMockHttpRouteWithOperationPolicies(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        model:Httproute prodhttpRouteWithAPIPolicies = getMockHttpRouteWithAPIPolicies(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        model:Httproute prodhttpRouteWithOperationInterceptorPolicy = getMockHttpRouteWithOperationInterceptorPolicy(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        string locationUrl = runtimeConfiguration.baseURl + "/apis/" + k8sapiUUID;

        CreatedAPI CreatedAPIWithOperationPolicies = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [{"policyName": "addHeader", "policyVersion": "v1", "parameters": {"headerName": "customadd", "headerValue": "customvalue"}}], "response": [{"policyName": "removeHeader", "policyVersion": "v1", "parameters": {"headerName": "content-length"}}]}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}]}, headers: {location: locationUrl}};
        CreatedAPI CreatedAPIWithAPIPolicies = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}], apiPolicies: {"request": [{"policyName": "addHeader", "policyVersion": "v1", "parameters": {"headerName": "customadd", "headerValue": "customvalue"}}], "response": [{"policyName": "removeHeader", "policyVersion": "v1", "parameters": {"headerName": "content-length"}}]}}, headers: {location: locationUrl}};
        CreatedAPI CreatedAPIWithOperationLevelInterceptorPolicy = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [{"policyName": "Interceptor", "policyVersion": "v1", "parameters": {"headersEnabled": true, "bodyEnabled": false, "contextEnabled": true, "backendUrl": "http://interceptor-backend1.interceptor:9082", "trailersEnabled": false}}], "response": [{"policyName": "Interceptor", "policyVersion": "v1", "parameters": {"headersEnabled": false, "bodyEnabled": true, "contextEnabled": true, "backendUrl": "http://interceptor-backend2.interceptor:9083", "trailersEnabled": false}}]}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}]}, headers: {location: locationUrl}};
        CreatedAPI CreatedAPIWithAPILevelInterceptorPolicy = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}], apiPolicies: {"request": [{"policyName": "Interceptor", "policyVersion": "v1", "parameters": {"headersEnabled": true, "bodyEnabled": false, "contextEnabled": true, "backendUrl": "http://interceptor-backend1.interceptor:9082", "trailersEnabled": false}}], "response": [{"policyName": "Interceptor", "policyVersion": "v1", "parameters": {"headersEnabled": false, "bodyEnabled": true, "contextEnabled": true, "backendUrl": "http://interceptor-backend2.interceptor:9083", "trailersEnabled": false}}]}}, headers: {location: locationUrl}};

        map<[string, string, API, model:ConfigMap,
    any, model:Httproute|(), any, model:Httproute|(),
    any, [model:Backend, any][], model:API, any, model:RuntimeAPI, any,
    model:RateLimitPolicy|(), any, model:APIPolicy|(), any, [model:InterceptorService, any][], string,
    anydata]> data = {

            "1": [
                apiUUID,
                backenduuid,
                check apiWithBothPolicies.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIErrorNameExist(),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                bothPoliciesPresentError.toBalString()
            ]
        ,
            "2": [
                apiUUID,
                backenduuid,
                check apiWithInvalidPolicyName.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIErrorNameExist(),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                invalidPolicyNameError.toBalString()
            ]
        ,
            "3": [
                apiUUID,
                backenduuid,
                check apiWithInvalidPolicyParameters.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIErrorNameExist(),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                invalidPolicyParametersError.toBalString()
            ]
        ,
            "4": [
                apiUUID,
                backenduuid,
                check apiWithOperationPolicies.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRouteWithOperationPolicies,
                getMockHttpRouteResponse(prodhttpRouteWithOperationPolicies.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(check apiWithOperationPolicies.cloneWithType(API), apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(check apiWithOperationPolicies.cloneWithType(API), apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                CreatedAPIWithOperationPolicies.toBalString()
            ]
        ,
            "5": [
                apiUUID,
                backenduuid,
                check apiWithAPIPolicies.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRouteWithAPIPolicies,
                getMockHttpRouteResponse(prodhttpRouteWithAPIPolicies.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(check apiWithAPIPolicies.cloneWithType(API), apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(check apiWithAPIPolicies.cloneWithType(API), apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                CreatedAPIWithAPIPolicies.toBalString()
            ]
        ,
            "6": [
                apiUUID,
                backenduuid,
                check apiWithOperationLevelInterceptorPolicy.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRouteWithOperationInterceptorPolicy,
                getMockHttpRouteResponse(prodhttpRouteWithOperationInterceptorPolicy.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), apiUUID, organiztion1, ())),
                (),
                (),
                getMockResourceLevelPolicy(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID),
                getMockAPIPolicyResponse(getMockResourceLevelPolicy(check apiWithOperationLevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID).clone()),
                interceptorServices,
                k8sapiUUID,
                CreatedAPIWithOperationLevelInterceptorPolicy.toBalString()
            ]
        ,
            "7": [
                apiUUID,
                backenduuid,
                check apiWithAPILevelInterceptorPolicy.cloneWithType(API),
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), apiUUID, organiztion1, ())),
                (),
                (),
                getMockAPILevelPolicy(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID),
                getMockAPIPolicyResponse(getMockAPILevelPolicy(check apiWithAPILevelInterceptorPolicy.cloneWithType(API), organiztion1, apiUUID).clone()),
                interceptorServices,
                k8sapiUUID,
                CreatedAPIWithAPILevelInterceptorPolicy.toBalString()
            ]
        };
        return data;
    } on fail var e {
        test:assertFail("tests failed===" + e.message());
    }
}

function getMockHttpRouteResponse(model:Httproute request) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    request.metadata.uid = uuid:createType1AsString();
    response.setJsonPayload(request.toJson());
    return response;
}

function getMockHttpRouteErrorResponse() returns http:Response {
    http:Response response = new;
    response.statusCode = 400;
    model:Status status = {code: 400, message: "httproute already exist"};
    response.setJsonPayload(status.toJson());
    return response;
}

function getMockHttpRouteWithBackend(API api, string apiUUID, string backenduuid, string 'type, commons:Organization organization) returns model:Httproute {
    string hostnames = 'type == PRODUCTION_TYPE ? string:concat(organization.uuid, ".", "gw.wso2.com") : string:concat(organization.uuid, ".", "sandbox.gw.wso2.com");
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organization)},
        "spec": {
            "hostnames": [
                hostnames
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithOperationPolicies(API api, string apiUUID, string backenduuid, string 'type, commons:Organization organization) returns model:Httproute {
    string hostnames = 'type == PRODUCTION_TYPE ? string:concat(organization.uuid, ".", "gw.wso2.com") : string:concat(organization.uuid, ".", "sandbox.gw.wso2.com");
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organization)},
        "spec": {
            "hostnames": [
                hostnames
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithAPIPolicies(API api, string apiUUID, string backenduuid, string 'type, commons:Organization organization) returns model:Httproute {
    string hostnames = 'type == PRODUCTION_TYPE ? string:concat(organization.uuid, ".", "gw.wso2.com") : string:concat(organization.uuid, ".", "sandbox.gw.wso2.com");
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organization)},
        "spec": {
            "hostnames": [
                hostnames
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "RequestHeaderModifier",
                            "requestHeaderModifier": {
                                "set": [
                                    {
                                        "name": "customadd",
                                        "value": "customvalue"
                                    }
                                ]
                            }
                        },
                        {
                            "type": "ResponseHeaderModifier",
                            "responseHeaderModifier": {
                                "remove": ["content-length"]
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithOperationRateLimits(API api, string apiUUID, string backenduuid, string 'type, commons:Organization organization) returns model:Httproute {
    string hostnames = 'type == PRODUCTION_TYPE ? string:concat(organization.uuid, ".", "gw.wso2.com") : string:concat(organization.uuid, ".", "sandbox.gw.wso2.com");
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organization)},
        "spec": {
            "hostnames": [
                hostnames
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "ExtensionRef",
                            "extensionRef": {
                                "group": "dp.wso2.com",
                                "kind": "RateLimitPolicy",
                                "name": "rate-limit-policy-ref-name"
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockHttpRouteWithOperationInterceptorPolicy(API api, string apiUUID, string backenduuid, string 'type, commons:Organization organization) returns model:Httproute {
    string hostnames = 'type == PRODUCTION_TYPE ? string:concat(organization.uuid, ".", "gw.wso2.com") : string:concat(organization.uuid, ".", "sandbox.gw.wso2.com");
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organization)},
        "spec": {
            "hostnames": [
                hostnames
            ],
            "rules": [
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "GET"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        },
                        {
                            "type": "ExtensionRef",
                            "extensionRef": {
                                "group": "dp.wso2.com",
                                "kind": "APIPolicy",
                                "name": "api-policy-ref-name"
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PUT"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "POST"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "DELETE"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                },
                {
                    "matches": [
                        {
                            "path": {
                                "type": "RegularExpression",
                                "value": "/pizzaAPI/1.0.0(.*)"
                            },
                            "method": "PATCH"
                        }
                    ],
                    "filters": [
                        {
                            "type": "URLRewrite",
                            "urlRewrite": {
                                "path": {
                                    "type": "ReplaceFullPath",
                                    "replaceFullPath": "\\1"
                                }
                            }
                        }
                    ],
                    "backendRefs": [
                        {
                            "group": "dp.wso2.com",
                            "kind": "Backend",
                            "name": backenduuid,
                            "namespace": "apk-platform"
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "default",
                    "sectionName": "httpslistener"
                }
            ]
        }
    };
}

function getMockResourceRateLimitPolicy(API api, commons:Organization organiztion, string apiUUID) returns model:RateLimitPolicy {
    return {
        "apiVersion": "dp.wso2.com/v1alpha1",
        "kind": "RateLimitPolicy",
        "metadata": {"name": "rate-limit-policy-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "default": {
                "api": {
                    "rateLimit": {
                        "requestsPerUnit": 10,
                        "unit": "Minute"
                    }
                },
                "type": "Api",
                "organization": organiztion.uuid

            },
            "targetRef": {
                "group": "dp.wso2.com",
                "kind": "Resource",
                "name": apiUUID,
                "namespace": "apk-platform"
            }
        }
    };
}

function getMockAPIRateLimitPolicy(API api, commons:Organization organiztion, string apiUUID) returns model:RateLimitPolicy {
    return {
        "apiVersion": "dp.wso2.com/v1alpha1",
        "kind": "RateLimitPolicy",
        "metadata": {"name": "rate-limit-policy-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "default": {
                "api": {
                    "rateLimit": {
                        "requestsPerUnit": 10,
                        "unit": "Minute"
                    }
                },
                "organization": organiztion.uuid,
                "type": "Api"
            },
            "targetRef": {
                "group": "gateway.networking.k8s.io",
                "kind": "API",
                "name": apiUUID,
                "namespace": "apk-platform"
            }
        }
    };
}

function getMockRateLimitResponse(model:RateLimitPolicy request) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    request.metadata.uid = uuid:createType1AsString();
    response.setJsonPayload(request.toJson());
    return response;
}

function getMockResourceLevelPolicy(API api, commons:Organization organiztion, string apiUUID) returns model:APIPolicy {
    return {
        "apiVersion": "dp.wso2.com/v1alpha1",
        "kind": "APIPolicy",
        "metadata": {"name": "api-policy-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "default": {
               "requestInterceptors": [
                    {
                        "name": getInterceptorServiceUid(api, organiztion, "request", 0),
                        "namespace": "apk-platform"
                    }
                    
                ],
                "responseInterceptors": [
                   {
                        "name": getInterceptorServiceUid(api, organiztion, "response", 0),
                        "namespace": "apk-platform"
                    }
                ]
            },
            "targetRef": {
                "group": "dp.wso2.com",
                "kind": "Resource",
                "name": apiUUID,
                "namespace": "apk-platform"
            }
        }
    };
}

function getMockAPILevelPolicy(API api, commons:Organization organiztion, string apiUUID) returns model:APIPolicy {
    return {
        "apiVersion": "dp.wso2.com/v1alpha1",
        "kind": "APIPolicy",
        "metadata": {"name": "api-policy-ref-name", "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
        "spec": {
            "default": {
                "requestInterceptors": [
                    {
                        "name": getInterceptorServiceUid(api, organiztion, "request", 0),
                        "namespace": "apk-platform"
                    }
                    
                ],
                "responseInterceptors": [
                   {
                        "name": getInterceptorServiceUid(api, organiztion, "response", 0),
                        "namespace": "apk-platform"
                    }
                ]
            },
            "targetRef": {
                "group": "dp.wso2.com",
                "kind": "API",
                "name": apiUUID,
                "namespace": "apk-platform"
            }
        }
    };
}

function getMockAPIPolicyResponse(model:APIPolicy request) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    request.metadata.uid = uuid:createType1AsString();
    response.setJsonPayload(request.toJson());
    return response;
}

function getMockInterceptorService(API api, commons:Organization organiztion, string apiUUID, string flow, string[] includes, string backendUrl) returns model:InterceptorService {
    return {
            "apiVersion": "dp.wso2.com/v1alpha1",
            "kind": "InterceptorService",
            "metadata": {"name": getInterceptorServiceUid(api, organiztion, flow, 0), "namespace": "apk-platform", "labels": getLabels(api, organiztion)},
            "spec": {
                "backendRef": {
                    "name": getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion, backendUrl),
                    "namespace": "apk-platform"
                },
                "includes": includes
            }
        };
}

function getMockInterceptorServiceResponse(model:InterceptorService request) returns http:Response {
    http:Response response = new;
    response.statusCode = 201;
    request.metadata.uid = uuid:createType1AsString();
    response.setJsonPayload(request.toJson());
    return response;
}

function createAPIDataProvider() returns map<[string, string, API, model:ConfigMap, any, model:Httproute?, any, model:Httproute?, any, [model:Backend, any][], model:API, any, model:RuntimeAPI, any, model:RateLimitPolicy?, any, model:APIPolicy?, any, [model:InterceptorService, any][], string, anydata]> {
    do {
        API api = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0",
            endpointConfig: {"production_endpoints": {"url": "https://localhost"}}
        };
        API produrlmissingAPI = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0",
            endpointConfig: {"production_endpoints": {}}
        };
        API sandboxOnlyAPI = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0",
            endpointConfig: {"sandbox_endpoints": {"url": "https://localhost"}}
        };
        API sandboxurlmissingapi = {
            name: "PizzaAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0",
            endpointConfig: {"sandbox_endpoints": {}}
        };
        API alreadyNameExist = {
            name: "pizzashackAPI",
            context: "/pizzaAPI/1.0.0",
            'version: "1.0.0"
        };
        commons:APKError nameAlreadyExistError = error commons:APKError(
            "API Name - " + alreadyNameExist.name + " already exist",
            code = 909011,
            message = "API Name - " + alreadyNameExist.name + " already exist",
            statusCode = 409,
            description = "API Name - " + alreadyNameExist.name + " already exist"
        );
        API contextAlreadyExist = {
            name: "PizzaAPI",
            context: "/pizzashack/1.0.0",
            'version: "1.0.0"
        };
        commons:APKError contextAlreadyExistError = error commons:APKError(
            "API Context - " + contextAlreadyExist.context + " already exist",
            code = 909012,
            message = "API Context - " + contextAlreadyExist.context + " already exist",
            statusCode = 409,
            description = "API Context - " + contextAlreadyExist.context + " already exist"
        );
        string apiUUID = getUniqueIdForAPI(api.name, api.'version, organiztion1);
        string backenduuid = getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion1);
        string backenduuid1 = getBackendServiceUid(api, (), SANDBOX_TYPE, organiztion1);
        string interceptorBackenduuid1 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend1.interceptor:9082");
        string interceptorBackenduuid2 = getInterceptorBackendUid(api, INTERCEPTOR_TYPE, organiztion1, "http://interceptor-backend2.interceptor:9083");
        string k8sapiUUID = uuid:createType1AsString();
        model:Backend backendService = {
            metadata: {name: backenduuid, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "localhost", port: 443}], protocol: "https"}
        };
        model:Backend backendService1 = {
            metadata: {name: backenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "localhost", port: 443}], protocol: "https"}
        };
        model:Backend interceptorBackendService1 = {
            metadata: {name: interceptorBackenduuid1, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend1.interceptor", port: 9082}], protocol: "http"}
        };
        model:Backend interceptorBackendService2 = {
            metadata: {name: interceptorBackenduuid2, namespace: "apk-platform", labels: getLabels(api, organiztion1)},
            spec: {services: [{host: "interceptor-backend2.interceptor", port: 9083}], protocol: "http"}
        };
        http:Response backendServiceResponse = getOKBackendServiceResponse(backendService);
        http:Response backendServiceResponse1 = getOKBackendServiceResponse(backendService);
        http:Response interceptorBackendServiceResponse1 = getOKBackendServiceResponse(interceptorBackendService1);
        http:Response interceptorBackendServiceResponse2 = getOKBackendServiceResponse(interceptorBackendService2);
        http:Response backendServiceErrorResponse = new;
        backendServiceErrorResponse.statusCode = 403;
        [model:Backend, any][] services = [];
        services.push([backendService, backendServiceResponse]);
        [model:Backend, any][] services1 = [];
        services.push([backendService1, backendServiceResponse1]);
        services.push([interceptorBackendService1, interceptorBackendServiceResponse1]);
        services.push([interceptorBackendService2, interceptorBackendServiceResponse2]);
        [model:Backend, any][] servicesError = [];
        servicesError.push([backendService, backendServiceErrorResponse]);
        model:ConfigMap configmap = check getMockConfigMap1(apiUUID, api);
        model:Httproute prodhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
        model:Httproute sandhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid1, SANDBOX_TYPE, organiztion1);
        string locationUrl = runtimeConfiguration.baseURl + "/apis/" + k8sapiUUID;

        CreatedAPI createdAPI = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"production_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}]}, headers: {location: locationUrl}};
        CreatedAPI createdSandboxOnlyAPI = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z", endpointConfig: {"sandbox_endpoints": {"url": "https://localhost"}}, operations: [{"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}, {"target": "/*", "verb": "PATCH", "authTypeEnabled": true, "scopes": [], "operationPolicies": {"request": [], "response": []}}]}, headers: {location: locationUrl}};

        commons:APKError productionEndpointNotSpecifiedError = error commons:APKError("Production endpoint not specified",
            code = 909014,
            message = "Production endpoint not specified",
            statusCode = 406,
            description = "Production endpoint not specified"
        );
        commons:APKError sandboxEndpointNotSpecifiedError = error commons:APKError("Sandbox endpoint not specified",
            code = 909013,
            message = "Sandbox endpoint not specified",
            statusCode = 406,
            description = "Sandbox endpoint not specified"
        );
        commons:APKError k8sLevelError = error("Internal error occured while deploying API", code = 909028, message
        = "Internal error occured while deploying API", statusCode = 500, description = "Internal error occured while deploying API", moreInfo = {});
        commons:APKError k8sLevelError1 = error commons:APKError("Internal server error", error("Internal server error"),
            code = 909022,
            message = "Internal server error",
            statusCode = 500,
            description = "Internal server error"
        );
        commons:APKError invalidAPINameError = error commons:APKError("API name PizzaAPI invalid",
            code = 909016,
            message = "API name PizzaAPI invalid",
            statusCode = 406,
            description = "API name PizzaAPI invalid"
        );
        map<[string, string, API, model:ConfigMap,
    any, model:Httproute|(), any, model:Httproute|(),
    any, [model:Backend, any][], model:API, any, model:RuntimeAPI, any,
    model:RateLimitPolicy|(), any, model:APIPolicy|(), any, [model:InterceptorService, any][], string,
    anydata]> data = {
            "1": [
                apiUUID,
                backenduuid,
                api,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                createdAPI.toBalString()
            ]
        ,
            "2": [
                apiUUID,
                backenduuid,
                alreadyNameExist,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                nameAlreadyExistError.toBalString()
            ],
            "3": [
                apiUUID,
                backenduuid,
                contextAlreadyExist,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                contextAlreadyExistError.toBalString()
            ],
            "4": [
                apiUUID,
                backenduuid,
                sandboxOnlyAPI,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                (),
                (),
                sandhttpRoute,
                getMockHttpRouteResponse(sandhttpRoute.clone()),
                services1,
                getMockAPI1(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI1(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(sandboxOnlyAPI, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(sandboxOnlyAPI, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                createdSandboxOnlyAPI.toBalString()
            ]
        ,
            "5": [
                apiUUID,
                backenduuid,
                produrlmissingAPI,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                productionEndpointNotSpecifiedError.toBalString()
            ],
            "6": [
                apiUUID,
                backenduuid,
                sandboxurlmissingapi,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                sandboxEndpointNotSpecifiedError.toBalString()
            ]
        ,
            "7": [
                apiUUID,
                backenduuid,
                api,
                configmap,
                getMockConfigMapErrorResponse(),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                k8sLevelError1.toBalString()
            ]
        ,
            "8": [
                apiUUID,
                backenduuid,
                sandboxOnlyAPI,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                (),
                (),
                sandhttpRoute,
                getMockHttpRouteErrorResponse(),
                services1,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                k8sLevelError1.toBalString()
            ]
        ,
            "9": [
                apiUUID,
                backenduuid,
                api,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteErrorResponse(),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                k8sLevelError1.toBalString()
            ]
        ,
            "10": [
                apiUUID,
                backenduuid,
                api,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                servicesError,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                k8sLevelError.toBalString()
            ]
        ,
            "11": [
                apiUUID,
                backenduuid,
                api,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIErrorResponse(),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                k8sLevelError.toBalString()
            ]
        ,
            "12": [
                apiUUID,
                backenduuid,
                api,
                configmap,
                getMockConfigMapResponse(configmap.clone()),
                prodhttpRoute,
                getMockHttpRouteResponse(prodhttpRoute.clone()),
                (),
                (),
                services,
                getMockAPI(api, apiUUID, organiztion1.uuid),
                getMockAPIErrorNameExist(),
                getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
                getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
                (),
                (),
                (),
                (),
                [],
                k8sapiUUID,
                invalidAPINameError.toBalString()
            ]
        };
        return data;
    } on fail var e {
        test:assertFail("tests failed===" + e.message());
    }
}

@test:Config {dataProvider: mediationPolicyByIdDataProvider}
public function testGetMediationPolicyById(string policyId, commons:Organization organization, anydata expectedData) {
    APIClient apiclient = new ();
    MediationPolicy|NotFoundError|commons:APKError mediationPolicyById = apiclient.getMediationPolicyById(policyId, organization);
    if mediationPolicyById is any {
        test:assertEquals(mediationPolicyById.toBalString(), expectedData);
    } else {
        test:assertEquals(mediationPolicyById.toBalString(), expectedData);
    }
}

public function mediationPolicyByIdDataProvider() returns map<[string, commons:Organization, anydata]> {
    MediationPolicy & readonly mediationPolicy1 = {
        id: "1",
        'type: MEDIATION_POLICY_TYPE_REQUEST_HEADER_MODIFIER,
        name: MEDIATION_POLICY_NAME_ADD_HEADER,
        displayName: "Add Header",
        description: "This policy allows you to add a new header to the request",
        applicableFlows: [MEDIATION_POLICY_FLOW_REQUEST],
        supportedApiTypes: [API_TYPE_REST],
        policyAttributes: [
            {
                name: "headerName",
                description: "Name of the header to be added",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            },
            {
                name: "headerValue",
                description: "Value of the header",
                'type: "String",
                required: true,
                validationRegex: "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$"
            }
        ]
    };
    commons:APKError notfound = error commons:APKError("6 not found",
        code = 909001,
        message = "6 not found",
        statusCode = 404,
        description = "6 not found"
    );
    map<[string, commons:Organization, anydata]> dataset = {
        "1": ["1", organiztion1, mediationPolicy1.toBalString()],
        "2": ["6", organiztion1, notfound.toBalString()]
    };
    return dataset;
}

@test:Config {dataProvider: getMediationPolicyListDataProvider}
public function testGetMediationPolicyList(string? query, int 'limit, int offset, string sortBy, string sortOrder, anydata expected) {
    APIClient apiclient = new ();
    any|error mediationPolicyList = apiclient.getMediationPolicyList(query, 'limit, offset, sortBy, sortOrder, organiztion1);
    if mediationPolicyList is any {
        test:assertEquals(mediationPolicyList.toBalString(), expected);
    } else {
        test:assertEquals(mediationPolicyList.toBalString(), expected);
    }
}

function getMediationPolicyListDataProvider() returns map<[string?, int, int, string, string, anydata]> {
    commons:APKError badRequestError = error commons:APKError("Invalid Sort By/Sort Order value",
        code = 909020,
        message = "Invalid Sort By/Sort Order value",
        statusCode = 406,
        description = "Invalid Sort By/Sort Order value"
    );
    commons:APKError badRequest = error commons:APKError("Invalid keyword name1",
        code = 909019,
        message = "Invalid keyword name1",
        statusCode = 406,
        description = "Invalid keyword name1"
    );
    map<[string?, int, int, string, string, anydata]> dataSet = {
        "1": [
            (),
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 5,
                "list": [
                    {
                        "id": "5",
                        "type": "Interceptor",
                        "name": "Interceptor",
                        "displayName": "Interceptor",
                        "description": "This policy allows you to engage an interceptor service",
                        "applicableFlows": [
                            "request",
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headersEnabled",
                                "description": "Indicates whether request/response header details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "bodyEnabled",
                                "description": "Indicates whether request/response body details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "contextEnabled",
                                "description": "Indicates whether context details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "trailersEnabled",
                                "description": "Indicates whether request/response trailer details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "backendUrl",
                                "description": "Backend URL of the interceptor service",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "2",
                        "type": "RequestHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "4",
                        "type": "ResponseHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "2": [
            (),
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_DESC,
            {
                "count": 5,
                "list": [
                    {
                        "id": "2",
                        "type": "RequestHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "4",
                        "type": "ResponseHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "5",
                        "type": "Interceptor",
                        "name": "Interceptor",
                        "displayName": "Interceptor",
                        "description": "This policy allows you to engage an interceptor service",
                        "applicableFlows": [
                            "request",
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headersEnabled",
                                "description": "Indicates whether request/response header details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "bodyEnabled",
                                "description": "Indicates whether request/response body details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "contextEnabled",
                                "description": "Indicates whether context details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "trailersEnabled",
                                "description": "Indicates whether request/response trailer details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "backendUrl",
                                "description": "Backend URL of the interceptor service",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "3": [
            (),
            10,
            0,
            SORT_BY_ID,
            SORT_ORDER_ASC,
            {
                "count": 5,
                "list": [
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "2",
                        "type": "RequestHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "4",
                        "type": "ResponseHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "5",
                        "type": "Interceptor",
                        "name": "Interceptor",
                        "displayName": "Interceptor",
                        "description": "This policy allows you to engage an interceptor service",
                        "applicableFlows": [
                            "request",
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headersEnabled",
                                "description": "Indicates whether request/response header details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "bodyEnabled",
                                "description": "Indicates whether request/response body details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "contextEnabled",
                                "description": "Indicates whether context details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "trailersEnabled",
                                "description": "Indicates whether request/response trailer details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "backendUrl",
                                "description": "Backend URL of the interceptor service",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "4": [
            (),
            10,
            0,
            SORT_BY_ID,
            SORT_ORDER_DESC,
            {
                "count": 5,
                "list": [
                    {
                        "id": "5",
                        "type": "Interceptor",
                        "name": "Interceptor",
                        "displayName": "Interceptor",
                        "description": "This policy allows you to engage an interceptor service",
                        "applicableFlows": [
                            "request",
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headersEnabled",
                                "description": "Indicates whether request/response header details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "bodyEnabled",
                                "description": "Indicates whether request/response body details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "contextEnabled",
                                "description": "Indicates whether context details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "trailersEnabled",
                                "description": "Indicates whether request/response trailer details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "backendUrl",
                                "description": "Backend URL of the interceptor service",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "4",
                        "type": "ResponseHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "2",
                        "type": "RequestHeaderModifier",
                        "name": "removeHeader",
                        "displayName": "Remove Header",
                        "description": "This policy allows you to remove a header from the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be removed",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "5": [(), 10, 0, "description", SORT_ORDER_DESC, badRequestError.toBalString()],
        "6": [
            (),
            2,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 2,
                "list": [
                    {
                        "id": "5",
                        "type": "Interceptor",
                        "name": "Interceptor",
                        "displayName": "Interceptor",
                        "description": "This policy allows you to engage an interceptor service",
                        "applicableFlows": [
                            "request",
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headersEnabled",
                                "description": "Indicates whether request/response header details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "bodyEnabled",
                                "description": "Indicates whether request/response body details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "contextEnabled",
                                "description": "Indicates whether context details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "trailersEnabled",
                                "description": "Indicates whether request/response trailer details should be sent to the interceptor service",
                                "required": false,
                                "type": "boolean"
                            },
                            {
                                "name": "backendUrl",
                                "description": "Backend URL of the interceptor service",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 2,
                    "total": 5,
                    "next": "/policies?limit=2&offset=2&sortBy=policyName&sortOrder=asc&query=",
                    "previous": ""
                }
            }.toBalString()
        ],
        "7": [
            (),
            2,
            2,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_DESC,
            {
                "count": 2,
                "list": [
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 2,
                    "limit": 2,
                    "total": 5,
                    "next": "/policies?limit=2&offset=4&sortBy=policyName&sortOrder=desc&query=",
                    "previous": "/policies?limit=2&offset=0&sortBy=policyName&sortOrder=desc&query="
                }
            }.toBalString()
        ],
        "8": [
            (),
            3,
            6,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 0,
                "list": [],
                "pagination": {
                    "offset": 6,
                    "limit": 3,
                    "total": 5,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "9": [
            "name:add",
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 2,
                "list": [
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 2,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "10": [
            "add",
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 2,
                "list": [
                    {
                        "id": "1",
                        "type": "RequestHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the request",
                        "applicableFlows": [
                            "request"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    },
                    {
                        "id": "3",
                        "type": "ResponseHeaderModifier",
                        "name": "addHeader",
                        "displayName": "Add Header",
                        "description": "This policy allows you to add a new header to the response",
                        "applicableFlows": [
                            "response"
                        ],
                        "supportedApiTypes": [
                            "REST"
                        ],
                        "policyAttributes": [
                            {
                                "name": "headerName",
                                "description": "Name of the header to be added",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            },
                            {
                                "name": "headerValue",
                                "description": "Value of the header",
                                "required": true,
                                "validationRegex": "^([a-zA-Z_][a-zA-Z\\d_\\-\\ ]*)$",
                                "type": "String"
                            }
                        ]
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 2,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "11": [
            "type:modify",
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 0,
                "list": [],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 0,
                    "next": "",
                    "previous": ""
                }
            }.toBalString()
        ],
        "12": [
            "name1:add",
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            badRequest.toBalString()
        ]
    };
    return dataSet;
}

function getOKBackendServiceResponse(model:Backend backendService) returns http:Response {
    http:Response backendServiceResponse = new;
    backendServiceResponse.statusCode = 201;
    model:Backend serviceClone = backendService.clone();
    serviceClone.metadata.uid = uuid:createType1AsString();
    backendServiceResponse.setJsonPayload(serviceClone.toJson());
    return backendServiceResponse;
}

function getMockAPIErrorResponse() returns http:Response {
    http:Response response = new;
    response.statusCode = 404;
    return response;
}

function getMockAPIErrorNameExist() returns http:Response {
    http:Response response = new;
    response.statusCode = 400;
    model:StatusCause[] causes = [{'field: "spec.apiDisplayName", message: "API Name already Exist", reason: "FieldValueDuplicate"}];
    model:Status status = {code: 400, details: {'causes: causes, group: "dp.wso2.com", kind: "API", name: uuid:createType1AsString()}};
    response.setJsonPayload(status.toJson());
    return response;
}
