import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.model.Scope` class.
@java:Binding {'class: "org.wso2.apk.runtime.model.Scope"}
public distinct class Scope {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.model.Scope` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.model.Scope` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.model.Scope` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_model_Scope_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_model_Scope_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getDescription` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getDescription() returns string? {
        return java:toString(org_wso2_apk_runtime_model_Scope_getDescription(self.jObj));
    }

    # The function that maps to the `getId` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getId() returns string? {
        return java:toString(org_wso2_apk_runtime_model_Scope_getId(self.jObj));
    }

    # The function that maps to the `getKey` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getKey() returns string? {
        return java:toString(org_wso2_apk_runtime_model_Scope_getKey(self.jObj));
    }

    # The function that maps to the `getName` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getName() returns string? {
        return java:toString(org_wso2_apk_runtime_model_Scope_getName(self.jObj));
    }

    # The function that maps to the `getRoles` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getRoles() returns string? {
        return java:toString(org_wso2_apk_runtime_model_Scope_getRoles(self.jObj));
    }

    # The function that maps to the `getUsageCount` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function getUsageCount() returns int {
        return org_wso2_apk_runtime_model_Scope_getUsageCount(self.jObj);
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_model_Scope_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.model.Scope`.
    public function notify() {
        org_wso2_apk_runtime_model_Scope_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.model.Scope`.
    public function notifyAll() {
        org_wso2_apk_runtime_model_Scope_notifyAll(self.jObj);
    }

    # The function that maps to the `setDescription` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setDescription(string arg0) {
        org_wso2_apk_runtime_model_Scope_setDescription(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setId` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setId(string arg0) {
        org_wso2_apk_runtime_model_Scope_setId(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setKey` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setKey(string arg0) {
        org_wso2_apk_runtime_model_Scope_setKey(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setName` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setName(string arg0) {
        org_wso2_apk_runtime_model_Scope_setName(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setRoles` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setRoles(string arg0) {
        org_wso2_apk_runtime_model_Scope_setRoles(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setUsageCount` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setUsageCount(int arg0) {
        org_wso2_apk_runtime_model_Scope_setUsageCount(self.jObj, arg0);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_Scope_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_Scope_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.Scope`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_Scope_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.model.Scope`.
#
# + return - The new `Scope` class generated.
public function newScope1() returns Scope {
    handle externalObj = org_wso2_apk_runtime_model_Scope_newScope1();
    Scope newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_runtime_model_Scope_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_model_Scope_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_getDescription(handle receiver) returns handle = @java:Method {
    name: "getDescription",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_getId(handle receiver) returns handle = @java:Method {
    name: "getId",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_getKey(handle receiver) returns handle = @java:Method {
    name: "getKey",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_getName(handle receiver) returns handle = @java:Method {
    name: "getName",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_getRoles(handle receiver) returns handle = @java:Method {
    name: "getRoles",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_getUsageCount(handle receiver) returns int = @java:Method {
    name: "getUsageCount",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_setDescription(handle receiver, handle arg0) = @java:Method {
    name: "setDescription",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_Scope_setId(handle receiver, handle arg0) = @java:Method {
    name: "setId",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_Scope_setKey(handle receiver, handle arg0) = @java:Method {
    name: "setKey",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_Scope_setName(handle receiver, handle arg0) = @java:Method {
    name: "setName",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_Scope_setRoles(handle receiver, handle arg0) = @java:Method {
    name: "setRoles",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_Scope_setUsageCount(handle receiver, int arg0) = @java:Method {
    name: "setUsageCount",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["int"]
} external;

function org_wso2_apk_runtime_model_Scope_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_Scope_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_model_Scope_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_runtime_model_Scope_newScope1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.model.Scope",
    paramTypes: []
} external;

