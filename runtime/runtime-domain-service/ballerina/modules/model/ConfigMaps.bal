public type ConfigMap record {|
    string kind = "ConfigMap";
    string apiVersion = "v1";
    Metadata metadata;
    map<string> data?;
    map<string> binaryData?;
|};
