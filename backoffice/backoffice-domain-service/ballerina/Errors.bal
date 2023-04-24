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
// Before adding another function for a new error code
// make sure there is no already existing error code for that.
// If there is an error code for that, reuse it.

import wso2/apk_common_lib as commons;

public isolated function e909601(error e) returns commons:APKError {
    return error commons:APKError( "Error while retrieving connection", e,
        code = 909601,
        message = "Error while retrieving connection",
        statusCode = 500,
        description = "Error while retrieving connection"
    ); 
}

public isolated function e909602() returns commons:APKError {
    return error commons:APKError( "API Definition Not Found for provided API ID",
        code = 909602,
        message = "API Definition Not Found for provided API ID",
        statusCode = 404,
        description = "API Definition Not Found for provided API ID"
    ); 
}

public isolated function e909603() returns commons:APKError {
    return error commons:APKError( "API not found in the database",
        code = 909603,
        message = "API not found in the database",
        statusCode = 404,
        description = "API not found in the database"
    ); 
}

public isolated function e909604() returns commons:APKError {
    return error commons:APKError( "Error while retrieving API",
        code = 909604,
        message = "Error while retrieving API",
        statusCode = 500,
        description = "Error while retrieving API"
    ); 
}

public isolated function e909605() returns commons:APKError {
    return error commons:APKError( "Internal Error while deleting API By Id",
        code = 909605,
        message = "Internal Error while deleting API By Id",
        statusCode = 500,
        description = "Internal Error while deleting API By Id"
    ); 
}

public isolated function e909606(string apiId) returns commons:APKError {
    return error commons:APKError( "API with " + apiId + " not found",
        code = 909606,
        message = "API with " + apiId + " not found",
        statusCode = 404,
        description = "API with " + apiId + " not found"
    ); 
}

public isolated function e909607(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error occured while retrieving APIs", e,
        code = 909607,
        message = "Internal Error occured while retrieving APIs",
        statusCode = 500,
        description = "Internal Error occured while retrieving APIs"
    ); 
}

public isolated function e909608(error e) returns commons:APKError {
    return error commons:APKError( "Error while updating LC state into Database", e,
        code = 909608,
        message = "Error while updating LC state into Database",
        statusCode = 500,
        description = "Error while updating LC state into Database"
    ); 
}

public isolated function e909609() returns commons:APKError {
    return error commons:APKError( "Invalid Lifecycle targetState",
        code = 909609,
        message = "Invalid Lifecycle targetState",
        statusCode = 400,
        description = "Invalid Lifecycle targetState"
    ); 
}

public isolated function e909610(error e) returns commons:APKError {
    return error commons:APKError( "Error while geting LC state from Database", e,
        code = 909610,
        message = "Error while geting LC state from Database",
        statusCode = 400,
        description = "Error while geting LC state from Database"
    ); 
}

public isolated function e909611(error e) returns commons:APKError {
    return error commons:APKError( "Error while inserting data into Database", e,
        code = 909611,
        message = "Error while inserting data into Database",
        statusCode = 500,
        description = "Error while inserting data into Database"
    ); 
}

public isolated function e909612(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error occured while retrieving LC event History", e,
        code = 909612,
        message = "Internal Error occured while retrieving LC event History",
        statusCode = 500,
        description = "Internal Error occured while retrieving LC event History"
    ); 
}

public isolated function e909613(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error while geting subscription infomation", e,
        code = 909613,
        message = "Internal Error while geting subscription infomation",
        statusCode = 500,
        description = "Internal Error while geting subscription infomation"
    ); 
}

public isolated function e909614(error e, string apiId) returns commons:APKError {
    return error commons:APKError( "Internal Error while geting API for provided apiId " + apiId, e,
        code = 909614,
        message = "Internal Error while geting API for provided apiId " + apiId,
        statusCode = 500,
        description = "Internal Error while geting API for provided apiId " + apiId
    ); 
}

public isolated function e909615(error e) returns commons:APKError {
    return error commons:APKError( "Error while changing status of the subscription in the Database", e,
        code = 909615,
        message = "Error while changing status of the subscription in the Database",
        statusCode = 500,
        description = "Error while changing status of the subscription in the Database"
    ); 
}

public isolated function e909616(error e) returns commons:APKError {
    return error commons:APKError( "Error while retriving API", e,
        code = 909616,
        message = "Error while retriving API",
        statusCode = 500,
        description = "Error while retriving API"
    ); 
}

public isolated function e909617(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error while retrieving API Definition", e,
        code = 909617,
        message = "Internal Error while retrieving API Definition",
        statusCode = 500,
        description = "Internal Error while retrieving API Definition"
    ); 
}

public isolated function e909618(error e) returns commons:APKError {
    return error commons:APKError( "Error while updating API data into Database", e,
        code = 909618,
        message = "Error while updating API data into Database",
        statusCode = 500,
        description = "Error while updating API data into Database"
    ); 
}

public isolated function e909619(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error occured while retrieving API Categories", e,
        code = 909619,
        message = "Internal Error occured while retrieving API Categories",
        statusCode = 500,
        description = "Internal Error occured while retrieving API Categories"
    ); 
}

public isolated function e909620(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error occured while retrieving Business Plans", e,
        code = 909620,
        message = "Internal Error occured while retrieving Business Plans",
        statusCode = 500,
        description = "Internal Error occured while retrieving Business Plans"
    ); 
}

public isolated function e909621() returns commons:APKError {
    return error commons:APKError( "Invalid Content Search Text Provided. Missing :",
        code = 909621,
        message = "Invalid Content Search Text Provided. Missing :",
        statusCode = 400,
        description = "Invalid Content Search Text Provided. Missing :"
    ); 
}

public isolated function e909622() returns commons:APKError {
    return error commons:APKError( "Invalid Content Search Text Provided. Missing content keyword",
        code = 909622,
        message = "Invalid Content Search Text Provided. Missing content keyword",
        statusCode = 400,
        description = "Invalid Content Search Text Provided. Missing content keyword"
    ); 
}

public isolated function e909623() returns commons:APKError {
    return error commons:APKError( "Invalid blockState provided",
        code = 909623,
        message = "Invalid blockState provided",
        statusCode = 400,
        description = "Invalid blockState provided"
    ); 
}