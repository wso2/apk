Feature: Backend JWT auth
  Scenario: Testing API level Endpoint backend jwt auth header
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/backend_jwt_auth_conf.yaml"
    And the definition file "artifacts/definitions/backend_jwt_auth_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-jwt-security/3.14/get" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response body should not contain "\"X-Jwt-Assertion\""
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-jwt-security/3.14/employee" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response body should contain "\"X-Jwt-Assertion\""

  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID            | expectedStatusCode |
      | backend-jwt-test | 202                |
