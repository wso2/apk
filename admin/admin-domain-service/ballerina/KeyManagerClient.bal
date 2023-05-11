import apk_keymanager_libs;
import ballerina/uuid;
import ballerina/lang.value;
import wso2/apk_common_lib as commons;

public class KeyManagerClient {

    public isolated function addKeyManagerEntryToOrganization(KeyManager keyManager, commons:Organization organization) returns KeyManager|commons:APKError {
        if keyManager.name.trim().length() == 0 {
            return e909434();
        }
        // validate existence
        if check checkKeyManagerExist(keyManager.name, organization) {
            return e909435(keyManager.name, organization.name);
        }
        // validate type
        apk_keymanager_libs:KeyManagerConfigurations? retrieveKeyManagerConfigByType = keyManagerInitializer.retrieveKeyManagerConfigByType(keyManager.'type);
        if retrieveKeyManagerConfigByType is () {
            return e909436(keyManager.'type);
        }
        // validate configurations.
        if (!self.validateKeyManagerConfigurations(keyManager, retrieveKeyManagerConfigByType)) {
            return e909437();
        }
        // add key manager entry.
        KeyManagerDaoEntry keyManagerDtoToInsert = self.fromKeyManagerModelToKeyManagerDaoEntry(keyManager);
        _ = check addKeyManagerEntry(keyManagerDtoToInsert, organization);
        return check self.getKeyManagerById(<string>keyManagerDtoToInsert.uuid, organization);
    }
    private isolated function fromKeyManagerModelToKeyManagerDaoEntry(KeyManager keyManager) returns KeyManagerDaoEntry {

        KeyManagerDaoEntry keyManagerDTO = {
            name: keyManager.name,
            'type: keyManager.'type,
            issuer: <string>keyManager.issuer,
            enabled: keyManager.enabled ?: true,
            description: keyManager.description
        };
        if keyManager.id is () {
            keyManagerDTO.uuid = uuid:createType1AsString();
        } else {
            keyManagerDTO.uuid = keyManager.id;
        }
        record {} additionalProperties = {};
        additionalProperties = keyManager.additionalProperties.clone() ?: {};
        KeyManagerEndpoint[]? endpoints = keyManager.endpoints;
        if endpoints is KeyManagerEndpoint[] {
            foreach KeyManagerEndpoint item in endpoints {
                record {} defineEndpoints = {};
                if (additionalProperties.hasKey("endpoints")) {
                    defineEndpoints = <record {|anydata...;|}>additionalProperties.get("endpoints");
                }else{
                    additionalProperties["endpoints"] = defineEndpoints;
                }
                defineEndpoints[item.name] = item.value;
            }
        }
        string[]? availableGrantTypes = keyManager.availableGrantTypes;
        if availableGrantTypes is string[] {
            foreach string grantType in availableGrantTypes {
                string[] grantTypes = [];
                if (additionalProperties.hasKey("grantTypes")) {
                    grantTypes = <string[]>additionalProperties.get("grantTypes");
                }
                grantTypes.push(grantType);
            }
        }
        if keyManager.consumerKeyClaim is string {
            additionalProperties["consumerKeyClaim"] = keyManager.consumerKeyClaim;
        }
        if keyManager.scopesClaim is string {
            additionalProperties["scopesClaim"] = keyManager.scopesClaim;
        }
        KeyManager_certificates? certificates = keyManager.certificates;
        if certificates is KeyManager_certificates {
            if certificates.'type is string {
                additionalProperties["certificate_type"] = certificates.'type;
            }
            if certificates.value is string {
                additionalProperties["certificate_value"] = certificates.value;
            }
        }
        additionalProperties["mapOAuthConsumerApps"] = keyManager.enableMapOAuthConsumerApps is boolean ? keyManager.enableMapOAuthConsumerApps : true;
        additionalProperties["enableTokenGeneration"] = keyManager.enableTokenGeneration is boolean ? keyManager.enableTokenGeneration : true;
        additionalProperties["enableOauthAppCreation"] = keyManager.enableOAuthAppCreation is boolean ? keyManager.enableOAuthAppCreation : true;
        keyManagerDTO.configuration = additionalProperties.toJsonString().toBytes();
        return keyManagerDTO;
    }
    private isolated function validateKeyManagerConfigurations(KeyManager keyManagerConfiguration, apk_keymanager_libs:KeyManagerConfigurations keyManagerConnectorConfigurations) returns boolean {
        KeyManager_certificates? certificates = keyManagerConfiguration.certificates;
        if certificates is KeyManager_certificates {
            if certificates.'type is () || certificates.value is () {
                return false;
            }
        } else {
            return false;
        }
        apk_keymanager_libs:EndpointConfiguration[] endpointsDefined = keyManagerConnectorConfigurations.endpoints;
        KeyManagerEndpoint[]? endpoints = keyManagerConfiguration.endpoints;
        if endpoints is KeyManagerEndpoint[] && endpoints.length() > 0 {
            foreach apk_keymanager_libs:EndpointConfiguration keymanagerEndpointConfiguration in endpointsDefined {
                if keymanagerEndpointConfiguration.required {
                    boolean found = false;
                    foreach KeyManagerEndpoint endpoint in endpoints {
                        if endpoint.name == keymanagerEndpointConfiguration.name {
                            found = true;
                            if endpoint.value.length() == 0 {
                                return false;
                            }
                            break;
                        }
                    }
                    if !found {
                        return false;
                    }
                }
            }
        } else {
            foreach apk_keymanager_libs:EndpointConfiguration keymanagerEndpointConfiguration in endpointsDefined {
                if keymanagerEndpointConfiguration.required {
                    return false;
                }
            }
        }
        // validate endpoint configurations.
        record {}? additionalProperties = keyManagerConfiguration.additionalProperties;
        apk_keymanager_libs:KeyManagerConfiguration[] endpointConfigurations = keyManagerConnectorConfigurations.endpointConfigurations;
        if additionalProperties is record {} && additionalProperties.length() > 0 {
            foreach apk_keymanager_libs:KeyManagerConfiguration endpointConfiguration in endpointConfigurations {
                if endpointConfiguration.required {
                    if !additionalProperties.hasKey(endpointConfiguration.name) {
                        return false;
                    }
                }
            }
        } else {
            foreach apk_keymanager_libs:KeyManagerConfiguration endpointConfiguration in endpointConfigurations {
                if endpointConfiguration.required {
                    return false;
                }
            }
        }
        return true;
    }

