# Feature: Test mTLS between client and gateway with client certificate sent in header
#     Scenario: Test API with mandatory mTLS and OAuth2 disabled
#         Given The system is ready
#         And I have a valid token with a client certificate "config-map-1.txt"
#         When I use the APK Conf file "artifacts/apk-confs/grpc/grpc_with_mtls_mandatory_oauth2_disabled.apk-conf"
#         And the definition file "artifacts/definitions/student.proto"
#         And make the API deployment request
#         Then the response status code should be 200
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 0
#         And the student response body should contain name: "Student" age: 10
#         And I have a valid token with a client certificate "invalid-cert.txt"
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 16

#     Scenario: Undeploy API
#         Given The system is ready
#         And I have a valid subscription
#         When I undeploy the API whose ID is "grpc-mtls-mandatory-oauth2-disabled"
#         Then the response status code should be 202
    
#     Scenario: Test mandatory mTLS and mandatory OAuth2
#         Given The system is ready
#         And I have a valid token with a client certificate "config-map-1.txt"
#         When I use the APK Conf file "artifacts/apk-confs/grpc/grpc_with_mtls_mandatory_oauth2_mandatory.apk-conf"
#         And the definition file "artifacts/definitions/student.proto"
#         And make the API deployment request
#         Then the response status code should be 200
#         Then I set headers
#             | Authorization | Bearer ${accessToken} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 16
#         Then I remove header "Authorization"
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 16
#         Then I remove header "X-WSO2-CLIENT-CERTIFICATE"
#         And I have a valid token with a client certificate "invalid-cert.txt"
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#             | Authorization | Bearer ${accessToken} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 16

#     Scenario: Undeploy API
#         Given The system is ready
#         And I have a valid subscription
#         When I undeploy the API whose ID is "grpc-mtls-mandatory-oauth2-mandatory"
#         Then the response status code should be 202

#     Scenario: Test optional mTLS and optional OAuth2
#         Given The system is ready
#         And I have a valid token with a client certificate "config-map-1.txt"
#         When I use the APK Conf file "artifacts/apk-confs/grpc/grpc_with_mtls_optional_oauth2_optional.apk-conf"
#         And the definition file "artifacts/definitions/student.proto"
#         And make the API deployment request
#         Then the response status code should be 200
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#             | Authorization | Bearer ${accessToken} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 0
#         And the student response body should contain name: "Student" age: 10
#         Then I set headers
#             | Authorization | Bearer ${accessToken} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 0
#         And the student response body should contain name: "Student" age: 10
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 0
#         And the student response body should contain name: "Student" age: 10
#         And I have a valid token with a client certificate "invalid-cert.txt"
#         Then I set headers
#             | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
#         And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
#         And the gRPC response status code should be 16

#     Scenario: Undeploy API
#         Given The system is ready
#         And I have a valid subscription
#         When I undeploy the API whose ID is "grpc-mtls-optional-oauth2-optional"
#         Then the response status code should be 202