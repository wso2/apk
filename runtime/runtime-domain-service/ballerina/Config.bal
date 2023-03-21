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
import wso2/apk_common_lib as commons;

public type ControlPlaneConfiguration record {
    string serviceBaseURl;
    commons:KeyStore certificate?;
    boolean enableAuthentication=true;
    boolean enableHostNameVerification=false;
    Header[] headers = [];
};

public type Header record {|
    string name;
    string value;
|};
# This Record contains the Runtime related configurations
#
# + serviceListingNamespaces - Namespaces for List Services  
# + apiCreationNamespace - Namespace for API creation.  
# + tokenIssuerConfiguration - Token Issuer Configuration for APIKey.  
# + keyStores - KeyStore Configuration  
# + k8sConfiguration - K8s Configuration (debug purpose only.)  
# + vhost - Field Description  
# + idpConfiguration - IDP configuration for JWT generated from Enforcer.  
# + controlPlane - Field Description  
# + orgResolver - Field Description 
# + gatewayConfiguration - Gateway Configuration with name and listener name  
public type RuntimeConfiguratation record {|
    
    (string[] & readonly) serviceListingNamespaces = [ALL_NAMESPACES];
    string apiCreationNamespace = CURRENT_NAMESPACE;
    (TokenIssuerConfiguration & readonly) tokenIssuerConfiguration = {};
    KeyStores keyStores;
    (K8sConfigurations & readonly) k8sConfiguration = {};
    (Vhost[] & readonly) vhost = [{name:"Default",hosts:["gw.wso2.com"],'type:PRODUCTION_TYPE},{name:"Default",hosts:["sandbox.gw.wso2.com"],'type:SANDBOX_TYPE}];
    commons:IDPConfiguration idpConfiguration;
    ControlPlaneConfiguration controlPlane;
    (GatewayConfigurations & readonly) gatewayConfiguration = {};

    string orgResolver = "controlPlane"; // controlPlane, k8s
|};

public type Vhost record {|
string name;
string[] hosts;
string 'type;
|};


public type TokenIssuerConfiguration record {|
    string issuer = "https://apim.wso2.com/publisher";
    string audience = "https://localhost:9443/oauth2/token";
    string keyId = "gateway_certificate_alias";
    decimal expTime = 3600;
|};


public type K8sConfigurations record {|
    string host = "kubernetes.default";
    string serviceAccountPath = "/var/run/secrets/kubernetes.io/serviceaccount";
    decimal readTimeout = 5;
|};

isolated function getNameSpace(string namespace) returns string {
    if namespace == CURRENT_NAMESPACE {
        return currentNameSpace;
    } else {
        return namespace;
    }
}
public type KeyStores record {|
    commons:KeyStore signing;
    commons:KeyStore tls;
|};

public type GatewayConfigurations record {|
    string name = "default";
    string listenerName = "gatewaylistener";
|};