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

type CommonRouteSpec record {
    ParentReference[] parentRefs?;
};

public type HTTPPathMatch record {
    string 'type;
    string value;

};

public type HTTPHeaderMatch record {
    string 'type;
    string name;
    string value;
};

public type HTTPQueryParamMatch record {
    string 'type;
    string name;
};

public type HTTPRouteMatch record {

    HTTPPathMatch path?;
    HTTPHeaderMatch headers?;
    HTTPQueryParamMatch queryParams?;
    string method?;

};

public type HTTPHeader record {
    string name;
    string value;
};

public type HTTPHeaderFilter record {
    HTTPHeader[] set?;
    HTTPHeader[] add?;
    string[] remove?;
};

public type BackendObjectReference record {
    string group;
    string kind;
    string name;
    string namespace?;
    int port?;
};

public type HTTPRequestMirrorFilter record {
    BackendObjectReference backendRef;
};

public type HTTPPathModifier record {
    string 'type;
    string replaceFullPath?;
    string replacePrefixMatch?;
};

public type HTTPRequestRedirectFilter record {
    string scheme?;
    string hostname?;
    HTTPPathModifier path?;
    string port?;
    int statusCode?;
};

public type HTTPURLRewriteFilter record {
    string hostname?;
    HTTPPathModifier path?;

};

public type LocalObjectReference record {
    string group;
    string kind;
    string name;
};

public type HTTPRouteFilter record {
    string 'type;
    HTTPHeaderFilter requestHeaderModifier?;
    HTTPHeaderFilter responseHeaderModifier?;
    HTTPRequestMirrorFilter requestMirror?;
    HTTPRequestRedirectFilter requestRedirect?;
    HTTPURLRewriteFilter urlRewrite?;
    LocalObjectReference extensionRef?;
};

public type BackendRef record {
    *BackendObjectReference;
    int weight?;
};

public type HTTPBackendRef record {
    *BackendRef;
    HTTPRouteFilter[] filters?;
};

public type HTTPRouteRule record {
    HTTPRouteMatch[] matches?;
    HTTPRouteFilter[] filters?;
    HTTPBackendRef[] backendRefs?;
};

public type HTTPRouteSpec record {
    *CommonRouteSpec;
    string[] hostnames?;
    HTTPRouteRule[] rules = [];
};

public type HTTPRoute record {|
    string apiVersion = "gateway.networking.k8s.io/v1beta1";
    string kind = "HTTPRoute";
    Metadata metadata;
    HTTPRouteSpec spec;
    anydata status = ();
|};

public type ParentReference record {|
    string group?;
    string kind?;
    string namespace?;
    string name?;
    string sectionName?;
    string port?;

|};

public type HTTPRouteList record {|
    string apiVersion = "gateway.networking.k8s.io/v1beta1";
    string kind = "HTTPRouteList";
    ListMeta metadata;
    HTTPRoute[] items;
|};
