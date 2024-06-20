Feature: Test HTTPRoute Filter Request Mirror functionality
    Scenario: Test request mirror functionality
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/request-mirror-filter.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization             | bearer ${accessToken} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/request-mirror-filter/3.14/employee/" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |

    Scenario: Undeploy the API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "api-with-request-mirror-filter"
        Then the response status code should be 202
    
        