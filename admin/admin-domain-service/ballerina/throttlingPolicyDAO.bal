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
import ballerinax/java.jdbc;
import ballerina/sql;

public isolated function addApplicationUsagePlanDAO(ApplicationRatePlan atp) returns string?|ApplicationRatePlan|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `INSERT INTO APPLICATION_USAGE_PLAN (NAME, DISPLAY_NAME, 
        ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID) 
        VALUES (${atp.policyName},${atp.displayName},${org},${atp.description},${atp.defaultLimit.'type},
        ${atp.defaultLimit.requestCount?.requestCount},${atp.defaultLimit.requestCount?.unitTime},
        ${atp.defaultLimit.requestCount?.timeUnit},${atp.isDeployed},${atp.policyId})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return atp;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}

public isolated function getApplicationUsagePlanByIdDAO(string policyId) returns string?|ApplicationRatePlan|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT * FROM APPLICATION_USAGE_PLAN WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        ApplicationRatePlan | sql:Error result =  dbClient->queryRow(query);
        check dbClient.close();
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is ApplicationRatePlan {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Application Usage Plan");
        }
    }
}

public isolated function getApplicationUsagePlansDAO(string org) returns ApplicationRatePlan[]|error? {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM APPLICATION_USAGE_PLAN WHERE ORGANIZATION =${org}`;
        stream<ApplicationRatePlan, sql:Error?> usagePlanStream = dbClient->query(query);
        ApplicationRatePlan[]? usagePlans = check from ApplicationRatePlan usagePlan in usagePlanStream select usagePlan;
        check usagePlanStream.close();
        check dbClient.close();
        return usagePlans;
    }
}

public isolated function updateApplicationUsagePlanDAO(ApplicationRatePlan atp) returns string?|ApplicationRatePlan|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `UPDATE APPLICATION_USAGE_PLAN SET DISPLAY_NAME = ${atp.displayName},
         DESCRIPTION = ${atp.description}, QUOTA_TYPE = ${atp.defaultLimit.'type}, QUOTA = ${atp.defaultLimit.requestCount?.requestCount}, 
         UNIT_TIME = ${atp.defaultLimit.requestCount?.unitTime}, TIME_UNIT = ${atp.defaultLimit.requestCount?.timeUnit} 
         WHERE UUID = ${atp.policyId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return atp;
        } else {
            log:printDebug(result.toString());
            return error("Error while updating data record in the Database");  
        }
    }
}

public isolated function deleteApplicationUsagePlanDAO(string policyId) returns string?|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `DELETE FROM APPLICATION_USAGE_PLAN WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            log:printDebug(result.toString());
            return error("Error while deleting data record in the Database");  
        }
    }
}

public isolated function addBusinessPlanDAO(BusinessPlan stp) returns string?|BusinessPlan|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `INSERT INTO BUSINESS_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, 
        QUOTA_TYPE, QUOTA, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID, 
        RATE_LIMIT_COUNT,RATE_LIMIT_TIME_UNIT,STOP_ON_QUOTA_REACH,MAX_DEPTH, MAX_COMPLEXITY,
        BILLING_PLAN,MONETIZATION_PLAN,CONNECTIONS_COUNT) VALUES (${stp.policyName},${stp.displayName},${org},${stp.description},${stp.defaultLimit.'type},
        ${stp.defaultLimit.requestCount?.requestCount},${stp.defaultLimit.requestCount?.unitTime},${stp.defaultLimit.requestCount?.timeUnit},
        ${stp.isDeployed},${stp.policyId},${stp.rateLimitCount},${stp.rateLimitTimeUnit},${stp.stopOnQuotaReach},0,
        0,${stp.billingPlan},${stp.monetization?.monetizationPlan},0)`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            log:printDebug(result.toString());
            return stp;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}

public isolated function getBusinessPlanByIdDAO(string policyId) returns string?|BusinessPlan|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT * FROM BUSINESS_PLAN WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        BusinessPlan | sql:Error result =  dbClient->queryRow(query);
        check dbClient.close();
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is BusinessPlan {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Business Plan");
        }
    }
}

