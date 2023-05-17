import wso2/apk_common_lib as commons;

# Description
public type KeyManagerClient isolated object {
    public isolated function registerOauthApplication(ClientRegistrationRequest clientRegistrationRequst) returns ClientRegistrationResponse|commons:APKError;
    public isolated function retrieveOauthApplicationByClientId(string clientId) returns ClientRegistrationResponse|commons:APKError;
    public isolated function updateOauthApplicationByClientId(string clientId, ClientUpdateRequest clientUpdateRequest) returns ClientRegistrationResponse|commons:APKError;
    public isolated function deleteOauthApplication(string clientId) returns boolean|commons:APKError;
    public isolated function generateAccessToken(TokenRequest tokenRequest) returns TokenResponse|commons:APKError;
};
