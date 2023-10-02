import ballerina/http;



public type CreatedAPI record {|
    *http:Created;
    API body;
|};


public type APIBody record {
    API apiProperties;
    # Content of the definition
    record {} Definition;
};

public type WSDLInfo record {
    # Indicates whether the WSDL is a single WSDL or an archive in ZIP format
    string 'type?;
};

public type API_threatProtectionPolicies_list record {
    string policyId?;
    int priority?;
};


public type API_serviceInfo record {
    string 'key?;
    string name?;
    string 'version?;
    boolean outdated?;
};

public type APIDefinition1 record {
    # Content of the definition
    record {} Definition;
};


public type API_additionalPropertiesMap record {
    string name?;
    string value?;
    boolean display?;
};

public type API_threatProtectionPolicies record {
    API_threatProtectionPolicies_list[] list?;
};
