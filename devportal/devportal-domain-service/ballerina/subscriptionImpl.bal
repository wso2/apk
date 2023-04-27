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
import ballerina/lang.value;
import wso2/notification_grpc_client;
import ballerina/time;
import ballerina/log;

isolated function addSubscription(Subscription payload, string org, string user) returns Subscription|APKError|NotFoundError|error {
    string apiId = "";
    string appId = "";
    string|NotFoundError|APKError subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId !is string {
        return subscriberId;
    } 
    string? apiUUID = payload.apiId;
    if apiUUID is string {
        API|NotFoundError|APKError api = getAPIByAPIId(apiUUID);
        if api is APKError|NotFoundError {
            return api;
        } else if api is API {
            string apiInString = api.toJsonString();
            json j = check value:fromJsonString(apiInString);
            apiId = check j.id.ensureType();
        }
    }
    string? appUUID = payload.applicationId;
    if appUUID is string {
        Application|APKError|NotFoundError application = getApplicationById(appUUID, org);
        if application is APKError|NotFoundError {
            return application;
        } else  if application is Application {
            string appInString = application.toJsonString();
            json j = check value:fromJsonString(appInString);
            appId = check j.applicationId.ensureType();
        }
    }
    string? businessPlan = payload.throttlingPolicy;
    if businessPlan is string {
        string|APKError|NotFoundError businessPlanID = getBusinessPlanByName(businessPlan);
        if businessPlanID is APKError|NotFoundError {
            return businessPlanID;
        }
        payload.requestedThrottlingPolicy = businessPlan;
    }
    string subscriptionId = uuid:createType1AsString();
    payload.subscriptionId = subscriptionId;
    payload.status = "UNBLOCKED";
    Subscription|APKError createdSub = addSubscriptionDAO(payload,user,apiId,appId);
    if createdSub is Subscription {
        string[]|APKError hostList = retrieveManagementServerHostsList();
        if hostList is string[] {
            string eventId = uuid:createType1AsString();
            time:Utc currTime = time:utcNow();
            string date = time:utcToString(currTime);
            SubscriptionGRPC createSubscriptionRequest = {eventId: eventId, applicationRef: createdSub.applicationId, 
            apiRef: <string>createdSub.apiId, policyId: createdSub.throttlingPolicy, subStatus:<string>createdSub.status,
            subscriber: user, uuid: subscriptionId, timeStamp: date, organization: org};
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = managementServerConfig.certPath;
            foreach string host in hostList {
                NotificationResponse|error subscriptionNotification = notification_grpc_client:createSubscription(createSubscriptionRequest,
                "https://" + host + ":8766",pubCertPath,devportalPubCert,devportalKeyCert);
                if subscriptionNotification is error {
                    string message = "Error while sending subscription create grpc event";
                    log:printError(subscriptionNotification.toString());
                    APKError e = error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = "500");
                    return e;
                }
            }
        } else {
            return hostList;
        }
    } else {
        return createdSub;
    }
    return createdSub;
}

isolated function getBusinessPlanByName(string policyName) returns string|APKError|NotFoundError {
    string|APKError|NotFoundError policy = getBusinessPlanByNameDAO(policyName);
    return policy;
}

isolated function addMultipleSubscriptions(Subscription[] subscriptions, string org, string user) returns Subscription[]|APKError|NotFoundError|error {
    Subscription[]|APKError addedSubs = [];
    foreach Subscription sub in subscriptions {
        Subscription|APKError|NotFoundError|error subscriptionResponse = check addSubscription(sub, org, user);
        if subscriptionResponse is Subscription {
            if addedSubs is Subscription[] {
                addedSubs.push(subscriptionResponse);
            }
        } else if subscriptionResponse is APKError|NotFoundError {
            return subscriptionResponse;
        }
    }
    return addedSubs;
}

isolated function getSubscriptionById(string subId, string org) returns Subscription|APKError|NotFoundError {
    Subscription|APKError|NotFoundError subscription = getSubscriptionByIdDAO(subId, org);
    return subscription;
}

