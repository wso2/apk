import ballerina/lang.value;
import ballerina/http;
import ballerina/log;
import ballerina/regex;

http:Service interceptorService = service object {
    # Handle Request
    #
    # + payload - Content of the request 
    # + return - Successful operation 
    isolated resource function post 'handle\-request(@http:Payload RequestHandlerRequestBody payload) returns OkRequestHandlerResponseBody {
        map<string> headers = {"Interceptor-header": "Interceptor-header-value"};
        InvocationContext? invocationContext = payload.invocationContext;
        if invocationContext is InvocationContext {
            string? apiProperties = invocationContext.apiProperties;
            if apiProperties is string {
                string replacedAPIProperties = regex:replaceAll(apiProperties, "'", "\"");
                do {
                    APIProperties apiPropertyJson = check value:fromJsonStringWithType(replacedAPIProperties, APIProperties);
                    foreach any key in apiPropertyJson.keys() {
                        if key is string {
                            headers["Interceptor-header-"+key] = <string>apiPropertyJson.get(key);
                        }
                    }
                } on fail var e {
                    log:printError("Error while parsing apiProperties: " + e.message());
                }
            }
        }
        OkRequestHandlerResponseBody okRequestHandlerResponseBody = {body: {headersToAdd: headers}};
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

public type APIProperties record {

};
