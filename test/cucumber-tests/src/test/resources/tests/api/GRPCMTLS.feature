Feature: Test mTLS between client and gateway with client certificate sent in header
    Scenario: Test API with mandatory mTLS and OAuth2 disabled
        Given The system is ready
        And I have a valid token with a client certificate "invalid-cert.txt"
        When I use the APK Conf file "artifacts/apk-confs/grpc/grpc_with_mtls_mandatory_oauth2_disabled.apk-conf"
        And the definition file "artifacts/definitions/student.proto"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
        And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
        And the gRPC response status code should be 0
        And the student response body should contain name: "Dineth" age: 10

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "grpc-mtls-mandatory-oauth2-disabled"
        Then the response status code should be 202

    Scenario: Test optional mTLS and optional OAuth2 with an invalid client certificate and invalid token in header
        Given The system is ready
        And I have a valid token with a client certificate "invalid-cert.txt"
        When I use the APK Conf file "artifacts/apk-confs/grpc/grpc_with_mtls_optional_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/student.proto"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
            | Authorization             | bearer {accessToken} |
        And I make grpc request GetStudent to "default.gw.wso2.com" with port 9095
        And the gRPC response status code should be 16

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "grpc-mtls-optional-oauth2-optional"
        Then the response status code should be 202
