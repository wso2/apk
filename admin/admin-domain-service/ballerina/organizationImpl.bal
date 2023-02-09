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
    string | APKError validateOrganization = validateOrganizationByNameDAO(payload.name);
    if validateOrganization is "true" {
        string message = "Organization already exists by name";
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    string | APKError validateOrganizationbyDisplayName = validateOrganizationByDisplayNameDAO(payload.displayName);
    if validateOrganizationbyDisplayName is "true" {
        string message = "Organization already exists by name";
        return error(message, message = message, description = message, code = 90911, statusCode = "400");
    }
    string organizationId = uuid:createType1AsString();
    payload.id = organizationId;
    Organization|APKError organization = addOrganizationDAO(payload);
    if organization is Organization {
        CreatedOrganization createdOrganization = {body: organization};
        return createdOrganization;
    } else {
        return organization;
    } 
}