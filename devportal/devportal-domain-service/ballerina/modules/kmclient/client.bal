import ballerina/http;

# Key Manager Proxy Service
public isolated client class Client {
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
    # Register Client.
    #
    # + return - Successful operation 
    resource isolated function post register(ClientRegistrationRequest payload) returns ClientRegistrationResponse|error {
        string resourcePath = string `/register`;
        http:Request request = new;
        json jsonBody = payload.toJson();
        request.setPayload(jsonBody, "application/json");
        ClientRegistrationResponse response = check self.clientEp->post(resourcePath, request);
        return response;
    }
    # Get Client.
    #
    # + client_id - Client Id 
    # + return - Successful operation 
    resource isolated function get register(string client_id) returns ClientRegistrationResponse|error {
        string resourcePath = string `/register/${getEncodedUri(client_id)}`;
        map<anydata> queryParam = {"client_id": client_id};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        ClientRegistrationResponse response = check self.clientEp->get(resourcePath);
        return response;
    }
    # Update Client
    #
    # + client_id - Client Id 
    # + return - Successful operation 
    resource isolated function put register(string client_id, ClientUpdateRequest payload) returns ClientRegistrationResponse|error {
        string resourcePath = string `/register/${getEncodedUri(client_id)}`;
        map<anydata> queryParam = {"client_id": client_id};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        http:Request request = new;
        json jsonBody = payload.toJson();
        request.setPayload(jsonBody, "application/json");
        ClientRegistrationResponse response = check self.clientEp->put(resourcePath, request);
        return response;
    }
    # Delete Client.
    #
    # + client_id - Client Id 
    # + return - Successful operation 
    resource isolated function delete register(string client_id) returns http:Response|error {
        string resourcePath = string `/register/${getEncodedUri(client_id)}`;
        map<anydata> queryParam = {"client_id": client_id};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        http:Response response = check self.clientEp-> delete(resourcePath);
        return response;
    }
    # Get Token.
    #
    # + return - Successful operation 
    resource isolated function post token(TokenRequest payload) returns TokenResponse|error {
        string resourcePath = string `/token`;
        http:Request request = new;
        string encodedRequestBody = createFormURLEncodedRequestBody(payload);
        request.setPayload(encodedRequestBody, "application/x-www-form-urlencoded");
        TokenResponse response = check self.clientEp->post(resourcePath, request);
        return response;
    }
}
