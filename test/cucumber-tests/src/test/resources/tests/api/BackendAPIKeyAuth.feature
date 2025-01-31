Feature: Backend API Key auth
  Scenario: Testing API level Endpoint backend api key auth header
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/backend_apikey_auth_conf.yaml"
    And the definition file "artifacts/definitions/backend_apikey_auth_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-api-key-security/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "\"Api-Key\": \"sampath\""


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | backend-api-key-test          | 202                 |

  Scenario: Testing undeployed API
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-api-key-security/3.14/employee/" with body ""
    And I eventually receive 404 response code, not accepting
      | 200 |

