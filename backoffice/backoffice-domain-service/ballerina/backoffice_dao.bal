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

import ballerina/sql;
import ballerina/time;
import ballerinax/postgresql;
import ballerina/io;
import wso2/apk_common_lib as commons;
import ballerina/log;

isolated function db_getAPIsDAO(string organization, string[] groups) returns APIInfo[]|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        do {
            sql:ParameterizedQuery GET_API;
            if (groups is string[] && groups.length() > 0) {
                GET_API = sql:queryConcat(
                    `SELECT
                        UUID AS ID, API_NAME as NAME, API_VERSION as VERSION, CONTEXT, ORGANIZATION, STATUS as STATE,
                        string_to_array(SDK::text,',')::text[] AS SDK, string_to_array(API_TIER::text,',') AS POLICIES,
                        ARTIFACT as ARTIFACT
                    FROM API WHERE ARTIFACT->>'accessControl' = 'NONE' OR ( ARTIFACT->>'accessControl' = 'RESTRICTED' AND
                        (
                         SELECT COUNT(DISTINCT groups)
                         FROM jsonb_array_elements_text(api.artifact->'accessControlGroups') AS groups
                         WHERE groups IN (`, sql:arrayFlattenQuery(groups), `)
                        ) >= 1 )
                    AND ORGANIZATION = ${organization}`
                );
            } else {
                GET_API =
                    `SELECT
                        UUID AS ID, API_NAME as NAME, API_VERSION as VERSION, CONTEXT, ORGANIZATION, STATUS as STATE,
                        string_to_array(SDK::text,',')::text[] AS SDK, string_to_array(API_TIER::text,',') AS POLICIES,
                        ARTIFACT as ARTIFACT
                    FROM API WHERE ARTIFACT->>'accessControl' = 'NONE' AND ORGANIZATION = ${organization}`;
            }
            stream<APIInfo, sql:Error?> apisStream = db_Client->query(GET_API);
            APIInfo[] apis = check from APIInfo api in apisStream
                select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            return e909607(e);
        }
    }
}

isolated function db_changeLCState(string targetState, string apiId) returns string|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        string newState = actionToLCState(targetState);
        if newState.equalsIgnoreCaseAscii("any") {
            return e909609();
        }
        sql:ParameterizedQuery UPDATE_API_LifeCycle_Prefix = `UPDATE api SET status = `;
        sql:ParameterizedQuery values = `${newState}
        WHERE uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_LifeCycle_Prefix, values);

        sql:ExecutionResult|sql:Error result = db_Client->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return targetState;
        } else {
            return e909608(result);
        }
    }
}

isolated function db_getCurrentLCStatus(string apiId) returns string|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_API_LifeCycle_Prefix = `SELECT status from api where uuid = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_API_LifeCycle_Prefix, values);

        string|sql:Error result = db_Client->queryRow(sqlQuery);

        if result is string {
            return result;
        } else {
            return e909610(result);
        }
    }
}

# Update LC event to the database 
#
# + apiId - API id Parameter
# + organization - organization
# + prev_state - prev_state 
# + new_state - new_state
# + return - API | error
isolated function db_AddLCEvent(string? apiId, string? prev_state, string? new_state, string organization) returns string|commons:APKError {
    postgresql:Client|error db_client = getConnection();
    time:Utc utc = time:utcNow();
    if db_client is error {
        return e909601(db_client);
    } else {
        sql:ParameterizedQuery values = `${apiId},
                                        ${prev_state}, 
                                        ${new_state},
                                        'apkuser',
                                        ${organization},
                                        ${utc}
                                    )`;
        sql:ParameterizedQuery ADD_LC_EVENT_Prefix = `INSERT INTO api_lc_event (api_uuid,previous_state,new_state,user_uuid,organization,event_date) VALUES (`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_LC_EVENT_Prefix, values);

        sql:ExecutionResult|sql:Error result = db_client->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return result.toString();
        } else {
            return e909611(result);
        }
    }
}

