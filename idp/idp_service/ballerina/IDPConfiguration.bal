import ballerina/uuid;

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
# Description
#
# + loginUiUrl - Field Description  
# + user - Field Description  
# + signingKeyStore - Field Description  
# + publicKey - Field Description  
# + dataSource - Field Description  
# + fileBaseApp - Field Description  
# + tokenIssuerConfiguration - Field Description
public type IDPConfiguration record {|

    string loginUiUrl;

    Users[] user = [];
    KeyStoreConfiguration signingKeyStore;
    KeyStoreConfiguration publicKey;
    DatasourceConfiguration dataSource;
    FileBaseOAuthapps[] fileBaseApp = [];
    TokenIssuerConfiguration tokenIssuerConfiguration = {};
|};

# Description
#
# + username - User name  
# + password - Password of user.
# + organizations - organizations belongs to user.  
# + superAdmin - User mark as super_admin.
public type Users record {|
    string username;
    string password;
    string[] organizations = [];
    boolean superAdmin;
|};

public type KeyStoreConfiguration record {|
    string id = uuid:createType1AsString();
    string path;
    string keyPassword?;
|};

public type DatasourceConfiguration record {
    string name = "jdbc/idpdb";
    string description;
    string url;
    string username;
    string password;
    string host;
    int port;
    string databaseName;
    int maxPoolSize = 50;
    int minIdle = 20;
    int maxLifeTime = 60000;
    int validationTimeout;
    boolean autoCommit = true;
    string testQuery;
    string driver;
};

public type FileBaseOAuthapps record {|
    string clientId;
    string clientSecret;
    string[] callbackUrls = [];
    string[] grantTypes = [];
|};

public type TokenIssuerConfiguration record {|
    string issuer = "https://localhost:9443/oauth2/token";
    string audience = "https://localhost:9443/oauth2/token";
    string keyId = "gateway_certificate_alias";
    decimal expTime = 3600;
|};
