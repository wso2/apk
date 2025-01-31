Feature: API backend based AI ratelimit Feature

  Scenario: backend based AI ratelimit token detail comes in the body.
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4999 |
    Then I see following strings in the enforcer logs
      |aiMetadata|
      |gpt-35-turbo|
      |AzureAI|
      |2024-06-01|
      |aiTokenUsage|
      |1000|
      |300|
      |500|
      |hour|
      |vendor_name|
      |vendor_version|
      |totalTokens|
      |promptTokens|
      |completionTokens|
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4699 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body&prompt_tokens=40000" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4399 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body&prompt_tokens=40000" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body&completion_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body&total_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 429

  Scenario: backend based AI ratelimit token detail comes in the header.
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4999 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4699 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header&prompt_tokens=40000" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4399 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header&prompt_tokens=40000" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header&completion_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header&total_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-header/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 429

  Scenario: backend based AI ratelimit token detail comes in the header but a body configured api checked.
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4999 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4998 |

  Scenario: apk conf backend based AI ratelimit token detail comes in the body.
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/backend_based_airl_conf.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I set headers
      | Authorization | Bearer ${accessToken} |
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4999 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4699 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body&prompt_tokens=40000" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining | 4399 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body&prompt_tokens=40000" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body&completion_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body&total_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/backend-based-airl/1.0.0/employee?send=body" with body ""
    Then the response status code should be 429

  Scenario Outline: Undeploy API
    Given The system is ready
    And I have a valid subscription
    When I undeploy the API whose ID is "<apiID>"
    Then the response status code should be <expectedStatusCode>

    Examples:
      | apiID              | expectedStatusCode |
      | backend-based-airl |                202 |
