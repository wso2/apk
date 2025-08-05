import ballerina/log;
import ballerina/http;

isolated service /oauth2 on ep0 {
    # Description
    #
    # + response_type - Expected response type 
    # + client_id - OAuth client identifier 
    # + redirect_uri - Clients redirection endpoint 
    # + scope - OAuth scopes 
    # + state - Opaque value used by the client to maintain state between the request and callback 
    # + return - Response from authorization endpoint 
    isolated resource function get authorize(string response_type, string client_id, string? redirect_uri, string? scope, string? state) returns http:Found {
        TokenUtil tokenUtil = new;
        return tokenUtil.handleAuthorizeRequest(response_type,client_id,redirect_uri,scope,state);
    }
    # Description
    #
    # + sessionKey - Session key. 
    # + return - Response from authorization endpoint 
    isolated resource function get 'auth\-callback(string sessionKey,http:Request request) returns http:Found {
        TokenUtil tokenUtil = new;
        return tokenUtil.handleOauthCallBackRequest(request,sessionKey);
    }
    # Description
    #
    # + authorization - Authentication scheme header 
    # + payload - parameter description 
    # + return - returns can be any of following types
    # OkTokenResponse (OK.Successful response from token endpoint.)
    # BadRequestTokenErrorResponse (Bad Request,Error response from token endpoint due to malformed request.)
    # UnauthorizedTokenErrorResponse (Unauthorized. Error response from token endpoint due to client authentication failure.)
    isolated resource function post token(@http:Header string? authorization, @http:Payload map<string> payload) returns OkTokenResponse|BadRequestTokenErrorResponse|UnauthorizedTokenErrorResponse {
        TokenUtil tokenUtil = new;
        do {
            Token_body tokenBody = check payload.cloneWithType(Token_body);
            return check tokenUtil.generateToken(authorization, tokenBody);

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
}
