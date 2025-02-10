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
        Then I set headers
            | Authorization | Bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And the response body should contain "\"name\":\"string\""
        Then I set headers
            | Authorization | Bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And the response body should contain "\"name\":\"string\""

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "graphql-without-sub"
        Then the response status code should be 202

    Scenario: Deploying GraphQL API with scopes
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_scopes.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | Bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 403 response code, not accepting
            | 429 |
            | 500 |
        Given I have a valid subscription with scopes
            | wso2 |
        Then I set headers
            | Authorization | Bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "graphql-scopes"
        Then the response status code should be 202

    Scenario: Deploying a ratelimited GraphQL API
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_rl.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | Bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
        Then the response status code should be 429

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "graphql-rl"
        Then the response status code should be 202

    Scenario: Deploying multiple versions of a GraphQL API
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_3.0.0.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_4.0.0.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | Bearer ${accessToken} |
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.0.0" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And the response body should contain "\"name\":\"string\""
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/4.0.0" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |

    Scenario Outline: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "<apiID>"
        Then the response status code should be <expectedStatusCode>

        Examples:
            | apiID      | expectedStatusCode |
            | graphql-v3 | 202                |
            | graphql-v4 | 202                |

    # Scenario: Deploying APK conf using a valid GraphQL API definition with mTLS mandatory and valid certificate
    #     Given The system is ready
    #     And I have a valid token with a client certificate "config-map-1.txt"
    #     When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_with_mtls.apk-conf"
    #     And the definition file "artifacts/definitions/graphql_sample_api.graphql"
    #     And make the API deployment request
    #     Then the response status code should be 200
    #     Then I set headers
    #         | Authorization             | Bearer ${accessToken} |
    #         | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
    #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
    #     And I eventually receive 200 response code, not accepting
    #         | 429 |
    #         | 500 |
    #     And the response body should contain "\"name\":\"string\""

    # Scenario: Undeploy API
    #     Given The system is ready
    #     And I have a valid subscription
    #     When I undeploy the API whose ID is "graphql-mtls"
    #     Then the response status code should be 202

    # Scenario: Deploying APK conf using a valid GraphQL API definition with mTLS mandatory and no certificate
    #     Given The system is ready
    #     And I have a valid subscription
    #     When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_with_mtls.apk-conf"
    #     And the definition file "artifacts/definitions/graphql_sample_api.graphql"
    #     And make the API deployment request
    #     Then the response status code should be 200
    #     Then I set headers
    #         | Authorization | Bearer ${accessToken} |
    #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
    #     And I eventually receive 401 response code, not accepting
    #         | 200 |
    #         | 429 |
    #         | 500 |

    # Scenario: Undeploy API
    #     Given The system is ready
    #     And I have a valid subscription
    #     When I undeploy the API whose ID is "graphql-mtls"
    #     Then the response status code should be 202

    # Scenario: Deploying APK conf using a valid GraphQL API definition with OAuth2 mandatory mTLS optional
    #     Given The system is ready
    #     And I have a valid token with a client certificate "config-map-1.txt"
    #     When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_with_mtls_optional_oauth2_mandatory.apk-conf"
    #     And the definition file "artifacts/definitions/graphql_sample_api.graphql"
    #     And make the API deployment request
    #     Then the response status code should be 200
    #     Then I set headers
    #         | Authorization             | Bearer ${accessToken} |
    #         | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
    #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
    #     And I eventually receive 200 response code, not accepting
    #         | 429 |
    #         | 500 |
    #     And the response body should contain "\"name\":\"string\""
    #     Then I set headers
    #         | Authorization | Bearer ${accessToken} |
    #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
    #     And I eventually receive 200 response code, not accepting
    #         | 429 |
    #         | 500 |
    #     And the response body should contain "\"name\":\"string\""
    #     And I have a valid token with a client certificate "invalid-cert.txt"
    #     Then I set headers
    #         | Authorization             | Bearer ${accessToken} |
    #         | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
    #     And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
    #     And I eventually receive 401 response code, not accepting
    #         | 429 |
    #         | 500 |

    # Scenario: Undeploy API
    #     Given The system is ready
    #     And I have a valid subscription
    #     When I undeploy the API whose ID is "graphql-mtls-optional"
    #     Then the response status code should be 202

    Scenario: Deploying GraphQL API with OAuth2 disabled
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/graphql/graphql_with_disabled_auth.apk-conf"
        And the definition file "artifacts/definitions/graphql_sample_api.graphql"
        And make the API deployment request
        Then the response status code should be 200
        And I send "POST" request to "https://default.gw.wso2.com:9095/graphql/3.14" with body "{\"query\":\"{ allHumans { name } }\"}"
        And I eventually receive 200 response code, not accepting
            | 429 |
            | 500 |
        And the response body should contain "\"name\":\"string\""

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "graphql-auth-disabled"
        Then the response status code should be 202

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

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "graphql-with-sub"
        Then the response status code should be 202