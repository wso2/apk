import ballerina/http;

public type UnauthorizedLoginErrorResponse record {|
    *http:Unauthorized;
    LoginErrorResponse body;
|};

public type Login_body record {
    # username
    string username;
    # password
    string password;
    # organization
    string organization?;
    string sessionKey;
};

public type LoginErrorResponse record {
    # Error code classifying the type of preProcessingError.
    string 'error;
    string error_description?;
};
