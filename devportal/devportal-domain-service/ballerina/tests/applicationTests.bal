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