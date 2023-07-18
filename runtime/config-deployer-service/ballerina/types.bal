import ballerina/http;
import ballerina/constraint;

public type OkAnydata record {|
    *http:Ok;
    anydata body;
|};

public type NotFoundError record {|
    *http:NotFound;
    Error body;
|};

public type BadRequestError record {|
    *http:BadRequest;
    Error body;
|};

public type AcceptedString record {|
    *http:Accepted;
    string body;
|};

public type InternalServerErrorError record {|
    *http:InternalServerError;
    Error body;
|};

public type ErrorListItem record {
    string code;
    # Description about individual errors occurred
    string message;
    # A detail description about the error message.
    string description?;
};

public type EndpointSecurity record {
    boolean enabled?;
    BasicEndpointSecurity securityType?;
};

public type CustomClaims record {
    string claim;
    string value;
};

public type K8sService record {
    string name?;
    string namespace?;
    int port?;
    string protocol?;
};

public type RateLimit record {
    # Number of requests allowed per specified unit of time
    int requestsPerUnit;
    # Unit of time
    string unit;
};

public type BasicEndpointSecurity record {
    string secretName?;
    string userNameKey?;
    string passwordKey?;
};

public type APKOperations record {
    string target?;
    string verb?;
    # Authentication mode for resource (true/false)
    boolean authTypeEnabled?;
    EndpointConfigurations endpointConfigurations?;
    APIOperationPolicies operationPolicies?;
    RateLimit operationRateLimit?;
    string[] scopes?;
};

public type APKOperationPolicy InterceptorPolicy|BackendJWTPolicy;

public type DeployApiBody record {
    # apk-configuration file
    record {byte[] fileContent; string fileName;} apkConfiguration?;
    # api definition (OAS/Graphql/WebSocket)
    record {byte[] fileContent; string fileName;} definitionFile?;
};

public type APIOperationPolicies record {
    APKOperationPolicy[] request?;
    APKOperationPolicy[] response?;
};

public type APKConf_additionalProperties record {
    string name?;
    string value?;
};

public type Resiliency record {
    Timeout timeout?;
    RetryPolicy retryPolicy?;
};

public type AuthenticationRequest JWTAuthentication|APIKeyAuthentication;

public type EndpointConfigurations record {
    EndpointConfiguration production?;
    EndpointConfiguration sandbox?;
};

public type JWTAuthentication record {
    *Authentication;
    boolean sendTokenToUpstream?;
    string headerName?;
    boolean headerEnable?;
};

public type Timeout record {
    int maxRouteTimeoutSeconds?;
    int routeIdleTimeoutSeconds?;
    int routeTimeoutSeconds?;
};

public type GenerateK8sResourcesBody record {
    # apk-configuration file
    record {byte[] fileContent; string fileName;} apkConfiguration?;
    # api definition (OAS/Graphql/WebSocket)
    record {byte[] fileContent; string fileName;} definitionFile?;
    # Type of API
    string apiType?;
};

public type Error record {
    int code;
    # Error message.
    string message;
    # A detail description about the error message.
    string description?;
    # Preferably an url with more details about the error.
    string moreInfo?;
    # If there are more than one error list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] 'error?;
};

public type DefinitionBody record {
    # api definition (OAS/Graphql/WebSocket)
    record {byte[] fileContent; string fileName;} definition?;
    # url of api definition
    string url?;
    # Type of API
    string apiType?;
};

# Map of virtual hosts of API
#
# + production - Field Description  
# + sandbox - Field Description
public type APKConf_vhosts record {
    string[] production?;
    string[] sandbox?;
};

# CORS Configuration of API
#
# + corsConfigurationEnabled - Field Description  
# + accessControlAllowOrigins - Field Description  
# + accessControlAllowCredentials - Field Description  
# + accessControlAllowHeaders - Field Description  
# + accessControlAllowMethods - Field Description  
# + accessControlExposeHeaders - Field Description
# + accessControlAllowMaxAge - Field Description
public type CORSConfiguration record {
    boolean corsConfigurationEnabled?;
    string[] accessControlAllowOrigins?;
    boolean accessControlAllowCredentials?;
    string[] accessControlAllowHeaders?;
    string[] accessControlAllowMethods?;
    int accessControlAllowMaxAge?;
    string[] accessControlExposeHeaders?;
};

public type BaseOperationPolicy record {
    string policyName;
    string policyVersion = "v1";
    string policyId?;
};

public type Authentication record {
    string authType?;
    boolean enabled?;
};

public type InterceptorPolicy record {
    *BaseOperationPolicy;
    InterceptorPolicy_parameters parameters?;
};

public type APKConf record {
    # UUID of the API
    string id?;
    @constraint:String {maxLength: 60, minLength: 1}
    string name;
    @constraint:String {maxLength: 232, minLength: 1}
    string context;
    @constraint:String {maxLength: 30, minLength: 1}
    string version;
    string 'type = "REST";
    # Organization of the API
    string organization?;
    # Is this the default version of the API
    boolean defaultVersion = false;
    EndpointConfigurations endpointConfigurations?;
    APKOperations[] operations?;
    APIOperationPolicies apiPolicies?;
    RateLimit apiRateLimit?;
    AuthenticationRequest[] authentication?;
    # Map of custom properties of API
    APKConf_additionalProperties[] additionalProperties?;
    # Map of virtual hosts of API
    APKConf_vhosts vhosts?;
    # CORS Configuration of API
    CORSConfiguration corsConfiguration?;
};

public type InterceptorPolicy_parameters record {
    string backendUrl?;
    boolean headersEnabled?;
    boolean bodyEnabled?;
    boolean trailersEnabled?;
    boolean contextEnabled?;
};

public type BackendJWTPolicy_parameters record {
    boolean enabled?;
    string encoding?;
    string signingAlgorithm?;
    string header?;
    int tokenTTL?;
    CustomClaims[] customClaims?;
};

public type EndpointConfiguration record {
    string|K8sService endpoint;
    EndpointSecurity endpointSecurity?;
    Certificate certificate?;
    Resiliency resiliency?;
};

public type BackendJWTPolicy record {
    *BaseOperationPolicy;
    BackendJWTPolicy_parameters parameters?;
};

public type RetryPolicy record {
    int count?;
    int baseIntervalInMillis?;
    int[] statusCodes?;
};

public type APIKeyAuthentication record {
    *Authentication;
    boolean sendTokenToUpstream?;
    string headerName?;
    string queryParamName?;
    boolean headerEnable?;
    boolean queryParamEnable?;
};

public type Certificate record {
    string secretName?;
    string secretKey?;
};
