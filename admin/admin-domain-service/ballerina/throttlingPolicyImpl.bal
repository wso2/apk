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

import ballerina/uuid;

isolated function addApplicationUsagePlan(ApplicationRatePlan body) returns ApplicationRatePlan|APKError {
    string policyId = uuid:createType1AsString();
    body.planId = policyId;
    match body.defaultLimit.'type {
        "REQUESTCOUNTLIMIT" => {
            body.defaultLimit.'type = "requestCount";
        }
        "BANDWIDTHLIMIT" => {
            body.defaultLimit.'type = "bandwidth";
        }
        "EVENTCOUNTLIMIT" => {
            body.defaultLimit.'type = "eventCount";
        }
    }
    ApplicationRatePlan|APKError policy = addApplicationUsagePlanDAO(body);
    return policy;
}

isolated function getApplicationUsagePlanById(string policyId) returns ApplicationRatePlan|APKError|NotFoundError {
    ApplicationRatePlan|APKError|NotFoundError policy = getApplicationUsagePlanByIdDAO(policyId);
    return policy;
}

isolated function getApplicationUsagePlans() returns ApplicationRatePlanList|APKError {
    string org = "carbon.super";
    ApplicationRatePlan[]|APKError usagePlans = getApplicationUsagePlansDAO(org);
    if usagePlans is ApplicationRatePlan[] {
        int count = usagePlans.length();
        ApplicationRatePlanList usagePlansList = {count: count, list: usagePlans};
        return usagePlansList;
    } else {
        return usagePlans;
    }
}

isolated function updateApplicationUsagePlan(string policyId, ApplicationRatePlan body) returns ApplicationRatePlan|NotFoundError|APKError {
    ApplicationRatePlan|APKError|NotFoundError existingPolicy = getApplicationUsagePlanByIdDAO(policyId);
    if existingPolicy is ApplicationRatePlan {
        body.planId = policyId;
        //body.policyName = existingPolicy.name;
    } else {
        return existingPolicy;
    }

    match body.defaultLimit.'type {
        "REQUESTCOUNTLIMIT" => {
            body.defaultLimit.'type = "requestCount";
        }
        "BANDWIDTHLIMIT" => {
            body.defaultLimit.'type = "bandwidth";
        }
        "EVENTCOUNTLIMIT" => {
            body.defaultLimit.'type = "eventCount";
        }
    }
    ApplicationRatePlan|APKError policy = updateApplicationUsagePlanDAO(body);
    return policy;
}

isolated function removeApplicationUsagePlan(string policyId) returns APKError|string {
    APKError|string status = deleteApplicationUsagePlanDAO(policyId);
    return status;
}

isolated function addBusinessPlan(BusinessPlan body) returns BusinessPlan|APKError {
    string policyId = uuid:createType1AsString();
    body.planId = policyId;
    match body.defaultLimit.'type {
        "REQUESTCOUNTLIMIT" => {
            body.defaultLimit.'type = "requestCount";
        }
        "BANDWIDTHLIMIT" => {
            body.defaultLimit.'type = "bandwidth";
        }
        "EVENTCOUNTLIMIT" => {
            body.defaultLimit.'type = "eventCount";
        }
    }
    BusinessPlan|APKError policy = addBusinessPlanDAO(body);
    return policy;
}

isolated function getBusinessPlanById(string policyId) returns BusinessPlan|APKError|NotFoundError {
    BusinessPlan|APKError|NotFoundError policy = getBusinessPlanByIdDAO(policyId);
    return policy;
}

isolated function getBusinessPlans() returns BusinessPlanList|APKError {
    string org = "carbon.super";
    BusinessPlan[]|APKError businessPlans = getBusinessPlansDAO(org);
    if businessPlans is BusinessPlan[] {
        int count = businessPlans.length();
        BusinessPlanList BusinessPlansList = {count: count, list: businessPlans};
        return BusinessPlansList;
    } else {
        return businessPlans;
    }
}

isolated function updateBusinessPlan(string policyId, BusinessPlan body) returns BusinessPlan|NotFoundError|APKError {
    BusinessPlan|APKError|NotFoundError existingPolicy = getBusinessPlanByIdDAO(policyId);
    if existingPolicy is BusinessPlan {
        body.planId = policyId;
        //body.policyName = existingPolicy.name;
    } else {
       return existingPolicy;
    }
    match body.defaultLimit.'type {
        "REQUESTCOUNTLIMIT" => {
            body.defaultLimit.'type = "requestCount";
        }
        "BANDWIDTHLIMIT" => {
            body.defaultLimit.'type = "bandwidth";
        }
        "EVENTCOUNTLIMIT" => {
            body.defaultLimit.'type = "eventCount";
        }
    }
    BusinessPlan|APKError policy = updateBusinessPlanDAO(body);
    return policy;
}

isolated function removeBusinessPlan(string policyId) returns APKError|string {
    APKError|string status = deleteBusinessPlanDAO(policyId);
    return status;
}

isolated function addDenyPolicy(BlockingCondition body) returns BlockingCondition|APKError {
    string policyId = uuid:createType1AsString();
    body.policyId = policyId;
    //Todo : need to validate each type
    match body.conditionType {
        "APPLICATION" => {
        }
        "API" => {
        }
        "IP" => {
        }
        "IPRANGE" => {
        }
        "USER" => {
        }
    }
    BlockingCondition|APKError policy = addDenyPolicyDAO(body);
    return policy;
}

isolated function getDenyPolicyById(string policyId) returns BlockingCondition|APKError|NotFoundError {
    BlockingCondition|APKError|NotFoundError policy = getDenyPolicyByIdDAO(policyId);
    return policy;
}

isolated function getAllDenyPolicies() returns BlockingConditionList|APKError {
    string org = "carbon.super";
    BlockingCondition[]|APKError denyPolicies = getDenyPoliciesDAO(org);
    if denyPolicies is BlockingCondition[] {
        int count = denyPolicies.length();
        BlockingConditionList denyPoliciesList = {count: count, list: denyPolicies};
        return denyPoliciesList;
    } else {
       return denyPolicies;
    }
}

isolated function updateDenyPolicy(string policyId, BlockingConditionStatus status) returns BlockingCondition|NotFoundError|APKError {
    BlockingCondition|APKError|NotFoundError existingPolicy = getDenyPolicyByIdDAO(policyId);
    if existingPolicy !is BlockingCondition {
        return existingPolicy;
    } else {
        status.policyId = policyId;
    }
    string|APKError response = updateDenyPolicyDAO(status);
    if response is APKError{
        return response;
    }
    BlockingCondition|APKError|NotFoundError updatedPolicy = getDenyPolicyByIdDAO(policyId);
    if updatedPolicy is BlockingCondition {
        return updatedPolicy;
    } else {
        return updatedPolicy;
    }
}

isolated function removeDenyPolicy(string policyId) returns APKError|string {
    APKError|string status = deleteDenyPolicyDAO(policyId);
    return status;
}

