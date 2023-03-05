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
    map<string> parameters?;
};

public type MediationPolicy record {
    string 'type;
    string id?;
    string name;
    string displayName?;
    string description?;
    string[] applicableFlows?;
    string[] supportedApiTypes?;
    boolean isApplicableforAPILevel?;
    boolean isApplicableforOperationLevel?;
    MediationPolicySpecAttribute[] policyAttributes?;
};

public type MediationPolicySpecAttribute record {
    # Name of the attibute
    string name?;
    # Description of the attibute
    string description?;
    # Is this option mandetory for the policy
    boolean required?;
    # UI validation regex for the attibute
    string validationRegex?;
    # Type of the attibute
    string 'type?;
    # Default value for the attribute
    string defaultValue?;
};

public type OperationPolicies record {|
    OperationPolicy[] request = [];
    OperationPolicy[] response=[];
|};
