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


import wso2/apk_common_lib as commons;
import ballerinax/postgresql;
import ballerina/sql;
import ballerina/time;

//This function is used to retrive the pending workflow requests 
// Using Workflow table
isolated function getApplicationCreationWorkflowListDAO(string? workflowType, commons:Organization organization) returns ApplciationWorkflowDTO[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            ApplciationWorkflowDTO[] appWorkflowList = [];
            sql:ParameterizedQuery query = 
                `SELECT app.name as applicationName, app.created_by as createdBy, wf.wf_reference as workflowReferenceId, wf.wf_type as workflowType,
                 wf.wf_status as workflowStatus, wf.wf_created_time as createdTime, wf.wf_updated_time as updatedTime
                 FROM WORKFLOWS as wf, APPLICATION as app
                 WHERE wf.wf_status = 'CREATED' AND wf.wf_type = ${workflowType}
                 AND wf.wf_reference = app.uuid
                 AND wf.organization = ${organization.uuid};`;
            stream<ApplciationWorkflowDTO, sql:Error?> workFlowStream = dbClient->query(query);
            check from ApplciationWorkflowDTO appworkflow in workFlowStream do {
                appWorkflowList.push(appworkflow);
            };
            return appWorkflowList;
        } on fail var e {
            return e909400(e);
        }
    }
}

isolated function getSubscriptionCreationWorkflowListDAO(string? workflowType, commons:Organization organization) returns SubscriptionWorkflowDTO[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            SubscriptionWorkflowDTO[] subWorkflowList = [];
            sql:ParameterizedQuery query = 
                `SELECT api.api_name as apiName, app.name as applicationName, sub.created_by as createdBy, wf.wf_reference as workflowReferenceId, wf.wf_type as workflowType,
                 wf.wf_status as workflowStatus, wf.wf_created_time as createdTime, wf.wf_updated_time as updatedTime
                 FROM WORKFLOWS as wf, SUBSCRIPTION as sub, APPLICATION as app, API as api
                 WHERE wf.wf_status = 'CREATED' AND wf.wf_type = ${workflowType}
                 AND wf.wf_reference = sub.uuid
                 AND wf.organization = ${organization.uuid}
				 AND sub.application_uuid = app.uuid
				 AND sub.api_uuid = api.uuid;`;
            stream<SubscriptionWorkflowDTO, sql:Error?> workFlowStream = dbClient->query(query);
            check from SubscriptionWorkflowDTO subworkflow in workFlowStream do {
                subWorkflowList.push(subworkflow);
            };
            return subWorkflowList;
        } on fail var e {
            return e909400(e);
        }
    }
}

isolated function updateApplciationWorkflowStatusDAO(string workflowReferenceId, WorkflowInfo payload, commons:Organization organization) returns WorkflowInfo|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `Update WORKFLOWS SET wf_status = 'COMPLETED', wf_updated_time = ${time:utcNow()} 
            WHERE wf_reference = ${workflowReferenceId} AND organization = ${organization.uuid};`;
            sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
            if result is sql:ExecutionResult {
                sql:ParameterizedQuery query2 = `Update APPLICATION SET status = 'APPROVED' 
                WHERE uuid = ${workflowReferenceId} AND organization = ${organization};`;
                sql:ExecutionResult | sql:Error result2 =  dbClient->execute(query2);
                if result2 is sql:ExecutionResult {
                    return payload;
                } else {
                    return e909400(result2);
                }
            } else {
                return e909400(result);
            }
        } on fail var e {
            return e909400(e);
        }
    }
}

isolated function updateSubscriptionWorkflowStatusDAO(string workflowReferenceId, WorkflowInfo payload, commons:Organization organization) returns WorkflowInfo|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `Update WORKFLOWS SET wf_status = 'COMPLETED', wf_updated_time = ${time:utcNow()} 
            WHERE wf_reference = ${workflowReferenceId} AND organization = ${organization};`;
            sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
            if result is sql:ExecutionResult {
                sql:ParameterizedQuery query2 = `Update SUBSCRIPTION SET status = 'APPROVED' 
                WHERE uuid = ${workflowReferenceId} AND organization = ${organization.uuid};`;
                sql:ExecutionResult | sql:Error result2 =  dbClient->execute(query2);
                if result2 is sql:ExecutionResult {
                    return payload;
                } else {
                    return e909400(result2);
                }
            } else {
                return e909400(result);
            }
        } on fail var e {
            return e909400(e);
        }
    }
}