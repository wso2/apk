import apk_keymanager_libs;
import ballerina/uuid;
import ballerina/lang.value;
import ballerina/log;
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
        KeyManagerDaoEntry keyManagerDtoToInsert = check self.fromKeyManagerModelToKeyManagerDaoEntry(keyManager);
        _ = check addKeyManagerEntry(keyManagerDtoToInsert, organization);
        return check self.getKeyManagerById(<string>keyManagerDtoToInsert.uuid, organization);
    }
    private isolated function fromKeyManagerModelToKeyManagerDaoEntry(KeyManager keyManager) returns KeyManagerDaoEntry|commons:APKError {
        do {
            KeyManagerDaoEntry keyManagerDTO = {
                name: keyManager.name,
                'type: keyManager.'type,
                issuer: <string>keyManager.issuer,
                enabled: keyManager.enabled,
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
                    } else {
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
            KeyManager_signingCertificate? certificates = keyManager.signingCertificate;
            if certificates is KeyManager_signingCertificate {
                if certificates.'type is string {
                    additionalProperties["signing_certificate_type"] = certificates.'type;
                }
                string? certificateValue = certificates.value;
                if certificateValue is string {
                    if certificates.'type == "JWKS" {
                        additionalProperties["signing_certificate_value"] = certificateValue;
                    } else {
                        byte[] encodedBytes = check commons:EncoderUtil_encodeBase64(certificateValue.toBytes());
                        additionalProperties["signing_certificate_value"] = check string:fromBytes(encodedBytes);
                    }
                    additionalProperties["signing_certificate_value"] = certificateValue;
                }
            }
            string? tlsCertificate = keyManager.tlsCertificate;
            if tlsCertificate is string {
                byte[] encodedBytes = check commons:EncoderUtil_encodeBase64(tlsCertificate.toBytes());
                additionalProperties["tls_certificate"] = check string:fromBytes(encodedBytes);
            }
            additionalProperties["mapOAuthConsumerApps"] = keyManager.enableMapOAuthConsumerApps;
            additionalProperties["enableTokenGeneration"] = keyManager.enableTokenGeneration;
            additionalProperties["enableOauthAppCreation"] = keyManager.enableOAuthAppCreation;
            additionalProperties["enableOauthAppValidation"] = keyManager.enableOauthAppValidation;

            return keyManagerDTO;
        } on fail var e {
            log:printError("Error while converting key manager model to key manager dto: " + e.message());
            return e909438(e);
        }
    }
    private isolated function validateKeyManagerConfigurations(KeyManager keyManagerConfiguration, apk_keymanager_libs:KeyManagerConfigurations keyManagerConnectorConfigurations) returns boolean {
        KeyManager_signingCertificate? certificates = keyManagerConfiguration.signingCertificate;
        if certificates is KeyManager_signingCertificate {
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
        KeyManagerDaoEntry keyManagerEntry = check getKeyManagerById(id, organization);
        KeyManagerDaoEntry updatedKeyManagerEntry = check self.fromKeyManagerModelToKeyManagerDaoEntry(updatedKeyManager);
        check updateKeyManager(id, updatedKeyManagerEntry, organization);
        return self.getKeyManagerById(id, organization);
    }
    public isolated function deleteKeyManager(string id, commons:Organization organization) returns commons:APKError? {
        KeyManagerDaoEntry keyManagerEntry = check getKeyManagerById(id, organization);
        check deleteKeyManager(id, organization);
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
            if additionalProperties.hasKey("signing_certificate_type") {
                string certificateType = <string>additionalProperties.get("signing_certificate_type");
                if additionalProperties.hasKey("signing_certificate_value") {
                    string certificateValue = <string>additionalProperties.get("signing_certificate_value");
                    if certificateType == "PEM" {
                        byte[] encodedBytes = check commons:EncoderUtil_decodeBase64(certificateValue.toBytes());
                        certificateValue = check string:fromBytes(encodedBytes);
                    }
                    KeyManager_signingCertificate certificates = {
                        'type: certificateType,
                        value: certificateValue
                    };
                    keymanager.signingCertificate = certificates;
                }
                _ = additionalProperties.removeIfHasKey("signing_certificate_type");
                _ = additionalProperties.removeIfHasKey("signing_certificate_value");
            }
            if additionalProperties.hasKey("tls_certificate") {
                string certificateValue = <string>additionalProperties.get("tls_certificate");
                byte[] encodedBytes = check commons:EncoderUtil_decodeBase64(certificateValue.toBytes());
                certificateValue = check string:fromBytes(encodedBytes);
                keymanager.tlsCertificate = certificateValue;
                _ = additionalProperties.removeIfHasKey("tls_certificate");
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
            if additionalProperties.hasKey("enableOauthAppValidation"){
                keymanager.enableOauthAppValidation = <boolean>additionalProperties.get("enableOauthAppValidation");
                _ = additionalProperties.removeIfHasKey("enableOauthAppValidation");
            }
            keymanager.additionalProperties = additionalProperties;
            return keymanager;
        } on fail var e {
            return e909438(e);
        }
    }
}

