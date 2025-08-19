# Feature: Test simple rate limit feature
#   Scenario: Test simple rate limit api level
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/simple_rl_conf.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|Bearer ${accessToken}|
#     And I wait for next minute
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 429
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 429

#   Scenario: Test simple rate limit api level for unsecured api
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/simple_rl_jwt_disabled_conf.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|invalid|
#     And I wait for next minute
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-jwt-disabled/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-jwt-disabled/3.14/employee/" with body ""
#     Then the response status code should be 429

#   Scenario: Test simple rate limit resource level
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/simple_rl_resource_conf.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|Bearer ${accessToken}|
#     And I send "POST" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "POST" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 429
#     And I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 429
#     And I wait for next minute
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 200


#   Scenario Outline: Undeploy API
#     Given The system is ready
#     And I have a valid subscription
#     When I undeploy the API whose ID is "<apiID>"
#     Then the response status code should be <expectedStatusCode>

#     Examples:
#       | apiID                       | expectedStatusCode |
#       | simple-rl-test              | 202                |
#       | simple-rl-r-test            | 202                |
#       | simple-rl-jwt-disabled-test | 202                |
