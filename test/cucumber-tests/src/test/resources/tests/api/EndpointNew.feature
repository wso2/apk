Feature: Endpoint
  Scenario: Testing API level and resource level endpoints
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/endpoint_conf_new.yaml"
    And the definition file "artifacts/definitions/employees_api_new.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken}         |
    And I send "GET" request to "https://default.gw.wso2.com:9095/endpoint/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://default.gw.wso2.com:9095/anything/employee"
    And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/endpoint/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://default.sandbox.gw.wso2.com:9095/anything/test/employee"
    And I send "POST" request to "https://default.gw.wso2.com:9095/endpoint/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://default.gw.wso2.com:9095/anything/test/employee"
    And I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/endpoint/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://default.sandbox.gw.wso2.com:9095/anything/employee"


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                  | expectedStatusCode  |
      | endpoint-test          | 202                 |
