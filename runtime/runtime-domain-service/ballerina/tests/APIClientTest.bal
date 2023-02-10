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

@test:Mock {functionName: "startAndAttachServices"}
function getMockStartandAttachServices() returns error? {
}

@test:Mock {functionName: "getBackendPolicyUid"}
function testgetBackendPolicyUid(API api, string? endpointType, string organization) returns string {
    return "backendpolicy-uuid";
}

@test:Mock {functionName: "getServiceMappingClient"}
function getMockServiceMappingClient(string resourceVersion) returns websocket:Client|error|() {
    string initialConectionId = uuid:createType1AsString();
    websocket:Client mock;
    if resourceVersion == "39433" {
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturnSequence(true, true, false);
        test:prepare(mock).when("getConnectionId").thenReturn(initialConectionId);
        test:prepare(mock).when("readMessage").thenReturn(getServiceMappingEvent());
    } else if resourceVersion == "5834" {
        string connectionId = uuid:createType1AsString();
        mock = test:mock(websocket:Client);
        test:prepare(mock).when("isOpen").thenReturn(true);
        test:prepare(mock).when("getConnectionId").thenReturn(connectionId);
        test:prepare(mock).when("readMessage").thenReturnSequence(getNextServiceMappingEvent(), ());
    } else {
        mock = test:mock(websocket:Client);
    }

    return mock;
}

@test:Mock {functionName: "getServiceClient"}
function getMockServiceClient(string resourceVersion) returns websocket:Client|error|() {
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
    } else {
        mock = test:mock(websocket:Client);
    }
    return mock;
}

@test:Mock {functionName: "getClient"}
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
    } else {
        mock = test:mock(websocket:Client);
    }
    return mock;
}

@test:Mock {
    functionName: "initializeK8sClient"
}
function getMockK8sClient() returns http:Client {
    http:Client mockK8sClient = test:mock(http:Client);
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/apis")
        .thenReturn(getMockAPIList());
    string fieldSlector = "metadata.namespace%21%3Dkube-system%2Cmetadata.namespace%21%3Dkubernetes-dashboard%2Cmetadata.namespace%21%3Dgateway-system%2Cmetadata.namespace%21%3Dingress-nginx%2Cmetadata.namespace%21%3Dapk-platform";
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/services?fieldSelector=" + fieldSlector)
        .thenReturn(getMockServiceList());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/servicemappings")
        .thenReturn(getMockServiceMappings());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b08-f2b1-1166-82d5-649ae706d29e").thenReturn(mock404Response());
    test:prepare(mockK8sClient).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk/apis/pizzashackAPI1").thenReturn(mock404Response());
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114a-definition").thenReturn(mockConfigMaps());
    http:ClientError clientError = error("Backend Failure");
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7b08-f2b1-1166-82d5-649ae706d29d-definition").thenReturn(mock404ConfigMap());
    test:prepare(mockK8sClient).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114b-definition").thenReturn(clientError);
    return mockK8sClient;
}

@test:Config {
    dataProvider: pathProvider
}
public function testretrievePathPrefix(string context, string 'version, string path, string expected) {
    APIClient apiclient = new ();
    string retrievePathPrefix = apiclient.retrievePathPrefix(context, 'version, path, "carbon.super");
    test:assertEquals(retrievePathPrefix, expected);
}

function pathProvider() returns map<[string, string, string, string]>|error {
    map<[string, string, string, string]> dataSet = {
        "1": ["/abc/1.0.0", "1.0.0", "/abc", "/t/carbon.super/abc/1.0.0/abc"],
        "2": ["/abc", "1.0.0", "/abc", "/t/carbon.super/abc/1.0.0/abc"],
        "3": ["/abc", "1.0.0", "/*", "/t/carbon.super/abc/1.0.0(.*)"],
        "4": ["/abc/1.0.0", "1.0.0", "/*", "/t/carbon.super/abc/1.0.0(.*)"],
        "5": ["/abc/1.0.0", "1.0.0", "/{path}/abcd", "/t/carbon.super/abc/1.0.0/(.*)/abcd"],
        "6": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}", "/t/carbon.super/abc/1.0.0/path1/(.*)"],
        "7": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}/path2", "/t/carbon.super/abc/1.0.0/path1/(.*)/path2"],
        "8": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}/{pathparam2}", "/t/carbon.super/abc/1.0.0/path1/(.*)/(.*)"],
        "9": ["/abc", "1.0.0", "/path1/*", "/t/carbon.super/abc/1.0.0/path1(.*)"]
    };
    return dataSet;
}

