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
public type RateLimitPolicy record {|
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "RateLimitPolicy";
    Metadata metadata;
    RateLimitSpec spec;
|};

public type RateLimitSpec record {|
    RateLimitData 'default?;
    RateLimitData override?;
    TargetRef targetRef;
|};

public type RateLimitData record {
    APIRateLimitDetails api?;
};

public type APIRateLimitDetails record {|
    int requestsPerUnit?;
    string unit?;
|};

public type RateLimitPolicyList record {
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "RateLimitPolicyList";
    ListMeta metadata;
    RateLimitPolicy[] items;
};

public type AIRateLimitPolicyList record {
    string apiVersion = "dp.wso2.com/v1alpha3";
    string kind = "AIRateLimitPolicyList";
    ListMeta metadata;
    AIRateLimitPolicy[] items;
};
