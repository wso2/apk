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
public type BackendPolicy record {
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "BackendPolicy";
    Metadata metadata;
    BackendPolicySpec spec;
};

public type BackendPolicySpec record {|
    BackendConfigs default?;
    BackendConfigs override?;
    TargetRef targetRef;
|};

public type BackendConfigs record {
    string protocol;
    TLSConfig tls?;
};

public type TLSConfig record {
    string certificateInline;
};
public type BackendPolicyList record {
    string apiVersion="dp.wso2.com/v1alpha1";
    string kind = "BackendPolicyList";
    BackendPolicy[] items;
    ListMeta metadata;
};