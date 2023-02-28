//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
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
function getOrganizationWatchEvent() returns string {
    json event = {
        "type": "ADDED",
        "object": {
            apiVersion: "cp.wso2.com/v1alpha1",
            kind: "Organization",
            metadata: {
                name: "org2",
                namespace: "apk-platform",
                resourceVersion: "28705",
                selfLink: "/apis/cp.wso2.com/v1alpha1/namespaces/apk-platform/organizations/org2",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b"
            },
            spec: {
                name: "org2",
                displayName: "org2",
                uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114c",
                organizationClaimValue: "org2",
                enabled: true
            }
        }
    };
    return event.toJsonString();
}
function getOrganizationWatchUpdateEvent() returns string {
    json event = {
        "type": "MODIFIED",
        "object": {
            apiVersion: "cp.wso2.com/v1alpha1",
            kind: "Organization",
            metadata: {
                name: "org3",
                namespace: "apk-platform",
                resourceVersion: "28730",
                selfLink: "/apis/cp.wso2.com/v1alpha1/namespaces/apk-platform/organizations/org2",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b"
            },
            spec: {
                name: "org3",
                displayName: "org3",
                uuid: "01ed7aca-eb6b-1178-a200-f604a4ce115c",
                organizationClaimValue: "org3",
                enabled: true
            }
        }
    };
    return event.toJsonString();
}
function getOrganizationWatchDeleteEvent() returns string {
    json event = {
        "type": "DELETED",
        "object": {
            apiVersion: "cp.wso2.com/v1alpha1",
            kind: "Organization",
            metadata: {
                name: "org3",
                namespace: "apk-platform",
                resourceVersion: "28731",
                selfLink: "/apis/cp.wso2.com/v1alpha1/namespaces/apk-platform/organizations/org2",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b"
            },
            spec: {
                name: "org3",
                displayName: "org3",
                uuid: "01ed7aca-eb6b-1178-a200-f604a4ce115c",
                organizationClaimValue: "org3",
                enabled: true
            }
        }
    };
    return event.toJsonString();
}
function getNextOrganizationEvent() returns string {
    json event = {
        "type": "MODIFIED",
        "object": {
            apiVersion: "cp.wso2.com/v1alpha1",
            kind: "Organization",
            metadata: {
                name: "org1",
                namespace: "apk-platform",
                resourceVersion: "28714",
                selfLink: "/apis/cp.wso2.com/v1alpha1/namespaces/apk-platform/organizations/org2",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b"
            },
            spec: {
                name: "org1",
                displayName: "org 1",
                uuid: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                organizationClaimValue: "org1",
                enabled: true
            }
        }
    };
    return event.toJsonString();
}
