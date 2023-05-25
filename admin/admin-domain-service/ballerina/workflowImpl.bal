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


//This function return the pending workflow list
isolated function getWorkflowList(string? workflowType, commons:Organization organization, int 'limit, int offset, string? accept) returns WorkflowList|commons:APKError{
    WorkflowList workflowList = {};
    if(workflowType == "APPLICATION_CREATION") {
        ApplciationWorkflowDTO[]|commons:APKError appWorkflowList = getApplicationCreationWorkflowListDAO(workflowType, organization);
        if(appWorkflowList is ApplciationWorkflowDTO[]) {
            WorkflowInfo[] workFlowInfoList = [];
            foreach ApplciationWorkflowDTO appWorkflow in appWorkflowList {
                WorkflowInfo workFlowInfo = {};
                string[] applicationProperty = [];
                workFlowInfo.workflowReferenceId = appWorkflow.workflowReferenceId;
                workFlowInfo.workflowType = appWorkflow.workflowType;
                workFlowInfo.workflowStatus = appWorkflow.workflowStatus;
                workFlowInfo.createdTime = appWorkflow.createdTime;
                workFlowInfo.updatedTime = appWorkflow.updatedTime;
                applicationProperty.push(appWorkflow.applicationName);
                applicationProperty.push(appWorkflow.createdBy);
                workFlowInfo.workflowProperties = applicationProperty;
                workFlowInfoList.push(workFlowInfo);
            }
            workflowList.list = workFlowInfoList;
        }
    } else if(workflowType == "SUBSCRIPTION_CREATION") {
        SubscriptionWorkflowDTO[]|commons:APKError subWorkflowList = getSubscriptionCreationWorkflowListDAO(workflowType, organization);
        if(subWorkflowList is SubscriptionWorkflowDTO[]) {
            WorkflowInfo[] workFlowInfoList = [];
            foreach SubscriptionWorkflowDTO subWorkflow in subWorkflowList {
                WorkflowInfo workFlowInfo = {};
                string[] subscriptionProperty = [];
                workFlowInfo.workflowReferenceId = subWorkflow.workflowReferenceId;
                workFlowInfo.workflowType = subWorkflow.workflowType;
                workFlowInfo.workflowStatus = subWorkflow.workflowStatus;
                workFlowInfo.createdTime = subWorkflow.createdTime;
                workFlowInfo.updatedTime = subWorkflow.updatedTime;
                subscriptionProperty.push(subWorkflow.applicationName);
                subscriptionProperty.push(subWorkflow.apiName);
                subscriptionProperty.push(subWorkflow.createdBy);
                workFlowInfo.workflowProperties = subscriptionProperty;
                workFlowInfoList.push(workFlowInfo);
            }
            workflowList.list = workFlowInfoList;
        }
    } 
    return workflowList;
}

// This function approvel/reject workflow request
isolated function updateWorkflowStatus(string workflowReferenceId, WorkflowInfo payload, commons:Organization organization) returns OkWorkflowInfo|commons:APKError {
    OkWorkflowInfo okWorkflowInfo = {
        body: {
            workflowReferenceId: ""
        }
    };
    if(payload.workflowType == "APPLICATION_CREATION") {
        WorkflowInfo|commons:APKError workflowInfo = updateApplciationWorkflowStatusDAO(workflowReferenceId, payload, organization);
        if workflowInfo is WorkflowInfo {
            okWorkflowInfo = {
                    body: {
                        workflowReferenceId: workflowInfo.workflowReferenceId
                    }
                };
        } else {
            return e909400(workflowInfo);
        }  
    } else if (payload.workflowType == "SUBSCRIPTION_CREATION") {
        WorkflowInfo|commons:APKError workflowInfo = updateSubscriptionWorkflowStatusDAO(workflowReferenceId, payload, organization);
        if workflowInfo is WorkflowInfo {
            okWorkflowInfo = {
                    body: {
                        workflowReferenceId: workflowInfo.workflowReferenceId
                    }
                };
        } else {
            return e909400(workflowInfo);
        }  
    }
    return okWorkflowInfo;
}