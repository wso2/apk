//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

import ballerina/test;
import ballerina/log;
import wso2/apk_common_lib as commons;
import ballerina/uuid;

@test:Mock {functionName: "generateToken"}
test:MockFunction generateTokenMock = new ();

@test:Mock {functionName: "updateApplication", moduleName: "wso2/notification_grpc_client"}
public isolated function updateApplicationMock(ApplicationGRPC updateApplicationRequest, string endpoint, string pubCert, string devCert, string devKey) returns error|NotificationResponse {
    NotificationResponse noti = {code: "OK"};
    return noti;
}

@test:Mock {functionName: "deleteApplication", moduleName: "wso2/notification_grpc_client"}
public isolated function deleteApplicationMock(ApplicationGRPC deleteApplicationRequest, string endpoint, string pubCert, string devCert, string devKey) returns error|NotificationResponse {
    NotificationResponse noti = {code: "OK"};
    return noti;
}

Application application = {name: "sampleApp", description: "sample application"};

@test:BeforeSuite
function beforeFunc1() {
    ApplicationRatePlan payload = {
        "planName": "25PerMin",
        "displayName": "25PerMin",
        "description": "25 Per Min",
        "defaultLimit": {
            "type": "REQUESTCOUNTLIMIT",
            "requestCount": {
                "requestCount": 25,
                "timeUnit": "min",
                "unitTime": 1
            }
        }
    };
    string applicationUsagePlanId = uuid:createType1AsString();
    payload.planId = applicationUsagePlanId;
    ApplicationRatePlan|commons:APKError createdAppPol = addApplicationUsagePlanDAO(payload, organiztion.uuid);
    if createdAppPol is ApplicationRatePlan {
        test:assertTrue(true, "Application usage plan added successfully");
        BusinessPlan payloadbp = {
            "planName": "MyBusinessPlan",
            "displayName": "MyBusinessPlan",
            "description": "test sub pol test",
            "defaultLimit": {
                "type": "REQUESTCOUNTLIMIT",
                "requestCount": {
                    "requestCount": 20,
                    "timeUnit": "min",
                    "unitTime": 1
                }
            },
            "rateLimitCount": 10,
            "rateLimitTimeUnit": "sec",
            "customAttributes": []
        };
        payloadbp.planId = uuid:createType1AsString();
        BusinessPlan|commons:APKError createdBusinessPlan = addBusinessPlanDAO(payloadbp, organiztion.uuid);
        if createdBusinessPlan is BusinessPlan {
            test:assertTrue(true, "Business Plan added successfully");
        } else if createdBusinessPlan is commons:APKError {
            test:assertFail("Error occured while adding Business Plan");
        }
    } else if createdAppPol is commons:APKError {
        log:printError(createdAppPol.toString());
        test:assertFail("Error occured while adding Application Usage Plan");
    }
}

@test:Config {}
function addApplicationTest() {
    string[] testHosts = ["http://localhost:9090"];
    test:when(retrieveManagementServerHostsListMock).thenReturn(testHosts);
    Application payload = {name: "sampleApp", description: "sample application"};
    NotFoundError|Application|commons:APKError createdApplication = addApplication(payload, organiztion, "apkuser");
    if createdApplication is Application {
        test:assertTrue(true, "Successfully added the application");
        application.applicationId = createdApplication.applicationId;
    } else if createdApplication is error {
        test:assertFail("Error occured while adding application");
    }
}

