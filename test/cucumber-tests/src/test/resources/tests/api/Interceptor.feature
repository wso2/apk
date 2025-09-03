Feature: API Deployment with Interceptor
  Scenario: Deploying an API without Interceptors
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/original.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should not contain "\"Interceptor-Header\""

  Scenario: Deploying an API with API level Lua Request Interceptor
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/withLuaRequestInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Lua-Interceptor-Header\": \"Lua-Interceptor-header-value\""

  Scenario: Deploying an API with API level Lua Response Interceptor
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/withLuaResponseInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "lua-interceptor-response-header" and value "Lua-Interceptor-Response-header-value"

  Scenario: Deploying an API with API level WASM Response Interceptor
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/withWASMResponseInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "x-wasm-custom" and value "FOO"

  Scenario: Deploying an API with API and resource level Lua Interceptor and WASM Interceptors
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/withLuaAndWASMInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/uuid" with body ""
    Then the response status code should be 200
    And the response body should contain
      | "Lua-Interceptor-Request-Global": "Lua-Interceptor-Request-global-header-value" |
      | "Lua-Request-Interceptor-Header": "Lua-Request-Interceptor-header-value"        |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/ai/spelling" with body ""
    Then the response status code should be 200
    And the response body should contain
      | "Lua-Interceptor-Request-Global": "Lua-Interceptor-Request-global-header-value" |
    Then the response headers contains key "x-wasm-custom" and value "FOO"
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/ip" with body ""
    Then the response status code should be 200
    And the response body should contain
      | "Lua-Interceptor-Request-Global": "Lua-Interceptor-Request-global-header-value" |
      | "Lua-Request-Interceptor-Header-1": "Lua-Request-Interceptor-header-value-1"    |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/user-agent" with body ""
    Then the response status code should be 200
    And the response body should contain
      | "Lua-Interceptor-Request-Global": "Lua-Interceptor-Request-global-header-value" |
    Then the response headers contains key "lua-interceptor-response-header" and value "Lua-Interceptor-Response-header-value"

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                    | expectedStatusCode |
      | 547961eeaafed989119c45ffc13f8b87bfda821d | 202                |