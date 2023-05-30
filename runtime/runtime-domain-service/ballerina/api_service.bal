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
import wso2/apk_common_lib as commons;

@display {
    label: "runtime-api-service",
    id: "runtime-api-service"
}

service /api/runtime on ep0 {
    # Retrieve/Search APIs
    #
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + sortBy - Criteria for sorting. 
    # + sortOrder - Order of sorting(ascending/descending). 
    # + query - parameter description 
    # + return - returns can be any of following types
    # APIList (OK. List of qualifying APIs is returned.)
    # InternalServerErrorError (Internal Server Error.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    isolated resource function get apis(http:RequestContext requestContext, string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc") returns APIList|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIList(query, 'limit, offset, sortBy, sortOrder, organization);
    }
    # Create a New API
    #
    # + payload - API object that needs to be added 
    # + return - returns can be any of following types
    # API (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis(http:RequestContext requestContext, @http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.createAPI(payload, (), organization, user);
    }
    # Get Details of an API
    #
    # + apiId - **API ID** consisting of the **Name** of the API. 
    # + return - returns can be any of following types
    # API (OK. Requested API is returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get apis/[string apiId](http:RequestContext requestContext) returns API|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIById(apiId, organization);
    }
    # Update an API
    #
    # + apiId - **API ID** consisting of the **Name** of the API. 
    # + payload - API object that needs to be added 
    # + return - returns can be any of following types
    # API (OK. Successful response with updated API object)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # ForbiddenError (Forbidden. The request must be conditional but no condition has been specified.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function put apis/[string apiId](http:RequestContext requestContext, @http:Payload API payload) returns API|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.updateAPI(apiId, payload, (), organization, user);
    }
    # Delete an API
    #
    # + apiId - **API ID** consisting of the **Name** of the API. 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # ForbiddenError (Forbidden. The request must be conditional but no condition has been specified.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function delete apis/[string apiId](http:RequestContext requestContext) returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.deleteAPIById(apiId, organization);
    }
    # Generate internal API Key to invoke APIS.
    #
    # + apiId - **API ID** consisting of the **Name** of the API. 
    # + return - returns can be any of following types
    # OkAPIKey (OK. apikey generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # ForbiddenError (Forbidden. The request must be conditional but no condition has been specified.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/[string apiId]/'generate\-key(http:RequestContext requestContext) returns http:Ok|BadRequestError|NotFoundError|ForbiddenError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.generateAPIKey(apiId, organization);
    }
    # Create API from a Service
    #
    # + serviceKey - ID of service that should be imported from Service Catalog 
    # + payload - parameter description 
    # + return - returns can be any of following types
    # API (Created. Successful response with the newly created object as entity in the body. Location header contains the URL of the newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/'import\-service(http:RequestContext requestContext, string serviceKey, @http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.createAPIFromService(serviceKey, payload, organization, user);
    }
    # Import an API Definition
    #
    # + request - parameter description 
    # + return - returns can be any of following types
    # API (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/'import\-definition(http:RequestContext requestContext, http:Request message) returns CreatedAPI|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.importDefinition(message, organization, user);
    }
    # Validate an OpenAPI Definition
    #
    # + returnContent - Specify whether to return the full content of the OpenAPI definition in the response. This is only applicable when using url based validation 
    # + request - parameter description 
    # + return - returns can be any of following types
    # OkAPIDefinitionValidationResponse (OK. API definition validation information is returned)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/'validate\-definition(http:RequestContext requestContext, http:Request message, boolean returnContent = false) returns http:Ok|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        return apiService.validateDefinition(message, returnContent);
    }
    # Check Given API Context Name already Exists
    #
    # + query - You can search in attributes by using an **"<attribute>:"** modifier. Eg."name:wso2" will match an API if the provider of the API is exactly "wso2". Supported attribute modifiers are [** version, context, name **]. If no advanced attribute modifier has been specified, search will match the given query string against API Name. 
    # + return - returns can be any of following types
    # http:Ok (OK. API definition validation information is returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/validate(http:RequestContext requestContext, string query) returns http:Ok|NotFoundError|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.validateAPIExistence(query, organization);
    }
    # Get API Definition
    #
    # + apiId - **API ID** consisting of the **Name** of the API. 
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + return - returns can be any of following types
    # string (OK. Requested definition document of the API is returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get apis/[string apiId]/definition(http:RequestContext requestContext,@http:Header string? accept = APPLICATION_JSON_MEDIA_TYPE) returns http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIDefinitionByID(apiId, organization, accept);
    }
    # Update API Definition
    #
    # + apiId - **API ID** consisting of the **Name** of the API. 
    # + request - parameter description 
    # + return - returns can be any of following types
    # string (OK. Successful response with updated Swagger definition)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # ForbiddenError (Forbidden. The request must be conditional but no condition has been specified.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function put apis/[string apiId]/definition(http:RequestContext requestContext, http:Request message) returns http:Response|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.updateAPIDefinition(apiId, message, organization, user);
    }
    # Export an API
    #
    # + apiId - Name of the API 
    # + name - API Name 
    # + version - Version of the API 
    # + format - Format of output documents. Can be YAML or JSON. 
    # + return - returns can be any of following types
    # anydata (OK. Export Successful.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get apis/export(http:RequestContext requestContext, string? apiId, string? name, string? 'version, string? format) returns http:Response|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.exportAPI(apiId, organization);
    }
    # Create a New API Version
    #
    # + newVersion - Version of the new API. 
    # + serviceId - Version of the Service that will used in creating new version 
    # + apiId - **API ID** consisting of the **UUID** of the API. The combination of the provider of the API, name of the API and the version is also accepted as a valid API I. Should be formatted as **provider-name-version**. 
    # + return - returns can be any of following types
    # API (Created. Successful response with the newly created API as entity in the body. Location header contains URL of newly created API.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function post apis/'copy\-api(http:RequestContext requestContext, string newVersion, string? serviceId, string apiId) returns CreatedAPI|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.copyAPI(newVersion, serviceId, apiId, organization, user);
    }
    # Retrieve/search services
    #
    # + query - Search K8s Services based on name or namespace. 
    # + sortBy - Criteria for sorting. 
    # + sortOrder - Order of sorting(ascending/descending). 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + return - returns can be any of following types
    # ServiceList (Paginated matched list of services returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get services(http:RequestContext requestContext, string? query, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns ServiceList|BadRequestError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServices(query, sortBy, sortOrder, 'limit, offset, organization);
    }
    # Get details of a service
    #
    # + serviceId - UUID (unique across all namespaces) of the service 
    # + return - returns can be any of following types
    # Service (Requested service in the Service Catalog is returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get services/[string serviceId](http:RequestContext requestContext) returns Service|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServiceById(serviceId, organization);
    }
    # Retrieve the usage of service
    #
    # + serviceId - UUID(unique id across cluster) of the service 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + sortBy - Criteria for sorting. 
    # + sortOrder - Order of sorting(ascending/descending). 
    # + query - parameter description 
    # + return - returns can be any of following types
    # APIList (List of APIs that uses the service in the Service Catalog is returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get services/[string serviceId]/usage(http:RequestContext requestContext, string? query, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns APIList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServiceUsageByServiceId(query, 'limit, offset, sortBy, sortOrder, serviceId, organization);
    }
    # Get all common mediation policies to all the APIs
    #
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + sortBy - Criteria for sorting. 
    # + sortOrder - Order of sorting(ascending/descending). 
    # + query - parameter description 
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + return - returns can be any of following types
    # MediationPolicyList (OK. List of qualifying policies is returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get policies(http:RequestContext requestContext, string? query, int 'limit = 25, int offset = 0, string sortBy = "id", string sortOrder = "asc", @http:Header string? accept = "application/json") returns MediationPolicyList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getMediationPolicyList(query, 'limit, offset, sortBy, sortOrder, organization);
    }
    # Get the details of a common mediation policy by providing mediation policy ID
    #
    # + policyId - Mediation policy Id 
    # + return - returns can be any of following types
    # MediationPolicy (OK. Mediation policy returned.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get policies/[string policyId](http:RequestContext requestContext) returns MediationPolicy|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getMediationPolicyById(policyId, organization);
    }
    # Retrieve/Search Uploaded Certificates
    #
    # + apiId - parameter description 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + endpoint - Endpoint of which the certificate is uploaded 
    # + return - returns can be any of following types
    # Certificates (OK. Successful response with the list of matching certificate information in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get apis/[string apiId]/'endpoint\-certificates(http:RequestContext requestContext, string? endpoint, int? 'limit, int? offset) returns Certificates|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        int finalLimit = 'limit is () ? 25 : 'limit;
        int finalOffset = offset is () ? 0 : offset;
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getCertificates(apiId, endpoint, finalLimit, finalOffset, organization);
    }
    # Upload a new Certificate.
    #
    # + apiId - parameter description 
    # + request - parameter description 
    # + return - returns can be any of following types
    # OkCertMetadata (OK. The Certificate added successfully.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function post apis/[string apiId]/'endpoint\-certificates(http:RequestContext requestContext, http:Request request) returns OkCertMetadata|BadRequestError|InternalServerErrorError|NotFoundError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.addCertificate(apiId, request, organization);
    }
    # Get the Certificate Information
    #
    # + apiId - parameter description 
    # + certificateId - parameter description 
    # + return - returns can be any of following types
    # CertificateInfo (OK.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get apis/[string apiId]/'endpoint\-certificates/[string certificateId](http:RequestContext requestContext) returns CertificateInfo|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getEndpointCertificateByID(apiId,certificateId, organization);
    }
    # Update a certificate.
    #
    # + apiId - parameter description 
    # + certificateId - parameter description 
    # + request - parameter description 
    # + return - returns can be any of following types
    # CertMetadata (OK. The Certificate updated successfully.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    resource function put apis/[string apiId]/'endpoint\-certificates/[string certificateId](http:RequestContext requestContext,http:Request request) returns OkCertMetadata|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.updateEndpointCertificate(apiId,certificateId, request, organization);
    }
    # Delete a certificate.
    #
    # + apiId - parameter description 
    # + certificateId - parameter description 
    # + return - returns can be any of following types
    # http:Ok (OK. The Certificate deleted successfully.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    resource function delete apis/[string apiId]/'endpoint\-certificates/[string certificateId](http:RequestContext requestContext) returns http:Ok|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.deleteEndpointCertificate(apiId,certificateId, organization);
    }
    # Download a Certificate
    #
    # + apiId - parameter description 
    # + certificateId - parameter description 
    # + return - returns can be any of following types
    # http:Ok (OK.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    resource function get apis/[string apiId]/'endpoint\-certificates/[string certificateId]/content(http:RequestContext requestContext) returns http:Response|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getEndpointCertificateContent(apiId,certificateId, organization);
    }
}
