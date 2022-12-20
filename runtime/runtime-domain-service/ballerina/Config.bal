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

# This Record contains the Runtime related configurations
#
# + serviceListingNamespaces - Namespaces for List Services
# + apiCreationNamespace - Namespace for API creation.  
# + tokenIssuerConfiguration - Token Issuer Configuration for APIKey.  
# + keyStores - KeyStore Configuration
# + k8sConfiguration - K8s Configuration (debug purpose only.)
public type RuntimeConfiguratation record {|

    (string[] & readonly) serviceListingNamespaces = [ALL_NAMESPACES];
    string apiCreationNamespace = CURRENT_NAMESPACE;
    (TokenIssuerConfiguration & readonly) tokenIssuerConfiguration = {};
    KeyStores keyStores;
    (K8sConfigurations & readonly) k8sConfiguration = {};
|};

public type TokenIssuerConfiguration record {|
    string issuer = "https://localhost:9443/oauth2/token";
    string audience = "https://localhost:9443/oauth2/token";
    string keyId = "gateway_certificate_alias";
    decimal expTime = 3600;
|};

public type KeyStores record {|
    KeyStore signing;
    KeyStore tls;
|};

public type KeyStore record {|
    string path;
    string keyPassword?;
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
