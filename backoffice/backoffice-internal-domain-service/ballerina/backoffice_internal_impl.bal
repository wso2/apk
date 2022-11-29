import ballerina/sql;
import ballerina/io;


# The function that maps createAPI service to dao layer.
#
# + body - The `API` record type parameter.
public function createAPI(API body) {
    sql:ExecutionResult | sql:Error db = db_createAPI(body);
    io:println(db);
}
