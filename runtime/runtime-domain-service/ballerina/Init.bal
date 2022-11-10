import ballerina/log;
import ballerina/http;

listener http:Listener ep0 = new (9443);

function init() {
    log:printInfo("Initializing Runtime Domain Service..");
}

