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

isolated function addApplication(Application application, commons:Organization org, string user) returns NotFoundError|Application|commons:APKError {
    string applicationId = uuid:createType1AsString();
    application.applicationId = applicationId;
    // TODO: Removed Validate the policy
    // string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    // if policyId is error {
    //     string message = "Invalid Policy";
    //     log:printError(message);
    //     return error(message, policyId, message = message, description = message, code = 909000, statusCode = 500);
    // }
    string|NotFoundError subscriberId = check getSubscriberIdDAO(user, org.uuid);
    if subscriberId is string {
        Application createdApp = check addApplicationDAO(application, subscriberId, org.uuid);
        string[] hostList = check retrieveManagementServerHostsList();
        string eventId = uuid:createType1AsString();
        time:Utc currTime = time:utcNow();
        string date = time:utcToString(currTime);
        ApplicationGRPC createApplicationRequest = {
            eventId: eventId,
            name: createdApp.name,
            uuid: applicationId,
            owner: user,
            policy: "unlimited",
            keys: [],
            attributes: [],
            timeStamp: date,
            organization: org.uuid
        };
        foreach string host in hostList {
            log:printDebug("Retrieved Mgt Host:" + host);
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = managementServerConfig.certPath;
            NotificationResponse|error applicationNotification = notification_grpc_client:createApplication(createApplicationRequest,
                    "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
            if applicationNotification is error {
                string message = "Error while sending application create grpc event";
                log:printError(applicationNotification.toString());
                return error(message, applicationNotification, message = message, description = message, code = 909000, statusCode = 500);
            }
        }
        return application;
    } else {
        return subscriberId;
    }
}

isolated function validateApplicationUsagePolicy(string policyName, commons:Organization org) returns string?|error {
    string?|error policy = getApplicationUsagePlanByNameDAO(policyName, org.uuid);
    return policy;
}

isolated function getApplicationById(string appId, commons:Organization org) returns Application|commons:APKError|NotFoundError {
    Application|commons:APKError|NotFoundError application = getApplicationByIdDAO(appId, org.uuid);
    return application;
}

isolated function getApplicationList(string? sortBy, string? groupId, string? query, string? sortOrder, int 'limit, int offset, commons:Organization org) returns ApplicationList|commons:APKError {
    Application[]|commons:APKError applications = getApplicationsDAO(org.uuid);
    if applications is Application[] {
        int count = applications.length();
        ApplicationList applicationsList = {count: count, list: applications};
        return applicationsList;
    } else {
        return applications;
    }
}

isolated function updateApplication(string appId, Application application, commons:Organization org, string user) returns Application|NotFoundError|commons:APKError {
    Application|commons:APKError|NotFoundError existingApp = getApplicationByIdDAO(appId, org.uuid);
    if existingApp is Application {
        application.applicationId = appId;
    } else {
        Error err = {code: 9010101, message: "Application Not Found"};
        NotFoundError nfe = {body: err};
        return nfe;
    }
    // TODO: Removed Validate the policy
    // string?|error policyId = validateApplicationUsagePolicy(application.throttlingPolicy, org);
    // if policyId is error {
    //     string message = "Invalid Policy";
    //     log:printError(message);
    //     return error(message, policyId, message = message, description = message, code = 909000, statusCode = 500);
    // }
    string|NotFoundError subscriberId = check getSubscriberIdDAO(user, org.uuid);
    if subscriberId is string {
        log:printDebug("subscriber id" + subscriberId.toString());
        Application updatedApp = check updateApplicationDAO(application, subscriberId, org.uuid);
        string[] hostList = check retrieveManagementServerHostsList();
        string eventId = uuid:createType1AsString();
        time:Utc currTime = time:utcNow();
        string date = time:utcToString(currTime);
        ApplicationGRPC createApplicationRequest = {
            eventId: eventId,
            name: updatedApp.name,
            uuid: appId,
            owner: user,
            policy: "unlimited",
            keys: [],
            attributes: [],
            timeStamp: date,
            organization: org.uuid
        };
        foreach string host in hostList {
            log:printDebug("Retrieved Host:" + host);
            string devportalPubCert = <string>keyStores.tls.certFilePath;
            string devportalKeyCert = <string>keyStores.tls.keyFilePath;
            string pubCertPath = managementServerConfig.certPath;
            NotificationResponse|error applicationNotification = notification_grpc_client:createApplication(createApplicationRequest,
                    "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
            if applicationNotification is error {
                string message = "Error while sending application create grpc event";
                log:printError(applicationNotification.toString());
                return error(message, applicationNotification, message = message, description = message, code = 909000, statusCode = 500);
            }
        }
        return updatedApp;
    } else {
        return subscriberId;
    }
}

isolated function deleteApplication(string appId, commons:Organization organization) returns boolean|commons:APKError {
    boolean status = check deleteApplicationDAO(appId, organization.uuid);
    string[] hostList = check retrieveManagementServerHostsList();
    string eventId = uuid:createType1AsString();
    time:Utc currTime = time:utcNow();
    string date = time:utcToString(currTime);
    ApplicationGRPC deleteApplicationRequest = {eventId: eventId, uuid: appId, timeStamp: date, organization: organization.uuid};
    string devportalPubCert = <string>keyStores.tls.certFilePath;
    string devportalKeyCert = <string>keyStores.tls.keyFilePath;
    string pubCertPath = <string>managementServerConfig.certPath;
    foreach string host in hostList {
        NotificationResponse|error applicationNotification = notification_grpc_client:deleteApplication(deleteApplicationRequest,
                "https://" + host + ":8766", pubCertPath, devportalPubCert, devportalKeyCert);
        if applicationNotification is error {
            string message = "Error while sending application delete grpc event";
            log:printError(applicationNotification.toString());
            return error(message, applicationNotification, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
    return status;
}

isolated function generateAPIKey(APIKeyGenerateRequest payload, string appId, string keyType, string user, commons:Organization org) returns APIKey|commons:APKError|NotFoundError {
    Application|NotFoundError application = check getApplicationById(appId, org);
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
                return error(message, payloadValPeriod, message = message, description = message, code = 909000, statusCode = 500);
            }
            record {} addProperties = {};
            record {}? payloadAddProperties = payload.additionalProperties;
            if payloadAddProperties is record {} {
                addProperties = payloadAddProperties;
            } else {
                string message = "Invalid Additional Properties";
                log:printError(message);
                return error(message, payloadAddProperties, message = message, description = message, code = 909000, statusCode = 500);
            }

            // retrieve subscribed APIs
            SubscriptionList|NotFoundError subscriptions = check getSubscriptions(null, appId, null, 0, 0, org);
            API[] apiList = [];
            if subscriptions is SubscriptionList {
                Subscription[]? subArray = subscriptions.list;
                if subArray is Subscription[] {
                    foreach Subscription item in subArray {
                        string? apiUUID = item.apiId;
                        if apiUUID is string {
                            API|NotFoundError api = check getAPIByAPIId(apiUUID);
                            if api is API {
                                apiList.push(api);
                            }
                        } else {
                            string message = "Invalid API UUID found:" + apiUUID.toString();
                            log:printError(message);
                            return error(message, apiUUID, message = message, description = message, code = 909000, statusCode = 500);
                        }
                    }
                }
            }
            APIKey|commons:APKError apiKey = generateAPIKeyForApplication(user, application, apiList, keyType, validityPeriod, addProperties);
            return apiKey;
        } else {
            string message = "User:" + user + " doesn't have permission to Application with application id:" + appId;
            log:printError(message);
            return error(message, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function generateAPIKeyForApplication(string username, Application application, API[] apiList, string keyType, int validityPeriod, record {} addProperties) returns APIKey|commons:APKError {
    if keyType !is "PRODUCTION"|"SANDBOX" {
        string message = "Invalid Key Type:" + keyType;
        log:printError(message);
        return error(message, message = message, description = message, code = 909000, statusCode = 400);
    }
    JWTTokenInfo jwtTokenInfoPayload = {application: application, subscriber: username, expireTime: "", keyType: keyType, permittedIP: "", permittedReferrer: "", subscribedAPIs: apiList};
    string|error token = generateToken(jwtTokenInfoPayload);
    if token is string {
        APIKey apiKey = {apikey: token, validityTime: 3600};
        return apiKey;
    } else {
        string message = "Error while generating token";
        log:printError(message);
        return error(message, token, message = message, description = message, code = 909000, statusCode = 400);
    }
}

isolated function checkUserAccessAllowedForApplication(Application application, string user) returns boolean {
    return true;
}

isolated function retrieveManagementServerHostsList() returns string[]|commons:APKError {
    string managementServerServiceName = managementServerConfig.serviceName;
    string managementServerNamespace = managementServerConfig.namespace;
    log:printDebug("Service:" + managementServerServiceName);
    log:printDebug("Namespace:" + managementServerNamespace);
    string[]|commons:APKError hostList = getPodFromNameAndNamespace(managementServerServiceName, managementServerNamespace);
    return hostList;
}

// # Description
// #
// # + application - Parameter Description  
// # + applicationKeyGenRequest - Parameter Description  
// # + organization - Parameter Description
// # + return - Return Value Description
// public isolated function generateKeysForApplication(Application application, ApplicationKeyGenerateRequest applicationKeyGenRequest, commons:Organization organization) returns commons:APKError {
//     string? keyManager = applicationKeyGenRequest.keyManager;
//     if keyManager is string {
//         KeyManagerDaoEntry keyManagerById = check getKeyManagerById(keyManager, organization);
//         KeyManager keyManagerEntry = check fromKeyManagerDaoEntryToKeyManagerModel(keyManagerById);
//         if !keyManagerEntry.enabled {
//             return error("Key Manager is disabled", message = "Key Manager is disabled", description = "Key Manager is disabled", code = 900951, statusCode = 400);
//         }
//         if !keyManagerEntry.enableOAuthAppCreation {
//             return error("OAuth App Creation is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, message = "OAuth App Creation is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, description = "OAuth App Creation is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, code = 900952, statusCode = 400);
//         }
//         string[]? availableGrantTypes = keyManagerEntry.availableGrantTypes;
//         if availableGrantTypes is string[] {
//             foreach string item in applicationKeyGenRequest.grantTypesToBeSupported {
//                 if availableGrantTypes.indexOf(item) is () {
//                     return error("Grant Type " + item + " is not supported by keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, message = "Grant Type " + item + " is not supported by keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, description = "Grant Type " + item + " is not supported by keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, code = 900953, statusCode = 400);
//                 }
//             }
//         } else {

//         }
//     } else {
//         return error("Key Manager is not provided", message = "Key Manager is not provided", description = "Key Manager is not provided", code = 900953, statusCode = 400);
//     }
// }

