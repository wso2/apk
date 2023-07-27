Feature: API Deployment with Interceptor
  Scenario: Deploying an API
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/interceptors/original.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "268309bb509758d0ec2ac03f96929cbb001e10cb"
    And I wait for 1 minute
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should not contain "\"Interceptor-Header\""
    Then I use the APK Conf file "artifacts/apk-confs/interceptors/withRequestInterceptor.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "268309bb509758d0ec2ac03f96929cbb001e10cb"
    And I wait for 1 minute
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/interceptor/1.0.0/get" with body ""
    Then the response status code should be 200
    And the response body should contain "\"Interceptor-Header\": \"Interceptor-header-value\""

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID | expectedStatusCode |
      | 268309bb509758d0ec2ac03f96929cbb001e10cb  | 202         |