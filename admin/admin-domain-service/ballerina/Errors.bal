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

public isolated function e909400(error e) returns commons:APKError {
    return error commons:APKError(e.message(), e,
        code = 909400,
        message = e.message(),
        statusCode = 500,
        description = e.message()
    );
}

public isolated function e909401(error e) returns commons:APKError {
    return error commons:APKError("Error while retrieving connection", e,
        code = 909401,
        message = "Error while retrieving connection",
        statusCode = 500,
        description = "Error while retrieving connection"
    );
}

public isolated function e909402(error e) returns commons:APKError {
    return error commons:APKError("Error while inserting data into Database", e,
        code = 909402,
        message = "Error while inserting data into Database",
        statusCode = 500,
        description = "Error while inserting data into Database"
    );
}

public isolated function e909403(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving API Categories", e,
        code = 909403,
        message = "Internal Error occured while retrieving API Categories",
        statusCode = 500,
        description = "Internal Error occured while retrieving API Categories"
    );
}

public isolated function e909404(error e) returns commons:APKError {
    return error commons:APKError("Error while checking API Category existence", e,
        code = 909404,
        message = "Error while checking API Category existence",
        statusCode = 500,
        description = "Error while checking API Category existence"
    );
}

public isolated function e909405(error e) returns commons:APKError {
    return error commons:APKError("Error while updating data record in the Database", e,
        code = 909405,
        message = "Error while updating data record in the Database",
        statusCode = 500,
        description = "Error while updating data record in the Database"
    );
}

public isolated function e909406(error e) returns commons:APKError {
    return error commons:APKError("Error while deleting data record in the Database", e,
        code = 909406,
        message = "Error while deleting data record in the Database",
        statusCode = 500,
        description = "Error while deleting data record in the Database"
    );
}

public isolated function e909407() returns commons:APKError {
    return error commons:APKError("Invalid query parameters. Only one of the query parameters can be provided.",
        code = 909407,
        message = "Invalid query parameters. Only one of the query parameters can be provided.",
        statusCode = 406,
        description = "Invalid query parameters. Only one of the query parameters can be provided."
    );
}

public isolated function e909408() returns commons:APKError {
    return error commons:APKError("Error while inserting vhosts data into Database",
        code = 909408,
        message = "Error while inserting vhosts data into Database",
        statusCode = 500,
        description = "Error while inserting vhosts data into Database"
    );
}

public isolated function e909409() returns commons:APKError {
    return error commons:APKError("Error while inserting organization claim data into Database",
        code = 909409,
        message = "Error while inserting organization claim data into Database",
        statusCode = 500,
        description = "Error while inserting organization claim data into Database"
    );
}

public isolated function e909410() returns commons:APKError {
    return error commons:APKError("Error while validating organization name in Database",
        code = 909410,
        message = "Error while validating organization name in Database",
        statusCode = 500,
        description = "Error while validating organization name in Database"
    );
}

public isolated function e909411() returns commons:APKError {
    return error commons:APKError("Error while validating organization id in Database",
        code = 909411,
        message = "Error while validating organization id in Database",
        statusCode = 500,
        description = "Error while validating organization id in Database"
    );
}

public isolated function e909412() returns commons:APKError {
    return error commons:APKError("Error while updating vhosts data into Database",
        code = 909412,
        message = "Error while updating vhosts data into Database",
        statusCode = 500,
        description = "Error while updating vhosts data into Database"
    );
}

public isolated function e909413() returns commons:APKError {
    return error commons:APKError("Error while updating organization data into Database",
        code = 909413,
        message = "Error while updating organization data into Database",
        statusCode = 500,
        description = "Error while updating organization data into Database"
    );
}

public isolated function e909414() returns commons:APKError {
    return error commons:APKError("Organization not found",
        code = 909414,
        message = "Organization not found",
        statusCode = 404,
        description = "Organization not found"
    );
}

public isolated function e909415(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving organization data from Database", e,
        code = 909415,
        message = "Internal Error occured while retrieving organization data from Database",
        statusCode = 500,
        description = "Internal Error occured while retrieving organization data from Database"
    );
}

public isolated function e909416() returns commons:APKError {
    return error commons:APKError("Error while deleting organization data from Database",
        code = 909416,
        message = "Error while deleting organization data from Database",
        statusCode = 500,
        description = "Error while deleting organization data from Database"
    );
}

public isolated function e909417() returns commons:APKError {
    return error commons:APKError("Error while retrieving organization data from Database",
        code = 909417,
        message = "Error while retrieving organization data from Database",
        statusCode = 500,
        description = "Error while retrieving organization data from Database"
    );
}

public isolated function e909418() returns commons:APKError {
    return error commons:APKError("Error while retrieving Application Usage Plan",
        code = 909418,
        message = "Error while retrieving Application Usage Plan",
        statusCode = 500,
        description = "Error while retrieving Application Usage Plan"
    );
}

public isolated function e909419(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving Application Usage Plans", e,
        code = 909419,
        message = "Internal Error occured while retrieving Application Usage Plans",
        statusCode = 500,
        description = "Internal Error occured while retrieving Application Usage Plans"
    );
}

public isolated function e909420(error e) returns commons:APKError {
    return error commons:APKError("Error while retrieving Business Plan", e,
        code = 909420,
        message = "Error while retrieving Business Plan",
        statusCode = 500,
        description = "Error while retrieving Business Plan"
    );
}

