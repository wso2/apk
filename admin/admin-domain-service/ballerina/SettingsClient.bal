import apk_keymanager_libs;
import wso2/apk_common_lib as commons;

public class SettingsClient {

    public isolated function getSettings(commons:Organization? organization) returns Settings {
        Settings settings = {};
        settings.keyManagerConfiguration = self.setKeyManagerConfigsToSettings();
        return settings;
    }
    private isolated function setKeyManagerConfigsToSettings() returns Settings_keyManagerConfiguration[] {
        Settings_keyManagerConfiguration[] keyManagerConfigs = [];
        apk_keymanager_libs:KeyManagerConfigurations[] kmconfigs = keyManagerInitializer.retrieveAllKeyManagerConfigs();
        foreach apk_keymanager_libs:KeyManagerConfigurations item in kmconfigs {
            Settings_keyManagerConfiguration keyManagerConfig = {};
            keyManagerConfig.'type = item.'type;
            keyManagerConfig.displayName = item.'display_name;
            keyManagerConfig.defaultConsumerKeyClaim = item.consumerKeyClaim;
            keyManagerConfig.defaultScopesClaim = item.scopesClaim;
            KeyManagerConfiguration[] endpointConfiguration = [];
            apk_keymanager_libs:EndpointConfiguration[] endpointConfigs = item.endpoints;
            foreach apk_keymanager_libs:EndpointConfiguration endpointConfig in endpointConfigs {
                KeyManagerConfiguration endpointConfigGenerated = {
                    name: endpointConfig.name,
                    tooltip: endpointConfig.toolTip,
                    label: endpointConfig.display_name,
                    'type: "input",
                    mask: false,
                    required: endpointConfig.required,
                    multiple: false
                };
                endpointConfiguration.push(endpointConfigGenerated);
            }
            keyManagerConfig.endpointConfigurations = endpointConfiguration;
            apk_keymanager_libs:KeyManagerConfiguration[] keyManagerConnectorConfigs = item.endpointConfigurations;
            KeyManagerConfiguration[] generatedConnectorConfigs = [];
            foreach apk_keymanager_libs:KeyManagerConfiguration keyManagerConnectorConfig in keyManagerConnectorConfigs {
                KeyManagerConfiguration generatedConnectorConfig = {
                    name: keyManagerConnectorConfig.name,
                    tooltip: keyManagerConnectorConfig.toolTip,
                    label: keyManagerConnectorConfig.display_name,
                    'type: keyManagerConnectorConfig.'type,
                    mask: keyManagerConnectorConfig.masked,
                    required: keyManagerConnectorConfig.required,
                    multiple: keyManagerConnectorConfig.multiple,
                    default: keyManagerConnectorConfig?.default,
                    values: keyManagerConnectorConfig.values
                };

                generatedConnectorConfigs.push(generatedConnectorConfig);
            }
            keyManagerConfigs.push(keyManagerConfig);
        }
        return keyManagerConfigs;
    }
}
