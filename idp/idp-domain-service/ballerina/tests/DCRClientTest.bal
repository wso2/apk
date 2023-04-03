
//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
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
import ballerina/http;
import ballerina/uuid;
import ballerina/test;
import ballerina/sql;
import ballerinax/postgresql;



@test:Mock {functionName: "getConnection"}
test:MockFunction testgetConnection = new ();
ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};

@test:Config {}
public function testCreateDCRApplication() returns error? {
    test:when(testgetConnection).callOriginal();
    RegistrationRequest registrationRequest = {client_name: "app1", redirect_uris: ["https://localhost"], grant_types: ["client_credentials", "authorization_code", "refresh_token"]};
    DCRClient dcrClient = check new ("https://localhost:9443",connectionConfig);
    Application createdApplication = check dcrClient->/register.post(registrationRequest);
    test:assertEquals(createdApplication.client_name, registrationRequest.client_name);
    test:assertEquals(createdApplication.grant_types, registrationRequest.grant_types);
    test:assertEquals(createdApplication.redirect_uris, registrationRequest.redirect_uris);
    test:assertTrue(createdApplication.client_id != (), "client_id is empty");
    test:assertTrue(createdApplication.client_secret != (), "client_secret is empty");
    test:assertTrue(createdApplication.client_secret_expires_at != (), "client_secret_expires_at is empty");
    Application retrieveApplication = check dcrClient->/register/[createdApplication.client_id.toString()];
    test:assertEquals(createdApplication.client_name, retrieveApplication.client_name);
    test:assertEquals(createdApplication.grant_types, retrieveApplication.grant_types);
    test:assertEquals(createdApplication.redirect_uris, retrieveApplication.redirect_uris);
    test:assertEquals(createdApplication.client_id, retrieveApplication.client_id);
    test:assertEquals(createdApplication.client_secret, retrieveApplication.client_secret);
    test:assertEquals(createdApplication.client_secret_expires_at, retrieveApplication.client_secret_expires_at);
    UpdateRequest updateRequest = {client_name: "updated_name", grant_types: registrationRequest.grant_types, redirect_uris: registrationRequest.redirect_uris};
    Application updatedApplicationByName = check dcrClient->/register/[createdApplication.client_id.toString()].put(updateRequest);
    test:assertNotEquals(updatedApplicationByName.client_name, retrieveApplication.client_name);
    test:assertEquals(updatedApplicationByName.client_name, updateRequest.client_name);
    test:assertEquals(updatedApplicationByName.grant_types, retrieveApplication.grant_types);
    test:assertEquals(updatedApplicationByName.redirect_uris, retrieveApplication.redirect_uris);
    test:assertEquals(updatedApplicationByName.client_id, retrieveApplication.client_id);
    test:assertEquals(updatedApplicationByName.client_secret, retrieveApplication.client_secret);
    test:assertEquals(updatedApplicationByName.client_secret_expires_at, retrieveApplication.client_secret_expires_at);
    updateRequest = {client_name: updatedApplicationByName.client_name, grant_types: ["client_credentials"], redirect_uris: registrationRequest.redirect_uris};
    Application updatedApplicationByGrantTypes = check dcrClient->/register/[createdApplication.client_id.toString()].put(updateRequest);
    test:assertEquals(updatedApplicationByGrantTypes.client_name, updateRequest.client_name);
    test:assertNotEquals(updatedApplicationByGrantTypes.grant_types, retrieveApplication.grant_types);
    test:assertEquals(updatedApplicationByGrantTypes.grant_types, updateRequest.grant_types);
    test:assertEquals(updatedApplicationByGrantTypes.redirect_uris, retrieveApplication.redirect_uris);
    test:assertEquals(updatedApplicationByGrantTypes.client_id, retrieveApplication.client_id);
    test:assertEquals(updatedApplicationByGrantTypes.client_secret, retrieveApplication.client_secret);
    test:assertEquals(updatedApplicationByGrantTypes.client_secret_expires_at, retrieveApplication.client_secret_expires_at);
    updateRequest = {client_name: updatedApplicationByName.client_name, grant_types: ["client_credentials"], redirect_uris: ["https://httpbin.org"]};
    Application updateApplicationByRedirectUri = check dcrClient->/register/[createdApplication.client_id.toString()].put(updateRequest);
    test:assertEquals(updateApplicationByRedirectUri.client_name, updateRequest.client_name);
    test:assertEquals(updateApplicationByRedirectUri.grant_types, updateRequest.grant_types);
    test:assertNotEquals(updateApplicationByRedirectUri.redirect_uris, retrieveApplication.redirect_uris);
    test:assertEquals(updateApplicationByRedirectUri.redirect_uris, updateRequest.redirect_uris);
    test:assertEquals(updateApplicationByRedirectUri.client_id, retrieveApplication.client_id);
    test:assertEquals(updateApplicationByRedirectUri.client_secret, retrieveApplication.client_secret);
    test:assertEquals(updateApplicationByRedirectUri.client_secret_expires_at, retrieveApplication.client_secret_expires_at);
    UpdateRequest invalidGrantTypes = {client_name: updatedApplicationByName.client_name, grant_types: ["client_credentials1"], redirect_uris: ["https://httpbin.org"]};
    Application|error invalidGrantTypeResponse = dcrClient->/register/[createdApplication.client_id.toString()].put(invalidGrantTypes);
    test:assertTrue(invalidGrantTypeResponse is error);
    if invalidGrantTypeResponse is error {
        test:assertTrue(invalidGrantTypeResponse.toString().includes("client_credentials1"));
    }
    invalidGrantTypes = {client_name: updatedApplicationByName.client_name, grant_types: [], redirect_uris: ["https://httpbin.org"]};
    invalidGrantTypeResponse = dcrClient->/register/[createdApplication.client_id.toString()].put(invalidGrantTypes);
    test:assertTrue(invalidGrantTypeResponse is error);
    if invalidGrantTypeResponse is error {
        test:assertTrue(invalidGrantTypeResponse.toString().includes("100151"));
    }
    UpdateRequest emptyClientName = {client_name: "", grant_types: [], redirect_uris: ["https://httpbin.org"]};
    invalidGrantTypeResponse = dcrClient->/register/[createdApplication.client_id.toString()].put(emptyClientName);
    test:assertTrue(invalidGrantTypeResponse is error);
    if invalidGrantTypeResponse is error {
        test:assertTrue(invalidGrantTypeResponse.toString().includes("100150"));
    }
    updateRequest = {client_name: updatedApplicationByName.client_name, grant_types: ["client_credentials"], redirect_uris: registrationRequest.redirect_uris};
    Application|error updateinvalidClient = dcrClient->/register/[uuid:createType1AsString()].put(updateRequest);
    test:assertTrue(updateinvalidClient is error);
    if updateinvalidClient is error {
        test:assertTrue(updateinvalidClient.toString().includes(""));
    }
    http:Response deletionResponse = check dcrClient->/register/[createdApplication.client_id.toString()].delete;
    test:assertEquals(deletionResponse.statusCode, 204);
    Application|error applicationNotFound = dcrClient->/register/[createdApplication.client_id.toString()];
    if applicationNotFound is Application {
        test:assertFail("Application found");
    }
}

