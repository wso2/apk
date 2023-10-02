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

import ballerina/constraint;

// Defines a record to create artifact.
type Artifact record {|
    string? id;
    string apiName;
    string context;
    string 'version;
    string? status;
|};

public type APIBody record {
    API apiProperties;
    # Content of the definition
    record {} Definition;
};

public type Policy record {
    # Id of plan
    string planId?;
    # Name of plan
    @constraint:String {maxLength: 60, minLength: 1}
    string planName;
    # Display name of the policy
    @constraint:String {maxLength: 512}
    string displayName?;
    # Description of the policy
    @constraint:String {maxLength: 1024}
    string description?;
    # Indicates whether the policy is deployed successfully or not.
    boolean isDeployed = false;
    # Indicates the type of throttle policy
    string 'type?;
};

public type GraphQLQuery record {
    # Maximum Complexity of the GraphQL query
    int graphQLMaxComplexity?;
    # Maximum Depth of the GraphQL query
    int graphQLMaxDepth?;
};

public type ThrottleLimitBase record {
    # Unit of the time. Allowed values are "sec", "min", "hour", "day"
    string timeUnit;
    # Time limit that the throttling limit applies.
    int unitTime;
};

public type RequestCountLimit record {
    *ThrottleLimitBase;
    # Maximum number of requests allowed
    int requestCount;
};

public type ThrottleLimit record {
    # Type of the throttling limit. Allowed values are "REQUESTCOUNTLIMIT" and "BANDWIDTHLIMIT".
    # Please see schemas of "RequestCountLimit" and "BandwidthLimit" throttling limit types in
    # Definitions section.
    string 'type;
    RequestCountLimit requestCount?;
    BandwidthLimit bandwidth?;
    EventCountLimit eventCount?;
};
public type ApplicationRatePlan record {
    *Policy;
    ThrottleLimit defaultLimit;
};

public type CustomAttribute record {
    # Name of the custom attribute
    string name;
    # Value of the custom attribute
    string value;
};

public type BusinessPlanPermission record {
    string permissionType;
    string[] roles;
};

public type BusinessPlan record {
    *Policy;
    *GraphQLQuery;
    ThrottleLimit defaultLimit;
    # Burst control request count
    int rateLimitCount?;
    # Burst control time unit
    string rateLimitTimeUnit?;
    # Number of subscriptions allowed
    int subscriberCount?;
    # Custom attributes added to the Subscription Throttling Policy
    CustomAttribute[] customAttributes?;
    BusinessPlanPermission permissions?;
};

public type BandwidthLimit record {
    *ThrottleLimitBase;
    # Amount of data allowed to be transfered
    int dataAmount;
    # Unit of data allowed to be transfered. Allowed values are "KB", "MB" and "GB"
    string dataUnit;
};

public type EventCountLimit record {
    *ThrottleLimitBase;
    # Maximum number of events allowed
    int eventCount;
};