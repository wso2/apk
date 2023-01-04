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

# This function used to connect API create service to database
#
# + body - API parameter
# + organization - organization
# + return - Return Value API | error
isolated function createAPI(APIBody body, string organization) returns API | error{
    transaction {
        API | error apiCr = db_createAPI(body, organization);
        if apiCr is API {
            API | error defCr = db_AddDefinition(body, organization);
            if defCr is API {
                string|error lcEveCr = db_AddLCEvent(body.apiProperties.id, "carbon.super");
                if lcEveCr is string {
                    check commit;
                } else {
                    rollback;
                    return error("Error while adding API LC event");
                }
            } else {
                rollback;
                return error("Error while adding API definition");
            }
        } else {
            rollback;
            return error("Error while adding API data", apiCr);
        }
        return apiCr;
    }    
}

# This function used to connect API update service to database
#
# + body - API parameter
# + apiId - API Id parameter
# + organization - organization
# + return - Return Value API | error
isolated function updateAPI(string apiId, APIBody body, string organization) returns API | error {
    API | error apiUp = db_updateAPI(apiId, body, organization);
    if apiUp is error {
        return error("Error while updating API data");
    }
    API | error defUp = db_updateDefinition(apiId, body);
    if defUp is error {
        return error("Error while updating API definition");
    }
    return apiUp;
}

# This function used to connect API update service to database
#
# + apiId - API Id parameter
# + return - Return Value string | error
isolated function deleteAPI(string apiId) returns string|error? {
    error?|string apiDel = db_deleteAPI(apiId);
    if apiDel is error {
        return error("Error while deleting API data");
    }
    error?|string defDel = db_deleteDefinition(apiId);
    if defDel is error {
        return error("Error while deleting API definition data");
    }
    return apiDel;
}

# This function used to connect API update service to database
#
# + apiId - API Id parameter
# + apiBody - ApiidDefinitionBody 
# + return - Return Value string | error
isolated function updateDefinition(APIDefinition apiBody, string apiId) returns APIDefinition|error? {
    APIDefinition | error apiUp = db_updateDefinitionbyId(apiId, apiBody);
    if apiUp is error {
        return error("Error while updating API definition data");
    }
    return apiUp;
}

# This function used to create artifact from API
#
# + apiID - API Id parameter
# + api - api object
# + return - Return Value json
isolated function createArtifact(string? apiID, API api) returns json {
    Artifact artifact = {
                    id: apiID,
                    apiName : api.name,
                    context : api.context,
                    'version : api.'version,
                    status: api.lifeCycleStatus,
                    providerName: api.provider
                    };
    json artifactJson = artifact;
    return artifactJson;
}
