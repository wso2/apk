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

configurable int DEVPORTAL_PORT = 9444;

listener http:Listener ep0 = new (DEVPORTAL_PORT);

@http:ServiceConfig {
    cors: {
        allowOrigins: ["*"],
        allowCredentials: true,
        allowHeaders: ["*"],
        exposeHeaders: ["*"],
        maxAge: 84900
    }
}

@display {
    label: "devportal-api-service",
    id: "devportal-api-service"
}

service /api/am/devportal/v2 on ep0 {
    resource function get apis(@http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns APIList|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns API|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/swagger(string? environmentName, @http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant) returns string|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/'graphql\-schema(@http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant) returns string|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/sdks/[string language](@http:Header string? 'x\-wso2\-tenant) returns json|BadRequestError|NotFoundError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/wsdl(string? environmentName, @http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant) returns http:Ok|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/documents(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns DocumentList|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/documents/[string documentId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns Document|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/documents/[string documentId]/content(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns http:Ok|http:SeeOther|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/thumbnail(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns http:Ok|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/ratings(@http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0) returns RatingList|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns Rating|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function put apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Payload Rating payload) returns Rating|BadRequestError|UnsupportedMediaTypeError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function delete apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-match) returns http:Ok|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/comments(@http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|NotFoundError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    // resource function post apis/[string apiId]/comments(string? replyTo, @http:Payload PostRequestBody payload) returns CreatedComment|BadRequestError|UnauthorizedError|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError|http:NotImplemented {
    //     http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
    //     return notImplementedError;
    // }
    resource function get apis/[string apiId]/comments/[string commentId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, boolean includeCommenterInfo = false, int replyLimit = 25, int replyOffset = 0) returns Comment|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function delete apis/[string apiId]/comments/[string commentId](@http:Header string? 'if\-match) returns http:Ok|UnauthorizedError|http:Forbidden|NotFoundError|http:MethodNotAllowed|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    // resource function patch apis/[string apiId]/comments/[string commentId](@http:Payload PatchRequestBody payload) returns Comment|BadRequestError|UnauthorizedError|http:Forbidden|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError|http:NotImplemented {
    //     http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
    //     return notImplementedError;
    // }
    resource function get apis/[string apiId]/comments/[string commentId]/replies(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|UnauthorizedError|NotFoundError|NotAcceptableError|http:NotImplemented|InternalServerErrorError {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/topics(@http:Header string? 'x\-wso2\-tenant) returns TopicList|NotFoundError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/'subscription\-policies(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns ThrottlingPolicy|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get applications(string groupId, string query, string sortBy, string sortOrder, int 'limit = 25, int offset = 0, string organization = "carbon.super") returns ApplicationList|http:NotModified|BadRequestError|NotAcceptableError|http:NotImplemented|error {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function post applications(@http:Payload Application payload) returns CreatedApplication|AcceptedWorkflowResponse|BadRequestError|ConflictError|UnsupportedMediaTypeError|error|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get applications/[string applicationId]() returns Application|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented|error {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function put applications/[string applicationId](@http:Payload Application payload) returns Application|BadRequestError|NotFoundError|PreconditionFailedError|error|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function delete applications/[string applicationId]() returns http:Ok|AcceptedWorkflowResponse|NotFoundError|PreconditionFailedError|string|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/'generate\-keys(@http:Header string? 'x\-wso2\-tenant, @http:Payload ApplicationKeyGenerateRequest payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function post applications/[string applicationId]/'map\-keys(@http:Header string? 'x\-wso2\-tenant, @http:Payload ApplicationKeyMappingRequest payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get applications/[string applicationId]/keys() returns ApplicationKeyList|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get applications/[string applicationId]/keys/[string keyType](string? groupId) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function put applications/[string applicationId]/keys/[string keyType](@http:Payload ApplicationKey payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/keys/[string keyType]/'regenerate\-secret() returns ApplicationKeyReGenerateResponse|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/keys/[string keyType]/'clean\-up(@http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function post applications/[string applicationId]/keys/[string keyType]/'generate\-token(@http:Header string? 'if\-match, @http:Payload ApplicationTokenGenerateRequest payload) returns ApplicationToken|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get applications/[string applicationId]/'oauth\-keys(@http:Header string? 'x\-wso2\-tenant) returns ApplicationKeyList|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get applications/[string applicationId]/'oauth\-keys/[string keyMappingId](string? groupId) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function put applications/[string applicationId]/'oauth\-keys/[string keyMappingId](@http:Payload ApplicationKey payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'regenerate\-secret() returns ApplicationKeyReGenerateResponse|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'clean\-up(@http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'generate\-token(@http:Header string? 'if\-match, @http:Payload ApplicationTokenGenerateRequest payload) returns ApplicationToken|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/'api\-keys/[string keyType]/generate(@http:Header string? 'if\-match, @http:Payload APIKeyGenerateRequest payload) returns APIKey|BadRequestError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/[string applicationId]/'api\-keys/[string keyType]/revoke(@http:Header string? 'if\-match, @http:Payload APIKeyRevokeRequest payload) returns http:Ok|BadRequestError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get applications/export(string appName, string appOwner, boolean? withKeys, string? format) returns json|BadRequestError|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post applications/'import(boolean? preserveOwner, boolean? skipSubscriptions, string? appOwner, boolean? skipApplicationKeys, boolean? update, @http:Payload json payload) returns ApplicationInfo|BadRequestError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get subscriptions(string? apiId, string? applicationId, string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns SubscriptionList|http:NotModified|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function post subscriptions(@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns CreatedSubscription|AcceptedWorkflowResponse|BadRequestError|UnsupportedMediaTypeError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function post subscriptions/multiple(@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription[] payload) returns Subscription[]|BadRequestError|UnsupportedMediaTypeError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get subscriptions/[string apiId]/additionalInfo(string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns AdditionalSubscriptionInfoList|http:NotFound|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get subscriptions/[string subscriptionId](@http:Header string? 'if\-none\-match) returns Subscription|http:NotModified|NotFoundError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function put subscriptions/[string subscriptionId](@http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns Subscription|AcceptedWorkflowResponse|http:NotModified|BadRequestError|http:NotFound|http:UnsupportedMediaType|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function delete subscriptions/[string subscriptionId](@http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get subscriptions/[string subscriptionId]/usage() returns APIMonetizationUsage|http:NotModified|NotFoundError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get 'throttling\-policies/[string policyLevel](@http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0) returns ThrottlingPolicyList|http:NotModified|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get 'throttling\-policies/[string policyLevel]/[string policyId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns ThrottlingPolicy|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get tags(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns TagList|http:NotModified|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get search(@http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns SearchResultList|http:NotModified|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get 'sdk\-gen/languages() returns json|NotFoundError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get webhooks/subscriptions(string? applicationId, string? apiId, @http:Header string? 'x\-wso2\-tenant) returns WebhookSubscriptionList|NotFoundError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get settings(@http:Header string? 'x\-wso2\-tenant) returns Settings|NotFoundError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get settings/'application\-attributes(@http:Header string? 'if\-none\-match) returns ApplicationAttributeList|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get tenants(string state = "active", int 'limit = 25, int offset = 0) returns TenantList|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get recommendations() returns Recommendations|NotFoundError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get 'api\-categories(@http:Header string? 'x\-wso2\-tenant) returns APICategoryList|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get 'key\-managers(@http:Header string? 'x\-wso2\-tenant) returns KeyManagerList|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/'graphql\-policies/complexity() returns GraphQLQueryComplexityInfo|http:NotFound|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get apis/[string apiId]/'graphql\-policies/complexity/types() returns GraphQLSchemaTypeList|http:NotFound|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function post me/'change\-password(@http:Payload CurrentAndNewPasswords payload) returns http:Ok|BadRequestError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
}
