Feature: Test mTLS between client and gateway with client certificate sent in header with OAuth2 optional

    # mTLS mandatory OAuth2 optional
    Scenario: Test mandatory mTLS and optional OAuth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_mandatory_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization             | bearer ${accessToken} |
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
        And I eventually receive 200 response code, not accepting
            | 401 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-optional"
        Then the response status code should be 202

    Scenario: Test mandatory mTLS and optional OAuth2 with an invalid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "invalid-cert.txt"
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_mandatory_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization             | bearer ${accessToken} |
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |
        Then I set headers
            | Authorization | bearer ${accessToken} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-optional"
        Then the response status code should be 202

    Scenario: Test mandatory mTLS and optional OAuth2 without a client certificate in header
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_mandatory_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | bearer ${accessToken} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-optional"
        Then the response status code should be 202

    # mTLS optional OAuth2 optional
    Scenario: Test optional mTLS and optional OAuth2 with a valid token and then a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_optional_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | bearer ${accessToken} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional"
        Then the response status code should be 202

    Scenario: Test optional mTLS and optional OAuth2 with an invalid client certificate and invalid token in header
        Given The system is ready
        And I have a valid token with a client certificate "invalid-cert.txt"
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_optional_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
            | Authorization             | bearer {accessToken} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional"
        Then the response status code should be 202

    Scenario: Test optional mTLS and optional OAuth2 with an invalid client certificate and valid token in header
        Given The system is ready
        And I have a valid token with a client certificate "invalid-cert.txt"
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_optional_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
            | Authorization             | bearer ${accessToken} |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional"
        Then the response status code should be 202

    Scenario: Test optional mTLS and optional OAuth2 with an invalid token in header
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_optional_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        Then I set headers
            | Authorization | bearer invalidToken |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional"
        Then the response status code should be 202

    Scenario: Test optional mTLS and optional OAuth2 with no client certificate or token in header
        Given The system is ready
        And I have a valid subscription
        When I use the APK Conf file "artifacts/apk-confs/mtls/mtls_optional_oauth2_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee/" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional"
        Then the response status code should be 202