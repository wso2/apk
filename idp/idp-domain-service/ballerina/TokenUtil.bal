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

import ballerina/log;
import ballerina/regex;
import ballerina/lang.array;
import ballerina/uuid;
import ballerina/jwt;
import ballerina/http;
import ballerina/time;

public class TokenUtil {

    public isolated function generateToken(string? authorization, Token_body payload) returns UnauthorizedTokenErrorResponse|BadRequestTokenErrorResponse|OkTokenResponse|error {
        if (authorization is ()) || (authorization.toString().trim().length() == 0) || (!authorization.toString().startsWith("Basic ")) {
            UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Unauthorized"}};
            return unauthorized;
        } else {
            string authorizationString = authorization.toString();
            string basicEncodedValue = authorizationString.substring(6, authorizationString.length()).trim();
            if basicEncodedValue.length() == 0 {
                UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Unauthorized"}};
                return unauthorized;
            } else {
                do {
                    byte[] base64DecodedBytes = check array:fromBase64(basicEncodedValue);
                    string base64DecodedString = check 'string:fromBytes(base64DecodedBytes);
                    string[] clientIdSecretToken = regex:split(base64DecodedString.trim(), ":");
                    if clientIdSecretToken.length() != 2 {
                        UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Unauthorized"}};
                        return unauthorized;
                    }
                    DCRMClient dcrmClient = new;
                    Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError application = dcrmClient.getApplicationIncludeFileBaseApps(clientIdSecretToken[0]);
                    if application is Application {
                        if application.client_secret != clientIdSecretToken[1] {
                            UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "unauthorized_client", error_description: "Invalide Client Id/Secret"}};
                            return unauthorized;
                        }
                        string grantType = payload.grant_type;
                        string[]? grantTypes = application.grant_types;
                        if grantTypes is string[] {
                            int? indexOf = grantTypes.indexOf(grantType);
                            if indexOf is () {
                                BadRequestTokenErrorResponse tokenError = {body: {'error: "unsupported_grant_type", error_description: grantType + " not supported by application."}};
                                return tokenError;
                            }
                        }
                        if grantType == CLIENT_CREDENTIALS_GRANT_TYPE {
                            return self.handleClientCredentialsGrant(payload, application);
                        } else if grantType == AUTHORIZATION_CODE_GRANT_TYPE {
                            return self.handleAuthorizationCodeGrant(payload, application);
                        } else if grantType == REFRESH_TOKEN_GRANT_TYPE {
                            return self.hanleRefreshTokenGrant(payload, application);
                        } else {
                            BadRequestTokenErrorResponse tokenError = {body: {'error: "unsupported_grant_type", error_description: grantType + " not supported by system."}};
                            return tokenError;
                        }
                    } else {
                        UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Invalide Client Id/Secret"}};
                        return unauthorized;
                    }
                }
                on fail var e {
                    log:printError("Error on decoding base64", e);
                    BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
                    return tokenError;
                }
            }
        }
    }
    public isolated function handleClientCredentialsGrant(Token_body payload, Application application) returns OkTokenResponse|BadRequestTokenErrorResponse|UnauthorizedTokenErrorResponse {
        string[] scopeArray = self.filterScopes(payload.scope);
        string|jwt:Error tokenResult = self.issueToken(application, (), scopeArray, (), ACCESS_TOKEN_TYPE);
        if tokenResult is string {
            TokenResponse tokenResponse = {
                access_token: tokenResult,
                token_type: TOKEN_TYPE_BEARER,
                expires_in: idpConfiguration.tokenIssuerConfiguration.expTime,
                scope: string:'join(" ", ...scopeArray)
            };
            return {body: tokenResponse};
        }
        else {
            log:printError("Error on Generating token", tokenResult);
            BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            return tokenError;

        }
    }
    public isolated function issueToken(Application application, string? username, string[] scopes, string? organization, string tokenType) returns string|jwt:Error {
        TokenIssuerConfiguration issuerConfiguration = idpConfiguration.tokenIssuerConfiguration;
        string jwtid = uuid:createType1AsString();
        decimal exptime = tokenType == ACCESS_TOKEN_TYPE ? issuerConfiguration.expTime : issuerConfiguration.refrshTokenValidity;
        jwt:IssuerConfig issuerConfig = {
            issuer: issuerConfiguration.issuer,
            expTime: exptime,
            jwtId: jwtid,
            keyId: issuerConfiguration.keyId,
            signatureConfig: {
                config: {keyFile: idpConfiguration.keyStores.signing.keyFile}
            }
        };
        if username is string {
            issuerConfig.username = username;
        } else {
            issuerConfig.username = application.client_id;
        }
        map<string> customClaims = {};
        customClaims[CLIENT_ID_CLAIM] = <string>application.client_id;
        customClaims[SCOPES_CLAIM] = string:'join(" ", ...scopes);
        if organization is string && organization.toString().trim().length() > 0 {
            customClaims[ORGANIZATION_CLAIM] = organization;
        }
        if tokenType == REFRESH_TOKEN_TYPE {
            customClaims[TOKEN_TYPE_CLAIM] = tokenType;
        }
        issuerConfig.customClaims = customClaims;
        return jwt:issue(issuerConfig);
    }
    public isolated function handleAuthorizationCodeGrant(Token_body payload, Application application) returns BadRequestTokenErrorResponse|OkTokenResponse|error {
        string? authorization_code = payload.code;
        string? redirectUri = payload.redirect_uri;

        if (authorization_code is () || authorization_code.toString().trim().length() == 0) || (redirectUri is () || redirectUri.toString().trim().length() == 0) {
            BadRequestTokenErrorResponse tokenError = {body: {'error: "invalid_request", error_description: "authorization_code|redirect_uri not available in request."}};
            return tokenError;
        }
        // authorization code available.
        jwt:Payload|jwt:Error validatedPayload = jwt:validate(authorization_code, getValidationConfig());
        if validatedPayload is jwt:Payload {
            // validating expiry.
            string tokenType = validatedPayload.hasKey(TOKEN_TYPE_CLAIM) ? <string>validatedPayload.get(TOKEN_TYPE_CLAIM) : "";
            if tokenType != AUTHORIZATION_CODE_TYPE {
                BadRequestTokenErrorResponse tokenError = {"body": {'error: "invalid_request", error_description: "Invalid authorization_code"}};
                return tokenError;
            }
            if validatedPayload.exp <= time:utcNow()[0] {
                BadRequestTokenErrorResponse tokenError = {"body": {'error: "invalid_grant", error_description: "authorization_code expired."}};
                return tokenError;
            }
            string requestRedirectUrl = <string>validatedPayload.get(REDIRECT_URI_CLAIM);
            string clientId = <string>validatedPayload.get(CLIENT_ID_CLAIM);
            json[] scopes = <json[]>validatedPayload.get(SCOPES_CLAIM);
            string sub = <string>validatedPayload.sub;
            string[]? redirectUris = application.redirect_uris;

            string? organization = validatedPayload.hasKey(ORGANIZATION_CLAIM) ? <string>validatedPayload.get(ORGANIZATION_CLAIM) : ();
            if requestRedirectUrl != redirectUri || application.client_id != clientId || (redirectUris is () || redirectUris.indexOf(redirectUri) is ()) {
                BadRequestTokenErrorResponse tokenError = {"body": {'error: "unauthorized_client", error_description: "redirectUrl not matched with application"}};
                return tokenError;
            }
            string[] scopesArray = [];
            foreach json scope in scopes {
                scopesArray.push(scope.toString());
            }
            do {
                string accessToken = check self.issueToken(application, sub, scopesArray, organization, ACCESS_TOKEN_TYPE);
                string refreshToken = check self.issueToken(application, sub, scopesArray, organization, REFRESH_TOKEN_TYPE);
                TokenResponse token = {access_token: accessToken, refresh_token: refreshToken, expires_in: idpConfiguration.tokenIssuerConfiguration.expTime, token_type: TOKEN_TYPE_BEARER, scope: string:'join(" ", ...scopesArray)};
                return {body: token};
            } on fail var e {
                log:printInfo("Error on generating token", e);
                return {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            }
        } else {
            log:printError("Error on validating authorization_code", validatedPayload);
            BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            return tokenError;
        }
    }
    public isolated function hanleRefreshTokenGrant(Token_body payload, Application application) returns BadRequestTokenErrorResponse|OkTokenResponse {
        string? refresh_token = payload.refresh_token;

        if (refresh_token is () || refresh_token.toString().trim().length() == 0) {
            BadRequestTokenErrorResponse tokenError = {body: {'error: "invalid_request", error_description: "refresh_token not available in request."}};
            return tokenError;
        }
        // authorization code available.
        jwt:Payload|jwt:Error validatedPayload = jwt:validate(refresh_token, getValidationConfig());
        if validatedPayload is jwt:Payload {
            // validating expiry.
            string tokenType = validatedPayload.hasKey(TOKEN_TYPE_CLAIM) ? <string>validatedPayload.get(TOKEN_TYPE_CLAIM) : "";
            if tokenType != REFRESH_TOKEN_TYPE {
                BadRequestTokenErrorResponse tokenError = {"body": {'error: "invalid_request", error_description: "Invalid refresh_token"}};
                return tokenError;
            }
            if validatedPayload.exp <= time:utcNow()[0] {
                BadRequestTokenErrorResponse tokenError = {"body": {'error: "invalid_grant", error_description: "refredh_token expired."}};
                return tokenError;
            }
            string clientId = <string>validatedPayload.get(CLIENT_ID_CLAIM);
            string scopes = <string>validatedPayload.get(SCOPES_CLAIM);
            string sub = <string>validatedPayload.sub;

            string? organization = validatedPayload.hasKey(ORGANIZATION_CLAIM) ? <string>validatedPayload.get(ORGANIZATION_CLAIM) : ();
            if application.client_id != clientId {
                BadRequestTokenErrorResponse tokenError = {"body": {'error: "invalid_request", error_description: "Invalid refresh_token"}};
                return tokenError;
            }
            do {
                string[] scopesArray = regex:split(scopes, " ");
                string accessToken = check self.issueToken(application, sub, scopesArray, organization, ACCESS_TOKEN_TYPE);
                string refreshToken = check self.issueToken(application, sub, scopesArray, organization, REFRESH_TOKEN_TYPE);
                TokenResponse token = {access_token: accessToken, refresh_token: refreshToken, expires_in: idpConfiguration.tokenIssuerConfiguration.expTime, token_type: TOKEN_TYPE_BEARER, scope: scopes};
                return {body: token};
            } on fail var e {
                log:printInfo("Error on generating token", e);
                return {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            }
        } else {
            log:printError("Error on validating authorization_code", validatedPayload);
            BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            return tokenError;
        }
    }
    public isolated function handleAuthorizeRequest(string response_type, string client_id, string? redirect_uri, string? scope, string? state) returns http:Found {
        do {
            if client_id.trim().length() > 0 && redirect_uri is string {
                DCRMClient dcrmClient = new;
                Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError application = dcrmClient.getApplicationIncludeFileBaseApps(client_id);
                if application is Application {
                    string[]? grantTypes = application.grant_types;
                    if grantTypes is string[] {
                        int? indexOf = grantTypes.indexOf(AUTHORIZATION_CODE_GRANT_TYPE);
                        if indexOf is () {
                            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=unauthorized_client&error_description=authorization_code grant not supported from application";
                            return {headers: {"Location": loginPageRedirect}};
                        }
                        if response_type != AUTHORIZATION_CODE_QUERY_PARAM {
                            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=unsupported_response_type&error_description=" + response_type + " not supported from authorization server.";
                            return {headers: {"Location": loginPageRedirect}};
                        }
                        string[] scopeArray = self.filterScopes(scope);
                        return self.redirectRequest(application, redirect_uri, scopeArray, state);
                    }
                } else if application is NotFoundClientRegistrationError {
                    string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=unauthorized_client&error_description=Client application not found in system";
                    return {headers: {"Location": loginPageRedirect}};
                } else {
                    string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Server Error";
                    return {headers: {"Location": loginPageRedirect}};
                }
            }
            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=invalid_request&error_description=authorization_code grant not supported from application";
            return {headers: {"Location": loginPageRedirect}};
        }
        on fail var e {
            log:printError("Error on authorizing request", e);
            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Server Error";
            return {headers: {"Location": loginPageRedirect}};
        }
    }
    public isolated function redirectRequest(Application application, string redirectUri, string[] scopes, string? state) returns http:Found {
        TokenIssuerConfiguration issuerConfiguration = idpConfiguration.tokenIssuerConfiguration;
        string jwtid = uuid:createType1AsString();
        jwt:IssuerConfig issuerConfig = {
            issuer: issuerConfiguration.issuer,
            expTime: 600,
            jwtId: jwtid,
            signatureConfig: {
                config: {keyFile: idpConfiguration.keyStores.signing.keyFile}
            }
        };
        issuerConfig.customClaims = {[REDIRECT_URI_CLAIM] : redirectUri, [SCOPES_CLAIM] : scopes, [CLIENT_ID_CLAIM] : application.client_id, [TOKEN_TYPE_CLAIM] : SESSION_KEY_TYPE};
        string|jwt:Error stateKey = jwt:issue(issuerConfig);
        if stateKey is string {
            string loginPageRedirect = idpConfiguration.loginPageURl + "?" + STATE_KEY_QUERY_PARAM + "=" + jwtid;

            http:CookieOptions cookieOption = {domain: idpConfiguration.hostname, secure: true, path: "/"};
            http:Cookie cookie = new (SESSION_KEY_PREFIX + jwtid, stateKey, cookieOption);
            return {
                headers: {
                    "Location": loginPageRedirect,
                    "Set-Cookie": cookie.toStringValue()
                }
            };
        } else {
            log:printInfo("Error on generating State", stateKey);
            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Server Error";
            return {headers: {"Location": loginPageRedirect}};
        }
    }
    public isolated function filterScopes(string? scopes) returns string[] {
        string[] scopeArray = ["default"];
        if scopes is string && scopes.trim().length() > 0 {
            scopeArray = regex:split(scopes, " ");
        }
        return scopeArray;
    }
    public isolated function handleOauthCallBackRequest(http:Request request, string sessionKey) returns http:Found {
        do {

            http:Cookie[] cookies = request.getCookies();
            http:Cookie? sessionCookieValue = ();
            foreach http:Cookie cookie in cookies {
                if cookie.name == SESSION_KEY_PREFIX + sessionKey {
                    sessionCookieValue = cookie;
                    break;
                }
            }
            if sessionCookieValue is http:Cookie {
                string sessionValue = sessionCookieValue.value;
                jwt:Payload validatedPayload = check jwt:validate(sessionValue, getValidationConfig());
                if validatedPayload.exp <= time:utcNow()[0] {
                    string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=expired_session&error_description=Login Session expired.";
                    return {headers: {"Location": loginPageRedirect}};
                }
                string? tokenType = validatedPayload.hasKey(TOKEN_TYPE_CLAIM) ? <string>validatedPayload.get(TOKEN_TYPE_CLAIM) : ();
                if tokenType is () || tokenType != SESSION_KEY_TYPE {
                    string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=invalid_request&error_description=Invalid login Session.";
                    return {headers: {"Location": loginPageRedirect}};
                }
                // non expired session.
                return self.generateOauthcodeResponse(validatedPayload);
            }
            else {
                string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=invalid_request&error_description=Invalid login Session.";
                return {headers: {"Location": loginPageRedirect}};
            }
        } on fail var e {
            log:printError("Error occured login user.", e);
            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Error Occurred.";
            return {headers: {"Location": loginPageRedirect}};
        }

    }
    public isolated function generateOauthcodeResponse(jwt:Payload payload) returns http:Found {
        string redirectUri = <string>payload.get(REDIRECT_URI_CLAIM);
        string clientId = <string>payload.get(CLIENT_ID_CLAIM);
        json[] scopes = <json[]>payload.get(SCOPES_CLAIM);
        string sub = <string>payload.sub;
        string? organization = payload.hasKey(ORGANIZATION_CLAIM) ? <string>payload.get(ORGANIZATION_CLAIM) : ();
        TokenIssuerConfiguration issuerConfiguration = idpConfiguration.tokenIssuerConfiguration;
        string jwtid = uuid:createType1AsString();
        jwt:IssuerConfig issuerConfig = {
            issuer: issuerConfiguration.issuer,
            expTime: 600,
            jwtId: jwtid,
            username: sub,
            keyId: issuerConfiguration.keyId,
            signatureConfig: {
                config: {keyFile: idpConfiguration.keyStores.signing.keyFile}
            }
        };
        issuerConfig.customClaims = {[REDIRECT_URI_CLAIM] : redirectUri, [SCOPES_CLAIM] : scopes, [CLIENT_ID_CLAIM] : clientId, [TOKEN_TYPE_CLAIM] : AUTHORIZATION_CODE_TYPE};
        if organization is string {
            issuerConfig.customClaims = {[REDIRECT_URI_CLAIM] : redirectUri, [SCOPES_CLAIM] : scopes, [CLIENT_ID_CLAIM] : clientId, [ORGANIZATION_CLAIM] : organization, [TOKEN_TYPE_CLAIM] : AUTHORIZATION_CODE_TYPE};
        }
        do {
            string oauthcode = check jwt:issue(issuerConfig);
            string redirectUrl = redirectUri.includes("?") ? (redirectUri + "&" + AUTHORIZATION_CODE_QUERY_PARAM + "=" + oauthcode) : (redirectUri + "?" + AUTHORIZATION_CODE_QUERY_PARAM + "=" + oauthcode);
            return {headers: {"Location": redirectUrl}};
        } on fail var e {
            log:printError("Error occured login user.", e);
            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Error Occurred.";
            return {headers: {"Location": loginPageRedirect}};
        }
    }
}
