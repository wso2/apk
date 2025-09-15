Feature: API Deployment
  Background:
    Given The system is ready
  Scenario: Import an API, Create Application, Generate Keys, Subscribe to an API
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/sample_api.json"
    And I use the OAS URL "https://raw.githubusercontent.com/wso2/docs-apk/refs/heads/1.3.0/en/docs/assets/files/get-started/SampleAPIDefinition.json"
    And make the import API Creation request using OAS "URL"
    Then the response status code should be 201
    And the response body should contain "Sample API"
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
    And I send "GET" request to "https://default.gw.wso2.com:9095/sample-api/1.0.0/uuid" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
  
  Scenario: Undeploying an already existing REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "Sample API"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/sample-api/1.0.0/uuid" with body ""
    And I eventually receive 404 response code, not accepting
      |200|
    Then I wait for 5 minute

  Scenario: Deploying a GraphQL API
    And I have a DCR application
    And I have a valid Publisher access token
    When the definition file "artifacts/definitions/schema_graphql.graphql"
    When I use the Payload file "artifacts/payloads/gqlPayload.json"
    Then I make the import GraphQLAPI Creation request
    Then the response status code should be 201
    And the response body should contain "StarwarsAPI"
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
    And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14/hero" with body "{\"query\":\"{ hero { name } }\"}"
    Then the response status code should be 200
    And I eventually receive 200 response code, not accepting
      | 404 |
      | 401 |

  Scenario: Undeploying an already existing GraphQL API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "TestApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "StarwarsAPI"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14/hero" with body "{\"query\":\"{ hero { name } }\"}"
    And I eventually receive 404 response code, not accepting
      |200|