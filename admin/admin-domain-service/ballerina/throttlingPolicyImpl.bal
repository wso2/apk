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

import ballerina/log;
import ballerina/uuid;

function addApplicationUsagePlan(ApplicationThrottlePolicy body) returns string?|ApplicationThrottlePolicy|error {
    string policyId = uuid:createType1AsString();
    body.policyId = policyId;
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
    string?|ApplicationThrottlePolicy|error policy = addApplicationUsagePlanDAO(body);
    return policy;
}

function getApplicationUsagePlanById(string policyId) returns string?|ApplicationThrottlePolicy|error {
    string?|ApplicationThrottlePolicy|error policy = getApplicationUsagePlanByIdDAO(policyId);
    return policy;
}

function getApplicationUsagePlans() returns string?|ApplicationThrottlePolicyList|error {
    string org = "carbon.super";
    ApplicationThrottlePolicy[]|error? usagePlans = getApplicationUsagePlansDAO(org);
    if usagePlans is ApplicationThrottlePolicy[] {
        int count = usagePlans.length();
        ApplicationThrottlePolicyList usagePlansList = {count: count, list: usagePlans};
        return usagePlansList;
    } else {
        return usagePlans;
    }
}

function updateApplicationUsagePlan(string policyId, ApplicationThrottlePolicy body) returns string?|ApplicationThrottlePolicy|NotFoundError|error {
    string?|ApplicationThrottlePolicy|error existingPolicy = getApplicationUsagePlanByIdDAO(policyId);
    if existingPolicy is ApplicationThrottlePolicy {
        body.policyId = policyId;
        //body.policyName = existingPolicy.name;
    } else {
        Error err = {code:9010101, message:"Policy Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
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
    string?|ApplicationThrottlePolicy|error policy = updateApplicationUsagePlanDAO(body);
    return policy;
}

function removeApplicationUsagePlan(string policyId) returns error?|string {
    error?|string status = deleteApplicationUsagePlanDAO(policyId);
    return status;
}

function addBusinessPlan(SubscriptionThrottlePolicy body) returns string?|SubscriptionThrottlePolicy|error {
    string policyId = uuid:createType1AsString();
    body.policyId = policyId;
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
    string?|SubscriptionThrottlePolicy|error policy = addBusinessPlanDAO(body);
    return policy;
}

function getBusinessPlanById(string policyId) returns string?|SubscriptionThrottlePolicy|error {
    string?|SubscriptionThrottlePolicy|error policy = getBusinessPlanByIdDAO(policyId);
    return policy;
}

function getBusinessPlans() returns string?|SubscriptionThrottlePolicyList|error {
    string org = "carbon.super";
    SubscriptionThrottlePolicy[]|error? businessPlans = getBusinessPlansDAO(org);
    if businessPlans is SubscriptionThrottlePolicy[] {
        int count = businessPlans.length();
        SubscriptionThrottlePolicyList BusinessPlansList = {count: count, list: businessPlans};
        return BusinessPlansList;
    } else {
        return businessPlans;
    }
}

function updateBusinessPlan(string policyId, SubscriptionThrottlePolicy body) returns string?|SubscriptionThrottlePolicy|NotFoundError|error {
    string?|SubscriptionThrottlePolicy|error existingPolicy = getBusinessPlanByIdDAO(policyId);
    if existingPolicy is SubscriptionThrottlePolicy {
        body.policyId = policyId;
        //body.policyName = existingPolicy.name;
    } else {
        Error err = {code:9010101, message:"Policy Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
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
    string?|SubscriptionThrottlePolicy|error policy = updateBusinessPlanDAO(body);
    return policy;
}

function removeBusinessPlan(string policyId) returns error?|string {
    error?|string status = deleteBusinessPlanDAO(policyId);
    return status;
}

function addDenyPolicy(BlockingCondition body) returns string?|BlockingCondition|error {
    string policyId = uuid:createType1AsString();
    body.conditionId = policyId;
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
    string?|BlockingCondition|error policy = addDenyPolicyDAO(body);
    return policy;
}

function getDenyPolicyById(string policyId) returns string?|BlockingCondition|error {
    string?|BlockingCondition|error policy = getDenyPolicyByIdDAO(policyId);
    return policy;
}

function getAllDenyPolicies() returns string?|BlockingConditionList|error {
    string org = "carbon.super";
    BlockingCondition[]|error? denyPolicies = getDenyPoliciesDAO(org);
    if denyPolicies is BlockingCondition[] {
        int count = denyPolicies.length();
        BlockingConditionList denyPoliciesList = {count: count, list: denyPolicies};
        return denyPoliciesList;
    } else {
        log:printError("Error");
        return error("Error while retrieving all deny polcies from DB");
    }
}

function updateDenyPolicy(string policyId, BlockingConditionStatus status) returns string?|BlockingCondition|NotFoundError|error {
    string?|BlockingCondition|error existingPolicy = getDenyPolicyByIdDAO(policyId);
    if existingPolicy !is BlockingCondition {
        Error err = {code:9010101, message:"Policy Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    } else {
        status.conditionId = policyId;
    }
    string?|error response = updateDenyPolicyDAO(status);
    if response is error {
        return response;
    }
    string?|BlockingCondition|error updatedPolicy = getDenyPolicyByIdDAO(policyId);
    if updatedPolicy is BlockingCondition {
        return updatedPolicy;
    } else {
        return updatedPolicy;
    }
}

function removeDenyPolicy(string policyId) returns error?|string {
    error?|string status = deleteDenyPolicyDAO(policyId);
    return status;
}

