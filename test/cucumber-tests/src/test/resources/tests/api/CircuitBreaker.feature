Feature: BackendRetry
  Scenario: Testing backend retry
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/circuit_breaker_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "POST" request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/set-retry-count" with body "{\"count\": 500}"
    And I eventually receive 200 response code, not accepting
      |429|
    And I send "POST" request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/reset" with body ""
    Then the response status code should be 200
    And I send "GET" async request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/retry" with body ""
    And I send "GET" async request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/retry" with body ""
    And I wait for 2 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/retry" with body ""
    Then the response status code should be 500
    And I send "POST" request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/set-retry-count" with body "{\"count\": 3}"
    Then the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/circuit-breaker/3.14/reset" with body ""
    Then the response status code should be 200



  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | circuit-breaker-test          | 202                 |
