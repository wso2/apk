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
import devportal_service.org.wso2.apk.devportal.sdk as sdk;

configurable int DEVPORTAL_PORT = 9443;

listener http:Listener ep0 = new (DEVPORTAL_PORT);

service /api/am/devportal on ep0 {
    isolated resource function get apis(@http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns APIList|BadRequestError|InternalServerErrorError {
        string organization = "carbon.super";
        APIList|APKError apiList = getAPIList('limit, offset, query, organization);
        if apiList is APIList {
            log:printDebug(apiList.toString());
            return apiList;
        } else {
            return handleAPKError(apiList);
        }
    }
    isolated resource function get apis/[string apiId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns API|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|json {
        string organization = "carbon.super";
        API|NotFoundError|APKError api = getAPIByAPIId(apiId, organization);
        if api is API|NotFoundError {
            log:printDebug(api.toString());
            return api;
        } else {
            return handleAPKError(api);
        }
    }
    isolated resource function get apis/[string apiId]/definition(@http:Header string? 'if\-none\-match) returns APIDefinition|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError {
        string organization = "carbon.super";
        APIDefinition|NotFoundError|APKError apiDefinition = getAPIDefinition(apiId, organization);
        if apiDefinition is APIDefinition|NotFoundError {
            log:printDebug(apiDefinition.toString());
            return apiDefinition;
        } else {
            return handleAPKError(apiDefinition);
        }
    }
    resource function get apis/[string apiId]/sdks/[string language](@http:Header string? 'x\-wso2\-tenant) returns http:Response|json|BadRequestError|NotFoundError|InternalServerErrorError {
        string organization = "carbon.super";
        NotFoundError|http:Response|sdk:APIClientGenerationException|APKError sdk = generateSDKImpl(apiId,language, organization);
        if sdk is http:Response {
            return sdk;
        } else if sdk is sdk:APIClientGenerationException {
            InternalServerErrorError internalError = {body: {code: 90931, message: "Internal Error while generating client SDK for given language:" + language + ". Error: " + sdk.message()}};
            return internalError;
        } else if sdk is APKError {
            return handleAPKError(sdk);
        }
    }
    // resource function get apis/[string apiId]/documents(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns DocumentList|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/documents/[string documentId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns Document|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/documents/[string documentId]/content(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns http:Ok|http:SeeOther|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/thumbnail(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns http:Ok|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/ratings(@http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0) returns RatingList|NotAcceptableError {
    // }
    // resource function get apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns Rating|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function put apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Payload Rating payload) returns Rating|BadRequestError|UnsupportedMediaTypeError {
    // }
    // resource function delete apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-match) returns http:Ok {
    // }
    // resource function get apis/[string apiId]/comments(@http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|NotFoundError|InternalServerErrorError {
    // }
    // resource function post apis/[string apiId]/comments(string? replyTo, @http:Payload 'postRequestBody payload) returns CreatedComment|BadRequestError|UnauthorizedError|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError {
    // }
    // resource function get apis/[string apiId]/comments/[string commentId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, boolean includeCommenterInfo = false, int replyLimit = 25, int replyOffset = 0) returns Comment|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError {
    // }
    // resource function delete apis/[string apiId]/comments/[string commentId](@http:Header string? 'if\-match) returns http:Ok|UnauthorizedError|http:Forbidden|NotFoundError|http:MethodNotAllowed|InternalServerErrorError {
    // }
    // resource function patch apis/[string apiId]/comments/[string commentId](@http:Payload 'patchRequestBody payload) returns Comment|BadRequestError|UnauthorizedError|http:Forbidden|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError {
    // }
    // resource function get apis/[string apiId]/comments/[string commentId]/replies(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError {
    // }
    // resource function get apis/[string apiId]/topics(@http:Header string? 'x\-wso2\-tenant) returns TopicList|NotFoundError|InternalServerErrorError {
    // }
    // resource function get apis/[string apiId]/'subscription\-policies(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns ThrottlingPolicy|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    isolated resource function get applications(string? groupId, string? query, string? sortBy, string? sortOrder, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns ApplicationList|http:NotModified|BadRequestError|NotAcceptableError|InternalServerErrorError{
        string organization = "carbon.super";
        ApplicationList|APKError applicationList = getApplicationList(sortBy, groupId, query, sortOrder, 'limit, offset, organization);
        if applicationList is ApplicationList {
            log:printDebug(applicationList.toString());
            return applicationList;
        } else {
            return handleAPKError(applicationList);
        }
    }
    isolated resource function post applications(@http:Payload Application payload) returns CreatedApplication|AcceptedWorkflowResponse|BadRequestError|ConflictError|NotFoundError|InternalServerErrorError|error|json {
        Application|NotFoundError|APKError application = addApplication(payload, "carbon.super", "apkuser");
        if application is Application {
            CreatedApplication createdApp = {body: check application.cloneWithType(Application)};
            log:printDebug(application.toString());
            return createdApp;
        } else if application is NotFoundError {
            return application;
        } else if application is APKError {
            return handleAPKError(application);
        } else {
            return {};
        }
    }
    isolated resource function get applications/[string applicationId](@http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant) returns Application|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError {
        Application|APKError|NotFoundError application = getApplicationById(applicationId, "carbon.super");
        if application is Application|NotFoundError {
            log:printDebug(application.toString());
            return application;
        } else {
            return handleAPKError(application);
        }
    }
    isolated resource function put applications/[string applicationId](@http:Header string? 'if\-match, @http:Payload Application payload) returns Application|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        string organization = "carbon.super";
        Application|NotFoundError|APKError application = updateApplication(applicationId, payload, organization,"apkuser");
        if application is Application|NotFoundError {
            log:printDebug(application.toString());
            return application;
        } else {
            return handleAPKError(application);
        }
    }
    isolated resource function delete applications/[string applicationId](@http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|BadRequestError|PreconditionFailedError|InternalServerErrorError {
        string organization = "carbon.super";
        string|APKError response = deleteApplication(applicationId,organization);
        if response is APKError {
            return handleAPKError(response);
        } else {
            return http:OK;
        }
    }
    // resource function post applications/[string applicationId]/'generate\-keys(@http:Header string? 'x\-wso2\-tenant, @http:Payload ApplicationKeyGenerateRequest payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/'map\-keys(@http:Header string? 'x\-wso2\-tenant, @http:Payload ApplicationKeyMappingRequest payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get applications/[string applicationId]/keys() returns ApplicationKeyList|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get applications/[string applicationId]/keys/[string keyType](string? groupId) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function put applications/[string applicationId]/keys/[string keyType](@http:Payload ApplicationKey payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/keys/[string keyType]/'regenerate\-secret() returns ApplicationKeyReGenerateResponse|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/keys/[string keyType]/'clean\-up(@http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/keys/[string keyType]/'generate\-token(@http:Header string? 'if\-match, @http:Payload ApplicationTokenGenerateRequest payload) returns ApplicationToken|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get applications/[string applicationId]/'oauth\-keys(@http:Header string? 'x\-wso2\-tenant) returns ApplicationKeyList|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get applications/[string applicationId]/'oauth\-keys/[string keyMappingId](string? groupId) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function put applications/[string applicationId]/'oauth\-keys/[string keyMappingId](@http:Payload ApplicationKey payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'regenerate\-secret() returns ApplicationKeyReGenerateResponse|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'clean\-up(@http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'generate\-token(@http:Header string? 'if\-match, @http:Payload ApplicationTokenGenerateRequest payload) returns ApplicationToken|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    isolated resource function post applications/[string applicationId]/'api\-keys/[string keyType]/generate(@http:Header string? 'if\-match, @http:Payload APIKeyGenerateRequest payload) returns APIKey|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        APIKey|APKError|NotFoundError apiKey = generateAPIKey(payload, applicationId, keyType, "apkuser", "carbon.super");
        if apiKey is APIKey|NotFoundError {
            return apiKey;
        } else {
            return handleAPKError(apiKey);
        }
    }
    // resource function post applications/[string applicationId]/'api\-keys/[string keyType]/revoke(@http:Header string? 'if\-match, @http:Payload APIKeyRevokeRequest payload) returns http:Ok|BadRequestError|PreconditionFailedError {
    // }
    // resource function get applications/export(string appName, string appOwner, boolean? withKeys, string? format) returns json|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    // resource function post applications/'import(boolean? preserveOwner, boolean? skipSubscriptions, string? appOwner, boolean? skipApplicationKeys, boolean? update, @http:Payload json payload) returns ApplicationInfo|BadRequestError|NotAcceptableError {
    // }
    isolated resource function get subscriptions(string? apiId, string? applicationId, string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns SubscriptionList|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError {
        SubscriptionList|APKError|NotFoundError subscriptionList = getSubscriptions(apiId, applicationId, groupId, offset, 'limit, "carbon.super");
        if subscriptionList is SubscriptionList|NotFoundError {
            log:printDebug(subscriptionList.toString());
            return subscriptionList;
        } else {
            return handleAPKError(subscriptionList);
        }
    }
    isolated resource function post subscriptions(@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns CreatedSubscription|AcceptedWorkflowResponse|BadRequestError|NotFoundError|InternalServerErrorError|error|json {
        Subscription|APKError|NotFoundError|error subscription = addSubscription(payload, "carbon.super", "apkuser");
        if subscription is APKError {
            return handleAPKError(subscription);
        } else {
            if subscription is Subscription {
            CreatedSubscription createdSub = {body: check subscription.cloneWithType(Subscription)};
            log:printDebug(subscription.toString());
            return createdSub;
            } else if subscription is NotFoundError|error {
                return subscription;
            }
        }
        return {};
    }
    isolated resource function post subscriptions/multiple(@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription[] payload) returns Subscription[]|BadRequestError|UnsupportedMediaTypeError|NotFoundError|InternalServerErrorError|error {
        Subscription[]|APKError|NotFoundError subscriptions = check addMultipleSubscriptions(payload, "carbon.super", "apkuser");
        if subscriptions is Subscription[]|NotFoundError  {
            log:printDebug(subscriptions.toString());
            return subscriptions;
        } else {
            return handleAPKError(subscriptions);
        }
    }
    // resource function get subscriptions/[string apiId]/additionalInfo(string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns AdditionalSubscriptionInfoList|http:NotFound {
    // }
    isolated resource function get subscriptions/[string subscriptionId](@http:Header string? 'if\-none\-match) returns Subscription|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|error {
        Subscription|APKError|NotFoundError subscription = getSubscriptionById(subscriptionId, "carbon.super");
        if subscription is Subscription|NotFoundError {
            log:printDebug(subscription.toString());
            return subscription;
        } else {
            return handleAPKError(subscription);
        }
    }
    isolated resource function put subscriptions/[string subscriptionId](@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns Subscription|AcceptedWorkflowResponse|http:NotModified|BadRequestError|NotFoundError|http:UnsupportedMediaType|InternalServerErrorError|error|json {
        Subscription|NotFoundError|APKError|error subscription = check updateSubscription(subscriptionId, payload, "carbon.super", "apkuser");
        if subscription is Subscription|NotFoundError  {
            log:printDebug(subscription.toString());
            return subscription;
        } else if subscription is APKError {
            return handleAPKError(subscription);
        }
    }
    isolated resource function delete subscriptions/[string subscriptionId](@http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|PreconditionFailedError|BadRequestError|InternalServerErrorError|error{
        string organization = "carbon.super";
        string|APKError response = deleteSubscription(subscriptionId,organization);
        if response is APKError {
            return handleAPKError(response);
        } else {
            return http:OK;
        }
    }
    // resource function get subscriptions/[string subscriptionId]/usage() returns APIMonetizationUsage|http:NotModified|NotFoundError {
    // }
    // resource function get 'throttling\-policies/[string policyLevel](@http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0) returns ThrottlingPolicyList|http:NotModified|NotAcceptableError {
    // }
    // resource function get 'throttling\-policies/[string policyLevel]/[string policyId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns ThrottlingPolicy|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get tags(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns TagList|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function get search(@http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns SearchResultList|http:NotModified|NotAcceptableError {
    // }
    isolated resource function get 'sdk\-gen/languages() returns json|NotFoundError|InternalServerErrorError|BadRequestError|APKError {
        string|json|APKError sdkLanguages = check getSDKLanguages();
        if sdkLanguages is string|json {
            return sdkLanguages;
        } else {
            return handleAPKError(sdkLanguages);
        }
    }
    // resource function get webhooks/subscriptions(string? applicationId, string? apiId, @http:Header string? 'x\-wso2\-tenant) returns WebhookSubscriptionList|NotFoundError|InternalServerErrorError {
    // }
    // resource function get settings(@http:Header string? 'x\-wso2\-tenant) returns Settings|NotFoundError {
    // }
    // resource function get settings/'application\-attributes(@http:Header string? 'if\-none\-match) returns ApplicationAttributeList|NotFoundError|NotAcceptableError {
    // }
    // resource function get tenants(string state = "active", int 'limit = 25, int offset = 0) returns TenantList|NotFoundError|NotAcceptableError {
    // }
    // resource function get recommendations() returns Recommendations|NotFoundError {
    // }
    // resource function get 'api\-categories(@http:Header string? 'x\-wso2\-tenant) returns APICategoryList {
    // }
    // resource function get 'key\-managers(@http:Header string? 'x\-wso2\-tenant) returns KeyManagerList {
    // }
    // resource function get apis/[string apiId]/'graphql\-policies/complexity() returns GraphQLQueryComplexityInfo|http:NotFound {
    // }
    // resource function get apis/[string apiId]/'graphql\-policies/complexity/types() returns GraphQLSchemaTypeList|http:NotFound {
    // }
    // resource function post me/'change\-password(@http:Payload CurrentAndNewPasswords payload) returns http:Ok|BadRequestError {
    // }

}

isolated function handleAPKError(APKError errorDetail) returns InternalServerErrorError|BadRequestError {
    ErrorHandler & readonly detail = errorDetail.detail();
    if detail.statusCode=="400" {
        BadRequestError badRequest = {body: {code: detail.code, message: detail.message}};
        return badRequest;
    }
    InternalServerErrorError internalServerError = {body: {code: detail.code, message: detail.message}};
    return internalServerError;
}