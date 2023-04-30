import ballerina/http;
import ballerina/lang.runtime;

listener http:Listener ep0 = new (8443, http1Settings = {keepAlive: "ALWAYS"}
);
listener http:Listener ep1 = new (8444, secureSocket = {
    'key: {
        certFile: "/home/ineterceptor/tls.pem",
        keyFile: "/home/ineterceptor/tls.key"
    }
}, http1Settings = {keepAlive: "ALWAYS"}
    );
listener http:Listener ep2 = new (8445, secureSocket = {
    'key: {
        certFile: "/home/ineterceptor/tls.pem",
        keyFile: "/home/ineterceptor/tls.key"
    }
}, http1Settings = {keepAlive: "ALWAYS"}
    );    

function init() returns error? {
    check ep0.attach(interceptorService, "/api/v1");
    check ep0.'start();
    runtime:registerListener(ep0);
    check ep1.attach(interceptorService, "/api/v1");
    check ep1.'start();
    runtime:registerListener(ep1);
    check ep2.attach(gwInterceptorService, "/api/v1");
    check ep2.'start();
    runtime:registerListener(ep2);

}
