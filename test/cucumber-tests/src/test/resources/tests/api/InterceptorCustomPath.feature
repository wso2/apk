Feature: Interceptor Custom Resource Path
  Scenario: Deploying an API with interceptor custom resource paths
    Given The system is ready
    And I have a valid subscription

    # Test 1: API-level request interceptor with custom path
    # backendUrl includes /custom/v2/handle-custom-request, so the interceptor call
    # hits /custom/v2/handle-custom-request instead of the default /api/v1/handle-request
    When I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathRequestInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    Then I wait for 10 seconds
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Custom-Interceptor-Header\": \"Custom-path-interceptor-value\""
    And the response body should not contain "\"Interceptor-Header\""

    # Test 2: API-level response interceptor with custom path
    # backendUrl includes /custom/v2/handle-custom-response, so the interceptor call
    # hits /custom/v2/handle-custom-response instead of the default /api/v1/handle-response
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathResponseInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    Then I wait for 10 seconds
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "custom-interceptor-response-header" and value "Custom-path-response-interceptor-value"

    # Test 3: API-level both request and response interceptors with custom paths
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathRequestAndResponse.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    Then I wait for 10 seconds
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Custom-Interceptor-Header\": \"Custom-path-interceptor-value\""
    Then the response headers contains key "custom-interceptor-response-header" and value "Custom-path-response-interceptor-value"

    # Test 4: API-level custom request path + default response path
    # Request interceptor uses custom /custom/v2/handle-custom-request
    # Response interceptor uses default backendUrl (no path), falling back to /api/v1/handle-response
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathRequestDefaultResponse.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    Then I wait for 10 seconds
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Custom-Interceptor-Header\": \"Custom-path-interceptor-value\""
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

    # Test 5: API-level default request path + custom response path
    # Request interceptor uses default backendUrl (no path), falling back to /api/v1/handle-request
    # Response interceptor uses custom /custom/v2/handle-custom-response
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withDefaultRequestCustomPathResponse.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    Then I wait for 10 seconds
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
    Then the response headers contains key "custom-interceptor-response-header" and value "Custom-path-response-interceptor-value"

    # Test 6: Verify default paths still work (backward compatibility)
    # backendUrl without path uses defaults /api/v1/handle-request and /api/v1/handle-response
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withRequestAndResponse.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    Then I wait for 10 seconds
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
    And the response body should not contain "\"Custom-Interceptor-Header\""
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                    | expectedStatusCode |
      | 547961eeaafed989119c45ffc13f8b87bfda821d | 202                |
