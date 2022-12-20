import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.util as javautil;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.model.OperationPolicy` class.
@java:Binding {'class: "org.wso2.apk.runtime.model.OperationPolicy"}
public distinct class OperationPolicy {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.model.OperationPolicy` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.model.OperationPolicy` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.model.OperationPolicy` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null"
    ;
    }

    # The function that maps to the `compareTo` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `OperationPolicy` value required to map with the Java method parameter.
    # + return - The `int` value returning from the Java mapping.
    public function compareTo(OperationPolicy arg0) returns int {
        return org_wso2_apk_runtime_model_OperationPolicy_compareTo(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_model_OperationPolicy_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_model_OperationPolicy_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getDirection` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getDirection() returns string? {
        return java:toString(org_wso2_apk_runtime_model_OperationPolicy_getDirection(self.jObj));
    }

    # The function that maps to the `getOrder` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function getOrder() returns int {
        return org_wso2_apk_runtime_model_OperationPolicy_getOrder(self.jObj);
    }

    # The function that maps to the `getParameters` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `javautil:Map` value returning from the Java mapping.
    public function getParameters() returns javautil:Map {
        handle externalObj = org_wso2_apk_runtime_model_OperationPolicy_getParameters(self.jObj);
        javautil:Map newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getPolicyId` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getPolicyId() returns string? {
        return java:toString(org_wso2_apk_runtime_model_OperationPolicy_getPolicyId(self.jObj));
    }

    # The function that maps to the `getPolicyName` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getPolicyName() returns string? {
        return java:toString(org_wso2_apk_runtime_model_OperationPolicy_getPolicyName(self.jObj));
    }

    # The function that maps to the `getPolicyVersion` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getPolicyVersion() returns string? {
        return java:toString(org_wso2_apk_runtime_model_OperationPolicy_getPolicyVersion(self.jObj));
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_model_OperationPolicy_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    public function notify() {
        org_wso2_apk_runtime_model_OperationPolicy_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    public function notifyAll() {
        org_wso2_apk_runtime_model_OperationPolicy_notifyAll(self.jObj);
    }

    # The function that maps to the `setDirection` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setDirection(string arg0) {
        org_wso2_apk_runtime_model_OperationPolicy_setDirection(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setOrder` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setOrder(int arg0) {
        org_wso2_apk_runtime_model_OperationPolicy_setOrder(self.jObj, arg0);
    }

    # The function that maps to the `setParameters` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `javautil:Map` value required to map with the Java method parameter.
    public function setParameters(javautil:Map arg0) {
        org_wso2_apk_runtime_model_OperationPolicy_setParameters(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setPolicyId` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setPolicyId(string arg0) {
        org_wso2_apk_runtime_model_OperationPolicy_setPolicyId(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setPolicyName` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setPolicyName(string arg0) {
        org_wso2_apk_runtime_model_OperationPolicy_setPolicyName(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setPolicyVersion` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setPolicyVersion(string arg0) {
        org_wso2_apk_runtime_model_OperationPolicy_setPolicyVersion(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_OperationPolicy_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_OperationPolicy_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.OperationPolicy`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_OperationPolicy_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.model.OperationPolicy`.
#
# + return - The new `OperationPolicy` class generated.
public function newOperationPolicy1() returns OperationPolicy {
    handle externalObj = org_wso2_apk_runtime_model_OperationPolicy_newOperationPolicy1();
    OperationPolicy newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_runtime_model_OperationPolicy_compareTo(handle receiver, handle arg0) returns int = @java:Method {
    name: "compareTo",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["org.wso2.apk.runtime.model.OperationPolicy"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getDirection(handle receiver) returns handle = @java:Method {
    name: "getDirection",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getOrder(handle receiver) returns int = @java:Method {
    name: "getOrder",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getParameters(handle receiver) returns handle = @java:Method {
    name: "getParameters",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getPolicyId(handle receiver) returns handle = @java:Method {
    name: "getPolicyId",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getPolicyName(handle receiver) returns handle = @java:Method {
    name: "getPolicyName",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_getPolicyVersion(handle receiver) returns handle = @java:Method {
    name: "getPolicyVersion",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_setDirection(handle receiver, handle arg0) = @java:Method {
    name: "setDirection",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_setOrder(handle receiver, int arg0) = @java:Method {
    name: "setOrder",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["int"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_setParameters(handle receiver, handle arg0) = @java:Method {
    name: "setParameters",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["java.util.Map"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_setPolicyId(handle receiver, handle arg0) = @java:Method {
    name: "setPolicyId",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_setPolicyName(handle receiver, handle arg0) = @java:Method {
    name: "setPolicyName",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_setPolicyVersion(handle receiver, handle arg0) = @java:Method {
    name: "setPolicyVersion",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_OperationPolicy_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_runtime_model_OperationPolicy_newOperationPolicy1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.model.OperationPolicy",
    paramTypes: []
} external;

