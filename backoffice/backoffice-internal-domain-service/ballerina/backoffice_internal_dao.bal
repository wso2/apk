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

import ballerina/sql;
import ballerinax/postgresql;
import ballerina/time;

# Add API details to the database 
#
# + apiBody - API Parameter
# + organization - organization
# + return - API | error
isolated function db_createAPI(APIBody apiBody, string organization) returns API | error {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        postgresql:JsonBinaryValue artifact = new (createArtifact(apiBody.apiProperties.id, apiBody.apiProperties));
        sql:ParameterizedQuery ADD_API_Suffix = `INSERT INTO api(api_uuid, api_name, api_version,context,api_provider,status,organization,artifact) VALUES (`;
        sql:ParameterizedQuery values = `${apiBody.apiProperties.id},
                                            ${apiBody.apiProperties.name}, 
                                            ${apiBody.apiProperties.'version}, 
                                            ${apiBody.apiProperties.context},
                                            ${apiBody.apiProperties.provider},
                                            ${apiBody.apiProperties.lifeCycleStatus}, 
                                            ${organization},
                                            ${artifact})`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_API_Suffix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return apiBody.apiProperties;
        } else {
            return error("Error while inserting data into Database", result);  
        }
    }
}

# Add API definition to the database 
#
# + apiBody - API Parameter
# + organization - organization
# + return - API | error
isolated function db_AddDefinition(APIBody apiBody, string organization) returns API | error {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery ADD_API_DEFINITION_Suffix = `INSERT INTO api_artifact(organization, api_uuid, api_definition,media_type) VALUES (`;
        sql:ParameterizedQuery values = `${organization},
                                        ${apiBody.apiProperties.id},
                                        ${apiBody.Definition.toString().toBytes()}, 
                                        ${apiBody.apiProperties.'type}
                                    )`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_API_DEFINITION_Suffix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return apiBody.apiProperties;
        } else {
            return error("Error while inserting data into Database");  
        }
    }
}


# Update API details to the database 
#
# + api - API Parameter
# + apiId - API Id parameter
# + organization - organization
# + return - API | error
isolated function db_updateAPI(string apiId, APIBody api, string organization) returns API | error {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery UPDATE_API_Suffix = `UPDATE api SET`;
        sql:ParameterizedQuery values = ` api_name = ${api.apiProperties.name}
        WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_Suffix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return api.apiProperties;
        } else {
            return error("Error while updating data into Database");  
        }
    }
}

# Update API details to the database 
#
# + api - API Parameter
# + apiId - API Id parameter
# + return - API | error
isolated function db_updateDefinition(string apiId, APIBody api) returns API | error {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery UPDATE_API_DEFINITION_Suffix = `UPDATE api_artifact SET`;
        sql:ParameterizedQuery values = ` api_definition = ${api.Definition.toString().toBytes()}
        WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_DEFINITION_Suffix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return api.apiProperties;
        } else {
            return error("Error while updating definition into Database");  
        }
    }
}

# Delete API details from the database 
#
# + apiId - API Id parameter
# + return - string | error
isolated function db_deleteAPI(string apiId) returns string | error? {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery DELETE_API_Suffix = `DELETE FROM api WHERE api_uuid = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(DELETE_API_Suffix, values);
        sql:ExecutionResult | sql:Error result =  db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            return error("Error while deleting api data record in the Database");  
        }
    }
}

# Delete API details from the database 
#
# + apiId - API Id parameter
# + return - string | error
isolated function db_deleteDefinition(string apiId) returns string | error? {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery DELETE_API_DEFINITION_Suffix = `DELETE FROM api_artifact WHERE api_uuid = `;
        sql:ParameterizedQuery values = `${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(DELETE_API_DEFINITION_Suffix, values);
        sql:ExecutionResult | sql:Error result =  db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            return error("Error while deleting definition record in the Database");  
        }
    }
}


# Update API details to the database 
#
# + api - API Parameter
# + apiId - API Id parameter
# + return - API | error
isolated function db_updateDefinitionbyId(string apiId, APIDefinition api) returns APIDefinition | error {
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery UPDATE_API_DEFINITION_Suffix = `UPDATE api_artifact SET`;
        sql:ParameterizedQuery values = ` api_definition = ${api.Definition.toString().toBytes()}
        WHERE api_uuid = ${apiId}`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(UPDATE_API_DEFINITION_Suffix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return api;
        } else {
            return error("Error while updating definition into Database");  
        }
    }
}

# Add LC event to the database 
#
# + apiId - API id Parameter
# + organization - organization
# + return - API | error
isolated function db_AddLCEvent(string? apiId, string organization) returns string | error {
    postgresql:Client | error db_client  = getConnection();
    time:Utc utc = time:utcNow();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        sql:ParameterizedQuery values = `${apiId},
                                        null, 
                                        'CREATED',
                                        'apkuser',
                                        ${organization},
                                        ${utc}
                                    )`;
        sql:ParameterizedQuery ADD_LC_EVENT_Prefix = `INSERT INTO api_lc_event (api_id,previous_state,new_state,user_id,organization,event_date) VALUES (`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_LC_EVENT_Prefix, values);

        sql:ExecutionResult | sql:Error result = db_client->execute(sqlQuery);
        
        if result is sql:ExecutionResult {
            return result.toString();
        } else {
            return error("Error while inserting data into Database");  
        }
    }
}