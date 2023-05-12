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

import wso2/apk_common_lib as commons;
import ballerina/log;
import ballerina/time;
import ballerina/uuid;
import wso2/notification_grpc_client as notification;
import ballerina/http;
import ballerina/mime;

# This function used to get API from database
#
# + return - Return Value string?|APIList|error
isolated function getAPIList(int 'limit, int offset, string? query, string organization) returns APIList|commons:APKError {
    if query !is string {
        API[]|commons:APKError apis = db_getAPIsDAO(organization);
        if apis is API[] {
            API[] limitSet = [];
            if apis.length() > offset {
                foreach int i in offset ... (apis.length() - 1) {
                    if limitSet.length() < 'limit {
                        limitSet.push(apis[i]);
                    }
                }
            }
            APIList apisList = {count: limitSet.length(), list: limitSet, pagination: {total: apis.length(), 'limit: 'limit, offset: offset}};
            return apisList;
        } else {
            return apis;
        }
    } else {
        boolean hasPrefix = query.startsWith("content");
        if hasPrefix {
            int? index = query.indexOf(":");
            if index is int {
                string modifiedQuery = "%" + query.substring(index + 1) + "%";
                API[]|commons:APKError apis = getAPIsByQueryDAO(modifiedQuery, organization);
                if apis is API[] {
                    API[] limitSet = [];
                    if apis.length() > offset {
                        foreach int i in offset ... (apis.length() - 1) {
                            if limitSet.length() < 'limit {
                                limitSet.push(apis[i]);
                            }
                        }
                    }
                    APIList apisList = {count: limitSet.length(), list: limitSet, pagination: {total: apis.length(), 'limit: 'limit, offset: offset}};
                    return apisList;
                } else {
                    return apis;
                }
            } else {
                return e909621();
            }
        } else {
            return e909622();
        }
    }
}

# This function used to change the lifecycle of API
#
# + targetState - lifecycle action
# + apiId - API Id
# + organization - organization
# + return - Return Value LifecycleState|error
isolated function changeLifeCyleState(string targetState, string apiId, string organization) returns LifecycleState|error {
    string prevLCState = check db_getCurrentLCStatus(apiId);
    transaction {
        string|error lcState = db_changeLCState(targetState, apiId);
        if lcState is string {
            string newvLCState = check db_getCurrentLCStatus(apiId);
            string|error lcEvent = db_AddLCEvent(apiId, prevLCState, newvLCState, organization);
            if lcEvent is string {
                check commit;
                json lcPayload = check getTransitionsFromState(targetState);
                LifecycleState lcCr = check lcPayload.cloneWithType(LifecycleState);
                return lcCr;
            } else {
                rollback;
                return error("error while adding LC event" + lcEvent.message());
            }
        } else {
            rollback;
            return error("error while updating LC state" + lcState.message());
        }
    }
}

# This function used to get current state of the API.
#
# + apiId - API Id parameter
# + organization - organization
# + return - Return Value LifecycleState|error
isolated function getLifeCyleState(string apiId) returns LifecycleState|error {
    string|error currentLCState = db_getCurrentLCStatus(apiId);
    if currentLCState is string {
        json lcPayload = check getTransitionsFromState(currentLCState);
        LifecycleState|error lcGet = lcPayload.cloneWithType(LifecycleState);
        if lcGet is error {
            return e909601(lcGet);
        }
        return lcGet;
    } else {
        return currentLCState;
    }
}

# This function used to map user action to LC state
#
# + v - any parameter object
# + return - Return LC state
isolated function actionToLCState(any v) returns string {
    if (v.toString().equalsIgnoreCaseAscii("published")) {
        return "PUBLISHED";
    } else if (v.toString().equalsIgnoreCaseAscii("created")) {
        return "CREATED";
    } else if (v.toString().equalsIgnoreCaseAscii("blocked")) {
        return "BLOCKED";
    } else if (v.toString().equalsIgnoreCaseAscii("deprecated")) {
        return "DEPRECATED";
    } else if (v.toString().equalsIgnoreCaseAscii("prototyped")) {
        return "PROTOTYPED";
    } else if (v.toString().equalsIgnoreCaseAscii("retired")) {
        return "RETIRED";
    } else {
        return "any";
    }
}

# This function used to get the availble event transitions from state
#
# + state - state parameter
# + return - Return Value jsons
isolated function getTransitionsFromState(string state) returns json|error {
    StatesList c = check lifeCycleStateTransitions.cloneWithType(StatesList);
    foreach States x in c.States {
        if (state.equalsIgnoreCaseAscii(x.State)) {
            return x.toJson();
        }
    }

}

# This function used to connect API create service to database
#
# + apiId - API Id parameter
# + return - Return Value LifecycleHistory
isolated function getLcEventHistory(string apiId) returns LifecycleHistory|commons:APKError {
    LifecycleHistoryItem[]|commons:APKError lcHistory = db_getLCEventHistory(apiId);
    if lcHistory is LifecycleHistoryItem[] {
        int count = lcHistory.length();
        LifecycleHistory eventList = {count: count, list: lcHistory};
        return eventList;
    } else {
        return lcHistory;
    }
}

