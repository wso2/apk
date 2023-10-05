import devportal_service.java.lang as javalang;
import devportal_service.java.util as javautil;

import ballerina/jballerina.java;

# Ballerina class mapping for the Java `org.wso2.apk.devportal.sdk.APIClientGenerationManager` class.
@java:Binding {'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager"}
public distinct class APIClientGenerationManager {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.devportal.sdk.APIClientGenerationManager` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.devportal.sdk.APIClientGenerationManager` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.devportal.sdk.APIClientGenerationManager` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "";
    }
    # The function that maps to the `cleanTempDirectory` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    public function cleanTempDirectory(string arg0) {
        org_wso2_apk_devportal_sdk_APIClientGenerationManager_cleanTempDirectory(self.jObj, java:fromString(arg0));
    }

    # The function that maps to the `equals` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_devportal_sdk_APIClientGenerationManager_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `generateSDK` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + arg0 - The `string` value required to map with the Java method parameter.
    # + arg1 - The `string` value required to map with the Java method parameter.
    # + arg2 - The `string` value required to map with the Java method parameter.
    # + arg3 - The `string` value required to map with the Java method parameter.
    # + arg4 - The `string` value required to map with the Java method parameter.
    # + arg5 - The `string` value required to map with the Java method parameter.
    # + arg6 - The `string` value required to map with the Java method parameter.
    # + arg7 - The `string` value required to map with the Java method parameter.
    # + return - The `javautil:Map` or the `APIClientGenerationException` value returning from the Java mapping.
    public isolated function generateSDK(string arg0, string arg1, string arg2, string arg3, string arg4, string arg5, string arg6, string arg7) returns javautil:Map|APIClientGenerationException {
        handle|error externalObj = org_wso2_apk_devportal_sdk_APIClientGenerationManager_generateSDK(self.jObj, java:fromString(arg0), java:fromString(arg1), java:fromString(arg2), java:fromString(arg3), java:fromString(arg4), java:fromString(arg5), java:fromString(arg6), java:fromString(arg7));
        if (externalObj is error) {
            APIClientGenerationException e = error APIClientGenerationException(APICLIENTGENERATIONEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            javautil:Map newObj = new (externalObj);
            return newObj;
        }
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_devportal_sdk_APIClientGenerationManager_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `getSupportedSDKLanguages` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getSupportedSDKLanguages() returns string {
        return java:toString(org_wso2_apk_devportal_sdk_APIClientGenerationManager_getSupportedSDKLanguages(self.jObj)) ?: "";
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_devportal_sdk_APIClientGenerationManager_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    public function notify() {
        org_wso2_apk_devportal_sdk_APIClientGenerationManager_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    public function notifyAll() {
        org_wso2_apk_devportal_sdk_APIClientGenerationManager_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_devportal_sdk_APIClientGenerationManager_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_devportal_sdk_APIClientGenerationManager_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_devportal_sdk_APIClientGenerationManager_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The constructor function to generate an object of `org.wso2.apk.devportal.sdk.APIClientGenerationManager`.
#
# + return - The new `APIClientGenerationManager` class generated.
public function newAPIClientGenerationManager1() returns APIClientGenerationManager {
    handle externalObj = org_wso2_apk_devportal_sdk_APIClientGenerationManager_newAPIClientGenerationManager1();
    APIClientGenerationManager newObj = new (externalObj);
    return newObj;
}

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_cleanTempDirectory(handle receiver, handle arg0) = @java:Method {
    name: "cleanTempDirectory",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: ["java.lang.Object"]
} external;

isolated function org_wso2_apk_devportal_sdk_APIClientGenerationManager_generateSDK(handle receiver, handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6, handle arg7) returns handle|error = @java:Method {
    name: "generateSDK",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

isolated function org_wso2_apk_devportal_sdk_APIClientGenerationManager_getSupportedSDKLanguages(handle receiver) returns handle = @java:Method {
    name: "getSupportedSDKLanguages",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: ["long"]
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_devportal_sdk_APIClientGenerationManager_newAPIClientGenerationManager1() returns handle = @java:Constructor {
    'class: "org.wso2.apk.devportal.sdk.APIClientGenerationManager",
    paramTypes: []
} external;

