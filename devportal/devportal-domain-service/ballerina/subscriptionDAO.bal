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
import ballerinax/postgresql;
import ballerina/sql;
import wso2/apk_common_lib as commons;
#  DAO for GET Subscription plan
#
# + policyName -   Policy Name
# + return -      Policy ID
public isolated function getBusinessPlanByNameDAO(string policyName, string org) returns string|commons:APKError|NotFoundError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode =500);
    } else {
        sql:ParameterizedQuery query = `SELECT UUID FROM BUSINESS_PLAN WHERE NAME =${policyName} AND ORGANIZATION =${org}`;
        string| sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90915, message: "Business Plan Not Found for provided Plan Name"}};
            return nfe;
        } else if result is string {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving Business Plan";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function addSubscriptionDAO(Subscription sub, string user, string apiId, string appId) returns Subscription|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        // check existing subscriptions
        sql:ParameterizedQuery existingCheckQuery = `SELECT SUB_STATUS, SUBS_CREATE_STATE FROM SUBSCRIPTION 
        WHERE API_UUID = ${apiId} AND APPLICATION_UUID = ${appId}`;
        Subscription | sql:Error existingCheckResult =  dbClient->queryRow(existingCheckQuery);
        if existingCheckResult is sql:NoRowsError {
            log:printDebug(existingCheckResult.toString());
        } else if existingCheckResult is Subscription {
            log:printDebug(existingCheckResult.toString());
            string message = "Subscription Already exists";
            return error(message, message = message, description = message, code = 909000, statusCode = 500);
        } else {
            log:printDebug(existingCheckResult.toString());
            string message = "Error while checking exisiting subscriptions";
            return error(message, existingCheckResult, message = message, description = message, code = 909000, statusCode = 500);
        }
        var Tier_ID = "unlimited";

        // Insert into SUBSCRIPTION table
        sql:ParameterizedQuery query = `INSERT INTO SUBSCRIPTION (TIER_ID,API_UUID,APPLICATION_UUID,
        SUB_STATUS,CREATED_BY,UUID, TIER_ID_PENDING) 
        VALUES (${Tier_ID},${apiId},${appId},
        ${sub.status},${user},${sub.subscriptionId},${Tier_ID})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            log:printDebug(result.toString());
            return sub;
        } else {
            log:printDebug(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function getSubscriptionByIdDAO(string subId, string org) returns Subscription|NotFoundError|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
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
        SUBS.SUB_STATUS AS SUB_STATUS, 
        SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, 
        SUBS.CREATED_TIME AS CREATED_TIME, 
        SUBS.UPDATED_TIME AS UPDATED_TIME, 
        API.UUID AS APIID
        FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
        WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID AND SUBS.UUID =${subId}`;
        Subscription | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90915, message: "Subscription Not Found for provided ID"}};
            return nfe;
        } else if result is Subscription {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving Subscription";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function deleteSubscriptionDAO(string subId, string org) returns commons:APKError|string {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `DELETE FROM SUBSCRIPTION WHERE UUID = ${subId}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            log:printDebug(result.toString());
            string message = "Error while deleting data record in the Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function updateSubscriptionDAO(Subscription sub, string user, string apiId, string appId) returns Subscription|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        var Tier_ID = "unlimited";
        // Update Policy of a subscription in SUBSCRIPTION table
        sql:ParameterizedQuery query = ` UPDATE SUBSCRIPTION SET TIER_ID_PENDING = ${Tier_ID}
        , TIER_ID = ${Tier_ID} , SUB_STATUS = ${sub.status}
        WHERE UUID = ${sub.subscriptionId}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            log:printDebug(result.toString());
            return sub;
        } else {
            log:printError(result.toString());
            string message = "Error while updating data record in Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function getSubscriptionByAPIandAppIdDAO(string apiId, string appId, string org) returns Subscription|commons:APKError|NotFoundError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
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
        SUBS.SUB_STATUS AS SUB_STATUS, 
        SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, 
        SUBS.CREATED_TIME AS CREATED_TIME, 
        SUBS.UPDATED_TIME AS UPDATED_TIME, 
        API.UUID AS APIID
        FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
        WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID AND API.UUID =${apiId} AND APP.UUID=${appId}`;
        Subscription | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90916, message: "Subscription not found"}};
            return nfe;
        } else if result is Subscription {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving Subscription";
            return error(message, result, message = message, description = message, code = 909007, statusCode = 500);
        }
    }
}

