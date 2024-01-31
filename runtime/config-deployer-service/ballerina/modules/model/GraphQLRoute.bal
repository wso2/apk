//
// Copyright (c) 2024, WSO2 LLC. (GraphQL://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// GraphQL://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
public type GQLRoute record {|
    string apiVersion = "dp.wso2.com/v1alpha2";
    string kind = "GQLRoute";
    Metadata metadata;
    GQLRouteSpec spec;
|};

public type GQLRouteList record {|
    string apiVersion = "dp.wso2.com/v1alpha2";
    string kind = "GQLRouteList";
    ListMeta metadata;
    GQLRoute[] items;
|};

public type GQLRouteMatch record {
    string 'type; //TODO: make enum
    string path;
};

enum GQLType {
    QUERY,
    MUTATION
};

public type GQLRouteFilter record {
    LocalObjectReference extensionRef?;
};

public type GQLRouteRule record {
    GQLRouteMatch[] matches?;
    GQLRouteFilter[] filters?;
};

public type GQLRouteSpec record {
    *CommonRouteSpec;
    string[] hostnames?;
    GQLRouteRule[] rules = [];
    BackendRef[] backendRefs?;
};

