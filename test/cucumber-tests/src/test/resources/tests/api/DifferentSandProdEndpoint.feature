Feature: API different endpoint resource level
  Scenario: Testing different endpoint resource level
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/different_sand_prod_endpoint.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-different-sand-prod-endpoint/3.14/endpoint1" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://backend/anything/prodr/endpoint1"
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-different-sand-prod-endpoint/endpoint2" with body ""
    Then the response status code should be 200
    And the response body should contain "https://backend/anything/prod/endpoint2"
    And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/test-different-sand-prod-endpoint/3.14/endpoint1" with body ""
    Then the response status code should be 200
    And the response body should contain "https://backend/anything/sandr/endpoint1"
    And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/test-different-sand-prod-endpoint/endpoint2" with body ""
    Then the response status code should be 200
    And the response body should contain "https://backend/anything/sand/endpoint2"


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | test-different-sand-prod-endpoint | 202 |
