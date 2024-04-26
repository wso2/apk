Feature: Testing the internal-key generation and invocation
  Background:
    Given The system is ready
  Scenario: Creating and invoking a REST API using Internal-Key
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api1.json"
    And I use the OAS URL "https://petstore3.swagger.io/api/v3/openapi.json"
    And make the import API Creation request using OAS "URL"
    Then the response status code should be 201
    And the response body should contain "SwaggerPetstore"
    And make the API Revision Deployment request
    Then the response status code should be 201
    And make the Change Lifecycle request
    Then the response status code should be 200
    Then I make an internal key generation request
    Then the response status code should be 200
    And the response body should contain "apikey"
    Then I set headers
        | Internal-Key  | ${internalKey} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
  
  Scenario: Undeploying an already existing REST API
    And I have a DCR application
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "SwaggerPetstore"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
    And I eventually receive 404 response code, not accepting
      |200|
