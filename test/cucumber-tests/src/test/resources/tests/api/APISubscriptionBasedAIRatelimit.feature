Feature: API subscription based AI ratelimit Feature
  Scenario: subscription based AI ratelimit token detail comes in the body.
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining      | 4999 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining      | 4699 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body&prompt_tokens=40000" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining      | 4399 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body&prompt_tokens=40000" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body&completion_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 429
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body&total_tokens=40000" with body ""
    Then the response status code should be 200
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=body" with body ""
    Then the response status code should be 429
  
  Scenario: subs based AI ratelimit token detail comes in the header but a body configured api checked.
    Given The system is ready
    And I have a valid subscription
    Then I set headers
      |Authorization|bearer ${accessToken}|
    And I wait for next minute strictly
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining      | 4999 |
    And I wait for 3 seconds
    And I send "GET" request to "https://default.gw.wso2.com:9095/llm-api-subs/v1.0.0/3.14/employee?send=header" with body ""
    Then the response status code should be 200
    And the response headers should contain
      | x-ratelimit-remaining      | 4998 |