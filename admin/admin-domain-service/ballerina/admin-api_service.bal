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

service /api/am/admin on ep0 {
    // resource function get policies/search(string? query) returns PolicyDetailsList {
    // }
    isolated resource function get 'application\-rate\-plans(@http:Header string? accept = "application/json") returns ApplicationRatePlanList|NotAcceptableError|InternalServerErrorError|error {
        string?|ApplicationRatePlanList|error appPolicyList = getApplicationUsagePlans();
        if appPolicyList is string {
            json j = check value:fromJsonString(appPolicyList);
            ApplicationRatePlanList polList = check j.cloneWithType(ApplicationRatePlanList);
            return polList;
        } else if appPolicyList is ApplicationRatePlanList {
            log:printDebug(appPolicyList.toString());
            return appPolicyList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving all Application Rate Plans"}};
            return internalError;
        }
    }
    isolated resource function post 'application\-rate\-plans(@http:Payload ApplicationRatePlan payload, @http:Header string 'content\-type = "application/json") returns CreatedApplicationRatePlan|BadRequestError|UnsupportedMediaTypeError|InternalServerErrorError|error {
        string?|ApplicationRatePlan|error createdAppPol = addApplicationUsagePlan(payload);
        if createdAppPol is string {
            json j = check value:fromJsonString(createdAppPol);
            CreatedApplicationRatePlan crPol = {body: check j.cloneWithType(ApplicationRatePlan)};
            return crPol;
        } else if createdAppPol is ApplicationRatePlan {
            log:printDebug(createdAppPol.toString());
            CreatedApplicationRatePlan crPol = {body: check createdAppPol.cloneWithType(ApplicationRatePlan)};
            return crPol;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while generating Application Rate Plan"}};
            return internalError;
        }
    }
    isolated resource function get 'application\-rate\-plans/[string planId]() returns ApplicationRatePlan|NotFoundError|NotAcceptableError|InternalServerErrorError|error {
        string?|ApplicationRatePlan|error appPolicy = getApplicationUsagePlanById(planId);
        if appPolicy is string {
            json j = check value:fromJsonString(appPolicy);
            ApplicationRatePlan policy = check j.cloneWithType(ApplicationRatePlan);
            return policy;
        } else if appPolicy is ApplicationRatePlan {
            log:printDebug(appPolicy.toString());
            return appPolicy;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving Application Rate Plan By Id"}};
            return internalError;
        }
    }
    isolated resource function put 'application\-rate\-plans/[string planId](@http:Payload ApplicationRatePlan payload, @http:Header string 'content\-type = "application/json") returns ApplicationRatePlan|BadRequestError|NotFoundError|InternalServerErrorError|error {
        string?|ApplicationRatePlan|NotFoundError|error appPolicy = updateApplicationUsagePlan(planId, payload);
        if appPolicy is string {
            json j = check value:fromJsonString(appPolicy);
            ApplicationRatePlan updatedPolicy = check j.cloneWithType(ApplicationRatePlan);
            return updatedPolicy;
        } else if appPolicy is ApplicationRatePlan|NotFoundError {
            log:printDebug(appPolicy.toString());
            return appPolicy;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while updating Application Rate Plan By Id"}};
            return internalError;
        }
    }
    isolated resource function delete 'application\-rate\-plans/[string planId]() returns http:Ok|NotFoundError|InternalServerErrorError|error {
        string|error? ex = removeApplicationUsagePlan(planId);
        if ex is error {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while deleting Application Rate Plan By Id"}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    isolated resource function get 'business\-plans(@http:Header string? accept = "application/json") returns BusinessPlanList|NotAcceptableError|InternalServerErrorError|error {
        string?|BusinessPlanList|error subPolicyList = getBusinessPlans();
        if subPolicyList is string {
            json j = check value:fromJsonString(subPolicyList);
            BusinessPlanList polList = check j.cloneWithType(BusinessPlanList);
            return polList;
        } else  if subPolicyList is BusinessPlanList {
            log:printDebug(subPolicyList.toString());
            return subPolicyList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving list of Business Plans"}};
            return internalError;
        }
    }
    isolated resource function post 'business\-plans(@http:Payload BusinessPlan payload, @http:Header string 'content\-type = "application/json") returns CreatedBusinessPlan|BadRequestError|UnsupportedMediaTypeError|InternalServerErrorError|error {
        string?|BusinessPlan|error createdSubPol = addBusinessPlan(payload);
        if createdSubPol is string {
            json j = check value:fromJsonString(createdSubPol);
            CreatedBusinessPlan crPol = {body: check j.cloneWithType(BusinessPlan)};
            return crPol;
        } else if createdSubPol is BusinessPlan {
            log:printDebug(createdSubPol.toString());
            CreatedBusinessPlan crPol = {body: check createdSubPol.cloneWithType(BusinessPlan)};
            return crPol;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while adding Business Plan"}};
            return internalError;
        }
    }
    isolated resource function get 'business\-plans/[string planId]() returns BusinessPlan|NotFoundError|NotAcceptableError|InternalServerErrorError|error {
        string?|BusinessPlan|error subPolicy = getBusinessPlanById(planId);
        if subPolicy is string {
            json j = check value:fromJsonString(subPolicy);
            BusinessPlan policy = check j.cloneWithType(BusinessPlan);
            return policy;
        } else if subPolicy is BusinessPlan {
            log:printDebug(subPolicy.toString());
            return subPolicy;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving Business Plan by Id"}};
            return internalError;
        }
    }
    isolated resource function put 'business\-plans/[string planId](@http:Payload BusinessPlan payload, @http:Header string 'content\-type = "application/json") returns BusinessPlan|BadRequestError|NotFoundError|InternalServerErrorError|error {
        string?|BusinessPlan|NotFoundError|error  subPolicy = updateBusinessPlan(planId, payload);
        if subPolicy is string {
            json j = check value:fromJsonString(subPolicy);
            BusinessPlan updatedPolicy = check j.cloneWithType(BusinessPlan);
            return updatedPolicy;
        } else if subPolicy is BusinessPlan | NotFoundError {
            return subPolicy;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while updating Business Plan by Id"}};
            return internalError;
        }
    }
    isolated resource function delete 'business\-plans/[string planId]() returns http:Ok|NotFoundError|InternalServerErrorError|error{
        string|error? ex = removeBusinessPlan(planId);
        if ex is error {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while deleting Business Plan by Id"}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    // resource function get throttling/policies/advanced(@http:Header string? accept = "application/json") returns AdvancedThrottlePolicyList|NotAcceptableError {
    // }
    // resource function post throttling/policies/advanced(@http:Payload AdvancedThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns CreatedAdvancedThrottlePolicy|BadRequestError|UnsupportedMediaTypeError {
    // }
    // resource function get throttling/policies/advanced/[string policyId]() returns AdvancedThrottlePolicy|NotFoundError|NotAcceptableError {
    // }
    // resource function put throttling/policies/advanced/[string policyId](@http:Payload AdvancedThrottlePolicy payload, @http:Header string 'content\-type = "application/json") returns AdvancedThrottlePolicy|BadRequestError|NotFoundError {
    // }
    // resource function delete throttling/policies/advanced/[string policyId]() returns http:Ok|NotFoundError {
    // }
    // resource function get throttling/policies/export(string? policyId, string? name, string? 'type, string? format) returns ExportPolicy|NotFoundError|InternalServerErrorError {
    // }
    // resource function post throttling/policies/'import(boolean? overwrite, @http:Payload json payload) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|InternalServerErrorError {
    // }
    isolated resource function get 'deny\-policies(@http:Header string? accept = "application/json") returns BlockingConditionList|NotAcceptableError|InternalServerErrorError|error {
        string?|BlockingConditionList|error conditionList = getAllDenyPolicies();
        if conditionList is string {
            json j = check value:fromJsonString(conditionList);
            BlockingConditionList list = check j.cloneWithType(BlockingConditionList);
            return list;
        } else if conditionList is BlockingConditionList {
            return conditionList;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving all Deny Policies"}};
            return internalError;
        }
    }
    isolated resource function post 'deny\-policies(@http:Payload BlockingCondition payload, @http:Header string 'content\-type = "application/json") returns CreatedBlockingCondition|BadRequestError|UnsupportedMediaTypeError|InternalServerErrorError|error {
        string?|BlockingCondition|error createdDenyPol = addDenyPolicy(payload);
        if createdDenyPol is string {
            json j = check value:fromJsonString(createdDenyPol);
            CreatedBlockingCondition condition = {body: check j.cloneWithType(BlockingCondition)};
            return condition;
        } else if createdDenyPol is BlockingCondition {
            log:printDebug(createdDenyPol.toString());
            CreatedBlockingCondition condition = {body: check createdDenyPol.cloneWithType(BlockingCondition)};
            return condition;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while adding Deny Policy"}};
            return internalError;
        }
    }
    isolated resource function get 'deny\-policies/[string policyId]() returns BlockingCondition|NotFoundError|NotAcceptableError|InternalServerErrorError|error {
        string?|BlockingCondition|error denyPolicy = getDenyPolicyById(policyId);
        if denyPolicy is string {
            json j = check value:fromJsonString(denyPolicy);
            BlockingCondition condition = check j.cloneWithType(BlockingCondition);
            return condition;
        } else if denyPolicy is BlockingCondition {
            log:printDebug(denyPolicy.toString());
            return denyPolicy;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while retrieving Deny Policy by Id"}};
            return internalError;
        }
    }
    isolated resource function delete 'deny\-policies/[string policyId]() returns http:Ok|NotFoundError|InternalServerErrorError|error {
        string|error? ex = removeDenyPolicy(policyId);
        if ex is error {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while deleting Deny Policy by Id"}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    isolated resource function patch 'deny\-policies/[string policyId](@http:Payload BlockingConditionStatus payload, @http:Header string 'content\-type = "application/json") returns BlockingCondition|BadRequestError|NotFoundError|InternalServerErrorError|error {
        string?|BlockingCondition|NotFoundError|error updatedPolicy = updateDenyPolicy(policyId, payload);
        if updatedPolicy is string {
            json j = check value:fromJsonString(updatedPolicy);
            BlockingCondition condition = check j.cloneWithType(BlockingCondition);
            return condition;
        } else if updatedPolicy is BlockingCondition|NotFoundError {
            log:printDebug(updatedPolicy.toString());
            return updatedPolicy;
        } else {
            InternalServerErrorError internalError = {body: {code: 90900, message: "Internal Error while updating Deny Policy Status by Id"}};
            return internalError;
        }
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
