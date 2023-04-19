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
import wso2/apk_common_lib as commons;

service /api/am/admin on ep0 {
    // resource function get policies/search(string? query) returns PolicyDetailsList {
    // }
    isolated resource function get 'application\-rate\-plans(@http:Header string? accept = "application/json") returns ApplicationRatePlanList|commons:APKError {
        ApplicationRatePlanList|commons:APKError appPolicyList = getApplicationUsagePlans();
        if appPolicyList is ApplicationRatePlanList {
            log:printDebug(appPolicyList.toString());
            return appPolicyList;
        } else {
            return handleAPKError(appPolicyList);
        }
    }
    isolated resource function post 'application\-rate\-plans(@http:Payload ApplicationRatePlan payload, @http:Header string 'content\-type = "application/json") returns ApplicationRatePlan|commons:APKError {
        ApplicationRatePlan|commons:APKError createdAppPol = addApplicationUsagePlan(payload);
        if createdAppPol is ApplicationRatePlan {
            log:printDebug(createdAppPol.toString());
            return createdAppPol;
        } else {
            return handleAPKError(createdAppPol);
        }
    }
    isolated resource function get 'application\-rate\-plans/[string planId]() returns ApplicationRatePlan|commons:APKError {
        ApplicationRatePlan|commons:APKError appPolicy = getApplicationUsagePlanById(planId);
        if appPolicy is ApplicationRatePlan {
            log:printDebug(appPolicy.toString());
            return appPolicy;
        } else {
            return handleAPKError(appPolicy);
        }
    }
    isolated resource function put 'application\-rate\-plans/[string planId](@http:Payload ApplicationRatePlan payload, @http:Header string 'content\-type = "application/json") returns ApplicationRatePlan|commons:APKError {
        ApplicationRatePlan|commons:APKError appPolicy = updateApplicationUsagePlan(planId, payload);
        if appPolicy is ApplicationRatePlan {
            log:printDebug(appPolicy.toString());
            return appPolicy;
        } else {
            return handleAPKError(appPolicy);
        }
    }
    isolated resource function delete 'application\-rate\-plans/[string planId]() returns http:Ok|commons:APKError {
        string|commons:APKError ex = removeApplicationUsagePlan(planId);
        if ex is commons:APKError {
            return handleAPKError(ex);
        } else {
            return http:OK;
        }
    }
    isolated resource function get 'business\-plans(@http:Header string? accept = "application/json") returns BusinessPlanList|commons:APKError {
        BusinessPlanList|commons:APKError subPolicyList = getBusinessPlans();
        if subPolicyList is BusinessPlanList {
            log:printDebug(subPolicyList.toString());
            return subPolicyList;
        } else {
            return handleAPKError(subPolicyList);
        }
    }
    isolated resource function post 'business\-plans(@http:Payload BusinessPlan payload, @http:Header string 'content\-type = "application/json") returns BusinessPlan|commons:APKError {
        BusinessPlan|commons:APKError createdSubPol = addBusinessPlan(payload);
        if createdSubPol is BusinessPlan {
            log:printDebug(createdSubPol.toString());
            return createdSubPol;
        } else {
            return handleAPKError(createdSubPol);
        }
    }
    isolated resource function get 'business\-plans/[string planId]() returns BusinessPlan|commons:APKError {
        BusinessPlan|commons:APKError subPolicy = getBusinessPlanById(planId);
        if subPolicy is BusinessPlan {
            log:printDebug(subPolicy.toString());
            return subPolicy;
        } else {
            return handleAPKError(subPolicy);
        }
    }
    isolated resource function put 'business\-plans/[string planId](@http:Payload BusinessPlan payload, @http:Header string 'content\-type = "application/json") returns BusinessPlan|commons:APKError {
        BusinessPlan|commons:APKError  subPolicy = updateBusinessPlan(planId, payload);
        if subPolicy is BusinessPlan {
            return subPolicy;
        } else {
            return handleAPKError(subPolicy);
        }
    }
    isolated resource function delete 'business\-plans/[string planId]() returns http:Ok|commons:APKError{
        string|commons:APKError ex = removeBusinessPlan(planId);
        if ex is commons:APKError {
            return handleAPKError(ex);
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
    isolated resource function get 'deny\-policies(@http:Header string? accept = "application/json") returns BlockingConditionList|commons:APKError {
        BlockingConditionList|commons:APKError conditionList = getAllDenyPolicies();
        if conditionList is BlockingConditionList {
            return conditionList;
        } else {
            return handleAPKError(conditionList);
        }
    }
    isolated resource function post 'deny\-policies(@http:Payload BlockingCondition payload, @http:Header string 'content\-type = "application/json") returns BlockingCondition|commons:APKError {
        BlockingCondition|commons:APKError createdDenyPol = addDenyPolicy(payload);
        if createdDenyPol is BlockingCondition {
            log:printDebug(createdDenyPol.toString());
            return createdDenyPol;
        } else {
            return handleAPKError(createdDenyPol);
        }
    }
    isolated resource function get 'deny\-policies/[string policyId]() returns BlockingCondition|commons:APKError {
        BlockingCondition|commons:APKError denyPolicy = getDenyPolicyById(policyId);
        if denyPolicy is BlockingCondition {
            log:printDebug(denyPolicy.toString());
            return denyPolicy;
        } else {
            return handleAPKError(denyPolicy);
        }
    }
    isolated resource function delete 'deny\-policies/[string policyId]() returns http:Ok|commons:APKError {
        string|commons:APKError ex = removeDenyPolicy(policyId);
        if ex is commons:APKError {
            return handleAPKError(ex);
        } else {
            return http:OK;
        }
    }
    isolated resource function patch 'deny\-policies/[string policyId](@http:Payload BlockingConditionStatus payload, @http:Header string 'content\-type = "application/json") returns BlockingCondition|commons:APKError {
        BlockingCondition|commons:APKError updatedPolicy = updateDenyPolicy(policyId, payload);
        if updatedPolicy is BlockingCondition {
            log:printDebug(updatedPolicy.toString());
            return updatedPolicy;
        } else {
            return handleAPKError(updatedPolicy);
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
    isolated resource function get 'api\-categories() returns APICategoryList|commons:APKError {
        APICategoryList|commons:APKError apiCategoryList = getAllCategoryList();
        if apiCategoryList is APICategoryList {
            return apiCategoryList;
        } else {
            return handleAPKError(apiCategoryList);
        }
    }
    isolated resource function post 'api\-categories(@http:Payload APICategory payload) returns APICategory|commons:APKError {
        APICategory|commons:APKError createdApiCategory = addAPICategory(payload);
        if createdApiCategory is APICategory {
            return createdApiCategory;
        } else {
            return handleAPKError(createdApiCategory);
        }
    }
    isolated resource function put 'api\-categories/[string apiCategoryId](@http:Payload APICategory payload) returns APICategory|commons:APKError {
        APICategory|commons:APKError  apiCategory = updateAPICategory(apiCategoryId, payload);
        if apiCategory is APICategory {
            return apiCategory;
        } else {
            return handleAPKError(apiCategory);
        }
    }
    isolated resource function delete 'api\-categories/[string apiCategoryId]() returns http:Ok|commons:APKError {
        string|commons:APKError ex = removeAPICategory(apiCategoryId);
        if ex is commons:APKError {
            return handleAPKError(ex);
        } else {
            return http:OK;
        }
    }
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
    isolated resource function get organizations() returns OrganizationList|commons:APKError {
        OrganizationList|commons:APKError getAllOrganizationList = getAllOrganization();
        if getAllOrganizationList is OrganizationList {
            return getAllOrganizationList;
        } else {
            return handleAPKError(getAllOrganizationList);
        }
    }


    isolated resource function post organizations(@http:Payload Organization payload) returns Organization|commons:APKError {
        Organization|commons:APKError createdOrganization = addOrganization(payload);
        if createdOrganization is Organization {
            return createdOrganization;
        } else {
            return handleAPKError(createdOrganization);
        }
    }

    isolated resource function get organizations/[string organizationId]() returns Organization|commons:APKError {
        Organization|commons:APKError getOrganization = getOrganizationById(organizationId);
        if getOrganization is Organization {
            return getOrganization;
        } else {
            return handleAPKError(getOrganization);
        }
    }

    isolated resource function put organizations/[string organizationId](@http:Payload Organization payload) returns Organization|commons:APKError {
        Organization|commons:APKError updateOrganization = updatedOrganization(organizationId, payload);
        if updateOrganization is Organization {
            return updateOrganization;
        } else {
            return handleAPKError(updateOrganization);
        }
    }

    resource function delete organizations/[string organizationId]() returns http:Ok|commons:APKError {
        boolean|commons:APKError deleteOrganization = removeOrganization(organizationId);
        if deleteOrganization is commons:APKError {
            return handleAPKError(deleteOrganization);
        } else {
            return http:OK;
        }
    }
    resource function get 'organization\-info() returns Organization|commons:APKError {
        Organization|commons:APKError getOrganization = getOrganizationByOrganizationClaim();
        if getOrganization is Organization {
            return getOrganization;
        } else {
            return handleAPKError(getOrganization);
        }
    }
}

isolated function handleAPKError(commons:APKError errorDetail) returns commons:APKError {
    commons:ErrorHandler & readonly detail = errorDetail.detail();
    return error commons:APKError( detail.message,
        code = 909400,
        message = detail.message,
        statusCode = detail.code,
        description = detail.message
    ); 
}
