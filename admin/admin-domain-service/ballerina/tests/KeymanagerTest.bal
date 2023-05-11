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
            {name:"introspection_endpoint",value:"https://localhost:9443/api/introspect"},
            {name: "token_endpoint", value: "https://localhost:9443/oauth2/token"}
        ],
        additionalProperties: {
            "client_id": "abcde",
            "client_secret": "abcde"
        },
        certificates: {'type: "JWKS",value: "https://localhost:9443/oauth2/jwks"}
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
