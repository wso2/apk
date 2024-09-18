Feature: CORS Policy

  Scenario: Testing CORS Policy
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/cors_API.apk-conf"
    And the definition file "artifacts/definitions/cors_api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And I wait for next minute
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 204 response code, not accepting
      | 429 |
    And the response headers should not contain
      | Access-Control-Allow-Origin      |
      | Access-Control-Allow-Credentials |
      | Access-Control-Allow-Methods     |
      | Access-Control-Allow-Headers     |
      | Access-Control-Max-Age           |
    Then I set headers
      | Origin | test.domain.com |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 204 response code, not accepting
      | 429 |
    And the response headers should not contain
      | Access-Control-Allow-Origin      |
      | Access-Control-Allow-Credentials |
      | Access-Control-Allow-Methods     |
      | Access-Control-Allow-Headers     |
      | Access-Control-Max-Age           |
    Then I set headers
      | Origin | abc.com |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 204 response code, not accepting
      | 429 |
    And the response headers should contain
      | Access-Control-Allow-Origin      | abc.com |
      | Access-Control-Allow-Credentials | true    |
    Then I set headers
      | Origin                        | abc.com |
      | Access-Control-Request-Method | GET     |
    And I send "OPTIONS" request to "https://default.gw.wso2.com:9095/test_cors/2.0.0/anything/" with body ""
    And I eventually receive 200 response code, not accepting
      | 429 |
    And the response headers should contain
      | Access-Control-Allow-Origin      | abc.com                     |
      | Access-Control-Allow-Credentials | true                        |
      | Access-Control-Allow-Methods     | GET, POST, PUT, DELETE      |
      | Access-Control-Allow-Headers     | Content-Type, Authorization |
      | Access-Control-Max-Age           | 3600                        |

  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                                         | expectedStatusCode |
      | cors-api-adff3dbc-2787-11ee-be56-0242ac120002 | 202                |
