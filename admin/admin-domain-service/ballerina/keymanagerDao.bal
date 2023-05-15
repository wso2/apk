import wso2/apk_common_lib as commons;
import ballerinax/postgresql;
import ballerina/sql;
import ballerina/log;

isolated function addKeyManagerEntry(KeyManagerDaoEntry keyManager, commons:Organization organization) returns commons:APKError? {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `INSERT INTO KEY_MANAGER (UUID,NAME,DISPLAY_NAME,DESCRIPTION,TYPE,ENABLED,ORGANIZATION,ISSUER,CONFIGURATION) VALUES (${keyManager.uuid},${keyManager.name},${keyManager.display_name},${keyManager.description},${keyManager.'type},${keyManager.enabled},${organization.uuid},${keyManager.issuer},${keyManager.configuration})`;
        sql:ExecutionResult|sql:Error result = dbClient->execute(query);
        if result is sql:ExecutionResult && result.affectedRowCount == 0 {
            return e909438(());
        } else if result is sql:Error {
            log:printError("Error while adding key manager entry", result);
            return e909402(result);
        }
    }
}

isolated function getAllKeyManagersByOrganization(commons:Organization organization) returns KeyManagerListingDaoEntry[]|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        do {
            sql:ParameterizedQuery query = `SELECT UUID,NAME,DESCRIPTION,TYPE,ENABLED FROM KEY_MANAGER WHERE ORGANIZATION = ${organization.uuid}`;
            stream<KeyManagerListingDaoEntry, sql:Error?> keyManagerListStream = dbClient->query(query);
            KeyManagerListingDaoEntry[] keYmanagerList = check from KeyManagerListingDaoEntry keyManager in keyManagerListStream
                select keyManager;
            check keyManagerListStream.close();
            return keYmanagerList;
        } on fail var e {
            return e909432(e);
        }
    }
}

isolated function checkKeyManagerExist(string name, commons:Organization organization) returns boolean|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `select exists(SELECT 1 FROM KEY_MANAGER WHERE ORGANIZATION = ${organization.uuid} AND NAME = ${name})`;
        boolean|sql:Error result = dbClient->queryRow(query);
        if result is boolean {
            return result;
        } else {
            log:printError("Error while checking key manager existance", result);
            return e909433(result);
        }
    }
}

isolated function getKeyManagerById(string id, commons:Organization organization) returns KeyManagerDaoEntry|commons:APKError {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909401(dbClient);
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

public isolated function updateKeyManager(string id, KeyManagerDaoEntry keyManager, commons:Organization organization) returns commons:APKError? {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `UPDATE KEY_MANAGER SET NAME = ${keyManager.name},DISPLAY_NAME = ${keyManager.display_name},DESCRIPTION = ${keyManager.description},TYPE = ${keyManager.'type},ENABLED = ${keyManager.enabled},ISSUER = ${keyManager.issuer},CONFIGURATION = ${keyManager.configuration} WHERE UUID = ${id} AND ORGANIZATION = ${organization.uuid}`;
        sql:ExecutionResult|sql:Error result = dbClient->execute(query);
        if result is sql:ExecutionResult && result.affectedRowCount == 0 {
            return e909438(());
        } else if result is sql:Error {
            log:printError("Error while updating key manager", result);
            return e909402(result);
        }
    }
}

public isolated function deleteKeyManager(string id, commons:Organization organization) returns commons:APKError? {
    postgresql:Client|error dbClient = getConnection();
    if dbClient is error {
        return e909401(dbClient);
    } else {
        sql:ParameterizedQuery query = `DELETE FROM KEY_MANAGER WHERE UUID = ${id} AND ORGANIZATION = ${organization.uuid}`;
        sql:ExecutionResult|sql:Error result = dbClient->execute(query);
        if result is sql:ExecutionResult && result.affectedRowCount == 0 {
            return e909440(id, organization.uuid, ());
        } else if result is sql:Error {
            log:printError(result.toString());
            return e909440(id, organization.uuid, result);
        }
    }
}

