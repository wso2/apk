import ballerina/http;

public type BadRequestTokenErrorResponse record {|
    *http:BadRequest;
    TokenErrorResponse body;
|};

public type OkTokenResponse record {|
    *http:Ok;
    TokenResponse body;
|};

public type UnauthorizedTokenErrorResponse record {|
    *http:Unauthorized;
    TokenErrorResponse body;
|};

public type TokenErrorResponse record {
    # Error code classifying the type of preProcessingError.
    string 'error;
    string error_description?;
};

public type TokenResponse record {
    # OAuth access tokn issues by authorization server.
    string access_token;
    # The type of the token issued.
    string token_type;
    # The lifetime in seconds of the access token.
    decimal expires_in?;
    # OPTIONAL.The refresh token, which can be used to obtain new access tokens.
    string refresh_token?;
    # The scope of the access token requested.
    string scope?;
};

public type JWKList_keys record {
    string kid?;
    string kty?;
    string use?;
    string[] key_ops?;
    string alg?;
    string x5u?;
    string[] x5c?;
    string x5t?;
    string 'x5t\#S256?;
    string e?;
    string n?;
    string x?;
    string y?;
    string d?;
    string p?;
    string q?;
    string dp?;
    string dq?;
    string qi?;
    string k?;
};

public type JWKList record {
    JWKList_keys keys?;
};

public type Token_body record {
    # Required OAuth grant type
    string grant_type;
    # Authorization code to be sent for authorization grant type
    string code?;
    # Clients redirection endpoint
    string redirect_uri?;
    # OAuth client identifier
    string client_id?;
    # OAuth client secret
    string client_secret?;
    # Refresh token issued to the client.
    string refresh_token?;
    # OAuth scopes
    string scope?;
    # username
    string username?;
    # password
    string password?;
    # Validity period of token
    int validity_period?;
};
