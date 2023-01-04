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

@test:Mock { functionName: "getAPIByIdDAO" }
test:MockFunction getAPIByIdDAOMock = new();

@test:Mock { functionName: "getAPIsDAO" }
test:MockFunction getAPIsDAOMock = new();

@test:Mock { functionName: "getAPIDefinitionDAO" }
test:MockFunction getAPIDefinitionDAOMock = new();

@test:Config {}
function getAPIByIdTest(){
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED"};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    string?|API|error apiResponse = getAPIByIdDAO("12sqwsqadasd","carbon.super");
    if apiResponse is API {
    test:assertTrue(true, "Successfully retrieved API");
    } else if apiResponse is  error {
        test:assertFail("Error occured while retrieving API");
    }
}

@test:Config {}
function getAPIListTest(){
    API[] apiList = [{name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED"}, 
    {name: "MyAPI2", context: "/myapi2", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED"}];
    test:when(getAPIsDAOMock).thenReturn(apiList);
    string?| APIList | error apiListReturned = getAPIList(0, 0, "", "carbon.super");
    if apiListReturned is APIList {
    test:assertTrue(true, "Successfully retrieved all APIs");
    } else if apiListReturned is  error {
        test:assertFail("Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIDefinitionByIdTest(){
    APIDefinition apiDefinition = {'type: "swagger",schemaDefinition: ""};
    test:when(getAPIDefinitionDAOMock).thenReturn(apiDefinition);
    APIDefinition|NotFoundError|error apiDefResponse = getAPIDefinition("12sqwsqadasd","carbon.super");
    if apiDefResponse is APIDefinition {
        test:assertTrue(true, "Successfully retrieved API Definition");
    } else if apiDefResponse is  error {
        test:assertFail("Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertFail("Definition Not Found Error");
    }
}

@test:Config {}
function getAPIDefinitionByIdNegativeTest(){
    NotFoundError nfe = {body:{code: 90915, message: "API Definition Not Found for provided API ID"}};
    test:when(getAPIDefinitionDAOMock).thenReturn(nfe);
    APIDefinition|NotFoundError|error apiDefResponse = getAPIDefinition("12sqwsqadasd","carbon.super");
    if apiDefResponse is APIDefinition {
        test:assertFail("Successfully retrieved API Definition");
    } else if apiDefResponse is  error {
        test:assertFail("Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertTrue(true,"Definition Not Found Error");
    }
}

@test:Config {}
function getAPIDefinitionByIdNegative2Test(){
    test:when(getAPIDefinitionDAOMock).thenReturn(error("Internal Error while retrieving API Definition"));
    APIDefinition|NotFoundError|error apiDefResponse = getAPIDefinition("12sqwsqadasd","carbon.super");
    if apiDefResponse is APIDefinition {
        test:assertFail("Successfully retrieved API Definition");
    } else if apiDefResponse is  error {
        test:assertTrue(true,"Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertFail("Definition Not Found Error");
    }
}