isolated function db_getLCEventHistory(string apiId) returns LifecycleHistoryItem[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT previous_state, new_state, user_uuid, event_date FROM api_lc_event WHERE api_uuid =${apiId}`;
            stream<LifecycleHistoryItem, sql:Error?> lcStream = dbClient->query(query);
            LifecycleHistoryItem[] lcItems = check from LifecycleHistoryItem lcitem in lcStream
                select lcitem;
            check lcStream.close();
            return lcItems;
        } on fail var e {
            return e909612(e);
        }
    }
}

isolated function db_getSubscriptionsForAPI(string apiId) returns Subscription[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {

        sql:ParameterizedQuery query = `SELECT uuid as api_id FROM api WHERE uuid =${apiId}`;
        string|sql:Error result = dbClient->queryRow(query);
        if result is string {
            do {
                sql:ParameterizedQuery query1 = `SELECT 
                    SUBS.UUID AS subscriptionId,
                    APP.UUID AS applicationId,
                    APP.name AS name,
                    SUBS.TIER_ID AS usagePlan, 
                    SUBS.sub_status AS subscriptionStatus
                    FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
                    WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID AND API.UUID = ${apiId}`;
                stream<Subscriptions, sql:Error?> result1 = dbClient->query(query1);
                Subscription[] subsList = [];
                check from Subscriptions subitem in result1
                    do {
                        Subscription sub = {applicationInfo: {}, subscriptionId: "", subscriptionStatus: <"BLOCKED"|"PROD_ONLY_BLOCKED"|"UNBLOCKED"|"ON_HOLD"|"REJECTED"|"TIER_UPDATE_PENDING"|"DELETE_PENDING">"", usagePlan: ""};
                        sub.subscriptionId = subitem.subscriptionId;
                        sub.subscriptionStatus = <"BLOCKED"|"PROD_ONLY_BLOCKED"|"UNBLOCKED"|"ON_HOLD"|"REJECTED"|"TIER_UPDATE_PENDING"|"DELETE_PENDING">subitem.subscriptionStatus;
                        sub.applicationInfo.applicationId = subitem.applicationId;
                        sub.usagePlan = subitem.usagePlan;
                        sub.applicationInfo.name = subitem.name;
                        subsList.push(sub);
                    };
                return subsList;
            } on fail var e {
                return e909613(e);
            }
        } else {
            return e909614(result, apiId);
        }
    }
}

isolated function getSubscriptionByIdDAO(string subId) returns SubscriptionInternal|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        sql:ParameterizedQuery query = `SELECT 
        SUBS.UUID AS SUBSCRIPTION_ID, 
        API.API_NAME AS API_NAME, 
        API.API_VERSION AS API_VERSION, 
        API.API_TYPE AS API_TYPE, 
        API.ORGANIZATION AS ORGANIZATION, 
        APP.UUID AS APPLICATIONID, 
        SUBS.TIER_ID AS THROTTLINGPOLICY, 
        SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, 
        SUBS.SUB_STATUS AS STATUS, 
        SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, 
        SUBS.UUID AS UUID, 
        SUBS.CREATED_TIME AS CREATED_TIME, 
        SUBS.UPDATED_TIME AS UPDATED_TIME, 
        API.UUID AS APIID
        FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
        WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID AND SUBS.UUID =${subId}`;
        SubscriptionInternal|sql:Error result = dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            string message = "Subscription Not Found for provided ID";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 404);
        } else if result is SubscriptionInternal {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving Subscription";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function db_blockSubscription(string subscriptionId, string blockState) returns string|commons:APKError {
    postgresql:Client|error db_client = getConnection();
    if db_client is error {
        return e909601(db_client);
    } else {
        sql:ParameterizedQuery SUBSCRIPTION_BLOCK_Prefix = `UPDATE subscription set sub_status = `;
        sql:ParameterizedQuery values = `${blockState} where uuid = ${subscriptionId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(SUBSCRIPTION_BLOCK_Prefix, values);
        sql:ExecutionResult|sql:Error result = db_client->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return "blocked";
        } else {
            return e909615(result);
        }
    }
}

isolated function db_unblockSubscription(string subscriptionId) returns string|commons:APKError {
    postgresql:Client|error db_client = getConnection();
    if db_client is error {
        return e909601(db_client);
    } else {
        sql:ParameterizedQuery SUBSCRIPTION_UNBLOCK_Prefix = `UPDATE subscription set sub_status = 'UNBLOCKED'`;
        sql:ParameterizedQuery values = ` where uuid = ${subscriptionId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(SUBSCRIPTION_UNBLOCK_Prefix, values);
        sql:ExecutionResult|sql:Error result = db_client->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return "Unblocked";
        } else {
            return e909615(result);
        }
    }
}

isolated function db_getAPI(string apiId) returns API|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_API_Prefix = `SELECT UUID AS ID,
        API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,string_to_array(SDK::text,',')::text[] AS SDK,string_to_array(API_TIER::text,',') AS POLICIES, ARTIFACT as ARTIFACT
        FROM API where UUID = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_API_Prefix, values);

        API|sql:Error result = db_Client->queryRow(sqlQuery);

        if result is API {
            return result;
        } else {
            return e909616(result);
        }
    }
}

