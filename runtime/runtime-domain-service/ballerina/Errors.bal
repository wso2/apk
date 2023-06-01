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

public isolated function e909002() returns commons:APKError {
    return error commons:APKError("Context/Name doesn't exist",
        code = 909002,
        message = "Context/Name doesn't exist",
        statusCode = 404,
        description = "Context/Name doesn't exist"
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

public isolated function e909004() returns commons:APKError {
    return error commons:APKError("Invalid property id in Request",
        code = 909004,
        message = "Invalid property id in Request",
        statusCode = 406,
        description = "Invalid property id in Request"
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
        message = "Atleast one of the field required",
        statusCode = 406,
        description = "Atleast one of the field required"
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

public isolated function e909012(string apiContext) returns commons:APKError {
    return error commons:APKError("API Context - " + apiContext + " already exist",
        code = 909012,
        message = "API Context - " + apiContext + " already exist",
        statusCode = 409,
        description = "API Context - " + apiContext + " already exist"
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

public isolated function e909023(error e) returns commons:APKError {
    return error commons:APKError("Internal error occured while retrieving definition", e,
        code = 909023,
        message = "Internal error occured while retrieving definition",
        statusCode = 500,
        description = "Internal error occured while retrieving definition"
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

public isolated function e909027(error e) returns commons:APKError {
    return error commons:APKError("Error while retrieving API", e,
        code = 909027,
        message = "Error while retrieving API",
        statusCode = 500,
        description = "Error while retrieving API"
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

public isolated function e909029(error e) returns commons:APKError {
    return error commons:APKError("Error while retrieving Mediation policy", e,
        code = 909029,
        message = "Error while retrieving Mediation policy",
        statusCode = 500,
        description = "Error while retrieving Mediation policy"
    );
}

public isolated function e909030() returns commons:APKError {
    return error commons:APKError("Certificate is expired",
        code = 909030,
        message = "Certificate is expired",
        statusCode = 400,
        description = "Certificate is expired"
    );
}

public isolated function e909031(error e) returns commons:APKError {
    return error commons:APKError("Error while adding certificate", e,
        code = 909031,
        message = "Error while adding certificate",
        statusCode = 500,
        description = "Error while adding certificate"
    );
}

public isolated function e909032() returns commons:APKError {
    return error commons:APKError("Host/Certificte is empty in payload",
        code = 909032,
        message = "Host/Certificte is empty in payload",
        statusCode = 500,
        description = "Host/Certificte is empty in payload"
    );
}

public isolated function e909033(error e) returns commons:APKError {
    return error commons:APKError("Error while retrieving endpoint certificate request", e,
        code = 909033,
        message = "Error while retrieving endpoint certificate request",
        statusCode = 500,
        description = "Error while retrieving endpoint certificate request"
    );
}

public isolated function e909034(string certificateId) returns commons:APKError {
    return error commons:APKError("Certificate " + certificateId + " not found",
        code = 909034,
        message = "Certificate " + certificateId + " not found",
        statusCode = 404,
        description = "Certificate " + certificateId + " not found"
    );
}

public isolated function e909035(error e) returns commons:APKError {
    return error commons:APKError("Error while deleting endpoint certificate", e,
        code = 909035,
        message = "Error while deleting endpoint certificate",
        statusCode = 500,
        description = "Error while deleting endpoint certificate"
    );
}

public isolated function e909036(error e) returns commons:APKError {
    return error commons:APKError("Error while getting endpoint certificate content", e,
        code = 909036,
        message = "Error while getting endpoint certificate content",
        statusCode = 500,
        description = "Error while getting endpoint certificate content"
    );
}

public isolated function e909037(error e) returns commons:APKError {
    return error commons:APKError("Error while getting endpoint certificate by id", e,
        code = 909037,
        message = "Error while getting endpoint certificate by id",
        statusCode = 500,
        description = "Error while getting endpoint certificate by id"
    );
}

public isolated function e909038(error e) returns commons:APKError {
    return error commons:APKError("Error while updating endpoint certificate", e,
        code = 909038,
        message = "Error while updating endpoint certificate",
        statusCode = 500,
        description = "Error while updating endpoint certificate"
    );
}

public isolated function e909039() returns commons:APKError {
    return error commons:APKError("Invalid value for offset",
        code = 909039,
        message = "Invalid value for offset",
        statusCode = 406,
        description = "Invalid value for offset"
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

public isolated function e909041() returns commons:APKError {
    return error commons:APKError("Accept header should be application/json or application/yaml",
        code = 909041,
        message = "Accept header should be application/json or application/yaml",
        statusCode = 406,
        description = "Accept header should be application/json or application/yaml"
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

public isolated function e909043() returns commons:APKError {
    return error commons:APKError("Error occured while generating openapi definition",
        code = 909043,
        message = "Error occured while generating openapi definition",
        statusCode = 500,
        description = "Error occured while generating openapi definition"
    );
}

public isolated function e909044() returns commons:APKError {
    return error commons:APKError("Retrieve definition from Url result",
        code = 909044,
        message = "Retrieve definition from Url result",
        statusCode = 406,
        description = "Retrieve definition from Url result"
    );
}

public isolated function e909045() returns commons:APKError {
    return error commons:APKError("New Version/APIID does not exist",
        code = 909045,
        message = "New Version/APIID does not exist",
        statusCode = 404,
        description = "New Version/APIID does not exist"
    );
}

public isolated function e909046(string newVersion) returns commons:APKError {
    return error commons:APKError("New version - " + newVersion + " already exist",
        code = 909046,
        message = "New version - " + newVersion + " already exist",
        statusCode = 409,
        description = "New version - " + newVersion + " already exist"
    );
}

public isolated function e909047(string serviceId) returns commons:APKError {
    return error commons:APKError(serviceId + " service does not exist",
        code = 909047,
        message = serviceId + " service does not exist",
        statusCode = 404,
        description = serviceId + " service does not exist"
    );
}

public isolated function e909048() returns commons:APKError {
    return error commons:APKError("API Type change not supported from update",
        code = 909048,
        message = "API Type change not supported from update",
        statusCode = 406,
        description = "API Type change not supported from update"
    );
}

public isolated function e909049() returns commons:APKError {
    return error commons:APKError("Context change not supported from update",
        code = 909049,
        message = "Context change not supported from update",
        statusCode = 406,
        description = "Context change not supported from update"
    );
}

public isolated function e909050() returns commons:APKError {
    return error commons:APKError("Version change not supported from update",
        code = 909050,
        message = "Version change not supported from update",
        statusCode = 406,
        description = "Version change not supported from update"
    );
}
