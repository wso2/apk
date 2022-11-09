import admin_service.admin as _;
import ballerina/http;

service on new http:Listener(8080) {
    resource function get health (http:Caller caller, http:Request req)
            returns error? {
        check caller->respond("Domain service started!");
    }
}
