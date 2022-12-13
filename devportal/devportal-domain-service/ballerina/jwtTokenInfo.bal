public type JWTTokenInfo record {
    Application application;
    API[] subscribedAPIs;
    string subscriber;
    string expireTime;
    string keyType;
    string permittedIP;
    string permittedReferrer;
};