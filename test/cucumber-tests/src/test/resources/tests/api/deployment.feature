Feature: API Deployment
  Scenario: Deploying an API
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "/home/sanoj/work/repos/apk-test/expectedAPK.apk-conf"
    And the definition file "/home/sanoj/work/repos/apk-test/api.yaml"
    And make the API deployment request
    Then the response status code should be 200
    And the response body should contain "id"

  Scenario Outline: Undeploy an API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID | expectedStatusCode |
      | b76df54c6c2e37adfca2c9afda3439949a34f73a  | 202         |
      | abcdeadsxzads | 404        |
