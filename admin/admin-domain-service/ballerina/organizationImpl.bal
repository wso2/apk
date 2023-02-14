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

import ballerina/uuid;

isolated function addOrganization(Organization payload) returns CreatedOrganization|APKError {
    boolean validateOrganization = check validateOrganizationByNameDAO(payload.name);
    if validateOrganization is true {
        string message = "Organization already exists by name:" + payload.name;
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    boolean validateOrganizationbyDisplayName = check validateOrganizationByDisplayNameDAO(payload.displayName);
    if validateOrganizationbyDisplayName is true {
        string message = "Organization already exists by displayName:" + payload.displayName;
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    string organizationId = uuid:createType1AsString();
    payload.id = organizationId;
    Organization|APKError organization = addOrganizationDAO(payload);
    if organization is Organization {
        Organization|APKError organization1 = addOrganizationClaimMappingDAO(payload);
        if organization1 is Organization{
            CreatedOrganization createdOrganization = {body: organization};
            return createdOrganization;
        } else {
            return organization1;
        } 
    } else {
        return organization;
    } 
}

isolated function updatedOrganization(string id, Organization payload) returns Organization|APKError {
    boolean validateOrganizationId = check validateOrganizationById(id);
    if validateOrganizationId is false {
        string message = "Organization ID not exist by:" + id;
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    
    boolean validateOrganizationClaimKey = check validateClaimKeys(payload.claimList);
    if validateOrganizationClaimKey is false {
        string message = "Organization claim key invalid";
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }

    Organization|APKError organization = updateOrganizationDAO(id, payload);
    if organization is Organization {
        Organization|APKError organization1 = updateOrganizationClaimMappingDAO(id, payload);
        if organization1 is Organization{
            return organization1;
        } else {
            return organization1;
        } 
    } else {
        return organization;
    } 
}

isolated function getAllOrganization() returns OrganizationList|APKError {
    Organization[]|APKError getOrgnizations = getAllOrganizationDAO();
    if getOrgnizations is Organization[] {
        int count = getOrgnizations.length();
        OrganizationList getOrgnizationsList = {count: count, list: getOrgnizations};
        return getOrgnizationsList;
    } else {
       return getOrgnizations;
    }
}

isolated function getOrganizationById(string id) returns Organization|APKError {
    boolean validateOrganizationId = check validateOrganizationById(id);
    if validateOrganizationId is false {
        string message = "Organization ID not exist by:" + id;
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    Organization|APKError organization = getOrganizationByIdDAO(id);
    return organization;
}

isolated function removeOrganization(string id) returns boolean|APKError {
    boolean validateOrganizationId = check validateOrganizationById(id);
    if validateOrganizationId is false {
        string message = "Organization ID not exist by:" + id;
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    boolean|APKError organization = removeOrganizationDAO(id);
    return organization;
}
