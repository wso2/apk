import ballerina/jballerina.java;
import runtime_domain_service.java.lang as javalang;
import runtime_domain_service.java.util as javautil;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.model.URITemplate` class.
@java:Binding {'class: "org.wso2.apk.runtime.model.URITemplate"}
public distinct class URITemplate {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.model.URITemplate` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.model.URITemplate` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.model.URITemplate` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public isolated function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `addAllScopes` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `javautil:List` value required to map with the Java method parameter.
    public function addAllScopes(javautil:List arg0) {
        org_wso2_apk_runtime_model_URITemplate_addAllScopes(self.jObj, arg0.jObj);
    }

    # The function that maps to the `addOperationPolicy` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `OperationPolicy` value required to map with the Java method parameter.
    public function addOperationPolicy(OperationPolicy arg0) {
        org_wso2_apk_runtime_model_URITemplate_addOperationPolicy(self.jObj, arg0.jObj);
    }

    # The function that maps to the `equals` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_runtime_model_URITemplate_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getAmznResourceName` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getAmznResourceName() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getAmznResourceName(self.jObj));
    }

    # The function that maps to the `getAmznResourceTimeout` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function getAmznResourceTimeout() returns int {
        return org_wso2_apk_runtime_model_URITemplate_getAmznResourceTimeout(self.jObj);
    }

    # The function that maps to the `getAuthType` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getAuthType() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getAuthType(self.jObj));
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_runtime_model_URITemplate_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getHTTPVerb` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getHTTPVerb() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getHTTPVerb(self.jObj));
    }

    # The function that maps to the `getId` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function getId() returns int {
        return org_wso2_apk_runtime_model_URITemplate_getId(self.jObj);
    }

    # The function that maps to the `getOperationPolicies` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `javautil:List` value returning from the Java mapping.
    public function getOperationPolicies() returns javautil:List {
        handle externalObj = org_wso2_apk_runtime_model_URITemplate_getOperationPolicies(self.jObj);
        javautil:List newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getResourceSandboxURI` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getResourceSandboxURI() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getResourceSandboxURI(self.jObj));
    }

    # The function that maps to the `getResourceURI` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getResourceURI() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getResourceURI(self.jObj));
    }

    # The function that maps to the `getScope` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `Scope` value returning from the Java mapping.
    public function getScope() returns Scope {
        handle externalObj = org_wso2_apk_runtime_model_URITemplate_getScope(self.jObj);
        Scope newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getScopes` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `javautil:List` value returning from the Java mapping.
    public function getScopes() returns javautil:List {
        handle externalObj = org_wso2_apk_runtime_model_URITemplate_getScopes(self.jObj);
        javautil:List newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getThrottlingTier` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getThrottlingTier() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getThrottlingTier(self.jObj));
    }

    # The function that maps to the `getUriTemplate` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `string` value returning from the Java mapping.
    public function getUriTemplate() returns string? {
        return java:toString(org_wso2_apk_runtime_model_URITemplate_getUriTemplate(self.jObj));
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_runtime_model_URITemplate_hashCode(self.jObj);
    }

    # The function that maps to the `isResourceSandboxURIExist` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isResourceSandboxURIExist() returns boolean {
        return org_wso2_apk_runtime_model_URITemplate_isResourceSandboxURIExist(self.jObj);
    }

    # The function that maps to the `isResourceURIExist` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public function isResourceURIExist() returns boolean {
        return org_wso2_apk_runtime_model_URITemplate_isResourceURIExist(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.runtime.model.URITemplate`.
    public function notify() {
        org_wso2_apk_runtime_model_URITemplate_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.runtime.model.URITemplate`.
    public function notifyAll() {
        org_wso2_apk_runtime_model_URITemplate_notifyAll(self.jObj);
    }

    # The function that maps to the `retrieveAllScopes` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `javautil:List` value returning from the Java mapping.
    public function retrieveAllScopes() returns javautil:List {
        handle externalObj = org_wso2_apk_runtime_model_URITemplate_retrieveAllScopes(self.jObj);
        javautil:List newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `setAmznResourceName` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setAmznResourceName(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setAmznResourceName(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setAmznResourceTimeout` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setAmznResourceTimeout(int arg0) {
        org_wso2_apk_runtime_model_URITemplate_setAmznResourceTimeout(self.jObj, arg0);
    }

    # The function that maps to the `setAuthType` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setAuthType(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setAuthType(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setHTTPVerb` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setHTTPVerb(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setHTTPVerb(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setId` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    public function setId(int arg0) {
        org_wso2_apk_runtime_model_URITemplate_setId(self.jObj, arg0);
    }

    # The function that maps to the `setOperationPolicies` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `javautil:List` value required to map with the Java method parameter.
    public isolated function setOperationPolicies(javautil:List arg0) {
        org_wso2_apk_runtime_model_URITemplate_setOperationPolicies(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setResourceSandboxURI` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setResourceSandboxURI(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setResourceSandboxURI(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setResourceURI` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function setResourceURI(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setResourceURI(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setScope` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `Scope` value required to map with the Java method parameter.
    public function setScope(Scope arg0) {
        org_wso2_apk_runtime_model_URITemplate_setScope(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setScopes` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `Scope` value required to map with the Java method parameter.
    public isolated function setScopes(Scope arg0) {
        org_wso2_apk_runtime_model_URITemplate_setScopes(self.jObj, arg0.jObj);
    }

    # The function that maps to the `setThrottlingTier` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setThrottlingTier(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setThrottlingTier(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `setUriTemplate` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public isolated function setUriTemplate(string arg0) {
        org_wso2_apk_runtime_model_URITemplate_setUriTemplate(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_URITemplate_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_URITemplate_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.runtime.model.URITemplate`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_runtime_model_URITemplate_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.runtime.model.URITemplate`.
#
# + return - The new `URITemplate` class generated.
public isolated function newURITemplate1() returns URITemplate {
    lock {
        handle externalObj = org_wso2_apk_runtime_model_URITemplate_newURITemplate1();
        URITemplate newObj = new (externalObj);
        return newObj;
    }
}

function org_wso2_apk_runtime_model_URITemplate_addAllScopes(handle receiver, handle arg0) = @java:Method {
    name: "addAllScopes",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.util.List"]
} external;

function org_wso2_apk_runtime_model_URITemplate_addOperationPolicy(handle receiver, handle arg0) = @java:Method {
    name: "addOperationPolicy",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["org.wso2.apk.runtime.model.OperationPolicy"]
} external;

function org_wso2_apk_runtime_model_URITemplate_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_runtime_model_URITemplate_getAmznResourceName(handle receiver) returns handle = @java:Method {
    name: "getAmznResourceName",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getAmznResourceTimeout(handle receiver) returns int = @java:Method {
    name: "getAmznResourceTimeout",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getAuthType(handle receiver) returns handle = @java:Method {
    name: "getAuthType",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getHTTPVerb(handle receiver) returns handle = @java:Method {
    name: "getHTTPVerb",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getId(handle receiver) returns int = @java:Method {
    name: "getId",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getOperationPolicies(handle receiver) returns handle = @java:Method {
    name: "getOperationPolicies",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getResourceSandboxURI(handle receiver) returns handle = @java:Method {
    name: "getResourceSandboxURI",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getResourceURI(handle receiver) returns handle = @java:Method {
    name: "getResourceURI",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getScope(handle receiver) returns handle = @java:Method {
    name: "getScope",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getScopes(handle receiver) returns handle = @java:Method {
    name: "getScopes",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getThrottlingTier(handle receiver) returns handle = @java:Method {
    name: "getThrottlingTier",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_getUriTemplate(handle receiver) returns handle = @java:Method {
    name: "getUriTemplate",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_isResourceSandboxURIExist(handle receiver) returns boolean = @java:Method {
    name: "isResourceSandboxURIExist",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_isResourceURIExist(handle receiver) returns boolean = @java:Method {
    name: "isResourceURIExist",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_retrieveAllScopes(handle receiver) returns handle = @java:Method {
    name: "retrieveAllScopes",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_setAmznResourceName(handle receiver, handle arg0) = @java:Method {
    name: "setAmznResourceName",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_URITemplate_setAmznResourceTimeout(handle receiver, int arg0) = @java:Method {
    name: "setAmznResourceTimeout",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["int"]
} external;

function org_wso2_apk_runtime_model_URITemplate_setAuthType(handle receiver, handle arg0) = @java:Method {
    name: "setAuthType",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_runtime_model_URITemplate_setHTTPVerb(handle receiver, handle arg0) = @java:Method {
    name: "setHTTPVerb",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_URITemplate_setId(handle receiver, int arg0) = @java:Method {
    name: "setId",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["int"]
} external;

isolated function org_wso2_apk_runtime_model_URITemplate_setOperationPolicies(handle receiver, handle arg0) = @java:Method {
    name: "setOperationPolicies",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.util.List"]
} external;

function org_wso2_apk_runtime_model_URITemplate_setResourceSandboxURI(handle receiver, handle arg0) = @java:Method {
    name: "setResourceSandboxURI",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_URITemplate_setResourceURI(handle receiver, handle arg0) = @java:Method {
    name: "setResourceURI",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_URITemplate_setScope(handle receiver, handle arg0) = @java:Method {
    name: "setScope",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["org.wso2.apk.runtime.model.Scope"]
} external;

isolated function org_wso2_apk_runtime_model_URITemplate_setScopes(handle receiver, handle arg0) = @java:Method {
    name: "setScopes",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["org.wso2.apk.runtime.model.Scope"]
} external;

isolated function org_wso2_apk_runtime_model_URITemplate_setThrottlingTier(handle receiver, handle arg0) = @java:Method {
    name: "setThrottlingTier",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

isolated function org_wso2_apk_runtime_model_URITemplate_setUriTemplate(handle receiver, handle arg0) = @java:Method {
    name: "setUriTemplate",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_runtime_model_URITemplate_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

function org_wso2_apk_runtime_model_URITemplate_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["long"]
} external;

function org_wso2_apk_runtime_model_URITemplate_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: ["long", "int"]
} external;

isolated function org_wso2_apk_runtime_model_URITemplate_newURITemplate1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.runtime.model.URITemplate",
    paramTypes: []
} external;