isolated function db_getAPIDefinition(string apiId) returns APIDefinition|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        sql:ParameterizedQuery query = `SELECT encode(API_DEFINITION, 'escape')::text AS schemaDefinition, MEDIA_TYPE as type
        FROM API_ARTIFACT WHERE API_UUID =${apiId}`;
        APIDefinition|sql:Error result = dbClient->queryRow(query);
        if result is sql:NoRowsError {
            return e909602();
        } else if result is APIDefinition {
            return result;
        } else {
            return e909617(result);
        }
    }
}

isolated function db_updateAPI(string apiId, ModifiableAPI payload) returns API|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        postgresql:JsonBinaryValue sdk = new (payload.sdk.toJson());
        postgresql:JsonBinaryValue categories = new (payload.categories.toJson());
        postgresql:JsonBinaryValue businessPlans = new (payload.policies.toJson());
        sql:ParameterizedQuery UPDATE_API_Suffix = `UPDATE api SET`;
        sql:ParameterizedQuery values = ` status= ${payload.state}, sdk = ${sdk},
        categories = ${categories}, api_tier=${businessPlans} WHERE uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_Suffix, values);

        sql:ExecutionResult|sql:Error result = dbClient->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return db_getAPI(apiId);
        } else {
            return e909618(result);
        }
    }
}

isolated function getAPICategoriesDAO(string org) returns APICategory[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID as ID, NAME, DESCRIPTION 
            FROM API_CATEGORIES WHERE ORGANIZATION =${org} ORDER BY NAME`;
            stream<APICategory, sql:Error?> apiCategoryStream = dbClient->query(query);
            APICategory[] apiCategoryList = check from APICategory apiCategory in apiCategoryStream
                select apiCategory;
            check apiCategoryStream.close();
            return apiCategoryList;
        } on fail var e {
            return e909619(e);
        }
    }
}

isolated function getAPIsByQueryDAO(string payload, string org, string[] groups) returns APIInfo[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query;
            if (groups is string[] && groups.length() > 0) {
                query = sql:queryConcat(`SELECT DISTINCT UUID AS ID,
                    API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,
                    ARTIFACT as ARTIFACT FROM API JOIN JSONB_EACH_TEXT(ARTIFACT) e ON true
                    WHERE e.value LIKE ${payload} AND ORGANIZATION = ${org} AND
                    ARTIFACT->>'accessControl' = 'NONE' OR ( ARTIFACT->>'accessControl' = 'RESTRICTED' AND
                    (
                     SELECT COUNT(DISTINCT groups)
                     FROM jsonb_array_elements_text(api.artifact->'accessControlGroups') AS groups
                     WHERE groups IN (`, sql:arrayFlattenQuery(groups), `)
                    ) >= 1 )`
                    );
            } else {
                query = `SELECT DISTINCT UUID AS ID,
                    API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,
                    ARTIFACT as ARTIFACT FROM API JOIN JSONB_EACH_TEXT(ARTIFACT) e ON true
                    WHERE e.value LIKE ${payload} AND ORGANIZATION = ${org} AND ARTIFACT->>'accessControl' = 'NONE'`;
            }

            stream<APIInfo, sql:Error?> apisStream = dbClient->query(query);
            APIInfo[] apis = check from APIInfo api in apisStream
                select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            io:print(e);
            return e909607(e);
        }
    }
}

isolated function db_getResourceByResourceCategory(string apiId, int resourceCategoryId) returns Resource|NotFoundError|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery sqlQuery = `SELECT UUID AS resourceUUID, API_UUID AS apiUuid, RESOURCE_CATEGORY_ID AS resourceCategoryId, DATA_TYPE AS dataType,
        RESOURCE_CONTENT AS resourceContent,  RESOURCE_BINARY_VALUE AS resourceBinaryValue  
        FROM API_RESOURCES where API_UUID = ${apiId} AND RESOURCE_CATEGORY_ID = ${resourceCategoryId}`;
        Resource|sql:Error result = db_Client->queryRow(sqlQuery);

        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body: {code: 90915, message: "Thumbnail Not Found for provided API ID"}};
            return nfe;
        } else if result is Resource {
            return result;
        } else {
            return e909626(result);
        }
    }
}

