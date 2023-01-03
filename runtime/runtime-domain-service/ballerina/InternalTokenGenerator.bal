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
import runtime_domain_service.model;

# This Class used to generate Runtime Token
public class InternalTokenGenerator {

    public isolated function generateToken(model:API api, string username) returns string|jwt:Error {
        TokenIssuerConfiguration issuerConfiguration = runtimeConfiguration.tokenIssuerConfiguration;
        KeyStore & readonly signingCert = runtimeConfiguration.keyStores.signing;
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
        issuerConfig.username = username;

        issuerConfig.customClaims = self.handleCustomClaims(api);
        return jwt:issue(issuerConfig);

    }

    # This function used to handle internal token custom tokens.
    #
    # + api - invoked API
    # + return - Return list of custom claims
    private isolated function handleCustomClaims(model:API api) returns map<json> {
        map<json> claims = {};
        claims["keytype"] = "PRODUCTION";
        claims["uuid"] = api.metadata.uid;
        claims["token_type"] = "InternalKey";
        claims["subscribedAPIs"] = [self.createSubscribedAPIJSON(api)];
        return claims;
    }
    # This GenerateSubscribedAPIS Element.
    #
    # + api - Invoke API.
    # + return - Return SubscribedAPI.
    private isolated function createSubscribedAPIJSON(model:API api) returns json {
        map<string> subscribedAPIs = {};
        subscribedAPIs["name"] = api.spec.apiDisplayName;
        subscribedAPIs["context"] = api.spec.context;
        subscribedAPIs["version"] = api.spec.apiVersion;
        subscribedAPIs["publisher"] = APK_USER;
        subscribedAPIs["uuid"] = <string>api.metadata.uid;
        return subscribedAPIs;
    }
}

