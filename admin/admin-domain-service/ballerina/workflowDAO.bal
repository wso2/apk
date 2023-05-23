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

//This function is used to retrive the pending workflow requests 
// Using Workflow table
isolated function getWorkflowListDAO(string? workflowType) returns WorkflowInfo[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            WorkflowInfo[] workflowList = [];
            sql:ParameterizedQuery query = 
                `SELECT wf_reference as workflowReferenceId, wf_type as workflowType, wf_status as workflowStatus, wf_created_time as createdTime, wf_updated_time as updatedTime
                 FROM WORKFLOWS WHERE wf_status = 'CREATED' AND wf_type = ${workflowType};`;
            stream<WorkflowInfo, sql:Error?> workFlowStream = dbClient->query(query);
            check from WorkflowInfo workflow in workFlowStream do {
                workflowList.push(workflow);
            };
            return workflowList;
        } on fail var e {
            return e909400(e);
        }
    }
}

isolated function getWorkflowDAO(string workflowReferenceId, WorkflowInfo payload) returns WorkflowInfo|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `Update WORKFLOWS SET wf_status = 'COMPLETED' WHERE wf_reference = ${workflowReferenceId};`;
            sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
            if result is sql:ExecutionResult {
                return payload;
            } else {
                return e909400(result);
            }
        } on fail var e {
            return e909400(e);
        }
    }
}