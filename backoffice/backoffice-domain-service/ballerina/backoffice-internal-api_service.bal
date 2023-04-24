import ballerina/http;
import wso2/apk_common_lib as commons;

configurable int BACKOFFICE_PORT_INT = 9444;
listener http:Listener ep1 = new (BACKOFFICE_PORT_INT, secureSocket = {
    'key: {
        certFile: <string>keyStores.tls.certFilePath,
        keyFile: <string>keyStores.tls.keyFilePath
    }
}, interceptors = [requestErrorInterceptor, responseErrorInterceptor]);

service /api/am/backoffice/internal on ep1 {
    isolated resource function post apis(@http:Payload json payload) returns CreatedAPI|error {
        APIBody apiBody = check payload.cloneWithType(APIBody);

        API|error? createdApi = createAPI(apiBody, "carbon.super");
        if createdApi is API {
            CreatedAPI crAPI = {body: check createdApi.cloneWithType(API)};
            return crAPI;
        }
        return error("Error while adding API", createdApi);
    }

    isolated resource function get apis/[string apiId](@http:Header string? 'if\-none\-match) returns API|commons:APKError|error {
        API | commons:APKError | error ? response = getAPI_internal(apiId, "carbon.super");
        if (response is API | commons:APKError) {
            return response;
        }
        return error("Error while retireving API");
    }

    isolated resource function put apis/[string apiId](@http:Header string? 'if\-match, @http:Payload json payload) returns API|commons:APKError|error {
        APIBody apiUpdateBody = check payload.cloneWithType(APIBody);

        API|commons:APKError |error? updatedAPI = updateAPI_internal(apiId, apiUpdateBody, "carbon.super");
        if updatedAPI is API {
            API upAPI = check updatedAPI.cloneWithType(API);
            return upAPI;
        } else if (updatedAPI is commons:APKError) {
            return updatedAPI;
        }
        return error("Error while updating API");
    }

    isolated resource function delete apis/[string apiId](@http:Header string? 'if\-match) returns http:Ok|commons:APKError {
        string|commons:APKError|error? response = deleteAPI(apiId, "carbon.super");
        if response is commons:APKError {
            return response;
        }
        else if response is error {
            return e909605();
        } else {
            return http:OK;
        }
    }
    isolated resource function put apis/[string apiId]/definition(@http:Header string? 'if\-match, @http:Payload APIDefinition1 payload) returns string|error {
        APIDefinition1|error? updateDef = updateDefinition(payload, apiId);
        if updateDef is APIDefinition1 {
            APIDefinition1 crAPI = check updateDef.cloneWithType(APIDefinition1);
            return crAPI.Definition.toString();
        }
        return error("Error while updating API definition");
    }
}
