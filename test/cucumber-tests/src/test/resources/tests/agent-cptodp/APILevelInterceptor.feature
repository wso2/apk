Feature: API Policy Addition(Interceptor Service)
  Background:
    Given The system is ready
  Scenario: Create a REST API and add policy for request flow over API Level
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api_policy/request_interceptor.json"
    When the definition file "artifacts/definitions/cors-definition.json"
    And make the import API Creation request using OAS "File"
    Then the response status code should be 201
    And the response body should contain "IntercepterServiceAPI"
    And make the API Revision Deployment request
    Then the response status code should be 201
    Then I wait for 40 seconds
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
        | Authorization             | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/intercepter/1.0.0/get" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
  
  Scenario: Undeploying the created REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "IntercepterServiceAPI"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    And I eventually receive 404 response code, not accepting
      |200|

  Scenario: Create a REST API and add policy for response flow over API Level
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api_policy/response_interceptor.json"
    When the definition file "artifacts/definitions/cors-definition.json"
    And make the import API Creation request using OAS "File"
    Then the response status code should be 201
    And the response body should contain "IntercepterServiceAPI"
    And make the API Revision Deployment request
    Then the response status code should be 201
    Then I wait for 40 seconds
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
        | Authorization             | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"
  
  Scenario: Undeploying the created REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "ResponseIntercepterServiceAPI"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    And I eventually receive 404 response code, not accepting
      |200|

  Scenario: Create a REST API and add policy for request and response flow over API Level
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api_policy/request_and_response.json"
    When the definition file "artifacts/definitions/cors-definition.json"
    And make the import API Creation request using OAS "File"
    Then the response status code should be 201
    And the response body should contain "ReqandResIntercepterServiceAPI"
    And make the API Revision Deployment request
    Then the response status code should be 201
    Then I wait for 40 seconds
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
        | Authorization             | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/intercepter/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"
  
  Scenario: Undeploying the created REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "ReqandResIntercepterServiceAPI"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    And I eventually receive 404 response code, not accepting
      |200|

  Scenario: Create a REST API and add interceptor with parameter variation
    And I have a DCR application
    And I have a valid Publisher access token
    When I use the Payload file "artifacts/payloads/api_policy/request_interceptor_param_variation.json"
    When the definition file "artifacts/definitions/cors-definition.json"
    And make the import API Creation request using OAS "File"
    Then the response status code should be 201
    And the response body should contain "ParamVarIntercepterServiceAPI"
    And make the API Revision Deployment request
    Then the response status code should be 201
    Then I wait for 40 seconds
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
        | Authorization             | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/intercepter/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
  
  Scenario: Undeploying the created REST API
    And I have a DCR application
    And I have a valid Devportal access token
    Then I delete the application "SampleApp" from devportal
    Then the response status code should be 200
    And I have a valid Publisher access token
    Then I find the apiUUID of the API created with the name "ParamVarIntercepterServiceAPI"
    Then I undeploy the selected API
    Then the response status code should be 200
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    And I eventually receive 404 response code, not accepting
      |200|