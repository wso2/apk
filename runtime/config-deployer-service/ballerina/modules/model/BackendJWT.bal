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
# Description.
#
# + apiVersion - field description  
# + kind - field description  
# + metadata - field description  
# + spec - field description
public type BackendJWT record {
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "BackendJWT";
    Metadata metadata;
    BackendJWTSpec spec;
};

public type BackendJWTSpec record {|
    string encoding?;
    string header?;
    string signingAlgorithm?;
    int tokenTTL?;
    CustomClaims[] customClaims?;
|};

public type CustomClaims record {|
    string claim?;
    string 'type?;
    string value;
|};

public type BackendJWTList record {
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "BackendJWTList";
    BackendJWT[] items;
    ListMeta metadata;
};
