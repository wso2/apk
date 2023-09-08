Feature: Security
  Scenario: Testing security related features
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/security.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    Then I remove headers
      |Authorization|
    And I send "POST" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body "{}"
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 401

#   Get another token

    And I have a valid subscription
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}u|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 401
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}u|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 401

#   Get another token

    And I have a valid subscription
    Then I set headers
      |Authorization|bearer ${accessToken}u|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 401
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}u|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 401
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-security/3.14/employee" with body ""
    Then the response status code should be 200



  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | security-test   | 202                 |
