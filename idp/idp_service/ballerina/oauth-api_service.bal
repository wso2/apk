import ballerina/log;
import ballerina/http;

service /oauth2 on ep0 {
    resource function get authorize(string response_type, string client_id, string? redirect_uri, string? scope, string? state) returns http:Found {
        TokenUtil tokenUtil = new;
        return tokenUtil.handleAuthorizeRequest(response_type,client_id,redirect_uri,scope,state);
    }
    isolated resource function post token(@http:Header string? authorization, http:Request request) returns TokenResponse|BadRequestTokenErrorResponse|UnauthorizedTokenErrorResponse {
        TokenUtil tokenUtil = new;
        do {
            map<string> formParams = check request.getFormParams();
            Token_body payload = check formParams.cloneWithType(Token_body);
            return check tokenUtil.generateToken(authorization, payload);

        } on fail var e {
            log:printError("Error occured on pasing payload", e);
            BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            return tokenError;

        }

    }
    resource function get keys() returns JWKList {
        JWKList jwklist = {};
        return jwklist;
    }
    isolated resource function get 'auth\-callback(string sessionKey,http:Request request) returns http:Found {
        TokenUtil tokenUtil = new;
        return tokenUtil.handleOauthCallBackRequest(request,sessionKey);
    }

}
