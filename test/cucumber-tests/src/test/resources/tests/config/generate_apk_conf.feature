# Feature: Generating APK conf
#   Scenario: Generating APK conf using the API definition
#     Given The system is ready
#     When I use the definition file "artifacts/definitions/sample_api.yaml" in resources
#     And generate the APK conf file for a "REST" API
#     Then the response status code should be 200
#     And the response body should be "artifacts/apk-confs/sampleAPK.apk-conf" in resources
#   Scenario: Generating APK conf using and invalid API definition
#     Given The system is ready
#     When I use the definition file "artifacts/definitions/invalid_api.yaml" in resources
#     And generate the APK conf file for a "REST" API
#     Then the response status code should be 400
