import ballerina/jballerina.java;
import admin_service.org.wso2.apk.apimgt.api as orgwso2apkapimgtapi;
import admin_service.java.lang as javalang;
import admin_service.java.util as javautil;
import admin_service.java.io as javaio;
import admin_service.org.wso2.apk.apimgt.admin.dto as orgwso2apkapimgtadmindto;

# Ballerina class mapping for the Java `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl` class.
@java:Binding {'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl"}
public distinct class ThrottlingCommonImpl {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    public function notify() {
        org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    public function notifyAll() {
        org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The function that maps to the `addAdvancedPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_addAdvancedPolicy(orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO arg0) returns orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addAdvancedPolicy(arg0.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `addApplicationThrottlePolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_addApplicationThrottlePolicy(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addApplicationThrottlePolicy(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `addDenyPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_addDenyPolicy(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addDenyPolicy(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `addSubscriptionThrottlePolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_addSubscriptionThrottlePolicy(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addSubscriptionThrottlePolicy(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `exportThrottlingPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtadmindto:ExportThrottlePolicyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_exportThrottlingPolicy(string arg0, string arg1) returns orgwso2apkapimgtadmindto:ExportThrottlePolicyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_exportThrottlingPolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtadmindto:ExportThrottlePolicyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAdvancedPolicyById` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getAdvancedPolicyById(string arg0) returns orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAdvancedPolicyById(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAllAdvancedPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + return - The `orgwso2apkapimgtadmindto:AdvancedThrottlePolicyListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getAllAdvancedPolicy() returns orgwso2apkapimgtadmindto:AdvancedThrottlePolicyListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAllAdvancedPolicy();
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtadmindto:AdvancedThrottlePolicyListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAllDenyPolicies` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getAllDenyPolicies() returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAllDenyPolicies();
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getAllSubscriptionThrottlePolicies` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getAllSubscriptionThrottlePolicies() returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAllSubscriptionThrottlePolicies();
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getApplicationThrottlePolicies` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getApplicationThrottlePolicies() returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getApplicationThrottlePolicies();
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getApplicationThrottlePolicyById` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getApplicationThrottlePolicyById(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getApplicationThrottlePolicyById(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getDenyPolicyById` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getDenyPolicyById(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getDenyPolicyById(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getSubscriptionThrottlePolicyById` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_getSubscriptionThrottlePolicyById(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getSubscriptionThrottlePolicyById(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `importThrottlingPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `javaio:InputStream` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `boolean` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `javautil:Map` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_importThrottlingPolicy(javaio:InputStream arg0, string arg1, boolean arg2, string arg3) returns javautil:Map|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_importThrottlingPolicy(arg0.jObj, java:fromString(arg1), arg2, java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        javautil:Map newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `removeAdvancedPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_removeAdvancedPolicy(string arg0, string arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeAdvancedPolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `removeApplicationThrottlePolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_removeApplicationThrottlePolicy(string arg0, string arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeApplicationThrottlePolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `removeDenyPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_removeDenyPolicy(string arg0) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeDenyPolicy(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `removeSubscriptionThrottlePolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_removeSubscriptionThrottlePolicy(string arg0, string arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeSubscriptionThrottlePolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `throttlingPolicySearch` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_throttlingPolicySearch(string arg0) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_throttlingPolicySearch(java:fromString(arg0));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `updateAdvancedPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_updateAdvancedPolicy(string arg0, orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO arg1) returns orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateAdvancedPolicy(java:fromString(arg0), arg1.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtadmindto:AdvancedThrottlePolicyDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `updateApplicationThrottlePolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_updateApplicationThrottlePolicy(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateApplicationThrottlePolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `updateDenyPolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_updateDenyPolicy(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateDenyPolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `updateSubscriptionThrottlePolicy` method of `org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ThrottlingCommonImpl_updateSubscriptionThrottlePolicy(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateSubscriptionThrottlePolicy(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addAdvancedPolicy(handle arg0) returns handle|error = @java:Method {
    name: "addAdvancedPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["org.wso2.apk.apimgt.admin.dto.AdvancedThrottlePolicyDTO"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addApplicationThrottlePolicy(handle arg0) returns handle|error = @java:Method {
    name: "addApplicationThrottlePolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addDenyPolicy(handle arg0) returns handle|error = @java:Method {
    name: "addDenyPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_addSubscriptionThrottlePolicy(handle arg0) returns handle|error = @java:Method {
    name: "addSubscriptionThrottlePolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_exportThrottlingPolicy(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "exportThrottlingPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAdvancedPolicyById(handle arg0) returns handle|error = @java:Method {
    name: "getAdvancedPolicyById",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAllAdvancedPolicy() returns handle|error = @java:Method {
    name: "getAllAdvancedPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAllDenyPolicies() returns handle|error = @java:Method {
    name: "getAllDenyPolicies",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getAllSubscriptionThrottlePolicies() returns handle|error = @java:Method {
    name: "getAllSubscriptionThrottlePolicies",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getApplicationThrottlePolicies() returns handle|error = @java:Method {
    name: "getApplicationThrottlePolicies",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getApplicationThrottlePolicyById(handle arg0) returns handle|error = @java:Method {
    name: "getApplicationThrottlePolicyById",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getDenyPolicyById(handle arg0) returns handle|error = @java:Method {
    name: "getDenyPolicyById",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_getSubscriptionThrottlePolicyById(handle arg0) returns handle|error = @java:Method {
    name: "getSubscriptionThrottlePolicyById",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_importThrottlingPolicy(handle arg0, handle arg1, boolean arg2, handle arg3) returns handle|error = @java:Method {
    name: "importThrottlingPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.io.InputStream", "java.lang.String", "boolean", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeAdvancedPolicy(handle arg0, handle arg1) returns error? = @java:Method {
    name: "removeAdvancedPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeApplicationThrottlePolicy(handle arg0, handle arg1) returns error? = @java:Method {
    name: "removeApplicationThrottlePolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeDenyPolicy(handle arg0) returns error? = @java:Method {
    name: "removeDenyPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_removeSubscriptionThrottlePolicy(handle arg0, handle arg1) returns error? = @java:Method {
    name: "removeSubscriptionThrottlePolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_throttlingPolicySearch(handle arg0) returns handle|error = @java:Method {
    name: "throttlingPolicySearch",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateAdvancedPolicy(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "updateAdvancedPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.admin.dto.AdvancedThrottlePolicyDTO"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateApplicationThrottlePolicy(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "updateApplicationThrottlePolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateDenyPolicy(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "updateDenyPolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_updateSubscriptionThrottlePolicy(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "updateSubscriptionThrottlePolicy",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["long"]
} external;

function org_wso2_apk_apimgt_admin_impl_ThrottlingCommonImpl_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.admin.impl.ThrottlingCommonImpl",
    paramTypes: ["long", "int"]
} external;

