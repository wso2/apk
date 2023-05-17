import ballerina/lang.value;
import wso2/apk_common_lib as commons;
import devportal_service.types;

isolated function fromKeyManagerDaoEntryToKeyManagerModel(KeyManagerDaoEntry keyManagerDaoEntry) returns types:KeyManager|commons:APKError {
    do {
        string additionalPropertiesString = check string:fromBytes(<byte[]>keyManagerDaoEntry.configuration);
        json additionalPropertiesJson = check value:fromJsonString(additionalPropertiesString);
        types:KeyManager keymanager = {
            id: keyManagerDaoEntry.uuid,
            name: keyManagerDaoEntry.name,
            'type: keyManagerDaoEntry.'type,
            description: keyManagerDaoEntry.description,
            issuer: keyManagerDaoEntry.issuer,
            enabled: keyManagerDaoEntry.enabled
        };
        map<string> endpoints = {};
        record {} additionalProperties = check additionalPropertiesJson.cloneWithType();
        if additionalProperties.hasKey("endpoints") {
            record {} endpointsInRecord = <record {|anydata...;|}>additionalProperties.get("endpoints");
            foreach string key in endpointsInRecord.keys() {
                endpoints[key] = <string>endpointsInRecord[key];
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
        if additionalProperties.hasKey("tls_certificate") {
            string certificateValue = <string>additionalProperties.get("tls_certificate");
            byte[] encodedBytes = check commons:EncoderUtil_decodeBase64(certificateValue.toBytes());
            certificateValue = check string:fromBytes(encodedBytes);
            keymanager.tlsCertficate = certificateValue;
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
        if additionalProperties.hasKey("enableOauthAppValidation") {
            keymanager.enableOauthAppValidation = <boolean>additionalProperties.get("enableOauthAppValidation");
            _ = additionalProperties.removeIfHasKey("enableOauthAppValidation");
        }
        keymanager.additionalProperties = additionalProperties;
        return keymanager;
    } on fail var e {
        return e909438(e);
    }
}

isolated function e909438(error? e) returns commons:APKError {
    return error commons:APKError("Internal Server Error", e,
        code = 909438,
        message = "Internal Server Error",
        statusCode = 500,
        description = "Internal Server Error"
    );
}
