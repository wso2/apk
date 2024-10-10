import config_deployer_service.java.lang as javalang;
import config_deployer_service.org.wso2.apk.config.api as orgwso2apkconfigapi;

import ballerina/jballerina.java;

# Ballerina class mapping for the Java `org.wso2.apk.config.APKConfValidator` class.
@java:Binding {'class: "org.wso2.apk.config.APKConfValidator"}
public isolated distinct class APKConfValidator {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.config.APKConfValidator` object.
    public final handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.config.APKConfValidator` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.config.APKConfValidator` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }

    # The function that maps to the `equals` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_config_APKConfValidator_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_config_APKConfValidator_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_config_APKConfValidator_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.config.APKConfValidator`.
    public function notify() {
        org_wso2_apk_config_APKConfValidator_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.config.APKConfValidator`.
    public function notifyAll() {
        org_wso2_apk_config_APKConfValidator_notifyAll(self.jObj);
    }

    # The function that maps to the `validate` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + return - The `orgwso2apkconfigapi:APKConfValidationResponse` or the `javalang:Exception` value returning from the Java mapping.
    public isolated function validate(string arg0) returns orgwso2apkconfigapi:APKConfValidationResponse|javalang:Exception {
        handle|error externalObj = org_wso2_apk_config_APKConfValidator_validate(self.jObj, java:fromString(arg0));
        if (externalObj is error) {
            javalang:Exception e = error javalang:Exception(javalang:EXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            orgwso2apkconfigapi:APKConfValidationResponse newObj = new (externalObj);
            return newObj;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_APKConfValidator_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_APKConfValidator_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.APKConfValidator`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_APKConfValidator_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.config.APKConfValidator`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + return - The new `APKConfValidator` class or `javalang:Exception` error generated.
public function newAPKConfValidator1(string arg0) returns APKConfValidator|javalang:Exception {
    handle|error externalObj = org_wso2_apk_config_APKConfValidator_newAPKConfValidator1(java:fromString(arg0));
    if (externalObj is error) {
        javalang:Exception e = error javalang:Exception(javalang:EXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        APKConfValidator newObj = new (externalObj);
        return newObj;
    }
}

function org_wso2_apk_config_APKConfValidator_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_config_APKConfValidator_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: []
} external;

function org_wso2_apk_config_APKConfValidator_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: []
} external;

function org_wso2_apk_config_APKConfValidator_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: []
} external;

function org_wso2_apk_config_APKConfValidator_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_APKConfValidator_validate(handle receiver, handle arg0) returns handle|error = @java:Method {
    name: "validate",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_APKConfValidator_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: []
} external;

function org_wso2_apk_config_APKConfValidator_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: ["long"]
} external;

function org_wso2_apk_config_APKConfValidator_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_config_APKConfValidator_newAPKConfValidator1(handle arg0) returns handle|error = @java:Constructor {
    'class: "org.wso2.apk.config.APKConfValidator",
    paramTypes: ["java.lang.String"]
} external;

