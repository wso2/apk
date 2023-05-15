import ballerina/lang.value;
import wso2/apk_common_lib as commons;

isolated function fromKeyManagerDaoEntryToKeyManagerModel(KeyManagerDaoEntry keyManagerDaoEntry) returns KeyManager|commons:APKError {
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
            _ = additionalProperties.removeIfHasKey("tls_certificate");
            string certificateValue = <string>additionalProperties.get("tls_certificate");
            byte[] encodedBytes = check commons:EncoderUtil_decodeBase64(certificateValue.toBytes());
            certificateValue = check string:fromBytes(encodedBytes);
            keymanager.tlsCertficate = certificateValue;
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

public type KeyManager record {
    string id?;
    string name;
    string displayName?;
    string 'type;
    string description?;
    KeyManagerEndpoint[] endpoints?;
    KeyManager_signingCertificate signingCertificate?;
    # PEM type certificate
    string tlsCertficate?;
    string issuer;
    string[] availableGrantTypes?;
    boolean enableTokenGeneration?;
    boolean enableMapOAuthConsumerApps = false;
    boolean enableOAuthAppCreation = true;
    string consumerKeyClaim?;
    string scopesClaim?;
    boolean enabled = true;
    record {} additionalProperties?;
};

public type KeyManagerEndpoint record {
    string name;
    string value;
};

public type KeyManager_signingCertificate record {
    string 'type?;
    string value?;
};

isolated function e909438(error? e) returns commons:APKError {
    return error commons:APKError("Internal Server Error", e,
        code = 909438,
        message = "Internal Server Error",
        statusCode = 500,
        description = "Internal Server Error"
    );
}
