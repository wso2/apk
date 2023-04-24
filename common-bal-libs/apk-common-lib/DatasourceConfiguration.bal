import ballerina/os;
public type DatasourceConfiguration record {
    string name = "jdbc/apkdb";
    string description;
    string url;
    string host;
    int port;
    string databaseName;
    string username;
    string password = os:getEnv("DB_PASSWORD");
    int maxPoolSize = 50;
    int minIdle = 20;
    int maxLifeTime = 60000;
    int validationTimeout;
    boolean autoCommit = true;
    string testQuery;
    string driver;
};