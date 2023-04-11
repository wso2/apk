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

public isolated function e909001(string id) returns commons:APKError {
    return error commons:APKError( id + " not found",
        code = 909001,
        message = id + " not found",
        statusCode = 404,
        description = id + " not found"
    ); 
}

public isolated function e909002() returns commons:APKError {
    return error commons:APKError( "Context/Name doesn't exist",
        code = 909002,
        message = "Context/Name doesn't exist",
        statusCode = 404,
        description = "Context/Name doesn't exist"
    ); 
}

public isolated function e909003() returns commons:APKError {
    return error commons:APKError( "apiId not found in request",
        code = 909003,
        message = "apiId not found in request",
        statusCode = 404,
        description = "apiId not found in request"
    ); 
}

public isolated function e909004(string serviceKey) returns commons:APKError {
    return error commons:APKError( "Service from " + serviceKey + " not found",
        code = 909004,
        message =  "Service from " + serviceKey + " not found",
        statusCode = 404,
        description =  "Service from " + serviceKey + " not found"
    ); 
}

public isolated function e909005() returns commons:APKError {
    return error commons:APKError( "type field unavailable",
        code = 909005,
        message = "type field unavailable",
        statusCode = 404,
        description = "type field unavailable"
    ); 
}

public isolated function e909006() returns commons:APKError {
    return error commons:APKError( "Unsupported API type",
        code = 909006,
        message = "Unsupported API type",
        statusCode = 406,
        description = "Unsupported API type"
    ); 
}

public isolated function e909007() returns commons:APKError {
    return error commons:APKError( "Multiple fields of url, file, inlineAPIDefinition given",
        code = 909007,
        message = "Multiple fields of url, file, inlineAPIDefinition given",
        statusCode = 406,
        description = "Multiple fields of url, file, inlineAPIDefinition given"
    ); 
}

public isolated function e909008() returns commons:APKError {
    return error commons:APKError( "Atleast one of the field required",
        code = 909008,
        message = "Atleast one of the field required",
        statusCode = 406,
        description = "Atleast one of the field required"
    ); 
}

public isolated function e909009() returns commons:APKError {
    return error commons:APKError( "Additional properties not provided",
        code = 909009,
        message = "Additional properties not provided",
        statusCode = 406,
        description = "Additional properties not provided"
    ); 
}

public isolated function e909010() returns commons:APKError {
    return error commons:APKError( "Invalid operation policy name",
        code = 909010,
        message = "Invalid operation policy name",
        statusCode = 406,
        description = "Invalid operation policy name"
    ); 
}

public isolated function e909011(string apiName) returns commons:APKError {
    return error commons:APKError( "API Name - " + apiName + " already exist",
        code = 909011,
        message = "API Name - " + apiName + " already exist",
        statusCode = 409,
        description = "API Name - " + apiName + " already exist"
    ); 
}

public isolated function e909012(string apiContext) returns commons:APKError {
    return error commons:APKError( "API Context - " + apiContext + " already exist",
        code = 909012,
        message = "API Context - " + apiContext + " already exist",
        statusCode = 409,
        description = "API Context - " + apiContext + " already exist"
    ); 
}

public isolated function e909013() returns commons:APKError {
    return error commons:APKError( "Sandbox endpoint not specified",
        code = 909013,
        message = "Sandbox endpoint not specified",
        statusCode = 406,
        description = "Sandbox endpoint not specified"
    ); 
}

public isolated function e909014() returns commons:APKError {
    return error commons:APKError( "Production endpoint not specified",
        code = 909014,
        message = "Production endpoint not specified",
        statusCode = 406,
        description = "Production endpoint not specified"
    ); 
}

public isolated function e909015(string apiContext) returns commons:APKError {
    return error commons:APKError("API context " + apiContext + " invalid",
        code = 909015,
        message = "API context " + apiContext + " invalid",
        statusCode = 406,
        description = "API context " + apiContext + " invalid"
    ); 
}

public isolated function e909016(string apiName) returns commons:APKError {
    return error commons:APKError("API name " + apiName + " invalid",
        code = 909016,
        message = "API name " + apiName + " invalid",
        statusCode = 406,
        description = "API name " + apiName + " invalid"
    ); 
}

public isolated function e909017() returns commons:APKError {
    return error commons:APKError("Invalid API request",
        code = 909017,
        message = "Invalid API request",
        statusCode = 406,
        description = "Invalid API request"
    ); 
}

public isolated function e909018() returns commons:APKError {
    return error commons:APKError("Error while generating token",
        code = 909018,
        message = "Error while generating token",
        statusCode = 500,
        description = "Error while generating token"
    ); 
}

public isolated function e909019(string keyWord) returns commons:APKError {
    return error commons:APKError("Invalid keyword " + keyWord,
        code = 909019,
        message = "Invalid keyword " + keyWord,
        statusCode = 406,
        description = "Invalid keyword " + keyWord
    ); 
}

public isolated function e909020() returns commons:APKError {
    return error commons:APKError("Invalid Sort By/Sort Order value",
        code = 909020,
        message = "Invalid Sort By/Sort Order value",
        statusCode = 406,
        description = "Invalid Sort By/Sort Order value"
    ); 
}

public isolated function e909021() returns commons:APKError {
    return error commons:APKError("Atleast one operation need to specified",
        code = 909021,
        message = "Atleast one operation need to specified",
        statusCode = 406,
        description = "Atleast one operation need to specified"
    ); 
}

public isolated function e909022(error e) returns commons:APKError {
    return error commons:APKError("Internal server error", e,
        code = 909022,
        message = "Internal server error",
        statusCode = 500,
        description = "Internal server error"
    ); 
}

// e90923 is here

public isolated function e909024(string policyName) returns commons:APKError {
    return error commons:APKError( "Invalid parameters provided for policy " + policyName,
        code = 909024,
        message = "Invalid parameters provided for policy " + policyName,
        statusCode = 406,
        description = "Invalid parameters provided for policy " + policyName
    ); 
}

public isolated function e909025() returns commons:APKError {
    return error commons:APKError( "Presence of both resource level and API level operation policies is not allowed",
        code = 909025,
        message = "Presence of both resource level and API level operation policies is not allowed",
        statusCode = 406,
        description = "Presence of both resource level and API level operation policies is not allowed"
    ); 
}

public isolated function e909026() returns commons:APKError {
    return error commons:APKError( "Presence of both resource level and API level rate limits is not allowed",
        code = 909026,
        message = "Presence of both resource level and API level rate limits is not allowed",
        statusCode = 406,
        description = "Presence of both resource level and API level rate limits is not allowed"
    ); 
}