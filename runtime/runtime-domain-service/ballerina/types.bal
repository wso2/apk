//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
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
import ballerina/constraint;

public type InternalServerErrorError record {|
    *http:InternalServerError;
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

public type ErrorListItem record {
    string code;
    # Description about individual errors occurred
    string message;
    # A detail description about the error message.
    string description?;
};

public type MediationPolicy record {
    string id;
    string 'type;
    string name;
    string displayName?;
    string description?;
    string[] applicableFlows?;
    string[] supportedApiTypes?;
    MediationPolicySpecAttribute[] policyAttributes?;
};

public type Apis_importdefinition_body record {
    # Type of Definition.
    string 'type?;
    # Definition to upload as a file
    string file?;
    # Definition url
    string url?;
    # Additional attributes specified as a stringified JSON with API's schema
    string additionalProperties?;
    # Inline content of the API definition
    string inlineAPIDefinition?;
};

public type Apis_validatedefinition_body record {
    # API definition definition url
    string url?;
    # API definition as a file
    string file?;
    # API definition type - OpenAPI/AsyncAPI/GraphQL
    string 'type?;
    # Inline content of the API definition
    string inlineAPIDefinition?;
};

public type Pagination record {
    int offset?;
    int 'limit?;
    int total?;
    # Link to the next subset of resources qualified.
    # Empty if no more resources are to be returned.
    string next?;
    # Link to the previous subset of resources qualified.
    # Empty if current subset is the first subset returned.
    string previous?;
};

public type GatewayList record {
    Gateway[] list?;
    Pagination pagination?;
};

public type ApiId_definition_body record {
    # API definition of the API
    string apiDefinition?;
    # API definition URL of the API
    string url?;
    # API definitio as a file
    string file?;
};

public type Gateway record {
    # Name of the Gateway
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    # Protocol of the Listener
    @constraint:String {maxLength: 50, minLength: 1}
    string protocol;
    # Port of the Listener
    decimal port;
};

public type APIOperations record {
    string target?;
    string verb?;
    # Authentication mode for resource (true/false)
    boolean authTypeEnabled?;
    # Endpoint configuration of the API. This can be used to provide different types of endpoints including Simple REST Endpoints, Loadbalanced and Failover.
    # 
    # `Simple REST Endpoint`
    #   {
    #     "endpoint_type": "http",
    #     "sandbox_endpoints":       {
    #        "url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"
    #     },
    #     "production_endpoints":       {
    #        "url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"
    #     }
    #   }
    record {} endpointConfig?;
    string[] scopes?;
    APIOperationPolicies operationPolicies?;
};

public type APIOperationPolicies record {
    OperationPolicy[] request?;
    OperationPolicy[] response?;
    OperationPolicy[] fault?;
};

public type MediationPolicySpecAttribute record {
    # Name of the attibute
    string name?;
    # Description of the attibute
    string description?;
    # Is this option mandetory for the policy
    boolean required?;
    # UI validation regex for the attibute
    string validationRegex?;
    # Type of the attibute
    string 'type?;
    # Default value for the attribute
    string defaultValue?;
};

public type APIList record {
    # Number of APIs returned.
    int count?;
    APIInfo[] list?;
    Pagination pagination?;
};

public type MediationPolicyList record {
    # Number of mediation policies returned.
    int count?;
    MediationPolicy[] list?;
    Pagination pagination?;
};

public type PortMapping record {
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    string protocol?;
    int targetport;
    int port;
};

# API definition information
public type APIDefinitionValidationResponse_info record {
    # Name of the API
    string name?;
    # Version of the API
    string 'version?;
    # Context of the API
    string context?;
    # Description of the API
    string description?;
    # OpenAPI Version.
    string openAPIVersion?;
    # contains host/servers specified in the API definition file/URL
    string[] endpoints?;
};

public type ServiceList record {
    Service[] list?;
    Pagination pagination?;
};

public type API_serviceInfo record {
    string name?;
    string namespace?;
};

public type APIInfo record {
    # UUID of the API
    string id?;
    string name?;
    string context?;
    string 'version?;
    string 'type?;
    string createdTime?;
    string updatedTime?;
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

public type Service record {
    @constraint:String {maxLength: 255, minLength: 1}
    string id;
    @constraint:String {maxLength: 255, minLength: 1}
    string name;
    @constraint:String {maxLength: 255}
    string namespace;
    string 'type;
    PortMapping[] portmapping?;
    string createdTime?;
};

public type SearchResult record {
    string id?;
    string name;
    # Accepted values are HTTP, WS, GRAPHQL
    string transportType?;
};

public type GraphQLSchema record {
    string name;
    string schemaDefinition?;
};

public type APIDefinitionValidationResponse record {
    # This attribute declares whether this definition is valid or not.
    boolean isValid;
    # OpenAPI definition content.
    string content?;
    # API definition information
    APIDefinitionValidationResponse_info info?;
    # If there are more than one error list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] errors?;
};

public type APIKey record {
    # API Key
    string apikey?;
    int validityTime?;
};

public type OperationPolicy record {
    string policyName;
    string policyVersion = "v1";
    string policyId?;
    OperationPolicyParameters[] parameters?;
};

public type OperationPolicyParameters record {
    string headerName?;
    string headerValue?;
};

public type Apis_import_body record {
    # Zip archive consisting on exported API configuration
    string file;
};

public type API record {
    # UUID of the API
    string id?;
    @constraint:String {maxLength: 60, minLength: 1}
    string name;
    @constraint:String {maxLength: 232, minLength: 1}
    string context;
    @constraint:String {maxLength: 30, minLength: 1}
    string 'version;
    # The api creation type to be used. Accepted values are REST, WS, GRAPHQL, WEBSUB, SSE, WEBHOOK, ASYNC
    string 'type = "REST";
    # Endpoint configuration of the API. This can be used to provide different types of endpoints including Simple REST Endpoints, Loadbalanced and Failover.
    # 
    # `Simple REST Endpoint`
    #   {
    #     "endpoint_type": "http",
    #     "sandbox_endpoints":       {
    #        "url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"
    #     },
    #     "production_endpoints":       {
    #        "url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"
    #     }
    #   }
    record {} endpointConfig?;
    APIOperations[] operations?;
    API_serviceInfo serviceInfo?;
    APIOperationPolicies apiPolicies?;
    string createdTime?;
    string lastUpdatedTime?;
};
