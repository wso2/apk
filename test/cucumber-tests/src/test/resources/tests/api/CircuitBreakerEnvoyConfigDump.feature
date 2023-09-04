Feature: circuitBreakerMaxRequest
  Scenario: Testing backend timeout
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/circuit-breaker-config-dump.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I wait for 5 seconds
    And I send "GET" request to "http://default.gw.wso2.com:9000/config_dump" with body ""
    Then the response status code should be 200
    And the response body should contain "\"max_connections\": 1111"
    And the response body should contain "\"max_pending_requests\": 1112"
    And the response body should contain "\"max_requests\": 1113"
    And the response body should contain "\"max_retries\": 1114"
    And the response body should contain "\"max_connection_pools\": 1115"

  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | circuit-breaker-config-dump          | 202                 |
