Feature: API Deployment and invocation
  Scenario: Deploying an API and basic http method invocations
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/employees_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "POST" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/test/3.14/test/" with body ""
    And the response status code should be 404
    And I send "PUT" request to "https://default.gw.wso2.com:9095/test/3.14/employee/12" with body ""
    And the response status code should be 200
    And I send "DELETE" request to "https://default.gw.wso2.com:9095/test/3.14/employee/12" with body ""
    And the response status code should be 200
    Then I set headers
      | Authorization | Bearer invalidToken |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And the response status code should be 401
    And I send "POST" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And the response status code should be 401
    And I send "PUT" request to "https://default.gw.wso2.com:9095/test/3.14/employee/12" with body ""
    And the response status code should be 401
    And I send "DELETE" request to "https://default.gw.wso2.com:9095/test/3.14/employee/12" with body ""
    And the response status code should be 401

  Scenario: Deploying an API with new version
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/new_version_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    When I use the APK Conf file "artifacts/apk-confs/new_version_conf2.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I wait for next minute
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-version/1.0/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-version/2.0/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |

  Scenario: Deploying an API with default version
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/default_version_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-default/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "POST" request to "https://default.gw.wso2.com:9095/test-default/3.14/employee/" with body ""
    And the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/test-default/3.14/test/" with body ""
    And the response status code should be 404
    And I send "PUT" request to "https://default.gw.wso2.com:9095/test-default/3.14/employee/12" with body ""
    And the response status code should be 200
    And I send "DELETE" request to "https://default.gw.wso2.com:9095/test-default/3.14/employee/12" with body ""
    And the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-default/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "POST" request to "https://default.gw.wso2.com:9095/test-default/employee/" with body ""
    And the response status code should be 200
    And I send "POST" request to "https://default.gw.wso2.com:9095/test-default/test/" with body ""
    And the response status code should be 404
    And I send "PUT" request to "https://default.gw.wso2.com:9095/test-default/employee/12" with body ""
    And the response status code should be 200
    And I send "DELETE" request to "https://default.gw.wso2.com:9095/test-default/employee/12" with body ""
    And the response status code should be 200

  Scenario: Scope Validation
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/employees_scope_test_conf.yaml"
    And the definition file "artifacts/definitions/employees_scope_test_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithoutscope/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithscope1/" with body ""
    And the response status code should be 403
    Given I have a valid subscription with scopes
      | scope1 |
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithoutscope/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithscope1/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithscope2/" with body ""
    And I eventually receive 403 response code, not accepting
      | 200 |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithscopes/" with body ""
    And I eventually receive 403 response code, not accepting
      | 429 |
    Given I have a valid subscription with scopes
      | scope1 |
      | scope2 |
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-scope/1.0.0/employeewithscopes/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                    | expectedStatusCode |
      | f7996dce4ac15e2af0f8ee14546c4f72988eddae | 202                |
      | default-version-api-test                 | 202                |
      | emp-api-test-scope                       | 202                |
      | version-api-test                         | 202                |
      | version-api-test2                        | 202                |
