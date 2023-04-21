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

isolated function db_getAPIsDAO() returns API[]|commons:APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        do {
            sql:ParameterizedQuery GET_API = `SELECT API_UUID AS ID, API_ID as APIID,
            API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION, STATUS as STATE, string_to_array(SDK::text,',')::text[] AS SDK,string_to_array(API_TIER::text,',') AS POLICIES, ARTIFACT as ARTIFACT
            FROM API `;
            stream<API, sql:Error?> apisStream = db_Client->query(GET_API);
            API[] apis = check from API api in apisStream select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            return e909607(e);
        }
    }
}

isolated function db_changeLCState(string targetState, string apiId, string organization) returns string|commons:APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        string newState = actionToLCState(targetState);
        if newState.equalsIgnoreCaseAscii("any") {
            string message = "Invalid Lifecycle targetState";
            return error(message, message = message, description = message, code = 909002, statusCode = 400);
        }
        sql:ParameterizedQuery UPDATE_API_LifeCycle_Prefix = `UPDATE api SET status = `;
        sql:ParameterizedQuery values = `${newState}
        WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_LifeCycle_Prefix, values);

        sql:ExecutionResult | sql:Error result = db_Client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return targetState;
        } else {
            return e909608(result);  
        }
    }
}

isolated function db_getCurrentLCStatus(string apiId, string organization) returns string|commons:APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_API_LifeCycle_Prefix = `SELECT status from api where api_uuid = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_API_LifeCycle_Prefix, values);

        string | sql:Error result =  db_Client->queryRow(sqlQuery);
        
        if result is string {
            return result;
        } else {
            string message = "Error while geting LC state from Database";
            return error(message, result, message = message, description = message, code = 909002, statusCode = 500);
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
    postgresql:Client | error db_client  = getConnection();
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
        sql:ParameterizedQuery ADD_LC_EVENT_Prefix = `INSERT INTO api_lc_event (api_id,previous_state,new_state,user_id,organization,event_date) VALUES (`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_LC_EVENT_Prefix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return result.toString();
        } else {
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909002, statusCode = 500);
        }
    }
}

isolated function db_getLCEventHistory(string apiId) returns LifecycleHistoryItem[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do{
            sql:ParameterizedQuery query = `SELECT previous_state, new_state, user_id, event_date FROM api_lc_event WHERE api_id =${apiId}`;
            stream<LifecycleHistoryItem, sql:Error?> lcStream = dbClient->query(query);
            LifecycleHistoryItem[] lcItems = check from LifecycleHistoryItem lcitem in lcStream select lcitem;
            check lcStream.close();
            return lcItems;
        } on fail var e {
        	string message = "Internal Error occured while retrieving LC event History";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}


isolated function db_getSubscriptionsForAPI(string apiId) returns Subscription[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        
        sql:ParameterizedQuery query = `SELECT api_id FROM api WHERE api_uuid =${apiId}`;
        int | sql:Error result =  dbClient->queryRow(query);
        if result is int {
            do {
                sql:ParameterizedQuery query1 = `SELECT 
                    SUBS.SUBSCRIPTION_ID AS subscriptionId, 
                    APP.UUID AS applicationId,
                    APP.name AS name,
                    SUBS.TIER_ID AS usagePlan, 
                    SUBS.sub_status AS subscriptionStatus
                    FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
                    WHERE APP.APPLICATION_ID=SUBS.APPLICATION_ID AND API.API_ID = SUBS.API_ID AND API.API_UUID = ${apiId}`;
                stream<Subscriptions, sql:Error?> result1 =  dbClient->query(query1);
                Subscription[] subsList = [];
                check from Subscriptions subitem in result1 do {
                    Subscription sub = {applicationInfo: {},subscriptionId: "",subscriptionStatus: "",usagePlan: ""};
                    sub.subscriptionId =subitem.subscriptionId;
                    sub.subscriptionStatus = subitem.subscriptionStatus;
                    sub.applicationInfo.applicationId = subitem.applicationId;
                    sub.usagePlan = subitem.usagePlan;
                    sub.applicationInfo.name = subitem.name;
                    subsList.push(sub);
                };
                return subsList;
            } on fail var e {
                string message = "Internal Error while geting subscription infomation";
                return error(message, e, message = message, description = message, code = 909003, statusCode = 500);
            }
        } else {
            string message = "Internal Error while geting API for provided apiId " + apiId;
            return error(message, result, message = message, description = message, code = 909003, statusCode = 500);
        }
    }
}