public isolated function getBusinessPlansDAO(string org) returns BusinessPlan[]|error? {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM BUSINESS_PLAN WHERE ORGANIZATION =${org}`;
        stream<BusinessPlan, sql:Error?> businessPlanStream = dbClient->query(query);
        BusinessPlan[]? businessPlans = check from BusinessPlan businessPlan in businessPlanStream select businessPlan;
        check businessPlanStream.close();
        check dbClient.close();
        return businessPlans;
    }
}

public isolated function updateBusinessPlanDAO(BusinessPlan stp) returns string?|BusinessPlan|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `UPDATE BUSINESS_PLAN SET DISPLAY_NAME = ${stp.displayName},
         DESCRIPTION = ${stp.description}, QUOTA_TYPE = ${stp.defaultLimit.'type}, QUOTA = ${stp.defaultLimit.requestCount?.requestCount}, 
         UNIT_TIME = ${stp.defaultLimit.requestCount?.unitTime}, TIME_UNIT = ${stp.defaultLimit.requestCount?.timeUnit},
         RATE_LIMIT_COUNT = ${stp.rateLimitCount} , RATE_LIMIT_TIME_UNIT = ${stp.rateLimitTimeUnit} ,STOP_ON_QUOTA_REACH = ${stp.stopOnQuotaReach},
         BILLING_PLAN = ${stp.billingPlan}, 
         MONETIZATION_PLAN = ${stp.monetization?.monetizationPlan}, CONNECTIONS_COUNT = 0  
         WHERE UUID = ${stp.policyId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return stp;
        } else {
            log:printDebug(result.toString());
            return error("Error while updating data record in the Database");  
        }
    }
}

public isolated function deleteBusinessPlanDAO(string policyId) returns string?|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `DELETE FROM BUSINESS_PLAN WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return ();
        } else {
            log:printDebug(result.toString());
            return error("Error while deleting data record in the Database");  
        }
    }
}

public isolated function addDenyPolicyDAO(BlockingCondition bc) returns string?|BlockingCondition|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `INSERT INTO BLOCK_CONDITION (TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID) 
        VALUES (${bc.conditionType},${bc.conditionValue},${bc.conditionStatus},${org},${bc.conditionId})`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return bc;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}

public isolated function getDenyPolicyByIdDAO(string policyId) returns string?|BlockingCondition|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT * FROM BLOCK_CONDITION WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        BlockingCondition | sql:Error result =  dbClient->queryRow(query);
        check dbClient.close();
        if result is sql:NoRowsError {
            log:printDebug(result.toString());
            return error("Not Found");
        } else if result is BlockingCondition {
            log:printDebug(result.toString());
            return result;
        } else {
            log:printDebug(result.toString());
            return error("Error while retrieving Deny Policy from DB");
        }
    }
}

public isolated function getDenyPoliciesDAO(string org) returns BlockingCondition[]|error? {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM BLOCK_CONDITION WHERE ORGANIZATION =${org}`;
        stream<BlockingCondition, sql:Error?> denyPoliciesStream = dbClient->query(query);
        BlockingCondition[]? denyPolicies = check from BlockingCondition denyPolicy in denyPoliciesStream select denyPolicy;
        check denyPoliciesStream.close();
        check dbClient.close();
        return denyPolicies;
    }
}

public isolated function updateDenyPolicyDAO(BlockingConditionStatus status) returns string?|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `UPDATE BLOCK_CONDITION SET ENABLED = ${status.conditionStatus} WHERE UUID = ${status.conditionId}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return status.conditionId;
        } else {
            log:printDebug(result.toString());
            return error("Error while inserting data into Database");  
        }
    }
}

public isolated function deleteDenyPolicyDAO(string policyId) returns string?|error {
    jdbc:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `DELETE FROM BLOCK_CONDITION WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        check dbClient.close();
        if result is sql:ExecutionResult {
            return ();
        } else {
            log:printDebug(result.toString());
            return error("Error while deleting data record in the Database");  
        }
    }
}