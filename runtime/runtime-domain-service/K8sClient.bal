import ballerina/http;
import ballerina/io;

const string K8S_API_ENDPOINT = "/api/v1";
final http:Client k8sApiServerEp = check initializeK8sClient();
configurable string k8sHost = "kubernetes.default";
configurable string saTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token";
configurable string token = check io:fileReadString(saTokenPath);
configurable string caCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt";

function init() {
}

function initializeK8sClient() returns http:Client|error {
    http:Client k8sApiClient = check new ("https://" + k8sHost + K8S_API_ENDPOINT,
    auth = {
        token: token
    },
    secureSocket = {
        cert: caCertPath

    }
    );
    return k8sApiClient;
}

isolated function getServicesListFromK8s(string namespace) returns ServiceList|error {
    Service[] serviceNames = [];
    string endpoint = "/namespaces/" + namespace + "/services";
    error|json serviceResp = k8sApiServerEp->get(endpoint, targetType = json);
    if (serviceResp is json) {
        json[] serviceArr = <json[]>check serviceResp.items;
        foreach json i in serviceArr {
            Service serviceData = {
                name: <string>check i.metadata.name,
                namespace: <string>check i.metadata.namespace,
                'type: <string>check i.spec.'type
            };
            serviceNames.push(serviceData);
        }
        ServiceList serviceList = {
            list: serviceNames
        };
        return serviceList;
    }
    return error("error while retrieving service list from K8s API server for namespace : " +
                namespace);
}

isolated function getServiceFromK8s(string name, string namespace) returns ServiceList|error {
    Service[] serviceNames = [];
    string endpoint = "/namespaces/" + namespace + "/services/" + name;
    error|json serviceResp = k8sApiServerEp->get(endpoint, targetType = json);
    if (serviceResp is json) {
        json[] serviceArr = <json[]>check serviceResp.items;
        foreach json i in serviceArr {
            Service serviceData = {
                name: <string>check i.metadata.name,
                namespace: <string>check i.metadata.namespace,
                'type: <string>check i.spec.'type
            };
            serviceNames.push(serviceData);
        }
        ServiceList serviceList = {
            list: serviceNames
        };
        return serviceList;
    }
    return error("error while retrieving service list from K8s API server for namespace : " +
                namespace);
}
