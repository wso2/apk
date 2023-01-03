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

isolated function addApplicationDAO(Application application,int subscriberId, string org) returns string?|Application|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `INSERT INTO APPLICATION (NAME, SUBSCRIBER_ID, APPLICATION_TIER,
        DESCRIPTION, APPLICATION_STATUS, GROUP_ID, CREATED_BY, CREATED_TIME, UPDATED_TIME,
        UUID, TOKEN_TYPE, ORGANIZATION) VALUES (${application.name},${subscriberId},${application.throttlingPolicy},
        ${application.description},${application.status},${application.groups},${application.owner},${application.createdTime},
        ${application.updatedTime},${application.applicationId},${application.tokenType},${org})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return application;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}

isolated function getSubscriberIdDAO(string user, string org) returns int?|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT SUBSCRIBER_ID FROM SUBSCRIBER WHERE USER_ID =${user} AND ORGANIZATION =${org}`;
        int|sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is int {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Subscriber ID");
        }
    }
}

isolated function getApplicationUsagePlanByNameDAO(string policyName, string org) returns string?|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT POLICY_ID FROM APPLICATION_USAGE_PLAN WHERE NAME =${policyName} AND ORGANIZATION =${org}`;
        string | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is string {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Application Usage Plan");
        }
    }
}

isolated function getApplicationByIdDAO(string appId, string org) returns string?|Application|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT NAME, APPLICATION_ID as ID, UUID as APPLICATIONID, DESCRIPTION, APPLICATION_TIER as THROTTLINGPOLICY, TOKEN_TYPE as TOKENTYPE, ORGANIZATION,
        APPLICATION_STATUS as STATUS FROM APPLICATION WHERE UUID =${appId} AND ORGANIZATION =${org}`;
        Application | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is Application {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Application");
        }
    }
}

isolated function getApplicationsDAO(string org) returns Application[]|error? {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT NAME, APPLICATION_ID as ID, UUID as APPLICATIONID, DESCRIPTION, APPLICATION_TIER as THROTTLINGPOLICY, TOKEN_TYPE as TOKENTYPE, ORGANIZATION,
        APPLICATION_STATUS as STATUS  FROM APPLICATION WHERE ORGANIZATION =${org}`;
        stream<Application, sql:Error?> applicationStream = dbClient->query(query);
        Application[]? applications = check from Application application in applicationStream select application;
        check applicationStream.close();
        return applications;
    }
}

isolated function updateApplicationDAO(Application application,int subscriberId, string org) returns string?|Application|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `UPDATE APPLICATION SET NAME = ${application.name},
         DESCRIPTION = ${application.description}, SUBSCRIBER_ID = ${subscriberId}, APPLICATION_TIER = ${application.throttlingPolicy}, 
         APPLICATION_STATUS = ${application.status}, GROUP_ID = ${application.groups},CREATED_BY = ${application.owner},
         CREATED_TIME = ${application.createdTime}, UPDATED_TIME = ${application.updatedTime}, TOKEN_TYPE = ${application.tokenType} 
         WHERE UUID = ${application.applicationId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return application;
        } else {
            log:printDebug(result.toString());
            return error("Error while updating data record in the Database");  
        }
    }
}

isolated function deleteApplicationDAO(string appId, string org) returns error?|string {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `DELETE FROM APPLICATION WHERE UUID = ${appId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            log:printDebug(result.toString());
            return error("Error while deleting data record in the Database");  
        }
    }
}

