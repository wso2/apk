import devportal_service.types;
import wso2/apk_common_lib as commons;
import ballerina/log;
import ballerina/http;
import ballerina/regex;
import devportal_service.kmclient;

public isolated class NonProdIdpKeyManagerClient {
    *kmclient:KeyManagerClient;
    private final Client dcrClient;
    private final http:Client tokenClient;
    public isolated function init(types:KeyManager keyManager) returns commons:APKError? {
        do {
            map<string> endpoints = keyManager.endpoints;
            string dcrEndpoint = "";
            string tokenEndpoint = "";
            string username;
            string password;
            if endpoints.hasKey("dcr_endpoint") {
                dcrEndpoint = endpoints.get("dcr_endpoint");
            } else {
                return error("DCR endpoint is not defined", code = 900960, message = "DCR endpoint is not defined", description = "DCR endpoint is not defined", statusCode = 400);
            }
            if endpoints.hasKey("token_endpoint") {
                tokenEndpoint = endpoints.get("token_endpoint");
            } else {
                return error("Token endpoint is not defined", code = 900960, message = "Token endpoint is not defined", description = "Token endpoint is not defined", statusCode = 400);
            }
            record {}? additionalProperties = keyManager.additionalProperties;
            if additionalProperties is record {} {
                if additionalProperties.hasKey("username") {
                    username = <string>additionalProperties.get("username");
                } else {
                    return error("Username is not defined", code = 900960, message = "Username is not defined", description = "Username is not defined", statusCode = 400);
                }
                if additionalProperties.hasKey("password") {
                    password = <string>additionalProperties.get("password");
                } else {
                    return error("Password is not defined", code = 900960, message = "Password is not defined", description = "Password is not defined", statusCode = 400);
                }
            } else {
                return error("Additional properties are not defined", code = 900960, message = "Additional properties are not defined", description = "Additional properties are not defined", statusCode = 400);
            }
            ConnectionConfig connectionConfig = {clientAuth: {username: username, password: password}, secureSocket: {enable: false}};
            self.dcrClient = check new (dcrEndpoint, connectionConfig);
            self.tokenClient = check new (tokenEndpoint, secureSocket = {enable: false});
        } on fail var e {
            log:printError("Error while initializing the DCR client", e);
            return error("Internal Server Error", e, code = 900901, message = "Internal Server Error", description = "Internal Server Error", statusCode = 500);
        }

    }
    public isolated function registerOauthApplication(kmclient:ClientRegistrationRequest clientRegistrationRequst) returns kmclient:ClientRegistrationResponse|commons:APKError {
        RegistrationRequest registrationRequest = {
            client_name: clientRegistrationRequst.client_name,
            grant_types: clientRegistrationRequst.grant_types,
            redirect_uris: clientRegistrationRequst.redirect_uris
        };
        Application|error applicationResult = self.dcrClient->/register.post(registrationRequest);
        if applicationResult is Application {
            kmclient:ClientRegistrationResponse clientRegistrationResponse = {
                client_id: applicationResult.client_id,
                client_secret: applicationResult.client_secret,
                grant_types: applicationResult.grant_types,
                redirect_uris: applicationResult.redirect_uris
            };
            return clientRegistrationResponse;
        } else {
            commons:APKError apkError = error("Error while registering the application", applicationResult, code = 900960, message = "Error while registering the application", description = "Error while registering the application", statusCode = 500);
            return apkError;
        }
    }

    public isolated function retrieveOauthApplicationByClientId(string clientId) returns kmclient:ClientRegistrationResponse|commons:APKError {
        Application|error retrievedApplication = self.dcrClient->/register/[clientId];
        if retrievedApplication is Application {
            kmclient:ClientRegistrationResponse clientRegistrationResponse = {
                client_id: retrievedApplication.client_id,
                client_secret: retrievedApplication.client_secret,
                grant_types: retrievedApplication.grant_types,
                redirect_uris: retrievedApplication.redirect_uris
            };
            return clientRegistrationResponse;
        } else {
            commons:APKError apkError = error("Error while retrieving the application", retrievedApplication, code = 900960, message = "Error while retrieving the application", description = "Error while retrieving the application", statusCode = 500);
            return apkError;
        }
    }

    public isolated function updateOauthApplicationByClientId(string clientId, kmclient:ClientUpdateRequest clientUpdateRequest) returns kmclient:ClientRegistrationResponse|commons:APKError {
        UpdateRequest updateRequest = {client_name: clientUpdateRequest.client_name, redirect_uris: clientUpdateRequest.redirect_uris, grant_types: clientUpdateRequest.grant_types};
        Application|error updatedApplication = self.dcrClient->/register/[clientId].put(updateRequest);
        if updatedApplication is Application {
            return {
                client_id: updatedApplication.client_id,
                client_secret: updatedApplication.client_secret,
                grant_types: updatedApplication.grant_types,
                redirect_uris: updatedApplication.redirect_uris
            };

        } else {
            commons:APKError apkError = error("Error while updating the application", updatedApplication, code = 900960, message = "Error while updating the application", description = "Error while updating the application", statusCode = 500);
            return apkError;
        }
    }

    public isolated function deleteOauthApplication(string clientId) returns boolean|commons:APKError {
        http:Response|error deletionResponse = self.dcrClient->/register/[clientId].delete;
        if deletionResponse is http:Response {
            if (deletionResponse.statusCode == 204) {
                return true;
            } else {
                commons:APKError apkError = error("Error while deleting the application", code = 900960, message = "Error while deleting the application", description = "Error while deleting the application", statusCode = 500);
                return apkError;
            }
        } else {
            commons:APKError apkError = error("Error while deleting the application", deletionResponse, code = 900960, message = "Error while deleting the application", description = "Error while deleting the application", statusCode = 500);
            return apkError;
        }
    }

    public isolated function generateAccessToken(kmclient:TokenRequest tokenRequest) returns kmclient:TokenResponse|commons:APKError {
        do {
            string resourcePath = string `/`;
            string collonSeperatedValue = <string>tokenRequest.client_id + ":" + <string>tokenRequest.client_secret;
            byte[] bytes = collonSeperatedValue.toBytes();
            string authorizationHeaderValue = bytes.toBase64();
            map<string|string[]> httpHeaders = {
                "Authorization": "Basic " + authorizationHeaderValue,
                "Content-Type": "application/x-www-form-urlencoded"
            };
            http:Request request = new;
            record {} payload = {"grant_type": "client_credentials", "scope": string:'join(" ", ...tokenRequest.scopes ?: ["default"])};
            string encodedRequestBody = createFormURLEncodedRequestBody(payload);
            request.setPayload(encodedRequestBody, "application/x-www-form-urlencoded");
            http:Response httpResponse = check self.tokenClient->post(resourcePath, request, httpHeaders);
            if httpResponse.statusCode == http:STATUS_OK {
                json jsonPayload = check httpResponse.getJsonPayload();
                kmclient:TokenResponse tokenResponse = {};
                string? accessToken = check jsonPayload?.access_token;
                string? refreshToken = check jsonPayload?.refresh_token;
                string? tokenType = check jsonPayload?.token_type;
                int? expiresIn = check jsonPayload?.expires_in;
                string? scopes = check jsonPayload?.scope;
                tokenResponse.access_token = accessToken;
                tokenResponse.refresh_token = refreshToken;
                tokenResponse.token_type = tokenType;
                tokenResponse.expires_in = expiresIn;
                if scopes is string {
                    tokenResponse.scopes = regex:split(scopes, ",");
                }
                return tokenResponse;
            } else {
                log:printError("Error while generating the access token", statusCode = httpResponse.statusCode);
                return error("Internal Server Error", code = 900901, message = "Internal Server Error", description = "Internal Server Error", statusCode = 500);
            }
        } on fail var e {
            return error("Internal Server Error", e, code = 900901, message = "Internal Server Error", description = "Internal Server Error", statusCode = 500);
        }
    }
}
