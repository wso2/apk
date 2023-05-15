import ballerina/http;

public type OkInline_response_200 record {|
    *http:Ok;
    Inline_response_200 body;
|};

public type InternalServerErrorError record {|
    *http:InternalServerError;
    Error body;
|};

public type ErrorListItem record {
    # Error code
    string code;
    # Description about individual errors occurred
    string message;
};

public type ClientRegistrationResponse record {
    *ClientRegistrationRequest;
    string client_secret?;
    string client_id?;
    int client_secret_expires_at?;
    string registration_access_token?;
};

public type ConfigurationInitialization record {
    record {} endpoints?;
    record {} configurations?;
};

public type Error record {
    # Error code
    int code;
    # Error message.
    string message;
    # A detail description about the error message.
    string description?;
    # Preferably an url with more details about the error.
    string moreInfo?;
    # If there are more than one error list them out.
    # For example, list out validation errors by each field.
    ErrorListItem[] 'error?;
};

public type ClientRegistrationRequest record {
    string[] redirect_uris?;
    string[] response_types?;
    string[] grant_types?;
    string application_type?;
    string client_name?;
    string logo_uri?;
    string client_uri?;
    string policy_uri?;
    string tos_uri?;
    string jwks_uri?;
    string subject_type?;
    string token_endpoint_auth_method?;
    record {} additional_properties?;
};

public type Inline_response_500 record {
    int code?;
    string message?;
    string description?;
};

public type ClientUpdateRequest record {
    *ClientRegistrationRequest;
    string client_secret?;
    string client_id?;
};

public type Inline_response_200 record {
    boolean initialized?;
};
