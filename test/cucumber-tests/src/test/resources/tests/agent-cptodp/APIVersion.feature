Feature: Creating new versions of the APIs
  Background:
    Given The system is ready
  Scenario: Create a new version of a REST API and try to invoke both old and newer versions
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api1.json"
    And I use the OAS URL "https://petstore3.swagger.io/api/v3/openapi.json"
    And make the import API Creation request using OAS "URL"
    Then the response status code should be 201
    And the response body should contain "SwaggerPetstore"
    And make the API Revision Deployment request
    Then the response status code should be 201
    Then I wait for 40 seconds
    And make the Change Lifecycle request
    Then the response status code should be 200
    And I have a valid Devportal access token
    And make the Application Creation request with the name "SampleApp"
    Then the response status code should be 201
    And the response body should contain "SampleApp"
    And I have a KeyManager
    And make the Generate Keys request
    Then the response status code should be 200
    And the response body should contain "consumerKey"
    And the response body should contain "consumerSecret"
    And make the Subscription request
    Then the response status code should be 201
    And the response body should contain "Unlimited"
    And I get "production" oauth keys for application
    Then the response status code should be 200
    And make the Access Token Generation request for "production"
    Then the response status code should be 200
    And the response body should contain "accessToken"
    Then I set headers
        | Authorization             | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And I create the new version "2.0.0" of the same API with default version set to "false"
    Then the response status code should be 201
    And the response body should contain "SwaggerPetstore"
    And the response body should contain "2.0.0"
    And make the API Revision Deployment request
    Then the response status code should be 201
    And make the Change Lifecycle request
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/2.0.0/pet/4" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
  
  Scenario: Undeploy the created REST APIs
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "1.0.0"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
    And the response status code should be 404
    Then I find the apiUUID of the API created with the name "2.0.0"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/2.0.0/pet/4" with body ""
    And the response status code should be 404

  
  Scenario: Create a new version of a GraphQL API and try to invoke both old and newer versions
    And I have a DCR application
    And I have a valid Publisher access token
    When the definition file "artifacts/definitions/schema_graphql.graphql"
    When I use the Payload file "artifacts/payloads/gql_with_scopes.json"
    Then I make the import GraphQLAPI Creation request
    Then the response status code should be 201
    And the response body should contain "StarWarsAPI"
    And make the API Revision Deployment request
    Then the response status code should be 201
    Then I wait for 40 seconds
    And make the Change Lifecycle request
    Then the response status code should be 200  
    And I have a valid Devportal access token
    And make the Application Creation request with the name "TestApp"
    Then the response status code should be 201
    And the response body should contain "TestApp"
    And I have a KeyManager
    And make the Generate Keys request
    Then the response status code should be 200
    And the response body should contain "consumerKey"
    And the response body should contain "consumerSecret"
    And make the Subscription request
    Then the response status code should be 201
    And the response body should contain "Unlimited"
    And I get "production" oauth keys for application
    Then the response status code should be 200
    And make the Access Token Generation request for "production"
    Then the response status code should be 200
    And the response body should contain "accessToken"
    Then I set headers
        | Authorization             | Bearer ${accessToken} |
    And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
    Then the response status code should be 200
    And I eventually receive 200 response code, not accepting
      | 404 |
      | 401 |
    And I create the new version "3.2" of the same API with default version set to "true"
    Then the response status code should be 201
    And the response body should contain "StarWarsAPI"
    And the response body should contain "3.2"
    And make the API Revision Deployment request
    Then the response status code should be 201
    And make the Change Lifecycle request
    Then the response status code should be 200
    And make the Subscription request
    Then the response status code should be 201
    And the response body should contain "Unlimited"
    And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.2" with body "{\"query\":\"{ hero { name } }\"}"
    And I eventually receive 200 response code, not accepting
      |429|

  Scenario: Undeploying the created GraphQL APIs
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "TestApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "3.14"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
    And I eventually receive 404 response code, not accepting
      |200|
    Then I find the apiUUID of the API created with the name "3.2"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.2" with body "{\"query\":\"{ hero { name } }\"}"
    And I eventually receive 404 response code, not accepting
      |200|