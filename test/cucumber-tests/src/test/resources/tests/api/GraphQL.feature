Feature: Generating APK conf for GraphQL API
    Scenario: Generating APK conf using a valid GraphQL API definition
        Given The system is ready
        When I use the definition file "artifacts/definitions/graphql_sample_api.graphql" in resources
        And generate the APK conf file for a "GRAPHQL" API
        Then the response status code should be 200

    Scenario: Deploying APK conf using a valid GraphQL API definition without a subscription resource
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_conf_without_sub.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Deploying APK conf using a valid GraphQL API definition containing a subscription resource
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_conf_with_sub.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Generating APK conf using an invalid GraphQL API definition
        Given The system is ready
        When I use the definition file "artifacts/definitions/invalid_graphql_api.graphql" in resources
        And generate the APK conf file for a "GRAPHQL" API
        Then the response status code should be 400
