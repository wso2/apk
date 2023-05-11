public type KeyManagerListingDaoEntry record {|
    string uuid;
    string name;
    string display_name?;
    string description?;
    string 'type;
    boolean enabled;
|};
public type KeyManagerDaoEntry record {|
string uuid?;
string name;
string display_name?;
string issuer;
string description?;
string 'type;
byte[] configuration?;
boolean enabled;
|};