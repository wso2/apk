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
import devportal_service.org.wso2.apk.devportal.sdk as sdk;
import ballerina/http;

@test:BeforeSuite
function beforeFunc() {
    APIBody body = {
        "apiProperties":{
            "id": "01ed75e2-b30b-18c8-wwf2-25da7edd2231",
            "name":"PizzaShask",
            "context":"pizzssa",
            "version":"1.0.0",
            "provider":"admin",
            "lifeCycleStatus":"PUBLISHED",
            "type":"HTTP"
        },
        "Definition" : {	  
        "openapi": "3.0.0",
        "info": {
            "title": "Sample API",
            "description": "Optional multiline or single-line description in [CommonMark](http://commonmark.org/help/) or HTML.",
            "version": "0.1.9"
        },
        "servers": [
            {
            "url": "http://api.example.com/v1",
            "description": "Optional server description, e.g. Main (production) server"
            },
            {
            "url": "http://staging-api.example.com",
            "description": "Optional server description, e.g. Internal staging server for testing"
            }
        ],
        "paths": {
            "/users": {
            "get": {
            "summary": "Returns a list of users.",
            "description": "Optional extended description in CommonMark or HTML.",
            "responses": {
            "200": {
                "description": "A JSON array of user names",
                "content": {
                "application/json": {
                    "schema": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                    }
                }
                }
            }
            }
            }
            }
        }
        }
    };
    API|APKError createdAPI = createAPIDAO(body,"carbon.super");
    if createdAPI is API {
    test:assertTrue(true, "Successfully created API");
        API|APKError createdAPIDefinition = addDefinitionDAO(body,"carbon.super");
        if createdAPIDefinition is API {
        test:assertTrue(true, "Successfully created API");
        } else if createdAPIDefinition is  APKError {
            log:printError(createdAPIDefinition.toString());
            test:assertFail("Error occured while creating API");
        }
    } else if createdAPI is  APKError {
        log:printError(createdAPI.toString());
        test:assertFail("Error occured while creating API");
    }

}

@test:Config {}
function getAPIByIdTest(){
    API|APKError|NotFoundError apiResponse = getAPIByAPIId("01ed75e2-b30b-18c8-wwf2-25da7edd2231","carbon.super");
    if apiResponse is API {
        test:assertTrue(true, "Successfully retrieved API");
    } else if apiResponse is  APKError {
        test:assertFail("Error occured while retrieving API");
    }
}

@test:Config {}
function getAPIListTest(){
    APIList | APKError apiListReturned = getAPIList(25, 0, null, "carbon.super");
    if apiListReturned is APIList {
        test:assertTrue(true, "Successfully retrieved all APIs");
    } else if apiListReturned is  APKError {
        test:assertFail("Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIListContentSearchTest1(){
    APIList | APKError apiListReturned = getAPIList(25, 0, "content:pizza", "carbon.super");
    if apiListReturned is APIList {
        test:assertTrue(true, "Successfully retrieved all APIs");
    } else if apiListReturned is  APKError {
        test:assertFail("Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIListContentSearchTest2(){
    //Invalid Search Query without "content:" keyword
    APIList | APKError apiListReturned = getAPIList(25, 0, "pizza", "carbon.super");
    if apiListReturned is APIList {
        test:assertFail("Successfully retrieved all APIs");
    } else if apiListReturned is  APKError {
        test:assertTrue(true,"Successfully Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIListContentSearchTest3(){
    //Invalid Search Query without ":" 
    APIList | APKError apiListReturned = getAPIList(25, 0, "contentpizza", "carbon.super");
    if apiListReturned is APIList {
        test:assertFail("Successfully retrieved all APIs");
    } else if apiListReturned is  APKError {
        test:assertTrue(true,"Successfully Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIDefinitionByIdTest(){
    APIDefinition|NotFoundError|APKError apiDefResponse = getAPIDefinition("01ed75e2-b30b-18c8-wwf2-25da7edd2231","carbon.super");
    if apiDefResponse is APIDefinition {
        test:assertTrue(true, "Successfully retrieved API Definition");
    } else if apiDefResponse is  APKError {
        log:printError(apiDefResponse.toString());
        test:assertFail("Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertFail("Definition Not Found Error");
    }
}

@test:Config {}
function getAPIDefinitionByIdNegativeTest(){
    APIDefinition|NotFoundError|APKError apiDefResponse = getAPIDefinition("12sqwsqadasd","carbon.super");
    if apiDefResponse is APIDefinition {
        test:assertFail("Successfully retrieved API Definition");
    } else if apiDefResponse is  APKError {
        test:assertFail("Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertTrue(true,"Definition Not Found Error");
    }
}

@test:Config {}
function generateSDKImplTest(){
    http:Response|sdk:APIClientGenerationException|NotFoundError|APKError sdk = generateSDKImpl("01ed75e2-b30b-18c8-wwf2-25da7edd2231","java","carbon.super");
    if sdk is http:Response {
        test:assertTrue(true, "Successfully generated API SDK");
    } else if sdk is sdk:APIClientGenerationException|APKError{
        test:assertFail("Error while generating API SDK");
    }
}

@test:Config {}
function generateSDKImplTestNegative(){
    http:Response|sdk:APIClientGenerationException|NotFoundError|APKError sdk = generateSDKImpl("12sqwsqadasd","java","carbon.super");
    if sdk is http:Response {
        test:assertFail("Successfully generated API SDK");
    } else if sdk is sdk:APIClientGenerationException|APKError {
        test:assertTrue(true,"Error while generating API SDK");
    }
}
