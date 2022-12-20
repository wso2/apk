import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.io as javaio;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.api.APIManagementException` class.
@java:Binding {'class: "org.wso2.apk.runtime.api.APIManagementException"}
public distinct class JAPIManagementException {

    *java:JObject;
    *javalang:JException;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.api.APIManagementException` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.api.APIManagementException` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.api.APIManagementException` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `addSuppressed` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `javalang:Throwable` value required to map with the Java method parameter.
    public function addSuppressed(javalang:Throwable arg0) {
        org_wso2_apk_runtime_api_APIManagementException_addSuppressed(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_api_APIManagementException_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `fillInStackTrace` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `javalang:Throwable` value returning from the Java mapping.
    public function fillInStackTrace() returns javalang:Throwable {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_fillInStackTrace(self.jObj);
        javalang:Throwable newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getCause` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `javalang:Throwable` value returning from the Java mapping.
    public function getCause() returns javalang:Throwable {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_getCause(self.jObj);
        javalang:Throwable newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getErrorHandler` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `ErrorHandler` value returning from the Java mapping.
    public isolated function getErrorHandler() returns ErrorHandler {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_getErrorHandler(self.jObj);
        ErrorHandler newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getLocalizedMessage` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getLocalizedMessage() returns string? {
        return java:toString(org_wso2_apk_runtime_api_APIManagementException_getLocalizedMessage(self.jObj));
    }

    # The function that maps to the `getMessage` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getMessage() returns string? {
        return java:toString(org_wso2_apk_runtime_api_APIManagementException_getMessage(self.jObj));
    }

    # The function that maps to the `getStackTrace` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `javalang:StackTraceElement[]` value returning from the Java mapping.
    public isolated function getStackTrace() returns javalang:StackTraceElement[]|error {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_getStackTrace(self.jObj);
        javalang:StackTraceElement[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:StackTraceElement element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `getSuppressed` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `javalang:Throwable[]` value returning from the Java mapping.
    public function getSuppressed() returns javalang:Throwable[]|error {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_getSuppressed(self.jObj);
        javalang:Throwable[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            javalang:Throwable element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_api_APIManagementException_hashCode(self.jObj);
    }

    # The function that maps to the `initCause` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `javalang:Throwable` value required to map with the Java method parameter.
    # + return - The `javalang:Throwable` value returning from the Java mapping.
    public function initCause(javalang:Throwable arg0) returns javalang:Throwable {
        handle externalObj = org_wso2_apk_runtime_api_APIManagementException_initCause(self.jObj, arg0.jObj);
        javalang:Throwable newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.api.APIManagementException`.
    public function notify() {
        org_wso2_apk_runtime_api_APIManagementException_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.api.APIManagementException`.
    public function notifyAll() {
        org_wso2_apk_runtime_api_APIManagementException_notifyAll(self.jObj);
    }

    # The function that maps to the `printStackTrace` method of `org.wso2.apk.runtime.api.APIManagementException`.
    public function printStackTrace() {
        org_wso2_apk_runtime_api_APIManagementException_printStackTrace(self.jObj);
    }

    # The function that maps to the `printStackTrace` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `javaio:PrintStream` value required to map with the Java method parameter.
    public function printStackTrace2(javaio:PrintStream arg0) {
        org_wso2_apk_runtime_api_APIManagementException_printStackTrace2(self.jObj, arg0.jObj);
    }

    # The function that maps to the `printStackTrace` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `javaio:PrintWriter` value required to map with the Java method parameter.
    public function printStackTrace3(javaio:PrintWriter arg0) {
        org_wso2_apk_runtime_api_APIManagementException_printStackTrace3(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setStackTrace` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `javalang:StackTraceElement[]` value required to map with the Java method parameter.
    # + return - The `error?` value returning from the Java mapping.
    public function setStackTrace(javalang:StackTraceElement[] arg0) returns error? {
        org_wso2_apk_runtime_api_APIManagementException_setStackTrace(self.jObj, check jarrays:toHandle(arg0, "java.lang.StackTraceElement"));
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_APIManagementException_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_APIManagementException_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.APIManagementException`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_APIManagementException_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIManagementException`.
#
# + arg0 - The `ErrorHandler` value required to map with the Java constructor parameter.
# + return - The new `JAPIManagementException` class generated.
public function newJAPIManagementException1(ErrorHandler arg0) returns JAPIManagementException {
    handle externalObj = org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException1(arg0.jObj);
    JAPIManagementException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIManagementException`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + return - The new `JAPIManagementException` class generated.
public function newJAPIManagementException2(string arg0) returns JAPIManagementException {
    handle externalObj = org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException2(java:fromString(arg0));
    JAPIManagementException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIManagementException`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `ErrorHandler` value required to map with the Java constructor parameter.
# + return - The new `JAPIManagementException` class generated.
public function newJAPIManagementException3(string arg0, ErrorHandler arg1) returns JAPIManagementException {
    handle externalObj = org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException3(java:fromString(arg0), arg1.jObj);
    JAPIManagementException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIManagementException`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `javalang:Throwable` value required to map with the Java constructor parameter.
# + return - The new `JAPIManagementException` class generated.
public function newJAPIManagementException4(string arg0, javalang:Throwable arg1) returns JAPIManagementException {
    handle externalObj = org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException4(java:fromString(arg0), arg1.jObj);
    JAPIManagementException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIManagementException`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `javalang:Throwable` value required to map with the Java constructor parameter.
# + arg2 - The `ErrorHandler` value required to map with the Java constructor parameter.
# + return - The new `JAPIManagementException` class generated.
public function newJAPIManagementException5(string arg0, javalang:Throwable arg1, ErrorHandler arg2) returns JAPIManagementException {
    handle externalObj = org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException5(java:fromString(arg0), arg1.jObj, arg2.jObj);
    JAPIManagementException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIManagementException`.
#
# + arg0 - The `javalang:Throwable` value required to map with the Java constructor parameter.
# + return - The new `JAPIManagementException` class generated.
public function newJAPIManagementException6(javalang:Throwable arg0) returns JAPIManagementException {
    handle externalObj = org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException6(arg0.jObj);
    JAPIManagementException newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_runtime_api_APIManagementException_addSuppressed(handle receiver, handle arg0) = @java:Method {
    name: "addSuppressed",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.Throwable"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_fillInStackTrace(handle receiver) returns handle = @java:Method {
    name: "fillInStackTrace",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_getCause(handle receiver) returns handle = @java:Method {
    name: "getCause",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_APIManagementException_getErrorHandler(handle receiver) returns handle = @java:Method {
    name: "getErrorHandler",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_APIManagementException_getLocalizedMessage(handle receiver) returns handle = @java:Method {
    name: "getLocalizedMessage",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_APIManagementException_getMessage(handle receiver) returns handle = @java:Method {
    name: "getMessage",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_APIManagementException_getStackTrace(handle receiver) returns handle = @java:Method {
    name: "getStackTrace",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_getSuppressed(handle receiver) returns handle = @java:Method {
    name: "getSuppressed",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_initCause(handle receiver, handle arg0) returns handle = @java:Method {
    name: "initCause",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.Throwable"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_printStackTrace(handle receiver) = @java:Method {
    name: "printStackTrace",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_printStackTrace2(handle receiver, handle arg0) = @java:Method {
    name: "printStackTrace",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.io.PrintStream"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_printStackTrace3(handle receiver, handle arg0) = @java:Method {
    name: "printStackTrace",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.io.PrintWriter"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_setStackTrace(handle receiver, handle arg0) = @java:Method {
    name: "setStackTrace",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["[Ljava.lang.StackTraceElement;"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIManagementException_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException1(handle arg0) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["org.wso2.apk.runtime.api.ErrorHandler"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException2(handle arg0) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException3(handle arg0, handle arg1) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.String", "org.wso2.apk.runtime.api.ErrorHandler"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException4(handle arg0, handle arg1) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.String", "java.lang.Throwable"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException5(handle arg0, handle arg1, handle arg2) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.String", "java.lang.Throwable", "org.wso2.apk.runtime.api.ErrorHandler"]
} external;

function org_wso2_apk_runtime_api_APIManagementException_newJAPIManagementException6(handle arg0) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIManagementException",
    paramTypes: ["java.lang.Throwable"]
} external;

