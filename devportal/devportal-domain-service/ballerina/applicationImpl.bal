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
import devportal_service.types;
import devportal_service.kmclient;
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
    check deleteOauthApps(appId, organization);
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

isolated function deleteOauthApps(string appId, commons:Organization organization) returns commons:APKError? {
    types:KeyMappingDaoEntry[] keyMappingEntriesByApplication = check getKeyMappingEntriesByApplication(appId);
    foreach types:KeyMappingDaoEntry item in keyMappingEntriesByApplication {
        KeyManagerDaoEntry keyManagerById = check getKeyManagerById(item.key_manager_uuid, organization);
        types:KeyManager keyManagerConfig = check fromKeyManagerDaoEntryToKeyManagerModel(keyManagerById);
        if keyManagerConfig.enabled && item.create_mode == "CREATED" {
            kmclient:KeyManagerClient kmClient = check getKmClient(keyManagerConfig);
            boolean _ = check kmClient.deleteOauthApplication(item.consumer_key);
        }
    }
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

# Description
#
# + application - Parameter Description  
# + applicationKeyGenRequest - Parameter Description  
# + organization - Parameter Description
# + return - Return Value Description
public isolated function generateKeysForApplication(Application application, ApplicationKeyGenerateRequest applicationKeyGenRequest, commons:Organization organization) returns OkApplicationKey|commons:APKError {
    string? keyManager = applicationKeyGenRequest.keyManager;
    if keyManager is string {
        if check isKeyMappingEntryByApplicationAndKeyManagerExist(<string>application.applicationId, keyManager,applicationKeyGenRequest.keyType) {
            return error("Key Mapping Entry already exists for application " + application.name + " and keyManager " + keyManager, message = "Key Mapping Entry already exists for application " + application.name + " and keyManager " + keyManager, description = "Key Mapping Entry already exists for application " + application.name + " and keyManager " + keyManager, code = 900950, statusCode = 400);
        }
        KeyManagerDaoEntry keyManagerById = check getKeyManagerById(keyManager, organization);
        types:KeyManager keyManagerEntry = check fromKeyManagerDaoEntryToKeyManagerModel(keyManagerById);
        if !keyManagerEntry.enabled {
            return error("Key Manager is disabled", message = "Key Manager is disabled", description = "Key Manager is disabled", code = 900951, statusCode = 400);
        }
        if !keyManagerEntry.enableOAuthAppCreation {
            return error("OAuth App Creation is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, message = "OAuth App Creation is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, description = "OAuth App Creation is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, code = 900952, statusCode = 400);
        }
        string[]? availableGrantTypes = keyManagerEntry.availableGrantTypes;
        if availableGrantTypes is string[] {
            foreach string item in applicationKeyGenRequest.grantTypesToBeSupported {
                if availableGrantTypes.indexOf(item) is () {
                    return error("Grant Type " + item + " is not supported by keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, message = "Grant Type " + item + " is not supported by keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, description = "Grant Type " + item + " is not supported by keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, code = 900953, statusCode = 400);
                }
            }
        }
        kmclient:KeyManagerClient kmClient = check getKmClient(keyManagerEntry);
        kmclient:ClientRegistrationResponse oauthApplicationCreationResponse = check registerOauthApplication(kmClient, application, applicationKeyGenRequest);
        types:KeyMappingDaoEntry keyMappingDaoEntry = {
            application_uuid: <string>application.applicationId,
            consumer_key: <string>oauthApplicationCreationResponse.client_id,
            key_manager_uuid: keyManager,
            key_type: applicationKeyGenRequest.keyType,
            uuid: uuid:createType1AsString(),
            app_info: oauthApplicationCreationResponse.toJsonString().toBytes(),
            create_mode: "CREATED"
        };
        check addKeyMappingEntryForApplication(keyMappingDaoEntry);
        OkApplicationKey okApplicationKey = {
            body: {
                consumerKey: oauthApplicationCreationResponse.client_id,
                consumerSecret: oauthApplicationCreationResponse.client_secret,
                callbackUrls: oauthApplicationCreationResponse.redirect_uris,
                supportedGrantTypes: oauthApplicationCreationResponse.grant_types,
                keyManager: keyManager,
                keyMappingId: keyMappingDaoEntry.uuid,
                keyType: applicationKeyGenRequest.keyType,
                keyState: "CREATED",
                mode: "CREATED",
                additionalProperties: oauthApplicationCreationResponse.additional_properties
            }
        };
        return okApplicationKey;
    } else {
        return error("Key Manager is not provided", message = "Key Manager is not provided", description = "Key Manager is not provided", code = 900953, statusCode = 400);
    }
}

isolated function registerOauthApplication(kmclient:KeyManagerClient keyManagerClient, Application application, ApplicationKeyGenerateRequest oauthAppRegistrationRequest) returns kmclient:ClientRegistrationResponse|commons:APKError {
    kmclient:ClientRegistrationRequest clientRegistrationRequest = {
        client_name: application.name,
        grant_types: oauthAppRegistrationRequest.grantTypesToBeSupported,
        redirect_uris: oauthAppRegistrationRequest.callbackUrls
    };
    return keyManagerClient.registerOauthApplication(clientRegistrationRequest);
}

public isolated function mapKeys(Application application, ApplicationKeyMappingRequest applicationKeyMappingRequest, commons:Organization organization) returns OkApplicationKey|commons:APKError {
    string? keyManager = applicationKeyMappingRequest.keyManager;
    if keyManager is string {
        KeyManagerDaoEntry keyManagerById = check getKeyManagerById(keyManager, organization);
        types:KeyManager keyManagerEntry = check fromKeyManagerDaoEntryToKeyManagerModel(keyManagerById);
        if !keyManagerEntry.enabled {
            commons:APKError e = error("Key Manager is disabled", message = "Key Manager is disabled", description = "Key Manager is disabled", code = 900951, statusCode = 400);
            return e;
        }
        if !keyManagerEntry.enableMapOAuthConsumerApps {
            commons:APKError e = error("OAuth App Mapping is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, message = "OAuth App Mapping is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, description = "OAuth App Mapping is disabled for keyManager " + keyManagerEntry.name + " in organization " + organization.uuid, code = 900952, statusCode = 400);
            return e;
        }
        if keyManagerEntry.enableOauthAppValidation {
            kmclient:KeyManagerClient kmClient = check getKmClient(keyManagerEntry);
            kmclient:ClientRegistrationResponse|commons:APKError retrievedOauthApp = kmClient.retrieveOauthApplicationByClientId(applicationKeyMappingRequest.consumerKey);
            if retrievedOauthApp is kmclient:ClientRegistrationResponse {
                if retrievedOauthApp.client_secret != applicationKeyMappingRequest.consumerSecret {
                    commons:APKError e = error("Consumer Secret is not valid", message = "Consumer Secret is not valid", description = "Consumer Secret is not valid", code = 900953, statusCode = 400);
                    return e;
                }
            } else {
                commons:APKError e = error("Consumer Key is not valid", message = "Consumer Key is not valid", description = "Consumer Key is not valid", code = 900953, statusCode = 400);
                return e;
            }
        }
        types:KeyMappingDaoEntry keyMappingDaoEntry = {
            application_uuid: <string>application.applicationId,
            consumer_key: <string>applicationKeyMappingRequest.consumerKey,
            key_manager_uuid: keyManager,
            key_type: applicationKeyMappingRequest.keyType,
            uuid: uuid:createType1AsString(),
            app_info: "".toBytes(),
            create_mode: "MAPPED"
        };
        check addKeyMappingEntryForApplication(keyMappingDaoEntry);
        OkApplicationKey okApplicationKey = {
            body: {
                consumerKey: applicationKeyMappingRequest.consumerKey,
                consumerSecret: applicationKeyMappingRequest.consumerSecret,
                keyManager: keyManager,
                keyMappingId: keyMappingDaoEntry.uuid,
                keyType: applicationKeyMappingRequest.keyType,
                keyState: "CREATED",
                mode: "MAPPED"
            }
        };
        return okApplicationKey;
    } else {
        commons:APKError e = error("Key Manager is not provided", message = "Key Manager is not provided", description = "Key Manager is not provided", code = 900953, statusCode = 400);
        return e;
    }
}

public isolated function oauthKeys(Application application, commons:Organization organization) returns ApplicationKeyList|commons:APKError {
    types:KeyMappingDaoEntry[] keyMAppingEntries = check getKeyMappingEntriesByApplication(<string>application.applicationId);
    ApplicationKeyList applicationKeyList = {};
    ApplicationKey[] applicationKeys = [];
    foreach types:KeyMappingDaoEntry item in keyMAppingEntries {
        applicationKeys.push(fromKeyMappingDaoEntryToApplicationKey(item, ()));
    }
    applicationKeyList.list = applicationKeys;
    applicationKeyList.count = applicationKeys.length();
    return applicationKeyList;
}

public isolated function oauthKeyByMappingId(Application application, string keyMappingId, commons:Organization organization) returns ApplicationKey|commons:APKError {
    types:KeyMappingDaoEntry keyMappingEntry = check getKeyMappingEntryByApplicationAndKeyMappingId(<string>application.applicationId, keyMappingId);
    types:KeyManager keyManagerEntry = check getKeymanagerByKeyManagerUUID(keyMappingEntry.key_manager_uuid, organization);
    if keyMappingEntry.create_mode == "CREATED" {
        kmclient:KeyManagerClient keyManagerClient = check getKmClient(keyManagerEntry);
        kmclient:ClientRegistrationResponse|commons:APKError retrieveOauthApplicationByClientId = keyManagerClient.retrieveOauthApplicationByClientId(keyMappingEntry.consumer_key);
        if retrieveOauthApplicationByClientId is kmclient:ClientRegistrationResponse {
            return fromKeyMappingDaoEntryToApplicationKey(keyMappingEntry, retrieveOauthApplicationByClientId);
        }
    }
    return fromKeyMappingDaoEntryToApplicationKey(keyMappingEntry, ());
}

isolated function fromKeyMappingDaoEntryToApplicationKey(types:KeyMappingDaoEntry item, kmclient:ClientRegistrationResponse? oauthAppResponse) returns ApplicationKey {
    ApplicationKey applicationKey = {
        keyMappingId: item.uuid,
        consumerKey: item.consumer_key,
        keyManager: item.key_manager_uuid,
        keyType: item.key_type,
        mode: item.create_mode,
        keyState: "CREATED"
    };
    if oauthAppResponse is kmclient:ClientRegistrationResponse {
        applicationKey.consumerSecret = oauthAppResponse.client_secret;
        applicationKey.supportedGrantTypes = oauthAppResponse.grant_types;
        applicationKey.callbackUrls = oauthAppResponse.redirect_uris;
        applicationKey.additionalProperties = oauthAppResponse.additional_properties;
    }
    return applicationKey;
}

public isolated function generateApplicationToken(Application application, string keyMappingId, ApplicationTokenGenerateRequest payload, commons:Organization organization) returns OkApplicationToken|commons:APKError {
    types:KeyMappingDaoEntry keyMappingEntry = check getKeyMappingEntryByApplicationAndKeyMappingId(<string>application.applicationId, keyMappingId);
    types:KeyManager keyManagerEntry = check getKeymanagerByKeyManagerUUID(keyMappingEntry.key_manager_uuid, organization);
    if <boolean>keyManagerEntry.enabled && <boolean>keyManagerEntry.enableTokenGeneration {
        kmclient:KeyManagerClient kmClient = check getKmClient(keyManagerEntry);
        kmclient:TokenRequest tokenRequest = {client_id: keyMappingEntry.consumer_key, client_secret: payload.consumerSecret, scopes: payload.scopes};
        kmclient:TokenResponse generateAccessToken = check kmClient.generateAccessToken(tokenRequest);
        OkApplicationToken applicationToken = {body: {accessToken: generateAccessToken.access_token, tokenScopes: generateAccessToken.scopes, validityTime: generateAccessToken.expires_in}};
        return applicationToken;
    } else {
        return error("Key Manager is disabled", message = "Key Manager is disabled", description = "Key Manager is disabled", code = 900951, statusCode = 400);
    }
}

isolated function getKeymanagerByKeyManagerUUID(string uuid, commons:Organization organization) returns types:KeyManager|commons:APKError {
    KeyManagerDaoEntry keyManagerById = check getKeyManagerById(uuid, organization);
    return check fromKeyManagerDaoEntryToKeyManagerModel(keyManagerById);
}

public isolated function updateOauthApp(Application application, string keyMappingId, ApplicationKey payload, commons:Organization organization) returns ApplicationKey|commons:APKError {
    do {
        types:KeyMappingDaoEntry keyMappingEntry = check getKeyMappingEntryByApplicationAndKeyMappingId(<string>application.applicationId, keyMappingId);
        types:KeyManager keyManagerEntry = check getKeymanagerByKeyManagerUUID(keyMappingEntry.key_manager_uuid, organization);
        if keyManagerEntry.enabled {
            kmclient:KeyManagerClient kmClient = check getKmClient(keyManagerEntry);
            kmclient:ClientRegistrationResponse retrieveOauthApplicationByClientId = check kmClient.retrieveOauthApplicationByClientId(keyMappingEntry.consumer_key);
            kmclient:ClientUpdateRequest clientUpDateRquest = {...retrieveOauthApplicationByClientId};
            clientUpDateRquest.client_id = payload.consumerKey;
            clientUpDateRquest.client_secret = payload.consumerSecret;
            clientUpDateRquest.grant_types = payload.supportedGrantTypes;
            clientUpDateRquest.redirect_uris = payload.callbackUrls;
            record {}? additionalProperties = payload.additionalProperties;
            if additionalProperties is record {} {
                if additionalProperties.hasKey("application_type") {
                    clientUpDateRquest.application_type = <string>additionalProperties["application_type"];
                    _ = additionalProperties.removeIfHasKey("application_type");
                }
                if additionalProperties.hasKey("client_name") {
                    clientUpDateRquest.client_name = <string>additionalProperties["client_name"];
                    _ = additionalProperties.removeIfHasKey("client_name");
                }
                if additionalProperties.hasKey("logo_uri") {
                    clientUpDateRquest.logo_uri = <string>additionalProperties["logo_uri"];
                    _ = additionalProperties.removeIfHasKey("logo_uri");
                }
                if additionalProperties.hasKey("client_uri") {
                    clientUpDateRquest.client_uri = <string>additionalProperties["client_uri"];
                    _ = additionalProperties.removeIfHasKey("client_uri");
                }
                if additionalProperties.hasKey("policy_uri") {
                    clientUpDateRquest.policy_uri = <string>additionalProperties["policy_uri"];
                    _ = additionalProperties.removeIfHasKey("policy_uri");
                }
                if additionalProperties.hasKey("tos_uri") {
                    clientUpDateRquest.tos_uri = <string>additionalProperties["tos_uri"];
                    _ = additionalProperties.removeIfHasKey("tos_uri");
                }
                if additionalProperties.hasKey("jwks_uri") {
                    clientUpDateRquest.jwks_uri = <string>additionalProperties["jwks_uri"];
                    _ = additionalProperties.removeIfHasKey("jwks_uri");
                }
                if additionalProperties.hasKey("subject_type") {
                    clientUpDateRquest.subject_type = <string>additionalProperties["subject_type"];
                    _ = additionalProperties.removeIfHasKey("subject_type");
                }
                if additionalProperties.hasKey("token_endpoint_auth_method") {
                    clientUpDateRquest.token_endpoint_auth_method = <string>additionalProperties["token_endpoint_auth_method"];
                    _ = additionalProperties.removeIfHasKey("token_endpoint_auth_method");
                }
                clientUpDateRquest.additional_properties = additionalProperties;
            }
            kmclient:ClientRegistrationResponse oauthApplicationByClientId = check kmClient.updateOauthApplicationByClientId(keyMappingEntry.consumer_key, clientUpDateRquest);
            types:KeyMappingDaoEntry updatedKeyMappingentry = keyMappingEntry.clone();
            updatedKeyMappingentry.app_info = oauthApplicationByClientId.toJsonString().toBytes();
            check updateKeyMappingEntry(updatedKeyMappingentry);
            return {
                keyMappingId: keyMappingEntry.uuid,
                consumerKey: oauthApplicationByClientId.client_id,
                consumerSecret: oauthApplicationByClientId.client_secret,
                keyManager: keyMappingEntry.key_manager_uuid,
                keyType: keyMappingEntry.key_type,
                mode: keyMappingEntry.create_mode,
                keyState: "CREATED",
                callbackUrls: oauthApplicationByClientId.redirect_uris,
                supportedGrantTypes: oauthApplicationByClientId.grant_types,
                additionalProperties: oauthApplicationByClientId.additional_properties
            };
        } else {
            return error("Key Manager is disabled", message = "Key Manager is disabled", description = "Key Manager is disabled", code = 900951, statusCode = 400);
        }
    } on fail var e {
        return error("Internal Server Error", e, message = "Internal Server Error", description = "Internal Server Error", code = 900952, statusCode = 500);
    }
}
