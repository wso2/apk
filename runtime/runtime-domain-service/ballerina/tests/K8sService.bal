// import ballerina/http;

// listener http:Listener ep1 = new (9090,
// secureSocket = {
//     key: {
//         certFile: "tests/resource/wso2carbon.crt",
//         keyFile: "tests/resource/wso2carbon.key"
//     }
// });

// // The `absolute resource path` can be omitted. Then, it defaults to `/`.
// service on ep1 {

//     resource function get apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/apis/[string name]() returns json {
//         json message = {};
//         return message;
//     }
//     resource function put apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/apis/[string name]() returns json {
//         json message = {};
//         return message;
//     }
//     resource function post apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/apis/[string name]() returns json {
//         json message = {};
//         return message;

//     }
//     resource function delete apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/apis/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get apis/'dp\.wso2\.com/v1alpha1/servicemappings() returns http:Ok {

//         http:Ok okResponse = {body: serviceMappingList};
//         return okResponse;
//     }
//     resource function post apis/'dp\.wso2\.com/v1alpha1/servicemappings() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/servicemappings/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function post apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/servicemappings/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function delete apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/servicemappings/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get api/v1/namespaces/[string namespace]/services() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get api/v1/namespaces/[string namespaces]/services/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get api/v1/services() returns http:Ok {

//         http:Ok okResponse = {body: servicesResponse};
//         return okResponse;
//     }
//     resource function get apis/'dp\.wso2\.com/v1alpha1/apis() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function post apis/'gateway\.networking\.k8s\.io/v1beta1/namespaces/[string namespaces]/httproutes() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get apis/'gateway\.networking\.k8s\.io/v1beta1/namespaces/[string namespaces]/httproutes/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function put apis/'gateway\.networking\.k8s\.io/v1beta1/namespaces/[string namespaces]/httproutes/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function delete apis/'gateway\.networking\.k8s\.io/v1beta1/namespaces/[string namespaces]/httproutes/[string name]() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function get apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/apis() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
//     resource function post apis/'dp\.wso2\.com/v1alpha1/namespaces/[string namespace]/apis() returns http:Ok {
//         http:Ok okResponse = {};
//         return okResponse;
//     }
// }