isolated function deleteSubscription(string subId, string organization) returns string|APKError {
    APKError|string status = deleteSubscriptionDAO(subId,organization);
    if status is string {
        string[]|APKError hostList = retrieveManagementServerHostsList();
        if hostList is string[] {
            string eventId = uuid:createType1AsString();
            time:Utc currTime = time:utcNow();
            string date = time:utcToString(currTime);
            SubscriptionGRPC deleteSubscriptionRequest = {eventId: eventId, uuid: subId, timeStamp: date, organization: organization};
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = managementServerConfig.certPath;
            foreach string host in hostList {
                NotificationResponse|error subscriptionNotification = notification_grpc_client:deleteSubscription(deleteSubscriptionRequest,
                "https://" + host + ":8766",pubCertPath,devportalPubCert,devportalKeyCert);
                if subscriptionNotification is error {
                    string message = "Error while sending subscription delete grpc event";
                    log:printError(subscriptionNotification.toString());
                    APKError e = error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = "500");
                    return e;
                }
            }
        } else {
            return hostList;
        }
    } else {
        return status;
    }
    return status;

}

isolated function updateSubscription(string subId, Subscription payload, string org, string user) returns Subscription|NotFoundError|APKError |error{
    Subscription|NotFoundError|APKError existingSub = getSubscriptionByIdDAO(subId, org);
    if existingSub is Subscription {
        payload.subscriptionId = subId;
    } else {
        return existingSub;
    }
    string apiId = "";
    string appId = "";
    string|NotFoundError|APKError subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is APKError|NotFoundError {
        return subscriberId;
    } 
    string? apiUUID = payload.apiId;
    if apiUUID is string {
        API|NotFoundError|APKError api = getAPIByAPIId(apiUUID);
        if api is NotFoundError|APKError {
            return api;
        } else if api is API {
            string apiInString = api.toJsonString();
            json j = check value:fromJsonString(apiInString);
            apiId = check j.id.ensureType();
        }
    }
    string? appUUID = payload.applicationId;
    if appUUID is string {
        Application|APKError|NotFoundError application = getApplicationById(appUUID, org);
        if application is APKError|NotFoundError {
            return application;
        } else if application is Application {
            string appInString = application.toJsonString();
            json j = check value:fromJsonString(appInString);
            appId = check j.applicationId.ensureType();
        }
    }
    string? businessPlan = payload.throttlingPolicy;
    if businessPlan is string {
        string|APKError|NotFoundError businessPlanID = getBusinessPlanByName(businessPlan);
        if businessPlanID is APKError|NotFoundError {
            return businessPlanID;
        }
        payload.requestedThrottlingPolicy = businessPlan;
    }
    payload.status = "UNBLOCKED";
    Subscription|APKError createdSub = updateSubscriptionDAO(payload,user,apiId,appId);
    if createdSub is Subscription {
        string[]|APKError hostList = retrieveManagementServerHostsList();
        if hostList is string[] {
            string eventId = uuid:createType1AsString();
            time:Utc currTime = time:utcNow();
            string date = time:utcToString(currTime);
            SubscriptionGRPC updateSubscriptionRequest = {eventId: eventId, applicationRef: createdSub.applicationId, 
            apiRef: <string>createdSub.apiId, policyId: createdSub.throttlingPolicy, subStatus:<string>createdSub.status,
            subscriber: user, uuid: subId, timeStamp: date, organization: org};
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = managementServerConfig.certPath;
            foreach string host in hostList {
                NotificationResponse|error subscriptionNotification = notification_grpc_client:updateSubscription(updateSubscriptionRequest,
                "https://" + host + ":8766",pubCertPath,devportalPubCert,devportalKeyCert);
                if subscriptionNotification is error {
                    string message = "Error while sending subscription update grpc event";
                    log:printError(subscriptionNotification.toString());
                    APKError e = error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = "500");
                    return e;
                }
            }
        } else {
            return hostList;
        }
    } else {
        return createdSub;
    }
    return createdSub;
}

isolated function getSubscriptions(string? apiId, string? applicationId, string? groupId, int offset, int limitCount, string org) returns SubscriptionList|APKError|NotFoundError {
    if apiId is string && applicationId is string {
        // Retrieve Subscriptions per given API Id and App Id
        Subscription|APKError|NotFoundError subscription = getSubscriptionByAPIandAppIdDAO(apiId,applicationId,org);
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
        Subscription[]|APKError subs = getSubscriptionsByAPIIdDAO(apiId,org);
        if subs is Subscription[] {
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subs;
        }
    } else if applicationId is string {
        // Retrieve Subscriptions per given APP Id
        Subscription[]|APKError subs = getSubscriptionsByAPPIdDAO(applicationId,org);
        if subs is Subscription[] {
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subs;
        }
    } else {
        // Retrieve All Subscriptions
        Subscription[]|APKError subs = getSubscriptionsList(org);
        if subs is Subscription[] {
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else {
            return subs;
        }
    }
}