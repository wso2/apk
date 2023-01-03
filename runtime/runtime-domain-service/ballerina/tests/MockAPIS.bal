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

public function getMockAPIList() returns json {

    json response = {
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
                    "apiType": "HTTP",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack/1.0.0",
                    "organization": "carbon.super",
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
                    "apiType": "HTTP",
                    "apiVersion": "1.0.0",
                    "context": "/pizzashack1/1.0.0",
                    "organization": "carbon.super",
                    "definitionFileRef": "01ed7b08-f2b1-1166-82d5-649ae706d29d-definition",
                    "prodHTTPRouteRef": "01ed7b08-f2b1-1166-82d5-649ae706d29d-production"
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
                "apiType": "HTTP",
                "apiVersion": "1.0.0",
                "context": "/pizzashack6/1.0.0",
                "organization": "carbon.super",
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
                "apiType": "HTTP",
                "apiVersion": "1.0.0",
                "context": "/demoapi/1.0.0",
                "definitionFileRef": "01ed7b16-90f7-1a88-8114-a7e71796d460-definition",
                "prodHTTPRouteRef": "01ed7b16-90f7-1a88-8114-a7e71796d460-production",
                "organization": "carbon.super"
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
            apiType: "HTTP",
            apiVersion: "1.0.0",
            context: "/pizzashack/1.0.0",
            definitionFileRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-definition",
            prodHTTPRouteRef: "01ed7aca-eb6b-1178-a200-f604a4ce114a-production",
            organization: "carbon.super"
        }
    };
    return k8sAPI.cloneReadOnly();
}

public function mock404Response() returns http:Response {
    http:Response response = new;
    response.statusCode = 404;
    return response;
}
