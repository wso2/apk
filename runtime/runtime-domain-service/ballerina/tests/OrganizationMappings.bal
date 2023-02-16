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
function getNextOrganizationEvent() returns string {
    json event = {
        "type": "Modified",
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
