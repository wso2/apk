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
import wso2/apk_common_lib as commons;
import ballerina/time;
import ballerina/log;

isolated function addSubscription(Subscription payload, commons:Organization org, string user) returns Subscription|NotFoundError|commons:APKError {
    do {
        string apiId = "";
        string appId = "";
        string|NotFoundError subscriberId = check getSubscriberIdDAO(user, org.uuid);
        if subscriberId !is string {
            return subscriberId;
        }
        string? apiUUID = payload.apiId;
        if apiUUID is string {
            API|NotFoundError api = check getAPIByAPIId(apiUUID);
            if api is NotFoundError {
                return api;
            } else if api is API {
                string apiInString = api.toJsonString();
                json j = check value:fromJsonString(apiInString);
                apiId = check j.id.ensureType();
            }
        }
        string? appUUID = payload.applicationId;
        if appUUID is string {
            Application|NotFoundError application = check getApplicationById(appUUID, org);
            if application is NotFoundError {
                return application;
            } else if application is Application {
                string appInString = application.toJsonString();
                json j = check value:fromJsonString(appInString);
                appId = check j.applicationId.ensureType();
            }
        }
        // TODO: Removed Validate the policy name
        // string? businessPlan = payload.throttlingPolicy;
        // if businessPlan is string {
        //     string|commons:APKError|NotFoundError businessPlanID = getBusinessPlanByName(businessPlan, org);
        //     if businessPlanID is APKError|NotFoundError {
        //         return businessPlanID;
        //     }
        //     payload.requestedThrottlingPolicy = businessPlan;
        // }
        string subscriptionId = uuid:createType1AsString();
        boolean|error isSubscriptionWorkflowEnable = isSubsciptionWorkflowEnabled(org.uuid);
        if isSubscriptionWorkflowEnable is error {
            string message = "Error while checking subscription workflow";
            return error(message, message = message, description = message, code = 909000, statusCode = 500);
        } else if (isSubscriptionWorkflowEnable) {
            string|error subworkflow = addSubscriptionCreationWorkflow(subscriptionId, org.uuid);
            if subworkflow is error {
                string message = "Error while creating subscription workflow";
                return error(message, subworkflow, message = message, description = message, code = 909000, statusCode = 500);
            }
            payload.subscriptionCreateState = "CREATED";
            Subscription createdSub = check addSubscriptionDAO(payload, user, apiId, appId);
            return createdSub;
        } else {
            payload.subscriptionCreateState = "APPROVED";
            payload.subscriptionId = subscriptionId;
            payload.status = "UNBLOCKED";
            Subscription createdSub = check addSubscriptionDAO(payload, user, apiId, appId);
            string[] hostList = check retrieveManagementServerHostsList();
            string eventId = uuid:createType1AsString();
            time:Utc currTime = time:utcNow();
            string date = time:utcToString(currTime);
            SubscriptionGRPC createSubscriptionRequest = {
                eventId: eventId,
                applicationRef: createdSub.applicationId,
                apiRef: <string>createdSub.apiId,
                policyId: "unlimited",
                subStatus: <string>createdSub.status,
                subscriber: user,
                uuid: subscriptionId,
                timeStamp: date,
                organization: org.uuid
            };
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = managementServerConfig.certPath;
            foreach string host in hostList {
                NotificationResponse|error subscriptionNotification = notification_grpc_client:createSubscription(createSubscriptionRequest,
                    "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
                if subscriptionNotification is error {
                    string message = "Error while sending subscription create grpc event";
                    log:printError(subscriptionNotification.toString());
                    return error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = 500);
                }
            }
            return createdSub;
        }
    } on fail var e {
        return error("Internal Error", e, code = 900900, description = "Internal Error", statusCode = 500, message = "Internal Error");
    }
}

isolated function getBusinessPlanByName(string policyName, commons:Organization org) returns string|commons:APKError|NotFoundError {
    string|commons:APKError|NotFoundError policy = getBusinessPlanByNameDAO(policyName, org.uuid);
    return policy;
}

isolated function addMultipleSubscriptions(Subscription[] subscriptions, commons:Organization org, string user) returns Subscription[]|commons:APKError|NotFoundError {
    Subscription[]|commons:APKError addedSubs = [];
    foreach Subscription sub in subscriptions {
        Subscription|commons:APKError|NotFoundError subscriptionResponse = check addSubscription(sub, org, user);
        if subscriptionResponse is Subscription {
            if addedSubs is Subscription[] {
                addedSubs.push(subscriptionResponse);
            }
        } else if subscriptionResponse is commons:APKError|NotFoundError {
            return subscriptionResponse;
        }
    }
    return addedSubs;
}

isolated function getSubscriptionById(string subId, commons:Organization org) returns Subscription|commons:APKError|NotFoundError {
    Subscription|NotFoundError subscription = check getSubscriptionByIdDAO(subId, org.uuid);
    return subscription;
}

