# Feature: Invoking APIs with scopes enabled
#   Background:
#     Given The system is ready
#   Scenario: Create API and enable scopes for resources and invoke them using keys with and without proper scopes
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
#     Then I make Access Token Generation request without scopes
#     Then the response status code should be 200
#     And the response body should contain "accessToken"
#     Then I set headers
#         | Authorization             | Bearer ${accessToken} |
#     And I send "GET" request to "https://default.gw.wso2.com:9095/petstore/1.0.0/pet/4" with body ""
#     Then the response status code should be 403
  
#   Scenario: Undeploy the created REST API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "SampleApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "SwaggerPetstore"
#     Then I undeploy the selected API
#     Then the response status code should be 200
  
#   Scenario: Deploying a GraphQL API with scopes and invoking it with and without scopes
#     And I have a DCR application
#     And I have a valid Publisher access token
#     When the definition file "artifacts/definitions/schema_graphql.graphql"
#     When I use the Payload file "artifacts/payloads/gql_with_scopes.json"
#     Then I make the import GraphQLAPI Creation request
#     Then the response status code should be 201
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
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 200
#     And I eventually receive 200 response code, not accepting
#       | 404 |
#       | 401 |
#     Then I make Access Token Generation request without scopes
#     Then the response status code should be 200
#     And the response body should contain "accessToken"
#     Then I set headers
#         | Authorization             | Bearer ${accessToken} |
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 403

#   Scenario: Undeploying an already existing GraphQL API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "TestApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "StarWarsAPI"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     And I eventually receive 404 response code, not accepting
#       |200|