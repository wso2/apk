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

import ballerina/test;

Organization  organization = {
        id: "01234567-0123-0123-0123-012345678901",
        name: "Finance",
        displayName: "Finance",
        organizationClaimValue: "01234567-0123-0123-0123",
        production: [],
        sandbox: [],
        serviceNamespaces: [
        "string"
        ]
    };
string orgId = "";

@test:Config {}
function addOrganizationTest() {
    Organization|APKError response = addOrganization(organization);
    if response is Organization {
        
        orgId = response.id.toString();
        test:assertTrue(true,"Organization added successfully");
    } else if response is APKError {
        test:assertFail("Error occured while adding Organization");
    }
    
}

@test:Config {dependsOn: [addOrganizationTest]}
function updateOrganizationTest() {
    Organization  updateOrganization = {
        id: "01234567-0123-0123-0123-012345678901",
        name: "Finance",
        displayName: "Finance",
        organizationClaimValue: "01234567-0123-0123-0123",
        serviceNamespaces: [
        "string"
        ]
    };
    Organization|APKError response = updatedOrganization(orgId, updateOrganization);
    if response is Organization {
        test:assertTrue(true,"Organization updated successfully");
    } else if response is APKError {
        test:assertFail("Error occured while updating Organization");
    }
    
}


@test:Config {dependsOn: [updateOrganizationTest]}
function getOrganizationsTest() {
    OrganizationList|APKError response = getAllOrganization();
    if response is OrganizationList {
        test:assertTrue(true,"Organization list retrieved successfully");
    } else if response is APKError {
        test:assertFail("Error occured while retrieving Organization list");
    }
}

@test:Config {dependsOn: [getOrganizationsTest]}
function getOrganizationTest() {
    Organization|APKError response = getOrganizationById(orgId);
    if response is Organization {
        test:assertTrue(true,"Organization retrieved successfully");
    } else if response is APKError {
        test:assertFail("Error occured while retrieving Organization");
    }
}

@test:Config {dependsOn: [getOrganizationTest]}
function deleteOrganizationTest() {
    boolean|APKError response = removeOrganization(orgId);
    if response is boolean {
        test:assertTrue(true,"Organization deleted successfully");
    } else if response is APKError {
        test:assertFail("Error occured while deleting Organization");
    }
}
