Feature: API Deployment with Interceptor
  Scenario: Deploying an API
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/original.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should not contain "\"Interceptor-Header\""
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withRequestInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withResponseInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withResponseInterceptorParameterVariation1.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withResponseInterceptorParameterVariation2.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withResponseInterceptorParameterVariation3.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withResponseInterceptorParameterVariation4.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"

    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withRequestAndResponse.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withRequestAndResponsetls.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "547961eeaafed989119c45ffc13f8b87bfda821d"
    And I wait for 1 minute
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""
    Then the response status code should be 200
    Then the response headers contains key "interceptor-response-header" and value "Interceptor-Response-header-value"
  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID | expectedStatusCode |
      | 547961eeaafed989119c45ffc13f8b87bfda821d  | 202         |