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

import ballerina/log;
import ballerinax/postgresql;
import ballerina/sql;
import wso2/apk_common_lib as commons;

public isolated function addDenyPolicyDAO(BlockingCondition bc, string org) returns BlockingCondition|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `INSERT INTO BLOCK_CONDITION (TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID) 
        VALUES (${bc.conditionType},${bc.conditionValue},${bc.conditionStatus},${org},${bc.policyId})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return bc;
        } else {
            log:printError(result.toString());
            return e909402(result);
        }
    }
}

public isolated function getDenyPolicyByIdDAO(string policyId, string org) returns BlockingCondition|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `SELECT UUID as POLICYID, TYPE as CONDITIONTYPE, BLOCK_CONDITION as CONDITIONVALUE, ENABLED::BOOLEAN as CONDITIONSTATUS FROM BLOCK_CONDITION WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        BlockingCondition | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return e909431();
        } else if result is BlockingCondition {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printError(result.toString());
            return e909422(result);
        }
    }
}

public isolated function getDenyPoliciesDAO(string org) returns BlockingCondition[]|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID as POLICYID, TYPE as CONDITIONTYPE, BLOCK_CONDITION as CONDITIONVALUE, ENABLED::BOOLEAN as CONDITIONSTATUS FROM BLOCK_CONDITION WHERE ORGANIZATION =${org}`;
            stream<BlockingCondition, sql:Error?> denyPoliciesStream = dbClient->query(query);
            BlockingCondition[] denyPolicies = check from BlockingCondition denyPolicy in denyPoliciesStream select denyPolicy;
            check denyPoliciesStream.close();
            return denyPolicies;
        } on fail var e {
            return e909423(e);
        }
    }
}

public isolated function updateDenyPolicyDAO(BlockingConditionStatus status) returns string|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `UPDATE BLOCK_CONDITION SET ENABLED = ${status.conditionStatus} WHERE UUID = ${status.policyId}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "";
        } else { 
            log:printError(result.toString());
            return e909402(result); 
        }
    }
}

public isolated function deleteDenyPolicyDAO(string policyId, string org) returns string|commons:APKError {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `DELETE FROM BLOCK_CONDITION WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "";
        } else {
            log:printError(result.toString());
            return e909406(result);
        }
    }
}