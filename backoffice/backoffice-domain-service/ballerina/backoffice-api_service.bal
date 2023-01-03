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
import ballerina/lang.value;

configurable int BACKOFFICE_PORT = 9443;

listener http:Listener ep0 = new (BACKOFFICE_PORT);

@http:ServiceConfig {
    cors: {
        allowOrigins: ["*"],
        allowCredentials: true,
        allowHeaders: ["*"],
        exposeHeaders: ["*"],
        maxAge: 84900
    }
}

service /api/am/backoffice on ep0 {
    
    isolated resource function get apis(string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns APIList|http:NotModified|NotAcceptableError|InternalServerErrorError|error {
        string?| APIList | error apiList = check getAPIList();
        if apiList is string {
            json j = check value:fromJsonString(apiList);
            APIList apiListObj = check j.cloneWithType(APIList);
            return apiListObj;
        } else if apiList is APIList {
            log:printDebug(apiList.toString());
            return apiList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90914, message: "Internal Error while retrieving all APIs"}};
            return internalError;
        }
    }
    // resource function get apis/[string apiId](@http:Header string? 'if\-none\-match) returns API|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function put apis/[string apiId](@http:Header string? 'if\-none\-match, @http:Payload ModifiableAPI payload) returns API|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get apis/[string apiId]/definition(@http:Header string? 'if\-none\-match) returns APIDefinition|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/'resource\-paths(@http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns ResourcePathList|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/thumbnail(@http:Header string? 'if\-none\-match, @http:Header string? accept = "application/json") returns http:Ok|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function put apis/[string apiId]/thumbnail(@http:Header string? 'if\-match, @http:Payload json payload) returns FileInfo|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get apis/[string apiId]/documents(@http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, @http:Header string? accept = "application/json") returns DocumentList|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function post apis/[string apiId]/documents(@http:Payload Document payload) returns CreatedDocument|BadRequestError|UnsupportedMediaTypeError {
    // }
    // resource function get apis/[string apiId]/documents/[string documentId](@http:Header string? 'if\-none\-match, @http:Header string? accept = "application/json") returns Document|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function put apis/[string apiId]/documents/[string documentId](@http:Header string? 'if\-match, @http:Payload Document payload) returns Document|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function delete apis/[string apiId]/documents/[string documentId](@http:Header string? 'if\-match) returns http:Ok|NotFoundError|PreconditionFailedError {
    // }
    // resource function get apis/[string apiId]/documents/[string documentId]/content(@http:Header string? 'if\-none\-match, @http:Header string? accept = "application/json") returns http:Ok|http:SeeOther|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function post apis/[string apiId]/documents/[string documentId]/content(@http:Header string? 'if\-match, @http:Payload json payload) returns Document|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get apis/[string apiId]/comments(int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|NotFoundError|InternalServerErrorError {
    // }
    // resource function post apis/[string apiId]/comments(string? replyTo, @http:Payload 'postRequestBody payload) returns CreatedComment|BadRequestError|UnauthorizedError|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError {
    // }
    // resource function get apis/[string apiId]/comments/[string commentId](@http:Header string? 'if\-none\-match, boolean includeCommenterInfo = false, int replyLimit = 25, int replyOffset = 0) returns Comment|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError {
    // }
    // resource function delete apis/[string apiId]/comments/[string commentId](@http:Header string? 'if\-match) returns http:Ok|UnauthorizedError|ForbiddenError|NotFoundError|http:MethodNotAllowed|InternalServerErrorError {
    // }
    // resource function patch apis/[string apiId]/comments/[string commentId](@http:Payload 'patchRequestBody payload) returns Comment|BadRequestError|UnauthorizedError|ForbiddenError|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError {
    // }
    // resource function get apis/[string apiId]/comments/[string commentId]/replies(@http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError {
    // }
    // resource function post apis/[string apiId]/monetize(@http:Payload APIMonetizationInfo payload) returns http:Created|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/monetization() returns http:Ok|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/revenue() returns APIRevenue|http:NotModified|NotFoundError {
    // }
    // resource function get apis/[string apiId]/'external\-stores(@http:Header string? 'if\-none\-match) returns APIExternalStoreList|NotFoundError|InternalServerErrorError {
    // }
    // resource function post apis/[string apiId]/'publish\-to\-external\-stores(string? externalStoreIds, @http:Header string? 'if\-match) returns APIExternalStoreList|NotFoundError|InternalServerErrorError {
    // }
    isolated resource function get subscriptions(string? apiId, @http:Header string? 'if\-none\-match, string? query, int 'limit = 25, int offset = 0) returns SubscriptionList|http:NotModified|NotAcceptableError|InternalServerErrorError|error {
        SubscriptionList|error subList = getSubscriptions(apiId);
        if subList is SubscriptionList {
            return subList;
        } else {
           return  error(subList.message());
        }
    }
    // resource function get subscriptions/[string subscriptionId]/usage() returns APIMonetizationUsage|http:NotModified|NotFoundError {
    // }
    // resource function get subscriptions/[string subscriptionId]/'subscriber\-info() returns SubscriberInfo|NotFoundError {
    // }
    isolated resource function post subscriptions/'block\-subscription(string subscriptionId, string blockState, @http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError| InternalServerErrorError {
        string|error response = blockSubscription(subscriptionId, blockState);
        if response is error {
            InternalServerErrorError internalError = {body: {code: 90912, message: response.message()}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    isolated resource function post subscriptions/'unblock\-subscription(string subscriptionId, @http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError| InternalServerErrorError {
        string|error response = unblockSubscription(subscriptionId);
        if response is error {
            InternalServerErrorError internalError = {body: {code: 90912, message: response.message()}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    // resource function get 'usage\-plans(@http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns UsagePlanList|http:NotModified|NotAcceptableError {
    // }
    // resource function get search(string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns SearchResultList|http:NotModified|NotAcceptableError {
    // }
    // resource function get 'external\-stores() returns ExternalStore|InternalServerErrorError {
    // }
    // resource function get settings() returns Settings|NotFoundError {
    // }
    // resource function get 'api\-categories() returns APICategoryList {
    // }
    isolated resource function post apis/'change\-lifecycle(string targetState, string apiId, @http:Header string? 'if\-match) returns LifecycleState|BadRequestError|UnauthorizedError|NotFoundError|ConflictError|InternalServerErrorError|error {
        LifecycleState | error  changeState = changeLifeCyleState(targetState, apiId, "carbon.super");
        if changeState is LifecycleState {
            return changeState;
        } else {
            return error("Error while updating LC state of API" + changeState.message());
        }
    }
    isolated resource function get apis/[string apiId]/'lifecycle\-history(@http:Header string? 'if\-none\-match) returns LifecycleHistory|UnauthorizedError|NotFoundError|InternalServerErrorError {
        LifecycleHistory|error? lcList = getLcEventHistory(apiId);
        if lcList is LifecycleHistory {
            return lcList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving LC events"}};
            return internalError;
        }
    }
    isolated resource function get apis/[string apiId]/'lifecycle\-state(@http:Header string? 'if\-none\-match) returns LifecycleState|UnauthorizedError|NotFoundError|InternalServerErrorError|error {
        LifecycleState | error currentState = getLifeCyleState(apiId, "carbon.super");
            if currentState is LifecycleState {
                return currentState;
        } else {
            return error("Error while getting LC state of API" + currentState.message());
        }
    }
}
