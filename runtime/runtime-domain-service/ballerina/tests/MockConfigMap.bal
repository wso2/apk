import runtime_domain_service.model;

public function getConfigMapEvent() returns string {
    json event = {
        "type": "ADDED",
        "object": {
            apiVersion: "v1",
            kind: "ConfigMap",
            metadata: {
                name: "org1",
                namespace: "apk-platform",
                resourceVersion: "28705",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                labels: {
                    "managed-by": "apk"
                }
            },
            data: {
                "org1": "org1"
            }
        }
    };
    return event.toJsonString();
}
public function getConfigMapUpdateEvent() returns string{
        json event = {
        "type": "MODIFIED",
        "object": {
            apiVersion: "v1",
            kind: "ConfigMap",
            metadata: {
                name: "org1",
                namespace: "apk-platform",
                resourceVersion: "28714",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                labels: {
                    "managed-by": "apk"
                }
            },
            data: {
                "org1": "org1"
            }
        }
    };
    return event.toJsonString();
}
public function getConfigMapDeleteEvent() returns string{
        json event = {
        "type": "DELETED",
        "object": {
            apiVersion: "v1",
            kind: "ConfigMap",
            metadata: {
                name: "org1",
                namespace: "apk-platform",
                resourceVersion: "28714",
                uid: "01ed7aca-eb6b-1178-a200-f604a4ce114b",
                labels: {
                    "managed-by": "apk"
                }
            },
            data: {
                "org1": "org1"
            }
        }
    };
    return event.toJsonString();
}
public function getMockLabelList() returns model:ConfigMapList {
    model:ConfigMapList configMapList = {
        metadata: {resourceVersion: "28702"},
        items: []
    };
    return configMapList;
}
