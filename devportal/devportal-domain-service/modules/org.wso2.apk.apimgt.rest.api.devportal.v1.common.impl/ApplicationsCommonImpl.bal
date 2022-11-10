import ballerina/jballerina.java;
import devportal_service.org.wso2.apk.apimgt.api as orgwso2apkapimgtapi;
import devportal_service.java.lang as javalang;
import devportal_service.java.util as javautil;
import devportal_service.java.io as javaio;
import devportal_service.org.wso2.apk.apimgt.api.model as orgwso2apkapimgtapimodel;
import devportal_service.org.wso2.apk.apimgt.rest.api.devportal.v1.common.models as orgwso2apkapimgtrestapidevportalv1commonmodels;
import devportal_service.org.wso2.apk.apimgt.rest.api.devportal.v1.dto as orgwso2apkapimgtrestapidevportalv1dto;

# Ballerina class mapping for the Java `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl` class.
@java:Binding {'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl"}
public distinct class ApplicationsCommonImpl {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_hashCode(self.jObj);
    }

    # The function that maps to the `importApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + arg0 - The `javaio:InputStream` value required to map with the Java method parameter.
    # + arg1 - The `javalang:Boolean` value required to map with the Java method parameter.
    # + arg2 - The `javalang:Boolean` value required to map with the Java method parameter.
    # + arg3 - The `string` value required to map with the Java method parameter.
    # + arg4 - The `javalang:Boolean` value required to map with the Java method parameter.
    # + arg5 - The `javalang:Boolean` value required to map with the Java method parameter.
    # + arg6 - The `string` value required to map with the Java method parameter.
    # + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
    public function importApplication(javaio:InputStream arg0, javalang:Boolean arg1, javalang:Boolean arg2, string arg3, javalang:Boolean arg4, javalang:Boolean arg5, string arg6) returns string?|orgwso2apkapimgtapi:APIManagementException {
        handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_importApplication(self.jObj, arg0.jObj, arg1.jObj, arg2.jObj, java:fromString(arg3), arg4.jObj, arg5.jObj, java:fromString(arg6));
        if (externalObj is error) {
            orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
            return e;
        } else {
            return java:toString(externalObj);
        }
    }

    # The function that maps to the `notify` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    public function notify() {
        org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    public function notifyAll() {
        org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The function that maps to the `addApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_addApplication(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_addApplication(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `cleanUpApplicationRegistrationByApplicationIdAndKeyMappingId` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_cleanUpApplicationRegistrationByApplicationIdAndKeyMappingId(string arg0, string arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_cleanUpApplicationRegistrationByApplicationIdAndKeyMappingId(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `cleanupApplicationRegistration` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_cleanupApplicationRegistration(string arg0, string arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_cleanupApplicationRegistration(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `deleteApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `int` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_deleteApplication(string arg0) returns int|orgwso2apkapimgtapi:APIManagementException {
    int|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_deleteApplication(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return externalObj;
    }
}

# The function that maps to the `exportApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `javalang:Boolean` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `javaio:File` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_exportApplication(string arg0, string arg1, javalang:Boolean arg2, string arg3) returns javaio:File|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_exportApplication(java:fromString(arg0), java:fromString(arg1), arg2.jObj, java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        javaio:File newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `generateAPIKey` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `orgwso2apkapimgtrestapidevportalv1dto:APIKeyGenerateRequestDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:APIKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_generateAPIKey(string arg0, string arg1, orgwso2apkapimgtrestapidevportalv1dto:APIKeyGenerateRequestDTO arg2) returns orgwso2apkapimgtrestapidevportalv1dto:APIKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateAPIKey(java:fromString(arg0), java:fromString(arg1), arg2.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:APIKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `generateApplicationToken` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenGenerateRequestDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_generateApplicationToken(string arg0, string arg1, orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenGenerateRequestDTO arg2) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateApplicationToken(java:fromString(arg0), java:fromString(arg1), arg2.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `generateKeys` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyGenerateRequestDTO` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_generateKeys(string arg0, orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyGenerateRequestDTO arg1, string arg2) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateKeys(java:fromString(arg0), arg1.jObj, java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `generateTokenByOauthKeysKeyMappingId` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenGenerateRequestDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_generateTokenByOauthKeysKeyMappingId(string arg0, string arg1, orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenGenerateRequestDTO arg2) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateTokenByOauthKeysKeyMappingId(java:fromString(arg0), java:fromString(arg1), arg2.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationTokenDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getApplicationById` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getApplicationById(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationById(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getApplicationKeyByAppIDAndKeyMapping` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` value returning from the Java mapping.
public function ApplicationsCommonImpl_getApplicationKeyByAppIDAndKeyMapping(string arg0, string arg1) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO {
    handle externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationKeyByAppIDAndKeyMapping(java:fromString(arg0), java:fromString(arg1));
    orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `getApplicationKeyByAppIDAndKeyType` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getApplicationKeyByAppIDAndKeyType(string arg0, string arg1) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationKeyByAppIDAndKeyType(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getApplicationKeysByApplicationId` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getApplicationKeysByApplicationId(string arg0) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationKeysByApplicationId(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getApplicationList` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + arg4 - The `int` value required to map with the Java method parameter.
# + arg5 - The `int` value required to map with the Java method parameter.
# + arg6 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getApplicationList(string arg0, string arg1, string arg2, string arg3, int arg4, int arg5, string arg6) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationList(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2), java:fromString(arg3), arg4, arg5, java:fromString(arg6));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getApplicationOauthKeys` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getApplicationOauthKeys(string arg0, string arg1) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationOauthKeys(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getExportedApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `javaio:InputStream` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1commonmodels:ExportedApplication` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getExportedApplication(javaio:InputStream arg0) returns orgwso2apkapimgtrestapidevportalv1commonmodels:ExportedApplication|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getExportedApplication(arg0.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1commonmodels:ExportedApplication newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getOwnerId` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `javautil:List` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `javalang:Boolean` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getOwnerId(javautil:List arg0, string arg1, javalang:Boolean arg2, string arg3) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getOwnerId(arg0.jObj, java:fromString(arg1), arg2.jObj, java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getSkippedAPIs` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `javautil:Set` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `javalang:Boolean` value required to map with the Java method parameter.
# + arg3 - The `javalang:Boolean` value required to map with the Java method parameter.
# + arg4 - The `orgwso2apkapimgtapimodel:Application` value required to map with the Java method parameter.
# + arg5 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:APIInfoListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_getSkippedAPIs(javautil:Set arg0, string arg1, javalang:Boolean arg2, javalang:Boolean arg3, orgwso2apkapimgtapimodel:Application arg4, string arg5) returns orgwso2apkapimgtrestapidevportalv1dto:APIInfoListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getSkippedAPIs(arg0.jObj, java:fromString(arg1), arg2.jObj, arg3.jObj, arg4.jObj, java:fromString(arg5));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:APIInfoListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `mapApplicationKeys` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyMappingRequestDTO` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_mapApplicationKeys(string arg0, orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyMappingRequestDTO arg1, string arg2) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_mapApplicationKeys(java:fromString(arg0), arg1.jObj, java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `preProcessApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationDTO` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `boolean` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapimodel:Application` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_preProcessApplication(string arg0, orgwso2apkapimgtrestapidevportalv1dto:ApplicationDTO arg1, string arg2, boolean arg3) returns orgwso2apkapimgtapimodel:Application|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_preProcessApplication(java:fromString(arg0), arg1.jObj, java:fromString(arg2), arg3);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtapimodel:Application newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `regenerateSecretApplicationOauthKeysKeyMapping` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_regenerateSecretApplicationOauthKeysKeyMapping(string arg0, string arg1) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_regenerateSecretApplicationOauthKeysKeyMapping(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `renewConsumerSecret` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_renewConsumerSecret(string arg0, string arg1) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_renewConsumerSecret(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `revokeAPIKey` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtrestapidevportalv1dto:APIKeyRevokeRequestDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_revokeAPIKey(string arg0, orgwso2apkapimgtrestapidevportalv1dto:APIKeyRevokeRequestDTO arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_revokeAPIKey(java:fromString(arg0), arg1.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `updateApplication` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_updateApplication(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_updateApplication(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `updateApplicationKeysKeyType` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_updateApplicationKeysKeyType(string arg0, string arg1, orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO arg2) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_updateApplicationKeysKeyType(java:fromString(arg0), java:fromString(arg1), arg2.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `updateApplicationOauthKeysKeyMapping` method of `org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApplicationsCommonImpl_updateApplicationOauthKeysKeyMapping(string arg0, string arg1, orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO arg2) returns orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_updateApplicationOauthKeysKeyMapping(java:fromString(arg0), java:fromString(arg1), arg2.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtrestapidevportalv1dto:ApplicationKeyDTO newObj = new (externalObj);
        return newObj;
    }
}

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_addApplication(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "addApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_cleanUpApplicationRegistrationByApplicationIdAndKeyMappingId(handle arg0, handle arg1) returns error? = @java:Method {
    name: "cleanUpApplicationRegistrationByApplicationIdAndKeyMappingId",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_cleanupApplicationRegistration(handle arg0, handle arg1) returns error? = @java:Method {
    name: "cleanupApplicationRegistration",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_deleteApplication(handle arg0) returns int|error = @java:Method {
    name: "deleteApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_exportApplication(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "exportApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.Boolean", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateAPIKey(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "generateAPIKey",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.APIKeyGenerateRequestDTO"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateApplicationToken(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "generateApplicationToken",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationTokenGenerateRequestDTO"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateKeys(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "generateKeys",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationKeyGenerateRequestDTO", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_generateTokenByOauthKeysKeyMappingId(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "generateTokenByOauthKeysKeyMappingId",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationTokenGenerateRequestDTO"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationById(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getApplicationById",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationKeyByAppIDAndKeyMapping(handle arg0, handle arg1) returns handle = @java:Method {
    name: "getApplicationKeyByAppIDAndKeyMapping",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationKeyByAppIDAndKeyType(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getApplicationKeyByAppIDAndKeyType",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationKeysByApplicationId(handle arg0) returns handle|error = @java:Method {
    name: "getApplicationKeysByApplicationId",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationList(handle arg0, handle arg1, handle arg2, handle arg3, int arg4, int arg5, handle arg6) returns handle|error = @java:Method {
    name: "getApplicationList",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String", "int", "int", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getApplicationOauthKeys(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getApplicationOauthKeys",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getExportedApplication(handle arg0) returns handle|error = @java:Method {
    name: "getExportedApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.io.InputStream"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getOwnerId(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "getOwnerId",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.util.List", "java.lang.String", "java.lang.Boolean", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_getSkippedAPIs(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5) returns handle|error = @java:Method {
    name: "getSkippedAPIs",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.util.Set", "java.lang.String", "java.lang.Boolean", "java.lang.Boolean", "org.wso2.apk.apimgt.api.model.Application", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_importApplication(handle receiver, handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5, handle arg6) returns handle|error = @java:Method {
    name: "importApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.io.InputStream", "java.lang.Boolean", "java.lang.Boolean", "java.lang.String", "java.lang.Boolean", "java.lang.Boolean", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_mapApplicationKeys(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "mapApplicationKeys",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationKeyMappingRequestDTO", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_preProcessApplication(handle arg0, handle arg1, handle arg2, boolean arg3) returns handle|error = @java:Method {
    name: "preProcessApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationDTO", "java.lang.String", "boolean"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_regenerateSecretApplicationOauthKeysKeyMapping(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "regenerateSecretApplicationOauthKeysKeyMapping",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_renewConsumerSecret(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "renewConsumerSecret",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_revokeAPIKey(handle arg0, handle arg1) returns error? = @java:Method {
    name: "revokeAPIKey",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.APIKeyRevokeRequestDTO"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_updateApplication(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "updateApplication",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_updateApplicationKeysKeyType(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "updateApplicationKeysKeyType",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationKeyDTO"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_updateApplicationOauthKeysKeyMapping(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "updateApplicationOauthKeysKeyMapping",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "org.wso2.apk.apimgt.rest.api.devportal.v1.dto.ApplicationKeyDTO"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["long"]
} external;

function org_wso2_apk_apimgt_rest_api_devportal_v1_common_impl_ApplicationsCommonImpl_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.rest.api.devportal.v1.common.impl.ApplicationsCommonImpl",
    paramTypes: ["long", "int"]
} external;

