import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import config_deployer_service.java.lang as javalang;

# Ballerina class mapping for the Java `org.wso2.apk.config.model.API` class.
@java:Binding {'class: "org.wso2.apk.config.model.API"}
public distinct class API {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.config.model.API` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.config.model.API` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.config.model.API` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_config_model_API_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getApiSecurity` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getApiSecurity() returns string {
        return java:toString(org_wso2_apk_config_model_API_getApiSecurity(self.jObj)) ?: "";
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_config_model_API_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getBasePath` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getBasePath() returns string {
        return java:toString(org_wso2_apk_config_model_API_getBasePath(self.jObj)) ?: "";
    }

    # The function that maps to the `getEndpoint` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getEndpoint() returns string {
        return java:toString(org_wso2_apk_config_model_API_getEndpoint(self.jObj)) ?: "";
    }

    # The function that maps to the `getEnvironment` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getEnvironment() returns string {
        return java:toString(org_wso2_apk_config_model_API_getEnvironment(self.jObj)) ?: "";
    }

    # The function that maps to the `getGraphQLSchema` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getGraphQLSchema() returns string {
        return java:toString(org_wso2_apk_config_model_API_getGraphQLSchema(self.jObj)) ?: "";
    }

    # The function that maps to the `getName` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getName() returns string {
        return java:toString(org_wso2_apk_config_model_API_getName(self.jObj)) ?: "";
    }

    # The function that maps to the `getScopes` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string[]` value returning from the Java mapping.
    public isolated function getScopes() returns string[]|error {
        handle externalObj = org_wso2_apk_config_model_API_getScopes(self.jObj);
        if java:isNull(externalObj) {
            return [];
        }
        return <string[]>check jarrays:fromHandle(externalObj, "string");
    }

    # The function that maps to the `getSwaggerDefinition` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getSwaggerDefinition() returns string {
        return java:toString(org_wso2_apk_config_model_API_getSwaggerDefinition(self.jObj)) ?: "";
    }

    # The function that maps to the `getType` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getType() returns string {
        return java:toString(org_wso2_apk_config_model_API_getType(self.jObj)) ?: "";
    }

    # The function that maps to the `getUriTemplates` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `URITemplate[]` value returning from the Java mapping.
    public isolated function getUriTemplates() returns URITemplate[]|error {
        handle externalObj = org_wso2_apk_config_model_API_getUriTemplates(self.jObj);
        URITemplate[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            URITemplate element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `getVersion` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getVersion() returns string {
        return java:toString(org_wso2_apk_config_model_API_getVersion(self.jObj)) ?: "";
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_config_model_API_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.config.model.API`.
    public function notify() {
        org_wso2_apk_config_model_API_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.config.model.API`.
    public function notifyAll() {
        org_wso2_apk_config_model_API_notifyAll(self.jObj);
    }

    # The function that maps to the `setApiSecurity` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setApiSecurity(string arg0) {
        org_wso2_apk_config_model_API_setApiSecurity(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setBasePath` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setBasePath(string arg0) {
        org_wso2_apk_config_model_API_setBasePath(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setEndpoint` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setEndpoint(string arg0) {
        org_wso2_apk_config_model_API_setEndpoint(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setEnvironment` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setEnvironment(string arg0) {
        org_wso2_apk_config_model_API_setEnvironment(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setGraphQLSchema` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setGraphQLSchema(string arg0) {
        org_wso2_apk_config_model_API_setGraphQLSchema(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setProtoDefinition` method of `org.wso2.apk.config.model.API`.
    #
    // # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setProtoDefinition(string arg0) {
        org_wso2_apk_config_model_API_setProtoDefinition(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setName` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setName(string arg0) {
        org_wso2_apk_config_model_API_setName(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setScopes` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string[]` value required to map with the Java method parameter.
    # + return - The `error?` value returning from the Java mapping.
    public isolated function setScopes(string[] arg0) returns error? {
        org_wso2_apk_config_model_API_setScopes(self.jObj, check jarrays:toHandle(arg0, "java.lang.String"));
    }

    # The function that maps to the `setSwaggerDefinition` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setSwaggerDefinition(string arg0) {
        org_wso2_apk_config_model_API_setSwaggerDefinition(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setType` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setType(string arg0) {
        org_wso2_apk_config_model_API_setType(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setUriTemplates` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `URITemplate[]` value required to map with the Java method parameter.
    # + return - The `error?` value returning from the Java mapping.
    public isolated function setUriTemplates(URITemplate[] arg0) returns error? {
        org_wso2_apk_config_model_API_setUriTemplates(self.jObj, check jarrays:toHandle(arg0, "org.wso2.apk.config.model.URITemplate"));
    }

    # The function that maps to the `setVersion` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setVersion(string arg0) {
        org_wso2_apk_config_model_API_setVersion(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.model.API`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_model_API_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_model_API_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.model.API`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_model_API_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.config.model.API`.
#
# + return - The new `API` class generated.
public isolated function newAPI1() returns API {
    handle externalObj = org_wso2_apk_config_model_API_newAPI1();
    API newObj = new (externalObj);
    return newObj;
}

# The constructor function to generate an object of `org.wso2.apk.config.model.API`.
#
# + arg0 - The `string` value required to map with the Java constructor parameter.
# + arg1 - The `string` value required to map with the Java constructor parameter.
# + arg2 - The `URITemplate[]` value required to map with the Java constructor parameter.
# + return - The new `API` class generated.
public function newAPI2(string arg0, string arg1, URITemplate[] arg2) returns API|error {
    handle externalObj = org_wso2_apk_config_model_API_newAPI2(java:fromString(arg0), java:fromString(arg1), check jarrays:toHandle(arg2, "org.wso2.apk.config.model.URITemplate"));
    API newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_config_model_API_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_config_model_API_getApiSecurity(handle receiver) returns handle = @java:Method {
    name: "getApiSecurity",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getBasePath(handle receiver) returns handle = @java:Method {
    name: "getBasePath",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getEndpoint(handle receiver) returns handle = @java:Method {
    name: "getEndpoint",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getEnvironment(handle receiver) returns handle = @java:Method {
    name: "getEnvironment",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getGraphQLSchema(handle receiver) returns handle = @java:Method {
    name: "getGraphQLSchema",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;


isolated function org_wso2_apk_config_model_API_getProtoDefinition(handle receiver) returns handle = @java:Method {
    name: "getProtoDefinition",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getName(handle receiver) returns handle = @java:Method {
    name: "getName",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getScopes(handle receiver) returns handle = @java:Method {
    name: "getScopes",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_getSwaggerDefinition(handle receiver) returns handle = @java:Method {
    name: "getSwaggerDefinition",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getType(handle receiver) returns handle = @java:Method {
    name: "getType",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getUriTemplates(handle receiver) returns handle = @java:Method {
    name: "getUriTemplates",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_model_API_getVersion(handle receiver) returns handle = @java:Method {
    name: "getVersion",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_setApiSecurity(handle receiver, handle arg0) = @java:Method {
    name: "setApiSecurity",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_model_API_setBasePath(handle receiver, handle arg0) = @java:Method {
    name: "setBasePath",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_model_API_setEndpoint(handle receiver, handle arg0) = @java:Method {
    name: "setEndpoint",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_model_API_setEnvironment(handle receiver, handle arg0) = @java:Method {
    name: "setEnvironment",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_API_setGraphQLSchema(handle receiver, handle arg0) = @java:Method {
    name: "setGraphQLSchema",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_API_setProtoDefinition(handle receiver, handle arg0) = @java:Method {
    name: "setProtoDefinition",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_API_setName(handle receiver, handle arg0) = @java:Method {
    name: "setName",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_API_setScopes(handle receiver, handle arg0) = @java:Method {
    name: "setScopes",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["[Ljava.lang.String;"]
} external;

function org_wso2_apk_config_model_API_setSwaggerDefinition(handle receiver, handle arg0) = @java:Method {
    name: "setSwaggerDefinition",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_API_setType(handle receiver, handle arg0) = @java:Method {
    name: "setType",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_config_model_API_setUriTemplates(handle receiver, handle arg0) = @java:Method {
    name: "setUriTemplates",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["[Lorg.wso2.apk.config.model.URITemplate;"]
} external;

isolated function org_wso2_apk_config_model_API_setVersion(handle receiver, handle arg0) = @java:Method {
    name: "setVersion",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_config_model_API_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["long"]
} external;

function org_wso2_apk_config_model_API_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["long", "int"]
} external;

isolated function org_wso2_apk_config_model_API_newAPI1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.config.model.API",
    paramTypes: []
} external;

function org_wso2_apk_config_model_API_newAPI2(handle arg0, handle arg1, handle arg2) returns handle = @java:Constructor {
    'class: "org.wso2.apk.config.model.API",
    paramTypes: ["java.lang.String", "java.lang.String", "[Lorg.wso2.apk.config.model.URITemplate;"]
} external;

