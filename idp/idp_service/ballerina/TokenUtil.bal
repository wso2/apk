import ballerina/log;
import ballerina/regex;
import ballerina/lang.array;
import ballerina/uuid;
import ballerina/jwt;

public class TokenUtil {

    public isolated function generateToken(string? authorization, Token_body payload) returns TokenResponse|BadRequestTokenErrorResponse|UnauthorizedTokenErrorResponse {
        if (authorization is ()) || (authorization.toString().trim().length() == 0) || (!authorization.toString().startsWith("Basic ")) {
            UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Unauthorized"}};
            return unauthorized;
        } else {
            string authorizationString = authorization.toString();
            string basicEncodedValue = authorizationString.substring(6, authorizationString.length()).trim();
            if basicEncodedValue.length() == 0 {
                UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Unauthorized"}};
                return unauthorized;
            } else {
                do {
                    byte[] base64DecodedBytes = check array:fromBase64(basicEncodedValue);
                    string base64DecodedString = check 'string:fromBytes(base64DecodedBytes);
                    string[] clientIdSecretToken = regex:split(base64DecodedString.trim(), ":");
                    if clientIdSecretToken.length() != 2 {
                        UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Unauthorized"}};
                        return unauthorized;
                    }
                    DCRMClient dcrmClient = new;
                    Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError application = dcrmClient.getApplication(clientIdSecretToken[0]);
                    if application is Application {
                        if application.client_secret != clientIdSecretToken[1] {
                            UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "unauthorized_client", error_description: "Invalide Client Id/Secret"}};
                            return unauthorized;
                        }
                        string grantType = payload.grant_type;
                        string[]? grantTypes = application.grant_types;
                        if grantTypes is string[] {
                            int? indexOf = grantTypes.indexOf(grantType);
                            if indexOf is () {
                                BadRequestTokenErrorResponse tokenError = {body: {'error: "unsupported_grant_type", error_description: grantType + " not supported by application."}};
                                return tokenError;
                            }
                        }
                        if grantType == CLIENT_CREDENTIALS_GRANT_TYPE {
                            return self.handleClientCredentialsGrant(payload, application);
                        }
                    } else if application is NotFoundClientRegistrationError {
                        UnauthorizedTokenErrorResponse unauthorized = {body: {'error: "access_denied", error_description: "Invalide Client Id/Secret"}};
                        return unauthorized;
                    } else {
                        BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
                        return tokenError;
                    }
                    BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
                    return tokenError;
                }
                on fail var e {
                    log:printError("Error on decoding base64", e);
                    BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
                    return tokenError;
                }
            }
        }
    }
    public isolated function handleClientCredentialsGrant(Token_body payload, Application application) returns TokenResponse|BadRequestTokenErrorResponse|UnauthorizedTokenErrorResponse {
        string? scope = payload.scope;
        string[] scopeArray = ["default"];
        if scope is string && scope.trim().length() > 0 {
            scopeArray = regex:split(scope, " ");
        }
        string|jwt:Error tokenResult = self.issueToken(application, (), scopeArray);
        if tokenResult is string {
            TokenResponse tokenResponse = {
                access_token: tokenResult,
                token_type: "Bearer",
                expires_in: idpConfiguration.tokenIssuerConfiguration.expTime,
                scope: string:'join(" ", ...scopeArray)
            };
            return tokenResponse;
        }
        else {
            log:printError("Error on Generating token", tokenResult);
            BadRequestTokenErrorResponse tokenError = {"body": {'error: "server_error", error_description: "Server Error occured on generating token"}};
            return tokenError;

        }
    }
    public isolated function issueToken(Application application, string? username, string[] scopes) returns string|jwt:Error {
        TokenIssuerConfiguration issuerConfiguration = idpConfiguration.tokenIssuerConfiguration;
        KeyStoreConfiguration signingCert = idpConfiguration.signingKeyStore;
        string jwtid = uuid:createType1AsString();
        jwt:IssuerConfig issuerConfig = {
            issuer: issuerConfiguration.issuer,
            expTime: issuerConfiguration.expTime,
            jwtId: jwtid,
            keyId: issuerConfiguration.keyId,
            signatureConfig: {
                config: {keyFile: signingCert.path}
            }
        };
        if username is string {
            issuerConfig.username = username;
        } else {
            issuerConfig.username = application.client_id;
        }

        issuerConfig.customClaims = self.handleCustomClaims(application, username, scopes);
        return jwt:issue(issuerConfig);
    }
    public isolated function handleCustomClaims(Application application, string? username, string[] scopes) returns map<json> {
        map<json> claims = {};
        claims = {"scope": string:'join(" ", ...scopes)};
        return claims;
    }
}
