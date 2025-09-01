Feature: Token revocation
 Scenario: Testing token revocation
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/employees_conf_new.yaml"
    And the definition file "artifacts/definitions/employees_api_new.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
        | Authorization | Bearer ${accessToken}         |
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
        | 429 |
    # Revoke the token
    Then I set headers
        |stsAuthKey|2jrmypak391zsqz974ugdddebf812ofx1b9t1oq27530ir02tc815eemrx435qvcp41ucgy7v5uuawzi4qcmjrx0k1zgox2s28cr|
    And I send "POST" request to "https://api.am.wso2.com:9095/api/notification/1.0.0/notify?type=TOKEN_REVOCATION" with body "{\"token\": \"${accessToken}\"}"
    And the response status code should be 200
    And I wait for 5 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/test/3.14/employee/" with body ""
    And the response status code should be 401


#  Scenario Outline: Undeploy API
#    Given The system is ready
#    And I have a valid subscription
#    When I undeploy the API whose ID is "<apiID>"
#    Then the response status code should be <expectedStatusCode>

#    Examples:
#      | apiID                 | expectedStatusCode  |
#      | jwt-basic-test          | 202                 |
