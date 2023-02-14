import ballerina/test;
import ballerina/http;
import ballerina/regex;
import ballerina/lang.array;


@test:Config
function testClientCredentialsTokenGenerationForFileBaseApp() returns error? {
    test:when(testgetConnection).callOriginal();
ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    TokenClient tokenClient = check new ("https://localhost:9443",connectionConfig);
    Token_body tokenRequest = {grant_type: "client_credentials"};
    string concatString = "45f1c5c8-a92e-11ed-afa1-0242ac120002:4fbd62ec-a92e-11ed-afa1-0242ac120002";
    string bse64encoded = array:toBase64(concatString.toBytes());
    string header = "Basic " + bse64encoded;
    TokenResponse tokenResponse = check tokenClient->/token.post(tokenRequest, header);
    test:assertEquals(tokenResponse.token_type, "Bearer");
    test:assertTrue(tokenResponse.access_token.length() >= 0);
}

@test:Config
function testClientCredentialsTokenGenerationForFileBaseAppNegative1() returns error? {
    test:when(testgetConnection).callOriginal();
    TokenClient tokenClient = check new ("https://localhost:9443",connectionConfig);
    Token_body tokenRequest = {grant_type: "client_credentials"};
    string concatString = "45f1c5c8-a92e-11ed-afa1-0242ac120002:4fbd62ec-a92e-11ed-afa1-0242ac120003";
    string bse64encoded = array:toBase64(concatString.toBytes());
    string header = "Basic " + bse64encoded;
    TokenResponse|error tokenResponse = tokenClient->/token.post(tokenRequest, header);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("unauthorized_client"));
    }
    tokenResponse = tokenClient->/token.post(tokenRequest);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("access_denied"));
    }
    tokenResponse = tokenClient->/token.post(tokenRequest, bse64encoded);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("access_denied"));
    }
    tokenResponse = tokenClient->/token.post(tokenRequest, "Basic ");
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("access_denied"));
    }
    string bse64encodedError = array:toBase64("concatString".toBytes());
    tokenResponse = tokenClient->/token.post(tokenRequest, "Basic " + bse64encodedError);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("access_denied"));
    }
    concatString = "45f1c5c8-a92e-11ed-afa1-0242ac120002:4fbd62ec-a92e-11ed-afa1-0242ac120002";
    bse64encoded = array:toBase64(concatString.toBytes());
    tokenRequest = {grant_type: "password"};
    tokenResponse = tokenClient->/token.post(tokenRequest, "Basic " + bse64encoded);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("unsupported_grant_type"));
    }
    concatString = "45f1c5c8-a92e-11ed-afa1-0242ac120005:4fbd62ec-a92e-11ed-afa1-0242ac120005";
    bse64encoded = array:toBase64(concatString.toBytes());
    tokenRequest = {grant_type: "password"};
    tokenResponse = tokenClient->/token.post(tokenRequest, "Basic " + bse64encoded);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("unsupported_grant_type"));
    }
    concatString = "45f1c5c8-a92e-11ed-afa1-0242ac120004:4fbd62ec-a92e-11ed-afa1-0242ac120005";
    bse64encoded = array:toBase64(concatString.toBytes());
    tokenRequest = {grant_type: "password"};
    tokenResponse = tokenClient->/token.post(tokenRequest, "Basic " + bse64encoded);
    test:assertTrue(tokenResponse is error);
    if tokenResponse is error {
        test:assertTrue(tokenResponse.toString().includes("access_denied"));
    }
}

