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

isolated function addApplication(Application application, string org, string user) returns string?|Application|error {
    string applicationId = uuid:createType1AsString();
    application.applicationId = applicationId;
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        log:printError("Invalid Policy");
        return error("Invalid Policy");
    }
    int|error? subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is int {
        log:printDebug("subscriber id" + subscriberId.toString());
        string?|Application|error createdApp = addApplicationDAO(application, subscriberId, org);
        return createdApp;
    } else {
        return subscriberId;
    }
}

isolated function validateApplicationUsagePolicy(string policyName, string org) returns string?|error {
    string?|error policy = getApplicationUsagePlanByNameDAO(policyName,org);
    return policy;
}

isolated function getApplicationById(string appId, string org) returns string?|Application|error {
    string?|Application|error application = getApplicationByIdDAO(appId, org);
    return application;
}

isolated function getApplicationList(string? sortBy, string? groupId, string? query, string? sortOrder, int 'limit, int offset, string org) returns string?|ApplicationList|error {
    Application[]|error? applications = getApplicationsDAO(org);
    if applications is Application[] {
        int count = applications.length();
        ApplicationList applicationsList = {count: count, list: applications};
        return applicationsList;
    } else {
        return applications;
    }
}

isolated function updateApplication(string appId, Application application, string org, string user) returns string?|Application|NotFoundError|error {
    string?|Application|error existingApp = getApplicationByIdDAO(appId, org);
    if existingApp is Application {
        application.applicationId = appId;
    } else {
        Error err = {code:9010101, message:"Application Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    }
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        log:printError("Invalid Policy");
        return error("Invalid Policy");
    }
    int|error? subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is int {
        log:printDebug("subscriber id" + subscriberId.toString());
        string?|Application|error updatedApp = updateApplicationDAO(application, subscriberId, org);
        return updatedApp;
    } else {
        return subscriberId;
    }
}

isolated function deleteApplication(string appId, string organization) returns string|error? {
    error?|string status = deleteApplicationDAO(appId,organization);
    return status;
}

isolated function generateAPIKey(APIKeyGenerateRequest payload, string appId, string keyType, string user, string org) returns APIKey|error {
    string?|Application|error application = getApplicationById(appId, org);
    if application !is Application {
        return error("Invalid Application. Application with application id:" + appId + " not found");
    } else {
        boolean userAllowed = checkUserAccessAllowedForApplication(application, user);
        if userAllowed {
            int validityPeriod = 0;
            int? payloadValPeriod = payload.validityPeriod;
            if payloadValPeriod is int {
                validityPeriod = payloadValPeriod;
            } else {
                return error("Invalid validity period");
            }
            record {} addProperties = {};
            record {}? payloadAddProperties = payload.additionalProperties;
            if payloadAddProperties is record {} {
                addProperties = payloadAddProperties;
            } else {
                return error("Invalid Additional Properties");
            }

            // retrieve subscribed APIs
            string?|SubscriptionList|error subscriptions =  getSubscriptions(null, appId, null, 0, 0, org);
            API[] apiList = [];
            if subscriptions is SubscriptionList {
                Subscription[]? subArray = subscriptions.list;
                if subArray is Subscription[] {
                    foreach Subscription item in subArray {
                        string? apiUUID = item.apiId;
                        if apiUUID is string {
                            string?|API|error api = getAPIByAPIId(apiUUID, org);
                            if api is API {
                                apiList.push(api);
                            }
                        } else {
                            log:printDebug("Invalid API UUID found:" + apiUUID.toString());
                            return error("Invalid API UUID found:" + apiUUID.toString());
                        }
                    }
                }
            } else if subscriptions is error {
                log:printDebug(subscriptions.message());
                return error(subscriptions.message());
            }

            APIKey|error apiKey = generateAPIKeyForApplication(user, application, apiList, keyType, validityPeriod, addProperties);
            return apiKey;
        } else {
            return error("User:"+ user +" doesn't have permission to Application with application id:" + appId);
        }
    }
}

isolated function generateAPIKeyForApplication(string username, Application application, API[] apiList, string keyType, int validityPeriod, record {} addProperties) returns APIKey|error {
    if keyType !is "PRODUCTION" | "SANDBOX" {
        return error("Invalid Key Type:" + keyType);
    }
    JWTTokenInfo jwtTokenInfoPayload = {application: application, subscriber: username, expireTime: "", keyType: keyType, permittedIP: "", permittedReferrer: "", subscribedAPIs: apiList };
    string|error token = generateToken(jwtTokenInfoPayload);
    if token is string {
        APIKey apiKey = {apikey:token,validityTime:3600};
        return apiKey;
    } else {
        log:printDebug(token.message());
        return error("Error while generating token");
    }
}

isolated function checkUserAccessAllowedForApplication(Application application, string user) returns boolean {
    return true;
}
