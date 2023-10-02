public type KeyManager record {
    string id?;
    string name;
    string displayName?;
    string 'type;
    string description?;
    map<string> endpoints = {};
    # PEM type certificate
    string tlsCertficate?;
    string issuer;
    string[] availableGrantTypes?;
    boolean enableTokenGeneration?;
    boolean enableMapOAuthConsumerApps = false;
    boolean enableOAuthAppCreation = true;
    boolean enableOauthAppValidation = true;
    string consumerKeyClaim?;
    string scopesClaim?;
    boolean enabled = true;
    record {} additionalProperties?;
};

public type KeyManagerEndpoint record {
    string name;
    string value;
};

public type KeyManager_signingCertificate record {
    string 'type?;
    string value?;
};

public type KeyMappingDaoEntry record {
    string uuid;
    string application_uuid;
    string consumer_key;
    string key_type;
    string create_mode;
    byte[] app_info;
    string key_manager_uuid;
};
