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

import ballerina/http;
import ballerina/log;

# Error handler for the APK service.
public isolated service class ResponseErrorInterceptor {
    *http:ResponseErrorInterceptor;

    isolated remote function interceptResponseError(error err) returns http:Response {
        return getErrorResponse(err);
    }
}

isolated function getErrorResponse(error err) returns http:Response {
    http:Response response = new;
    if err is APKError {
        APKError apkError = <APKError>err;
        ErrorHandler & readonly detail = apkError.detail();
        ErrorDto errorDto = {code: detail.code, message: detail.message, description: detail.description, moreInfo: detail.moreInfo};
        response.statusCode = detail.statusCode;
        response.setJsonPayload(errorDto);
    } else {
        log:printError("Exception Occured", err);
        response.statusCode = 500;
        ErrorDto errorDto = {code: 900900, message: "Internal Server Error", description: "Internal Server Error"};
        response.setJsonPayload(errorDto);
    }
    return response;
}
