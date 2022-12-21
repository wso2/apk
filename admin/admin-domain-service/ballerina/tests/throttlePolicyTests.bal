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

@test:Mock { functionName: "addApplicationUsagePlanDAO" }
test:MockFunction addApplicationUsagePlanDAOMock = new();

@test:Mock { functionName: "getApplicationUsagePlanByIdDAO" }
test:MockFunction getApplicationUsagePlanByIdDAOMock = new();

@test:Mock { functionName: "getApplicationUsagePlansDAO" }
test:MockFunction getApplicationUsagePlansDAOMock = new();

@test:Mock { functionName: "updateApplicationUsagePlanDAO" }
test:MockFunction updateApplicationUsagePlanDAOMock = new();

@test:Mock { functionName: "deleteApplicationUsagePlanDAO" }
test:MockFunction deleteApplicationUsagePlanDAOMock = new();

@test:Mock { functionName: "addBusinessPlanDAO" }
test:MockFunction addBusinessPlanDAOMock = new();

@test:Mock { functionName: "getBusinessPlanByIdDAO" }
test:MockFunction getBusinessPlanByIdDAOMock = new();

@test:Mock { functionName: "getBusinessPlansDAO" }
test:MockFunction getBusinessPlansDAOMock = new();

@test:Mock { functionName: "updateBusinessPlanDAO" }
test:MockFunction updateBusinessPlanDAOMock = new();

@test:Mock { functionName: "deleteBusinessPlanDAO" }
test:MockFunction deleteBusinessPlanDAOMock = new();

@test:Mock { functionName: "addDenyPolicyDAO" }
test:MockFunction addDenyPolicyDAOMock = new();

@test:Mock { functionName: "getDenyPolicyByIdDAO" }
test:MockFunction getDenyPolicyByIdDAOMock = new();

@test:Mock { functionName: "getDenyPoliciesDAO" }
test:MockFunction getDenyPoliciesDAOMock = new();

@test:Mock { functionName: "updateDenyPolicyDAO" }
test:MockFunction updateDenyPolicyDAOMock = new();

@test:Mock { functionName: "deleteDenyPolicyDAO" }
test:MockFunction deleteDenyPolicyDAOMock = new();

@test:Config {}
function addApplicationUsagePlanTest() {
    string?|ApplicationRatePlan|error  applicationUsagePlan = {planName: "15PerMin", description: "15 Per Minute",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}};
    ApplicationRatePlan payload = {planName: "15PerMin", description: "15 Per Minute",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"},planId: "1234455"};
    test:when(addApplicationUsagePlanDAOMock).thenReturn(applicationUsagePlan);
    string?|ApplicationRatePlan|error createdAppPol = addApplicationUsagePlan(payload);
    if createdAppPol is ApplicationRatePlan {
        test:assertTrue(true,"Application usage plan added successfully");
    } else if createdAppPol is error {
        test:assertFail("Error occured while adding Application Usage Plan");
    }
}

