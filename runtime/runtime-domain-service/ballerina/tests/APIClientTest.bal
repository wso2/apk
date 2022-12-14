import ballerina/test;
import ballerina/http;

@test:Mock {functionName: "getServiceMappingClient"}
test:MockFunction websocketServiceMappingClient = new ();
@test:Mock {functionName: "getServiceClient"}
test:MockFunction websocketServiceClient = new ();
@test:Mock {functionName: "getClient"}
test:MockFunction websocketAPIClient = new ();

@test:Mock {
    functionName: "initializeK8sClient"
}
function getMockK8sClient() returns http:Client|error {
    return test:mock(http:Client);
}

@test:BeforeSuite
public function before() {
    test:prepare(k8sApiServerEp).when("get").thenReturn(getMockAPIList());
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/apis")
        .thenReturn(getMockAPIList());
    test:prepare(k8sApiServerEp).when("get").withArguments("/api/v1/services")
        .thenReturn(getMockServiceList());
    test:prepare(k8sApiServerEp).when("get").withArguments("/apis/dp.wso2.com/v1alpha1/servicemappings")
        .thenReturn(getMockServiceMappings());
}

@test:Config {}
public function testretrievePathPrefix() {
    APIClient apiclient = new ();
    string retrievePathPrefix = apiclient.retrievePathPrefix("/abc/1.0.0", "1.0.0", "/abc");
    test:assertEquals(retrievePathPrefix, "/abc/1.0.0/abc");

}
