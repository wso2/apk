# Feature: Testing Ratelimit feature
#   Background:
#     Given The system is ready
#   Scenario: Testing API level rate limiiting for REST API
#     And I have a DCR application
#     And I have a valid Adminportal access token
#     Then I set new API throttling policy allowing "2" requests per every "1" minute
#     Then the response status code should be 201
#     And I have a valid Publisher access token
#     When I use the Payload file "artifacts/payloads/ratelimit_api.json"
#     When the definition file "artifacts/definitions/employees_api.json"
#     And make the import API Creation request using OAS "File"
#     Then the response status code should be 201
#     And the response body should contain "SimpleRateLimitAPI"
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
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 429
#     Then I wait for next minute strictly
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 200


#   Scenario: Undeploy the created REST API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "SampleApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "SimpleRateLimitAPI"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 404
#     And I send "GET" request to "https://sandbox.default.gw.wso2.com:9095/simple-rl/3.14/employee/" with body ""
#     Then the response status code should be 404
#     And I have a valid Adminportal access token
#     Then I delete the created API throttling policy

#   Scenario: Testing Resource level rate limiiting for REST API
#     And I have a DCR application
#     And I have a valid Adminportal access token
#     Then I set new API throttling policy allowing "2" requests per every "1" minute
#     Then the response status code should be 201
#     And I have a valid Publisher access token
#     When I use the Payload file "artifacts/payloads/resource_level_rl.json"
#     When the definition file "artifacts/definitions/employee_with_rl_r.json"
#     And make the import API Creation request using OAS "File"
#     Then the response status code should be 201
#     And the response body should contain "SimpleRateLimitResourceLevelAPI"
#     And the response body should contain "\"throttlingPolicy\":\"TestRatelimit\""
#     And make the API Revision Deployment request
#     Then the response status code should be 201
#     Then I wait for 40 seconds
#     And make the Change Lifecycle request
#     Then the response status code should be 200
#     And I have a valid Devportal access token
#     And make the Application Creation request with the name "ResourceLevelApp"
#     Then the response status code should be 201
#     And the response body should contain "ResourceLevelApp"
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
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 429
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 429
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/withoutrl/" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/withoutrl/" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/withoutrl/" with body ""
#     Then the response status code should be 200
#     Then I wait for next minute strictly
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 200
#     And I send "GET" request to "https://default.gw.wso2.com:9095/simple-rl-r/3.14/employee/" with body ""
#     Then the response status code should be 200


#   Scenario: Undeploy the created REST API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "ResourceLevelApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "SimpleRateLimitResourceLevelAPI"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I have a valid Adminportal access token
#     Then I delete the created API throttling policy

#   Scenario: Testing API level rate limiiting for GraphQL API
#     And I have a DCR application
#     And I have a valid Adminportal access token
#     Then I set new API throttling policy allowing "2" requests per every "1" minute
#     Then the response status code should be 201
#     And I have a valid Publisher access token
#     When the definition file "artifacts/definitions/schema_graphql.graphql"
#     Then I use the Payload file "artifacts/payloads/gqlPayload.json"
#     Then I make the import GraphQLAPI Creation request
#     Then the response status code should be 201
#     And the response body should contain "StarwarsAPI"
#     Then I use the Payload file "artifacts/payloads/gql_api_level_rl.json"
#     And I update the API settings
#     Then the response status code should be 200
#     And the response body should contain "StarwarsAPI"
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
#     And I eventually receive 200 response code, not accepting
#       |429|
#       |401|
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 200
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 429
#     Then I wait for next minute strictly
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 200


#   Scenario: Undeploy the created GraphQL API
#     And I have a DCR application
#     And I have a valid Devportal access token
#     Then I delete the application "TestApp" from devportal
#     Then the response status code should be 200
#     And I have a valid Publisher access token
#     Then I find the apiUUID of the API created with the name "StarwarsAPI"
#     Then I undeploy the selected API
#     Then the response status code should be 200
#     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 404
#     And I send "POST" request to "https://sandbox.default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
#     Then the response status code should be 404
#     And I have a valid Adminportal access token
#     Then I delete the created API throttling policy
#     Then the response status code should be 200

#  # NOTE: In the current implementation, APK only supports API level ratelimitting hence this test case
#  #       commented out. Uncomment this after the proper implementation.
# #   Scenario: Testing Resource level rate limiiting for REST API
# #     And I have a DCR application
# #     And I have a valid Adminportal access token
# #     Then I set new API throttling policy allowing "2" requests per every "1" minute
# #     Then the response status code should be 201
# #     And I have a valid Publisher access token
# #     When the definition file "artifacts/definitions/schema_graphql.graphql"
# #     Then I use the Payload file "artifacts/payloads/gqlPayload.json"
# #     Then I make the import GraphQLAPI Creation request
# #     Then the response status code should be 201
# #     And the response body should contain "StarwarsAPI"
# #     Then I use the Payload file "artifacts/payloads/gql_resource_level_rl.json"
# #     And I update the GQL API settings
# #     Then the response status code should be 200
# #     And the response body should contain "StarwarsAPI"
# #     And make the API Revision Deployment request
# #     Then the response status code should be 201
# #     And make the Change Lifecycle request
# #     Then the response status code should be 200  
# #     And I have a valid Devportal access token
# #     And make the Application Creation request with the name "TestApp"
# #     Then the response status code should be 201
# #     And the response body should contain "TestApp"
# #     And I have a KeyManager
# #     And make the Generate Keys request
# #     Then the response status code should be 200
# #     And the response body should contain "consumerKey"
# #     And the response body should contain "consumerSecret"
# #     And make the Subscription request
# #     Then the response status code should be 201
# #     And the response body should contain "Unlimited"
# #     And I get "production" oauth keys for application
# #     Then the response status code should be 200
# #     And make the Access Token Generation request for "production"
# #     Then the response status code should be 200
# #     And the response body should contain "accessToken"
# #     Then I set headers
# #         | Authorization             | Bearer ${accessToken} |
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
# #     And I eventually receive 200 response code, not accepting
# #       |429|
# #       |401|
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
# #     Then the response status code should be 200
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
# #     Then the response status code should be 429
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
# #     Then the response status code should be 429
# #     #From here onwards, it should query an endpoint without rate limit
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ human(id:1000){ id name }}\"}";
# #     Then the response status code should be 200
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ human(id:1000){ id name }}\"}";
# #     Then the response status code should be 200
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ human(id:1000){ id name }}\"}";
# #     Then the response status code should be 200
# #     Then I wait for next minute strictly
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
# #     Then the response status code should be 200
# #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ hero { name } }\"}"
# #     Then the response status code should be 200


# #   Scenario: Undeploy the created GraphQL API
# #     And I have a DCR application
# #     And I have a valid Devportal access token
# #     Then I delete the application "TestApp" from devportal
# #     Then the response status code should be 200
# #     And I have a valid Publisher access token
# #     Then I find the apiUUID of the API created with the name "StarwarsAPI"
# #     Then I undeploy the selected API
# #     Then the response status code should be 200
# #     And I have a valid Adminportal access token
# #     Then I delete the created API throttling policy
# #     Then the response status code should be 200

