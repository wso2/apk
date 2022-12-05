import ballerina/http;

configurable int BACKOFFICE_PORT = 9443;

listener http:Listener ep0 = new (BACKOFFICE_PORT);

service /api/am/backoffice/internal on ep0 {
    resource function post apis(@http:Payload API payload, string openAPIVersion = "v3") returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError|error {
        API | error ? createdApi = createAPI(payload);
        if createdApi is API {
            CreatedAPI crAPI = {body: check createdApi.cloneWithType(API)};
            return crAPI;
        }
        return error("Error while adding API");
    }
    resource function put apis/[string apiId](@http:Header string? 'if\-match, @http:Payload API payload) returns API|BadRequestError|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError|error {
        API | error ? updatedAPI = updateAPI(apiId, payload);
        if updatedAPI is API {
            API upAPI = check updatedAPI.cloneWithType(API);
            return upAPI;
        }
        return error("Error while updating API");
    }
    resource function delete apis/[string apiId](@http:Header string? 'if\-match) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError|http:InternalServerError {
        string|error? response = deleteAPI(apiId);
        if response is error {
            http:InternalServerError internalError = {body: {code: 90912, message: "Internal Error while deleting API By Id"}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    // resource function put apis/[string apiId]/definition(@http:Header string? 'if\-match, @http:Payload json payload) returns string|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post apis/'validate\-openapi(@http:Payload json payload, boolean returnContent = false) returns OpenAPIDefinitionValidationResponse|BadRequestError|NotFoundError {
    // }
    // resource function post apis/'validate\-wsdl(@http:Payload json payload) returns WSDLValidationResponse|BadRequestError|NotFoundError {
    // }
    // resource function post apis/'validate\-graphql\-schema(@http:Payload json payload) returns GraphQLValidationResponse|BadRequestError|NotFoundError {
    // }
}
