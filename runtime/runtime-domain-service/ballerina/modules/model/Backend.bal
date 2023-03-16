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
public type Backend record {
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "Backend";
    Metadata metadata;
    BackendSpec spec;
};

public type BackendSpec record {|
    BackendService[] services;
    string protocol;
    TLSConfig tls?;
    SecurityConfig[] security?;
|};

public type BackendService record {
    string host;
    int port;
};

public type BasicSecurityConfig record {
   string username;
   string password; 
};

public type SecurityConfig record {
    string 'type;
    BasicSecurityConfig basic;
};


public type RefConfig record {
    string key;
    string name;
};

public type TLSConfig record {
    string[] allowedSANs?;
    string certificateInline?;
    RefConfig configMapRef?;
    RefConfig secretRef?;

};
public type BackendList record {
    string apiVersion="dp.wso2.com/v1alpha1";
    string kind = "BackendList";
    Backend [] items;
    ListMeta metadata;
};