import ballerina/jballerina.java;

# Ballerina class mapping for the Java `org.wso2.apk.runtime.api.ErrorHandler` interface.
@java:Binding {'class: "org.wso2.apk.runtime.api.ErrorHandler"}
public distinct class ErrorHandler {

    *java:JObject;

    # The `handle` field that stores the reference to the `org.wso2.apk.runtime.api.ErrorHandler` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.runtime.api.ErrorHandler` Java interface.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public isolated function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.runtime.api.ErrorHandler` Java interface.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `getErrorCode` method of `org.wso2.apk.runtime.api.ErrorHandler`.
    #
    # + return - The `int` value returning from the Java mapping.
    public isolated function getErrorCode() returns int {
        return org_wso2_apk_runtime_api_ErrorHandler_getErrorCode(self.jObj);
    }

    # The function that maps to the `getErrorDescription` method of `org.wso2.apk.runtime.api.ErrorHandler`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getErrorDescription() returns string? {
        return java:toString(org_wso2_apk_runtime_api_ErrorHandler_getErrorDescription(self.jObj));
    }

    # The function that maps to the `getErrorMessage` method of `org.wso2.apk.runtime.api.ErrorHandler`.
    #
    # + return - The `string` value returning from the Java mapping.
    public isolated function getErrorMessage() returns string? {
        return java:toString(org_wso2_apk_runtime_api_ErrorHandler_getErrorMessage(self.jObj))
    ;
    }

    # The function that maps to the `getHttpStatusCode` method of `org.wso2.apk.runtime.api.ErrorHandler`.
    #
    # + return - The `int` value returning from the Java mapping.
    public isolated function getHttpStatusCode() returns int {
        return org_wso2_apk_runtime_api_ErrorHandler_getHttpStatusCode(self.jObj);
    }

    # The function that maps to the `printStackTrace` method of `org.wso2.apk.runtime.api.ErrorHandler`.
    #
    # + return - The `boolean` value returning from the Java mapping.
    public isolated function printStackTrace() returns boolean {
        return org_wso2_apk_runtime_api_ErrorHandler_printStackTrace(self.jObj);
    }

}

isolated function org_wso2_apk_runtime_api_ErrorHandler_getErrorCode(handle receiver) returns int = @java:Method {
    name: "getErrorCode",
    'class: "org.wso2.apk.runtime.api.ErrorHandler",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorHandler_getErrorDescription(handle receiver) returns handle = @java:Method {
    name: "getErrorDescription",
    'class: "org.wso2.apk.runtime.api.ErrorHandler",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorHandler_getErrorMessage(handle receiver) returns handle = @java:Method {
    name: "getErrorMessage",
    'class: "org.wso2.apk.runtime.api.ErrorHandler",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorHandler_getHttpStatusCode(handle receiver) returns int = @java:Method {
    name: "getHttpStatusCode",
    'class: "org.wso2.apk.runtime.api.ErrorHandler",
    paramTypes: []
} external;

isolated function org_wso2_apk_runtime_api_ErrorHandler_printStackTrace(handle receiver) returns boolean = @java:Method {
    name: "printStackTrace",
    'class: "org.wso2.apk.runtime.api.ErrorHandler",
    paramTypes: []
} external;

