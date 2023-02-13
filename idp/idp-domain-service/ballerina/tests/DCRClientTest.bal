import ballerina/http;
import ballerina/uuid;
import ballerina/test;

@test:Config {}
public function testCreateDCRApplication() returns error? {

    RegistrationRequest registrationRequest = {client_name: "app1", redirect_uris: ["https://localhost"], grant_types: ["client_credentials", "authorization_code", "refresh_token"]};
    DCRClient dcrClient = check new ("http://localhost:9443");
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
function testCreateApplicationNegativeTests() returns error? {
    DCRClient dcrClient = check new ("http://localhost:9443");
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
