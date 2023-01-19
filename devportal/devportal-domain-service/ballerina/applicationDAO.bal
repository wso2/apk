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

isolated function addApplicationDAO(Application application,int subscriberId, string org) returns Application|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
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
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = "500");
        }
    }
}

isolated function getSubscriberIdDAO(string user, string org) returns int|NotFoundError|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `SELECT SUBSCRIBER_ID FROM SUBSCRIBER WHERE USER_ID =${user} AND ORGANIZATION =${org}`;
        int|sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90916, message: "Subscriber Id not found"}};
            return nfe;
        } else if result is int {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving Subscriber ID";
            return error(message, result, message = message, description = message, code = 909007, statusCode = "500");
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

isolated function getApplicationByIdDAO(string appId, string org) returns Application|APKError|NotFoundError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `SELECT NAME, APPLICATION_ID as ID, UUID as APPLICATIONID, DESCRIPTION, APPLICATION_TIER as THROTTLINGPOLICY, TOKEN_TYPE as TOKENTYPE, ORGANIZATION,
        APPLICATION_STATUS as STATUS FROM APPLICATION WHERE UUID =${appId} AND ORGANIZATION =${org}`;
        Application | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            NotFoundError nfe = {body:{code: 90916, message: "Application Not Found"}};
            return nfe;
        } else if result is Application {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            string message = "Error while retrieving Application";
            return error(message, result, message = message, description = message, code = 909007, statusCode = "500");
        }
    }
}

isolated function getApplicationsDAO(string org) returns Application[]|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT NAME, APPLICATION_ID as ID, UUID as APPLICATIONID, DESCRIPTION, APPLICATION_TIER as THROTTLINGPOLICY, TOKEN_TYPE as TOKENTYPE, ORGANIZATION,
            APPLICATION_STATUS as STATUS  FROM APPLICATION WHERE ORGANIZATION =${org}`;
            stream<Application, sql:Error?> applicationStream = dbClient->query(query);
            Application[] applications = check from Application application in applicationStream select application;
            check applicationStream.close();
            return applications;
        } on fail var e {
            string message = "Internal Error occured while retrieving Applications";
            return error(message, e, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function updateApplicationDAO(Application application,int subscriberId, string org) returns Application|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
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
            log:printError(result.toString());
            string message = "Error while updating data record in the Database";
            return error(message, result, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function deleteApplicationDAO(string appId, string org) returns APKError|string {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `DELETE FROM APPLICATION WHERE UUID = ${appId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            log:printError(result.toString());
            string message = "Error while deleting data record in the Database";
            return error(message, result, message = message, description = message, code = 909001, statusCode = "500"); 
        }
    }
}

