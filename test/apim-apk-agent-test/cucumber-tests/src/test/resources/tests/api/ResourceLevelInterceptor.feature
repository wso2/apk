# Feature: API Policy Addition(Interceptor Service)
#   Background:
#     Given The system is ready
#   Scenario: Create a REST API and add policy for request flow over Resource Level
#     And I have a DCR application
#     And I have a valid Publisher access token
#     When I use the Payload file "artifacts/payloads/api1.json"
#     When the definition file "artifacts/definitions/cors-definition.json"
#     And make the import API Creation request using OAS "File"
#     Then the response status code should be 201
#     And the response body should contain "SwaggerPetstore"
#     And make the API Revision Deployment request
#     Then the response status code should be 201
#     And make the Change Lifecycle request
#     Then the response status code should be 200
#     And I have a valid Devportal access token
#     And make the Application Creation request with the name "SampleApp"
#     Then the response status code should be 201
#     And the response body should contain "SampleApp"
#     And I have a KeyManager
#     And make the Generate Keys request
#     Then the response status code should be 200
#     And the response body should contain "consumerKey"
#     And the response body should contain "consumerSecret"
#     And make the Subscription request
#     Then the response status code should be 201
#     And the response body should contain "Unlimited"
#     And I get "production" oauth keys for application
#     Then the response status code should be 200
#     And make the Access Token Generation request for "production"
#     Then the response status code should be 200
#     And the response body should contain "accessToken"
#     And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/employee" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
  
#   Scenario: Undeploying an already existing REST API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "SampleApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "SwaggerPetstore"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/employee" with body ""
#     And I eventually receive 404 response code, not accepting
#       |200|


# Feature: API Deployment with Resource Interceptor
#   Scenario: Deploying an API
#     Given The system is ready
#     And I have a valid subscription
#     When I use the APK Conf file "artifacts/apk-confs/interceptors/resourceLevelInterptor.apk-conf"
#     And the definition file "artifacts/definitions/cors_api.yaml"
#     And make the API deployment request
#     Then the response status code should be 200
#     And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
#     And I wait for 1 minute
#     Then I set headers
#       |Authorization|bearer ${accessToken}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
#     And the response body should not contain "\"Interceptor-Header\""
#     Then the response status code should be 200
#     Then the response headers not contains key "interceptor-response-header"
#     Then I set headers
#       |Authorization|bearer ${accessToken}|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/headers" with body ""
#     And the response body should contain
#       |"Interceptor-Header": "Interceptor-header-value"|
#       |"Interceptor-Header-Apigroup": "Gold"|
#       |"Interceptor-Header-Apitier": "Unlimited"|
#     Then the response status code should be 200
#     Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"
#   Scenario Outline: Undeploy an API
#     Given The system is ready
#     And I have a valid subscription
#     When I undeploy the API whose ID is "<apiID>"
#     Then the response status code should be <expectedStatusCode>

#     Examples:
#       | apiID | expectedStatusCode |
#       | 547961eeaafed989119c45ffc13f8b87bfda821d  | 202         |