public isolated function e909421(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving Business Plans", e,
        code = 909421,
        message = "Internal Error occured while retrieving Business Plans",
        statusCode = 500,
        description = "Internal Error occured while retrieving Business Plans"
    );
}

public isolated function e909422(error e) returns commons:APKError {
    return error commons:APKError("Error while retrieving Deny Policy from DB", e,
        code = 909422,
        message = "Error while retrieving Deny Policy from DB",
        statusCode = 500,
        description = "Error while retrieving Deny Policy from DB"
    );
}

public isolated function e909423(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving Deny Policies", e,
        code = 909423,
        message = "Internal Error occured while retrieving Deny Policies",
        statusCode = 500,
        description = "Internal Error occured while retrieving Deny Policies"
    );
}

public isolated function e909424(string payloadName) returns commons:APKError {
    return error commons:APKError("API Category already exists by name:" + payloadName,
        code = 909424,
        message = "API Category already exists by name:" + payloadName,
        statusCode = 400,
        description = "API Category already exists by name:" + payloadName
    );
}

public isolated function e909425(string name) returns commons:APKError {
    return error commons:APKError("API Category already exists by name:" + name,
        code = 909425,
        message = "API Category already exists by name:" + name,
        statusCode = 400,
        description = "API Category already exists by name:" + name
    );
}

public isolated function e909426() returns commons:APKError {
    return error commons:APKError("API Category not found",
        code = 909426,
        message = "API Category not found",
        statusCode = 404,
        description = "API Category not found"
    );
}

public isolated function e909427(string name) returns commons:APKError {
    return error commons:APKError("Organization already exists by name:" + name,
        code = 909427,
        message = "Organization already exists by name:" + name,
        statusCode = 409,
        description = "Organization already exists by name:" + name
    );
}

public isolated function e909428(string id) returns commons:APKError {
    return error commons:APKError("Organization ID not exist by:" + id,
        code = 909428,
        message = "Organization ID not exist by:" + id,
        statusCode = 400,
        description = "Organization ID not exist by:" + id
    );
}

public isolated function e909429() returns commons:APKError {
    return error commons:APKError("Application Usage Plan not found",
        code = 909429,
        message = "Application Usage Plan not found",
        statusCode = 404,
        description = "Application Usage Plan not found"
    );
}

public isolated function e909430() returns commons:APKError {
    return error commons:APKError("Business Plan not found",
        code = 909430,
        message = "Business Plan not found",
        statusCode = 404,
        description = "Business Plan not found"
    );
}

public isolated function e909431() returns commons:APKError {
    return error commons:APKError("Deny Policy not found",
        code = 909431,
        message = "Deny Policy not found",
        statusCode = 404,
        description = "Deny Policy not found"
    );
}

public isolated function e909432(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving KeyManagers", e,
        code = 909432,
        message = "Internal Error occured while retrieving KeyManagers",
        statusCode = 500,
        description = "Internal Error occured while retrieving KeyManagers"
    );
}

public isolated function e909433(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving KeyManager", e,
        code = 909433,
        message = "Internal Error occured while retrieving KeyManager",
        statusCode = 500,
        description = "Internal Error occured while retrieving KeyManager"
    );
}

public isolated function e909434() returns commons:APKError {
    return error commons:APKError("KeyManager Name can't be empty",
        code = 909434,
        message = "KeyManager Name can't be empty",
        statusCode = 400,
        description = "KeyManager Name can't be empty"
    );
}

isolated function e909435(string name, string organization) returns commons:APKError {
    return error commons:APKError("KeyManager Name " + name + " already exist in organization " + organization,
        code = 909434,
        message = "KeyManager Name " + name + " already exist in organization " + organization,
        statusCode = 400,
        description = "KeyManager Name " + name + " already exist in organization " + organization
    );
}

isolated function e909436(string 'type) returns commons:APKError {
    return error commons:APKError("KeyManager Type " + 'type + " not exist",
        code = 909436,
        message = "KeyManager Type " + 'type + " not exist",
        statusCode = 400,
        description = "KeyManager Type " + 'type + " not exist"
    );
}

isolated function e909437() returns commons:APKError {
    return error commons:APKError("Required Keymanager Endpoint or Configuration not available",
        code = 909437,
        message = "Required Keymanager Endpoint or Configuration not available",
        statusCode = 400,
        description = "Required Keymanager Endpoint or Configuration not available"
    );
}

isolated function e909438(error? e) returns commons:APKError {
    return error commons:APKError("Internal Server Error",e,
        code = 909438,
        message = "Internal Server Error",
        statusCode = 500,
        description = "Internal Server Error"
    );
}

isolated function e909439(string id, string organization) returns commons:APKError {
    return error commons:APKError("KeyManager from " + id + " not exist in organization " + organization,
        code = 909439,
        message = "KeyManager from " + id + " not exist in organization " + organization,
        statusCode = 404,
        description = "KeyManager from " + id + " not exist in organization " + organization
    );
}
isolated function e909440(string id, string organization,error? e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while deleting keymanager " + id + " from organization " + organization,e,
        code = 909440,
        message = "Internal Error occured while deleting keymanager " + id + " from organization " + organization,
        statusCode = 500,
        description = "Internal Error occured while deleting keymanager " + id + " from organization " + organization
    );
}
