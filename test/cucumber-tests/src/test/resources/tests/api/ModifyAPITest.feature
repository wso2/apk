Feature: API Deployment and invocation
  Scenario: Deploying an API and basic http method invocations
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/modifyAPI/originalAPI.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee1/" with body ""
    And the response status code should be 404
    When I use the APK Conf file "artifacts/apk-confs/modifyAPI/aNewResourceAddedAPI.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee1/" with body ""
    And the response status code should be 200
    When I use the APK Conf file "artifacts/apk-confs/modifyAPI/aResourceRemovedAPI.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And I eventually receive 404 response code, not accepting
      | 429 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee1/" with body ""
    And the response status code should be 200

  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                    | expectedStatusCode |
      | modify-api-test          | 202                |

