import ballerina/jballerina.java;
import ballerina/jballerina.java.arrays as jarrays;
import devportal_service.org.'json.simple as orgjsonsimple;
import devportal_service.org.wso2.apk.apimgt.api as orgwso2apkapimgtapi;
import devportal_service.java.lang as javalang;
import devportal_service.java.util as javautil;
import devportal_service.java.io as javaio;
import devportal_service.org.wso2.apk.apimgt.api.model as orgwso2apkapimgtapimodel;
import devportal_service.org.wso2.apk.apimgt.devportal.dto as orgwso2apkapimgtdevportaldto;

# Ballerina class mapping for the Java `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl` class.
@java:Binding {'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl"}
public distinct class ApisCommonImpl {

    *java:JObject;
    *javalang:Object;

    # The `handle` field that stores the reference to the `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl` object.
    public handle jObj;

    # The init function of the Ballerina class mapping the `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl` Java class.
    #
    # + obj - The `handle` value containing the Java reference of the object.
    public function init(handle obj) {
        self.jObj = obj;
    }

    # The function to retrieve the string representation of the Ballerina class mapping the `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl` Java class.
    #
    # + return - The `string` form of the Java object instance.
    public function toString() returns string {
        return java:toString(self.jObj) ?: "null";
    }
    # The function that maps to the `equals` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    #
    # + arg0 - The `javalang:Object` value required to map with the Java method parameter.
    # + return - The `boolean` value returning from the Java mapping.
    public function 'equals(javalang:Object arg0) returns boolean {
        return org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_equals(self.jObj, arg0.jObj);
    }

    # The function that maps to the `getClass` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    #
    # + return - The `javalang:Class` value returning from the Java mapping.
    public function getClass() returns javalang:Class {
        handle externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getClass(self.jObj);
        javalang:Class newObj = new (externalObj);
        return newObj;
    }

    # The function that maps to the `hashCode` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    #
    # + return - The `int` value returning from the Java mapping.
    public function hashCode() returns int {
        return org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_hashCode(self.jObj);
    }

    # The function that maps to the `notify` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    public function notify() {
        org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_notify(self.jObj);
    }

    # The function that maps to the `notifyAll` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    public function notifyAll() {
        org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_notifyAll(self.jObj);
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    #
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function 'wait() returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_wait(self.jObj);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait2(int arg0) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_wait2(self.jObj, arg0);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

    # The function that maps to the `wait` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
    #
    # + arg0 - The `int` value required to map with the Java method parameter.
    # + arg1 - The `int` value required to map with the Java method parameter.
    # + return - The `javalang:InterruptedException` value returning from the Java mapping.
    public function wait3(int arg0, int arg1) returns javalang:InterruptedException? {
        error|() externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_wait3(self.jObj, arg0, arg1);
        if (externalObj is error) {
            javalang:InterruptedException e = error javalang:InterruptedException(javalang:INTERRUPTEDEXCEPTION, externalObj, message = externalObj.message());
            return e;
        }
    }

}

# The function that maps to the `addCommentToAPI` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtdevportaldto:PostRequestBodyDTO` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:CommentDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_addCommentToAPI(string arg0, orgwso2apkapimgtdevportaldto:PostRequestBodyDTO arg1, string arg2, string arg3) returns orgwso2apkapimgtdevportaldto:CommentDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_addCommentToAPI(java:fromString(arg0), arg1.jObj, java:fromString(arg2), java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:CommentDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `deleteAPIUserRating` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_deleteAPIUserRating(string arg0, string arg1) returns orgwso2apkapimgtapi:APIManagementException? {
    error|() externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_deleteAPIUserRating(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    }
}

