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

import ballerina/io;
import ballerinax/postgresql;
import ballerina/sql;


public function getConnection() returns postgresql:Client | error {
    //Todo: Need to read database config from toml
    postgresql:Client|sql:Error dbClient = 
                                check new ("localhost", "sampath", "1qaz2wsx@Q", 
                                     "apklat1", 5432);
    return dbClient;
}

public function addApplicationUsagePlanDAO(ApplicationThrottlePolicy atp) returns string?|ApplicationThrottlePolicy {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `INSERT INTO APPLICATION_USAGE_PLAN (NAME, DISPLAY_NAME, 
        ORGANIZATION, DESCRIPTION, QUOTA_TYPE, QUOTA, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID) 
        VALUES (${atp.policyName},${atp.displayName},${org},${atp.description},${atp.defaultLimit.'type},
        ${atp.defaultLimit.requestCount?.requestCount},${atp.defaultLimit.requestCount?.unitTime},
        ${atp.defaultLimit.requestCount?.timeUnit},${atp.isDeployed},${atp.policyId})`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return atp;
        } else {
            io:println(result);
            return "Error while inserting data into Database";  
        }
    }
}

public function getApplicationUsagePlanByIdDAO(string policyId) returns string?|ApplicationThrottlePolicy {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT * FROM APPLICATION_USAGE_PLAN WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        io:println(query);
        ApplicationThrottlePolicy | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            io:println(result);
            return "Not Found";
        } else if result is ApplicationThrottlePolicy {
            io:println(result);
            return result;
        } else {
            io:println(result);
            return ();
        }
    }
}

public function getApplicationUsagePlansDAO(string org) returns ApplicationThrottlePolicy[]|error? {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM APPLICATION_USAGE_PLAN WHERE ORGANIZATION =${org}`;
        io:println(query);
        stream<ApplicationThrottlePolicy, sql:Error?> usagePlanStream = dbClient->query(query);
        ApplicationThrottlePolicy[]? usagePlans = check from ApplicationThrottlePolicy usagePlan in usagePlanStream select usagePlan;
        check usagePlanStream.close();
        return usagePlans;
    }
}

public function updateApplicationUsagePlanDAO(ApplicationThrottlePolicy atp) returns string?|ApplicationThrottlePolicy {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `UPDATE APPLICATION_USAGE_PLAN SET DISPLAY_NAME = ${atp.displayName},
         DESCRIPTION = ${atp.description}, QUOTA_TYPE = ${atp.defaultLimit.'type}, QUOTA = ${atp.defaultLimit.requestCount?.requestCount}, 
         UNIT_TIME = ${atp.defaultLimit.requestCount?.unitTime}, TIME_UNIT = ${atp.defaultLimit.requestCount?.timeUnit} 
         WHERE UUID = ${atp.policyId} AND ORGANIZATION = ${org}`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return atp;
        } else {
            io:println(result);
            return "Error while updating data record in the Database";  
        }
    }
}

public function deleteApplicationUsagePlanDAO(string policyId) returns string?|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `DELETE FROM APPLICATION_USAGE_PLAN WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            io:println(result);
            return error("Error while deleting data record in the Database");  
        }
    }
}

public function addBusinessPlanDAO(SubscriptionThrottlePolicy stp) returns string?|SubscriptionThrottlePolicy {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `INSERT INTO BUSINESS_PLAN (NAME, DISPLAY_NAME, ORGANIZATION, DESCRIPTION, 
        QUOTA_TYPE, QUOTA, UNIT_TIME, TIME_UNIT, IS_DEPLOYED, UUID, 
        RATE_LIMIT_COUNT,RATE_LIMIT_TIME_UNIT,STOP_ON_QUOTA_REACH,MAX_DEPTH, MAX_COMPLEXITY,
        BILLING_PLAN,MONETIZATION_PLAN,CONNECTIONS_COUNT) VALUES (${stp.policyName},${stp.displayName},${org},${stp.description},${stp.defaultLimit.'type},
        ${stp.defaultLimit.requestCount?.requestCount},${stp.defaultLimit.requestCount?.unitTime},${stp.defaultLimit.requestCount?.timeUnit},
        ${stp.isDeployed},${stp.policyId},${stp.rateLimitCount},${stp.rateLimitTimeUnit},${stp.stopOnQuotaReach},${stp.graphQLMaxDepth},
        ${stp.graphQLMaxComplexity},${stp.billingPlan},${stp.monetization?.monetizationPlan},${stp.subscriberCount})`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return stp;
        } else {
            io:println(result);
            return "Error while inserting data into Database";  
        }
    }
}

public function getBusinessPlanByIdDAO(string policyId) returns string?|SubscriptionThrottlePolicy {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT * FROM BUSINESS_PLAN WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        io:println(query);
        SubscriptionThrottlePolicy | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            io:println(result);
            return "Not Found";
        } else if result is SubscriptionThrottlePolicy {
            io:println(result);
            return result;
        } else {
            io:println(result);
            return ();
        }
    }
}

