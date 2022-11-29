import ballerina/http;
import ballerina/constraint;

public type UnsupportedMediaTypeError record {|
    *http:UnsupportedMediaType;
    Error body;
|};

public type ForbiddenError record {|
    *http:Forbidden;
    Error body;
|};

public type ConflictError record {|
    *http:Conflict;
    Error body;
|};

public type CreatedAPI record {|
    *http:Created;
    API body;
|};

public type PreconditionFailedError record {|
    *http:PreconditionFailed;
    Error body;
|};

public type NotFoundError record {|
    *http:NotFound;
    Error body;
|};

public type BadRequestError record {|
    *http:BadRequest;
    Error body;
|};

public type ApiidDefinitionBody record {
    # Swagger definition of the API
    string apiDefinition?;
    # Swagger definition URL of the API
    string url?;
    # Swagger definitio as a file
    string file?;
};

public type ApisValidategraphqlschemaBody record {
    # Definition to upload as a file
    string file;
};

# Summary of the GraphQL including the basic information
public type GraphqlvalidationresponseGraphqlinfo record {
    GraphQLSchema graphQLSchema?;
};

public type ErrorListItem record {
    string code;
    # Description about individual errors occurred
    string message;
    # A detail description about the error message.
    string description?;
};

public type AsyncAPISpecificationValidationResponse record {
    # This attribute declares whether this definition is valid or not.
    boolean isValid;
    # AsyncAPI specification content
    string content?;
    # API definition information
    AsyncapispecificationvalidationresponseInfo info?;
    # If there are more than one error list them out. For example, list out validation error by each field.
    ErrorListItem[] errors?;
};