isolated function getSubscriptions(string? apiId) returns SubscriptionList|commons:APKError {
    Subscription[]|commons:APKError subcriptions;
    subcriptions = check db_getSubscriptionsForAPI(apiId.toString());
    if subcriptions is Subscription[] {
        int count = subcriptions.length();
        SubscriptionList subsList = {count: count, list: subcriptions};
        return subsList;
    } else {
        return subcriptions;
    }
}

isolated function blockSubscription(string subscriptionId, string blockState) returns string|commons:APKError {
    if ("blocked".equalsIgnoreCaseAscii(blockState) || "prod_only_blocked".equalsIgnoreCaseAscii(blockState)) {
        commons:APKError|string blockSub = db_blockSubscription(subscriptionId, blockState);
        if blockSub is commons:APKError {
            return blockSub;
        } else {
            SubscriptionInternal|commons:APKError updatedSub = getSubscriptionByIdDAO(subscriptionId);
            if updatedSub is SubscriptionInternal {
                string[]|commons:APKError hostList = retrieveManagementServerHostsList();
                if hostList is string[] {
                    string eventId = uuid:createType1AsString();
                    time:Utc currTime = time:utcNow();
                    string date = time:utcToString(currTime);
                    SubscriptionGRPC updateSubscriptionRequest = {
                        eventId: eventId,
                        applicationRef: updatedSub.applicationId,
                        apiRef: <string>updatedSub.apiId,
                        policyId: updatedSub.throttlingPolicy,
                        subStatus: <string>updatedSub.status,
                        subscriber: "user",
                        uuid: subscriptionId,
                        timeStamp: date,
                        organization: "org"
                    };
                    string backofficePubCert = <string>keyStores.tls.certFilePath;
                    string backofficeKeyCert = <string>keyStores.tls.keyFilePath;
                    string pubCertPath = managementServerConfig.certPath;
                    foreach string host in hostList {
                        NotificationResponse|error subscriptionNotification = notification:updateSubscription(updateSubscriptionRequest,
                        "https://" + host + ":8766", pubCertPath, backofficePubCert, backofficeKeyCert);
                        if subscriptionNotification is error {
                            string message = "Error while sending subscription update grpc event";
                            log:printError(subscriptionNotification.toString());
                            commons:APKError e = error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = 500);
                            return e;
                        }
                    }
                } else {
                    return hostList;
                }
            } else {
                return updatedSub;
            }
            return blockSub;
        }
    } else {
        return e909623();
    }
}

isolated function unblockSubscription(string subscriptionId) returns string|commons:APKError {
    commons:APKError|string unblockSub = db_unblockSubscription(subscriptionId);
    if unblockSub is commons:APKError {
        return unblockSub;
    } else {
        SubscriptionInternal|commons:APKError updatedSub = getSubscriptionByIdDAO(subscriptionId);
        if updatedSub is SubscriptionInternal {
            string[]|commons:APKError hostList = retrieveManagementServerHostsList();
            if hostList is string[] {
                string eventId = uuid:createType1AsString();
                time:Utc currTime = time:utcNow();
                string date = time:utcToString(currTime);
                SubscriptionGRPC updateSubscriptionRequest = {
                    eventId: eventId,
                    applicationRef: updatedSub.applicationId,
                    apiRef: <string>updatedSub.apiId,
                    policyId: updatedSub.throttlingPolicy,
                    subStatus: <string>updatedSub.status,
                    subscriber: "user",
                    uuid: subscriptionId,
                    timeStamp: date,
                    organization: "org"
                };
                string backofficePubCert = <string>keyStores.tls.certFilePath;
                string backofficeKeyCert = <string>keyStores.tls.keyFilePath;
                string pubCertPath = managementServerConfig.certPath;
                foreach string host in hostList {
                    NotificationResponse|error subscriptionNotification = notification:updateSubscription(updateSubscriptionRequest,
                    "https://" + host + ":8766", pubCertPath, backofficePubCert, backofficeKeyCert);
                    if subscriptionNotification is error {
                        string message = "Error while sending subscription update grpc event";
                        log:printError(subscriptionNotification.toString());
                        commons:APKError e = error(message, subscriptionNotification, message = message, description = message, code = 909000, statusCode = 500);
                        return e;
                    }
                }
            } else {
                return hostList;
            }
        } else {
            return updatedSub;
        }
        return unblockSub;
    }
}

isolated function getAPI(string apiId) returns API|commons:APKError {
    API|commons:APKError getAPI = check db_getAPI(apiId);
    return getAPI;
}

isolated function getAPIDefinition(string apiId) returns APIDefinition|commons:APKError {
    APIDefinition|commons:APKError apiDefinition = db_getAPIDefinition(apiId);
    return apiDefinition;
}

isolated function updateAPI(string apiId, ModifiableAPI payload) returns API|commons:APKError {
    API|commons:APKError api = db_updateAPI(apiId, payload);
    return api;
}

