# Feature: Test mTLS between client and gateway with client certificate sent in header
#     Background:
#         Given The system is ready
#     Scenario: Deploying GraphQL API with mTLS mandatory and valid certificate
#         And I have a DCR application
#         And I have a valid Publisher access token
#         When the definition file "artifacts/definitions/schema_graphql.graphql"
#         When I use the Payload file "artifacts/payloads/mtls/graphql_with_mtls.json"
#         Then I make the import GraphQLAPI Creation request
#         Then the response status code should be 201
#         And the response body should contain "GraphQLAPImTLS"
#         And I have a client certificate "config-map-1.crt"
#         Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
#         Then the response status code should be 201
#         And make the API Revision Deployment request
#         Then the response status code should be 201
#         Then I wait for 40 seconds
#         And make the Change Lifecycle request
#         Then the response status code should be 200  
#         And I have a valid Devportal access token
#         And make the Application Creation request with the name "TestApp"
#         Then the response status code should be 201
#         And the response body should contain "TestApp"
#         And I have a KeyManager
#         And make the Generate Keys request
#         Then the response status code should be 200
#         And the response body should contain "consumerKey"
#         And the response body should contain "consumerSecret"
#         And make the Subscription request
#         Then the response status code should be 201
#         And the response body should contain "Unlimited"
#         And I get "production" oauth keys for application
#         Then the response status code should be 200
#         And make the Access Token Generation request for "production"
#         Then the response status code should be 200
#         And the response body should contain "accessToken"
#         Then I set headers
#             | Authorization             | Bearer ${accessToken} |
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
#         And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
#         And I eventually receive 200 response code, not accepting
#             | 429 |
#             | 500 |
#         And the response body should contain "\"name\":\"string\""

#     Scenario: Undeploy API
#         And I have a DCR application
#         And I have a valid Devportal access token
#         Then I delete the application "TestApp" from devportal
#         Then the response status code should be 200
#         And I have a valid Publisher access token
#         Then I find the apiUUID of the API created with the name "GraphQLAPImTLS"
#         Then I undeploy the selected API
#         Then the response status code should be 200
#         And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
#         And I eventually receive 404 response code, not accepting
#         |200|


#     Scenario: Deploying GraphQL API with mTLS mandatory and no certificate
#         And I have a DCR application
#         And I have a valid Publisher access token
#         When the definition file "artifacts/definitions/schema_graphql.graphql"
#         When I use the Payload file "artifacts/payloads/mtls/graphql_with_mtls.json"
#         Then I make the import GraphQLAPI Creation request
#         Then the response status code should be 201
#         And the response body should contain "GraphQLAPImTLS"
#         And I have a client certificate "config-map-1.crt"
#         Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
#         Then the response status code should be 201
#         And make the API Revision Deployment request
#         Then the response status code should be 201
#         Then I wait for 40 seconds
#         And make the Change Lifecycle request
#         Then the response status code should be 200  
#         And I have a valid Devportal access token
#         And make the Application Creation request with the name "TestApp"
#         Then the response status code should be 201
#         And the response body should contain "TestApp"
#         And I have a KeyManager
#         And make the Generate Keys request
#         Then the response status code should be 200
#         And the response body should contain "consumerKey"
#         And the response body should contain "consumerSecret"
#         And make the Subscription request
#         Then the response status code should be 201
#         And the response body should contain "Unlimited"
#         And I get "production" oauth keys for application
#         Then the response status code should be 200
#         And make the Access Token Generation request for "production"
#         Then the response status code should be 200
#         And the response body should contain "accessToken"
#         Then I set headers
#             | Authorization             | Bearer ${accessToken} |
#         And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
#         And I eventually receive 401 response code, not accepting
#             | 200 |
#             | 429 |
#             | 500 |

#     Scenario: Undeploy API
#         And I have a DCR application
#         And I have a valid Devportal access token
#         Then I delete the application "TestApp" from devportal
#         Then the response status code should be 200
#         And I have a valid Publisher access token
#         Then I find the apiUUID of the API created with the name "GraphQLAPImTLS"
#         Then I undeploy the selected API
#         Then the response status code should be 200
#         And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
#         And I eventually receive 404 response code, not accepting
#         |200|

#     Scenario: Deploying GraphQL API with OAuth2 mandatory mTLS optional
#         And I have a DCR application
#         And I have a valid Publisher access token
#         When the definition file "artifacts/definitions/schema_graphql.graphql"
#         When I use the Payload file "artifacts/payloads/mtls/graphql_with_mtls_optional_oauth2_mandatory.json"
#         Then I make the import GraphQLAPI Creation request
#         Then the response status code should be 201
#         And the response body should contain "GraphQLAPImTLS"
#         And I have a client certificate "config-map-1.crt"
#         Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
#         Then the response status code should be 201
#         And make the API Revision Deployment request
#         Then the response status code should be 201
#         Then I wait for 40 seconds
#         And make the Change Lifecycle request
#         Then the response status code should be 200  
#         And I have a valid Devportal access token
#         And make the Application Creation request with the name "TestApp"
#         Then the response status code should be 201
#         And the response body should contain "TestApp"
#         And I have a KeyManager
#         And make the Generate Keys request
#         Then the response status code should be 200
#         And the response body should contain "consumerKey"
#         And the response body should contain "consumerSecret"
#         And make the Subscription request
#         Then the response status code should be 201
#         And the response body should contain "Unlimited"
#         And I get "production" oauth keys for application
#         Then the response status code should be 200
#         And make the Access Token Generation request for "production"
#         Then the response status code should be 200
#         And the response body should contain "accessToken"
#         Then I set headers
#             | Authorization             | Bearer ${accessToken} |        
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
#         And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
#         And I eventually receive 200 response code, not accepting
#             | 429 |
#             | 500 |
#         And the response body should contain "\"name\":\"string\""
#         Then I remove the header "X-WSO2-CLIENT-CERTIFICATE"
#         And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
#         And I eventually receive 200 response code, not accepting
#             | 429 |
#             | 500 |
#         And the response body should contain "\"name\":\"string\""
#         And I have a client certificate "invalid-cert.crt"
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
#         And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
#         And I eventually receive 401 response code, not accepting
#             | 429 |
#             | 500 |

#     Scenario: Undeploy API
#         And I have a DCR application
#         And I have a valid Devportal access token
#         Then I delete the application "TestApp" from devportal
#         Then the response status code should be 200
#         And I have a valid Publisher access token
#         Then I find the apiUUID of the API created with the name "GraphQLAPImTLS"
#         Then I undeploy the selected API
#         Then the response status code should be 200
#         And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
#         And I eventually receive 404 response code, not accepting
#         |200|

