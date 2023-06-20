import ballerina/http;

# Provides a set of configurations for controlling the behaviours when communicating with a remote HTTP endpoint.
@display {label: "Connection Config"}
public type ConnectionConfig record {|
    # The HTTP version understood by the client
    http:HttpVersion httpVersion = http:HTTP_2_0;
    # Configurations related to HTTP/1.x protocol
    ClientHttp1Settings http1Settings?;
    # Configurations related to HTTP/2 protocol
    http:ClientHttp2Settings http2Settings?;
    # The maximum time to wait (in seconds) for a response before closing the connection
    decimal timeout = 60;
    # The choice of setting `forwarded`/`x-forwarded` header
    string forwarded = "disable";
    # Configurations associated with request pooling
    http:PoolConfiguration poolConfig?;
    # HTTP caching related configurations
    http:CacheConfig cache?;
    # Specifies the way of handling compression (`accept-encoding`) header
    http:Compression compression = http:COMPRESSION_AUTO;
    # Configurations associated with the behaviour of the Circuit Breaker
    http:CircuitBreakerConfig circuitBreaker?;
    # Configurations associated with retrying
    http:RetryConfig retryConfig?;
    # Configurations associated with inbound response size limits
    http:ResponseLimitConfigs responseLimits?;
    # SSL/TLS-related options
    http:ClientSecureSocket secureSocket?;
    # Proxy server related options
    http:ProxyConfig proxy?;
    # Enables the inbound payload validation functionality which provided by the constraint package. Enabled by default
    boolean validation = true;
    # Specify additional headers to send with the request
    map<string> headers = {};
|};

# Provides settings related to HTTP/1.x protocol.
public type ClientHttp1Settings record {|
    # Specifies whether to reuse a connection for multiple requests
    http:KeepAlive keepAlive = http:KEEPALIVE_AUTO;
    # The chunking behaviour of the request
    http:Chunking chunking = http:CHUNKING_AUTO;
    # Proxy server related options
    ProxyConfig proxy?;
|};

# Proxy server configurations to be used with the HTTP client endpoint.
public type ProxyConfig record {|
    # Host name of the proxy server
    string host = "";
    # Proxy server port
    int port = 0;
    # Proxy server username
    string userName = "";
    # Proxy server password
    @display {label: "", kind: "password"}
    string password = "";
|};

public type Partition record {
    # Partition Name
    string name?;
    # Partition Namespace
    string namespace?;
    # Number of APIs deployed in the partition
    int apiCount?;
};

public type ErrorListItem record {
    string code;
    # A description on the individual errors that occurred.
    string message;
    # A detailed description of the error message.
    string description?;
};

public type Event record {
    # Event Type
    string eventType?;
    # UUID of API 
    string apiId?;
    # API Name
    string apiName?;
    # API Context 
    string apiContext?;
    # API Version 
    string apiVersion?;
    # Organization
    string organization?;
    # Partition Name
    string partition?;
    # API Deployed Vhosts
    string[] vhosts?;
};

public type Error record {
    int code;
    # Error message.
    string message;
    # A detailed description of the error message.
    string description?;
    # Preferably a URL with more details about the error.
    string moreInfo?;
    # If there is more than one error, list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] 'error?;
};
