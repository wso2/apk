# Feature: Test HTTPRoute Filter Header Modifier functionality
#     Scenario: Test request and response header modification functionality
#         Given The system is ready
#         And I have a valid subscription
#         When I use the APK Conf file "artifacts/apk-confs/httproute-filters/header-modifier-filter.apk-conf"
#         And the definition file "artifacts/definitions/employees_api.json"
#         And make the API deployment request
#         Then the response status code should be 200
#         Then I set headers
#             | Authorization             | Bearer ${accessToken} |
#         And I send "GET" request to "https://default.gw.wso2.com:9095/header-modifier-filters/3.14/employee/" with body ""
#         And I eventually receive 200 response code, not accepting
#             | 401 |
#         And the response body should contain "\"Test-Request-Header\": \"Test-Value\""
#         And the response body should contain "\"Set-Request-Header\": \"Test-Value\""
#         And the response body should not contain "\"Authorization\""
#         Then the response headers contains key "Set-Response-Header" and value "Test-Value"
#         Then the response headers contains key "Test-Response-Header" and value "Test-Value"
#         And the response headers should not contain
#         | content-type |

#     Scenario: Undeploy the API
#         Given The system is ready
#         And I have a valid subscription
#         When I undeploy the API whose ID is "api-with-header-modifier-filters"
#         Then the response status code should be 202

#     Scenario: Test request and response header modification functionality
#         Given The system is ready
#         And I have a valid subscription
#         When I use the APK Conf file "artifacts/apk-confs/httproute-filters/api-level-header.apk-conf"
#         And the definition file "artifacts/definitions/employees_api.json"
#         And make the API deployment request
#         Then the response status code should be 200
#         Then I set headers
#             | Authorization             | Bearer ${accessToken} |
#         And I send "GET" request to "https://default.gw.wso2.com:9095/header-modifier-filters/3.14/employee/" with body ""
#         And I eventually receive 200 response code, not accepting
#             | 401 |
#         And the response body should contain "\"Test-Request-Header\": \"Test-Value\""
#         And the response body should contain "\"Set-Request-Header\": \"Test-Value\""
#         And the response body should not contain "\"Authorization\""
#         Then the response headers contains key "Set-Response-Header" and value "Test-Value"
#         Then the response headers contains key "Test-Response-Header" and value "Test-Value"
#         And the response headers should not contain
#         | content-type |
#         And I send "POST" request to "https://default.gw.wso2.com:9095/header-modifier-filters/3.14/employee/" with body ""
#         And I eventually receive 200 response code, not accepting
#             | 401 |
#         And the response body should contain "\"Test-Request-Header\": \"Test-Value\""
#         And the response body should contain "\"Set-Request-Header\": \"Test-Value\""
#         And the response body should not contain "\"Authorization\""
#         Then the response headers contains key "Set-Response-Header" and value "Test-Value"
#         Then the response headers contains key "Test-Response-Header" and value "Test-Value"
#         And the response headers should not contain
#         | content-type |
#         And I send "PUT" request to "https://default.gw.wso2.com:9095/header-modifier-filters/3.14/employee/1" with body ""
#         And I eventually receive 200 response code, not accepting
#             | 401 |
#         And the response body should contain "\"Test-Request-Header\": \"Test-Value\""
#         And the response body should contain "\"Set-Request-Header\": \"Test-Value\""
#         And the response body should not contain "\"Authorization\""
#         Then the response headers contains key "Set-Response-Header" and value "Test-Value"
#         Then the response headers contains key "Test-Response-Header" and value "Test-Value"
#         And the response headers should not contain
#         | content-type |
#         And I send "DELETE" request to "https://default.gw.wso2.com:9095/header-modifier-filters/3.14/employee/1" with body ""
#         And I eventually receive 200 response code, not accepting
#             | 401 |
#         And the response body should contain "\"Test-Request-Header\": \"Test-Value\""
#         And the response body should contain "\"Set-Request-Header\": \"Test-Value\""
#         And the response body should not contain "\"Authorization\""
#         Then the response headers contains key "Set-Response-Header" and value "Test-Value"
#         Then the response headers contains key "Test-Response-Header" and value "Test-Value"
#         And the response headers should not contain
#         | content-type |

#     Scenario: Undeploy the API
#         Given The system is ready
#         And I have a valid subscription
#         When I undeploy the API whose ID is "api-with-header-modifier-filters"
#         Then the response status code should be 202
    
        