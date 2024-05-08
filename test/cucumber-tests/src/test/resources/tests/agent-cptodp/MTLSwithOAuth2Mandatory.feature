Feature: Test mTLS between client and gateway with client certificate sent in header
    Background:
        Given The system is ready
    # mTLS mandatory OAuth2 mandatory
    Scenario: Test mandatory mTLS and mandatory OAuth2 with a valid client certificate in header
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_enabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200
        And I have a valid Devportal access token
        And make the Application Creation request with the name "SampleApp"
        Then the response status code should be 201
        And the response body should contain "SampleApp"
        And I have a KeyManager
        And make the Generate Keys request
        Then the response status code should be 200
        And the response body should contain "consumerKey"
        And the response body should contain "consumerSecret"
        And make the Subscription request
        Then the response status code should be 201
        And the response body should contain "Unlimited"
        And I get "production" oauth keys for application
        Then the response status code should be 200
        And make the Access Token Generation request for "production"
        Then the response status code should be 200
        And the response body should contain "accessToken"      
        Then I set headers
              | Authorization             | Bearer ${accessToken} |
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
            | 403 |
            | 429 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Devportal access token
        Then I delete the application "SampleApp" from devportal
        Then the response status code should be 200
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|


    Scenario: Test mandatory mTLS and mandatory OAuth2 with a invalid client certificate in header
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_enabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "invalid-cert.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 500

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200

    Scenario: Test mandatory mTLS and mandatory OAuth2 without client certificate in header
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_enabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200
        And I have a valid Devportal access token
        And make the Application Creation request with the name "SampleApp"
        Then the response status code should be 201
        And the response body should contain "SampleApp"
        And I have a KeyManager
        And make the Generate Keys request
        Then the response status code should be 200
        And the response body should contain "consumerKey"
        And the response body should contain "consumerSecret"
        And make the Subscription request
        Then the response status code should be 201
        And the response body should contain "Unlimited"
        And I get "production" oauth keys for application
        Then the response status code should be 200
        And make the Access Token Generation request for "production"
        Then the response status code should be 200
        And the response body should contain "accessToken"
        Then I set headers
              | Authorization             | Bearer ${accessToken} |     
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Devportal access token
        Then I delete the application "SampleApp" from devportal
        Then the response status code should be 200
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|

    
    # mTLS optional OAuth2 mandatory
    Scenario: Test optional mTLS and mandatory OAuth2 with a valid client certificate in header
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_enabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200
        And I have a valid Devportal access token
        And make the Application Creation request with the name "SampleApp"
        Then the response status code should be 201
        And the response body should contain "SampleApp"
        And I have a KeyManager
        And make the Generate Keys request
        Then the response status code should be 200
        And the response body should contain "consumerKey"
        And the response body should contain "consumerSecret"
        And make the Subscription request
        Then the response status code should be 201
        And the response body should contain "Unlimited"
        And I get "production" oauth keys for application
        Then the response status code should be 200
        And make the Access Token Generation request for "production"
        Then the response status code should be 200
        And the response body should contain "accessToken"      
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
            | 403 |
            | 429 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Devportal access token
        Then I delete the application "SampleApp" from devportal
        Then the response status code should be 200
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|

    Scenario: Test optional mTLS and mandatory OAuth2 without a token
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_optional_oauth2_mandatory.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200  
        Then I set headers
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|

    
    Scenario: Test optional mTLS and mandatory OAuth2 with an invalid token in header
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_optional_oauth2_mandatory.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200
        And I set "invlaidToken" as the new access token
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|

    # #TODO: ERROR Exists IN THE APK SIDE that NEED TO BE FIXED. Enable this test case after that
    # Scenario: Test mandatory mTLS and disabled OAuth2 with a valid client certificate in header
    #     And I have a DCR application
    #     And I have a valid Publisher access token
    #     When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_disabled.json"
    #     When the definition file "artifacts/definitions/cors-definition.json"
    #     And make the import API Creation request using OAS "File"
    #     Then the response status code should be 201
    #     And the response body should contain "EmployeeServiceAPI"
    #     And I have a client certificate "config-map-1.crt"
    #     Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
    #     Then the response status code should be 201
    #     And make the API Revision Deployment request
    #     Then the response status code should be 201
    #     Then I wait for 40 seconds
    #     And make the Change Lifecycle request
    #     Then the response status code should be 200  
    #     Then I set headers
    #         | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
    #     And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
    #     And I eventually receive 200 response code, not accepting
    #         | 401 |

    # Scenario: Undeploy the created REST API
    #     And I have a DCR application
    #     And I have a valid Devportal access token
    #     Then I delete the application "SampleApp" from devportal
    #     Then the response status code should be 200
    #     And I have a valid Publisher access token
    #     Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
    #     Then I undeploy the selected API
    #     Then the response status code should be 200
    #     And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
    #     And I eventually receive 404 response code, not accepting
    #     |200|
      
    # Scenario: Test mandatory mTLS and disabled OAuth2 with an invalid client certificate in header

    Scenario: Test mandatory mTLS and disabled OAuth2 without a client certificate in header
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_disabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200  
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|

    # Multiple certificates test cases
    Scenario: Test an API with mTLS enabled and one associated certificate with multiple certificates existing in system
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_enabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200
        And I have a valid Devportal access token
        And make the Application Creation request with the name "SampleApp"
        Then the response status code should be 201
        And the response body should contain "SampleApp"
        And I have a KeyManager
        And make the Generate Keys request
        Then the response status code should be 200
        And the response body should contain "consumerKey"
        And the response body should contain "consumerSecret"
        And make the Subscription request
        Then the response status code should be 201
        And the response body should contain "Unlimited"
        And I get "production" oauth keys for application
        Then the response status code should be 200
        And make the Access Token Generation request for "production"
        Then the response status code should be 200
        And the response body should contain "accessToken"      
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
        And I have a client certificate "config-map-2.crt"
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |
        And I have a client certificate "config-map-3.crt"
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Devportal access token
        Then I delete the application "SampleApp" from devportal
        Then the response status code should be 200
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|

    Scenario: Test an API with mTLS enabled and multiple certificates configured
        And I have a DCR application
        And I have a valid Publisher access token
        When I use the Payload file "artifacts/payloads/mtls/mtls_mandatory_oauth2_enabled.json"
        When the definition file "artifacts/definitions/cors-definition.json"
        And make the import API Creation request using OAS "File"
        Then the response status code should be 201
        And the response body should contain "EmployeeServiceAPI"
        And I have a client certificate "config-map-2.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap2"
        Then the response status code should be 201
        And I have a client certificate "config-map-1.crt"
        Then I update the API with mtls certificate data with the alias "mtls-test-configmap"
        Then the response status code should be 201
        And make the API Revision Deployment request
        Then the response status code should be 201
        Then I wait for 40 seconds
        And make the Change Lifecycle request
        Then the response status code should be 200
        And I have a valid Devportal access token
        And make the Application Creation request with the name "SampleApp"
        Then the response status code should be 201
        And the response body should contain "SampleApp"
        And I have a KeyManager
        And make the Generate Keys request
        Then the response status code should be 200
        And the response body should contain "consumerKey"
        And the response body should contain "consumerSecret"
        And make the Subscription request
        Then the response status code should be 201
        And the response body should contain "Unlimited"
        And I get "production" oauth keys for application
        Then the response status code should be 200
        And make the Access Token Generation request for "production"
        Then the response status code should be 200
        And the response body should contain "accessToken"      
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
        And I have a client certificate "config-map-2.crt"
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 200 response code, not accepting
            | 401 |
        And I have a client certificate "config-map-3.crt"
        Then I set headers
              | Authorization             | Bearer ${accessToken} | 
            | X-WSO2-CLIENT-CERTIFICATE | ${clientCertificate}  |
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 401 response code, not accepting
            | 200 |

    Scenario: Undeploy the created REST API
        And I have a DCR application
        And I have a valid Devportal access token
        Then I delete the application "SampleApp" from devportal
        Then the response status code should be 200
        And I have a valid Publisher access token
        Then I find the apiUUID of the API created with the name "EmployeeServiceAPI"
        Then I undeploy the selected API
        Then the response status code should be 200
        And I send "GET" request to "https://default.gw.wso2.com:9095/mtls/3.14/employee" with body ""
        And I eventually receive 404 response code, not accepting
        |200|
