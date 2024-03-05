Feature: API Deployment
  Background:
    Given The system is ready
  Scenario: Import an API, Create Application, Generate Keys, Subscribe to an API
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api1.json"
    And I use the OAS URL "https://petstore3.swagger.io/api/v3/openapi.json"
    And make the import API Creation request
    Then the response status code should be 201
    And the response body should contain "SwaggerPetstore"
    And make the API Revision Deployment request
    Then the response status code should be 201
    And make the Change Lifecycle request
    Then the response status code should be 200
    And I have a valid Devportal access token
    And make the Application Creation request
    Then the response status code should be 201
    And the response body should contain "PetstoreApp"
    And I have a KeyManager
    And make the Generate Keys request
    Then the response status code should be 200
    And the response body should contain "consumerKey"
    And the response body should contain "consumerSecret"
    And make the Subscription request
    Then the response status code should be 201
    And the response body should contain "Gold"
    And I get oauth keys for application
    Then the response status code should be 200
    And make the Access Token Generation request
    Then the response status code should be 200
    And the response body should contain "accessToken"
    And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/5" with body ""
    And I eventually receive 200 response code, not accepting
      |429|

#  Scenario Outline: Undeploy an API
#    Given The system is ready
#    And I have a valid subscription
#    When I undeploy the API whose ID is "<apiID>"
#    Then the response status code should be <expectedStatusCode>

#    Examples:
#      | apiID | expectedStatusCode |
#      | cors-api-adff3dbc-2787-11ee-be56-0242ac120002  | 202         |
