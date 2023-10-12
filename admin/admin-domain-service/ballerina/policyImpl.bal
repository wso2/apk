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

import ballerina/uuid;
import wso2/apk_common_lib as commons;

isolated function addDenyPolicy(BlockingCondition body, commons:Organization org) returns BlockingCondition|commons:APKError {
    string policyId = uuid:createType1AsString();
    body.policyId = policyId;
    //Todo : need to validate each type
    match body.conditionType {
        "APPLICATION" => {
        }
        "API" => {
        }
        "IP" => {
        }
        "IPRANGE" => {
        }
        "USER" => {
        }
    }
    BlockingCondition|commons:APKError policy = addDenyPolicyDAO(body, org.uuid);
    return policy;
}

isolated function getDenyPolicyById(string policyId, commons:Organization org) returns BlockingCondition|commons:APKError {
    BlockingCondition|commons:APKError policy = getDenyPolicyByIdDAO(policyId, org.uuid);
    return policy;
}

isolated function getAllDenyPolicies(commons:Organization org) returns BlockingConditionList|commons:APKError {
    BlockingCondition[]|commons:APKError denyPolicies = getDenyPoliciesDAO(org.uuid);
    if denyPolicies is BlockingCondition[] {
        int count = denyPolicies.length();
        BlockingConditionList denyPoliciesList = {count: count, list: denyPolicies};
        return denyPoliciesList;
    } else {
       return denyPolicies;
    }
}

isolated function updateDenyPolicy(string policyId, BlockingConditionStatus status, commons:Organization org) returns BlockingCondition|commons:APKError {
    BlockingCondition|commons:APKError existingPolicy = getDenyPolicyByIdDAO(policyId, org.uuid);
    if existingPolicy !is BlockingCondition {
        return existingPolicy;
    } else {
        status.policyId = policyId;
    }
    string|commons:APKError response = updateDenyPolicyDAO(status);
    if response is commons:APKError{
        return response;
    }
    BlockingCondition|commons:APKError updatedPolicy = getDenyPolicyByIdDAO(policyId, org.uuid);
    if updatedPolicy is BlockingCondition {
        return updatedPolicy;
    } else {
        return updatedPolicy;
    }
}

isolated function removeDenyPolicy(string policyId, commons:Organization org) returns commons:APKError|string {
    commons:APKError|string status = deleteDenyPolicyDAO(policyId, org.uuid);
    return status;
}

