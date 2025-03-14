Feature: WeightedRouting
  Scenario Outline: Testing Weighted Routing across Mulitple Endpoints with Different Weights
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/weighted-routing/weighted_routing_sample_<sample_config_number>.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I clear all stored responses
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    Then I send 200 requests to "https://default.gw.wso2.com:9095/weightedrouting/1.0/demo" to test weighted routing and count the responses from each endpoint
    And I ensure that the weights of the endpoints increase in the order of API_version: "<api1>", "<api2>", "<api3>" from the response counts
    Then I clear all stored responses and weight counts
    When I undeploy the API whose ID is "weighted-routing-sample"
    Then the response status code should be 202

    Examples:
      | sample_config_number | api1 | api2 | api3 |
      | 1                    | 2.0  | 1.0  | 3.0  |
      | 2                    | 3.0  | 2.0  | 1.0  |

  Scenario: Testing Weighted Routing across Mulitple Endpoints with Different Weights and an Endpoint with Zero Weight
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/weighted-routing/weighted_routing_zero_weight_endpoint.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I clear all stored responses
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    Then I send 100 requests to "https://default.gw.wso2.com:9095/weightedrouting/1.0/demo" to test weighted routing and count the responses from each endpoint
    And I ensure that the weights of the endpoints increase in the order of API_version: "2.0", "3.0", "1.0" from the response counts
    And I ensure that the response count of the endpoint "2.0" with zero weight is zero
    Then I clear all stored responses and weight counts
    When I undeploy the API whose ID is "weighted-routing-sample"
    Then the response status code should be 202

  Scenario: Testing Weighted Routing across Mulitple Endpoints with Equal Weights (weight greater than 1)
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/weighted-routing/weighted_routing_equal_weights.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I clear all stored responses
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    Then I send 1000 requests to "https://default.gw.wso2.com:9095/weightedrouting/1.0/demo" to test weighted routing and count the responses from each endpoint
    And I ensure that the responses are distributed approximately equally among the endpoints
    Then I clear all stored responses and weight counts
    When I undeploy the API whose ID is "weighted-routing-equal-weights"
    Then the response status code should be 202

  Scenario: Testing Weighted Routing across Mulitple Endpoints with Equal Weights (weights equal to 1)
    Given The system is ready
    And I have a valid subscription
    When I use the APK Conf file "artifacts/apk-confs/weighted-routing/weighted_routing_equal_weights_all_one.yaml"
    And the definition file "artifacts/definitions/employees_api.json"
    And make the API deployment request
    Then the response status code should be 200
    Then I clear all stored responses
    Then I set headers
      |Authorization|Bearer ${accessToken}|
    Then I send 10 requests to "https://default.gw.wso2.com:9095/weightedrouting/1.0/demo" to test weighted routing and count the responses from each endpoint
    And I ensure that all the responses are from one of the endpoints
    Then I clear all stored responses and weight counts
    When I undeploy the API whose ID is "weighted-routing-equal-weights"
    Then the response status code should be 202
