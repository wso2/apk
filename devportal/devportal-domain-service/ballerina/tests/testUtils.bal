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

import ballerina/log;
import ballerinax/postgresql;
import ballerina/sql;
import wso2/apk_common_lib as commons;
import ballerina/time;

# Add API details to the database
#
# + apiBody - API Parameter
# + organization - organization
# + return - API | error
isolated function createAPIDAO(APIBody apiBody, string organization) returns API | commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        postgresql:JsonBinaryValue artifact = new (createArtifact(apiBody.apiProperties.id, apiBody.apiProperties));
        sql:ParameterizedQuery ADD_API_Suffix = `INSERT INTO api(uuid, api_name, api_version,context,status,organization,artifact) VALUES (`;
        sql:ParameterizedQuery values = `${apiBody.apiProperties.id},
                                            ${apiBody.apiProperties.name},
                                            ${apiBody.apiProperties.'version},
                                            ${apiBody.apiProperties.context},
                                            ${apiBody.apiProperties.lifeCycleStatus},
                                            ${organization},
                                            ${artifact})`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_API_Suffix, values);

        sql:ExecutionResult | sql:Error result = dbClient->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return apiBody.apiProperties;
        } else {
            log:printDebug(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
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
                    status: api.lifeCycleStatus
                    };
    json artifactJson = artifact;
    return artifactJson;
}

# Add API definition to the database
#
# + apiBody - API Parameter
# + organization - organization
# + return - API | error
isolated function addDefinitionDAO(APIBody apiBody, string organization) returns API | commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery ADD_API_DEFINITION_Suffix = `INSERT INTO api_artifact(organization, api_uuid, api_definition,media_type) VALUES (`;
        sql:ParameterizedQuery values = `${organization},
                                        ${apiBody.apiProperties.id},
                                        ${apiBody.Definition.toString().toBytes()},
                                        ${apiBody.apiProperties.'type}
                                    )`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_API_DEFINITION_Suffix, values);

        sql:ExecutionResult | sql:Error result = dbClient->execute(sqlQuery);

        if result is sql:ExecutionResult {
            return apiBody.apiProperties;
        } else {
            log:printDebug(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

public isolated function addApplicationUsagePlanDAO(ApplicationRatePlan atp, string org) returns ApplicationRatePlan|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `INSERT INTO APPLICATION_USAGE_PLAN (NAME, DISPLAY_NAME, 
        ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID) 
        VALUES (${atp.planName},${atp.displayName},${org},${atp.description},${atp.defaultLimit.'type},
        ${atp.defaultLimit.requestCount?.requestCount},${atp.defaultLimit.requestCount?.unitTime},
        ${atp.defaultLimit.requestCount?.timeUnit},${atp.isDeployed},${atp.planId})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return atp;
        } else {
            log:printDebug(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

public isolated function addBusinessPlanDAO(BusinessPlan stp, string org) returns BusinessPlan|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `INSERT INTO BUSINESS_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, 
        QUOTA_TYPE, QUOTA, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID, 
        RATE_LIMIT_COUNT,RATE_LIMIT_TIME_UNIT,MAX_DEPTH, MAX_COMPLEXITY,
        BILLING_PLAN,CONNECTIONS_COUNT) VALUES (${stp.planName},${stp.displayName},${org},${stp.description},${stp.defaultLimit.'type},
        ${stp.defaultLimit.requestCount?.requestCount},${stp.defaultLimit.requestCount?.unitTime},${stp.defaultLimit.requestCount?.timeUnit},
        ${stp.isDeployed},${stp.planId},${stp.rateLimitCount},${stp.rateLimitTimeUnit},0,
        0,'FREE',0)`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            log:printDebug(result.toString());
            return stp;
        } else { 
            log:printError(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

isolated function addResourceDAO(Resource resourceItem) returns Resource|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        time:Utc utc = time:utcNow();
        sql:ParameterizedQuery values = `${resourceItem.resourceUUID},
                                        ${resourceItem.apiUuid},
                                        ${resourceItem.resourceCategoryId},
                                        ${resourceItem.dataType},
                                       to_tsvector(${resourceItem.resourceContent}),
                                        bytea(${resourceItem.resourceBinaryValue}),
                                        'apkuser',
                                        ${utc},
                                        'apkuser',
                                        ${utc}
                                    )`;
        sql:ParameterizedQuery ADD_THUMBNAIL_Prefix = `INSERT INTO API_RESOURCES (UUID, API_UUID, RESOURCE_CATEGORY_ID, DATA_TYPE, RESOURCE_CONTENT, RESOURCE_BINARY_VALUE, CREATED_BY, CREATED_TIME, UPDATED_BY, LAST_UPDATED_TIME) VALUES (`;
        sql:ParameterizedQuery sqlQuery = sql:queryConcat(ADD_THUMBNAIL_Prefix, values);
        sql:ExecutionResult | sql:Error result = dbClient->execute(sqlQuery);
        if result is sql:ExecutionResult {
            log:printDebug("Resource added successfully");
            return resourceItem;
        } else {
            log:printError(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

