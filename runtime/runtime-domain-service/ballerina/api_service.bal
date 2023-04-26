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

service /api/am/runtime on ep0 {
    isolated resource function get apis(http:RequestContext requestContext, string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc") returns APIList|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIList(query, 'limit, offset, sortBy, sortOrder, organization);
    }
    isolated resource function post apis(http:RequestContext requestContext, @http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.createAPI(payload, (), organization, user);
    }
    isolated resource function get apis/[string apiId](http:RequestContext requestContext) returns API|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIById(apiId, organization);
    }
    isolated resource function put apis/[string apiId](http:RequestContext requestContext, @http:Payload API payload) returns API|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.updateAPI(apiId, payload, (), organization, user);
    }
    isolated resource function delete apis/[string apiId](http:RequestContext requestContext) returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.deleteAPIById(apiId, organization);
    }
    isolated resource function post apis/[string apiId]/'generate\-key(http:RequestContext requestContext) returns http:Ok|BadRequestError|NotFoundError|ForbiddenError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.generateAPIKey(apiId, organization);
    }
    isolated resource function post apis/'import\-service(http:RequestContext requestContext, string serviceKey, @http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.createAPIFromService(serviceKey, payload, organization, user);
    }
    isolated resource function post apis/'import\-definition(http:RequestContext requestContext, http:Request message) returns CreatedAPI|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.importDefinition(message, organization, user);
    }
    isolated resource function post apis/'validate\-definition(http:RequestContext requestContext, http:Request message, boolean returnContent = false) returns http:Ok|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        return apiService.validateDefinition(message, returnContent);
    }
    isolated resource function post apis/validate(http:RequestContext requestContext, string query) returns http:Ok|NotFoundError|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.validateAPIExistence(query, organization);
    }
    isolated resource function get apis/[string apiId]/definition(http:RequestContext requestContext,@http:Header string? accept = APPLICATION_JSON_MEDIA_TYPE) returns http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIDefinitionByID(apiId, organization, accept);
    }
    isolated resource function put apis/[string apiId]/definition(http:RequestContext requestContext, http:Request message) returns http:Response|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.updateAPIDefinition(apiId, message, organization, user);
    }
    isolated resource function get apis/export(http:RequestContext requestContext, string? apiId, string? name, string? 'version, string? format) returns http:Response|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.exportAPI(apiId, organization);
    }
    isolated resource function post apis/'copy\-api(http:RequestContext requestContext, string newVersion, string? serviceId, string apiId) returns CreatedAPI|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string user = authenticatedUserContext.username;
        return apiService.copyAPI(newVersion, serviceId, apiId, organization, user);
    }
    isolated resource function get services(http:RequestContext requestContext, string? query, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns ServiceList|BadRequestError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServices(query, sortBy, sortOrder, 'limit, offset, organization);
    }
    isolated resource function get services/[string serviceId](http:RequestContext requestContext) returns Service|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServiceById(serviceId, organization);
    }
    isolated resource function get services/[string serviceId]/usage(http:RequestContext requestContext, string? query, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns APIList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServiceUsageByServiceId(query, 'limit, offset, sortBy, sortOrder, serviceId, organization);
    }
    isolated resource function get policies(http:RequestContext requestContext, string? query, int 'limit = 25, int offset = 0, string sortBy = "id", string sortOrder = "asc", @http:Header string? accept = "application/json") returns MediationPolicyList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getMediationPolicyList(query, 'limit, offset, sortBy, sortOrder, organization);
    }
    isolated resource function get policies/[string policyId](http:RequestContext requestContext) returns MediationPolicy|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getMediationPolicyById(policyId, organization);
    }
    isolated resource function get apis/[string apiId]/'endpoint\-certificates(http:RequestContext requestContext, string? endpoint, int? 'limit, int? offset) returns Certificates|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        int finalLimit = 'limit is () ? 25 : 'limit;
        int finalOffset = offset is () ? 0 : offset;
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getCertificates(apiId, endpoint, finalLimit, finalOffset, organization);
    }
    isolated resource function post apis/[string apiId]/'endpoint\-certificates(http:RequestContext requestContext, http:Request request) returns OkCertMetadata|BadRequestError|InternalServerErrorError|NotFoundError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.addCertificate(apiId, request, organization);
    }
    isolated resource function get apis/[string apiId]/'endpoint\-certificates/[string certificateId](http:RequestContext requestContext) returns CertificateInfo|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getEndpointCertificateByID(apiId,certificateId, organization);
    }
    resource function put apis/[string apiId]/'endpoint\-certificates/[string certificateId](http:RequestContext requestContext,http:Request request) returns OkCertMetadata|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.updateEndpointCertificate(apiId,certificateId, request, organization);
    }
    resource function delete apis/[string apiId]/'endpoint\-certificates/[string certificateId](http:RequestContext requestContext) returns http:Ok|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.deleteEndpointCertificate(apiId,certificateId, organization);
    }
    resource function get apis/[string apiId]/'endpoint\-certificates/[string certificateId]/content(http:RequestContext requestContext) returns http:Response|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getEndpointCertificateContent(apiId,certificateId, organization);
    }
}