@test:Config {dataProvider: contextVersionDataProvider}
public function testValidateContextAndVersion(string context, string 'version, boolean expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.validateContextAndVersion(context, 'version, "carbon.super"), expected);
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
public function testValidateName(string name, string organization, boolean expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.validateName(name, organization), expected);
}

function nameDataProvider() returns map<[string, string, boolean]>|error {
    map<[string, string, boolean]> dataSet = {
        "1": ["pizzashackAPI1", "carbon.super", true],
        "2": ["pizzashackAPInew", "carbon.super", false],
        "3": ["pizzashackAPI1", "wso2.com", false]

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
    test:assertEquals(getAPIByNameAndNamespace(name, namespace, "carbon.super"), expected);
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
public function testGetAPIById(string id, string organization, model:API & readonly|error expected) returns error? {
    test:assertEquals(getAPI(id, organization), check expected);
}

function apiIDDataprovider() returns map<[string, string, model:API & readonly|error]>|error {

    map<[string, string, model:API & readonly|error]> dataSet = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", "carbon.super", getMockPizzaShakK8sAPI()]
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
    test:assertEquals(apiclient.getAPIDefinitionByID(apiid, "carbon.super"), expectedResponse);
}

public function apiDefinitionDataProvider() returns map<[string, json|NotFoundError|PreconditionFailedError|InternalServerErrorError]> {
    NotFoundError notfound = {body: {code: 909100, message: "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found."}};
    InternalServerErrorError internalError = {body: {code: 909000, message: "Internal Error Occured while retrieving definition"}};

    map<[string, json|NotFoundError|PreconditionFailedError|InternalServerErrorError]> dataSet = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", mockOpenAPIJson()],
        "2": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e9", notfound],
        "3": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f9", mockpizzashackAPI11Definition()],
        "4": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f1", mockPizzashackAPI12Definition()],
        "5": ["7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1", mockPizzaShackAPI1Definition()],
        "6": ["c5ab2423-b9e8-432b-92e8-35e6907ed5f3", internalError]

    };
    return dataSet;
}

@test:Config {dataProvider: apiByIdDataProvider}
public function testgetApiById(string apiid, string organization, anydata expectedData) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getAPIById(apiid, organization), expectedData);
}

public function apiByIdDataProvider() returns map<[string, string, API|NotFoundError]> {
    API & readonly api1 = {name: "pizzashackAPI", context: "/t/carbon.super/pizzashack/1.0.0", 'version: "1.0.0", id: "c5ab2423-b9e8-432b-92e8-35e6907ed5e8", createdTime: "2022-12-13T09:45:47Z"};
    NotFoundError notfound = {body: {code: 909100, message: "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found."}};
    map<[string, string, API & readonly|NotFoundError]> dataset = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", "carbon.super", api1],
        "2": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e9", "carbon.super", notfound]
    };
    return dataset;
}

@test:Config {dataProvider: getApilistDataProvider}
public function testGetAPIList(string? query, int 'limit, int offset, string sortBy, string sortOrder, anydata expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getAPIList(query, 'limit, offset, sortBy, sortOrder, "carbon.super"), expected);
}

function getApilistDataProvider() returns map<[string?, int, int, string, string, APIList|InternalServerErrorError|BadRequestError]> {
    BadRequestError badRequestError = {"body": {"code": 90912, "message": "Invalid Sort By/Sort Order Value "}};
    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord type1"}};

    map<[string?, int, int, string, string, APIList|InternalServerErrorError|BadRequestError]> dataSet = {
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
                        "context": "/t/carbon.super/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
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
            }
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
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/t/carbon.super/demoapi/1.0.0",
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
            }
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
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/t/carbon.super/demoapi/1.0.0",
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
            }
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
                        "context": "/t/carbon.super/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
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
            }
        ],
        "5": [(), 10, 0, "description", SORT_ORDER_DESC, badRequestError],
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
                        "context": "/t/carbon.super/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
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
            }
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
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
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
            }
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
            }
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
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
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
            }
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
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
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
            }
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
                        "context": "/t/carbon.super/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/t/carbon.super/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/t/carbon.super/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/t/carbon.super/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/t/carbon.super/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "REST",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/t/carbon.super/pizzashack13/1.0.0",
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
            }
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
            }
        ],
        "13": [
            "type1:WS",
            10,
            0,
            SORT_BY_API_NAME,
            SORT_ORDER_ASC,
            badRequest
        ]
    };
    return dataSet;
}

@test:Config {dataProvider: testDataGeneratedSwaggerDefinition}
public function testRetrieveGeneratedSwaggerDefinition(API api, string? definition, anydata expectedOutput) {
    APIClient apiclient = new;
    test:assertEquals(apiclient.retrieveGeneratedSwaggerDefinition(api, definition), expectedOutput);
}

