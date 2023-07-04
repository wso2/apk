import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import config_deployer_service.java.lang as javalang;

# Ballerina class mapping for the Java `org.wso2.apk.config.api.APKConfValidationResponse` class.
@java:Binding {'class: "org.wso2.apk.config.api.APKConfValidationResponse"}
public distinct class APKConfValidationResponse {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.config.api.APKConfValidationResponse` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.config.api.APKConfValidationResponse` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.config.api.APKConfValidationResponse` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_config_api_APKConfValidationResponse_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_config_api_APKConfValidationResponse_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getErrorItems` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + return - The `ErrorHandler[]` value returning from the Java mapping.
    public isolated function getErrorItems() returns ErrorHandler[]|error {
        handle externalObj = org_wso2_apk_config_api_APKConfValidationResponse_getErrorItems(self.jObj);
        ErrorHandler[] newObj = [];
        handle[] anyObj = <handle[]>check jarrays:fromHandle(externalObj, "handle");
        int count = anyObj.length();
        foreach int i in 0 ... count - 1 {
            ErrorHandler element = new (anyObj[i]);
            newObj[i] = element;
        }
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_config_api_APKConfValidationResponse_hashCode(self.jObj);
    }

    # The function that maps to the `isValidated` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public isolated function isValidated() returns boolean {
        return org_wso2_apk_config_api_APKConfValidationResponse_isValidated(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    public function notify() {
        org_wso2_apk_config_api_APKConfValidationResponse_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    public function notifyAll() {
        org_wso2_apk_config_api_APKConfValidationResponse_notifyAll(self.jObj);
    }

    # The function that maps to the `setErrorItems` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + arg0 - The `ErrorHandler[]` value required to map with the Java method parameter.
    # + return - The `error?` value returning from the Java mapping.
    public function setErrorItems(ErrorHandler[] arg0) returns error? {
        org_wso2_apk_config_api_APKConfValidationResponse_setErrorItems(self.jObj, check jarrays:toHandle(arg0, "org.wso2.apk.config.api.ErrorHandler"));
    }

    # The function that maps to the `setValidated` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + arg0 - The `boolean` value required to map with the Java method parameter.
    public function setValidated(boolean arg0) {
        org_wso2_apk_config_api_APKConfValidationResponse_setValidated(self.jObj, arg0);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_api_APKConfValidationResponse_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_api_APKConfValidationResponse_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.config.api.APKConfValidationResponse`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_config_api_APKConfValidationResponse_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.config.api.APKConfValidationResponse`.
#
# + arg0 - The `boolean` value required to map with the Java constructor parameter.
# + return - The new `APKConfValidationResponse` class generated.
public function newAPKConfValidationResponse1(boolean arg0) returns APKConfValidationResponse {
    handle externalObj = org_wso2_apk_config_api_APKConfValidationResponse_newAPKConfValidationResponse1(arg0);
    APKConfValidationResponse newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_config_api_APKConfValidationResponse_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APKConfValidationResponse_getErrorItems(handle receiver) returns handle = @java:Method {
    name: "getErrorItems",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

isolated function org_wso2_apk_config_api_APKConfValidationResponse_isValidated(handle receiver) returns boolean = @java:Method {
    name: "isValidated",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_setErrorItems(handle receiver, handle arg0) = @java:Method {
    name: "setErrorItems",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: ["[Lorg.wso2.apk.config.api.ErrorHandler;"]
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_setValidated(handle receiver, boolean arg0) = @java:Method {
    name: "setValidated",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: ["boolean"]
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: []
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: ["long"]
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_config_api_APKConfValidationResponse_newAPKConfValidationResponse1(boolean arg0) returns handle = @java:Constructor {
    'class: "org.wso2.apk.config.api.APKConfValidationResponse",
    paramTypes: ["boolean"]
} external;