@test:Config {}
function testCreateApplicationNegativeTests1() returns error? {
    test:when(testgetConnection).callOriginal();
    DCRClient dcrClient = check new ("https://localhost:9443", connectionConfig);
    RegistrationRequest registrationRequest = {client_name: "", redirect_uris: ["https://localhost"], grant_types: ["client_credentials", "authorization_code", "refresh_token"]};
    Application|error createdApplication = dcrClient->/register.post(registrationRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100150", 0));
    }
    registrationRequest = {client_name: "app1", redirect_uris: ["https://localhost"], grant_types: ["client_credentials1"]};
    createdApplication = dcrClient->/register.post(registrationRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100153", 0));
    }
    registrationRequest = {client_name: "app1", redirect_uris: ["https://localhost"], grant_types: []};
    createdApplication = dcrClient->/register.post(registrationRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100151", 0));
    }
}

@test:Config {}
function testCreateApplicationNegativeTests2() returns error? {
    sql:DatabaseError error1 = error("Error while connecting to db", errorCode = 0, sqlState = ());
    test:when(testgetConnection).thenReturn(error1);
    ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    DCRClient dcrClient = check new ("https://localhost:9443", connectionConfig);
    RegistrationRequest registrationRequest = {client_name: "abcde", redirect_uris: ["https://localhost"], grant_types: ["client_credentials", "authorization_code", "refresh_token"]};
    Application|error createdApplication = dcrClient->/register.post(registrationRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
}

@test:Config {}
function testCreateApplicationNegativeTests3() returns error? {
    postgresql:Client mockClient = test:mock(postgresql:Client);
    test:when(testgetConnection).thenReturn(mockClient);
    sql:DatabaseError dataBaseError = error("error while executing query", errorCode = 90100, sqlState = ());
    sql:ExecutionResult execution = {lastInsertId: 1, affectedRowCount: 0};

    test:prepare(mockClient).when("execute").thenReturnSequence(dataBaseError, execution);
    ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    DCRClient dcrClient = check new ("https://localhost:9443", connectionConfig);
    RegistrationRequest registrationRequest = {client_name: "abcde", redirect_uris: ["https://localhost"], grant_types: ["client_credentials", "authorization_code", "refresh_token"]};
    Application|error createdApplication = dcrClient->/register.post(registrationRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
    createdApplication = dcrClient->/register.post(registrationRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
}

@test:Config {}
function testUpdateApplicationNegativeTests4() returns error? {
    string clientId = uuid:createType1AsString();
    sql:DatabaseError error1 = error("Error while connecting to db", errorCode = 0, sqlState = ());
    postgresql:Client mockClient = test:mock(postgresql:Client);
    test:when(testgetConnection).thenReturn(error1);
    ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    DCRClient dcrClient = check new ("https://localhost:9443", connectionConfig);
    UpdateRequest updateRequest = {client_name: "abcde", redirect_uris: ["https://localhost"], grant_types: ["client_credentials", "authorization_code", "refresh_token"]};
    Application|error createdApplication = dcrClient->/register/[clientId].put(updateRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
    sql:DatabaseError dataBaseError = error("error while executing query", errorCode = 90100, sqlState = ());
    test:when(testgetConnection).thenReturn(mockClient);
    sql:ExecutionResult execution = {lastInsertId: 1, affectedRowCount: 0};
    OauthAppSqlEntry oauthAppentry = {consumer_key: clientId, callback_url: "", app_name: "", consumer_secret: uuid:createType1AsString(), grant_types: ""};
    test:prepare(mockClient).when("queryRow").thenReturn(oauthAppentry);
    test:prepare(mockClient).when("execute").thenReturnSequence(dataBaseError, execution);
    createdApplication = dcrClient->/register/[clientId].put(updateRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
    createdApplication = dcrClient->/register/[clientId].put(updateRequest);
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100152", 0));
    }
}

@test:Config {}
function testGetApplicationNegativeTests5() returns error? {
    string clientId = uuid:createType1AsString();
    sql:DatabaseError error1 = error("Error while connecting to db", errorCode = 0, sqlState = ());
    postgresql:Client mockClient = test:mock(postgresql:Client);
    test:when(testgetConnection).thenReturn(error1);
    ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    DCRClient dcrClient = check new ("https://localhost:9443", connectionConfig);
    Application|error createdApplication = dcrClient->/register/[clientId];
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
    sql:DatabaseError dataBaseError = error("error while executing query", errorCode = 90100, sqlState = ());
    test:when(testgetConnection).thenReturn(mockClient);
    test:prepare(mockClient).when("queryRow").thenReturn(dataBaseError);
    createdApplication = dcrClient->/register/[clientId];
    test:assertTrue(createdApplication is error);
    if createdApplication is error {
        test:assertTrue(createdApplication.toString().includes("100100", 0));
    }
}

@test:Config {}
function testDeleteApplicationNegativeTests6() returns error? {
    string clientId = uuid:createType1AsString();
    sql:DatabaseError error1 = error("Error while connecting to db", errorCode = 0, sqlState = ());
    postgresql:Client mockClient = test:mock(postgresql:Client);
    test:when(testgetConnection).thenReturn(error1);
    ConnectionConfig connectionConfig = {secureSocket: {enable: true, cert: "tests/resources/wso2carbon.crt"}};
    DCRClient dcrClient = check new ("https://localhost:9443", connectionConfig);
    http:Response|error createdApplication = dcrClient->/register/[clientId].delete;
    test:assertTrue(createdApplication is http:Response);
    if createdApplication is http:Response {
        test:assertEquals(createdApplication.statusCode, 500);
        test:assertEquals(check createdApplication.getJsonPayload(), {'error: INTERNAL_ERROR, error_description: "Internal Error"});
    }
    sql:DatabaseError dataBaseError = error("error while executing query", errorCode = 90100, sqlState = ());
    sql:ExecutionResult executionResult = {affectedRowCount: 0, lastInsertId: 1};
    test:when(testgetConnection).thenReturn(mockClient);
    test:prepare(mockClient).when("execute").thenReturnSequence(dataBaseError, executionResult);
    createdApplication = dcrClient->/register/[clientId].delete;
    test:assertTrue(createdApplication is http:Response);
    if createdApplication is http:Response {
        test:assertEquals(createdApplication.statusCode, 500);
        test:assertEquals(check createdApplication.getJsonPayload(), {'error: INTERNAL_ERROR, error_description: "Internal Error"});
    }
    createdApplication = dcrClient->/register/[clientId].delete;
    test:assertTrue(createdApplication is http:Response);
    if createdApplication is http:Response {
        test:assertEquals(createdApplication.statusCode, 404);
        test:assertEquals(check createdApplication.getJsonPayload(), {'error: CLIENT_ID_NOT_FOUND_ERROR, error_description: clientId + " not found in system."});
    }

}
