import ballerina/log;
import ballerina/http;

public service class RequestErrorInterceptor {
    *http:RequestErrorInterceptor;

    resource function 'default [string... path](error err, http:RequestContext ctx) returns http:Response {
        http:Response response = new;
        if err is APKError {
            APKError apkError = <APKError>err;
            ErrorHandler & readonly detail = apkError.detail();
            ErrorDto errorDto = {code: detail.code, message: detail.message, description: detail.description, moreInfo: detail.moreInfo};
            response.statusCode = detail.statusCode;
            response.setJsonPayload(errorDto);
        } else {
            log:printError("Exception Occured", err);
            response.statusCode = 500;
            ErrorDto errorDto = {code: 900900, message: "Internal Server Error", description: "Internal Server Error"};
            response.setJsonPayload(errorDto);
        }
        return response;
    }
}

public type ErrorDto record {|
    int code;
    string message;
    string description;
    map<string> moreInfo = {};
|};