isolated function deleteSubscription(string subId, commons:Organization organization) returns commons:APKError? {
    check deleteSubscriptionDAO(subId, organization.uuid);
    string[] hostList = check retrieveManagementServerHostsList();
    string eventId = uuid:createType1AsString();
    time:Utc currTime = time:utcNow();
    string date = time:utcToString(currTime);
    SubscriptionGRPC deleteSubscriptionRequest = {eventId: eventId, uuid: subId, timeStamp: date, organization: organization.uuid};
    string devportalPubCert = <string>keyStores.tls.certFilePath;
    string devportalKeyCert = <string>keyStores.tls.keyFilePath;
    string pubCertPath = managementServerConfig.certPath;
    foreach string host in hostList {
        NotificationResponse|error subscriptionNotification = notification_grpc_client:deleteSubscription(deleteSubscriptionRequest,
                "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
        if subscriptionNotification is error {
            string message = "Error while sending subscription delete grpc event";
            log:printError(subscriptionNotification.toString());
            return error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function updateSubscription(string subId, Subscription payload, commons:Organization org, string user) returns Subscription|NotFoundError|commons:APKError {
    do {
        Subscription|NotFoundError existingSub = check getSubscriptionByIdDAO(subId, org.uuid);
        if existingSub is Subscription {
            payload.subscriptionId = subId;
        } else {
            return existingSub;
        }
        string apiId = "";
        string appId = "";
        string|NotFoundError subscriberId = check getSubscriberIdDAO(user, org.uuid);
        if subscriberId is NotFoundError {
            return subscriberId;
        }
        string? apiUUID = payload.apiId;
        if apiUUID is string {
            API|NotFoundError api = check getAPIByAPIId(apiUUID);
            if api is NotFoundError {
                return api;
            } else if api is API {
                string apiInString = api.toJsonString();
                json j = check value:fromJsonString(apiInString);
                apiId = check j.id.ensureType();
            }
        }
        string? appUUID = payload.applicationId;
        if appUUID is string {
            Application|NotFoundError application = check getApplicationById(appUUID, org);
            if application is NotFoundError {
                return application;
            } else if application is Application {
                string appInString = application.toJsonString();
                json j = check value:fromJsonString(appInString);
                appId = check j.applicationId.ensureType();
            }
        }
        // TODO: Removed Validate the policy name
        // string? businessPlan = payload.throttlingPolicy;
        // if businessPlan is string {
        //     string|commons:APKError|NotFoundError businessPlanID = getBusinessPlanByName(businessPlan, org);
        //     if businessPlanID is APKError|NotFoundError {
        //         return businessPlanID;
        //     }
        //     payload.requestedThrottlingPolicy = businessPlan;
        // }
        payload.status = "UNBLOCKED";
        Subscription createdSub = check updateSubscriptionDAO(payload, user, apiId, appId);
        string[] hostList = check retrieveManagementServerHostsList();
        string eventId = uuid:createType1AsString();
        time:Utc currTime = time:utcNow();
        string date = time:utcToString(currTime);
        SubscriptionGRPC updateSubscriptionRequest = {
            eventId: eventId,
            applicationRef: createdSub.applicationId,
            apiRef: <string>createdSub.apiId,
            policyId: "unlimited",
            subStatus: <string>createdSub.status,
            subscriber: user,
            uuid: subId,
            timeStamp: date,
            organization: org.uuid
        };
        string devportalPubCert = <string>keyStores.tls.certFilePath;
        string devportalKeyCert = <string>keyStores.tls.keyFilePath;
        string pubCertPath = managementServerConfig.certPath;
        foreach string host in hostList {
            NotificationResponse|error subscriptionNotification = notification_grpc_client:updateSubscription(updateSubscriptionRequest,
                "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
            if subscriptionNotification is error {
                string message = "Error while sending subscription update grpc event";
                log:printError(subscriptionNotification.toString());
                return error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = 500);
            }
        }
        return createdSub;
    } on fail var e {
        return error("Internal Error", e, code = 900900, description = "Internal Error", statusCode = 500, message = "Internal Error");
    }
}

isolated function getSubscriptions(string? apiId, string? applicationId, string? groupId, int offset, int limitCount, commons:Organization org) returns SubscriptionList|commons:APKError|NotFoundError {
    if apiId is string && applicationId is string {
        // Retrieve Subscriptions per given API Id and App Id
        Subscription|NotFoundError subscription = check getSubscriptionByAPIandAppIdDAO(apiId, applicationId, org.uuid);
        if subscription is Subscription {
            Subscription[] subs = [subscription];
            int count = subs.length();
            SubscriptionList subList = {count: count, list: subs};
            return subList;
        } else if subscription is NotFoundError  {
            return subscription;
        } else {
            string message = "Internal Error occured while retrieving Subscription";
            return error(message, message = message, description = message, code = 909001, statusCode = 500);
        }
    } else if apiId is string {
        // Retrieve Subscriptions per given API Id
        Subscription[] subs = check getSubscriptionsByAPIIdDAO(apiId, org.uuid);
        int count = subs.length();
        SubscriptionList subList = {count: count, list: subs};
        return subList;
    } else if applicationId is string {
        // Retrieve Subscriptions per given APP Id
        Subscription[] subs = check getSubscriptionsByAPPIdDAO(applicationId, org.uuid);
        int count = subs.length();
        SubscriptionList subList = {count: count, list: subs};
        return subList;
    } else {
        // Retrieve All Subscriptions
        Subscription[] subs = check getSubscriptionsList(org.uuid);
        int count = subs.length();
        SubscriptionList subList = {count: count, list: subs};
        return subList;
    }
}
