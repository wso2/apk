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
# + return - Return Value API | error
public function createAPI(API body) returns API | error{
    API | error db = db_createAPI(body);
    return db;
}

# This function used to create artifact from API
#
# + apiID - API Id parameter
# + api - api object
# + return - Return Value json
function createArtifact(string apiID, API api) returns json {
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
