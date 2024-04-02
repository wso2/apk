Feature: Testing Ratelimit feature
  Background:
    Given The system is ready

  # Scenario: Undeploy the created REST API
  #   And I have a DCR application
  #   And I have a valid Devportal access token
  #   Then I delete the application "SampleApp" from devportal
  #   Then the response status code should be 200
  #   Then I delete the application "ResourceLevelApp" from devportal
  #   Then the response status code should be 200
  #   And I have a valid Publisher access token
  #   Then I find the apiUUID of the API created with the name "SimpleRateLimitAPI"
  #   Then I undeploy the selected API
  #   Then I find the apiUUID of the API created with the name "SimpleRateLimitResourceLevelAPI"
  #   Then I undeploy the selected API
  #   Then the response status code should be 200

  Scenario: Testing API level rate limiiting for REST API
# #   For API Level
# #   PUT https://am.wso2.com/api/am/publisher/v4/apis/0a56786a-4ee6-4d6b-991c-1c9406b2b062
# #   Payload => api-level-ratelimit.json -> should get 200
    And I have a DCR application
    And I have a valid Adminportal access token
    # Then I set new API throttling policy allowing "2" requests per every "1" minute
    # Then the response status code should be 201
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/ratelimit_api.json"
    And I use the OAS URL "https://raw.githubusercontent.com/O-sura/JunkYard/main/employees_api.json"
    And make the import API Creation request using OAS "URL"
    Then the response status code should be 201
    And the response body should contain "SimpleRateLimitAPI"
    And make the API Revision Deployment request
    Then the response status code should be 201
    And make the Change Lifecycle request
    Then the response status code should be 200
    And I have a valid Devportal access token
    And make the Application Creation request with the name "SampleApp"
    Then the response status code should be 201
    And the response body should contain "SampleApp"
    And I have a KeyManager
    And make the Generate Keys request
    Then the response status code should be 200
    And the response body should contain "consumerKey"
    And the response body should contain "consumerSecret"
    And make the Subscription request
    Then the response status code should be 201
    And the response body should contain "Unlimited"
    And I get "production" oauth keys for application
    Then the response status code should be 200
    And make the Access Token Generation request for "production"
    Then the response status code should be 200
    And the response body should contain "accessToken"
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
    Then the response status code should be 429
    And I get "sandbox" oauth keys for application
    Then the response status code should be 200
    And make the Access Token Generation request for "sandbox"
    Then the response status code should be 200
    And the response body should contain "accessToken"
    And I send "GET" request to "https://sandbox.default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
    Then the response status code should be 429
    Then I wait for next minute strictly
    And I send "GET" request to "https://sandbox.default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
    Then the response status code should be 200


  # Scenario: Undeploy the created REST API
  #   And I have a DCR application
  #   And I have a valid Devportal access token
  #   Then I delete the application "SampleApp" from devportal
  #   Then the response status code should be 200
  #   And I have a valid Publisher access token
  #   Then I find the apiUUID of the API created with the name "SimpleRateLimitAPI"
  #   Then I undeploy the selected API
  #   Then the response status code should be 200
  #   And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
  #   Then the response status code should be 404
  #   And I send "GET" request to "https://sandbox.default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
  #   Then the response status code should be 404

  Scenario: Testing Resource level rate limiiting for REST API
