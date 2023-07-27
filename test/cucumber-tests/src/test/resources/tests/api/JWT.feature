Feature: Test JWT related functionalities
  Scenario: Test JWT authentication with valid and invalid access token
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/jwt_basic_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    Then I set headers
      |Authorization|bearer invalidToken|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And the response status code should be 401

  Scenario: Test disabled JWT configuration
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/jwt_disabled_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer invalidToken|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-disabled/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|

  Scenario: Test customized JWT headers
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/jwt_custom_header_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-custom-header/3.14/employee/" with body ""
    And I eventually receive 401 response code, not accepting
      |429|
      |200|
    Then I set headers
      |testAuth|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-custom-header/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    And the response body should contain "\"X-Jwt-Assertion\""
    And the decoded "X-Jwt-Assertion" jwt should contain
      | claim    | value    |
      | claim1 | value1 |
      | claim2 | value2 |

  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | jwt-basic-test          | 202                 |
      | jwt-disabled-test      | 202                 |
      | jwt-custom-header-test      | 202                 |
