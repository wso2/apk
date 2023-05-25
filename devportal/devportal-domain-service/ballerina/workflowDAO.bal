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

import ballerinax/postgresql;
import ballerina/sql;
import ballerina/uuid;
import ballerina/log;
import ballerina/time;


// This function is used to check the workflow enabled or not
isolated function isApplicationWorkflowEnabled(string organization) returns boolean|error {
    boolean isWorkflowEnabled = false;
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    }
    do {
        sql:ParameterizedQuery query = `SELECT encode(WORKFLOWS, 'escape')::text
                FROM ORGANIZATION where ORGANIZATION.UUID = ${organization}`;
        string | sql:Error result =  dbClient->queryRow(query);
        if(result is sql:Error) {
            log:printError(result.message());
            string message = "Error while retrieving workflows";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
            
        } else {
            Internal_Organization organization_workflow = {
                workflows : result.length() > 0 ? check result.fromJsonStringWithType() : []
            };
            if(organization_workflow.length() > 0) {
                foreach WorkflowProperties i in organization_workflow.workflows {
                    if(i.name == "ApplicationCreation") {
                        isWorkflowEnabled = i.enabled;
                        break;
                    }
                }
            }
        }
        return isWorkflowEnabled;
    }
}


isolated function isSubsciptionWorkflowEnabled(string organization) returns boolean|error {
    boolean isWorkflowEnabled = false;
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    }
    do {
        sql:ParameterizedQuery query = `SELECT encode(WORKFLOWS, 'escape')::text
                FROM ORGANIZATION where ORGANIZATION.UUID = ${organization}`;
        string | sql:Error result =  dbClient->queryRow(query);
        if(result is sql:Error) {
            log:printError(result.message());
            string message = "Error while retrieving workflows";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
            
        } else {
            Internal_Organization organization_workflow = {
                workflows : result.length() > 0 ? check result.fromJsonStringWithType() : []
            };
            if(organization_workflow.length() > 0) {
                foreach WorkflowProperties i in organization_workflow.workflows {
                    if(i.name == "SubscriptionCreation") {
                        isWorkflowEnabled = i.enabled;
                        break;
                    }
                }
            }
        }
        return isWorkflowEnabled;
    }
}

isolated function addApplicationCreationWorkflow(string applicationID, string organization) returns string|error {
    postgresql:Client | error dbClient  = getConnection();
    string uuid = uuid:createType1AsString();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    }
    do {
        sql:ParameterizedQuery query = `INSERT INTO WORKFLOWS(uuid, wf_reference, wf_type, wf_status,
                wf_created_time, wf_updated_time, organization) VALUES (${uuid}, ${applicationID}, 'APPLICATION_CREATION', 
                'CREATED', ${time:utcNow()}, ${time:utcNow()} , ${organization})`;
        sql:ExecutionResult|sql:Error result = dbClient->execute(query);
        if(result is sql:Error) {
            string message = "Error while adding application creation workflow";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        } 
        return uuid;
    }   
}

public type WorkflowProperties record {
    string name;
    boolean enabled;
    string[] properties?;
};

public type Internal_Organization record {
    WorkflowProperties[] workflows;
};
