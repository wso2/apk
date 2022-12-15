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

@test:Mock { functionName: "getBusinessPlanByNameDAO" }
test:MockFunction getBusinessPlanByNameDAOMock = new();

@test:Mock { functionName: "addSubscriptionDAO" }
test:MockFunction addSubscriptionDAOMock = new();

@test:Mock { functionName: "getSubscriptionByIdDAO" }
test:MockFunction getSubscriptionByIdDAOMock = new();

@test:Mock { functionName: "deleteSubscriptionDAO" }
test:MockFunction deleteSubscriptionDAOMock = new();

@test:Mock { functionName: "updateSubscriptionDAO" }
test:MockFunction updateSubscriptionDAOMock = new();

@test:Mock { functionName: "getSubscriptionByAPIandAppIdDAO" }
test:MockFunction getSubscriptionByAPIandAppIdDAOMock = new();

@test:Mock { functionName: "getSubscriptionsByAPIIdDAO" }
test:MockFunction getSubscriptionsByAPIIdDAOMock = new();

@test:Mock { functionName: "getSubscriptionsByAPPIdDAO" }
test:MockFunction getSubscriptionsByAPPIdDAOMock = new();

@test:Mock { functionName: "getSubscriptionsList" }
test:MockFunction getSubscriptionsListMock = new();

@test:Config {}
function addSubscriptionTest() {
    test:when(getSubscriberIdDAOMock).thenReturn(1);
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED", id:"123456wew",apiId: 1};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    Application application = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq",id: 1};
    test:when(getApplicationByIdDAOMock).thenReturn(application);
    string businessPlanName = "MySubPol5";
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when( addSubscriptionDAOMock).thenReturn(payload);
    string?|Subscription|error subscription = addSubscription(payload, "carbon.super", "apkuser");
    if subscription is Subscription {
        test:assertTrue(true, "Succesfully added a subscription");
    } else if subscription is error {
        log:printDebug(subscription.message());
        test:assertFail("Error occured while adding subscription");
    }
}

@test:Config {}
function addSubscriptionNegativeTest1() {
    test:when(getSubscriberIdDAOMock).thenReturn(1);
    // API ID is not found or API Id is not returned
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED"};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    Application application = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq",id: 1};
    test:when(getApplicationByIdDAOMock).thenReturn(application);
    string businessPlanName = "MySubPol5";
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when( addSubscriptionDAOMock).thenReturn(payload);
    string?|Subscription|error subscription = addSubscription(payload, "carbon.super", "apkuser");
    if subscription is Subscription {
        test:assertFail("Succesfully added a subscription for a invalid API");
    } else if subscription is error {
        test:assertTrue(true, "Sucessfully validated API not available while adding a subscription");
    }
}

@test:Config {}
function addSubscriptionNegativeTest2() {
    test:when(getSubscriberIdDAOMock).thenReturn(1);
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED", id:"123456wew",apiId: 1};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    // APP ID is not found or APP Id is not returned
    Application application = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application"};
    test:when(getApplicationByIdDAOMock).thenReturn(application);
    string businessPlanName = "MySubPol5";
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when( addSubscriptionDAOMock).thenReturn(payload);
    string?|Subscription|error subscription = addSubscription(payload, "carbon.super", "apkuser");
    if subscription is Subscription {
        test:assertFail("Succesfully added a subscription for a invalid Application");
    } else if subscription is error {
        test:assertTrue(true, "Sucessfully validated Application not available while adding a subscription");
    }
}

@test:Config {}
function addSubscriptionNegativeTest3() {
    test:when(getSubscriberIdDAOMock).thenReturn(1);
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED", id:"123456wew",apiId: 1};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    Application application = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq",id: 1};
    test:when(getApplicationByIdDAOMock).thenReturn(application);
    // Policy Not Found
    error businessPlanName = error("policy not found");
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when( addSubscriptionDAOMock).thenReturn(payload);
    string?|Subscription|error subscription = addSubscription(payload, "carbon.super", "apkuser");
    if subscription is Subscription {
        test:assertFail("Succesfully added a subscription for a invalid Policy");
    } else if subscription is error {
        test:assertTrue(true, "Sucessfully validated Policy not available while adding a subscription");
    }
}

@test:Config {}
function getBusinessPlanByNameTest(){
    string businessPlanName = "MySubPol5";
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    string?|error businessPlanID = getBusinessPlanByName(businessPlanName);
    if businessPlanID is string {
    test:assertTrue(true, "Successfully retrieved business plan id");
    } else if businessPlanID is  error {
        test:assertFail("Error occured while retrieving business plan id");
    }
}

