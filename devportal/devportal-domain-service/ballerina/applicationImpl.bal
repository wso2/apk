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

isolated function addApplication(Application application, string org, string user) returns NotFoundError|Application|APKError {
    string applicationId = uuid:createType1AsString();
    application.applicationId = applicationId;
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        string message = "Invalid Policy";
        log:printError(message);
        return error(message, policyId, message = message, description = message, code = 909000, statusCode = "500");
    }
    int|NotFoundError|APKError subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is int {
        log:printDebug("subscriber id" + subscriberId.toString());
        Application|APKError createdApp = addApplicationDAO(application, subscriberId, org);
        return createdApp;
    } else {
        return subscriberId;
    }
}

isolated function validateApplicationUsagePolicy(string policyName, string org) returns string?|error {
    string?|error policy = getApplicationUsagePlanByNameDAO(policyName,org);
    return policy;
}

isolated function getApplicationById(string appId, string org) returns Application|APKError|NotFoundError {
    Application|APKError|NotFoundError application = getApplicationByIdDAO(appId, org);
    return application;
}

isolated function getApplicationList(string? sortBy, string? groupId, string? query, string? sortOrder, int 'limit, int offset, string org) returns ApplicationList|APKError {
    Application[]|APKError applications = getApplicationsDAO(org);
    if applications is Application[] {
        int count = applications.length();
        ApplicationList applicationsList = {count: count, list: applications};
        return applicationsList;
    } else {
        return applications;
    }
}

isolated function updateApplication(string appId, Application application, string org, string user) returns Application|NotFoundError|APKError {
    Application|APKError|NotFoundError existingApp = getApplicationByIdDAO(appId, org);
    if existingApp is Application {
        application.applicationId = appId;
    } else {
        Error err = {code:9010101, message:"Application Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    }
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        string message = "Invalid Policy";
        log:printError(message);
        return error(message, policyId, message = message, description = message, code = 909000, statusCode = "500");
    }
    int|NotFoundError|APKError subscriberId = getSubscriberIdDAO(user,org);
    if subscriberId is int {
        log:printDebug("subscriber id" + subscriberId.toString());
        Application|APKError updatedApp = updateApplicationDAO(application, subscriberId, org);
        return updatedApp;
    } else {
        return subscriberId;
    }
}

isolated function deleteApplication(string appId, string organization) returns string|APKError {
    APKError|string status = deleteApplicationDAO(appId,organization);
    return status;
}

isolated function generateAPIKey(APIKeyGenerateRequest payload, string appId, string keyType, string user, string org) returns APIKey|APKError|NotFoundError {
    Application|APKError|NotFoundError application = getApplicationById(appId, org);
    if application !is Application {
        return application;
    } else {
        boolean userAllowed = checkUserAccessAllowedForApplication(application, user);
        if userAllowed {
            int validityPeriod = 0;
            int? payloadValPeriod = payload.validityPeriod;
            if payloadValPeriod is int {
                validityPeriod = payloadValPeriod;
            } else {
                string message = "Invalid validity period";
                log:printError(message);
                return error(message, payloadValPeriod, message = message, description = message, code = 909000, statusCode = "500");
            }
            record {} addProperties = {};
            record {}? payloadAddProperties = payload.additionalProperties;
            if payloadAddProperties is record {} {
                addProperties = payloadAddProperties;
            } else {
                string message = "Invalid Additional Properties";
                log:printError(message);
                return error(message, payloadAddProperties, message = message, description = message, code = 909000, statusCode = "500");
            }

            // retrieve subscribed APIs
            SubscriptionList|APKError|NotFoundError subscriptions =  getSubscriptions(null, appId, null, 0, 0, org);
            API[] apiList = [];
            if subscriptions is SubscriptionList {
                Subscription[]? subArray = subscriptions.list;
                if subArray is Subscription[] {
                    foreach Subscription item in subArray {
                        string? apiUUID = item.apiId;
                        if apiUUID is string {
                            API|NotFoundError|APKError api = getAPIByAPIId(apiUUID, org);
                            if api is API {
                                apiList.push(api);
                            }
                        } else {
                            string message = "Invalid API UUID found:" + apiUUID.toString();
                            log:printError(message);
                            return error(message, apiUUID, message = message, description = message, code = 909000, statusCode = "500");
                        }
                    }
                }
            } else if subscriptions is APKError {
                return subscriptions;
            }
            APIKey|APKError apiKey = generateAPIKeyForApplication(user, application, apiList, keyType, validityPeriod, addProperties);
            return apiKey;
        } else {
            string message ="User:"+ user +" doesn't have permission to Application with application id:" + appId;
            log:printError(message);
            return error(message, message = message, description = message, code = 909000, statusCode = "500");
        }
    }
}

isolated function generateAPIKeyForApplication(string username, Application application, API[] apiList, string keyType, int validityPeriod, record {} addProperties) returns APIKey|APKError {
    if keyType !is "PRODUCTION" | "SANDBOX" {
        string message = "Invalid Key Type:" + keyType;
        log:printError(message);
        return error(message, message = message, description = message, code = 909000, statusCode = "400");
    }
    JWTTokenInfo jwtTokenInfoPayload = {application: application, subscriber: username, expireTime: "", keyType: keyType, permittedIP: "", permittedReferrer: "", subscribedAPIs: apiList };
    string|error token = generateToken(jwtTokenInfoPayload);
    if token is string {
        APIKey apiKey = {apikey:token,validityTime:3600};
        return apiKey;
    } else {
        string message = "Error while generating token";
        log:printError(message);
        return error(message, token, message = message, description = message, code = 909000, statusCode = "400");
    }
}

isolated function checkUserAccessAllowedForApplication(Application application, string user) returns boolean {
    return true;
}
