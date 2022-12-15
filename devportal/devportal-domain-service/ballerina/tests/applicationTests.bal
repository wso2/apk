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

@test:Mock { functionName: "validateApplicationUsagePolicy" }
test:MockFunction validateApplicationUsagePolicyMock = new();

@test:Mock { functionName: "getSubscriberIdDAO" }
test:MockFunction getSubscriberIdDAOMock = new();

@test:Mock { functionName: "addApplicationDAO" }
test:MockFunction addApplicationDAOMock = new();

@test:Mock { functionName: "getApplicationByIdDAO" }
test:MockFunction getApplicationByIdDAOMock = new();

@test:Mock { functionName: "getApplicationsDAO" }
test:MockFunction getApplicationsDAOMock = new();

@test:Mock { functionName: "updateApplicationDAO" }
test:MockFunction updateApplicationDAOMock = new();

@test:Mock { functionName: "deleteApplicationDAO" }
test:MockFunction deleteApplicationDAOMock = new();

@test:Mock { functionName: "generateToken" }
test:MockFunction generateTokenMock = new();

@test:Config {}
function addApplicationTest() {
    string?|Application|error application  ={name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application"};
    Application payload = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq"};
    test:when(validateApplicationUsagePolicyMock).withArguments("30PerMin", "carbon.super").thenReturn("20");
    test:when(getSubscriberIdDAOMock).withArguments("apkuser", "carbon.super").thenReturn(20);
    test:when(addApplicationDAOMock).thenReturn(application);
    if application is Application {
        test:assertEquals(addApplication(payload, "carbon.super", "apkuser"), application);
    } else if application is error {
        test:assertFail("Error occured while adding application");
    }
}

@test:Config {}
function getApplicationByIdTest(){
    Application app = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq"};
    test:when(getApplicationByIdDAOMock).thenReturn(app);
    string?|Application|error returnedResponse = getApplicationById("12sqwsq","carbon.super");
    if returnedResponse is Application {
    test:assertTrue(true, "Successfully retrieved application");
    } else if returnedResponse is  error {
        test:assertFail("Error occured while retrieving application");
    }
}

@test:Config {}
function getApplicationListTest(){
    Application[] appList = [{name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq"},
    {name:"sampleApp2",throttlingPolicy:"30PerMin",description: "sample application 2",applicationId: "heuqwsqasdsada"}];
    test:when(getApplicationsDAOMock).thenReturn(appList);
     string?|ApplicationList|error applicationList = getApplicationList("","","","",0,0,"carbon.super");
    if applicationList is ApplicationList {
    test:assertTrue(true, "Successfully retrieved all applications");
    } else if applicationList is  error {
        test:assertFail("Error occured while retrieving all applications");
    }
}

@test:Config {}
function updateApplicationTest() {
    string?|Application|error application  ={name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application"};
    Application payload = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq"};
    test:when(getApplicationByIdDAOMock).thenReturn(payload);
    test:when(validateApplicationUsagePolicyMock).withArguments("30PerMin", "carbon.super").thenReturn("20");
    test:when(getSubscriberIdDAOMock).withArguments("apkuser", "carbon.super").thenReturn(20);
    test:when(updateApplicationDAOMock).thenReturn(application);
    if application is Application {
        test:assertEquals(updateApplication("12sqwsq", payload, "carbon.super", "apkuser"), application);
    } else if application is error {
        test:assertFail("Error occured while updating application");
    }
}

@test:Config {}
function deleteApplicationTest(){
    test:when(deleteApplicationDAOMock).withArguments("12sqwsq", "carbon.super").thenReturn("");
    error?|string status = deleteApplication("12sqwsq", "carbon.super");
    if status is string {
    test:assertTrue(true, "Successfully deleted application");
    } else if status is  error {
        test:assertFail("Error occured while deleting application");
    }
}

@test:Config {}
function generateAPIKeyTest(){
    Application app = {name:"sampleApp",throttlingPolicy:"30PerMin",description: "sample application",applicationId: "12sqwsq"};
    test:when(getApplicationByIdDAOMock).thenReturn(app);
    APIKeyGenerateRequest payload = {
        validityPeriod: 3600,
        additionalProperties: {}
    };
    Subscription[] subList = [{ apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"},
    { apiId: "8e3a1ca4-b649-4e57-9a57-e43b6b545af0",applicationId: "01ed716f-9f85-1ade-b634-be97dee7ceb4",throttlingPolicy: "MySubPol4"}];
    test:when(getSubscriptionsByAPPIdDAOMock).thenReturn(subList);

    API api = {name: "MyAPI1", context: "/myapi", 'version: "1.0", provider: "apkuser", lifeCycleStatus: "PUBLISHED"};
    test:when(getAPIByIdDAOMock).thenReturn(api);

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

    APIKey|error key = generateAPIKey(payload, "01ed716f-9f85-1ade-b634-be97dee7ceb4", "PRODUCTION", "apkuser", "carbon.super");
    if key is APIKey {
    test:assertTrue(true, "API Key Successfully Generated");
    } else {
        test:assertFail("Error occured while generating API Key");
    }
}