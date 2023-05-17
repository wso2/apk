import ballerina/http;

# Provides a set of configurations for controlling the behaviours when communicating with a remote HTTP endpoint.
#
# + httpVersion - Field Description  
# + http1Settings - Field Description  
# + http2Settings - Field Description  
# + timeout - Field Description  
# + forwarded - Field Description  
# + poolConfig - Field Description  
# + cache - Field Description  
# + compression - Field Description  
# + circuitBreaker - Field Description  
# + retryConfig - Field Description  
# + responseLimits - Field Description  
# + secureSocket - Field Description  
# + proxy - Field Description  
# + validation - Field Description  
# + clientAuth - Field Description
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
    http:ClientAuthConfig clientAuth?;
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

public type TokenErrorResponse record {
    # Error code classifying the type of preProcessingError.
    string preProcessingError;
};

public type TokenResponse record {
    # OAuth access tokn issues by authorization server.
    string access_token;
    # The type of the token issued.
    string token_type;
    # The lifetime in seconds of the access token.
    int expires_in?;
    # OPTIONAL.
    # The refresh token, which can be used to obtain new access tokens.
    string refresh_token?;
    # The scope of the access token requested.
    string scope?;
};

public type ClientRegistrationError record {
    string 'error?;
    string error_description?;
};

public type UpdateRequest record {
    string[] redirect_uris?;
    string client_name?;
    string[] grant_types?;
};

public type JWKList_keys record {
    string kid?;
    string kty?;
    string use?;
    string[] key_ops?;
    string alg?;
    string x5u?;
    string[] x5c?;
    string x5t?;
    string 'x5t\#S256?;
    string e?;
    string n?;
    string x?;
    string y?;
    string d?;
    string p?;
    string q?;
    string dp?;
    string dq?;
    string qi?;
    string k?;
};

public type JWKList record {
    JWKList_keys keys?;
};

public type RegistrationRequest record {
    string[] redirect_uris?;
    string client_name?;
    string[] grant_types?;
};

public type Application record {
    string client_id?;
    string client_secret?;
    string[] redirect_uris?;
    string[] grant_types?;
    string client_name?;
    int client_secret_expires_at?;
};
