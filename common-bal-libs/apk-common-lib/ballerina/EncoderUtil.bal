import apk_common_lib.java.lang as javalang;

import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;

# Ballerina class mapping for the Java `org.wso2.apk.common.EncoderUtil` class.
@java:Binding {'class: "org.wso2.apk.common.EncoderUtil"}
public distinct class EncoderUtil {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.common.EncoderUtil` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.common.EncoderUtil` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.common.EncoderUtil` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }

    # The function that maps to the `equals` method of `org.wso2.apk.common.EncoderUtil`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_common_EncoderUtil_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.common.EncoderUtil`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_common_EncoderUtil_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.common.EncoderUtil`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_common_EncoderUtil_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.common.EncoderUtil`.
    public function notify() {
        org_wso2_apk_common_EncoderUtil_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.common.EncoderUtil`.
    public function notifyAll() {
        org_wso2_apk_common_EncoderUtil_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.common.EncoderUtil`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_common_EncoderUtil_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.common.EncoderUtil`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_common_EncoderUtil_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.common.EncoderUtil`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_common_EncoderUtil_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.common.EncoderUtil`.
#
# + return - The new `EncoderUtil` class generated.
public isolated function newEncoderUtil1() returns EncoderUtil {
    handle externalObj = org_wso2_apk_common_EncoderUtil_newEncoderUtil1();
    EncoderUtil newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `decodeBase64` method of `org.wso2.apk.common.EncoderUtil`.
#
# + arg0 - The `byte[]` value required to map with the Java method parameter.
# + return - The `byte[]` value returning from the Java mapping.
public isolated function EncoderUtil_decodeBase64(byte[] arg0) returns byte[]|error {
    handle externalObj = org_wso2_apk_common_EncoderUtil_decodeBase64(check jarrays:toHandle(arg0, "byte"));
    return <byte[]>check jarrays:fromHandle(externalObj, "byte");
}

# The function that maps to the `encodeBase64` method of `org.wso2.apk.common.EncoderUtil`.
#
# + arg0 - The `byte[]` value required to map with the Java method parameter.
# + return - The `byte[]` value returning from the Java mapping.
public isolated function EncoderUtil_encodeBase64(byte[] arg0) returns byte[]|error {
    handle externalObj = org_wso2_apk_common_EncoderUtil_encodeBase64(check jarrays:toHandle(arg0, "byte"));
    return <byte[]>check jarrays:fromHandle(externalObj, "byte");
}

# The function that maps to the `encodeBase64` method of `org.wso2.apk.common.EncoderUtil`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `byte[]` value returning from the Java mapping.
public isolated function EncoderUtil_encodeBase642(string arg0) returns byte[]|error {
    handle externalObj = org_wso2_apk_common_EncoderUtil_encodeBase642(java:fromString(arg0));
    return <byte[]>check jarrays:fromHandle(externalObj, "byte");
}

isolated function org_wso2_apk_common_EncoderUtil_decodeBase64(handle arg0) returns handle = @java:Method {
    name: "decodeBase64",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: ["[B"]
} external;

isolated function org_wso2_apk_common_EncoderUtil_encodeBase64(handle arg0) returns handle = @java:Method {
    name: "encodeBase64",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: ["[B"]
} external;

isolated function org_wso2_apk_common_EncoderUtil_encodeBase642(handle arg0) returns handle = @java:Method {
    name: "encodeBase64",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_common_EncoderUtil_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_common_EncoderUtil_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_EncoderUtil_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_EncoderUtil_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_EncoderUtil_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_EncoderUtil_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: []
} external;

function org_wso2_apk_common_EncoderUtil_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: ["long"]
} external;

function org_wso2_apk_common_EncoderUtil_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: ["long", "int"]
} external;

isolated function org_wso2_apk_common_EncoderUtil_newEncoderUtil1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.common.EncoderUtil",
    paramTypes: []
} external;

