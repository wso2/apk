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

@test:Mock {functionName: "startAndAttachServices"}
function getMockStartandAttachServices() returns error? {
}

@test:Mock {functionName: "getBackendPolicyUid"}
function testgetBackendPolicyUid(API api, string? endpointType, commons:Organization organization) returns string {
    return "backendpolicy-uuid";
}

@test:Mock {functionName: "retrieveHttpRouteRefName"}
function testRetrieveHttpRouteRefName(API api, string 'type, commons:Organization organization) returns string {
    return "http-route-ref-name";
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

int orgWatchIndex = 0;

@test:Mock {functionName: "getOrganizationWatchClient"}
function getMockOrganiationClient(string resourceVersion) returns websocket:Client|error {
    string initialConectionId = uuid:createType1AsString();
    if resourceVersion == "28702" {
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(getOrganizationWatchEvent());
        return mock;
    } else if resourceVersion == "28705" {
        string connectionId = uuid:createType1AsString();
        websocket:Client mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(connectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getNextOrganizationEvent(), ());
        return mock;
    } else if resourceVersion == "28714" {
        if orgWatchIndex == 0 {
            websocket:Error websocketError = error("Error", message = "Error");
            orgWatchIndex += 1;
            return websocketError;
        } else {
            initialConectionId = uuid:createType1AsString();
            websocket:Client mock = test:mock(websocket:Client);
            test:prepare(mock).when("isOpen").thenReturn(true);
            test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
            test:prepare(mock).when("readMessage").thenReturnSequence(getOrganizationWatchDeleteEvent(), ());
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
function getMockK8sClient() returns http:Client {
    http:Client mockK8sClient = test:mock(http:Client);
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis")
        .thenReturn(getMockAPIList());
    string fieldSlector = "metadata.namespace%21%3Dkube-system%2Cmetadata.namespace%21%3Dkubernetes-dashboard%2Cmetadata.namespace%21%3Dgateway-system%2Cmetadata.namespace%21%3Dingress-nginx%2Cmetadata.namespace%21%3Dapk-platform";
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/services?fieldSelector=" + fieldSlector)
        .thenReturn(getMockServiceList());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/servicemappings")
        .thenReturn(getMockServiceMappings());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b08-f2b1-1166-82d5-649ae706d29e").thenReturn(mock404Response());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk/apis/pizzashackAPI1").thenReturn(mock404Response());
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114a-definition").thenReturn(mockConfigMaps());
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
    model:API|error aPI = getAPI(id, organization);
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
        "1": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/order/{orderId}", verb: "POST"}, "/v3/f77cc767/order/\\1"],
        "2": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/menu", verb: "GET"}, "/v3/f77cc767/menu"],
        "3": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/menu", verb: "GET"}, "/v3/f77cc767/menu"],
        "4": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/*", verb: "GET"}, "/v3/f77cc767/\\1"],
        "5": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: true, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/*", verb: "GET"}, "\\1"],
        "6": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/order/{orderId}/details/{item}", verb: "GET"}, "/v3/f77cc767/order/\\1/details/\\2"]
    };
    return dataSet;
}

@test:Config {dataProvider: apiDefinitionDataProvider}
public function testGetAPIDefinitionByID(string apiid, anydata expectedResponse) returns error? {
    APIClient apiclient = new ();
    http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError aPIDefinitionByID = apiclient.getAPIDefinitionByID(apiid, organiztion1);
    if aPIDefinitionByID is http:Response {
        json jsonPayload = check aPIDefinitionByID.getJsonPayload();
        test:assertEquals(jsonPayload.toBalString(), expectedResponse);
    } else {
        test:assertEquals(aPIDefinitionByID.toBalString(), expectedResponse);
    }
}

public function apiDefinitionDataProvider() returns map<[string, anydata]> {
    NotFoundError notfound = {body: {code: 909100, message: "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found."}};
    InternalServerErrorError internalError = {body: {code: 909000, message: "Internal Error Occured while retrieving definition"}};

    map<[string, anydata]> dataSet = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", mockOpenAPIJson().toBalString()],
        "2": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e9", notfound.toBalString()],
        "3": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f9", mockpizzashackAPI11Definition().toBalString()],
        "4": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f1", mockPizzashackAPI12Definition().toBalString()],
        "5": ["7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1", mockPizzaShackAPI1Definition(organiztion1.uuid).toBalString()],
        "6": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f3", internalError.toBalString()]

    };
    return dataSet;
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
        endpointConfig: {"endpoint_type": "http", "sandbox_endpoints": {"url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"}, "production_endpoints": {"url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"}},
        operations: [
            {"target": "/*", "verb": "GET", "authTypeEnabled": true, "scopes": []},
            {"target": "/*", "verb": "PUT", "authTypeEnabled": true, "scopes": []},
            {"target": "/*", "verb": "POST", "authTypeEnabled": true, "scopes": []},
            {"target": "/*", "verb": "DELETE", "authTypeEnabled": true, "scopes": []}
        ],
        createdTime: "2022-12-13T09:45:47Z"
    };
    NotFoundError notfound = {body: {code: 909100, message: "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found."}};
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
    BadRequestError badRequestError = {"body": {"code": 90912, "message": "Invalid Sort By/Sort Order Value "}};
    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord type1"}};

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
                    "total": 6
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
                    "total": 6
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
                    "total": 6
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
                    "total": 6
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
                    "total": 6
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
                    "total": 6
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
                    "total": 6
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
                    "total": 5
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
                    "total": 5
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
                    "total": 6
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
                    "total": 0
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
                        "https://localhost:9443/am/sample/pizzashack/v1/api/"
                    ],
                    "type": "http"
                },
                "x-wso2-sandbox-endpoints": {
                    "urls": [
                        "https://localhost:9443/am/sample/pizzashack/v1/api/"
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
                        "https://localhost:9443/am/sample/pizzashack/v1/api/"
                    ],
                    "type": "http"
                },
                "x-wso2-sandbox-endpoints": {
                    "urls": [
                        "https://localhost:9443/am/sample/pizzashack/v1/api/"
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
    test:assertEquals(apiClient.validateAPIExistence(query, organiztion1).toBalString(), expected);
}

function validateExistenceDataProvider() returns map<[string, anydata]> {
    http:Ok ok = {};
    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord type"}};
    NotFoundError notFound = {body: {code: 900914, message: "context/name doesn't exist"}};
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
function testCreateAPIFromService(string serviceUUId, string apiUUID, [model:ConfigMap, any] configmapResponse, [model:Httproute, any] httproute, [model:K8sServiceMapping, any] servicemapping, [model:API, any] k8sAPI, [model:RuntimeAPI, any] runtimeAPI, API api, string k8sapiUUID, anydata expected) returns error? {
    APIClient apiClient = new;
    http:Response configmapResponse404 = new;
    configmapResponse404.statusCode = 404;
    http:ApplicationResponseError internalApiResponse = error("", statusCode = 404, body = "Not Found", headers = {});
    model:HttprouteList httpRouteList = {metadata: {}, items: []};
    model:ServiceMappingList serviceMappingList = {metadata: {}, items: []};
    model:AuthenticationList authenticationList = {metadata: {}, items: []};
    model:BackendPolicyList backendPolicyList = {metadata: {}, items: []};
    model:ServiceList serviceList = {metadata: {}, items: []};
    model:ScopeList scopeList = {metadata: {}, items: []};
    test:prepare(k8sApiServerEp).when("post").withArguments("/api/v1/namespaces/apk-platform/configmaps", configmapResponse[0]).thenReturn(configmapResponse[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", httproute[0]).thenReturn(httproute[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings", servicemapping[0]).thenReturn(servicemapping[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis", k8sAPI[0]).thenReturn(k8sAPI[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis", runtimeAPI[0]).thenReturn(runtimeAPI[1]);
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/" + apiClient.retrieveDefinitionName(apiUUID)).thenReturn(configmapResponse404);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes/?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(httpRouteList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(serviceMappingList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/authentications?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(authenticationList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backendpolicies?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(backendPolicyList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/scopes?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(scopeList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/namespaces/apk-platform/services?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(serviceList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sAPI[0].metadata.name).thenReturn(internalApiResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sAPI[0].metadata.name).thenReturn(runtimeAPI[0]);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/" + k8sAPI[0].metadata.name).thenReturn(configmapResponse404);
    any|error aPIFromService = apiClient.createAPIFromService(serviceUUId, api, organiztion1);
    if aPIFromService is any {
        test:assertEquals(aPIFromService.toBalString(), expected);
    } else {
        test:assertEquals(aPIFromService.toBalString(), expected);
    }
}

function createApiFromServiceDataProvider() returns map<[string, string, [model:ConfigMap, any], [model:Httproute, any], [model:K8sServiceMapping, any], [model:API, any], [model:RuntimeAPI, any], API, string, anydata]> {
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
    string apiUUID = getUniqueIdForAPI(api.name, api.'version, organiztion1);
    model:ConfigMap configmap = getMockConfigMap1(apiUUID, api);
    http:Response mockConfigMapResponse = getMockConfigMapResponse(configmap.clone());
    model:Httproute httpRoute = getMockHttpRoute(api, apiUUID, organiztion1);
    http:Response httpRouteResponse = getMockHttpRouteResponse(httpRoute.clone());
    model:K8sServiceMapping mockServiceMappingRequest = getMockServiceMappingRequest(api, apiUUID);
    model:API mockAPI = getMockAPI(api, apiUUID, organiztion1.uuid);
    http:Response mockAPIResponse = getMockAPIResponse(mockAPI.clone(), k8sAPIUUID1);
    Service serviceRecord = {
        name: "backend",
        namespace: "apk",
        id: "275b00d1-722c-4df2-b65a-9b14677abe4b",
        'type: "ClusterIP"
    };
    model:RuntimeAPI mockRuntimeAPI = getMockRuntimeAPI(api, apiUUID, organiztion1, serviceRecord);
    http:Response mockRuntimeResponse = getMockRuntimeAPIResponse(mockRuntimeAPI.clone());
    http:Response serviceMappingResponse = getMockServiceMappingResponse(mockServiceMappingRequest.clone());
    BadRequestError nameAlreadyExistError = {body: {code: 90911, message: "API Name - " + alreadyNameExist.name + " already exist.", description: "API Name - " + alreadyNameExist.name + " already exist."}};
    API contextAlreadyExist = {
        name: "PizzaAPI",
        context: "/pizzashack/1.0.0",
        'version: "1.0.0"
    };
    BadRequestError contextAlreadyExistError = {body: {code: 90911, message: "API Context - " + contextAlreadyExist.context + " already exist.", description: "API Context " + contextAlreadyExist.context + " already exist."}};
    BadRequestError serviceNotExist = {body: {code: 90913, message: "Service from 275b00d1-722c-4df2-b65a-9b14677abe4a not found."}};

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
        }
    };
    map<[string, string, [model:ConfigMap, any], [model:Httproute, any], [model:K8sServiceMapping, any], [model:API, any], [model:RuntimeAPI, any], API, string, anydata]> data = {
        "1": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], api, k8sAPIUUID1, createdAPI.toBalString()],
        "2": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], alreadyNameExist, k8sAPIUUID1, nameAlreadyExistError.toBalString()],
        "3": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], contextAlreadyExist, k8sAPIUUID1, contextAlreadyExistError.toBalString()],
        "4": ["275b00d1-722c-4df2-b65a-9b14677abe4a", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], [mockRuntimeAPI, mockRuntimeResponse], api, k8sAPIUUID1, serviceNotExist.toBalString()]
    };
    return data;
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
    model:API k8sapi = {
        "kind": "API",
        "apiVersion": "dp.wso2.com/v1alpha1",
        "metadata": {"name": apiUUID, "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
        "spec": {
            "apiDisplayName": api.name,
            "apiType": "REST",
            "apiVersion": api.'version,
            "context": apiClient.returnFullContext(api.context, api.'version),
            "organization": organization,
            "definitionFileRef": apiUUID + "-definition",
            "prodHTTPRouteRefs": ["http-route-ref-name"]
        },
        "status"
                : null
    };
    return k8sapi;
}

function getMockRuntimeAPI(API api, string apiUUID, commons:Organization organization, Service? serviceEntry) returns model:RuntimeAPI {
    APIClient apiClient = new;
    model:RuntimeAPI runtimeAPI = apiClient.generateRuntimeAPIArtifact(api, serviceEntry, organization);
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
    model:API k8sapi = {
        "kind": "API",
        "apiVersion": "dp.wso2.com/v1alpha1",
        "metadata": {"name": apiUUID, "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
        "spec": {
            "apiDisplayName": api.name,
            "apiType": "REST",
            "apiVersion": api.'version,
            "context": apiClient.returnFullContext(api.context, api.'version),
            "organization": organization,
            "definitionFileRef": apiUUID + "-definition",
            "sandHTTPRouteRefs": ["http-route-ref-name"]
        },
        "status": null
    };
    return k8sapi;
}

function getMockServiceMappingRequest(API api, string apiUUID) returns model:K8sServiceMapping {
    model:K8sServiceMapping serviceMapping = {"kind": "ServiceMapping", "apiVersion": "dp.wso2.com/v1alpha1", "metadata": {"name": apiUUID + "-servicemapping", "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}}, "spec": {"serviceRef": {"name": "backend", "namespace": "apk"}, "apiRef": {"name": apiUUID, "namespace": "apk-platform"}}};
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
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
        "spec": {
            "hostnames": [string:concat(organiztion.uuid, ".", "gw.wso2.com")],
            "rules": [
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "GET"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "PUT"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "POST"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "DELETE"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/pizzaAPI/1.0.0(.*)"}, "method": "PATCH"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                }
            ],
            "parentRefs": [{"group": "gateway.networking.k8s.io", "kind": "Gateway", "name": "Default"}]
        }
    };
}

function getMockConfigMap1(string apiUniqueId, API api) returns model:ConfigMap {
    model:ConfigMap configmap = {
        "apiVersion": "v1",
        "data": {
            "openapi.json": "{\"openapi\":\"3.0.1\", \"info\":{\"title\":\"" + api.name + "\", \"version\":\"" + api.'version + "\"}, \"security\":[{\"default\":[]}], \"paths\":{\"/*\":{\"get\":{\"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-auth-type\":true, \"x-throttling-tier\":\"Unlimited\"}, \"put\":{\"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-auth-type\":true, \"x-throttling-tier\":\"Unlimited\"}, \"post\":{\"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-auth-type\":true, \"x-throttling-tier\":\"Unlimited\"}, \"delete\":{\"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-auth-type\":true, \"x-throttling-tier\":\"Unlimited\"}, \"patch\":{\"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-auth-type\":true, \"x-throttling-tier\":\"Unlimited\"}}}, \"components\":{\"securitySchemes\":{\"default\":{\"type\":\"oauth2\", \"flows\":{\"implicit\":{\"authorizationUrl\":\"https://test.com\", \"scopes\":{}}}}}}}"
        },
        "kind": "ConfigMap",
        "metadata": {
            "labels": {
                "api-name": api.name,
                "api-version": api.'version
            },
            "name": apiUniqueId + "-definition",
            "namespace": "apk-platform"
        }
    };
    return configmap;
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
        [model:Service, any][] backendServices, [model:BackendPolicy, any][] backendPolicies,
        model:API k8sApi, any k8sapiResponse, model:RuntimeAPI runtimeAPI, any runtimeAPIResponse
, string k8sapiUUID, anydata expected) returns error? {
    APIClient apiClient = new;
    test:prepare(k8sApiServerEp).when("post").withArguments("/api/v1/namespaces/apk-platform/configmaps", configmap).thenReturn(configmapDeployingResponse);
    if prodhttpRoute is model:Httproute {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", prodhttpRoute).thenReturn(prodhttpResponse);
    }
    if sandHttpRoute is model:Httproute {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", sandHttpRoute).thenReturn(sandhttpResponse);
    }
    foreach [model:Service, any] servicesResponse in backendServices {
        test:prepare(k8sApiServerEp).when("post").withArguments("/api/v1/namespaces/apk-platform/services", servicesResponse[0]).thenReturn(servicesResponse[1]);
    }
    foreach [model:BackendPolicy, any] backendPolicy in backendPolicies {
        test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backendpolicies", backendPolicy[0]).thenReturn(backendPolicy[1]);
    }
    http:Response configmapResponse = new;
    configmapResponse.statusCode = 404;
    http:ApplicationResponseError internalApiResponse = error("internal api not found", statusCode = 404, body = {}, headers = {});
    http:Response internalAPIDeletionResponse = new;
    internalAPIDeletionResponse.statusCode = 200;
    model:HttprouteList httpRouteList = {metadata: {}, items: []};
    model:ServiceMappingList serviceMappingList = {metadata: {}, items: []};
    model:AuthenticationList authenticationList = {metadata: {}, items: []};
    model:BackendPolicyList backendPolicyList = {metadata: {}, items: []};
    model:ServiceList serviceList = {metadata: {}, items: []};
    model:ScopeList scopeList = {metadata: {}, items: []};
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/scopes?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(scopeList);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis", k8sApi).thenReturn(k8sapiResponse);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis", runtimeAPI).thenReturn(runtimeAPIResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/" + apiClient.retrieveDefinitionName(apiUUID)).thenReturn(configmapResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes/?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(httpRouteList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(serviceMappingList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/authentications?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(authenticationList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/backendpolicies?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(backendPolicyList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/namespaces/apk-platform/services?labelSelector=" + check generateUrlEncodedLabelSelector(api.name, api.'version)).thenReturn(serviceList);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sApi.metadata.name).thenReturn(internalApiResponse);
    test:prepare(k8sApiServerEp).when("delete").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sApi.metadata.name).thenReturn(internalAPIDeletionResponse);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/runtimeapis/" + k8sApi.metadata.name).thenReturn(runtimeAPI);
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/" + k8sApi.metadata.name).thenReturn(configmapResponse);
    any|error aPI = apiClient.createAPI(api, (), organiztion1);
    if aPI is any {
        test:assertEquals(aPI.toBalString(), expected);
    } else {
        test:assertEquals(aPI.toBalString(), expected);
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
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "Default"
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
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "Default"
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
        "metadata": {"name": "http-route-ref-name", "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
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
                            "weight": 1,
                            "group": "",
                            "kind": "Service",
                            "name": backenduuid,
                            "namespace": "apk-platform",
                            "port": 443
                        }
                    ]
                }
            ],
            "parentRefs": [
                {
                    "group": "gateway.networking.k8s.io",
                    "kind": "Gateway",
                    "name": "Default"
                }
            ]
        }
    };
}

function createAPIDataProvider() returns map<[string, string, API, model:ConfigMap, any, model:Httproute?, any, model:Httproute?, any, [model:Service, any][], [model:BackendPolicy, any][], model:API, any, model:RuntimeAPI, any, string, anydata]> {
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
    BadRequestError nameAlreadyExistError = {body: {code: 90911, message: "API Name - " + alreadyNameExist.name + " already exist.", description: "API Name - " + alreadyNameExist.name + " already exist."}};
    API contextAlreadyExist = {
        name: "PizzaAPI",
        context: "/pizzashack/1.0.0",
        'version: "1.0.0"
    };
    BadRequestError contextAlreadyExistError = {body: {code: 90911, message: "API Context - " + contextAlreadyExist.context + " already exist.", description: "API Context " + contextAlreadyExist.context + " already exist."}};
    API apiWithOperationPolicies = {
        name: "PizzaAPI",
        context: "/pizzaAPI/1.0.0",
        'version: "1.0.0",
        endpointConfig: {"production_endpoints": {"url": "https://localhost"}},
        operations: [
        {
            "target": "/*",
            "verb": "GET",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000,
            "operationPolicies": {
                "request": [
                    {
                        "policyName": "addHeader",
                        "parameters": [{
                            "headerName": "customadd",
                            "headerValue": "customvalue"
                        }]
                    }
                ],
                "response": [
                    {
                        "policyName": "removeHeader",
                        "parameters": [{
                            "headerName": "content-length"
                        }]
                    }
                ]
            }
        },
        {
            "target": "/*",
            "verb": "PUT",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "POST",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "DELETE",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "PATCH",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        }]
    };
    API apiWithAPIPolicies = {
        name: "PizzaAPI",
        context: "/pizzaAPI/1.0.0",
        'version: "1.0.0",
        endpointConfig: {"production_endpoints": {"url": "https://localhost"}},
        operations: [
        {
            "target": "/*",
            "verb": "GET",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "PUT",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "POST",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "DELETE",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        },
        {
            "target": "/*",
            "verb": "PATCH",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000
        }],
        apiPolicies: {
            "request": [
                {
                    "policyName": "addHeader",
                    "parameters": [{
                        "headerName": "customadd",
                        "headerValue": "customvalue"
                    }]
                }
            ],
            "response": [
                {
                    "policyName": "removeHeader",
                    "parameters": [{
                        "headerName": "content-length"
                    }]
                }
            ]
        }
    };
    API apiWithBothPolicies = {
        name: "PizzaAPIPolicies",
        context: "/pizzaAPIPolcies/1.0.0",
        'version: "1.0.0",
        endpointConfig: {"production_endpoints": {"url": "https://localhost"}},
        operations: [{
            "target": "/menu",
            "verb": "GET",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000,
            "operationPolicies": {
                "request": [
                    {
                        "policyName": "addHeader",
                        "parameters": [{
                            "headerName": "customadd",
                            "headerValue": "customvalue"
                        }]
                    }
                ]
            }
        }],
        apiPolicies: {
            "request": [
                {
                    "policyName": "addHeader",
                    "parameters": [{
                        "headerName": "customadd",
                        "headerValue": "customvalue"
                    }]
                }
            ]
        }
    };
    BadRequestError bothPoliciesPresentError = {body: {code: 90917, message: "Presence of both resource level and API level operation policies is not allowed"}};
    API apiWithInvalidPolicyName = {
        name: "PizzaAPIOps",
        context: "/pizzaAPIOps/1.0.0",
        'version: "1.0.0",
        endpointConfig: {"production_endpoints": {"url": "https://localhost"}},
        operations: [
        {
            "target": "/menu",
            "verb": "GET",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000,
            "operationPolicies": {
                "request": [
                    {
                        "policyName": "addHeader1",
                        "parameters": [{
                            "headerName": "customadd",
                            "headerValue": "customvalue"
                        }]
                    }
                ]
            }
        }]
    };
    BadRequestError invalidPolicyNameError = {body: {code: 90915, message: "Invalid operation policy name"}};
    API apiWithInvalidPolicyParameters = {
        name: "PizzaAPIOps",
        context: "/pizzaAPIOps/1.0.0",
        'version: "1.0.0",
        endpointConfig: {"production_endpoints": {"url": "https://localhost"}},
        operations: [
        {
            "target": "/menu",
            "verb": "GET",
            "authTypeEnabled": true,
            "throttlingPolicy": 1000,
            "operationPolicies": {
                "request": [
                    {
                        "policyName": "addHeader",
                        "parameters": [{
                            "headerName1": "customadd",
                            "headerValue": "customvalue"
                        }]
                    }
                ]
            }
        }]
    };
    BadRequestError invalidPolicyParametersError = {body: {code: 90916, message: "Invalid parameters provided for policy " + "addHeader"}};
    string apiUUID = getUniqueIdForAPI(api.name, api.'version, organiztion1);
    string backenduuid = getBackendServiceUid(api, (), PRODUCTION_TYPE, organiztion1);
    string backenduuid1 = getBackendServiceUid(api, (), SANDBOX_TYPE, organiztion1);
    string k8sapiUUID = uuid:createType1AsString();
    model:Service backendService = {
        metadata: {name: backenduuid, namespace: "apk-platform", labels: {"api-name": api.name, "api-version": api.'version}},
        spec: {externalName: "localhost", 'type: "ExternalName"}
    };
    model:Service backendService1 = {
        metadata: {name: backenduuid1, namespace: "apk-platform", labels: {"api-name": api.name, "api-version": api.'version}},
        spec: {externalName: "localhost", 'type: "ExternalName"}
    };
    http:Response backendServiceResponse = getOKBackendServiceResponse(backendService);
    http:Response backendServiceResponse1 = getOKBackendServiceResponse(backendService);
    http:Response backendServiceErrorResponse = new;
    backendServiceErrorResponse.statusCode = 403;
    [model:Service, any][] services = [];
    services.push([backendService, backendServiceResponse]);
    [model:Service, any][] services1 = [];
    services.push([backendService1, backendServiceResponse1]);
    [model:Service, any][] servicesError = [];
    servicesError.push([backendService, backendServiceErrorResponse]);
    [model:BackendPolicy, any][] backendPolicies = [];
    model:BackendPolicy backendPolicy = {
        metadata: {name: "backendpolicy-uuid", namespace: "apk-platform", labels: {"api-name": api.name, "api-version": api.'version}},
        spec: {
            default: {protocol: "https"},
            targetRef: {
                kind: "Service",
                name: backendService.metadata.name,
                namespace: backendService.metadata.namespace,
                group: ""
            }
        }
    };
    model:BackendPolicy backendPolicy1 = {
        metadata: {name: "backendpolicy-uuid", namespace: "apk-platform", labels: {"api-name": api.name, "api-version": api.'version}},
        spec: {
            default: {protocol: "https"},
            targetRef: {
                kind: "Service",
                name: backendService1.metadata.name,
                namespace: backendService1.metadata.namespace,
                group: ""
            }
        }
    };
    http:Response backendPolicyResponse = getOKBackendPolicyResponse(backendPolicy);
    http:Response backendPolicy1Response = getOKBackendPolicyResponse(backendPolicy1);
    backendPolicies.push([backendPolicy, backendPolicyResponse]);
    backendPolicies.push([backendPolicy1, backendPolicy1Response]);
    model:ConfigMap configmap = getMockConfigMap1(apiUUID, api);
    model:Httproute prodhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
    model:Httproute sandhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid1, SANDBOX_TYPE, organiztion1);
    model:Httproute prodhttpRouteWithOperationPolicies = getMockHttpRouteWithOperationPolicies(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);
    model:Httproute prodhttpRouteWithAPIPolicies = getMockHttpRouteWithAPIPolicies(api, apiUUID, backenduuid, PRODUCTION_TYPE, organiztion1);

    CreatedAPI createdAPI = {body: {name: "PizzaAPI", context: "/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID, createdTime: "2023-01-17T11:23:49Z"}};
    commons:APKError productionEndpointNotSpecifiedError = error("Production Endpoint Not specified", message = "Endpoint Not specified", description = "Production Endpoint Not specified", code = 90911, statusCode = 400);
    commons:APKError sandboxEndpointNotSpecifiedError = error("Sandbox Endpoint Not specified", message = "Endpoint Not specified", description = "Sandbox Endpoint Not specified", code = 90911, statusCode = 400);
    commons:APKError k8sLevelError = error("Internal Error occured while deploying API", code = 909000, message
        = "Internal Error occured while deploying API", statusCode = 500, description = "Internal Error occured while deploying API", moreInfo = {});
    commons:APKError k8sLevelError1 = error("Internal Server Error", code = 900900, message
        = "Internal Server Error", statusCode = 500, description = "Internal Server Error", moreInfo = {});
    commons:APKError invalidAPINameError = error("Invalid API Name", code = 90911, message = "Invalid API Name", statusCode = 400, description = "API Name PizzaAPI Invalid", moreInfo = {});
    map<[string, string, API, model:ConfigMap,
    any, model:Httproute|(), any, model:Httproute|(),
    any, [model:Service, any][], [model:BackendPolicy, any][], model:API, any, model:RuntimeAPI, any, string,
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI1(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI1(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(sandboxOnlyAPI, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(sandboxOnlyAPI, apiUUID, organiztion1, ())),
            k8sapiUUID,
            createdAPI.toBalString()
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIErrorResponse(),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
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
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIErrorNameExist(),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
            k8sapiUUID,
            invalidAPINameError.toBalString()
        ]
        ,
        "13": [
            apiUUID,
            backenduuid,
            apiWithBothPolicies,
            configmap,
            getMockConfigMapResponse(configmap.clone()),
            prodhttpRoute,
            getMockHttpRouteResponse(prodhttpRoute.clone()),
            (),
            (),
            services,
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIErrorNameExist(),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
            k8sapiUUID,
            bothPoliciesPresentError.toBalString()
        ]
        ,
        "14": [
            apiUUID,
            backenduuid,
            apiWithInvalidPolicyName,
            configmap,
            getMockConfigMapResponse(configmap.clone()),
            prodhttpRoute,
            getMockHttpRouteResponse(prodhttpRoute.clone()),
            (),
            (),
            services,
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIErrorNameExist(),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
            k8sapiUUID,
            invalidPolicyNameError.toBalString()
        ]
        ,
        "15": [
            apiUUID,
            backenduuid,
            apiWithInvalidPolicyParameters,
            configmap,
            getMockConfigMapResponse(configmap.clone()),
            prodhttpRoute,
            getMockHttpRouteResponse(prodhttpRoute.clone()),
            (),
            (),
            services,
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIErrorNameExist(),
            getMockRuntimeAPI(api, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(api, apiUUID, organiztion1, ())),
            k8sapiUUID,
            invalidPolicyParametersError.toBalString()
        ]
        ,
        "16": [
            apiUUID,
            backenduuid,
            apiWithOperationPolicies,
            configmap,
            getMockConfigMapResponse(configmap.clone()),
            prodhttpRouteWithOperationPolicies,
            getMockHttpRouteResponse(prodhttpRouteWithOperationPolicies.clone()),
            (),
            (),
            services,
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(apiWithOperationPolicies, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(apiWithOperationPolicies, apiUUID, organiztion1, ())),
            k8sapiUUID,
            createdAPI.toBalString()
        ]
        ,
        "17": [
            apiUUID,
            backenduuid,
            apiWithAPIPolicies,
            configmap,
            getMockConfigMapResponse(configmap.clone()),
            prodhttpRouteWithAPIPolicies,
            getMockHttpRouteResponse(prodhttpRouteWithAPIPolicies.clone()),
            (),
            (),
            services,
            backendPolicies,
            getMockAPI(api, apiUUID, organiztion1.uuid),
            getMockAPIResponse(getMockAPI(api, apiUUID, organiztion1.uuid), k8sapiUUID),
            getMockRuntimeAPI(apiWithAPIPolicies, apiUUID, organiztion1, ()),
            getMockRuntimeAPIResponse(getMockRuntimeAPI(apiWithAPIPolicies, apiUUID, organiztion1, ())),
            k8sapiUUID,
            createdAPI.toBalString()
        ]
    };
    return data;
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
    MediationPolicy & readonly mediationPolicy1 =    {
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
    NotFoundError notfound = {body: {code: 909100, message: "6 not found."}};
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
    BadRequestError badRequestError = {"body": {"code": 90912, "message": "Invalid Sort By/Sort Order Value "}};
    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord name1"}};

    map<[string?, int, int, string, string, anydata]> dataSet = {
        "1": [
            (),
            10,
            0,
            SORT_BY_POLICY_NAME,
            SORT_ORDER_ASC,
            {
                "count": 4,
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
                    "total": 4
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
                "count": 4,
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
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 4
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
                "count": 4,
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
                    }
                ],
                "pagination": {
                    "offset": 0,
                    "limit": 10,
                    "total": 4
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
                "count": 4,
                "list": [
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
                    "total": 4
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
                    "limit": 2,
                    "total": 4
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
                    "total": 4
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
                    "total": 4
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
                    "total": 2
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
                    "total": 2
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
                    "total": 0
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

function getOKBackendServiceResponse(model:Service backendService) returns http:Response {
    http:Response backendServiceResponse = new;
    backendServiceResponse.statusCode = 201;
    model:Service serviceClone = backendService.clone();
    serviceClone.metadata.uid = uuid:createType1AsString();
    backendServiceResponse.setJsonPayload(serviceClone.toJson());
    return backendServiceResponse;
}

function getOKBackendPolicyResponse(model:BackendPolicy backendPolicy) returns http:Response {
    http:Response backendServiceResponse = new;
    backendServiceResponse.statusCode = 201;
    model:BackendPolicy backendPolicyClone = backendPolicy.clone();
    backendPolicyClone.metadata.uid = uuid:createType1AsString();
    backendServiceResponse.setJsonPayload(backendPolicyClone.toJson());
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
