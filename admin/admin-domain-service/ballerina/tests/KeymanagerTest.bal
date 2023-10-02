import wso2/apk_common_lib as commons;
import ballerina/uuid;
import ballerina/log;
import ballerina/test;

@test:Config {}
public function addKeyManager() {
    KeyManagerClient keyManagerClient = new;
    commons:Organization organization = {enabled: true, uuid: uuid:createType1AsString(), name: "default", displayName: "default", organizationClaimValue: ""};
    KeyManager keyManager = {
        name: "abcde",
        'type: "Okta",
        availableGrantTypes: [],
        enabled: true,
        issuer: "https://localhost:9443/oauth2/token",
        endpoints: [
            {name: "dcr_endpoint", value: "https://localhost:9443/api/dcr"},
            {name: "introspection_endpoint", value: "https://localhost:9443/api/introspect"},
            {name: "token_endpoint", value: "https://localhost:9443/oauth2/token"}
        ],
        additionalProperties: {
            "client_id": "abcde",
            "client_secret": "abcde"
        },
        signingCertificate: {'type: "JWKS", value: "https://localhost:9443/oauth2/jwks"}
    };
    KeyManager|commons:APKError keyManagerEntryToOrganization = keyManagerClient.addKeyManagerEntryToOrganization(keyManager, organization);
    if keyManagerEntryToOrganization is KeyManager {
        test:assertTrue(!(keyManagerEntryToOrganization.id is ()));
        test:assertEquals(keyManagerEntryToOrganization.name, keyManager.name);
    } else {
        log:printError("failed to insert", keyManagerEntryToOrganization);
        test:assertFail();
    }
    KeyManagerList|commons:APKError allKeyManagersByOrganization = keyManagerClient.getAllKeyManagersByOrganization(organization);
    if allKeyManagersByOrganization is KeyManagerList {
        test:assertEquals(allKeyManagersByOrganization.count, 1);
    } else {
        log:printError("failed to retrieve all", allKeyManagersByOrganization);
        test:assertFail();
    }
    if keyManagerEntryToOrganization is KeyManager {
        KeyManager|error keyManagerById = keyManagerClient.getKeyManagerById(keyManagerEntryToOrganization.id ?: "", organization);
        if keyManagerById is KeyManager {
            test:assertEquals(keyManagerById, keyManagerEntryToOrganization);
        }
        // update KeyManager 
        KeyManager updatedKeyManager = keyManagerEntryToOrganization.clone();
        updatedKeyManager.description = "updated text";
        KeyManager|commons:APKError keyManagerResult = keyManagerClient.updateKeyManager(<string>keyManagerEntryToOrganization.id, updatedKeyManager, organization);
        if keyManagerResult is KeyManager {
            test:assertEquals(keyManagerResult, updatedKeyManager);
        } else {
            test:assertFail();
        }
        commons:APKError? deleteKeyManagerResponse = keyManagerClient.deleteKeyManager(<string>keyManagerEntryToOrganization.id, organization);
        if deleteKeyManagerResponse is commons:APKError {
            test:assertFail();
        }
        KeyManager|commons:APKError keyManagerByIdResult = keyManagerClient.getKeyManagerById(<string>keyManagerEntryToOrganization.id, organization);
        if keyManagerByIdResult is commons:APKError {
            test:assertEquals(keyManagerByIdResult.detail().code, 909439);
        }
    }
}

@test:Config {}
public function addKeyManagerNegative() {
    KeyManagerClient keyManagerClient = new;
    commons:Organization organization = {enabled: true, uuid: uuid:createType1AsString(), name: "default", displayName: "default", organizationClaimValue: ""};
    KeyManager keyManager = {
        name: "abcde",
        'type: "Okta",
        availableGrantTypes: [],
        enabled: true,
        issuer: "https://localhost:9443/oauth2/token",
        endpoints: [
            {name: "token_endpoint", value: "https://localhost:9443/oauth2/token"}
        ],
        additionalProperties: {
            "client_id": "abcde"
        }
    };
    KeyManager|commons:APKError keyManagerEntryToOrganization = keyManagerClient.addKeyManagerEntryToOrganization(keyManager, organization);
    if keyManagerEntryToOrganization is commons:APKError {
        test:assertEquals(keyManagerEntryToOrganization.detail().code, 909437);
    }
}

@test:Config {}
public function addKeyManagerWithTlsCert() {
    commons:Organization organization = {enabled: true, uuid: uuid:createType1AsString(), name: "default", displayName: "default", organizationClaimValue: ""};

    KeyManager keyManager = {
        "name": "nonprod-idp2",
        "displayName": "Non production IDP",
        "type": "Okta",
        "description": "This is a key manager for Developers",
        "endpoints": [
            {
                "name": "token_endpoint",
                "value": "https://keymanager-wso2-apk-idp-ds-service:9443/oauth2/token"
            },
            {
                "name": "dcr_endpoint",
                "value": "https://keymanager-wso2-apk-idp-ds-service:9443/dcr"
            }
        ],
        "signingCertificate": {
            "type": "JWKS",
            "value": "https://keymanager-wso2-apk-idp-ds-service:9443/oauth2/jwks"
        },
        "tlsCertificate": "-----BEGIN CERTIFICATE----- MIIC/TCCAeWgAwIBAgIUd4njv8ySPgo7t0F1e2aJEo9TpQ4wDQYJKoZIhvcNAQEL BQAwDjEMMAoGA1UEAwwDYXBrMB4XDTIzMDMwOTA4MTYzNFoXDTMzMDMwNjA4MTYz NFowDjEMMAoGA1UEAwwDYXBrMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC AQEA1G9BAcVa3aMemxERY+J8UmH/W8JpcEOtdeNX5YNihnNOXtlvhSzFzOPaK97P r4IqGASbVNiR+1J6WEoi/b6ZJTx0q3YUn0YlQJrz7g20TdoGJjxGVWzn0EW4beHX Gq60vXLf4t3mlLCLGIK3kJTWAoRzd74djV7+5v0Bm/6KBBAWcu5UbOD9KRpOsxGM n3Z0103oAGViyq84QtFvhXVNWttDLe2jU/7o42ddaJozRL9z+1AepdoWPyJZIZqU bXcGAk7idk7c/8dKMxwAm3CV/WvgWrVK5R+YTiGqRf5pd9WWCydEVQkNqCZgTPNy BTRvHo52onPnT6ALtMI0mnWLtQIDAQABo1MwUTAdBgNVHQ4EFgQUzxcA8ceCF5t+ vPeOpYbi11CWjwcwHwYDVR0jBBgwFoAUzxcA8ceCF5t+vPeOpYbi11CWjwcwDwYD VR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAFCKd5x8z64p1j9BnoSAl JJuNz4uWjv/YK94B3s3KUkXOXjHuLuazkSss2UHFb+KeLd6vhimM5QqkFsUEMsnC 8RyBZP58kghcLzJld5Uiwlnr/RuOANvcv4eKSavwu5/ABXuMUQb/GvKtWPrr2VbJ ULg3p7NGigXHHg84eVMA7oNX1Z5R2cS4ISklWXm5SpMPh+SNCgqwqhxRNYJ2J0EZ qlp4ofQG3GJ72J+DRHlNujWEskP5IJjw6w8Q0zjXx26yelGe2+TM6BB7PpCN6kNU zHo2k/575bu2iZztnYVmE74H1W3cXJ7c0q82uUFvdW+FlRtm+OPIiIGK74lFiZNB OQ== -----END CERTIFICATE-----",
        "issuer": "https://idp.am.wso2.com/token",
        "availableGrantTypes": [
            "client_credentials"
        ],
        "enableTokenGeneration": true,
        "enableMapOAuthConsumerApps": true,
        "enableOAuthAppCreation": true,
        "consumerKeyClaim": "azp",
        "scopesClaim": "scopes",
        "enabled": true,
        "additionalProperties": {
            "client_id": "abcde",
            "client_secret": "abcde"
        }
    };
    KeyManagerClient keyManagerClient = new;
    KeyManager|commons:APKError keyManagerEntryToOrganization = keyManagerClient.addKeyManagerEntryToOrganization(keyManager, organization);
    if keyManagerEntryToOrganization is KeyManager {
        test:assertEquals(keyManagerEntryToOrganization.name, "nonprod-idp2");
    }
}

