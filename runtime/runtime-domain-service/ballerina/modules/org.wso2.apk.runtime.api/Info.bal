import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.util as javautil;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.api.Info` class.
@java:Binding {'class: "org.wso2.apk.runtime.api.Info"}
public distinct class Info {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.api.Info` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.api.Info` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.api.Info` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_api_Info_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_api_Info_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getContext` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getContext() returns string? {
        return java:toString(org_wso2_apk_runtime_api_Info_getContext(self.jObj));
    }

    # The function that maps to the `getDescription` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getDescription() returns string? {
        return java:toString(org_wso2_apk_runtime_api_Info_getDescription(self.jObj));
    }

    # The function that maps to the `getEndpoints` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `javautil:List` value returning from the Java mapping.
    public isolated function getEndpoints() returns javautil:List {
        handle externalObj = org_wso2_apk_runtime_api_Info_getEndpoints(self.jObj);
        javautil:List newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getName` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getName() returns string? {
        return java:toString(org_wso2_apk_runtime_api_Info_getName(self.jObj));
    }

    # The function that maps to the `getOpenAPIVersion` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getOpenAPIVersion() returns string? {
        return java:toString(org_wso2_apk_runtime_api_Info_getOpenAPIVersion(self.jObj));
    }

    # The function that maps to the `getVersion` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getVersion() returns string? {
        return java:toString(org_wso2_apk_runtime_api_Info_getVersion(self.jObj));
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_api_Info_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.api.Info`.
    public function notify() {
        org_wso2_apk_runtime_api_Info_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.api.Info`.
    public function notifyAll() {
        org_wso2_apk_runtime_api_Info_notifyAll(self.jObj);
    }

    # The function that maps to the `setContext` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setContext(string arg0) {
        org_wso2_apk_runtime_api_Info_setContext(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setDescription` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setDescription(string arg0) {
        org_wso2_apk_runtime_api_Info_setDescription(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setEndpoints` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `javautil:List` value required to map with the Java method parameter.
    public function setEndpoints(javautil:List arg0) {
        org_wso2_apk_runtime_api_Info_setEndpoints(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setName` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setName(string arg0) {
        org_wso2_apk_runtime_api_Info_setName(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setOpenAPIVersion` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setOpenAPIVersion(string arg0) {
        org_wso2_apk_runtime_api_Info_setOpenAPIVersion(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setVersion` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setVersion(string arg0) {
        org_wso2_apk_runtime_api_Info_setVersion(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_Info_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_Info_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.Info`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_Info_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.Info`.
#
# + return - The new `Info` class generated.
public function newInfo1() returns Info {
    handle externalObj = org_wso2_apk_runtime_api_Info_newInfo1();
    Info newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_runtime_api_Info_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_api_Info_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_Info_getContext(handle receiver) returns handle = @java:Method {
    name: "getContext",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_Info_getDescription(handle receiver) returns handle = @java:Method {
    name: "getDescription",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_Info_getEndpoints(handle receiver) returns handle = @java:Method {
    name: "getEndpoints",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_Info_getName(handle receiver) returns handle = @java:Method {
    name: "getName",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_Info_getOpenAPIVersion(handle receiver) returns handle = @java:Method {
    name: "getOpenAPIVersion",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_Info_getVersion(handle receiver) returns handle = @java:Method {
    name: "getVersion",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_Info_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_Info_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_Info_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_Info_setContext(handle receiver, handle arg0) = @java:Method {
    name: "setContext",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_Info_setDescription(handle receiver, handle arg0) = @java:Method {
    name: "setDescription",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_Info_setEndpoints(handle receiver, handle arg0) = @java:Method {
    name: "setEndpoints",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.util.List"]
} external;

function org_wso2_apk_runtime_api_Info_setName(handle receiver, handle arg0) = @java:Method {
    name: "setName",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_Info_setOpenAPIVersion(handle receiver, handle arg0) = @java:Method {
    name: "setOpenAPIVersion",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_Info_setVersion(handle receiver, handle arg0) = @java:Method {
    name: "setVersion",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_Info_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_Info_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_api_Info_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_runtime_api_Info_newInfo1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.Info",
    paramTypes: []
} external;

