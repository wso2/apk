import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import backoffice_domain_service.org.wso2.apk.apimgt.api as orgwso2apkapimgtapi;
import backoffice_domain_service.java.lang as javalang;
import backoffice_domain_service.java.io as javaio;
import backoffice_domain_service.org.wso2.apk.apimgt.api.model as orgwso2apkapimgtapimodel;

# Ballerina class mapping for the Java `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl` class.
@java:Binding {'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl"}
public distinct class ApisApiCommonImpl {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    public function notify() {
        org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    public function notifyAll() {
        org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The function that maps to the `getAPI` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_getAPI(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAPI(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getAPIResourcePaths` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `int` value required to map with the Java method parameter.
# + arg2 - The `int` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_getAPIResourcePaths(string arg0, int arg1, int arg2) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAPIResourcePaths(java:fromString(arg0), arg1, arg2);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getAPIThumbnail` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtapi:APIProvider` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_getAPIThumbnail(string arg0, orgwso2apkapimgtapi:APIProvider arg1, string arg2) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAPIThumbnail(java:fromString(arg0), arg1.jObj, java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getAllAPIs` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `int` value required to map with the Java method parameter.
# + arg1 - The `int` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + arg4 - The `string` value required to map with the Java method parameter.
# + arg5 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_getAllAPIs(int arg0, int arg1, string arg2, string arg3, string arg4, string arg5) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAllAPIs(arg0, arg1, java:fromString(arg2), java:fromString(arg3), java:fromString(arg4), java:fromString(arg5));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getApiDefinition` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `orgwso2apkapimgtapimodel:API` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_getApiDefinition(orgwso2apkapimgtapimodel:API arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getApiDefinition(arg0.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `updateAPI` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string[]` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_updateAPI(string arg0, string arg1, string[] arg2, string arg3) returns string?|orgwso2apkapimgtapi:APIManagementException|error {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_updateAPI(java:fromString(arg0), java:fromString(arg1), check jarrays:toHandle(arg2, "java.lang.String"), java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `updateAPIThumbnail` method of `org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `javaio:InputStream` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + arg4 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisApiCommonImpl_updateAPIThumbnail(string arg0, javaio:InputStream arg1, string arg2, string arg3, string arg4) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_updateAPIThumbnail(java:fromString(arg0), arg1.jObj, java:fromString(arg2), java:fromString(arg3), java:fromString(arg4));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that retrieves the value of the public field `MESSAGE`.
#
# + return - The `string` value of the field.
public function ApisApiCommonImpl_getMESSAGE() returns string? {
    return java:toString(org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getMESSAGE());
}

# The function that retrieves the value of the public field `ERROR_WHILE_UPDATING_API`.
#
# + return - The `string` value of the field.
public function ApisApiCommonImpl_getERROR_WHILE_UPDATING_API() returns string? {
    return java:toString(org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getERROR_WHILE_UPDATING_API());
}

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAPI(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getAPI",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAPIResourcePaths(handle arg0, int arg1, int arg2) returns handle|error = @java:Method {
    name: "getAPIResourcePaths",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["java.lang.String", "int", "int"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAPIThumbnail(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "getAPIThumbnail",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.api.APIProvider", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getAllAPIs(int arg0, int arg1, handle arg2, handle arg3, handle arg4, handle arg5) returns handle|error = @java:Method {
    name: "getAllAPIs",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["int", "int", "java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getApiDefinition(handle arg0) returns handle|error = @java:Method {
    name: "getApiDefinition",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["org.wso2.apk.apimgt.api.model.API"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_updateAPI(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "updateAPI",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "[Ljava.lang.String;", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_updateAPIThumbnail(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4) returns handle|error = @java:Method {
    name: "updateAPIThumbnail",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["java.lang.String", "java.io.InputStream", "java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["long"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl",
    paramTypes: ["long", "int"]
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getMESSAGE() returns handle = @java:FieldGet {
    name: "MESSAGE",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl"
} external;

function org_wso2_apk_apimgt_rest_api_backoffice_v1_common_impl_ApisApiCommonImpl_getERROR_WHILE_UPDATING_API() returns handle = @java:FieldGet {
    name: "ERROR_WHILE_UPDATING_API",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl"
} external;

function HILE_UPDATING_API() returns handle = @java:FieldGet {
    name: "ERROR_WHILE_UPDATING_API",
    'class: "org.wso2.apk.apimgt.rest.api.backoffice.v1.common.impl.ApisApiCommonImpl"
} external;

