Feature: API Visibility and Access Control
    Scenario: View APIs with a groups claim and API view scope in the access token
        Given The system is ready
        And I have a valid subscription with groups
            | group1 |
        When I make the GET APIs call to the backoffice
        Then the response status code should be 200
        And the response body should contain "\"count\":1"

    Scenario: View APIs without a groups claim in the access token
        Given The system is ready
        And I have a valid subscription with scopes
            | apk:api_view |
        When I make the GET APIs call to the backoffice
        Then the response status code should be 200
        And the response body should contain "\"count\":0"
