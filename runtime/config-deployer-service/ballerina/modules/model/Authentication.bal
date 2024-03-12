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
public type Authentication record {
    string apiVersion = "dp.wso2.com/v1alpha2";
    string kind = "Authentication";
    Metadata metadata;
    AuthenticationSpec spec;
};

public type AuthenticationSpec record {
    AuthenticationData override?;
    AuthenticationData default?;
    TargetRef targetRef;

};

public type AuthenticationData record {
    AuthenticationExtensionType authTypes?;
    boolean disabled?;
};

public type AuthenticationExtensionType record {
    OAuth2Authentication oauth2?;
    APIKey[] apiKey = [];
    MutualSSL mtls?;
    JWTAuthentication jwt?;
};

public type MutualSSL record {
    string required;
    boolean disabled;
    RefConfig[] configMapRefs?;
    RefConfig[] secretRefs?;
    string[] certificatesInline?;
};

public type OAuth2Authentication record {
    string required;
    string header?;
    boolean sendTokenToUpstream?;
    boolean disabled;
};

public type JWTAuthentication record {
    string header?;
    boolean sendTokenToUpstream?;
    boolean disabled = true;
    string[] audience = [];
};

public type InternalKey record {
    string header?;
    string sendTokenToUpstream?;
};

public type APIKey record {
    string 'in?;
    string name?;
    boolean sendTokenToUpstream?;
};

public type AuthenticationList record {
    string apiVersion = "dp.wso2.com/v1alpha2";
    string kind = "AuthenticationList";
    Authentication[] items;
    ListMeta metadata;
};