isolated function db_blockSubscription(string subscriptionId, string blockState) returns string|commons:APKError {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return e909601(db_client);
    } else {
        sql:ParameterizedQuery SUBSCRIPTION_BLOCK_Prefix = `UPDATE subscription set sub_status = `; 
        sql:ParameterizedQuery values = `${blockState} where uuid = ${subscriptionId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(SUBSCRIPTION_BLOCK_Prefix, values);
        sql:ExecutionResult | sql:Error result =  db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return "blocked";
        } else {
            string message = "Error while changing status of the subscription in the Database";
            return error(message, result, message = message, description = message, code = 909003, statusCode = 500);
        }
    }
}

isolated function db_unblockSubscription(string subscriptionId) returns string|commons:APKError {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return e909601(db_client);
    } else {
        sql:ParameterizedQuery SUBSCRIPTION_UNBLOCK_Prefix = `UPDATE subscription set sub_status = 'UNBLOCKED'`; 
        sql:ParameterizedQuery values = ` where uuid = ${subscriptionId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(SUBSCRIPTION_UNBLOCK_Prefix, values);
        sql:ExecutionResult | sql:Error result =  db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return "Unblocked";
        } else {
            string message = "Error while changing status of the subscription in the Database";
            return error(message, result, message = message, description = message, code = 909003, statusCode = 500);
        }
    }
}

isolated function db_getAPI(string apiId) returns API|commons:APKError {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return e909601(db_Client);
    } else {
        sql:ParameterizedQuery GET_API_Prefix = `SELECT API_UUID AS ID, API_ID as APIID,
        API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS,string_to_array(SDK::text,',')::text[] AS SDK,string_to_array(API_TIER::text,',') AS POLICIES, ARTIFACT as ARTIFACT
        FROM API where API_UUID = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_API_Prefix, values);

        API | sql:Error result =  db_Client->queryRow(sqlQuery);
        
        if result is API {
            return result;
        } else {
            string message = "Error while retriving API";
            return error(message, result, message = message, description = message, code = 909004, statusCode = 500);
        }
    }
}

isolated function db_getAPIDefinition(string apiId) returns APIDefinition|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        sql:ParameterizedQuery query = `SELECT encode(API_DEFINITION, 'escape')::text AS schemaDefinition, MEDIA_TYPE as type
        FROM API_ARTIFACT WHERE API_UUID =${apiId}`;
        APIDefinition | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            return e909602();
        } else if result is APIDefinition {
            return result;
        } else {
            string message = "Internal Error while retrieving API Definition";
            return error(message, result, message = message, description = message, code = 909005, statusCode = 500);
        }
    }
}

isolated function db_updateAPI(string apiId, ModifiableAPI payload, string organization) returns API|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        postgresql:JsonBinaryValue sdk = new (payload.sdk.toJson());
        postgresql:JsonBinaryValue categories = new (payload.categories.toJson());
        postgresql:JsonBinaryValue businessPlans = new (payload.policies.toJson());
        sql:ParameterizedQuery UPDATE_API_Suffix = `UPDATE api SET`;
        sql:ParameterizedQuery values = ` status= ${payload.state}, sdk = ${sdk},
        categories = ${categories}, api_tier=${businessPlans} WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_Suffix, values);

        sql:ExecutionResult | sql:Error result = dbClient->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return db_getAPI(apiId);
        } else {
            string message = "Error while updating API data into Database";
            return error(message, result, message = message, description = message, code = 909005, statusCode = 500); 
        }
    }
}


