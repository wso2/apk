Feature: Custom ratelimit
  Scenario: Testing custom ratelimit
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/custom_ratelimit_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
      | user_id       | bob                   |
      | org_id        | wso2                  |
    And I wait for next minute
# Request 1
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/3.14/employee" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
# Request 2
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 3
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 4
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 5 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Request 6 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Request 7 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Request 8 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Request 9 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Request 10 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
    Then I set headers
      | user_id       | dummy                 |  
      | org_id        | wso2                  |
# Starting from Request 5 the org_id descriptor should not be counted
# Request 5 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 6 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 7 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 8 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 9 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 10 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 11 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Request 12 - for org_id descriptor
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429
# Test org_id only 
    And I wait for next minute strictly
# Request 1
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/3.14/employee" with body ""
    Then the response status code should be 200
# Request 2
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 3
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 4
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 5
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 6 
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 7 
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 8
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 9
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 10
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 200
# Request 11 - should be limitted
    And I send "GET" request to "https://default.gw.wso2.com:9095/test-custom-ratelimit/employee" with body ""
    Then the response status code should be 429


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                | expectedStatusCode |
      | custom-ratelimit-api                 | 202                |
