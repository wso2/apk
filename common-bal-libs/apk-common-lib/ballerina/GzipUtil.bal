import apk_common_lib.java.io as javaio;
import apk_common_lib.java.lang as javalang;

import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;

# Ballerina class mapping for the Java `org.wso2.apk.common.GzipUtil` class.
@java:Binding {'class: "org.wso2.apk.common.GzipUtil"}
public distinct class GzipUtil {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.common.GzipUtil` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.common.GzipUtil` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.common.GzipUtil` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }

    # The function that maps to the `equals` method of `org.wso2.apk.common.GzipUtil`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_common_GzipUtil_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.common.GzipUtil`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_common_GzipUtil_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.common.GzipUtil`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_common_GzipUtil_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.common.GzipUtil`.
    public function notify() {
        org_wso2_apk_common_GzipUtil_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.common.GzipUtil`.
    public function notifyAll() {
        org_wso2_apk_common_GzipUtil_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.common.GzipUtil`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_common_GzipUtil_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.common.GzipUtil`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_common_GzipUtil_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.common.GzipUtil`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_common_GzipUtil_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.common.GzipUtil`.
#
# + return - The new `GzipUtil` class generated.
public isolated function newGzipUtil1() returns GzipUtil {
    handle externalObj = org_wso2_apk_common_GzipUtil_newGzipUtil1();
    GzipUtil newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `compressGzipFile` method of `org.wso2.apk.common.GzipUtil`.
#
# + arg0 - The `byte[]` value required to map with the Java method parameter.
# + return - The `byte[]` or the `javaio:IOException` value returning from the Java mapping.
public isolated function GzipUtil_compressGzipFile(byte[] arg0) returns byte[]|javaio:IOException|error {
    handle|error externalObj = org_wso2_apk_common_GzipUtil_compressGzipFile(check jarrays:toHandle(arg0, "byte"));
    if (externalObj is error) {
        javaio:IOException e = error javaio:IOException(javaio:IOEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return <byte[]>check jarrays:fromHandle(externalObj, "byte");
    }
}

# The function that maps to the `compressGzipFile` method of `org.wso2.apk.common.GzipUtil`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `javaio:IOException` value returning from the Java mapping.
public isolated function GzipUtil_compressGzipFile2(string arg0, string arg1) returns javaio:IOException? {
    error|() externalObj = org_wso2_apk_common_GzipUtil_compressGzipFile2(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        javaio:IOException e = error javaio:IOException(javaio:IOEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `decompressGzipFile` method of `org.wso2.apk.common.GzipUtil`.
#
# + arg0 - The `byte[]` value required to map with the Java method parameter.
# + return - The `byte[]` or the `javaio:IOException` value returning from the Java mapping.
public isolated function GzipUtil_decompressGzipFile(byte[] arg0) returns byte[]|javaio:IOException|error {
    handle|error externalObj = org_wso2_apk_common_GzipUtil_decompressGzipFile(check jarrays:toHandle(arg0, "byte"));
    if (externalObj is error) {
        javaio:IOException e = error javaio:IOException(javaio:IOEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return <byte[]>check jarrays:fromHandle(externalObj, "byte");
    }
}

# The function that maps to the `decompressGzipFile` method of `org.wso2.apk.common.GzipUtil`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `javaio:IOException` value returning from the Java mapping.
public isolated function GzipUtil_decompressGzipFile2(string arg0, string arg1) returns javaio:IOException? {
    error|() externalObj = org_wso2_apk_common_GzipUtil_decompressGzipFile2(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        javaio:IOException e = error javaio:IOException(javaio:IOEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

isolated function org_wso2_apk_common_GzipUtil_compressGzipFile(handle arg0) returns handle|error = @java:Method {
    name: "compressGzipFile",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["[B"]
} external;

isolated function org_wso2_apk_common_GzipUtil_compressGzipFile2(handle arg0, handle arg1) returns error? = @java:Method {
    name: "compressGzipFile",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

isolated function org_wso2_apk_common_GzipUtil_decompressGzipFile(handle arg0) returns handle|error = @java:Method {
    name: "decompressGzipFile",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["[B"]
} external;

isolated function org_wso2_apk_common_GzipUtil_decompressGzipFile2(handle arg0, handle arg1) returns error? = @java:Method {
    name: "decompressGzipFile",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_common_GzipUtil_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_common_GzipUtil_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_GzipUtil_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_GzipUtil_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_GzipUtil_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_GzipUtil_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_GzipUtil_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["long"]
} external;

function org_wso2_apk_common_GzipUtil_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: ["long", "int"]
} external;

isolated function org_wso2_apk_common_GzipUtil_newGzipUtil1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.common.GzipUtil",
    paramTypes: []
} external;

