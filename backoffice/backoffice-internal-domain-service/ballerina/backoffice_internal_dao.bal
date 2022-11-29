import ballerinax/postgresql;
import ballerina/sql;


public function getConnection() returns postgresql:Client | error {
    //Todo: Need to read database config from toml
    postgresql:Client|sql:Error dbClient = 
                                check new ("localhost", "admin", "admin", 
                                     "APKDB", 5432);
    return dbClient;   
}

public function db_createAPI(API api) returns sql:ExecutionResult | sql:Error{
    postgresql:Client | error db_client  = getConnection();
    if db_client is error {
        return error("Issue while conecting to databse");
    } else {
        //Todo: query need to improve
        sql:ParameterizedQuery query = `INSERT INTO am_api(api_name, api_version,context,api_provider,status,artifact)
                                  VALUES (${api.name}, ${api.'version}, ${api.context},${api.provider},${api.lifeCycleStatus}, '{}')`;
        sql:ExecutionResult result = check db_client->execute(query);
        return result;
    }
}