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
import wso2/apk_common_lib as commons;

isolated function addApplicationUsagePlan(ApplicationRatePlan body, commons:Organization org) returns ApplicationRatePlan|commons:APKError {
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
    ApplicationRatePlan|commons:APKError policy = addApplicationUsagePlanDAO(body , org.uuid);
    return policy;
}

isolated function getApplicationUsagePlanById(string policyId, commons:Organization org) returns ApplicationRatePlan|commons:APKError {
    ApplicationRatePlan|commons:APKError policy = getApplicationUsagePlanByIdDAO(policyId, org.uuid);
    return policy;
}

isolated function getApplicationUsagePlans(commons:Organization org) returns ApplicationRatePlanList|commons:APKError {
    ApplicationRatePlan[]|commons:APKError usagePlans = getApplicationUsagePlansDAO(org.uuid);
    if usagePlans is ApplicationRatePlan[] {
        int count = usagePlans.length();
        ApplicationRatePlanList usagePlansList = {count: count, list: usagePlans};
        return usagePlansList;
    } else {
        return usagePlans;
    }
}

isolated function updateApplicationUsagePlan(string policyId, ApplicationRatePlan body, commons:Organization org) returns ApplicationRatePlan|commons:APKError {
    ApplicationRatePlan|commons:APKError existingPolicy = getApplicationUsagePlanByIdDAO(policyId, org.uuid);
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
    ApplicationRatePlan|commons:APKError policy = updateApplicationUsagePlanDAO(body, org.uuid);
    return policy;
}

isolated function removeApplicationUsagePlan(string policyId, commons:Organization org) returns commons:APKError|string {
    commons:APKError|string status = deleteApplicationUsagePlanDAO(policyId, org.uuid);
    return status;
}

isolated function addBusinessPlan(BusinessPlan body, commons:Organization org) returns BusinessPlan|commons:APKError {
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
    BusinessPlan|commons:APKError policy = addBusinessPlanDAO(body, org.uuid);
    return policy;
}

isolated function getBusinessPlanById(string policyId, commons:Organization org) returns BusinessPlan|commons:APKError {
    BusinessPlan|commons:APKError policy = getBusinessPlanByIdDAO(policyId, org.uuid);
    return policy;
}

isolated function getBusinessPlans(commons:Organization org) returns BusinessPlanList|commons:APKError {
    BusinessPlan[]|commons:APKError businessPlans = getBusinessPlansDAO(org.uuid);
    if businessPlans is BusinessPlan[] {
        int count = businessPlans.length();
        BusinessPlanList BusinessPlansList = {count: count, list: businessPlans};
        return BusinessPlansList;
    } else {
        return businessPlans;
    }
}

isolated function updateBusinessPlan(string policyId, BusinessPlan body, commons:Organization org) returns BusinessPlan|commons:APKError {
    BusinessPlan|commons:APKError existingPolicy = getBusinessPlanByIdDAO(policyId, org.uuid);
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
    BusinessPlan|commons:APKError policy = updateBusinessPlanDAO(body, org.uuid);
    return policy;
}

isolated function removeBusinessPlan(string policyId, commons:Organization org) returns commons:APKError|string {
    commons:APKError|string status = deleteBusinessPlanDAO(policyId, org.uuid);
    return status;
}

isolated function addDenyPolicy(BlockingCondition body, commons:Organization org) returns BlockingCondition|commons:APKError {
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
    BlockingCondition|commons:APKError policy = addDenyPolicyDAO(body, org.uuid);
    return policy;
}

isolated function getDenyPolicyById(string policyId, commons:Organization org) returns BlockingCondition|commons:APKError {
    BlockingCondition|commons:APKError policy = getDenyPolicyByIdDAO(policyId, org.uuid);
    return policy;
}

isolated function getAllDenyPolicies(commons:Organization org) returns BlockingConditionList|commons:APKError {
    BlockingCondition[]|commons:APKError denyPolicies = getDenyPoliciesDAO(org.uuid);
    if denyPolicies is BlockingCondition[] {
        int count = denyPolicies.length();
        BlockingConditionList denyPoliciesList = {count: count, list: denyPolicies};
        return denyPoliciesList;
    } else {
       return denyPolicies;
    }
}

isolated function updateDenyPolicy(string policyId, BlockingConditionStatus status, commons:Organization org) returns BlockingCondition|commons:APKError {
    BlockingCondition|commons:APKError existingPolicy = getDenyPolicyByIdDAO(policyId, org.uuid);
    if existingPolicy !is BlockingCondition {
        return existingPolicy;
    } else {
        status.policyId = policyId;
    }
    string|commons:APKError response = updateDenyPolicyDAO(status);
    if response is commons:APKError{
        return response;
    }
    BlockingCondition|commons:APKError updatedPolicy = getDenyPolicyByIdDAO(policyId, org.uuid);
    if updatedPolicy is BlockingCondition {
        return updatedPolicy;
    } else {
        return updatedPolicy;
    }
}

isolated function removeDenyPolicy(string policyId, commons:Organization org) returns commons:APKError|string {
    commons:APKError|string status = deleteDenyPolicyDAO(policyId, org.uuid);
    return status;
}

