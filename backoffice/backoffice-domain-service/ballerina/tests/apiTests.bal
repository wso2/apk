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

import ballerina/log;
import ballerina/test;
import wso2/apk_common_lib as commons;

@test:Config {dependsOn: [createAPITest]}
function getAPITest() {
    APIList|commons:APKError getAPI = getAPIList(25,0,"content:pizza","carbon.super");
    if getAPI is APIList {
        test:assertTrue(true, "Successfully retrieve APIs");
        log:printInfo(getAPI.toString());
    } else if getAPI is  commons:APKError {
        log:printError(getAPI.toString());
        test:assertFail("Error occured while creating API");
    }
}

@test:Config {dependsOn: [createAPITest]}
function getAPIByIdTest() {
    API|commons:APKError getAPIById = getAPI("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if getAPIById is API {
        test:assertTrue(true, "Successfully retrieve API");
        log:printInfo(getAPIById.toString());
    } else if getAPIById is  commons:APKError {
        log:printError(getAPIById.toString());
        test:assertFail("Error occured while creating API");
    }
}

@test:Config {dependsOn: [createAPITest]}
function getAPIDefinitionTest() {
    APIDefinition|commons:APKError getAPIDef = getAPIDefinition("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if getAPIDef is API {
        test:assertTrue(true, "Successfully retrieve API Definition");
        log:printInfo(getAPIDef.toString());
    } else if getAPIDef is  commons:APKError {
        log:printError(getAPIDef.toString());
        test:assertFail("Error occured while retrieve API Definition");
    }
}

@test:Config {dependsOn: [createAPITest]}
function updateAPITest() {
    ModifiableAPI payload = {
            "name": "PizzaShask",
            "description": "chnage description",
            "sdk": [
                "java", "android"
            ],
            "categories": [
                "cloud","open"
            ],
            "Policies": [
                "mysub1","mysub2"
            ]
        };
    API|commons:APKError updateAPICr = updateAPI("01ed75e2-b30b-18c8-wwf2-25da7edd2231", payload, "carbon.super");
    if updateAPICr is API {
        test:assertTrue(true, "Successfully Update API");
        log:printInfo(updateAPICr.toString());
    } else if updateAPICr is  commons:APKError {
        log:printError(updateAPICr.toString());
        test:assertFail("Error occured while updating API");
    }
}
