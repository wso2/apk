public type IDPConfiguration record {|
string issuer = "wso2.org/products/am";
string jwksUrl? ;
string organizationClaim ="organization";
string authorizationHeader = "X-JWT-Assertion";
KeyStore publicKey;
|};

public type KeyStore record {|
    string path;
    string keyPassword?;
|};
