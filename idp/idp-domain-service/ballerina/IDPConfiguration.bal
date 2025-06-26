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

import ballerina/io;
import ballerina/log;

# Description
#
# + hostname - Field Description  
# + loginPageURl - Field Description  
# + loginErrorPageUrl - Field Description  
# + loginCallBackURl - Field Description  
# + user - Field Description  
# + keyStores - Field Description  
# + fileBaseApp - Field Description  
# + tokenIssuerConfiguration - Field Description
public type IDPConfiguration record {|
    string hostname = "localhost:9443";
    string loginPageURl;
    string loginErrorPageUrl;
    string loginCallBackURl;
    User[] user = [];
    KeyStores keyStores = {};
    FileBaseOAuthapps[] fileBaseApp = [];
    TokenIssuerConfiguration tokenIssuerConfiguration = {};
|};

public type KeyStores record {|
    CertKey signing = {keyFile: "/home/wso2kgw/idp/security/wso2carbon.key", certFile: "/home/wso2kgw/idp/security/wso2carbon.pem"};
    CertKey tls = {keyFile: "/home/wso2kgw/idp/security/idp.key",certFile: "/home/wso2kgw/idp/security/idp.crt"};
|};

# Description
#
# + username - User name  
# + password - Password of user.
# + organizations - organizations belongs to user.  
# + superAdmin - User mark as super_admin.
public type User record {|
    string username;
    string password;
    string[] organizations = [];
    boolean superAdmin = false;
|};

# Represents combination of certificate, private key and private key password if encrypted.
#
# + certFile - A file containing the certificate
# + keyFile - A file containing the private key in PKCS8 format
# + keyPassword - Password of the private key if it is encrypted
public type CertKey record {|
    string certFile;
    string keyFile;
    string keyPassword?;
|};

public type DatasourceConfiguration record {
    string name = "jdbc/idpdb";
    string description;
    string url;
    string username;
    string password = getPassword();
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
    decimal refrshTokenValidity = 86400;
|};

public isolated function getPassword() returns string {
    string|error password = io:fileReadString("/home/wso2kgw/idp/security/database/db-password");
    if (password is error) {
        log:printError("Error while reading the password");
        return "";
    }
    return password;
};