    public isolated function getAllKeyManagersByOrganization(commons:Organization organization) returns KeyManagerList|commons:APKError {
        KeyManagerListingDaoEntry[] allKeyManagersByOrganization = check getAllKeyManagersByOrganization(organization);
        KeyManagerList keyManagerList = {};
        KeyManagerInfo[] keyManagerInfoList = [];
        foreach KeyManagerListingDaoEntry item in allKeyManagersByOrganization {
            KeyManagerInfo keyManagerInfo = {
                name: item.name,
                'type: item.'type,
                id: item.uuid,
                description: item.description,
                enabled: item.enabled
            };
            keyManagerInfoList.push(keyManagerInfo);
        }
        keyManagerList.list = keyManagerInfoList;
        keyManagerList.count = keyManagerInfoList.length();
        return keyManagerList;
    }
    public isolated function getKeyManagerById(string id, commons:Organization organization) returns KeyManager|commons:APKError {
        KeyManagerDaoEntry keyManagerEntry = check getKeyManagerById(id, organization);
        return self.fromKeyManagerDaoEntryToKeyManagerModel(keyManagerEntry);
    }
    public isolated function updateKeyManager(string id, KeyManager updatedKeyManager, commons:Organization organization) returns KeyManager|commons:APKError {
        KeyManagerDaoEntry keyManagerEntry = check getKeyManagerById(id,organization);
        KeyManagerDaoEntry updatedKeyManagerEntry = self.fromKeyManagerModelToKeyManagerDaoEntry(updatedKeyManager);
        check updateKeyManager(id,updatedKeyManagerEntry, organization);
        return self.getKeyManagerById(id, organization);
    }
    public isolated function deleteKeyManager(string id,commons:Organization organization) returns commons:APKError?{
        KeyManagerDaoEntry keyManagerEntry = check getKeyManagerById(id,organization);
        check deleteKeyManager(id,organization);
    }
    private isolated function fromKeyManagerDaoEntryToKeyManagerModel(KeyManagerDaoEntry keyManagerDaoEntry) returns KeyManager|commons:APKError {
        do {
            string additionalPropertiesString = check string:fromBytes(<byte[]>keyManagerDaoEntry.configuration);
            json additionalPropertiesJson = check value:fromJsonString(additionalPropertiesString);
            KeyManager keymanager = {
                id: keyManagerDaoEntry.uuid,
                name: keyManagerDaoEntry.name,
                'type: keyManagerDaoEntry.'type,
                description: keyManagerDaoEntry.description,
                issuer: keyManagerDaoEntry.issuer,
                enabled: keyManagerDaoEntry.enabled
            };
            KeyManagerEndpoint[] endpoints = [];
            record {} additionalProperties = check additionalPropertiesJson.cloneWithType();
            if additionalProperties.hasKey("endpoints") {
                record {} endpointsInRecord = <record {|anydata...;|}>additionalProperties.get("endpoints");
                foreach string key in endpointsInRecord.keys() {
                    KeyManagerEndpoint endpoint = {
                        name: key,
                        value: <string>endpointsInRecord.get(key)
                    };
                    endpoints.push(endpoint);
                }
                _ = additionalProperties.removeIfHasKey("endpoints");
                keymanager.endpoints = endpoints;
            }
            if additionalProperties.hasKey("grantTypes") {
                string[] grantTypes = <string[]>additionalProperties.get("grantTypes");
                keymanager.availableGrantTypes = grantTypes;
                _ = additionalProperties.removeIfHasKey("grantTypes");
            }
            if additionalProperties.hasKey("consumerKeyClaim") {
                keymanager.consumerKeyClaim = <string>additionalProperties.get("consumerKeyClaim");
                _ = additionalProperties.removeIfHasKey("consumerKeyClaim");
            }
            if additionalProperties.hasKey("scopesClaim") {
                keymanager.scopesClaim = <string>additionalProperties.get("scopesClaim");
                _ = additionalProperties.removeIfHasKey("scopesClaim");
            }
            if additionalProperties.hasKey("certificate_type") {
                string certificateType = <string>additionalProperties.get("certificate_type");
                if additionalProperties.hasKey("certificate_value") {
                    string certificateValue = <string>additionalProperties.get("certificate_value");
                    KeyManager_certificates certificates = {
                        'type: certificateType,
                        value: certificateValue
                    };
                    keymanager.certificates = certificates;
                    _ = additionalProperties.removeIfHasKey("certificate_type");
                    _ = additionalProperties.removeIfHasKey("certificate_value");
                }
            }
            if additionalProperties.hasKey("mapOAuthConsumerApps") {
                keymanager.enableMapOAuthConsumerApps = <boolean>additionalProperties.get("mapOAuthConsumerApps");
                _ = additionalProperties.removeIfHasKey("mapOAuthConsumerApps");
            }
            if additionalProperties.hasKey("enableTokenGeneration") {
                keymanager.enableTokenGeneration = <boolean>additionalProperties.get("enableTokenGeneration");
                _ = additionalProperties.removeIfHasKey("enableTokenGeneration");
            }
            if additionalProperties.hasKey("enableOauthAppCreation") {
                keymanager.enableOAuthAppCreation = <boolean>additionalProperties.get("enableOauthAppCreation");
                _ = additionalProperties.removeIfHasKey("enableOauthAppCreation");
            }
            keymanager.additionalProperties = additionalProperties;
            return keymanager;
        } on fail var e {
            return e909438(e);
        }
    }

}

