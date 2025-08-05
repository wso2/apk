import ballerina/http;

# OAuth API
public isolated client class LoginClientModule {
    final http:Client clientEp;

    public isolated function init(http:Client clientEp) {
        self.clientEp = clientEp;
    }

    #
    # + return - Response from authorization endpoint 
    resource isolated function post login(Login_body payload) returns http:Response|error {
        string resourcePath = string `/commonauth/login`;
        http:Request request = new;
        string encodedRequestBody = createFormURLEncodedRequestBody(payload);
        request.setPayload(encodedRequestBody, "application/x-www-form-urlencoded");
        http:Response response = check self.clientEp->post(resourcePath, request);
        return response;
    }
}
