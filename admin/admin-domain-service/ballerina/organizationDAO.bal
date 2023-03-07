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

isolated function addOrganizationDAO(Internal_Organization payload) returns Internal_Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `INSERT INTO ORGANIZATION(UUID, NAME, 
        DISPLAY_NAME,STATUS) VALUES (${payload.id},${payload.name},
        ${payload.displayName},${payload.enabled})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult && result.affectedRowCount == 1 {
           return addOrganizationClaimMappingDAO(dbClient, payload);    
        } else { 
            string message = "Error while inserting organization data into Database";
            return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
        }
    }
}

isolated function addOrganizationClaimMappingDAO(postgresql:Client dbClient, Internal_Organization payload) returns Internal_Organization|APKError {
    foreach OrganizationClaim e in payload.claimList {
        sql:ParameterizedQuery query = `INSERT INTO ORGANIZATION_CLAIM_MAPPING(UUID, CLAIM_KEY, 
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

isolated function validateOrganizationByNameDAO(string name) returns boolean|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } 
    sql:ParameterizedQuery query = `SELECT * FROM ORGANIZATION WHERE NAME = ${name}`;
    Organization | sql:Error result =  dbClient->queryRow(query);
    if result is sql:NoRowsError {
        return false;
    } else if result is APICategory {
        return true;
    } else {
        string message = "Error while validating organization name in Database";
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

isolated function updateOrganizationDAO(string id, Internal_Organization payload) returns Internal_Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `UPDATE ORGANIZATION SET NAME =${payload.name},
         DISPLAY_NAME = ${payload.displayName}, STATUS=${payload.enabled} WHERE UUID = ${id}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult && result.affectedRowCount == 1 {
                return updateOrganizationClaimMappingDAO(dbClient, id, payload);
        } else { 
            string message = "Error while updating organization data into Database";
            return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
        }
    }
}

isolated function updateOrganizationClaimMappingDAO(postgresql:Client dbClient, string id, Internal_Organization payload) returns Internal_Organization|APKError {
    sql:ParameterizedQuery query = `DELETE FROM ORGANIZATION_CLAIM_MAPPING WHERE UUID = ${id}`;
    sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
    if result is sql:ExecutionResult {
        foreach OrganizationClaim e in payload.claimList {
            sql:ParameterizedQuery query1 = `INSERT INTO ORGANIZATION_CLAIM_MAPPING(UUID, CLAIM_KEY, 
            CLAIM_VALUE) VALUES (${id},${e.claimKey},
            ${e.claimValue})`;
            sql:ExecutionResult | sql:Error result1 =  dbClient->execute(query1);
            if result1 is sql:ExecutionResult {
                continue;
            } else { 
                string message = "Error while inserting organization claim data into Database";
                return error(message, result1, message = message, description = message, code = 909000, statusCode = "500"); 
            }
        }
    }
    return payload;
}

public isolated function getAllOrganizationDAO() returns Internal_Organization[]|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT ORGANIZATION.UUID as id, NAME as name, DISPLAY_NAME as displayName, claim_key as claimKey, claim_value as claimValue FROM ORGANIZATION, ORGANIZATION_CLAIM_MAPPING where ORGANIZATION.UUID = ORGANIZATION_CLAIM_MAPPING.UUID`;
            stream<Organizations, sql:Error?> orgStream = dbClient->query(query);
            map<Internal_Organization> organization = {};

            check from Organizations org in orgStream do {
                if organization.hasKey(org.id) {
                    OrganizationClaim claim = {claimKey: org.claimKey, claimValue: org.claimValue};
                    organization.get(org.id).claimList.push(claim);
                } else {
                    OrganizationClaim claim = {claimKey: org.claimKey, claimValue: org.claimValue};
                    Internal_Organization organizationData = {id: org.id, name: org.name, displayName: org.displayName, enabled: org.enabled,  claimList: [claim]};
                    organization[org.id] = organizationData;
                }
            };
            check orgStream.close();
            return organization.toArray();
        } on fail var e {
        	string message = "Internal Error occured while retrieving organization data from Database";
            return error(message, e, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function getOrganizationByIdDAO(string id) returns Internal_Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT ORGANIZATION.UUID as id, NAME as name, DISPLAY_NAME as displayName, claim_key as claimKey, 
                    claim_value as claimValue FROM ORGANIZATION, ORGANIZATION_CLAIM_MAPPING where ORGANIZATION.UUID = ORGANIZATION_CLAIM_MAPPING.UUID and ORGANIZATION.UUID =${id}`;
            stream<Organizations, sql:Error?> orgStream = dbClient->query(query);
            Internal_Organization organization1 = {
                id: "",
                name: "",
                displayName: "",
                enabled: true,
                claimList: []
            };
            check from Organizations org in orgStream do {
                if (organization1.id == "") {
                    organization1 = {
                        id:id,
                        name:org.name,
                        displayName:org.displayName,
                        enabled: org.enabled,
                        claimList:[{
                            claimKey:org.claimKey,
                            claimValue: org.claimValue
                        }]
                    };
                } else {
                    organization1.claimList.push({
                        claimKey:org.claimKey,
                        claimValue: org.claimValue
                    });
                }
            }; 
            if (organization1.id == "") {
                string message = "Organization not found";
                return error(message, message = message, description = message, code = 909000, statusCode = "404");
            } else {
                 return organization1;
            }

            } on fail var e {
        	string message = "Internal Error occured while retrieving organization data from Database";
            return error(message, e, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}

isolated function removeOrganizationDAO(string id) returns boolean|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        sql:ParameterizedQuery query = `DELETE FROM ORGANIZATION WHERE UUID = ${id}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return true;
        } else { 
            string message = "Error while deleting organization data from Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = "500"); 
        }
    }
}

