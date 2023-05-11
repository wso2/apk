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
import ballerina/test;
import runtime_domain_service.model;
import runtime_domain_service.java.lang;
import ballerina/log;
import wso2/apk_common_lib as commons;
public function getMockAPIList() returns model:APIList {
    model:EnvConfig[]? pizzashackAPIEndpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];
    model:EnvConfig[]? pizzashackAPI1Endpoint = [{httpRouteRefs: ["01ed7b08-f2b1-1166-82d5-649ae706d29d-production"]}];
    model:EnvConfig[]? pizzashackAPI11Endpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];
    model:EnvConfig[]? pizzashackAPI12Endpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];
    model:EnvConfig[]? pizzashackAPI13Endpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];

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
                                    "f:production": {}
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
                    "production": pizzashackAPIEndpoint
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
                                    "f:production": {}
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
                    "production": pizzashackAPI1Endpoint
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
                                    "f:production": {}
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
                    "production": pizzashackAPI11Endpoint
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
                                    "f:production": {}
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
                    "production": pizzashackAPI12Endpoint,
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
                                    "f:production": {}
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
                    "production": pizzashackAPI13Endpoint,
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
                "managedFields": [
                    {
                        "apiVersion": "dp.wso2.com/v1alpha1",
                        "fieldsType": "FieldsV1",
                        "fieldsV1": {
                            "f:spec": {".": {}, "f:apiDisplayName": {}, "f:apiType": {}, "f:apiVersion": {}, "f:context": {}, "f:definitionFileRef": {}, "f:production": {}}
                        },
                        "manager": "ballerina",
                        "operation": "Update",
                        "time": "2022-12-13T18:51:26Z"
                    }
                ],
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
                "production": [
                    {
                        "httpRouteRefs": ["01ed7b16-90f7-1a88-8113-a7e71796d460-production"]
                    }
                ]
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
                "managedFields": [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsV1": {"f:spec": {".": {}, "f:apiDisplayName": {}, "f:apiType": {}, "f:apiVersion": {}, "f:context": {}, "f:definitionFileRef": {}, "f:production": {}}}, "manager": "ballerina", "time": "2022-12-13T09:45:47Z", "operation": "Update", "fieldsType": "FieldsV1"}],
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
                "production": [
                    {
                        "httpRouteRefs": ["01ed7b16-90f7-1a88-8114-a7e71796d460-production"]
                    }
                ],
                "organization": "01ed7aca-eb6b-1178-a200-f604a4ce114b"
            }
        }
    };
    return message.toString();

}

public function getMockPizzaShakK8sAPI() returns model:API & readonly {
    model:EnvConfig[]? pizzashackAPIEndpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];
    model:API k8sAPI = {
        metadata: {
            name: "01ed7aca-eb6b-1178-a200-f604a4ce114a",
            namespace: "apk-platform",
            uid: "c5ab2423-b9e8-432b-92e8-35e6907ed5e8",
            creationTimestamp: "2022-12-13T09:45:47Z",
            generation: 1,
            selfLink: "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/apis/01ed7aca-eb6b-1178-a200-f604a4ce114a",
            resourceVersion: "5833",
            managedFields: [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsV1": {"f:spec": {".": {}, "f:apiDisplayName": {}, "f:apiType": {}, "f:apiVersion": {}, "f:context": {}, "f:definitionFileRef": {}, "f:production": {}}}, "manager": "ballerina", "time": "2022-12-13T09:45:47Z", "operation": "Update", "fieldsType": "FieldsV1"}]
        },
        spec: {
            apiDisplayName: "pizzashackAPI",
            apiType: "REST",
            apiVersion: "1.0.0",
            context: "/pizzashack/1.0.0",
            definitionFileRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-definition",
            production: pizzashackAPIEndpoint,
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

public function mockConfigMaps() returns http:Response|error {
    json definition = {"openapi": "3.0.1", "info": {"title": "pizza1234567", "version": "1.0.0"}, "security": [{"default": []}], "paths": {"/menu": {"get": {"responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-throttling-tier": "Unlimited"}}, "/order/{orderId}": {"post": {"parameters": [{"name": "orderId", "in": "path", "required": true, "schema": {"type": "string"}}], "responses": {"200": {"description": "OK"}}, "security": [{"default": []}], "x-throttling-tier": "Unlimited"}}}, "components": {"securitySchemes": {"default": {"type": "oauth2", "flows": {"implicit": {"authorizationUrl": "https://test.com", "scopes": {}}}}}}};
    http:Response response = new;
    json configmap = {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "creationTimestamp": "2023-01-05T05:34:44Z",
            "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a-definition",
            "namespace": "apk-platform",
            "resourceVersion": "113573",
            "uid": "ce2915d6-cdeb-4c70-8cdd-e03c158105ba"
        },
        binaryData: {
            [CONFIGMAP_DEFINITION_KEY] : check getBase64EncodedGzipContent(definition.toJsonString().toBytes())
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

public function convertJsonToYaml(string jsonString) returns string|error {
    commons:YamlUtil yamlUtil = commons:newYamlUtil1();
     string|lang:Exception unionResult = check yamlUtil.fromJsonStringToYaml(jsonString) ?: "";
     if unionResult is string {
            return unionResult;
     } else {
            log:printError(unionResult.message());
            return unionResult;
     }
}

public function mockpizzashackAPI11Definition() returns json {
    model:EnvConfig[]? pizzashackAPIEndpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];
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
            production: pizzashackAPIEndpoint
        }
    };
    APIClient apiclient = new ();
    return apiclient.retrieveDefaultDefinition(api);
}

public function mockPizzashackAPI12Definition() returns json {
    model:EnvConfig[]? pizzashackAPIEndpoint = [{httpRouteRefs: ["01ed7aca-eb6b-1178-a200-f604a4ce114a-production"]}];
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
            production: pizzashackAPIEndpoint,
            definitionFileRef: ""
        }
    };
    APIClient apiclient = new ();
    return apiclient.retrieveDefaultDefinition(api);
}

