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

isolated function getAPIByAPIId(string apiId, string organization) returns string?|API|error {
    string?|API|error api = getAPIByIdDAO(apiId, organization);
    return api;
}

isolated function getAPIList(int 'limit, int  offset, string? query, string organization) returns string?|APIList|error {
    API[]|error? apis = getAPIsDAO(organization);
    if apis is API[] {
        int count = apis.length();
        ApplicationList apisList = {count: count, list: apis};
        return apisList;
    } else {
        return apis;
    }
}

isolated function getAPIDefinition(string apiId, string organization) returns APIDefinition|NotFoundError|error {
    APIDefinition|NotFoundError|error apiDefinition = getAPIDefinitionDAO(apiId,organization);
    return apiDefinition;
}
