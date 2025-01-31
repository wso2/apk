Feature: circuitBreakerMaxRequest
  Scenario: Testing backend timeout
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/circuit-breaker-max-request-test.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request/3.14/anything/test" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |500|
    And I send "GET" async request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request/3.14/delay/10" with body ""
    And I send "GET" async request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request/3.14/delay/10" with body ""
    And I wait for 2 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request/3.14/delay/10" with body ""
    Then the response status code should be 503

    When I use the APK Conf file "artifacts/apk-confs/circuit-breaker-max-request-test-v1.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request-v1/3.14/anything/test" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |500|
    And I send "GET" async request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request-v1/3.14/delay/10" with body ""
    And I send "GET" async request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request-v1/3.14/delay/10" with body ""
    And I wait for 2 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/circuit-breaker-max-request-v1/3.14/delay/10" with body ""
    Then the response status code should be 200



  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | circuit-breaker-max-request-test          | 202                 |
