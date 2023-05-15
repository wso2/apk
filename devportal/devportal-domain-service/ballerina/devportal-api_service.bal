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
import ballerina/log;
import wso2/apk_common_lib as commons;

isolated service /api/am/devportal on ep0 {
    # Retrieve/Search APIs
    #
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + query - You can search in attributes by using an attribute modifier. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # APIList (OK. List of qualifying APIs is returned.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get apis(http:RequestContext requestContext, @http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns APIList|BadRequestError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        APIList apiList = check getAPIList('limit, offset, query, organization);
        log:printDebug(apiList.toString());
        return apiList;
    }
    # Get Details of an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # API (OK. Requested API is returned)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get apis/[string apiId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns API|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|json|commons:APKError {
        API|NotFoundError api = check getAPIByAPIId(apiId);
        log:printDebug(api.toString());
        return api;
    }
    # Get the API Definition
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # APIDefinition (OK. Requested definition document of the API is returned)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get apis/[string apiId]/definition(@http:Header string? 'if\-none\-match) returns APIDefinition|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|commons:APKError {

        APIDefinition|NotFoundError apiDefinition = check getAPIDefinition(apiId);
        log:printDebug(apiDefinition.toString());
        return apiDefinition;
    }
    # Generate a SDK for an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + language - Programming language of the SDK that is required. Languages supported by default are **Java**, **Javascript**, **Android** and **JMeter**. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # anydata (OK. SDK generated successfully.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get apis/[string apiId]/sdks/[string language](@http:Header string? 'x\-wso2\-tenant) returns NotFoundError|http:Response|commons:APKError|InternalServerErrorError {
        return check generateSDKImpl(apiId, language);
    }
    # Get a List of Documents of an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # DocumentList (OK. Document list is returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get apis/[string apiId]/documents(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns DocumentList|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Get a Document of an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + documentId - Document Identifier 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # Document (OK. Document returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get apis/[string apiId]/documents/[string documentId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns Document|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Get the Content of an API Document
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + documentId - Document Identifier 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # http:Ok (OK. File or inline content returned.)
    # http:SeeOther (See Other. Source can be retrieved from the URL specified at the Location header.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get apis/[string apiId]/documents/[string documentId]/content(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns http:Ok|http:SeeOther|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Get Thumbnail Image
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # http:Ok (OK. Thumbnail image returned)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    resource function get apis/[string apiId]/thumbnail(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns http:Response|http:NotModified|NotFoundError|NotAcceptableError|APKError {
        return getThumbnail(apiId);
    }
    # Retrieve API Ratings
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # RatingList (OK. Rating list returned.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get apis/[string apiId]/ratings(@http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0) returns RatingList|NotAcceptableError {
    // }
    # Retrieve API Rating of User
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # Rating (OK. Rating returned.)
    # http:NotModified (Not Modified. Empty body because the client already has the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns Rating|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Add or Update Logged in User's Rating for an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + payload - Rating object that should to be added 
    # + return - returns can be any of following types
    # Rating (OK. Successful response with the newly created or updated object as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    // resource function put apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Payload Rating payload) returns Rating|BadRequestError|UnsupportedMediaTypeError {
    // }
    # Delete User API Rating
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + return - OK. Resource successfully deleted. 
    // resource function delete apis/[string apiId]/'user\-rating(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-match) returns http:Ok {
    // }
    # Retrieve API Comments
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + includeCommenterInfo - Whether we need to display commenter details. 
    # + return - returns can be any of following types
    # CommentList (OK. Comments list is returned.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function get apis/[string apiId]/comments(@http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|NotFoundError|InternalServerErrorError {
    // }
    # Add an API Comment
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + replyTo - ID of the parent comment. 
    # + payload - Comment object that should to be added 
    # + return - returns can be any of following types
    # Comment (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnauthorizedError (Unauthorized. The user is not authorized.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function post apis/[string apiId]/comments(string? replyTo, @http:Payload 'postRequestBody payload) returns Comment|BadRequestError|UnauthorizedError|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError {
    // }
    # Get Details of an API Comment
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + commentId - Comment Id 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + includeCommenterInfo - Whether we need to display commenter details. 
    # + replyLimit - Maximum size of replies array to return. 
    # + replyOffset - Starting point within the complete list of replies. 
    # + return - returns can be any of following types
    # Comment (OK. Comment returned.)
    # UnauthorizedError (Unauthorized. The user is not authorized.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function get apis/[string apiId]/comments/[string commentId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, boolean includeCommenterInfo = false, int replyLimit = 25, int replyOffset = 0) returns Comment|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError {
    // }
    # Delete an API Comment
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + commentId - Comment Id 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # UnauthorizedError (Unauthorized. The user is not authorized.)
    # http:Forbidden (Forbidden. The request must be conditional but no condition has been specified.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # http:MethodNotAllowed (MethodNotAllowed. Request method is known by the server but is not supported by the target resource.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function delete apis/[string apiId]/comments/[string commentId](@http:Header string? 'if\-match) returns http:Ok|UnauthorizedError|http:Forbidden|NotFoundError|http:MethodNotAllowed|InternalServerErrorError {
    // }
    # Edit a comment
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + commentId - Comment Id 
    # + payload - Comment object that should to be updated 
    # + return - returns can be any of following types
    # Comment (OK. Comment updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnauthorizedError (Unauthorized. The user is not authorized.)
    # http:Forbidden (Forbidden. The request must be conditional but no condition has been specified.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function patch apis/[string apiId]/comments/[string commentId](@http:Payload 'patchRequestBody payload) returns Comment|BadRequestError|UnauthorizedError|http:Forbidden|NotFoundError|UnsupportedMediaTypeError|InternalServerErrorError {
    // }
    # Get replies of a comment
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + commentId - Comment Id 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + includeCommenterInfo - Whether we need to display commenter details. 
    # + return - returns can be any of following types
    # CommentList (OK. Comment returned.)
    # UnauthorizedError (Unauthorized. The user is not authorized.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function get apis/[string apiId]/comments/[string commentId]/replies(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0, boolean includeCommenterInfo = false) returns CommentList|UnauthorizedError|NotFoundError|NotAcceptableError|InternalServerErrorError {
    // }
    # Get a list of available topics for a given Async API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # TopicList (OK. Topic list returned.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function get apis/[string apiId]/topics(@http:Header string? 'x\-wso2\-tenant) returns TopicList|NotFoundError|InternalServerErrorError {
    // }
    # Get Details of the Subscription Throttling Policies of an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # ThrottlingPolicy (OK. Throttling Policy returned)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get apis/[string apiId]/'subscription\-policies(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns ThrottlingPolicy|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Retrieve/Search Applications
    #
    # + groupId - Application Group Id 
    # + query - You can search for an application by specifying the name as "query" attribute. 
    # + sortBy - parameter description 
    # + sortOrder - parameter description 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # ApplicationList (OK. Application list returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get applications(http:RequestContext requestContext, string? groupId, string? query, string? sortBy, string? sortOrder, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns ApplicationList|http:NotModified|BadRequestError|NotAcceptableError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        ApplicationList applicationList = check getApplicationList(sortBy, groupId, query, sortOrder, 'limit, offset, organization);
        log:printDebug(applicationList.toString());
        return applicationList;
    }
    # Create a New Application
    #
    # + payload - Application object that is to be created. 
    # + return - returns can be any of following types
    # Application (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # AcceptedWorkflowResponse (Accepted. The request has been accepted.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # ConflictError (Conflict. Specified resource already exists.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    isolated resource function post applications(http:RequestContext requestContext, @http:Payload Application payload) returns CreatedApplication|AcceptedWorkflowResponse|BadRequestError|ConflictError|NotFoundError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        Application|NotFoundError application = check addApplication(payload, organization, <string>authenticatedUserContext.userId);
        if application is Application {
            CreatedApplication createdApp = {body: application};
            log:printDebug(application.toString());
            return createdApp;
        } else {
            return <NotFoundError>application;
        }
    }
    # Get Details of an Application
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # Application (OK. Application returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get applications/[string applicationId](http:RequestContext requestContext, @http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant) returns Application|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return getApplicationById(applicationId, organization);
    }
    # Update an Application
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + payload - Application object that needs to be updated 
    # + return - returns can be any of following types
    # Application (OK. Application updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    isolated resource function put applications/[string applicationId](http:RequestContext requestContext, @http:Header string? 'if\-match, @http:Payload Application payload) returns Application|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        Application|NotFoundError application = check updateApplication(applicationId, payload, organization, <string>authenticatedUserContext.userId);
        if application is Application|NotFoundError {
            log:printDebug(application.toString());
            return application;
        }
    }
    # Remove an Application
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # AcceptedWorkflowResponse (Accepted. The request has been accepted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    isolated resource function delete applications/[string applicationId](http:RequestContext requestContext, @http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        _ = check deleteApplication(applicationId, organization);
        return http:OK;
    }

    # Generate Application Keys
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + payload - Application key generation request object 
    # + return - returns can be any of following types
    # OkApplicationKey (OK. Keys are generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/'generate\-keys(@http:Header string? 'x\-wso2\-tenant, @http:Payload ApplicationKeyGenerateRequest payload, http:RequestContext requestContext) returns OkApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError|commons:APKError {
    //     commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
    //     commons:Organization organization = authenticatedUserContext.organization;
    //     Application|NotFoundError applicationById = check getApplicationById(applicationId, organization);
    //     if applicationById is NotFoundError {
    //         return applicationById;
    //     }
    //     // generateKeysForApplication(<Application>applicationById, payload, organization);
    // }
    # Map Application Keys
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + payload - Application key mapping request object 
    # + return - returns can be any of following types
    # OkApplicationKey (OK. Keys are mapped.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/'map\-keys(@http:Header string? 'x\-wso2\-tenant, @http:Payload ApplicationKeyMappingRequest payload) returns OkApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Retrieve All Application Keys
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + return - returns can be any of following types
    # ApplicationKeyList (OK. Keys are returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function get applications/[string applicationId]/keys() returns ApplicationKeyList|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Get Key Details of a Given Type
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + groupId - Application Group Id 
    # + return - returns can be any of following types
    # ApplicationKey (OK. Keys of given type are returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function get applications/[string applicationId]/keys/[string keyType](string? groupId) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Update Grant Types and Callback Url of an Application
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + payload - Grant types/Callback URL update request object 
    # + return - returns can be any of following types
    # ApplicationKey (Ok. Grant types or/and callback url is/are updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function put applications/[string applicationId]/keys/[string keyType](@http:Payload ApplicationKey payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Re-Generate Consumer Secret
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + return - returns can be any of following types
    # OkApplicationKeyReGenerateResponse (OK. Keys are re generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/keys/[string keyType]/'regenerate\-secret() returns OkApplicationKeyReGenerateResponse|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Clean-Up Application Keys
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + return - returns can be any of following types
    # http:Ok (OK. Clean up is performed)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/keys/[string keyType]/'clean\-up(@http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Generate Application Token
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + payload - Application token generation request object 
    # + return - returns can be any of following types
    # OkApplicationToken (OK. Token is generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/keys/[string keyType]/'generate\-token(@http:Header string? 'if\-match, @http:Payload ApplicationTokenGenerateRequest payload) returns OkApplicationToken|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Retrieve All Application Keys
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # ApplicationKeyList (OK. Keys are returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function get applications/[string applicationId]/'oauth\-keys(@http:Header string? 'x\-wso2\-tenant) returns ApplicationKeyList|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Get Key Details of a Given Type
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyMappingId - OAuth Key Identifier consisting of the UUID of the Oauth Key Mapping. 
    # + groupId - Application Group Id 
    # + return - returns can be any of following types
    # ApplicationKey (OK. Keys of given type are returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function get applications/[string applicationId]/'oauth\-keys/[string keyMappingId](string? groupId) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Update Grant Types and Callback URL of an Application
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyMappingId - OAuth Key Identifier consisting of the UUID of the Oauth Key Mapping. 
    # + payload - Grant types/Callback URL update request object 
    # + return - returns can be any of following types
    # ApplicationKey (Ok. Grant types or/and callback url is/are updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function put applications/[string applicationId]/'oauth\-keys/[string keyMappingId](@http:Payload ApplicationKey payload) returns ApplicationKey|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Re-Generate Consumer Secret
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyMappingId - OAuth Key Identifier consisting of the UUID of the Oauth Key Mapping. 
    # + return - returns can be any of following types
    # OkApplicationKeyReGenerateResponse (OK. Keys are re generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'regenerate\-secret() returns OkApplicationKeyReGenerateResponse|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Clean-Up Application Keys
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyMappingId - OAuth Key Identifier consisting of the UUID of the Oauth Key Mapping. 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + return - returns can be any of following types
    # http:Ok (OK. Clean up is performed)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'clean\-up(@http:Header string? 'if\-match) returns http:Ok|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Generate Application Token
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyMappingId - OAuth Key Identifier consisting of the UUID of the Oauth Key Mapping. 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + payload - Application token generation request object 
    # + return - returns can be any of following types
    # OkApplicationToken (OK. Token is generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/'oauth\-keys/[string keyMappingId]/'generate\-token(@http:Header string? 'if\-match, @http:Payload ApplicationTokenGenerateRequest payload) returns OkApplicationToken|BadRequestError|NotFoundError|PreconditionFailedError {
    // }
    # Generate API Key
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + payload - API Key generation request object 
    # + return - returns can be any of following types
    # OkAPIKey (OK. apikey generated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    isolated resource function post applications/[string applicationId]/'api\-keys/[string keyType]/generate(http:RequestContext requestContext, @http:Header string? 'if\-match, @http:Payload APIKeyGenerateRequest payload) returns APIKey|BadRequestError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        APIKey|NotFoundError apiKey = check generateAPIKey(payload, applicationId, keyType, <string>authenticatedUserContext.userId, organization);
        return apiKey;
    }
    # Revoke API Key
    #
    # + applicationId - Application Identifier consisting of the UUID of the Application. 
    # + keyType - **Application Key Type** standing for the type of the keys (i.e. Production or Sandbox). 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + payload - API Key revoke request object 
    # + return - returns can be any of following types
    # http:Ok (OK. apikey revoked successfully.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    // resource function post applications/[string applicationId]/'api\-keys/[string keyType]/revoke(@http:Header string? 'if\-match, @http:Payload APIKeyRevokeRequest payload) returns http:Ok|BadRequestError|PreconditionFailedError {
    // }
    # Export an Application
    #
    # + appName - Application Name 
    # + appOwner - Owner of the Application 
    # + withKeys - Export application keys 
    # + format - Format of output documents. Can be YAML or JSON. 
    # + return - returns can be any of following types
    # anydata (OK. Export Successful.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get applications/export(string appName, string appOwner, boolean? withKeys, string? format) returns anydata {
    // }
    # Import an Application
    #
    # + preserveOwner - Preserve Original Creator of the Application 
    # + skipSubscriptions - Skip importing Subscriptions of the Application 
    # + appOwner - Expected Owner of the Application in the Import Environment 
    # + skipApplicationKeys - Skip importing Keys of the Application 
    # + update - Update if application exists 
    # + request - parameter description 
    # + return - returns can be any of following types
    # OkApplicationInfo (OK. Successful response with the updated object information as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function post applications/'import(boolean? preserveOwner, boolean? skipSubscriptions, string? appOwner, boolean? skipApplicationKeys, boolean? update, http:Request request) returns OkApplicationInfo|BadRequestError|NotAcceptableError {
    // }
    # Get All Subscriptions
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + applicationId - **Application Identifier** consisting of the UUID of the Application. 
    # + groupId - Application Group Id 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'limit - Maximum size of resource array to return. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # SubscriptionList (OK. Subscription list returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get subscriptions(http:RequestContext requestContext, string? apiId, string? applicationId, string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns SubscriptionList|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        SubscriptionList|NotFoundError subscriptionList = check getSubscriptions(apiId, applicationId, groupId, offset, 'limit, organization);
        log:printDebug(subscriptionList.toString());
        return subscriptionList;
    }
    # Add a New Subscription
    #
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + payload - Subscription object that should to be added 
    # + return - returns can be any of following types
    # Subscription (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # AcceptedWorkflowResponse (Accepted. The request has been accepted.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    isolated resource function post subscriptions(http:RequestContext requestContext, @http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns CreatedSubscription|AcceptedWorkflowResponse|BadRequestError|NotFoundError|InternalServerErrorError|json|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        Subscription|NotFoundError subscription = check addSubscription(payload, organization, <string>authenticatedUserContext.userId);
        if subscription is Subscription {
            CreatedSubscription createdSub = {body: subscription};
            log:printDebug(subscription.toString());
            return createdSub;
        } else if subscription is NotFoundError {
            return subscription;
        }
    }

    # Add New Subscriptions
    #
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + payload - Subscription objects that should to be added 
    # + return - returns can be any of following types
    # OkSubscription (OK. Successful response with the newly created objects as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    isolated resource function post subscriptions/multiple(http:RequestContext requestContext, @http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription[] payload) returns Subscription[]|BadRequestError|UnsupportedMediaTypeError|NotFoundError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        Subscription[]|NotFoundError subscriptions = check addMultipleSubscriptions(payload, organization, <string>authenticatedUserContext.userId);
        log:printDebug(subscriptions.toString());
        return subscriptions;
    }
    # Get Additional Information of subscriptions attached to an API.
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + groupId - Application Group Id 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'limit - Maximum size of resource array to return. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # AdditionalSubscriptionInfoList (OK. Types and fields returned successfully.)
    # http:NotFound (Not Found. Retrieving types and fields failed.)
    // resource function get subscriptions/[string apiId]/additionalInfo(string? groupId, @http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int offset = 0, int 'limit = 25) returns AdditionalSubscriptionInfoList|http:NotFound {
    // }
    # Get Details of a Subscription
    #
    # + subscriptionId - Subscription Id 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # Subscription (OK. Subscription returned)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function get subscriptions/[string subscriptionId](http:RequestContext requestContext, @http:Header string? 'if\-none\-match) returns Subscription|http:NotModified|NotFoundError|BadRequestError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        Subscription|NotFoundError subscription = check getSubscriptionById(subscriptionId, organization);
        log:printDebug(subscription.toString());
        return subscription;
    }
    # Update Existing Subscription
    #
    # + subscriptionId - Subscription Id 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + payload - Subscription object that should to be added 
    # + return - returns can be any of following types
    # Subscription (Subscription Updated. Successful response with the updated object as entity in the body. Location header contains URL of newly updates entity.)
    # AcceptedWorkflowResponse (Accepted. The request has been accepted.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # http:NotFound (Not Found. Requested Subscription does not exist.)
    # http:UnsupportedMediaType (Unsupported media type. The entity of the request was in a not supported format.)
    isolated resource function put subscriptions/[string subscriptionId](http:RequestContext requestContext, @http:Header string? 'x\-wso2\-tenant, @http:Payload Subscription payload) returns Subscription|AcceptedWorkflowResponse|http:NotModified|BadRequestError|NotFoundError|http:UnsupportedMediaType|InternalServerErrorError|json|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        Subscription|NotFoundError subscription = check updateSubscription(subscriptionId, payload, organization, <string>authenticatedUserContext.userId);
        log:printDebug(subscription.toString());
        return subscription;
    }
    # Remove a Subscription
    #
    # + subscriptionId - Subscription Id 
    # + 'if\-match - Validator for conditional requests; based on ETag. 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # AcceptedWorkflowResponse (Accepted. The request has been accepted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # PreconditionFailedError (Precondition Failed. The request has not been performed because one of the preconditions is not met.)
    isolated resource function delete subscriptions/[string subscriptionId](http:RequestContext requestContext, @http:Header string? 'if\-match) returns http:Ok|AcceptedWorkflowResponse|NotFoundError|PreconditionFailedError|BadRequestError|InternalServerErrorError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string response = check deleteSubscription(subscriptionId, organization);
        return http:OK;
    }
    # Get Details of a Pending Invoice for a Monetized Subscription with Metered Billing.
    #
    # + subscriptionId - Subscription Id 
    # + return - returns can be any of following types
    # APIMonetizationUsage (OK. Details of a pending invoice returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource (Will be supported in future).)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function get subscriptions/[string subscriptionId]/usage() returns APIMonetizationUsage|http:NotModified|NotFoundError {
    // }
    # Get All Available Throttling Policies
    #
    # + policyLevel - List Application or Subscription type thro. 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # ThrottlingPolicyList (OK. List of throttling policies returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get 'throttling\-policies/[string policyLevel](@http:Header string? 'if\-none\-match, @http:Header string? 'x\-wso2\-tenant, int 'limit = 25, int offset = 0) returns ThrottlingPolicyList|http:NotModified|NotAcceptableError {
    // }
    # Get Details of a Throttling Policy
    #
    # + policyLevel - List Application or Subscription type thro. 
    # + policyId - The name of the policy 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # ThrottlingPolicy (OK. Throttling Policy returned)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get 'throttling\-policies/[string policyLevel]/[string policyId](@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match) returns ThrottlingPolicy|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Get All Tags
    #
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # TagList (OK. Tag list is returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get tags(@http:Header string? 'x\-wso2\-tenant, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns TagList|http:NotModified|NotFoundError|NotAcceptableError {
    // }
    # Retrieve/Search APIs and API Documents by Content
    #
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + query - You can search by using providing the search term in the query parameters. 
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # SearchResultList (OK. List of qualifying APIs and docs is returned.)
    # http:NotModified (Not Modified. Empty body because the client has already the latest version of the requested resource (Will be supported in future).)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get search(@http:Header string? 'x\-wso2\-tenant, string? query, @http:Header string? 'if\-none\-match, int 'limit = 25, int offset = 0) returns SearchResultList|http:NotModified|NotAcceptableError {
    // }
    # Get a List of Supported SDK Languages
    #
    # + return - returns can be any of following types
    # json (OK. List of supported languages for generating SDKs.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    isolated resource function get 'sdk\-gen/languages() returns json|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        string|json sdkLanguages = check getSDKLanguages();
        return sdkLanguages;
    }
    # Get available web hook subscriptions for a given application.
    #
    # + applicationId - **Application Identifier** consisting of the UUID of the Application. 
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # WebhookSubscriptionList (OK. Topic list returned.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function get webhooks/subscriptions(string? applicationId, string? apiId, @http:Header string? 'x\-wso2\-tenant) returns WebhookSubscriptionList|NotFoundError|InternalServerErrorError {
    // }
    # Retrieve Developer Portal settings
    #
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - returns can be any of following types
    # Settings (OK. Settings returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function get settings(@http:Header string? 'x\-wso2\-tenant) returns Settings|NotFoundError {
    // }
    # Get All Application Attributes from Configuration
    #
    # + 'if\-none\-match - Validator for conditional requests; based on the ETag of the formerly retrieved variant of the resource. 
    # + return - returns can be any of following types
    # ApplicationAttributeList (OK. Application attributes returned.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get settings/'application\-attributes(@http:Header string? 'if\-none\-match) returns ApplicationAttributeList|NotFoundError|NotAcceptableError {
    // }
    # Get Tenants by State
    #
    # + state - The state represents the current state of the tenant 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + return - returns can be any of following types
    # TenantList (OK. Tenant names returned.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get tenants(string state = "active", int 'limit = 25, int offset = 0) returns TenantList|NotFoundError|NotAcceptableError {
    // }
    # Give API Recommendations for a User
    #
    # + return - returns can be any of following types
    # Recommendations (OK. Requested recommendations are returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function get recommendations() returns Recommendations|NotFoundError {
    // }
    # Get All API Categories
    #
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - OK. Categories returned 
    // resource function get 'api\-categories(@http:Header string? 'x\-wso2\-tenant) returns APICategoryList {
    // }
    # Get All Key Managers
    #
    # + 'x\-wso2\-tenant - For cross-tenant invocations, this is used to specify the tenant/organization domain, where the resource need to be   retrieved from. 
    # + return - OK. Key Manager list returned 
    // resource function get 'key\-managers(@http:Header string? 'x\-wso2\-tenant) returns KeyManagerList {
    // }
    # Get the Complexity Related Details of an API
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + return - returns can be any of following types
    # GraphQLQueryComplexityInfo (OK. Requested complexity details returned.)
    # http:NotFound (Not Found. Requested API does not contain any complexity details.)
    // resource function get apis/[string apiId]/'graphql\-policies/complexity() returns GraphQLQueryComplexityInfo|http:NotFound {
    // }
    # Retrieve Types and Fields of a GraphQL Schema
    #
    # + apiId - **API ID** consisting of the **UUID** of the API. 
    # + return - returns can be any of following types
    # GraphQLSchemaTypeList (OK. Types and fields returned successfully.)
    # http:NotFound (Not Found. Retrieving types and fields failed.)
    // resource function get apis/[string apiId]/'graphql\-policies/complexity/types() returns GraphQLSchemaTypeList|http:NotFound {
    // }
    # Change the Password of the user
    #
    # + payload - Current and new password of the user 
    # + return - returns can be any of following types
    # http:Ok (OK. User password changed successfully)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    isolated resource function post me/'change\-password(@http:Payload CurrentAndNewPasswords payload) returns http:Ok|BadRequestError {
        BadRequestError badRequest = {body: {code: 400, message: "Invalid request or validation error."}};
        return badRequest;
    }

}
