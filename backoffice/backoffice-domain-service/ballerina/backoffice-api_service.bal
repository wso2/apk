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

configurable int BACKOFFICE_PORT = 9443;

listener http:Listener ep0 = new (BACKOFFICE_PORT, secureSocket = {
    'key: {
        certFile: <string>keyStores.tls.certFilePath,
        keyFile: <string>keyStores.tls.keyFilePath
    }
}, interceptors = [jwtValidationInterceptor, requestErrorInterceptor]);

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

    isolated resource function get apis(string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns APIList|http:NotModified|NotAcceptableError|InternalServerErrorError|BadRequestError {
        APIList|APKError apiList = getAPIList('limit, offset, query, "carbon.super");
        if apiList is APIList {
            return apiList;
        } else {
            return handleAPKError(apiList);
        }
    }

    isolated resource function get apis/[string apiId](@http:Header string? 'if\-none\-match) returns API|http:NotModified|NotFoundError|NotAcceptableError|BadRequestError|InternalServerErrorError {
        API|APKError response = getAPI(apiId);
        if response is API {
            return response;
        } else {
            return handleAPKError(response);
        }
    }
    resource function put apis/[string apiId](@http:Header string? 'if\-none\-match, @http:Payload ModifiableAPI payload) returns API|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|BadRequestError|InternalServerErrorError {
        API|APKError updatedAPI = updateAPI(apiId, payload, "carbon.super");
        if updatedAPI is API {
            return updatedAPI;
        }
        return handleAPKError(updatedAPI);
    }

    isolated resource function get apis/[string apiId]/definition(@http:Header string? 'if\-none\-match) returns APIDefinition|http:NotModified|NotFoundError|NotAcceptableError|BadRequestError|InternalServerErrorError {
        APIDefinition|NotFoundError|APKError apiDefinition = getAPIDefinition(apiId);
        if apiDefinition is APIDefinition|NotFoundError {
            log:printDebug(apiDefinition.toString());
            return apiDefinition;
        } else {
            return handleAPKError(apiDefinition);
        }
    }
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
    isolated resource function get subscriptions(string? apiId, @http:Header string? 'if\-none\-match, string? query, int 'limit = 25, int offset = 0) returns SubscriptionList|http:NotModified|NotAcceptableError|BadRequestError|InternalServerErrorError {
        SubscriptionList|APKError subList = getSubscriptions(apiId);
        if subList is SubscriptionList {
            return subList;
        } else {
            return handleAPKError(subList);
        }
    }
    // resource function get subscriptions/[string subscriptionId]/'subscriber\-info() returns SubscriberInfo|NotFoundError {
    // }
    isolated resource function post subscriptions/'block\-subscription(string subscriptionId, string blockState, @http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        string|APKError response = blockSubscription(subscriptionId, blockState);
        if response is APKError {
            return handleAPKError(response);
        } else {
            return http:OK;
        }
    }
    isolated resource function post subscriptions/'unblock\-subscription(string subscriptionId, @http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        string|error response = unblockSubscription(subscriptionId);
        if response is APKError {
            return handleAPKError(response);
        } else {
            return http:OK;
        }
    }
    // resource function get 'usage\-plans(@http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns UsagePlanList|http:NotModified|NotAcceptableError {
    // }
    // resource function get search(string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns SearchResultList|http:NotModified|NotAcceptableError {
    // }
    // resource function get settings() returns Settings|NotFoundError {
    // }

    isolated resource function get 'api\-categories() returns APICategoryList|BadRequestError|InternalServerErrorError {
        APICategoryList|APKError apiCategoryList = getAllCategoryList();
        if apiCategoryList is APICategoryList {
            return apiCategoryList;
        } else {
            return handleAPKError(apiCategoryList);
        }
    }

    isolated resource function post apis/'change\-lifecycle(string targetState, string apiId, @http:Header string? 'if\-match) returns LifecycleState|BadRequestError|UnauthorizedError|NotFoundError|ConflictError|InternalServerErrorError|BadRequestError|error {
        LifecycleState|error changeState = changeLifeCyleState(targetState, apiId, "carbon.super");
        if changeState is LifecycleState {
            return changeState;
        } else {
            return error("Error while updating LC state of API" + changeState.message());
        }
    }
    isolated resource function get apis/[string apiId]/'lifecycle\-history(@http:Header string? 'if\-none\-match) returns LifecycleHistory|UnauthorizedError|NotFoundError|InternalServerErrorError|BadRequestError {
        LifecycleHistory|APKError lcList = getLcEventHistory(apiId);
        if lcList is LifecycleHistory {
            return lcList;
        } else {
            return handleAPKError(lcList);
        }
    }
    isolated resource function get apis/[string apiId]/'lifecycle\-state(@http:Header string? 'if\-none\-match) returns LifecycleState|UnauthorizedError|NotFoundError|InternalServerErrorError|BadRequestError|error {
        LifecycleState|error currentState = getLifeCyleState(apiId, "carbon.super");
        if currentState is LifecycleState {
            return currentState;
        } else {
            return error("Error while getting LC state of API" + currentState.message());
        }
    }
    resource function get 'business\-plans(@http:Header string? accept = "application/json") returns BusinessPlanList|InternalServerErrorError|BadRequestError {
        BusinessPlanList|APKError subPolicyList = getBusinessPlans();
        if subPolicyList is BusinessPlanList {
            log:printDebug(subPolicyList.toString());
            return subPolicyList;
        } else {
            return handleAPKError(subPolicyList);
        }
    }
}