@test:Config {}
function addMultipleSubscriptionsTest() {
    Subscription[] multiSub = [{ apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"},
    { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "dqwdqw2232qsdaaawsdqw",throttlingPolicy: "MySubPol4"}];
    test:when(getSubscriberIdDAOMock).thenReturn(1);
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED", id:"123456wew",apiId: 1};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    Application application = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq",id: 1};
    test:when(getApplicationByIdDAOMock).thenReturn(application);
    string businessPlanName = "MySubPol5";
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when( addSubscriptionDAOMock).thenReturn(payload);
    string?|Subscription[]|error subscriptions = addMultipleSubscriptions(multiSub, "carbon.super", "apkuser");
    if subscriptions is Subscription[] {
        test:assertTrue(true,"Succesfully added multiple subscriptions");
    } else if subscriptions is error {
        test:assertFail("Error occured while adding multiple subscriptions");
    }
}

@test:Config {}
function getSubscriptionByIdTest() {
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when(getSubscriptionByIdDAOMock).thenReturn(payload);
    string?|Subscription|error returnedResponse = getSubscriptionById("12sqwsq","carbon.super");
    if returnedResponse is Subscription {
    test:assertTrue(true, "Successfully retrieved subscription");
    } else if returnedResponse is  error {
        test:assertFail("Error occured while retrieving subscription");
    }
}

@test:Config {}
function deleteSubscriptionTest(){
    test:when(deleteSubscriptionDAOMock).withArguments("12sqwsq", "carbon.super").thenReturn("");
    error?|string status = deleteSubscription("12sqwsq", "carbon.super");
    if status is string {
    test:assertTrue(true, "Successfully deleted subscription");
    } else if status is  error {
        test:assertFail("Error occured while deleting subscription");
    }
}

@test:Config {}
function updateSubscriptionTest() {
    Subscription existingSub = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol5"};
    test:when(getSubscriptionByIdDAOMock).thenReturn(existingSub);
    test:when(getSubscriberIdDAOMock).thenReturn(1);
    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED", id:"123456wew",apiId: 1};
    test:when(getAPIByIdDAOMock).thenReturn(api);
    Application application = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq",id: 1};
    test:when(getApplicationByIdDAOMock).thenReturn(application);
    string businessPlanName = "MySubPol5";
    test:when(getBusinessPlanByNameDAOMock).thenReturn(businessPlanName);
    Subscription payload = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"};
    test:when(updateSubscriptionDAOMock).thenReturn(payload);
    string?|Subscription|NotFoundError|error subscription = updateSubscription("12sqwsq", payload, "carbon.super", "apkuser");
    if subscription is Subscription {
        test:assertTrue(true, "Succesfully updated the subscription");
    } else if subscription is error {
        log:printDebug(subscription.message());
        test:assertFail("Error occured while updating subscription");
    }
}

@test:Config {}
function getSubscriptionListTest1() {
    Subscription subscription = { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"};
    test:when(getSubscriptionByAPIandAppIdDAOMock).thenReturn(subscription);
    // Providing both API ID and App Id
    string?|SubscriptionList|error subscriptionList = getSubscriptions("8e3a1ca4-b649-4e57-9a57-e43b6b545af0","01ed716f-9f85-1ade-b634-be97dee7ceb4","",0,0,"carbon.super");
    if subscriptionList is ApplicationList {
    test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
    } else if subscriptionList is  error {
        test:assertFail("Error occured while retrieving all subscriptions");
    }
}

@test:Config {}
function getSubscriptionListTest2() {
    Subscription[] subList = [{ apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"},
    { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"}];
    test:when(getSubscriptionsByAPIIdDAOMock).thenReturn(subList);
    // Providing only API ID
    string?|SubscriptionList|error subscriptionList = getSubscriptions("8e3a1ca4-b649-4e57-9a57-e43b6b545af0",null,"",0,0,"carbon.super");
    if subscriptionList is ApplicationList {
    test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
    } else if subscriptionList is  error {
        test:assertFail("Error occured while retrieving all subscriptions");
    }
}

@test:Config {}
function getSubscriptionListTest3() {
    Subscription[] subList = [{ apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"},
    { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"}];
    test:when(getSubscriptionsByAPPIdDAOMock).thenReturn(subList);
    // Providing only App ID
    string?|SubscriptionList|error subscriptionList = getSubscriptions(null,"01ed716f-9f85-1ade-b634-be97dee7ceb4","",0,0,"carbon.super");
    if subscriptionList is ApplicationList {
    test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
    } else if subscriptionList is  error {
        test:assertFail("Error occured while retrieving all subscriptions");
    }
}

@test:Config {}
function getSubscriptionListTest4() {
    Subscription[] subList = [{ apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"},
    { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"}];
    test:when(getSubscriptionsListMock).thenReturn(subList);
    // Providing nothing and retrieving all subscriptions
    string?|SubscriptionList|error subscriptionList = getSubscriptions(null,null,"",0,0,"carbon.super");
    if subscriptionList is ApplicationList {
    test:assertTrue(true, "Successfully retrieved all subscriptions by API ID and App ID");
    } else if subscriptionList is  error {
        test:assertFail("Error occured while retrieving all subscriptions");
    }
}