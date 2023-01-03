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

isolated function db_getAPIsDAO() returns API[]|error? {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        sql:ParameterizedQuery GET_API = `SELECT API_UUID AS ID, API_ID as APIID,
        API_PROVIDER as PROVIDER, API_NAME as NAME, API_VERSION as VERSION,CONTEXT, ORGANIZATION,STATUS, ARTIFACT as ARTIFACT
        FROM API `;
        stream<API, sql:Error?> apisStream = db_Client->query(GET_API);
        API[]? apis = check from API api in apisStream select api;
        check apisStream.close();
        return apis;
    }
}

isolated function db_changeLCState(string targetState, string apiId, string organization) returns string|error {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        string newState = actionToLCState(targetState);
        if newState.equalsIgnoreCaseAscii("any") {
            return error(" Invalid Lifecycle targetState"); 
        }
        sql:ParameterizedQuery UPDATE_API_LifeCycle_Prefix = `UPDATE api SET status = `;
        sql:ParameterizedQuery values = `${newState}
        WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_LifeCycle_Prefix, values);

        sql:ExecutionResult | sql:Error result = db_Client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return targetState;
        } else {
            return error("Error while updating LC state into Database" + result.message());  
        }
    }
}

isolated function db_getCurrentLCStatus(string apiId, string organization) returns string|error {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        sql:ParameterizedQuery GET_API_LifeCycle_Prefix = `SELECT status from api where api_uuid = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(GET_API_LifeCycle_Prefix, values);

        string | sql:Error result =  db_Client->queryRow(sqlQuery);
        
        if result is string {
            return result;
        } else {
            return error("Error while geting LC state from Database" + result.message());  
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
isolated function db_AddLCEvent(string? apiId, string? prev_state, string? new_state, string organization) returns string | error {
    postgresql:Client | error db_client  = getConnection();
    time:Utc utc = time:utcNow();
    if db_client is error {
        return error("Issue while conecting to databse", db_client);
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
            return error("Error while inserting data into Database" + result.message());  
        }
    }
}

isolated function db_getLCEventHistory(string apiId) returns LifecycleHistoryItem[]?|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection", dbClient);
    } else {
        sql:ParameterizedQuery query = `SELECT previous_state, new_state, user_id, event_date FROM api_lc_event WHERE api_id =${apiId}`;
        stream<LifecycleHistoryItem, sql:Error?> lcStream = dbClient->query(query);
        LifecycleHistoryItem[]? lcItems = check from LifecycleHistoryItem lcitem in lcStream select lcitem;
        check lcStream.close();
        return lcItems;
    }
}


isolated function db_getSubscriptionsForAPI(string apiId) returns Subscription[]|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection", dbClient);
    } else {
        sql:ParameterizedQuery query = `SELECT api_id FROM api WHERE api_uuid =${apiId}`;
        int | sql:Error result =  dbClient->queryRow(query);
        
        if result is int {
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
            
        } else {
            return error("Error while geting subscription infomation" + result.message());  
        }
    }
}


isolated function db_blockSubscription(string subscriptionId, string blockState) returns error|string{
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery SUBSCRIPTION_BLOCK_Prefix = `UPDATE subscription set sub_status = `; 
        sql:ParameterizedQuery values = `${blockState} where uuid = ${subscriptionId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(SUBSCRIPTION_BLOCK_Prefix, values);
        sql:ExecutionResult | sql:Error result =  db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return "blocked";
        } else {
            return error("Error while changing status of the subscription in the Database");  
        }
    }
}

isolated function db_unblockSubscription(string subscriptionId) returns error|string{
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery SUBSCRIPTION_UNBLOCK_Prefix = `UPDATE subscription set sub_status = 'UNBLOCKED'`; 
        sql:ParameterizedQuery values = ` where uuid = ${subscriptionId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(SUBSCRIPTION_UNBLOCK_Prefix, values);
        sql:ExecutionResult | sql:Error result =  db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return "Unblocked";
        } else {
            return error("Error while changing status of the subscription in the Database");  
        }
    }
}
