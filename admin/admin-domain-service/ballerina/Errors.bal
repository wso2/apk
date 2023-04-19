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

public isolated function e909401(error e) returns commons:APKError {
    return error commons:APKError( "Error while retrieving connection", e,
        code = 909401,
        message = "Error while retrieving connection",
        statusCode = 500,
        description = "Error while retrieving connection"
    ); 
}

public isolated function e909402(error e) returns commons:APKError {
    return error commons:APKError( "Error while inserting data into Database", e,
        code = 909402,
        message = "Error while inserting data into Database",
        statusCode = 500,
        description = "Error while inserting data into Database"
    ); 
}

public isolated function e909403(error e) returns commons:APKError {
    return error commons:APKError( "Internal Error occured while retrieving API Categories", e,
        code = 909403,
        message = "Internal Error occured while retrieving API Categories",
        statusCode = 500,
        description = "Internal Error occured while retrieving API Categories"
    ); 
}

public isolated function e909404(error e) returns commons:APKError {
    return error commons:APKError( "Error while checking API Category existence", e,
        code = 909404,
        message = "Error while checking API Category existence",
        statusCode = 500,
        description = "Error while checking API Category existence"
    ); 
}

public isolated function e909405(error e) returns commons:APKError {
    return error commons:APKError( "Error while updating data record in the Database", e,
        code = 909405,
        message = "Error while updating data record in the Database",
        statusCode = 500,
        description = "Error while updating data record in the Database"
    ); 
}

public isolated function e909406(error e) returns commons:APKError {
    return error commons:APKError( "Error while deleting data record in the Database", e,
        code = 909406,
        message = "Error while deleting data record in the Database",
        statusCode = 500,
        description = "Error while deleting data record in the Database"
    ); 
}
