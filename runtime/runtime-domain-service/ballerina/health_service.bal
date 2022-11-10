import ballerina/http;

service / on ep0 {
    resource function get health() returns http:Ok {
        json status = {"health": "Ok"};
        return {body: status};
    }
}
