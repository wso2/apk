import ballerina/http;

# Key Manager Proxy Service
public isolated client class Client {
    final http:Client clientEp;
    # Gets invoked to initialize the `connector`.
    #
    # + config - The configurations to be used when initializing the `connector` 
    # + serviceUrl - URL of the target service 
    # + return - An error if connector initialization failed 
    public isolated function init(string serviceUrl, ConnectionConfig config = {}) returns error? {
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
    # Initialize Key Manager
    #
    # + xOrganizationid - Organization Id 
    # + return - Successful operation 
    resource isolated function post initialize(string xOrganizationid, ConfigurationInitialization payload) returns Inline_response_200|error {
        string resourcePath = string `/initialize`;
        map<any> headerValues = {"X-OrganizationId": xOrganizationid};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        http:Request request = new;
        json jsonBody = payload.toJson();
        request.setPayload(jsonBody, "application/json");
        Inline_response_200 response = check self.clientEp->post(resourcePath, request, httpHeaders);
        return response;
    }
    # Check Key Manager initialized
    #
    # + xOrganizationid - Organization Id 
    # + return - Successful operation 
    resource isolated function head initialize(string xOrganizationid) returns Inline_response_200|error {
        string resourcePath = string `/initialize`;
        map<any> headerValues = {"X-OrganizationId": xOrganizationid};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        http:Response response = check self.clientEp->head(resourcePath, httpHeaders);
        json jsonPayload = check response.getJsonPayload();
        return jsonPayload.cloneWithType(Inline_response_200);
    }
    # Register Client.
    #
    # + xOrganizationid - Organization Id 
    # + return - Successful operation 
    resource isolated function post register(string xOrganizationid, ClientRegistrationRequest payload) returns ClientRegistrationResponse|error {
        string resourcePath = string `/register`;
        map<any> headerValues = {"X-OrganizationId": xOrganizationid};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        http:Request request = new;
        json jsonBody = payload.toJson();
        request.setPayload(jsonBody, "application/json");
        ClientRegistrationResponse response = check self.clientEp->post(resourcePath, request, httpHeaders);
        return response;
    }
    # Get Client.
    #
    # + xOrganizationid - Organization Id 
    # + client_id - Client Id 
    # + return - Successful operation 
    resource isolated function get register(string xOrganizationid, string client_id) returns ClientRegistrationResponse|error {
        string resourcePath = string `/register/${getEncodedUri(client_id)}`;
        map<anydata> queryParam = {"client_id": client_id};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        map<any> headerValues = {"X-OrganizationId": xOrganizationid};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        ClientRegistrationResponse response = check self.clientEp->get(resourcePath, httpHeaders);
        return response;
    }
    # Update Client
    #
    # + xOrganizationid - Organization Id 
    # + client_id - Client Id 
    # + return - Successful operation 
    resource isolated function put register(string xOrganizationid, string client_id, ClientUpdateRequest payload) returns ClientRegistrationResponse|error {
        string resourcePath = string `/register/${getEncodedUri(client_id)}`;
        map<anydata> queryParam = {"client_id": client_id};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        map<any> headerValues = {"X-OrganizationId": xOrganizationid};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        http:Request request = new;
        json jsonBody = payload.toJson();
        request.setPayload(jsonBody, "application/json");
        ClientRegistrationResponse response = check self.clientEp->put(resourcePath, request, httpHeaders);
        return response;
    }
    # Delete Client.
    #
    # + xOrganizationid - Organization Id 
    # + client_id - Client Id 
    # + return - Successful operation 
    resource isolated function delete register(string xOrganizationid, string client_id) returns http:Response|error {
        string resourcePath = string `/register/${getEncodedUri(client_id)}`;
        map<anydata> queryParam = {"client_id": client_id};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        map<any> headerValues = {"X-OrganizationId": xOrganizationid};
        map<string|string[]> httpHeaders = getMapForHeaders(headerValues);
        http:Response response = check self.clientEp->delete(resourcePath, headers = httpHeaders);
        return response;
    }
}
