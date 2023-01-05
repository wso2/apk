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

@display {
    label: "runtime-api-service",
    id: "runtime-api-service"
}

http:Service runtimeService = service object {

    isolated resource function get apis(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc") returns APIList|InternalServerErrorError|BadRequestError {
        APIClient apiService = new ();
        return apiService.getAPIList(query, 'limit, offset, sortBy, sortOrder);
    }
    isolated resource function post apis(@http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError {
        APIClient apiService = new();
        APKError|CreatedAPI|BadRequestError createdAPI = apiService.createAPI(payload);
        if createdAPI is APKError {
           return  handleAPKError(createdAPI);
        } else{
            return createdAPI;
        } 
    }
    isolated resource function get apis/[string apiId]() returns API|NotFoundError|InternalServerErrorError {
        APIClient apiService = new ();
        return apiService.getAPIById(apiId);
    }
    isolated resource function put apis/[string apiId](@http:Payload API payload) returns API|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        BadRequestError badRequest = {body: {code: 900910, message: "Not implemented"}};
        return badRequest;
    }
    isolated resource function delete apis/[string apiId]() returns http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|BadRequestError {
        APIClient apiService = new ();
        http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError|APKError apiDeletionResponse = apiService.deleteAPIById(apiId);
        if apiDeletionResponse is http:Ok|ForbiddenError|NotFoundError|InternalServerErrorError {
            return apiDeletionResponse;
        } else {
            log:printError("Internal Error occured deleting API", apiDeletionResponse);
            return handleAPKError(apiDeletionResponse);
       }
    }
    isolated resource function post apis/[string apiId]/'generate\-key() returns APIKey|BadRequestError|NotFoundError|ForbiddenError|InternalServerErrorError {
        APIClient apiService = new ();
        return apiService.generateAPIKey(apiId);
    }
    isolated resource function post apis/'import\-service(string serviceKey, @http:Payload API payload) returns CreatedAPI|BadRequestError|InternalServerErrorError {
        APIClient apiService = new ();
        CreatedAPI|BadRequestError|InternalServerErrorError|APKError aPIFromService = apiService.createAPIFromService(serviceKey, payload);
        if aPIFromService is CreatedAPI|BadRequestError|InternalServerErrorError {
            return aPIFromService;
        } else {
            log:printError("Internal Error occured deploying API", aPIFromService);
            InternalServerErrorError internalEror = {body: {code: 90900, message: "Internal Error occured deploying API"}};
            return internalEror;

        }

    }
    isolated resource function post apis/'import\-definition(@http:Payload json payload) returns CreatedAPI|BadRequestError|PreconditionFailedError|InternalServerErrorError {
        BadRequestError badRequest = {body: {code: 900910, message: "Not implemented"}};
        return badRequest;
    }
    isolated resource function post apis/'validate\-definition(http:Request message, boolean returnContent = false) returns APIDefinitionValidationResponse|BadRequestError|NotFoundError|InternalServerErrorError {
        APIClient apiService = new ();
        do {
            APIDefinitionValidationResponse|BadRequestError|NotFoundError|InternalServerErrorError|error validateDefinition = apiService.validateDefinition(message, returnContent);
            if validateDefinition is APIDefinitionValidationResponse|BadRequestError|NotFoundError|InternalServerErrorError {
                return validateDefinition;
            } else {
                InternalServerErrorError internalError = {body: {code: 90900, message: ""}};
                return internalError;
            }
        }
    }
    isolated resource function post apis/validate() returns http:Ok|BadRequestError|PreconditionFailedError|InternalServerErrorError {
        BadRequestError badRequest = {body: {code: 900910, message: "Not implemented"}};
        return badRequest;
    }
    isolated resource function get apis/[string apiId]/definition() returns json|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        APIClient apiService = new ();
        return apiService.getAPIDefinitionByID(apiId);
    }
    isolated resource function put apis/[string apiId]/definition(@http:Payload json payload) returns string|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|InternalServerErrorError {
        BadRequestError badRequest = {body: {code: 900910, message: "Not implemented"}};
        return badRequest;
    }
    isolated resource function get apis/export(string? apiId, string? name, string? 'version, string? format) returns json|NotFoundError|InternalServerErrorError {
        InternalServerErrorError internalError = {body: {code: 900910, message: "Not implemented"}};
        return internalError;
    }
    isolated resource function post apis/'import(boolean? overwrite, @http:Payload json payload) returns http:Ok|ForbiddenError|ConflictError|PreconditionFailedError|InternalServerErrorError {
        InternalServerErrorError internalError = {body: {code: 900910, message: "Not implemented"}};
        return internalError;
    }
    isolated resource function get services(string? name, string? namespace, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns ServiceList|BadRequestError|InternalServerErrorError {
        ServiceClient serviceClient = new ();
        return serviceClient.getServices(name, namespace, sortBy, sortOrder, 'limit, offset);
    }
    isolated resource function get services/[string serviceId]() returns Service|BadRequestError|NotFoundError|InternalServerErrorError {
        ServiceClient serviceClient = new ();
        return serviceClient.getServiceById(serviceId);
    }
    isolated resource function get services/[string serviceId]/usage() returns APIList|BadRequestError|NotFoundError|InternalServerErrorError {
        ServiceClient serviceClient = new ();
        return serviceClient.getServiceUsageByServiceId(serviceId);
    }
    isolated resource function get policies(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns MediationPolicyList|InternalServerErrorError {
        InternalServerErrorError internalError = {body: {code: 900910, message: "Not implemented"}};
        return internalError;
    }
    isolated resource function get policies/[string policyId]() returns MediationPolicy|NotFoundError|InternalServerErrorError {
        InternalServerErrorError internalError = {body: {code: 900910, message: "Not implemented"}};
        return internalError;
    }
};
public isolated function handleAPKError(APKError errorDetail) returns InternalServerErrorError|BadRequestError{
                ErrorHandler & readonly detail = errorDetail.detail();
            if detail.statusCode=="400"{
                BadRequestError badRequest = {body: {code: detail.code, message: detail.message}};
                return badRequest;
            }
                InternalServerErrorError internalServerError = {body: {code: detail.code, message: detail.message}};
                return internalServerError;

}