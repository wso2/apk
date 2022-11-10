import ballerina/http;

import ballerina/io;

const string K8S_API_ENDPOINT = "/api/v1";
final http:Client k8sApiServerEp = check initializeK8sClient();
configurable string k8sHost = "kubernetes.default";
configurable string saTokenPath = "var/run/secrets/kubernetes.io/serviceaccount/token";
configurable string token = check io:fileReadString(saTokenPath);
configurable string caCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt";

function initializeK8sClient() returns http:Client|error {
    http:Client k8sApiClient = check new ("https://" + k8sHost ,
    auth = {
        token: token
    },
    secureSocket = {
        cert: caCertPath

    }
    );
    return k8sApiClient;
}

isolated function getServicesListInNamespace(string namespace) returns ServiceList|error {
    Service[] serviceNames = [];
    string endpoint = K8S_API_ENDPOINT + "/namespaces/" + namespace + "/services";
    error|json serviceResp = k8sApiServerEp->get(endpoint, targetType = json);
    if (serviceResp is json) {
        json[] serviceArr = <json[]>check serviceResp.items;
        foreach json i in serviceArr {
            Service serviceData = {
                id: <string>check i.metadata.uid,
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

isolated function getServicesListFromK8s() returns ServiceList|error {
    Service[] serviceNames = [];
    string endpoint = K8S_API_ENDPOINT + "/services";
    error|json serviceResp = k8sApiServerEp->get(endpoint, targetType = json);
    if (serviceResp is json) {
        json[] serviceArr = <json[]>check serviceResp.items;
        foreach json i in serviceArr {
            Service serviceData = {
                id: <string>check i.metadata.uid,
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
    return error("error while retrieving service list from K8s API server for namespace");
}

isolated function getServiceFromK8s(string name, string namespace) returns ServiceList|error {
    Service[] serviceNames = [];
    string endpoint = K8S_API_ENDPOINT + "/namespaces/" + namespace + "/services/" + name;
    error|json serviceResp = k8sApiServerEp->get(endpoint, targetType = json);
    if (serviceResp is json) {
        json[] serviceArr = <json[]>check serviceResp.items;
        foreach json i in serviceArr {
            Service serviceData = {
                id: <string>check i.metadata.uid,
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

function getAPIListInNamespace(string namespace) returns APIList|error {
    API[] APINames = [];
    string endpoint = "/apis/dp.wso2.com/v1alpha1/namespaces/" + namespace +"/apis";
    error|json APIResp = k8sApiServerEp->get(endpoint, targetType = json);
    if (APIResp is json) {
        json[] serviceArr = <json[]>check APIResp.items;
        foreach json i in serviceArr {
            API APIData = {
                context: <string>check i.spec.context,
                name: <string>check i.metadata.name,
            'version: <string>check i.spec.'apiVersion};
            APINames.push(APIData);
        }
        APIList APIList = {
            list: APINames
        };
        return APIList;
    }
    return error("error while retrieving API list from K8s API server for namespace : " +
                namespace);
}

//Get APIs deployed in default namespace by APIId.
function getAPIById(string id) returns API|InternalServerErrorError|BadRequestError| error {
    boolean APIIDAvailable = id.length() > 0 ? true : false;
    if (APIIDAvailable && string:length(id.toString()) > 0) 
    {
        //TODO replace default namespace to work with any namespace. As of now API contract sends only query to this API and 
        //hence default namespace hard coded in the implementation
        string endpoint = "/apis/dp.wso2.com/v1alpha1/namespaces/" + "default" +"/apis/"+id;
        error|json APIResp = k8sApiServerEp->get(endpoint, targetType = json);
        if APIResp is error {
            InternalServerErrorError internalError = {body: {code: 900910, message: "APIResp.message()"}};
            return internalError;
        } 
        else 
        {
            API APIData = {
                        context: <string>check APIResp.spec.context,
                        name: <string>check APIResp.metadata.name,
                        'version: <string>check APIResp.spec.'apiVersion};
            return APIData;
        }
    }
    BadRequestError badRequestError = {body: {code: 900910, message: "missing required attributes"}};
    return badRequestError;
} 

//Get all deployed APIs in namespace with specific search query
function getAPIListInNamespaceWithQuery(string? query, int 'limit = 25, int offset = 0, string sortBy = "createdTime", string sortOrder = "desc") returns  APIList|InternalServerErrorError|BadRequestError|error {
        boolean queryAvailable = query == () ? false : true;
        if (queryAvailable && string:length(query.toString()) > 0) 
        {
            API[] APINames = [];
            //TODO replace default namespace to work with any namespace. As of now API contract sends only query to this API and 
            //hence default namespace hard coded in the implementation
            string endpoint = "/apis/dp.wso2.com/v1alpha1/namespaces/" + "default" + "/apis?" + query.toString();
            error|json APIResp = k8sApiServerEp->get(endpoint, targetType = json);
            if (APIResp is json) {
            json[] serviceArr = <json[]>check APIResp.items;
                foreach json i in serviceArr 
                {
                    API APIData = {
                        context: <string>check i.spec.context,
                        name: <string>check i.metadata.name,
                        'version: <string>check i.spec.'apiVersion};
                        APINames.push(APIData);
                }
                APIList APIList = {
                    list: APINames
                };
                return APIList;
            }
            else {
                InternalServerErrorError internalError = {body: {code: 900910, message: APIResp.message()}};
                return internalError;
            }
        }
        BadRequestError badRequestError = {body: {code: 900910, message: "missing required attributes"}};
        return badRequestError;
}
