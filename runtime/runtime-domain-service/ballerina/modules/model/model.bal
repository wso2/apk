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

# Data Model for K8s API CR.
#
# + uuid - uuid of API artifact deployed in k8s.
# + apiDisplayName - displayname of API.
# + apiType - Type of API (HTTP,GRAPHQL,WEBSOCKET,ASYNC,etc).  
# + apiVersion - Version of API.  
# + context - Context of API.  
# + definitionFileRef - API Definition Reference.  
# + prodHTTPRouteRef - Production Endpoint Http route file Ref.  
# + sandHTTPRouteRef - Sandbox Endpoint Http route file Ref.   
# + namespace - Namespace of API Deployed.   
# + creationTimestamp - Created Time of API.  
# + k8sName - K8s Internal Name.
public type K8sAPI record {
    string uuid;
    string apiDisplayName;
    string apiType;
    string apiVersion;
    string context;
    string definitionFileRef;
    string prodHTTPRouteRef?;
    string sandHTTPRouteRef?;
    string namespace;
    string creationTimestamp;
    string k8sName;
};
