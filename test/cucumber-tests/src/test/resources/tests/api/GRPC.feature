Feature: Generating APK conf for gRPC API
    Scenario: Generating APK conf using a valid GRPC API definition
        Given The system is ready
        When I use the definition file "artifacts/definitions/student.proto" in resources
        And generate the APK conf file for a "gRPC" API
        Then the response status code should be 200
