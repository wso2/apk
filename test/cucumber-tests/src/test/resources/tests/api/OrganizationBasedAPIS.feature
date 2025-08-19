# Feature: Organization Base API Deployment
#   Scenario: Deploying an API and basic http method invocations
#     Given The system is ready
#     And I have a valid token for organization "org1"
#     When I use the APK Conf file "artifacts/apk-confs/employees_conf.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request for organization "org1"
#     Then the response status code should be 200
#     And I have a valid token for organization "org2"
#     When I use the APK Conf file "artifacts/apk-confs/employees_conf.yaml"
#     And the definition file "artifacts/definitions/employees_api.json"
#     And make the API deployment request for organization "org2"
#     Then the response status code should be 200
#     Then I set headers
#       |Authorization|Bearer ${org1}|
#     And I send "GET" request to "https://org1.gw.wso2.com:9095/test/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     Then I set headers
#       |Authorization|Bearer ${org2}|
#     And I send "GET" request to "https://org2.gw.wso2.com:9095/test/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     When I undeploy the API whose ID is "432bf873bf751ad714ddc635b0a4d9d194b39eb3" and organization "org1"
#     Then the response status code should be 202
#     Then I set headers
#       |Authorization|Bearer ${org1}|
#     And I send "GET" request to "https://org1.gw.wso2.com:9095/test/3.14/employee/" with body ""
#     And I eventually receive 404 response code, not accepting
#       |429|
#     Then I set headers
#       |Authorization|Bearer ${org2}|
#     And I send "GET" request to "https://org2.gw.wso2.com:9095/test/3.14/employee/" with body ""
#     And I eventually receive 200 response code, not accepting
#       |429|
#     When I undeploy the API whose ID is "d45e79bec8c4b2d3a7543e09530e0a995ea68691" and organization "org2"
#     Then the response status code should be 202
#     Then I set headers
#       |Authorization|Bearer ${org1}|
#     And I send "GET" request to "https://org1.gw.wso2.com:9095/test/3.14/employee/" with body ""
#     And I eventually receive 404 response code, not accepting
#       |429|
#     Then I set headers
#       |Authorization|Bearer ${org2}|
#     And I send "GET" request to "https://org2.gw.wso2.com:9095/test/3.14/employee/" with body ""
#     And I eventually receive 404 response code, not accepting
#       |429|