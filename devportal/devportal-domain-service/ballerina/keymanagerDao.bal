import wso2/apk_common_lib as commons;
import ballerinax/postgresql;
import ballerina/sql;
import ballerina/log;

isolated function getAllKeyManagersByOrganization(commons:Organization organization) returns KeyManagerListingDaoEntry[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID,NAME,DESCRIPTION,TYPE,ENABLED FROM KEY_MANAGER WHERE ORGANIZATION = ${organization.uuid}`;
            stream<KeyManagerListingDaoEntry, sql:Error?> keyManagerListStream = dbClient->query(query);
            KeyManagerListingDaoEntry[] keYmanagerList = check from KeyManagerListingDaoEntry keyManager in keyManagerListStream
                select keyManager;
            check keyManagerListStream.close();
            return keYmanagerList;
        } on fail var e {
            return error("Internal Error occured while retrieving KeyManagers", e,
        code = 909432,
        message = "Internal Error occured while retrieving KeyManagers",
        statusCode = 500,
        description = "Internal Error occured while retrieving KeyManagers");
        }
    }
}

isolated function getKeyManagerById(string id, commons:Organization organization) returns KeyManagerDaoEntry|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `SELECT UUID,NAME,DISPLAY_NAME,ISSUER,DESCRIPTION,TYPE,CONFIGURATION,ENABLED FROM KEY_MANAGER WHERE ORGANIZATION = ${organization.uuid} AND UUID = ${id}`;
        KeyManagerDaoEntry|sql:Error result = dbClient->queryRow(query);
        if result is sql:NoRowsError {
            return e909439(id, organization.uuid);
        } else if result is KeyManagerDaoEntry {
            return result;
        } else {
            log:printError("Error while getting key manager by id", result);
            return e909433(result);
        }
    }
}



public type KeyManagerListingDaoEntry record {|
    string uuid;
    string name;
    string display_name?;
    string description?;
    string 'type;
    boolean enabled;
|};

public type KeyManagerDaoEntry record {|
    string uuid?;
    string name;
    string display_name?;
    string issuer;
    string description?;
    string 'type;
    byte[] configuration?;
    boolean enabled;
|};

isolated function e909439(string id, string organization) returns commons:APKError {
    return error commons:APKError("KeyManager from " + id + " not exist in organization " + organization,
        code = 909439,
        message = "KeyManager from " + id + " not exist in organization " + organization,
        statusCode = 404,
        description = "KeyManager from " + id + " not exist in organization " + organization
    );
}

public isolated function e909433(error e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while retrieving KeyManager", e,
        code = 909433,
        message = "Internal Error occured while retrieving KeyManager",
        statusCode = 500,
        description = "Internal Error occured while retrieving KeyManager"
    );
}

isolated function e909440(string id, string organization, error? e) returns commons:APKError {
    return error commons:APKError("Internal Error occured while deleting keymanager " + id + " from organization " + organization, e,
        code = 909440,
        message = "Internal Error occured while deleting keymanager " + id + " from organization " + organization,
        statusCode = 500,
        description = "Internal Error occured while deleting keymanager " + id + " from organization " + organization
    );
}
