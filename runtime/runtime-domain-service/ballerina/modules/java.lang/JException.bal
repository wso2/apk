import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import runtime_domain_service.java.io as javaio;

# Ballerina class mapping for the Java `java.lang.Exception` class.
@java:Binding {'class: "java.lang.Exception"}
public distinct class JException {

    *java:JObject;
    *Throwable;

    # The `handle` field that stores the reference to the `java.lang.Exception` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `java.lang.Exception` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `java.lang.Exception` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `addSuppressed` method of `java.lang.Exception`.
    #
    # + arg0 - The `Throwable` value required to map with the Java method parameter.
    public function addSuppressed(Throwable arg0) {
        java_lang_Exception_addSuppressed(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `java.lang.Exception`.
    #
    # + arg0 - The `Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(Object arg0) returns boolean {
        return java_lang_Exception_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `fillInStackTrace` method of `java.lang.Exception`.
    #
    # + return - The `Throwable` value returning from the Java mapping.
    public function fillInStackTrace() returns Throwable {
        handle externalObj = java_lang_Exception_fillInStackTrace(self.jObj);
        Throwable newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getCause` method of `java.lang.Exception`.
    #
    # + return - The `Throwable` value returning from the Java mapping.
    public function getCause() returns Throwable {
        handle externalObj = java_lang_Exception_getCause(self.jObj);
        Throwable newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getClass` method of `java.lang.Exception`.
    #
    # + return - The `Class` value returning from the Java mapping.
    public function getClass() returns Class {
        handle externalObj = java_lang_Exception_getClass(self.jObj);
        Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getLocalizedMessage` method of `java.lang.Exception`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getLocalizedMessage() returns string? {
        return java:toString(java_lang_Exception_getLocalizedMessage(self.jObj));
    }

    # The function that maps to the `getMessage` method of `java.lang.Exception`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getMessage() returns string? {
        return java:toString(java_lang_Exception_getMessage(self.jObj));
    }

    # The function that maps to the `getStackTrace` method of `java.lang.Exception`.
    #
    # + return - The `StackTraceElement[]` value returning from the Java mapping.
    public function getStackTrace() returns StackTraceElement[]|error {
        handle externalObj = java_lang_Exception_getStackTrace(self.jObj);
        StackTraceElement[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            StackTraceElement element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `getSuppressed` method of `java.lang.Exception`.
    #
    # + return - The `Throwable[]` value returning from the Java mapping.
    public function getSuppressed() returns Throwable[]|error {
        handle externalObj = java_lang_Exception_getSuppressed(self.jObj);
        Throwable[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            Throwable element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `hashCode` method of `java.lang.Exception`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return java_lang_Exception_hashCode(self.jObj);
    }

    # The function that maps to the `initCause` method of `java.lang.Exception`.
    #
    # + arg0 - The `Throwable` value required to map with the Java method parameter.
    # + return - The `Throwable` value returning from the Java mapping.
    public function initCause(Throwable arg0) returns Throwable {
        handle externalObj = java_lang_Exception_initCause(self.jObj, arg0.jObj);
        Throwable newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `notify` method of `java.lang.Exception`.
    public function notify() {
        java_lang_Exception_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `java.lang.Exception`.
    public function notifyAll() {
        java_lang_Exception_notifyAll(self.jObj);
    }

    # The function that maps to the `printStackTrace` method of `java.lang.Exception`.
    public function printStackTrace() {
        java_lang_Exception_printStackTrace(self.jObj);
    }

    # The function that maps to the `printStackTrace` method of `java.lang.Exception`.
    #
    # + arg0 - The `javaio:PrintStream` value required to map with the Java method parameter.
    public function printStackTrace2(javaio:PrintStream arg0) {
        java_lang_Exception_printStackTrace2(self.jObj, arg0.jObj);
    }

    # The function that maps to the `printStackTrace` method of `java.lang.Exception`.
    #
    # + arg0 - The `javaio:PrintWriter` value required to map with the Java method parameter.
    public function printStackTrace3(javaio:PrintWriter arg0) {
        java_lang_Exception_printStackTrace3(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setStackTrace` method of `java.lang.Exception`.
    #
    # + arg0 - The `StackTraceElement[]` value required to map with the Java method parameter.
    # + return - The `error?` value returning from the Java mapping.
    public function setStackTrace(StackTraceElement[] arg0) returns error? {
        java_lang_Exception_setStackTrace(self.jObj, check jarrays:toHandle(arg0, "java.lang.StackTraceElement"));
    }

    # The function that maps to the `wait` method of `java.lang.Exception`.
    #
    # + return - The `InterruptedException` value returning from the Java mapping.
    public function 'wait() returns InterruptedException? {
        error|() externalObj = java_lang_Exception_wait(self.jObj);
        if (externalObj is error) {
            InterruptedException e = error InterruptedException(INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.lang.Exception`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns InterruptedException? {
        error|() externalObj = java_lang_Exception_wait2(self.jObj, arg0);
        if (externalObj is error) {
            InterruptedException e = error InterruptedException(INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `java.lang.Exception`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns InterruptedException? {
        error|() externalObj = java_lang_Exception_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            InterruptedException e = error InterruptedException(INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `java.lang.Exception`.
#
# + return - The new `JException` class generated.
public function newJException1() returns JException {
    handle externalObj = java_lang_Exception_newJException1();
    JException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `java.lang.Exception`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + return - The new `JException` class generated.
public function newJException2(string arg0) returns JException {
    handle externalObj = java_lang_Exception_newJException2(java:fromString(arg0));
    JException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `java.lang.Exception`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `Throwable` value required to map with the Java constructor parameter.
# + return - The new `JException` class generated.
public function newJException3(string arg0, Throwable arg1) returns JException {
    handle externalObj = java_lang_Exception_newJException3(java:fromString(arg0), arg1.jObj);
    JException newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `java.lang.Exception`.
#
# + arg0 - The `Throwable` value required to map with the Java constructor parameter.
# + return - The new `JException` class generated.
public function newJException4(Throwable arg0) returns JException {
    handle externalObj = java_lang_Exception_newJException4(arg0.jObj);
    JException newObj = new (externalObj);
    return newObj;
}

function java_lang_Exception_addSuppressed(handle receiver, handle arg0) = @java:Method {
    name: "addSuppressed",
    'class: "java.lang.Exception",
    paramTypes: ["java.lang.Throwable"]
} external;

function java_lang_Exception_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "java.lang.Exception",
    paramTypes: ["java.lang.Object"]
} external;

function java_lang_Exception_fillInStackTrace(handle receiver) returns handle = @java:Method {
    name: "fillInStackTrace",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_getCause(handle receiver) returns handle = @java:Method {
    name: "getCause",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_getLocalizedMessage(handle receiver) returns handle = @java:Method {
    name: "getLocalizedMessage",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_getMessage(handle receiver) returns handle = @java:Method {
    name: "getMessage",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_getStackTrace(handle receiver) returns handle = @java:Method {
    name: "getStackTrace",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_getSuppressed(handle receiver) returns handle = @java:Method {
    name: "getSuppressed",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_initCause(handle receiver, handle arg0) returns handle = @java:Method {
    name: "initCause",
    'class: "java.lang.Exception",
    paramTypes: ["java.lang.Throwable"]
} external;

function java_lang_Exception_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_printStackTrace(handle receiver) = @java:Method {
    name: "printStackTrace",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_printStackTrace2(handle receiver, handle arg0) = @java:Method {
    name: "printStackTrace",
    'class: "java.lang.Exception",
    paramTypes: ["java.io.PrintStream"]
} external;

function java_lang_Exception_printStackTrace3(handle receiver, handle arg0) = @java:Method {
    name: "printStackTrace",
    'class: "java.lang.Exception",
    paramTypes: ["java.io.PrintWriter"]
} external;

function java_lang_Exception_setStackTrace(handle receiver, handle arg0) = @java:Method {
    name: "setStackTrace",
    'class: "java.lang.Exception",
    paramTypes: ["[Ljava.lang.StackTraceElement;"]
} external;

function java_lang_Exception_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "java.lang.Exception",
    paramTypes: ["long"]
} external;

function java_lang_Exception_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "java.lang.Exception",
    paramTypes: ["long", "int"]
} external;

function java_lang_Exception_newJException1() returns handle = @java:Constructor {
    'class: "java.lang.Exception",
    paramTypes: []
} external;

function java_lang_Exception_newJException2(handle arg0) returns handle = @java:Constructor {
    'class: "java.lang.Exception",
    paramTypes: ["java.lang.String"]
} external;

function java_lang_Exception_newJException3(handle arg0, handle arg1) returns handle = @java:Constructor {
    'class: "java.lang.Exception",
    paramTypes: ["java.lang.String", "java.lang.Throwable"]
} external;

function java_lang_Exception_newJException4(handle arg0) returns handle = @java:Constructor {
    'class: "java.lang.Exception",
    paramTypes: ["java.lang.Throwable"]
} external;

