#Feature: Deploy APIs in multiple environments
#  Scenario: Deploying an API without specifing an Environment and token issuer has no environments.
#    Given The system is ready
#    And I have a valid subscription
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${accessToken}|
#    And I send "GET" request to "https://default.gw.wso2.com:9095/withoutenv/3.14/employee/" with body ""
#    And I eventually receive 200 response code, not accepting
#      |429|
#    When I undeploy the API whose ID is "without-env-api"
#    Then the response status code should be 202

#  Scenario: Deploying an API without specifing an Environment and token issuer has all(*) environments.
#    Given The system is ready
#    And I have a valid token for organization "org3"
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf.yaml"
#   And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request for organization "org3"
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${org3}|
#    And I send "GET" request to "https://org3.gw.wso2.com:9095/withoutenv/3.14/employee/" with body ""
#    And I eventually receive 200 response code, not accepting
#      |429|
#    When I undeploy the API whose ID is "without-env-api" and organization "org3"
#    Then the response status code should be 202

#  Scenario: Deploying an API without specifing an Environment and token issuer has only dev environment.
#    Given The system is ready
#    And I have a valid token for organization "org4"
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request for organization "org4"
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${org4}|
#    And I send "GET" request to "https://org4.gw.wso2.com:9095/withoutenv/3.14/employee/" with body ""
#    And I eventually receive 401 response code, not accepting
#      |200|
#   When I undeploy the API whose ID is "without-env-api" and organization "org4"
#    Then the response status code should be 202

#  Scenario: Deploying APIs in Dev and QA environments and token issuer has no environments.
#    Given The system is ready
#    And I have a valid subscription
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf_dev.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${accessToken}|
#    And I send "GET" request to "https://default-dev.gw.wso2.com:9095/multienv/3.14/employee/" with body ""
#    And I eventually receive 200 response code, not accepting
#      |429|
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf_qa.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${accessToken}|
#    And I send "GET" request to "https://default-qa.gw.wso2.com:9095/multienv/3.14/employee/" with body ""
#    And I eventually receive 200 response code, not accepting
#     |429|
#    When I undeploy the API whose ID is "multi-env-dev-api"
#    Then the response status code should be 202
#    When I undeploy the API whose ID is "multi-env-qa-api"
#    Then the response status code should be 202
    
#  Scenario: Deploying an API in QA environment and token issuer has all(*) environments.
#    Given The system is ready
#    And I have a valid token for organization "org3"
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf_qa.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request for organization "org3"
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${org3}|
#    And I send "GET" request to "https://org3-qa.gw.wso2.com:9095/multienv/3.14/employee/" with body ""
#    And I eventually receive 200 response code, not accepting
#      |401|
#    When I undeploy the API whose ID is "multi-env-qa-api" and organization "org3"
#    Then the response status code should be 202

#  Scenario: Deploying an API in QA environment and token issuer has only Dev environment.
#    Given The system is ready
#    And I have a valid token for organization "org4"
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf_qa.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request for organization "org4"
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${org4}|
#    And I send "GET" request to "https://org4-qa.gw.wso2.com:9095/multienv/3.14/employee/" with body ""
#    And I eventually receive 401 response code, not accepting
#      |200|
#    When I undeploy the API whose ID is "multi-env-qa-api" and organization "org4"
#    Then the response status code should be 202

#  Scenario: Deploying an API in Dev environment and token issuer has only Dev environment.
#    Given The system is ready
#    And I have a valid token for organization "org4"
#    When I use the APK Conf file "artifacts/apk-confs/multi-env/employees_conf_dev.yaml"
#    And the definition file "artifacts/definitions/employees_api.json"
#    And make the API deployment request for organization "org4"
#    Then the response status code should be 200
#    Then I set headers
#      |Authorization|Bearer ${org4}|
#    And I send "GET" request to "https://org4-dev.gw.wso2.com:9095/multienv/3.14/employee/" with body ""
#    And I eventually receive 200 response code, not accepting
#      |401|
#    When I undeploy the API whose ID is "multi-env-dev-api" and organization "org4"
#    Then the response status code should be 202
  