isolated function getSubscriptionsByAPIIdDAO(string apiId, string org) returns Subscription[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT 
            SUBS.UUID AS SUBSCRIPTION_ID, 
            API.API_NAME AS API_NAME, 
            API.API_VERSION AS API_VERSION, 
            API.API_TYPE AS API_TYPE, 
            API.ORGANIZATION AS ORGANIZATION, 
            APP.UUID AS APPLICATIONID, 
            SUBS.TIER_ID AS THROTTLINGPOLICY, 
            SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, 
            SUBS.SUB_STATUS AS SUB_STATUS, 
            SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, 
            SUBS.CREATED_TIME AS CREATED_TIME, 
            SUBS.UPDATED_TIME AS UPDATED_TIME, 
            API.UUID AS APIID
            FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
            WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID AND API.UUID =${apiId}`;
            stream<Subscription, sql:Error?> subscriptionStream = dbClient->query(query);
            Subscription[] subscriptions = check from Subscription subscription in subscriptionStream select subscription;
            check subscriptionStream.close();
            return subscriptions;
        } on fail var e {
            string message = "Internal Error occured while retrieving Subscription By API Id";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }     
    }
}

isolated function getSubscriptionsByAPPIdDAO(string appId, string org) returns Subscription[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT 
            SUBS.UUID AS SUBSCRIPTION_ID, 
            API.API_NAME AS API_NAME, 
            API.API_VERSION AS API_VERSION, 
            API.API_TYPE AS API_TYPE, 
            API.ORGANIZATION AS ORGANIZATION, 
            APP.UUID AS APPLICATIONID, 
            SUBS.TIER_ID AS THROTTLINGPOLICY, 
            SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, 
            SUBS.SUB_STATUS AS SUB_STATUS, 
            SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, 
            SUBS.CREATED_TIME AS CREATED_TIME, 
            SUBS.UPDATED_TIME AS UPDATED_TIME, 
            API.UUID AS APIID
            FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
            WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID AND APP.UUID=${appId}`;
            stream<Subscription, sql:Error?> subscriptionStream = dbClient->query(query);
            Subscription[] subscriptions = check from Subscription subscription in subscriptionStream select subscription;
            check subscriptionStream.close();
            return subscriptions;
        } on fail var e {
            string message = "Internal Error occured while retrieving Subscription By Application Id";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }  
    }
}

isolated function getSubscriptionsList(string org) returns Subscription[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT 
            SUBS.UUID AS SUBSCRIPTION_ID, 
            API.API_NAME AS API_NAME, 
            API.API_VERSION AS API_VERSION, 
            API.API_TYPE AS API_TYPE, 
            API.ORGANIZATION AS ORGANIZATION, 
            APP.UUID AS APPLICATIONID, 
            SUBS.TIER_ID AS THROTTLINGPOLICY, 
            SUBS.TIER_ID_PENDING AS TIER_ID_PENDING, 
            SUBS.SUB_STATUS AS SUB_STATUS, 
            SUBS.SUBS_CREATE_STATE AS SUBS_CREATE_STATE, 
            SUBS.UUID AS UUID, 
            SUBS.CREATED_TIME AS CREATED_TIME, 
            SUBS.UPDATED_TIME AS UPDATED_TIME, 
            API.UUID AS APIID
            FROM SUBSCRIPTION SUBS, API API, APPLICATION APP 
            WHERE APP.UUID=SUBS.APPLICATION_UUID AND API.UUID = SUBS.API_UUID`;
            stream<Subscription, sql:Error?> subscriptionStream = dbClient->query(query);
            Subscription[] subscriptions = check from Subscription subscription in subscriptionStream select subscription;
            check subscriptionStream.close();
            return subscriptions;
        } on fail var e {
            string message = "Internal Error occured while retrieving Subscriptions";
            return error(message, e, message = message, description = message, code = 909001, statusCode = 500);
        }  
    }
}

