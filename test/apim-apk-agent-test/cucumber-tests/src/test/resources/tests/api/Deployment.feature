Feature: API Deployment
  Scenario: Import an API
    Given The system is ready
    And I have a DCR application for Publisher
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api1.json"
    And I use the OAS URL "https://petstore3.swagger.io/api/v3/openapi.json"
    And make the import API Creation request
    Then the response status code should be 201
    And the response body should contain "SwaggerPetstore"
    And make the API Revision Deployment request
    Then the response status code should be 201

#  Scenario: Deploying an API
#    Given The system is ready
#    And I have a valid subscription
#    When I use the APK Conf file "artifacts/apk-confs/cors_API.apk-conf"
#    And the definition file "artifacts/definitions/cors_api.yaml"
#    And make the API deployment request
#    Then the response status code should be 200
#    And the response body should contain "cors-api-adff3dbc-2787-11ee-be56-0242ac120002"
#
#  Scenario Outline: Undeploy an API
#    Given The system is ready
#    And I have a valid subscription
#    When I undeploy the API whose ID is "<apiID>"
#    Then the response status code should be <expectedStatusCode>

#    Examples:
#      | apiID | expectedStatusCode |
#      | cors-api-adff3dbc-2787-11ee-be56-0242ac120002  | 202         |
