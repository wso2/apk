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
    string apiVersion = "dp.wso2.com/v1alpha2";
    string kind = "Backend";
    Metadata metadata;
    BackendSpec spec;
};

public type BackendSpec record {|
    BackendService[] services;
    string basePath?;
    string protocol;
    TLSConfig tls?;
    SecurityConfig security?;
    CircuitBreaker circuitBreaker?;
    Timeout timeout?;
    Retry 'retry?;
|};

public type BackendService record {
    string host;
    int port;
};

public type BasicSecurityConfig record {
   SecretRefConfig secretRef;
};

public type SecurityConfig record {
    BasicSecurityConfig basic?;
    APIKeySecurityConfig apiKey?;
};

public type APIKeySecurityConfig record {
    string name;
    string 'in;
    ValueFrom valueFrom;
};

public type ValueFrom record {
    string name;
    string valueKey;
};

public type SecretRefConfig record {
   string name;
   string usernameKey = "username";
   string passwordKey = "password"; 
};

public type RefConfig record {
    string key;
    string name;
};

public type Timeout record {
    int downstreamRequestIdleTimeout?;
    int upstreamResponseTimeout?;
};

public type Retry record {
    int count?;
    int baseIntervalMillis?;
    int[] statusCodes?;
};

public type CircuitBreaker record {
    int maxConnectionPools?;
    int maxConnections?;
    int maxPendingRequests?;
    int maxRequests?;
    int maxRetries?;
};

public type TLSConfig record {
    string[] allowedSANs?;
    string certificateInline?;
    RefConfig configMapRef?;
    RefConfig secretRef?;

};
public type BackendList record {
    string apiVersion="dp.wso2.com/v1alpha2";
    string kind = "BackendList";
    Backend [] items;
    ListMeta metadata;
};