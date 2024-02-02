Feature: Generating APK conf for GraphQL API
    Scenario: Generating APK conf using a valid GraphQL API definition
        Given The system is ready
        When I use the definition file "artifacts/definitions/graphql_sample_api.graphql" in resources
        And generate the APK conf file for a "GRAPHQL" API
        Then the response status code should be 200
        And the response body should be "artifacts/apk-confs/graphql_conf.apk-conf" in resources

    Scenario: Generating APK conf using an invalid GraphQL API definition
        Given The system is ready
        When I use the definition file "artifacts/definitions/invalid_graphql_api.graphql" in resources
        And generate the APK conf file for a "GRAPHQL" API
        Then the response status code should be 400
