Feature: Test JWT related functionalities
  Scenario: Test JWT authentication with valid and invalid access token
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/jwt_basic_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    Then I set headers
      |Authorization|Bearer invalidToken|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And the response status code should be 401
    Then I remove header "Authorization"
    Then I set headers
      | custom-jwt | ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And the response status code should be 200
  Scenario: Test JWT Token from different issuer with JWKS
    Given The system is ready
    Then I generate JWT token from idp1 with kid "123-456"
    Then I set headers
      |Authorization|Bearer ${idp-1-token}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    Then I set headers
      |Authorization|Bearer "${idp-1-token}h"|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And I eventually receive 401 response code, not accepting
      |429|
      |200|
    Then I set headers
      |Authorization|Bearer ${idp-1-token}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    Then I generate JWT token from idp1 with kid "456-789"
    And I send "DELETE" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/1234" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    Then I set headers
      |Authorization|Bearer ${idp-1-token}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-basic/3.14/employee/" with body ""
    And I eventually receive 401 response code, not accepting
      |429|
      |200|
  Scenario: Test disabled JWT configuration
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/jwt_disabled_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer invalidToken|
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
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-custom-header/3.14/employee/" with body ""
    And I eventually receive 401 response code, not accepting
      |429|
      |200|
    Then I set headers
      |testAuth|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-custom-header/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    And the response body should contain "\"X-Jwt-Assertion\""
    And the "X-Jwt-Assertion" jwt should validate from JWKS "https://api.am.wso2.com:9095/.wellknown/jwks" and contain
      | claim    | value    |
      | claim1 | value1 |
      | claim2 | value2 |

    Scenario: Test customized JWT headers with Resource Endpoint
      Given The system is ready
      And I have a valid subscription
      When I use the APK Conf file "artifacts/apk-confs/jwt_custom_header_resource_endpoint_conf.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      Then I set headers
        |Authorization|Bearer ${accessToken}|
      And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-custom-header-resource/3.14/employee/" with body ""
      And I eventually receive 401 response code, not accepting
        |429|
        |200|
      Then I set headers
        |testAuth|Bearer ${accessToken}|
      And I send "GET" request to "https://default.gw.wso2.com:9095/jwt-custom-header-resource/3.14/employee/" with body ""
      And I eventually receive 200 response code, not accepting
        |429|
        |401|

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
      | jwt-custom-header-resource-test      | 202                 |
