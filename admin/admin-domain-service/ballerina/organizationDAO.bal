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

import ballerinax/postgresql;
import ballerina/sql;

isolated function addOrganizationDAO(Organization payload) returns Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `INSERT INTO ORGANIZATION(UUID, NAME, 
        DISPLAY_NAME) VALUES (${payload.id},${payload.name},
        ${payload.displayName})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return payload;
        } else { 
            string message = "Error while inserting organization data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = "500"); 
        }
    }
}

isolated function addOrganizationClaimMappingDAO(Organization payload) returns Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        foreach OrganizationClaim e in payload.claimList {
            sql:ParameterizedQuery query = `INSERT INTO ORGANIZATION_CLIAM_MAPPING(UUID, CLAIM_KEY, 
            CLAIM_VALUE) VALUES (${payload.id},${e.claimKey},
            ${e.claimValue})`;
            sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
            if result is sql:ExecutionResult {
                continue;
            } else { 
                string message = "Error while inserting organization claim data into Database";
                return error(message, result, message = message, description = message, code = 909000, statusCode = "500"); 
            }
        }
        return payload;
    }
}

isolated function validateOrganizationByNameDAO(string name) returns boolean|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } 
    sql:ParameterizedQuery query = `select exists(SELECT 1 FROM ORGANIZATION WHERE NAME = ${name})`;
    boolean | sql:Error result =  dbClient->queryRow(query);
    if result is boolean {
        return result;
    } else { 
        string message = "Error while validating organization name in Database";
        return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
    }
    
}

isolated function validateOrganizationByDisplayNameDAO(string displayname) returns boolean|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } 
    sql:ParameterizedQuery query = `select exists(SELECT 1 FROM ORGANIZATION WHERE DISPLAY_NAME = ${displayname})`;
    boolean | sql:Error result =  dbClient->queryRow(query);
    if result is boolean {
        return result;
    } else { 
        string message = "Error while validating organization display name in Database";
        return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
    }   
}

isolated function validateOrganizationById(string? id) returns boolean|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } 
    sql:ParameterizedQuery query = `select exists(SELECT 1 FROM ORGANIZATION WHERE UUID = ${id})`;
    boolean | sql:Error result =  dbClient->queryRow(query);
    if result is boolean {
        return result;
    } else { 
        string message = "Error while validating organization id in Database";
        return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
    }   
}

isolated function validateClaimKeys(OrganizationClaim[] claims) returns boolean|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } 
    foreach OrganizationClaim e in claims {
        sql:ParameterizedQuery query = `select exists(SELECT 1 FROM organization_cliam_mapping WHERE CLAIM_KEY = ${e.claimKey})`;
        boolean | sql:Error result =  dbClient->queryRow(query);
        if result is true {
            continue;
        } else if result is false {
            return false;
        } else { 
            string message = "Error while validating claim key in Database";
            return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
        }  
    } 
    return true;
}

isolated function updateOrganizationDAO(string id, Organization payload) returns Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `UPDATE ORGANIZATION SET NAME =${payload.name},
         DISPLAY_NAME = ${payload.displayName} WHERE UUID = ${id}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return payload;
        } else { 
            string message = "Error while updating organization data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = "500"); 
        }
    }
}

isolated function updateOrganizationClaimMappingDAO(string id, Organization payload) returns Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        foreach OrganizationClaim e in payload.claimList {
            sql:ParameterizedQuery query = `UPDATE ORGANIZATION_CLIAM_MAPPING SET CLAIM_VALUE=${e.claimValue} WHERE CLAIM_KEY=${e.claimKey}`;
            sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
            if result is sql:ExecutionResult {
                continue;
            } else { 
                string message = "Error while updating organization claim data into Database";
                return error(message, result, message = message, description = message, code = 909000, statusCode = "500"); 
            }
        }
        return payload;
    }
}