public type ApiAdditionalpropertiesmap record {
    string name?;
    string value?;
    boolean display?;
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

public type ApiServiceinfo record {
    string 'key?;
    string name?;
    string 'version?;
    boolean outdated?;
};

public type GraphQLSchema record {
    string name;
    string schemaDefinition?;
};

public type ApisValidateopenapiBody record {
    # OpenAPI definition url
    string url?;
    # OpenAPI definition as a file
    string file?;
    # Inline content of the OpenAPI definition
    string inlineAPIDefinition?;
};

public type ApisValidatewsdlBody record {
    # Definition url
    string url?;
    # Definition to upload as a file
    string file?;
};

# API definition information
public type AsyncapispecificationvalidationresponseInfo record {
    string name?;
    string 'version?;
    string context?;
    string description?;
    string asyncAPIVersion?;
    string protocol?;
    # contains host/servers specified in the AsyncAPI file/URL
    string[] endpoints?;
    string gatewayVendor?;
    # contains available transports for an async API
    string[] asyncTransportProtocols?;
};

# API definition information
public type OpenapidefinitionvalidationresponseInfo record {
    string name?;
    string 'version?;
    string context?;
    string description?;
    string openAPIVersion?;
    # contains host/servers specified in the OpenAPI file/URL
    string[] endpoints?;
};

public type WSDLInfo record {
    # Indicates whether the WSDL is a single WSDL or an archive in ZIP format
    string 'type?;
};

public type OpenAPIDefinitionValidationResponse record {
    # This attribute declares whether this definition is valid or not.
    boolean isValid;
    # OpenAPI definition content.
    string content?;
    # API definition information
    OpenapidefinitionvalidationresponseInfo info?;
    # If there are more than one error list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] errors?;
};

public type WsdlvalidationresponseWsdlinfoEndpoints record {
    # Name of the endpoint
    string name?;
    # Endpoint URL
    string location?;
};

public type WSDLValidationResponse record {
    # This attribute declares whether this definition is valid or not.
    boolean isValid;
    # If there are more than one error list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] errors?;
    # Summary of the WSDL including the basic information
    WsdlvalidationresponseWsdlinfo wsdlInfo?;
};

public type API record {
    # UUID of the api registry artifact
    string id?;
    @constraint:String {maxLength: 60, minLength: 1}
    string name;
    @constraint:String {maxLength: 32766}
    string description?;
    @constraint:String {maxLength: 232, minLength: 1}
    string context;
    @constraint:String {maxLength: 30, minLength: 1}
    string 'version;
    # If the provider value is not given user invoking the api will be used as the provider.
    @constraint:String {maxLength: 50}
    string provider?;
    string lifeCycleStatus?;
    WSDLInfo wsdlInfo?;
    string wsdlUrl?;
    boolean responseCachingEnabled?;
    int cacheTimeout?;
    boolean hasThumbnail?;
    boolean isDefaultVersion?;
    boolean isRevision?;
    # UUID of the api registry artifact
    string revisionedApiId?;
    int revisionId?;
    boolean enableSchemaValidation?;
    # The api creation type to be used. Accepted values are HTTP, WS, SOAPTOREST, GRAPHQL, WEBSUB, SSE, WEBHOOK, ASYNC
    string 'type = "HTTP";
    # The audience of the API. Accepted values are PUBLIC, SINGLE
    string audience?;
    # Supported transports for the API (http and/or https).
    string[] transport?;
    string[] tags?;
    string[] policies?;
    # The API level throttling policy selected for the particular API
    string apiThrottlingPolicy?;
    # Name of the Authorization header used for invoking the API. If it is not set, Authorization header name specified
    # in tenant or system level will be used.
    string authorizationHeader?;
    # Types of API security, the current API secured with. It can be either OAuth2 or mutual SSL or both. If
    # it is not set OAuth2 will be set as the security for the current API.
    string[] securityScheme?;
    # The visibility level of the API. Accepts one of the following. PUBLIC, PRIVATE, RESTRICTED.
    string visibility = "PUBLIC";
    # The user roles that are able to access the API in Developer Portal
    string[] visibleRoles?;
    string[] visibleTenants?;
    # The subscription availability. Accepts one of the following. CURRENT_TENANT, ALL_TENANTS or SPECIFIC_TENANTS.
    string subscriptionAvailability = "CURRENT_TENANT";
    string[] subscriptionAvailableTenants?;
    # Map of custom properties of API
    ApiAdditionalproperties[] additionalProperties?;
    record {} additionalPropertiesMap?;
    # Is the API is restricted to certain set of publishers or creators or is it visible to all the
    # publishers and creators. If the accessControl restriction is none, this API can be modified by all the
    # publishers and creators, if not it can only be viewable/modifiable by certain set of publishers and creators,
    #  based on the restriction.
    string accessControl = "NONE";
    # The user roles that are able to view/modify as API publisher or creator.
    string[] accessControlRoles?;
    string workflowStatus?;
    string createdTime?;
    string lastUpdatedTime?;
    # Endpoint configuration of the API. This can be used to provide different types of endpoints including Simple REST Endpoints, Loadbalanced and Failover.
    # 
    # `Simple REST Endpoint`
    #   {
    #     "endpoint_type": "http",
    #     "sandbox_endpoints":       {
    #        "url": "https://localhost:9443/am/sample/pizzashack/v3/api/"
    #     },
    #     "production_endpoints":       {
    #        "url": "https://localhost:9443/am/sample/pizzashack/v3/api/"
    #     }
    #   }
    # 
    # `Loadbalanced Endpoint`
    # 
    #   {
    #     "endpoint_type": "load_balance",
    #     "algoCombo": "org.apache.synapse.endpoints.algorithms.RoundRobin",
    #     "sessionManagement": "",
    #     "sandbox_endpoints":       [
    #                 {
    #           "url": "https://localhost:9443/am/sample/pizzashack/v3/api/1"
    #        },
    #                 {
    #           "endpoint_type": "http",
    #           "template_not_supported": false,
    #           "url": "https://localhost:9443/am/sample/pizzashack/v3/api/2"
    #        }
    #     ],
    #     "production_endpoints":       [
    #                 {
    #           "url": "https://localhost:9443/am/sample/pizzashack/v3/api/3"
    #        },
    #                 {
    #           "endpoint_type": "http",
    #           "template_not_supported": false,
    #           "url": "https://localhost:9443/am/sample/pizzashack/v3/api/4"
    #        }
    #     ],
    #     "sessionTimeOut": "",
    #     "algoClassName": "org.apache.synapse.endpoints.algorithms.RoundRobin"
    #   }
    # 
    # `Failover Endpoint`
    # 
    #   {
    #     "production_failovers":[
    #        {
    #           "endpoint_type":"http",
    #           "template_not_supported":false,
    #           "url":"https://localhost:9443/am/sample/pizzashack/v3/api/1"
    #        }
    #     ],
    #     "endpoint_type":"failover",
    #     "sandbox_endpoints":{
    #        "url":"https://localhost:9443/am/sample/pizzashack/v3/api/2"
    #     },
    #     "production_endpoints":{
    #        "url":"https://localhost:9443/am/sample/pizzashack/v3/api/3"
    #     },
    #     "sandbox_failovers":[
    #        {
    #           "endpoint_type":"http",
    #           "template_not_supported":false,
    #           "url":"https://localhost:9443/am/sample/pizzashack/v3/api/4"
    #        }
    #     ]
    #   }
    # 
    # `Default Endpoint`
    # 
    #   {
    #     "endpoint_type":"default",
    #     "sandbox_endpoints":{
    #        "url":"default"
    #     },
    #     "production_endpoints":{
    #        "url":"default"
    #     }
    #   }
    # 
    # `Endpoint from Endpoint Registry`
    #   {
    #     "endpoint_type": "Registry",
    #     "endpoint_id": "{registry-name:entry-name:version}",
    #   }
    record {} endpointConfig?;
    string endpointImplementationType = "ENDPOINT";
    ApiThreatprotectionpolicies threatProtectionPolicies?;
    # API categories
    string[] categories?;
    # API Key Managers
    record {} keyManagers?;
    ApiServiceinfo serviceInfo?;
    string gatewayVendor?;
    # The gateway type selected for the API policies. Accepts one of the following. wso2/synapse, wso2/choreo-connect.
    string gatewayType?;
    # Supported transports for the async API (http and/or https).
    string[] asyncTransportProtocols?;
};

# Summary of the WSDL including the basic information
public type WsdlvalidationresponseWsdlinfo record {
    # WSDL version
    string 'version?;
    # A list of endpoints the service exposes
    WsdlvalidationresponseWsdlinfoEndpoints[] endpoints?;
};

public type GraphQLValidationResponse record {
    # This attribute declares whether this definition is valid or not.
    boolean isValid;
    # This attribute declares the validation error message
    string errorMessage;
    # Summary of the GraphQL including the basic information
    GraphqlvalidationresponseGraphqlinfo graphQLInfo?;
};

public type ApiThreatprotectionpoliciesList record {
    string policyId?;
    int priority?;
};

public type ApiAdditionalproperties record {
    string name?;
    string value?;
    boolean display?;
};

public type ApiThreatprotectionpolicies record {
    ApiThreatprotectionpoliciesList[] list?;
};