function testDataGeneratedSwaggerDefinition() returns map<[API, string?, json|APKError]> {
    map<[API, string?, json|APKError]> data = {
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
    test:assertEquals(apiClient.validateAPIExistence(query), expected);
}

function validateExistenceDataProvider() returns map<[string, NotFoundError|BadRequestError|http:Ok]> {
    http:Ok ok = {};
    BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord type"}};
    NotFoundError notFound = {body: {code: 900914, message: "context/name doesn't exist"}};
    map<[string, NotFoundError|BadRequestError|http:Ok]> data = {
        "1": ["name:pizzashackAPI", ok],
        "2": ["name:mockapi", notFound],
        "3": ["context:/api/v1", notFound],
        "4": ["context:/pizzashack/1.0.0", ok],
        "5": ["pizzashackAPI", ok],
        "6": ["type:pizzashackAPI", badRequest]
    };
    return data;
}

@test:Config {dataProvider: createApiFromServiceDataProvider}
function testCreateAPIFromService(string serviceUUId, string apiUUID, [model:ConfigMap, any] configmapResponse, [model:Httproute, any] httproute, [model:K8sServiceMapping, any] servicemapping, [model:API, any] k8sAPI, API api, string k8sapiUUID, anydata expected) {
    test:prepare(k8sApiServerEp).when("post").withArguments("/api/v1/namespaces/apk-platform/configmaps", configmapResponse[0]).thenReturn(configmapResponse[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/gateway.networking.k8s.io/v1beta1/namespaces/apk-platform/httproutes", httproute[0]).thenReturn(httproute[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings", servicemapping[0]).thenReturn(servicemapping[1]);
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis", k8sAPI[0]).thenReturn(k8sAPI[1]);
    APIClient apiClient = new;
    test:assertEquals(apiClient.createAPIFromService(serviceUUId, api, "carbon.super"), expected);
}

function createApiFromServiceDataProvider() returns map<[string, string, [model:ConfigMap, any], [model:Httproute, any], [model:K8sServiceMapping, any], [model:API, any], API, string, CreatedAPI|BadRequestError|InternalServerErrorError|APKError]> {
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
    string apiUUID = getUniqueIdForAPI(api.name, api.'version, "carbon.super");
    model:ConfigMap configmap = getMockConfigMap1(apiUUID, api);
    http:Response mockConfigMapResponse = getMockConfigMapResponse(configmap.clone());
    model:Httproute httpRoute = getMockHttpRoute(api, apiUUID);
    http:Response httpRouteResponse = getMockHttpRouteResponse(httpRoute.clone());
    model:K8sServiceMapping mockServiceMappingRequest = getMockServiceMappingRequest(api, apiUUID);
    model:API mockAPI = getMockAPI(api, apiUUID, "carbon.super");
    http:Response mockAPIResponse = getMockAPIResponse(mockAPI.clone(), k8sAPIUUID1);
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
            "id": k8sAPIUUID1,
            "name": "PizzaAPI",
            "context": "/t/carbon.super/pizzaAPI/1.0.0",
            "version": "1.0.0",
            "type": "REST"
        }
    };
    map<[string, string, [model:ConfigMap, any], [model:Httproute, any], [model:K8sServiceMapping, any], [model:API, any], API, string, CreatedAPI|BadRequestError|InternalServerErrorError|APKError]> data = {
        "1": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], api, k8sAPIUUID1, createdAPI],
        "2": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], alreadyNameExist, k8sAPIUUID1, nameAlreadyExistError],
        "3": ["275b00d1-722c-4df2-b65a-9b14677abe4b", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], contextAlreadyExist, k8sAPIUUID1, contextAlreadyExistError],
        "4": ["275b00d1-722c-4df2-b65a-9b14677abe4a", apiUUID, [configmap, mockConfigMapResponse], [httpRoute, httpRouteResponse], [mockServiceMappingRequest, serviceMappingResponse], [mockAPI, mockAPIResponse], api, k8sAPIUUID1, serviceNotExist]
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
            "context": apiClient.returnFullContext(api.context, api.'version, organization),
            "organization": organization,
            "definitionFileRef": apiUUID + "-definition",
            "prodHTTPRouteRef": apiUUID + "-production"
        },
        "status": null
    };
    return k8sapi;
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
            "context": apiClient.returnFullContext(api.context, api.'version, organization),
            "organization": organization,
            "definitionFileRef": apiUUID + "-definition",
            "sandHTTPRouteRef": apiUUID + "-sandbox"
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

function getMockHttpRoute(API api, string apiUUID) returns model:Httproute {
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": apiUUID + "-production", "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
        "spec": {
            "hostnames": ["gw.wso2.com"],
            "rules": [
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"}, "method": "GET"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"}, "method": "PUT"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"}, "method": "POST"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"}, "method": "DELETE"}],
                    "filters": [{"type": "URLRewrite", "urlRewrite": {"path": {"type": "ReplaceFullPath", "replaceFullPath": "\\1"}}}],
                    "backendRefs": [{"weight": 1, "group": "", "kind": "Service", "name": "backend", "namespace": "apk", "port": 80}]
                },
                {
                    "matches": [{"path": {"type": "RegularExpression", "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"}, "method": "PATCH"}],
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
        model:API k8sApi, any k8sapiResponse,
        string k8sapiUUID, anydata expected) {
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
    test:prepare(k8sApiServerEp).when("post").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis", k8sApi).thenReturn(k8sapiResponse);
    APKError|CreatedAPI|BadRequestError aPI = apiClient.createAPI(api, (), "carbon.super");
    if aPI is BadRequestError || aPI is CreatedAPI {
        test:assertEquals(aPI, expected);
    } else if aPI is APKError {
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

function getMockHttpRouteWithBackend(API api, string apiUUID, string backenduuid, string 'type) returns model:Httproute {
    string hostnames = 'type == PRODUCTION_TYPE ? "gw.wso2.com" : "sandbox.gw.wso2.com";
    return {
        "apiVersion": "gateway.networking.k8s.io/v1beta1",
        "kind": "HTTPRoute",
        "metadata": {"name": apiUUID + "-" + 'type, "namespace": "apk-platform", "labels": {"api-name": api.name, "api-version": api.'version}},
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
                                "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"
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
                                "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"
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
                                "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"
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
                                "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"
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
                                "value": "/t/carbon.super/pizzaAPI/1.0.0(.*)"
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

function createAPIDataProvider() returns map<[string, string, API, model:ConfigMap, any, model:Httproute?, any, model:Httproute?, any, [model:Service, any][], [model:BackendPolicy, any][], model:API, any, string, string|CreatedAPI|BadRequestError]> {
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
    string apiUUID = getUniqueIdForAPI(api.name, api.'version, "carbon.super");
    string backenduuid = getBackendServiceUid(api, (), PRODUCTION_TYPE, "carbon.super");
    string backenduuid1 = getBackendServiceUid(api, (), SANDBOX_TYPE, "carbon.super");
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
    model:Httproute prodhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid, PRODUCTION_TYPE);
    model:Httproute sandhttpRoute = getMockHttpRouteWithBackend(api, apiUUID, backenduuid1, SANDBOX_TYPE);

    CreatedAPI createdAPI = {body: {name: "PizzaAPI", context: "/t/carbon.super/pizzaAPI/1.0.0", 'version: "1.0.0", id: k8sapiUUID}};
    APKError productionEndpointNotSpecifiedError = error("Production Endpoint Not specified", message = "Endpoint Not specified", description = "Production Endpoint Not specified", code = 90911, statusCode = "400");
    APKError sandboxEndpointNotSpecifiedError = error("Sandbox Endpoint Not specified", message = "Endpoint Not specified", description = "Sandbox Endpoint Not specified", code = 90911, statusCode = "400");
    APKError k8sLevelError = error("Internal Error occured while deploying API", code = 909000, message
        = "Internal Error occured while deploying API", statusCode = "500", description = "Internal Error occured while deploying API", moreInfo = {});
    APKError invalidAPINameError = error("Invalid API Name", code = 90911, message = "Invalid API Name", statusCode = "400", description = "API Name PizzaAPI Invalid", moreInfo = {});
    map<[string, string, API, model:ConfigMap,
    any, model:Httproute|(), any, model:Httproute|(),
    any, [model:Service, any][], [model:BackendPolicy, any][], model:API, any, string,
    string|CreatedAPI|BadRequestError]> data = {
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            createdAPI
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            nameAlreadyExistError
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            contextAlreadyExistError
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
            getMockAPI1(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI1(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            createdAPI
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            k8sLevelError.toBalString()
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            k8sLevelError.toBalString()
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
            k8sapiUUID,
            k8sLevelError.toBalString()
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIResponse(getMockAPI(api, apiUUID, "carbon.super"), k8sapiUUID),
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIErrorResponse(),
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
            getMockAPI(api, apiUUID, "carbon.super"),
            getMockAPIErrorNameExist(),
            k8sapiUUID,
            invalidAPINameError.toBalString()
        ]
    };
    return data;
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
