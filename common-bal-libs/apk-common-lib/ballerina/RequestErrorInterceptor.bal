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
import ballerina/http;

# Returns the error response.
public isolated service class RequestErrorInterceptor {
    *http:RequestErrorInterceptor;

    isolated resource function 'default [string... path](error err, http:RequestContext ctx) returns http:Response {
        return getErrorResponse(err);
    }
}

public type ErrorDto record {|
    int code;
    string message;
    string description;
    map<string> moreInfo = {};
|};
