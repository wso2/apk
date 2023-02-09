import ballerina/http;

public type NotFoundClientRegistrationError record {|
    *http:NotFound;
    ClientRegistrationError body;
|};

public type InternalServerErrorClientRegistrationError record {|
    *http:InternalServerError;
    ClientRegistrationError body;
|};

public type ConflictClientRegistrationError record {|
    *http:Conflict;
    ClientRegistrationError body;
|};

public type BadRequestClientRegistrationError record {|
    *http:BadRequest;
    ClientRegistrationError body;
|};

public type CreatedApplication record {|
    *http:Created;
    Application body;
|};

public type ClientRegistrationError record {
    string 'error?;
    string error_description?;
};

public type UpdateRequest record {
    string[] redirect_uris?;
    string client_name?;
    string[] grant_types?;
};



public type RegistrationRequest record {
    string[] redirect_uris?;
    string client_name?;
    string[] grant_types?;
};

public type Application record {
    string client_id?;
    string client_secret?;
    string[] redirect_uris?;
    string[] grant_types?;
    string client_name?;
    int client_secret_expires_at?;
};
