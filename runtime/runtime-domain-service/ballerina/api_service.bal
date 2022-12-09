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

@display {
    label: "runtime-api-service",
    id: "runtime-api-service"
}

http:Service runtimeService = service object {
    APIClient apiService = new ();
    resource function get apis(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns APIList|BadRequestError|UnauthorizedError|InternalServerErrorError|error {
        APIClient apiService = new ();
        return apiService.getAPIList(query, 'limit, offset, sortBy, sortOrder);
    }
    resource function get apis/[string apiId]() returns API|BadRequestError|InternalServerErrorError|NotFoundError {
        APIClient apiService = new ();
        return apiService.getAPIById(apiId);
    }
    resource function post apis(@http:Payload API payload) returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function put apis/[string apiId](@http:Payload API payload) returns http:NotImplemented|API|BadRequestError|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function delete apis/[string apiId]() returns http:Ok|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
        APIClient apiService = new ();
        return apiService.deleteAPIById(apiId);
    }
    resource function post apis/'import\-service(string serviceKey, @http:Payload API payload) returns CreatedAPI|NotFoundError|InternalServerErrorError|ConflictError {
        APIClient apiService = new ();
        return apiService.createAPIFromService(serviceKey, payload);
    }
    resource function post apis/'import\-definition(@http:Payload json payload) returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post apis/'validate\-definition(@http:Payload json payload, boolean returnContent = false) returns APIDefinitionValidationResponse|BadRequestError|NotFoundError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post apis/validate() returns http:Ok|BadRequestError|NotFoundError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/[string apiId]/definition() returns string|NotFoundError|NotAcceptableError {
        APIClient apiService = new ();
        return apiService.getAPIDefinitionByID(apiId);
    }
    resource function put apis/[string apiId]/definition(@http:Payload json payload) returns string|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get apis/export(string? apiId, string? name, string? 'version, string? format) returns json|NotFoundError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post apis/'import(boolean? overwrite, @http:Payload json payload) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|InternalServerErrorError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function get services(string? name, string? namespace, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns ServiceList|BadRequestError|UnauthorizedError|InternalServerErrorError {
        ServiceClient serviceClient = new ();
        return serviceClient.getServices(name, namespace, sortBy, sortOrder, 'limit, offset);
    }
    resource function get services/[string serviceId](string? namespace) returns Service|BadRequestError|NotFoundError|InternalServerErrorError {
        ServiceClient serviceClient = new ();
        return serviceClient.getServiceById(serviceId);
    }

    resource function get services/[string serviceId]/usage(string? namespace) returns APIList|BadRequestError|NotFoundError|InternalServerErrorError {
        ServiceClient serviceClient = new ();
        return serviceClient.getServiceUsageByServiceId(serviceId);
    }
    resource function get policies(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns MediationPolicyList|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;

    }
    resource function get policies/[string policyId]() returns MediationPolicy|NotFoundError|NotAcceptableError|http:NotImplemented {
        http:NotImplemented notImplementedError = {body: {code: 900910, message: "Not implemented"}};
        return notImplementedError;
    }
    resource function post apis/[string apiId]/'generate\-key() returns APIKey|BadRequestError|NotFoundError|InternalServerErrorError {
        APIClient apiService = new ();
        return apiService.generateAPIKey(apiId);
    }

};