@test:Config {}
function testAuthorizationCodeGrant() returns error? {
    test:when(testgetConnection).callOriginal();
    DCRClient dcrClient = check new ("https://localhost:9443",connectionConfig);
    RegistrationRequest registrationRequest = {client_name: "authorizationApp", grant_types: ["authorization_code", "refresh_token"], redirect_uris: ["http://httpbin.org"]};
    Application createdApp = check dcrClient->/register.post(registrationRequest);
    ConnectionConfig connectionConfig = {cookieConfig: {enabled: true},secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    TokenClient tokenClient = check new ("https://localhost:9443", connectionConfig);
    LoginClientModule loginClient = new (tokenClient.clientEp);
    http:Response authCodeResponse = check tokenClient->/authorize("code", <string>createdApp.client_id, "http://httpbin.org", "default", ());
    test:assertEquals(authCodeResponse.statusCode, 302);
    string|http:HeaderNotFoundError header = authCodeResponse.getHeader("Location", http:LEADING);
    test:assertTrue(header is string);
    if header is string {
        test:assertTrue(header.startsWith("https://localhost:9443/login"));
        string queryParam = string:substring(header, "https://localhost:9443".length());
        string[] queryParamSplit = regex:split(queryParam, "=");
        string stateKey = queryParamSplit[1];
        Login_body loginPayLoad = {
            password: "admin",
            sessionKey: stateKey,
            username: "admin",
            organization: "org1"
        };
        http:Response|error loginResponse = loginClient->/login.post(loginPayLoad);
        if loginResponse is http:Response {
            test:assertEquals(loginResponse.statusCode, 302);
            string|http:HeaderNotFoundError locationHeaderFromLogin = loginResponse.getHeader("Location", http:LEADING);
            test:assertTrue(locationHeaderFromLogin is string);
            if locationHeaderFromLogin is string {
                test:assertTrue(locationHeaderFromLogin.startsWith("https://localhost:9443/login-callback"));
                string queryParamFromLogin = string:substring(locationHeaderFromLogin, "https://localhost:9443/login-callback".length());
                string[] queryParamSplit1 = regex:split(queryParamFromLogin, "=");
                string stateKeyFromLogin = queryParamSplit1[1];
                http:Response|error authCallBackResponse = tokenClient->/auth\-callback(stateKeyFromLogin);
                if authCallBackResponse is http:Response {
                    test:assertEquals(authCallBackResponse.statusCode, 302);
                    string|http:HeaderNotFoundError locationHeaderFromAuthCallBack = authCallBackResponse.getHeader("Location", http:LEADING);
                    test:assertTrue(locationHeaderFromAuthCallBack is string);
                    if locationHeaderFromAuthCallBack is string {
                        test:assertTrue(locationHeaderFromAuthCallBack.startsWith("http://httpbin.org"));
                        string queryParamFromAuthCallBack = string:substring(locationHeaderFromAuthCallBack, "http://httpbin.org".length());
                        string[] queryParamSplit2 = regex:split(queryParamFromAuthCallBack, "=");
                        string authcode = queryParamSplit2[1].trim();
                        string concatString = createdApp.client_id.toString() + ":" + createdApp.client_secret.toString();
                        string bse64encoded = array:toBase64(concatString.toBytes());
                        string authorizationHeader = "Basic " + bse64encoded;
                        Token_body tokenBody = {grant_type: "authorization_code", redirect_uri: "http://httpbin.org", code: authcode};
                        TokenResponse|error tokenResponse = tokenClient->/token.post(tokenBody, authorizationHeader);
                        if tokenResponse is TokenResponse {
                            test:assertTrue(tokenResponse.access_token.length() > 0);
                            test:assertTrue(tokenResponse.refresh_token is string);
                            test:assertTrue(tokenResponse.refresh_token.toString().length() > 0);
                            Token_body refreshTokenRequest  = {grant_type: "refresh_token",refresh_token: tokenResponse.refresh_token,scope: "default"};
                            TokenResponse|error refreshTokenResponse = tokenClient->/token.post(refreshTokenRequest,authorizationHeader);
                            if refreshTokenResponse is TokenResponse {
                            test:assertTrue(tokenResponse.access_token.length() > 0);
                            test:assertTrue(tokenResponse.refresh_token is string);
                            test:assertTrue(tokenResponse.refresh_token.toString().length() > 0);
                            }
                        }
                    }
                }
            }
        }
    }
}
