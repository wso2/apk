import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.io as javaio;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.ZIPUtils` class.
@java:Binding {'class: "org.wso2.apk.runtime.ZIPUtils"}
public distinct class ZIPUtils {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.ZIPUtils` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.ZIPUtils` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.ZIPUtils` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.runtime.ZIPUtils`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_ZIPUtils_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.ZIPUtils`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_ZIPUtils_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.ZIPUtils`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_ZIPUtils_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.ZIPUtils`.
    public function notify() {
        org_wso2_apk_runtime_ZIPUtils_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.ZIPUtils`.
    public function notifyAll() {
        org_wso2_apk_runtime_ZIPUtils_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.ZIPUtils`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_ZIPUtils_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.ZIPUtils`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_ZIPUtils_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.ZIPUtils`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_ZIPUtils_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.ZIPUtils`.
#
# + return - The new `ZIPUtils` class generated.
public function newZIPUtils1() returns ZIPUtils {
    handle externalObj = org_wso2_apk_runtime_ZIPUtils_newZIPUtils1();
    ZIPUtils newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `zipDir` method of `org.wso2.apk.runtime.ZIPUtils`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `javaio:IOException` value returning from the Java mapping.
public isolated function ZIPUtils_zipDir(string arg0, string arg1) returns javaio:IOException? {
    error|() externalObj = org_wso2_apk_runtime_ZIPUtils_zipDir(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        javaio:IOException e = error javaio:IOException(javaio:IOEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

function org_wso2_apk_runtime_ZIPUtils_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_ZIPUtils_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: []
} external;

function org_wso2_apk_runtime_ZIPUtils_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: []
} external;

function org_wso2_apk_runtime_ZIPUtils_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: []
} external;

function org_wso2_apk_runtime_ZIPUtils_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: []
} external;

function org_wso2_apk_runtime_ZIPUtils_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: []
} external;

function org_wso2_apk_runtime_ZIPUtils_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_ZIPUtils_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: ["long", "int"]
} external;

isolated function org_wso2_apk_runtime_ZIPUtils_zipDir(handle arg0, handle arg1) returns error? = @java:Method {
    name: "zipDir",
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_runtime_ZIPUtils_newZIPUtils1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.ZIPUtils",
    paramTypes: []
} external;

