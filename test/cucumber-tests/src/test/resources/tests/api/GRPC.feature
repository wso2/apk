Feature: Generating APK conf for gRPC API
    Scenario: Generating APK conf using a valid GRPC API definition
        Given The system is ready
        When I use the definition file "artifacts/definitions/student.proto" in resources
        And generate the APK conf file for a "gRPC" API
        Then the response status code should be 200

    Scenario: Deploying APK conf using a valid gRPC API definition
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/grpc/grpc.apk-conf"
        And the definition file "artifacts/definitions/student.proto"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | bearer ${accessToken} |
        And I make grpc request
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And the response body should contain "\"name\":\"string\""
        Then I set headers
            | Authorization | bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/grpc"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And the response body should contain "\"name\":\"string\""

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "grpc-basic-api"
        Then the response status code should be 202