public function mockPizzaShackAPI1Definition(string organization) returns json {
    model:EnvConfig[]? pizzashackAPIEndpoint = [{httpRouteRefs: ["01ed7b08-f2b1-1166-82d5-649ae706d29d-production"]}];
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
            organization: organization,
            definitionFileRef: "01ed7b08-f2b1-1166-82d5-649ae706d29d-definition",
            production: pizzashackAPIEndpoint
        }
    };
    APIClient apiclient = new ();
    return apiclient.retrieveDefaultDefinition(api);
}

public function getMockInternalAPI() returns model:RuntimeAPI {
    model:RuntimeAPI runtimeAPI = {
        metadata: {name: "01ed7aca-eb6b-1178-a200-f604a4ce114a", namespace: "apk-platform"},
        spec: {
            name: "pizzashackAPI",
            context: "",
            'type: "REST",
            'version: "1.0.0",
            endpointConfig: {
                "endpoint_type": "http",
                "sandbox_endpoints": {
                    "url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"
                },
                "production_endpoints": {
                    "url": "https://pizzashack-service:8080/am/sample/pizzashack/v3/api/"
                }
            },
            operations: [
                {verb: "GET", target: "/*", authTypeEnabled: true},
                {verb: "PUT", target: "/*", authTypeEnabled: true},
                {verb: "POST", target: "/*", authTypeEnabled: true},
                {verb: "DELETE", target: "/*", authTypeEnabled: true}
            ]
        }
    };
    return runtimeAPI;
}

@test:Config {}
public function testConvertion() {
    json backendList = {
        "apiVersion": "v1",
        "items": [
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "creationTimestamp": "2023-03-16T08:05:48Z",
                    "generation": 1,
                    "labels": {
                        "api-name": "testAPIV2",
                        "api-version": "1.0.0"
                    },
                    "name": "backend-3e5fb4d0a2a8d53915dfa179f3089821f0f6abc5-api",
                    "namespace": "backend-cr",
                    "resourceVersion": "1166701",
                    "uid": "0b16ae85-9a69-4256-a8d2-76b2c32b2522"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend.test-apk.svc.cluster.local",
                            "port": 80
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "creationTimestamp": "2023-03-16T08:01:21Z",
                    "generation": 1,
                    "labels": {
                        "api-name": "testAPI",
                        "api-version": "1.0.0"
                    },
                    "name": "backend-cbe2237e4e97924e88d9300fc28719944ce634d2-api",
                    "namespace": "backend-cr",
                    "resourceVersion": "1166436",
                    "uid": "73766dd5-dacf-4aea-ab60-bd87b25c2693"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend.test-apk.svc.cluster.local",
                            "port": 80
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-admin-ds-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147899",
                    "uid": "dffb8f27-47b4-4a6c-bb70-df192b9713e7"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-admin-ds-service.backend-cr",
                            "port": 9443
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-backoffice-ds-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147897",
                    "uid": "25b17252-039f-4d24-893b-0e27aab1f675"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-backoffice-ds-service.backend-cr",
                            "port": 9443
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-devportal-ds-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147895",
                    "uid": "69db6be4-5db0-40f7-adaf-bdf62bd42081"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-devportal-ds-service.backend-cr",
                            "port": 9443
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-idp-ds-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147900",
                    "uid": "64544deb-1b09-4055-bc17-8513d5dd895d"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-idp-ds-service.backend-cr",
                            "port": 9443
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-idp-ui-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147896",
                    "uid": "17083578-47f5-4135-b154-69431ab07f9c"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-idp-ui-service.backend-cr",
                            "port": 9443
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-internal-admin-ds-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147901",
                    "uid": "f05bedec-f9b5-4951-b750-266d657ad3ce"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-admin-ds-service.backend-cr",
                            "port": 9444
                        }
                    ]
                }
            },
            {
                "apiVersion": "dp.wso2.com/v1alpha1",
                "kind": "Backend",
                "metadata": {
                    "annotations": {
                        "meta.helm.sh/release-name": "backend-cr",
                        "meta.helm.sh/release-namespace": "backend-cr"
                    },
                    "creationTimestamp": "2023-03-15T18:11:43Z",
                    "generation": 1,
                    "labels": {
                        "app.kubernetes.io/managed-by": "Helm"
                    },
                    "name": "backend-cr-wso2-apk-runtime-ds-backend",
                    "namespace": "backend-cr",
                    "resourceVersion": "1147898",
                    "uid": "cba43be3-6f91-4511-9330-12bad6d19837"
                },
                "spec": {
                    "protocol": "http",
                    "services": [
                        {
                            "host": "backend-cr-wso2-apk-runtime-ds-service.backend-cr",
                            "port": 9443
                        }
                    ]
                }
            }
        ],
        "kind": "List",
        "metadata": {
            "resourceVersion": ""
        }
    };
    model:BackendList|error backends = backendList.cloneWithType(model:BackendList);
    test:assertTrue(backends is model:BackendList);
}
