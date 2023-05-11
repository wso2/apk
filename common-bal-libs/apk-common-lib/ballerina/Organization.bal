# Organization AdditonalProperties Definition.
#
# + key - property key.  
# + value - property value.
public type OrganizationProperties record {
    string key;
    string value;
};

public type Organization record {|
    string uuid;
    string name;
    string displayName;
    string organizationClaimValue;
    boolean enabled;
    string[] serviceListingNamespaces = ["*"];
    OrganizationProperties[] properties = [];
|};