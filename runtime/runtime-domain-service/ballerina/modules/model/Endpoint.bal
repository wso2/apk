public type Endpoint record {|
string url?;
string namespace?;
string name?;
boolean serviceEntry = false;

|};

public type EndpointSecurity record {|
    boolean enabled = false;
    string 'type?;
    string username?;
    string password?;
    string secretRefName?;
|};
