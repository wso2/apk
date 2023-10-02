import wso2/apk_common_lib as commons;
import ballerinax/postgresql;
import ballerina/sql;
import devportal_service.types;
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

public isolated function addKeyMappingEntryForApplication(types:KeyMappingDaoEntry keyMappingDaoEntry) returns commons:APKError? {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        // Insert into SUBSCRIPTION table
        sql:ParameterizedQuery query = `INSERT INTO APPLICATION_KEY_MAPPING (UUID,APPLICATION_UUID,CONSUMER_KEY,
        KEY_TYPE,CREATE_MODE, APP_INFO,KEY_MANAGER_UUID) VALUES (${keyMappingDaoEntry.uuid},${keyMappingDaoEntry.application_uuid},${keyMappingDaoEntry.consumer_key},
        ${keyMappingDaoEntry.key_type},${keyMappingDaoEntry.create_mode},${keyMappingDaoEntry.app_info},${keyMappingDaoEntry.key_manager_uuid})`;
        sql:ExecutionResult|sql:Error result = dbClient->execute(query);
        if result is sql:ExecutionResult {
            if result.affectedRowCount == 0 {
                return error("Error while inserting data into Database", message = "Error while inserting data into Database", description = "Error while inserting data into Database", code = 900954, statusCode = 500);
            } else {
                return;
            }
        } else {
            log:printDebug(result.toString());
            string message = "Error while inserting data into Database";
            return error(message, result, message = message, description = message, code = 909000, statusCode = 500);
        }
    }
}

public isolated function getKeyMappingEntriesByApplication(string applicationId) returns types:KeyMappingDaoEntry[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID,APPLICATION_UUID,CONSUMER_KEY,KEY_TYPE,CREATE_MODE,APP_INFO,KEY_MANAGER_UUID FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_UUID = ${applicationId}`;
            stream<types:KeyMappingDaoEntry, sql:Error?> keyMappingListStream = dbClient->query(query);
            types:KeyMappingDaoEntry[] keyMappingList = check from types:KeyMappingDaoEntry keyMapping in keyMappingListStream
                select keyMapping;
            check keyMappingListStream.close();
            return keyMappingList;
        } on fail var e {
            return error("Internal Error occured while retrieving KeyMapping", e, message = "Internal Error occured while retrieving KeyMapping", description = "Internal Error occured while retrieving KeyMapping", code = 909432, statusCode = 500);
        }
    }
}

public isolated function getKeyMappingEntryByApplicationAndKeyMappingId(string applicationId, string keyMappingId) returns types:KeyMappingDaoEntry|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID,APPLICATION_UUID,CONSUMER_KEY,KEY_TYPE,CREATE_MODE,APP_INFO,KEY_MANAGER_UUID FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_UUID = ${applicationId} AND UUID = ${keyMappingId}`;
            types:KeyMappingDaoEntry|sql:Error keyMappingEntry = dbClient->queryRow(query);
            if keyMappingEntry is types:KeyMappingDaoEntry {
                return keyMappingEntry;
            } else if keyMappingEntry is sql:NoRowsError {
                return error("KeyMapping from " + keyMappingId + " not exist in application " + applicationId,
                    code = 909439,
                    message = "KeyMapping from " + keyMappingId + " not exist in application " + applicationId,
                    statusCode = 404,
                    description = "KeyMapping from " + keyMappingId + " not exist in application " + applicationId);
            } else {
                return error("Internal Error occured while retrieving KeyMapping", keyMappingEntry, message = "Internal Error occured while retrieving KeyMapping", description = "Internal Error occured while retrieving KeyMapping", code = 909432, statusCode = 500);
            }
        } on fail var e {
            return error("Internal Error occured while retrieving KeyMapping", e, message = "Internal Error occured while retrieving KeyMapping", description = "Internal Error occured while retrieving KeyMapping", code = 909432, statusCode = 500);
        }
    }
}

public isolated function updateKeyMappingEntry(types:KeyMappingDaoEntry keyMappingentry) returns commons:APKError? {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        do {
            sql:ParameterizedQuery query = `UPDATE APPLICATION_KEY_MAPPING SET APP_INFO = ${keyMappingentry.app_info} WHERE APPLICATION_UUID = ${keyMappingentry.application_uuid} AND UUID = ${keyMappingentry.uuid}`;
            sql:ExecutionResult result = check dbClient->execute(query);
            if result.affectedRowCount == 0 {
                return error("KeyMapping from " + keyMappingentry.uuid + " not exist in application " + keyMappingentry.application_uuid,
                    code = 909439,
                    message = "KeyMapping from " + keyMappingentry.uuid + " not exist in application " + keyMappingentry.application_uuid,
                    statusCode = 404,
                    description = "KeyMapping from " + keyMappingentry.uuid + " not exist in application " + keyMappingentry.application_uuid);
            }
        } on fail var e {
            return error("Internal Error occured while retrieving KeyMapping", e, message = "Internal Error occured while retrieving KeyMapping", description = "Internal Error occured while retrieving KeyMapping", code = 909432, statusCode = 500);
        }
    }
}

public isolated function isKeyMappingEntryByApplicationAndKeyManagerExist(string applicationId, string keyManagerId,string keyType) returns boolean|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        string message = "Error while retrieving connection";
        return error(message, dbClient, message = message, description = message, code = 909000, statusCode = 500);
    } else {
        sql:ParameterizedQuery query = `select exists(SELECT 1 FROM APPLICATION_KEY_MAPPING WHERE APPLICATION_UUID = ${applicationId} AND KEY_MANAGER_UUID = ${keyManagerId} AND KEY_TYPE = ${keyType})`;
        boolean|sql:Error result = dbClient->queryRow(query);
        if result is boolean {
            return result;
        } else {
            return error("Internal Error occured while retrieving KeyMapping", result, message = "Internal Error occured while retrieving KeyMapping", description = "Internal Error occured while retrieving KeyMapping", code = 909432, statusCode = 500);
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
