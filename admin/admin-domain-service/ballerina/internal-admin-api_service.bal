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
service /api/am/internal/v1 on internalAdminEp {
    resource function get organizations(string? organizationName, string? organizationClaimValue) returns OrganizationList|error|BadRequestError {
        if organizationName is string && organizationClaimValue is () {
            Internal_Organization organizationByNameDAO = check getOrganizationByNameDAO(organizationName);
            OrganizationList organizationList = {
                count: 1,
                list: [createOrganizationFromInternal(organizationByNameDAO)]
            };
            return organizationList;
        } else if organizationClaimValue is string && organizationName is () {
            Internal_Organization organizationByClaimDAO = check getOrganizationByOrganizationClaimDAO(organizationClaimValue);
            OrganizationList organizationList = {
                count: 1,
                list: [createOrganizationFromInternal(organizationByClaimDAO)]
            };
            return organizationList;
        } else if organizationName is string && organizationClaimValue is string {
            BadRequestError badRequestError = {body: {code: 900900, message: "Invalid query parameters. Only one of the query parameters can be provided."}};
            return badRequestError;
        }
        return check getAllOrganization();
    }
}