isolated function db_getResourceByResourceId(string resourceId) returns Resource|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_RESOURCE_Prefix = `SELECT UUID AS resourceUUID, API_UUID AS apiUuid, RESOURCE_CATEGORY_ID AS resourceCategoryId, DATA_TYPE AS dataType,
        RESOURCE_CONTENT AS resourceContent,  RESOURCE_BINARY_VALUE AS resourceBinaryValue  
        FROM API_RESOURCES where UUID = `;
        sql:ParameterizedQuery values = `${resourceId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_RESOURCE_Prefix, values);
        Resource|sql:Error result = db_Client->queryRow(sqlQuery);
        if result is Resource {
            return result;
        } else {
            return e909626(result);
        }
    }
}

isolated function db_addResource(Resource resourceItem) returns Resource|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        time:Utc utc = time:utcNow();
        sql:ParameterizedQuery values = `${resourceItem.resourceUUID},
                                        ${resourceItem.apiUuid},
                                        ${resourceItem.resourceCategoryId},
                                        ${resourceItem.dataType},
                                       to_tsvector(${resourceItem.resourceContent}),
                                        bytea(${resourceItem.resourceBinaryValue}),
                                        'apkuser',
                                        ${utc},
                                        'apkuser',
                                        ${utc}
                                    )`;
        sql:ParameterizedQuery ADD_THUMBNAIL_Prefix = `INSERT INTO API_RESOURCES (UUID, API_UUID, RESOURCE_CATEGORY_ID, DATA_TYPE, RESOURCE_CONTENT, RESOURCE_BINARY_VALUE, CREATED_BY, CREATED_TIME, UPDATED_BY, LAST_UPDATED_TIME) VALUES (`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_THUMBNAIL_Prefix, values);
        sql:ExecutionResult|sql:Error result = dbClient->execute(sqlQuery);
        if result is sql:ExecutionResult {
            log:printDebug("Resource added successfully");
            return resourceItem;
        } else {
            return e909624(result);
        }
    }
}

isolated function db_updateResource(Resource resourceItem) returns Resource|commons:APKError {
    log:printInfo(resourceItem.toBalString());
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        time:Utc utc = time:utcNow();
        string user = "apkuser";
        sql:ParameterizedQuery UPDATE_RESOURCE_Suffix = `UPDATE API_RESOURCES SET`;
        sql:ParameterizedQuery values = ` API_UUID= ${resourceItem.apiUuid}, RESOURCE_CATEGORY_ID = ${resourceItem.resourceCategoryId}, DATA_TYPE = ${resourceItem.dataType}, 
        RESOURCE_CONTENT = to_tsvector(${resourceItem.resourceContent}),
        RESOURCE_BINARY_VALUE = bytea(${resourceItem.resourceBinaryValue}), UPDATED_BY =${user}, LAST_UPDATED_TIME =${utc} WHERE UUID = ${resourceItem.resourceUUID}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_RESOURCE_Suffix, values);
        sql:ExecutionResult|sql:Error result = dbClient->execute(sqlQuery);
        if result is sql:ExecutionResult {
            return resourceItem;
        } else {
            return e909625(result);
        }
    }
}

isolated function db_getResourceCategoryIdByCategoryType(string resourceType) returns int|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_RESOURCE_CATEGORY_Prefix = `SELECT RESOURCE_CATEGORY_ID FROM RESOURCE_CATEGORIES where RESOURCE_CATEGORY = `;
        sql:ParameterizedQuery values = `${resourceType}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_RESOURCE_CATEGORY_Prefix, values);
        int|sql:Error result = db_Client->queryRow(sqlQuery);
        if result is int {
            return result;
        } else {
            return e909626(result);
        }
    }
}

isolated function db_addDocumentMetaData(DocumentMetaData documentMetaData, string apiId) returns DocumentMetaData|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        time:Utc utc = time:utcNow();
        sql:ParameterizedQuery values = `${documentMetaData.documentId},
                                        ${documentMetaData.resourceId},
                                        ${apiId},
                                        ${documentMetaData.name},
                                        ${documentMetaData.summary},
                                        ${documentMetaData.documentType},
                                        ${documentMetaData.otherTypeName},
                                        ${documentMetaData.sourceUrl},
                                        ${documentMetaData.fileName},
                                        ${documentMetaData.sourceType},
                                        ${documentMetaData.visibility},
                                        'apkuser',
                                        ${utc},
                                        'apkuser',
                                        ${utc}
                                    )`;
        sql:ParameterizedQuery ADD_DOCUMENT_Prefix = `INSERT INTO API_DOC_META_DATA (UUID, RESOURCE_UUID, API_UUID, NAME, SUMMARY, TYPE, OTHER_TYPE_NAME, SOURCE_URL, FILE_NAME,  SOURCE_TYPE,
         VISIBILITY, CREATED_BY, CREATED_TIME, UPDATED_BY, LAST_UPDATED_TIME) VALUES (`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_DOCUMENT_Prefix, values);
        sql:ExecutionResult|sql:Error result = dbClient->execute(sqlQuery);
        if result is sql:ExecutionResult {
            log:printDebug("Resource added successfully");
            return documentMetaData;
        } else {
            return e909624(result);
        }
    }
}

isolated function db_updateDocumentMetaData(DocumentMetaData documentMetaData, string apiId) returns DocumentMetaData|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        time:Utc utc = time:utcNow();
        string user = "apkuser";
        sql:ParameterizedQuery UPDATE_RESOURCE_Suffix = `UPDATE API_DOC_META_DATA SET`;
        sql:ParameterizedQuery values = ` NAME= ${documentMetaData.name}, SUMMARY = ${documentMetaData.summary}, TYPE = ${documentMetaData.documentType}, 
        OTHER_TYPE_NAME = ${documentMetaData.otherTypeName}, SOURCE_URL = ${documentMetaData.sourceUrl}, FILE_NAME = ${documentMetaData.fileName}, SOURCE_TYPE = ${documentMetaData.sourceType},
        VISIBILITY = ${documentMetaData.visibility}, UPDATED_BY =${user}, LAST_UPDATED_TIME =${utc} WHERE UUID = ${documentMetaData.documentId} AND API_UUID = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_RESOURCE_Suffix, values);
        sql:ExecutionResult|sql:Error result = dbClient->execute(sqlQuery);
        if result is sql:ExecutionResult {
            return documentMetaData;
        } else {
            return e909625(result);
        }
    }
}

