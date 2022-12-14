import ballerina/http;

public function getMockServiceMappings() returns http:Response {
    http:Response mockResponse = new ();

    json response = {
        "apiVersion": "dp.wso2.com/v1alpha1",
        "items": [
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "ServiceMapping",
                "metadata": {
                    "creationTimestamp": "2022-12-13T17:09:49Z",
                    "generation": 1,
                    "managedFields": [
                        {
                            "apiVersion": "dp.wso2.com/v1alpha1",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {
                                "f:spec": {
                                    ".": {},
                                    "f:apiRef": {
                                        ".": {},
                                        "f:name": {},
                                        "f:namespace": {}
                                    },
                                    "f:serviceRef": {
                                        ".": {},
                                        "f:name": {},
                                        "f:namespace": {}
                                    }
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T17:09:49Z"
                        }
                    ],
                    "name": "01ed7b08-f2b1-1166-82d5-649ae706d29d-servicemapping",
                    "namespace": "apk-platform",
                    "resourceVersion": "23555",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings/01ed7b08-f2b1-1166-82d5-649ae706d29d-servicemapping",
                    "uid": "f074fb18-7924-4ddb-a65d-0ac2aa8f4953"
                },
                "spec": {
                    "apiRef": {
                        "name": "01ed7b08-f2b1-1166-82d5-649ae706d29d",
                        "namespace": "apk-platform"
                    },
                    "serviceRef": {
                        "name": "backend",
                        "namespace": "apk"
                    }
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "ServiceMapping",
                "metadata": {
                    "creationTimestamp": "2022-12-13T18:07:58Z",
                    "generation": 1,
                    "managedFields": [
                        {
                            "apiVersion": "dp.wso2.com/v1alpha1",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {
                                "f:spec": {
                                    ".": {},
                                    "f:apiRef": {
                                        ".": {},
                                        "f:name": {},
                                        "f:namespace": {}
                                    },
                                    "f:serviceRef": {
                                        ".": {},
                                        "f:name": {},
                                        "f:namespace": {}
                                    }
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T18:07:58Z"
                        }
                    ],
                    "name": "01ed7b11-0b25-12ee-927c-cd10449788c2-servicemapping",
                    "namespace": "apk-platform",
                    "resourceVersion": "26482",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings/01ed7b11-0b25-12ee-927c-cd10449788c2-servicemapping",
                    "uid": "47bfb27d-0ea5-44eb-9cfb-b319d353a7fc"
                },
                "spec": {
                    "apiRef": {
                        "name": "01ed7b11-0b25-12ee-927c-cd10449788c2",
                        "namespace": "apk-platform"
                    },
                    "serviceRef": {
                        "name": "backend",
                        "namespace": "apk"
                    }
                }
            }
        ],
        "kind": "ServiceMappingList",
        "metadata": {
            "continue": "",
            "resourceVersion": "39433",
            "selfLink": "/apis/dp.wso2.com/v1alpha1/servicemappings"
        }
    };
    mockResponse.setPayload(response);
    return mockResponse;
}
// import ballerina/websocket;

// service /apis/'dp\.wso2\.com/v1alpha1/watch/servicemappings on new websocket:Listener(ep1) {
//     resource function get .() returns websocket:Service|websocket:Error {
//         // Accept the WebSocket upgrade by returning a `websocket:Service`.
//         return new WatchServiceMappins();
//     }
// }

// service class WatchServiceMappins {
//     *websocket:Service;
//     remote function onMessage(websocket:Caller caller, string chatMessage) returns websocket:Error? {
//         json message = {
//             "type": "ADDED",
//             "object": {
//                 "apiVersion": "dp.wso2.com/v1alpha1",
//                 "kind": "ServiceMapping",
//                 "metadata": {
//                     "creationTimestamp": "2022-12-13T09:45:47Z",
//                     "generation": 1,
//                     "managedFields": [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsType": "FieldsV1", "fieldsV1": {"f:spec": {".": {}, "f:apiRef": {".": {}, "f:name": {}, "f:namespace": {}}, "f:serviceRef": {".": {}, "f:name": {}, "f:namespace": {}}}}, "manager": "ballerina", "operation": "Update", "time": "2022-12-13T09:45:47Z"}],
//                     "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a-servicemapping",
//                     "namespace": "apk-platform",
//                     "resourceVersion": "5834",
//                     "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings/01ed7aca-eb6b-1178-a200-f604a4ce114a-servicemapping",
//                     "uid": "79d280d6-df31-4017-911f-3229955fdb55"
//                 },
//                 "spec": {
//                     "apiRef": {
//                         "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a",
//                         "namespace": "apk-platform"
//                     },
//                     "serviceRef": {
//                         "name": "backend",
//                         "namespace": "apk"
//                     }
//                 }
//             }
//         };

//         check caller->writeMessage(message);
//     }
// }
