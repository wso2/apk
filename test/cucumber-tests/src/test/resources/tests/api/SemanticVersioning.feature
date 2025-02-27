Feature: Semantic Versioning Based Intelligent Routing

  Scenario: API version with Major and Minor version 1.0
      Given The system is ready
      And I have a valid subscription
      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v1-0.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120005"
      Then I set headers
      |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120005-token}|

      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1.0/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.0\""

      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.0\""

  Scenario: API version with Major and Minor version 1.1
      Given The system is ready
      And I have a valid subscription
      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v1-1.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120005"
      Then I set headers
      |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120005-token}|
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1.1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.1\""
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.1\""

  Scenario: API version with Major and Minor version 1.5
      Given The system is ready
      And I have a valid subscription
      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v1-5.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120005"
      Then I set headers
      |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120005-token}|
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1.5/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.5\""
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.5\""

      When I undeploy the API whose ID is "sem-api-v1-5"
      Then the response status code should be 202
      And I wait for 2 seconds
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      And the response body should contain "\"version\":\"v1.1\""

      When I undeploy the API whose ID is "sem-api-v1-1"
      Then the response status code should be 202
      And I wait for 2 seconds
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      And the response body should contain "\"version\":\"v1.0\""

      When I undeploy the API whose ID is "sem-api-v1-0"
      Then the response status code should be 202

  Scenario: Multiple Major and minor versions for an API
      Given The system is ready
      And I have a valid subscription
      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v1-0.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v1-1.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v1-5.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200
      Then I generate JWT token from idp1 with kid "123-456" and consumer_key "45f1c5c8-a92e-11ed-afa1-0242ac120005"
      Then I set headers
      |Authorization|Bearer ${idp-1-45f1c5c8-a92e-11ed-afa1-0242ac120005-token}|
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.5\""

      When I undeploy the API whose ID is "sem-api-v1-1"
      Then the response status code should be 202
      And I wait for 2 seconds
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.5\""

      When I use the APK Conf file "artifacts/apk-confs/semantic-versioning/sem_api_v2-1.yaml"
      And the definition file "artifacts/definitions/employees_api.json"
      And make the API deployment request
      Then the response status code should be 200

      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v2/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v2.1\""
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v2.1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v2.1\""

      When I undeploy the API whose ID is "sem-api-v1-0"
      Then the response status code should be 202
      And I wait for 2 seconds
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v1.5\""

      When I undeploy the API whose ID is "sem-api-v1-5"
      Then the response status code should be 202
      And I wait for 2 seconds
      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v1/employee/" with body ""
      Then the response status code should be 404

      And I send "GET" request to "https://default.gw.wso2.com:9095/sem-api/v2/employee/" with body ""
      Then the response status code should be 200
      And the response body should contain "\"version\":\"v2.1\""

      When I undeploy the API whose ID is "sem-api-v2-1"
      Then the response status code should be 202