@test:Config {}
function getApplicationUsagePlanByIdTest(){
    ApplicationRatePlan aup = {planName: "15PerMin", description: "15 Per Minute",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"},planId: "1234455"};
    test:when(getApplicationUsagePlanByIdDAOMock).withArguments("1234455").thenReturn(aup);
    string?|ApplicationRatePlan|error policy = getApplicationUsagePlanByIdDAO("1234455");
    if policy is ApplicationRatePlan {
        test:assertTrue(true, "Successfully retrieved Application Usage Plan");
    } else if policy is  error {
        test:assertFail("Error occured while retrieving Application Usage Plan");
    }
}

@test:Config {}
function getApplicationUsagePlansTest(){
    ApplicationRatePlan[] aupList = [{planName: "15PerMin", description: "15 Per Minute",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"},planId: "1234455"},
    {planName: "30PerMin", description: "30 Per Minute",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"},planId: "asqw1234dsd"}];
    test:when(getApplicationUsagePlansDAOMock).withArguments("carbon.super").thenReturn(aupList);
    string?|ApplicationRatePlanList|error appPolicyList = getApplicationUsagePlans();
    if appPolicyList is ApplicationRatePlanList {
    test:assertTrue(true, "Successfully retrieved all Application Usage Plans");
    } else if appPolicyList is  error {
        test:assertFail("Error occured while retrieving all Application Usage Plans");
    }
}

@test:Config {}
function updateApplicationUsagePlanTest() {
    string?|ApplicationRatePlan|error  applicationUsagePlan = {planName: "15PerMin", description: "15 Per Minute",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"},planId: "1234455"};
    ApplicationRatePlan payload = {planName: "15PerMin", description: "15 Per Minute New",'type:"ApplicationThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"},planId: "1234455"};
    test:when(getApplicationUsagePlanByIdDAOMock).withArguments("1234455").thenReturn(applicationUsagePlan);
    test:when(updateApplicationUsagePlanDAOMock).thenReturn(payload);
    string?|ApplicationRatePlan|error createdAppPol = updateApplicationUsagePlanDAO(payload);
    if createdAppPol is ApplicationRatePlan {
        test:assertTrue(true,"Application usage plan updated successfully");
    } else if createdAppPol is error {
        test:assertFail("Error occured while updating Application Usage Plan");
    }
}

@test:Config {}
function removeApplicationUsagePlanTest(){
    test:when(deleteApplicationUsagePlanDAOMock).withArguments("1234455").thenReturn("");
    error?|string status = removeApplicationUsagePlan("1234455");
    if status is string {
    test:assertTrue(true, "Successfully deleted Application Usage Plan");
    } else if status is  error {
        test:assertFail("Error occured while deleting Application Usage Plan");
    }
}

@test:Config {}
function addBusinessPlanTest() {
    string?|BusinessPlan|error  businessPlan = {planName: "MySubPol1", description: "test sub pol",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 12, rateLimitCount: 10,rateLimitTimeUnit: "sec"};
    BusinessPlan payload = {planName: "MySubPol1", description: "test sub pol",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 12, rateLimitCount: 10,rateLimitTimeUnit: "sec", planId: "123456"};
    test:when(addBusinessPlanDAOMock).thenReturn(businessPlan);
    string?|BusinessPlan|error createdBusinessPlan = addBusinessPlan(payload);
    if createdBusinessPlan is BusinessPlan {
        test:assertTrue(true,"Business Plan added successfully");
    } else if createdBusinessPlan is error {
        test:assertFail("Error occured while adding Business Plan");
    }
}

@test:Config {}
function getBusinessPlanByIdTest() {
    string?|BusinessPlan|error  businessPlan = {planName: "MySubPol1", description: "test sub pol",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 12, rateLimitCount: 10,rateLimitTimeUnit: "sec", planId: "123456"};
    test:when(getBusinessPlanByIdDAOMock).thenReturn(businessPlan);
    string?|BusinessPlan|error businessPlanResponse = getBusinessPlanById("123456");
    if businessPlanResponse is BusinessPlan {
        test:assertTrue(true,"Successfully retrieved Business Plan");
    } else if businessPlanResponse is error {
        test:assertFail("Error occured while retrieving Business Plan");
    }
}

@test:Config {}
function getBusinessPlansTest() {
    BusinessPlan[]  businessPlans = [{planName: "MySubPol1", description: "test sub pol",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 12, rateLimitCount: 10,rateLimitTimeUnit: "sec", planId: "123456"},{planName: "MySubPol2", description: "test sub pol 2",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 15, rateLimitCount: 14,rateLimitTimeUnit: "sec", planId: "qwe13123456asd"}];
    test:when(getBusinessPlansDAOMock).withArguments("carbon.super").thenReturn(businessPlans);
    string?|BusinessPlanList|error businessPlansResponse = getBusinessPlans();
    if businessPlansResponse is BusinessPlanList {
        test:assertTrue(true,"Successfully retrieved all Business Plans");
    } else if businessPlansResponse is error {
        test:assertFail("Error occured while retrieving all Business Plans");
    }
}

@test:Config {}
function updateBusinessPlanTest() {
    string?|BusinessPlan|error  businessPlan = {planName: "MySubPol1", description: "test sub pol",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 12, rateLimitCount: 10,rateLimitTimeUnit: "sec"};
    BusinessPlan payload = {planName: "MySubPol1", description: "test sub pol New",'type:"SubscriptionThrottlePolicy",defaultLimit: {'type: "REQUESTCOUNTLIMIT"}, 
    subscriberCount: 12, rateLimitCount: 5,rateLimitTimeUnit: "sec", planId: "123456"};
     test:when(getBusinessPlanByIdDAOMock).thenReturn(businessPlan);
    test:when(updateBusinessPlanDAOMock).thenReturn(payload);
    string?|BusinessPlan|NotFoundError|error updatedBusinessPlan = updateBusinessPlan("123456",payload);
    if updatedBusinessPlan is BusinessPlan {
        test:assertTrue(true,"Business Plan updated successfully");
    } else if updatedBusinessPlan is error {
        test:assertFail("Error occured while updating Business Plan");
    }
}

@test:Config {}
function removeBusinessPlanTest(){
    test:when(deleteBusinessPlanDAOMock).withArguments("123456").thenReturn("");
    error?|string status = removeBusinessPlan("123456");
    if status is string {
    test:assertTrue(true, "Successfully deleted Business Plan");
    } else if status is  error {
        test:assertFail("Error occured while deleting Business Plan");
    }
}

@test:Config {}
function addDenyPolicyTest() {
    string?|BlockingCondition|error  denyPolicy = {policyId: "123456",conditionType: "APPLICATION",conditionValue: "admin:MyApp5",conditionStatus: true};
    BlockingCondition payload = {policyId: "123456",conditionType: "APPLICATION",conditionValue: "admin:MyApp5",conditionStatus: true};
    test:when(addDenyPolicyDAOMock).thenReturn(denyPolicy);
    string?|BlockingCondition|error createdDenyPolicy = addDenyPolicy(payload);
    if createdDenyPolicy is BlockingCondition {
        test:assertTrue(true,"Deny Policy added successfully");
    } else if createdDenyPolicy is error {
        test:assertFail("Error occured while adding Deny Policy");
    }
}

@test:Config {}
function getDenyPolicyByIdTest() {
    BlockingCondition denyPolicy = {policyId: "123456",conditionType: "APPLICATION",conditionValue: "admin:MyApp5",conditionStatus: true};
    test:when(getDenyPolicyByIdDAOMock).withArguments("123456").thenReturn(denyPolicy);
    string?|BlockingCondition|error denyPolicyResponse = getDenyPolicyById("123456");
    if denyPolicyResponse is BlockingCondition {
        test:assertTrue(true,"Successfully retrieved Deny Policy");
    } else if denyPolicyResponse is error {
        test:assertFail("Error occured while retrieving Deny Policy");
    }
}

@test:Config {}
function getAllDenyPoliciesTest() {
    BlockingCondition[] denyPolicies = [{policyId: "123456",conditionType: "APPLICATION",conditionValue: "admin:MyApp5",conditionStatus: true},
    {policyId: "asqe123456ad",conditionType: "APPLICATION",conditionValue: "admin:MyApp5",conditionStatus: true}];
    test:when(getDenyPoliciesDAOMock).withArguments("carbon.super").thenReturn(denyPolicies);
    string?|BlockingConditionList|error denyPoliciesResponse = getAllDenyPolicies();
    if denyPoliciesResponse is BlockingConditionList {
        test:assertTrue(true,"Successfully retrieved all Deny Policy");
    } else if denyPoliciesResponse is error {
        test:assertFail("Error occured while retrieving all Deny Policy");
    }
}

@test:Config {}
function updateDenyPolicyTest() {
    BlockingCondition denyPolicy = {policyId: "123456",conditionType: "APPLICATION",conditionValue: "admin:MyApp5",conditionStatus: true};
    BlockingConditionStatus status = {conditionStatus: false, policyId: "123456"};
    test:when(getDenyPolicyByIdDAOMock).withArguments("123456").thenReturn(denyPolicy);
    test:when(updateDenyPolicyDAOMock).withArguments(status).thenReturn("");
    string?|BlockingCondition|NotFoundError|error denyPolicyResponse = updateDenyPolicy("123456", status);
    if denyPolicyResponse is BlockingCondition {
        test:assertTrue(true,"Successfully updated Deny Policy Status");
    } else if denyPolicyResponse is error {
        test:assertFail("Error occured while updating Deny Policy Status");
    }
}

@test:Config {}
function removeDenyPolicyTest(){
    test:when(deleteDenyPolicyDAOMock).withArguments("123456").thenReturn("");
    error?|string status = removeDenyPolicy("123456");
    if status is string {
    test:assertTrue(true, "Successfully deleted Deny Policy");
    } else if status is  error {
        test:assertFail("Error occured while deleting Deny Policy");
    }
}
