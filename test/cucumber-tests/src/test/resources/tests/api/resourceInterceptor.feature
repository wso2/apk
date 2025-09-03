Feature: API Deployment with Resource Interceptor
  Scenario: Deploying an API with Resource Level Lua Interceptor
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/resourceLevelInterceptorLua.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/lua-interceptor/1.0.0/get" with body ""
    And the response body should not contain "\"Lua-Interceptor-Header\""
    Then the response status code should be 200
    Then the response headers not contains key "lua-interceptor-response-header"
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/lua-interceptor/1.0.0/headers" with body ""
    And the response body should contain
      | "Lua-Interceptor-Header": "Lua-Interceptor-header-value" |
    Then the response status code should be 200
    Then the response headers contains key "lua-interceptor-response-header" and value "Lua-Interceptor-Response-header-value"

  Scenario: Deploying an API with Resource Level WASM Interceptor
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/resourceLevelInterceptorWASM.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821e"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/wasm-interceptor/1.0.0/get" with body ""
    And the response status code should be 200
    Then the response headers not contains key "x-wasm-custom"
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/wasm-interceptor/1.0.0/headers" with body ""
    And the response status code should be 200
    Then the response headers contains key "x-wasm-custom" and value "FOO"

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                    | expectedStatusCode |
      | 547961eeaafed989119c45ffc13f8b87bfda821d | 202                |
      | 547961eeaafed989119c45ffc13f8b87bfda821e | 202                |
