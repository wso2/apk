import ballerina/http;

# API for Partition Service
public isolated client class Client {
    private final http:Client clientEp;
    private final map<string> & readonly headers;
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
            self.headers = config.headers.cloneReadOnly();
        }
        http:Client httpEp = check new (serviceUrl, httpClientConfig);
        self.clientEp = httpEp;
        return;
    }
    # Create Event
    #
    # + payload - Create Event 
    # + return - Accepted. Created Event 
    resource isolated function post 'api\-deployment(Event payload) returns http:Response|error {
        string resourcePath = string `/api-deployment`;
        http:Request request = new;
        json jsonBody = payload.toJson();
        request.setPayload(jsonBody, "application/json");
        foreach string headerName in self.headers.keys() {
            request.setHeader(headerName, self.headers.get(headerName));
        }
        http:Response response = check self.clientEp->post(resourcePath, request);
        return response;
    }
    # Get Event
    #
    # + apiId - API Id. 
    # + return - OK. Event 
    resource isolated function get 'api\-deployment/[string apiId]() returns Partition|error {
        string resourcePath = string `/api-deployment/${getEncodedUri(apiId)}`;
        Partition response = check self.clientEp->get(resourcePath, headers = self.headers);
        return response;
    }
    # Get Active Partition
    #
    # + apiId - API Id. 
    # + return - OK. Active Partition 
    resource isolated function get 'deployable\-partition(string? apiId = ()) returns Partition|error {
        string resourcePath = string `/deployable-partition`;
        map<anydata> queryParam = {"apiId": apiId};
        resourcePath = resourcePath + check getPathForQueryParam(queryParam);
        Partition response = check self.clientEp->get(resourcePath, headers = self.headers);
        return response;
    }
}
