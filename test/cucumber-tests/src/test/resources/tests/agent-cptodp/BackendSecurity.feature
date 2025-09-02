# Feature: Backend Security for APIs
#   Background:
#     Given The system is ready
#   Scenario: Enable basic auth for endpoint and verfy the authorization for REST API
#     And I have a DCR application
#     And I have a valid Publisher access token
#     When I use the Payload file "artifacts/payloads/api_with_basic_auth.json"
#     And I use the OAS URL "https://petstore3.swagger.io/api/v3/openapi.json"
#     And make the import API Creation request using OAS "URL"
#     Then the response status code should be 201
#     And the response body should contain "SwaggerPetstore"
#     And make the API Revision Deployment request
#     Then the response status code should be 201
#     And make the Change Lifecycle request
#     Then the response status code should be 200
#     And I have a valid Devportal access token
#     And make the Application Creation request with the name "SampleApp"
#     Then the response status code should be 201
#     Then I wait for 40 seconds
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
#     And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/5" with body ""
#     Then the response status code should be 200
#     # And the response body should contain "\"Authorization\": \"Basic YWRtaW46YWRtaW4=\""

#   Scenario: Undeploying an already existing REST API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "SampleApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "SwaggerPetstore"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/basic-auth/3.14/employee/" with body ""
#     And I eventually receive 404 response code, not accepting
#       |200|

#   Scenario: Enable basic auth for endpoint and verfy the authorization for a GraphQL API
#     And I have a DCR application
#     And I have a valid Publisher access token
#     When the definition file "artifacts/definitions/schema_graphql.graphql"
#     When I use the Payload file "artifacts/payloads/gql_with_basic_auth.json"
#     Then I make the import GraphQLAPI Creation request
#     Then the response status code should be 201
#     And the response body should contain "StarWarsAPI"
#     Then I use the Payload file "artifacts/payloads/gql_with_basic_auth2.json"
#     And I update the API settings
#     Then the response status code should be 200
#     And the response body should contain "StarWarsAPI"
#     And make the API Revision Deployment request
#     Then the response status code should be 201
#     Then I wait for 40 seconds
#     And make the Change Lifecycle request
#     Then the response status code should be 200  
#     And I have a valid Devportal access token
#     And make the Application Creation request with the name "TestApp"
#     Then the response status code should be 201
#     And the response body should contain "TestApp"
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
#     And I send "POST" request to "https://default.gw.wso2.com:9095/gql_basic_auth/2.0.0/" with body "{\"query\":\"{ hero { id } }\"}"
#     Then the response status code should be 200
#     # And the response body should contain "\"Authorization\": \"Basic YWRtaW46YWRtaW4=\""
  
#   Scenario: Undeploying an already existing GraphQL API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "TestApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "StarWarsAPI"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "POST" request to "https://default.gw.wso2.com:9095/gql_basic_auth/2.0.0/" with body "{\"query\":\"{ hero { id } }\"}"
#     And I eventually receive 404 response code, not accepting
#       |200|