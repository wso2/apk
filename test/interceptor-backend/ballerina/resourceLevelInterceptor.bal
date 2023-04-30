import ballerina/http;

http:Service interceptorService = service object {
    # Handle Request
    #
    # + payload - Content of the request 
    # + return - Successful operation 
    isolated resource function post 'handle\-request(@http:Payload RequestHandlerRequestBody payload) returns OkRequestHandlerResponseBody {
        OkRequestHandlerResponseBody okRequestHandlerResponseBody = {body: {headersToAdd: {"Interceptor-header": "Interceptor-header-value"}}};
        return okRequestHandlerResponseBody;
    }
    # Handle Response
    #
    # + payload - Content of the request 
    # + return - Successful operation 
    isolated resource function post 'handle\-response(@http:Payload ResponseHandlerRequestBody payload) returns OkResponseHandlerResponseBody {
        return {body: {headersToAdd: {"Interceptor-Response-header": "Interceptor-Response-header-value"}}};
    }
    isolated resource function get health() returns http:Ok {
        json status = {"health": "Ok"};
        return {body: status};
    }
};
