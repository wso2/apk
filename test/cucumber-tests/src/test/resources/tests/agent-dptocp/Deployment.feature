Feature: API Deploying in DP to CP Flow
  Scenario: Deploy a REST API to APK DP and invoke it using APIM Devportal
    Given The system is ready
    And I have a valid subscription
    When I use the apk conf file "artifacts/apk-confs/endpoint_conf.yaml" in resources
    And I use the definition file "artifacts/definitions/employees_api.json"
    Then I generate and apply the K8Artifacts belongs to that API
    Then I wait for 40 seconds  
    And I have a DCR application
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "APIResourceEndpoint"
    Then the response status code should be 200
    And the response body should contain "APIResourceEndpoint"
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
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://carbon.super.gw.wso2.com:9095/endpoint/3.14/employee" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://backend/anything/employee"
    And I send "POST" request to "https://carbon.super.gw.wso2.com:9095/endpoint/3.14/employee" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "https://backend/anything/test/employee"


  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "APIResourceEndpoint"
    Then I undeploy the selected API
    Then the response status code should be 200

    Examples:
      | apiID                 | expectedStatusCode  |
      | endpoint-test          | 202                 |
