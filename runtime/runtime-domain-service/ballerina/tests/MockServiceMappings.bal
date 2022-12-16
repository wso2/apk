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

public function getMockServiceMappings() returns json {
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
    return response;
}

public function getServiceMappingEvent() returns string {
    json message = {
        "type": "ADDED",
        "object": {
            "apiVersion": "dp.wso2.com/v1alpha1",
            "kind": "ServiceMapping",
            "metadata": {
                "creationTimestamp": "2022-12-13T09:45:47Z",
                "generation": 1,
                "managedFields": [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsType": "FieldsV1", "fieldsV1": {"f:spec": {".": {}, "f:apiRef": {".": {}, "f:name": {}, "f:namespace": {}}, "f:serviceRef": {".": {}, "f:name": {}, "f:namespace": {}}}}, "manager": "ballerina", "operation": "Update", "time": "2022-12-13T09:45:47Z"}],
                "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a-servicemapping",
                "namespace": "apk-platform",
                "resourceVersion": "5834",
                "selfLink": "/apis/dp.wso2.com/v1alpha1/namespaces/apk-platform/servicemappings/01ed7aca-eb6b-1178-a200-f604a4ce114a-servicemapping",
                "uid": "79d280d6-df31-4017-911f-3229955fdb55"
            },
            "spec": {
                "apiRef": {
                    "name": "01ed7aca-eb6b-1178-a200-f604a4ce114a",
                    "namespace": "apk-platform"
                },
                "serviceRef": {
                    "name": "backend",
                    "namespace": "apk"
                }
            }
        }
    };
    return message.toString();
}

public function getNextServiceMappingEvent() returns string {
    json message = {
        "type": "ADDED",
        "object": {
            "apiVersion": "dp.wso2.com/v1alpha1",
            "kind": "ServiceMapping",
            "metadata": {
                "creationTimestamp": "2022-12-13T17:09:49Z",
                "generation": 1,
                "managedFields": [{"apiVersion": "dp.wso2.com/v1alpha1", "fieldsType": "FieldsV1", "fieldsV1": {"f:spec": {".": {}, "f:apiRef": {".": {}, "f:name": {}, "f:namespace": {}}, "f:serviceRef": {".": {}, "f:name": {}, "f:namespace": {}}}}, "manager": "ballerina", "operation": "Update", "time": "2022-12-13T17:09:49Z"}],
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
        }
    };
    return message.toString();
}
