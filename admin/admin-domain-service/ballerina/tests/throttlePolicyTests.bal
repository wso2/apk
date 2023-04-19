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

import ballerina/test;
import ballerina/log;
import wso2/apk_common_lib as commons;

ApplicationRatePlan  applicationUsagePlan = {planName: "25PerMin", description: "25 Per Minute",
'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}};

BusinessPlan  businessPlan = {planName: "MySubPol1", description: "test sub pol",
'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
subscriberCount: 12, rateLimitCount: 10,rateLimitTimeUnit: "sec", planId: "123456"};

BlockingCondition  denyPolicy = {policyId: "123456",conditionType: "APPLICATION",
conditionValue: "admin:MyApp5",conditionStatus: true};

@test:Config {}
function addApplicationUsagePlanTest() {
    ApplicationRatePlan payload = {
        "planName": "25PerMin",
        "displayName": "25PerMin",
        "description": "25 Per Min",
        "defaultLimit": {
            "type": "REQUESTCOUNTLIMIT",
            "requestCount": {
            "requestCount": 25,
            "timeUnit": "min",
            "unitTime": 1
            }
        }
    };
    ApplicationRatePlan|commons:APKError createdAppPol = addApplicationUsagePlan(payload);
    if createdAppPol is ApplicationRatePlan {
        test:assertTrue(true,"Application usage plan added successfully");
        applicationUsagePlan = createdAppPol;
        log:printInfo(createdAppPol.toString());
    } else if createdAppPol is commons:APKError {
        log:printError(createdAppPol.toString());
        test:assertFail("Error occured while adding Application Usage Plan");
    }
}

@test:Config {dependsOn: [addApplicationUsagePlanTest]}
function getApplicationUsagePlanByIdTest(){
    string? planId = applicationUsagePlan.planId;
    if planId is string {
        ApplicationRatePlan|commons:APKError policy = getApplicationUsagePlanById(planId);
        if policy is ApplicationRatePlan {
            test:assertTrue(true, "Successfully retrieved Application Usage Plan");
            log:printInfo(policy.toString());
        } else {
            test:assertFail("Error occured while retrieving Application Usage Plan");
        }
    } else {
        test:assertFail("Plan ID isn't a string");
    }
}

@test:Config {dependsOn: [getApplicationUsagePlanByIdTest]}
function getApplicationUsagePlansTest(){
    ApplicationRatePlanList|commons:APKError appPolicyList = getApplicationUsagePlans();
    if appPolicyList is ApplicationRatePlanList {
    test:assertTrue(true, "Successfully retrieved all Application Usage Plans");
    log:printInfo(appPolicyList.toString());
    } else if appPolicyList is commons:APKError {
        test:assertFail("Error occured while retrieving all Application Usage Plans");
    }
}

@test:Config {dependsOn: [getApplicationUsagePlansTest]}
function updateApplicationUsagePlanTest() {
    ApplicationRatePlan payload = {
        "planName": "25PerMin",
        "displayName": "25PerMin",
        "description": "25 Per Min Updated",
        "defaultLimit": {
            "type": "REQUESTCOUNTLIMIT",
            "requestCount": {
            "requestCount": 26,
            "timeUnit": "min",
            "unitTime": 1
            }
        }
    };
    string? planId = applicationUsagePlan.planId;
    if planId is string {
        ApplicationRatePlan|commons:APKError createdAppPol = updateApplicationUsagePlan(planId,payload);
        if createdAppPol is ApplicationRatePlan {
            test:assertTrue(true,"Application usage plan updated successfully");
        } else if createdAppPol is commons:APKError {
            log:printError(createdAppPol.toString());
            test:assertFail("Error occured while updating Application Usage Plan");
        }
    } else {
        test:assertFail("Plan ID isn't a string");
    }
}

@test:Config {dependsOn: [updateApplicationUsagePlanTest]}
function removeApplicationUsagePlanTest(){
    string? planId = applicationUsagePlan.planId;
    if planId is string {
        error?|string status = removeApplicationUsagePlan(planId);
        if status is string {
        test:assertTrue(true, "Successfully deleted Application Usage Plan");
        } else if status is  error {
            test:assertFail("Error occured while deleting Application Usage Plan");
        }
    } else {
        test:assertFail("Plan ID isn't a string");
    }
}

@test:Config
function addBusinessPlanTest() {
    BusinessPlan payload = {
        "planName": "BusinessPlan2",
        "displayName": "BusinessPlan2",
        "description": "test sub pol test2",
        "defaultLimit": {
            "type": "REQUESTCOUNTLIMIT",
            "requestCount": {
            "requestCount": 20,
            "timeUnit": "min",
            "unitTime": 1
            }
        },
        "rateLimitCount": 10,
        "rateLimitTimeUnit": "sec",
        "customAttributes": []
    };
    BusinessPlan|commons:APKError createdBusinessPlan = addBusinessPlan(payload);
    if createdBusinessPlan is BusinessPlan {
        test:assertTrue(true,"Business Plan added successfully");
        businessPlan.planId = createdBusinessPlan.planId;
    } else if createdBusinessPlan is commons:APKError {
        test:assertFail("Error occured while adding Business Plan");
    }
}

@test:Config {dependsOn: [addBusinessPlanTest]}
function getBusinessPlanByIdTest() {
    string? planId = businessPlan.planId;
    if planId is string {
        BusinessPlan|commons:APKError businessPlanResponse = getBusinessPlanById(planId);
        if businessPlanResponse is BusinessPlan {
            test:assertTrue(true,"Successfully retrieved Business Plan");
        } else {
            test:assertFail("Error occured while retrieving Business Plan");
        }
    } else {
        test:assertFail("Plan ID isn't a string");
    }
}

@test:Config {dependsOn: [getBusinessPlanByIdTest]}
function getBusinessPlansTest() {
    BusinessPlanList|commons:APKError businessPlansResponse = getBusinessPlans();
    if businessPlansResponse is BusinessPlanList {
        test:assertTrue(true,"Successfully retrieved all Business Plans");
    } else if businessPlansResponse is commons:APKError {
        test:assertFail("Error occured while retrieving all Business Plans");
    }
}

@test:Config {dependsOn: [getBusinessPlansTest]}
function updateBusinessPlanTest() {
    BusinessPlan payload = {
        "planName": "BusinessPlan2",
        "displayName": "BusinessPlan2",
        "description": "test sub pol test2 updated",
        "defaultLimit": {
            "type": "REQUESTCOUNTLIMIT",
            "requestCount": {
            "requestCount": 25,
            "timeUnit": "min",
            "unitTime": 1
            }
        },
        "rateLimitCount": 10,
        "rateLimitTimeUnit": "sec",
        "customAttributes": []
    };
    string? planId = businessPlan.planId;
    if planId is string {
        BusinessPlan|commons:APKError updatedBusinessPlan = updateBusinessPlan(planId,payload);
        if updatedBusinessPlan is BusinessPlan {
            test:assertTrue(true,"Business Plan updated successfully");
        } else if updatedBusinessPlan is commons:APKError {
            test:assertFail("Error occured while updating Business Plan");
        }
    } else {
        test:assertFail("Plan ID isn't a string");
    }
}

@test:Config {dependsOn: [updateBusinessPlanTest]}
function removeBusinessPlanTest(){
    string? planId = businessPlan.planId;
    if planId is string {
        commons:APKError|string status = removeBusinessPlan(planId);
        if status is string {
            test:assertTrue(true, "Successfully deleted Business Plan");
        } else if status is  commons:APKError {
            test:assertFail("Error occured while deleting Business Plan");
        }
    } else {
        test:assertFail("Plan ID isn't a string");
    }
}

@test:Config {}
function addDenyPolicyTest() {
    BlockingCondition payload = {
        "conditionType": "APPLICATION",
        "conditionValue": "admin:MyApp6",
        "conditionStatus": true
    };
    BlockingCondition|commons:APKError createdDenyPolicy = addDenyPolicy(payload);
    if createdDenyPolicy is BlockingCondition {
        test:assertTrue(true,"Deny Policy added successfully");
        denyPolicy.policyId = createdDenyPolicy.policyId;
    } else if createdDenyPolicy is commons:APKError {
        test:assertFail("Error occured while adding Deny Policy");
    }
}

@test:Config {dependsOn: [addDenyPolicyTest]}
function getDenyPolicyByIdTest() {
    string? policyId = denyPolicy.policyId;
    if policyId is string {
        BlockingCondition|commons:APKError denyPolicyResponse = getDenyPolicyById(policyId);
        if denyPolicyResponse is BlockingCondition {
            test:assertTrue(true,"Successfully retrieved Deny Policy");
        } else  {
            test:assertFail("Error occured while retrieving Deny Policy");
        }
    } else {
        test:assertFail("Policy ID isn't a string");
    }
}

@test:Config {dependsOn: [getDenyPolicyByIdTest]}
function getAllDenyPoliciesTest() {
    BlockingConditionList|commons:APKError denyPoliciesResponse = getAllDenyPolicies();
    if denyPoliciesResponse is BlockingConditionList {
        test:assertTrue(true,"Successfully retrieved all Deny Policy");
    } else if denyPoliciesResponse is commons:APKError {
        test:assertFail("Error occured while retrieving all Deny Policy");
    }
}

@test:Config {dependsOn: [getAllDenyPoliciesTest]}
function updateDenyPolicyTest() {
    string? policyId = denyPolicy.policyId;
    if policyId is string {
        BlockingConditionStatus status = {conditionStatus: false, policyId: policyId};
        string?|BlockingCondition|commons:APKError denyPolicyResponse = updateDenyPolicy(policyId, status);
        if denyPolicyResponse is BlockingCondition {
            test:assertTrue(true,"Successfully updated Deny Policy Status");
        } else if denyPolicyResponse is commons:APKError {
            test:assertFail("Error occured while updating Deny Policy Status");
        }
    } else {
        test:assertFail("Policy ID isn't a string");
    }
}

@test:Config {dependsOn: [updateDenyPolicyTest]}
function removeDenyPolicyTest(){
    string? policyId = denyPolicy.policyId;
    if policyId is string {
        commons:APKError|string status = removeDenyPolicy(policyId);
        if status is string {
        test:assertTrue(true, "Successfully deleted Deny Policy");
        } else if status is  commons:APKError {
            test:assertFail("Error occured while deleting Deny Policy");
        }
    } else {
        test:assertFail("Policy ID isn't a string");
    }
}
