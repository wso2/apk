Feature: API Deployment with Global Interceptor
  Scenario: Deploying an API
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/globalInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "579ba27a1e03e2fdf099d1b6745e265f2d495606"
    And I wait for api deployment
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/globalinterceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Gw-Header\": \"GW-header-value\""
    Then the response headers contains key "gw-response-header" and value "GW-response-header-value"
  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID | expectedStatusCode |
      | 579ba27a1e03e2fdf099d1b6745e265f2d495606  | 202         |