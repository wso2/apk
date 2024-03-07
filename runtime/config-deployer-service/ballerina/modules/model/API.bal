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

public type API record {
    string kind = "API";
    string apiVersion = "dp.wso2.com/v1alpha2";
    Metadata metadata;
    APISpec spec;
    APIStatus? status = ();
};

public type APISpec record {|
    string apiName;
    string apiType;
    string apiVersion;
    string basePath;
    string organization;
    boolean isDefaultVersion?;
    string definitionFileRef?;
    string environment?;
    string definitionPath?;
    EnvConfig[]|() production = ();
    EnvConfig[]|() sandbox = ();
    boolean systemAPI?;
    APIProperties[] apiProperties = [];
|};

public type APIProperties record {|
    string name;
    string value;
|};

public type DeploymentStatus record {
    boolean accepted;
    string[] events;
    string message;
    string status;
    string transitionTime;
};

public type APIStatus record {
    DeploymentStatus deploymentStatus;
};

public type EnvConfig record {
    string[] routeRefs;
};

public type APIList record {
    string apiVersion = "dp.wso2.com/v1alpha2";
    string kind = "APIList";
    API[] items;
    ListMeta metadata;
};
