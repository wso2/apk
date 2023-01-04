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

isolated function addApplicationUsagePlan(ApplicationRatePlan body) returns string?|ApplicationRatePlan|error {
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
    string?|ApplicationRatePlan|error policy = addApplicationUsagePlanDAO(body);
    return policy;
}

isolated function getApplicationUsagePlanById(string policyId) returns string?|ApplicationRatePlan|error {
    string?|ApplicationRatePlan|error policy = getApplicationUsagePlanByIdDAO(policyId);
    return policy;
}

isolated function getApplicationUsagePlans() returns string?|ApplicationRatePlanList|error {
    string org = "carbon.super";
    ApplicationRatePlan[]|error? usagePlans = getApplicationUsagePlansDAO(org);
    if usagePlans is ApplicationRatePlan[] {
        int count = usagePlans.length();
        ApplicationRatePlanList usagePlansList = {count: count, list: usagePlans};
        return usagePlansList;
    } else {
        return usagePlans;
    }
}

isolated function updateApplicationUsagePlan(string policyId, ApplicationRatePlan body) returns string?|ApplicationRatePlan|NotFoundError|error {
    string?|ApplicationRatePlan|error existingPolicy = getApplicationUsagePlanByIdDAO(policyId);
    if existingPolicy is ApplicationRatePlan {
        body.planId = policyId;
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
    string?|ApplicationRatePlan|error policy = updateApplicationUsagePlanDAO(body);
    return policy;
}

isolated function removeApplicationUsagePlan(string policyId) returns error?|string {
    error?|string status = deleteApplicationUsagePlanDAO(policyId);
    return status;
}

isolated function addBusinessPlan(BusinessPlan body) returns string?|BusinessPlan|error {
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
    string?|BusinessPlan|error policy = addBusinessPlanDAO(body);
    return policy;
}

isolated function getBusinessPlanById(string policyId) returns string?|BusinessPlan|error {
    string?|BusinessPlan|error policy = getBusinessPlanByIdDAO(policyId);
    return policy;
}

isolated function getBusinessPlans() returns string?|BusinessPlanList|error {
    string org = "carbon.super";
    BusinessPlan[]|error? businessPlans = getBusinessPlansDAO(org);
    if businessPlans is BusinessPlan[] {
        int count = businessPlans.length();
        BusinessPlanList BusinessPlansList = {count: count, list: businessPlans};
        return BusinessPlansList;
    } else {
        return businessPlans;
    }
}

isolated function updateBusinessPlan(string policyId, BusinessPlan body) returns string?|BusinessPlan|NotFoundError|error {
    string?|BusinessPlan|error existingPolicy = getBusinessPlanByIdDAO(policyId);
    if existingPolicy is BusinessPlan {
        body.planId = policyId;
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
    string?|BusinessPlan|error policy = updateBusinessPlanDAO(body);
    return policy;
}

isolated function removeBusinessPlan(string policyId) returns error?|string {
    error?|string status = deleteBusinessPlanDAO(policyId);
    return status;
}

isolated function addDenyPolicy(BlockingCondition body) returns string?|BlockingCondition|error {
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
    string?|BlockingCondition|error policy = addDenyPolicyDAO(body);
    return policy;
}

isolated function getDenyPolicyById(string policyId) returns string?|BlockingCondition|error {
    string?|BlockingCondition|error policy = getDenyPolicyByIdDAO(policyId);
    return policy;
}

isolated function getAllDenyPolicies() returns string?|BlockingConditionList|error {
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

isolated function updateDenyPolicy(string policyId, BlockingConditionStatus status) returns string?|BlockingCondition|NotFoundError|error {
    string?|BlockingCondition|error existingPolicy = getDenyPolicyByIdDAO(policyId);
    if existingPolicy !is BlockingCondition {
        Error err = {code:9010101, message:"Policy Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    } else {
        status.policyId = policyId;
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

isolated function removeDenyPolicy(string policyId) returns error?|string {
    error?|string status = deleteDenyPolicyDAO(policyId);
    return status;
}

