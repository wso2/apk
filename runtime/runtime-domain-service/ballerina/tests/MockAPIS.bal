//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
import ballerina/http;
import runtime_domain_service.model;

public function getMockAPIList() returns model:APIList {

    model:APIList response = {
        "apiVersion": "dp.wso2.com/v1alpha1",
        "items": [
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "API",
                "metadata": {
                    "creationTimestamp": "2022-12-13T09:45:47Z",
                    "generation": 1,
                    "managedFields": [
                        {
                            "apiVersion": "dp.wso2.com/v1alpha1",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {
                                "f:spec": {
                                    ".": {},
                                    "f:apiDisplayName": {},
                                    "f:apiType": {},
                                    "f:apiVersion": {},
                                    "f:context": {},
                                    "f:definitionFileRef": {},
                                    "f:prodHTTPRouteRef": {}
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T09:45:47Z"
                        }
                    ],
                    "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a",
                    "namespace": "apk-platform",
                    "resourceVersion": "5833",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7aca-eb6b-1178-a200-f604a4ce114a",
                    "uid": "c5ab2423-b9e8-432b-92e8-35e6907ed5e8"
                },
                "spec": {
                    "apiDisplayName": "pizzashackAPI",
                    "apiType": "REST",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack/1.0.0",
                    "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                    "definitionFileRef": "01ed7aca-eb6b-1178-a200-f604a4ce114a-definition",
                    "prodHTTPRouteRef": "01ed7aca-eb6b-1178-a200-f604a4ce114a-production"
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "API",
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
                                    "f:apiDisplayName": {},
                                    "f:apiType": {},
                                    "f:apiVersion": {},
                                    "f:context": {},
                                    "f:definitionFileRef": {},
                                    "f:prodHTTPRouteRef": {}
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T17:09:49Z"
                        }
                    ],
                    "name": "01ed7b08-f2b1-1166-82d5-649ae706d29d",
                    "namespace": "apk-platform",
                    "resourceVersion": "23554",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b08-f2b1-1166-82d5-649ae706d29d",
                    "uid": "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1"
                },
                "spec": {
                    "apiDisplayName": "pizzashackAPI1",
                    "apiType": "REST",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack1/1.0.0",
                    "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                    "definitionFileRef": "01ed7b08-f2b1-1166-82d5-649ae706d29d-definition",
                    "prodHTTPRouteRef": "01ed7b08-f2b1-1166-82d5-649ae706d29d-production"
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "API",
                "metadata": {
                    "creationTimestamp": "2022-12-13T09:45:47Z",
                    "generation": 1,
                    "managedFields": [
                        {
                            "apiVersion": "dp.wso2.com/v1alpha1",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {
                                "f:spec": {
                                    ".": {},
                                    "f:apiDisplayName": {},
                                    "f:apiType": {},
                                    "f:apiVersion": {},
                                    "f:context": {},
                                    "f:definitionFileRef": {},
                                    "f:prodHTTPRouteRef": {}
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T09:45:47Z"
                        }
                    ],
                    "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a",
                    "namespace": "apk-platform",
                    "resourceVersion": "5833",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                    "uid": "c5ab2423-b9e8-432b-92e8-35e6907ed5f9"
                },
                "spec": {
                    "apiDisplayName": "pizzashackAPI11",
                    "apiType": "REST",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack11/1.0.0",
                    "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                    "prodHTTPRouteRef": "01ed7aca-eb6b-1178-a200-f604a4ce114a-production"
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "API",
                "metadata": {
                    "creationTimestamp": "2022-12-13T09:45:47Z",
                    "generation": 1,
                    "managedFields": [
                        {
                            "apiVersion": "dp.wso2.com/v1alpha1",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {
                                "f:spec": {
                                    ".": {},
                                    "f:apiDisplayName": {},
                                    "f:apiType": {},
                                    "f:apiVersion": {},
                                    "f:context": {},
                                    "f:definitionFileRef": {},
                                    "f:prodHTTPRouteRef": {}
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T09:45:47Z"
                        }
                    ],
                    "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a",
                    "namespace": "apk-platform",
                    "resourceVersion": "5833",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                    "uid": "c5ab2423-b9e8-432b-92e8-35e6907ed5f1"
                },
                "spec": {
                    "apiDisplayName": "pizzashackAPI12",
                    "apiType": "REST",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack12/1.0.0",
                    "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                    "prodHTTPRouteRef": "01ed7aca-eb6b-1178-a200-f604a4ce114a-production",
                    "definitionFileRef": ""
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "API",
                "metadata": {
                    "creationTimestamp": "2022-12-13T09:45:47Z",
                    "generation": 1,
                    "managedFields": [
                        {
                            "apiVersion": "dp.wso2.com/v1alpha1",
                            "fieldsType": "FieldsV1",
                            "fieldsV1": {
                                "f:spec": {
                                    ".": {},
                                    "f:apiDisplayName": {},
                                    "f:apiType": {},
                                    "f:apiVersion": {},
                                    "f:context": {},
                                    "f:definitionFileRef": {},
                                    "f:prodHTTPRouteRef": {}
                                }
                            },
                            "manager": "ballerina",
                            "operation": "Update",
                            "time": "2022-12-13T09:45:47Z"
                        }
                    ],
                    "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a",
                    "namespace": "apk-platform",
                    "resourceVersion": "5833",
                    "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
                    "uid": "c5ab2423-b9e8-432b-92e8-35e6907ed5f3"
                },
                "spec": {
                    "apiDisplayName": "pizzashackAPI13",
                    "apiType": "REST",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack13/1.0.0",
                    "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                    "prodHTTPRouteRef": "01ed7aca-eb6b-1178-a200-f604a4ce114a-production",
                    "definitionFileRef": "01ed7aca-eb6b-1178-a200-f604a4ce114b-definition"
                }
            }
        ],
        "kind": "APIList",
        "metadata": {
            "continue": "",
            "resourceVersion": "40316",
            "selfLink": "/apis/dp.wso2.com/v1alpha1/apis"
        }
    };
    return response;
}

public function getMockWatchAPIEvent() returns string {
    json message = {
        "type":
        "ADDED",
        "object": {
            "apiVersion": "dp.wso2.com/v1alpha1",
            "kind": "API",
            "metadata": {
                "creationTimestamp": "2022-12-13T18:51:26Z",
                "generation": 1,
                "managedFields": [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsType": "FieldsV1", "fieldsV1": {"f:spec": {".": {}, "f:apiDisplayName": {}, "f:apiType": {}, "f:apiVersion": {}, "f:context": {}, "f:definitionFileRef": {}, "f:prodHTTPRouteRef": {}}}, "manager": "ballerina", "operation": "Update", "time": "2022-12-13T18:51:26Z"}],
                "name": "01ed7b16-90f7-1a88-8113-a7e71796d460",
                "namespace": "apk-platform",
                "resourceVersion": "28702",
                "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b16-90f7-1a88-8113-a7e71796d460",
                "uid": "8a1eb4f9-efab-4682-a051-4df4050812d2"
            },
            "spec": {
                "apiDisplayName": "pizzashackAPI6",
                "apiType": "REST",
                "apiVersion": "1.0.0",
                "context": "/pizzashack6/1.0.0",
                "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                "definitionFileRef": "01ed7b16-90f7-1a88-8113-a7e71796d460-definition",
                "prodHTTPRouteRef": "01ed7b16-90f7-1a88-8113-a7e71796d460-production"
            }
        }
    };
    return message.toString();

}

public function getNextMockWatchAPIEvent() returns string {
    json message = {
        "type":
        "ADDED",
        "object": {
            "apiVersion": "dp.wso2.com/v1alpha1",
            "kind": "API",
            "metadata": {
                "creationTimestamp": "2022-12-14T18:51:26Z",
                "generation": 1,
                "managedFields": [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsV1": {"f:spec": {".": {}, "f:apiDisplayName": {}, "f:apiType": {}, "f:apiVersion": {}, "f:context": {}, "f:definitionFileRef": {}, "f:prodHTTPRouteRef": {}}}, "manager": "ballerina", "time": "2022-12-13T09:45:47Z", "operation": "Update", "fieldsType": "FieldsV1"}],
                "name": "01ed7b16-90f7-1a88-8114-a7e71796d460",
                "namespace": "apk-platform",
                "resourceVersion": "28712",
                "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7b16-90f7-1a88-8114-a7e71796d460",
                "uid": "8a1eb4f9-efab-4682-a051-4df4050812d2"
            },
            "spec": {
                "apiDisplayName": "DemoAPI",
                "apiType": "REST",
                "apiVersion": "1.0.0",
                "context": "/demoapi/1.0.0",
                "definitionFileRef": "01ed7b16-90f7-1a88-8114-a7e71796d460-definition",
                "prodHTTPRouteRef": "01ed7b16-90f7-1a88-8114-a7e71796d460-production",
                "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b"
            }
        }
    };
    return message.toString();

}

public function getMockPizzaShakK8sAPI() returns model:API & readonly {
    model:API k8sAPI = {
        metadata: {
            name: "01ed7aca-eb6b-1178-a200-f604a4ce114a",
            namespace: "apk-platform",
            uid: "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
            creationTimestamp: "2022-12-13T09:45:47Z",
            generation: 1,
            selfLink: "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7aca-eb6b-1178-a200-f604a4ce114a",
            resourceVersion: "5833",
            managedFields: [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsV1": {"f:spec": {".": {}, "f:apiDisplayName": {}, "f:apiType": {}, "f:apiVersion": {}, "f:context": {}, "f:definitionFileRef": {}, "f:prodHTTPRouteRef": {}}}, "manager": "ballerina", "time": "2022-12-13T09:45:47Z", "operation": "Update", "fieldsType": "FieldsV1"}]
        },
        spec: {
            apiDisplayName: "pizzashackAPI",
            apiType: "REST",
            apiVersion: "1.0.0",
            context: "/pizzashack/1.0.0",
            definitionFileRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-definition",
            prodHTTPRouteRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-production",
            organization: "01ed7aca-eb6b-1178-a200-f604a4ce114b"
        }
    };
    return k8sAPI.cloneReadOnly();
}

public function mock404Response() returns http:Response {
    http:Response response = new;
    response.statusCode = 404;
    return response;
}

public function mockConfigMaps() returns http:Response {
    http:Response response = new;
    json configmap = {
        "apiVersion": "v1",
        "data": {
            "openapi.json": "{\"openapi\":\"3.0.1\", \"info\":{\"title\":\"pizza1234567\", \"version\":\"1.0.0\"}, \"security\":[{\"default\":[]}], \"paths\":{\"/menu\":{\"get\":{\"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-throttling-tier\":\"Unlimited\"}}, \"/order/{orderId}\":{\"post\":{\"parameters\":[{\"name\":\"orderId\", \"in\":\"path\", \"required\":true, \"schema\":{\"type\":\"string\"}}], \"responses\":{\"200\":{\"description\":\"OK\"}}, \"security\":[{\"default\":[]}], \"x-throttling-tier\":\"Unlimited\"}}}, \"components\":{\"securitySchemes\":{\"default\":{\"type\":\"oauth2\", \"flows\":{\"implicit\":{\"authorizationUrl\":\"https://test.com\", \"scopes\":{}}}}}}}"
        },
        "kind": "ConfigMap",
        "metadata": {
            "creationTimestamp": "2023-01-05T05:34:44Z",
            "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a-definition",
            "namespace": "apk-platform",
            "resourceVersion": "113573",
            "uid": "ce2915d6-cdeb-4c70-8cdd-e03c158105ba"
        }
    };
    response.setJsonPayload(configmap);
    response.statusCode = 200;
    return response;
}

public function mock404ConfigMap() returns http:Response {
    http:Response response = new;

    json body = {
        "kind": "Status",
        "apiVersion": "v1"
    ,
        "metadata": {},
        "status": "Failure",
        "message": "configmaps \"01ed7b08-f2b1-1166-82d5-649ae706d29d-definition\" not found",
        "reason": "NotFound",
        "details": {"name": "01ed7b08-f2b1-1166-82d5-649ae706d29d-definition", "kind": "configmaps"},
        "code": 404
    };
    response.setJsonPayload(body);
    response.statusCode = 404;
    return response;
}

public function mockOpenAPIJson() returns json {
    json openapi = {
        "openapi": "3.0.1",
        "info": {
            "title": "pizza1234567",
            "version": "1.0.0"
        },
        "security": [
            {
                "default": []
            }
        ],
        "paths": {
            "/menu": {
                "get": {
                    "responses": {
                        "200": {
                            "description": "OK"
                        }
                    },
                    "security": [
                        {
                            "default": []
                        }
                    ],
                    "x-throttling-tier": "Unlimited"
                }
            },
            "/order/{orderId}": {
                "post": {
                    "parameters": [
                        {
                            "name": "orderId",
                            "in": "path",
                            "required": true,
                            "schema": {
                                "type": "string"
                            }
                        }
                    ],
                    "responses": {
                        "200": {
                            "description": "OK"
                        }
                    },
                    "security": [
                        {
                            "default": []
                        }
                    ],
                    "x-throttling-tier": "Unlimited"
                }
            }
        },
        "components": {
            "securitySchemes": {
                "default": {
                    "type": "oauth2",
                    "flows": {
                        "implicit": {
                            "authorizationUrl": "https://test.com",
                            "scopes": {}
                        }
                    }
                }
            }
        }
    };
    return openapi;
}

public function mockpizzashackAPI11Definition() returns json {
    model:API api = {
        metadata: {
            name: "01ed7aca-eb6b-1178-a200-f604a4ce114a",
            namespace: "apk-platform",
            uid: "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
            creationTimestamp: ()
        },
        spec: {
            apiDisplayName: "pizzashackAPI11",
            apiType: "REST",
            apiVersion: "1.0.0",
            context: "/pizzashack11/1.0.0",
            organization: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
            prodHTTPRouteRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-production"
        }
    };
    APIClient apiclient = new ();
    return apiclient.retrieveDefaultDefinition(api);
}

public function mockPizzashackAPI12Definition() returns json {
    model:API api = {
        metadata: {
            name: "01ed7aca-eb6b-1178-a200-f604a4ce114a",
            namespace: "apk-platform",
            uid: "c5ab2423-b9e8-432b-92e8-35e6907ed5f9",
            creationTimestamp: ()
        },
        spec: {
            apiDisplayName: "pizzashackAPI12",
            apiType: "REST",
            apiVersion: "1.0.0",
            context: "/pizzashack12/1.0.0",
            organization: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
            prodHTTPRouteRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-production",
            definitionFileRef: ""
        }
    };
    APIClient apiclient = new ();
    return apiclient.retrieveDefaultDefinition(api);
}

public function mockPizzaShackAPI1Definition() returns json {
    model:API api = {
        kind: "API",
        metadata: {
            creationTimestamp: "2022-12-13T17:09:49Z",
            name: "01ed7b08-f2b1-1166-82d5-649ae706d29d",
            namespace: "apk-platform",
            uid: "7b7db1f0-0a9a-4f72-9f9d-5a1696d590c1"
        },
        spec: {
            apiDisplayName: "pizzashackAPI1",
            apiType: "REST",
            apiVersion: "1.0.0",
            context: "/pizzashack1/1.0.0",
            organization: "carbon.super",
            definitionFileRef: "01ed7b08-f2b1-1166-82d5-649ae706d29d-definition",
            prodHTTPRouteRef: "01ed7b08-f2b1-1166-82d5-649ae706d29d-production"
        }
    };
    APIClient apiclient = new ();
    return apiclient.retrieveDefaultDefinition(api);
}
