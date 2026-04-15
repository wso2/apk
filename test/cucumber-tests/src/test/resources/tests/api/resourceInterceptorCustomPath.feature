Feature: Interceptor Custom Resource Path
  Scenario: Deploying an API with interceptor custom resource paths
    Given The system is ready
    And I have a valid subscription

    # Test 1: Operation-level request interceptor with custom path
    # Interceptor applied at operation level (operationPolicies) on GET /get only
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathOperationRequestInterceptor.apk-conf"
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
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/headers" with body ""
    Then the response status code should be 200
    And the response body should not contain "\"Custom-Interceptor-Header\""
    And the response body should not contain "\"Interceptor-Header\""

    # Test 2: Operation-level response interceptor with custom path
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathOperationResponseInterceptor.apk-conf"
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
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/headers" with body ""
    Then the response status code should be 200
    Then the response headers not contains key "custom-interceptor-response-header"

    # Test 3: Operation-level both request and response interceptors with custom paths
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathOperationRequestAndResponse.apk-conf"
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
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/headers" with body ""
    Then the response status code should be 200
    And the response body should not contain "\"Custom-Interceptor-Header\""
    Then the response headers not contains key "custom-interceptor-response-header"

    # Test 4: Operation-level custom request path + default response path
    # Request interceptor uses custom path, response interceptor uses default /api/v1/handle-response
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withCustomPathOperationRequestDefaultResponse.apk-conf"
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
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/headers" with body ""
    Then the response status code should be 200
    And the response body should not contain "\"Custom-Interceptor-Header\""
    And the response body should not contain "\"Interceptor-Header\""
    Then the response headers not contains key "interceptor-response-header"

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                    | expectedStatusCode |
      | 547961eeaafed989119c45ffc13f8b87bfda821d | 202                |
