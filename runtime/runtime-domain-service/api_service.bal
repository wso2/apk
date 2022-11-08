import ballerina/http;

listener http:Listener ep0 = new (9090, config = {host: "localhost"});

service / on ep0 {
    // resource function get apis(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns APIList|NotAcceptableError {
    // }
    // resource function post apis(@http:Payload API payload) returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError {
    // }
    // resource function get apis/[string apiId]() returns API|NotFoundError|NotAcceptableError {
    // }
    // resource function put apis/[string apiId](@http:Payload API payload) returns API|BadRequestError|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
    // }
    // resource function delete apis/[string apiId]() returns http:Ok|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
    // }
    // resource function post apis/'import\-service(string serviceKey, @http:Payload API payload) returns CreatedAPI|NotFoundError|InternalServerErrorError {
    // }
    // resource function post apis/'import\-definition(@http:Payload json payload) returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError {
    // }
    // resource function post apis/'validate\-definition(@http:Payload json payload, boolean returnContent = false) returns APIDefinitionValidationResponse|BadRequestError|NotFoundError {
    // }
    // resource function post apis/validate() returns http:Ok|BadRequestError|NotFoundError {
    // }
    // resource function get apis/[string apiId]/definition() returns string|NotFoundError|NotAcceptableError {
    // }
    // resource function put apis/[string apiId]/definition(@http:Payload json payload) returns string|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError {
    // }
    // resource function get apis/export(string? apiId, string? name, string? 'version, string? format) returns json|NotFoundError|InternalServerErrorError {
    // }
    // resource function post apis/'import(boolean? overwrite, @http:Payload json payload) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|InternalServerErrorError {
    // }
    // resource function get services(string? name, string? namespace, string sortBy = "createdTime", string sortOrder = "desc", int 'limit = 25, int offset = 0) returns ServiceList|BadRequestError|UnauthorizedError|InternalServerErrorError {
    // }
    // resource function get services/[string serviceId](string? namespace) returns Service|BadRequestError|UnauthorizedError|NotFoundError|InternalServerErrorError {
    // }
    // resource function get services/[string serviceId]/usage(string? namespace) returns APIList|BadRequestError|UnauthorizedError|NotFoundError|InternalServerErrorError {
    // }
    // resource function get policies(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc", @http:Header string? accept = "application/json") returns MediationPolicyList|NotAcceptableError {
    // }
    // resource function get policies/[string mediationPolicyId]() returns MediationPolicy|NotFoundError|NotAcceptableError {
    // }
}
