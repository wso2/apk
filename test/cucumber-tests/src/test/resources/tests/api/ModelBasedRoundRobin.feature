Feature: AI Model Based Round Robin
  Scenario: Testing API level API Policy with Model Based Round Robin and Multi Endpoints
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/modelroundrobin/api-level-policy-multi-endpoints.yaml"
    And the definition file "artifacts/definitions/backend_apikey_auth_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/chat-service-prod-sand/1.0/get" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    Then I send "POST" request to "https://default.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"
    Then I send "POST" request to "https://default.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"
    Then I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/chat-service-prod-sand/1.0/get" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    Then I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"
    Then I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"

  Scenario Outline: Undeploy API Level API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | chat-service-api-prod-sand          | 202                 |

  Scenario: Testing Resource level API Policy with Model Based Round Robin and Multi Endpoints
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/modelroundrobin/resource-level-policy-multi-endpoints.yaml"
    And the definition file "artifacts/definitions/backend_apikey_auth_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    And I send "GET" request to "https://default.gw.wso2.com:9095/chat-service-prod-sand/1.0/get" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    Then I send "POST" request to "https://default.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"
    Then I send "POST" request to "https://default.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"
    Then I send "GET" request to "https://default.sandbox.gw.wso2.com:9095/chat-service-prod-sand/1.0/get" with body ""
    And I eventually receive 200 response code, not accepting
      |429|
    Then I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"
    Then I send "POST" request to "https://default.sandbox.gw.wso2.com:9095/chat-service-prod-sand/1.0/chat/completions" with body "{\"model\": \"gpt-4.5\", \"input\": \"Hello\"}"
    And I eventually receive 200 response code, not accepting
      |429|
    And the response body should contain "gpt-"

  Scenario Outline: Undeploy Resource Level API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID                 | expectedStatusCode  |
      | chat-service-api-prod-sand          | 202                 |