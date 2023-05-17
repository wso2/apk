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
import ballerina/uuid;
import wso2/apk_common_lib as commons;

Subscription sub = {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: "21212"};
Application applicationNew = {name: "sampleAppNew", description: "sample application"};

@test:Mock {functionName: "retrieveManagementServerHostsList"}
test:MockFunction retrieveManagementServerHostsListMock = new ();

@test:Mock {functionName: "createApplication", moduleName: "wso2/notification_grpc_client"}
public isolated function createApplicationMock(ApplicationGRPC createApplicationRequest, string endpoint, string pubCert, string devCert, string devKey) returns error|NotificationResponse {
    NotificationResponse noti = {code: "OK"};
    return noti;
}

@test:Mock {functionName: "createSubscription", moduleName: "wso2/notification_grpc_client"}
public isolated function createSubscriptionMock(SubscriptionGRPC createSubscriptionRequest, string endpoint, string pubCert, string devCert, string devKey) returns error|NotificationResponse {
    NotificationResponse noti = {code: "OK"};
    return noti;
}

@test:Mock {functionName: "updateSubscription", moduleName: "wso2/notification_grpc_client"}
public isolated function updateSubscriptionMock(SubscriptionGRPC updateSubscriptionRequest, string endpoint, string pubCert, string devCert, string devKey) returns error|NotificationResponse {
    NotificationResponse noti = {code: "OK"};
    return noti;
}

@test:Mock {functionName: "deleteSubscription", moduleName: "wso2/notification_grpc_client"}
public isolated function deleteSubscriptionMock(SubscriptionGRPC deleteSubscriptionRequest, string endpoint, string pubCert, string devCert, string devKey) returns error|NotificationResponse {
    NotificationResponse noti = {code: "OK"};
    return noti;
}

