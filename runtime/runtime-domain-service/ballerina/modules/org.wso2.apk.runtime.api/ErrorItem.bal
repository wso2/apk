import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.api.ErrorItem` class.
@java:Binding {'class: "org.wso2.apk.runtime.api.ErrorItem"}
public distinct class ErrorItem {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.api.ErrorItem` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.api.ErrorItem` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.api.ErrorItem` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_api_ErrorItem_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_api_ErrorItem_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getErrorCode` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `int` value returning from the Java mapping.
    public isolated function getErrorCode() returns int {
        return org_wso2_apk_runtime_api_ErrorItem_getErrorCode(self.jObj)
    ;
    }

    # The function that maps to the `getErrorDescription` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getErrorDescription() returns string? {
        return java:toString(org_wso2_apk_runtime_api_ErrorItem_getErrorDescription(self.jObj));
    }

    # The function that maps to the `getErrorMessage` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getErrorMessage() returns string? {
        return java:toString(org_wso2_apk_runtime_api_ErrorItem_getErrorMessage(self.jObj));
    }

    # The function that maps to the `getHttpStatusCode` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `int` value returning from the Java mapping.
    public isolated function getHttpStatusCode() returns int {
        return org_wso2_apk_runtime_api_ErrorItem_getHttpStatusCode(self.jObj);
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_api_ErrorItem_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.api.ErrorItem`.
    public function notify() {
        org_wso2_apk_runtime_api_ErrorItem_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.api.ErrorItem`.
    public function notifyAll() {
        org_wso2_apk_runtime_api_ErrorItem_notifyAll(self.jObj);
    }

    # The function that maps to the `printStackTrace` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function printStackTrace() returns boolean {
        return org_wso2_apk_runtime_api_ErrorItem_printStackTrace(self.jObj);
    }

    # The function that maps to the `setDescription` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setDescription(string arg0) {
        org_wso2_apk_runtime_api_ErrorItem_setDescription(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setErrorCode` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setErrorCode(int arg0) {
        org_wso2_apk_runtime_api_ErrorItem_setErrorCode(self.jObj, arg0);
    }

    # The function that maps to the `setMessage` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setMessage(string arg0) {
        org_wso2_apk_runtime_api_ErrorItem_setMessage(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setStatusCode` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setStatusCode(int arg0) {
        org_wso2_apk_runtime_api_ErrorItem_setStatusCode(self.jObj, arg0);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_ErrorItem_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_ErrorItem_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.ErrorItem`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_ErrorItem_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.ErrorItem`.
#
# + return - The new `ErrorItem` class generated.
public function newErrorItem1() returns ErrorItem {
    handle externalObj = org_wso2_apk_runtime_api_ErrorItem_newErrorItem1();
    ErrorItem newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.ErrorItem`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `string` value required to map with the Java constructor parameter.
# + arg2 - The `int` value required to map with the Java constructor parameter.
# + arg3 - The `int` value required to map with the Java constructor parameter.
# + return - The new `ErrorItem` class generated.
public function newErrorItem2(string arg0, string arg1, int arg2, int arg3) returns ErrorItem {
    handle externalObj = org_wso2_apk_runtime_api_ErrorItem_newErrorItem2(java:fromString(arg0), java:fromString(arg1), arg2, arg3);
    ErrorItem newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.ErrorItem`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `string` value required to map with the Java constructor parameter.
# + arg2 - The `int` value required to map with the Java constructor parameter.
# + arg3 - The `int` value required to map with the Java constructor parameter.
# + arg4 - The `boolean` value required to map with the Java constructor parameter.
# + return - The new `ErrorItem` class generated.
public function newErrorItem3(string arg0, string arg1, int arg2, int arg3, boolean arg4) returns ErrorItem {
    handle externalObj = org_wso2_apk_runtime_api_ErrorItem_newErrorItem3(java:fromString(arg0), java:fromString(arg1), arg2, arg3, arg4);
    ErrorItem newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_runtime_api_ErrorItem_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorItem_getErrorCode(handle receiver) returns int = @java:Method {
    name: "getErrorCode",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorItem_getErrorDescription(handle receiver) returns handle = @java:Method {
    name: "getErrorDescription",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorItem_getErrorMessage(handle receiver) returns handle = @java:Method {
    name: "getErrorMessage",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorItem_getHttpStatusCode(handle receiver) returns int = @java:Method {
    name: "getHttpStatusCode",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_printStackTrace(handle receiver) returns boolean = @java:Method {
    name: "printStackTrace",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_setDescription(handle receiver, handle arg0) = @java:Method {
    name: "setDescription",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_setErrorCode(handle receiver, int arg0) = @java:Method {
    name: "setErrorCode",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_setMessage(handle receiver, handle arg0) = @java:Method {
    name: "setMessage",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_setStatusCode(handle receiver, int arg0) = @java:Method {
    name: "setStatusCode",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["int"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_newErrorItem1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_ErrorItem_newErrorItem2(handle arg0, handle arg1, int arg2, int arg3) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["java.lang.String", "java.lang.String", "long", "int"]
} external;

function org_wso2_apk_runtime_api_ErrorItem_newErrorItem3(handle arg0, handle arg1, int arg2, int arg3, boolean arg4) returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.ErrorItem",
    paramTypes: ["java.lang.String", "java.lang.String", "long", "int", "boolean"]
} external;

