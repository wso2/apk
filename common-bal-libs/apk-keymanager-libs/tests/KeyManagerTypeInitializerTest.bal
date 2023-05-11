import ballerina/test;

@test:Config {}
public function testReadKeyManagerYamlFiles() returns error? {
    KeyManagerTypeInitializer keyManagerTypeInitializer = new ();
    _ = check keyManagerTypeInitializer.initialize("./tests/resources/validKmConfigs");
    KeyManagerConfigurations? retrieveKeyManagerConfigByType = keyManagerTypeInitializer.retrieveKeyManagerConfigByType("Okta");
    if retrieveKeyManagerConfigByType is KeyManagerConfigurations {
        test:assertEquals(retrieveKeyManagerConfigByType.'type, "Okta");
        test:assertEquals(retrieveKeyManagerConfigByType.consumerKeyClaim, "azp");
        test:assertEquals(retrieveKeyManagerConfigByType.scopesClaim, "scp");
        test:assertEquals(retrieveKeyManagerConfigByType.'display_name, "Okta");
    } else {
        test:assertFail("Error while retrieving key manager configurations");
    }
}
