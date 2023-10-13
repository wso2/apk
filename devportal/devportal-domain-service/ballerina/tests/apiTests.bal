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
import wso2/apk_common_lib as commons;
import ballerina/io;

commons:Organization organiztion = {
    name: "org1",
    displayName: "org1",
    uuid: "a3b58ccf-6ecc-4557-b5bb-0a35cce38256",
    organizationClaimValue: "org1",
    enabled: true,
    serviceListingNamespaces: ["*"],
    properties: []
};


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
    API|commons:APKError createdAPI = createAPIDAO(body,organiztion.uuid);
    if createdAPI is API {
    test:assertTrue(true, "Successfully created API");
        API|commons:APKError createdAPIDefinition = addDefinitionDAO(body,organiztion.uuid);
        if createdAPIDefinition is API {
        test:assertTrue(true, "Successfully created API");
        } else if createdAPIDefinition is  commons:APKError {
            log:printError(createdAPIDefinition.toString());
            test:assertFail("Error occured while creating API");
        }
    } else if createdAPI is  commons:APKError {
        log:printError(createdAPI.toString());
        test:assertFail("Error occured while creating API");
    }

}


@test:BeforeSuite
function beforeFunc2() returns error? {
     // Add thumbnail
    int|commons:APKError thumbnailCategoryId = getResourceCategoryIdByCategoryTypeDAO(RESOURCE_TYPE_THUMBNAIL);
    if thumbnailCategoryId is int {
        Resource thumbnail = {
            resourceUUID: "02ad95e2-b30b-10c8-wwf2-65da7edd2219",
            apiUuid: "01ed75e2-b30b-18c8-wwf2-25da7edd2231",
            resourceCategoryId: thumbnailCategoryId,
            dataType: "image/png",
            resourceContent: "thumbnail.png",
            resourceBinaryValue: check io:fileReadBytes("./tests/resources/thumbnail.png")

        };
        Resource|commons:APKError addedThumbnail = addResourceDAO(thumbnail);
       if addedThumbnail is Resource {
            test:assertTrue(true, "Successfully added Thumbnail");
        } else if addedThumbnail is commons:APKError {
            log:printError(addedThumbnail.toString());
            test:assertFail("Error occured while adding Thumbnail");
        }
    }
}

@test:Config {}
function getAPIByIdTest(){
    API|commons:APKError|NotFoundError apiResponse =getAPIByAPIId("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if apiResponse is API {
        test:assertTrue(true, "Successfully retrieved API");
    } else if apiResponse is  commons:APKError {
        test:assertFail("Error occured while retrieving API");
    }
}

@test:Config {}
function getAPIListTest(){
    APIList | commons:APKError apiListReturned = getAPIList(25, 0, null, organiztion, []);
    if apiListReturned is APIList {
        test:assertTrue(true, "Successfully retrieved all APIs");
    } else if apiListReturned is  commons:APKError {
        test:assertFail("Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIListContentSearchTest1(){
    APIList | commons:APKError apiListReturned = getAPIList(25, 0, "content:pizza", organiztion, []);
    if apiListReturned is APIList {
        test:assertTrue(true, "Successfully retrieved all APIs");
    } else if apiListReturned is  commons:APKError {
        test:assertFail("Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIListContentSearchTest2(){
    //Invalid Search Query without "content:" keyword
    APIList | commons:APKError apiListReturned = getAPIList(25, 0, "pizza", organiztion, []);
    if apiListReturned is APIList {
        test:assertFail("Successfully retrieved all APIs");
    } else if apiListReturned is  commons:APKError {
        test:assertTrue(true,"Successfully Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIListContentSearchTest3(){
    //Invalid Search Query without ":" 
    APIList | commons:APKError apiListReturned = getAPIList(25, 0, "contentpizza", organiztion, []);
    if apiListReturned is APIList {
        test:assertFail("Successfully retrieved all APIs");
    } else if apiListReturned is  commons:APKError {
        test:assertTrue(true,"Successfully Error occured while retrieving all APIs");
    }
}

@test:Config {}
function getAPIDefinitionByIdTest(){
    APIDefinition|NotFoundError|commons:APKError apiDefResponse = getAPIDefinition("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if apiDefResponse is APIDefinition {
        test:assertTrue(true, "Successfully retrieved API Definition");
    } else if apiDefResponse is  commons:APKError {
        log:printInfo(apiDefResponse.toBalString());
        log:printError(apiDefResponse.toString());
        test:assertFail("Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertFail("Definition Not Found Error");
    }
}

@test:Config {}
function getAPIDefinitionByIdNegativeTest(){
    APIDefinition|NotFoundError|commons:APKError apiDefResponse = getAPIDefinition("12sqwsqadasd");
    if apiDefResponse is APIDefinition {
        test:assertFail("Successfully retrieved API Definition");
    } else if apiDefResponse is  commons:APKError {
        test:assertFail("Error occured while retrieving API");
    } else if apiDefResponse is  NotFoundError {
        test:assertTrue(true,"Definition Not Found Error");
    }
}

@test:Config {}
function generateSDKImplTest(){
    http:Response|sdk:APIClientGenerationException|NotFoundError|commons:APKError sdk = generateSDKImpl("01ed75e2-b30b-18c8-wwf2-25da7edd2231","java");
    if sdk is http:Response {
        test:assertTrue(true, "Successfully generated API SDK");
    } else if sdk is sdk:APIClientGenerationException|commons:APKError{
        test:assertFail("Error while generating API SDK");
    }
}

@test:Config {}
function generateSDKImplTestNegative(){
    http:Response|sdk:APIClientGenerationException|NotFoundError|commons:APKError sdk = generateSDKImpl("12sqwsqadasd","java");
    if sdk is http:Response {
        test:assertFail("Successfully generated API SDK");
    } else if sdk is sdk:APIClientGenerationException|commons:APKError {
        log:printInfo(sdk.toBalString());
        test:assertTrue(true,"Error while generating API SDK");
    }
}

@test:Config {}
function gethumbnailTest() {
    http:Response|NotFoundError|commons:APKError thumbnail = getThumbnail("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if thumbnail is http:Response {
        test:assertTrue(true, "Successfully getting the thumbnail");
    } else {
        test:assertFail("Error occured while getting the thumbnail");
    }
}
