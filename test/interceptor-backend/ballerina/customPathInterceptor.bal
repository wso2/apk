import ballerina/http;

http:Service customPathInterceptorService = service object {
    # Handle Request at custom path
    #
    # + payload - Content of the request
    # + return - Successful operation
    isolated resource function post 'handle\-custom\-request(@http:Payload RequestHandlerRequestBody payload) returns OkRequestHandlerResponseBody {
        OkRequestHandlerResponseBody okRequestHandlerResponseBody = {body: {headersToAdd: {"Custom-Interceptor-Header": "Custom-path-interceptor-value"}}};
        return okRequestHandlerResponseBody;
    }
    # Handle Response at custom path
    #
    # + payload - Content of the request
    # + return - Successful operation
    isolated resource function post 'handle\-custom\-response(@http:Payload ResponseHandlerRequestBody payload) returns OkResponseHandlerResponseBody {
        return {body: {headersToAdd: {"Custom-Interceptor-Response-Header": "Custom-path-response-interceptor-value"}}};
    }

    isolated resource function get health() returns http:Ok {
        json status = {"health": "Ok"};
        return {body: status};
    }
};
