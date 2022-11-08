import ballerina/io;
import apk.org.wso2.apk.apimgt.init as apkinit;
import apk.org.wso2.apk.apimgt.api as api;

configurable DatasourceConfiguration datasourceConfiguration = ?;
configurable ThrottlingConfiguration throttleConfig = ?;

function init() {
    io:println("Starting APK Admin Domain Service...");
    APKConfiguration apkConfig = {
        throttlingConfiguration: throttleConfig,
        datasourceConfiguration: datasourceConfiguration
    };
    string configJson = apkConfig.toJson().toJsonString();
    io:println(configJson);
    api:APIManagementException? err = apkinit:APKComponent_activate(configJson);
    if (err != ()) {
        io:println(err);
    }
}
