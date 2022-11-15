import ballerina/grpc;


listener grpc:Listener ep = new (9090);

@display {
    label: "management-server",
    id: "management-server"
}

@grpc:Descriptor {value: APIDS_DESC}
service "APIService" on ep {

    remote function createAPI(API value) returns Response|error {
    }
    remote function updateAPI(API value) returns Response|error {
    }
    remote function deleteAPI(API value) returns Response|error {
    }
}

