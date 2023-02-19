public type RuntimeAPI record {|
    string apiVersion = "dp.wso2.com/v1alpha1";
    string kind = "RuntimeAPI";
    Metadata metadata;
    RuntimeAPISpec spec;
|};

public type ServiceInfo record {
    string name;
    string namespace;
};

public type RuntimeAPISpec record {|
    string name;
    string context;
    string 'type;
    string 'version;
    anydata endpointConfig?;
    Operations[] operations?;
    OperationPolicies apiPolicies?;
    ServiceInfo serviceInfo?;
|};

public type Operations record {|
    boolean authTypeEnabled = true;
    anydata endpointConfig?;
    string[] scopes = [];
    string target;
    string verb;
    OperationPolicies operationPolicies?;
|};

public type OperationPolicy record {
    string policyName;
    map<string> parameters?;
};

public type OperationPolicies record {|
    OperationPolicy[] request?;
    OperationPolicy[] response?;
|};
