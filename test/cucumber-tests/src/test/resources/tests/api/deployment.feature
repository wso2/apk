Feature: API Deployment
  Scenario: Deploying an API without api create scope
    Given The system is ready
    And I have a valid subscription without api deploy permission
    When I use the APK Conf file "artifacts/apk-confs/cors_API.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 403
    
  Scenario: Deploying an API
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/cors_API.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "cors-api-adff3dbc-2787-11ee-be56-0242ac120002"

  Scenario: Deploying an API with invalid APK Conf file
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/invalid_cors_API.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 400
    And the response body should contain
      |"#/corsConfiguration/corsConfigurationEnabled: expected type: Boolean, found: String"|
  
  Scenario Outline: Undeploy an API without api create scope
    Given The system is ready
    And I have a valid subscription without api deploy permission
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be 403

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID | expectedStatusCode |
      | cors-api-adff3dbc-2787-11ee-be56-0242ac120002  | 202         |
      | abcdeadsxzads | 404        |
