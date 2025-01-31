//
// Copyright (c) 2023, WSO2 LLC (http://www.wso2.com).
//
// WSO2 LLC licenses this file to you under the Apache License,
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
public type APIPolicy record {|
    string apiVersion = "dp.wso2.com/v1alpha4";
    string kind = "APIPolicy";
    Metadata metadata;
    APIPolicySpec spec;
|};

public type APIPolicySpec record {|
    APIPolicyData default?;
    APIPolicyData override?;
    TargetRef targetRef;
|};

public type APIPolicyData record {
    InterceptorReference[] requestInterceptors?;
    InterceptorReference[] responseInterceptors?;
    CORSPolicy cORSPolicy?;
    BackendJwtReference backendJwtPolicy?;
    boolean subscriptionValidation?;
    AIProviderReference aiProvider?;
    ModelBasedRoundRobin modelBasedRoundRobin?;
};

public type InterceptorReference record {
    string name;
};

public type BackendJwtReference record {
    string name?;
};

public type AIProviderReference record {
    string name?;
};

public type ModelBasedRoundRobin record {
    int onQuotaExceedSuspendDuration?;
    ModelWeight[] models;
};

public type ModelWeight record {
    string model;
    int weight;
};

public type APIPolicyList record {
    string apiVersion = "dp.wso2.com/v1alpha4";
    string kind = "APIPolicyList";
    ListMeta metadata;
    APIPolicy[] items;
};

public type CORSPolicy record {
    boolean enabled = true;
    boolean accessControlAllowCredentials = false;
    string[] accessControlAllowOrigins = [];
    string[] accessControlAllowHeaders = [];
    string[] accessControlAllowMethods = [];
    string[] accessControlExposeHeaders = [];
    int accessControlMaxAge?;
};
