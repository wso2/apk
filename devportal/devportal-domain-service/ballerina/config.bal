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

type ThrottlingConfiguration record {
    BlockCondition blockCondition = {
        enabled: true
    };
    boolean enableUnlimitedTier = true;
    boolean enableHeaderConditions = false;
    boolean enableJWTClaimConditions = false;
    boolean enableQueryParamConditions = false;
    boolean enablePolicyDeployment = true;
};

type BlockCondition record {
    boolean enabled = true;
};

type DatasourceConfiguration record {
    string name = "jdbc/apkdb";
    string description;
    string url;
    string host;
    int port;
    string databaseName;
    string username;
    string password;
    int maxPoolSize = 50;
    int minIdle = 20;
    int maxLifeTime = 60000;
    int validationTimeout;
    boolean autoCommit = true;
    string testQuery;
    string driver;
};

type APKConfiguration record {
    ThrottlingConfiguration throttlingConfiguration;
    DatasourceConfiguration datasourceConfiguration;
    TokenIssuerConfiguration tokenIssuerConfiguration;
    KeyStores keyStores;
};

public type TokenIssuerConfiguration record {|
    string issuer = "https://apim.wso2.com/oauth2/token";
    string audience = "https://apim.wso2.com/oauth2/token";
    string keyId = "gateway_certificate_alias";
    decimal expTime = 3600;
|};


public type SDKConfiguration record {|
    string groupId = "org.wso2";
    string artifactId = "org.wso2.client.";
    string modelPackage = "org.wso2.client.model.";
    string apiPackage = "org.wso2.client.api.";
|};
public type KeyStores record{|
        commons:KeyStore tls;
        commons:KeyStore signing;
|};