import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import config_deployer_service.java.lang as javalang;
import config_deployer_service.java.util as javautil;

# Ballerina class mapping for the Java `org.wso2.apk.config.model.URITemplate` class.
@java:Binding {'class: "org.wso2.apk.config.model.URITemplate"}
public distinct class URITemplate {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.config.model.URITemplate` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.config.model.URITemplate` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.config.model.URITemplate` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }
    # The function that maps to the `addAllScopes` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `javautil:List` value required to map with the Java method parameter.
    public function addAllScopes(javautil:List arg0) {
        org_wso2_apk_config_model_URITemplate_addAllScopes(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_config_model_URITemplate_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_config_model_URITemplate_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getEndpoint` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getEndpoint() returns string {
        return java:toString(org_wso2_apk_config_model_URITemplate_getEndpoint(self.jObj)) ?: "";
    }

    # The function that maps to the `getHTTPVerb` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getHTTPVerb() returns string {
        return java:toString(org_wso2_apk_config_model_URITemplate_getHTTPVerb(self.jObj)) ?: "";
    }

    # The function that maps to the `getId` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function getId() returns int {
        return org_wso2_apk_config_model_URITemplate_getId(self.jObj);
    }

    # The function that maps to the `getResourceURI` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getResourceURI() returns string {
        return java:toString(org_wso2_apk_config_model_URITemplate_getResourceURI(self.jObj)) ?: "";
    }

    # The function that maps to the `getScopes` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `string[]` value returning from the Java mapping.
    public isolated function getScopes() returns string[]|error {
        handle externalObj = org_wso2_apk_config_model_URITemplate_getScopes(self.jObj);
        if java:isNull(externalObj) {
            return [];
        }
        return <string[]>check jarrays:fromHandle(externalObj, "string");
    }

    # The function that maps to the `getUriTemplate` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getUriTemplate() returns string {
        return java:toString(org_wso2_apk_config_model_URITemplate_getUriTemplate(self.jObj)) ?: "";
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_config_model_URITemplate_hashCode(self.jObj);
    }

    # The function that maps to the `isAuthEnabled` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public isolated function isAuthEnabled() returns boolean {
        return org_wso2_apk_config_model_URITemplate_isAuthEnabled(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.config.model.URITemplate`.
    public function notify() {
        org_wso2_apk_config_model_URITemplate_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.config.model.URITemplate`.
    public function notifyAll() {
        org_wso2_apk_config_model_URITemplate_notifyAll(self.jObj);
    }

    # The function that maps to the `retrieveAllScopes` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `string[]` value returning from the Java mapping.
    public isolated function retrieveAllScopes() returns string[]|error {
        handle externalObj = org_wso2_apk_config_model_URITemplate_retrieveAllScopes(self.jObj);
        if java:isNull(externalObj) {
            return [];
        }
        return <string[]>check jarrays:fromHandle(externalObj, "string");
    }

    # The function that maps to the `setAuthEnabled` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `boolean` value required to map with the Java method parameter.
    public isolated function setAuthEnabled(boolean arg0) {
        org_wso2_apk_config_model_URITemplate_setAuthEnabled(self.jObj, arg0);
    }

    # The function that maps to the `setEndpoint` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setEndpoint(string arg0) {
        org_wso2_apk_config_model_URITemplate_setEndpoint(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setHTTPVerb` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setHTTPVerb(string arg0) {
        org_wso2_apk_config_model_URITemplate_setHTTPVerb(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setId` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setId(int arg0) {
        org_wso2_apk_config_model_URITemplate_setId(self.jObj, arg0);
    }

    # The function that maps to the `setResourceURI` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setResourceURI(string arg0) {
        org_wso2_apk_config_model_URITemplate_setResourceURI(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setScopes` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setScopes(string arg0) {
        org_wso2_apk_config_model_URITemplate_setScopes(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setUriTemplate` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setUriTemplate(string arg0) {
        org_wso2_apk_config_model_URITemplate_setUriTemplate(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_model_URITemplate_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_model_URITemplate_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_model_URITemplate_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.config.model.URITemplate`.
#
# + return - The new `URITemplate` class generated.
public isolated function newURITemplate1() returns URITemplate {
    handle externalObj = org_wso2_apk_config_model_URITemplate_newURITemplate1();
    URITemplate newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_config_model_URITemplate_addAllScopes(handle receiver, handle arg0) = @java:Method {
    name: "addAllScopes",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.util.List"]
} external;

function org_wso2_apk_config_model_URITemplate_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_config_model_URITemplate_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_getEndpoint(handle receiver) returns handle = @java:Method {
    name: "getEndpoint",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_getHTTPVerb(handle receiver) returns handle = @java:Method {
    name: "getHTTPVerb",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_config_model_URITemplate_getId(handle receiver) returns int = @java:Method {
    name: "getId",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_config_model_URITemplate_getResourceURI(handle receiver) returns handle = @java:Method {
    name: "getResourceURI",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_getScopes(handle receiver) returns handle = @java:Method {
    name: "getScopes",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_getUriTemplate(handle receiver) returns handle = @java:Method {
    name: "getUriTemplate",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_config_model_URITemplate_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_isAuthEnabled(handle receiver) returns boolean = @java:Method {
    name: "isAuthEnabled",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_config_model_URITemplate_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_config_model_URITemplate_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_retrieveAllScopes(handle receiver) returns handle = @java:Method {
    name: "retrieveAllScopes",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_URITemplate_setAuthEnabled(handle receiver, boolean arg0) = @java:Method {
    name: "setAuthEnabled",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["boolean"]
} external;

isolated function org_wso2_apk_config_model_URITemplate_setEndpoint(handle receiver, handle arg0) = @java:Method {
    name: "setEndpoint",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_URITemplate_setHTTPVerb(handle receiver, handle arg0) = @java:Method {
    name: "setHTTPVerb",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_model_URITemplate_setId(handle receiver, int arg0) = @java:Method {
    name: "setId",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["int"]
} external;

function org_wso2_apk_config_model_URITemplate_setResourceURI(handle receiver, handle arg0) = @java:Method {
    name: "setResourceURI",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_URITemplate_setScopes(handle receiver, handle arg0) = @java:Method {
    name: "setScopes",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_URITemplate_setUriTemplate(handle receiver, handle arg0) = @java:Method {
    name: "setUriTemplate",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_model_URITemplate_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_config_model_URITemplate_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["long"]
} external;

function org_wso2_apk_config_model_URITemplate_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: ["long", "int"]
} external;

isolated function org_wso2_apk_config_model_URITemplate_newURITemplate1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.config.model.URITemplate",
    paramTypes: []
} external;

