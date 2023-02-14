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
import ballerina/http;
import ballerina/jwt;
import ballerina/log;
import ballerina/time;
import ballerina/uuid;

public class LoginClient {
    public isolated function handleLoginRequest(http:Request request) returns http:Found {
        do {
            map<string> formParams = check request.getFormParams();
            Login_body payload = check formParams.cloneWithType(Login_body);
            http:Cookie[] cookies = request.getCookies();
            string sessionKey = payload.sessionKey;
            http:Cookie? sessionCookieValue = ();
            foreach http:Cookie cookie in cookies {
                if cookie.name == SESSION_KEY_PREFIX + sessionKey {
                    sessionCookieValue = cookie;
                    break;
                    // session available.
                }
            }
            if sessionCookieValue is http:Cookie {
                string sessionValue = sessionCookieValue.value;
                jwt:Payload validatedPayload = check jwt:validate(sessionValue, getValidationConfig());
                if validatedPayload.exp <= time:utcNow()[0] {
                    string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=expired_session&error_description=Login Session expired.";
                    return {headers: {"Location": loginPageRedirect}};
                }
                // non expired session.
                string? tokenType = validatedPayload.hasKey(TOKEN_TYPE_CLAIM) ? <string>validatedPayload.get(TOKEN_TYPE_CLAIM) : ();
                if tokenType is () || tokenType != SESSION_KEY_TYPE {
                    string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=invalid_request&error_description=Invalid login Session.";
                    return {headers: {"Location": loginPageRedirect}};
                }
                return self.validateUserAndOrg(payload, validatedPayload);

            } else {
                string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=invalid_request&error_description=Invalid login Session.";
                return {headers: {"Location": loginPageRedirect}};
            }
        } on fail var e {
            log:printError("Error occured login user.", e);
            string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Error Occurred.";
            return {headers: {"Location": loginPageRedirect}};
        }

    }
    private isolated function validateUserAndOrg(Login_body loginPayLoad, jwt:Payload sessionData) returns http:Found {
        User[] & readonly users = idpConfiguration.user;
        User? derivedUser = ();
        foreach User & readonly user in users {
            if user.username == loginPayLoad.username && user.password == loginPayLoad.password {
                derivedUser = user;
                // user found.
                break;
            }
        }
        if derivedUser is User {
            string? organization = loginPayLoad.organization;
            if organization is string {
                int? organizationIndex = derivedUser.organizations.indexOf(organization, 0);
                if organizationIndex is int {
                    // login successed for organization.
                    [string, string|jwt:Error] generateSucessSession = self.generateSucessSessionData(derivedUser.username, organization, sessionData);
                    string sessionKey = generateSucessSession[0];
                    string|jwt:Error cookieValue = generateSucessSession[1];
                    if cookieValue is string {
                        string loginPageRedirect = idpConfiguration.loginCallBackURl + "?" + STATE_KEY_QUERY_PARAM + "=" + sessionKey;
                        http:CookieOptions cookieOption = {domain: idpConfiguration.hostname, secure: true, path: "/"};
                        http:Cookie cookie = new (SESSION_KEY_PREFIX + sessionKey, cookieValue, cookieOption);
                        return {
                            headers: {
                                "Location": loginPageRedirect,
                                "Set-Cookie": cookie.toStringValue()
                            }
                        };
                    } else {
                        log:printError("Error occured generating success Session.", cookieValue);
                        string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Error Occurred.";
                        return {headers: {"Location": loginPageRedirect}};
                    }
                } else {
                    string loginPageRedirect = idpConfiguration.loginPageURl + "?" + STATE_KEY_QUERY_PARAM + "=" + loginPayLoad.sessionKey + "&authFailure=true&authFailureMsg=invalid.organization&errorCode=" + INVALID_ORGANIZATION;
                    return {
                        headers: {
                            "Location": loginPageRedirect
                        }
                    };
                }

            } else {
                if derivedUser.superAdmin {
                    // login successed for organization.
                    [string, string|jwt:Error] generateSucessSession = self.generateSucessSessionData(derivedUser.username, (), sessionData);
                    string sessionKey = generateSucessSession[0];
                    string|jwt:Error cookieValue = generateSucessSession[1];
                    if cookieValue is string {
                        string loginPageRedirect = idpConfiguration.loginCallBackURl + "?" + STATE_KEY_QUERY_PARAM + "=" + sessionKey;
                        http:CookieOptions cookieOption = {domain: idpConfiguration.hostname, secure: true, path: "/"};
                        http:Cookie cookie = new (SESSION_KEY_PREFIX + sessionKey, cookieValue, cookieOption);
                        return {
                            headers: {
                                "Location": loginPageRedirect,
                                "Set-Cookie": cookie.toStringValue()
                            }
                        };
                    } else {
                        log:printError("Error occured generating success Session.", cookieValue);
                        string loginPageRedirect = idpConfiguration.loginErrorPageUrl + "?error=server_error&error_description=Internal Error Occurred.";
                        return {headers: {"Location": loginPageRedirect}};
                    }
                } else {
                    string loginPageRedirect = idpConfiguration.loginPageURl + "?" + STATE_KEY_QUERY_PARAM + "=" + loginPayLoad.sessionKey + "&authFailure=true&authFailureMsg=login.fail.message&errorCode=" + INVALID_PERMISSION;
                    return {
                        headers: {
                            "Location": loginPageRedirect
                        }
                    };
                }
            }

        } else {
            string loginPageRedirect = idpConfiguration.loginPageURl + "?" + STATE_KEY_QUERY_PARAM + "=" + loginPayLoad.sessionKey + "&authFailure=true&authFailureMsg=login.fail.message&errorCode=" + INVALID_USERNAME_OR_PASSWORD;
            return {
                headers: {
                    "Location": loginPageRedirect
                }
            };
        }
    }
    public isolated function generateSucessSessionData(string user, string? organization, jwt:Payload requestSessionPayload) returns [string, string|jwt:Error] {

        TokenIssuerConfiguration issuerConfiguration = idpConfiguration.tokenIssuerConfiguration;
        string jwtid = uuid:createType1AsString();
        jwt:IssuerConfig issuerConfig = {
            issuer: issuerConfiguration.issuer,
            expTime: 600,
            jwtId: jwtid,
            username: user,
            signatureConfig: {
                config: {keyFile: idpConfiguration.keyStores.signing.keyFile}
            }
        };
        string redirectUri = <string>requestSessionPayload.get(REDIRECT_URI_CLAIM);
        string clientId = <string>requestSessionPayload.get(CLIENT_ID_CLAIM);
        json[] scopes = <json[]>requestSessionPayload.get(SCOPES_CLAIM);
        issuerConfig.customClaims = {[REDIRECT_URI_CLAIM] : redirectUri, [SCOPES_CLAIM] : scopes, [CLIENT_ID_CLAIM] : clientId, [TOKEN_TYPE_CLAIM] : SESSION_KEY_TYPE};
        if organization is string {
            issuerConfig.customClaims = {[REDIRECT_URI_CLAIM] : redirectUri, [SCOPES_CLAIM] : scopes, [CLIENT_ID_CLAIM] : clientId, [ORGANIZATION_CLAIM] : organization, [TOKEN_TYPE_CLAIM] : SESSION_KEY_TYPE};
        }
        return [jwtid, jwt:issue(issuerConfig)];
    }
}