# The function that maps to the `deleteComment` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string[]` value required to map with the Java method parameter.
# + return - The `orgjsonsimple:JSONObject` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_deleteComment(string arg0, string arg1, string arg2, string[] arg3) returns orgjsonsimple:JSONObject|orgwso2apkapimgtapi:APIManagementException|error {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_deleteComment(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2), check jarrays:toHandle(arg3, "java.lang.String"));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgjsonsimple:JSONObject newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `editCommentOfAPI` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + arg4 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:CommentDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_editCommentOfAPI(string arg0, string arg1, string arg2, string arg3, string arg4) returns orgwso2apkapimgtdevportaldto:CommentDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_editCommentOfAPI(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2), java:fromString(arg3), java:fromString(arg4));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:CommentDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAPI` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapimodel:API` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getAPI(string arg0, string arg1) returns orgwso2apkapimgtapimodel:API|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPI(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtapimodel:API newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAPIByAPIId` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getAPIByAPIId(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIByAPIId(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getAPIList` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `int` value required to map with the Java method parameter.
# + arg1 - The `int` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getAPIList(int arg0, int arg1, string arg2, string arg3) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIList(arg0, arg1, java:fromString(arg2), java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getAPIRating` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg2 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:RatingListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getAPIRating(string arg0, javalang:Integer arg1, javalang:Integer arg2, string arg3) returns orgwso2apkapimgtdevportaldto:RatingListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIRating(java:fromString(arg0), arg1.jObj, arg2.jObj, java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:RatingListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAPIThrottlePolicies` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `javautil:List` value required to map with the Java method parameter.
# + arg1 - The `javautil:List` value required to map with the Java method parameter.
# + return - The `javautil:List` value returning from the Java mapping.
public function ApisCommonImpl_getAPIThrottlePolicies(javautil:List arg0, javautil:List arg1) returns javautil:List {
    handle externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIThrottlePolicies(arg0.jObj, arg1.jObj);
    javautil:List newObj = new (externalObj);
    return newObj;
}

# The function that maps to the `getAPITopicList` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:TopicListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getAPITopicList(string arg0, string arg1) returns orgwso2apkapimgtdevportaldto:TopicListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPITopicList(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:TopicListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getAllCommentsOfAPI` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg3 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg4 - The `javalang:Boolean` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:CommentListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getAllCommentsOfAPI(string arg0, string arg1, javalang:Integer arg2, javalang:Integer arg3, javalang:Boolean arg4) returns orgwso2apkapimgtdevportaldto:CommentListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAllCommentsOfAPI(java:fromString(arg0), java:fromString(arg1), arg2.jObj, arg3.jObj, arg4.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:CommentListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getCommentOfAPI` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `javalang:Boolean` value required to map with the Java method parameter.
# + arg4 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg5 - The `javalang:Integer` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:CommentDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getCommentOfAPI(string arg0, string arg1, string arg2, javalang:Boolean arg3, javalang:Integer arg4, javalang:Integer arg5) returns orgwso2apkapimgtdevportaldto:CommentDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getCommentOfAPI(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2), arg3.jObj, arg4.jObj, arg5.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:CommentDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getDocumentContent` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapimodel:DocumentationContent` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getDocumentContent(string arg0, string arg1, string arg2) returns orgwso2apkapimgtapimodel:DocumentationContent|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getDocumentContent(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtapimodel:DocumentationContent newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getDocumentation` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:DocumentDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getDocumentation(string arg0, string arg1, string arg2) returns orgwso2apkapimgtdevportaldto:DocumentDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getDocumentation(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:DocumentDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getDocumentationList` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg2 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:DocumentListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getDocumentationList(string arg0, javalang:Integer arg1, javalang:Integer arg2, string arg3) returns orgwso2apkapimgtdevportaldto:DocumentListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getDocumentationList(java:fromString(arg0), arg1.jObj, arg2.jObj, java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:DocumentListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getGraphqlPoliciesComplexity` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:GraphQLQueryComplexityInfoDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getGraphqlPoliciesComplexity(string arg0, string arg1) returns orgwso2apkapimgtdevportaldto:GraphQLQueryComplexityInfoDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getGraphqlPoliciesComplexity(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:GraphQLQueryComplexityInfoDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getGraphqlPoliciesComplexityTypes` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:GraphQLSchemaTypeListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getGraphqlPoliciesComplexityTypes(string arg0, string arg1) returns orgwso2apkapimgtdevportaldto:GraphQLSchemaTypeListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getGraphqlPoliciesComplexityTypes(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:GraphQLSchemaTypeListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getGraphqlSchemaDefinition` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getGraphqlSchemaDefinition(string arg0, string arg1) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getGraphqlSchemaDefinition(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getOpenAPIDefinitionForEnvironment` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `string` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getOpenAPIDefinitionForEnvironment(string arg0, string arg1, string arg2) returns string?|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getOpenAPIDefinitionForEnvironment(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        return java:toString(externalObj);
    }
}

