import ballerina/file;
import ballerina/io;
import ballerina/log;
import wso2/apk_common_lib as commons;

public isolated class KeyManagerTypeInitializer {
    private final map<KeyManagerConfigurations> keyManagerConfigurationsMap = {};
    public isolated function initialize(string path) returns error? {
        file:MetaData[] & readonly keyManagerConfigs = check file:readDir(path);
        foreach file:MetaData & readonly item in keyManagerConfigs {
            do {
                string content = check io:fileReadString(item.absPath);
                json? jsonStringContent = check commons:fromYamlStringToJson(content);
                if jsonStringContent is json {
                    KeyManagerConfigurations|error keyManagerConfigurations = jsonStringContent.cloneWithType(KeyManagerConfigurations);
                    if keyManagerConfigurations is KeyManagerConfigurations {
                        lock{
                        self.keyManagerConfigurationsMap[keyManagerConfigurations.'type] = keyManagerConfigurations.clone();
                        }
                    }
                }
            } on fail var e {
                log:printError("Error while reading the key manager configurations from the file: " + item.absPath + " " + e.message(), e);
            }
        }
    }
    public isolated function retrieveKeyManagerConfigByType(string 'type) returns KeyManagerConfigurations|() {
        lock{
        if self.keyManagerConfigurationsMap.hasKey('type) {
            return self.keyManagerConfigurationsMap.get('type).clone();
        }
        }
        return ();
    }
    public isolated function retrieveAllKeyManagerConfigs() returns KeyManagerConfigurations[] {
        lock{
        return self.keyManagerConfigurationsMap.toArray().clone();
        }
    }
}