@test:Config {dependsOn: [addApplicationTest]}
function getApplicationByIdTest() {
    string? appId = application.applicationId;
    if appId is string {
        Application|commons:APKError|NotFoundError returnedResponse = getApplicationById(appId, organiztion);
        if returnedResponse is Application {
            test:assertTrue(true, "Successfully retrieved application");
        } else if returnedResponse is commons:APKError|NotFoundError {
            test:assertFail("Error occured while retrieving application");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [getApplicationByIdTest]}
function getApplicationListTest() {
    ApplicationList|commons:APKError applicationList = getApplicationList("", "", "", "", 0, 0, organiztion);
    if applicationList is ApplicationList {
        test:assertTrue(true, "Successfully retrieved all applications");
    } else if applicationList is commons:APKError {
        test:assertFail("Error occured while retrieving all applications");
    }
}

@test:Config {dependsOn: [getApplicationListTest]}
function updateApplicationTest() {
    string[] testHosts = ["http://localhost:9090"];
    test:when(retrieveManagementServerHostsListMock).thenReturn(testHosts);
    Application payload = {name: "sampleApp", description: "sample application updated"};
    string? appId = application.applicationId;
    if appId is string {
        NotFoundError|Application|commons:APKError createdApplication = updateApplication(appId, payload, organiztion, "apkuser");
        if createdApplication is Application {
            test:assertTrue(true, "Successfully added the application");
            application.applicationId = createdApplication.applicationId;
        } else if createdApplication is commons:APKError {
            test:assertFail("Error occured while updating application");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [updateApplicationTest]}
function generateAPIKeyTest() {
    APIKeyGenerateRequest payload = {
        validityPeriod: 3600,
        additionalProperties: {}
    };
    string token = "eyJhbGciOiJSUzI1NiIsICJ0eXAiOiJKV1QiLCAia2lkIjoiZ2F0ZXdheV9jZXJ0aWZpY2F0ZV9hbGlhcyJ9." +
    "eyJpc3MiOiJodHRwczovL2FwaW0ud3NvMi5jb20vb2F1dGgyL3Rva2VuIiwgInN1YiI6IiIsICJhdWQiOiJodHRwczovL2FwaW0u" +
    "d3NvMi5jb20vb2F1dGgyL3Rva2VuIiwgImV4cCI6MTY3MTA5MDg0NCwgIm5iZiI6MTY3MTA4NzI0NCwgImlhdCI6MTY3MTA4NzI0NCw" +
    "gImp0aSI6IjAxZWQ3YzQ1LTQzMmQtMThlMC05MzBmLTUyN2I4ODM0NDM3MyIsICJrZXl0eXBlIjoiUFJPRFVDVElPTiIsICJwZXJtaXR0" +
    "ZWRSZWZlcmVyIjoiIiwgInBlcm1pdHRlZElwIjoiIiwgInRva2VuX3R5cGUiOiJBUElLZXkiLCAidGllckluZm8iOiIiLCAic3Vic2" +
    "NyaWJlZEFQSXMiOlt7Im5hbWUiOiJodHRwLWJpbi1hcGkiLCAiY29udGV4dCI6Ii9odHRwLWJpbi1hcGkvMS4wLjgiLCAidmVyc2lvbi" +
    "I6IjEuMC44IiwgInB1Ymxpc2hlciI6ImFwa3VzZXIiLCAidXVpZCI6IjRjM2NkYTRkLTI0YTQtNDI0OC1iODI0LTliMDAwM2U5YTUxMSJ9X" +
    "SwgImFwcGxpY2F0aW9uIjp7Im5hbWUiOiJTYW1wbGVBcHAxIiwgInV1aWQiOiI4Y2M1NWQ1Ny0xMmFhLTQyMzgtOGEyMS03OTMyYjZjMmJ" +
    "hYjIiLCAiaWQiOiIxIiwgInRpZXIiOiI2MFBlck1pbiIsICJ0aWVyUXVvdGFUeXBlIjoiIn19.WiuN7lT7JG-66XOY-Dzam7_QpKzhZPm" +
    "i3UaF1ri2950jR5ghthgYZIZ4WzUnhRvDmRLTwdgekEon_JRwe7bQYdHAgXB-_dJYdf-wGq9qmOKfB3T2c26ngO4Ca4PgU_lpl9xyHT" +
    "8LFIFE_GWes43CU_SVYL4X6yoSMuu2qN9VXsOlCWnK6v5xoNZjzcqr4qZtLuck3rcR70OF-yKZ1FRc-UlmDM_4nI9LTiOYvXGvJ8V" +
    "qXevZIdvWOm0qapK3hUFwK4uwxp_6qvTcwNOxGm1CoLv0JC-t2ds9EfkbI0qpWEjwJMFQqfgTLEZGT7I1m0re253Xv3Mg3I-eLtV9HcC2FA";
    test:when(generateTokenMock).thenReturn(token);
    string? appId = application.applicationId;
    if appId is string {
        APIKey|commons:APKError|NotFoundError key = generateAPIKey(payload, appId, "PRODUCTION", "apkuser", organiztion);
        if key is APIKey {
            test:assertTrue(true, "API Key Successfully Generated");
        } else {
            test:assertFail("Error occured while generating API Key");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}

@test:Config {dependsOn: [generateAPIKeyTest]}
function deleteApplicationTest() {
    string? appId = application.applicationId;
    if appId is string {
        error?|boolean status = deleteApplication(appId, organiztion);
        if status is boolean {
            test:assertTrue(true, "Successfully deleted application");
        } else if status is error {
            test:assertFail("Error occured while deleting application");
        }
    } else {
        test:assertFail("App ID isn't a string");
    }
}