isolated function db_getResourceIdByDocumentId(string documentId) returns string|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_RESOURCE_ID_Prefix = `SELECT RESOURCE_UUID FROM API_DOC_META_DATA where UUID = `;
        sql:ParameterizedQuery values = `${documentId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_RESOURCE_ID_Prefix, values);
        string|sql:Error result = db_Client->queryRow(sqlQuery);
        if result is string {
            return result;
        } else {
            return e909626(result);
        }
    }
}

isolated function db_deleteDocumentMetaData(string documentId, string apiId) returns string|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery sqlQuery = `DELETE FROM API_DOC_META_DATA WHERE UUID = ${documentId} AND API_UUID = ${apiId}`;
        sql:ExecutionResult|sql:Error result = db_Client->execute(sqlQuery);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            return e909626(result);
        }
    }
}

isolated function db_getDocumentByDocumentId(string documentId, string apiId) returns DocumentMetaData|NotFoundError|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_DOCUMENT_Prefix = `SELECT UUID AS documentId, RESOURCE_UUID AS resourceId, NAME AS name, SUMMARY AS summary,
        TYPE AS documentType, OTHER_TYPE_NAME AS otherTypeName, SOURCE_URL AS sourceUrl, FILE_NAME AS fileName,
        SOURCE_TYPE AS sourceType, VISIBILITY AS visibility FROM API_DOC_META_DATA where UUID = `;
        sql:ParameterizedQuery values = `${documentId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_DOCUMENT_Prefix, values);
        DocumentMetaData|sql:Error result = db_Client->queryRow(sqlQuery);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body: {code: 90915, message: "Document Not Found for provided Document ID"}};
            return nfe;
        } else if result is DocumentMetaData {
            return result;
        } else {
            return e909629(result);
        }
    }
}

isolated function db_deleteResource(string resourceId) returns string|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery DELETE_DOCUMENT_Prefix = `DELETE FROM API_RESOURCES WHERE UUID = `;
        sql:ParameterizedQuery values = `${resourceId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(DELETE_DOCUMENT_Prefix, values);
        sql:ExecutionResult|sql:Error result = db_Client->execute(sqlQuery);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            return e909626(result);
        }
    }
}

isolated function db_getDocuments(string apiId) returns Document[]|commons:APKError {
    postgresql:Client|error db_Client = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        do {
            sql:ParameterizedQuery GET_DOCUMENTS_Prefix = `SELECT UUID AS documentId, NAME AS name, SUMMARY AS summary,
        TYPE AS documentType, OTHER_TYPE_NAME AS otherTypeName, SOURCE_URL AS sourceUrl, FILE_NAME AS fileName,
        SOURCE_TYPE AS sourceType, VISIBILITY AS visibility FROM API_DOC_META_DATA where API_UUID = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_DOCUMENTS_Prefix, values);
            stream<Document, sql:Error?> documentStream = db_Client->query(sqlQuery);
            Document[] documents = check from Document document in documentStream
                select document;
            check documentStream.close();
            return documents;
        } on fail var e {
            return e909630(e);
        }
    }
}
