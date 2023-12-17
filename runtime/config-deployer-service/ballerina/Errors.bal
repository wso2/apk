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

public isolated function e909000(int code, string msg) returns commons:APKError {
    return error commons:APKError(msg,
        code = 909000,
        message = msg,
        statusCode = code,
        description = msg
    );
}

public isolated function e909001(string id) returns commons:APKError {
    return error commons:APKError(id + " not found",
        code = 909001,
        message = id + " not found",
        statusCode = 404,
        description = id + " not found"
    );
}

public isolated function e909003() returns commons:APKError {
    return error commons:APKError("apiId not found in request",
        code = 909003,
        message = "apiId not found in request",
        statusCode = 404,
        description = "apiId not found in request"
    );
}

public isolated function e909005(string 'field) returns commons:APKError {
    string msg = 'field + " field(s) unavailable";
    return error commons:APKError(msg,
        code = 909005,
        message = msg,
        statusCode = 404,
        description = msg
    );
}

public isolated function e909006() returns commons:APKError {
    return error commons:APKError("Unsupported API type",
        code = 909006,
        message = "Unsupported API type",
        statusCode = 406,
        description = "Unsupported API type"
    );
}

public isolated function e909007() returns commons:APKError {
    return error commons:APKError("Multiple fields of url, file, inlineAPIDefinition given",
        code = 909007,
        message = "Multiple fields of url, file, inlineAPIDefinition given",
        statusCode = 406,
        description = "Multiple fields of url, file, inlineAPIDefinition given"
    );
}

public isolated function e909008() returns commons:APKError {
    return error commons:APKError("Atleast one of the field required",
        code = 909008,
        message = "Atleast one of the fields required",
        statusCode = 406,
        description = "Atleast one of the fields required"
    );
}

public isolated function e909009() returns commons:APKError {
    return error commons:APKError("Additional properties not provided",
        code = 909009,
        message = "Additional properties not provided",
        statusCode = 406,
        description = "Additional properties not provided"
    );
}

public isolated function e909010() returns commons:APKError {
    return error commons:APKError("Invalid operation policy name",
        code = 909010,
        message = "Invalid operation policy name",
        statusCode = 406,
        description = "Invalid operation policy name"
    );
}

public isolated function e909011(string apiName) returns commons:APKError {
    return error commons:APKError("API Name - " + apiName + " already exist",
        code = 909011,
        message = "API Name - " + apiName + " already exist",
        statusCode = 409,
        description = "API Name - " + apiName + " already exist"
    );
}

public isolated function e909012(string basePath) returns commons:APKError {
    return error commons:APKError("API Base Path - " + basePath + " already exist",
        code = 909012,
        message = "API Base Path - " + basePath + " already exist",
        statusCode = 409,
        description = "Base Path - " + basePath + " already exist"
    );
}

public isolated function e909013() returns commons:APKError {
    return error commons:APKError("Sandbox endpoint not specified",
        code = 909013,
        message = "Sandbox endpoint not specified",
        statusCode = 406,
        description = "Sandbox endpoint not specified"
    );
}

public isolated function e909014() returns commons:APKError {
    return error commons:APKError("Production endpoint not specified",
        code = 909014,
        message = "Production endpoint not specified",
        statusCode = 406,
        description = "Production endpoint not specified"
    );
}

public isolated function e909015(string basePath) returns commons:APKError {
    return error commons:APKError("API Base Path " + basePath + " invalid",
        code = 909015,
        message = "API Base Path " + basePath + " invalid",
        statusCode = 406,
        description = "API Base Path " + basePath + " invalid"
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

public isolated function e909018(string msg) returns commons:APKError {
    return error commons:APKError("Invalid API request",
        code = 909018,
        message = "Invalid API request",
        statusCode = 406,
        description = "Invalid API request" + msg
    );
}

public isolated function e909019() returns commons:APKError {
    return error commons:APKError("Invalid authtypes provided",
        code = 909018,
        message = "Invalid authtypes provided",
        statusCode = 406,
        description = "Invalid authtypes provided"
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

public isolated function e909022(string msg, error? e) returns commons:APKError {
    if e is error {
        return error commons:APKError(msg, e,
        code = 909022,
        message = "Internal server error",
        statusCode = 500,
        description = "Internal server error"
    );
    } else {
        return error commons:APKError(msg,
        code = 909022,
        message = "Internal server error",
        statusCode = 500,
        description = "Internal server error"
    );
    }
}

public isolated function e909029(map<string> errorItems) returns commons:APKError {
    return error commons:APKError("Apk-conf validation failed",
        code = 909029,
        message = "Invalid apk-conf provided",
        statusCode = 400,
        description = "Invalid apk-conf provided",
        moreInfo = errorItems
    );
}

public isolated function e909024(string policyName) returns commons:APKError {
    return error commons:APKError("Invalid parameters provided for policy " + policyName,
        code = 909024,
        message = "Invalid parameters provided for policy " + policyName,
        statusCode = 406,
        description = "Invalid parameters provided for policy " + policyName
    );
}

public isolated function e909025() returns commons:APKError {
    return error commons:APKError("Presence of both resource level and API level operation policies is not allowed",
        code = 909025,
        message = "Presence of both resource level and API level operation policies is not allowed",
        statusCode = 406,
        description = "Presence of both resource level and API level operation policies is not allowed"
    );
}

public isolated function e909026() returns commons:APKError {
    return error commons:APKError("Presence of both resource level and API level rate limits is not allowed",
        code = 909026,
        message = "Presence of both resource level and API level rate limits is not allowed",
        statusCode = 406,
        description = "Presence of both resource level and API level rate limits is not allowed"
    );
}

public isolated function e909040() returns commons:APKError {
    return error commons:APKError("Internal Error Occured while converting json to yaml",
        code = 909040,
        message = "Internal Error Occured while converting json to yaml",
        statusCode = 500,
        description = "Internal Error Occured while converting json to yaml"
    );
}

public isolated function e909042() returns commons:APKError {
    return error commons:APKError("Unsupported API type",
        code = 909042,
        message = "Unsupported API type",
        statusCode = 400,
        description = "Unsupported API type"
    );
}

public isolated function e909052(error e) returns commons:APKError {
    return error commons:APKError("Error while generating k8s artifact", e,
        code = 909052,
        message = "Error while generating k8s artifact",
        statusCode = 500,
        description = "Error while generating k8s artifact:" + e.message()
    );
}

public isolated function e909043() returns commons:APKError {
    return error commons:APKError("Error occured while generating openapi definition",
        code = 909043,
        message = "Error occured while generating openapi definition",
        statusCode = 500,
        description = "Error occured while generating openapi definition"
    );
}

public isolated function e909044() returns commons:APKError {
    return error commons:APKError("Invalid url provided for openapi definition",
        code = 909044,
        message = "Invalid url provided for openapi definition",
        statusCode = 406,
        description = "Invalid url provided for openapi definition"
    );
}

public isolated function e9090445() returns commons:APKError {
    return error commons:APKError("Atleast one Vhost need per production or/and sandbox",
        code = 9090445,
        message = "Atleast one Vhost need per production or/and sandbox",
        statusCode = 406,
        description = "Atleast one Vhost need per production or/and sandbox"
    );
}

public isolated function e909028() returns commons:APKError {
    return error commons:APKError("Internal error occured while deploying API",
        code = 909028,
        message = "Internal error occured while deploying API",
        statusCode = 500,
        description = "Internal error occured while deploying API"
    );
}
