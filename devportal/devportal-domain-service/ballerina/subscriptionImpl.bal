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
import ballerina/lang.value;

isolated function addSubscription(Subscription payload, string org, string user) returns string?|Subscription|error {
    int apiId = 0;
    int appId = 0;
    int|error? subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId !is int {
        string err = "Error while retrieving user by name " + user;
        log:printError(err);
        return error(err);
    } 
    string? apiUUID = payload.apiId;
    if apiUUID is string {
        string?|API|error api = getAPIByAPIId(apiUUID, org);
        if api !is API {
            string err = "Error while retrieving API by provided id" + apiUUID;
            log:printError(err);
            return error(err);
        }
        string apiInString = api.toJsonString();
        json j = check value:fromJsonString(apiInString);
        apiId = check j.apiId.ensureType();
    }
    string? appUUID = payload.applicationId;
    if appUUID is string {
        string?|Application|error application = getApplicationById(appUUID, org);
        if application !is Application {
            string err = "Error while retrieving Application by provided id" + appUUID;
            log:printError(err);
            return error(err);
        }
        string appInString = application.toJsonString();
        json j = check value:fromJsonString(appInString);
        appId = check j.id.ensureType();
    }
    string? businessPlan = payload.throttlingPolicy;
    if businessPlan is string {
        string?|error businessPlanID = getBusinessPlanByName(businessPlan);
        if businessPlanID !is string {
            string err = "Error while retrieving BusinessPlan by provided name" + businessPlan;
            log:printError(err);
            return error(err);
        }
        payload.requestedThrottlingPolicy = businessPlan;
    }
    string subscriptionId = uuid:createType1AsString();
    payload.subscriptionId = subscriptionId;
    payload.status = "UNBLOCKED";
    string?|Subscription|error createdSub = addSubscriptionDAO(payload,user,apiId,appId);
    return createdSub;
}

isolated function getBusinessPlanByName(string policyName) returns string?|error {
    string?|error policy = getBusinessPlanByNameDAO(policyName);
    return policy;
}

isolated function addMultipleSubscriptions(Subscription[] subscriptions, string org, string user) returns Subscription[]|error? {
    Subscription[]|error? addedSubs = [];
    foreach Subscription sub in subscriptions {
        string?|Subscription|error subscriptionResponse = check addSubscription(sub, org, user);
        if subscriptionResponse is Subscription {
            if addedSubs is Subscription[] {
                addedSubs.push(subscriptionResponse);
            }
        } else if subscriptionResponse is error {
            return subscriptionResponse;
        }
    }
    return addedSubs;
}

isolated function getSubscriptionById(string subId, string org) returns string?|Subscription|error {
    string?|Subscription|error subscription = getSubscriptionByIdDAO(subId, org);
    return subscription;
}

isolated function deleteSubscription(string subId, string organization) returns string|error? {
    error?|string status = deleteSubscriptionDAO(subId,organization);
    return status;
}

isolated function updateSubscription(string subId, Subscription payload, string org, string user) returns string?|Subscription|NotFoundError|error {
    string?|Subscription|error existingSub = getSubscriptionByIdDAO(subId, org);
    if existingSub is Subscription {
        payload.subscriptionId = subId;
    } else {
        Error err = {code:9010101, message:"Subscription Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    }
    int apiId = 0;
    int appId = 0;
    int|error? subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId !is int {
        string err = "Error while retrieving user by name " + user;
        log:printError(err);
        return error(err);
    } 
    string? apiUUID = payload.apiId;
    if apiUUID is string {
        string?|API|error api = getAPIByAPIId(apiUUID, org);
        if api !is API {
            string err = "Error while retrieving API by provided id" + apiUUID;
            log:printError(err);
            return error(err);
        }
        string apiInString = api.toJsonString();
        json j = check value:fromJsonString(apiInString);
        apiId = check j.apiId.ensureType();
    }
    string? appUUID = payload.applicationId;
    if appUUID is string {
        string?|Application|error application = getApplicationById(appUUID, org);
        if application !is Application {
            string err = "Error while retrieving Application by provided id" + appUUID;
            log:printError(err);
            return error(err);
        }
        string appInString = application.toJsonString();
        json j = check value:fromJsonString(appInString);
        appId = check j.id.ensureType();
    }
    string? businessPlan = payload.throttlingPolicy;
    if businessPlan is string {
        string?|error businessPlanID = getBusinessPlanByName(businessPlan);
        if businessPlanID !is string {
            string err = "Error while retrieving BusinessPlan by provided name" + businessPlan;
            log:printError(err);
            return error(err);
        }
        payload.requestedThrottlingPolicy = businessPlan;
    }
    payload.status = "UNBLOCKED";
    string?|Subscription|error createdSub = updateSubscriptionDAO(payload,user,apiId,appId);
    return createdSub;
}

isolated function getSubscriptions(string? apiId, string? applicationId, string? groupId, int offset, int limitCount, string org) returns string?|SubscriptionList|error {
    if apiId is string && applicationId is string {
        // Retrieve Subscriptions per given API Id and App Id
        string?|Subscription|error subscription = getSubscriptionByAPIandAppIdDAO(apiId,applicationId,org);
        if subscription is Subscription {
            Subscription[] subs = [subscription];
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subscription;
        }
    } else if apiId is string {
        // Retrieve Subscriptions per given API Id
        Subscription[]|error? subs = getSubscriptionsByAPIIdDAO(apiId,org);
        if subs is Subscription[] {
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subs;
        }
    } else if applicationId is string {
        // Retrieve Subscriptions per given APP Id
        Subscription[]|error? subs = getSubscriptionsByAPPIdDAO(applicationId,org);
        if subs is Subscription[] {
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subs;
        }
    } else {
        // Retrieve All Subscriptions
        Subscription[]|error? subs = getSubscriptionsList(org);
        if subs is Subscription[] {
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subs;
        }
    }
}