@test:BeforeSuite
function beforeFunc3() {
    string[] testHosts = ["http://localhost:9090"];
    test:when(retrieveManagementServerHostsListMock).thenReturn(testHosts);
    Application payload = {name: "sampleAppNew", description: "sample application"};
    NotFoundError|Application|commons:APKError createdApplication = addApplication(payload, organiztion, "apkuser");
    if createdApplication is Application {
        test:assertTrue(true, "Successfully added the application");
        applicationNew.applicationId = createdApplication.applicationId;
        BusinessPlan payloadbp = {
            "planName": "MyBusinessPlan3",
            "displayName": "MyBusinessPlan3",
            "description": "test sub pol test",
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
        payloadbp.planId = uuid:createType1AsString();
        BusinessPlan|commons:APKError createdBusinessPlan = addBusinessPlanDAO(payloadbp, organiztion.uuid);
        if createdBusinessPlan is commons:APKError {
            test:assertFail("Error occured while adding Business Plan");
        }
    } else if createdApplication is error {
        test:assertFail("Error occured while adding application");
    }
}

@test:Config {}
function addSubscriptionTest() {
    string? appId = applicationNew.applicationId;
    if appId is string {
        Subscription payload = {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: appId};
        Subscription|commons:APKError|NotFoundError subscription = addSubscription(payload, organiztion, "apkuser");
        if subscription is Subscription {
            test:assertTrue(true, "Succesfully added a subscription");
            sub.subscriptionId = subscription.subscriptionId;
        } else if subscription is commons:APKError {
            log:printError(subscription.toString());
            test:assertFail("Error occured while adding subscription");
        } else if subscription is NotFoundError {
            log:printError(subscription.toString());
            test:assertFail("Error occured while adding subscription");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [addSubscriptionTest]}
function addSubscriptionNegativeTest1() {
    // API ID is not found or API Id is not returned
    string? appId = applicationNew.applicationId;
    if appId is string {
        Subscription payload = {apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0", applicationId: appId};
        Subscription|commons:APKError|NotFoundError subscription = addSubscription(payload, organiztion, "apkuser");
        if subscription is Subscription {
            test:assertFail("Succesfully added a subscription for a invalid API");
        } else {
            test:assertTrue(true, "Sucessfully validated API not available while adding a subscription");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [addSubscriptionNegativeTest1]}
function addSubscriptionNegativeTest2() {
    // APP ID is not found or APP Id is not returned
    Subscription payload = {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4"};
    Subscription|commons:APKError|NotFoundError subscription = addSubscription(payload, organiztion, "apkuser");
    if subscription is Subscription {
        test:assertFail("Succesfully added a subscription for a invalid Application");
    } else {
        test:assertTrue(true, "Sucessfully validated Application not available while adding a subscription");
    }
}

@test:Config {dependsOn: [addSubscriptionNegativeTest2]}
function addSubscriptionNegativeTest3() {
    // Policy Not Found
    string? appId = applicationNew.applicationId;
    if appId is string {
        Subscription payload = {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: appId};
        Subscription|commons:APKError|NotFoundError subscription = addSubscription(payload, organiztion, "apkuser");
        if subscription is Subscription {
            test:assertFail("Succesfully added a subscription for a invalid Policy");
        } else {
            test:assertTrue(true, "Sucessfully validated Policy not available while adding a subscription");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [addSubscriptionNegativeTest3]}
function addMultipleSubscriptionsTest() {
    // Add 2 new app
    string? newappId1 = "";
    string? newappId2 = "";
    Application payload = {name: "sampleAppNew1", description: "sample application"};
    NotFoundError|Application|commons:APKError createdApplication = addApplication(payload, organiztion, "apkuser");
    if createdApplication is Application {
        test:assertTrue(true, "Successfully added the application");
        newappId1 = createdApplication.applicationId;
    } else if createdApplication is error {
        test:assertFail("Error occured while adding application");
    }
    Application payload2 = {name: "sampleAppNew2", description: "sample application"};
    NotFoundError|Application|commons:APKError createdApplication2 = addApplication(payload2, organiztion, "apkuser");
    if createdApplication2 is Application {
        test:assertTrue(true, "Successfully added the application");
        newappId2 = createdApplication2.applicationId;
    } else if createdApplication2 is error {
        test:assertFail("Error occured while adding application");
    }

    if newappId1 is string && newappId2 is string {
        Subscription[] multiSub = [
            {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: newappId1},
            {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: newappId2}
        ];

        Subscription[]|commons:APKError|NotFoundError subscriptions = addMultipleSubscriptions(multiSub, organiztion, "apkuser");
        if subscriptions is Subscription[] {
            test:assertTrue(true, "Succesfully added multiple subscriptions");
        } else {
            test:assertFail("Error occured while adding multiple subscriptions");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [addMultipleSubscriptionsTest]}
function getSubscriptionByIdTest() {
    string? subId = sub.subscriptionId;
    if subId is string {
        Subscription|commons:APKError|NotFoundError returnedResponse = getSubscriptionById(subId, organiztion);
        if returnedResponse is Subscription {
            test:assertTrue(true, "Successfully retrieved subscription");
        } else {
            test:assertFail("Error occured while retrieving subscription");
        }
    } else {
        test:assertFail("Sub ID isn't a string");
    }
}

@test:Config {dependsOn: [getSubscriptionByIdTest]}
function updateSubscriptionTest() {
    // add a new policy
    BusinessPlan payloadbp = {
        "planName": "MyBusinessPlan2",
        "displayName": "MyBusinessPlan2",
        "description": "test sub pol test",
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
    payloadbp.planId = uuid:createType1AsString();
    BusinessPlan|commons:APKError createdBusinessPlan = addBusinessPlanDAO(payloadbp, organiztion.uuid);
    if createdBusinessPlan is BusinessPlan {
        test:assertTrue(true, "Business Plan added successfully");
        string? appId = applicationNew.applicationId;
        string? subId = sub.subscriptionId;
        if appId is string && subId is string {
            // Use new policy
            Subscription payload = {apiId: "01ed75e2-b30b-18c8-wwf2-25da7edd2231", applicationId: appId};
            string?|Subscription|NotFoundError|error subscription = updateSubscription(subId, payload, organiztion, "apkuser");
            if subscription is Subscription {
                test:assertTrue(true, "Succesfully updated the subscription");
            } else if subscription is error {
                test:assertFail("Error occured while updating subscription");
            }
        } else {
            test:assertFail("App ID isn't a string");
        }
    } else if createdBusinessPlan is commons:APKError {
        test:assertFail("Error occured while adding Business Plan");
    }
}

@test:Config {dependsOn: [updateSubscriptionTest]}
function getSubscriptionListTest1() {
    // Providing both API ID and App Id
    string? appId = applicationNew.applicationId;
    if appId is string {
        SubscriptionList|commons:APKError|NotFoundError subscriptionList = getSubscriptions("01ed75e2-b30b-18c8-wwf2-25da7edd2231", appId, "", 0, 0, organiztion);
        if subscriptionList is ApplicationList {
            test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
        } else {
            test:assertFail("Error occured while retrieving all subscriptions");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [getSubscriptionListTest1]}
function getSubscriptionListTest2() {
    // Providing only API ID
    SubscriptionList|commons:APKError|NotFoundError subscriptionList = getSubscriptions("01ed75e2-b30b-18c8-wwf2-25da7edd2231", null, "", 0, 0, organiztion);
    if subscriptionList is ApplicationList {
        test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
    } else {
        test:assertFail("Error occured while retrieving all subscriptions");
    }
}

@test:Config {dependsOn: [getSubscriptionListTest2]}
function getSubscriptionListTest3() {
    // Providing only App ID
    string? appId = applicationNew.applicationId;
    if appId is string {
        SubscriptionList|commons:APKError|NotFoundError subscriptionList = getSubscriptions(null, appId, "", 0, 0, organiztion);
        if subscriptionList is ApplicationList {
            test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
        } else {
            test:assertFail("Error occured while retrieving all subscriptions");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [getSubscriptionListTest3]}
function getSubscriptionListTest4() {
    // Providing nothing and retrieving all subscriptions
    SubscriptionList|commons:APKError|NotFoundError subscriptionList = getSubscriptions(null, null, "", 0, 0, organiztion);
    if subscriptionList is ApplicationList {
        test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
    } else {
        test:assertFail("Error occured while retrieving all subscriptions");
    }
}

@test:Config {dependsOn: [getSubscriptionListTest4]}
function deleteSubscriptionTest() {
    string? subId = sub.subscriptionId;
    if subId is string {
        commons:APKError? status = deleteSubscription(subId, organiztion);
        if status is () {
            test:assertTrue(true, "Successfully deleted subscription");
        } else {
            test:assertFail("Error occured while deleting subscription");
        }
    } else {
        test:assertFail("Sub ID isn't a string");
    }
}
