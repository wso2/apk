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
    http:Client mock = test:mock(http:Client);
    test:prepare(mock).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/apis")
        .thenReturn(getMockAPIList());
    test:prepare(mock).when("get").withArguments("/api/v1/services")
        .thenReturn(getMockServiceList());
    test:prepare(mock).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/servicemappings")
        .thenReturn(getMockServiceMappings());
    test:prepare(mock).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b08-f2b1-1166-82d5-649ae706d29e").thenReturn(mock404Response());
    test:prepare(mock).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/namespaces/apk/apis/pizzashackAPI1").thenReturn(mock404Response());
    test:prepare(mock).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114a-definition").thenReturn(mockConfigMaps());
    http:ClientError clientError = error("Backend Failure");
    test:prepare(mock).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7b08-f2b1-1166-82d5-649ae706d29d-definition").thenReturn(mock404ConfigMap());
    test:prepare(mock).when("get").withArguments("/api/v1/namespaces/apk-platform/configmaps/01ed7aca-eb6b-1178-a200-f604a4ce114b-definition").thenReturn(clientError);
    return mock;
}

@test:Config {
    dataProvider: pathProvider
}
public function testretrievePathPrefix(string context, string 'version, string path, string expected) {
    APIClient apiclient = new ();
    string retrievePathPrefix = apiclient.retrievePathPrefix(context, 'version, path);
    test:assertEquals(retrievePathPrefix, expected);
}

function pathProvider() returns map<[string, string, string, string]>|error {
    map<[string, string, string, string]> dataSet = {
        "1": ["/abc/1.0.0", "1.0.0", "/abc", "/abc/1.0.0/abc"],
        "2": ["/abc", "1.0.0", "/abc", "/abc/1.0.0/abc"],
        "3": ["/abc", "1.0.0", "/*", "/abc/1.0.0"],
        "4": ["/abc/1.0.0", "1.0.0", "/*", "/abc/1.0.0"],
        "5": ["/abc/1.0.0", "1.0.0", "/{path}/abcd", "/abc/1.0.0/.*/abcd"],
        "6": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}", "/abc/1.0.0/path1/.*"],
        "7": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}/path2", "/abc/1.0.0/path1/.*/path2"],
        "8": ["/abc/1.0.0", "1.0.0", "/path1/{bcd}/{pathparam2}", "/abc/1.0.0/path1/.*/.*"],
        "9": ["/abc", "1.0.0", "/path1/*", "/abc/1.0.0/path1"]
    };
    return dataSet;
}

@test:Config {dataProvider: contextVersionDataProvider}
public function testValidateContextAndVersion(string context, string 'version, boolean expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.validateContextAndVersion(context, 'version), expected);
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
public function testValidateName(string name, boolean expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.validateName(name), expected);
}

function nameDataProvider() returns map<[string, boolean]>|error {
    map<[string, boolean]> dataSet = {
        "1": ["pizzashackAPI1", true],
        "2": ["pizzashackAPInew", false]
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
    test:assertEquals(getAPIByNameAndNamespace(name, namespace), expected);
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
public function testGetAPIById(string id, model:API & readonly|error expected) returns error? {
    test:assertEquals(getAPI(id), check expected);
}

function apiIDDataprovider() returns map<[string, model:API & readonly|error]>|error {
    map<[string, model:API & readonly|error]> dataSet = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", getMockPizzaShakK8sAPI()]
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
        "1": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/order/{orderId}", verb: "POST"}, "/v3/f77cc767/order"],
        "2": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767"}, {target: "/menu", verb: "GET"}, "/v3/f77cc767/menu"],
        "3": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/menu", verb: "GET"}, "/v3/f77cc767/menu"],
        "4": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: false, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/*", verb: "GET"}, "/v3/f77cc767/"],
        "5": [{name: "pizzaAPI", context: "/pizza1234", 'version: "1.0.0"}, {name: "service1", namespace: "apk-platform", port: 443, serviceEntry: true, url: "https://run.mocky.io/v3/f77cc767/"}, {target: "/*", verb: "GET"}, "/"]
    };
    return dataSet;
}

@test:Config {dataProvider: apiDefinitionDataProvider}
public function testGetAPIDefinitionByID(string apiid, anydata expectedResponse) returns error? {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getAPIDefinitionByID(apiid), expectedResponse);
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
public function testgetApiById(string apiid, anydata expectedData) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getAPIById(apiid), expectedData);
}

public function apiByIdDataProvider() returns map<[string, API|NotFoundError]> {
    API & readonly api1 = {name: "pizzashackAPI", context: "/pizzashack/1.0.0", 'version: "1.0.0", id: "c5ab2423-b9e8-432b-92e8-35e6907ed5e8", createdTime: "2022-12-13T09:45:47Z"};
    NotFoundError notfound = {body: {code: 909100, message: "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found."}};
    map<[string, API & readonly|NotFoundError]> dataset = {
        "1": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e8", api1],
        "2": ["c5ab2423-b9e8-432b-92e8-35e6907ed5e9", notfound]
    };
    return dataset;
}

@test:Config {dataProvider: getApilistDataProvider}
public function testGetAPIList(string? query, int 'limit, int offset, string sortBy, string sortOrder, anydata expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getAPIList(query, 'limit, offset, sortBy, sortOrder), expected);
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
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "8a1eb4f9-efab-4682-a051-4df4050812d2",
                        "name": "DemoAPI",
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/demoapi/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
            "type:HTTP",
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
                        "type": "HTTP",
                        "createdTime": "2022-12-14T18:51:26Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
                        "name": "pizzashackAPI",
                        "context": "/pizzashack/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",
                        "name": "pizzashackAPI1",
                        "context": "/pizzashack1/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T17:09:49Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                        "name": "pizzashackAPI11",
                        "context": "/pizzashack11/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1",
                        "name": "pizzashackAPI12",
                        "context": "/pizzashack12/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
                        "createdTime": "2022-12-13T09:45:47Z"
                    },
                    {
                        "id": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3",
                        "name": "pizzashackAPI13",
                        "context": "/pizzashack13/1.0.0",
                        "version": "1.0.0",
                        "type": "HTTP",
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
public function testRetrieveGeneratedSwaggerDefinition(API api, anydata expectedOutput) {
    APIClient apiclient = new;
    test:assertEquals(apiclient.retrieveGeneratedSwaggerDefinition(api), expectedOutput);
}

function testDataGeneratedSwaggerDefinition() returns map<[API, json|APKError]> {
    map<[API, json|APKError]> data = {
        "1": [
            {
                "name": "demoAPI",
                "context": "/demoAPI/1.0.0",
                "version": "1.0.0"
            },
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
    }
;
    return data;
}

@test:Config{dataProvider:validateExistenceDataProvider}
function testValidateAPIExistence(string query,anydata expected) {
    APIClient apiClient = new;
    test:assertEquals(apiClient.validateAPIExistence(query),expected);
}
function validateExistenceDataProvider() returns map<[string,NotFoundError|BadRequestError|http:Ok]>{
http:Ok ok = {};
BadRequestError badRequest = {body: {code: 90912, message: "Invalid KeyWord type"}};
NotFoundError notFound = {body: {code: 900914, message: "context/name doesn't exist"}};
map<[string,NotFoundError|BadRequestError|http:Ok]> data = {
"1":["name:pizzashackAPI",ok],
"2":["name:mockapi",notFound],
"3":["context:/api/v1",notFound],
"4":["context:/pizzashack/1.0.0",ok],
"5":["pizzashackAPI",ok],
"6":["type:pizzashackAPI",badRequest]
};
return data;
}