#     #   For Resource Level
#     #   PUT https://am.wso2.com/api/am/publisher/v4/apis/41a6ede6-4943-4528-8abd-9f3b63eaf19b/swagger
#     #   form-data; name="apiDefinition"
#     #   Payload => resource_level_rl_def.json -> should get 200
#     #   should contain -> "EmployeeAPI"
#     #   should contain -> "x-throttling-tier": "10KPerMin"

    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/resource_level_rl.json"
    And I use the OAS URL "https://raw.githubusercontent.com/O-sura/JunkYard/main/employee_with_rl_r.json"
    And make the import API Creation request using OAS "URL"
    Then the response status code should be 201
    And the response body should contain "SimpleRateLimitResourceLevelAPI"
    And the response body should contain "\"throttlingPolicy\":\"TestRatelimit\""
    And make the API Revision Deployment request
    Then the response status code should be 201
    And make the Change Lifecycle request
    Then the response status code should be 200
    And I have a valid Devportal access token
    And make the Application Creation request with the name "ResourceLevelApp"
    Then the response status code should be 201
    And the response body should contain "ResourceLevelApp"
    And I have a KeyManager
    And make the Generate Keys request
    Then the response status code should be 200
    And the response body should contain "consumerKey"
    And the response body should contain "consumerSecret"
    And make the Subscription request
    Then the response status code should be 201
    And the response body should contain "Unlimited"
    And I get "production" oauth keys for application
    Then the response status code should be 200
    And make the Access Token Generation request for "production"
    Then the response status code should be 200
    And the response body should contain "accessToken"
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
      |401|
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
    Then the response status code should be 429
    And I send "GET" request to "https://sandbox.default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
    Then the response status code should be 429
    And I wait for next minute
    And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
    Then the response status code should be 200
    And I send "GET" request to "https://sandbox.default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
    Then the response status code should be 200


  Scenario: Undeploy the created REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    Then I delete the application "ResourceLevelApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "SimpleRateLimitAPI"
    Then I undeploy the selected API
    Then I find the apiUUID of the API created with the name "SimpleRateLimitResourceLevelAPI"
    Then I undeploy the selected API
    Then the response status code should be 200







# #   Scenario: Test simple rate limit api level for an REST API
# #     Given The system is ready
# #     And I have a valid subscription
# #     When I use the APK Conf file "artifacts/apk-confs/simple_rl_conf.yaml"
# #     And the definition file "artifacts/definitions/employees_api.json"
# #     And make the API deployment request
# #     Then the response status code should be 200
# #     Then I set headers
# #       |Authorization|bearer ${accessToken}|
# #     And I wait for next minute
# #     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
# #     And I eventually receive 200 response code, not accepting
# #       |429|
# #       |401|
# #     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
# #     Then the response status code should be 429
# #     And I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
# #     Then the response status code should be 429

# #   Scenario: Test simple rate limit api level for unsecured api
# #     Given The system is ready
# #     And I have a valid subscription
# #     When I use the APK Conf file "artifacts/apk-confs/simple_rl_jwt_disabled_conf.yaml"
# #     And the definition file "artifacts/definitions/employees_api.json"
# #     And make the API deployment request
# #     Then the response status code should be 200
# #     Then I set headers
# #       |Authorization|invalid|
# #     And I wait for next minute
# #     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-jwt-disabled/3.14/employee/" with body ""
# #     And I eventually receive 200 response code, not accepting
# #       |429|
# #       |401|
# #     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-jwt-disabled/3.14/employee/" with body ""
# #     Then the response status code should be 429

# #   Scenario: Test simple rate limit resource level
# #     Given The system is ready
# #     And I have a valid subscription
# #     When I use the APK Conf file "artifacts/apk-confs/simple_rl_resource_conf.yaml"
# #     And the definition file "artifacts/definitions/employees_api.json"
# #     And make the API deployment request
# #     Then the response status code should be 200
# #     Then I set headers
# #       |Authorization|bearer ${accessToken}|
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
# #     And I eventually receive 200 response code, not accepting
# #       |429|
# #       |401|
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
# #     Then the response status code should be 429
# #     And I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
# #     And I eventually receive 200 response code, not accepting
# #       |429|
# #       |401|
# #     And I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
# #     Then the response status code should be 429
# #     And I wait for next minute
# #     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
# #     Then the response status code should be 200
# #     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
# #     Then the response status code should be 200


# #   Scenario Outline: Undeploy API
# #     Given The system is ready
# #     And I have a valid subscription
# #     When I undeploy the API whose ID is "<apiID>"
# #     Then the response status code should be <expectedStatusCode>

# #     Examples:
# #       | apiID                       | expectedStatusCode |
# #       | simple-rl-test              | 202                |
# #       | simple-rl-r-test            | 202                |
# #       | simple-rl-jwt-disabled-test | 202                |