isolated function getAPICategoriesDAO(string org) returns APICategory[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID as ID, NAME, DESCRIPTION 
            FROM API_CATEGORIES WHERE ORGANIZATION =${org} ORDER BY NAME`;
            stream<APICategory, sql:Error?> apiCategoryStream = dbClient->query(query);
            APICategory[] apiCategoryList = check from APICategory apiCategory in apiCategoryStream select apiCategory;
            check apiCategoryStream.close();
            return apiCategoryList;
        } on fail var e {
        	string message = "Internal Error occured while retrieving API Categories";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}

isolated function getAPIsByQueryDAO(string payload, string org) returns API[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT DISTINCT API_UUID AS ID, API_ID as APIID,
            API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS, 
            ARTIFACT as ARTIFACT FROM API JOIN JSONB_EACH_TEXT(ARTIFACT) e ON true 
            WHERE e.value LIKE ${payload}`;
            stream<API, sql:Error?> apisStream = dbClient->query(query);
            API[] apis = check from API api in apisStream select api;
            check apisStream.close();
            return apis;
        } on fail var e {
            io:print(e);
            return e909607(e)
        }
    }
}

public isolated function getBusinessPlansDAO(string org) returns BusinessPlan[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909601(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT NAME as PLANNAME, DISPLAY_NAME as DISPLAYNAME, DESCRIPTION, 
            UUID as PLANID, IS_DEPLOYED as ISDEPLOYED, 
            QUOTA_TYPE as DefaulLimitType, QUOTA , TIME_UNIT as TIMEUNIT, UNIT_TIME as 
            UNITTIME, RATE_LIMIT_COUNT as RATELIMITCOUNT, RATE_LIMIT_TIME_UNIT as RATELIMITTIMEUNIT FROM BUSINESS_PLAN WHERE ORGANIZATION =${org}`;
            stream<BusinessPlanDAO, sql:Error?> businessPlanStream = dbClient->query(query);
            BusinessPlanDAO[] businessPlansDAO = check from BusinessPlanDAO businessPlan in businessPlanStream select businessPlan;
            check businessPlanStream.close();
            BusinessPlan[] businessPlans =[];
            if businessPlansDAO is BusinessPlanDAO[] {
                foreach BusinessPlanDAO result in businessPlansDAO {
                    if result.defaulLimitType == "requestCount" {
                        BusinessPlan bp = {planName: result.planName, displayName: result.displayName, 
                        description: result.description, planId: result.planId, isDeployed: result.isDeployed, 
                        rateLimitCount: result.rateLimitCount, rateLimitTimeUnit: result.rateLimitTimeUnit,
                        defaultLimit: {'type: result.defaulLimitType, requestCount: 
                        {requestCount: result.quota, timeUnit: result.timeUnit, unitTime: result.unitTime}
                        }};
                        businessPlans.push(bp);
                    } else if result.defaulLimitType == "bandwidth" {
                        BusinessPlan bp = {planName: result.planName, displayName: result.displayName, 
                        description: result.description, planId: result.planId, isDeployed: result.isDeployed, 
                        rateLimitCount: result.rateLimitCount, rateLimitTimeUnit: result.rateLimitTimeUnit,
                        defaultLimit: {'type: result.defaulLimitType, bandwidth: 
                        {dataAmount: result.quota, dataUnit: <string>result.dataUnit, timeUnit: result.timeUnit, unitTime: result.unitTime}
                        }};
                        businessPlans.push(bp);
                    } else {
                        BusinessPlan bp = {planName: result.planName, displayName: result.displayName, 
                        description: result.description, planId: result.planId, isDeployed: result.isDeployed, 
                        rateLimitCount: result.rateLimitCount, rateLimitTimeUnit: result.rateLimitTimeUnit,
                        defaultLimit: {'type: result.defaulLimitType, eventCount: 
                        {eventCount:result.quota, timeUnit: result.timeUnit, unitTime: result.unitTime}
                        }};
                        businessPlans.push(bp);
                    }
                }
            }
            return businessPlans;
        } on fail var e {
        	string message = "Internal Error occured while retrieving Business Plans";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }
    }
}
