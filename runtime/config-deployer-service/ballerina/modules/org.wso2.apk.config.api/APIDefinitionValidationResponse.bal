import config_deployer_service.java.lang as javalang;
import config_deployer_service.java.util as javautil;

import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;

# Ballerina class mapping for the Java `org.wso2.apk.config.api.APIDefinitionValidationResponse` class.
@java:Binding {'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse"}
public distinct class APIDefinitionValidationResponse {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.config.api.APIDefinitionValidationResponse` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.config.api.APIDefinitionValidationResponse` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.config.api.APIDefinitionValidationResponse` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_config_api_APIDefinitionValidationResponse_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getContent` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getContent() returns string {
        return java:toString(org_wso2_apk_config_api_APIDefinitionValidationResponse_getContent(self.jObj)) ?: "";
    }

    # The function that maps to the `getErrorItems` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `javautil:ArrayList` value returning from the Java mapping.
    public isolated function getErrorItems() returns javautil:ArrayList {
        handle externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_getErrorItems(self.jObj);
        javautil:ArrayList newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getInfo` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `Info` value returning from the Java mapping.
    public function getInfo() returns Info {
        handle externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_getInfo(self.jObj);
        Info newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getJsonContent` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getJsonContent() returns string {
        return java:toString(org_wso2_apk_config_api_APIDefinitionValidationResponse_getJsonContent(self.jObj)) ?: "";
    }

    # The function that maps to the `getParser` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `APIDefinition` value returning from the Java mapping.
    public isolated function getParser() returns APIDefinition {
        handle externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_getParser(self.jObj);
        APIDefinition newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getProtoContent` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `byte[]` value returning from the Java mapping.
    public isolated function getProtoContent() returns byte[]|error {
        handle externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_getProtoContent(self.jObj);
        return <byte[]>check jarrays:fromHandle(externalObj, "byte");
    }

    # The function that maps to the `getProtocol` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getProtocol() returns string {
        return java:toString(org_wso2_apk_config_api_APIDefinitionValidationResponse_getProtocol(self.jObj)) ?: "";
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_config_api_APIDefinitionValidationResponse_hashCode(self.jObj);
    }

    # The function that maps to the `isInit` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isInit() returns boolean {
        return org_wso2_apk_config_api_APIDefinitionValidationResponse_isInit(self.jObj);
    }

    # The function that maps to the `isValid` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public isolated function isValid() returns boolean {
        return org_wso2_apk_config_api_APIDefinitionValidationResponse_isValid(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    public function notify() {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    public function notifyAll() {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_notifyAll(self.jObj);
    }

    # The function that maps to the `setContent` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setContent(string arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setContent(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setErrorItems` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `javautil:ArrayList` value required to map with the Java method parameter.
    public isolated function setErrorItems(javautil:ArrayList arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setErrorItems(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setInfo` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `Info` value required to map with the Java method parameter.
    public function setInfo(Info arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setInfo(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setInit` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `boolean` value required to map with the Java method parameter.
    public function setInit(boolean arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setInit(self.jObj, arg0);
    }

    # The function that maps to the `setJsonContent` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setJsonContent(string arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setJsonContent(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setParser` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `APIDefinition` value required to map with the Java method parameter.
    public function setParser(APIDefinition arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setParser(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setProtoContent` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `byte[]` value required to map with the Java method parameter.
    # + return - The `error?` value returning from the Java mapping.
    public function setProtoContent(byte[] arg0) returns error? {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setProtoContent(self.jObj, check jarrays:toHandle(arg0, "byte"));
    }

    # The function that maps to the `setProtocol` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setProtocol(string arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setProtocol(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setValid` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `boolean` value required to map with the Java method parameter.
    public function setValid(boolean arg0) {
        org_wso2_apk_config_api_APIDefinitionValidationResponse_setValid(self.jObj, arg0);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.config.api.APIDefinitionValidationResponse`.
#
# + return - The new `APIDefinitionValidationResponse` class generated.
public function newAPIDefinitionValidationResponse1() returns APIDefinitionValidationResponse {
    handle externalObj = org_wso2_apk_config_api_APIDefinitionValidationResponse_newAPIDefinitionValidationResponse1();
    APIDefinitionValidationResponse newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_config_api_APIDefinitionValidationResponse_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_getContent(handle receiver) returns handle = @java:Method {
    name: "getContent",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_getErrorItems(handle receiver) returns handle = @java:Method {
    name: "getErrorItems",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_getInfo(handle receiver) returns handle = @java:Method {
    name: "getInfo",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_getJsonContent(handle receiver) returns handle = @java:Method {
    name: "getJsonContent",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_getParser(handle receiver) returns handle = @java:Method {
    name: "getParser",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_getProtoContent(handle receiver) returns handle = @java:Method {
    name: "getProtoContent",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_getProtocol(handle receiver) returns handle = @java:Method {
    name: "getProtocol",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_isInit(handle receiver) returns boolean = @java:Method {
    name: "isInit",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_isValid(handle receiver) returns boolean = @java:Method {
    name: "isValid",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_setContent(handle receiver, handle arg0) = @java:Method {
    name: "setContent",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_api_APIDefinitionValidationResponse_setErrorItems(handle receiver, handle arg0) = @java:Method {
    name: "setErrorItems",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["java.util.ArrayList"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setInfo(handle receiver, handle arg0) = @java:Method {
    name: "setInfo",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["org.wso2.apk.config.api.Info"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setInit(handle receiver, boolean arg0) = @java:Method {
    name: "setInit",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["boolean"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setJsonContent(handle receiver, handle arg0) = @java:Method {
    name: "setJsonContent",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setParser(handle receiver, handle arg0) = @java:Method {
    name: "setParser",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["org.wso2.apk.config.api.APIDefinition"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setProtoContent(handle receiver, handle arg0) = @java:Method {
    name: "setProtoContent",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["[B"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setProtocol(handle receiver, handle arg0) = @java:Method {
    name: "setProtocol",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_setValid(handle receiver, boolean arg0) = @java:Method {
    name: "setValid",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["boolean"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["long"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_config_api_APIDefinitionValidationResponse_newAPIDefinitionValidationResponse1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.config.api.APIDefinitionValidationResponse",
    paramTypes: []
} external;

