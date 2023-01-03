import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.org.wso2.apk.runtime.model as orgwso2apkruntimemodel;
import runtime_domain_service.java.util as javautil;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.api.APIDefinition` class.
@java:Binding {'class: "org.wso2.apk.runtime.api.APIDefinition"}
public distinct class APIDefinition {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.api.APIDefinition` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.api.APIDefinition` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.api.APIDefinition` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `canHandleDefinition` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function canHandleDefinition(string arg0) returns boolean {
        return org_wso2_apk_runtime_api_APIDefinition_canHandleDefinition(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `equals` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_api_APIDefinition_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `generateAPIDefinition` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `orgwso2apkruntimemodel:API` value required to map with the Java method parameter.
    # + return - The `string` or the `APIManagementException` value returning from the Java mapping.
    public function generateAPIDefinition(orgwso2apkruntimemodel:API arg0) returns string?|APIManagementException {
        handle|error externalObj = org_wso2_apk_runtime_api_APIDefinition_generateAPIDefinition(self.jObj, arg0.jObj);
        if (externalObj is error) {
            APIManagementException e = error APIManagementException(APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            return java:toString(externalObj);
        }
    }

    # The function that maps to the `generateAPIDefinition` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `orgwso2apkruntimemodel:API` value required to map with the Java method parameter.
    # + arg1 - The `string` value required to map with the Java method parameter.
    # + return - The `string` or the `APIManagementException` value returning from the Java mapping.
    public function generateAPIDefinition2(orgwso2apkruntimemodel:API arg0, string arg1) returns string?|APIManagementException {
        handle|error externalObj = org_wso2_apk_runtime_api_APIDefinition_generateAPIDefinition2(self.jObj, arg0.jObj, java:fromString(arg1));
        if (externalObj is error) {
            APIManagementException e = error APIManagementException(APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            return java:toString(externalObj);
        }
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_api_APIDefinition_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getPathParamNames` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + return - The `javautil:List` value returning from the Java mapping.
    public function getPathParamNames(string arg0) returns javautil:List {
        handle externalObj = org_wso2_apk_runtime_api_APIDefinition_getPathParamNames(self.jObj, java:fromString(arg0));
        javautil:List newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getResourceMap` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `orgwso2apkruntimemodel:API` value required to map with the Java method parameter.
    # + return - The `javautil:Map` value returning from the Java mapping.
    public function getResourceMap(orgwso2apkruntimemodel:API arg0) returns javautil:Map {
        handle externalObj = org_wso2_apk_runtime_api_APIDefinition_getResourceMap(self.jObj, arg0.jObj);
        javautil:Map newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getScopes` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + return - The `javautil:Set` or the `APIManagementException` value returning from the Java mapping.
    public function getScopes(string arg0) returns javautil:Set|APIManagementException {
        handle|error externalObj = org_wso2_apk_runtime_api_APIDefinition_getScopes(self.jObj, java:fromString(arg0));
        if (externalObj is error) {
            APIManagementException e = error APIManagementException(APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            javautil:Set newObj = new (externalObj);
            return newObj;
        }
    }

    # The function that maps to the `getType` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getType() returns string? {
        return java:toString(org_wso2_apk_runtime_api_APIDefinition_getType(self.jObj));
    }

    # The function that maps to the `getURITemplates` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + return - The `javautil:Set` or the `APIManagementException` value returning from the Java mapping.
    public function getURITemplates(string arg0) returns javautil:Set|APIManagementException {
        handle|error externalObj = org_wso2_apk_runtime_api_APIDefinition_getURITemplates(self.jObj, java:fromString(arg0));
        if (externalObj is error) {
            APIManagementException e = error APIManagementException(APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            javautil:Set newObj = new (externalObj);
            return newObj;
        }
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_api_APIDefinition_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.api.APIDefinition`.
    public function notify() {
        org_wso2_apk_runtime_api_APIDefinition_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.api.APIDefinition`.
    public function notifyAll() {
        org_wso2_apk_runtime_api_APIDefinition_notifyAll(self.jObj);
    }

    # The function that maps to the `processOtherSchemeScopes` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + return - The `string` or the `APIManagementException` value returning from the Java mapping.
    public function processOtherSchemeScopes(string arg0) returns string?|APIManagementException {
        handle|error externalObj = org_wso2_apk_runtime_api_APIDefinition_processOtherSchemeScopes(self.jObj, java:fromString(arg0));
        if (externalObj is error) {
            APIManagementException e = error APIManagementException(APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            return java:toString(externalObj);
        }
    }

    # The function that maps to the `validateAPIDefinition` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + arg1 - The `boolean` value required to map with the Java method parameter.
    # + return - The `APIDefinitionValidationResponse` or the `APIManagementException` value returning from the Java mapping.
    public function validateAPIDefinition(string arg0, boolean arg1) returns APIDefinitionValidationResponse|APIManagementException {
        handle|error externalObj = org_wso2_apk_runtime_api_APIDefinition_validateAPIDefinition(self.jObj, java:fromString(arg0), arg1);
        if (externalObj is error) {
            APIManagementException e = error APIManagementException(APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            APIDefinitionValidationResponse newObj = new (externalObj);
            return newObj;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_APIDefinition_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_APIDefinition_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.api.APIDefinition`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_api_APIDefinition_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.api.APIDefinition`.
#
# + return - The new `APIDefinition` class generated.
public function newAPIDefinition1() returns APIDefinition {
    handle externalObj = org_wso2_apk_runtime_api_APIDefinition_newAPIDefinition1();
    APIDefinition newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_runtime_api_APIDefinition_canHandleDefinition(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "canHandleDefinition",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_generateAPIDefinition(handle receiver, handle arg0) returns handle|error = @java:Method {
    name: "generateAPIDefinition",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["org.wso2.apk.runtime.model.API"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_generateAPIDefinition2(handle receiver, handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "generateAPIDefinition",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["org.wso2.apk.runtime.model.API", "java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIDefinition_getPathParamNames(handle receiver, handle arg0) returns handle = @java:Method {
    name: "getPathParamNames",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_getResourceMap(handle receiver, handle arg0) returns handle = @java:Method {
    name: "getResourceMap",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["org.wso2.apk.runtime.model.API"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_getScopes(handle receiver, handle arg0) returns handle|error = @java:Method {
    name: "getScopes",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_getType(handle receiver) returns handle = @java:Method {
    name: "getType",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIDefinition_getURITemplates(handle receiver, handle arg0) returns handle|error = @java:Method {
    name: "getURITemplates",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIDefinition_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIDefinition_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIDefinition_processOtherSchemeScopes(handle receiver, handle arg0) returns handle|error = @java:Method {
    name: "processOtherSchemeScopes",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_validateAPIDefinition(handle receiver, handle arg0, boolean arg1) returns handle|error = @java:Method {
    name: "validateAPIDefinition",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["java.lang.String", "boolean"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

function org_wso2_apk_runtime_api_APIDefinition_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_runtime_api_APIDefinition_newAPIDefinition1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.api.APIDefinition",
    paramTypes: []
} external;

