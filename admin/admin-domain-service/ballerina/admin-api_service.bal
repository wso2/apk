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
import ballerina/io;
import ballerina/lang.value;

service /api/am/admin on ep0 {
    // resource function get throttling/policies/search(string? query) returns ThrottlePolicyDetailsList {
    // }
    resource function get throttling/policies/application(@http:Header string? accept = "application/json") returns ApplicationThrottlePolicyList|NotAcceptableError|error {
        string?|ApplicationThrottlePolicyList appPolicyList = getApplicationUsagePlans();
        if appPolicyList is string {
            json j = check value:fromJsonString(appPolicyList);
            ApplicationThrottlePolicyList polList = check j.cloneWithType(ApplicationThrottlePolicyList);
            return polList;
        }
        if appPolicyList is ApplicationThrottlePolicyList {
            return appPolicyList;
        }
        io:print(appPolicyList);
        return {};
    }
    resource function post throttling/policies/application(@http:Payload ApplicationThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns CreatedApplicationThrottlePolicy|BadRequestError|UnsupportedMediaTypeError|error {
        string?|ApplicationThrottlePolicy createdAppPol = addApplicationUsagePlan(payload);
        io:println(createdAppPol);
        if createdAppPol is string {
            json j = check value:fromJsonString(createdAppPol);
            CreatedApplicationThrottlePolicy crPol = {body: check j.cloneWithType(ApplicationThrottlePolicy)};
            return crPol;
        }
        if createdAppPol is ApplicationThrottlePolicy {
            CreatedApplicationThrottlePolicy crPol = {body: check createdAppPol.cloneWithType(ApplicationThrottlePolicy)};
            return crPol;
        }
        return error("Error while adding Application Policy");
    }
    resource function get throttling/policies/application/[string policyId]() returns ApplicationThrottlePolicy|NotFoundError|NotAcceptableError|error {
        string?|ApplicationThrottlePolicy appPolicy = getApplicationUsagePlanById(policyId);
        if appPolicy is string {
            json j = check value:fromJsonString(appPolicy);
            ApplicationThrottlePolicy policy = check j.cloneWithType(ApplicationThrottlePolicy);
            return policy;
        }
        if appPolicy is ApplicationThrottlePolicy {
            return appPolicy;
        }
        io:println(appPolicy);
        return error("Error while gettting Application Policy");
    }
    resource function put throttling/policies/application/[string policyId](@http:Payload ApplicationThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns ApplicationThrottlePolicy|BadRequestError|NotFoundError|error {
        string?|ApplicationThrottlePolicy|NotFoundError appPolicy = updateApplicationUsagePlan(policyId, payload);
        if appPolicy is string {
            json j = check value:fromJsonString(appPolicy);
            ApplicationThrottlePolicy updatedPolicy = check j.cloneWithType(ApplicationThrottlePolicy);
            return updatedPolicy;
        }
        if appPolicy is ApplicationThrottlePolicy|NotFoundError {
            return appPolicy;
        } 
        io:println(appPolicy);
        return error("Error while updating Application Policy");
    }
    resource function delete throttling/policies/application/[string policyId]() returns http:Ok|NotFoundError|error {
        string|error? ex = removeApplicationUsagePlan(policyId);
        if ex is error {
            return ex;
        } else {
            return http:OK;
        }
    }
    resource function get throttling/policies/subscription(@http:Header string? accept = "application/json") returns SubscriptionThrottlePolicyList|NotAcceptableError|error {
        string?|SubscriptionThrottlePolicyList subPolicyList = getBusinessPlans();
        if subPolicyList is string {
            json j = check value:fromJsonString(subPolicyList);
            SubscriptionThrottlePolicyList polList = check j.cloneWithType(SubscriptionThrottlePolicyList);
            return polList;
        } else  if subPolicyList is SubscriptionThrottlePolicyList {
            return subPolicyList;
        }
        io:print(subPolicyList);
        return error("Error while getting Business plans");
    }
    resource function post throttling/policies/subscription(@http:Payload SubscriptionThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns CreatedSubscriptionThrottlePolicy|BadRequestError|UnsupportedMediaTypeError|error {
        string?|SubscriptionThrottlePolicy createdSubPol = addBusinessPlan(payload);
        if createdSubPol is string {
            json j = check value:fromJsonString(createdSubPol);
            CreatedSubscriptionThrottlePolicy crPol = {body: check j.cloneWithType(SubscriptionThrottlePolicy)};
            return crPol;
        } else if createdSubPol is SubscriptionThrottlePolicy {
            CreatedSubscriptionThrottlePolicy crPol = {body: check createdSubPol.cloneWithType(SubscriptionThrottlePolicy)};
            return crPol;
        }
        io:println(createdSubPol);
        return error("Error while adding Business plan");
    }
    resource function get throttling/policies/subscription/[string policyId]() returns SubscriptionThrottlePolicy|NotFoundError|NotAcceptableError|error {
        string?|SubscriptionThrottlePolicy subPolicy = getBusinessPlanById(policyId);
        if subPolicy is string {
            json j = check value:fromJsonString(subPolicy);
            SubscriptionThrottlePolicy policy = check j.cloneWithType(SubscriptionThrottlePolicy);
            return policy;
        } else if subPolicy is SubscriptionThrottlePolicy {
            return subPolicy;
        }
        io:println(subPolicy);
        return error("Error while getting Business plan");
    }
    resource function put throttling/policies/subscription/[string policyId](@http:Payload SubscriptionThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns SubscriptionThrottlePolicy|BadRequestError|NotFoundError|error {
        string?|SubscriptionThrottlePolicy|NotFoundError  subPolicy = updateBusinessPlan(policyId, payload);
        if subPolicy is string {
            json j = check value:fromJsonString(subPolicy);
            SubscriptionThrottlePolicy updatedPolicy = check j.cloneWithType(SubscriptionThrottlePolicy);
            return updatedPolicy;
        } else if subPolicy is SubscriptionThrottlePolicy | NotFoundError {
            return subPolicy;
        }
        io:println(subPolicy);
        return error("Error while updating Business plan");
    }
    resource function delete throttling/policies/subscription/[string policyId]() returns http:Ok|NotFoundError|error {
        string|error? ex = removeBusinessPlan(policyId);
        if ex is error {
            return ex;
        } else {
            return http:OK;
        }
    }
    // resource function get throttling/policies/advanced(@http:Header string? accept = "application/json") returns AdvancedThrottlePolicyList|NotAcceptableError {
    //     return advancedPolicyList;
    // }
    // resource function post throttling/policies/advanced(@http:Payload AdvancedThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns CreatedAdvancedThrottlePolicy|BadRequestError|UnsupportedMediaTypeError {
    //     io:println("Created Advanced Policy: " + payload.get("policyName").toString());
    //     return {body: policyCreated};
    // }
    // resource function get throttling/policies/advanced/[string policyId]() returns AdvancedThrottlePolicy|NotFoundError|NotAcceptableError {
    // }
    // resource function put throttling/policies/advanced/[string policyId](@http:Payload AdvancedThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns AdvancedThrottlePolicy|BadRequestError|NotFoundError {
    // }
    // resource function delete throttling/policies/advanced/[string policyId]() returns http:Ok|NotFoundError {
    // }
    // resource function get throttling/policies/export(string? policyId, string? name, string? 'type, string? format) returns ExportThrottlePolicy|NotFoundError|InternalServerErrorError {
    // }
    // resource function post throttling/policies/'import(boolean? overwrite, @http:Payload json payload) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|InternalServerErrorError {
    // }
    resource function get throttling/'deny\-policies(@http:Header string? accept = "application/json") returns BlockingConditionList|NotAcceptableError|error {
        string?|BlockingConditionList conditionList = getAllDenyPolicies();
        if conditionList is string {
            json j = check value:fromJsonString(conditionList);
            BlockingConditionList list = check j.cloneWithType(BlockingConditionList);
            return list;
        } else if conditionList is BlockingConditionList {
            return conditionList;
        }
        io:println(conditionList);
        return error("Error while getting block conditions");
    }
    resource function post throttling/'deny\-policies(@http:Payload BlockingCondition payload, @http:Header string 'content\-type = "application/json") returns CreatedBlockingCondition|BadRequestError|UnsupportedMediaTypeError|error {
        string?|BlockingCondition createdDenyPol = addDenyPolicy(payload);
        if createdDenyPol is string {
            json j = check value:fromJsonString(createdDenyPol);
            CreatedBlockingCondition condition = {body: check j.cloneWithType(BlockingCondition)};
            return condition;
        } else if createdDenyPol is BlockingCondition {
            CreatedBlockingCondition condition = {body: check createdDenyPol.cloneWithType(BlockingCondition)};
            return condition;
        }
        io:println(createdDenyPol);
        return error("Error while adding deny policy");
    }
    resource function get throttling/'deny\-policy/[string conditionId]() returns BlockingCondition|NotFoundError|NotAcceptableError|error {
        string?|BlockingCondition denyPolicy = getDenyPolicyById(conditionId);
        if denyPolicy is string {
            json j = check value:fromJsonString(denyPolicy);
            BlockingCondition condition = check j.cloneWithType(BlockingCondition);
            return condition;
        } else if denyPolicy is BlockingCondition {
            return denyPolicy;
        }
        io:println(denyPolicy);
        return error("Error while getting deny policy");
    }
    resource function delete throttling/'deny\-policy/[string conditionId]() returns http:Ok|NotFoundError|error {
        string|error? ex = removeDenyPolicy(conditionId);
        if ex is error {
            return ex;
        } else {
            return http:OK;
        }
    }
    resource function patch throttling/'deny\-policy/[string conditionId](@http:Payload BlockingConditionStatus payload, @http:Header string 'content\-type = "application/json") returns BlockingCondition|BadRequestError|NotFoundError|error {
        string?|BlockingCondition|NotFoundError updatedPolicy = updateDenyPolicy(conditionId, payload);
        if updatedPolicy is string {
            json j = check value:fromJsonString(updatedPolicy);
            BlockingCondition condition = check j.cloneWithType(BlockingCondition);
            return condition;
        } else if updatedPolicy is BlockingCondition|NotFoundError {
            return updatedPolicy;
        }
        io:println(updatedPolicy);
        return error("Error while updating deny policy");
    }
    // resource function get applications(string? user, string? name, string? tenantDomain, int 'limit = 25, int offset = 0, @http:Header string? accept = "application/json", string sortBy = "name", string sortOrder = "asc") returns ApplicationList|BadRequestError|NotAcceptableError {
    // }
    // resource function get applications/[string applicationId]() returns Application|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    // resource function delete applications/[string applicationId]() returns http:Ok|AcceptedWorkflowResponse|NotFoundError {
    // }
    // resource function post applications/[string applicationId]/'change\-owner(string owner) returns http:Ok|BadRequestError|NotFoundError {
    // }
    // resource function get environments() returns EnvironmentList {
    // }
    // resource function post environments(@http:Payload Environment payload) returns CreatedEnvironment|BadRequestError {
    // }
    // resource function put environments/[string environmentId](@http:Payload Environment payload) returns Environment|BadRequestError|NotFoundError {
    // }
    // resource function delete environments/[string environmentId]() returns http:Ok|NotFoundError {
    // }
    // resource function get 'bot\-detection\-data() returns BotDetectionDataList|NotFoundError {
    // }
    // resource function post monetization/'publish\-usage() returns PublishStatus|AcceptedPublishStatus|NotFoundError|InternalServerErrorError {
    // }
    // resource function get monetization/'publish\-usage/status() returns MonetizationUsagePublishInfo {
    // }
    // resource function get workflows(string? workflowType, int 'limit = 25, int offset = 0, @http:Header string? accept = "application/json") returns WorkflowList|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    // resource function get workflows/[string externalWorkflowRef]() returns WorkflowInfo|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    // resource function post workflows/'update\-workflow\-status(string workflowReferenceId, @http:Payload Workflow payload) returns Workflow|BadRequestError|NotFoundError {
    // }
    // resource function get 'tenant\-info/[string username]() returns TenantInfo|NotFoundError|NotAcceptableError {
    // }
    // resource function get 'custom\-urls/[string tenantDomain]() returns CustomUrlInfo|NotFoundError|NotAcceptableError {
    // }
    // resource function get 'api\-categories() returns APICategoryList {
    // }
    // resource function post 'api\-categories(@http:Payload APICategory payload) returns CreatedAPICategory|BadRequestError {
    // }
    // resource function put 'api\-categories/[string apiCategoryId](@http:Payload APICategory payload) returns APICategory|BadRequestError|NotFoundError {
    // }
    // resource function delete 'api\-categories/[string apiCategoryId]() returns http:Ok|NotFoundError {
    // }
    // resource function get settings() returns Settings|NotFoundError {
    //     return settingPayload;
    // }
    // resource function get 'system\-scopes/[string scopeName](string? username) returns ScopeSettings|BadRequestError|NotFoundError {
    // }
    // resource function get 'system\-scopes() returns ScopeList|InternalServerErrorError {
    // }
    // resource function put 'system\-scopes(@http:Payload ScopeList payload) returns ScopeList|BadRequestError|InternalServerErrorError {
    // }
    // resource function get 'system\-scopes/'role\-aliases() returns RoleAliasList|NotFoundError {
    // }
    // resource function put 'system\-scopes/'role\-aliases(@http:Payload RoleAliasList payload) returns RoleAliasList|BadRequestError|InternalServerErrorError {
    // }
    // resource function head roles/[string roleId]() returns http:Ok|NotFoundError|InternalServerErrorError {
    // }
    // resource function get 'tenant\-theme() returns json|ForbiddenError|NotFoundError|InternalServerErrorError {
    // }
    // resource function put 'tenant\-theme(@http:Payload json payload) returns http:Ok|ForbiddenError|PayloadTooLargeError|InternalServerErrorError {
    // }
    // resource function get 'tenant\-config() returns string|ForbiddenError|NotFoundError|InternalServerErrorError {
    // }
    // resource function put 'tenant\-config(@http:Payload string payload) returns string|ForbiddenError|PayloadTooLargeError|InternalServerErrorError {
    // }
    // resource function get 'tenant\-config\-schema() returns string|ForbiddenError|NotFoundError|InternalServerErrorError {
    // }
    // resource function get 'key\-managers() returns KeyManagerList {
    // }
    // resource function post 'key\-managers(@http:Payload KeyManager payload) returns CreatedKeyManager|BadRequestError {
    // }
    // resource function get 'key\-managers/[string keyManagerId]() returns KeyManager|NotFoundError|NotAcceptableError {
    // }
    // resource function put 'key\-managers/[string keyManagerId](@http:Payload KeyManager payload) returns KeyManager|BadRequestError|NotFoundError {
    // }
    // resource function delete 'key\-managers/[string keyManagerId]() returns http:Ok|NotFoundError {
    // }
    // resource function post 'key\-managers/discover(@http:Payload json payload) returns KeyManagerWellKnownResponse {
    // }
}