public function getBusinessPlansDAO(string org) returns SubscriptionThrottlePolicy[]|error? {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM BUSINESS_PLAN WHERE ORGANIZATION =${org}`;
        io:println(query);
        stream<SubscriptionThrottlePolicy, sql:Error?> businessPlanStream = dbClient->query(query);
        SubscriptionThrottlePolicy[]? businessPlans = check from SubscriptionThrottlePolicy businessPlan in businessPlanStream select businessPlan;
        check businessPlanStream.close();
        return businessPlans;
    }
}

public function updateBusinessPlanDAO(SubscriptionThrottlePolicy stp) returns string?|SubscriptionThrottlePolicy {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `UPDATE BUSINESS_PLAN SET DISPLAY_NAME = ${stp.displayName},
         DESCRIPTION = ${stp.description}, QUOTA_TYPE = ${stp.defaultLimit.'type}, QUOTA = ${stp.defaultLimit.requestCount?.requestCount}, 
         UNIT_TIME = ${stp.defaultLimit.requestCount?.unitTime}, TIME_UNIT = ${stp.defaultLimit.requestCount?.timeUnit},
         RATE_LIMIT_COUNT = ${stp.rateLimitCount} , RATE_LIMIT_TIME_UNIT = ${stp.rateLimitTimeUnit} ,STOP_ON_QUOTA_REACH = ${stp.stopOnQuotaReach},
         MAX_DEPTH = ${stp.graphQLMaxDepth}, MAX_COMPLEXITY = ${stp.graphQLMaxComplexity}, BILLING_PLAN = ${stp.billingPlan}, 
         MONETIZATION_PLAN = ${stp.monetization?.monetizationPlan}, CONNECTIONS_COUNT = ${stp.subscriberCount}  
         WHERE UUID = ${stp.policyId} AND ORGANIZATION = ${org}`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return stp;
        } else {
            io:println(result);
            return "Error while updating data record in the Database";  
        }
    }
}

public function deleteBusinessPlanDAO(string policyId) returns string?|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `DELETE FROM BUSINESS_PLAN WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            io:println(result);
            return error("Error while deleting data record in the Database");  
        }
    }
}

public function addDenyPolicyDAO(BlockingCondition bc) returns string?|BlockingCondition {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `INSERT INTO BLOCK_CONDITION (TYPE,BLOCK_CONDITION,ENABLED,ORGANIZATION,UUID) 
        VALUES (${bc.conditionType},${bc.conditionValue},${bc.conditionStatus},${org},${bc.conditionId})`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return bc;
        } else {
            io:println(result);
            return "Error while inserting data into Database";  
        }
    }
}

public function getDenyPolicyByIdDAO(string policyId) returns string?|BlockingCondition {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `SELECT * FROM BLOCK_CONDITION WHERE UUID =${policyId} AND ORGANIZATION =${org}`;
        io:println(query);
        BlockingCondition | sql:Error result =  dbClient->queryRow(query);
        if result is sql:NoRowsError {
            io:println(result);
            return "Not Found";
        } else if result is BlockingCondition {
            io:println(result);
            return result;
        } else {
            io:println(result);
            return ();
        }
    }
}

public function getDenyPoliciesDAO(string org) returns BlockingCondition[]|error? {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return error("Error while retrieving connection");
    } else {
        sql:ParameterizedQuery query = `SELECT * FROM BLOCK_CONDITION WHERE ORGANIZATION =${org}`;
        io:println(query);
        stream<BlockingCondition, sql:Error?> denyPoliciesStream = dbClient->query(query);
        BlockingCondition[]? denyPolicies = check from BlockingCondition denyPolicy in denyPoliciesStream select denyPolicy;
        check denyPoliciesStream.close();
        return denyPolicies;
    }
}

public function updateDenyPolicyDAO(BlockingConditionStatus status) returns string? {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        sql:ParameterizedQuery query = `UPDATE BLOCK_CONDITION SET ENABLED = ${status.conditionStatus} WHERE UUID = ${status.conditionId}`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return status.conditionId;
        } else {
            io:println(result);
            return "Error while inserting data into Database";  
        }
    }
}

public function deleteDenyPolicyDAO(string policyId) returns string?|error {
    postgresql:Client | error dbClient  = getConnection();
    if dbClient is error {
        return "Error while retrieving connection";
    } else {
        string org = "carbon.super";
        sql:ParameterizedQuery query = `DELETE FROM BLOCK_CONDITION WHERE UUID = ${policyId} AND ORGANIZATION = ${org}`;
        io:println(query);
        sql:ExecutionResult | sql:Error result =  dbClient->execute(query);
        if result is sql:ExecutionResult {
            return "deleted";
        } else {
            io:println(result);
            return error("Error while deleting data record in the Database");  
        }
    }
}