isolated function getOrganizationByOrganizationClaimDAO(string claim) returns Internal_Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else { 
        sql:ParameterizedQuery query = `SELECT UUID as id FROM ORGANIZATION_CLAIM_MAPPING where claim_value =${claim}`;
        string | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            string message = "No organization found";
            return error(message, message = message, description = message, code = 909000, statusCode = "404"); 
        } else if result is string {
            return getOrganizationByIdDAO(result);
        } else { 
            string message = "Error while retrieving organization data from Database";
            return error(message, message = message, description = message, code = 909000, statusCode = "500"); 
        }
    }
    
}
isolated function getOrganizationByNameDAO(string name) returns Internal_Organization|APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = "500");
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT ORGANIZATION.UUID as id, NAME as name, DISPLAY_NAME as displayName, claim_key as claimKey, 
                    claim_value as claimValue FROM ORGANIZATION, ORGANIZATION_CLAIM_MAPPING where ORGANIZATION.UUID = ORGANIZATION_CLAIM_MAPPING.UUID and ORGANIZATION.NAME =${name}`;
            stream<Organizations, sql:Error?> orgStream = dbClient->query(query);
            Internal_Organization organization1 = {
                id: "",
                name: "",
                displayName: "",
                enabled: true,
                claimList: []
            };
            check from Organizations org in orgStream do {
                if (organization1.id == "") {
                    organization1 = {
                        id:org.id,
                        name:org.name,
                        displayName:org.displayName,
                        enabled: org.enabled,
                        claimList:[{
                            claimKey:org.claimKey,
                            claimValue: org.claimValue
                        }]
                    };
                } else {
                    organization1.claimList.push({
                        claimKey:org.claimKey,
                        claimValue: org.claimValue
                    });
                }
            }; 
            if (organization1.id == "") {
                string message = "Organization not found";
                return error(message, message = message, description = message, code = 909000, statusCode = "404");
            } else {
                 return organization1;
            }

            } on fail var e {
        	string message = "Internal Error occured while retrieving organization data from Database";
            return error(message, e, message = message, description = message, code = 909001, statusCode = "500");
        }
    }
}
