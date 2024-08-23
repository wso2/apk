import config_deployer_service.java.lang as javalang;
import config_deployer_service.java.util as javautil;
import config_deployer_service.org.wso2.apk.config.api as orgwso2apkconfigapi;
import config_deployer_service.org.wso2.apk.config.model as orgwso2apkconfigmodel;

import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;

# Ballerina class mapping for the Java `org.wso2.apk.config.RuntimeAPICommonUtil` class.
@java:Binding {'class: "org.wso2.apk.config.RuntimeAPICommonUtil"}
public distinct class RuntimeAPICommonUtil {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.config.RuntimeAPICommonUtil` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.config.RuntimeAPICommonUtil` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.config.RuntimeAPICommonUtil` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_config_RuntimeAPICommonUtil_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_config_RuntimeAPICommonUtil_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    public function notify() {
        org_wso2_apk_config_RuntimeAPICommonUtil_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    public function notifyAll() {
        org_wso2_apk_config_RuntimeAPICommonUtil_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The function that maps to the `generateDefinition` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
#
# + arg0 - The `orgwso2apkconfigmodel:API` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkconfigapi:APIManagementException` value returning from the Java mapping.
public isolated function RuntimeAPICommonUtil_generateDefinition(orgwso2apkconfigmodel:API arg0) returns orgwso2apkconfigapi:APIManagementException|string? {
    handle|error externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_generateDefinition(arg0.jObj);
    if (externalObj is error) {
        orgwso2apkconfigapi:APIManagementException e = error orgwso2apkconfigapi:APIManagementException(orgwso2apkconfigapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `generateDefinition` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
#
# + arg0 - The `orgwso2apkconfigmodel:API` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkconfigapi:APIManagementException` value returning from the Java mapping.
public isolated function RuntimeAPICommonUtil_generateDefinition2(orgwso2apkconfigmodel:API arg0, string arg1) returns orgwso2apkconfigapi:APIManagementException|string? {
    handle|error externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_generateDefinition2(arg0.jObj, java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkconfigapi:APIManagementException e = error orgwso2apkconfigapi:APIManagementException(orgwso2apkconfigapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `generateUriTemplatesFromAPIDefinition` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `javautil:Set` or the `orgwso2apkconfigapi:APIManagementException` value returning from the Java mapping.
public function RuntimeAPICommonUtil_generateUriTemplatesFromAPIDefinition(string arg0, string arg1) returns javautil:Set|orgwso2apkconfigapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_generateUriTemplatesFromAPIDefinition(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkconfigapi:APIManagementException e = error orgwso2apkconfigapi:APIManagementException(orgwso2apkconfigapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        javautil:Set newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAPIFromDefinition` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkconfigmodel:API` or the `orgwso2apkconfigapi:APIManagementException` value returning from the Java mapping.
public isolated function RuntimeAPICommonUtil_getAPIFromDefinition(string arg0, string arg1) returns orgwso2apkconfigmodel:API|orgwso2apkconfigapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_getAPIFromDefinition(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkconfigapi:APIManagementException e = error orgwso2apkconfigapi:APIManagementException(orgwso2apkconfigapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkconfigmodel:API newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `validateOpenAPIDefinition` method of `org.wso2.apk.config.RuntimeAPICommonUtil`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `byte[]` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + arg4 - The `boolean` value required to map with the Java method parameter.
# + return - The `orgwso2apkconfigapi:APIDefinitionValidationResponse` or the `orgwso2apkconfigapi:APIManagementException` value returning from the Java mapping.
public isolated function RuntimeAPICommonUtil_validateOpenAPIDefinition(string arg0, byte[] arg1, string arg2, string arg3, boolean arg4) returns orgwso2apkconfigapi:APIDefinitionValidationResponse|orgwso2apkconfigapi:APIManagementException|error {
    handle|error externalObj = org_wso2_apk_config_RuntimeAPICommonUtil_validateOpenAPIDefinition(java:fromString(arg0), check jarrays:toHandle(arg1, "byte"), java:fromString(arg2), java:fromString(arg3), arg4);
    if (externalObj is error) {
        orgwso2apkconfigapi:APIManagementException e = error orgwso2apkconfigapi:APIManagementException(orgwso2apkconfigapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkconfigapi:APIDefinitionValidationResponse newObj = new (externalObj);
        return newObj;
    }
}

function org_wso2_apk_config_RuntimeAPICommonUtil_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["java.lang.Object"]
} external;

isolated function org_wso2_apk_config_RuntimeAPICommonUtil_generateDefinition(handle arg0) returns handle|error = @java:Method {
    name: "generateDefinition",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["org.wso2.apk.config.model.API"]
} external;

isolated function org_wso2_apk_config_RuntimeAPICommonUtil_generateDefinition2(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "generateDefinition",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["org.wso2.apk.config.model.API", "java.lang.String"]
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_generateUriTemplatesFromAPIDefinition(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "generateUriTemplatesFromAPIDefinition",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

isolated function org_wso2_apk_config_RuntimeAPICommonUtil_getAPIFromDefinition(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getAPIFromDefinition",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: []
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: []
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: []
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_RuntimeAPICommonUtil_validateOpenAPIDefinition(handle arg0, handle arg1, handle arg2, handle arg3, boolean arg4) returns handle|error = @java:Method {
    name: "validateOpenAPIDefinition",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["java.lang.String", "[B", "java.lang.String", "java.lang.String", "boolean"]
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: []
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["long"]
} external;

function org_wso2_apk_config_RuntimeAPICommonUtil_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.RuntimeAPICommonUtil",
    paramTypes: ["long", "int"]
} external;

