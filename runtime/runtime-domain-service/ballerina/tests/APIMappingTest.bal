import ballerina/test;

@test:Config {dataProvider: fromAPIToAPKConfDataProvider}
public function testAPItoApkConfMapping(API api, APKConf apkConf) {
    test:assertEquals(fromAPIToAPKConf(api), apkConf, "API to APKConf mapping failed");
}

public function fromAPIToAPKConfDataProvider() returns map<[API, APKConf]> {
    map<[API, APKConf]> data = {
        "1": [{name: "simple", context: "/simple", version: "1.0.0"}, {name: "simple", context: "/simple", version: "1.0.0"}],
        "2": [
            {
                "id": "c358c21f-547e-4e51-8410-9a5275531aa3",
                "name": "testAPIV2",
                "context": "/testAPIV2/2.0.0",
                "version": "2.0.0",
                "type": "REST",
                "operations": [
                    {
                        "target": "/*",
                        "verb": "GET",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "PUT",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "POST",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "DELETE",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "PATCH",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    }
                ],
                "serviceInfo": {
                    "name": "backend",
                    "namespace": "test-apk"
                },
                "createdTime": "2023-06-01T16:05:40Z",
                "lastUpdatedTime": "2023-06-01T16:05:40Z"
            },
            {
                "id": "c358c21f-547e-4e51-8410-9a5275531aa3",
                "name": "testAPIV2",
                "context": "/testAPIV2/2.0.0",
                "version": "2.0.0",
                "type": "REST",
                "operations": [
                    {
                        "target": "/*",
                        "verb": "GET",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "PUT",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "POST",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "DELETE",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    },
                    {
                        "target": "/*",
                        "verb": "PATCH",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    }
                ],
                "serviceInfo": {
                    "name": "backend",
                    "namespace": "test-apk"
                }
            }
        ],
        "3": [
            {
                "id": "123-456",
                "name": "testAPIV5",
                "context": "/testAPIV5/1.0.0",
                "version": "1.0.0",
                "type": "REST",
                "operations": [
                    {
                        "target": "/headers",
                        "verb": "GET",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    }
                ],
                "endpointConfig": {
                    "endpoint_type": "http",
                    "sandbox_endpoints": {
                        "url": "http://backend.test-apk.svc.cluster.local"
                    },
                    "production_endpoints": {
                        "url": "http://backend.test-apk.svc.cluster.local"
                    },
                    "endpoint_security": {
                        "production": {
                            "enabled": true,
                            "type": "Basic",
                            "username": "admin123",
                            "password": "admin123"
                        },
                        "sandbox": {
                            "enabled": false,
                            "type": "Basic",
                            "username": "admin1",
                            "password": "admin1"
                        }
                    }
                }
            },
            {
                "id": "123-456",
                "name": "testAPIV5",
                "context": "/testAPIV5/1.0.0",
                "version": "1.0.0",
                "type": "REST",
                "operations": [
                    {
                        "target": "/headers",
                        "verb": "GET",
                        "authTypeEnabled": true,
                        "scopes": [],
                        "operationPolicies": {
                            "request": [],
                            "response": []
                        }
                    }
                ],
                "endpointConfig": {
                    sandbox: {
                        endpointURL: "http://backend.test-apk.svc.cluster.local",
                        endpointSecurity: {
                            enable: false,
                            securityType: "Basic",
                            securityProperties: {
                                "username": "admin1",
                                "password": "admin1"
                            }
                        }
                    },
                    production: {
                        endpointURL: "http://backend.test-apk.svc.cluster.local",
                        endpointSecurity: {
                            enable: true,
                            securityType: "Basic",
                            securityProperties: {
                                "username": "admin123",
                                "password": "admin123"
                            }
                        }
                    }
                }
            }
        ]
    };

    return data;
}
