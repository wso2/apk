Feature: Test all valid security combinations
    Scenario: Test mandatory mtls and mandatory oauth2 and mandatory apikey with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_mandatory_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-mandatory-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_mandatory_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-mandatory-apikey-optional"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_mandatory_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-mandatory-apikey-disabled"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_optional_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-optional-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_optional_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-optional-apikey-optional"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_optional_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-optional-apikey-disabled"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_disabled_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-disabled-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_disabled_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-disabled-apikey-optional"
        Then the response status code should be 202

    Scenario: Test mandatory mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_mandatory_oauth2_disabled_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-mandatory-oauth2-disabled-apikey-disabled"
        Then the response status code should be 202

    Scenario: Test optional mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_mandatory_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-mandatory-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test optional mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_mandatory_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-mandatory-apikey-optional"
        Then the response status code should be 202

    Scenario: Test optional mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_mandatory_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-mandatory-apikey-disabled"
        Then the response status code should be 202

    Scenario: Test optional mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_optional_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test optional mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_optional_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional-apikey-optional"
        Then the response status code should be 202

    Scenario: Test optional mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_optional_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-optional-apikey-disabled"
        Then the response status code should be 202

    Scenario: Test optional mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_disabled_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-disabled-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test optional mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_disabled_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-disabled-apikey-optional"
        Then the response status code should be 202

    Scenario: Test optional mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_optional_oauth2_disabled_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 406

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-optional-oauth2-disabled-apikey-disabled"
        Then the response status code should be 404

    Scenario: Test disabled mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_mandatory_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-mandatory-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test disabled mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_mandatory_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-mandatory-apikey-optional"
        Then the response status code should be 202

    Scenario: Test disabled mtls and mandatory oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_mandatory_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-mandatory-apikey-disabled"
        Then the response status code should be 202

    Scenario: Test disabled mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_optional_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-optional-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test disabled mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_optional_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 406

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-optional-apikey-optional"
        Then the response status code should be 404

    Scenario: Test disabled mtls and optional oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_optional_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 406

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-optional-apikey-disabled"
        Then the response status code should be 404

    Scenario: Test disabled mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_disabled_apikey_mandatory.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 200

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-disabled-apikey-mandatory"
        Then the response status code should be 202

    Scenario: Test disabled mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_disabled_apikey_optional.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 406

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-disabled-apikey-optional"
        Then the response status code should be 404

    Scenario: Test disabled mtls and disabled oauth2 with a valid client certificate in header
        Given The system is ready
        And I have a valid token with a client certificate "config-map-1.txt"
        When I use the APK Conf file "artifacts/apk-confs/apikey/mtls_disabled_oauth2_disabled_apikey_disabled.apk-conf"
        And the definition file "artifacts/definitions/employees_api.json"
        And make the API deployment request
        Then the response status code should be 406

    Scenario: Undeploy API
        Given The system is ready
        And I have a valid subscription
        When I undeploy the API whose ID is "mtls-disabled-oauth2-disabled-apikey-disabled"
        Then the response status code should be 404

