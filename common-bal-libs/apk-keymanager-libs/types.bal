public type KeyManagerConfigurations record {|
    string 'type;
    string 'display_name?;
    EndpointConfiguration[] endpoints = [];
    KeyManagerConfiguration[] endpointConfigurations = [];
    KeyManagerConfiguration[] applicationConfigurations =[];
    string consumerKeyClaim;
    string scopesClaim;
|};

public type EndpointConfiguration record {|
    string name;
    string display_name;
    string toolTip;
    boolean required;
|};

public type KeyManagerConfiguration record {|
    string name;
    string 'display_name;
    string 'type;
    string toolTip;
    boolean required = false;
    string default?;
    boolean masked = false;
    string[] values = [];
    boolean multiple = false;
|};
