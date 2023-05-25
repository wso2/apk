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

public type ApplciationWorkflowDTO record {
    string workflowReferenceId?;
    # Type of the Workflow Request. It shows which type of request is it.
    string workflowType?;
    # Show the Status of the the workflow request whether it is approved or created.
    string workflowStatus?;
    string applicationName;
    string createdBy;
    # Time of the the workflow request created.
    string createdTime?;
    # Time of the the workflow request updated.
    string updatedTime?;
    # description is a message with basic details about the workflow request.
    string description?;
};

public type SubscriptionWorkflowDTO record {
    string workflowReferenceId?;
    # Type of the Workflow Request. It shows which type of request is it.
    string workflowType?;
    # Show the Status of the the workflow request whether it is approved or created.
    string workflowStatus?;
    string applicationName;
    string apiName;
    string createdBy;
    # Time of the the workflow request created.
    string createdTime?;
    # Time of the the workflow request updated.
    string updatedTime?;
    # description is a message with basic details about the workflow request.
    string description?;
};
