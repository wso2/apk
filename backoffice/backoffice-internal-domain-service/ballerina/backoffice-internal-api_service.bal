import ballerina/http;

configurable int BACKOFFICE_PORT = 9443;

listener http:Listener ep0 = new (BACKOFFICE_PORT);

service /api/am/backoffice/internal on ep0 {
    isolated resource function post apis(@http:Payload json payload) returns CreatedAPI|BadRequestError|UnsupportedMediaTypeError|error {
        APIBody apiBody = check payload.cloneWithType(APIBody);
         
        API | error ? createdApi = createAPI(apiBody, "carbon.super");
        if createdApi is API {
            CreatedAPI crAPI = {body: check createdApi.cloneWithType(API)};
            return crAPI;
        }
        return error("Error while adding API", createdApi);
    }


    isolated resource function put apis/[string apiId](@http:Header string? 'if\-match, @http:Payload json payload) returns API|BadRequestError|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError|error {
        APIBody apiUpdateBody = check payload.cloneWithType(APIBody);
        
        API | error ? updatedAPI = updateAPI(apiId, apiUpdateBody, "carbon.super");
        if updatedAPI is API {
            API upAPI = check updatedAPI.cloneWithType(API);
            return upAPI;
        }
        return error("Error while updating API");
    }

    isolated resource function delete apis/[string apiId](@http:Header string? 'if\-match) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|PreconditionFailedError|http:InternalServerError {
        string|error? response = deleteAPI(apiId);
        if response is error {
            http:InternalServerError internalError = {body: {code: 90912, message: "Internal Error while deleting API By Id"}};
            return internalError;
        } else {
            return http:OK;
        }
    }
    isolated resource function put apis/[string apiId]/definition(@http:Header string? 'if\-match, @http:Payload APIDefinition payload) returns string|BadRequestError|ForbiddenError|NotFoundError|PreconditionFailedError|error { 
        APIDefinition | error ? updateDef = updateDefinition(payload, apiId);
        if updateDef is APIDefinition {
            APIDefinition crAPI = check updateDef.cloneWithType(APIDefinition);
            return crAPI.Definition.toString();
        }
        return error("Error while updating API definition");
    }
}
