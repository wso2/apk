Feature: CORS Policy handling
  Background:
    Given The system is ready
  Scenario: Testing CORS Policy for a REST API
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/cors_api.json"
    When the definition file "artifacts/definitions/cors-definition.json"
    And make the import API Creation request using OAS "File"
    Then the response status code should be 201
    And the response body should contain "test-cors"
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
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should not contain
      | Access-Control-Allow-Origin      |
      | Access-Control-Allow-Credentials |
      | Access-Control-Allow-Methods     |
      | Access-Control-Allow-Headers     |
      | Access-Control-Max-Age           |
    Then I set headers
      | Origin | test.domain.com |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should not contain
      | Access-Control-Allow-Origin      |
      | Access-Control-Allow-Credentials |
      | Access-Control-Allow-Methods     |
      | Access-Control-Allow-Headers     |
      | Access-Control-Max-Age           |
    Then I set headers
      | Origin | abc.com |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should contain
      | Access-Control-Allow-Origin      | abc.com |
      | Access-Control-Allow-Credentials | true    |
    Then I set headers
      | Origin                        | abc.com |
      | Access-Control-Request-Method | GET     |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should contain
      | Access-Control-Allow-Origin      | abc.com                     |
      | Access-Control-Allow-Credentials | true                        |
      | Access-Control-Allow-Methods     | GET, PUT, POST, DELETE      |
      | Access-Control-Allow-Headers     | authorization, Content-Type |
  
  Scenario: Undeploying an already existing REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "cors"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 404 response code, not accepting
      |200|

  Scenario: Testing CORS Policy for a GraphQL API
    And I have a DCR application
    And I have a valid Publisher access token
    When the definition file "artifacts/definitions/schema_graphql.graphql"
    When I use the Payload file "artifacts/payloads/gql_cors.json"
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
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/" with body "{\"query\":\"{ anything }\"}"
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should not contain
      | Access-Control-Allow-Origin      |
      | Access-Control-Allow-Credentials |
      | Access-Control-Allow-Methods     |
      | Access-Control-Allow-Headers     |
      | Access-Control-Max-Age           |
    Then I set headers
      | Origin | test.domain.com |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/" with body "{\"query\":\"{ anything }\"}"
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should not contain
      | Access-Control-Allow-Origin      |
      | Access-Control-Allow-Credentials |
      | Access-Control-Allow-Methods     |
      | Access-Control-Allow-Headers     |
      | Access-Control-Max-Age           |
    Then I set headers
      | Origin | abc.com |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/" with body "{\"query\":\"{ anything }\"}"
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should contain
      | Access-Control-Allow-Origin      | abc.com |
      | Access-Control-Allow-Credentials | true    |
    Then I set headers
      | Origin                        | abc.com |
      | Access-Control-Request-Method | GET     |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/" with body "{\"query\":\"{ anything }\"}"
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should contain
      | Access-Control-Allow-Origin      | abc.com                     |
      | Access-Control-Allow-Credentials | true                        |
      | Access-Control-Allow-Methods     | GET, PUT, POST, DELETE      |
      | Access-Control-Allow-Headers     | authorization, Access-Control-Allow-Origin |


  
  Scenario: Undeploying an already existing GraphQL API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "TestApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "StarWarsAPI"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/" with body "{\"query\":\"{ anything }\"}"
    And I eventually receive 404 response code, not accepting
      |200|