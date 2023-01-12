//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
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

@test:Mock { functionName: "checkAPICategoryExistsByNameDAO" }
test:MockFunction checkAPICategoryExistsByNameDAOMock = new();

@test:Mock { functionName: "addAPICategoryDAO" }
test:MockFunction addAPICategoryDAOMock = new();

@test:Mock { functionName: "getAPICategoriesDAO" }
test:MockFunction getAPICategoriesDAOMock = new();

@test:Mock { functionName: "getAPICategoryByIdDAO" }
test:MockFunction getAPICategoryByIdDAOMock = new();

@test:Mock { functionName: "updateAPICategoryDAO" }
test:MockFunction updateAPICategoryDAOMock = new();

@test:Mock { functionName: "deleteAPICategoryDAO" }
test:MockFunction deleteAPICategoryDAOMock = new();

@test:Config {}
function addAPICategoryTest() {
    APICategory|error  apiCategory = {name: "MyCat1", description: "My Desc 1", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    APICategory payload = {name: "MyCat1", description: "My Desc 1"};
    test:when(checkAPICategoryExistsByNameDAOMock).thenReturn(false);
    test:when(addAPICategoryDAOMock).thenReturn(apiCategory);
    CreatedAPICategory|error createdApiCategory = addAPICategory(payload);
    if createdApiCategory is CreatedAPICategory {
        test:assertTrue(true,"API Category added successfully");
    } else if createdApiCategory is error {
        test:assertFail("Error occured while adding API Category");
    }
}

@test:Config {}
function addAPICategoryTestNegative1() {
    APICategory|error  apiCategory = {name: "MyCat1", description: "My Desc 1", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    APICategory payload = {name: "MyCat1", description: "My Desc 1"};
    //API Category Name alrady exisitng
    test:when(checkAPICategoryExistsByNameDAOMock).thenReturn(true);
    test:when(addAPICategoryDAOMock).thenReturn(apiCategory);
    CreatedAPICategory|error createdApiCategory = addAPICategory(payload);
    if createdApiCategory is CreatedAPICategory {
        test:assertFail("API Category added successfully");
    } else if createdApiCategory is error {
        test:assertTrue(true, "Error occured while adding API Category");
    }
}

@test:Config {}
function addAPICategoryTestNegative2() {
    //API Category is an error
    error  apiCategory = error("error");
    APICategory payload = {name: "MyCat1", description: "My Desc 1"};
    test:when(checkAPICategoryExistsByNameDAOMock).thenReturn(false);
    test:when(addAPICategoryDAOMock).thenReturn(apiCategory);
    CreatedAPICategory|error createdApiCategory = addAPICategory(payload);
    if createdApiCategory is CreatedAPICategory {
        test:assertFail("API Category added successfully");
    } else if createdApiCategory is error {
        test:assertTrue(true, "Error occured while adding API Category");
    }
}

@test:Config {}
function getAllCategoryListTest() {
    APICategory[]  apiCategories = [{name: "MyCat1", description: "My Desc 1",
     id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2},
    {name: "MyCat2", description: "My Desc 2", 
    id: "01ed9235-49f3-1f4e-a3ac-e440022f0c5e", numberOfAPIs: 0}];
    test:when(getAPICategoriesDAOMock).thenReturn(apiCategories);
    APICategoryList|error apiCategoryList = getAllCategoryList();
    if apiCategoryList is APICategoryList {
        test:assertTrue(true,"API Category list retrieved successfully");
    } else if apiCategoryList is error {
        test:assertFail("Error occured while retrieving API Category List");
    }
}

@test:Config {}
function getAllCategoryListTestNegative1() {
    //API Category is an error
    error  apiCategories = error("error");
    test:when(getAPICategoriesDAOMock).thenReturn(apiCategories);
    APICategoryList|error apiCategoryList = getAllCategoryList();
    if apiCategoryList is APICategoryList {
        test:assertFail("API Category list retrieved successfully");
    } else if apiCategoryList is error {
        test:assertTrue(true,"Error occured while retrieving API Category List");
    }
}

@test:Config {}
function updateAPICategoryTest() {
    APICategory|error  apiCategory = {name: "MyCat1", description: "My Desc 1 new", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    APICategory|error  exisitingApiCategory = {name: "MyCat1", description: "My Desc 1", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    APICategory payload = {name: "MyCat1", description: "My Desc 1 new"};
    test:when(getAPICategoryByIdDAOMock).thenReturn(exisitingApiCategory);
    test:when(checkAPICategoryExistsByNameDAOMock).thenReturn(false);
    test:when(updateAPICategoryDAOMock).thenReturn(apiCategory);
    APICategory|NotFoundError|error createdApiCategory = updateAPICategory("01ed9241-2d5d-1b98-8ecb-40f85676b081",payload);
    if createdApiCategory is APICategory {
        test:assertTrue(true,"API Category updated successfully");
    } else if createdApiCategory is NotFoundError {
        test:assertFail("Not Found Error");
    } else if createdApiCategory is error {
        test:assertFail("Error occured while adding API Category");
    }
}

@test:Config {}
function updateAPICategoryTestNegative1() {
    APICategory|error  apiCategory = {name: "MyCat1", description: "My Desc 1 new", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    // Exisiting API Category is an error
    error exisitingApiCategory = error("error");
    APICategory payload = {name: "MyCat1", description: "My Desc 1 new"};
    test:when(getAPICategoryByIdDAOMock).thenReturn(exisitingApiCategory);
    test:when(checkAPICategoryExistsByNameDAOMock).thenReturn(false);
    test:when(updateAPICategoryDAOMock).thenReturn(apiCategory);
    APICategory|NotFoundError|error createdApiCategory = updateAPICategory("01ed9241-2d5d-1b98-8ecb-40f85676b081",payload);
    if createdApiCategory is APICategory {
        test:assertFail("API Category updated successfully");
    } else if createdApiCategory is NotFoundError {
        test:assertTrue(true, "Not Found Error");
    } else if createdApiCategory is error {
        test:assertFail("Error occured while adding API Category");
    }
}

@test:Config {}
function updateAPICategoryTestNegative2() {
    APICategory|error  apiCategory = {name: "MyCat1", description: "My Desc 1 new", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    APICategory|error  exisitingApiCategory = {name: "MyCat1", description: "My Desc 1", id: "01ed9241-2d5d-1b98-8ecb-40f85676b081", numberOfAPIs: 2};
    // New Name
    APICategory payload = {name: "MyCat1New", description: "My Desc 1 new"};
    test:when(getAPICategoryByIdDAOMock).thenReturn(exisitingApiCategory);
    // Another API Category by same name exists
    test:when(checkAPICategoryExistsByNameDAOMock).thenReturn(true);
    test:when(updateAPICategoryDAOMock).thenReturn(apiCategory);
    APICategory|NotFoundError|error createdApiCategory = updateAPICategory("01ed9241-2d5d-1b98-8ecb-40f85676b081",payload);
    if createdApiCategory is APICategory {
        test:assertFail("API Category updated successfully");
    } else if createdApiCategory is NotFoundError {
        test:assertFail("Not Found Error");
    } else if createdApiCategory is error {
        test:assertTrue(true,"Error occured while adding API Category");
    }
}

@test:Config {}
function removeAPICategoryTest(){
    test:when(deleteAPICategoryDAOMock).withArguments("01ed9241-2d5d-1b98-8ecb-40f85676b081","carbon.super").thenReturn("");
    error?|string status = removeAPICategory("01ed9241-2d5d-1b98-8ecb-40f85676b081");
    if status is string {
    test:assertTrue(true, "Successfully deleted API Category");
    } else if status is  error {
        test:assertFail("Error occured while deleting API Category");
    }
}

@test:Config {}
function removeAPICategoryTestNegative(){
    test:when(deleteAPICategoryDAOMock).withArguments("01ed9241-2d5d-1b98-8ecb-40f85676b081","carbon.super").thenReturn(error("error"));
    error?|string status = removeAPICategory("01ed9241-2d5d-1b98-8ecb-40f85676b081");
    if status is string {
    test:assertFail("Successfully deleted API Category");
    } else if status is  error {
        test:assertTrue(true,"Error occured while deleting API Category");
    }
}
