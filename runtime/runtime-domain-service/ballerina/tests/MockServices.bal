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
public function getMockServiceList() returns json {

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
                        "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"name\":\"backend\",\"namespace\":\"apk\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"targetPort\":80}],\"selector\":{\"app\":\"httpbin\"}}}"
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
    return response;
}

public function getServiceEvent() returns string {
    json message = {
        "type": "ADDED",
        "object": {
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "backend",
                "namespace": "apk",
                "selfLink": "/api/v1/namespaces/apk/services/backend",
                "uid": "275b00d1-722c-4df2-b65a-9b14677abe4b",
                "resourceVersion": "1514",
                "creationTimestamp": "2022-12-13T08:25:09Z",
                "annotations": {"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"name\":\"backend\",\"namespace\":\"apk\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"targetPort\":80}],\"selector\":{\"app\":\"httpbin\"}}}"},
                "managedFields": [{"manager": "kubectl-client-side-apply", "operation": "Update", "apiVersion": "v1", "time": "2022-12-13T08:25:09Z", "fieldsType": "FieldsV1", "fieldsV1": {"f:metadata": {"f:annotations": {".": {}, "f:kubectl.kubernetes.io/last-applied-configuration": {}}}, "f:spec": {"f:ports": {".": {}, "k:{\"port\":80,\"protocol\":\"TCP\"}": {".": {}, "f:name": {}, "f:port": {}, "f:protocol": {}, "f:targetPort": {}}}, "f:selector": {".": {}, "f:app": {}}, "f:sessionAffinity": {}, "f:type": {}}}}]
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
                "selector": {"app": "httpbin"},
                "clusterIP": "10.98.200.176",
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {"loadBalancer": {}}
        }
    };
    return message.toJsonString();
}

public function getNextMockServiceEvent() returns string {
    json message = {
        "type": "ADDED",
        "object": {
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "backend-1",
                "namespace": "apk",
                "selfLink": "/api/v1/namespaces/apk/services/backend",
                "uid": "275b00d1-722c-4df2-b65a-9b14678abe4b",
                "resourceVersion": "1517",
                "creationTimestamp": "2022-12-13T08:30:09Z",
                "annotations": {"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"name\":\"backend-1\",\"namespace\":\"apk\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"targetPort\":80}],\"selector\":{\"app\":\"httpbin\"}}}\n"},
                "managedFields": [{"manager": "kubectl-client-side-apply", "operation": "Update", "apiVersion": "v1", "time": "2022-12-13T08:25:09Z", "fieldsType": "FieldsV1", "fieldsV1": {"f:metadata": {"f:annotations": {".": {}, "f:kubectl.kubernetes.io/last-applied-configuration": {}}}, "f:spec": {"f:ports": {".": {}, "k:{\"port\":80,\"protocol\":\"TCP\"}": {".": {}, "f:name": {}, "f:port": {}, "f:protocol": {}, "f:targetPort": {}}}, "f:selector": {".": {}, "f:app": {}}, "f:sessionAffinity": {}, "f:type": {}}}}]
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
                "selector": {"app": "httpbin"},
                "clusterIP": "10.98.200.176",
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {"loadBalancer": {}}
        }
    };
    return message.toJsonString();
}
