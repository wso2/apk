import ballerinax/postgresql;
import ballerina/log;
import ballerina/uuid;
import ballerina/regex;
import ballerina/sql;
import ballerina/http;

public class DCRMClient {

    public isolated function createDCRApplication(RegistrationRequest payload) returns CreatedApplication|BadRequestClientRegistrationError|ConflictClientRegistrationError|InternalServerErrorClientRegistrationError {
        if !(payload.client_name is string && payload.client_name.toString().length() > 0) {
            BadRequestClientRegistrationError badClient = {body: {'error: CLIENT_NAME_EMPTY_ERROR, error_description: "Client Name is Empty."}};
            return badClient;
        }
        string[]? grantTypes = payload.grant_types;
        if (grantTypes is () || grantTypes.length()==0) {
            BadRequestClientRegistrationError badClient = {body: {'error: GRANT_TYPES_EMPTY_ERROR, error_description: "grant type list is empty"}};
            return badClient;
        }
        BadRequestClientRegistrationError? validateGrantType = self.validateGrantTypes(grantTypes);
        if validateGrantType is BadRequestClientRegistrationError {
            return validateGrantType;
        }
        string clientID = uuid:createType1AsString();
        string clientSecret = uuid:createType1AsString();
        string clientName = <string>payload.client_name;
        string grantTypesArray = string:'join(",", ...grantTypes);
        string callBackurls = "";
        string[]? redirectUris = payload.redirect_uris;
        if redirectUris is string[] {
            callBackurls = string:'join("|", ...redirectUris);
        }
        postgresql:Client|error db_client = getConnection();
        if db_client is error {
            string message = "Error while retrieving connection";
            log:printError(message, db_client);
            InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
            return internalError;
        } else {
            sql:ParameterizedQuery values = `${clientID},
                                        ${clientSecret}, 
                                        ${clientName},
                                        ${callBackurls},
                                        ${grantTypesArray}
                                    )`;
            sql:ParameterizedQuery CREATE_OAUTH_APPLICATION_PREFIX = `INSERT INTO CONSUMER_APPS (CONSUMER_KEY,CONSUMER_SECRET,APP_NAME,CALLBACK_URL,GRANT_TYPES) VALUES (`;
            sql:ParameterizedQuery sqlQuery = sql:queryConcat(CREATE_OAUTH_APPLICATION_PREFIX, values);

            sql:ExecutionResult|sql:Error result = db_client->execute(sqlQuery);

            if result is sql:ExecutionResult {
                if result.affectedRowCount > 0 {
                    CreatedApplication createdApp = {body: {client_id: clientID, client_name: clientName, grant_types: grantTypes, client_secret: clientSecret, redirect_uris: redirectUris, client_secret_expires_at: int:MAX_VALUE}};
                    return createdApp;
                } else {
                    log:printWarn("Entry not inserted to db", sqlResult = result.toString());
                    InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
                    return internalError;
                }
            } else {
                string message = "Error while inserting data into Database";
                log:printError(message, result);
                InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
                return internalError;
            }
        }
    }
    public isolated function updateDCRApplication(string clientId, UpdateRequest payload) returns Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError|BadRequestClientRegistrationError {
        if !(payload.client_name is string && payload.client_name.toString().length() > 0) {
            BadRequestClientRegistrationError badClient = {body: {'error: CLIENT_NAME_EMPTY_ERROR, error_description: "Client Name is Empty."}};
            return badClient;
        }
        string[]? grantTypes = payload.grant_types;
        if (grantTypes is () || grantTypes.length()==0) {
            BadRequestClientRegistrationError badClient = {body: {'error: GRANT_TYPES_EMPTY_ERROR, error_description: "grant type list is empty"}};
            return badClient;
        }
        BadRequestClientRegistrationError? validateGrantType = self.validateGrantTypes(grantTypes);
        if validateGrantType is BadRequestClientRegistrationError {
            return validateGrantType;
        }
        string clientName = <string>payload.client_name;
        string grantTypesArray = string:'join(",", ...grantTypes);
        string callBackurls = "";
        string[]? redirectUris = payload.redirect_uris;
        if redirectUris is string[] {
            callBackurls = string:'join("|", ...redirectUris);
        }
        Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError application = self.getApplication(clientId);
        if application is Application {
            postgresql:Client|error db_client = getConnection();
            if db_client is error {
                string message = "Error while retrieving connection";
                log:printError(message, db_client);
                InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
                return internalError;
            } else {
                sql:ParameterizedQuery sqlQuery = `UPDATE CONSUMER_APPS SET APP_NAME=${clientName},CALLBACK_URL=${callBackurls},GRANT_TYPES=${grantTypesArray} WHERE CONSUMER_KEY=${clientId}`;

                sql:ExecutionResult|sql:Error result = db_client->execute(sqlQuery);

                if result is sql:ExecutionResult {
                    if result.affectedRowCount > 0 {
                        return self.getApplication(clientId);
                    } else {
                        NotFoundClientRegistrationError badRequest = {body: {'error: CLIENT_ID_NOT_FOUND_ERROR, error_description: clientId + " not found in system."}};
                        return badRequest;
                    }
                } else {
                    string message = "Error while inserting data into Database";
                    log:printError(message, result);
                    InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
                    return internalError;
                }
            }
        } else  {
            return application;
        }
    }
    public isolated function getApplication(string consumerKey) returns Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError {
        postgresql:Client|error db_Client = getConnection();
        if db_Client is error {
            string message = "Error while retrieving connection";
            log:printError(message, db_Client);
            InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
            return internalError;
        } else {
            sql:ParameterizedQuery query = `SELECT * FROM CONSUMER_APPS WHERE CONSUMER_KEY = ${consumerKey}`;
            stream<OauthAppSqlEntry, sql:Error?> resultStream = db_Client->query(query);
            do {
                check from OauthAppSqlEntry oauthAppEntry in resultStream
                    do {
                        string[] callBackUrls = [];
                        string callbackUrl = oauthAppEntry.callback_url;
                        if callbackUrl.trim().length() > 0 {
                            callBackUrls = regex:split(callbackUrl, "\\|");
                        }
                        string[] grantTypes = [];
                        if oauthAppEntry.grant_types.trim().length() > 0 {
                            grantTypes = regex:split(oauthAppEntry.grant_types, "\\,");
                        }
                        return {client_id: oauthAppEntry.consumer_key, client_secret: oauthAppEntry.consumer_secret, client_name: oauthAppEntry.app_name, grant_types: grantTypes, redirect_uris: callBackUrls, client_secret_expires_at: int:MAX_VALUE};
                    };
                NotFoundClientRegistrationError notFound = {body: {'error: CLIENT_ID_NOT_FOUND_ERROR, error_description: consumerKey + " not found in system."}};
                return notFound;
            } on fail var e {
                log:printError("Internal Error", e);
                InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
                return internalError;
            }
        }
    }

    public isolated function deleteApplication(string consumerKey) returns http:NoContent|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError {
        postgresql:Client|error db_Client = getConnection();
        if db_Client is error {
            string message = "Error while retrieving connection";
            log:printError(message, db_Client);
            InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
            return internalError;
        } else {
            sql:ParameterizedQuery query = ` DELETE FROM CONSUMER_APPS WHERE CONSUMER_KEY = ${consumerKey}`;
            sql:ExecutionResult|sql:Error result = db_Client->execute(query);
            if result is sql:ExecutionResult {
                if result.affectedRowCount > 0 {
                    return {};
                } else {
                    NotFoundClientRegistrationError clientNotFound = {body: {'error: CLIENT_ID_NOT_FOUND_ERROR, error_description: consumerKey + " not found in system."}};
                    return clientNotFound;
                }
            } else {
                string message = "Error while invoking Delete Application";
                log:printError(message, result);
                InternalServerErrorClientRegistrationError internalError = {body: {'error: INTERNAL_ERROR, error_description: "Internal Error"}};
                return internalError;
            }
        }
    }
    public isolated function getApplicationIncludeFileBaseApps(string clientId) returns Application|Application|NotFoundClientRegistrationError|InternalServerErrorClientRegistrationError {
        foreach FileBaseOAuthapps oauthApp in idpConfiguration.fileBaseApp {
            if oauthApp.clientId == clientId {
                return {
                    client_id: oauthApp.clientId,
                    client_secret: oauthApp.clientSecret,
                    grant_types: oauthApp.grantTypes,
                    redirect_uris: oauthApp.callbackUrls
                };
            }
        }
        return self.getApplication(clientId);
    }
    isolated function validateGrantTypes(string[] grantTypes) returns BadRequestClientRegistrationError? {
        foreach string grantType in grantTypes {
            lock {
                int? available = ALLOWED_GRANT_TYPES.indexOf(grantType);
                if available is () {
                    BadRequestClientRegistrationError badRequest = {body: {'error: UNSUPPORTED_GRANT_TYPE_ERROR, error_description: grantType + " grant type not supported."}};
                    return badRequest.cloneReadOnly();
                }
            }
        }
        return;
    }
}

type OauthAppSqlEntry record {|
    string consumer_key;
    string consumer_secret;
    string app_name;
    string callback_url;
    string grant_types;
|};
