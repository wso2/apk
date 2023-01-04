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

configurable int DEVPORTAL_PORT = 9443;

listener http:Listener ep0 = new (DEVPORTAL_PORT);

service /api/am/devportal on ep0 {
    isolated resource function get apis(@http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns APIList|NotAcceptableError|InternalServerErrorError|error {
        string organization = "carbon.super";
        string?| APIList | error apiList = check getAPIList('limit, offset, query, organization);
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
    isolated resource function get apis/[string apiId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns API|http:NotModified|NotFoundError|NotAcceptableError|InternalServerErrorError|error|json {
        string organization = "carbon.super";
        string?|API|error api = check getAPIByAPIId(apiId, organization);
        if api is string {
            json j = check value:fromJsonString(api);
            API clonedAPI = check j.cloneWithType(API);
            log:printDebug(clonedAPI.toString());
            return clonedAPI;
        } else if api is API {
            log:printDebug(api.toString());
            return api;
        } else {
            InternalServerErrorError internalError = {body: {code: 90913, message: "Internal Error while retrieving API By Id"}};
            return internalError;
        }
    }
    isolated resource function get apis/[string apiId]/definition(@http:Header string? 'if\-none\-match) returns APIDefinition|http:NotModified|NotFoundError|NotAcceptableError|InternalServerErrorError|error {
        string organization = "carbon.super";
        APIDefinition|NotFoundError|error apiDefinition = check getAPIDefinition(apiId, organization);
        if apiDefinition is APIDefinition|NotFoundError {
            log:printDebug(apiDefinition.toString());
            return apiDefinition;
        } else {
            InternalServerErrorError internalError = {body: {code: 90914, message: "Internal Error while retrieving API Definition By API Id"}};
            return internalError;
        }
    }
    // resource function get apis/[string apiId]/sdks/[string language](@http:Header string? 'x\-wso2\-tenant) returns json|BadRequestError|NotFoundError|InternalServerErrorError {
    // }
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
    isolated resource function get applications(string? groupId, string? query, string? sortBy, string? sortOrder, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns ApplicationList|http:NotModified|BadRequestError|NotAcceptableError|InternalServerErrorError|error{
        string organization = "carbon.super";
        string?|ApplicationList|error applicationList = check getApplicationList(sortBy, groupId, query, sortOrder, 'limit, offset, organization);
        if applicationList is string {
            json j = check value:fromJsonString(applicationList);
            ApplicationList appList = check j.cloneWithType(ApplicationList);
            log:printDebug(appList.toString());
            return appList;
        } else if applicationList is ApplicationList {
            log:printDebug(applicationList.toString());
            return applicationList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving all Applications"}};
            return internalError;
        }
    }
    isolated resource function post applications(@http:Payload Application payload) returns CreatedApplication|AcceptedWorkflowResponse|BadRequestError|ConflictError|UnsupportedMediaTypeError|InternalServerErrorError|error {
        string?|Application|error application = check addApplication(payload, "carbon.super", "apkuser");
        if application is string {
            json j = check value:fromJsonString(application);
            CreatedApplication createdApp = {body: check j.cloneWithType(Application)};
            return createdApp;
        } else if application is Application {
            CreatedApplication createdApp = {body: check application.cloneWithType(Application)};
            log:printDebug(application.toString());
            return createdApp;
        } else {
            InternalServerErrorError internalError = {body: {code: 90910, message: "Internal Error while adding Application"}};
            return internalError;
        }
    }
    isolated resource function get applications/[string applicationId](@http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant) returns Application|http:NotModified|NotFoundError|NotAcceptableError|InternalServerErrorError|error {
        string?|Application|error application = check getApplicationById(applicationId, "carbon.super");
        if application is string {
            json j = check value:fromJsonString(application);
            Application app = check j.cloneWithType(Application);
            log:printDebug(app.toString());
            return app;
        } else if application is Application {
            log:printDebug(application.toString());
            return application;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving Application By Id"}};
            return internalError;
        }
    }
    isolated resource function put applications/[string applicationId](@http:Header string? 'if\-match, @http:Payload Application payload) returns Application|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError|error {
        string organization = "carbon.super";
        string?|Application|NotFoundError|error application = check updateApplication(applicationId, payload, organization,"apkuser");
        if application is string {
            json j = check value:fromJsonString(application);
            Application app = check j.cloneWithType(Application);
            log:printDebug(app.toString());
            return app;
        } else if application is Application|NotFoundError {
            log:printDebug(application.toString());
            return application;
        } else {
            InternalServerErrorError internalError = {body: {code: 90911, message: "Internal Error while updating Application"}};
            return internalError;
        }
    }
    isolated resource function delete applications/[string applicationId](@http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|PreconditionFailedError|InternalServerErrorError|error {
        string organization = "carbon.super";
        string|error? response = check deleteApplication(applicationId,organization);
        if response is error {
            InternalServerErrorError internalError = {body: {code: 90912, message: "Internal Error while deleting Application By Id"}};
            return internalError;
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
    isolated resource function post applications/[string applicationId]/'api\-keys/[string keyType]/generate(@http:Header string? 'if\-match, @http:Payload APIKeyGenerateRequest payload) returns APIKey|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError|error {
        APIKey|error apiKey = check generateAPIKey(payload, applicationId, keyType, "apkuser", "carbon.super");
        if apiKey is APIKey {
            return apiKey;
        } else {
            InternalServerErrorError internalError = {body: {code: 909123, message: "Internal Error while generating API Key"}};
            return internalError;
        }
    }
    // resource function post applications/[string applicationId]/'api\-keys/[string keyType]/revoke(@http:Header string? 'if\-match, @http:Payload APIKeyRevokeRequest payload) returns http:Ok|BadRequestError|PreconditionFailedError {
    // }
    // resource function get applications/export(string appName, string appOwner, boolean? withKeys, string? format) returns json|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    // resource function post applications/'import(boolean? preserveOwner, boolean? skipSubscriptions, string? appOwner, boolean? skipApplicationKeys, boolean? update, @http:Payload json payload) returns ApplicationInfo|BadRequestError|NotAcceptableError {
    // }
    isolated resource function get subscriptions(string? apiId, string? applicationId, string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns SubscriptionList|http:NotModified|NotAcceptableError|InternalServerErrorError|error {
        string?|SubscriptionList|error subscriptionList = check getSubscriptions(apiId, applicationId, groupId, offset, 'limit, "carbon.super");
        if subscriptionList is string {
            json j = check value:fromJsonString(subscriptionList);
            SubscriptionList sub = check j.cloneWithType(SubscriptionList);
            log:printDebug(sub.toString());
            return sub;
        } else if subscriptionList is SubscriptionList {
            log:printDebug(subscriptionList.toString());
            return subscriptionList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90922, message: "Internal Error while retrieving All Subscriptions of an API or Application or both"}};
            return internalError;
        }
    }
    isolated resource function post subscriptions(@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns CreatedSubscription|AcceptedWorkflowResponse|BadRequestError|UnsupportedMediaTypeError|InternalServerErrorError|error {
        string?|Subscription|error subscription = check addSubscription(payload, "carbon.super", "apkuser");
        if subscription is string {
            json j = check value:fromJsonString(subscription);
            CreatedSubscription createdSub = {body: check j.cloneWithType(Subscription)};
            return createdSub;
        } else if subscription is Subscription  {
            CreatedSubscription createdSub = {body: check subscription.cloneWithType(Subscription)};
            log:printDebug(subscription.toString());
            return createdSub;
        } else {
            InternalServerErrorError internalError = {body: {code: 90921, message: "Internal Error while adding Subscription"}};
            return internalError;
        }
    }
    isolated resource function post subscriptions/multiple(@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription[] payload) returns Subscription[]|BadRequestError|UnsupportedMediaTypeError|InternalServerErrorError|error {
        Subscription[]|error? subscriptions = check addMultipleSubscriptions(payload, "carbon.super", "apkuser");
        if subscriptions is Subscription[]  {
            log:printDebug(subscriptions.toString());
            return subscriptions;
        } else {
            InternalServerErrorError internalError = {body: {code: 90921, message: "Internal Error while adding Subscriptions"}};
            return internalError;
        }
    }
    // resource function get subscriptions/[string apiId]/additionalInfo(string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns AdditionalSubscriptionInfoList|http:NotFound {
    // }
    isolated resource function get subscriptions/[string subscriptionId](@http:Header string? 'if\-none\-match) returns Subscription|http:NotModified|NotFoundError|InternalServerErrorError|error {
        string?|Subscription|error subscription = check getSubscriptionById(subscriptionId, "carbon.super");
        if subscription  is string {
            json j = check value:fromJsonString(subscription );
            Subscription sub = check j.cloneWithType(Subscription);
            log:printDebug(sub.toString());
            return sub;
        } else if subscription is Subscription {
            log:printDebug(subscription.toString());
            return subscription;
        } else {
            InternalServerErrorError internalError = {body: {code: 90922, message: "Internal Error while retrieving Subscription By Id"}};
            return internalError;
        }
    }
    isolated resource function put subscriptions/[string subscriptionId](@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns Subscription|AcceptedWorkflowResponse|http:NotModified|BadRequestError|NotFoundError|http:UnsupportedMediaType|InternalServerErrorError|error {
        string?|Subscription|NotFoundError|error subscription = check updateSubscription(subscriptionId, payload, "carbon.super", "apkuser");
        if subscription is string {
            json j = check value:fromJsonString(subscription);
            Subscription updatedSub = check j.cloneWithType(Subscription);
            return updatedSub;
        } else if subscription is Subscription  {
            log:printDebug(subscription.toString());
            return subscription;
        } else {
            InternalServerErrorError internalError = {body: {code: 90921, message: "Internal Error while updating Subscription Tier"}};
            return internalError;
        }
    }
    isolated resource function delete subscriptions/[string subscriptionId](@http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|PreconditionFailedError|InternalServerErrorError|error{
        string organization = "carbon.super";
        string|error? response = check deleteSubscription(subscriptionId,organization);
        if response is error {
            InternalServerErrorError internalError = {body: {code: 90923, message: "Internal Error while deleting Subscription By Id"}};
            return internalError;
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
    // resource function get 'sdk\-gen/languages() returns json|NotFoundError|InternalServerErrorError {
    // }
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
