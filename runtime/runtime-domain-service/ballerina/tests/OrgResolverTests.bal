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
import wso2/apk_common_lib as commons;

@test:Config {dataProvider: organizationFromIdpClaimValueDataProvider}
function testRetrieveOrganizationFromIDPClaimValue(string orgclaimValue, commons:Organization? expectedOrg) {
    K8sBaseOrgResolver k8sBaseOrgResolver = new;
    map<anydata> claims = {};
    test:assertEquals(k8sBaseOrgResolver.retrieveOrganizationFromIDPClaimValue(claims, orgclaimValue), expectedOrg);
}

function organizationFromIdpClaimValueDataProvider() returns map<[string, commons:Organization?]> {
    commons:Organization org2 = {
        displayName: "org2",
        name: "org2",
        organizationClaimValue: "org2",
        uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114c",
        enabled: true,
        serviceListingNamespaces: ["*"]
    };
    map<[string, commons:Organization?]> data = {
        "1": ["org2", org2],
        "2": ["org5", ()]
    };
    return data;
}

@test:Config {dataProvider: retrieveOrganizationByNameDataProvider}
function testRetrieveOrganizationByName(string orgclaimValue, commons:Organization? expectedOrg) {
    K8sBaseOrgResolver k8sBaseOrgResolver = new;
    test:assertEquals(k8sBaseOrgResolver.retrieveOrganizationByName(orgclaimValue), expectedOrg);
}

function retrieveOrganizationByNameDataProvider() returns map<[string, commons:Organization?]> {
    commons:Organization org2 = {
        displayName: "org2",
        name: "org2",
        organizationClaimValue: "org2",
        uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114c",
        enabled: true,
        serviceListingNamespaces: ["*"]
    };
    map<[string, commons:Organization?]> data = {
        "1": ["org2", org2],
        "2": ["org5", ()]
    };
    return data;
}
