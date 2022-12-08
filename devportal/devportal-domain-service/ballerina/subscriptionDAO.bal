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
import ballerinax/java.jdbc;
import ballerina/sql;

public function getBusinessPlanByNameDAO(string policyName) returns string?|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT UUID FROM BUSINESS_PLAN WHERE NAME =${policyName} AND ORGANIZATION =${org}`;
        string| sql:Error result =  dbClient->queryRow(query);
        check dbClient.close();
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is string {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Business Plan");
        }
    }
}

function addSubscriptionDAO(Subscription sub, string user, int apiId, int appId) returns string?|Subscription|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        // check existing subscriptions
        sql:ParameterizedQuery existingCheckQuery = `SELECT SUB_STATUS, SUBS_CREATE_STATE FROM SUBSCRIPTION 
        WHERE API_ID = ${apiId} AND APPLICATION_ID = ${appId}`;
        Subscription | sql:Error existingCheckResult =  dbClient->queryRow(existingCheckQuery);
        if existingCheckResult is sql:NoRowsError {
            log:printDebug(existingCheckResult.toString());
        } else if existingCheckResult is Subscription {
            log:printDebug(existingCheckResult.toString());
            return error("Subscription Already exists");
        } else {
            log:printDebug(existingCheckResult.toString());
            return error("Error while checking exisiting subscriptions");
        }

        // Insert into SUBSCRIPTION table
        sql:ParameterizedQuery query = `INSERT INTO SUBSCRIPTION (TIER_ID,API_ID,APPLICATION_ID,
        SUB_STATUS,CREATED_BY,UUID, TIER_ID_PENDING) 
        VALUES (${sub.throttlingPolicy},${apiId},${appId},
        ${sub.status},${user},${sub.subscriptionId},${sub.requestedThrottlingPolicy})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            log:printDebug(result.toString());
            return sub;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}
