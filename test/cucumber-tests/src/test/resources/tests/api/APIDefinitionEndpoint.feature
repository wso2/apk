# Feature: API Definition Endpoint
#   Scenario: Testing default API definition endpoint
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/api_definition_default.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|Bearer ${accessToken}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-default/3.14/api-definition" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-default/api-definition" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|

#   Scenario: Testing custom API definition endpoint
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/api_definition_custom.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|Bearer ${accessToken}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-custom/3.14/docs" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-custom/docs" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-custom/api-definition" with body ""
#     And I eventually receive 404 response code, not accepting
#       |429|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-custom/3.14/api-definition" with body ""
#     And I eventually receive 404 response code, not accepting
#       |429|

#   Scenario: Testing a deleted production endpoint
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/api_definition_default_without_production.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     And I wait for 1 minute
#     Then I set headers
#       | Authorization | Bearer ${accessToken} |
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-default/3.14/api-definition" with body ""
#     Then the response status code should be 404
#     And I send "GET" request to "https://default.gw.wso2.com:9095/test-definition-default/api-definition" with body ""
#     Then the response status code should be 404

#   Scenario: Testing a deleted production endpoint
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/api_definition_default.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     When I use the APK Conf file "artifacts/apk-confs/api_definition_default_without_sandbox.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     And I wait for 1 minute
#     Then I set headers
#       | Authorization | Bearer ${accessToken} |
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/test-definition-default/3.14/api-definition" with body ""
#     Then the response status code should be 404
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/test-definition-default/api-definition" with body ""
#     Then the response status code should be 404


#   Scenario Outline: Undeploy API
#     Given The system is ready
#     And I have a valid subscription
#     When I undeploy the API whose ID is "<apiID>"
#     Then the response status code should be <expectedStatusCode>

#     Examples:
#       | apiID                 | expectedStatusCode  |
#       | custom-api-definition-endpoint-test   | 202                 |
#       | default-api-definition-endpoint-test   | 202                 |
