import ballerina/io;
import admin_service.org.wso2.apk.apimgt.api as api;
import admin_service.org.wso2.apk.apimgt.init as apkinit;

configurable DatasourceConfiguration datasourceConfiguration = ?;
configurable ThrottlingConfiguration throttleConfig = ?;

function init() {
    io:println("Starting APK Admin Domain Service...");
    APKConfiguration apkConfig = {
        throttlingConfiguration: throttleConfig,
        datasourceConfiguration: datasourceConfiguration
    };
    string configJson = apkConfig.toJson().toJsonString();
    // Pass the configurations to java init component
    api:APIManagementException? err = apkinit:APKComponent_activate(configJson);
    if (err != ()) {
        io:println(err);
    }
}
