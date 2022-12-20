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

import ballerina/jwt;
import ballerina/uuid;

isolated function generateToken(JWTTokenInfo jwtInfo) returns string|error {
        TokenIssuerConfiguration & readonly issuerConfiguration = apkConfig.tokenIssuerConfiguration;
        KeyStore & readonly signingCert = apkConfig.keyStores.signing;
        string jwtid = uuid:createType1AsString();
        jwt:IssuerConfig issuerConfig = {
            issuer: issuerConfiguration.issuer,
            audience: issuerConfiguration.audience,
            expTime: issuerConfiguration.expTime,
            jwtId: jwtid,
            keyId: issuerConfiguration.keyId,
            signatureConfig: {
                config: {keyFile: signingCert.path}
            }
        };
        issuerConfig.username = jwtInfo.subscriber;

        issuerConfig.customClaims =  handleCustomClaims(jwtInfo);
        return jwt:issue(issuerConfig);
}

# This function used to handle internal token custom tokens.
#
# + jwtInfo- invoked API
# + return - Return list of custom claims
isolated function handleCustomClaims(JWTTokenInfo jwtInfo) returns map<json> {
    map<json> claims = {};
    claims["keytype"] = jwtInfo.keyType;
    claims["permittedReferer"] = jwtInfo.permittedReferrer;
    claims["permittedIp"] = jwtInfo.permittedIP;
    claims["token_type"] = "APIKey";
    claims["tierInfo"] = "";
    claims["subscribedAPIs"] = createSubscribedAPIJSON(jwtInfo.subscribedAPIs);
    claims["application"] = createApplicationJSON(jwtInfo.application);
    return claims;
}

# This GenerateSubscribedAPIS Element.
#
# + apis - subscribed APIs.
# + return - Return SubscribedAPI.
isolated function createSubscribedAPIJSON(API[] apis) returns json {
    //return apis.toJson();
    map<string>[] strippedAPIs = [];
    foreach API api in apis {
        map<string> subscribedAPI = {};
        subscribedAPI["name"] = api.name;
        subscribedAPI["context"] = api.context;
        subscribedAPI["version"] = api.'version;
        subscribedAPI["publisher"] = "apkuser";
        string? uuid = api.id;
        if uuid is string {
            subscribedAPI["uuid"] = uuid;
        }
        strippedAPIs.push(subscribedAPI);
    }
    return strippedAPIs.toJson();
}

# This GenerateApplication Element.
#
# + app - Application.
# + return - Return Application.
isolated function createApplicationJSON(Application app) returns json {
    map<string> application = {};
    string? uuid = app.applicationId;
    int? id = app.id;
    string? owner = app.owner;
    string? tier = app.throttlingPolicy;
    application["name"] = app.name;
    if uuid is string{
        application["uuid"] = uuid;
    }
    if id is int {
        application["id"] = id.toString();
    }
    if owner is string {
        application["owner"] = owner;
    }
    if tier is string {
        application["tier"] = tier;
    }
    application["tierQuotaType"] = "";
    return application;
}
