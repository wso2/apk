import ballerina/http;

public function getMockServiceList() returns http:Response {
    http:Response mockResponse = new;

    json response = {
        "kind": "ServiceList",
        "apiVersion": "v1",
        "metadata": {
            "selfLink": "/api/v1/services",
            "resourceVersion": "39691"
        },
        "items": [
            {
                "metadata": {
                    "name": "backend",
                    "namespace": "apk",
                    "selfLink": "/api/v1/namespaces/apk/services/backend",
                    "uid": "275b00d1-722c-4df2-b65a-9b14677abe4b",
                    "resourceVersion": "1514",
                    "creationTimestamp": "2022-12-13T08:25:09Z",
                    "annotations": {
                        "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"name\":\"backend\",\"namespace\":\"apk\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"targetPort\":80}],\"selector\":{\"app\":\"httpbin\"}}}\n"
                    },
                    "managedFields": [
                        {
                            "manager": "kubectl-client-side-apply",
                            "operation": "Update",
                            "apiVersion": "v1",
                            "time": "2022-12-13T08:25:09Z",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {"f:metadata": {"f:annotations": {".": {}, "f:kubectl.kubernetes.io/last-applied-configuration": {}}}, "f:spec": {"f:ports": {".": {}, "k:{\"port\":80,\"protocol\":\"TCP\"}": {".": {}, "f:name": {}, "f:port": {}, "f:protocol": {}, "f:targetPort": {}}}, "f:selector": {".": {}, "f:app": {}}, "f:sessionAffinity": {}, "f:type": {}}}
                        }
                    ]
                },
                "spec": {
                    "ports": [
                        {
                            "name": "http",
                            "protocol": "TCP",
                            "port": 80,
                            "targetPort": 80
                        }
                    ],
                    "selector": {
                        "app": "httpbin"
                    },
                    "clusterIP": "10.98.200.176",
                    "type": "ClusterIP",
                    "sessionAffinity": "None"
                },
                "status": {
                    "loadBalancer": {

                    }
                }
            }
        ]
    };
    mockResponse.setPayload(response);
    return mockResponse;
}
// import ballerina/websocket;

// service /api/v1/watch/services on new websocket:Listener(ep1) {
//     resource function get .() returns websocket:Service|websocket:Error {
//         // Accept the WebSocket upgrade by returning a `websocket:Service`.
//         return new WatchServices();
//     }
// }

// service class WatchServices {
//     *websocket:Service;
//     remote function onMessage(websocket:Caller caller, string chatMessage) returns websocket:Error? {
//         json message = {
//             "type": "ADDED",
//             "object": {
//                 "kind": "Service",
//                 "apiVersion": "v1",
//                 "metadata": {
//                     "name": "backend",
//                     "namespace": "apk",
//                     "selfLink": "/api/v1/namespaces/apk/services/backend",
//                     "uid": "275b00d1-722c-4df2-b65a-9b14677abe4b",
//                     "resourceVersion": "1514",
//                     "creationTimestamp": "2022-12-13T08:25:09Z",
//                     "annotations": {"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"name\":\"backend\",\"namespace\":\"apk\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"targetPort\":80}],\"selector\":{\"app\":\"httpbin\"}}}\n"},
//                     "managedFields": [{"manager": "kubectl-client-side-apply", "operation": "Update", "apiVersion": "v1", "time": "2022-12-13T08:25:09Z", "fieldsType": "FieldsV1", "fieldsV1": {"f:metadata": {"f:annotations": {".": {}, "f:kubectl.kubernetes.io/last-applied-configuration": {}}}, "f:spec": {"f:ports": {".": {}, "k:{\"port\":80,\"protocol\":\"TCP\"}": {".": {}, "f:name": {}, "f:port": {}, "f:protocol": {}, "f:targetPort": {}}}, "f:selector": {".": {}, "f:app": {}}, "f:sessionAffinity": {}, "f:type": {}}}}]
//                 },
//                 "spec": {
//                     "ports": [
//                         {
//                             "name": "http",
//                             "protocol": "TCP",
//                             "port": 80,
//                             "targetPort": 80
//                         }
//                     ],
//                     "selector": {"app": "httpbin"},
//                     "clusterIP": "10.98.200.176",
//                     "type": "ClusterIP",
//                     "sessionAffinity": "None"
//                 },
//                 "status": {"loadBalancer": {}}
//             }
//         };

//         check caller->writeMessage(message);
//     }
// }
