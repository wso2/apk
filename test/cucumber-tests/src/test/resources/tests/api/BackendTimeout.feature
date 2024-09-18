Feature: BackendTimeout
  Scenario: Testing backend timeout
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/backend_timeout_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-timeout/3.14/delay/2" with body ""
    And I eventually receive 504 response code, not accepting
      |429|
      |200|
    And the response body should contain "timeout"
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-timeout/3.14/delay/0" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |500|
    And I send "POST" request to "https://default.gw.wso2.com:9095/backend-timeout/3.14/delay/3" with body ""
    And I eventually receive 504 response code, not accepting
      |429|
      |200|
    And the response body should contain "timeout"
    And I send "POST" request to "https://default.gw.wso2.com:9095/backend-timeout/3.14/delay/1" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |500|
    And I send "PUT" request to "https://default.gw.wso2.com:9095/backend-timeout/3.14/delay/4" with body ""
    And I eventually receive 504 response code, not accepting
      |429|
      |200|
    And the response body should contain "timeout"
    And I send "PUT" request to "https://default.gw.wso2.com:9095/backend-timeout/3.14/delay/2" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |500|


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | backend-timeout-test          | 202                 |
