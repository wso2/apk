//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
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

import ballerina/http;
import ballerinax/postgresql;
import ballerina/sql;
import ballerina/log;
import ballerina/uuid;

final postgresql:Client|sql:Error dbClient2;

public isolated service class UserRegistrationInterceptor {
    *http:RequestInterceptor;

    public isolated function init(DatasourceConfiguration datasourceConfiguration) {
        dbClient2 =
        new (host = datasourceConfiguration.host,
        username = datasourceConfiguration.username,
        password = datasourceConfiguration.password,
        database = datasourceConfiguration.databaseName,
        port = datasourceConfiguration.port,
            connectionPool = {maxOpenConnections: datasourceConfiguration.maxPoolSize}
            );
        if dbClient2 is error {
            return log:printError("Error while connecting to database");
        }

    }
    isolated resource function 'default [string... path](http:RequestContext ctx, http:Request request, http:Caller caller) returns http:NextService|error? {
        if path[0] == "health" {
            return ctx.next();
        }
        if ctx.hasKey(VALIDATED_USER_CONTEXT) {
            UserContext userContext = check ctx.getWithType(VALIDATED_USER_CONTEXT, UserContext);
            string userId = check self.RegisterUser(userContext);
            userContext.userId = userId;
            ctx.set(VALIDATED_USER_CONTEXT, userContext.clone());
            return ctx.next();
        }
        return;
    }
    isolated function RegisterUser(UserContext userContext) returns string|APKError {
        string|APKError|() userId = self.retrieveUserFromIDPClaim(userContext.username);
        if userId is string {
            return userId;
        } else if userId is APKError {
                APKError apkError = error("Error while registering user", userId, code = 900900, description = "Internal Server Error.", statusCode = 500, message = "Internal Server Error.");
                return apkError;
        } else {
            string userIdUUID = uuid:createType1AsString();
            User payload = {
                IDPUserName: userContext.username,
                uuid: userIdUUID
            };
            User|APKError addedUser = check self.addUsertoDB(payload);
            if addedUser is User {
                log:printDebug("User added to the DB " + addedUser.toBalString());
                return addedUser.uuid;
            } else {
                APKError apkError = error("Error while adding user to the DB", userId, code = 900900, description = "Internal Server Error.", statusCode = 500, message = "Internal Server Error.");
                return apkError;
            }
        }
    }

    public isolated function retrieveUserFromIDPClaim(string userId) returns string|APKError|() {
        postgresql:Client|sql:Error dbClient = self.getConnection();
        if dbClient is sql:Error {
            return;
        } else {
            sql:ParameterizedQuery query = `SELECT UUID 
                FROM INTERNAL_USER WHERE IDP_USER_NAME = ${userId}`;
            string|sql:Error result = dbClient->queryRow(query);
            if result is sql:NoRowsError {
                log:printDebug("no rows found for User " + userId);
                return;
            } else if result is string {
                return result;
            } else {
                log:printError("Error while getting user " + userId, result);
                APKError apkError = error("Error while getting user", result, code = 900900, description = "Internal Server Error.", statusCode = 500, message = "Internal Server Error.");
                return apkError;
            }
        }
    }

    isolated function addUsertoDB(User payload) returns User|APKError {
        postgresql:Client|error dbClient = self.getConnection();
        if dbClient is error {
            string message = "Error while retrieving connection";
            return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
        } else {
            sql:ParameterizedQuery query = `INSERT INTO INTERNAL_USER(UUID, IDP_USER_NAME) VALUES (${payload.uuid},${payload.IDPUserName})`;
            sql:ExecutionResult | sql:Error result = dbClient->execute(query);
            if result is sql:ExecutionResult {
                return payload;
            } else {
                log:printError(result.toString());
                string message = "Error while inserting data into Database";
                return error(message, result, message = message, description = message, code = 909000, statusCode = 500); 
            }
        }
    }

    private isolated function getConnection() returns postgresql:Client|sql:Error {
        return dbClient2;
    }
}