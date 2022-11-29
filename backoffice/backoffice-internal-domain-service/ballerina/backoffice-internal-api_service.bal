import backoffice_internal_service.org.wso2.apk.apimgt.api as api;
import ballerina/http;
import ballerina/lang.value;

configurable int BACKOFFICE_PORT = 9443;

listener http:Listener ep0 = new (BACKOFFICE_PORT);

service /api/am/backoffice/internal on ep0 {
    resource function post apis(@http:Payload API payload, string openAPIVersion = "v3") returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError|error {
        string|api:APIManagementException? createdApi = createAPI(payload);
        if createdApi is string {
            json j = check value:fromJsonString(createdApi);
            CreatedAPI crAPI = {body: check j.cloneWithType(API)};
            return crAPI;
        }
        return error("Error while adding API");
    }
    // resource function put apis/[string apiId](@http:Header string? 'if\-match, @http:Payload API payload) returns API|BadRequestError|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
    // }
    // resource function delete apis/[string apiId](@http:Header string? 'if\-match) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError {
    // }
    // resource function put apis/[string apiId]/definition(@http:Header string? 'if\-match, @http:Payload json payload) returns string|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError {
    // }
    // resource function post apis/'validate\-openapi(@http:Payload json payload, boolean returnContent = false) returns OpenAPIDefinitionValidationResponse|BadRequestError|NotFoundError {
    // }
    // resource function post apis/'validate\-wsdl(@http:Payload json payload) returns WSDLValidationResponse|BadRequestError|NotFoundError {
    // }
    // resource function post apis/'validate\-graphql\-schema(@http:Payload json payload) returns GraphQLValidationResponse|BadRequestError|NotFoundError {
    // }
}
