Feature: API different endpoint resource level
  Scenario: Testing different endpoint resource level
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/different_endpoint_resource_level_new.yaml"
    And the definition file "artifacts/definitions/employees_api_new.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-different-endpoint-resource-level/3.14/endpoint1" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response body should contain "https://default.gw.wso2.com:9095/anything/base1/endpoint1"
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-different-endpoint-resource-level/endpoint2" with body ""
    Then the response status code should be 200
    And the response body should contain "https://default.gw.wso2.com:9095/anything/base2/endpoint2"


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                  | expectedStatusCode |
      | different-endpoint-resource-level-test | 202                |