# The function that maps to the `getRating` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `orgwso2apkapimgtapi:APIConsumer` value required to map with the Java method parameter.
# + arg1 - The `int` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + return - The `orgjsonsimple:JSONObject` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getRating(orgwso2apkapimgtapi:APIConsumer arg0, int arg1, string arg2, string arg3) returns orgjsonsimple:JSONObject|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getRating(arg0.jObj, arg1, java:fromString(arg2), java:fromString(arg3));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgjsonsimple:JSONObject newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getRepliesOfComment` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg3 - The `string` value required to map with the Java method parameter.
# + arg4 - The `javalang:Integer` value required to map with the Java method parameter.
# + arg5 - The `javalang:Boolean` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:CommentListDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getRepliesOfComment(string arg0, string arg1, javalang:Integer arg2, string arg3, javalang:Integer arg4, javalang:Boolean arg5) returns orgwso2apkapimgtdevportaldto:CommentListDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getRepliesOfComment(java:fromString(arg0), java:fromString(arg1), arg2.jObj, java:fromString(arg3), arg4.jObj, arg5.jObj);
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:CommentListDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getSdkArtifacts` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `javaio:File` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getSdkArtifacts(string arg0, string arg1, string arg2) returns javaio:File|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getSdkArtifacts(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        javaio:File newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getSubscriptionPolicies` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `javautil:List` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getSubscriptionPolicies(string arg0, string arg1) returns javautil:List|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getSubscriptionPolicies(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        javautil:List newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getThumbnail` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapimodel:ResourceFile` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getThumbnail(string arg0, string arg1) returns orgwso2apkapimgtapimodel:ResourceFile|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getThumbnail(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtapimodel:ResourceFile newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getUserRating` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:RatingDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getUserRating(string arg0, string arg1) returns orgwso2apkapimgtdevportaldto:RatingDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getUserRating(java:fromString(arg0), java:fromString(arg1));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:RatingDTO newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `getWSDLOfAPI` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `string` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtapimodel:ResourceFile` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_getWSDLOfAPI(string arg0, string arg1, string arg2) returns orgwso2apkapimgtapimodel:ResourceFile|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getWSDLOfAPI(java:fromString(arg0), java:fromString(arg1), java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtapimodel:ResourceFile newObj = new (externalObj);
        return newObj;
    }
}

# The function that maps to the `updateUserRating` method of `org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl`.
#
# + arg0 - The `string` value required to map with the Java method parameter.
# + arg1 - The `orgwso2apkapimgtdevportaldto:RatingDTO` value required to map with the Java method parameter.
# + arg2 - The `string` value required to map with the Java method parameter.
# + return - The `orgwso2apkapimgtdevportaldto:RatingDTO` or the `orgwso2apkapimgtapi:APIManagementException` value returning from the Java mapping.
public function ApisCommonImpl_updateUserRating(string arg0, orgwso2apkapimgtdevportaldto:RatingDTO arg1, string arg2) returns orgwso2apkapimgtdevportaldto:RatingDTO|orgwso2apkapimgtapi:APIManagementException {
    handle|error externalObj = org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_updateUserRating(java:fromString(arg0), arg1.jObj, java:fromString(arg2));
    if (externalObj is error) {
        orgwso2apkapimgtapi:APIManagementException e = error orgwso2apkapimgtapi:APIManagementException(orgwso2apkapimgtapi:APIMANAGEMENTEXCEPTION, externalObj, message = externalObj.message());
        return e;
    } else {
        orgwso2apkapimgtdevportaldto:RatingDTO newObj = new (externalObj);
        return newObj;
    }
}

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_addCommentToAPI(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "addCommentToAPI",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.devportal.dto.PostRequestBodyDTO", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_deleteAPIUserRating(handle arg0, handle arg1) returns error? = @java:Method {
    name: "deleteAPIUserRating",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_deleteComment(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "deleteComment",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String", "[Ljava.lang.String;"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_editCommentOfAPI(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4) returns handle|error = @java:Method {
    name: "editCommentOfAPI",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_equals(handle receiver, handle arg0) returns boolean = @java:Method {
    name: "equals",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.Object"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPI(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getAPI",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIByAPIId(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getAPIByAPIId",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIList(int arg0, int arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "getAPIList",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["int", "int", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIRating(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "getAPIRating",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.Integer", "java.lang.Integer", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPIThrottlePolicies(handle arg0, handle arg1) returns handle = @java:Method {
    name: "getAPIThrottlePolicies",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.util.List", "java.util.List"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAPITopicList(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getAPITopicList",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getAllCommentsOfAPI(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4) returns handle|error = @java:Method {
    name: "getAllCommentsOfAPI",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.Integer", "java.lang.Integer", "java.lang.Boolean"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getClass(handle receiver) returns handle = @java:Method {
    name: "getClass",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getCommentOfAPI(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5) returns handle|error = @java:Method {
    name: "getCommentOfAPI",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String", "java.lang.Boolean", "java.lang.Integer", "java.lang.Integer"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getDocumentContent(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "getDocumentContent",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getDocumentation(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "getDocumentation",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getDocumentationList(handle arg0, handle arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "getDocumentationList",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.Integer", "java.lang.Integer", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getGraphqlPoliciesComplexity(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getGraphqlPoliciesComplexity",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getGraphqlPoliciesComplexityTypes(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getGraphqlPoliciesComplexityTypes",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getGraphqlSchemaDefinition(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getGraphqlSchemaDefinition",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getOpenAPIDefinitionForEnvironment(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "getOpenAPIDefinitionForEnvironment",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getRating(handle arg0, int arg1, handle arg2, handle arg3) returns handle|error = @java:Method {
    name: "getRating",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["org.wso2.apk.apimgt.api.APIConsumer", "int", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getRepliesOfComment(handle arg0, handle arg1, handle arg2, handle arg3, handle arg4, handle arg5) returns handle|error = @java:Method {
    name: "getRepliesOfComment",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.Integer", "java.lang.String", "java.lang.Integer", "java.lang.Boolean"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getSdkArtifacts(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "getSdkArtifacts",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getSubscriptionPolicies(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getSubscriptionPolicies",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getThumbnail(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getThumbnail",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getUserRating(handle arg0, handle arg1) returns handle|error = @java:Method {
    name: "getUserRating",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_getWSDLOfAPI(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "getWSDLOfAPI",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "java.lang.String", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_hashCode(handle receiver) returns int = @java:Method {
    name: "hashCode",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_notify(handle receiver) = @java:Method {
    name: "notify",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_notifyAll(handle receiver) = @java:Method {
    name: "notifyAll",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_updateUserRating(handle arg0, handle arg1, handle arg2) returns handle|error = @java:Method {
    name: "updateUserRating",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["java.lang.String", "org.wso2.apk.apimgt.devportal.dto.RatingDTO", "java.lang.String"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_wait(handle receiver) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: []
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_wait2(handle receiver, int arg0) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["long"]
} external;

function org_wso2_apk_apimgt_devportal_impl_ApisCommonImpl_wait3(handle receiver, int arg0, int arg1) returns error? = @java:Method {
    name: "wait",
    'class: "org.wso2.apk.apimgt.devportal.impl.ApisCommonImpl",
    paramTypes: ["long", "int"]
} external;

