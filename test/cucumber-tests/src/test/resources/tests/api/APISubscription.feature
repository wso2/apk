# Feature: API Subscription Feature
#   Scenario: testing api subscriptions.
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/subscription-api.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|Bearer ${accessToken}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 403 response code, not accepting
#       |429|
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     Then the response status code should be 403
#     Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120005"
#     Then I set headers
#       |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120005-token}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     Then the response status code should be 403

#     Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120006"
#     Then I set headers
#       |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120006-token}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 403 response code, not accepting
#       |200|
#       |429|
#     Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120007"
#     Then I set headers
#       |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120007-token}|
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 403 response code, not accepting
#       |200|
#       |429|
#     Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120008"
#     Then I set headers
#       |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120008-token}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 403 response code, not accepting
#       |200|
#       |401|
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     And I eventually receive 403 response code, not accepting
#       |200|
#       |429|
#     Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120009"
#     Then I set headers
#       |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120009-token}|
#     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     Then the response status code should be 403
#     And I send "GET" request to "https://default.gw.wso2.com:9095/subscription-api/1.0.0/endpoint1" with body ""
#     Then the response status code should be 403
#   Scenario Outline: Undeploy API
#     Given The system is ready
#     And I have a valid subscription
#     When I undeploy the API whose ID is "<apiID>"
#     Then the response status code should be <expectedStatusCode>

#     Examples:
#       | apiID                 | expectedStatusCode  |
#       | subscription-api | 202 |
