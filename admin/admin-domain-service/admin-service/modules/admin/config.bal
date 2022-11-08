type ThrottlingConfiguration record {
    BlockCondition blockCondition = {
        enabled: true
    };
    boolean enableUnlimitedTier = true;
    boolean enableHeaderConditions = false;
    boolean enableJWTClaimConditions = false;
    boolean enableQueryParamConditions = false;
    boolean enablePolicyDeployment = true;
};

type BlockCondition record {
    boolean enabled = true;
};

type DatasourceConfiguration record {
    string name = "jdbc/apkdb";
    string description;
    string url;
    string username;
    string password;
    int maxPoolSize = 50;
    int minIdleTime = 60000;
    int maxLifeTime = 60000;
    int validationTimeout;
    boolean setAutocommit = false;
    string testQuery;
};

type APKConfiguration record {
    ThrottlingConfiguration throttlingConfiguration;
    DatasourceConfiguration datasourceConfiguration;
};