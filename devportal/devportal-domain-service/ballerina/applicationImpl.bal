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
import wso2/apk_common_lib as commons;
import wso2/notification_grpc_client;
import ballerina/time;


isolated function addApplication(Application application, commons:Organization org, string user) returns NotFoundError|Application|APKError {
    string applicationId = uuid:createType1AsString();
    application.applicationId = applicationId;
    string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    if policyId is error {
        string message = "Invalid Policy";
        log:printError(message);
        return error(message, policyId, message = message, description = message, code = 909000, statusCode = "500");
    }
    string|NotFoundError|APKError subscriberId = getSubscriberIdDAO(user,org.uuid);
    if subscriberId is string {
        Application|APKError createdApp = addApplicationDAO(application, subscriberId, org.uuid);
        if createdApp is Application {
            string[]|APKError hostList = retrieveManagementServerHostsList();
            if hostList is string[] {
                string eventId = uuid:createType1AsString();
                time:Utc currTime = time:utcNow();
                string date = time:utcToString(currTime);
                ApplicationGRPC createApplicationRequest = {eventId: eventId, name: createdApp.name, uuid: applicationId, 
                owner: user, policy: createdApp.throttlingPolicy, keys: [],  
                attributes: [], timeStamp: date, organization: org};
                foreach string host in hostList {
                    log:printDebug("Retrieved Mgt Host:"+host);
                    string devportalPubCert = <string>keyStores.tls.certFilePath;
                    string devportalKeyCert = <string>keyStores.tls.keyFilePath;
                    string pubCertPath = managementServerConfig.certPath;
                    NotificationResponse|error applicationNotification = notification_grpc_client:createApplication(createApplicationRequest,
                    "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
                    if applicationNotification is error {
                        string message = "Error while sending application create grpc event";
                        log:printError(applicationNotification.toString());
                        return error(message, applicationNotification, message = message, description = message, code = 909000, statusCode = "500");
                    }
                }
            } else {
                return hostList;
            }
        }
        return application;

    } else {
        return subscriberId;
    }
}

isolated function validateApplicationUsagePolicy(string policyName, commons:Organization org) returns string?|error {
    string?|error policy = getApplicationUsagePlanByNameDAO(policyName,org.uuid);
    return policy;
}

isolated function getApplicationById(string appId, commons:Organization org) returns Application|APKError|NotFoundError {
    Application|APKError|NotFoundError application = getApplicationByIdDAO(appId, org.uuid);
    return application;
}

isolated function getApplicationList(string? sortBy, string? groupId, string? query, string? sortOrder, int 'limit, int offset, commons:Organization org) returns ApplicationList|APKError {
    Application[]|APKError applications = getApplicationsDAO(org.uuid);
    if applications is Application[] {
        int count = applications.length();
        ApplicationList applicationsList = {count: count, list: applications};
        return applicationsList;
    } else {
        return applications;
    }
}

isolated function updateApplication(string appId, Application application, commons:Organization org, string user) returns Application|NotFoundError|APKError {
    Application|APKError|NotFoundError existingApp = getApplicationByIdDAO(appId, org.uuid);
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
    string|NotFoundError|APKError subscriberId = getSubscriberIdDAO(user,org.uuid);
    if subscriberId is string {
        log:printDebug("subscriber id" + subscriberId.toString());
        Application|APKError updatedApp = updateApplicationDAO(application, subscriberId, org.uuid);
        if updatedApp is Application {
            string[]|APKError hostList = retrieveManagementServerHostsList();
            if hostList is string[] {
                string eventId = uuid:createType1AsString();
                time:Utc currTime = time:utcNow();
                string date = time:utcToString(currTime);
                ApplicationGRPC createApplicationRequest = {eventId: eventId, name: updatedApp.name, uuid: appId, 
                owner: user, policy: updatedApp.throttlingPolicy, keys: [],  
                attributes: [], timeStamp: date, organization: org};
                foreach string host in hostList {
                    log:printDebug("Retrieved Host:"+host);
                    string devportalPubCert = <string>keyStores.tls.certFilePath;
                    string devportalKeyCert = <string>keyStores.tls.keyFilePath;
                    string pubCertPath = managementServerConfig.certPath;
                    NotificationResponse|error applicationNotification = notification_grpc_client:createApplication(createApplicationRequest,
                    "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
                    if applicationNotification is error {
                        string message = "Error while sending application create grpc event";
                        log:printError(applicationNotification.toString());
                        return error(message, applicationNotification, message = message, description = message, code = 909000, statusCode = "500");
                    }
                }
            } else {
                return hostList;
            }
        } else {
            return updatedApp;
        }
        return updatedApp;
    } else {
        return subscriberId;
    }
}

isolated function deleteApplication(string appId, commons:Organization organization) returns string|APKError {
    APKError|string status = deleteApplicationDAO(appId,organization.uuid);
    if status is string {
        string[]|APKError hostList = retrieveManagementServerHostsList();
        if hostList is string[] {
            string eventId = uuid:createType1AsString();
            time:Utc currTime = time:utcNow();
            string date = time:utcToString(currTime);
            ApplicationGRPC deleteApplicationRequest = {eventId: eventId, uuid: appId, timeStamp: date, organization: organization};
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = <string>managementServerConfig.certPath;
            foreach string host in hostList {
                NotificationResponse|error applicationNotification = notification_grpc_client:deleteApplication(deleteApplicationRequest,
                "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
                if applicationNotification is error {
                    string message = "Error while sending application delete grpc event";
                    log:printError(applicationNotification.toString());
                    return error(message, applicationNotification, message = message, description = message, code = 909000, statusCode = "500");
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

isolated function generateAPIKey(APIKeyGenerateRequest payload, string appId, string keyType, string user, commons:Organization org) returns APIKey|APKError|NotFoundError {
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
                            API|NotFoundError|APKError api = getAPIByAPIId(apiUUID);
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

isolated function retrieveManagementServerHostsList() returns string[]|APKError {
    string managementServerServiceName = managementServerConfig.serviceName;
    string managementServerNamespace = managementServerConfig.namespace;
    log:printDebug("Service:" + managementServerServiceName);
    log:printDebug("Namespace:" + managementServerNamespace);
    string[]|APKError hostList = getPodFromNameAndNamespace(managementServerServiceName,managementServerNamespace);
    return hostList;
 }