isolated function getAllCategoryList(string organization) returns APICategoryList|commons:APKError {
    APICategory[]|commons:APKError categories = getAPICategoriesDAO(organization);
    if categories is APICategory[] {
        int count = categories.length();
        APICategoryList apiCategoriesList = {count: count, list: categories};
        return apiCategoriesList;
    } else {
        return categories;
    }
}

isolated function getBusinessPlans(string organization) returns BusinessPlanList|commons:APKError {
    BusinessPlan[]|commons:APKError businessPlans = getBusinessPlansDAO(organization);
    if businessPlans is BusinessPlan[] {
        int count = businessPlans.length();
        BusinessPlanList BusinessPlansList = {count: count, list: businessPlans};
        return BusinessPlansList;
    } else {
        return businessPlans;
    }
}

isolated function retrieveManagementServerHostsList() returns string[]|commons:APKError {
    string managementServerServiceName = managementServerConfig.serviceName;
    string managementServerNamespace = managementServerConfig.namespace;
    log:printDebug("Service:" + managementServerServiceName);
    log:printDebug("Namespace:" + managementServerNamespace);
    string[]|commons:APKError hostList = getPodFromNameAndNamespace(managementServerServiceName, managementServerNamespace);
    return hostList;
}

isolated function updateThumbnail(string apiId, http:Request message) returns FileInfo|NotFoundError|PreconditionFailedError|commons:APKError|error {
    API|commons:APKError getApi = check db_getAPI(apiId);
    if getApi is commons:APKError|NotFoundError {
        return getApi;
    } else if getApi is API {
        string|() fileName = ();
        byte[]|() fileContent = ();
        string imageType = "";
        mime:Entity[]|http:ClientError payLoadParts = message.getBodyParts();
        if payLoadParts is mime:Entity[] {
            foreach mime:Entity payLoadPart in payLoadParts {
                mime:ContentDisposition contentDisposition = payLoadPart.getContentDisposition();
                string fieldName = contentDisposition.name;
                if fieldName == "file" {
                    fileName = contentDisposition.fileName;
                    fileContent = check payLoadPart.getByteArray();
                    imageType = payLoadPart.getContentType();
                }
            }
        }
        if fileName is () || fileContent is () {
            string msg = "Thumbnail is not provided";
            commons:APKError e = error(msg, (), message = msg, description = msg, code = 909000, statusCode = 500);
            return e;
        } else {
            if !isThumbnailHasValidFileExtention(imageType) {
                PreconditionFailedError pfe = {
                    body: {
                        code: 90915,
                        message: "Thumbnail file extension is not allowed. Supported extensions are .jpg, .png, .jpeg .svg and .gif"
                    }
                };
                return pfe;
            }
            if isFileSizeGreaterThan1MB(fileContent) {
                PreconditionFailedError pfe = {body: {code: 90915, message: "Thumbnail size should be less than 1MB"}};
                return pfe;
            }
            int|commons:APKError thumbnailCategoryId = db_getResourceCategoryIdByCategoryType(RESOURCE_TYPE_THUMBNAIL);
            if thumbnailCategoryId is int {
                Resource thumbnailResource = {
                    resourceUUID: "",
                    apiUuid: apiId,
                    resourceCategoryId: thumbnailCategoryId,
                    dataType: imageType,
                    resourceContent: fileName,
                    resourceBinaryValue: fileContent
                };
                Resource|NotFoundError|commons:APKError thumbnail = db_getResourceByResourceCategory(apiId, thumbnailCategoryId);
                if thumbnail is Resource {
                    thumbnailResource.resourceUUID = thumbnail.resourceUUID;
                    Resource|commons:APKError updatedThumbnail = db_updateResource(thumbnailResource);
                    if updatedThumbnail is Resource {
                        return {fileName: updatedThumbnail.resourceContent, mediaType: updatedThumbnail.dataType};
                    } else {
                        return updatedThumbnail;
                    }
                } else if thumbnail is NotFoundError {
                    string resourceUUID = uuid:createType1AsString();
                    thumbnailResource.resourceUUID = resourceUUID;
                    Resource|commons:APKError addedThumbnail = db_addResource(thumbnailResource);
                    if addedThumbnail is Resource {
                        return {fileName: addedThumbnail.resourceContent, mediaType: addedThumbnail.dataType};
                    } else {
                        return addedThumbnail;
                    }
                } else {
                    return thumbnail;
                }
            } else {
                return thumbnailCategoryId;
            }
        }
    }
}

isolated function getThumbnail(string apiId) returns http:Response|NotFoundError|commons:APKError {
    API|commons:APKError getApi = check db_getAPI(apiId);
    if getApi is API {
        int|commons:APKError thumbnailCategoryId = db_getResourceCategoryIdByCategoryType(RESOURCE_TYPE_THUMBNAIL);
        if thumbnailCategoryId is int {
            Resource|NotFoundError|commons:APKError thumbnail = db_getResourceByResourceCategory(apiId, thumbnailCategoryId);
            if thumbnail is Resource {
                http:Response outResponse = new;
                outResponse.setBinaryPayload(thumbnail.resourceBinaryValue, thumbnail.dataType);
                return outResponse;
            } else {
                return thumbnail;
            }
        }
        return thumbnailCategoryId;
    } else {
        return getApi;
    }
}
