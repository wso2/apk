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

# DCR API
public isolated client class DCRClient {
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
        }
        http:Client httpEp = check new (serviceUrl, httpClientConfig);
        self.clientEp = httpEp;
        return;
    }
    # Registers an OAuth2 application
    #
    # + payload - Application information to register. 
    # + return - Created 
    resource isolated function post register(RegistrationRequest payload) returns Application|error {
        string resourcePath = string `/dcr/register`;
        http:Request request = new;
        json jsonBody = check payload.cloneWithType(json);
        request.setPayload(jsonBody, "application/json");
        Application response = check self.clientEp->post(resourcePath, request);
        return response;
    }
    # Get OAuth2 application information
    #
    # + client_id - Unique identifier of the OAuth2 client application. 
    # + return - Successfully Retrieved 
    resource isolated function get register/[string client_id]() returns Application|error {
        string resourcePath = string `/dcr/register/${getEncodedUri(client_id)}`;
        Application response = check self.clientEp->get(resourcePath);
        return response;
    }
    # Updates an OAuth2 application
    #
    # + client_id - Unique identifier for the OAuth2 client application. 
    # + payload - Application information to update. 
    # + return - Successfully updated 
    resource isolated function put register/[string client_id](UpdateRequest payload) returns Application|error {
        string resourcePath = string `/dcr/register/${getEncodedUri(client_id)}`;
        http:Request request = new;
        json jsonBody = check payload.cloneWithType(json);
        request.setPayload(jsonBody, "application/json");
        Application response = check self.clientEp->put(resourcePath, request);
        return response;
    }
    # Delete OAuth2 application
    #
    # + client_id - Unique identifier of the OAuth2 client application. 
    # + return - Successfully deleted 
    resource isolated function delete register/[string client_id]() returns http:Response|error {
        string resourcePath = string `/dcr/register/${getEncodedUri(client_id)}`;
        http:Response response = check self.clientEp-> delete(resourcePath);
        return response;
    }
}
