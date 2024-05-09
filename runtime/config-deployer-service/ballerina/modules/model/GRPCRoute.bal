//
// Copyright (c) 2024, WSO2 LLC. (http://www.wso2.com).
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

public type GRPCRouteSpec record {
    *CommonRouteSpec;
    string[] hostnames?;
    GRPCRouteRule[] rules = [];
};

public type GRPCRouteRule record {
    GRPCRouteMatch[] matches;
    GRPCRouteFilter[] filters?;
    GRPCBackendRef[] backendRefs?;
};
public type GRPCRouteMatch record {
    GRPCMethodMatch method;
    GRPCHeaderMatch[] headers?;
};

public type GRPCHeaderMatch record {
    string 'type;
    string name;
    string value;
};

public type GRPCMethodMatch record {
    string 'type;
    string 'service;
    string method;
};

public type GRPCRouteFilter record {
    string 'type;
    HTTPHeaderFilter requestHeaderModifier?;
    HTTPHeaderFilter responseHeaderModifier?;
    HTTPRequestMirrorFilter requestMirror?;
    LocalObjectReference extensionRef?;
};

public type GRPCBackendRef record {
    *BackendRef;
    GRPCRouteFilter[] filters?;
};

public type GRPCRoute record {|
    string apiVersion = "gateway.networking.k8s.io/v1alpha2";
    string kind = "GRPCRoute";
    Metadata metadata;
    GRPCRouteSpec spec;
|};


public type GRPCRouteList record {|
    string apiVersion = "gateway.networking.k8s.io/v1alpha2";
    string kind = "GRPCRouteList";
    ListMeta metadata;
    GRPCRoute[] items;
|};


