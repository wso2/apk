import ballerina/http;

http:Service gwInterceptorService = service object {
    # Handle Request
    #
    # + payload - Content of the request 
    # + return - Successful operation 
    isolated resource function post 'handle\-request(@http:Payload RequestHandlerRequestBody payload) returns OkRequestHandlerResponseBody {
        Headers? requestHeaders = payload.requestHeaders;
        string? userId = ();
        string? orgId = ();
        if requestHeaders is Headers {
            if requestHeaders.hasKey("user_id") {
                userId = requestHeaders.get("user_id");
            }
            if requestHeaders.hasKey("org_id") {
                orgId = requestHeaders.get("org_id");
            }
        }
        RateLimitKeys ratelimitKeys = {};
        if userId is string {
            ratelimitKeys["user_key"] = userId;
        }
        if orgId is string {
            ratelimitKeys["org_key"] = orgId;
        }
        OkRequestHandlerResponseBody okRequestHandlerResponseBody = {body: {headersToAdd: {"GW-header": "GW-header-value"}, rateLimitKeys: ratelimitKeys}};
        return okRequestHandlerResponseBody;
    }
    # Handle Response
    #
    # + payload - Content of the request 
    # + return - Successful operation 
    isolated resource function post 'handle\-response(@http:Payload ResponseHandlerRequestBody payload) returns OkResponseHandlerResponseBody {
        return {body: {headersToAdd: {"GW-response-header": "GW-response-header-value"}}};
    }
    isolated resource function get health() returns http:Ok {
        json status = {"health": "Ok"};
        return {body: status};
    }
};
