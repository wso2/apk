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

public isolated function e90901(string id) returns commons:APKError {
    return error commons:APKError( id + " not found",
        code = 90902,
        message = id + " not found",
        statusCode = 404,
        description = id + " not found"
    ); 
}

public isolated function e90902() returns commons:APKError {
    return error commons:APKError( "Context/Name doesn't exist",
        code = 90902,
        message = "Context/Name doesn't exist",
        statusCode = 404,
        description = "Context/Name doesn't exist"
    ); 
}

public isolated function e90903() returns commons:APKError {
    return error commons:APKError( "apiId not found in request",
        code = 90903,
        message = "apiId not found in request",
        statusCode = 404,
        description = "apiId not found in request"
    ); 
}

public isolated function e90904(string serviceKey) returns commons:APKError {
    return error commons:APKError( "Service from " + serviceKey + " not found.",
        code = 90904,
        message =  "Service from " + serviceKey + " not found.",
        statusCode = 404,
        description =  "Service from " + serviceKey + " not found."
    ); 
}

public isolated function e90905() returns commons:APKError {
    return error commons:APKError( "type field unavailable",
        code = 90905,
        message = "type field unavailable",
        statusCode = 404,
        description = "type field unavailable"
    ); 
}

public isolated function e90906() returns commons:APKError {
    return error commons:APKError( "Unsupported API type",
        code = 90906,
        message = "Unsupported API type",
        statusCode = 400,
        description = "Unsupported API type"
    ); 
}

public isolated function e90907() returns commons:APKError {
    return error commons:APKError( "Multiple fields of url, file, inlineAPIDefinition given",
        code = 90907,
        message = "Multiple fields of url, file, inlineAPIDefinition given",
        statusCode = 400,
        description = "Multiple fields of url, file, inlineAPIDefinition given"
    ); 
}

public isolated function e90908() returns commons:APKError {
    return error commons:APKError( "Atleast one of the field required",
        code = 90908,
        message = "Atleast one of the field required",
        statusCode = 400,
        description = "Atleast one of the field required"
    ); 
}

public isolated function e90909() returns commons:APKError {
    return error commons:APKError( "Additional properties not provided",
        code = 90909,
        message = "Additional properties not provided",
        statusCode = 400,
        description = "Additional properties not provided"
    ); 
}

public isolated function e90910() returns commons:APKError {
    return error commons:APKError( "Invalid operation policy name",
        code = 90910,
        message = "Invalid operation policy name",
        statusCode = 400,
        description = "Invalid operation policy name"
    ); 
}

public isolated function e90911(string apiName) returns commons:APKError {
    return error commons:APKError( "API Name - " + apiName + " already exist",
        code = 90911,
        message = "API Name - " + apiName + " already exist",
        statusCode = 409,
        description = "API Name - " + apiName + " already exist"
    ); 
}

public isolated function e90912(string apiContext) returns commons:APKError {
    return error commons:APKError( "API Context - " + apiContext + " already exist",
        code = 90912,
        message = "API Context - " + apiContext + " already exist",
        statusCode = 409,
        description = "API Context - " + apiContext + " already exist"
    ); 
}

public isolated function e90913() returns commons:APKError {
    return error commons:APKError( "Sandbox endpoint not specified",
        code = 90913,
        message = "Sandbox endpoint not specified",
        statusCode = 400,
        description = "Sandbox endpoint not specified"
    ); 
}

public isolated function e90914() returns commons:APKError {
    return error commons:APKError( "Production endpoint not specified",
        code = 90914,
        message = "Production endpoint not specified",
        statusCode = 400,
        description = "Production endpoint not specified"
    ); 
}

public isolated function e90915(string apiContext) returns commons:APKError {
    return error commons:APKError("API context " + apiContext + " invalid",
        code = 90915,
        message = "API context " + apiContext + " invalid",
        statusCode = 400,
        description = "API context " + apiContext + " invalid"
    ); 
}

public isolated function e90916(string apiName) returns commons:APKError {
    return error commons:APKError("API name " + apiName + " invalid",
        code = 90916,
        message = "API name " + apiName + " invalid",
        statusCode = 400,
        description = "API name " + apiName + " invalid"
    ); 
}

public isolated function e90917() returns commons:APKError {
    return error commons:APKError("Invalid API request",
        code = 90917,
        message = "Invalid API request",
        statusCode = 400,
        description = "Invalid API request"
    ); 
}

public isolated function e90918() returns commons:APKError {
    return error commons:APKError("Error while generating token",
        code = 90918,
        message = "Error while generating token",
        statusCode = 400,
        description = "Error while generating token"
    ); 
}

public isolated function e90919(string keyWord) returns commons:APKError {
    return error commons:APKError("Invalid KeyWord " + keyWord,
        code = 90919,
        message = "Invalid KeyWord " + keyWord,
        statusCode = 400,
        description = "Invalid KeyWord " + keyWord
    ); 
}

public isolated function e90920() returns commons:APKError {
    return error commons:APKError("Invalid Sort By/Sort Order value",
        code = 90920,
        message = "Invalid Sort By/Sort Order value",
        statusCode = 400,
        description = "Invalid Sort By/Sort Order value"
    ); 
}

public isolated function e90921() returns commons:APKError {
    return error commons:APKError("Atleast one operation need to specified",
        code = 90921,
        message = "Atleast one operation need to specified",
        statusCode = 400,
        description = "Atleast one operation need to specified"
    ); 
}

public isolated function e90922() returns commons:APKError {
    return error commons:APKError("Internal server error",
        code = 90922,
        message = "Internal server error",
        statusCode = 500,
        description = "Internal server error"
    ); 
}

public isolated function e90923() returns commons:APKError {
    return error commons:APKError("Internal server error",
        code = 90923,
        message = "Internal server error",
        statusCode = 500,
        description = "Internal server error"
    ); 
}

public isolated function e90924(string policyName) returns commons:APKError {
    return error commons:APKError( "Invalid parameters provided for policy " + policyName,
        code = 90924,
        message = "Invalid parameters provided for policy " + policyName,
        statusCode = 400,
        description = "Invalid parameters provided for policy " + policyName
    ); 
}

public isolated function e90925() returns commons:APKError {
    return error commons:APKError( "Presence of both resource level and API level operation policies is not allowed",
        code = 90925,
        message = "Presence of both resource level and API level operation policies is not allowed",
        statusCode = 400,
        description = "Presence of both resource level and API level operation policies is not allowed"
    ); 
}

public isolated function e90926() returns commons:APKError {
    return error commons:APKError( "Presence of both resource level and API level rate limits is not allowed",
        code = 90926,
        message = "Presence of both resource level and API level rate limits is not allowed",
        statusCode = 400,
        description = "Presence of both resource level and API level rate limits is not allowed"
    ); 
}