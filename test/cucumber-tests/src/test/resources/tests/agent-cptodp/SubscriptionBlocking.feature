# Feature: Blocking the subscription for a selected application
#   Background:
#     Given The system is ready
#   Scenario: Checking the subscription blocking for an REST API
#     And I have a DCR application
#     And I have a valid Publisher access token
#     When I use the Payload file "artifacts/payloads/api1.json"
#     And I use the OAS URL "https://petstore3.swagger.io/api/v3/openapi.json"
#     And make the import API Creation request using OAS "URL"
#     Then the response status code should be 201
#     And the response body should contain "SwaggerPetstore"
#     And make the API Revision Deployment request
#     Then the response status code should be 201
#     Then I wait for 40 seconds
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
#     Then I set headers
#         | Authorization             | Bearer ${accessToken} |
#     And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     Then I send the subcription blocking request
#     And the response status code should be 200
#     And the response body should contain "BLOCKED"
#     And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
#     And I eventually receive 403 response code, not accepting
#       |200|
#       |201|
#       |429|
#       |500|
  
#   Scenario: Undeploying the created REST API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "SampleApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "SwaggerPetstore"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
#     And I eventually receive 404 response code, not accepting
#       |200|