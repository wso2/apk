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

# OAuth API
public isolated client class TokenClient {
    final http:Client clientEp;
    # Gets invoked to initialize the `connector`.
    #
    # + config - The configurations to be used when initializing the `connector` 
    # + serviceUrl - URL of the target service 
    # + return - An error if connector initialization failed 
    public isolated function init(string serviceUrl, ConnectionConfig config =  {}) returns error? {
        http:ClientConfiguration httpClientConfig = {httpVersion: config.httpVersion, timeout: config.timeout, forwarded: config.forwarded, poolConfig: config.poolConfig, compression: config.compression, circuitBreaker: config.circuitBreaker, retryConfig: config.retryConfig, validation: config.validation};
        do {
            if config.http1Settings is ClientHttp1Settings {
                ClientHttp1Settings settings = check config.http1Settings.ensureType(ClientHttp1Settings);
                httpClientConfig.http1Settings = {...settings};
            }
            if config.http2Settings is http:ClientHttp2Settings {
                httpClientConfig.http2Settings = check config.http2Settings.ensureType(http:ClientHttp2Settings);
            }
            if config.cache is http:CacheConfig {
                httpClientConfig.cache = check config.cache.ensureType(http:CacheConfig);
            }
            if config.responseLimits is http:ResponseLimitConfigs {
                httpClientConfig.responseLimits = check config.responseLimits.ensureType(http:ResponseLimitConfigs);
            }
            if config.secureSocket is http:ClientSecureSocket {
                httpClientConfig.secureSocket = check config.secureSocket.ensureType(http:ClientSecureSocket);
            }
            if config.proxy is http:ProxyConfig {
                httpClientConfig.proxy = check config.proxy.ensureType(http:ProxyConfig);
            }
            if config.cookieConfig is http:CookieConfig{
                httpClientConfig.cookieConfig = check config.cookieConfig.ensureType(http:CookieConfig);
            }
        }
        http:Client httpEp = check new (serviceUrl, httpClientConfig);
        self.clientEp = httpEp;
        return;
    }
    #
    # + response_type - Expected response type 
    # + client_id - OAuth client identifier 
    # + redirect_uri - Clients redirection endpoint 
    # + scope - OAuth scopes 
    # + state - Opaque value used by the client to maintain state between the request and callback 
    # + return - Response from authorization endpoint 
    resource isolated function get authorize(string response_type, string client_id, string? redirect_uri = (), string? scope = (), string? state = ()) returns http:Response|error {
        string resourcePath = string `/oauth2/authorize`;
        map<anydata> queryParam = {"response_type": response_type, "client_id": client_id, "redirect_uri": redirect_uri, "scope": scope, "state": state};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        http:Response response = check self.clientEp->get(resourcePath);
        return response;
    }
    #
    # + sessionKey - Session key. 
    # + return - Response from authorization endpoint 
    resource isolated function get 'auth\-callback(string sessionKey) returns http:Response|error {
        string resourcePath = string `/oauth2/auth-callback`;
        map<anydata> queryParam = {"sessionKey": sessionKey};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        http:Response response = check self.clientEp->get(resourcePath);
        return response;
    }
    #
    # + authorization - Authentication scheme header 
    # + return - OK. Successful response from token endpoint. 
    resource isolated function post token(Token_body payload, string? authorization = ()) returns TokenResponse|error {
        string resourcePath = string `/oauth2/token`;
        map<any> headerValues = {"Authorization": authorization};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        http:Request request = new;
        string encodedRequestBody = createFormURLEncodedRequestBody(payload);
        request.setPayload(encodedRequestBody, "application/x-www-form-urlencoded");
        TokenResponse response = check self.clientEp->post(resourcePath, request, httpHeaders);
        return response;
    }
    #
    # + return - Signing key List 
    resource isolated function get keys() returns JWKList|error {
        string resourcePath = string `/oauth2/keys`;
        JWKList response = check self.clientEp->get(resourcePath);
        return response;
    }
}
