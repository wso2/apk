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

function db_getAPIsDAO() returns API[]|error? {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        stream<API, sql:Error?> apisStream = db_Client->query(GET_API);
        API[]? apis = check from API api in apisStream select api;
        check apisStream.close();
        return apis;
    }
}

function db_changeLCState(string targetState, string apiId, string organization) returns string|error {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
        string newState = actionToLCState(targetState);
        if newState.equalsIgnoreCaseAscii("any") {
            return error(" Invalid Lifecycle targetState"); 
        }
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

function db_getCurrentLCStatus(string apiId, string organization) returns string|error {
    postgresql:Client | error db_Client  = getConnection();
    if db_Client is error {
        return error("Error while retrieving connection", db_Client);
    } else {
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
public function db_AddLCEvent(string? apiId, string? prev_state, string? new_state, string organization) returns string | error {
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
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_LC_EVENT_Prefix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return result.toString();
        } else {
            return error("Error while inserting data into Database" + result.message());  
        }
    }
}

public function db_getLCEventHistory(string apiId) returns LifecycleHistoryItem[]?|error {
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
