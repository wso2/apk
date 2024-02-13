Feature: Basic auth
  Scenario: Testing API level and resource level basic auth header
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/basic_auth_conf.yaml"
    And the definition file "artifacts/definitions/basic_auth_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/basic-auth/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "\"Authorization\": \"Basic YWRtaW46YWRtaW4=\""
    And I send "GET" request to "https://default.gw.wso2.com:9095/basic-auth/3.14/get" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "\"Authorization\": \"Basic ZHNmZHNmc2Rmc2RmOmFkbWlu\""
    And I send "POST" request to "https://default.gw.wso2.com:9095/basic-auth/3.14/post" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "\"Authorization\": \"Basic YWRtaW46YWRtaW4=\""


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | basic-auth-api-test          | 202                 |

  Scenario: Testing undeployed API
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      | Authorization | bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/basic-auth/3.14/employee/" with body ""
    And I eventually receive 404 response code, not accepting
      | 429 |

  Scenario Outline: Undeploy API finally
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID               | expectedStatusCode |
      | basic-auth-api-test | 202                |