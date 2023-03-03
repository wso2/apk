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

http:Service runtimeService = service object {
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
        return apiService.createAPI(payload, (), organization);
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
        return apiService.updateAPI(apiId, payload,(), organization);
    }
    isolated resource function delete apis/[string apiId](http:RequestContext requestContext) returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.deleteAPIById(apiId, organization);
    }
    isolated resource function post apis/[string apiId]/'generate\-key(http:RequestContext requestContext) returns APIKey|BadRequestError|NotFoundError|ForbiddenError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.generateAPIKey(apiId, organization);
    }
    isolated resource function post apis/'import\-service(http:RequestContext requestContext, string serviceKey, @http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.createAPIFromService(serviceKey, payload, organization);
    }
    isolated resource function post apis/'import\-definition(http:RequestContext requestContext, http:Request message) returns CreatedAPI|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.importDefinition(message, organization);
    }
    isolated resource function post apis/'validate\-definition(http:RequestContext requestContext, http:Request message, boolean returnContent = false) returns APIDefinitionValidationResponse|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        return apiService.validateDefinition(message, returnContent);
    }
    isolated resource function post apis/validate(http:RequestContext requestContext, string query) returns http:Ok|NotFoundError|BadRequestError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.validateAPIExistence(query, organization);
    }
    isolated resource function get apis/[string apiId]/definition(http:RequestContext requestContext) returns http:Response|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getAPIDefinitionByID(apiId, organization);
    }
    isolated resource function put apis/[string apiId]/definition(http:RequestContext requestContext, http:Request message) returns http:Response|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError|commons:APKError {
         final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.updateAPIDefinition(apiId,message,organization);
    }
    isolated resource function get apis/export(http:RequestContext requestContext, string? apiId, string? name, string? 'version, string? format) returns http:Response|NotFoundError|InternalServerErrorError|BadRequestError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.exportAPI(apiId, organization);
    }
    isolated resource function post apis/'import(boolean? overwrite, @http:Payload json payload) returns http:Ok|ForbiddenError|ConflictError|PreconditionFailedError|InternalServerErrorError {
        InternalServerErrorError internalError = {body: {code: 900910, message: "Not implemented"}};
        return internalError;
    }
    isolated resource function post apis/'copy\-api(http:RequestContext requestContext, string newVersion, string? serviceId, string apiId) returns CreatedAPI|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.copyAPI(newVersion, serviceId, apiId, organization);
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
    isolated resource function get services/[string serviceId]/usage(http:RequestContext requestContext) returns APIList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final ServiceClient serviceClient = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return serviceClient.getServiceUsageByServiceId(serviceId, organization);
    }
    isolated resource function get policies(http:RequestContext requestContext, string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns MediationPolicyDataList|BadRequestError|NotFoundError|InternalServerErrorError|commons:APKError {
        final APIClient apiService = new ();
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return apiService.getMediationPolicyList(query, 'limit, offset, sortBy, sortOrder, organization);
    }
    isolated resource function get policies/[string policyId]() returns MediationPolicy|NotFoundError|InternalServerErrorError {
        InternalServerErrorError internalError = {body: {code: 900910, message: "Not implemented"}};
        return internalError;
    }
};
