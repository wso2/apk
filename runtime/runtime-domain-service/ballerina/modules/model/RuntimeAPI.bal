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
    record {} endpointConfig?;
    Operations[] operations?;
    OperationPolicies apiPolicies?;
    ServiceInfo serviceInfo?;
|};

public type Operations record {|
    boolean authTypeEnabled = true;
    record {} endpointConfig?;
    string[] scopes = [];
    string target;
    string verb;
    OperationPolicies operationPolicies?;
|};

public type OperationPolicy record {
    string policyName;
    OperationPolicyParameters[] parameters?;
};

public type OperationPolicyParameters record {
    string headerName?;
    string headerValue?;
};

public type MediationPolicy record {
    string id;
    string 'type;
    string name;
    string displayName?;
    string description?;
    string[] applicableFlows?;
    string[] supportedApiTypes?;
    MediationPolicySpecAttribute[] policyAttributes?;
};

public type MediationPolicySpecAttribute record {
    string name?;
    string description?;
    boolean required?;
    string validationRegex?;
    string 'type?;
    string defaultValue?;
};

public type OperationPolicies record {|
    OperationPolicy[] request = [];
    OperationPolicy[] response=[];
|};
