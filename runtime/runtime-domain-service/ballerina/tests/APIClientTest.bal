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
public function testGetDomainPath(string url, string expectedDomain, string expectedPath,int expectedPort,string expectedHost) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getDomain(url), expectedDomain);
    test:assertEquals(apiclient.getPath(url), expectedPath);
    test:assertEquals(apiclient.getPort(url), expectedPort);
    test:assertEquals(apiclient.gethost(url), expectedHost);
}

function hostnameDataProvider() returns map<[string, string, string,int,string]>|error {
    map<[string, string, string,int,string]> dataSet = {
        "1": ["https://localhost/api.json", "https://localhost", "/api.json",443,"localhost"],
        "2": ["http://localhost/api.json", "http://localhost", "/api.json",80,"localhost"],
        "3": ["https://localhost:443/api.json", "https://localhost:443", "/api.json",443,"localhost"],
        "4": ["http://localhost:80/api.json", "http://localhost:80", "/api.json",80,"localhost"],
        "5": ["https://localhost", "https://localhost", "",443,"localhost"],
        "6": ["http://localhost", "http://localhost", "",80,"localhost"],
        "7": ["https://localhost:443", "https://localhost:443", "",443,"localhost"],
        "8": ["http://localhost:80", "http://localhost:80", "",80,"localhost"],
        "9": ["tcp://localhost:443", "", "",-1,""]
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

@test:Config{dataProvider:prefixMatchDataProvider}
public function testGeneratePrefixMatch(API api,model:Endpoint endpoint,APIOperations apiOperation,string expected) {
    APIClient apiclient = new ();
    test:assertEquals(apiclient.generatePrefixMatch(api,endpoint,apiOperation,PRODUCTION_TYPE),expected);

}
function prefixMatchDataProvider() returns map<[API,model:Endpoint,APIOperations,string]>{
map<[API,model:Endpoint,APIOperations,string]> dataSet = {
    "1":[{name: "pizzaAPI",context: "/pizza1234",'version: "1.0.0"},{name: "service1",namespace: "apk-platform",port: 443,serviceEntry: false,url: "https://run.mocky.io/v3/f77cc767"},{target: "/order/{orderId}",verb: "POST"},"/v3/f77cc767/order"],
    "2":[{name: "pizzaAPI",context: "/pizza1234",'version: "1.0.0"},{name: "service1",namespace: "apk-platform",port: 443,serviceEntry: false,url: "https://run.mocky.io/v3/f77cc767"},{target: "/menu",verb: "GET"},"/v3/f77cc767/menu"],
    "3":[{name: "pizzaAPI",context: "/pizza1234",'version: "1.0.0"},{name: "service1",namespace: "apk-platform",port: 443,serviceEntry: false,url: "https://run.mocky.io/v3/f77cc767/"},{target: "/menu",verb: "GET"},"/v3/f77cc767/menu"],
    "4":[{name: "pizzaAPI",context: "/pizza1234",'version: "1.0.0"},{name: "service1",namespace: "apk-platform",port: 443,serviceEntry: false,url: "https://run.mocky.io/v3/f77cc767/"},{target: "/*",verb: "GET"},"/v3/f77cc767/"],
    "5":[{name: "pizzaAPI",context: "/pizza1234",'version: "1.0.0"},{name: "service1",namespace: "apk-platform",port: 443,serviceEntry: true,url: "https://run.mocky.io/v3/f77cc767/"},{target: "/*",verb: "GET"},"/"]
};
return dataSet;
}

@test:Config{dataProvider:apiDefinitionDataProvider}
public function testGetAPIDefinitionByID(string apiid,anydata expectedResponse) returns error?{
    APIClient apiclient = new ();
    test:assertEquals(apiclient.getAPIDefinitionByID(apiid),expectedResponse);
}
public function apiDefinitionDataProvider() returns map<[string,json|NotFoundError|PreconditionFailedError|InternalServerErrorError]>{
    NotFoundError notfound = {body: {code: 909100, message: "c5ab2423-b9e8-432b-92e8-35e6907ed5e9 not found."}};
    InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error Occured while retrieving definition"}};

    map<[string,json|NotFoundError|PreconditionFailedError|InternalServerErrorError]> dataSet = {
        "1":["c5ab2423-b9e8-432b-92e8-35e6907ed5e8",mockOpenAPIJson()],
        "2":["c5ab2423-b9e8-432b-92e8-35e6907ed5e9",notfound],
        "3":["c5ab2423-b9e8-432b-92e8-35e6907ed5f9",mockpizzashackAPI11Definition()],
        "4":["c5ab2423-b9e8-432b-92e8-35e6907ed5f1",mockPizzashackAPI12Definition()],
        "5":["7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1",mockPizzaShackAPI1Definition()],
        "6":["c5ab2423-b9e8-432b-92e8-35e6907ed5f3",internalError]

    };
    return dataSet;
}